package tomox

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/tomox/tomox_state"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

const (
	OrderStatusNew           = "NEW"
	OrderStatusOpen          = "OPEN"
	OrderStatusPartialFilled = "PARTIAL_FILLED"
	OrderStatusFilled        = "FILLED"
	OrderStatusCancelled     = "CANCELLED"
	OrderStatusRejected 	= "REJECTED"
)


type SignatureRecord struct {
	V byte   `json:"V" bson:"V"`
	R string `json:"R" bson:"R"`
	S string `json:"S" bson:"S"`
}

type OrderItemBSON struct {
	Quantity        string           `json:"quantity,omitempty" bson:"quantity"`
	Price           string           `json:"price,omitempty" bson:"price"`
	ExchangeAddress string           `json:"exchangeAddress,omitempty" bson:"exchangeAddress"`
	UserAddress     string           `json:"userAddress,omitempty" bson:"userAddress"`
	BaseToken       string           `json:"baseToken,omitempty" bson:"baseToken"`
	QuoteToken      string           `json:"quoteToken,omitempty" bson:"quoteToken"`
	Status          string           `json:"status,omitempty" bson:"status"`
	Side            string           `json:"side,omitempty" bson:"side"`
	Type            string           `json:"type,omitempty" bson:"type"`
	Hash            string           `json:"hash,omitempty" bson:"hash"`
	Signature       *SignatureRecord `json:"signature,omitempty" bson:"signature"`
	FilledAmount    string           `json:"filledAmount,omitempty" bson:"filledAmount"`
	Nonce           string           `json:"nonce,omitempty" bson:"nonce"`
	PairName        string           `json:"pairName,omitempty" bson:"pairName"`
	CreatedAt       time.Time        `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt,omitempty" bson:"updatedAt"`
	OrderID         string           `json:"orderID,omitempty" bson:"orderID"`
	NextOrder       string           `json:"nextOrder,omitempty" bson:"nextOrder"`
	PrevOrder       string           `json:"prevOrder,omitempty" bson:"prevOrder"`
	OrderList       string           `json:"orderList,omitempty" bson:"orderList"`
	Key             string           `json:"key" bson:"key"`
}

type Order struct {
	Item *tomox_state.OrderItem
	Key  []byte `json:"orderID"`
}

func (order *Order) String() string {

	return fmt.Sprintf("orderID : %s, price: %s, quantity :%s, relayerID: %s",
		new(big.Int).SetBytes(order.Key), order.Item.Price, order.Item.Quantity, order.Item.ExchangeAddress.Hex())
}

func (order *Order) GetNextOrder(orderList *OrderList, dryrun bool, blockHash common.Hash) *Order {
	nextOrder := orderList.GetOrder(order.Item.NextOrder, dryrun, blockHash)

	return nextOrder
}

func (order *Order) GetPrevOrder(orderList *OrderList, dryrun bool, blockHash common.Hash) *Order {
	prevOrder := orderList.GetOrder(order.Item.PrevOrder, dryrun, blockHash)

	return prevOrder
}

// NewOrder : create new order with quote ( can be ethereum address )
func NewOrder(orderItem *tomox_state.OrderItem, orderListKey []byte) *Order {
	key := GetKeyFromBig(new(big.Int).SetUint64(orderItem.OrderID))
	orderItem.NextOrder = EmptyKey()
	orderItem.PrevOrder = EmptyKey()
	orderItem.OrderList = orderListKey
	// key should be Hash for compatible with smart contract
	order := &Order{
		Key:  key,
		Item: orderItem,
	}

	return order
}

// UpdateQuantity : update quantity of the order
func (order *Order) UpdateQuantity(orderList *OrderList, newQuantity *big.Int, dryrun bool, blockHash common.Hash) error {
	if newQuantity.Cmp(order.Item.Quantity) > 0 && !bytes.Equal(orderList.Item.TailOrder, order.Key) {
		if err := orderList.MoveToTail(order, dryrun, blockHash); err != nil {
			return err
		}
	}
	// update volume and modified timestamp
	orderList.Item.Volume = Sub(orderList.Item.Volume, Sub(order.Item.Quantity, newQuantity))
	order.Item.Quantity = CloneBigInt(newQuantity)
	log.Debug("QUANTITY", order.Item.Quantity.String())
	if err := orderList.SaveOrder(order, dryrun, blockHash); err != nil {
		return err
	}
	if err := orderList.Save(dryrun, blockHash); err != nil {
		return err
	}
	return nil
}
