package tomox

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

const (
	LimitDepthPrint = 20
)

type OrderListItem struct {
	HeadOrder []byte   `json:"headOrder"`
	TailOrder []byte   `json:"tailOrder"`
	Length    uint64   `json:"length"`
	Volume    *big.Int `json:"volume"`
	Price     *big.Int `json:"price"`
}

type OrderListItemBSON struct {
	HeadOrder string `json:"headOrder" bson:"headOrder"`
	TailOrder string `json:"tailOrder" bson:"tailOrder"`
	Length    string `json:"length" bson:"length"`
	Volume    string `json:"volume" bson:"volume"`
	Price     string `json:"price" bson:"price"`
}

type OrderListItemRecord struct {
	Key   string
	Value *OrderListItem
}

type OrderListItemRecordBSON struct {
	Key   string
	Value *OrderListItemBSON
}

// OrderList : order list
type OrderList struct {
	orderTree *OrderTree
	slot      *big.Int
	Item      *OrderListItem
	Key       []byte
}

// NewOrderList : return new OrderList
// each orderlist will store information of order in a seperated domain
func NewOrderList(price *big.Int, orderTree *OrderTree) *OrderList {
	item := &OrderListItem{
		HeadOrder: EmptyKey(),
		TailOrder: EmptyKey(),
		Length:    uint64(0),
		Volume:    Zero(),
		Price:     CloneBigInt(price),
	}

	return NewOrderListWithItem(item, orderTree)

}

func NewOrderListWithItem(item *OrderListItem, orderTree *OrderTree) *OrderList {
	key := orderTree.getKeyFromPrice(item.Price)

	orderList := &OrderList{
		Item:      item,
		Key:       key,
		orderTree: orderTree,
	}

	// priceKey will be slot of order tree + plus price key
	// we can use orderList slot as orderbook slot to store sequential of orders
	if orderTree.orderBook != nil {
		orderList.slot = orderTree.orderBook.Slot
	} else {
		orderList.slot = new(big.Int).SetBytes(crypto.Keccak256(key))
	}

	return orderList
}

func (orderList *OrderList) GetOrder(key []byte, dryrun bool) *Order {
	// re-use method from orderbook, because orderlist has the same slot as orderbook
	storedKey := orderList.GetOrderIDFromKey(key)
	log.Debug("Get order from key", "storedKey", storedKey)
	return orderList.orderTree.orderBook.GetOrder(storedKey, key, dryrun)
}

func (orderList *OrderList) isEmptyKey(key []byte) bool {
	return orderList.orderTree.PriceTree.IsEmptyKey(key)
}

func (orderList *OrderList) Head(dryrun bool) *Order {
	return orderList.GetOrder(orderList.Item.HeadOrder, dryrun)
}

func (orderList *OrderList) Tail(dryrun bool) *Order {
	return orderList.GetOrder(orderList.Item.TailOrder, dryrun)
}

// String : travel the list to print it in nice format
func (orderList *OrderList) String(startDepth int, dryrun bool) string {

	if orderList == nil {
		return "<nil>"
	}

	var buffer bytes.Buffer
	tabs := strings.Repeat("\t", startDepth)
	buffer.WriteString(fmt.Sprintf("{\n\t%sLength: %d\n\t%sVolume: %v\n\t%sPrice: %v",
		tabs, orderList.Item.Length, tabs, orderList.Item.Volume, tabs, orderList.Item.Price))

	buffer.WriteString("\n\t")
	buffer.WriteString(tabs)
	buffer.WriteString("Head:")
	linkedList := orderList.Head(dryrun)
	depth := 0
	for linkedList != nil {
		depth++
		spaces := strings.Repeat(" ", depth)
		if depth > LimitDepthPrint {
			buffer.WriteString(fmt.Sprintf("\n\t%s%s |-> %s %d left", tabs, spaces, "...",
				orderList.Item.Length-LimitDepthPrint))
			break
		}
		buffer.WriteString(fmt.Sprintf("\n\t%s%s |-> %s", tabs, spaces, linkedList.String()))
		linkedList = orderList.GetOrder(linkedList.Item.NextOrder, dryrun)
	}
	if depth == 0 {
		buffer.WriteString(" <nil>")
	}
	buffer.WriteString("\n\t")
	buffer.WriteString(tabs)
	buffer.WriteString("Tail:")
	linkedList = orderList.Tail(dryrun)
	depth = 0
	for linkedList != nil {
		depth++
		spaces := strings.Repeat(" ", depth)
		if depth > LimitDepthPrint {
			buffer.WriteString(fmt.Sprintf("\n\t%s%s <-| %s %d left", tabs, spaces, "...",
				orderList.Item.Length-LimitDepthPrint))
			break
		}
		buffer.WriteString(fmt.Sprintf("\n\t%s%s <-| %s", tabs, spaces, linkedList.String()))
		linkedList = orderList.GetOrder(linkedList.Item.PrevOrder, dryrun)
	}
	if depth == 0 {
		buffer.WriteString(" <nil>")
	}
	buffer.WriteString("\n")
	buffer.WriteString(tabs)
	buffer.WriteString("}")
	return buffer.String()
}

// Less : compare if this order list is less than compared object
func (orderList *OrderList) Less(than *OrderList) bool {
	// cast to OrderList pointer
	return IsStrictlySmallerThan(orderList.Item.Price, than.Item.Price)
}

func (orderList *OrderList) Save(dryrun bool) error {
	return orderList.orderTree.SaveOrderList(orderList, dryrun)
}

// return the input orderID
func (orderList *OrderList) GetOrderIDFromList(key []byte) uint64 {
	orderSlot := new(big.Int).SetBytes(key)
	return Sub(orderSlot, orderList.slot).Uint64()
}

// GetOrderIDFromKey
// If we allow the same orderid belongs to many pricelist, we must use slot
// otherwise just use 1 db for storing all orders of all pricelists
// currently we use auto increase ment id so no need slot
func (orderList *OrderList) GetOrderIDFromKey(key []byte) []byte {
	orderSlot := new(big.Int).SetBytes(key)
	return common.BigToHash(Add(orderList.slot, orderSlot)).Bytes()
}

// GetOrderID return the real slot key of order in this linked list
func (orderList *OrderList) GetOrderID(order *Order) []byte {
	return orderList.GetOrderIDFromKey(order.Key)
}

// OrderExist search order in orderlist
func (orderList *OrderList) OrderExist(key []byte, dryrun bool) bool {
	orderKey := orderList.GetOrderIDFromKey(key)
	found, _ := orderList.orderTree.orderDB.Has(orderKey, dryrun)
	return found
}

func (orderList *OrderList) SaveOrder(order *Order, dryrun bool) error {
	key := orderList.GetOrderID(order)
	log.Debug("Save order ", "key", hex.EncodeToString(key), "value", ToJSON(order.Item))

	return orderList.orderTree.orderDB.Put(key, order.Item, dryrun)
}

// AppendOrder : append order into the order list
func (orderList *OrderList) AppendOrder(order *Order, dryrun bool) error {

	if orderList.Item.Length == uint64(0) {
		order.Item.NextOrder = EmptyKey()
		order.Item.PrevOrder = EmptyKey()
	} else {
		order.Item.PrevOrder = orderList.Item.TailOrder
		order.Item.NextOrder = EmptyKey()
	}

	// save into database first
	if err := orderList.SaveOrder(order, dryrun); err != nil {
		return err
	}

	if orderList.Item.Length == uint64(0) {
		orderList.Item.HeadOrder = order.Key
		orderList.Item.TailOrder = order.Key
	} else {
		tailOrder := orderList.GetOrder(orderList.Item.TailOrder, dryrun)
		if tailOrder != nil {
			tailOrder.Item.NextOrder = order.Key
			orderList.Item.TailOrder = order.Key
			if err := orderList.SaveOrder(tailOrder, dryrun); err != nil {
				return err
			}
		}
	}
	orderList.Item.Length++
	orderList.Item.Volume = Add(orderList.Item.Volume, order.Item.Quantity)
	return nil
}

func (orderList *OrderList) DeleteOrder(order *Order, dryrun bool) error {
	key := orderList.GetOrderID(order)
	return orderList.orderTree.orderDB.Delete(key, dryrun)
}

// RemoveOrder : remove order from the order list
func (orderList *OrderList) RemoveOrder(order *Order, dryrun bool) error {

	if orderList.Item.Length == uint64(0) {
		// empty mean nothing to delete
		return nil
	}

	if order.Item.Status == Cancel {
		// only CANCELLED order will be put back to DB
		if err := orderList.SaveOrder(order, dryrun); err != nil {
			return err
		}
	} else {
		if err := orderList.DeleteOrder(order, dryrun); err != nil {
			return err
		}
	}

	nextOrder := orderList.GetOrder(order.Item.NextOrder, dryrun)
	prevOrder := orderList.GetOrder(order.Item.PrevOrder, dryrun)

	orderList.Item.Volume = Sub(orderList.Item.Volume, order.Item.Quantity)
	orderList.Item.Length--

	if nextOrder != nil && prevOrder != nil {
		nextOrder.Item.PrevOrder = prevOrder.Key
		prevOrder.Item.NextOrder = nextOrder.Key

		if err := orderList.SaveOrder(nextOrder, dryrun); err != nil {
			return err
		}
		if err := orderList.SaveOrder(prevOrder, dryrun); err != nil {
			return err
		}
	} else if nextOrder != nil {
		// this might be wrong
		nextOrder.Item.PrevOrder = EmptyKey()
		orderList.Item.HeadOrder = nextOrder.Key

		if err := orderList.SaveOrder(nextOrder, dryrun); err != nil {
			return err
		}
	} else if prevOrder != nil {
		prevOrder.Item.NextOrder = EmptyKey()
		orderList.Item.TailOrder = prevOrder.Key

		if err := orderList.SaveOrder(prevOrder, dryrun); err != nil {
			return err
		}
	} else {
		// empty
		orderList.Item.HeadOrder = EmptyKey()
		orderList.Item.TailOrder = EmptyKey()
	}

	return nil
}

// MoveToTail : move order to the end of the order list
func (orderList *OrderList) MoveToTail(order *Order, dryrun bool) error {
	if !orderList.isEmptyKey(order.Item.PrevOrder) { // This Order is not the first Order in the OrderList
		prevOrder := orderList.GetOrder(order.Item.PrevOrder, dryrun)
		if prevOrder != nil {
			prevOrder.Item.NextOrder = order.Item.NextOrder // Link the previous Order to the next Order, then move the Order to tail
			if err := orderList.SaveOrder(prevOrder, dryrun); err != nil {
				return err
			}
		}

	} else { // This Order is the first Order in the OrderList
		orderList.Item.HeadOrder = order.Item.NextOrder // Make next order the first
	}

	nextOrder := orderList.GetOrder(order.Item.NextOrder, dryrun)
	if nextOrder != nil {
		nextOrder.Item.PrevOrder = order.Item.PrevOrder
		if err := orderList.SaveOrder(nextOrder, dryrun); err != nil {
			return err
		}
	}

	// Move Order to the last position. Link up the previous last position Order.
	tailOrder := orderList.GetOrder(orderList.Item.TailOrder, dryrun)
	if tailOrder != nil {
		tailOrder.Item.NextOrder = order.Key
		if err := orderList.SaveOrder(tailOrder, dryrun); err != nil {
			return err
		}
	}

	orderList.Item.TailOrder = order.Key
	return orderList.Save(dryrun)
}

func (orderList *OrderList) Hash() (common.Hash, error) {
	olEncoded, err := EncodeBytesItem(orderList.Item)
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(olEncoded), nil
}

func GetOrderListCommonKey(key []byte) []byte {
	return append([]byte("OL"), key...)
}
