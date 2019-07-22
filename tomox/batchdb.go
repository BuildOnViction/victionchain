package tomox

import (
	"bytes"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
)

const (
	defaultCacheLimit = 1024
)

type BatchItem struct {
	Value interface{}
}

type BatchDatabase struct {
	db             *ethdb.LDBDatabase
	emptyKey       []byte
	cacheItems     *lru.Cache // Cache for reading
	dryRunCache    *lru.Cache
	Debug          bool
}

// NewBatchDatabase use rlp as encoding
func NewBatchDatabase(datadir string, cacheLimit int) *BatchDatabase {
	return NewBatchDatabaseWithEncode(datadir, cacheLimit)
}

// batchdatabase is a fast cache db to retrieve in-mem object
func NewBatchDatabaseWithEncode(datadir string, cacheLimit int) *BatchDatabase {
	db, err := ethdb.NewLDBDatabase(datadir, 128, 1024)
	if err != nil {
		log.Error("Can't create new DB", "error", err)
		return nil
	}
	itemCacheLimit := defaultCacheLimit
	if cacheLimit > 0 {
		itemCacheLimit = cacheLimit
	}

	cacheItems, _ := lru.New(itemCacheLimit)
	dryRunCache, _ := lru.New(itemCacheLimit)

	batchDB := &BatchDatabase{
		db:             db,
		cacheItems:     cacheItems,
		emptyKey:       EmptyKey(), // pre alloc for comparison
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
			log.Debug("Debug get from dry-run cache", "cacheKey", cacheKey, "val", value)
			return value, nil
		}
	}


	if cached, ok := db.cacheItems.Get(cacheKey); ok && !dryrun {
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
		log.Debug("Debug put to dry-run cache", "cacheKey", cacheKey, "val", val)
		db.dryRunCache.Add(cacheKey, val)
		return nil
	}

	log.Debug("Debug DB put to cacheItems", "cacheKey", cacheKey,  "val", val)
	db.cacheItems.Add(cacheKey, val)
	value, err := EncodeBytesItem(val)
	if err != nil {
		return err
	}
	return db.db.Put(key, value)
}

func (db *BatchDatabase) Delete(key []byte, dryrun bool) error {
	// by default, we force delete both db and cache,
	// for better performance, we can mark a Deleted flag, to do batch delete
	cacheKey := db.getCacheKey(key)

	//mark it to nil in dryrun cache
	if dryrun {
		log.Debug("Debug DB delete from dry-run cache", "cacheKey", cacheKey)
		db.dryRunCache.Add(cacheKey, nil)
		return nil
	}

	log.Debug("Debug DB delete ", "cacheKey", cacheKey)
	db.cacheItems.Remove(cacheKey)
	return db.db.Delete(key)
}

func (db *BatchDatabase) InitDryRunMode() {
	log.Debug("Start dry-run mode, clear old data")
	db.dryRunCache.Purge()
}

func (db *BatchDatabase) SaveDryRunResult() error {

	batch := db.db.NewBatch()
	for _, cacheKey := range db.dryRunCache.Keys() {
		key, err := hex.DecodeString(cacheKey.(string))
		if err != nil {
			log.Error("Can't save dry-run result (hex.DecodeString)", "err", err)
			return err
		}
		val, ok := db.dryRunCache.Get(cacheKey)
		if !ok {
			err := errors.New("can't get item from dryrun cache")
			log.Error("Can't save dry-run result (db.dryRunCache.Get)", "err", err)
			return err
		}
		if val == nil {
			if err := db.db.Delete(key); err != nil {
				log.Error("Can't save dry-run result (db.db.Delete)", "err", err)
				return err
			}
			continue
		}

		value, err := EncodeBytesItem(val)
		if err != nil {
			log.Error("Can't save dry-run result (EncodeBytesItem)", "err", err)
			return err
		}

		if err := batch.Put(key, value); err != nil {
			log.Error("Can't save dry-run result (batch.Put)", "err", err)
			return err
		}
		log.Debug("Saved dry-run result to DB", "cacheKey", hex.EncodeToString(key), "value", ToJSON(val))
	}
	// purge cache data
	db.dryRunCache.Purge()
	return batch.Write()
}
