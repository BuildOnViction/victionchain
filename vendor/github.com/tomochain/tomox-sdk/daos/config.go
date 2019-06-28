package daos

import (
	"strconv"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

const (
	schemaVersionKey        = "swap_schema_version"
	ethereumAddressIndexKey = "ethereum_address_index"
	ethereumLastBlockKey    = "ethereum_last_block"
	bitcoinAddressIndexKey  = "bitcoin_address_index"
	bitcoinLastBlockKey     = "bitcoin_last_block"
	defaultBlockIndex       = 0
)

// ConfigDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type ConfigDao struct {
	collectionName string
	dbName         string
}

// NewBalanceDao returns a new instance of AddressDao
func NewConfigDao() *ConfigDao {
	dbName := app.Config.DBName
	// we save deposit information and use config for retrieving params.
	collection := "config"
	index := mgo.Index{
		Key:    []string{"key"},
		Unique: true,
	}

	err := db.Session.DB(dbName).C(collection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return &ConfigDao{collection, dbName}
}

func (dao *ConfigDao) GetSchemaVersion() uint64 {
	// get version
	value, err := dao.getUint64ValueFromKey(schemaVersionKey)
	if err != nil {
		return 0
	}
	return value
}

func (dao *ConfigDao) getAddressIndexKey(chain types.Chain) (string, error) {
	switch chain {
	case types.ChainEthereum:
		return ethereumAddressIndexKey, nil
	case types.ChainBitcoin:
		return bitcoinAddressIndexKey, nil
	default:
		return "", errors.New("Invalid chain")
	}
}

func (dao *ConfigDao) getValueFromKey(key string) (interface{}, error) {
	var response types.KeyValue
	err := db.GetOne(dao.dbName, dao.collectionName, bson.M{"key": key}, &response)
	if err != nil {
		logger.Errorf("Got error: %v", err)
		return nil, errors.Errorf("Value not found for key: %s", key)
	}

	return response.Value, nil
}

func (dao *ConfigDao) getUint64ValueFromKey(key string) (uint64, error) {
	value, err := dao.getValueFromKey(key)
	if err != nil {
		return 0, err
	}
	switch v := value.(type) {
	case int:
		return uint64(value.(int)), nil
	case int64:
		return uint64(value.(int64)), nil
	case string:
		return strconv.ParseUint(value.(string), 10, 64)
	default:
		return 0, errors.Errorf("Can not process type %T!\n", v)
	}
}

func (dao *ConfigDao) GetAddressIndex(chain types.Chain) (uint64, error) {
	key, err := dao.getAddressIndexKey(chain)
	if err != nil {
		return 0, err
	}

	return dao.getUint64ValueFromKey(key)
}

func (dao *ConfigDao) IncrementAddressIndex(chain types.Chain) error {
	// update database
	key, err := dao.getAddressIndexKey(chain)
	if err != nil {
		return err
	}

	err = db.Update(dao.dbName, dao.collectionName, bson.M{"key": key}, bson.M{
		"$inc": bson.M{
			"value": 1,
		},
	})

	return err
}

func (dao *ConfigDao) GetBlockToProcess(chain types.Chain) (uint64, error) {
	switch chain {
	case types.ChainEthereum:
		return dao.getUint64ValueFromKey(ethereumLastBlockKey)
	case types.ChainBitcoin:
		return dao.getUint64ValueFromKey(bitcoinLastBlockKey)
	default:
		return 0, errors.New("Invalid chain")
	}

}

func (dao *ConfigDao) SaveLastProcessedBlock(chain types.Chain, block uint64) error {
	// update database
	var key string
	switch chain {
	case types.ChainEthereum:
		key = ethereumLastBlockKey
	case types.ChainBitcoin:
		key = bitcoinLastBlockKey
	default:
		return errors.New("Invalid chain")
	}

	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"key": key}, bson.M{
		"$set": bson.M{
			"value": block,
		},
	})

	return err
}

// Drop drops all the order documents in the current database
func (dao *ConfigDao) Drop() {
	db.DropCollection(dao.dbName, dao.collectionName)
}

// ResetBlockCounters changes last processed bitcoin and ethereum block to default value.
// Used in stress tests.
func (dao *ConfigDao) ResetBlockCounters() error {
	err := dao.SaveLastProcessedBlock(types.ChainEthereum, defaultBlockIndex)
	if err != nil {
		return err
	}
	return dao.SaveLastProcessedBlock(types.ChainBitcoin, defaultBlockIndex)
}
