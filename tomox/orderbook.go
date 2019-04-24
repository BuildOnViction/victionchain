package tomox

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

const (
	Ask    = "SELL"
	Bid    = "BUY"
	Market = "market"
	Limit  = "limit"
	Cancel = "CANCEL"

	// we use a big number as segment for storing order, order list from order tree slot.
	// as sequential id
	SlotSegment = common.AddressLength
)

type OrderBookItem struct {
	Timestamp     uint64 `json:"time"`
	NextOrderID   uint64 `json:"nextOrderID"`
	MaxPricePoint uint64 `json:"maxVolume"` // maximum
	Name          string `json:"name"`
}

type OrderBookItemBSON struct {
	Timestamp     string `json:"time" bson:"time"`
	NextOrderID   string `json:"nextOrderID" bson:"nextOrderID"`
	MaxPricePoint string `json:"maxVolume" bson:"maxVolume"`
	Name          string `json:"name" bson:"name"`
}

type OrderBookItemRecord struct {
	Key   string
	Value *OrderBookItem
}

type OrderBookItemRecordBSON struct {
	Key   string
	Value *OrderBookItemBSON
}

// OrderBook : list of orders
type OrderBook struct {
	db   OrderDao   // this is for orderBook
	Bids *OrderTree `json:"bids"`
	Asks *OrderTree `json:"asks"`
	Item *OrderBookItem

	Key  []byte
	slot *big.Int
}

// NewOrderBook : return new order book
func NewOrderBook(name string, db OrderDao) *OrderBook {

	item := &OrderBookItem{
		NextOrderID: 0,
		Name:        strings.ToLower(name),
	}

	// do slot with hash to prevent collision

	// we convert to lower case, so even with name as contract address, it is still correct
	// without converting back from hex to bytes
	key := crypto.Keccak256([]byte(item.Name))
	slot := new(big.Int).SetBytes(key)

	// we just increase the segment at the most byte at address length level to avoid conflict
	// somehow it is like 2 hashes has the same common prefix and it is very difficult to resolve
	// the order id start at orderbook slot
	// the price of order tree start at order tree slot
	bidsKey := GetSegmentHash(key, 1, SlotSegment)
	asksKey := GetSegmentHash(key, 2, SlotSegment)

	orderBook := &OrderBook{
		db:   db,
		Item: item,
		slot: slot,
		Key:  key,
	}

	bids := NewOrderTree(db, bidsKey, orderBook)
	asks := NewOrderTree(db, asksKey, orderBook)

	// set asks and bids
	orderBook.Bids = bids
	orderBook.Asks = asks

	// no need to update when there is no operation yet
	orderBook.UpdateTime()

	return orderBook
}

func (orderBook *OrderBook) Save() error {

	orderBook.Asks.Save()
	orderBook.Bids.Save()

	return orderBook.db.Put(orderBook.Key, orderBook.Item)
}

// commit everything by trigger db.Commit, later we can map custom encode and decode based on item
func (orderBook *OrderBook) Commit() error {
	return orderBook.db.Commit()
}

func (orderBook *OrderBook) Restore() error {
	orderBook.Asks.Restore()
	orderBook.Bids.Restore()

	val, err := orderBook.db.Get(orderBook.Key, orderBook.Item)
	if err == nil {
		orderBook.Item = val.(*OrderBookItem)
	}

	return err
}

func (orderBook *OrderBook) GetOrderIDFromBook(key []byte) uint64 {
	orderSlot := new(big.Int).SetBytes(key)
	return Sub(orderSlot, orderBook.slot).Uint64()
}

func (orderBook *OrderBook) GetOrderIDFromKey(key []byte) []byte {
	orderSlot := new(big.Int).SetBytes(key)
	return common.BigToHash(Add(orderBook.slot, orderSlot)).Bytes()
}

func (orderBook *OrderBook) GetOrder(key []byte) *Order {
	if orderBook.db.IsEmptyKey(key) {
		return nil
	}
	storedKey := orderBook.GetOrderIDFromKey(key)
	orderItem := &OrderItem{}
	val, err := orderBook.db.Get(storedKey, orderItem)
	if err != nil {
		log.Error("Key not found", "key", storedKey, "err", err)
		return nil
	}

	order := &Order{
		Item: val.(*OrderItem),
		Key:  key,
	}
	return order
}

func (orderBook *OrderBook) String(startDepth int) string {
	tabs := strings.Repeat("\t", startDepth)
	return fmt.Sprintf("%s{\n\t%sName: %s\n\t%sTimestamp: %d\n\t%sNextOrderID: %d\n\t%sBids: %s\n\t%sAsks: %s\n%s}\n",
		tabs,
		tabs, orderBook.Item.Name, tabs, orderBook.Item.Timestamp, tabs, orderBook.Item.NextOrderID,
		tabs, orderBook.Bids.String(startDepth+1), tabs, orderBook.Asks.String(startDepth+1),
		tabs)
}

// UpdateTime : update time for order book
func (orderBook *OrderBook) UpdateTime() {
	timestamp := uint64(time.Now().Unix())
	orderBook.Item.Timestamp = timestamp
}

// BestBid : get the best bid of the order book
func (orderBook *OrderBook) BestBid() (value *big.Int) {
	return orderBook.Bids.MaxPrice()
}

// BestAsk : get the best ask of the order book
func (orderBook *OrderBook) BestAsk() (value *big.Int) {
	return orderBook.Asks.MinPrice()
}

// WorstBid : get the worst bid of the order book
func (orderBook *OrderBook) WorstBid() (value *big.Int) {
	return orderBook.Bids.MinPrice()
}

// WorstAsk : get the worst ask of the order book
func (orderBook *OrderBook) WorstAsk() (value *big.Int) {
	return orderBook.Asks.MaxPrice()
}

// processMarketOrder : process the market order
func (orderBook *OrderBook) processMarketOrder(order *OrderItem, verbose bool) []map[string]string {
	var trades []map[string]string
	quantityToTrade := order.Quantity
	side := order.Side
	var newTrades []map[string]string
	// speedup the comparison, do not assign because it is pointer
	zero := Zero()
	if side == Bid {
		for quantityToTrade.Cmp(zero) > 0 && orderBook.Asks.NotEmpty() {
			bestPriceAsks := orderBook.Asks.MinPriceList()
			quantityToTrade, newTrades = orderBook.processOrderList(Ask, bestPriceAsks, quantityToTrade, order, verbose)
			trades = append(trades, newTrades...)
		}
	} else {
		for quantityToTrade.Cmp(zero) > 0 && orderBook.Bids.NotEmpty() {
			bestPriceBids := orderBook.Bids.MaxPriceList()
			quantityToTrade, newTrades = orderBook.processOrderList(Bid, bestPriceBids, quantityToTrade, order, verbose)
			trades = append(trades, newTrades...)
		}
	}
	return trades
}

// processLimitOrder : process the limit order, can change the quote
// If not care for performance, we should make a copy of quote to prevent further reference problem
func (orderBook *OrderBook) processLimitOrder(order *OrderItem, verbose bool) ([]map[string]string, *OrderItem) {
	var trades []map[string]string
	quantityToTrade := order.Quantity
	side := order.Side
	price := order.Price

	var newTrades []map[string]string
	var orderInBook *OrderItem
	// speedup the comparison, do not assign because it is pointer
	zero := Zero()

	if side == Bid {
		minPrice := orderBook.Asks.MinPrice()
		for quantityToTrade.Cmp(zero) > 0 && orderBook.Asks.NotEmpty() && price.Cmp(minPrice) >= 0 {
			bestPriceAsks := orderBook.Asks.MinPriceList()
			quantityToTrade, newTrades = orderBook.processOrderList(Ask, bestPriceAsks, quantityToTrade, order, verbose)
			trades = append(trades, newTrades...)
			minPrice = orderBook.Asks.MinPrice()
		}

		if quantityToTrade.Cmp(zero) > 0 {
			order.OrderID = orderBook.Item.NextOrderID
			order.Quantity = quantityToTrade
			orderBook.Bids.InsertOrder(order)
			orderInBook = order
		}

	} else {
		maxPrice := orderBook.Bids.MaxPrice()
		for quantityToTrade.Cmp(zero) > 0 && orderBook.Bids.NotEmpty() && price.Cmp(maxPrice) <= 0 {
			bestPriceBids := orderBook.Bids.MaxPriceList()
			quantityToTrade, newTrades = orderBook.processOrderList(Bid, bestPriceBids, quantityToTrade, order, verbose)
			trades = append(trades, newTrades...)
			maxPrice = orderBook.Bids.MaxPrice()
		}

		if quantityToTrade.Cmp(zero) > 0 {
			order.OrderID = orderBook.Item.NextOrderID
			order.Quantity = quantityToTrade
			orderBook.Asks.InsertOrder(order)
			orderInBook = order
		}
	}
	return trades, orderInBook
}

// ProcessOrder : process the order
func (orderBook *OrderBook) ProcessOrder(order *OrderItem, verbose bool) ([]map[string]string, *OrderItem) {
	orderType := order.Type
	var orderInBook *OrderItem
	var trades []map[string]string

	//orderBook.UpdateTime()
	//// if we do not use auto-increment orderid, we must set price slot to avoid conflict
	//orderBook.Item.NextOrderID++

	if orderType == Market {
		trades = orderBook.processMarketOrder(order, verbose)
	} else {
		trades, orderInBook = orderBook.processLimitOrder(order, verbose)
	}

	// update orderBook
	orderBook.Save()

	return trades, orderInBook
}

// processOrderList : process the order list
func (orderBook *OrderBook) processOrderList(side string, orderList *OrderList, quantityStillToTrade *big.Int, order *OrderItem, verbose bool) (*big.Int, []map[string]string) {
	quantityToTrade := CloneBigInt(quantityStillToTrade)
	var trades []map[string]string
	// speedup the comparison, do not assign because it is pointer
	zero := Zero()
	for orderList.Item.Length > 0 && quantityToTrade.Cmp(zero) > 0 {

		headOrder := orderList.GetOrder(orderList.Item.HeadOrder)
		if headOrder == nil {
			panic("headOrder is null")
		}

		tradedPrice := CloneBigInt(headOrder.Item.Price)

		var newBookQuantity *big.Int
		var tradedQuantity *big.Int

		if IsStrictlySmallerThan(quantityToTrade, headOrder.Item.Quantity) {
			tradedQuantity = CloneBigInt(quantityToTrade)
			// Do the transaction
			newBookQuantity = Sub(headOrder.Item.Quantity, quantityToTrade)
			headOrder.UpdateQuantity(orderList, newBookQuantity, headOrder.Item.UpdatedAt)
			quantityToTrade = Zero()

		} else if IsEqual(quantityToTrade, headOrder.Item.Quantity) {
			tradedQuantity = CloneBigInt(quantityToTrade)
			if side == Bid {
				orderBook.Bids.RemoveOrder(headOrder)
			} else {
				orderBook.Asks.RemoveOrder(headOrder)
			}
			quantityToTrade = Zero()

		} else {
			tradedQuantity = CloneBigInt(headOrder.Item.Quantity)
			if side == Bid {
				orderBook.Bids.RemoveOrder(headOrder)
			} else {
				orderBook.Asks.RemoveOrderFromOrderList(headOrder, orderList)
			}
		}

		if verbose {
			log.Info("TRADE", "Timestamp", orderBook.Item.Timestamp, "Price", tradedPrice, "Quantity", tradedQuantity, "TradeID", headOrder.Item.ExchangeAddress.Hex(), "Matching TradeID", order.ExchangeAddress.Hex())
		}

		transactionRecord := make(map[string]string)
		transactionRecord["timestamp"] = strconv.FormatUint(orderBook.Item.Timestamp, 10)
		transactionRecord["price"] = tradedPrice.String()
		transactionRecord["quantity"] = tradedQuantity.String()

		trades = append(trades, transactionRecord)
	}
	return quantityToTrade, trades
}

// CancelOrder : cancel the order, just need ID, side and price, of course order must belong
// to a price point as well
func (orderBook *OrderBook) CancelOrder(side string, orderID uint64, price *big.Int) error {
	orderBook.UpdateTime()
	key := GetKeyFromBig(big.NewInt(int64(orderID)))
	var err error
	if side == Bid {
		order := orderBook.Bids.GetOrder(key, price)
		if order != nil {
			_, err = orderBook.Bids.RemoveOrder(order)
		}
	} else {

		order := orderBook.Asks.GetOrder(key, price)
		if order != nil {
			_, err = orderBook.Asks.RemoveOrder(order)
		}
	}

	return err
}

func (orderBook *OrderBook) UpdateOrder(order *OrderItem) error {
	return orderBook.ModifyOrder(order, order.OrderID, order.Price)
}

// ModifyOrder : modify the order
func (orderBook *OrderBook) ModifyOrder(order *OrderItem, orderID uint64, price *big.Int) error {
	orderBook.UpdateTime()

	key := GetKeyFromBig(new(big.Int).SetUint64(order.OrderID))
	if order.Side == Bid {

		if orderBook.Bids.OrderExist(key, price) {
			return orderBook.Bids.UpdateOrder(order)
		}
	} else {

		if orderBook.Asks.OrderExist(key, price) {
			return orderBook.Asks.UpdateOrder(order)
		}
	}

	return nil
}

// VolumeAtPrice : get volume at the current price
func (orderBook *OrderBook) VolumeAtPrice(side string, price *big.Int) *big.Int {
	volume := Zero()
	if side == Bid {
		if orderBook.Bids.PriceExist(price) {
			orderList := orderBook.Bids.PriceList(price)
			// incase we use cache for PriceList
			volume = CloneBigInt(orderList.Item.Volume)
		}
	} else {
		// other case
		if orderBook.Asks.PriceExist(price) {
			orderList := orderBook.Asks.PriceList(price)
			volume = CloneBigInt(orderList.Item.Volume)
		}
	}

	return volume
}

// Save order pending into orderbook tree.
func (orderBook *OrderBook) SaveOrderPending(order *OrderItem) error {
	quantityToTrade := order.Quantity
	side := order.Side
	zero := Zero()

	orderBook.UpdateTime()
	// if we do not use auto-increment orderid, we must set price slot to avoid conflict
	orderBook.Item.NextOrderID++

	if side == Bid {
		if quantityToTrade.Cmp(zero) > 0 {
			order.OrderID = orderBook.Item.NextOrderID
			order.Quantity = quantityToTrade
			return orderBook.Bids.InsertOrder(order)
		}
	} else {
		if quantityToTrade.Cmp(zero) > 0 {
			order.OrderID = orderBook.Item.NextOrderID
			order.Quantity = quantityToTrade
			return orderBook.Asks.InsertOrder(order)
		}
	}

	return nil
}
