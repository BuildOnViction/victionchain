package tomox

import (
	"errors"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type MongoDatabase struct {
	Session        *mgo.Session
	collectionName string
	dbName         string
}

// Global instance of Database struct for singleton use
var db *MongoDatabase

// InitSession initializes a new session with mongodb
func NewMongoDatabase(session *mgo.Session, mongoURL string) (*MongoDatabase, error) {
	dbName := "tomodex"
	collection := "orders"
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
			collectionName: collection,
			dbName:         dbName,
		}
	}

	return db, nil
}

func (m *MongoDatabase) IsEmptyKey(key []byte) bool {
	//TODO: put implementation here
	return false
}

func (m *MongoDatabase) Has(key []byte) (bool, error) {
	//TODO: put implementation here
	return false, nil
}

func (m *MongoDatabase) Get(key []byte, val interface{}) (interface{}, error) {
	//TODO: put implementation here
	return nil, nil
}

func (m *MongoDatabase) Put(key []byte, val interface{}) error {
	sc := db.Session.Copy()
	defer sc.Close()

	/* TODO: Refactor this later! ASSUME that val is OrderItem type, val might be TradeItem type in the future */
	o, ok := val.(*OrderItem)

	if ok == false {
		return errors.New("val is not OrderItem type")
	}

	query := bson.M{"hash": o.Hash.Hex()}

	_, err := sc.DB(m.dbName).C(m.collectionName).Upsert(query, o)

	if err != nil {
		return err
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
