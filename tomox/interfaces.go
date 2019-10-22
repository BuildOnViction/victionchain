package tomox

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

type OrderDao interface {
	IsEmptyKey(key []byte) bool
	HasObject(key []byte, dryrun bool, blockHash common.Hash) (bool, error)
	GetObject(key []byte, val interface{}, dryrun bool, blockHash common.Hash) (interface{}, error)
	PutObject(key []byte, val interface{}, dryrun bool, blockHash common.Hash) error
	DeleteObject(key []byte, dryrun bool, blockHash common.Hash) error // won't return error if key not found
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	Close()
	NewBatch() ethdb.Batch
}
