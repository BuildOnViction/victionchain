package daos

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/types"
)

type NotificationDao struct {
	collectionName string
	dbName         string
}

func NewNotificationDao() *NotificationDao {
	dbName := app.Config.DBName
	collection := "notifications"

	return &NotificationDao{
		collectionName: collection,
		dbName:         dbName,
	}
}

// Create function performs the DB insertion task for notification collection
// It accepts 1 or more notifications as input.
// All the notifications are inserted in one query itself.
func (dao *NotificationDao) Create(notifications ...*types.Notification) ([]*types.Notification, error) {
	y := make([]interface{}, len(notifications))
	result := make([]*types.Notification, len(notifications))

	for _, notification := range notifications {
		notification.ID = bson.NewObjectId()
		notification.CreatedAt = time.Now()
		notification.UpdatedAt = time.Now()
		y = append(y, notification)
		result = append(result, notification)
	}

	err := db.Create(dao.dbName, dao.collectionName, y...)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return result, nil
}

// GetAll function fetches all the notifications in the notification collection of mongodb.
func (dao *NotificationDao) GetAll() ([]types.Notification, error) {
	var response []types.Notification

	err := db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return response, nil
}

// GetByUserAddress function fetches list of orders from order collection based on user address.
// Returns array of Order type struct
func (dao *NotificationDao) GetByUserAddress(addr common.Address, limit int, offset int) ([]*types.Notification, error) {
	if limit == 0 {
		limit = 10 // Get last 10 records
	}

	var res []*types.Notification
	q := bson.M{"recipient": addr.Hex()}

	err := db.Get(dao.dbName, dao.collectionName, q, offset, limit, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if res == nil {
		return []*types.Notification{}, nil
	}

	return res, nil
}

// GetByID function fetches details of a notification based on its mongo id
func (dao *NotificationDao) GetByID(id bson.ObjectId) (*types.Notification, error) {
	var response *types.Notification

	err := db.GetByID(dao.dbName, dao.collectionName, id, &response)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return response, nil
}

func (dao *NotificationDao) FindAndModify(id bson.ObjectId, n *types.Notification) (*types.Notification, error) {
	n.UpdatedAt = time.Now()
	query := bson.M{"_id": id}
	updated := &types.Notification{}
	change := mgo.Change{
		Update:    types.NotificationBSONUpdate{Notification: n},
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

func (dao *NotificationDao) Update(n *types.Notification) error {
	n.UpdatedAt = time.Now()
	err := db.Update(dao.dbName, dao.collectionName, bson.M{"_id": n.ID}, n)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *NotificationDao) Upsert(id bson.ObjectId, n *types.Notification) error {
	n.UpdatedAt = time.Now()

	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"_id": id}, n)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *NotificationDao) Delete(notifications ...*types.Notification) error {
	ids := make([]bson.ObjectId, 0)
	for _, n := range notifications {
		ids = append(ids, n.ID)
	}

	err := db.RemoveAll(dao.dbName, dao.collectionName, bson.M{"_id": bson.M{"$in": ids}})

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *NotificationDao) DeleteByIds(ids ...bson.ObjectId) error {
	err := db.RemoveAll(dao.dbName, dao.collectionName, bson.M{"_id": bson.M{"$in": ids}})

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Aggregate function calls the aggregate pipeline of mongodb
func (dao *NotificationDao) Aggregate(q []bson.M) ([]*types.Notification, error) {
	var res []*types.Notification

	err := db.Aggregate(dao.dbName, dao.collectionName, q, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// Drop drops all the order documents in the current database
func (dao *NotificationDao) Drop() {
	db.DropCollection(dao.dbName, dao.collectionName)
}
