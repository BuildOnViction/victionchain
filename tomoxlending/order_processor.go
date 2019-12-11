package tomoxlending

import (
	"fmt"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/consensus"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/tomoxlending/lendingstate"
	"math/big"
)

func (l *Lending) CommitOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tomoXstatedb *lendingstate.TomoXStateDB, orderBook common.Hash, order *lendingstate.LendingItem) ([]map[string]string, []*lendingstate.LendingItem, error) {
	tomoxSnap := tomoXstatedb.Snapshot()
	dbSnap := statedb.Snapshot()
	trades, rejects, err := l.ApplyOrder(coinbase, chain, statedb, tomoXstatedb, orderBook, order)
	if err != nil {
		tomoXstatedb.RevertToSnapshot(tomoxSnap)
		statedb.RevertToSnapshot(dbSnap)
		return nil, nil, err
	}
	return trades, rejects, err
}

func (l *Lending) ApplyOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tomoXstatedb *lendingstate.TomoXStateDB, orderBook common.Hash, order *lendingstate.LendingItem) ([]map[string]string, []*lendingstate.LendingItem, error) {
	var (
		rejects []*lendingstate.LendingItem
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
	if order.Status == lendingstate.LendingStatusCancelled {
		err, reject := l.ProcessCancelOrder(tomoXstatedb, statedb, chain, coinbase, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
		if reject {
			rejects = append(rejects, order)
		}
		log.Debug("Exchange add user nonce:", "address", order.UserAddress, "status", order.Status, "nonce", nonce+1)
		tomoXstatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
		return trades, rejects, nil
	}
	if order.Type != lendingstate.Market {
		if order.Interest.Sign() == 0 || common.BigToHash(order.Interest).Big().Cmp(order.Interest) != 0 {
			log.Debug("Reject order Interest invalid", "Interest", order.Interest)
			rejects = append(rejects, order)
			tomoXstatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
			return trades, rejects, nil
		}
	}
	if order.Quantity.Sign() == 0 || common.BigToHash(order.Quantity).Big().Cmp(order.Quantity) != 0 {
		log.Debug("Reject order quantity invalid", "quantity", order.Quantity)
		rejects = append(rejects, order)
		tomoXstatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
		return trades, rejects, nil
	}
	orderType := order.Type
	// if we do not use auto-increment orderid, we must set Interest slot to avoid conflict
	if orderType == lendingstate.Market {
		log.Debug("Process maket order", "side", order.Side, "quantity", order.Quantity, "Interest", order.Interest)
		trades, rejects, err = l.processMarketOrder(coinbase, chain, statedb, tomoXstatedb, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
	} else {
		log.Debug("Process limit order", "side", order.Side, "quantity", order.Quantity, "Interest", order.Interest)
		trades, rejects, err = l.processLimitOrder(coinbase, chain, statedb, tomoXstatedb, orderBook, order)
		if err != nil {
			return nil, nil, err
		}
	}

	log.Debug("Exchange add user nonce:", "address", order.UserAddress, "status", order.Status, "nonce", nonce+1)
	tomoXstatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
	return trades, rejects, nil
}

// processMarketOrder : process the market order
func (l *Lending) processMarketOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tomoXstatedb *lendingstate.TomoXStateDB, orderBook common.Hash, order *lendingstate.LendingItem) ([]map[string]string, []*lendingstate.LendingItem, error) {
	var (
		trades     []map[string]string
		newTrades  []map[string]string
		rejects    []*lendingstate.LendingItem
		newRejects []*lendingstate.LendingItem
		err        error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	// speedup the comparison, do not assign because it is pointer
	zero := lendingstate.Zero
	if side == lendingstate.Bid {
		bestInterest, volume := tomoXstatedb.GetBestAskInterest(orderBook)
		log.Debug("processMarketOrder ", "side", side, "bestInterest", bestInterest, "quantityToTrade", quantityToTrade, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && bestInterest.Cmp(zero) > 0 {
			quantityToTrade, newTrades, newRejects, err = l.processOrderList(coinbase, chain, statedb, tomoXstatedb, lendingstate.Ask, orderBook, bestInterest, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			bestInterest, volume = tomoXstatedb.GetBestAskInterest(orderBook)
			log.Debug("processMarketOrder ", "side", side, "bestInterest", bestInterest, "quantityToTrade", quantityToTrade, "volume", volume)
		}
	} else {
		bestInterest, volume := tomoXstatedb.GetBestBidInterest(orderBook)
		log.Debug("processMarketOrder ", "side", side, "bestInterest", bestInterest, "quantityToTrade", quantityToTrade, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && bestInterest.Cmp(zero) > 0 {
			quantityToTrade, newTrades, newRejects, err = l.processOrderList(coinbase, chain, statedb, tomoXstatedb, lendingstate.Bid, orderBook, bestInterest, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			bestInterest, volume = tomoXstatedb.GetBestBidInterest(orderBook)
			log.Debug("processMarketOrder ", "side", side, "bestInterest", bestInterest, "quantityToTrade", quantityToTrade, "volume", volume)
		}
	}
	return trades, newRejects, nil
}

// processLimitOrder : process the limit order, can change the quote
// If not care for performance, we should make a copy of quote to prevent further reference problem
func (l *Lending) processLimitOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tomoXstatedb *lendingstate.TomoXStateDB, orderBook common.Hash, order *lendingstate.LendingItem) ([]map[string]string, []*lendingstate.LendingItem, error) {
	var (
		trades     []map[string]string
		newTrades  []map[string]string
		rejects    []*lendingstate.LendingItem
		newRejects []*lendingstate.LendingItem
		err        error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	Interest := order.Interest

	// speedup the comparison, do not assign because it is pointer
	zero := lendingstate.Zero

	if side == lendingstate.Bid {
		minInterest, volume := tomoXstatedb.GetBestAskInterest(orderBook)
		log.Debug("processLimitOrder ", "side", side, "minInterest", minInterest, "orderInterest", Interest, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && Interest.Cmp(minInterest) >= 0 && minInterest.Cmp(zero) > 0 {
			log.Debug("Min Interest in asks tree", "Interest", minInterest.String())
			quantityToTrade, newTrades, newRejects, err = l.processOrderList(coinbase, chain, statedb, tomoXstatedb, lendingstate.Ask, orderBook, minInterest, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			log.Debug("New trade found", "newTrades", newTrades, "quantityToTrade", quantityToTrade)
			minInterest, volume = tomoXstatedb.GetBestAskInterest(orderBook)
			log.Debug("processLimitOrder ", "side", side, "minInterest", minInterest, "orderInterest", Interest, "volume", volume)
		}
	} else {
		maxInterest, volume := tomoXstatedb.GetBestBidInterest(orderBook)
		log.Debug("processLimitOrder ", "side", side, "maxInterest", maxInterest, "orderInterest", Interest, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && Interest.Cmp(maxInterest) <= 0 && maxInterest.Cmp(zero) > 0 {
			log.Debug("Max Interest in bids tree", "Interest", maxInterest.String())
			quantityToTrade, newTrades, newRejects, err = l.processOrderList(coinbase, chain, statedb, tomoXstatedb, lendingstate.Bid, orderBook, maxInterest, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			log.Debug("New trade found", "newTrades", newTrades, "quantityToTrade", quantityToTrade)
			maxInterest, volume = tomoXstatedb.GetBestBidInterest(orderBook)
			log.Debug("processLimitOrder ", "side", side, "maxInterest", maxInterest, "orderInterest", Interest, "volume", volume)
		}
	}
	if quantityToTrade.Cmp(zero) > 0 {
		orderId := tomoXstatedb.GetNonce(orderBook)
		order.LendingId = orderId + 1
		order.Quantity = quantityToTrade
		tomoXstatedb.SetNonce(orderBook, orderId+1)
		orderIdHash := common.BigToHash(new(big.Int).SetUint64(order.LendingId))
		tomoXstatedb.InsertLendingItem(orderBook, orderIdHash, *order)
		log.Debug("After matching, order (unmatched part) is now added to tree", "side", order.Side, "order", order)
	}
	return trades, rejects, nil
}

// processOrderList : process the order list
func (l *Lending) processOrderList(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, tomoXstatedb *lendingstate.TomoXStateDB, side string, orderBook common.Hash, Interest *big.Int, quantityStillToTrade *big.Int, order *lendingstate.LendingItem) (*big.Int, []map[string]string, []*lendingstate.LendingItem, error) {
	quantityToTrade := lendingstate.CloneBigInt(quantityStillToTrade)
	log.Debug("Process matching between order and orderlist", "quantityToTrade", quantityToTrade)
	var (
		trades []map[string]string

		rejects []*lendingstate.LendingItem
	)
	for quantityToTrade.Sign() > 0 {
		orderId, amount, _ := tomoXstatedb.GetBestOrderIdAndAmount(orderBook, Interest, side)
		var oldestOrder lendingstate.LendingItem
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
			maxTradedQuantity = lendingstate.CloneBigInt(quantityToTrade)
		} else {
			maxTradedQuantity = lendingstate.CloneBigInt(amount)
		}
		var quoteInterest *big.Int
		if oldestOrder.CollateralToken.String() != common.TomoNativeAddress {
			quoteInterest = tomoXstatedb.GetInterest(lendingstate.GetOrderBookHash(oldestOrder.CollateralToken, common.HexToAddress(common.TomoNativeAddress)))
			log.Debug("TryGet quoteInterest CollateralToken/TOMO", "quoteInterest", quoteInterest)
			if (quoteInterest == nil || quoteInterest.Sign() == 0) && oldestOrder.LendingToken.String() != common.TomoNativeAddress {
				inverseInterest := tomoXstatedb.GetInterest(lendingstate.GetOrderBookHash(common.HexToAddress(common.TomoNativeAddress), oldestOrder.CollateralToken))
				CollateralTokenDecimal, err := l.GetTokenDecimal(chain, statedb, coinbase, oldestOrder.CollateralToken)
				if err != nil || CollateralTokenDecimal.Sign() == 0 {
					return nil, nil, nil, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", oldestOrder.CollateralToken.String(), err)
				}
				log.Debug("TryGet inverseInterest TOMO/CollateralToken", "inverseInterest", inverseInterest)
				if inverseInterest != nil && inverseInterest.Sign() > 0 {
					quoteInterest = new(big.Int).Div(lendingstate.BaseInterest, inverseInterest)
					quoteInterest = new(big.Int).Mul(quoteInterest, CollateralTokenDecimal)
					log.Debug("TryGet quoteInterest after get inverseInterest TOMO/CollateralToken", "quoteInterest", quoteInterest, "CollateralTokenDecimal", CollateralTokenDecimal)
				}
			}
		}
		tradedQuantity, rejectMaker, err := l.getTradeQuantity(quoteInterest, coinbase, chain, statedb, order, &oldestOrder, maxTradedQuantity)
		if err != nil && err == lendingstate.ErrQuantityTradeTooSmall {
			if tradedQuantity.Cmp(maxTradedQuantity) == 0 {
				if quantityToTrade.Cmp(amount) == 0 { // reject Taker & maker
					rejects = append(rejects, order)
					quantityToTrade = lendingstate.Zero
					rejects = append(rejects, &oldestOrder)
					err = tomoXstatedb.CancelOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					break
				} else if quantityToTrade.Cmp(amount) < 0 { // reject Taker
					rejects = append(rejects, order)
					quantityToTrade = lendingstate.Zero
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
					quantityToTrade = lendingstate.Zero
					break
				}
			}
		} else if err != nil {
			return nil, nil, nil, err
		}
		if tradedQuantity.Sign() == 0 && !rejectMaker {
			log.Debug("Reject order Taker ", "tradedQuantity", tradedQuantity, "rejectMaker", rejectMaker)
			rejects = append(rejects, order)
			quantityToTrade = lendingstate.Zero
			break
		}
		if tradedQuantity.Sign() > 0 {
			quantityToTrade = lendingstate.Sub(quantityToTrade, tradedQuantity)
			tomoXstatedb.SubAmountLendingItem(orderBook, orderId, Interest, tradedQuantity, side)
			tomoXstatedb.SetInterest(orderBook, Interest)
			log.Debug("Update quantity for orderId", "orderId", orderId.Hex())
			log.Debug("TRADE", "orderBook", orderBook, "Taker Interest", Interest, "maker Interest", order.Interest, "Amount", tradedQuantity, "orderId", orderId, "side", side)

			tradeRecord := make(map[string]string)
			//tradeRecord[lendingstate2.TradeTakerOrderHash] = order.Hash.Hex()
			//tradeRecord[lendingstate2.TradeMakerOrderHash] = oldestOrder.Hash.Hex()
			//tradeRecord[lendingstate2.TradeTimestamp] = strconv.FormatInt(time.Now().Unix(), 10)
			//tradeRecord[lendingstate2.TradeQuantity] = tradedQuantity.String()
			//tradeRecord[lendingstate2.TradeMakerExchange] = oldestOrder.Relayer.String()
			//tradeRecord[lendingstate2.TradeMaker] = oldestOrder.UserAddress.String()
			//tradeRecord[lendingstate2.TradeLendingToken] = oldestOrder.LendingToken.String()
			//tradeRecord[lendingstate2.TradeCollateralToken] = oldestOrder.CollateralToken.String()
			//// maker Interest is actual Interest
			//// Taker Interest is offer Interest
			//// tradedInterest is always actual Interest
			//tradeRecord[lendingstate2.TradeInterest] = oldestOrder.Interest.String()
			//tradeRecord[lendingstate2.MakerOrderType] = oldestOrder.Type
			trades = append(trades, tradeRecord)
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

func (l *Lending) getTradeQuantity(quoteInterest *big.Int, coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, takerOrder *lendingstate.LendingItem, makerOrder *lendingstate.LendingItem, quantityToTrade *big.Int) (*big.Int, bool, error) {
	LendingTokenDecimal, err := l.GetTokenDecimal(chain, statedb, coinbase, makerOrder.LendingToken)
	if err != nil || LendingTokenDecimal.Sign() == 0 {
		return lendingstate.Zero, false, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", makerOrder.LendingToken.String(), err)
	}
	CollateralTokenDecimal, err := l.GetTokenDecimal(chain, statedb, coinbase, makerOrder.CollateralToken)
	if err != nil || CollateralTokenDecimal.Sign() == 0 {
		return lendingstate.Zero, false, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", makerOrder.CollateralToken.String(), err)
	}
	if makerOrder.CollateralToken.String() == common.TomoNativeAddress {
		quoteInterest = CollateralTokenDecimal
	}
	if takerOrder.Relayer.String() == makerOrder.Relayer.String() {
		if err := lendingstate.CheckRelayerFee(takerOrder.Relayer, new(big.Int).Mul(common.RelayerFee, big.NewInt(2)), statedb); err != nil {
			log.Debug("Reject order Taker Exchnage = Maker Exchange , relayer not enough fee ", "err", err)
			return lendingstate.Zero, false, nil
		}
	} else {
		if err := lendingstate.CheckRelayerFee(takerOrder.Relayer, common.RelayerFee, statedb); err != nil {
			log.Debug("Reject order Taker , relayer not enough fee ", "err", err)
			return lendingstate.Zero, false, nil
		}
		if err := lendingstate.CheckRelayerFee(makerOrder.Relayer, common.RelayerFee, statedb); err != nil {
			log.Debug("Reject order maker , relayer not enough fee ", "err", err)
			return lendingstate.Zero, true, nil
		}
	}
	takerFeeRate := lendingstate.GetExRelayerFee(takerOrder.Relayer, statedb)
	makerFeeRate := lendingstate.GetExRelayerFee(makerOrder.Relayer, statedb)
	var takerBalance, makerBalance *big.Int
	switch takerOrder.Side {
	case lendingstate.Bid:
		takerBalance = lendingstate.GetTokenBalance(takerOrder.UserAddress, makerOrder.CollateralToken, statedb)
		makerBalance = lendingstate.GetTokenBalance(makerOrder.UserAddress, makerOrder.LendingToken, statedb)
	case lendingstate.Ask:
		takerBalance = lendingstate.GetTokenBalance(takerOrder.UserAddress, makerOrder.LendingToken, statedb)
		makerBalance = lendingstate.GetTokenBalance(makerOrder.UserAddress, makerOrder.CollateralToken, statedb)
	default:
		takerBalance = big.NewInt(0)
		makerBalance = big.NewInt(0)
	}
	quantity, rejectMaker := GetTradeQuantity(takerOrder.Side, takerFeeRate, takerBalance, makerOrder.Interest, makerFeeRate, makerBalance, LendingTokenDecimal, quantityToTrade)
	log.Debug("GetTradeQuantity", "side", takerOrder.Side, "takerBalance", takerBalance, "makerBalance", makerBalance, "LendingToken", makerOrder.LendingToken, "CollateralToken", makerOrder.CollateralToken, "quantity", quantity, "rejectMaker", rejectMaker, "quoteInterest", quoteInterest)
	if quantity.Sign() > 0 {
		// Apply Match Order
		settleBalanceResult, err := lendingstate.GetSettleBalance(quoteInterest, takerOrder.Side, takerFeeRate, makerOrder.LendingToken, makerOrder.CollateralToken, makerOrder.Interest, makerFeeRate, LendingTokenDecimal, CollateralTokenDecimal, quantity)
		log.Debug("GetSettleBalance", "settleBalanceResult", settleBalanceResult, "err", err)
		if err == nil {
			err = DoSettleBalance(coinbase, takerOrder, makerOrder, settleBalanceResult, statedb)
		}
		return quantity, rejectMaker, err
	}
	return quantity, rejectMaker, nil
}

func GetTradeQuantity(takerSide string, takerFeeRate *big.Int, takerBalance *big.Int, makerInterest *big.Int, makerFeeRate *big.Int, makerBalance *big.Int, LendingTokenDecimal *big.Int, quantityToTrade *big.Int) (*big.Int, bool) {
	if takerSide == lendingstate.Bid {
		// maker InQuantity CollateralTokenQuantity=(quantityToTrade*maker.Interest/LendingTokenDecimal)
		CollateralTokenQuantity := new(big.Int).Mul(quantityToTrade, makerInterest)
		CollateralTokenQuantity = CollateralTokenQuantity.Div(CollateralTokenQuantity, LendingTokenDecimal)
		// Fee
		// charge on the token he/she has before the trade, in this case: CollateralToken
		// charge on the token he/she has before the trade, in this case: LendingToken
		// takerFee = CollateralTokenQuantity*takerFeeRate/baseFee=(quantityToTrade*maker.Interest/LendingTokenDecimal) * makerFeeRate/baseFee
		takerFee := big.NewInt(0).Mul(CollateralTokenQuantity, takerFeeRate)
		takerFee = big.NewInt(0).Div(takerFee, common.TomoXBaseFee)
		//takerOutTotal= CollateralTokenQuantity + takerFee =  quantityToTrade*maker.Interest/LendingTokenDecimal + quantityToTrade*maker.Interest/LendingTokenDecimal * takerFeeRate/baseFee
		// = quantityToTrade *  maker.Interest/LendingTokenDecimal ( 1 +  takerFeeRate/baseFee)
		// = quantityToTrade * maker.Interest * (baseFee + takerFeeRate ) / ( LendingTokenDecimal * baseFee)
		takerOutTotal := new(big.Int).Add(CollateralTokenQuantity, takerFee)
		makerOutTotal := quantityToTrade
		if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			return quantityToTrade, false
		} else if takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			newQuantityTrade := new(big.Int).Mul(takerBalance, LendingTokenDecimal)
			newQuantityTrade = newQuantityTrade.Mul(newQuantityTrade, common.TomoXBaseFee)
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, new(big.Int).Add(common.TomoXBaseFee, takerFeeRate))
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, makerInterest)
			if newQuantityTrade.Sign() == 0 {
				log.Debug("Reject order Taker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "takerOutTotal", takerOutTotal)
			}
			return newQuantityTrade, false
		} else if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) < 0 {
			log.Debug("Reject order maker , not enough balance ", "makerBalance", makerBalance, " makerOutTotal", makerOutTotal)
			return makerBalance, true
		} else {
			// takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) < 0
			newQuantityTrade := new(big.Int).Mul(takerBalance, LendingTokenDecimal)
			newQuantityTrade = newQuantityTrade.Mul(newQuantityTrade, common.TomoXBaseFee)
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, new(big.Int).Add(common.TomoXBaseFee, takerFeeRate))
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, makerInterest)
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
		// CollateralTokenQuantity = quantityToTrade * makerInterest / LendingTokenDecimal
		CollateralTokenQuantity := new(big.Int).Mul(quantityToTrade, makerInterest)
		CollateralTokenQuantity = CollateralTokenQuantity.Div(CollateralTokenQuantity, LendingTokenDecimal)
		// maker InQuantity

		// Fee
		// charge on the token he/she has before the trade, in this case: LendingToken
		// makerFee = CollateralTokenQuantity * makerFeeRate / baseFee = quantityToTrade * makerInterest / LendingTokenDecimal * makerFeeRate / baseFee
		// charge on the token he/she has before the trade, in this case: CollateralToken
		makerFee := new(big.Int).Mul(CollateralTokenQuantity, makerFeeRate)
		makerFee = makerFee.Div(makerFee, common.TomoXBaseFee)

		takerOutTotal := quantityToTrade
		// makerOutTotal = CollateralTokenQuantity + makerFee  = quantityToTrade * makerInterest / LendingTokenDecimal + quantityToTrade * makerInterest / LendingTokenDecimal * makerFeeRate / baseFee
		// =  quantityToTrade * makerInterest / LendingTokenDecimal * (1+makerFeeRate / baseFee)
		// = quantityToTrade  * makerInterest * (baseFee + makerFeeRate) / ( LendingTokenDecimal * baseFee )
		makerOutTotal := new(big.Int).Add(CollateralTokenQuantity, makerFee)
		if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			return quantityToTrade, false
		} else if takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			if takerBalance.Sign() == 0 {
				log.Debug("Reject order Taker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "takerOutTotal", takerOutTotal)
			}
			return takerBalance, false
		} else if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) < 0 {
			newQuantityTrade := new(big.Int).Mul(makerBalance, LendingTokenDecimal)
			newQuantityTrade = newQuantityTrade.Mul(newQuantityTrade, common.TomoXBaseFee)
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, new(big.Int).Add(common.TomoXBaseFee, makerFeeRate))
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, makerInterest)
			log.Debug("Reject order maker , not enough balance ", "makerBalance", makerBalance, " makerOutTotal", makerOutTotal)
			return newQuantityTrade, true
		} else {
			// takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) < 0
			newQuantityTrade := new(big.Int).Mul(makerBalance, LendingTokenDecimal)
			newQuantityTrade = newQuantityTrade.Mul(newQuantityTrade, common.TomoXBaseFee)
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, new(big.Int).Add(common.TomoXBaseFee, makerFeeRate))
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, makerInterest)
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

func DoSettleBalance(coinbase common.Address, takerOrder, makerOrder *lendingstate.LendingItem, settleBalance *lendingstate.SettleBalance, statedb *state.StateDB) error {
	takerExOwner := lendingstate.GetRelayerOwner(takerOrder.Relayer, statedb)
	makerExOwner := lendingstate.GetRelayerOwner(makerOrder.Relayer, statedb)
	matchingFee := big.NewInt(0)
	// masternodes charges fee of both 2 relayers. If maker and Taker are on same relayer, that relayer is charged fee twice
	matchingFee = matchingFee.Add(matchingFee, common.RelayerFee)
	matchingFee = matchingFee.Add(matchingFee, common.RelayerFee)

	if common.EmptyHash(takerExOwner.Hash()) || common.EmptyHash(makerExOwner.Hash()) {
		return fmt.Errorf("Echange owner empty , Taker: %v , maker : %v ", takerExOwner, makerExOwner)
	}
	mapBalances := map[common.Address]map[common.Address]*big.Int{}
	//Checking balance
	newTakerInTotal, err := lendingstate.CheckAddTokenBalance(takerOrder.UserAddress, settleBalance.Taker.InTotal, settleBalance.Taker.InToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Taker.InToken] == nil {
		mapBalances[settleBalance.Taker.InToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Taker.InToken][takerOrder.UserAddress] = newTakerInTotal
	newTakerOutTotal, err := lendingstate.CheckSubTokenBalance(takerOrder.UserAddress, settleBalance.Taker.OutTotal, settleBalance.Taker.OutToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Taker.OutToken] == nil {
		mapBalances[settleBalance.Taker.OutToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Taker.OutToken][takerOrder.UserAddress] = newTakerOutTotal
	newMakerInTotal, err := lendingstate.CheckAddTokenBalance(makerOrder.UserAddress, settleBalance.Maker.InTotal, settleBalance.Maker.InToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Maker.InToken] == nil {
		mapBalances[settleBalance.Maker.InToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Maker.InToken][makerOrder.UserAddress] = newMakerInTotal
	newMakerOutTotal, err := lendingstate.CheckSubTokenBalance(makerOrder.UserAddress, settleBalance.Maker.OutTotal, settleBalance.Maker.OutToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Maker.OutToken] == nil {
		mapBalances[settleBalance.Maker.OutToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Maker.OutToken][makerOrder.UserAddress] = newMakerOutTotal
	newTakerFee, err := lendingstate.CheckAddTokenBalance(takerExOwner, settleBalance.Taker.Fee, makerOrder.CollateralToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[makerOrder.CollateralToken] == nil {
		mapBalances[makerOrder.CollateralToken] = map[common.Address]*big.Int{}
	}
	mapBalances[makerOrder.CollateralToken][takerExOwner] = newTakerFee
	newMakerFee, err := lendingstate.CheckAddTokenBalance(makerExOwner, settleBalance.Maker.Fee, makerOrder.CollateralToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	mapBalances[makerOrder.CollateralToken][makerExOwner] = newMakerFee

	mapRelayerFee := map[common.Address]*big.Int{}
	newRelayerTakerFee, err := lendingstate.CheckSubRelayerFee(takerOrder.Relayer, common.RelayerFee, statedb, mapRelayerFee)
	if err != nil {
		return err
	}
	mapRelayerFee[takerOrder.Relayer] = newRelayerTakerFee
	newRelayerMakerFee, err := lendingstate.CheckSubRelayerFee(makerOrder.Relayer, common.RelayerFee, statedb, mapRelayerFee)
	if err != nil {
		return err
	}
	mapRelayerFee[makerOrder.Relayer] = newRelayerMakerFee
	lendingstate.SetSubRelayerFee(takerOrder.Relayer, newRelayerTakerFee, common.RelayerFee, statedb)
	lendingstate.SetSubRelayerFee(makerOrder.Relayer, newRelayerMakerFee, common.RelayerFee, statedb)

	masternodeOwner := statedb.GetOwner(coinbase)
	statedb.AddBalance(masternodeOwner, matchingFee)

	lendingstate.SetTokenBalance(takerOrder.UserAddress, newTakerInTotal, settleBalance.Taker.InToken, statedb)
	lendingstate.SetTokenBalance(takerOrder.UserAddress, newTakerOutTotal, settleBalance.Taker.OutToken, statedb)

	lendingstate.SetTokenBalance(makerOrder.UserAddress, newMakerInTotal, settleBalance.Maker.InToken, statedb)
	lendingstate.SetTokenBalance(makerOrder.UserAddress, newMakerOutTotal, settleBalance.Maker.OutToken, statedb)

	// add balance for relayers
	//log.Debug("ApplyTomoXMatchedTransaction settle fee for relayers",
	//	"takerRelayerOwner", takerExOwner,
	//	"takerFeeToken", CollateralToken, "takerFee", settleBalanceResult[takerAddr][tomox.Fee].(*big.Int),
	//	"makerRelayerOwner", makerExOwner,
	//	"makerFeeToken", CollateralToken, "makerFee", settleBalanceResult[makerAddr][tomox.Fee].(*big.Int))
	// takerFee
	lendingstate.SetTokenBalance(takerExOwner, newTakerFee, makerOrder.CollateralToken, statedb)
	lendingstate.SetTokenBalance(makerExOwner, newMakerFee, makerOrder.CollateralToken, statedb)
	return nil
}

func (l *Lending) ProcessCancelOrder(tomoXstatedb *lendingstate.TomoXStateDB, statedb *state.StateDB, chain consensus.ChainContext, coinbase common.Address, orderBook common.Hash, order *lendingstate.LendingItem) (error, bool) {
	if err := lendingstate.CheckRelayerFee(order.Relayer, common.RelayerCancelFee, statedb); err != nil {
		log.Debug("Relayer not enough fee when cancel order", "err", err)
		return nil, true
	}
	LendingTokenDecimal, err := l.GetTokenDecimal(chain, statedb, coinbase, order.LendingToken)
	if err != nil || LendingTokenDecimal.Sign() == 0 {
		log.Debug("Fail to get tokenDecimal ", "Token", order.LendingToken.String(), "err", err)
		return err, false
	}
	var tokenBalance *big.Int
	switch order.Side {
	case lendingstate.Ask:
		tokenBalance = lendingstate.GetTokenBalance(order.UserAddress, order.LendingToken, statedb)
	case lendingstate.Bid:
		tokenBalance = lendingstate.GetTokenBalance(order.UserAddress, order.CollateralToken, statedb)
	default:
		log.Debug("Not found order side", "Side", order.Side)
		return nil, true
	}
	log.Debug("ProcessCancelOrder", "LendingToken", order.LendingToken, "CollateralToken", order.CollateralToken, "makerInterest", order.Interest, "LendingTokenDecimal", LendingTokenDecimal, "quantity", order.Quantity)
	feeRate := lendingstate.GetExRelayerFee(order.Relayer, statedb)
	tokenCancelFee := getCancelFee(LendingTokenDecimal, feeRate, order)
	if tokenBalance.Cmp(tokenCancelFee) < 0 {
		log.Debug("User not enough balance when cancel order", "Side", order.Side, "Interest", order.Interest, "Quantity", order.Quantity, "balance", tokenBalance, "fee", tokenCancelFee)
		return nil, true
	}
	err = tomoXstatedb.CancelOrder(orderBook, order)
	if err != nil {
		log.Debug("Error when cancel order", "order", order)
		return err, false
	}
	lendingstate.SubRelayerFee(order.Relayer, common.RelayerCancelFee, statedb)
	switch order.Side {
	case lendingstate.Ask:
		lendingstate.SubTokenBalance(order.UserAddress, tokenCancelFee, order.LendingToken, statedb)
	case lendingstate.Bid:
		lendingstate.SubTokenBalance(order.UserAddress, tokenCancelFee, order.CollateralToken, statedb)
	default:
	}
	masternodeOwner := statedb.GetOwner(coinbase)
	statedb.AddBalance(masternodeOwner, common.RelayerCancelFee)
	return nil, false
}

func getCancelFee(LendingTokenDecimal *big.Int, feeRate *big.Int, order *lendingstate.LendingItem) *big.Int {
	cancelFee := big.NewInt(0)
	if order.Side == lendingstate.Ask {
		// SELL 1 BTC => TOMO ,,
		// order.Quantity =1 && fee rate =2
		// ==> cancel fee = 2/10000
		LendingTokenQuantity := new(big.Int).Mul(order.Quantity, LendingTokenDecimal)
		cancelFee = new(big.Int).Mul(LendingTokenQuantity, feeRate)
		cancelFee = new(big.Int).Div(cancelFee, common.TomoXBaseCancelFee)
	} else {
		// BUY 1 BTC => TOMO with Interest : 10000
		// CollateralTokenQuantity = 10000 && fee rate =2
		// => cancel fee =2
		CollateralTokenQuantity := new(big.Int).Mul(order.Quantity, order.Interest)
		CollateralTokenQuantity = CollateralTokenQuantity.Div(CollateralTokenQuantity, LendingTokenDecimal)
		// Fee
		// makerFee = CollateralTokenQuantity * feeRate / baseFee = quantityToTrade * makerInterest / LendingTokenDecimal * feeRate / baseFee
		cancelFee = new(big.Int).Mul(CollateralTokenQuantity, feeRate)
		cancelFee = new(big.Int).Div(CollateralTokenQuantity, common.TomoXBaseCancelFee)
	}
	return cancelFee
}
