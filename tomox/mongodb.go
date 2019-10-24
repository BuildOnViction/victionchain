package tomox

import (
	"bytes"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/tomox/tomox_state"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	lru "github.com/hashicorp/golang-lru"
	"strings"
)

type MongoItem struct {
	Value interface{}
}

type MongoItemRecord struct {
	Key   string
	Value string
}

type MongoDatabase struct {
	Session    *mgo.Session
	dbName     string
	emptyKey   []byte
	cacheItems *lru.Cache // Cache for reading
	orderBulk  *mgo.Bulk
	tradeBulk  *mgo.Bulk
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
		Session:    session,
		dbName:     dbName,
		cacheItems: cacheItems,
	}

	return db, nil
}

func (db *MongoDatabase) IsEmptyKey(key []byte) bool {
	return key == nil || len(key) == 0 || bytes.Equal(key, db.emptyKey)
}

func (db *MongoDatabase) getCacheKey(key []byte) string {
	return hex.EncodeToString(key)
}

func (db *MongoDatabase) HasObject(key []byte) (bool, error) {
	if db.IsEmptyKey(key) {
		return false, nil
	}
	cacheKey := db.getCacheKey(key)
	if db.cacheItems.Contains(cacheKey) {
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

func (db *MongoDatabase) GetObject(key []byte, val interface{}) (interface{}, error) {

	if db.IsEmptyKey(key) {
		return nil, nil
	}

	cacheKey := db.getCacheKey(key)
	if cached, ok := db.cacheItems.Get(cacheKey); ok {
		return cached, nil
	} else {
		sc := db.Session.Copy()
		defer sc.Close()

		query := bson.M{"key": cacheKey}

		switch val.(type) {
		case *tomox_state.OrderItem:
			var oi *tomox_state.OrderItem
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
			db.cacheItems.Add(cacheKey, val)

			return val, nil
		}
	}
}

func (db *MongoDatabase) PutObject(key []byte, val interface{}) error {
	cacheKey := db.getCacheKey(key)
	db.cacheItems.Add(cacheKey, val)

	switch val.(type) {
	case *Trade:
		// PutObject trade into "trades" collection
		if err := db.CommitTrade(val.(*Trade)); err != nil {
			log.Error(err.Error())
			return err
		}
	case *tomox_state.OrderItem:
		// PutObject order into "orders" collection
		if err := db.CommitOrder(cacheKey, val.(*tomox_state.OrderItem)); err != nil {
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

func (db *MongoDatabase) DeleteObject(key []byte) error {
	cacheKey := db.getCacheKey(key)
	db.cacheItems.Remove(cacheKey)

	sc := db.Session.Copy()
	defer sc.Close()

	query := bson.M{"key": cacheKey}

	found, err := db.HasObject(key)
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

func (db *MongoDatabase) CommitOrder(cacheKey string, o *tomox_state.OrderItem) error {
	// Store the key
	if len(o.Key) == 0 {
		o.Key = cacheKey
	}
	query := bson.M{"hash": o.Hash.Hex()}
	db.orderBulk.Upsert(query, o)
	return nil
}

func (db *MongoDatabase) CommitTrade(t *Trade) error {
	query := bson.M{"hash": t.Hash.Hex()}
	db.tradeBulk.Upsert(query, t)
	return nil
}

func (db *MongoDatabase) InitBulk() *mgo.Session {
	sc := db.Session.Copy()
	db.orderBulk = sc.DB(db.dbName).C("orders").Bulk()
	db.tradeBulk = sc.DB(db.dbName).C("trades").Bulk()
	return sc
}

func (db *MongoDatabase) CommitBulk(sc *mgo.Session) error {
	defer func() {
		sc.Close()
	}()
	if _, err := db.orderBulk.Run(); err != nil {
		return err
	}
	if _, err := db.tradeBulk.Run(); err != nil {
		return err
	}
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

func (db *MongoDatabase) Put(key []byte, val []byte) error {
	cacheKey := db.getCacheKey(key)
	db.cacheItems.Add(cacheKey, val)
	sc := db.Session.Copy()
	defer sc.Close()
	r := &MongoItemRecord{
		Key:   cacheKey,
		Value: common.Bytes2Hex(val),
	}
	query := bson.M{"key": cacheKey}
	if _, err := sc.DB(db.dbName).C("items").Upsert(query, r); err != nil {
		return err
	}
	return nil
}

func (db *MongoDatabase) Delete(key []byte) error {
	cacheKey := db.getCacheKey(key)
	db.cacheItems.Remove(cacheKey)

	sc := db.Session.Copy()
	defer sc.Close()

	query := bson.M{"key": cacheKey}

	found, err := db.Has(key)
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

func (db *MongoDatabase) Has(key []byte) (bool, error) {
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

func (db *MongoDatabase) Get(key []byte) ([]byte, error) {
	if db.IsEmptyKey(key) {
		return nil, nil
	}
	cacheKey := db.getCacheKey(key)
	if cached, ok := db.cacheItems.Get(cacheKey); ok {
		return cached.([]byte), nil
	} else {
		sc := db.Session.Copy()
		defer sc.Close()
		query := bson.M{"key": cacheKey}
		var oi []byte
		err := sc.DB(db.dbName).C("orders").Find(query).One(&oi)
		if err != nil {
			return nil, err
		}
		db.cacheItems.Add(cacheKey, oi)
		return oi, nil
	}
}

func (db *MongoDatabase) DeleteTradeByTxHash(txhash common.Hash) {
	sc := db.Session.Copy()
	defer sc.Close()

	query := bson.M{"txHash": txhash.Hex()}

	err := sc.DB(db.dbName).C("trades").Remove(query)
	if err != nil && err != mgo.ErrNotFound {
		log.Error("Error when deleting order", "error", err)
	}
}

func (db *MongoDatabase) GetOrderByTxHash(txhash common.Hash) []*tomox_state.OrderItem {
	var result []*tomox_state.OrderItem
	sc := db.Session.Copy()
	defer sc.Close()

	query := bson.M{"txHash": txhash.Hex()}

	if err := sc.DB(db.dbName).C("orders").Find(query).All(&result); err != nil && err != mgo.ErrNotFound {
		log.Error("failed to GetOrderByTxHash", "err", err, "Txhash", txhash)
	}
	return result
}

func (db *MongoDatabase) GetListOrderByHashes(hashes []string) []*tomox_state.OrderItem {
	var result []*tomox_state.OrderItem
	sc := db.Session.Copy()
	defer sc.Close()

	query := bson.M{"hash": bson.M{"$in": hashes}}

	if err := sc.DB(db.dbName).C("orders").Find(query).All(&result); err != nil && err != mgo.ErrNotFound {
		log.Error("failed to GetListOrderByHashes", "err", err, "hashes", hashes)
	}
	return result
}

func (db *MongoDatabase) Close() {
	db.Close()
}

func (db *MongoDatabase) NewBatch() ethdb.Batch {
	return &Batch{db: db, b: []keyvalue{}, size: 0, collection: ""}
}

type keyvalue struct {
	key   []byte
	value []byte
}
type Batch struct {
	db         *MongoDatabase
	collection string
	b          []keyvalue
	size       int
}

func (b *Batch) SetCollection(collection string) {
	b.collection = collection
}

func (b *Batch) Put(key, value []byte) error {
	b.b = append(b.b, keyvalue{key: key, value: value})
	b.size += len(value)
	return nil
}

func (b *Batch) Write() error {
	sc := b.db.Session.Copy()
	defer sc.Close()
	for _, keyvalue := range b.b {
		cacheKey := b.db.getCacheKey(keyvalue.key)
		b.db.cacheItems.Add(cacheKey, keyvalue.value)
		r := &MongoItemRecord{
			Key:   cacheKey,
			Value: common.Bytes2Hex(keyvalue.value),
		}
		query := bson.M{"key": cacheKey}
		if _, err := sc.DB(b.db.dbName).C(b.collection).Upsert(query, r); err != nil {
			return err
		}
	}
	return nil
}

func (b *Batch) ValueSize() int {
	return b.size
}
func (b *Batch) Reset() {
	b.size = 0
	b.b = b.b[:0]
}
