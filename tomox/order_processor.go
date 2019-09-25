package tomox

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/tomox/tomox_state"
	"math/big"
	"strconv"
	"time"
)

func ProcessOrder(statedb *state.StateDB, tomoXstatedb *tomox_state.TomoXStateDB, orderBook common.Hash, order *tomox_state.OrderItem) ([]map[string]string, *tomox_state.OrderItem, error) {
	var (
		orderInBook *tomox_state.OrderItem
		trades      []map[string]string
		err         error
	)
	if order.Status == OrderStatusCancelled {
		err := tomoXstatedb.CancerOrder(orderBook, order)
		if err != nil {
			log.Debug("Error when cancel order", "order", order)
			return nil, nil, err
		}
	}
	orderType := order.Type
	// if we do not use auto-increment orderid, we must set price slot to avoid conflict
	if orderType == Market {
		log.Debug("Process maket order", "side", order.Side, "quantity", order.Quantity, "price", order.Price)
		trades, orderInBook, err = processMarketOrder(statedb, tomoXstatedb, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
	} else {
		log.Debug("Process limit order", "side", order.Side, "quantity", order.Quantity, "price", order.Price)
		trades, orderInBook, err = processLimitOrder(statedb, tomoXstatedb, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
	}

	return trades, orderInBook, nil
}

// processMarketOrder : process the market order
func processMarketOrder(statedb *state.StateDB, tomoXstatedb *tomox_state.TomoXStateDB, orderBook common.Hash, order *tomox_state.OrderItem) ([]map[string]string, *tomox_state.OrderItem, error) {
	var (
		trades      []map[string]string
		newTrades   []map[string]string
		orderInBook *tomox_state.OrderItem
		err         error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	// speedup the comparison, do not assign because it is pointer
	zero := Zero()
	if side == Bid {
		bestPrice, volume := tomoXstatedb.GetBestAskPrice(orderBook)
		log.Debug("processMarketOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && bestPrice.Cmp(zero) > 0 {
			quantityToTrade, newTrades, orderInBook, err = processOrderList(statedb, tomoXstatedb, Ask, orderBook, bestPrice, quantityToTrade, order)
			if err != nil {
				return nil, orderInBook, err
			}
			trades = append(trades, newTrades...)
			bestPrice, volume = tomoXstatedb.GetBestAskPrice(orderBook)
			log.Debug("processMarketOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade, "volume", volume)
		}
	} else {
		bestPrice, volume := tomoXstatedb.GetBestBidPrice(orderBook)
		log.Debug("processMarketOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && bestPrice.Cmp(zero) > 0 {
			quantityToTrade, newTrades, orderInBook, err = processOrderList(statedb, tomoXstatedb, Bid, orderBook, bestPrice, quantityToTrade, order)
			if err != nil {
				return nil, orderInBook, err
			}
			trades = append(trades, newTrades...)
			bestPrice, volume = tomoXstatedb.GetBestBidPrice(orderBook)
			log.Debug("processMarketOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade, "volume", volume)
		}
	}
	return trades, orderInBook, nil
}

// processLimitOrder : process the limit order, can change the quote
// If not care for performance, we should make a copy of quote to prevent further reference problem
func processLimitOrder(statedb *state.StateDB, tomoXstatedb *tomox_state.TomoXStateDB, orderBook common.Hash, order *tomox_state.OrderItem) ([]map[string]string, *tomox_state.OrderItem, error) {
	var (
		trades      []map[string]string
		newTrades   []map[string]string
		orderInBook *tomox_state.OrderItem
		err         error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	price := order.Price

	// speedup the comparison, do not assign because it is pointer
	zero := Zero()

	if side == Bid {
		minPrice, volume := tomoXstatedb.GetBestAskPrice(orderBook)
		log.Debug("processLimitOrder ", "side", side, "minPrice", minPrice, "orderPrice", price, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && price.Cmp(minPrice) >= 0 && minPrice.Cmp(zero) > 0 {
			log.Debug("Min price in asks tree", "price", minPrice.String())
			quantityToTrade, newTrades, orderInBook, err = processOrderList(statedb, tomoXstatedb, Ask, orderBook, minPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			log.Debug("New trade found", "newTrades", newTrades, "orderInBook", orderInBook, "quantityToTrade", quantityToTrade)
			minPrice, volume = tomoXstatedb.GetBestAskPrice(orderBook)
			log.Debug("processLimitOrder ", "side", side, "minPrice", minPrice, "orderPrice", price, "volume", volume)
		}
	} else {
		maxPrice, volume := tomoXstatedb.GetBestBidPrice(orderBook)
		log.Debug("processLimitOrder ", "side", side, "maxPrice", maxPrice, "orderPrice", price, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && price.Cmp(maxPrice) <= 0 && maxPrice.Cmp(zero) > 0 {
			log.Debug("Max price in bids tree", "price", maxPrice.String())
			quantityToTrade, newTrades, orderInBook, err = processOrderList(statedb, tomoXstatedb, Bid, orderBook, maxPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			log.Debug("New trade found", "newTrades", newTrades, "orderInBook", orderInBook, "quantityToTrade", quantityToTrade)
			maxPrice, volume = tomoXstatedb.GetBestBidPrice(orderBook)
			log.Debug("processLimitOrder ", "side", side, "maxPrice", maxPrice, "orderPrice", price, "volume", volume)
		}
	}
	if quantityToTrade.Cmp(zero) > 0 {
		orderId := tomoXstatedb.GetNonce(orderBook)
		order.OrderID = orderId + 1
		order.Quantity = quantityToTrade
		tomoXstatedb.SetNonce(orderBook, orderId+1)
		tomoXstatedb.InsertOrderItem(orderBook, *order)
		log.Debug("After matching, order (unmatched part) is now added to tree", "side", order.Side, "order", order)
		orderInBook = order
	}
	return trades, orderInBook, nil
}

// processOrderList : process the order list
func processOrderList(statedb *state.StateDB, tomoXstatedb *tomox_state.TomoXStateDB, side string, orderBook common.Hash, price *big.Int, quantityStillToTrade *big.Int, order *tomox_state.OrderItem) (*big.Int, []map[string]string, *tomox_state.OrderItem, error) {
	quantityToTrade := CloneBigInt(quantityStillToTrade)
	log.Debug("Process matching between order and orderlist", "quantityToTrade", quantityToTrade)
	var (
		trades      []map[string]string
		orderInBook *tomox_state.OrderItem
	)
	// speedup the comparison, do not assign because it is pointer
	zero := Zero()
	orderId, amount, err := tomoXstatedb.GetBestOrderIdAndAmount(orderBook, price,side)
	if err != nil {
		return nil, nil, nil, err
	}
	oldestOrder := tomoXstatedb.GetOrder(orderBook, orderId)
	log.Debug("found order ", "orderId ", orderId, "side", oldestOrder.Side, "amount", amount)
	for amount.Cmp(zero) > 0 && quantityToTrade.Cmp(zero) > 0 {
		var (
			tradedQuantity *big.Int
		)
		if quantityToTrade.Cmp(amount) <= 0 {
			tradedQuantity = CloneBigInt(quantityToTrade)
			quantityToTrade = Zero()
			orderInBook = &oldestOrder
		} else {
			tradedQuantity = CloneBigInt(amount)
			quantityToTrade = Sub(quantityToTrade, tradedQuantity)
		}
		tomoXstatedb.SubAmountOrderItem(orderBook, orderId, price, amount, side)
		log.Debug("Update quantity for orderId", "orderId", orderId.Hex())
		log.Debug("TRADE", "orderBook", orderBook, "Price 1", price, "Price 2", order.Price, "Amount", tradedQuantity, "orderId", orderId, "side", side)

		transactionRecord := make(map[string]string)
		transactionRecord[TradeTakerOrderHash] = hex.EncodeToString(order.Hash.Bytes())
		transactionRecord[TradeMakerOrderHash] = hex.EncodeToString(oldestOrder.Hash.Bytes())
		transactionRecord[TradeTimestamp] = strconv.FormatInt(time.Now().Unix(), 10)
		transactionRecord[TradeQuantity] = tradedQuantity.String()
		transactionRecord[TradeMakerExchange] = oldestOrder.ExchangeAddress.String()
		transactionRecord[TradeMaker] = oldestOrder.UserAddress.String()
		transactionRecord[TradeBaseToken] = oldestOrder.BaseToken.String()
		transactionRecord[TradeQuoteToken] = oldestOrder.QuoteToken.String()
		// maker price is actual price
		// taker price is offer price
		// tradedPrice is always actual price
		transactionRecord[TradePrice] = oldestOrder.Price.String()

		trades = append(trades, transactionRecord)
		orderId, amount, err = tomoXstatedb.GetBestOrderIdAndAmount(orderBook, price,side)
		if err != nil {
			return nil, nil, nil, err
		}
		oldestOrder = tomoXstatedb.GetOrder(orderBook, orderId)
		log.Debug("found order ", "orderId ", orderId, "side", oldestOrder.Side, "amount", amount)
	}
	return quantityToTrade, trades, orderInBook, nil
}
