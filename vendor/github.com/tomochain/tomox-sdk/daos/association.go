package daos

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/types"
)

// AssociationDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type AssociationDao struct {
	collectionName string
	dbName         string
}

// NewBalanceDao returns a new instance of AddressDao
func NewAssociationDao() *AssociationDao {
	dbName := app.Config.DBName
	// we save deposit information and use config for retrieving params.
	collection := "associations"
	index := mgo.Index{
		Key:    []string{"chain", "address"},
		Unique: true,
	}

	// chain and associatedAddress also is uniqued
	index1 := mgo.Index{
		Key:    []string{"chain", "associatedAddress"},
		Unique: true,
	}

	err := db.Session.DB(dbName).C(collection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dbName).C(collection).EnsureIndex(index1)
	if err != nil {
		panic(err)
	}

	return &AssociationDao{collection, dbName}
}

// return the lowercase of the key
func (dao *AssociationDao) getAddressKey(address common.Address) string {
	return "0x" + common.Bytes2Hex(address.Bytes())
}

// Drop drops all the order documents in the current database
func (dao *AssociationDao) Drop() {
	db.DropCollection(dao.dbName, dao.collectionName)
}

// SaveDepositTransaction update the transaction envelope for association item
func (dao *AssociationDao) SaveDepositTransaction(chain types.Chain, sourceAccount common.Address, txEnvelope string) error {
	// txEnvolope is rlp of result
	err := db.Update(dao.dbName, dao.collectionName, bson.M{
		"chain":   chain.String(),
		"address": dao.getAddressKey(sourceAccount),
	}, bson.M{
		"$addToSet": bson.M{
			"txEnvelopes": txEnvelope,
		},
	})
	return err
}

func (dao *AssociationDao) GetAssociationByChainAddress(chain types.Chain, userAddress common.Address) (*types.AddressAssociationRecord, error) {
	var response types.AddressAssociationRecord
	err := db.GetOne(dao.dbName, dao.collectionName, bson.M{
		"chain":   chain.String(),
		"address": dao.getAddressKey(userAddress),
	}, &response)

	// if not found, just return nil instead of error
	if err == mgo.ErrNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (dao *AssociationDao) GetAssociationByChainAssociatedAddress(chain types.Chain, associatedAddress common.Address) (*types.AddressAssociationRecord, error) {
	var response types.AddressAssociationRecord
	err := db.GetOne(dao.dbName, dao.collectionName, bson.M{
		"chain":             chain.String(),
		"associatedAddress": dao.getAddressKey(associatedAddress),
	}, &response)

	// if not found, just return nil instead of error
	if err == mgo.ErrNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SaveAssociation using upsert to update for existing users, only update allowed fields
func (dao *AssociationDao) SaveAssociation(record *types.AddressAssociationRecord) error {
	associatedAddress := strings.ToLower(record.AssociatedAddress)
	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{
		"chain":             record.Chain,
		"associatedAddress": associatedAddress,
	}, bson.M{
		"$set": bson.M{
			"associatedAddress": associatedAddress,
			"chain":             record.Chain,
			"addressIndex":      record.AddressIndex,
			"address":           strings.ToLower(record.Address),
			"pairName":          record.PairName,
			"status":            record.Status,
			"quoteTokenAddress": record.QuoteTokenAddress,
			"baseTokenAddress":  record.BaseTokenAddress,
		},
	})
	return err
}

func (dao *AssociationDao) SaveAssociationStatus(chain types.Chain, sourceAccount common.Address, status string) error {
	err := db.Update(dao.dbName, dao.collectionName, bson.M{
		"chain":   chain.String(),
		"address": dao.getAddressKey(sourceAccount),
	}, bson.M{
		"$set": bson.M{
			"status": status,
		},
	})
	return err
}
