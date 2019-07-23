package tomox

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
)

const (
	Ask    = "SELL"
	Bid    = "BUY"
	Market = "MO"
	Limit  = "LO"
	Cancel = "CANCELLED"

	// we use a big number as segment for storing order, order list from order tree slot.
	// as sequential id
	SlotSegment = common.AddressLength
)

var ErrDoesNotExist = errors.New("order doesn't exist in ordertree")

type OrderBookItem struct {
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
	Timestamp     uint64 `json:"time"`

	Key  []byte
	Slot *big.Int
}

// NewOrderBook : return new order book
func NewOrderBook(name string, db OrderDao) *OrderBook {

	item := &OrderBookItem{
		NextOrderID: 0,
		Name:        name, //name is already in lower format
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
		Slot: slot,
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

func (orderBook *OrderBook) Save(dryrun bool) error {

	log.Debug("save orderbook asks")
	if err := orderBook.Asks.Save(dryrun); err != nil {
		log.Error("can't save orderbook asks", "err", err)
		return err
	}

	log.Debug("save orderbook bids")
	if err := orderBook.Bids.Save(dryrun); err != nil {
		log.Error("can't save orderbook bids", "err", err)
		return err
	}

	log.Debug("save orderbook", "key", hex.EncodeToString(orderBook.Key))
	return orderBook.db.Put(orderBook.Key, orderBook.Item, dryrun)
}

func (orderBook *OrderBook) Restore(dryrun bool) error {
	log.Debug("restore orderbook asks")
	if err := orderBook.Asks.Restore(dryrun); err != nil {
		log.Error("can't restore orderbook asks", "err", err)
		return err
	}

	log.Debug("restore orderbook bids")
	if err := orderBook.Bids.Restore(dryrun); err != nil {
		log.Error("can't restore orderbook bids", "err", err)
		return err
	}

	val, err := orderBook.db.Get(orderBook.Key, orderBook.Item, dryrun)
	if err == nil {
		orderBook.Item = val.(*OrderBookItem)
		log.Debug("orderbook restored", "orderBook.Item", orderBook.Item)
	}

	return err
}

func (orderBook *OrderBook) GetOrderIDFromBook(key []byte) uint64 {
	orderSlot := new(big.Int).SetBytes(key)
	return Sub(orderSlot, orderBook.Slot).Uint64()
}

func (orderBook *OrderBook) GetOrderIDFromKey(key []byte) []byte {
	orderSlot := new(big.Int).SetBytes(key)
	return common.BigToHash(Add(orderBook.Slot, orderSlot)).Bytes()
}

func (orderBook *OrderBook) GetOrder(storedKey, key []byte, dryrun bool) *Order {
	if orderBook.db.IsEmptyKey(key) || orderBook.db.IsEmptyKey(storedKey) {
		return nil
	}
	orderItem := &OrderItem{}
	val, err := orderBook.db.Get(storedKey, orderItem, dryrun)
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

// UpdateTime : update time for order book
func (orderBook *OrderBook) UpdateTime() {
	timestamp := uint64(time.Now().Unix())
	orderBook.Timestamp = timestamp
}

// BestBid : get the best bid of the order book
func (orderBook *OrderBook) BestBid(dryrun bool) (value *big.Int) {
	return orderBook.Bids.MaxPrice(dryrun)
}

// BestAsk : get the best ask of the order book
func (orderBook *OrderBook) BestAsk(dryrun bool) (value *big.Int) {
	return orderBook.Asks.MinPrice(dryrun)
}

// WorstBid : get the worst bid of the order book
func (orderBook *OrderBook) WorstBid(dryrun bool) (value *big.Int) {
	return orderBook.Bids.MinPrice(dryrun)
}

// WorstAsk : get the worst ask of the order book
func (orderBook *OrderBook) WorstAsk(dryrun bool) (value *big.Int) {
	return orderBook.Asks.MaxPrice(dryrun)
}

// processMarketOrder : process the market order
func (orderBook *OrderBook) processMarketOrder(order *OrderItem, verbose bool, dryrun bool) ([]map[string]string, error) {
	var (
		trades    []map[string]string
		newTrades []map[string]string
		err       error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	// speedup the comparison, do not assign because it is pointer
	zero := Zero()
	if side == Bid {
		for quantityToTrade.Cmp(zero) > 0 && orderBook.Asks.NotEmpty() {
			bestPriceAsks := orderBook.Asks.MinPriceList(dryrun)
			quantityToTrade, newTrades, err = orderBook.processOrderList(Ask, bestPriceAsks, quantityToTrade, order, verbose, dryrun)
			if err != nil {
				return nil, err
			}
			trades = append(trades, newTrades...)
		}
	} else {
		for quantityToTrade.Cmp(zero) > 0 && orderBook.Bids.NotEmpty() {
			bestPriceBids := orderBook.Bids.MaxPriceList(dryrun)
			quantityToTrade, newTrades, err = orderBook.processOrderList(Bid, bestPriceBids, quantityToTrade, order, verbose, dryrun)
			if err != nil {
				return nil, err
			}
			trades = append(trades, newTrades...)
		}
	}
	return trades, nil
}

// processLimitOrder : process the limit order, can change the quote
// If not care for performance, we should make a copy of quote to prevent further reference problem
func (orderBook *OrderBook) processLimitOrder(order *OrderItem, verbose bool, dryrun bool) ([]map[string]string, *OrderItem, error) {
	var (
		trades      []map[string]string
		newTrades   []map[string]string
		orderInBook *OrderItem
		err         error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	price := order.Price

	// speedup the comparison, do not assign because it is pointer
	zero := Zero()

	if side == Bid {
		minPrice := orderBook.Asks.MinPrice(dryrun)
		for quantityToTrade.Cmp(zero) > 0 && orderBook.Asks.NotEmpty() && price.Cmp(minPrice) >= 0 {
			bestPriceAsks := orderBook.Asks.MinPriceList(dryrun)
			quantityToTrade, newTrades, err = orderBook.processOrderList(Ask, bestPriceAsks, quantityToTrade, order, verbose, dryrun)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			minPrice = orderBook.Asks.MinPrice(dryrun)
		}

		if quantityToTrade.Cmp(zero) > 0 {
			order.OrderID = orderBook.Item.NextOrderID
			order.Quantity = quantityToTrade
			orderBook.Bids.InsertOrder(order, dryrun)
			orderInBook = order
		}

	} else {
		maxPrice := orderBook.Bids.MaxPrice(dryrun)
		for quantityToTrade.Cmp(zero) > 0 && orderBook.Bids.NotEmpty() && price.Cmp(maxPrice) <= 0 {
			bestPriceBids := orderBook.Bids.MaxPriceList(dryrun)
			quantityToTrade, newTrades, err = orderBook.processOrderList(Bid, bestPriceBids, quantityToTrade, order, verbose, dryrun)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			maxPrice = orderBook.Bids.MaxPrice(dryrun)
		}

		if quantityToTrade.Cmp(zero) > 0 {
			order.OrderID = orderBook.Item.NextOrderID
			order.Quantity = quantityToTrade
			orderBook.Asks.InsertOrder(order, dryrun)
			orderInBook = order
		}
	}
	return trades, orderInBook, nil
}

// ProcessOrder : process the order
func (orderBook *OrderBook) ProcessOrder(order *OrderItem, verbose bool, dryrun bool) ([]map[string]string, *OrderItem, error) {
	var (
		orderInBook *OrderItem
		trades      []map[string]string
		err         error
	)
	orderType := order.Type
	orderBook.UpdateTime()
	// if we do not use auto-increment orderid, we must set price slot to avoid conflict
	orderBook.Item.NextOrderID++

	if orderType == Market {
		trades, err = orderBook.processMarketOrder(order, verbose, dryrun)
		if err != nil {
			return nil, nil, err
		}
	} else {
		trades, orderInBook, err = orderBook.processLimitOrder(order, verbose, dryrun)
		if err != nil {
			return nil, nil, err
		}
	}

	// update orderBook
	if err := orderBook.Save(dryrun); err != nil {
		return nil, nil, err
	}

	return trades, orderInBook, nil
}

// processOrderList : process the order list
func (orderBook *OrderBook) processOrderList(side string, orderList *OrderList, quantityStillToTrade *big.Int, order *OrderItem, verbose bool, dryrun bool) (*big.Int, []map[string]string, error) {
	quantityToTrade := CloneBigInt(quantityStillToTrade)
	var trades []map[string]string
	// speedup the comparison, do not assign because it is pointer
	zero := Zero()
	for orderList.Item.Length > uint64(0) && quantityToTrade.Cmp(zero) > 0 {

		headOrder := orderList.GetOrder(orderList.Item.HeadOrder, dryrun)
		if headOrder == nil {
			return nil, nil, fmt.Errorf("headOrder is null")
		}

		tradedPrice := CloneBigInt(headOrder.Item.Price)

		var newBookQuantity *big.Int
		var tradedQuantity *big.Int

		if IsStrictlySmallerThan(quantityToTrade, headOrder.Item.Quantity) {
			tradedQuantity = CloneBigInt(quantityToTrade)
			// Do the transaction
			newBookQuantity = Sub(headOrder.Item.Quantity, quantityToTrade)
			if err := headOrder.UpdateQuantity(orderList, newBookQuantity, headOrder.Item.UpdatedAt, dryrun); err != nil {
				return nil, nil, err
			}
			quantityToTrade = Zero()
		} else if IsEqual(quantityToTrade, headOrder.Item.Quantity) {
			tradedQuantity = CloneBigInt(quantityToTrade)
			if side == Bid {
				if err := orderBook.Bids.RemoveOrderFromOrderList(headOrder, orderList, dryrun); err != nil {
					return nil, nil, err
				}
			} else {
				if err := orderBook.Asks.RemoveOrderFromOrderList(headOrder, orderList, dryrun); err != nil {
					return nil, nil, err
				}
			}
			quantityToTrade = Zero()

		} else {
			tradedQuantity = CloneBigInt(headOrder.Item.Quantity)
			if side == Bid {
				if err := orderBook.Bids.RemoveOrderFromOrderList(headOrder, orderList, dryrun); err != nil {
					return nil, nil, err
				}
			} else {
				if err := orderBook.Asks.RemoveOrderFromOrderList(headOrder, orderList, dryrun); err != nil {
					return nil, nil, err
				}
			}
		}

		if verbose {
			log.Info("TRADE", "Timestamp", orderBook.Timestamp, "Price", tradedPrice, "Quantity", tradedQuantity, "TradeID", headOrder.Item.ExchangeAddress.Hex(), "Matching TradeID", order.ExchangeAddress.Hex())
		}

		transactionRecord := make(map[string]string)
		transactionRecord["timestamp"] = strconv.FormatUint(orderBook.Timestamp, 10)
		transactionRecord["quantity"] = tradedQuantity.String()
		transactionRecord["exAddr"] = headOrder.Item.ExchangeAddress.String()
		transactionRecord["uAddr"] = headOrder.Item.UserAddress.String()
		transactionRecord["bToken"] = headOrder.Item.BaseToken.String()
		transactionRecord["qToken"] = headOrder.Item.QuoteToken.String()

		trades = append(trades, transactionRecord)
	}
	return quantityToTrade, trades, nil
}

// CancelOrder : cancel the order, just need ID, side and price, of course order must belong
// to a price point as well
func (orderBook *OrderBook) CancelOrder(order *OrderItem, dryrun bool) error {
	key := GetKeyFromBig(big.NewInt(int64(order.OrderID)))
	if order.Side == Bid {
		orderInDB := orderBook.Bids.GetOrder(key, order.Price, dryrun)
		if orderInDB == nil || orderInDB.Item.Hash != order.Hash {
			return ErrDoesNotExist
		}
		orderInDB.Item.Status = Cancel
		if err := orderBook.Bids.RemoveOrder(orderInDB, dryrun); err != nil {
			return err
		}
	} else {
		orderInDB := orderBook.Asks.GetOrder(key, order.Price, dryrun)
		if orderInDB == nil || orderInDB.Item.Hash != order.Hash {
			return ErrDoesNotExist
		}
		orderInDB.Item.Status = Cancel
		if err := orderBook.Asks.RemoveOrder(orderInDB, dryrun); err != nil {
			return err
		}
	}

	// snapshot orderbook
	orderBook.UpdateTime()
	if err := orderBook.Save(dryrun); err != nil {
		return err
	}

	return nil
}

func (orderBook *OrderBook) UpdateOrder(order *OrderItem, dryrun bool) error {
	return orderBook.ModifyOrder(order, order.OrderID, order.Price, dryrun)
}

// ModifyOrder : modify the order
func (orderBook *OrderBook) ModifyOrder(order *OrderItem, orderID uint64, price *big.Int, dryrun bool) error {
	orderBook.UpdateTime()

	key := GetKeyFromBig(new(big.Int).SetUint64(order.OrderID))
	if order.Side == Bid {

		if orderBook.Bids.OrderExist(key, price, dryrun) {
			return orderBook.Bids.UpdateOrder(order, dryrun)
		}
	} else {

		if orderBook.Asks.OrderExist(key, price, dryrun) {
			return orderBook.Asks.UpdateOrder(order, dryrun)
		}
	}

	return nil
}

// VolumeAtPrice : get volume at the current price
func (orderBook *OrderBook) VolumeAtPrice(side string, price *big.Int, dryrun bool) *big.Int {
	volume := Zero()
	if side == Bid {
		if orderBook.Bids.PriceExist(price, dryrun) {
			orderList := orderBook.Bids.PriceList(price, dryrun)
			// incase we use cache for PriceList
			volume = CloneBigInt(orderList.Item.Volume)
		}
	} else {
		// other case
		if orderBook.Asks.PriceExist(price, dryrun) {
			orderList := orderBook.Asks.PriceList(price, dryrun)
			volume = CloneBigInt(orderList.Item.Volume)
		}
	}

	return volume
}

func (orderBook *OrderBook) Hash() (common.Hash, error) {
	obEncoded, err := EncodeBytesItem(orderBook.Item)
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(obEncoded), nil
}
