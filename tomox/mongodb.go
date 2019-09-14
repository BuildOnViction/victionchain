package tomox

import (
	"bytes"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/hashicorp/golang-lru"
	"strings"
	"time"
)

type MongoItem struct {
	Value interface{}
}

type MongoItemRecord struct {
	Key   string
	Value string
}

type MongoDatabase struct {
	Session      *mgo.Session
	dbName       string
	emptyKey     []byte
	dryRunCaches map[common.Hash]*lru.Cache
	cacheItems   *lru.Cache // Cache for reading
}

// InitSession initializes a new session with mongodb
func NewMongoDatabase(session *mgo.Session, dbName string, mongoURL string, replicaSetName string, cacheLimit int) (*MongoDatabase, error) {
	if session == nil {
		// in case of multiple database instances
		hosts := strings.Split(mongoURL, ",")
		dbInfo := &mgo.DialInfo{
			Addrs:          hosts,
			Database:       dbName,
			ReplicaSetName: replicaSetName,
		}
		ns, err := mgo.DialWithInfo(dbInfo)
		if err != nil {
			return nil, err
		}
		session = ns
	}
	itemCacheLimit := defaultCacheLimit
	if cacheLimit > 0 {
		itemCacheLimit = cacheLimit
	}
	cacheItems, _ := lru.New(itemCacheLimit)

	db := &MongoDatabase{
		Session:      session,
		dbName:       dbName,
		cacheItems:   cacheItems,
		dryRunCaches: make(map[common.Hash]*lru.Cache),
	}

	return db, nil
}

func (db *MongoDatabase) IsEmptyKey(key []byte) bool {
	return key == nil || len(key) == 0 || bytes.Equal(key, db.emptyKey)
}

func (db *MongoDatabase) getCacheKey(key []byte) string {
	return hex.EncodeToString(key)
}

func (db *MongoDatabase) Has(key []byte, dryrun bool, blockHash common.Hash) (bool, error) {
	if db.IsEmptyKey(key) {
		return false, nil
	}
	cacheKey := db.getCacheKey(key)
	if db.cacheItems.Contains(cacheKey) {
		// for dry-run mode, do not read cacheItems
		return true, nil
	}

	sc := db.Session.Copy()
	defer sc.Close()

	query := bson.M{"key": cacheKey}

	// Find key in "items" collection
	numItems, err := sc.DB(db.dbName).C("items").Find(query).Limit(1).Count()

	if err != nil {
		return false, err
	}

	if numItems == 1 {
		return true, nil
	}

	// Find key in "orders" collection
	numOrders, err := sc.DB(db.dbName).C("orders").Find(query).Limit(1).Count()

	if err != nil {
		return false, err
	}

	if numOrders == 1 {
		return true, nil
	}

	return false, nil
}

func (db *MongoDatabase) Get(key []byte, val interface{}, dryrun bool, blockHash common.Hash) (interface{}, error) {

	if db.IsEmptyKey(key) {
		return nil, nil
	}

	cacheKey := db.getCacheKey(key)
	if cached, ok := db.cacheItems.Get(cacheKey); ok && !dryrun {
		return cached, nil
	} else {
		sc := db.Session.Copy()
		defer sc.Close()

		query := bson.M{"key": cacheKey}

		switch val.(type) {
		case *OrderItem:
			var oi *OrderItem
			err := sc.DB(db.dbName).C("orders").Find(query).One(&oi)
			if err != nil {
				return nil, err
			}
			db.cacheItems.Add(cacheKey, oi)
			return oi, nil
		default:
			var i *MongoItemRecord
			err := sc.DB(db.dbName).C("items").Find(query).One(&i)
			if err != nil {
				return nil, err
			}
			err = DecodeBytesItem(common.Hex2Bytes(i.Value), val)
			if err != nil {
				return nil, err
			}
			if !dryrun {
				db.cacheItems.Add(cacheKey, val)
			}
			return val, nil
		}
	}
}

func (db *MongoDatabase) Put(key []byte, val interface{}, dryrun bool, blockHash common.Hash) error {
	cacheKey := db.getCacheKey(key)
	db.cacheItems.Add(cacheKey, val)

	switch val.(type) {
	case *Trade:
		// Put trade into "trades" collection
		if err := db.CommitTrade(val.(*Trade)); err != nil {
			log.Error(err.Error())
			return err
		}
	case *OrderItem:
		// Put order into "orders" collection
		if err := db.CommitOrder(cacheKey, val.(*OrderItem)); err != nil {
			log.Error(err.Error())
			return err
		}
	default:
		// put general item into "items" collection
		if err := db.CommitItem(cacheKey, val); err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}

func (db *MongoDatabase) Delete(key []byte, dryrun bool, blockHash common.Hash) error {
	cacheKey := db.getCacheKey(key)
	db.cacheItems.Remove(cacheKey)

	sc := db.Session.Copy()
	defer sc.Close()

	query := bson.M{"key": cacheKey}

	found, err := db.Has(key, dryrun, blockHash)
	if err != nil {
		return err
	}

	if found {
		err := sc.DB(db.dbName).C("items").Remove(query)
		if err != nil && err != mgo.ErrNotFound {
			log.Error("Error when deleting item", "error", err)
			return err
		}

		err = sc.DB(db.dbName).C("orders").Remove(query)
		if err != nil && err != mgo.ErrNotFound {
			log.Error("Error when deleting order", "error", err)
			return err
		}
	}

	return nil
}

func (db *MongoDatabase) InitDryRunMode(hash common.Hash) {
	// SDK node (which running with mongodb) doesn't run Matching engine
	// dry-run cache is useless for sdk node
}

func (db *MongoDatabase) SaveDryRunResult(hash common.Hash) error {
	// SDK node (which running with mongodb) doesn't run Matching engine
	// dry-run cache is useless for sdk node
	return nil
}

func (db *MongoDatabase) CancelOrder(orderHash common.Hash) error {
	sc := db.Session.Copy()
	defer sc.Close()
	query := bson.M{"hash": orderHash.Hex()}
	var result *OrderItem
	if err := sc.DB(db.dbName).C("orders").Find(query).One(&result); err != nil {
		if err == mgo.ErrNotFound {
			//cancel done
			return nil
		}
		return err
	}
	result.Status = Cancel
	if _, err := sc.DB(db.dbName).C("orders").Upsert(query, result); err != nil {
		return err
	}
	return nil
}

func (db *MongoDatabase) CommitOrder(cacheKey string, o *OrderItem) error {

	sc := db.Session.Copy()
	defer sc.Close()

	// Store the key
	if len(o.Key) == 0 {
		o.Key = cacheKey
	}

	query := bson.M{"key": cacheKey}

	_, err := sc.DB(db.dbName).C("orders").Upsert(query, o)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Debug("Save orderItem", "cacheKey", cacheKey, "value", ToJSON(o))

	return nil
}

func (db *MongoDatabase) CommitTrade(t *Trade) error {

	sc := db.Session.Copy()
	defer sc.Close()

	t.ID = bson.NewObjectId()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	query := bson.M{"hash": t.Hash}

	_, err := sc.DB(db.dbName).C("trades").Upsert(query, t)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Debug("Saved trade", "trade", ToJSON(t))

	return nil
}

func (db *MongoDatabase) CommitItem(cacheKey string, val interface{}) error {
	sc := db.Session.Copy()
	defer sc.Close()

	data, err := EncodeBytesItem(val)
	if err != nil {
		return err
	}

	r := &MongoItemRecord{
		Key:   cacheKey,
		Value: common.Bytes2Hex(data),
	}

	query := bson.M{"key": cacheKey}
	if _, err := sc.DB(db.dbName).C("items").Upsert(query, r); err != nil {
		return err
	}
	log.Debug("Save", "cacheKey", cacheKey, "value", ToJSON(common.Bytes2Hex(data)))
	return nil
}
