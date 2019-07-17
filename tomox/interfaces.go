package tomox

type OrderDao interface {
	IsEmptyKey(key []byte) bool
	Has(key []byte, dryrun bool) (bool, error)
	Get(key []byte, val interface{}, dryrun bool) (interface{}, error)
	Put(key []byte, val interface{}, dryrun bool) error
	Delete(key []byte, force bool, dryrun bool) error
	Commit() error
	InitDryRunMode()
	SaveDryRunResult() error
}
