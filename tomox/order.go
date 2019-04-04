package tomox

import (
	"bytes"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// Signature struct
type Signature struct {
	V byte
	R common.Hash
	S common.Hash
}

// OrderItem : info that will be store in database
type OrderItem struct {
	Quantity        *big.Int       `json:"quantity,omitempty"`
	Price           *big.Int       `json:"price,omitempty"`
	ExchangeAddress common.Address `json:"exchangeAddress,omitempty"`
	UserAddress     common.Address `json:"userAddress,omitempty"`
	BaseToken       common.Address `json:"baseToken,omitempty"`
	QuoteToken      common.Address `json:"quoteToken,omitempty"`
	Status          string         `json:"status,omitempty"`
	Side            string         `json:"side,omitempty"`
	Type            string         `json:"type,omitempty"`
	Hash            common.Hash    `json:"hash,omitempty"`
	Signature       *Signature     `json:"signature,omitempty"`
	FilledAmount    *big.Int       `json:"filledAmount,omitempty"`
	Nonce           *big.Int       `json:"nonce,omitempty"`
	MakeFee         *big.Int       `json:"makeFee,omitempty"`
	TakeFee         *big.Int       `json:"takeFee,omitempty"`
	PairName        string         `json:"pairName,omitempty"`
	CreatedAt       uint64         `json:"createdAt,omitempty"`
	UpdatedAt       uint64         `json:"updatedAt,omitempty"`
	OrderID         uint64         `json:"orderID,omitempty"`
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

	return fmt.Sprintf("orderID : %s, price: %s, quantity :%s, relayerID: %s",
		new(big.Int).SetBytes(order.Key), order.Item.Price, order.Item.Quantity, order.Item.ExchangeAddress.Hex())
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
func NewOrder(orderItem *OrderItem, orderList []byte) *Order {
	key := GetKeyFromBig(new(big.Int).SetUint64(orderItem.OrderID))
	orderItem.NextOrder = EmptyKey()
	orderItem.PrevOrder = EmptyKey()
	orderItem.OrderList = orderList
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
	order.Item.UpdatedAt = newTimestamp
	order.Item.Quantity = CloneBigInt(newQuantity)
	log.Debug("QUANTITY", order.Item.Quantity.String())
	orderList.SaveOrder(order)
	orderList.Save()
}
