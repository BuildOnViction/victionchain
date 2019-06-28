package daos

import (
	"reflect"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/utils"
)

// Database struct contains the pointer to mgo.session
// It is a wrapper over mgo to help utilize mgo connection pool
type Database struct {
	Session *mgo.Session
}

// Global instance of Database struct for singleton use
var db *Database
var logger = utils.Logger

// InitSession initializes a new session with mongodb
func InitSession(session *mgo.Session) (*mgo.Session, error) {
	if db == nil {
		if session == nil {
			db1, err := mgo.Dial(app.Config.MongoURL)
			if err != nil {
				logger.Error(err)
				return nil, err
			}

			session = db1
		}

		db = &Database{session}
	}
	return db.Session, nil
}

func (d *Database) InitDatabase(session *mgo.Session) {
	d.Session = session
}

func (d *Database) GetCollection(dbName, collection string) *mgo.Collection {
	s := d.Session

	return s.DB(dbName).C(collection)
}

// Create is a wrapper for mgo.Insert function.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) Create(dbName, collection string, data ...interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Insert(data...)
	return
}

func (d *Database) Count(dbName, collection string, query interface{}) (int, error) {
	sc := d.Session.Copy()
	defer sc.Close()

	return sc.DB(dbName).C(collection).Find(query).Count()
}

// GetByID is a wrapper for mgo.FindId function.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) GetByID(dbName, collection string, id bson.ObjectId, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).FindId(id).One(response)
	return
}

// Get is a wrapper for mgo.Find function.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) Get(dbName, collection string, query interface{}, offset, limit int, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).Skip(offset).Limit(limit).All(response)
	return
}

// GetOne return one document
func (d *Database) GetOne(dbName, collection string, query interface{}, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).One(response)
	// logger.Debug(err, query)
	return
}

func (d *Database) Query(dbName, collection string, query interface{}, selector interface{}, offset, limit int, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).Skip(offset).Limit(limit).Select(selector).All(response)
	return
}

// GetAndSort is a wrapper for mgo.Find function with SORT function in pipeline.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) GetAndSort(dbName, collection string, query interface{}, sort []string, offset, limit int, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).Sort(sort...).Skip(offset).Limit(limit).All(response)
	return
}

// Update is a wrapper for mgo.Update function.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) Update(dbName, collection string, query interface{}, update interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	err := sc.DB(dbName).C(collection).Update(query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (d *Database) Upsert(dbName, collection string, query interface{}, update interface{}) (interface{}, error) {
	sc := d.Session.Copy()
	defer sc.Close()

	changed, err := sc.DB(dbName).C(collection).Upsert(query, update)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return changed.UpsertedId, nil
}

func (d *Database) UpdateAll(dbName, collection string, query interface{}, update interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	_, err := sc.DB(dbName).C(collection).UpdateAll(query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (d *Database) FindAndModify(dbName, collection string, query interface{}, change mgo.Change, response interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	_, err := sc.DB(dbName).C(collection).Find(query).Apply(change, response)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Aggregate is a wrapper for mgo.Pipe function.
// It is used to make mongo aggregate pipeline queries
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) Aggregate(dbName, collection string, query []bson.M, response interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	result := reflect.ValueOf(response).Interface()
	err := sc.DB(dbName).C(collection).Pipe(query).All(result)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Remove removes one document matching a certain query
func (d *Database) Remove(dbName, collection string, query []bson.M) error {
	sc := d.Session.Copy()
	defer sc.Close()

	err := sc.DB(dbName).C(collection).Remove(query)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// RemoveAll removes all the documents from a collection matching a certain query
func (d *Database) RemoveAll(dbName, collection string, query interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	_, err := sc.DB(dbName).C(collection).RemoveAll(query)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// DropCollection drops all the documents in a collection
func (d *Database) DropCollection(dbName, collection string) error {
	sc := d.Session.Copy()
	defer sc.Close()

	err := sc.DB(dbName).C(collection).DropCollection()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
