package orderbook

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
)

const (
	defaultCacheLimit = 1024
	defaultMaxPending = 1024
)

// BatchItem : currently we do not support deletion batch, so ignore it
type BatchItem struct {
	Value interface{}
	// Deleted bool
}

type BatchDatabase struct {
	db             *ethdb.LDBDatabase
	itemCacheLimit int
	itemMaxPending int
	emptyKey       []byte
	pendingItems   map[string]*BatchItem
	cacheItems     *lru.Cache // Cache for reading
	Debug          bool

	EncodeToBytes EncodeToBytes
	DecodeBytes   DecodeBytes
}

// NewBatchDatabase use rlp as encoding
func NewBatchDatabase(datadir string, cacheLimit, maxPending int) *BatchDatabase {
	return NewBatchDatabaseWithEncode(datadir, cacheLimit, maxPending, rlp.EncodeToBytes, rlp.DecodeBytes)
}

// batchdatabase is a fast cache db to retrieve in-mem object
func NewBatchDatabaseWithEncode(datadir string, cacheLimit, maxPending int, encode EncodeToBytes, decode DecodeBytes) *BatchDatabase {
	db, _ := ethdb.NewLDBDatabase(datadir, 128, 1024)
	itemCacheLimit := defaultCacheLimit
	if cacheLimit > 0 {
		itemCacheLimit = cacheLimit
	}
	itemMaxPending := defaultMaxPending
	if maxPending > 0 {
		itemMaxPending = maxPending
	}

	cacheItems, _ := lru.New(defaultCacheLimit)

	batchDB := &BatchDatabase{
		db:             db,
		EncodeToBytes:  encode,
		DecodeBytes:    decode,
		itemCacheLimit: itemCacheLimit,
		itemMaxPending: itemMaxPending,
		cacheItems:     cacheItems,
		emptyKey:       EmptyKey(), // pre alloc for comparison
		pendingItems:   make(map[string]*BatchItem),
	}

	return batchDB

}

func (db *BatchDatabase) IsEmptyKey(key []byte) bool {
	return key == nil || len(key) == 0 || bytes.Equal(key, db.emptyKey)
}

func (db *BatchDatabase) getCacheKey(key []byte) string {
	return hex.EncodeToString(key)
}

func (db *BatchDatabase) Has(key []byte) (bool, error) {
	if db.IsEmptyKey(key) {
		return false, nil
	}
	cacheKey := db.getCacheKey(key)

	// has in pending and is not deleted
	if _, ok := db.pendingItems[cacheKey]; ok {
		return true, nil
	}

	if db.cacheItems.Contains(cacheKey) {
		return true, nil
	}

	return db.db.Has(key)
}

func (db *BatchDatabase) Get(key []byte, val interface{}) (interface{}, error) {

	if db.IsEmptyKey(key) {
		return nil, nil
	}

	cacheKey := db.getCacheKey(key)

	if pendingItem, ok := db.pendingItems[cacheKey]; ok {
		// we get value from the pending item
		return pendingItem.Value, nil
	}

	if cached, ok := db.cacheItems.Get(cacheKey); ok {
		val = cached
		if db.Debug {
			fmt.Println("Cache hit :", cacheKey)
		}
	} else {

		// we can use lru for retrieving cache item, by default leveldb support get data from cache
		// but it is raw bytes
		bytes, err := db.db.Get(key)
		if err != nil {
			if db.Debug {
				fmt.Printf("Key not found :%x\n", key)
			}
			return nil, err
		}

		err = db.DecodeBytes(bytes, val)

		// has problem here
		if err != nil {
			return nil, err
		}

		// update cache when reading
		db.cacheItems.Add(cacheKey, val)

	}

	return val, nil
}

func (db *BatchDatabase) Put(key []byte, val interface{}) error {

	cacheKey := db.getCacheKey(key)

	db.pendingItems[cacheKey] = &BatchItem{Value: val}

	if len(db.pendingItems) >= db.itemMaxPending {
		return db.Commit()
	}

	return nil
}

func (db *BatchDatabase) Delete(key []byte, force bool) error {

	// by default, we force delete both db and cache,
	// for better performance, we can mark a Deleted flag, to do batch delete
	cacheKey := db.getCacheKey(key)

	// force delete everything
	if force {
		delete(db.pendingItems, cacheKey)
		db.cacheItems.Remove(cacheKey)
	} else {
		if _, ok := db.pendingItems[cacheKey]; ok {
			// item.Deleted = true
			db.db.Delete(key)

			// delete from pending Items
			delete(db.pendingItems, cacheKey)
			// remove cache key as well
			db.cacheItems.Remove(cacheKey)
			return nil
		}
	}

	// cache not found, or force delete, must delete from database
	return db.db.Delete(key)
}

func (db *BatchDatabase) Commit() error {

	batch := db.db.NewBatch()
	for cacheKey, item := range db.pendingItems {
		key, _ := hex.DecodeString(cacheKey)
		value, err := db.EncodeToBytes(item.Value)
		if err != nil {
			fmt.Println(err)
			return err
		}

		batch.Put(key, value)

		if db.Debug {
			fmt.Printf("Save %x, value :%s\n", key, ToJSON(item.Value))
		}
	}
	// commit pending items does not affect the cache
	db.pendingItems = make(map[string]*BatchItem)
	return batch.Write()
}
