package tomox

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/ethereum/go-ethereum/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type MongoItem struct {
	Value interface{}
}

type MongoDatabase struct {
	Session        *mgo.Session
	dbName         string
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
			pendingItems: make(map[string]*MongoItem),
		}
	}

	return db, nil
}

func (m *MongoDatabase) IsEmptyKey(key []byte) bool {
	return key == nil || len(key) == 0 || bytes.Equal(key, EmptyKey())
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
	cacheKey := m.getCacheKey(key)

	if pendingItem, ok := db.pendingItems[cacheKey]; ok {
		// we get value from the pending item
		return pendingItem.Value, nil
	}

	sc := db.Session.Copy()
	defer sc.Close()

	switch val.(type) {
	case *Item:
		var res *ItemRecord
		query := bson.M{"key": cacheKey}

		err := sc.DB(m.dbName).C("items").Find(query).One(&res)

		if err != nil {
			return nil, err
		}

		val = res.Value

		break
	case *OrderItem:
		oi, ok := val.(*OrderItem)

		if ok == false {
			log.Error("val is not OrderItem type")
			return nil, errors.New("val is not OrderItem type")
		}

		query := bson.M{"key": cacheKey}

		err := sc.DB(m.dbName).C("orders").Find(query).One(oi)

		if err != nil {
			return nil, err
		}

		break
	case *OrderListItem:
		var res *OrderListItemRecord
		query := bson.M{"key": cacheKey}

		err := sc.DB(m.dbName).C("items").Find(query).One(&res)

		if err != nil {
			return nil, err
		}

		val = res.Value

		break
	case *OrderTreeItem:
		var res *OrderTreeItemRecord
		query := bson.M{"key": cacheKey}

		err := sc.DB(m.dbName).C("items").Find(query).One(&res)

		if err != nil {
			return nil, err
		}

		val = res.Value

		break
	case *OrderBookItem:
		var res *OrderBookItemRecord
		query := bson.M{"key": cacheKey}

		err := sc.DB(m.dbName).C("items").Find(query).One(&res)

		if err != nil {
			return nil, err
		}

		val = res.Value

		break
	default:
		log.Error("Can't recognize value")
		break
	}

	return val, nil
}

func (m *MongoDatabase) Put(key []byte, val interface{}) error {
	cacheKey := m.getCacheKey(key)

	db.pendingItems[cacheKey] = &MongoItem{Value: val}

	if len(db.pendingItems) >= db.itemMaxPending {
		return db.Commit()
	} else {
		switch val.(type) {
		case *OrderItem:
			return db.Commit()
		default:
			return nil
		}
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
		switch item.Value.(type) {
		case *Item:
			i, ok := item.Value.(*Item)

			if ok == false {
				log.Error("val is not OrderListItem type")
				return errors.New("val is not OrderListItem type")
			}

			r := &ItemRecord{
				Key:   cacheKey,
				Value: i,
			}

			query := bson.M{"key": cacheKey}

			_, err := sc.DB(m.dbName).C("items").Upsert(query, r)

			if err != nil {
				return err
			}

			break
		case *OrderItem:
			oi, ok := item.Value.(*OrderItem)

			if ok == false {
				return errors.New("val is not OrderItem type")
			}

			// Store the key
			oi.Key = cacheKey

			query := bson.M{"key": cacheKey}

			_, err := sc.DB(m.dbName).C("orders").Upsert(query, oi)

			if err != nil {
				return err
			}

			break
		case *OrderListItem:
			oli, ok := item.Value.(*OrderListItem)

			if ok == false {
				log.Error("val is not OrderListItem type")
				return errors.New("val is not OrderListItem type")
			}

			r := &OrderListItemRecord{
				Key:   cacheKey,
				Value: oli,
			}

			query := bson.M{"key": cacheKey}

			_, err := sc.DB(m.dbName).C("items").Upsert(query, r)

			if err != nil {
				return err
			}

			break
		case *OrderTreeItem:
			oti, ok := item.Value.(*OrderTreeItem)

			if ok == false {
				log.Error("val is not OrderTreeItem type")
				return errors.New("val is not OrderTreeItem type")
			}

			r := &OrderTreeItemRecord{
				Key:   cacheKey,
				Value: oti,
			}

			query := bson.M{"key": cacheKey}

			_, err := sc.DB(m.dbName).C("items").Upsert(query, r)

			if err != nil {
				return err
			}

			break
		case *OrderBookItem:
			obi, ok := item.Value.(*OrderBookItem)

			if ok == false {
				log.Error("val is not OrderBookItem type")
				return errors.New("val is not OrderBookItem type")
			}

			r := &OrderBookItemRecord{
				Key:   cacheKey,
				Value: obi,
			}

			query := bson.M{"key": cacheKey}

			_, err := sc.DB(m.dbName).C("items").Upsert(query, r)

			if err != nil {
				return err
			}

			break
		default:
			log.Error("Can't recognize value")
			break
		}

		log.Debug("Save", "cacheKey", cacheKey, "value", ToJSON(item.Value))
	}

	// Reset the object db.pendingItems
	db.pendingItems = make(map[string]*MongoItem)

	return nil
}
