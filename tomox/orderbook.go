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
	SlotSegment         = common.AddressLength
	orderbookItemPrefix = "OB"
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
	db        OrderDao   // this is for orderBook
	Bids      *OrderTree `json:"bids"`
	Asks      *OrderTree `json:"asks"`
	Item      *OrderBookItem
	Timestamp uint64 `json:"time"`

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

func (orderBook *OrderBook) Save(dryrun bool, blockHash common.Hash) error {

	log.Debug("save orderbook asks")
	if err := orderBook.Asks.Save(dryrun, blockHash); err != nil {
		log.Error("can't save orderbook asks", "err", err)
		return err
	}

	log.Debug("save orderbook bids")
	if err := orderBook.Bids.Save(dryrun, blockHash); err != nil {
		log.Error("can't save orderbook bids", "err", err)
		return err
	}

	orderBookItemKey := append([]byte(orderbookItemPrefix), orderBook.Key...)
	log.Debug("save orderbook", "key", hex.EncodeToString(orderBookItemKey))
	return orderBook.db.Put(orderBookItemKey, orderBook.Item, dryrun, blockHash)
}

func (orderBook *OrderBook) Restore(dryrun bool, blockHash common.Hash) error {
	if err := orderBook.Asks.Restore(dryrun, blockHash); err != nil {
		log.Error("can't restore orderbook asks", "err", err)
		return err
	}
	// rootKey might change once restoring thus to be safe, we should save changes to DB
	if err := orderBook.Asks.Save(dryrun, blockHash); err != nil {
		log.Error("can't save asks changes to DB", "err", err)
		return err
	}
	log.Debug("restored orderbook asks", "asks.Item", orderBook.Asks.Item)

	if err := orderBook.Bids.Restore(dryrun, blockHash); err != nil {
		log.Error("can't restore orderbook bids", "err", err)
		return err
	}
	// rootKey might change once restoring thus to be safe, we should save changes to DB
	if err := orderBook.Bids.Save(dryrun, blockHash); err != nil {
		log.Error("can't save bids changes to DB", "err", err)
		return err
	}
	log.Debug("restored orderbook bids", "bids.Item", orderBook.Bids.Item)

	orderBookItemKey := append([]byte(orderbookItemPrefix), orderBook.Key...)
	val, err := orderBook.db.Get(orderBookItemKey, orderBook.Item, dryrun, blockHash)
	if err == nil {
		orderBook.Item = val.(*OrderBookItem)
		log.Debug("restored orderbook", "orderBook.Item", orderBook.Item)
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

func (orderBook *OrderBook) GetOrder(storedKey, key []byte, dryrun bool, blockHash common.Hash) *Order {
	if orderBook.db.IsEmptyKey(key) || orderBook.db.IsEmptyKey(storedKey) {
		return nil
	}
	orderItem := &OrderItem{}
	val, err := orderBook.db.Get(storedKey, orderItem, dryrun, blockHash)
	if err != nil || val == nil {
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
func (orderBook *OrderBook) BestBid(dryrun bool, blockHash common.Hash) (value *big.Int) {
	return orderBook.Bids.MaxPrice(dryrun, blockHash)
}

// BestAsk : get the best ask of the order book
func (orderBook *OrderBook) BestAsk(dryrun bool, blockHash common.Hash) (value *big.Int) {
	return orderBook.Asks.MinPrice(dryrun, blockHash)
}

// WorstBid : get the worst bid of the order book
func (orderBook *OrderBook) WorstBid(dryrun bool, blockHash common.Hash) (value *big.Int) {
	return orderBook.Bids.MinPrice(dryrun, blockHash)
}

// WorstAsk : get the worst ask of the order book
func (orderBook *OrderBook) WorstAsk(dryrun bool, blockHash common.Hash) (value *big.Int) {
	return orderBook.Asks.MaxPrice(dryrun, blockHash)
}

// processMarketOrder : process the market order
func (orderBook *OrderBook) processMarketOrder(order *OrderItem, verbose bool, dryrun bool, blockHash common.Hash) ([]map[string]string, *OrderItem, error) {
	var (
		trades      []map[string]string
		newTrades   []map[string]string
		orderInBook *OrderItem
		err         error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	// speedup the comparison, do not assign because it is pointer
	zero := Zero()
	if side == Bid {
		for quantityToTrade.Cmp(zero) > 0 && orderBook.Asks.NotEmpty() {
			bestPriceAsks := orderBook.Asks.MinPriceList(dryrun, blockHash)
			quantityToTrade, newTrades, orderInBook, err = orderBook.processOrderList(Ask, bestPriceAsks, quantityToTrade, order, verbose, dryrun, blockHash)
			if err != nil {
				return nil, orderInBook, err
			}
			trades = append(trades, newTrades...)
		}
	} else {
		for quantityToTrade.Cmp(zero) > 0 && orderBook.Bids.NotEmpty() {
			bestPriceBids := orderBook.Bids.MaxPriceList(dryrun, blockHash)
			quantityToTrade, newTrades, orderInBook, err = orderBook.processOrderList(Bid, bestPriceBids, quantityToTrade, order, verbose, dryrun, blockHash)
			if err != nil {
				return nil, orderInBook, err
			}
			trades = append(trades, newTrades...)
		}
	}
	return trades, orderInBook, nil
}

// processLimitOrder : process the limit order, can change the quote
// If not care for performance, we should make a copy of quote to prevent further reference problem
func (orderBook *OrderBook) processLimitOrder(order *OrderItem, verbose bool, dryrun bool, blockHash common.Hash) ([]map[string]string, *OrderItem, error) {
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
		minPrice := orderBook.Asks.MinPrice(dryrun, blockHash)
		for quantityToTrade.Cmp(zero) > 0 && orderBook.Asks.NotEmpty() && price.Cmp(minPrice) >= 0 {
			log.Debug("Min price in asks tree", "price", minPrice.String())
			bestPriceAsks := orderBook.Asks.MinPriceList(dryrun, blockHash)
			log.Debug("Orderlist at min price", "orderlist", bestPriceAsks.Item)
			quantityToTrade, newTrades, orderInBook, err = orderBook.processOrderList(Ask, bestPriceAsks, quantityToTrade, order, verbose, dryrun, blockHash)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			log.Debug("New trade found", "newTrades", newTrades, "orderInBook", orderInBook, "quantityToTrade", quantityToTrade)
			minPrice = orderBook.Asks.MinPrice(dryrun, blockHash)
		}

		if quantityToTrade.Cmp(zero) > 0 {
			order.OrderID = orderBook.Item.NextOrderID
			order.Quantity = quantityToTrade
			log.Debug("After matching, order (unmatched part) is now added to bids tree", "order", order)
			if err := orderBook.Bids.InsertOrder(order, dryrun, blockHash); err != nil {
				log.Error("Failed to insert order to bidTree", "pairName", order.PairName, "orderItem", order, "err", err)
				return nil, nil, err
			}
			orderInBook = order
		}

	} else {
		maxPrice := orderBook.Bids.MaxPrice(dryrun, blockHash)
		for quantityToTrade.Cmp(zero) > 0 && orderBook.Bids.NotEmpty() && price.Cmp(maxPrice) <= 0 {
			log.Debug("Max price in bids tree", "price", maxPrice.String())
			bestPriceBids := orderBook.Bids.MaxPriceList(dryrun, blockHash)
			log.Debug("Orderlist at max price", "orderlist", bestPriceBids.Item)
			quantityToTrade, newTrades, orderInBook, err = orderBook.processOrderList(Bid, bestPriceBids, quantityToTrade, order, verbose, dryrun, blockHash)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			log.Debug("New trade found", "newTrades", newTrades, "orderInBook", orderInBook, "quantityToTrade", quantityToTrade)
			maxPrice = orderBook.Bids.MaxPrice(dryrun, blockHash)
		}

		if quantityToTrade.Cmp(zero) > 0 {
			order.OrderID = orderBook.Item.NextOrderID
			order.Quantity = quantityToTrade
			log.Debug("After matching, order (unmatched part) is now back to asks tree", "order", order)
			if err := orderBook.Asks.InsertOrder(order, dryrun, blockHash); err != nil {
				log.Error("Failed to insert order to askTree", "pairName", order.PairName, "orderItem", order, "err", err)
				return nil, nil, err
			}
			orderInBook = order
		}
	}
	return trades, orderInBook, nil
}

// ProcessOrder : process the order
func (orderBook *OrderBook) ProcessOrder(order *OrderItem, verbose bool, dryrun bool, blockHash common.Hash) ([]map[string]string, *OrderItem, error) {
	var (
		orderInBook *OrderItem
		trades      []map[string]string
		err         error
	)
	orderType := order.Type
	orderBook.UpdateTime()
	// if we do not use auto-increment orderid, we must set price slot to avoid conflict
	orderBook.Item.NextOrderID++

	if order.Status == OrderStatusCancelled {
		err := orderBook.CancelOrder(order, dryrun, blockHash)
		switch err {
		case ErrDoesNotExist:
			log.Debug("Order doesn't exist in tree", "order", order)
			return nil, nil, nil
		case nil:
			log.Debug("Cancelled order", "order", order)
			return nil, nil, nil
		default:
			log.Error("Can't cancel order", "order", order, "err", err)
			return nil, nil, err
		}
	}

	if orderType == Market {
		log.Debug("Process market order", "order", order)
		trades, orderInBook, err = orderBook.processMarketOrder(order, verbose, dryrun, blockHash)
		if err != nil {
			return nil, nil, err
		}
	} else {
		log.Debug("Process limit order", "order", order)
		trades, orderInBook, err = orderBook.processLimitOrder(order, verbose, dryrun, blockHash)
		if err != nil {
			return nil, nil, err
		}
	}

	// update orderBook
	if err := orderBook.Save(dryrun, blockHash); err != nil {
		return nil, nil, err
	}

	return trades, orderInBook, nil
}

// processOrderList : process the order list
func (orderBook *OrderBook) processOrderList(side string, orderList *OrderList, quantityStillToTrade *big.Int, order *OrderItem, verbose bool, dryrun bool, blockHash common.Hash) (*big.Int, []map[string]string, *OrderItem, error) {
	log.Debug("Process matching between order and orderlist")
	quantityToTrade := CloneBigInt(quantityStillToTrade)
	var (
		trades      []map[string]string
		orderInBook *OrderItem
	)
	// speedup the comparison, do not assign because it is pointer
	zero := Zero()
	for orderList.Item.Length > uint64(0) && quantityToTrade.Cmp(zero) > 0 {

		headOrder := orderList.GetOrder(orderList.Item.HeadOrder, dryrun, blockHash)
		if headOrder == nil {
			return nil, nil, nil, fmt.Errorf("headOrder is null")
		}
		log.Debug("Get head order in the orderlist", "headOrder", headOrder.Item)

		tradedPrice := CloneBigInt(headOrder.Item.Price)

		var (
			newBookQuantity *big.Int
			tradedQuantity  *big.Int
		)

		if IsStrictlySmallerThan(quantityToTrade, headOrder.Item.Quantity) {
			tradedQuantity = CloneBigInt(quantityToTrade)
			// Do the transaction
			newBookQuantity = Sub(headOrder.Item.Quantity, quantityToTrade)
			if err := headOrder.UpdateQuantity(orderList, newBookQuantity, dryrun, blockHash); err != nil {
				return nil, nil, nil, err
			}
			log.Debug("Update quantity for head order", "headOrder", headOrder.Item)
			quantityToTrade = Zero()
			orderInBook = headOrder.Item
		} else if IsEqual(quantityToTrade, headOrder.Item.Quantity) {
			tradedQuantity = CloneBigInt(quantityToTrade)
			if side == Bid {
				if err := orderBook.Bids.RemoveOrderFromOrderList(headOrder, orderList, dryrun, blockHash); err != nil {
					return nil, nil, nil, err
				}
				log.Debug("Removed headOrder from bids orderlist", "headOrder", headOrder.Item, "orderlist", orderList.Item, "side", side)
			} else {
				if err := orderBook.Asks.RemoveOrderFromOrderList(headOrder, orderList, dryrun, blockHash); err != nil {
					return nil, nil, nil, err
				}
				log.Debug("Removed headOrder from asks orderlist", "headOrder", headOrder.Item, "orderlist", orderList.Item, "side", side)
			}
			quantityToTrade = Zero()

		} else {
			tradedQuantity = CloneBigInt(headOrder.Item.Quantity)
			if side == Bid {
				if err := orderBook.Bids.RemoveOrderFromOrderList(headOrder, orderList, dryrun, blockHash); err != nil {
					return nil, nil, nil, err
				}
				log.Debug("Removed headOrder from bids orderlist", "headOrder", headOrder.Item, "orderlist", orderList.Item, "side", side)
			} else {
				if err := orderBook.Asks.RemoveOrderFromOrderList(headOrder, orderList, dryrun, blockHash); err != nil {
					return nil, nil, nil, err
				}
				log.Debug("Removed headOrder from asks orderlist", "headOrder", headOrder.Item, "orderlist", orderList.Item, "side", side)
			}
			quantityToTrade = Sub(quantityToTrade, tradedQuantity)
		}

		if verbose {
			log.Info("TRADE", "Timestamp", orderBook.Timestamp, "Price", tradedPrice, "Quantity", tradedQuantity, "TradeID", headOrder.Item.ExchangeAddress.Hex(), "Matching TradeID", order.ExchangeAddress.Hex())
		}

		transactionRecord := make(map[string]string)
		transactionRecord[TradeTakerOrderHash] = hex.EncodeToString(order.Hash.Bytes())
		transactionRecord[TradeMakerOrderHash] = hex.EncodeToString(headOrder.Item.Hash.Bytes())
		transactionRecord[TradeTimestamp] = strconv.FormatUint(orderBook.Timestamp, 10)
		transactionRecord[TradeQuantity] = tradedQuantity.String()
		transactionRecord[TradeMakerExchange] = headOrder.Item.ExchangeAddress.String()
		transactionRecord[TradeMaker] = headOrder.Item.UserAddress.String()
		transactionRecord[TradeBaseToken] = headOrder.Item.BaseToken.String()
		transactionRecord[TradeQuoteToken] = headOrder.Item.QuoteToken.String()
		// maker price is actual price
		// taker price is offer price
		// tradedPrice is always actual price
		transactionRecord[TradePrice] = headOrder.Item.Price.String()

		trades = append(trades, transactionRecord)
	}
	return quantityToTrade, trades, orderInBook, nil
}

// CancelOrder : cancel the order, just need ID, side and price, of course order must belong
// to a price point as well
func (orderBook *OrderBook) CancelOrder(order *OrderItem, dryrun bool, blockHash common.Hash) error {
	key := GetKeyFromBig(big.NewInt(int64(order.OrderID)))
	if order.Side == Bid {
		orderInDB := orderBook.Bids.GetOrder(key, order.Price, dryrun, blockHash)
		if orderInDB == nil || orderInDB.Item.Hash != order.Hash {
			return ErrDoesNotExist
		}
		if err := orderBook.Bids.RemoveOrder(orderInDB, dryrun, blockHash); err != nil {
			return err
		}
	} else {
		orderInDB := orderBook.Asks.GetOrder(key, order.Price, dryrun, blockHash)
		if orderInDB == nil || orderInDB.Item.Hash != order.Hash {
			return ErrDoesNotExist
		}
		if err := orderBook.Asks.RemoveOrder(orderInDB, dryrun, blockHash); err != nil {
			return err
		}
	}

	// snapshot orderbook
	orderBook.UpdateTime()
	if err := orderBook.Save(dryrun, blockHash); err != nil {
		return err
	}

	return nil
}

func (orderBook *OrderBook) UpdateOrder(order *OrderItem, dryrun bool, blockHash common.Hash) error {
	return orderBook.ModifyOrder(order, order.OrderID, order.Price, dryrun, blockHash)
}

// ModifyOrder : modify the order
func (orderBook *OrderBook) ModifyOrder(order *OrderItem, orderID uint64, price *big.Int, dryrun bool, blockHash common.Hash) error {
	orderBook.UpdateTime()

	key := GetKeyFromBig(new(big.Int).SetUint64(order.OrderID))
	if order.Side == Bid {

		if orderBook.Bids.OrderExist(key, price, dryrun, blockHash) {
			return orderBook.Bids.UpdateOrder(order, dryrun, blockHash)
		}
	} else {

		if orderBook.Asks.OrderExist(key, price, dryrun, blockHash) {
			return orderBook.Asks.UpdateOrder(order, dryrun, blockHash)
		}
	}

	return nil
}

// VolumeAtPrice : get volume at the current price
func (orderBook *OrderBook) VolumeAtPrice(side string, price *big.Int, dryrun bool, blockHash common.Hash) *big.Int {
	volume := Zero()
	if side == Bid {
		if orderBook.Bids.PriceExist(price, dryrun, blockHash) {
			orderList := orderBook.Bids.PriceList(price, dryrun, blockHash)
			// incase we use cache for PriceList
			volume = CloneBigInt(orderList.Item.Volume)
		}
	} else {
		// other case
		if orderBook.Asks.PriceExist(price, dryrun, blockHash) {
			orderList := orderBook.Asks.PriceList(price, dryrun, blockHash)
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
