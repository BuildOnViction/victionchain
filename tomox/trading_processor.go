package tomox

import (
	"github.com/tomochain/tomochain/consensus"
	"math/big"
	"strconv"
	"time"

	"fmt"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/tomox/trading_state"
)

func (tomox *TomoX) CommitTradingOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tradingStatedb *trading_state.TradingStateDB, orderBook common.Hash, order *trading_state.OrderItem) ([]map[string]string, []*trading_state.OrderItem, error) {
	tomoxSnap := tradingStatedb.Snapshot()
	dbSnap := statedb.Snapshot()
	trades, rejects, err := tomox.ApplyTradingOrder(coinbase, chain, statedb, tradingStatedb, orderBook, order)
	if err != nil {
		tradingStatedb.RevertToSnapshot(tomoxSnap)
		statedb.RevertToSnapshot(dbSnap)
		return nil, nil, err
	}
	return trades, rejects, err
}

func (tomox *TomoX) ApplyTradingOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tradingStatedb *trading_state.TradingStateDB, orderBook common.Hash, order *trading_state.OrderItem) ([]map[string]string, []*trading_state.OrderItem, error) {
	var (
		rejects []*trading_state.OrderItem
		trades  []map[string]string
		err     error
	)
	nonce := tradingStatedb.GetNonce(order.UserAddress.Hash())
	log.Debug("ApplyTradingOrder", "addr", order.UserAddress, "statenonce", nonce, "ordernonce", order.Nonce)
	if big.NewInt(int64(nonce)).Cmp(order.Nonce) == -1 {
		return nil, nil, ErrNonceTooHigh
	} else if big.NewInt(int64(nonce)).Cmp(order.Nonce) == 1 {
		return nil, nil, ErrNonceTooLow
	}
	if order.Status == OrderStatusCancelled {
		err, reject := tomox.ProcessCancelTradingOrder(tradingStatedb, statedb, chain, coinbase, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
		if reject {
			rejects = append(rejects, order)
		}
		log.Debug("Exchange add user nonce:", "address", order.UserAddress, "status", order.Status, "nonce", nonce+1)
		tradingStatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
		return trades, rejects, nil
	}
	if order.Type != trading_state.Market {
		if order.Price.Sign() == 0 || common.BigToHash(order.Price).Big().Cmp(order.Price) != 0 {
			log.Debug("Reject order price invalid", "price", order.Price)
			rejects = append(rejects, order)
			tradingStatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
			return trades, rejects, nil
		}
	}
	if order.Quantity.Sign() == 0 || common.BigToHash(order.Quantity).Big().Cmp(order.Quantity) != 0 {
		log.Debug("Reject order quantity invalid", "quantity", order.Quantity)
		rejects = append(rejects, order)
		tradingStatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
		return trades, rejects, nil
	}
	orderType := order.Type
	// if we do not use auto-increment orderid, we must set price slot to avoid conflict
	if orderType == trading_state.Market {
		log.Debug("Process maket order", "side", order.Side, "quantity", order.Quantity, "price", order.Price)
		trades, rejects, err = tomox.processMarketTradingOrder(coinbase, chain, statedb, tradingStatedb, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
	} else {
		log.Debug("Process limit order", "side", order.Side, "quantity", order.Quantity, "price", order.Price)
		trades, rejects, err = tomox.processLimitTradingOrder(coinbase, chain, statedb, tradingStatedb, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
	}

	log.Debug("Exchange add user nonce:", "address", order.UserAddress, "status", order.Status, "nonce", nonce+1)
	tradingStatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
	return trades, rejects, nil
}

func (tomox *TomoX) processMarketTradingOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tradingStatedb *trading_state.TradingStateDB, orderBook common.Hash, order *trading_state.OrderItem) ([]map[string]string, []*trading_state.OrderItem, error) {
	var (
		trades     []map[string]string
		newTrades  []map[string]string
		rejects    []*trading_state.OrderItem
		newRejects []*trading_state.OrderItem
		err        error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	// speedup the comparison, do not assign because it is pointer
	zero := trading_state.Zero
	if side == trading_state.Bid {
		bestPrice, _ := tradingStatedb.GetBestAskPrice(orderBook)
		log.Debug("processMarketTradingOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade)
		for quantityToTrade.Cmp(zero) > 0 && bestPrice.Cmp(zero) > 0 {
			quantityToTrade, newTrades, newRejects, err = tomox.processTradingOrderList(coinbase, chain, statedb, tradingStatedb, trading_state.Ask, orderBook, bestPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			bestPrice, _ = tradingStatedb.GetBestAskPrice(orderBook)
			log.Debug("processMarketTradingOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade)
		}
	} else {
		bestPrice, _ := tradingStatedb.GetBestBidPrice(orderBook)
		log.Debug("processMarketTradingOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade)
		for quantityToTrade.Cmp(zero) > 0 && bestPrice.Cmp(zero) > 0 {
			quantityToTrade, newTrades, newRejects, err = tomox.processTradingOrderList(coinbase, chain, statedb, tradingStatedb, trading_state.Bid, orderBook, bestPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			bestPrice, _ = tradingStatedb.GetBestBidPrice(orderBook)
			log.Debug("processMarketTradingOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade)
		}
	}
	return trades, newRejects, nil
}

func (tomox *TomoX) processLimitTradingOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tradingStatedb *trading_state.TradingStateDB, orderBook common.Hash, order *trading_state.OrderItem) ([]map[string]string, []*trading_state.OrderItem, error) {
	var (
		trades     []map[string]string
		newTrades  []map[string]string
		rejects    []*trading_state.OrderItem
		newRejects []*trading_state.OrderItem
		err        error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	price := order.Price

	// speedup the comparison, do not assign because it is pointer
	zero := trading_state.Zero

	if side == trading_state.Bid {
		minPrice, _ := tradingStatedb.GetBestAskPrice(orderBook)
		log.Debug("processLimitTradingOrder ", "side", side, "minPrice", minPrice, "orderPrice", price)
		for quantityToTrade.Cmp(zero) > 0 && price.Cmp(minPrice) >= 0 && minPrice.Cmp(zero) > 0 {
			log.Debug("Min price in asks tree", "price", minPrice.String())
			quantityToTrade, newTrades, newRejects, err = tomox.processTradingOrderList(coinbase, chain, statedb, tradingStatedb, trading_state.Ask, orderBook, minPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			log.Debug("New trade found", "newTrades", newTrades, "quantityToTrade", quantityToTrade)
			minPrice, _ = tradingStatedb.GetBestAskPrice(orderBook)
			log.Debug("processLimitTradingOrder ", "side", side, "minPrice", minPrice, "orderPrice", price)
		}
	} else {
		maxPrice, _ := tradingStatedb.GetBestBidPrice(orderBook)
		log.Debug("processLimitTradingOrder ", "side", side, "maxPrice", maxPrice, "orderPrice", price)
		for quantityToTrade.Cmp(zero) > 0 && price.Cmp(maxPrice) <= 0 && maxPrice.Cmp(zero) > 0 {
			log.Debug("Max price in bids tree", "price", maxPrice.String())
			quantityToTrade, newTrades, newRejects, err = tomox.processTradingOrderList(coinbase, chain, statedb, tradingStatedb, trading_state.Bid, orderBook, maxPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			log.Debug("New trade found", "newTrades", newTrades, "quantityToTrade", quantityToTrade)
			maxPrice, _ = tradingStatedb.GetBestBidPrice(orderBook)
			log.Debug("processLimitTradingOrder ", "side", side, "maxPrice", maxPrice, "orderPrice", price)
		}
	}
	if quantityToTrade.Cmp(zero) > 0 {
		orderId := tradingStatedb.GetNonce(orderBook)
		order.OrderID = orderId + 1
		order.Quantity = quantityToTrade
		tradingStatedb.SetNonce(orderBook, orderId+1)
		orderIdHash := common.BigToHash(new(big.Int).SetUint64(order.OrderID))
		tradingStatedb.InsertOrderItem(orderBook, orderIdHash, *order)
		log.Debug("After matching, order (unmatched part) is now added to tree", "side", order.Side, "order", order)
	}
	return trades, rejects, nil
}

func (tomox *TomoX) processTradingOrderList(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tradingStatedb *trading_state.TradingStateDB, side string, orderBook common.Hash, price *big.Int, quantityStillToTrade *big.Int, order *trading_state.OrderItem) (*big.Int, []map[string]string, []*trading_state.OrderItem, error) {
	quantityToTrade := trading_state.CloneBigInt(quantityStillToTrade)
	log.Debug("Process matching between order and orderlist", "quantityToTrade", quantityToTrade)
	var (
		trades []map[string]string

		rejects []*trading_state.OrderItem
	)
	for quantityToTrade.Sign() > 0 {
		orderId, amount, _ := tradingStatedb.GetBestOrderIdAndAmount(orderBook, price, side)
		var oldestOrder trading_state.OrderItem
		if amount.Sign() > 0 {
			oldestOrder = tradingStatedb.GetOrder(orderBook, orderId)
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
			maxTradedQuantity = trading_state.CloneBigInt(quantityToTrade)
		} else {
			maxTradedQuantity = trading_state.CloneBigInt(amount)
		}
		var quotePrice *big.Int
		if oldestOrder.QuoteToken.String() != common.TomoNativeAddress {
			quotePrice = tradingStatedb.GetPrice(trading_state.GetOrderBookHash(oldestOrder.QuoteToken, common.HexToAddress(common.TomoNativeAddress)))
			log.Debug("TryGet quotePrice QuoteToken/TOMO", "quotePrice", quotePrice)
			if (quotePrice == nil || quotePrice.Sign() == 0) && oldestOrder.BaseToken.String() != common.TomoNativeAddress {
				inversePrice := tradingStatedb.GetPrice(trading_state.GetOrderBookHash(common.HexToAddress(common.TomoNativeAddress), oldestOrder.QuoteToken))
				quoteTokenDecimal, err := tomox.GetTokenDecimal(chain, statedb, coinbase, oldestOrder.QuoteToken)
				if err != nil || quoteTokenDecimal.Sign() == 0 {
					return nil, nil, nil, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", oldestOrder.QuoteToken.String(), err)
				}
				log.Debug("TryGet inversePrice TOMO/QuoteToken", "inversePrice", inversePrice)
				if inversePrice != nil && inversePrice.Sign() > 0 {
					quotePrice = new(big.Int).Div(common.BasePrice, inversePrice)
					quotePrice = new(big.Int).Mul(quotePrice, quoteTokenDecimal)
					log.Debug("TryGet quotePrice after get inversePrice TOMO/QuoteToken", "quotePrice", quotePrice, "quoteTokenDecimal", quoteTokenDecimal)
				}
			}
		}
		tradedQuantity, rejectMaker, err := tomox.getTradingQuantity(quotePrice, coinbase, chain, statedb, order, &oldestOrder, maxTradedQuantity)
		if err != nil && err == trading_state.ErrQuantityTradeTooSmall {
			if tradedQuantity.Cmp(maxTradedQuantity) == 0 {
				if quantityToTrade.Cmp(amount) == 0 { // reject Taker & maker
					rejects = append(rejects, order)
					quantityToTrade = trading_state.Zero
					rejects = append(rejects, &oldestOrder)
					err = tradingStatedb.CancelOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					break
				} else if quantityToTrade.Cmp(amount) < 0 { // reject Taker
					rejects = append(rejects, order)
					quantityToTrade = trading_state.Zero
					break
				} else { // reject maker
					rejects = append(rejects, &oldestOrder)
					err = tradingStatedb.CancelOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					continue
				}
			} else {
				if rejectMaker { // reject maker
					rejects = append(rejects, &oldestOrder)
					err = tradingStatedb.CancelOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					continue
				} else { // reject Taker
					rejects = append(rejects, order)
					quantityToTrade = trading_state.Zero
					break
				}
			}
		} else if err != nil {
			return nil, nil, nil, err
		}
		if tradedQuantity.Sign() == 0 && !rejectMaker {
			log.Debug("Reject order Taker ", "tradedQuantity", tradedQuantity, "rejectMaker", rejectMaker)
			rejects = append(rejects, order)
			quantityToTrade = trading_state.Zero
			break
		}
		if tradedQuantity.Sign() > 0 {
			quantityToTrade = trading_state.Sub(quantityToTrade, tradedQuantity)
			tradingStatedb.SubAmountOrderItem(orderBook, orderId, price, tradedQuantity, side)
			tradingStatedb.SetPrice(orderBook, price)
			log.Debug("Update quantity for orderId", "orderId", orderId.Hex())
			log.Debug("TRADE", "orderBook", orderBook, "Taker price", price, "maker price", order.Price, "Amount", tradedQuantity, "orderId", orderId, "side", side)

			tradeRecord := make(map[string]string)
			tradeRecord[TradeTakerOrderHash] = order.Hash.Hex()
			tradeRecord[TradeMakerOrderHash] = oldestOrder.Hash.Hex()
			tradeRecord[TradeTimestamp] = strconv.FormatInt(time.Now().Unix(), 10)
			tradeRecord[TradeQuantity] = tradedQuantity.String()
			tradeRecord[TradeMakerExchange] = oldestOrder.ExchangeAddress.String()
			tradeRecord[TradeMaker] = oldestOrder.UserAddress.String()
			tradeRecord[TradeBaseToken] = oldestOrder.BaseToken.String()
			tradeRecord[TradeQuoteToken] = oldestOrder.QuoteToken.String()
			// maker price is actual price
			// Taker price is offer price
			// tradedPrice is always actual price
			tradeRecord[TradePrice] = oldestOrder.Price.String()
			tradeRecord[MakerOrderType] = oldestOrder.Type
			trades = append(trades, tradeRecord)
		}
		if rejectMaker {
			rejects = append(rejects, &oldestOrder)
			err := tradingStatedb.CancelOrder(orderBook, &oldestOrder)
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}
	return quantityToTrade, trades, rejects, nil
}

func (tomox *TomoX) getTradingQuantity(quotePrice *big.Int, coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, takerOrder *trading_state.OrderItem, makerOrder *trading_state.OrderItem, quantityToTrade *big.Int) (*big.Int, bool, error) {
	baseTokenDecimal, err := tomox.GetTokenDecimal(chain, statedb, coinbase, makerOrder.BaseToken)
	if err != nil || baseTokenDecimal.Sign() == 0 {
		return trading_state.Zero, false, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", makerOrder.BaseToken.String(), err)
	}
	quoteTokenDecimal, err := tomox.GetTokenDecimal(chain, statedb, coinbase, makerOrder.QuoteToken)
	if err != nil || quoteTokenDecimal.Sign() == 0 {
		return trading_state.Zero, false, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", makerOrder.QuoteToken.String(), err)
	}
	if makerOrder.QuoteToken.String() == common.TomoNativeAddress {
		quotePrice = quoteTokenDecimal
	}
	if takerOrder.ExchangeAddress.String() == makerOrder.ExchangeAddress.String() {
		if err := trading_state.CheckRelayerFee(takerOrder.ExchangeAddress, new(big.Int).Mul(common.RelayerFee, big.NewInt(2)), statedb); err != nil {
			log.Debug("Reject order Taker Exchnage = Maker Exchange , relayer not enough fee ", "err", err)
			return trading_state.Zero, false, nil
		}
	} else {
		if err := trading_state.CheckRelayerFee(takerOrder.ExchangeAddress, common.RelayerFee, statedb); err != nil {
			log.Debug("Reject order Taker , relayer not enough fee ", "err", err)
			return trading_state.Zero, false, nil
		}
		if err := trading_state.CheckRelayerFee(makerOrder.ExchangeAddress, common.RelayerFee, statedb); err != nil {
			log.Debug("Reject order maker , relayer not enough fee ", "err", err)
			return trading_state.Zero, true, nil
		}
	}
	takerFeeRate := trading_state.GetExRelayerFee(takerOrder.ExchangeAddress, statedb)
	makerFeeRate := trading_state.GetExRelayerFee(makerOrder.ExchangeAddress, statedb)
	var takerBalance, makerBalance *big.Int
	switch takerOrder.Side {
	case trading_state.Bid:
		takerBalance = trading_state.GetTokenBalance(takerOrder.UserAddress, makerOrder.QuoteToken, statedb)
		makerBalance = trading_state.GetTokenBalance(makerOrder.UserAddress, makerOrder.BaseToken, statedb)
	case trading_state.Ask:
		takerBalance = trading_state.GetTokenBalance(takerOrder.UserAddress, makerOrder.BaseToken, statedb)
		makerBalance = trading_state.GetTokenBalance(makerOrder.UserAddress, makerOrder.QuoteToken, statedb)
	default:
		takerBalance = big.NewInt(0)
		makerBalance = big.NewInt(0)
	}
	quantity, rejectMaker := getTradingQuantity(takerOrder.Side, takerFeeRate, takerBalance, makerOrder.Price, makerFeeRate, makerBalance, baseTokenDecimal, quantityToTrade)
	log.Debug("getTradingQuantity", "side", takerOrder.Side, "takerBalance", takerBalance, "makerBalance", makerBalance, "BaseToken", makerOrder.BaseToken, "QuoteToken", makerOrder.QuoteToken, "quantity", quantity, "rejectMaker", rejectMaker, "quotePrice", quotePrice)
	if quantity.Sign() > 0 {
		// Apply Match Order
		settleBalanceResult, err := trading_state.GetSettleBalance(quotePrice, takerOrder.Side, takerFeeRate, makerOrder.BaseToken, makerOrder.QuoteToken, makerOrder.Price, makerFeeRate, baseTokenDecimal, quoteTokenDecimal, quantity)
		log.Debug("GetSettleBalance", "settleBalanceResult", settleBalanceResult, "err", err)
		if err == nil {
			err = doSettleTradingBalance(coinbase, takerOrder, makerOrder, settleBalanceResult, statedb)
		}
		return quantity, rejectMaker, err
	}
	return quantity, rejectMaker, nil
}

func getTradingQuantity(takerSide string, takerFeeRate *big.Int, takerBalance *big.Int, makerPrice *big.Int, makerFeeRate *big.Int, makerBalance *big.Int, baseTokenDecimal *big.Int, quantityToTrade *big.Int) (*big.Int, bool) {
	if takerSide == trading_state.Bid {
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

func doSettleTradingBalance(coinbase common.Address, takerOrder, makerOrder *trading_state.OrderItem, settleBalance *trading_state.SettleBalance, statedb *state.StateDB) error {
	takerExOwner := trading_state.GetRelayerOwner(takerOrder.ExchangeAddress, statedb)
	makerExOwner := trading_state.GetRelayerOwner(makerOrder.ExchangeAddress, statedb)
	matchingFee := big.NewInt(0)
	// masternodes charges fee of both 2 relayers. If maker and Taker are on same relayer, that relayer is charged fee twice
	matchingFee = matchingFee.Add(matchingFee, common.RelayerFee)
	matchingFee = matchingFee.Add(matchingFee, common.RelayerFee)

	if common.EmptyHash(takerExOwner.Hash()) || common.EmptyHash(makerExOwner.Hash()) {
		return fmt.Errorf("Echange owner empty , Taker: %v , maker : %v ", takerExOwner, makerExOwner)
	}
	mapBalances := map[common.Address]map[common.Address]*big.Int{}
	//Checking balance
	newTakerInTotal, err := trading_state.CheckAddTokenBalance(takerOrder.UserAddress, settleBalance.Taker.InTotal, settleBalance.Taker.InToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Taker.InToken] == nil {
		mapBalances[settleBalance.Taker.InToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Taker.InToken][takerOrder.UserAddress] = newTakerInTotal
	newTakerOutTotal, err := trading_state.CheckSubTokenBalance(takerOrder.UserAddress, settleBalance.Taker.OutTotal, settleBalance.Taker.OutToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Taker.OutToken] == nil {
		mapBalances[settleBalance.Taker.OutToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Taker.OutToken][takerOrder.UserAddress] = newTakerOutTotal
	newMakerInTotal, err := trading_state.CheckAddTokenBalance(makerOrder.UserAddress, settleBalance.Maker.InTotal, settleBalance.Maker.InToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Maker.InToken] == nil {
		mapBalances[settleBalance.Maker.InToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Maker.InToken][makerOrder.UserAddress] = newMakerInTotal
	newMakerOutTotal, err := trading_state.CheckSubTokenBalance(makerOrder.UserAddress, settleBalance.Maker.OutTotal, settleBalance.Maker.OutToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Maker.OutToken] == nil {
		mapBalances[settleBalance.Maker.OutToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Maker.OutToken][makerOrder.UserAddress] = newMakerOutTotal
	newTakerFee, err := trading_state.CheckAddTokenBalance(takerExOwner, settleBalance.Taker.Fee, makerOrder.QuoteToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[makerOrder.QuoteToken] == nil {
		mapBalances[makerOrder.QuoteToken] = map[common.Address]*big.Int{}
	}
	mapBalances[makerOrder.QuoteToken][takerExOwner] = newTakerFee
	newMakerFee, err := trading_state.CheckAddTokenBalance(makerExOwner, settleBalance.Maker.Fee, makerOrder.QuoteToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	mapBalances[makerOrder.QuoteToken][makerExOwner] = newMakerFee

	mapRelayerFee := map[common.Address]*big.Int{}
	newRelayerTakerFee, err := trading_state.CheckSubRelayerFee(takerOrder.ExchangeAddress, common.RelayerFee, statedb, mapRelayerFee)
	if err != nil {
		return err
	}
	mapRelayerFee[takerOrder.ExchangeAddress] = newRelayerTakerFee
	newRelayerMakerFee, err := trading_state.CheckSubRelayerFee(makerOrder.ExchangeAddress, common.RelayerFee, statedb, mapRelayerFee)
	if err != nil {
		return err
	}
	mapRelayerFee[makerOrder.ExchangeAddress] = newRelayerMakerFee
	trading_state.SetSubRelayerFee(takerOrder.ExchangeAddress, newRelayerTakerFee, common.RelayerFee, statedb)
	trading_state.SetSubRelayerFee(makerOrder.ExchangeAddress, newRelayerMakerFee, common.RelayerFee, statedb)

	masternodeOwner := statedb.GetOwner(coinbase)
	statedb.AddBalance(masternodeOwner, matchingFee)

	trading_state.SetTokenBalance(takerOrder.UserAddress, newTakerInTotal, settleBalance.Taker.InToken, statedb)
	trading_state.SetTokenBalance(takerOrder.UserAddress, newTakerOutTotal, settleBalance.Taker.OutToken, statedb)

	trading_state.SetTokenBalance(makerOrder.UserAddress, newMakerInTotal, settleBalance.Maker.InToken, statedb)
	trading_state.SetTokenBalance(makerOrder.UserAddress, newMakerOutTotal, settleBalance.Maker.OutToken, statedb)

	// add balance for relayers
	//log.Debug("ApplyTomoXMatchedTransaction settle fee for relayers",
	//	"takerRelayerOwner", takerExOwner,
	//	"takerFeeToken", quoteToken, "takerFee", settleBalanceResult[takerAddr][tomox.Fee].(*big.Int),
	//	"makerRelayerOwner", makerExOwner,
	//	"makerFeeToken", quoteToken, "makerFee", settleBalanceResult[makerAddr][tomox.Fee].(*big.Int))
	// takerFee
	trading_state.SetTokenBalance(takerExOwner, newTakerFee, makerOrder.QuoteToken, statedb)
	trading_state.SetTokenBalance(makerExOwner, newMakerFee, makerOrder.QuoteToken, statedb)
	return nil
}

func (tomox *TomoX) ProcessCancelTradingOrder(tradingStatedb *trading_state.TradingStateDB, statedb *state.StateDB, chain consensus.ChainContext, coinbase common.Address, orderBook common.Hash, order *trading_state.OrderItem) (error, bool) {
	if err := trading_state.CheckRelayerFee(order.ExchangeAddress, common.RelayerCancelFee, statedb); err != nil {
		log.Debug("Relayer not enough fee when cancel order", "err", err)
		return nil, true
	}
	baseTokenDecimal, err := tomox.GetTokenDecimal(chain, statedb, coinbase, order.BaseToken)
	if err != nil || baseTokenDecimal.Sign() == 0 {
		log.Debug("Fail to get tokenDecimal ", "Token", order.BaseToken.String(), "err", err)
		return err, false
	}
	originOrder := tradingStatedb.GetOrder(orderBook, common.BigToHash(new(big.Int).SetUint64(order.OrderID)))

	var tokenBalance *big.Int
	switch originOrder.Side {
	case trading_state.Ask:
		tokenBalance = trading_state.GetTokenBalance(order.UserAddress, order.BaseToken, statedb)
	case trading_state.Bid:
		tokenBalance = trading_state.GetTokenBalance(order.UserAddress, order.QuoteToken, statedb)
	default:
		log.Debug("Not found order side", "Side", originOrder.Side)
		return nil, true
	}

	log.Debug("ProcessCancelOrder", "baseToken", order.BaseToken, "quoteToken", order.QuoteToken)
	feeRate := trading_state.GetExRelayerFee(order.ExchangeAddress, statedb)
	tokenCancelFee := getTradingCancelOrderFee(baseTokenDecimal, feeRate, order)
	if tokenBalance.Cmp(tokenCancelFee) < 0 {
		log.Debug("User not enough balance when cancel order", "Side", originOrder.Side, "balance", tokenBalance, "fee", tokenCancelFee)
		return nil, true
	}
	err = tradingStatedb.CancelOrder(orderBook, order)
	if err != nil {
		log.Debug("Error when cancel order", "order", order)
		return err, false
	}
	trading_state.SubRelayerFee(order.ExchangeAddress, common.RelayerCancelFee, statedb)
	switch originOrder.Side {
	case trading_state.Ask:
		trading_state.SubTokenBalance(order.UserAddress, tokenCancelFee, order.BaseToken, statedb)
	case trading_state.Bid:
		trading_state.SubTokenBalance(order.UserAddress, tokenCancelFee, order.QuoteToken, statedb)
	default:
	}
	masternodeOwner := statedb.GetOwner(coinbase)
	statedb.AddBalance(masternodeOwner, common.RelayerCancelFee)
	return nil, false
}

func getTradingCancelOrderFee(baseTokenDecimal *big.Int, feeRate *big.Int, order *trading_state.OrderItem) *big.Int {
	cancelFee := big.NewInt(0)
	if order.Side == trading_state.Ask {
		// SELL 1 BTC => TOMO ,,
		// order.Quantity =1 && fee rate =2
		// ==> cancel fee = 2/10000
		baseTokenQuantity := new(big.Int).Mul(order.Quantity, baseTokenDecimal)
		cancelFee = new(big.Int).Mul(baseTokenQuantity, feeRate)
		cancelFee = new(big.Int).Div(cancelFee, common.TomoXBaseCancelFee)
	} else {
		// BUY 1 BTC => TOMO with Price : 10000
		// quoteTokenQuantity = 10000 && fee rate =2
		// => cancel fee =2
		quoteTokenQuantity := new(big.Int).Mul(order.Quantity, order.Price)
		quoteTokenQuantity = quoteTokenQuantity.Div(quoteTokenQuantity, baseTokenDecimal)
		// Fee
		// makerFee = quoteTokenQuantity * feeRate / baseFee = quantityToTrade * makerPrice / baseTokenDecimal * feeRate / baseFee
		cancelFee = new(big.Int).Mul(quoteTokenQuantity, feeRate)
		cancelFee = new(big.Int).Div(quoteTokenQuantity, common.TomoXBaseCancelFee)
	}
	return cancelFee
}
