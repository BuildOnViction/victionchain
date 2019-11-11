package tomox

import (
	"encoding/json"
	"errors"
	"github.com/tomochain/tomochain/consensus"
	"math/big"
	"strconv"
	"time"

	"fmt"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/tomox/tomox_state"
)

var emptyAddress = common.StringToAddress("")
var errQuantityTradeTooSmall = errors.New("Quantity trade too small")

func (tomox *TomoX) CommitOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tomoXstatedb *tomox_state.TomoXStateDB, orderBook common.Hash, order *tomox_state.OrderItem) ([]map[string]string, []*tomox_state.OrderItem, error) {
	snap := tomoXstatedb.Snapshot()
	trades, rejects, err := tomox.ApplyOrder(coinbase, chain, statedb, tomoXstatedb, orderBook, order)
	if err != nil {
		tomoXstatedb.RevertToSnapshot(snap)
		return nil, nil, err
	}
	return trades, rejects, err
}

func (tomox *TomoX) ApplyOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tomoXstatedb *tomox_state.TomoXStateDB, orderBook common.Hash, order *tomox_state.OrderItem) ([]map[string]string, []*tomox_state.OrderItem, error) {
	var (
		rejects []*tomox_state.OrderItem
		trades  []map[string]string
		err     error
	)
	nonce := tomoXstatedb.GetNonce(order.UserAddress.Hash())
	log.Debug("ApplyOrder", "addr", order.UserAddress, "statenonce", nonce, "ordernonce", order.Nonce)
	if big.NewInt(int64(nonce)).Cmp(order.Nonce) == -1 {
		return nil, nil, ErrNonceTooHigh
	} else if big.NewInt(int64(nonce)).Cmp(order.Nonce) == 1 {
		return nil, nil, ErrNonceTooLow
	}
	if order.Price.Sign() == 0 || common.BigToHash(order.Price).Big().Cmp(order.Price) != 0 {
		log.Debug("Reject order price invalid", "price", order.Price)
		rejects = append(rejects, order)
		tomoXstatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
		return trades, rejects, nil
	}
	if order.Quantity.Sign() == 0 || common.BigToHash(order.Quantity).Big().Cmp(order.Quantity) != 0 {
		log.Debug("Reject order quantity invalid", "quantity", order.Quantity)
		rejects = append(rejects, order)
		tomoXstatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
		return trades, rejects, nil
	}

	if order.Status == OrderStatusCancelled {
		err := tomoXstatedb.CancelOrder(orderBook, order)
		if err != nil {
			log.Debug("Error when cancel order", "order", order)
			return nil, nil, err
		}
		log.Debug("Exchange add user nonce:", "address", order.UserAddress, "status", order.Status, "nonce", nonce+1)
		tomoXstatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
		return trades, rejects, nil
	}
	orderType := order.Type
	// if we do not use auto-increment orderid, we must set price slot to avoid conflict
	if orderType == tomox_state.Market {
		log.Debug("Process maket order", "side", order.Side, "quantity", order.Quantity, "price", order.Price)
		trades, rejects, err = tomox.processMarketOrder(coinbase, chain, statedb, tomoXstatedb, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
	} else {
		log.Debug("Process limit order", "side", order.Side, "quantity", order.Quantity, "price", order.Price)
		trades, rejects, err = tomox.processLimitOrder(coinbase, chain, statedb, tomoXstatedb, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
	}

	log.Debug("Exchange add user nonce:", "address", order.UserAddress, "status", order.Status, "nonce", nonce+1)
	tomoXstatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
	return trades, rejects, nil
}

// processMarketOrder : process the market order
func (tomox *TomoX) processMarketOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tomoXstatedb *tomox_state.TomoXStateDB, orderBook common.Hash, order *tomox_state.OrderItem) ([]map[string]string, []*tomox_state.OrderItem, error) {
	var (
		trades     []map[string]string
		newTrades  []map[string]string
		rejects    []*tomox_state.OrderItem
		newRejects []*tomox_state.OrderItem
		err        error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	// speedup the comparison, do not assign because it is pointer
	zero := tomox_state.Zero
	if side == tomox_state.Bid {
		bestPrice, volume := tomoXstatedb.GetBestAskPrice(orderBook)
		log.Debug("processMarketOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && bestPrice.Cmp(zero) > 0 {
			quantityToTrade, newTrades, newRejects, err = tomox.processOrderList(coinbase, chain, statedb, tomoXstatedb, tomox_state.Ask, orderBook, bestPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			bestPrice, volume = tomoXstatedb.GetBestAskPrice(orderBook)
			log.Debug("processMarketOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade, "volume", volume)
		}
	} else {
		bestPrice, volume := tomoXstatedb.GetBestBidPrice(orderBook)
		log.Debug("processMarketOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && bestPrice.Cmp(zero) > 0 {
			quantityToTrade, newTrades, newRejects, err = tomox.processOrderList(coinbase, chain, statedb, tomoXstatedb, tomox_state.Bid, orderBook, bestPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			bestPrice, volume = tomoXstatedb.GetBestBidPrice(orderBook)
			log.Debug("processMarketOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade, "volume", volume)
		}
	}
	return trades, newRejects, nil
}

// processLimitOrder : process the limit order, can change the quote
// If not care for performance, we should make a copy of quote to prevent further reference problem
func (tomox *TomoX) processLimitOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tomoXstatedb *tomox_state.TomoXStateDB, orderBook common.Hash, order *tomox_state.OrderItem) ([]map[string]string, []*tomox_state.OrderItem, error) {
	var (
		trades     []map[string]string
		newTrades  []map[string]string
		rejects    []*tomox_state.OrderItem
		newRejects []*tomox_state.OrderItem
		err        error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	price := order.Price

	// speedup the comparison, do not assign because it is pointer
	zero := tomox_state.Zero

	if side == tomox_state.Bid {
		minPrice, volume := tomoXstatedb.GetBestAskPrice(orderBook)
		log.Debug("processLimitOrder ", "side", side, "minPrice", minPrice, "orderPrice", price, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && price.Cmp(minPrice) >= 0 && minPrice.Cmp(zero) > 0 {
			log.Debug("Min price in asks tree", "price", minPrice.String())
			quantityToTrade, newTrades, newRejects, err = tomox.processOrderList(coinbase, chain, statedb, tomoXstatedb, tomox_state.Ask, orderBook, minPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			log.Debug("New trade found", "newTrades", newTrades, "quantityToTrade", quantityToTrade)
			minPrice, volume = tomoXstatedb.GetBestAskPrice(orderBook)
			log.Debug("processLimitOrder ", "side", side, "minPrice", minPrice, "orderPrice", price, "volume", volume)
		}
	} else {
		maxPrice, volume := tomoXstatedb.GetBestBidPrice(orderBook)
		log.Debug("processLimitOrder ", "side", side, "maxPrice", maxPrice, "orderPrice", price, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && price.Cmp(maxPrice) <= 0 && maxPrice.Cmp(zero) > 0 {
			log.Debug("Max price in bids tree", "price", maxPrice.String())
			quantityToTrade, newTrades, newRejects, err = tomox.processOrderList(coinbase, chain, statedb, tomoXstatedb, tomox_state.Bid, orderBook, maxPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			log.Debug("New trade found", "newTrades", newTrades, "quantityToTrade", quantityToTrade)
			maxPrice, volume = tomoXstatedb.GetBestBidPrice(orderBook)
			log.Debug("processLimitOrder ", "side", side, "maxPrice", maxPrice, "orderPrice", price, "volume", volume)
		}
	}
	if quantityToTrade.Cmp(zero) > 0 {
		orderId := tomoXstatedb.GetNonce(orderBook)
		order.OrderID = orderId + 1
		order.Quantity = quantityToTrade
		tomoXstatedb.SetNonce(orderBook, orderId+1)
		orderIdHash := common.BigToHash(new(big.Int).SetUint64(order.OrderID))
		tomoXstatedb.InsertOrderItem(orderBook, orderIdHash, *order)
		log.Debug("After matching, order (unmatched part) is now added to tree", "side", order.Side, "order", order)
	}
	return trades, rejects, nil
}

// processOrderList : process the order list
func (tomox *TomoX) processOrderList(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tomoXstatedb *tomox_state.TomoXStateDB, side string, orderBook common.Hash, price *big.Int, quantityStillToTrade *big.Int, order *tomox_state.OrderItem) (*big.Int, []map[string]string, []*tomox_state.OrderItem, error) {
	quantityToTrade := tomox_state.CloneBigInt(quantityStillToTrade)
	log.Debug("Process matching between order and orderlist", "quantityToTrade", quantityToTrade)
	var (
		trades []map[string]string

		rejects []*tomox_state.OrderItem
	)
	for quantityToTrade.Sign() > 0 {
		orderId, amount, _ := tomoXstatedb.GetBestOrderIdAndAmount(orderBook, price, side)
		var oldestOrder tomox_state.OrderItem
		if amount.Sign() > 0 {
			oldestOrder = tomoXstatedb.GetOrder(orderBook, orderId)
		}
		log.Debug("found order ", "orderId ", orderId, "side", oldestOrder.Side, "amount", amount)
		if oldestOrder.Quantity == nil || oldestOrder.Quantity.Sign() == 0 && amount.Sign() == 0 {
			break
		}
		var (
			tradedQuantity    *big.Int
			maxTradedQuantity *big.Int
		)
		if quantityToTrade.Cmp(amount) <= 0 {
			maxTradedQuantity = tomox_state.CloneBigInt(quantityToTrade)
		} else {
			maxTradedQuantity = tomox_state.CloneBigInt(amount)
		}
		var quotePrice *big.Int
		if oldestOrder.QuoteToken.String() != common.TomoNativeAddress {
			quotePrice = tomoXstatedb.GetPrice(tomox_state.GetOrderBookHash(oldestOrder.QuoteToken, common.HexToAddress(common.TomoNativeAddress)))
		}
		tradedQuantity, rejectMaker, err := tomox.getTradeQuantity(quotePrice, coinbase, chain, statedb, order, &oldestOrder, maxTradedQuantity)
		if err != nil && err == errQuantityTradeTooSmall {
			if tradedQuantity.Cmp(maxTradedQuantity) == 0 {
				if quantityToTrade.Cmp(amount) == 0 { // reject Taker & maker
					rejects = append(rejects, order)
					quantityToTrade = tomox_state.Zero
					rejects = append(rejects, &oldestOrder)
					err = tomoXstatedb.CancelOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					break
				} else if quantityToTrade.Cmp(amount) < 0 { // reject Taker
					rejects = append(rejects, order)
					quantityToTrade = tomox_state.Zero
					break
				} else { // reject maker
					rejects = append(rejects, &oldestOrder)
					err = tomoXstatedb.CancelOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					continue
				}
			} else {
				if rejectMaker { // reject maker
					rejects = append(rejects, &oldestOrder)
					err = tomoXstatedb.CancelOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					continue
				} else { // reject Taker
					rejects = append(rejects, order)
					quantityToTrade = tomox_state.Zero
					break
				}
			}
		} else if err != nil {
			return nil, nil, nil, err
		}
		if tradedQuantity.Sign() == 0 && !rejectMaker {
			log.Debug("Reject order Taker ", "tradedQuantity", tradedQuantity, "rejectMaker", rejectMaker)
			rejects = append(rejects, order)
			quantityToTrade = tomox_state.Zero
			break
		}
		if tradedQuantity.Sign() > 0 {
			quantityToTrade = tomox_state.Sub(quantityToTrade, tradedQuantity)
			tomoXstatedb.SubAmountOrderItem(orderBook, orderId, price, tradedQuantity, side)
			if oldestOrder.QuoteToken.String() == common.TomoNativeAddress {
				tomoXstatedb.SetPrice(orderBook, price)
			}
			log.Debug("Update quantity for orderId", "orderId", orderId.Hex())
			log.Debug("TRADE", "orderBook", orderBook, "Taker price", price, "maker price", order.Price, "Amount", tradedQuantity, "orderId", orderId, "side", side)

			transactionRecord := make(map[string]string)
			transactionRecord[TradeTakerOrderHash] = order.Hash.Hex()
			transactionRecord[TradeMakerOrderHash] = oldestOrder.Hash.Hex()
			transactionRecord[TradeTimestamp] = strconv.FormatInt(time.Now().Unix(), 10)
			transactionRecord[TradeQuantity] = tradedQuantity.String()
			transactionRecord[TradeMakerExchange] = oldestOrder.ExchangeAddress.String()
			transactionRecord[TradeMaker] = oldestOrder.UserAddress.String()
			transactionRecord[TradeBaseToken] = oldestOrder.BaseToken.String()
			transactionRecord[TradeQuoteToken] = oldestOrder.QuoteToken.String()
			// maker price is actual price
			// Taker price is offer price
			// tradedPrice is always actual price
			transactionRecord[TradePrice] = oldestOrder.Price.String()

			trades = append(trades, transactionRecord)
		}
		if rejectMaker {
			rejects = append(rejects, &oldestOrder)
			err := tomoXstatedb.CancelOrder(orderBook, &oldestOrder)
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}
	return quantityToTrade, trades, rejects, nil
}

func (tomox *TomoX) getTradeQuantity(quotePrice *big.Int, coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, takerOrder *tomox_state.OrderItem, makerOrder *tomox_state.OrderItem, quantityToTrade *big.Int) (*big.Int, bool, error) {
	baseTokenDecimal, err := tomox.GetTokenDecimal(chain, statedb, coinbase, makerOrder.BaseToken)
	if err != nil || baseTokenDecimal.Sign() == 0 {
		return tomox_state.Zero, false, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", makerOrder.BaseToken.String(), err)
	}
	quoteTokenDecimal, err := tomox.GetTokenDecimal(chain, statedb, coinbase, makerOrder.QuoteToken)
	if err != nil || quoteTokenDecimal.Sign() == 0 {
		return tomox_state.Zero, false, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", makerOrder.QuoteToken.String(), err)
	}
	if makerOrder.QuoteToken.String() == common.TomoNativeAddress {
		quotePrice = quoteTokenDecimal
	}
	if takerOrder.ExchangeAddress.String() == makerOrder.ExchangeAddress.String() {
		if err := tomox_state.CheckRelayerFee(takerOrder.ExchangeAddress, new(big.Int).Mul(common.RelayerFee, big.NewInt(2)), statedb); err != nil {
			log.Debug("Reject order Taker Exchnage = Maker Exchange , relayer not enough fee ", "err", err)
			return tomox_state.Zero, false, nil
		}
	} else {
		if err := tomox_state.CheckRelayerFee(takerOrder.ExchangeAddress, common.RelayerFee, statedb); err != nil {
			log.Debug("Reject order Taker , relayer not enough fee ", "err", err)
			return tomox_state.Zero, false, nil
		}
		if err := tomox_state.CheckRelayerFee(makerOrder.ExchangeAddress, common.RelayerFee, statedb); err != nil {
			log.Debug("Reject order maker , relayer not enough fee ", "err", err)
			return tomox_state.Zero, true, nil
		}
	}
	takerFeeRate := tomox_state.GetExRelayerFee(takerOrder.ExchangeAddress, statedb)
	makerFeeRate := tomox_state.GetExRelayerFee(makerOrder.ExchangeAddress, statedb)
	var takerBalance, makerBalance *big.Int
	switch takerOrder.Side {
	case tomox_state.Bid:
		takerBalance = tomox_state.GetTokenBalance(takerOrder.UserAddress, makerOrder.QuoteToken, statedb)
		makerBalance = tomox_state.GetTokenBalance(makerOrder.UserAddress, makerOrder.BaseToken, statedb)
	case tomox_state.Ask:
		takerBalance = tomox_state.GetTokenBalance(takerOrder.UserAddress, makerOrder.BaseToken, statedb)
		makerBalance = tomox_state.GetTokenBalance(makerOrder.UserAddress, makerOrder.QuoteToken, statedb)
	default:
		takerBalance = big.NewInt(0)
		makerBalance = big.NewInt(0)
	}
	quantity, rejectMaker := GetTradeQuantity(takerOrder.Side, takerFeeRate, takerBalance, makerOrder.Price, makerFeeRate, makerBalance, baseTokenDecimal, quantityToTrade)
	log.Debug("GetTradeQuantity", "side", takerOrder.Side, "takerBalance", takerBalance, "makerBalance", makerBalance, "BaseToken", makerOrder.BaseToken, "QuoteToken", makerOrder.QuoteToken, "quantity", quantity, "rejectMaker", rejectMaker, "quotePrice", quotePrice)
	if quantity.Sign() > 0 {
		// Apply Match Order
		setteBalance, err := GetSettleBalance(quotePrice, takerOrder.Side, takerFeeRate, makerOrder.BaseToken, makerOrder.QuoteToken, makerOrder.Price, makerFeeRate, baseTokenDecimal, quoteTokenDecimal, quantity)
		log.Debug("GetSettleBalance", "setteBalance", setteBalance, "err", err)
		if err == nil {
			err = SetteBalance(coinbase, takerOrder, makerOrder, setteBalance, statedb)
		}
		return quantity, rejectMaker, err
	}
	return quantity, rejectMaker, nil
}

func GetTradeQuantity(takerSide string, takerFeeRate *big.Int, takerBalance *big.Int, makerPrice *big.Int, makerFeeRate *big.Int, makerBalance *big.Int, baseTokenDecimal *big.Int, quantityToTrade *big.Int) (*big.Int, bool) {
	if takerSide == tomox_state.Bid {
		// maker InQuantity quoteTokenQuantity=(quantityToTrade*maker.Price/baseTokenDecimal)
		quoteTokenQuantity := new(big.Int).Mul(quantityToTrade, makerPrice)
		quoteTokenQuantity = quoteTokenQuantity.Div(quoteTokenQuantity, baseTokenDecimal)
		// Fee
		// charge on the token he/she has before the trade, in this case: quoteToken
		// charge on the token he/she has before the trade, in this case: baseToken
		// takerFee = quoteTokenQuantity*takerFeeRate/baseFee=(quantityToTrade*maker.Price/baseTokenDecimal) * makerFeeRate/baseFee
		takerFee := big.NewInt(0).Mul(quoteTokenQuantity, takerFeeRate)
		takerFee = big.NewInt(0).Div(takerFee, common.TomoXBaseFee)
		//takerOutTotal= quoteTokenQuantity + takerFee =  quantityToTrade*maker.Price/baseTokenDecimal + quantityToTrade*maker.Price/baseTokenDecimal * takerFeeRate/baseFee
		// = quantityToTrade *  maker.Price/baseTokenDecimal ( 1 +  takerFeeRate/baseFee)
		// = quantityToTrade * maker.Price * (baseFee + takerFeeRate ) / ( baseTokenDecimal * baseFee)
		takerOutTotal := new(big.Int).Add(quoteTokenQuantity, takerFee)
		makerOutTotal := quantityToTrade
		if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			return quantityToTrade, false
		} else if takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			newQuantityTrade := new(big.Int).Mul(takerBalance, baseTokenDecimal)
			newQuantityTrade = newQuantityTrade.Mul(newQuantityTrade, common.TomoXBaseFee)
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, new(big.Int).Add(common.TomoXBaseFee, takerFeeRate))
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, makerPrice)
			if newQuantityTrade.Sign() == 0 {
				log.Debug("Reject order Taker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "takerOutTotal", takerOutTotal)
			}
			return newQuantityTrade, false
		} else if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) < 0 {
			log.Debug("Reject order maker , not enough balance ", "makerBalance", makerBalance, " makerOutTotal", makerOutTotal)
			return makerBalance, true
		} else {
			// takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) < 0
			newQuantityTrade := new(big.Int).Mul(takerBalance, baseTokenDecimal)
			newQuantityTrade = newQuantityTrade.Mul(newQuantityTrade, common.TomoXBaseFee)
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, new(big.Int).Add(common.TomoXBaseFee, takerFeeRate))
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, makerPrice)
			if newQuantityTrade.Cmp(makerBalance) <= 0 {
				if newQuantityTrade.Sign() == 0 {
					log.Debug("Reject order Taker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "makerBalance", makerBalance, " newQuantityTrade ", newQuantityTrade)
				}
				return newQuantityTrade, false
			}
			log.Debug("Reject order maker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "makerBalance", makerBalance, " newQuantityTrade ", newQuantityTrade)
			return makerBalance, true
		}
	} else {
		// Taker InQuantity
		// quoteTokenQuantity = quantityToTrade * makerPrice / baseTokenDecimal
		quoteTokenQuantity := new(big.Int).Mul(quantityToTrade, makerPrice)
		quoteTokenQuantity = quoteTokenQuantity.Div(quoteTokenQuantity, baseTokenDecimal)
		// maker InQuantity

		// Fee
		// charge on the token he/she has before the trade, in this case: baseToken
		// makerFee = quoteTokenQuantity * makerFeeRate / baseFee = quantityToTrade * makerPrice / baseTokenDecimal * makerFeeRate / baseFee
		// charge on the token he/she has before the trade, in this case: quoteToken
		makerFee := new(big.Int).Mul(quoteTokenQuantity, makerFeeRate)
		makerFee = makerFee.Div(makerFee, common.TomoXBaseFee)

		takerOutTotal := quantityToTrade
		// makerOutTotal = quoteTokenQuantity + makerFee  = quantityToTrade * makerPrice / baseTokenDecimal + quantityToTrade * makerPrice / baseTokenDecimal * makerFeeRate / baseFee
		// =  quantityToTrade * makerPrice / baseTokenDecimal * (1+makerFeeRate / baseFee)
		// = quantityToTrade  * makerPrice * (baseFee + makerFeeRate) / ( baseTokenDecimal * baseFee )
		makerOutTotal := new(big.Int).Add(quoteTokenQuantity, makerFee)
		if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			return quantityToTrade, false
		} else if takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			if takerBalance.Sign() == 0 {
				log.Debug("Reject order Taker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "takerOutTotal", takerOutTotal)
			}
			return takerBalance, false
		} else if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) < 0 {
			newQuantityTrade := new(big.Int).Mul(makerBalance, baseTokenDecimal)
			newQuantityTrade = newQuantityTrade.Mul(newQuantityTrade, common.TomoXBaseFee)
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, new(big.Int).Add(common.TomoXBaseFee, makerFeeRate))
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, makerPrice)
			log.Debug("Reject order maker , not enough balance ", "makerBalance", makerBalance, " makerOutTotal", makerOutTotal)
			return newQuantityTrade, true
		} else {
			// takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) < 0
			newQuantityTrade := new(big.Int).Mul(makerBalance, baseTokenDecimal)
			newQuantityTrade = newQuantityTrade.Mul(newQuantityTrade, common.TomoXBaseFee)
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, new(big.Int).Add(common.TomoXBaseFee, makerFeeRate))
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, makerPrice)
			if newQuantityTrade.Cmp(takerBalance) <= 0 {
				log.Debug("Reject order maker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "makerBalance", makerBalance, " newQuantityTrade ", newQuantityTrade)
				return newQuantityTrade, true
			}
			if takerBalance.Sign() == 0 {
				log.Debug("Reject order Taker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "makerBalance", makerBalance, " newQuantityTrade ", newQuantityTrade)
			}
			return takerBalance, false
		}
	}
}

type TradeResult struct {
	Fee         *big.Int
	InToken     common.Address
	InQuantity  *big.Int
	InTotal     *big.Int
	OutToken    common.Address
	OutQuantity *big.Int
	OutTotal    *big.Int
}
type SettleBalance struct {
	Taker TradeResult
	Maker TradeResult
}

func (settleBalance *SettleBalance) String() string {
	json, _ := json.Marshal(settleBalance)
	return string(json)
}

func GetSettleBalance(quotePrice *big.Int, takerSide string, takerFeeRate *big.Int, baseToken, quoteToken common.Address, makerPrice *big.Int, makerFeeRate *big.Int, baseTokenDecimal *big.Int, quoteTokenDecimal *big.Int, quantityToTrade *big.Int) (*SettleBalance, error) {
	log.Debug("GetSettleBalance", "takerSide", takerSide, "takerFeeRate", takerFeeRate, "baseToken", baseToken, "quoteToken", quoteToken, "makerPrice", makerPrice, "makerFeeRate", makerFeeRate, "baseTokenDecimal", baseTokenDecimal, "quantityToTrade", quantityToTrade)
	var result *SettleBalance
	//result = map[common.Address]map[string]interface{}{}
	if takerSide == tomox_state.Bid {
		// maker InQuantity quoteTokenQuantity=(quantityToTrade*maker.Price/baseTokenDecimal)
		quoteTokenQuantity := new(big.Int).Mul(quantityToTrade, makerPrice)
		quoteTokenQuantity = quoteTokenQuantity.Div(quoteTokenQuantity, baseTokenDecimal)
		// Fee
		// charge on the token he/she has before the trade, in this case: quoteToken
		// charge on the token he/she has before the trade, in this case: baseToken
		// takerFee = quoteTokenQuantity*takerFeeRate/baseFee=(quantityToTrade*maker.Price/baseTokenDecimal) * makerFeeRate/baseFee
		takerFee := new(big.Int).Mul(quoteTokenQuantity, takerFeeRate)
		takerFee = new(big.Int).Div(takerFee, common.TomoXBaseFee)
		// charge on the token he/she has before the trade, in this case: baseToken
		makerFee := new(big.Int).Mul(quoteTokenQuantity, makerFeeRate)
		makerFee = new(big.Int).Div(makerFee, common.TomoXBaseFee)
		if quoteTokenQuantity.Cmp(makerFee) <= 0 {
			log.Debug("quantity trade too small", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee)
			return result, errQuantityTradeTooSmall
		}
		if quotePrice != nil && quotePrice.Cmp(common.Big0) > 0 {
			exMakerReceivedFee := new(big.Int).Mul(makerFee, quotePrice)
			exMakerReceivedFee = exMakerReceivedFee.Div(exMakerReceivedFee, quoteTokenDecimal)
			log.Debug("exMakerReceivedFee", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "quotePrice", quotePrice)
			if exMakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("makerFee too small", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "quotePrice", quotePrice)
				return result, errQuantityTradeTooSmall
			}
			exTakerReceivedFee := new(big.Int).Mul(takerFee, quotePrice)
			exTakerReceivedFee = exTakerReceivedFee.Div(exTakerReceivedFee, quoteTokenDecimal)
			log.Debug("exTakerReceivedFee", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "quotePrice", quotePrice)
			if exTakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("takerFee too small", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "quotePrice", quotePrice)
				return result, errQuantityTradeTooSmall
			}
		}
		inTotal := new(big.Int).Sub(quoteTokenQuantity, makerFee)
		//takerOutTotal= quoteTokenQuantity + takerFee =  quantityToTrade*maker.Price/baseTokenDecimal + quantityToTrade*maker.Price/baseTokenDecimal * takerFeeRate/baseFee
		// = quantityToTrade *  maker.Price/baseTokenDecimal ( 1 +  takerFeeRate/baseFee)
		// = quantityToTrade * maker.Price * (baseFee + takerFeeRate ) / ( baseTokenDecimal * baseFee)
		takerOutTotal := new(big.Int).Add(quoteTokenQuantity, takerFee)

		result = &SettleBalance{
			Taker: TradeResult{
				Fee:         takerFee,
				InToken:     baseToken,
				InQuantity:  quantityToTrade,
				InTotal:     quantityToTrade,
				OutToken:    quoteToken,
				OutQuantity: quoteTokenQuantity,
				OutTotal:    takerOutTotal,
			},
			Maker: TradeResult{
				Fee:         makerFee,
				InToken:     quoteToken,
				InQuantity:  quoteTokenQuantity,
				InTotal:     inTotal,
				OutToken:    baseToken,
				OutQuantity: quantityToTrade,
				OutTotal:    quantityToTrade,
			},
		}
	} else {
		// Taker InQuantity
		// quoteTokenQuantity = quantityToTrade * makerPrice / baseTokenDecimal
		quoteTokenQuantity := new(big.Int).Mul(quantityToTrade, makerPrice)
		quoteTokenQuantity = quoteTokenQuantity.Div(quoteTokenQuantity, baseTokenDecimal)
		// maker InQuantity

		// Fee
		// charge on the token he/she has before the trade, in this case: baseToken
		// makerFee = quoteTokenQuantity * makerFeeRate / baseFee = quantityToTrade * makerPrice / baseTokenDecimal * makerFeeRate / baseFee
		// charge on the token he/she has before the trade, in this case: quoteToken
		makerFee := new(big.Int).Mul(quoteTokenQuantity, makerFeeRate)
		makerFee = makerFee.Div(makerFee, common.TomoXBaseFee)

		// charge on the token he/she has before the trade, in this case: baseToken
		takerFee := new(big.Int).Mul(quoteTokenQuantity, takerFeeRate)
		takerFee = new(big.Int).Div(takerFee, common.TomoXBaseFee)
		if quoteTokenQuantity.Cmp(takerFee) <= 0 {
			log.Debug("quantity trade too small", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee)
			return result, errQuantityTradeTooSmall
		}
		if quotePrice != nil && quotePrice.Cmp(common.Big0) > 0 {
			exMakerReceivedFee := new(big.Int).Mul(makerFee, quotePrice)
			exMakerReceivedFee = exMakerReceivedFee.Div(exMakerReceivedFee, quoteTokenDecimal)
			log.Debug("exMakerReceivedFee", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "quotePrice", quotePrice)
			if exMakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("makerFee too small", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "quotePrice", quotePrice)
				return result, errQuantityTradeTooSmall
			}
			exTakerReceivedFee := new(big.Int).Mul(takerFee, quotePrice)
			exTakerReceivedFee = exTakerReceivedFee.Div(exTakerReceivedFee, quoteTokenDecimal)
			log.Debug("exTakerReceivedFee", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "quotePrice", quotePrice)
			if exTakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("takerFee too small", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "quotePrice", quotePrice)
				return result, errQuantityTradeTooSmall
			}
		}
		inTotal := new(big.Int).Sub(quoteTokenQuantity, takerFee)
		// makerOutTotal = quoteTokenQuantity + makerFee  = quantityToTrade * makerPrice / baseTokenDecimal + quantityToTrade * makerPrice / baseTokenDecimal * makerFeeRate / baseFee
		// =  quantityToTrade * makerPrice / baseTokenDecimal * (1+makerFeeRate / baseFee)
		// = quantityToTrade  * makerPrice * (baseFee + makerFeeRate) / ( baseTokenDecimal * baseFee )
		makerOutTotal := new(big.Int).Add(quoteTokenQuantity, makerFee)
		// Fee
		result = &SettleBalance{
			Taker: TradeResult{
				Fee:         takerFee,
				InToken:     quoteToken,
				InQuantity:  quoteTokenQuantity,
				InTotal:     inTotal,
				OutToken:    baseToken,
				OutQuantity: quantityToTrade,
				OutTotal:    quantityToTrade,
			},
			Maker: TradeResult{
				Fee:         makerFee,
				InToken:     baseToken,
				InQuantity:  quantityToTrade,
				InTotal:     quantityToTrade,
				OutToken:    quoteToken,
				OutQuantity: quoteTokenQuantity,
				OutTotal:    makerOutTotal,
			},
		}
	}
	return result, nil
}

func SetteBalance(coinbase common.Address, takerOrder, makerOrder *tomox_state.OrderItem, settleBalance *SettleBalance, statedb *state.StateDB) error {
	takerExOwner := tomox_state.GetRelayerOwner(takerOrder.ExchangeAddress, statedb)
	makerExOwner := tomox_state.GetRelayerOwner(makerOrder.ExchangeAddress, statedb)
	matchingFee := big.NewInt(0)
	// masternodes charges fee of both 2 relayers. If maker and Taker are on same relayer, that relayer is charged fee twice
	matchingFee = matchingFee.Add(matchingFee, common.RelayerFee)
	matchingFee = matchingFee.Add(matchingFee, common.RelayerFee)

	if common.EmptyHash(takerExOwner.Hash()) || common.EmptyHash(makerExOwner.Hash()) {
		return fmt.Errorf("Echange owner empty , Taker: %v , maker : %v ", takerExOwner, makerExOwner)
	}
	mapBalances := map[common.Address]map[common.Address]*big.Int{}
	//Checking balance
	newTakerInTotal, err := tomox_state.CheckAddTokenBalance(takerOrder.UserAddress, settleBalance.Taker.InTotal, settleBalance.Taker.InToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Taker.InToken] == nil {
		mapBalances[settleBalance.Taker.InToken] = map[common.Address]*big.Int{}
		mapBalances[settleBalance.Taker.InToken][takerOrder.UserAddress] = newTakerInTotal
	}
	newTakerOutTotal, err := tomox_state.CheckSubTokenBalance(takerOrder.UserAddress, settleBalance.Taker.OutTotal, settleBalance.Taker.OutToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Taker.OutToken] == nil {
		mapBalances[settleBalance.Taker.OutToken] = map[common.Address]*big.Int{}
		mapBalances[settleBalance.Taker.OutToken][takerOrder.UserAddress] = newTakerOutTotal
	}
	newMakerInTotal, err := tomox_state.CheckAddTokenBalance(makerOrder.UserAddress, settleBalance.Maker.InTotal, settleBalance.Maker.InToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Maker.InToken] == nil {
		mapBalances[settleBalance.Maker.InToken] = map[common.Address]*big.Int{}
		mapBalances[settleBalance.Maker.InToken][makerOrder.UserAddress] = newMakerInTotal
	}
	newMakerOutTotal, err := tomox_state.CheckSubTokenBalance(makerOrder.UserAddress, settleBalance.Maker.OutTotal, settleBalance.Maker.OutToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Maker.OutToken] == nil {
		mapBalances[settleBalance.Maker.OutToken] = map[common.Address]*big.Int{}
		mapBalances[settleBalance.Maker.OutToken][makerOrder.UserAddress] = newMakerOutTotal
	}
	newTakerFee, err := tomox_state.CheckAddTokenBalance(takerExOwner, settleBalance.Taker.Fee, makerOrder.QuoteToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[makerOrder.QuoteToken] == nil {
		mapBalances[makerOrder.QuoteToken] = map[common.Address]*big.Int{}
		mapBalances[makerOrder.QuoteToken][takerExOwner] = newTakerFee
	}
	newMakerFee, err := tomox_state.CheckAddTokenBalance(makerExOwner, settleBalance.Maker.Fee, makerOrder.QuoteToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	mapBalances[makerOrder.QuoteToken][makerExOwner] = newMakerFee

	mapRelayerFee := map[common.Address]*big.Int{}
	newRelayerTakerFee, err := tomox_state.CheckSubRelayerFee(takerOrder.ExchangeAddress, common.RelayerFee, statedb, mapRelayerFee)
	if err != nil {
		return err
	}
	mapRelayerFee[takerOrder.ExchangeAddress] = newRelayerTakerFee
	newRelayerMakerFee, err := tomox_state.CheckSubRelayerFee(makerOrder.ExchangeAddress, common.RelayerFee, statedb, mapRelayerFee)
	if err != nil {
		return err
	}
	mapRelayerFee[makerOrder.ExchangeAddress] = newRelayerMakerFee
	tomox_state.SetSubRelayerFee(takerOrder.ExchangeAddress, newRelayerTakerFee, common.RelayerFee, statedb)
	tomox_state.SetSubRelayerFee(makerOrder.ExchangeAddress, newRelayerMakerFee, common.RelayerFee, statedb)

	masternodeOwner := statedb.GetOwner(coinbase)
	statedb.AddBalance(masternodeOwner, matchingFee)

	tomox_state.SetTokenBalance(takerOrder.UserAddress, newTakerInTotal, settleBalance.Taker.InToken, statedb)
	tomox_state.SetTokenBalance(takerOrder.UserAddress, newTakerOutTotal, settleBalance.Taker.OutToken, statedb)

	tomox_state.SetTokenBalance(makerOrder.UserAddress, newMakerInTotal, settleBalance.Maker.InToken, statedb)
	tomox_state.SetTokenBalance(makerOrder.UserAddress, newMakerOutTotal, settleBalance.Maker.OutToken, statedb)

	// add balance for relayers
	//log.Debug("ApplyTomoXMatchedTransaction settle fee for relayers",
	//	"takerRelayerOwner", takerExOwner,
	//	"takerFeeToken", quoteToken, "takerFee", settleBalanceResult[takerAddr][tomox.Fee].(*big.Int),
	//	"makerRelayerOwner", makerExOwner,
	//	"makerFeeToken", quoteToken, "makerFee", settleBalanceResult[makerAddr][tomox.Fee].(*big.Int))
	// takerFee
	tomox_state.SetTokenBalance(takerExOwner, newTakerFee, makerOrder.QuoteToken, statedb)
	tomox_state.SetTokenBalance(makerExOwner, newMakerFee, makerOrder.QuoteToken, statedb)
	return nil
}
