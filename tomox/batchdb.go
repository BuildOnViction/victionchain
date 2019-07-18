package tomox

import (
	"bytes"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	lru "github.com/hashicorp/golang-lru"
	"sync"
)

const (
	defaultCacheLimit = 1024
	defaultMaxPending = 1024
)

type BatchItem struct {
	Value interface{}
}

type BatchDatabase struct {
	db             *ethdb.LDBDatabase
	lock           sync.RWMutex
	itemCacheLimit int
	itemMaxPending int
	emptyKey       []byte
	pendingItems   map[string]*BatchItem
	cacheItems     *lru.Cache // Cache for reading
	dryRunCache    *lru.Cache
	Debug          bool
}

// NewBatchDatabase use rlp as encoding
func NewBatchDatabase(datadir string, cacheLimit, maxPending int) *BatchDatabase {
	return NewBatchDatabaseWithEncode(datadir, cacheLimit, maxPending)
}

// batchdatabase is a fast cache db to retrieve in-mem object
func NewBatchDatabaseWithEncode(datadir string, cacheLimit, maxPending int) *BatchDatabase {
	db, err := ethdb.NewLDBDatabase(datadir, 128, 1024)
	if err != nil {
		log.Error("Can't create new DB", "error", err)
		return nil
	}
	itemCacheLimit := defaultCacheLimit
	if cacheLimit > 0 {
		itemCacheLimit = cacheLimit
	}
	itemMaxPending := defaultMaxPending
	if maxPending > 0 {
		itemMaxPending = maxPending
	}

	cacheItems, _ := lru.New(defaultCacheLimit)
	dryRunCache, _ := lru.New(defaultCacheLimit)

	batchDB := &BatchDatabase{
		db:             db,
		itemCacheLimit: itemCacheLimit,
		itemMaxPending: itemMaxPending,
		cacheItems:     cacheItems,
		emptyKey:       EmptyKey(), // pre alloc for comparison
		pendingItems:   make(map[string]*BatchItem),
		dryRunCache:    dryRunCache,
	}

	return batchDB

}

func (db *BatchDatabase) IsEmptyKey(key []byte) bool {
	return key == nil || len(key) == 0 || bytes.Equal(key, db.emptyKey)
}

func (db *BatchDatabase) getCacheKey(key []byte) string {
	return hex.EncodeToString(key)
}

func (db *BatchDatabase) Has(key []byte, dryrun bool) (bool, error) {
	if db.IsEmptyKey(key) {
		return false, nil
	}
	cacheKey := db.getCacheKey(key)

	if dryrun {
		if val, ok := db.dryRunCache.Get(cacheKey); ok {
			if val == nil {
				return false, nil
			}
			return true, nil
		}
	}

	// has in pending and is not deleted
	db.lock.Lock()
	if _, ok := db.pendingItems[cacheKey]; ok {
		db.lock.Unlock()
		return true, nil
	}
	db.lock.Unlock()

	if db.cacheItems.Contains(cacheKey) {
		return true, nil
	}

	return db.db.Has(key)
}

func (db *BatchDatabase) Get(key []byte, val interface{}, dryrun bool) (interface{}, error) {

	if db.IsEmptyKey(key) {
		return nil, nil
	}

	cacheKey := db.getCacheKey(key)

	if dryrun {
		if value, ok := db.dryRunCache.Get(cacheKey); ok {
			log.Debug("Debug DB get from dry-run cache", "cacheKey", cacheKey, "val", value)
			return value, nil
		}
	}

	db.lock.Lock()
	if pendingItem, ok := db.pendingItems[cacheKey]; ok {
		// we get value from the pending item
		db.lock.Unlock()
		return pendingItem.Value, nil
	}
	log.Debug("Debug DB get", "pending map", db.pendingItems, "cacheKey", cacheKey)
	db.lock.Unlock()

	if cached, ok := db.cacheItems.Get(cacheKey); ok {
		val = cached
	} else {

		// we can use lru for retrieving cache item, by default leveldb support get data from cache
		// but it is raw bytes
		b, err := db.db.Get(key)
		if err != nil {
			log.Debug("Key not found", "key", hex.EncodeToString(key), "err", err)
			return nil, err
		}

		err = DecodeBytesItem(b, val)

		// has problem here
		if err != nil {
			return nil, err
		}

		// update cache when reading
		db.cacheItems.Add(cacheKey, val)

	}

	return val, nil
}

func (db *BatchDatabase) Put(key []byte, val interface{}, dryrun bool) error {
	cacheKey := db.getCacheKey(key)
	if dryrun {
		log.Debug("Debug DB put to dry-run cache", "cacheKey", cacheKey, "val", val)
		db.dryRunCache.Add(cacheKey, val)
		return nil
	}

	db.lock.Lock()
	db.pendingItems[cacheKey] = &BatchItem{Value: val}
	log.Debug("Debug DB put", "pending map", db.pendingItems, "cacheKey", cacheKey)
	db.lock.Unlock()

	if len(db.pendingItems) >= db.itemMaxPending {
		return db.Commit()
	}

	return nil
}

func (db *BatchDatabase) Delete(key []byte, force bool, dryrun bool) error {
	// by default, we force delete both db and cache,
	// for better performance, we can mark a Deleted flag, to do batch delete
	cacheKey := db.getCacheKey(key)

	if dryrun {
		log.Debug("Debug DB delete from dry-run cache", "cacheKey", cacheKey)
		db.dryRunCache.Add(cacheKey, nil)
		return nil
	}
	log.Debug("Debug DB delete ", "cacheKey", cacheKey)
	db.lock.Lock()
	defer db.lock.Unlock()
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
	db.lock.Lock()
	for cacheKey, item := range db.pendingItems {
		key, _ := hex.DecodeString(cacheKey)
		value, err := EncodeBytesItem(item.Value)
		if err != nil {
			log.Error("Can't commit", "err", err)
			db.lock.Unlock()
			return err
		}

		batch.Put(key, value)
		log.Debug("Committed to DB", "key", cacheKey, "value", ToJSON(item.Value))
	}
	// commit pending items does not affect the cache
	db.pendingItems = make(map[string]*BatchItem)
	db.lock.Unlock()
	return batch.Write()
}

func (db *BatchDatabase) InitDryRunMode() {
	log.Debug("Start dry-run mode, clear old data")
	db.dryRunCache.Purge()
}

func (db *BatchDatabase) SaveDryRunResult() error {

	batch := db.db.NewBatch()
	db.lock.Lock()
	for _, cacheKey := range db.dryRunCache.Keys() {
		key, err := hex.DecodeString(cacheKey.(string))
		if err != nil {
			log.Error("Can't save dry-run result", "err", err)
			db.lock.Unlock()
			return err
		}
		val, ok := db.dryRunCache.Get(cacheKey)
		if !ok {
			continue
		}
		if val == nil {
			db.db.Delete(key)
			continue
		}

		value, err := EncodeBytesItem(val)
		if err != nil {
			log.Error("Can't save dry-run result", "err", err)
			db.lock.Unlock()
			return err
		}

		batch.Put(key, value)
		log.Debug("Saved dry-run result to DB", "cacheKey", hex.EncodeToString(key), "value", ToJSON(val))
	}
	// purge cache data
	db.dryRunCache.Purge()
	db.lock.Unlock()
	return batch.Write()
}