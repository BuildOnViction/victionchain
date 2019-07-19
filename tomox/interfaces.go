package tomox

type OrderDao interface {
	IsEmptyKey(key []byte) bool
	Has(key []byte, dryrun bool) (bool, error)
	Get(key []byte, val interface{}, dryrun bool) (interface{}, error)
	Put(key []byte, val interface{}, dryrun bool) error
	Delete(key []byte, dryrun bool) error
	InitDryRunMode()
	SaveDryRunResult() error
}
