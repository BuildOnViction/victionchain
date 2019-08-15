package tomox

import "github.com/ethereum/go-ethereum/common"

type OrderDao interface {
	IsEmptyKey(key []byte) bool
	Has(key []byte, dryrun bool) (bool, error)
	Get(key []byte, val interface{}, dryrun bool) (interface{}, error)
	Put(key []byte, val interface{}, dryrun bool) error
	Delete(key []byte, dryrun bool) error // won't return error if key not found
	InitDryRunMode()
	SaveDryRunResult() error
	CancelOrder(hash common.Hash) error
}
