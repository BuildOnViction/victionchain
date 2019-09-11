package tomox

import (
	"fmt"
	"math/big"
	"strings"
	// rbt "github.com/emirpasic/gods/trees/redblacktree"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type OrderTreeItem struct {
	Volume        *big.Int `json:"volume"`        // Contains total quantity from all Orders in tree
	NumOrders     uint64   `json:"numOrders"`     // Contains count of Orders in tree
	PriceTreeKey  []byte   `json:"priceTreeKey"`  // Root Key of price tree
	PriceTreeSize uint64   `json:"priceTreeSize"` // Number of nodes, currently it is Depth
}

type OrderTreeItemBSON struct {
	Volume        string `json:"volume" bson:"volume"`               // Contains total quantity from all Orders in tree
	NumOrders     string `json:"numOrders" bson:"numOrders"`         // Contains count of Orders in tree
	PriceTreeKey  string `json:"priceTreeKey"`                       // Root Key of price tree
	PriceTreeSize string `json:"priceTreeSize" bson:"priceTreeSize"` // Number of nodes, currently it is Depth
}

type OrderTreeItemRecord struct {
	Key   string
	Value *OrderTreeItem
}

type OrderTreeItemRecordBSON struct {
	Key   string
	Value *OrderTreeItemBSON
}

// OrderTree : order tree structure for travelling
type OrderTree struct {
	PriceTree *RedBlackTreeExtended `json:"priceTree"`
	orderBook *OrderBook
	orderDB   OrderDao // this is for order
	slot      *big.Int
	Key       []byte
	Item      *OrderTreeItem
}

// NewOrderTree create new order tree
func NewOrderTree(orderDB OrderDao, key []byte, orderBook *OrderBook) *OrderTree {
	priceTree := NewRedBlackTreeExtended(orderDB)
	item := &OrderTreeItem{
		Volume:    Zero(),
		NumOrders: 0,
		// Depth:     0,
		PriceTreeSize: 0,
	}

	slot := new(big.Int).SetBytes(key)

	// we will need a lru for cache hit, and internal cache for orderbook db to do the batch update
	orderTree := &OrderTree{
		orderDB:   orderDB,
		PriceTree: priceTree,
		Key:       key,
		slot:      slot,
		Item:      item,
		orderBook: orderBook,
	}

	return orderTree
}

func (orderTree *OrderTree) Save(dryrun bool, blockHash common.Hash) error {

	// update tree meta information, make sure item existed instead of checking rootKey
	priceTreeRoot := orderTree.PriceTree.Root(dryrun, blockHash)
	orderTree.Item.PriceTreeSize = orderTree.Depth()

	if priceTreeRoot != nil {
		orderTree.Item.PriceTreeKey = priceTreeRoot.Key
	} else if orderTree.Depth() == 0 {
		orderTree.Item.PriceTreeKey = EmptyKey()
		orderTree.Item.PriceTreeSize = 0
		orderTree.PriceTree.Clear()
	} else {
		return fmt.Errorf("Ordertree root not found")
	}

	log.Debug("Save ordertree", "key", hex.EncodeToString(orderTree.GetCommonKey()), "order.Item", orderTree.Item)
	return orderTree.orderDB.Put(orderTree.GetCommonKey(), orderTree.Item, dryrun, blockHash)
}

func (orderTree *OrderTree) Restore(dryrun bool, blockHash common.Hash) error {
	val, err := orderTree.orderDB.Get(orderTree.GetCommonKey(), orderTree.Item, dryrun, blockHash)

	if err == nil {
		orderTree.Item = val.(*OrderTreeItem)

		// update root key for pricetree
		orderTree.PriceTree.SetRootKey(orderTree.Item.PriceTreeKey, orderTree.Item.PriceTreeSize)
	}

	return err
}

func (orderTree *OrderTree) String(startDepth int, dryrun bool, blockHash common.Hash) string {
	tabs := strings.Repeat("\t", startDepth)
	return fmt.Sprintf("{\n\t%sMinPriceList: %s\n\t%sMaxPriceList: %s\n\t%sVolume: %v\n\t%sNumOrders: %d\n\t%sDepth: %d\n%s}",
		tabs, orderTree.MinPriceList(dryrun, blockHash).String(startDepth+1, dryrun, blockHash), tabs, orderTree.MaxPriceList(dryrun, blockHash).String(startDepth+1, dryrun, blockHash), tabs,
		orderTree.Item.Volume, tabs, orderTree.Item.NumOrders, tabs, orderTree.Depth(), tabs)
}

func (orderTree *OrderTree) Length() uint64 {
	return orderTree.Item.NumOrders
}

// Check the order database is emtpy or not
func (orderTree *OrderTree) NotEmpty() bool {
	return orderTree.Item.NumOrders > 0
}

func (orderTree *OrderTree) GetOrder(key []byte, price *big.Int, dryrun bool, blockHash common.Hash) *Order {
	orderList := orderTree.PriceList(price, dryrun, blockHash)
	if orderList == nil {
		return nil
	}

	// we can use orderID incremental way, so we just need a big slot from price of order tree
	return orderList.GetOrder(key, dryrun, blockHash)
}

func (orderTree *OrderTree) getSlotFromPrice(price *big.Int) *big.Int {
	return Add(orderTree.slot, price)
}

// next time this price will be big.Int
func (orderTree *OrderTree) getKeyFromPrice(price *big.Int) []byte {
	orderListKey := orderTree.getSlotFromPrice(price)
	return GetKeyFromBig(orderListKey)
}

// PriceList : get the price list from the price map using price as key
func (orderTree *OrderTree) PriceList(price *big.Int, dryrun bool, blockHash common.Hash) *OrderList {

	key := orderTree.getKeyFromPrice(price)
	bytes, found := orderTree.PriceTree.Get(GetOrderListCommonKey(key, orderTree.orderBook.Item.Name), dryrun, blockHash)
	log.Debug("Got orderlist by price", "key", hex.EncodeToString(key), "found", found)

	if found {
		orderList, err := orderTree.decodeOrderList(bytes)
		if err != nil {
			log.Error("Can't get price list", "price", price, "err", err)
			return nil
		}
		log.Debug("Decoded orderlist", "orderlist", orderList)
		return orderList
	}

	return nil
}

// CreatePrice : create new price list into PriceTree and PriceMap
func (orderTree *OrderTree) CreatePrice(price *big.Int, dryrun bool, blockHash common.Hash) *OrderList {

	// orderTree.Item.Depth++
	newList := NewOrderList(price, orderTree)
	// put new price list into tree
	newList.Save(dryrun, blockHash)

	// should use batch to optimize the performance
	orderTree.Save(dryrun, blockHash)

	return newList
}

func (orderTree *OrderTree) SaveOrderList(orderList *OrderList, dryrun bool, blockHash common.Hash) error {
	value, err := EncodeBytesItem(orderList.Item)
	if err != nil {
		log.Error("Can't encode", "orderList.Item", orderList.Item, "err", err)
		return err
	}
	olKey := GetOrderListCommonKey(orderList.Key, orderTree.orderBook.Item.Name)
	log.Debug("Save orderlist", "key", hex.EncodeToString(olKey), "value", value)
	return orderTree.PriceTree.Put(olKey, value, dryrun, blockHash)
}

func (orderTree *OrderTree) Depth() uint64 {
	return orderTree.PriceTree.Size()
}

// RemovePrice : delete a list by price
func (orderTree *OrderTree) RemovePrice(price *big.Int, dryrun bool, blockHash common.Hash) error {
	if orderTree.Depth() > 0 {
		priceKey := orderTree.getKeyFromPrice(price)
		log.Debug("Remove price", "price", price.String(), "priceKey", hex.EncodeToString(priceKey))
		orderListKey := GetOrderListCommonKey(priceKey, orderTree.orderBook.Item.Name)
		log.Debug("Remove price", "price", price.String(), "orderListKey", hex.EncodeToString(orderListKey))
		orderTree.PriceTree.Remove(orderListKey, dryrun, blockHash)
		log.Debug("Removed price from price tree", "price", price.String(), "orderListKey", hex.EncodeToString(orderListKey))
		// should use batch to optimize the performance
		if err := orderTree.Save(dryrun, blockHash); err != nil {
			return err
		}
		log.Debug("Remove price - saved ordertree")
	}
	return nil
}

// PriceExist : check price existed
func (orderTree *OrderTree) PriceExist(price *big.Int, dryrun bool, blockHash common.Hash) bool {

	priceKey := orderTree.getKeyFromPrice(price)
	orderListKey := GetOrderListCommonKey(priceKey, orderTree.orderBook.Item.Name)
	found, _ := orderTree.PriceTree.Has(orderListKey, dryrun, blockHash)

	return found
}

// OrderExist : check order existed, only support for a specific price
func (orderTree *OrderTree) OrderExist(key []byte, price *big.Int, dryrun bool, blockHash common.Hash) bool {
	orderList := orderTree.PriceList(price, dryrun, blockHash)
	if orderList == nil {
		return false
	}

	return orderList.OrderExist(key, dryrun, blockHash)
}

func (orderTree *OrderTree) InsertOrder(order *OrderItem, dryrun bool, blockHash common.Hash) error {

	price := order.Price

	var orderList *OrderList

	if !orderTree.PriceExist(price, dryrun, blockHash) {
		// create and save
		log.Debug("create price list", "detail", price.String())
		orderList = orderTree.CreatePrice(price, dryrun, blockHash)
	} else {
		orderList = orderTree.PriceList(price, dryrun, blockHash)
	}

	// order will be inserted to order list
	if orderList != nil {

		order := NewOrder(order, GetOrderListCommonKey(orderList.Key, orderTree.orderBook.Item.Name))

		if orderList.OrderExist(order.Key, dryrun, blockHash) {
			if err := orderTree.RemoveOrder(order, dryrun, blockHash); err != nil {
				return err
			}
		}

		// append order to order list
		log.Debug("debug orderlist", "before", orderList.Item)
		if err := orderList.AppendOrder(order, dryrun, blockHash); err != nil {
			return err
		}

		log.Debug("debug orderlist", "after", orderList.Item)
		// snapshot order list
		if err := orderList.Save(dryrun, blockHash); err != nil {
			return err
		}
		orderTree.Item.Volume = Add(orderTree.Item.Volume, order.Item.Quantity)

		// increase num of orders. FIXME: should be big.Int ?
		orderTree.Item.NumOrders++

		// finally, snapshot order tree
		if err := orderTree.Save(dryrun, blockHash); err != nil {
			return err
		}
	}

	return nil
}

// UpdateOrder : update an order
func (orderTree *OrderTree) UpdateOrder(orderItem *OrderItem, dryrun bool, blockHash common.Hash) error {

	price := orderItem.Price
	orderList := orderTree.PriceList(price, dryrun, blockHash)

	if orderList == nil {
		// create a price list for this price
		orderList = orderTree.CreatePrice(price, dryrun, blockHash)
	}

	orderID := new(big.Int).SetUint64(orderItem.OrderID)
	key := GetKeyFromBig(orderID)

	order := orderList.GetOrder(key, dryrun, blockHash)

	originalQuantity := CloneBigInt(order.Item.Quantity)

	if !IsEqual(price, order.Item.Price) {
		if err := orderList.RemoveOrder(order, dryrun, blockHash); err != nil {
			return err
		}
		if orderList.Item.Length == uint64(0) {
			if err := orderTree.RemovePrice(price, dryrun, blockHash); err != nil {
				return err
			}
		}
		orderTree.InsertOrder(orderItem, dryrun, blockHash)
	} else {
		quantity := orderItem.Quantity
		//timestamp, _ := strconv.ParseInt(quote["timestamp"], 10, 64)
		timestamp := orderItem.CreatedAt
		if err := order.UpdateQuantity(orderList, quantity, timestamp, dryrun, blockHash); err != nil {
			return err
		}
	}

	orderTree.Item.Volume = Add(orderTree.Item.Volume, Sub(order.Item.Quantity, originalQuantity))

	// should use batch to optimize the performance
	return orderTree.Save(dryrun, blockHash)
}

func (orderTree *OrderTree) RemoveOrderFromOrderList(order *Order, orderList *OrderList, dryrun bool, blockHash common.Hash) error {
	// next update orderList
	if err := orderList.RemoveOrder(order, dryrun, blockHash); err != nil {
		return err
	}
	log.Debug("Removed order from orderlist", "order", order.Item, "orderlist", orderList.Item)

	// snapshot order list
	if err := orderList.Save(dryrun, blockHash); err != nil {
		return err
	}
	log.Debug("Saved orderlist", "orderlist", orderList.Item)

	// no items left than safety remove
	if orderList.Item.Length == uint64(0) {
		if err := orderTree.RemovePrice(order.Item.Price, dryrun, blockHash); err != nil {
			return err
		}
		log.Debug("Removed price list", "price", order.Item.Price.String())
	}

	// update orderTree
	orderTree.Item.Volume = Sub(orderTree.Item.Volume, order.Item.Quantity)

	// delete(orderTree.OrderMap, orderID)
	orderTree.Item.NumOrders--

	// should use batch to optimize the performance
	return nil
}

func (orderTree *OrderTree) RemoveOrder(order *Order, dryrun bool, blockHash common.Hash) error {
	// get orderList by price. If there is orderlist existed, update it
	orderList := orderTree.PriceList(order.Item.Price, dryrun, blockHash)
	if orderList == nil {
		return fmt.Errorf("Orderlist not found")
	}
	if orderList != nil {
		if err := orderTree.RemoveOrderFromOrderList(order, orderList, dryrun, blockHash); err != nil {
			return err
		}

		// snapshot ordertree
		if err := orderTree.Save(dryrun, blockHash); err != nil {
			return err
		}
	}

	return nil
}

func (orderTree *OrderTree) getOrderListItem(bytes []byte) (*OrderListItem, error) {
	item := &OrderListItem{}
	// rlp.DecodeBytes(bytes, item)
	//orderTree.orderDB.DecodeBytes(bytes, item)
	err := DecodeBytesItem(bytes, item)
	if err != nil {
		log.Error("Can't decode", "bytes", bytes, "item", item)
		return nil, err
	}
	return item, nil
}

func (orderTree *OrderTree) decodeOrderList(bytes []byte) (*OrderList, error) {
	item, err := orderTree.getOrderListItem(bytes)
	if err != nil {
		return nil, err
	}
	orderList := NewOrderListWithItem(item, orderTree)

	return orderList, nil
}

// MaxPrice : get the max price
func (orderTree *OrderTree) MaxPrice(dryrun bool, blockHash common.Hash) *big.Int {
	if orderTree.Depth() > 0 {
		if bytes, found := orderTree.PriceTree.GetMax(dryrun, blockHash); found {
			item, err := orderTree.getOrderListItem(bytes)
			if err != nil {
				return Zero()
			}
			if item != nil {
				return CloneBigInt(item.Price)
			}
		}
	}
	return Zero()
}

// MinPrice : get the min price
func (orderTree *OrderTree) MinPrice(dryrun bool, blockHash common.Hash) *big.Int {
	if orderTree.Depth() > 0 {
		log.Debug("ordertree isn't empty")
		if bytes, found := orderTree.PriceTree.GetMin(dryrun, blockHash); found {
			item, err := orderTree.getOrderListItem(bytes)
			if err != nil {
				log.Error("Failed to get orderlist with min price", "err", err)
				return Zero()
			}
			if item != nil {
				return CloneBigInt(item.Price)
			}
		} else {
			log.Debug("Min not found")
		}
	}
	return Zero()
}

// MaxPriceList : get max price list
func (orderTree *OrderTree) MaxPriceList(dryrun bool, blockHash common.Hash) *OrderList {
	if orderTree.Depth() > 0 {
		if bytes, found := orderTree.PriceTree.GetMax(dryrun, blockHash); found {
			orderList, _ := orderTree.decodeOrderList(bytes)
			return orderList
		}
	}
	return nil

}

// MinPriceList : get min price list
func (orderTree *OrderTree) MinPriceList(dryrun bool, blockHash common.Hash) *OrderList {
	if orderTree.Depth() > 0 {
		if bytes, found := orderTree.PriceTree.GetMin(dryrun, blockHash); found {
			orderList, _ := orderTree.decodeOrderList(bytes)
			return orderList
		}
	}
	return nil
}

func (orderTree *OrderTree) Hash() (common.Hash, error) {
	olEncoded, err := EncodeBytesItem(orderTree.Item)
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(olEncoded), nil
}

func (orderTree *OrderTree) GetCommonKey() []byte {
	return append([]byte("OT"), orderTree.Key...)
}
