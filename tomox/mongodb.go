package tomox

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type MongoDatabase struct {
	Session *mgo.Session
	dbName  string
}

type MongoRecord struct {
	ID    bson.ObjectId
	Key   string
	Value interface{}
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
			Session: session,
			dbName:  dbName,
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
	//TODO: put implementation here
	return false, nil
}

func (m *MongoDatabase) Get(key []byte, val interface{}) (interface{}, error) {
	cacheKey := db.getCacheKey(key)

	sc := db.Session.Copy()
	defer sc.Close()

	var res MongoRecord

	switch val.(type) {
	case *Item:
		query := bson.M{"key": cacheKey}

		err := sc.DB(m.dbName).C("node_items").Find(query).One(res)

		if err != nil {
			return nil, err
		}

		val = res.Value

		break
	case *OrderItem:
		oi, ok := val.(*OrderItem)

		if ok == false {
			fmt.Println("val is not OrderItem type")
			return nil, errors.New("val is not OrderItem type")
		}
		query := bson.M{"hash": oi.Hash.Hex()}

		err := sc.DB(m.dbName).C("orders").Find(query).One(oi)

		if err != nil {
			return nil, err
		}

		break
	case *OrderListItem:
		query := bson.M{"key": cacheKey}

		err := sc.DB(m.dbName).C("order_list_items").Find(query).One(res)

		if err != nil {
			return nil, err
		}

		val = res.Value

		break
	case *OrderTreeItem:
		query := bson.M{"key": cacheKey}

		err := sc.DB(m.dbName).C("order_tree_items").Find(query).One(res)

		if err != nil {
			return nil, err
		}

		val = res.Value

		break
	case *OrderBookItem:
		query := bson.M{"key": cacheKey}

		err := sc.DB(m.dbName).C("order_book_items").Find(query).One(res)

		if err != nil {
			return nil, err
		}

		val = res.Value

		break
	default:
		fmt.Println("Can't recognize value")
		break
	}

	return val, nil
}

func (m *MongoDatabase) Put(key []byte, val interface{}) error {
	cacheKey := db.getCacheKey(key)

	fmt.Println("In Put function")
	sc := db.Session.Copy()
	defer sc.Close()

	switch val.(type) {
	case *Item:
		i, ok := val.(*Item)

		if ok == false {
			fmt.Println("val is not OrderListItem type")
			return errors.New("val is not OrderListItem type")
		}

		ib, err := EncodeItem(i)

		if err != nil {
			return err
		}

		query := bson.M{"key": cacheKey}

		_, err = sc.DB(m.dbName).C("node_items").Upsert(query, ib)

		if err != nil {
			return err
		}

		break
	case *OrderItem:
		oi, ok := val.(*OrderItem)

		if ok == false {
			fmt.Println("val is not OrderItem type")
			return errors.New("val is not OrderItem type")
		}

		oib, err := EncodeItem(oi)

		if err != nil {
			return err
		}

		query := bson.M{"hash": oi.Hash.Hex()}

		_, err = sc.DB(m.dbName).C("orders").Upsert(query, oib)

		if err != nil {
			return err
		}

		break
	case *OrderListItem:
		oli, ok := val.(*OrderListItem)

		if ok == false {
			fmt.Println("val is not OrderListItem type")
			return errors.New("val is not OrderListItem type")
		}

		olib, err := EncodeItem(oli)

		if err != nil {
			return err
		}

		query := bson.M{"key": cacheKey}

		_, err = sc.DB(m.dbName).C("order_list_items").Upsert(query, olib)

		if err != nil {
			return err
		}

		break
	case *OrderTreeItem:
		oti, ok := val.(*OrderTreeItem)

		if ok == false {
			fmt.Println("val is not OrderTreeItem type")
			return errors.New("val is not OrderTreeItem type")
		}

		otib, err := EncodeItem(oti)

		if err != nil {
			return err
		}

		query := bson.M{"key": cacheKey}

		_, err = sc.DB(m.dbName).C("order_tree_items").Upsert(query, otib)

		if err != nil {
			return err
		}

		break
	case *OrderBookItem:
		obi, ok := val.(*OrderBookItem)

		if ok == false {
			fmt.Println("val is not OrderBookItem type")
			return errors.New("val is not OrderBookItem type")
		}

		obib, err := EncodeItem(obi)

		if err != nil {
			return err
		}

		query := bson.M{"key": cacheKey}

		_, err = sc.DB(m.dbName).C("order_book_items").Upsert(query, obib)

		if err != nil {
			return err
		}

		break
	default:
		fmt.Println("Can't recognize value")
		break
	}

	return nil
}

func (m *MongoDatabase) Delete(key []byte, force bool) error {
	//TODO: put implementation here
	return nil
}

func (m *MongoDatabase) Commit() error {
	//TODO: put implementation here
	return nil
}
