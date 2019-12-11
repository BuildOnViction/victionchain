package tomoxlending

import (
	"github.com/globalsign/mgo"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/ethdb"
	"github.com/tomochain/tomochain/tomoxlending/lendingstate"
)

type OrderDao interface {
	// for both leveldb and mongodb
	IsEmptyKey(key []byte) bool
	Close()

	// mongodb methods
	HasObject(hash common.Hash, val interface{}) (bool, error)
	GetObject(hash common.Hash, val interface{}) (interface{}, error)
	PutObject(hash common.Hash, val interface{}) error
	DeleteObject(hash common.Hash, val interface{}) error // won't return error if key not found
	GetOrderByTxHash(txhash common.Hash) []*lendingstate.LendingItem
	GetListOrderByHashes(hashes []string) []*lendingstate.LendingItem
	DeleteTradeByTxHash(txhash common.Hash)
	InitBulk() *mgo.Session
	CommitBulk() error

	// leveldb methods
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	NewBatch() ethdb.Batch
}
