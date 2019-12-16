package tomox

import (
	"github.com/tomochain/tomochain/consensus"
	"github.com/tomochain/tomochain/tomox/lending_state"
	"math/big"
	"strconv"
	"time"

	"fmt"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/log"
)

func (tomox *TomoX) CommitLendingOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, lendingStatedb *lending_state.LendingStateDB, orderBook common.Hash, order *lending_state.LendingItem) ([]map[string]string, []*lending_state.LendingItem, error) {
	tomoxSnap := lendingStatedb.Snapshot()
	dbSnap := statedb.Snapshot()
	trades, rejects, err := tomox.ApplyLendingOrder(coinbase, chain, statedb, lendingStatedb, orderBook, order)
	if err != nil {
		lendingStatedb.RevertToSnapshot(tomoxSnap)
		statedb.RevertToSnapshot(dbSnap)
		return nil, nil, err
	}
	return trades, rejects, err
}

func (tomox *TomoX) ApplyLendingOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, lendingStatedb *lending_state.LendingStateDB, orderBook common.Hash, order *lending_state.LendingItem) ([]map[string]string, []*lending_state.LendingItem, error) {
	var (
		rejects []*lending_state.LendingItem
		trades  []map[string]string
		err     error
	)
	nonce := lendingStatedb.GetNonce(order.UserAddress.Hash())
	log.Debug("ApplyLendingOrder", "addr", order.UserAddress, "statenonce", nonce, "ordernonce", order.Nonce)
	if big.NewInt(int64(nonce)).Cmp(order.Nonce) == -1 {
		return nil, nil, ErrNonceTooHigh
	} else if big.NewInt(int64(nonce)).Cmp(order.Nonce) == 1 {
		return nil, nil, ErrNonceTooLow
	}
	if order.Status == OrderStatusCancelled {
		err, reject := tomox.ProcessCancelLendingOrder(lendingStatedb, statedb, chain, coinbase, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
		if reject {
			rejects = append(rejects, order)
		}
		log.Debug("Exchange add user nonce:", "address", order.UserAddress, "status", order.Status, "nonce", nonce+1)
		lendingStatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
		return trades, rejects, nil
	}
	if order.Type != lending_state.Market {
		if order.Price.Sign() == 0 || common.BigToHash(order.Price).Big().Cmp(order.Price) != 0 {
			log.Debug("Reject order price invalid", "price", order.Price)
			rejects = append(rejects, order)
			lendingStatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
			return trades, rejects, nil
		}
	}
	if order.Quantity.Sign() == 0 || common.BigToHash(order.Quantity).Big().Cmp(order.Quantity) != 0 {
		log.Debug("Reject order quantity invalid", "quantity", order.Quantity)
		rejects = append(rejects, order)
		lendingStatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
		return trades, rejects, nil
	}
	orderType := order.Type
	// if we do not use auto-increment orderid, we must set price slot to avoid conflict
	if orderType == lending_state.Market {
		log.Debug("Process maket order", "side", order.Side, "quantity", order.Quantity, "price", order.Price)
		trades, rejects, err = tomox.processMarketLendingOrder(coinbase, chain, statedb, lendingStatedb, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
	} else {
		log.Debug("Process limit order", "side", order.Side, "quantity", order.Quantity, "price", order.Price)
		trades, rejects, err = tomox.processLimitLendingOrder(coinbase, chain, statedb, lendingStatedb, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
	}

	log.Debug("Exchange add user nonce:", "address", order.UserAddress, "status", order.Status, "nonce", nonce+1)
	lendingStatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
	return trades, rejects, nil
}

func (tomox *TomoX) processMarketLendingOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, lendingStatedb *lending_state.LendingStateDB, orderBook common.Hash, order *lending_state.LendingItem) ([]map[string]string, []*lending_state.LendingItem, error) {
	var (
		trades     []map[string]string
		newTrades  []map[string]string
		rejects    []*lending_state.LendingItem
		newRejects []*lending_state.LendingItem
		err        error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	// speedup the comparison, do not assign because it is pointer
	zero := lending_state.Zero
	if side == lending_state.BORROWING {
		bestPrice := lendingStatedb.GetBestInvestingRate(orderBook)
		log.Debug("processMarketLendingOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade)
		for quantityToTrade.Cmp(zero) > 0 && bestPrice.Cmp(zero) > 0 {
			quantityToTrade, newTrades, newRejects, err = tomox.processLendingOrderList(coinbase, chain, statedb, lendingStatedb, lending_state.INVESTING, orderBook, bestPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			bestPrice = lendingStatedb.GetBestInvestingRate(orderBook)
			log.Debug("processMarketLendingOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade)
		}
	} else {
		bestPrice := lendingStatedb.GetBestBorrowRate(orderBook)
		log.Debug("processMarketLendingOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade)
		for quantityToTrade.Cmp(zero) > 0 && bestPrice.Cmp(zero) > 0 {
			quantityToTrade, newTrades, newRejects, err = tomox.processLendingOrderList(coinbase, chain, statedb, lendingStatedb, lending_state.BORROWING, orderBook, bestPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			bestPrice = lendingStatedb.GetBestBorrowRate(orderBook)
			log.Debug("processMarketLendingOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade)
		}
	}
	return trades, newRejects, nil
}

func (tomox *TomoX) processLimitLendingOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, lendingStatedb *lending_state.LendingStateDB, orderBook common.Hash, order *lending_state.LendingItem) ([]map[string]string, []*lending_state.LendingItem, error) {
	var (
		trades     []map[string]string
		newTrades  []map[string]string
		rejects    []*lending_state.LendingItem
		newRejects []*lending_state.LendingItem
		err        error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	price := order.Price

	// speedup the comparison, do not assign because it is pointer
	zero := lending_state.Zero

	if side == lending_state.BORROWING {
		minPrice := lendingStatedb.GetBestInvestingRate(orderBook)
		log.Debug("processLimitLendingOrder ", "side", side, "minPrice", minPrice, "orderPrice", price)
		for quantityToTrade.Cmp(zero) > 0 && price.Cmp(minPrice) >= 0 && minPrice.Cmp(zero) > 0 {
			log.Debug("Min price in asks tree", "price", minPrice.String())
			quantityToTrade, newTrades, newRejects, err = tomox.processLendingOrderList(coinbase, chain, statedb, lendingStatedb, lending_state.INVESTING, orderBook, minPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			log.Debug("New trade found", "newTrades", newTrades, "quantityToTrade", quantityToTrade)
			minPrice = lendingStatedb.GetBestInvestingRate(orderBook)
			log.Debug("processLimitLendingOrder ", "side", side, "minPrice", minPrice, "orderPrice", price)
		}
	} else {
		maxPrice := lendingStatedb.GetBestBorrowRate(orderBook)
		log.Debug("processLimitLendingOrder ", "side", side, "maxPrice", maxPrice, "orderPrice", price)
		for quantityToTrade.Cmp(zero) > 0 && price.Cmp(maxPrice) <= 0 && maxPrice.Cmp(zero) > 0 {
			log.Debug("Max price in bids tree", "price", maxPrice.String())
			quantityToTrade, newTrades, newRejects, err = tomox.processLendingOrderList(coinbase, chain, statedb, lendingStatedb, lending_state.BORROWING, orderBook, maxPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			log.Debug("New trade found", "newTrades", newTrades, "quantityToTrade", quantityToTrade)
			maxPrice = lendingStatedb.GetBestBorrowRate(orderBook)
			log.Debug("processLimitLendingOrder ", "side", side, "maxPrice", maxPrice, "orderPrice", price)
		}
	}
	if quantityToTrade.Cmp(zero) > 0 {
		orderId := lendingStatedb.GetNonce(orderBook)
		order.OrderID = orderId + 1
		order.Quantity = quantityToTrade
		lendingStatedb.SetNonce(orderBook, orderId+1)
		orderIdHash := common.BigToHash(new(big.Int).SetUint64(order.OrderID))
		lendingStatedb.InsertLendingItem(orderBook, orderIdHash, *order)
		log.Debug("After matching, order (unmatched part) is now added to tree", "side", order.Side, "order", order)
	}
	return trades, rejects, nil
}

// processTradingOrderList : process the order list
func (tomox *TomoX) processLendingOrderList(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, lendingStatedb *lending_state.LendingStateDB, side string, orderBook common.Hash, price *big.Int, quantityStillToTrade *big.Int, order *lending_state.LendingItem) (*big.Int, []map[string]string, []*lending_state.LendingItem, error) {
	quantityToTrade := lending_state.CloneBigInt(quantityStillToTrade)
	log.Debug("Process matching between order and orderlist", "quantityToTrade", quantityToTrade)
	var (
		trades []map[string]string

		rejects []*lending_state.LendingItem
	)
	for quantityToTrade.Sign() > 0 {
		orderId, amount, _ := lendingStatedb.GetBestLendingIdAndAmount(orderBook, price, side)
		var oldestOrder lending_state.LendingItem
		if amount.Sign() > 0 {
			oldestOrder = lendingStatedb.GetLendingOrder(orderBook, orderId)
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
			maxTradedQuantity = lending_state.CloneBigInt(quantityToTrade)
		} else {
			maxTradedQuantity = lending_state.CloneBigInt(amount)
		}
		//var quotePrice *big.Int
		//if oldestOrder.QuoteToken.String() != common.TomoNativeAddress {
		//	quotePrice = lendingStatedb.GetPrice(lending_state.GetOrderBookHash(oldestOrder.QuoteToken, common.HexToAddress(common.TomoNativeAddress)))
		//	log.Debug("TryGet quotePrice QuoteToken/TOMO", "quotePrice", quotePrice)
		//	if (quotePrice == nil || quotePrice.Sign() == 0) && oldestOrder.BaseToken.String() != common.TomoNativeAddress {
		//		inversePrice := lendingStatedb.GetPrice(lending_state.GetOrderBookHash(common.HexToAddress(common.TomoNativeAddress), oldestOrder.QuoteToken))
		//		quoteTokenDecimal, err := tomox.GetTokenDecimal(chain, statedb, coinbase, oldestOrder.QuoteToken)
		//		if err != nil || quoteTokenDecimal.Sign() == 0 {
		//			return nil, nil, nil, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", oldestOrder.QuoteToken.String(), err)
		//		}
		//		log.Debug("TryGet inversePrice TOMO/QuoteToken", "inversePrice", inversePrice)
		//		if inversePrice != nil && inversePrice.Sign() > 0 {
		//			quotePrice = new(big.Int).Div(common.BasePrice, inversePrice)
		//			quotePrice = new(big.Int).Mul(quotePrice, quoteTokenDecimal)
		//			log.Debug("TryGet quotePrice after get inversePrice TOMO/QuoteToken", "quotePrice", quotePrice, "quoteTokenDecimal", quoteTokenDecimal)
		//		}
		//	}
		//}
		tradedQuantity, rejectMaker, err := tomox.getLendingQuantity(nil, coinbase, chain, statedb, order, &oldestOrder, maxTradedQuantity)
		if err != nil && err == lending_state.ErrQuantityTradeTooSmall {
			if tradedQuantity.Cmp(maxTradedQuantity) == 0 {
				if quantityToTrade.Cmp(amount) == 0 { // reject Taker & maker
					rejects = append(rejects, order)
					quantityToTrade = lending_state.Zero
					rejects = append(rejects, &oldestOrder)
					err = lendingStatedb.CancelLendingOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					break
				} else if quantityToTrade.Cmp(amount) < 0 { // reject Taker
					rejects = append(rejects, order)
					quantityToTrade = lending_state.Zero
					break
				} else { // reject maker
					rejects = append(rejects, &oldestOrder)
					err = lendingStatedb.CancelLendingOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					continue
				}
			} else {
				if rejectMaker { // reject maker
					rejects = append(rejects, &oldestOrder)
					err = lendingStatedb.CancelLendingOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					continue
				} else { // reject Taker
					rejects = append(rejects, order)
					quantityToTrade = lending_state.Zero
					break
				}
			}
		} else if err != nil {
			return nil, nil, nil, err
		}
		if tradedQuantity.Sign() == 0 && !rejectMaker {
			log.Debug("Reject order Taker ", "tradedQuantity", tradedQuantity, "rejectMaker", rejectMaker)
			rejects = append(rejects, order)
			quantityToTrade = lending_state.Zero
			break
		}
		if tradedQuantity.Sign() > 0 {
			quantityToTrade = lending_state.Sub(quantityToTrade, tradedQuantity)
			lendingStatedb.SubAmountLendingItem(orderBook, orderId, price, tradedQuantity, side)
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
			err := lendingStatedb.CancelLendingOrder(orderBook, &oldestOrder)
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}
	return quantityToTrade, trades, rejects, nil
}

func (tomox *TomoX) getLendingQuantity(quotePrice *big.Int, coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, takerOrder *lending_state.LendingItem, makerOrder *lending_state.LendingItem, quantityToTrade *big.Int) (*big.Int, bool, error) {
	baseTokenDecimal, err := tomox.GetTokenDecimal(chain, statedb, coinbase, makerOrder.BaseToken)
	if err != nil || baseTokenDecimal.Sign() == 0 {
		return lending_state.Zero, false, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", makerOrder.BaseToken.String(), err)
	}
	quoteTokenDecimal, err := tomox.GetTokenDecimal(chain, statedb, coinbase, makerOrder.QuoteToken)
	if err != nil || quoteTokenDecimal.Sign() == 0 {
		return lending_state.Zero, false, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", makerOrder.QuoteToken.String(), err)
	}
	if makerOrder.QuoteToken.String() == common.TomoNativeAddress {
		quotePrice = quoteTokenDecimal
	}
	if takerOrder.ExchangeAddress.String() == makerOrder.ExchangeAddress.String() {
		if err := lending_state.CheckRelayerFee(takerOrder.ExchangeAddress, new(big.Int).Mul(common.RelayerFee, big.NewInt(2)), statedb); err != nil {
			log.Debug("Reject order Taker Exchnage = Maker Exchange , relayer not enough fee ", "err", err)
			return lending_state.Zero, false, nil
		}
	} else {
		if err := lending_state.CheckRelayerFee(takerOrder.ExchangeAddress, common.RelayerFee, statedb); err != nil {
			log.Debug("Reject order Taker , relayer not enough fee ", "err", err)
			return lending_state.Zero, false, nil
		}
		if err := lending_state.CheckRelayerFee(makerOrder.ExchangeAddress, common.RelayerFee, statedb); err != nil {
			log.Debug("Reject order maker , relayer not enough fee ", "err", err)
			return lending_state.Zero, true, nil
		}
	}
	takerFeeRate := lending_state.GetExRelayerFee(takerOrder.ExchangeAddress, statedb)
	makerFeeRate := lending_state.GetExRelayerFee(makerOrder.ExchangeAddress, statedb)
	var takerBalance, makerBalance *big.Int
	switch takerOrder.Side {
	case lending_state.BORROWING:
		takerBalance = lending_state.GetTokenBalance(takerOrder.UserAddress, makerOrder.QuoteToken, statedb)
		makerBalance = lending_state.GetTokenBalance(makerOrder.UserAddress, makerOrder.BaseToken, statedb)
	case lending_state.INVESTING:
		takerBalance = lending_state.GetTokenBalance(takerOrder.UserAddress, makerOrder.BaseToken, statedb)
		makerBalance = lending_state.GetTokenBalance(makerOrder.UserAddress, makerOrder.QuoteToken, statedb)
	default:
		takerBalance = big.NewInt(0)
		makerBalance = big.NewInt(0)
	}
	quantity, rejectMaker := getLendingQuantity(takerOrder.Side, takerFeeRate, takerBalance, makerOrder.Price, makerFeeRate, makerBalance, baseTokenDecimal, quantityToTrade)
	log.Debug("getLendingQuantity", "side", takerOrder.Side, "takerBalance", takerBalance, "makerBalance", makerBalance, "BaseToken", makerOrder.BaseToken, "QuoteToken", makerOrder.QuoteToken, "quantity", quantity, "rejectMaker", rejectMaker, "quotePrice", quotePrice)
	if quantity.Sign() > 0 {
		// Apply Match Order
		settleBalanceResult, err := lending_state.GetSettleBalance(quotePrice, takerOrder.Side, takerFeeRate, makerOrder.BaseToken, makerOrder.QuoteToken, makerOrder.Price, makerFeeRate, baseTokenDecimal, quoteTokenDecimal, quantity)
		log.Debug("GetSettleBalance", "settleBalanceResult", settleBalanceResult, "err", err)
		if err == nil {
			err = doSettleLendingBalance(coinbase, takerOrder, makerOrder, settleBalanceResult, statedb)
		}
		return quantity, rejectMaker, err
	}
	return quantity, rejectMaker, nil
}

func getLendingQuantity(takerSide string, takerFeeRate *big.Int, takerBalance *big.Int, makerPrice *big.Int, makerFeeRate *big.Int, makerBalance *big.Int, baseTokenDecimal *big.Int, quantityToTrade *big.Int) (*big.Int, bool) {
	if takerSide == lending_state.BORROWING {
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

func doSettleLendingBalance(coinbase common.Address, takerOrder, makerOrder *lending_state.LendingItem, settleBalance *lending_state.SettleBalance, statedb *state.StateDB) error {
	takerExOwner := lending_state.GetRelayerOwner(takerOrder.ExchangeAddress, statedb)
	makerExOwner := lending_state.GetRelayerOwner(makerOrder.ExchangeAddress, statedb)
	matchingFee := big.NewInt(0)
	// masternodes charges fee of both 2 relayers. If maker and Taker are on same relayer, that relayer is charged fee twice
	matchingFee = matchingFee.Add(matchingFee, common.RelayerFee)
	matchingFee = matchingFee.Add(matchingFee, common.RelayerFee)

	if common.EmptyHash(takerExOwner.Hash()) || common.EmptyHash(makerExOwner.Hash()) {
		return fmt.Errorf("Echange owner empty , Taker: %v , maker : %v ", takerExOwner, makerExOwner)
	}
	mapBalances := map[common.Address]map[common.Address]*big.Int{}
	//Checking balance
	newTakerInTotal, err := lending_state.CheckAddTokenBalance(takerOrder.UserAddress, settleBalance.Taker.InTotal, settleBalance.Taker.InToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Taker.InToken] == nil {
		mapBalances[settleBalance.Taker.InToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Taker.InToken][takerOrder.UserAddress] = newTakerInTotal
	newTakerOutTotal, err := lending_state.CheckSubTokenBalance(takerOrder.UserAddress, settleBalance.Taker.OutTotal, settleBalance.Taker.OutToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Taker.OutToken] == nil {
		mapBalances[settleBalance.Taker.OutToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Taker.OutToken][takerOrder.UserAddress] = newTakerOutTotal
	newMakerInTotal, err := lending_state.CheckAddTokenBalance(makerOrder.UserAddress, settleBalance.Maker.InTotal, settleBalance.Maker.InToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Maker.InToken] == nil {
		mapBalances[settleBalance.Maker.InToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Maker.InToken][makerOrder.UserAddress] = newMakerInTotal
	newMakerOutTotal, err := lending_state.CheckSubTokenBalance(makerOrder.UserAddress, settleBalance.Maker.OutTotal, settleBalance.Maker.OutToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Maker.OutToken] == nil {
		mapBalances[settleBalance.Maker.OutToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Maker.OutToken][makerOrder.UserAddress] = newMakerOutTotal
	newTakerFee, err := lending_state.CheckAddTokenBalance(takerExOwner, settleBalance.Taker.Fee, makerOrder.QuoteToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[makerOrder.QuoteToken] == nil {
		mapBalances[makerOrder.QuoteToken] = map[common.Address]*big.Int{}
	}
	mapBalances[makerOrder.QuoteToken][takerExOwner] = newTakerFee
	newMakerFee, err := lending_state.CheckAddTokenBalance(makerExOwner, settleBalance.Maker.Fee, makerOrder.QuoteToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	mapBalances[makerOrder.QuoteToken][makerExOwner] = newMakerFee

	mapRelayerFee := map[common.Address]*big.Int{}
	newRelayerTakerFee, err := lending_state.CheckSubRelayerFee(takerOrder.ExchangeAddress, common.RelayerFee, statedb, mapRelayerFee)
	if err != nil {
		return err
	}
	mapRelayerFee[takerOrder.ExchangeAddress] = newRelayerTakerFee
	newRelayerMakerFee, err := lending_state.CheckSubRelayerFee(makerOrder.ExchangeAddress, common.RelayerFee, statedb, mapRelayerFee)
	if err != nil {
		return err
	}
	mapRelayerFee[makerOrder.ExchangeAddress] = newRelayerMakerFee
	lending_state.SetSubRelayerFee(takerOrder.ExchangeAddress, newRelayerTakerFee, common.RelayerFee, statedb)
	lending_state.SetSubRelayerFee(makerOrder.ExchangeAddress, newRelayerMakerFee, common.RelayerFee, statedb)

	masternodeOwner := statedb.GetOwner(coinbase)
	statedb.AddBalance(masternodeOwner, matchingFee)

	lending_state.SetTokenBalance(takerOrder.UserAddress, newTakerInTotal, settleBalance.Taker.InToken, statedb)
	lending_state.SetTokenBalance(takerOrder.UserAddress, newTakerOutTotal, settleBalance.Taker.OutToken, statedb)

	lending_state.SetTokenBalance(makerOrder.UserAddress, newMakerInTotal, settleBalance.Maker.InToken, statedb)
	lending_state.SetTokenBalance(makerOrder.UserAddress, newMakerOutTotal, settleBalance.Maker.OutToken, statedb)

	// add balance for relayers
	//log.Debug("ApplyTomoXMatchedTransaction settle fee for relayers",
	//	"takerRelayerOwner", takerExOwner,
	//	"takerFeeToken", quoteToken, "takerFee", settleBalanceResult[takerAddr][tomox.Fee].(*big.Int),
	//	"makerRelayerOwner", makerExOwner,
	//	"makerFeeToken", quoteToken, "makerFee", settleBalanceResult[makerAddr][tomox.Fee].(*big.Int))
	// takerFee
	lending_state.SetTokenBalance(takerExOwner, newTakerFee, makerOrder.QuoteToken, statedb)
	lending_state.SetTokenBalance(makerExOwner, newMakerFee, makerOrder.QuoteToken, statedb)
	return nil
}

func (tomox *TomoX) ProcessCancelLendingOrder(lendingStatedb *lending_state.LendingStateDB, statedb *state.StateDB, chain consensus.ChainContext, coinbase common.Address, orderBook common.Hash, order *lending_state.LendingItem) (error, bool) {
	if err := lending_state.CheckRelayerFee(order.ExchangeAddress, common.RelayerCancelFee, statedb); err != nil {
		log.Debug("Relayer not enough fee when cancel order", "err", err)
		return nil, true
	}
	baseTokenDecimal, err := tomox.GetTokenDecimal(chain, statedb, coinbase, order.BaseToken)
	if err != nil || baseTokenDecimal.Sign() == 0 {
		log.Debug("Fail to get tokenDecimal ", "Token", order.BaseToken.String(), "err", err)
		return err, false
	}
	var tokenBalance *big.Int
	switch order.Side {
	case lending_state.INVESTING:
		tokenBalance = lending_state.GetTokenBalance(order.UserAddress, order.BaseToken, statedb)
	case lending_state.BORROWING:
		tokenBalance = lending_state.GetTokenBalance(order.UserAddress, order.QuoteToken, statedb)
	default:
		log.Debug("Not found order side", "Side", order.Side)
		return nil, true
	}
	log.Debug("ProcessCancelLendingOrder", "baseToken", order.BaseToken, "quoteToken", order.QuoteToken, "makerPrice", order.Price, "baseTokenDecimal", baseTokenDecimal, "quantity", order.Quantity)
	feeRate := lending_state.GetExRelayerFee(order.ExchangeAddress, statedb)
	tokenCancelFee := getLendingCancelOrderFee(baseTokenDecimal, feeRate, order)
	if tokenBalance.Cmp(tokenCancelFee) < 0 {
		log.Debug("User not enough balance when cancel order", "Side", order.Side, "Price", order.Price, "Quantity", order.Quantity, "balance", tokenBalance, "fee", tokenCancelFee)
		return nil, true
	}
	err = lendingStatedb.CancelLendingOrder(orderBook, order)
	if err != nil {
		log.Debug("Error when cancel order", "order", order)
		return err, false
	}
	lending_state.SubRelayerFee(order.ExchangeAddress, common.RelayerCancelFee, statedb)
	switch order.Side {
	case lending_state.INVESTING:
		lending_state.SubTokenBalance(order.UserAddress, tokenCancelFee, order.BaseToken, statedb)
	case lending_state.BORROWING:
		lending_state.SubTokenBalance(order.UserAddress, tokenCancelFee, order.QuoteToken, statedb)
	default:
	}
	masternodeOwner := statedb.GetOwner(coinbase)
	statedb.AddBalance(masternodeOwner, common.RelayerCancelFee)
	return nil, false
}

func getLendingCancelOrderFee(baseTokenDecimal *big.Int, feeRate *big.Int, order *lending_state.LendingItem) *big.Int {
	cancelFee := big.NewInt(0)
	if order.Side == lending_state.INVESTING {
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
