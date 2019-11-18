package tomox

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	lru "github.com/hashicorp/golang-lru"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/ethdb"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/tomox/tomox_state"
	"strings"
)

const (
	ordersCollection = "orders"
	tradesCollection = "trades"
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
	if err := db.EnsureIndexes(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *MongoDatabase) IsEmptyKey(key []byte) bool {
	return key == nil || len(key) == 0 || bytes.Equal(key, db.emptyKey)
}

func (db *MongoDatabase) getCacheKey(key []byte) string {
	return hex.EncodeToString(key)
}

func (db *MongoDatabase) HasObject(hash common.Hash) (bool, error) {
	if db.IsEmptyKey(hash.Bytes()) {
		return false, nil
	}
	cacheKey := db.getCacheKey(hash.Bytes())
	if db.cacheItems.Contains(cacheKey) {
		return true, nil
	}

	sc := db.Session.Copy()
	defer sc.Close()
	var (
		count int
		err   error
	)
	query := bson.M{"hash": hash.Hex()}

	// Find key in "orders" collection
	count, err = sc.DB(db.dbName).C("orders").Find(query).Limit(1).Count()

	if err != nil {
		return false, err
	}

	if count == 1 {
		return true, nil
	}
	// Find key in "trades" collection
	count, err = sc.DB(db.dbName).C("trades").Find(query).Limit(1).Count()

	if err != nil {
		return false, err
	}

	if count == 1 {
		return true, nil
	}

	return false, nil
}

func (db *MongoDatabase) GetObject(hash common.Hash, val interface{}) (interface{}, error) {

	if db.IsEmptyKey(hash.Bytes()) {
		return nil, nil
	}

	cacheKey := db.getCacheKey(hash.Bytes())
	if cached, ok := db.cacheItems.Get(cacheKey); ok {
		return cached, nil
	} else {
		sc := db.Session.Copy()
		defer sc.Close()

		query := bson.M{"hash": hash.Hex()}

		switch val.(type) {
		case *tomox_state.OrderItem:
			var oi *tomox_state.OrderItem
			err := sc.DB(db.dbName).C("orders").Find(query).One(&oi)
			if err != nil {
				return nil, err
			}
			db.cacheItems.Add(cacheKey, oi)
			return oi, nil
		case *Trade:
			var t *Trade
			err := sc.DB(db.dbName).C("trades").Find(query).One(&t)
			if err != nil {
				return nil, err
			}
			db.cacheItems.Add(cacheKey, t)
			return t, nil
		default:
			return nil, nil
		}
	}
}

func (db *MongoDatabase) PutObject(hash common.Hash, val interface{}) error {
	cacheKey := db.getCacheKey(hash.Bytes())
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
		// Store the key
		o := val.(*tomox_state.OrderItem)
		if len(o.Key) == 0 {
			o.Key = cacheKey
		}
		if err := db.CommitOrder(o); err != nil {
			log.Error(err.Error())
			return err
		}
	default:
		log.Error("PutObject: object is neither order nor trade", "val", val)
	}

	return nil
}

func (db *MongoDatabase) DeleteObject(hash common.Hash) error {
	cacheKey := db.getCacheKey(hash.Bytes())
	db.cacheItems.Remove(cacheKey)

	sc := db.Session.Copy()
	defer sc.Close()

	query := bson.M{"hash": hash.Hex()}

	found, err := db.HasObject(hash)
	if err != nil {
		return err
	}

	if found {
		err := sc.DB(db.dbName).C("trades").Remove(query)
		if err != nil && err != mgo.ErrNotFound {
			log.Error("Error when deleting trades", "error", err)
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

func (db *MongoDatabase) CommitOrder(o *tomox_state.OrderItem) error {
	if o.Status == OrderStatusOpen {
		db.orderBulk.Insert(o)
	} else {
		query := bson.M{"hash": o.Hash.Hex()}
		db.orderBulk.Upsert(query, o)
	}
	return nil
}

func (db *MongoDatabase) CommitTrade(t *Trade) error {
	// for trades: insert only, no update
	// Hence, insert is better than upsert
	db.tradeBulk.Insert(t)
	return nil
}

func (db *MongoDatabase) InitBulk() *mgo.Session {
	sc := db.Session.Copy()
	db.orderBulk = sc.DB(db.dbName).C("orders").Bulk()
	db.tradeBulk = sc.DB(db.dbName).C("trades").Bulk()
	return sc
}

func (db *MongoDatabase) CommitBulk() error {
	if _, err := db.orderBulk.Run(); err != nil && !mgo.IsDup(err) {
		return err
	}
	if _, err := db.tradeBulk.Run(); err != nil && !mgo.IsDup(err) {
		return err
	}
	return nil
}

func (db *MongoDatabase) Put(key []byte, val []byte) error {
	// for levelDB only
	return nil
}

func (db *MongoDatabase) Delete(key []byte) error {
	// for levelDB only
	return nil
}

func (db *MongoDatabase) Has(key []byte) (bool, error) {
	// for levelDB only
	return false, nil
}

func (db *MongoDatabase) Get(key []byte) ([]byte, error) {
	// for levelDB only
	return nil, nil
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
		return []*tomox_state.OrderItem{}
	}
	return result
}

func (db *MongoDatabase) EnsureIndexes() error {
	orderHashIndex := mgo.Index{
		Key:        []string{"hash"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
		Name:       "index_order_hash",
	}
	tradeHashIndex := mgo.Index{
		Key:        []string{"hash"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
		Name:       "index_trade_hash",
	}
	sc := db.Session.Copy()
	defer sc.Close()
	if err := sc.DB(db.dbName).C(ordersCollection).EnsureIndex(orderHashIndex); err != nil {
		return fmt.Errorf("failed to index orders.hash . Err: %v", err)
	}
	if err := sc.DB(db.dbName).C(tradesCollection).EnsureIndex(tradeHashIndex); err != nil {
		return fmt.Errorf("failed to index trades.hash . Err: %v", err)
	}
	return nil
}

func (db *MongoDatabase) Close() {
	db.Close()
}

func (db *MongoDatabase) NewBatch() ethdb.Batch {
	// for levelDB only
	return nil
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
	// for levelDB only
}

func (b *Batch) Put(key, value []byte) error {
	// for levelDB only
	return nil
}

func (b *Batch) Write() error {
	// for levelDB only
	return nil
}

func (b *Batch) ValueSize() int {
	// for levelDB only
	return int(0)
}
func (b *Batch) Reset() {
	// for levelDB only
}
