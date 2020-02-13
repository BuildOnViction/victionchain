package tomoxlending

import (
	"fmt"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/consensus"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/tomox/tradingstate"
	"github.com/tomochain/tomochain/tomoxlending/lendingstate"
	"math/big"
)

func (l *Lending) CommitOrder(createdBlockTime uint64, coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, lendingStateDB *lendingstate.LendingStateDB, tradingStateDb *tradingstate.TradingStateDB, lendingOrderBook common.Hash, order *lendingstate.LendingItem) ([]*lendingstate.LendingTrade, []*lendingstate.LendingItem, error) {
	lendingSnap := lendingStateDB.Snapshot()
	tradingSnap := tradingStateDb.Snapshot()
	dbSnap := statedb.Snapshot()
	trades, rejects, err := l.ApplyOrder(createdBlockTime, coinbase, chain, statedb, lendingStateDB, tradingStateDb, lendingOrderBook, order)
	if err != nil {
		lendingStateDB.RevertToSnapshot(lendingSnap)
		tradingStateDb.RevertToSnapshot(tradingSnap)
		statedb.RevertToSnapshot(dbSnap)
		return nil, nil, err
	}
	return trades, rejects, err
}

func (l *Lending) ApplyOrder(createdBlockTime uint64, coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, lendingStateDB *lendingstate.LendingStateDB, tradingStateDb *tradingstate.TradingStateDB, lendingOrderBook common.Hash, order *lendingstate.LendingItem) ([]*lendingstate.LendingTrade, []*lendingstate.LendingItem, error) {
	var (
		rejects []*lendingstate.LendingItem
		trades  []*lendingstate.LendingTrade
		err     error
	)
	nonce := lendingStateDB.GetNonce(order.UserAddress.Hash())
	log.Debug("ApplyOrder", "addr", order.UserAddress, "statenonce", nonce, "ordernonce", order.Nonce)
	if big.NewInt(int64(nonce)).Cmp(order.Nonce) == -1 {
		return nil, nil, ErrNonceTooHigh
	} else if big.NewInt(int64(nonce)).Cmp(order.Nonce) == 1 {
		return nil, nil, ErrNonceTooLow
	}
	switch order.Status {
	case lendingstate.LendingStatusCancelled:
		err, reject := l.ProcessCancelOrder(lendingStateDB, statedb, tradingStateDb, chain, coinbase, lendingOrderBook, order)
		if err != nil {
			return nil, nil, err
		}
		if reject {
			rejects = append(rejects, order)
		}
		log.Debug("Exchange add user nonce:", "address", order.UserAddress, "status", order.Status, "nonce", nonce+1)
		lendingStateDB.SetNonce(order.UserAddress.Hash(), nonce+1)
		return trades, rejects, nil
	case lendingstate.Deposit:
		err, reject, newLendingTrade := l.ProcessDeposit(lendingStateDB, statedb, tradingStateDb, order)
		if err != nil {
			return nil, nil, err
		}
		if reject {
			rejects = append(rejects, order)
		}
		log.Debug("Exchange add user nonce:", "address", order.UserAddress, "status", order.Status, "nonce", nonce+1)
		lendingStateDB.SetNonce(order.UserAddress.Hash(), nonce+1)
		trades = append(trades, newLendingTrade)
		return trades, rejects, nil
	case lendingstate.Payment:
		lendingTradeId := order.LendingTradeId
		_, err := l.ProcessPayment(createdBlockTime, lendingStateDB, statedb, tradingStateDb, lendingOrderBook, lendingTradeId)
		if err != nil {
			log.Debug("Can not process payment", "err", err)
			rejects = append(rejects, order)
		}
		log.Debug("Exchange add user nonce:", "address", order.UserAddress, "status", order.Status, "nonce", nonce+1)
		lendingStateDB.SetNonce(order.UserAddress.Hash(), nonce+1)
		trades = append(trades, &lendingstate.LendingTrade{
			Status: lendingstate.TradeStatusClosed,
		})
		return trades, rejects, nil
	default:
	}
	if order.Type != lendingstate.Market {
		if order.Interest.Sign() == 0 || common.BigToHash(order.Interest).Big().Cmp(order.Interest) != 0 {
			log.Debug("Reject order Interest invalid", "Interest", order.Interest)
			rejects = append(rejects, order)
			lendingStateDB.SetNonce(order.UserAddress.Hash(), nonce+1)
			return trades, rejects, nil
		}
	}
	if order.Quantity.Sign() == 0 || common.BigToHash(order.Quantity).Big().Cmp(order.Quantity) != 0 {
		log.Debug("Reject order quantity invalid", "quantity", order.Quantity)
		rejects = append(rejects, order)
		lendingStateDB.SetNonce(order.UserAddress.Hash(), nonce+1)
		return trades, rejects, nil
	}
	orderType := order.Type
	// if we do not use auto-increment orderid, we must set Interest slot to avoid conflict
	if orderType == lendingstate.Market {
		log.Debug("Process maket order", "side", order.Side, "quantity", order.Quantity, "Interest", order.Interest)
		trades, rejects, err = l.processMarketOrder(createdBlockTime, coinbase, chain, statedb, lendingStateDB, tradingStateDb, lendingOrderBook, order)
		if err != nil {
			return nil, nil, err
		}
	} else {
		log.Debug("Process limit order", "side", order.Side, "quantity", order.Quantity, "Interest", order.Interest)
		trades, rejects, err = l.processLimitOrder(createdBlockTime, coinbase, chain, statedb, lendingStateDB, tradingStateDb, lendingOrderBook, order)
		if err != nil {
			return nil, nil, err
		}
	}

	log.Debug("Exchange add user nonce:", "address", order.UserAddress, "status", order.Status, "nonce", nonce+1)
	lendingStateDB.SetNonce(order.UserAddress.Hash(), nonce+1)
	return trades, rejects, nil
}

// processMarketOrder : process the market order
func (l *Lending) processMarketOrder(createdBlockTime uint64, coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, lendingStateDB *lendingstate.LendingStateDB, tradingStateDb *tradingstate.TradingStateDB, lendingOrderBook common.Hash, order *lendingstate.LendingItem) ([]*lendingstate.LendingTrade, []*lendingstate.LendingItem, error) {
	var (
		trades     []*lendingstate.LendingTrade
		newTrades  []*lendingstate.LendingTrade
		rejects    []*lendingstate.LendingItem
		newRejects []*lendingstate.LendingItem
		err        error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	// speedup the comparison, do not assign because it is pointer
	zero := lendingstate.Zero
	if side == lendingstate.Borrowing {
		bestInterest, volume := lendingStateDB.GetBestInvestingRate(lendingOrderBook)
		log.Debug("processMarketOrder ", "side", side, "bestInterest", bestInterest, "quantityToTrade", quantityToTrade, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && bestInterest.Cmp(zero) > 0 {
			quantityToTrade, newTrades, newRejects, err = l.processOrderList(createdBlockTime, coinbase, chain, statedb, lendingStateDB, tradingStateDb, lendingstate.Investing, lendingOrderBook, bestInterest, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			bestInterest, volume = lendingStateDB.GetBestInvestingRate(lendingOrderBook)
			log.Debug("processMarketOrder ", "side", side, "bestInterest", bestInterest, "quantityToTrade", quantityToTrade, "volume", volume)
		}
	} else {
		bestInterest, volume := lendingStateDB.GetBestBorrowRate(lendingOrderBook)
		log.Debug("processMarketOrder ", "side", side, "bestInterest", bestInterest, "quantityToTrade", quantityToTrade, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && bestInterest.Cmp(zero) > 0 {
			quantityToTrade, newTrades, newRejects, err = l.processOrderList(createdBlockTime, coinbase, chain, statedb, lendingStateDB, tradingStateDb, lendingstate.Borrowing, lendingOrderBook, bestInterest, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			bestInterest, volume = lendingStateDB.GetBestBorrowRate(lendingOrderBook)
			log.Debug("processMarketOrder ", "side", side, "bestInterest", bestInterest, "quantityToTrade", quantityToTrade, "volume", volume)
		}
	}
	return trades, newRejects, nil
}

// processLimitOrder : process the limit order, can change the quote
// If not care for performance, we should make a copy of quote to prevent further reference problem
func (l *Lending) processLimitOrder(createdBlockTime uint64, coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, lendingStateDB *lendingstate.LendingStateDB, tradingStateDb *tradingstate.TradingStateDB, lendingOrderBook common.Hash, order *lendingstate.LendingItem) ([]*lendingstate.LendingTrade, []*lendingstate.LendingItem, error) {
	var (
		trades     []*lendingstate.LendingTrade
		newTrades  []*lendingstate.LendingTrade
		rejects    []*lendingstate.LendingItem
		newRejects []*lendingstate.LendingItem
		err        error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	Interest := order.Interest

	// speedup the comparison, do not assign because it is pointer
	zero := lendingstate.Zero
	if side == lendingstate.Borrowing {
		minInterest, volume := lendingStateDB.GetBestInvestingRate(lendingOrderBook)
		log.Debug("processLimitOrder ", "side", side, "minInterest", minInterest, "orderInterest", Interest, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && Interest.Cmp(minInterest) >= 0 && minInterest.Cmp(zero) > 0 {
			log.Debug("Min Interest in Investing tree", "Interest", minInterest.String())
			quantityToTrade, newTrades, newRejects, err = l.processOrderList(createdBlockTime, coinbase, chain, statedb, lendingStateDB, tradingStateDb, lendingstate.Investing, lendingOrderBook, minInterest, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			log.Debug("New trade found", "newTrades", newTrades, "quantityToTrade", quantityToTrade)
			minInterest, volume = lendingStateDB.GetBestInvestingRate(lendingOrderBook)
			log.Debug("processLimitOrder ", "side", side, "minInterest", minInterest, "orderInterest", Interest, "volume", volume)
		}
	} else {
		maxInterest, volume := lendingStateDB.GetBestBorrowRate(lendingOrderBook)
		log.Debug("processLimitOrder ", "side", side, "maxInterest", maxInterest, "orderInterest", Interest, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && Interest.Cmp(maxInterest) <= 0 && maxInterest.Cmp(zero) > 0 {
			log.Debug("Max Interest in Borrowing tree", "Interest", maxInterest.String())
			quantityToTrade, newTrades, newRejects, err = l.processOrderList(createdBlockTime, coinbase, chain, statedb, lendingStateDB, tradingStateDb, lendingstate.Borrowing, lendingOrderBook, maxInterest, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			log.Debug("New trade found", "newTrades", newTrades, "quantityToTrade", quantityToTrade)
			maxInterest, volume = lendingStateDB.GetBestBorrowRate(lendingOrderBook)
			log.Debug("processLimitOrder ", "side", side, "maxInterest", maxInterest, "orderInterest", Interest, "volume", volume)
		}
	}
	if quantityToTrade.Cmp(zero) > 0 {
		orderId := lendingStateDB.GetNonce(lendingOrderBook)
		order.LendingId = orderId + 1
		order.Quantity = quantityToTrade
		lendingStateDB.SetNonce(lendingOrderBook, orderId+1)
		orderIdHash := common.BigToHash(new(big.Int).SetUint64(order.LendingId))
		lendingStateDB.InsertLendingItem(lendingOrderBook, orderIdHash, *order)
		log.Debug("After matching, order (unmatched part) is now added to tree", "side", order.Side, "order", order)
	}
	return trades, rejects, nil
}

// processOrderList : process the order list
func (l *Lending) processOrderList(createdBlockTime uint64, coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, lendingStateDB *lendingstate.LendingStateDB, tradingStateDb *tradingstate.TradingStateDB, side string, lendingOrderBook common.Hash, Interest *big.Int, quantityStillToTrade *big.Int, order *lendingstate.LendingItem) (*big.Int, []*lendingstate.LendingTrade, []*lendingstate.LendingItem, error) {
	quantityToTrade := lendingstate.CloneBigInt(quantityStillToTrade)
	log.Debug("Process matching between order and orderlist", "quantityToTrade", quantityToTrade)
	var (
		trades  []*lendingstate.LendingTrade
		rejects []*lendingstate.LendingItem
	)
	for quantityToTrade.Sign() > 0 {
		orderId, amount, err := lendingStateDB.GetBestLendingIdAndAmount(lendingOrderBook, Interest, side)
		if err != nil {
			return nil, nil, nil, err
		}
		var oldestOrder lendingstate.LendingItem
		if amount.Sign() > 0 {
			oldestOrder = lendingStateDB.GetLendingOrder(lendingOrderBook, orderId)
		}
		log.Debug("found order ", "orderId ", orderId, "side", oldestOrder.Side, "amount", amount, "side", side, "Interest", Interest)
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
		collateralToken := order.CollateralToken
		borrowFee := lendingstate.GetFee(statedb, order.Relayer)
		if order.Side == lendingstate.Investing {
			collateralToken = oldestOrder.CollateralToken
			borrowFee = lendingstate.GetFee(statedb, oldestOrder.Relayer)
		}
		collateralPrice := common.BasePrice
		depositRate, liquidationRate, _ := lendingstate.GetCollateralDetail(statedb, collateralToken)
		lendTokenTOMOPrice, collateralPrice, err := l.GetCollateralPrices(chain, statedb, tradingStateDb, order.CollateralToken, order.LendingToken)
		if err != nil {
			return nil, nil, nil, err
		}
		tradedQuantity, collateralLockedAmount, rejectMaker, settleBalanceResult, err := l.getLendQuantity(lendTokenTOMOPrice, collateralPrice, depositRate, borrowFee, coinbase, chain, statedb, order, &oldestOrder, maxTradedQuantity)
		if err != nil && err == lendingstate.ErrQuantityTradeTooSmall {
			if tradedQuantity.Cmp(maxTradedQuantity) == 0 {
				if quantityToTrade.Cmp(amount) == 0 { // reject Taker & maker
					rejects = append(rejects, order)
					quantityToTrade = lendingstate.Zero
					rejects = append(rejects, &oldestOrder)
					err = lendingStateDB.CancelLendingOrder(lendingOrderBook, &oldestOrder)
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
					err = lendingStateDB.CancelLendingOrder(lendingOrderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					continue
				}
			} else {
				if rejectMaker { // reject maker
					rejects = append(rejects, &oldestOrder)
					err = lendingStateDB.CancelLendingOrder(lendingOrderBook, &oldestOrder)
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
			lendingStateDB.SubAmountLendingItem(lendingOrderBook, orderId, Interest, tradedQuantity, side)
			log.Debug("Update quantity for orderId", "orderId", orderId.Hex())
			log.Debug("LEND", "lendingOrderBook", lendingOrderBook, "Taker Interest", Interest, "maker Interest", order.Interest, "Amount", tradedQuantity, "orderId", orderId, "side", side)
			tradingId := lendingStateDB.GetTradeNonce(lendingOrderBook) + 1
			liquidationTime := createdBlockTime + order.Term
			liquidationPrice := new(big.Int).Mul(collateralPrice, liquidationRate)
			liquidationPrice = new(big.Int).Div(liquidationPrice, depositRate)
			lendingTrade := lendingstate.LendingTrade{
				TradeId:                tradingId,
				Term:                   oldestOrder.Term,
				LendingToken:           oldestOrder.LendingToken,
				CollateralToken:        collateralToken,
				Amount:                 tradedQuantity,
				LiquidationTime:        liquidationTime,
				LiquidationPrice:       liquidationPrice,
				Interest:               oldestOrder.Interest.Uint64(),
				DepositRate:            depositRate,
				CollateralLockedAmount: collateralLockedAmount,
			}
			lendingTrade.Status = lendingstate.TradeStatusOpen
			lendingTrade.TakerOrderSide = order.Side
			lendingTrade.TakerOrderType = order.Type
			lendingTrade.MakerOrderType = oldestOrder.Type
			lendingTrade.InvestingFee = lendingstate.Zero // current design: no investing fee
			lendingTrade.CollateralPrice = collateralPrice
			
			if order.Side == lendingstate.Borrowing {
				// taker is a borrower
				lendingTrade.BorrowingOrderHash = order.Hash
				lendingTrade.InvestingOrderHash = oldestOrder.Hash
				lendingTrade.BorrowingRelayer = order.Relayer
				lendingTrade.InvestingRelayer = oldestOrder.Relayer
				lendingTrade.Borrower = order.UserAddress
				lendingTrade.Investor = oldestOrder.UserAddress

				// fee
				if settleBalanceResult != nil {
					lendingTrade.BorrowingFee = settleBalanceResult.Taker.Fee
				}
			} else if order.Side == lendingstate.Investing {
				// taker is an investor
				lendingTrade.BorrowingOrderHash = oldestOrder.Hash
				lendingTrade.InvestingOrderHash = order.Hash
				lendingTrade.BorrowingRelayer = oldestOrder.Relayer
				lendingTrade.InvestingRelayer = order.Relayer
				lendingTrade.Borrower = oldestOrder.UserAddress
				lendingTrade.Investor = order.UserAddress
				// fee
				if settleBalanceResult != nil {
					lendingTrade.BorrowingFee = settleBalanceResult.Maker.Fee
				}
			}
			lendingStateDB.InsertTradingItem(lendingOrderBook, tradingId, lendingTrade)
			lendingStateDB.InsertLiquidationTime(lendingOrderBook, new(big.Int).SetUint64(liquidationTime), tradingId)
			lendingStateDB.SetTradeNonce(lendingOrderBook, tradingId+1)
			tradingStateDb.InsertLiquidationPrice(tradingstate.GetTradingOrderBookHash(order.CollateralToken, order.LendingToken), liquidationPrice, lendingOrderBook, tradingId)

			trades = append(trades, &lendingTrade)

		}
		if rejectMaker {
			rejects = append(rejects, &oldestOrder)
			err := lendingStateDB.CancelLendingOrder(lendingOrderBook, &oldestOrder)
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}
	return quantityToTrade, trades, rejects, nil
}

func (l *Lending) getLendQuantity(
	lendTokenTOMOPrice,
	collateralPrice,
	depositRate,
	borrowFee *big.Int,
	coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, takerOrder *lendingstate.LendingItem, makerOrder *lendingstate.LendingItem, quantityToTrade *big.Int) (*big.Int, *big.Int, bool, *lendingstate.LendingSettleBalance, error) {
	if collateralPrice == nil || collateralPrice.Sign() == 0 {
		if takerOrder.Side == lendingstate.Borrowing {
			log.Debug("Reject lending order taker , can not found  collateral price ")
			return lendingstate.Zero, lendingstate.Zero, false, nil, nil
		} else {
			log.Debug("Reject lending order maker , can not found  collateral price ")
			return lendingstate.Zero, lendingstate.Zero, true, nil, nil
		}
	}
	LendingTokenDecimal, err := l.tomox.GetTokenDecimal(chain, statedb, coinbase, makerOrder.LendingToken)
	if err != nil || LendingTokenDecimal.Sign() == 0 {
		return lendingstate.Zero, lendingstate.Zero, false, nil, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", makerOrder.LendingToken.String(), err)
	}
	collateralTokenDecimal, err := l.tomox.GetTokenDecimal(chain, statedb, coinbase, makerOrder.CollateralToken)
	if err != nil || collateralTokenDecimal.Sign() == 0 {
		return lendingstate.Zero, lendingstate.Zero, false, nil, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", makerOrder.CollateralToken.String(), err)
	}
	if takerOrder.Relayer.String() == makerOrder.Relayer.String() {
		if err := lendingstate.CheckRelayerFee(takerOrder.Relayer, new(big.Int).Mul(common.RelayerLendingFee, big.NewInt(2)), statedb); err != nil {
			log.Debug("Reject order Taker Exchnage = Maker Exchange , relayer not enough fee ", "err", err)
			return lendingstate.Zero, lendingstate.Zero, false, nil, nil
		}
	} else {
		if err := lendingstate.CheckRelayerFee(takerOrder.Relayer, common.RelayerLendingFee, statedb); err != nil {
			log.Debug("Reject order Taker , relayer not enough fee ", "err", err)
			return lendingstate.Zero, lendingstate.Zero, false, nil, nil
		}
		if err := lendingstate.CheckRelayerFee(makerOrder.Relayer, common.RelayerLendingFee, statedb); err != nil {
			log.Debug("Reject order maker , relayer not enough fee ", "err", err)
			return lendingstate.Zero, lendingstate.Zero, true, nil, nil
		}
	}
	var takerBalance, makerBalance *big.Int
	var lendToken, collateralToken common.Address
	switch takerOrder.Side {
	case lendingstate.Borrowing:
		takerBalance = lendingstate.GetTokenBalance(takerOrder.UserAddress, takerOrder.CollateralToken, statedb)
		makerBalance = lendingstate.GetTokenBalance(makerOrder.UserAddress, takerOrder.LendingToken, statedb)
		lendToken = takerOrder.LendingToken
		collateralToken = takerOrder.CollateralToken
	case lendingstate.Investing:
		takerBalance = lendingstate.GetTokenBalance(takerOrder.UserAddress, makerOrder.LendingToken, statedb)
		makerBalance = lendingstate.GetTokenBalance(makerOrder.UserAddress, makerOrder.CollateralToken, statedb)
		lendToken = makerOrder.LendingToken
		collateralToken = makerOrder.CollateralToken
	default:
		takerBalance = big.NewInt(0)
		makerBalance = big.NewInt(0)
	}
	quantity, rejectMaker := GetLendQuantity(takerOrder.Side, collateralTokenDecimal, depositRate, collateralPrice, takerBalance, makerBalance, quantityToTrade)
	log.Debug("GetLendQuantity", "side", takerOrder.Side, "takerBalance", takerBalance, "makerBalance", makerBalance, "LendingToken", makerOrder.LendingToken, "CollateralToken", makerOrder.CollateralToken, "quantity", quantity, "rejectMaker", rejectMaker)
	if quantity.Sign() > 0 {
		// Apply Match Order
		settleBalanceResult, err := lendingstate.GetSettleBalance(takerOrder.Side, lendTokenTOMOPrice, collateralPrice, depositRate, borrowFee, lendToken, collateralToken, LendingTokenDecimal, collateralTokenDecimal, quantity)
		log.Debug("GetSettleBalance", "settleBalanceResult", settleBalanceResult, "err", err)
		if err == nil {
			err = DoSettleBalance(coinbase, takerOrder, makerOrder, settleBalanceResult, statedb)
		}
		return quantity, settleBalanceResult.CollateralLockedAmount, rejectMaker, settleBalanceResult, err
	}
	return quantity, lendingstate.Zero, rejectMaker, nil, nil
}

func GetLendQuantity(takerSide string, collateralTokenDecimal *big.Int, depositRate *big.Int, collateralPrice *big.Int, takerBalance *big.Int, makerBalance *big.Int, quantityToLend *big.Int) (*big.Int, bool) {
	if takerSide == lendingstate.Borrowing {
		// taker = Borrower : takerOutTotal = CollateralLockedAmount = quantityToLend * collateral Token Decimal/ CollateralPrice  * deposit rate
		takerOutTotal := new(big.Int).Mul(quantityToLend, collateralTokenDecimal)
		takerOutTotal = new(big.Int).Mul(takerOutTotal, depositRate)
		takerOutTotal = new(big.Int).Div(takerOutTotal, collateralPrice)
		// Investor : makerOutTotal = quantityToLend
		makerOutTotal := quantityToLend
		if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			return quantityToLend, false
		} else if takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			newQuantityLend := new(big.Int).Mul(takerBalance, collateralPrice)
			newQuantityLend = new(big.Int).Div(newQuantityLend, depositRate)
			newQuantityLend = new(big.Int).Div(newQuantityLend, collateralTokenDecimal)
			if newQuantityLend.Sign() == 0 {
				log.Debug("Reject lending order Taker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "takerOutTotal", takerOutTotal)
			}
			return newQuantityLend, false
		} else if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) < 0 {
			log.Debug("Reject lending order maker , not enough balance ", "makerBalance", makerBalance, " makerOutTotal", makerOutTotal)
			return makerBalance, true
		} else {
			// takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) < 0
			newQuantityLend := new(big.Int).Mul(takerBalance, collateralPrice)
			newQuantityLend = new(big.Int).Div(newQuantityLend, depositRate)
			newQuantityLend = new(big.Int).Div(newQuantityLend, collateralTokenDecimal)
			if newQuantityLend.Cmp(makerBalance) <= 0 {
				if newQuantityLend.Sign() == 0 {
					log.Debug("Reject lending order Taker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "makerBalance", makerBalance, " newQuantityLend ", newQuantityLend)
				}
				return newQuantityLend, false
			}
			log.Debug("Reject lending order maker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "makerBalance", makerBalance, " newQuantityLend ", newQuantityLend)
			return makerBalance, true
		}
	} else {
		// maker =  Borrower : makerOutTotal = CollateralLockedAmount = quantityToLend * collateral Token Decimal / CollateralPrice  * deposit rate
		makerOutTotal := new(big.Int).Mul(quantityToLend, collateralTokenDecimal)
		makerOutTotal = new(big.Int).Mul(makerOutTotal, depositRate)
		makerOutTotal = new(big.Int).Div(makerOutTotal, collateralPrice)
		// Investor : makerOutTotal = quantityToLend
		takerOutTotal := quantityToLend
		if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			return quantityToLend, false
		} else if takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			if takerBalance.Sign() == 0 {
				log.Debug("Reject lending order Taker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "takerOutTotal", takerOutTotal)
			}
			return takerBalance, false
		} else if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) < 0 {
			newQuantityLend := new(big.Int).Mul(makerBalance, collateralPrice)
			newQuantityLend = new(big.Int).Div(newQuantityLend, depositRate)
			newQuantityLend = new(big.Int).Div(newQuantityLend, collateralTokenDecimal)
			log.Debug("Reject lending order maker , not enough balance ", "makerBalance", makerBalance, " makerOutTotal", makerOutTotal)
			return newQuantityLend, true
		} else {
			// takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) < 0
			newQuantityLend := new(big.Int).Mul(makerBalance, collateralPrice)
			newQuantityLend = new(big.Int).Div(newQuantityLend, depositRate)
			newQuantityLend = new(big.Int).Div(newQuantityLend, collateralTokenDecimal)
			if newQuantityLend.Cmp(takerBalance) <= 0 {
				log.Debug("Reject lending order maker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "makerBalance", makerBalance, " newQuantityLend ", newQuantityLend)
				return newQuantityLend, true
			}
			if takerBalance.Sign() == 0 {
				log.Debug("Reject lending order Taker , not enough balance ", "takerSide", takerSide, "takerBalance", takerBalance, "makerBalance", makerBalance, " newQuantityLend ", newQuantityLend)
			}
			return takerBalance, false
		}
	}
}

func DoSettleBalance(coinbase common.Address, takerOrder, makerOrder *lendingstate.LendingItem, settleBalance *lendingstate.LendingSettleBalance, statedb *state.StateDB) error {
	takerExOwner := lendingstate.GetRelayerOwner(takerOrder.Relayer, statedb)
	makerExOwner := lendingstate.GetRelayerOwner(makerOrder.Relayer, statedb)
	matchingFee := big.NewInt(0)
	// masternodes only charge borrower relayer fee
	matchingFee = new(big.Int).Add(matchingFee, common.RelayerLendingFee)

	if common.EmptyHash(takerExOwner.Hash()) || common.EmptyHash(makerExOwner.Hash()) {
		return fmt.Errorf("Echange owner empty , Taker: %v , maker : %v ", takerExOwner, makerExOwner)
	}
	mapBalances := map[common.Address]map[common.Address]*big.Int{}
	//Checking balance
	if takerOrder.Side == lendingstate.Borrowing {
		relayerFee, err := lendingstate.CheckSubRelayerFee(takerOrder.Relayer, common.RelayerLendingFee, statedb, map[common.Address]*big.Int{})
		if err != nil {
			return err
		}
		lendingstate.SetSubRelayerFee(takerOrder.Relayer, relayerFee, common.RelayerLendingFee, statedb)
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
		newMakerOutTotal, err := lendingstate.CheckSubTokenBalance(makerOrder.UserAddress, settleBalance.Maker.OutTotal, settleBalance.Maker.OutToken, statedb, mapBalances)
		if err != nil {
			return err
		}
		if mapBalances[settleBalance.Maker.OutToken] == nil {
			mapBalances[settleBalance.Maker.OutToken] = map[common.Address]*big.Int{}
		}
		mapBalances[settleBalance.Maker.OutToken][makerOrder.UserAddress] = newMakerOutTotal
		newTakerFee, err := lendingstate.CheckAddTokenBalance(takerExOwner, settleBalance.Taker.Fee, settleBalance.Taker.InToken, statedb, mapBalances)
		if err != nil {
			return err
		}
		mapBalances[settleBalance.Taker.InToken][takerExOwner] = newTakerFee

		newCollateralTokenLock, err := lendingstate.CheckAddTokenBalance(common.HexToAddress(common.LendingLockAddress), settleBalance.Taker.OutTotal, settleBalance.Taker.OutToken, statedb, mapBalances)
		if err != nil {
			return err
		}
		mapBalances[settleBalance.Taker.OutToken][common.HexToAddress(common.LendingLockAddress)] = newCollateralTokenLock
	} else {
		relayerFee, err := lendingstate.CheckSubRelayerFee(makerOrder.Relayer, common.RelayerLendingFee, statedb, map[common.Address]*big.Int{})
		if err != nil {
			return err
		}
		lendingstate.SetSubRelayerFee(makerOrder.Relayer, relayerFee, common.RelayerLendingFee, statedb)
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
		newMakerFee, err := lendingstate.CheckAddTokenBalance(makerExOwner, settleBalance.Maker.Fee, settleBalance.Maker.InToken, statedb, mapBalances)
		if err != nil {
			return err
		}
		mapBalances[settleBalance.Maker.InToken][makerExOwner] = newMakerFee

		newCollateralTokenLock, err := lendingstate.CheckAddTokenBalance(common.HexToAddress(common.LendingLockAddress), settleBalance.Maker.OutTotal, settleBalance.Maker.OutToken, statedb, mapBalances)
		if err != nil {
			return err
		}
		mapBalances[settleBalance.Maker.OutToken][common.HexToAddress(common.LendingLockAddress)] = newCollateralTokenLock
	}
	masternodeOwner := statedb.GetOwner(coinbase)
	statedb.AddBalance(masternodeOwner, matchingFee)
	for token, balances := range mapBalances {
		for adrr, value := range balances {
			lendingstate.SetTokenBalance(adrr, value, token, statedb)
		}
	}
	return nil
}

func (l *Lending) ProcessCancelOrder(lendingStateDB *lendingstate.LendingStateDB, statedb *state.StateDB, tradingStateDb *tradingstate.TradingStateDB, chain consensus.ChainContext, coinbase common.Address, lendingOrderBook common.Hash, order *lendingstate.LendingItem) (error, bool) {
	originOrder := lendingStateDB.GetLendingOrder(lendingOrderBook, common.BigToHash(new(big.Int).SetUint64(order.LendingId)))
	if originOrder == lendingstate.EmptyLendingOrder {
		return fmt.Errorf("lendingOrder not found. Id: %v. LendToken: %s . Term: %v. CollateralToken: %v", order.LendingId, order.LendingToken.Hex(), order.Term, order.CollateralToken.Hex()), false
	}
	if err := lendingstate.CheckRelayerFee(originOrder.Relayer, common.RelayerCancelFee, statedb); err != nil {
		log.Debug("Relayer not enough fee when cancel order", "err", err)
		return nil, true
	}
	lendTokenDecimal, err := l.tomox.GetTokenDecimal(chain, statedb, coinbase, originOrder.LendingToken)
	if err != nil || lendTokenDecimal.Sign() == 0 {
		log.Debug("Fail to get tokenDecimal ", "Token", originOrder.LendingToken.String(), "err", err)
		return err, false
	}
	var tokenBalance *big.Int
	switch originOrder.Side {
	case lendingstate.Investing:
		tokenBalance = lendingstate.GetTokenBalance(originOrder.UserAddress, originOrder.LendingToken, statedb)
	case lendingstate.Borrowing:
		tokenBalance = lendingstate.GetTokenBalance(originOrder.UserAddress, originOrder.CollateralToken, statedb)
	default:
		log.Debug("Not found order side", "Side", originOrder.Side)
		return nil, true
	}
	log.Debug("ProcessCancelOrder", "LendingToken", originOrder.LendingToken, "CollateralToken", originOrder.CollateralToken, "makerInterest", originOrder.Interest, "lendTokenDecimal", lendTokenDecimal, "quantity", originOrder.Quantity)
	borrowFee := lendingstate.GetFee(statedb, originOrder.Relayer)
	collateralPrice := common.BasePrice

	if originOrder.Side == lendingstate.Borrowing {
		_, collateralPrice, err = l.GetCollateralPrices(chain, statedb, tradingStateDb, originOrder.CollateralToken, originOrder.LendingToken)
		if err != nil {
			return err, false
		}
	}
	tokenCancelFee := getCancelFee(lendTokenDecimal, collateralPrice, borrowFee, &originOrder)
	if tokenBalance.Cmp(tokenCancelFee) < 0 {
		log.Debug("User not enough balance when cancel order", "Side", originOrder.Side, "Interest", originOrder.Interest, "Quantity", originOrder.Quantity, "balance", tokenBalance, "fee", tokenCancelFee)
		return nil, true
	}
	err = lendingStateDB.CancelLendingOrder(lendingOrderBook, order)
	if err != nil {
		log.Debug("Error when cancel order", "order", order)
		return err, false
	}
	// relayers pay TOMO for masternode
	lendingstate.SubRelayerFee(originOrder.Relayer, common.RelayerCancelFee, statedb)
	masternodeOwner := statedb.GetOwner(coinbase)
	statedb.AddBalance(masternodeOwner, common.RelayerCancelFee)

	relayerOwner := lendingstate.GetRelayerOwner(originOrder.Relayer, statedb)
	switch originOrder.Side {
	case lendingstate.Investing:
		// users pay token for relayer
		lendingstate.SubTokenBalance(originOrder.UserAddress, tokenCancelFee, originOrder.LendingToken, statedb)
		lendingstate.AddTokenBalance(relayerOwner, tokenCancelFee, originOrder.LendingToken, statedb)
	case lendingstate.Borrowing:
		// users pay token for relayer
		lendingstate.SubTokenBalance(originOrder.UserAddress, tokenCancelFee, originOrder.CollateralToken, statedb)
		lendingstate.AddTokenBalance(relayerOwner, tokenCancelFee, originOrder.CollateralToken, statedb)
	default:
	}

	return nil, false
}

func (l *Lending) ProcessDeposit(lendingStateDB *lendingstate.LendingStateDB, statedb *state.StateDB, tradingStateDb *tradingstate.TradingStateDB, order *lendingstate.LendingItem) (error, bool, *lendingstate.LendingTrade) {
	lendingTradeId := common.Uint64ToHash(order.LendingTradeId)
	lendingBook := lendingstate.GetLendingOrderBookHash(order.LendingToken, order.Term)
	lendingTrade := lendingStateDB.GetLendingTrade(lendingBook, lendingTradeId)
	if lendingTrade == lendingstate.EmptyLendingTrade {
		return fmt.Errorf("process deposit for emptyLendingTrade is not allowed. lendingTradeId: %v", lendingTradeId.Hex()), true, nil
	}
	if order.Quantity.Sign() <= 0 || lendingTrade.TradeId != lendingTradeId.Big().Uint64() {
		log.Debug("invalid order deposit", "Quantity", order.Quantity, "lendingTradeId", lendingTradeId.Hex())
		return nil, true, nil
	}
	tokenBalance := lendingstate.GetTokenBalance(order.UserAddress, lendingTrade.CollateralToken, statedb)
	if tokenBalance.Cmp(order.Quantity) < 0 {
		log.Debug("not enough balance deposit", "Quantity", order.Quantity, "tokenBalance", tokenBalance)
		return nil, true, nil
	}
	tradingStateDb.RemoveLiquidationPrice(tradingstate.GetTradingOrderBookHash(lendingTrade.CollateralToken, lendingTrade.LendingToken), lendingTrade.LiquidationPrice, lendingBook, lendingTrade.TradeId)

	lendingstate.SubTokenBalance(order.UserAddress, order.Quantity, lendingTrade.CollateralToken, statedb)
	lendingstate.AddTokenBalance(common.HexToAddress(common.LendingLockAddress), order.Quantity, lendingTrade.CollateralToken, statedb)
	oldLockedAmount := lendingTrade.CollateralLockedAmount
	newLockedAmount := new(big.Int).Add(order.Quantity, oldLockedAmount)
	newLiquidationPrice := new(big.Int).Mul(lendingTrade.LiquidationPrice, oldLockedAmount)
	newLiquidationPrice = new(big.Int).Div(newLiquidationPrice, newLockedAmount)
	lendingStateDB.UpdateLiquidationPrice(lendingBook, lendingTrade.TradeId, newLiquidationPrice)
	lendingStateDB.UpdateCollateralLockedAmount(lendingBook, lendingTrade.TradeId, newLockedAmount)
	tradingStateDb.InsertLiquidationPrice(tradingstate.GetTradingOrderBookHash(lendingTrade.CollateralToken, lendingTrade.LendingToken), newLiquidationPrice, lendingBook, lendingTrade.TradeId)
	newLendingTrade := lendingTrade
	newLendingTrade.LiquidationPrice = newLiquidationPrice
	newLendingTrade.CollateralLockedAmount = newLockedAmount
	return nil, false, &newLendingTrade
}

// return hash: hash of lendingTrade
func (l *Lending) ProcessPayment(time uint64, lendingStateDB *lendingstate.LendingStateDB, statedb *state.StateDB, tradingstateDB *tradingstate.TradingStateDB, lendingBook common.Hash, lendingTradeId uint64) (hash common.Hash, err error) {
	lendingTradeIdHash := common.Uint64ToHash(lendingTradeId)
	lendingTrade := lendingStateDB.GetLendingTrade(lendingBook, lendingTradeIdHash)
	if lendingTrade == lendingstate.EmptyLendingTrade {
		return common.Hash{}, fmt.Errorf("process payment for emptyLendingTrade is not allowed. lendingTradeId: %v", lendingTradeId)
	}
	tokenBalance := lendingstate.GetTokenBalance(lendingTrade.Borrower, lendingTrade.LendingToken, statedb)
	interestRate := lendingstate.CalculateInterestRate(time, lendingTrade.LiquidationTime, lendingTrade.Term, lendingTrade.Interest)

	paymentBalance := new(big.Int).Mul(lendingTrade.Amount, new(big.Int).Add(common.BaseLendingInterest, interestRate))
	paymentBalance = new(big.Int).Div(paymentBalance, common.BaseLendingInterest)
	if tokenBalance.Cmp(paymentBalance) < 0 {
		if lendingTrade.LiquidationTime > time {
			return common.Hash{}, fmt.Errorf("Not enough balance need : %s , have : %s ", paymentBalance, tokenBalance)
		}
		hash, err = l.LiquidationTrade(lendingStateDB, statedb, tradingstateDB, lendingBook, lendingTradeId)
		return hash, err
	} else {
		lendingstate.SubTokenBalance(lendingTrade.Borrower, paymentBalance, lendingTrade.LendingToken, statedb)
		lendingstate.AddTokenBalance(lendingTrade.Investor, paymentBalance, lendingTrade.LendingToken, statedb)

		lendingstate.SubTokenBalance(common.HexToAddress(common.LendingLockAddress), lendingTrade.CollateralLockedAmount, lendingTrade.CollateralToken, statedb)
		lendingstate.AddTokenBalance(lendingTrade.Borrower, lendingTrade.CollateralLockedAmount, lendingTrade.CollateralToken, statedb)

		err=lendingStateDB.RemoveLiquidationTime(lendingBook, lendingTradeId, common.Uint64ToHash(lendingTrade.LiquidationTime))
		if err != nil {
			log.Debug("ProcessPayment RemoveLiquidationTime", "err", err)
		}
		err=tradingstateDB.RemoveLiquidationPrice(tradingstate.GetTradingOrderBookHash(lendingTrade.CollateralToken, lendingTrade.LendingToken), lendingTrade.LiquidationPrice, lendingBook, lendingTradeId)
		if err != nil {
			log.Debug("ProcessPayment RemoveLiquidationPrice", "err", err)
		}
		lendingStateDB.CancelLendingTrade(lendingBook, lendingTradeId)
		if err != nil {
			log.Debug("ProcessPayment CancelLendingTrade", "err", err)
		}
	}
	return lendingTrade.Hash, nil
}

// return hash of liquidatedTrade
func (l *Lending) LiquidationTrade(lendingStateDB *lendingstate.LendingStateDB, statedb *state.StateDB, tradingstateDB *tradingstate.TradingStateDB, lendingBook common.Hash, lendingTradeId uint64) (common.Hash, error) {
	lendingTradeIdHash := common.Uint64ToHash(lendingTradeId)
	lendingTrade := lendingStateDB.GetLendingTrade(lendingBook, lendingTradeIdHash)
	if lendingTrade.TradeId != lendingTradeId {
		return common.Hash{}, fmt.Errorf("Lending Trade Id not found : %d ", lendingTradeId)
	}
	lendingstate.SubTokenBalance(common.HexToAddress(common.LendingLockAddress), lendingTrade.CollateralLockedAmount, lendingTrade.CollateralToken, statedb)
	lendingstate.AddTokenBalance(lendingTrade.Investor, lendingTrade.CollateralLockedAmount, lendingTrade.CollateralToken, statedb)

	err := lendingStateDB.RemoveLiquidationTime(lendingBook, lendingTradeId, common.Uint64ToHash(lendingTrade.LiquidationTime))
	if err != nil {
		log.Debug("LiquidationTrade RemoveLiquidationTime", "err", err)
	}
	err = tradingstateDB.RemoveLiquidationPrice(tradingstate.GetTradingOrderBookHash(lendingTrade.CollateralToken, lendingTrade.LendingToken), lendingTrade.LiquidationPrice, lendingBook, lendingTradeId)
	if err != nil {
		log.Debug("LiquidationTrade RemoveLiquidationPrice", "err", err)
	}
	err = lendingStateDB.CancelLendingTrade(lendingBook, lendingTradeId)
	if err != nil {
		log.Debug("LiquidationTrade CancelLendingTrade", "err", err)
	}
	return lendingTrade.Hash, nil
}
func getCancelFee(lendTokenDecimal *big.Int, collateralPrice, borrowFee *big.Int, order *lendingstate.LendingItem) *big.Int {
	cancelFee := big.NewInt(0)
	if order.Side == lendingstate.Investing {
		// cancel fee = quantityToLend*borrowFee/LendingCancelFee
		cancelFee = new(big.Int).Mul(order.Quantity, borrowFee)
		cancelFee = new(big.Int).Div(cancelFee, common.LendingCancelFee)
	} else {
		// Fee ==  quantityToLend/base lend token decimal *price*borrowFee/LendingCancelFee
		cancelFee = new(big.Int).Mul(order.Quantity, collateralPrice)
		cancelFee = new(big.Int).Mul(cancelFee, borrowFee)
		cancelFee = new(big.Int).Div(cancelFee, lendTokenDecimal)
		cancelFee = new(big.Int).Div(cancelFee, common.LendingCancelFee)
	}
	return cancelFee
}

func (l *Lending) getMediumTradePriceLastEpoch(chain consensus.ChainContext, statedb *state.StateDB, tradingStateDb *tradingstate.TradingStateDB, baseToken common.Address, quoteToken common.Address) (*big.Int, error) {
	price := tradingStateDb.GetMediumPriceLastEpoch(tradingstate.GetTradingOrderBookHash(baseToken, quoteToken))
	if price != nil && price.Sign() > 0 {
		return price, nil
	} else {
		inversePrice := tradingStateDb.GetMediumPriceLastEpoch(tradingstate.GetTradingOrderBookHash(quoteToken, baseToken))
		if inversePrice != nil && inversePrice.Sign() > 0 {
			quoteTokenDecimal, err := l.tomox.GetTokenDecimal(chain, statedb, common.Address{}, quoteToken)
			if err != nil || quoteTokenDecimal.Sign() == 0 {
				return nil, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", quoteToken.String(), err)
			}
			baseTokenDecimal, err := l.tomox.GetTokenDecimal(chain, statedb, common.Address{}, baseToken)
			if err != nil || baseTokenDecimal.Sign() == 0 {
				return nil, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", baseToken, err)
			}
			price = new(big.Int).Mul(baseTokenDecimal, quoteTokenDecimal)
			price = new(big.Int).Div(price, inversePrice)
			return price, nil
		}
	}
	return nil, nil
}

func (l *Lending) GetCollateralPrices(chain consensus.ChainContext, statedb *state.StateDB, tradingStateDb *tradingstate.TradingStateDB, collateralToken common.Address, lendingToken common.Address) (*big.Int, *big.Int, error) {
	// lendTokenTOMOPrice: price of ticker lendToken/TOMO
	// collateralTOMOPrice: price of ticker collateralToken/TOMO
	// collateralPrice: price of ticker collateralToken/lendToken

	_, _, collateralTOMOBasePrice := lendingstate.GetCollateralDetail(statedb, collateralToken)
	_, _, lendingTOMOBasePrice := lendingstate.GetCollateralDetail(statedb, lendingToken)
	collateralPrice := big.NewInt(0)
	lendTokenTOMOPrice, err := l.getMediumTradePriceLastEpoch(chain, statedb, tradingStateDb, lendingToken, common.HexToAddress(common.TomoNativeAddress))
	if err != nil {
		return lendTokenTOMOPrice, collateralPrice, err
	}
	if lendTokenTOMOPrice == nil || lendTokenTOMOPrice.Sign() == 0 {
		lendTokenTOMOPrice = lendingTOMOBasePrice
	}
	if collateralToken.String() != common.TomoNativeAddress {
		lastMediumPrice, err := l.getMediumTradePriceLastEpoch(chain, statedb, tradingStateDb, collateralToken, lendingToken)
		if err != nil {
			return lendTokenTOMOPrice, collateralPrice, err
		}
		log.Debug("GetCollateralPrices", "lastMediumPrice", lastMediumPrice)
		if lastMediumPrice != nil && lastMediumPrice.Sign() > 0 {
			collateralPrice = lastMediumPrice
		} else {
			collateralTOMOPrice, err := l.getMediumTradePriceLastEpoch(chain, statedb, tradingStateDb, collateralToken, common.HexToAddress(common.TomoNativeAddress))
			log.Debug("GetCollateralPrices", "collateralTOMOPrice", collateralTOMOPrice)
			if err != nil {
				return collateralPrice, lendTokenTOMOPrice, err
			}
			if collateralTOMOPrice == nil || collateralTOMOPrice.Sign() == 0 {
				collateralTOMOPrice = collateralTOMOBasePrice
			}
			if lendTokenTOMOPrice != nil && lendTokenTOMOPrice.Sign() > 0 {
				collateralPrice = new(big.Int).Div(collateralTOMOPrice, lendTokenTOMOPrice)
				lendingTokenDecimal, err := l.tomox.GetTokenDecimal(chain, statedb, common.Address{}, lendingToken)
				if err != nil {
					return nil, nil, err
				}
				collateralPrice = new(big.Int).Mul(collateralPrice, lendingTokenDecimal)
				log.Debug("GetCollateralPrices: Calculate collateral/LendToken price from collateral/TOMO, lendToken/TOMO", "collateralPrice", collateralPrice)
			}
		}
	}
	return lendTokenTOMOPrice, collateralPrice, nil
}
