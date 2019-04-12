package tomox

import (
	"github.com/globalsign/mgo"
)

type MongoDatabase struct {
	Session *mgo.Session
}

// Global instance of Database struct for singleton use
var db *MongoDatabase

// InitSession initializes a new session with mongodb
func NewMongoDatabase(session *mgo.Session, mongoURL string) (*MongoDatabase, error) {
	if db == nil {
		if session == nil {
			// Initialize new session
			ns, err := mgo.Dial(mongoURL)
			if err != nil {
				return nil, err
			}

			session = ns
		}

		db = &MongoDatabase{session}
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

	err := sc.DB(dbName).C(collection).Insert(val)

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
