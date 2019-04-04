package tomox

type MongoDatabase struct {
	//TODO: put implementation here
}

func (m *MongoDatabase) IsEmptyKey(key []byte) bool {
	//TODO: put implementation here
	return false
}

func (m *MongoDatabase) Has(key []byte) (bool, error) {
	//TODO: put implementation here
	return false, nil
}

func (m *MongoDatabase) Get(key []byte, val interface{}) (interface{}, error) {
	//TODO: put implementation here
	return nil, nil
}

func (m *MongoDatabase) Put(key []byte, val interface{}) error {
	//TODO: put implementation here
	return nil
}

func (m *MongoDatabase) Delete(key []byte, force bool) error {
	//TODO: put implementation here
	return nil
}

func (m *MongoDatabase) Commit() error {
	//TODO: put implementation here
	return nil
}