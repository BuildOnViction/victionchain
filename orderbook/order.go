package orderbook

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
)

// OrderItem : info that will be store in database
type OrderItem struct {
	Timestamp uint64   `json:"timestamp"`
	Quantity  *big.Int `json:"quantity"`
	Price     *big.Int `json:"price"`
	TradeID string `json:"tradeID"`
	// *OrderMeta
	NextOrder []byte `json:"-"`
	PrevOrder []byte `json:"-"`
	OrderList []byte `json:"-"`
}

type Order struct {
	Item *OrderItem
	Key  []byte `json:"orderID"`
}

func (order *Order) String() string {

	return fmt.Sprintf("orderID : %s, price: %s, quantity :%s, tradeID: %s",
		new(big.Int).SetBytes(order.Key), order.Item.Price, order.Item.Quantity, order.Item.TradeID)
}

func (order *Order) GetNextOrder(orderList *OrderList) *Order {
	nextOrder := orderList.GetOrder(order.Item.NextOrder)

	return nextOrder
}

func (order *Order) GetPrevOrder(orderList *OrderList) *Order {
	prevOrder := orderList.GetOrder(order.Item.PrevOrder)

	return prevOrder
}

// NewOrder : create new order with quote ( can be ethereum address )
func NewOrder(quote map[string]string, orderList []byte) *Order {
	timestamp, _ := strconv.ParseUint(quote["timestamp"], 10, 64)
	quantity := ToBigInt(quote["quantity"])
	price := ToBigInt(quote["price"])
	orderID := ToBigInt(quote["order_id"])
	key := GetKeyFromBig(orderID)
	tradeID := quote["trade_id"]
	orderItem := &OrderItem{
		Timestamp: timestamp,
		Quantity:  quantity,
		Price:     price,
		TradeID:   tradeID,
		NextOrder: EmptyKey(),
		PrevOrder: EmptyKey(),
		OrderList: orderList,
	}

	// key should be Hash for compatible with smart contract
	order := &Order{
		Key:  key,
		Item: orderItem,
	}

	return order
}

// UpdateQuantity : update quantity of the order
func (order *Order) UpdateQuantity(orderList *OrderList, newQuantity *big.Int, newTimestamp uint64) {
	if newQuantity.Cmp(order.Item.Quantity) > 0 && !bytes.Equal(orderList.Item.TailOrder, order.Key) {
		orderList.MoveToTail(order)
	}
	// update volume and modified timestamp
	orderList.Item.Volume = Sub(orderList.Item.Volume, Sub(order.Item.Quantity, newQuantity))
	order.Item.Timestamp = newTimestamp
	order.Item.Quantity = CloneBigInt(newQuantity)
	fmt.Println("QUANTITY", order.Item.Quantity.String())
	orderList.SaveOrder(order)
	orderList.Save()
}
