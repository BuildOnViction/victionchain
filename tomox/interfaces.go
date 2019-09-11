package tomox

import "github.com/ethereum/go-ethereum/common"

type OrderDao interface {
	IsEmptyKey(key []byte) bool
	Has(key []byte, dryrun bool, blockHash common.Hash) (bool, error)
	Get(key []byte, val interface{}, dryrun bool, blockHash common.Hash) (interface{}, error)
	Put(key []byte, val interface{}, dryrun bool, blockHash common.Hash) error
	Delete(key []byte, dryrun bool, blockHash common.Hash) error // won't return error if key not found
	InitDryRunMode(blockHash common.Hash)
	SaveDryRunResult(blockHash common.Hash) error
	CancelOrder(hash common.Hash) error
}
