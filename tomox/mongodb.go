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

// Global instance of Database struct for singleton use
var db *MongoDatabase

// InitSession initializes a new session with mongodb
func NewMongoDatabase(session *mgo.Session, mongoURL string) (*MongoDatabase, error) {
	dbName := "tomodex"
	mongoURL = "mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs0"
	if db == nil {
		if session == nil {
			// Initialize new session
			ns, err := mgo.Dial(mongoURL)
			if err != nil {
				return nil, err
			}

			session = ns
		}

		db = &MongoDatabase{
			Session:        session,
			dbName:         dbName,
			itemMaxPending: defaultMaxPending,
			pendingItems:   make(map[string]*MongoItem),
		}
	}

	return db, nil
}

func (m *MongoDatabase) IsEmptyKey(key []byte) bool {
	return key == nil || len(key) == 0 || bytes.Equal(key, db.emptyKey)
}

func (m *MongoDatabase) getCacheKey(key []byte) string {
	return hex.EncodeToString(key)
}

func (m *MongoDatabase) Has(key []byte) (bool, error) {
	if m.IsEmptyKey(key) {
		return false, nil
	}

	cacheKey := m.getCacheKey(key)

	// has in pending and is not deleted
	if _, ok := db.pendingItems[cacheKey]; ok {
		return true, nil
	}

	sc := m.Session.Copy()
	defer sc.Close()

	query := bson.M{"key": cacheKey}

	// Find key in "items" collection
	numItems, err := sc.DB(m.dbName).C("items").Find(query).Limit(1).Count()

	if err != nil {
		return false, err
	}

	if numItems == 1 {
		return true, nil
	}

	// Find key in "orders" collection
	numOrders, err := sc.DB(m.dbName).C("orders").Find(query).Limit(1).Count()

	if err != nil {
		return false, err
	}

	if numOrders == 1 {
		return true, nil
	}

	return false, nil
}

func (m *MongoDatabase) Get(key []byte, val interface{}) (interface{}, error) {

	if m.IsEmptyKey(key) {
		return nil, nil
	}

	cacheKey := m.getCacheKey(key)

	if pendingItem, ok := m.pendingItems[cacheKey]; ok {
		// we get value from the pending item
		return pendingItem.Value, nil
	}

	sc := db.Session.Copy()
	defer sc.Close()

	var i *MongoItemRecord

	query := bson.M{"key": cacheKey}

	err := sc.DB(m.dbName).C("items").Find(query).One(&i)

	if err != nil {
		return nil, err
	}

	err = DecodeBytesItem(common.Hex2Bytes(i.Value), val)

	// has problem here
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (m *MongoDatabase) Put(key []byte, val interface{}) error {
	cacheKey := m.getCacheKey(key)

	switch val.(type) {
	case *types.Trade:
		err := db.CommitTrade(val.(*types.Trade)) // Put order into "orders" collection

		if err != nil {
			log.Error(err.Error())
			return err
		}

		break
	case *OrderItem:
		err := db.CommitOrder(cacheKey, val.(*OrderItem)) // Put order into "orders" collection

		if err != nil {
			log.Error(err.Error())
			return err
		}
		// There is no break here
	default:
		db.pendingItems[cacheKey] = &MongoItem{Value: val}
		break
	}

	// Put everything (includes order) into "items" collection
	if len(db.pendingItems) >= db.itemMaxPending {
		return db.Commit()
	}

	return nil
}

func (m *MongoDatabase) Delete(key []byte, force bool) error {
	cacheKey := m.getCacheKey(key)

	// Delete from object db.pendingItems
	delete(db.pendingItems, cacheKey)

	sc := m.Session.Copy()
	defer sc.Close()

	query := bson.M{"key": cacheKey}

	found, err := m.Has(key)

	if err != nil {
		return err
	}

	if found {
		err := sc.DB(m.dbName).C("items").Remove(query)
		if err != nil && err != mgo.ErrNotFound {
			log.Error("Error when deleting item", "error", err)
			return err
		}

		err = sc.DB(m.dbName).C("orders").Remove(query)
		if err != nil && err != mgo.ErrNotFound {
			log.Error("Error when deleting order", "error", err)
			return err
		}
	}

	return nil
}

func (m *MongoDatabase) Commit() error {

	sc := m.Session.Copy()
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

		_, err = sc.DB(m.dbName).C("items").Upsert(query, r)

		if err != nil {
			return err
		}

		log.Warn("Save", "cacheKey", cacheKey, "value", ToJSON(item.Value))
	}

	// Reset the object db.pendingItems
	db.pendingItems = make(map[string]*MongoItem)

	return nil
}

func (m *MongoDatabase) CommitOrder(cacheKey string, o *OrderItem) error {

	sc := m.Session.Copy()
	defer sc.Close()

	// Store the key
	o.Key = cacheKey

	query := bson.M{"key": cacheKey}

	_, err := sc.DB(m.dbName).C("orders").Upsert(query, o)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Debug("Save", "cacheKey", cacheKey, "value", ToJSON(o))

	return nil
}

func (m *MongoDatabase) CommitTrade(t *types.Trade) error {

	sc := m.Session.Copy()
	defer sc.Close()

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	query := bson.M{"hash": t.Hash}

	_, err := sc.DB(m.dbName).C("trades").Upsert(query, t)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Debug("Saved trade", "trade", ToJSON(t))

	return nil
}
