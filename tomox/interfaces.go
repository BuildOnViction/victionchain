package tomox

type OrderDao interface {
	IsEmptyKey(key []byte) bool
	Has(key []byte) (bool, error)
	Get(key []byte, val interface{}) (interface{}, error)
	Put(key []byte, val interface{}) error
	Delete(key []byte, force bool) error
	Commit() error
}
