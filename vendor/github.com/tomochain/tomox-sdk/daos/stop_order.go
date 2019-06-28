package daos

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/types"
)

// StopOrderDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type StopOrderDao struct {
	collectionName string
	dbName         string
}

type StopOrderDaoOption = func(*StopOrderDao) error

// NewOrderDao returns a new instance of OrderDao
func NewStopOrderDao(opts ...StopOrderDaoOption) *StopOrderDao {
	dao := &StopOrderDao{}
	dao.collectionName = "stop_orders"
	dao.dbName = app.Config.DBName

	for _, op := range opts {
		err := op(dao)
		if err != nil {
			panic(err)
		}
	}

	index := mgo.Index{
		Key:    []string{"hash"},
		Unique: true,
	}

	i1 := mgo.Index{
		Key: []string{"userAddress"},
	}

	i2 := mgo.Index{
		Key: []string{"status"},
	}

	i3 := mgo.Index{
		Key: []string{"baseToken"},
	}

	i4 := mgo.Index{
		Key: []string{"quoteToken"},
	}

	i5 := mgo.Index{
		Key:       []string{"stopPrice"},
		Collation: &mgo.Collation{NumericOrdering: true, Locale: "en"},
	}

	i6 := mgo.Index{
		Key: []string{"baseToken", "quoteToken", "stopPrice"},
	}

	i7 := mgo.Index{
		Key: []string{"side", "status"},
	}

	i8 := mgo.Index{
		Key: []string{"baseToken", "quoteToken", "side", "status"},
	}

	err := db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i1)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i2)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i3)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i4)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i5)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i6)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i7)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i8)
	if err != nil {
		panic(err)
	}

	return dao
}

// Create function performs the DB insertion task for Order collection
func (dao *StopOrderDao) Create(o *types.StopOrder) error {
	o.ID = bson.NewObjectId()
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()

	if o.Status == "" {
		o.Status = types.OrderStatusOpen
	}

	err := db.Create(dao.dbName, dao.collectionName, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Update function performs the DB updations task for Order collection
// corresponding to a particular order ID
func (dao *StopOrderDao) Update(id bson.ObjectId, so *types.StopOrder) error {
	so.UpdatedAt = time.Now()

	err := db.Update(dao.dbName, dao.collectionName, bson.M{"_id": id}, so)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

//UpdateByHash updates fields that are considered updateable for an order.
func (dao *StopOrderDao) UpdateByHash(h common.Hash, so *types.StopOrder) error {
	so.UpdatedAt = time.Now()
	query := bson.M{"hash": h.Hex()}
	update := bson.M{"$set": bson.M{
		"stopPrice":    so.StopPrice.String(),
		"limitPrice":   so.LimitPrice.String(),
		"amount":       so.Amount.String(),
		"status":       so.Status,
		"filledAmount": so.FilledAmount.String(),
		"makeFee":      so.MakeFee.String(),
		"takeFee":      so.TakeFee.String(),
		"updatedAt":    so.UpdatedAt,
	}}

	err := db.Update(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *StopOrderDao) Upsert(id bson.ObjectId, o *types.StopOrder) error {
	o.UpdatedAt = time.Now()

	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"_id": id}, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *StopOrderDao) UpsertByHash(h common.Hash, so *types.StopOrder) error {
	_, err := db.Upsert(
		dao.dbName,
		dao.collectionName,
		bson.M{"hash": h.Hex()},
		types.StopOrderBSONUpdate{
			StopOrder: so,
		},
	)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *StopOrderDao) UpdateAllByHash(h common.Hash, so *types.StopOrder) error {
	so.UpdatedAt = time.Now()

	err := db.Update(dao.dbName, dao.collectionName, bson.M{"hash": h.Hex()}, so)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// GetByHash function fetches a single document from stop_order collection based on the hash.
// Returns StopOrder type struct
func (dao *StopOrderDao) GetByHash(hash common.Hash) (*types.StopOrder, error) {
	q := bson.M{"hash": hash.Hex()}
	res := []types.StopOrder{}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &res[0], nil
}

func (dao *StopOrderDao) FindAndModify(h common.Hash, so *types.StopOrder) (*types.StopOrder, error) {
	so.UpdatedAt = time.Now()
	query := bson.M{"hash": h.Hex()}
	updated := &types.StopOrder{}
	change := mgo.Change{
		Update: types.StopOrderBSONUpdate{
			StopOrder: so,
		},
		Upsert:    true,
		Remove:    false,
		ReturnNew: true,
	}

	err := db.FindAndModify(dao.dbName, dao.collectionName, query, change, &updated)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return updated, nil
}

func (dao *StopOrderDao) GetTriggeredStopOrders(baseToken, quoteToken common.Address, lastPrice *big.Int) ([]*types.StopOrder, error) {
	var stopOrders []*types.StopOrder

	q := bson.M{
		"$or": []bson.M{
			bson.M{
				"baseToken":  baseToken.Hex(),
				"quoteToken": quoteToken.Hex(),
				"direction":  -1,
				"stopPrice": bson.M{
					"$gte": lastPrice,
				},
				"status": types.StopOrderStatusOpen,
			},
			bson.M{
				"baseToken":  baseToken.Hex(),
				"quoteToken": quoteToken.Hex(),
				"direction":  1,
				"stopPrice": bson.M{
					"$lte": lastPrice,
				},
				"status": types.StopOrderStatusOpen,
			},
		},
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &stopOrders)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return stopOrders, nil
}

// Drop drops all the order documents in the current database
func (dao *StopOrderDao) Drop() error {
	err := db.DropCollection(dao.dbName, dao.collectionName)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
