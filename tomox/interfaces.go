package tomox

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/tomox/tomox_state"
	"github.com/globalsign/mgo"
)

type OrderDao interface {
	IsEmptyKey(key []byte) bool
	HasObject(key []byte) (bool, error)
	GetObject(key []byte, val interface{}) (interface{}, error)
	PutObject(key []byte, val interface{}) error
	DeleteObject(key []byte) error // won't return error if key not found
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	GetOrderByTxHash(txhash common.Hash) []*tomox_state.OrderItem
	GetListOrderByHashes(hashes []string) []*tomox_state.OrderItem
	DeleteTradeByTxHash(txhash common.Hash)
	InitBulk() *mgo.Session
	CommitBulk(sc *mgo.Session) error
	Close()
	NewBatch() ethdb.Batch
}
