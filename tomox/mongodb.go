package tomox

import (
	"bytes"
	"encoding/hex"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/types"
)

type MongoItem struct {
	Value interface{}
}

type MongoItemRecord struct {
	Key   string
	Value string
}

type MongoDatabase struct {
	Session        *mgo.Session
	dbName         string
	emptyKey       []byte
	itemMaxPending int
	pendingItems   map[string]*MongoItem
}

// InitSession initializes a new session with mongodb
func NewMongoDatabase(session *mgo.Session, mongoURL string) (*MongoDatabase, error) {
	dbName := "tomodex"
	mongoURL = "mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs0"
	if session == nil {
		// Initialize new session
		ns, err := mgo.Dial(mongoURL)
		if err != nil {
			return nil, err
		}

		session = ns
	}

	db := &MongoDatabase{
		Session:        session,
		dbName:         dbName,
		itemMaxPending: defaultMaxPending,
		pendingItems:   make(map[string]*MongoItem),
	}

	return db, nil
}

func (db *MongoDatabase) IsEmptyKey(key []byte) bool {
	return key == nil || len(key) == 0 || bytes.Equal(key, db.emptyKey)
}

func (db *MongoDatabase) getCacheKey(key []byte) string {
	return hex.EncodeToString(key)
}

func (db *MongoDatabase) Has(key []byte) (bool, error) {
	if db.IsEmptyKey(key) {
		return false, nil
	}

	cacheKey := db.getCacheKey(key)

	// has in pending and is not deleted
	if _, ok := db.pendingItems[cacheKey]; ok {
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

func (db *MongoDatabase) Get(key []byte, val interface{}) (interface{}, error) {

	if db.IsEmptyKey(key) {
		return nil, nil
	}

	cacheKey := db.getCacheKey(key)

	if pendingItem, ok := db.pendingItems[cacheKey]; ok {
		// we get value from the pending item
		return pendingItem.Value, nil
	}
	log.Debug("Cache info (DB get)", "pending map", db.pendingItems, "cacheKey", cacheKey)

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
		return val, nil
	}
}

func (db *MongoDatabase) Put(key []byte, val interface{}) error {
	cacheKey := db.getCacheKey(key)
	switch val.(type) {
	case *types.Trade:
		err := db.CommitTrade(val.(*types.Trade)) // Put trade into "trades" collection

		if err != nil {
			log.Error(err.Error())
			return err
		}
	case *OrderItem:
		err := db.CommitOrder(cacheKey, val.(*OrderItem)) // Put order into "orders" collection

		if err != nil {
			log.Error(err.Error())
			return err
		}
		db.pendingItems[cacheKey] = &MongoItem{Value: val}
	default:
		db.pendingItems[cacheKey] = &MongoItem{Value: val}
		log.Debug("Cache info (DB put - trades & orders committed directly)", "pending map", db.pendingItems, "cacheKey", cacheKey)
	}

	// Put everything (includes order) into "items" collection
	if len(db.pendingItems) >= db.itemMaxPending {
		return db.Commit()
	}

	return nil
}

func (db *MongoDatabase) Delete(key []byte, force bool) error {
	cacheKey := db.getCacheKey(key)

	// Delete from object m.pendingItems
	delete(db.pendingItems, cacheKey)

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

func (db *MongoDatabase) Commit() error {

	sc := db.Session.Copy()
	defer sc.Close()

	for cacheKey, item := range db.pendingItems {
		valByte, err := EncodeBytesItem(item.Value)

		if err != nil {
			log.Error(err.Error())
			continue
		}

		r := &MongoItemRecord{
			Key:   cacheKey,
			Value: common.Bytes2Hex(valByte),
		}

		query := bson.M{"key": cacheKey}

		_, err = sc.DB(db.dbName).C("items").Upsert(query, r)

		if err != nil {
			return err
		}

		log.Warn("Save", "cacheKey", cacheKey, "value", ToJSON(item.Value))
	}

	// Reset the object db.pendingItems
	db.pendingItems = make(map[string]*MongoItem)

	return nil
}

func (db *MongoDatabase) CommitOrder(cacheKey string, o *OrderItem) error {

	sc := db.Session.Copy()
	defer sc.Close()

	// Store the key
	o.Key = cacheKey

	query := bson.M{"key": cacheKey}

	_, err := sc.DB(db.dbName).C("orders").Upsert(query, o)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Debug("Save", "cacheKey", cacheKey, "value", ToJSON(o))

	return nil
}

func (db *MongoDatabase) CommitTrade(t *types.Trade) error {

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
