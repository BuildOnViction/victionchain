package chancoinx

import (
	"math/big"
	"strconv"
	"time"

	"github.com/chancoin-core/chancoin-gold/consensus"

	"fmt"

	"github.com/chancoin-core/chancoin-gold/common"
	"github.com/chancoin-core/chancoin-gold/core/state"
	"github.com/chancoin-core/chancoin-gold/log"
	"github.com/chancoin-core/chancoin-gold/chancoinx/chancoinx_state"
)

func (chancoinx *ChancoinX) CommitOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, chancoinXstatedb *chancoinx_state.ChancoinXStateDB, orderBook common.Hash, order *chancoinx_state.OrderItem) ([]map[string]string, []*chancoinx_state.OrderItem, error) {
	chancoinxSnap := chancoinXstatedb.Snapshot()
	dbSnap := statedb.Snapshot()
	trades, rejects, err := chancoinx.ApplyOrder(coinbase, chain, statedb, chancoinXstatedb, orderBook, order)
	if err != nil {
		chancoinXstatedb.RevertToSnapshot(chancoinxSnap)
		statedb.RevertToSnapshot(dbSnap)
		return nil, nil, err
	}
	return trades, rejects, err
}

func (chancoinx *ChancoinX) ApplyOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, chancoinXstatedb *chancoinx_state.ChancoinXStateDB, orderBook common.Hash, order *chancoinx_state.OrderItem) ([]map[string]string, []*chancoinx_state.OrderItem, error) {
	var (
		rejects []*chancoinx_state.OrderItem
		trades  []map[string]string
		err     error
	)
	nonce := chancoinXstatedb.GetNonce(order.UserAddress.Hash())
	log.Debug("ApplyOrder", "addr", order.UserAddress, "statenonce", nonce, "ordernonce", order.Nonce)
	if big.NewInt(int64(nonce)).Cmp(order.Nonce) == -1 {
		log.Debug("ApplyOrder ErrNonceTooHigh", "nonce", order.Nonce)
		return nil, nil, ErrNonceTooHigh
	} else if big.NewInt(int64(nonce)).Cmp(order.Nonce) == 1 {
		log.Debug("ApplyOrder ErrNonceTooLow", "nonce", order.Nonce)
		return nil, nil, ErrNonceTooLow
	}
	// increase nonce
	log.Debug("ApplyOrder setnonce", "nonce", nonce+1, "addr", order.UserAddress.Hex(), "status", order.Status, "oldnonce", nonce)
	chancoinXstatedb.SetNonce(order.UserAddress.Hash(), nonce+1)
	chancoinxSnap := chancoinXstatedb.Snapshot()
	dbSnap := statedb.Snapshot()
	defer func() {
		if err != nil {
			chancoinXstatedb.RevertToSnapshot(chancoinxSnap)
			statedb.RevertToSnapshot(dbSnap)
		}
	}()
	if order.Status == OrderStatusCancelled {
		err, reject := chancoinx.ProcessCancelOrder(chancoinXstatedb, statedb, chain, coinbase, orderBook, order)
		if err != nil || reject {
			log.Debug("Reject cancelled order", "err", err)
			rejects = append(rejects, order)
		}
		return trades, rejects, nil
	}
	if order.Type != chancoinx_state.Market {
		if order.Price.Sign() == 0 || common.BigToHash(order.Price).Big().Cmp(order.Price) != 0 {
			log.Debug("Reject order price invalid", "price", order.Price)
			rejects = append(rejects, order)
			return trades, rejects, nil
		}
	}
	if order.Quantity.Sign() == 0 || common.BigToHash(order.Quantity).Big().Cmp(order.Quantity) != 0 {
		log.Debug("Reject order quantity invalid", "quantity", order.Quantity)
		rejects = append(rejects, order)
		return trades, rejects, nil
	}
	orderType := order.Type
	// if we do not use auto-increment orderid, we must set price slot to avoid conflict
	if orderType == chancoinx_state.Market {
		log.Debug("Process maket order", "side", order.Side, "quantity", order.Quantity, "price", order.Price)
		trades, rejects, err = chancoinx.processMarketOrder(coinbase, chain, statedb, chancoinXstatedb, orderBook, order)
		if err != nil {
			log.Debug("Reject market order", "err", err, "order", chancoinx_state.ToJSON(order))
			trades = []map[string]string{}
			rejects = append(rejects, order)
		}
	} else {
		log.Debug("Process limit order", "side", order.Side, "quantity", order.Quantity, "price", order.Price)
		trades, rejects, err = chancoinx.processLimitOrder(coinbase, chain, statedb, chancoinXstatedb, orderBook, order)
		if err != nil {
			log.Debug("Reject limit order", "err", err, "order", chancoinx_state.ToJSON(order))
			trades = []map[string]string{}
			rejects = append(rejects, order)
		}
	}

	return trades, rejects, nil
}

// processMarketOrder : process the market order
func (chancoinx *ChancoinX) processMarketOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, chancoinXstatedb *chancoinx_state.ChancoinXStateDB, orderBook common.Hash, order *chancoinx_state.OrderItem) ([]map[string]string, []*chancoinx_state.OrderItem, error) {
	var (
		trades     []map[string]string
		newTrades  []map[string]string
		rejects    []*chancoinx_state.OrderItem
		newRejects []*chancoinx_state.OrderItem
		err        error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	// speedup the comparison, do not assign because it is pointer
	zero := chancoinx_state.Zero
	if side == chancoinx_state.Bid {
		bestPrice, volume := chancoinXstatedb.GetBestAskPrice(orderBook)
		log.Debug("processMarketOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && bestPrice.Cmp(zero) > 0 {
			quantityToTrade, newTrades, newRejects, err = chancoinx.processOrderList(coinbase, chain, statedb, chancoinXstatedb, chancoinx_state.Ask, orderBook, bestPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			bestPrice, volume = chancoinXstatedb.GetBestAskPrice(orderBook)
			log.Debug("processMarketOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade, "volume", volume)
		}
	} else {
		bestPrice, volume := chancoinXstatedb.GetBestBidPrice(orderBook)
		log.Debug("processMarketOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && bestPrice.Cmp(zero) > 0 {
			quantityToTrade, newTrades, newRejects, err = chancoinx.processOrderList(coinbase, chain, statedb, chancoinXstatedb, chancoinx_state.Bid, orderBook, bestPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			bestPrice, volume = chancoinXstatedb.GetBestBidPrice(orderBook)
			log.Debug("processMarketOrder ", "side", side, "bestPrice", bestPrice, "quantityToTrade", quantityToTrade, "volume", volume)
		}
	}
	return trades, newRejects, nil
}

// processLimitOrder : process the limit order, can change the quote
// If not care for performance, we should make a copy of quote to prevent further reference problem
func (chancoinx *ChancoinX) processLimitOrder(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, chancoinXstatedb *chancoinx_state.ChancoinXStateDB, orderBook common.Hash, order *chancoinx_state.OrderItem) ([]map[string]string, []*chancoinx_state.OrderItem, error) {
	var (
		trades     []map[string]string
		newTrades  []map[string]string
		rejects    []*chancoinx_state.OrderItem
		newRejects []*chancoinx_state.OrderItem
		err        error
	)
	quantityToTrade := order.Quantity
	side := order.Side
	price := order.Price

	// speedup the comparison, do not assign because it is pointer
	zero := chancoinx_state.Zero

	if side == chancoinx_state.Bid {
		minPrice, volume := chancoinXstatedb.GetBestAskPrice(orderBook)
		log.Debug("processLimitOrder ", "side", side, "minPrice", minPrice, "orderPrice", price, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && price.Cmp(minPrice) >= 0 && minPrice.Cmp(zero) > 0 {
			log.Debug("Min price in asks tree", "price", minPrice.String())
			quantityToTrade, newTrades, newRejects, err = chancoinx.processOrderList(coinbase, chain, statedb, chancoinXstatedb, chancoinx_state.Ask, orderBook, minPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			log.Debug("New trade found", "newTrades", newTrades, "quantityToTrade", quantityToTrade)
			minPrice, volume = chancoinXstatedb.GetBestAskPrice(orderBook)
			log.Debug("processLimitOrder ", "side", side, "minPrice", minPrice, "orderPrice", price, "volume", volume)
		}
	} else {
		maxPrice, volume := chancoinXstatedb.GetBestBidPrice(orderBook)
		log.Debug("processLimitOrder ", "side", side, "maxPrice", maxPrice, "orderPrice", price, "volume", volume)
		for quantityToTrade.Cmp(zero) > 0 && price.Cmp(maxPrice) <= 0 && maxPrice.Cmp(zero) > 0 {
			log.Debug("Max price in bids tree", "price", maxPrice.String())
			quantityToTrade, newTrades, newRejects, err = chancoinx.processOrderList(coinbase, chain, statedb, chancoinXstatedb, chancoinx_state.Bid, orderBook, maxPrice, quantityToTrade, order)
			if err != nil {
				return nil, nil, err
			}
			trades = append(trades, newTrades...)
			rejects = append(rejects, newRejects...)
			log.Debug("New trade found", "newTrades", newTrades, "quantityToTrade", quantityToTrade)
			maxPrice, volume = chancoinXstatedb.GetBestBidPrice(orderBook)
			log.Debug("processLimitOrder ", "side", side, "maxPrice", maxPrice, "orderPrice", price, "volume", volume)
		}
	}
	if quantityToTrade.Cmp(zero) > 0 {
		orderId := chancoinXstatedb.GetNonce(orderBook)
		order.OrderID = orderId + 1
		order.Quantity = quantityToTrade
		chancoinXstatedb.SetNonce(orderBook, orderId+1)
		orderIdHash := common.BigToHash(new(big.Int).SetUint64(order.OrderID))
		chancoinXstatedb.InsertOrderItem(orderBook, orderIdHash, *order)
		log.Debug("After matching, order (unmatched part) is now added to tree", "side", order.Side, "order", order)
	}
	return trades, rejects, nil
}

// processOrderList : process the order list
func (chancoinx *ChancoinX) processOrderList(coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, chancoinXstatedb *chancoinx_state.ChancoinXStateDB, side string, orderBook common.Hash, price *big.Int, quantityStillToTrade *big.Int, order *chancoinx_state.OrderItem) (*big.Int, []map[string]string, []*chancoinx_state.OrderItem, error) {
	quantityToTrade := chancoinx_state.CloneBigInt(quantityStillToTrade)
	log.Debug("Process matching between order and orderlist", "quantityToTrade", quantityToTrade)
	var (
		trades []map[string]string

		rejects []*chancoinx_state.OrderItem
	)
	for quantityToTrade.Sign() > 0 {
		orderId, amount, _ := chancoinXstatedb.GetBestOrderIdAndAmount(orderBook, price, side)
		var oldestOrder chancoinx_state.OrderItem
		if amount.Sign() > 0 {
			oldestOrder = chancoinXstatedb.GetOrder(orderBook, orderId)
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
			maxTradedQuantity = chancoinx_state.CloneBigInt(quantityToTrade)
		} else {
			maxTradedQuantity = chancoinx_state.CloneBigInt(amount)
		}
		var quotePrice *big.Int
		if oldestOrder.QuoteToken.String() != common.ChancoinNativeAddress {
			quotePrice = chancoinXstatedb.GetPrice(chancoinx_state.GetOrderBookHash(oldestOrder.QuoteToken, common.HexToAddress(common.ChancoinNativeAddress)))
			log.Debug("TryGet quotePrice QuoteToken/CHANCOIN", "quotePrice", quotePrice)
			if (quotePrice == nil || quotePrice.Sign() == 0) && oldestOrder.BaseToken.String() != common.ChancoinNativeAddress {
				inversePrice := chancoinXstatedb.GetPrice(chancoinx_state.GetOrderBookHash(common.HexToAddress(common.ChancoinNativeAddress), oldestOrder.QuoteToken))
				quoteTokenDecimal, err := chancoinx.GetTokenDecimal(chain, statedb, coinbase, oldestOrder.QuoteToken)
				if err != nil || quoteTokenDecimal.Sign() == 0 {
					return nil, nil, nil, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", oldestOrder.QuoteToken.String(), err)
				}
				log.Debug("TryGet inversePrice CHANCOIN/QuoteToken", "inversePrice", inversePrice)
				if inversePrice != nil && inversePrice.Sign() > 0 {
					quotePrice = new(big.Int).Div(common.BasePrice, inversePrice)
					quotePrice = new(big.Int).Mul(quotePrice, quoteTokenDecimal)
					log.Debug("TryGet quotePrice after get inversePrice CHANCOIN/QuoteToken", "quotePrice", quotePrice, "quoteTokenDecimal", quoteTokenDecimal)
				}
			}
		}
		tradedQuantity, rejectMaker, err := chancoinx.getTradeQuantity(quotePrice, coinbase, chain, statedb, order, &oldestOrder, maxTradedQuantity)
		if err != nil && err == chancoinx_state.ErrQuantityTradeTooSmall {
			if tradedQuantity.Cmp(maxTradedQuantity) == 0 {
				if quantityToTrade.Cmp(amount) == 0 { // reject Taker & maker
					rejects = append(rejects, order)
					quantityToTrade = chancoinx_state.Zero
					rejects = append(rejects, &oldestOrder)
					err = chancoinXstatedb.CancelOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					break
				} else if quantityToTrade.Cmp(amount) < 0 { // reject Taker
					rejects = append(rejects, order)
					quantityToTrade = chancoinx_state.Zero
					break
				} else { // reject maker
					rejects = append(rejects, &oldestOrder)
					err = chancoinXstatedb.CancelOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					continue
				}
			} else {
				if rejectMaker { // reject maker
					rejects = append(rejects, &oldestOrder)
					err = chancoinXstatedb.CancelOrder(orderBook, &oldestOrder)
					if err != nil {
						return nil, nil, nil, err
					}
					continue
				} else { // reject Taker
					rejects = append(rejects, order)
					quantityToTrade = chancoinx_state.Zero
					break
				}
			}
		} else if err != nil {
			return nil, nil, nil, err
		}
		if tradedQuantity.Sign() == 0 && !rejectMaker {
			log.Debug("Reject order Taker ", "tradedQuantity", tradedQuantity, "rejectMaker", rejectMaker)
			rejects = append(rejects, order)
			quantityToTrade = chancoinx_state.Zero
			break
		}
		if tradedQuantity.Sign() > 0 {
			quantityToTrade = chancoinx_state.Sub(quantityToTrade, tradedQuantity)
			chancoinXstatedb.SubAmountOrderItem(orderBook, orderId, price, tradedQuantity, side)
			chancoinXstatedb.SetPrice(orderBook, price)
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
			err := chancoinXstatedb.CancelOrder(orderBook, &oldestOrder)
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}
	return quantityToTrade, trades, rejects, nil
}

func (chancoinx *ChancoinX) getTradeQuantity(quotePrice *big.Int, coinbase common.Address, chain consensus.ChainContext, statedb *state.StateDB, takerOrder *chancoinx_state.OrderItem, makerOrder *chancoinx_state.OrderItem, quantityToTrade *big.Int) (*big.Int, bool, error) {
	baseTokenDecimal, err := chancoinx.GetTokenDecimal(chain, statedb, coinbase, makerOrder.BaseToken)
	if err != nil || baseTokenDecimal.Sign() == 0 {
		return chancoinx_state.Zero, false, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", makerOrder.BaseToken.String(), err)
	}
	quoteTokenDecimal, err := chancoinx.GetTokenDecimal(chain, statedb, coinbase, makerOrder.QuoteToken)
	if err != nil || quoteTokenDecimal.Sign() == 0 {
		return chancoinx_state.Zero, false, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", makerOrder.QuoteToken.String(), err)
	}
	if makerOrder.QuoteToken.String() == common.ChancoinNativeAddress {
		quotePrice = quoteTokenDecimal
	}
	if takerOrder.ExchangeAddress.String() == makerOrder.ExchangeAddress.String() {
		if err := chancoinx_state.CheckRelayerFee(takerOrder.ExchangeAddress, new(big.Int).Mul(common.RelayerFee, big.NewInt(2)), statedb); err != nil {
			log.Debug("Reject order Taker Exchnage = Maker Exchange , relayer not enough fee ", "err", err)
			return chancoinx_state.Zero, false, nil
		}
	} else {
		if err := chancoinx_state.CheckRelayerFee(takerOrder.ExchangeAddress, common.RelayerFee, statedb); err != nil {
			log.Debug("Reject order Taker , relayer not enough fee ", "err", err)
			return chancoinx_state.Zero, false, nil
		}
		if err := chancoinx_state.CheckRelayerFee(makerOrder.ExchangeAddress, common.RelayerFee, statedb); err != nil {
			log.Debug("Reject order maker , relayer not enough fee ", "err", err)
			return chancoinx_state.Zero, true, nil
		}
	}
	takerFeeRate := chancoinx_state.GetExRelayerFee(takerOrder.ExchangeAddress, statedb)
	makerFeeRate := chancoinx_state.GetExRelayerFee(makerOrder.ExchangeAddress, statedb)
	var takerBalance, makerBalance *big.Int
	switch takerOrder.Side {
	case chancoinx_state.Bid:
		takerBalance = chancoinx_state.GetTokenBalance(takerOrder.UserAddress, makerOrder.QuoteToken, statedb)
		makerBalance = chancoinx_state.GetTokenBalance(makerOrder.UserAddress, makerOrder.BaseToken, statedb)
	case chancoinx_state.Ask:
		takerBalance = chancoinx_state.GetTokenBalance(takerOrder.UserAddress, makerOrder.BaseToken, statedb)
		makerBalance = chancoinx_state.GetTokenBalance(makerOrder.UserAddress, makerOrder.QuoteToken, statedb)
	default:
		takerBalance = big.NewInt(0)
		makerBalance = big.NewInt(0)
	}
	quantity, rejectMaker := GetTradeQuantity(takerOrder.Side, takerFeeRate, takerBalance, makerOrder.Price, makerFeeRate, makerBalance, baseTokenDecimal, quantityToTrade)
	log.Debug("GetTradeQuantity", "side", takerOrder.Side, "takerBalance", takerBalance, "makerBalance", makerBalance, "BaseToken", makerOrder.BaseToken, "QuoteToken", makerOrder.QuoteToken, "quantity", quantity, "rejectMaker", rejectMaker, "quotePrice", quotePrice)
	if quantity.Sign() > 0 {
		// Apply Match Order
		settleBalanceResult, err := chancoinx_state.GetSettleBalance(quotePrice, takerOrder.Side, takerFeeRate, makerOrder.BaseToken, makerOrder.QuoteToken, makerOrder.Price, makerFeeRate, baseTokenDecimal, quoteTokenDecimal, quantity)
		log.Debug("GetSettleBalance", "settleBalanceResult", settleBalanceResult, "err", err)
		if err == nil {
			err = DoSettleBalance(coinbase, takerOrder, makerOrder, settleBalanceResult, statedb)
		}
		return quantity, rejectMaker, err
	}
	return quantity, rejectMaker, nil
}

func GetTradeQuantity(takerSide string, takerFeeRate *big.Int, takerBalance *big.Int, makerPrice *big.Int, makerFeeRate *big.Int, makerBalance *big.Int, baseTokenDecimal *big.Int, quantityToTrade *big.Int) (*big.Int, bool) {
	if takerSide == chancoinx_state.Bid {
		// maker InQuantity quoteTokenQuantity=(quantityToTrade*maker.Price/baseTokenDecimal)
		quoteTokenQuantity := new(big.Int).Mul(quantityToTrade, makerPrice)
		quoteTokenQuantity = quoteTokenQuantity.Div(quoteTokenQuantity, baseTokenDecimal)
		// Fee
		// charge on the token he/she has before the trade, in this case: quoteToken
		// charge on the token he/she has before the trade, in this case: baseToken
		// takerFee = quoteTokenQuantity*takerFeeRate/baseFee=(quantityToTrade*maker.Price/baseTokenDecimal) * makerFeeRate/baseFee
		takerFee := big.NewInt(0).Mul(quoteTokenQuantity, takerFeeRate)
		takerFee = big.NewInt(0).Div(takerFee, common.ChancoinXBaseFee)
		//takerOutTotal= quoteTokenQuantity + takerFee =  quantityToTrade*maker.Price/baseTokenDecimal + quantityToTrade*maker.Price/baseTokenDecimal * takerFeeRate/baseFee
		// = quantityToTrade *  maker.Price/baseTokenDecimal ( 1 +  takerFeeRate/baseFee)
		// = quantityToTrade * maker.Price * (baseFee + takerFeeRate ) / ( baseTokenDecimal * baseFee)
		takerOutTotal := new(big.Int).Add(quoteTokenQuantity, takerFee)
		makerOutTotal := quantityToTrade
		if takerBalance.Cmp(takerOutTotal) >= 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			return quantityToTrade, false
		} else if takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) >= 0 {
			newQuantityTrade := new(big.Int).Mul(takerBalance, baseTokenDecimal)
			newQuantityTrade = newQuantityTrade.Mul(newQuantityTrade, common.ChancoinXBaseFee)
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, new(big.Int).Add(common.ChancoinXBaseFee, takerFeeRate))
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
			newQuantityTrade = newQuantityTrade.Mul(newQuantityTrade, common.ChancoinXBaseFee)
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, new(big.Int).Add(common.ChancoinXBaseFee, takerFeeRate))
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
		makerFee = makerFee.Div(makerFee, common.ChancoinXBaseFee)

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
			newQuantityTrade = newQuantityTrade.Mul(newQuantityTrade, common.ChancoinXBaseFee)
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, new(big.Int).Add(common.ChancoinXBaseFee, makerFeeRate))
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, makerPrice)
			log.Debug("Reject order maker , not enough balance ", "makerBalance", makerBalance, " makerOutTotal", makerOutTotal)
			return newQuantityTrade, true
		} else {
			// takerBalance.Cmp(takerOutTotal) < 0 && makerBalance.Cmp(makerOutTotal) < 0
			newQuantityTrade := new(big.Int).Mul(makerBalance, baseTokenDecimal)
			newQuantityTrade = newQuantityTrade.Mul(newQuantityTrade, common.ChancoinXBaseFee)
			newQuantityTrade = newQuantityTrade.Div(newQuantityTrade, new(big.Int).Add(common.ChancoinXBaseFee, makerFeeRate))
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

func DoSettleBalance(coinbase common.Address, takerOrder, makerOrder *chancoinx_state.OrderItem, settleBalance *chancoinx_state.SettleBalance, statedb *state.StateDB) error {
	takerExOwner := chancoinx_state.GetRelayerOwner(takerOrder.ExchangeAddress, statedb)
	makerExOwner := chancoinx_state.GetRelayerOwner(makerOrder.ExchangeAddress, statedb)
	matchingFee := big.NewInt(0)
	// masternodes charges fee of both 2 relayers. If maker and Taker are on same relayer, that relayer is charged fee twice
	matchingFee = matchingFee.Add(matchingFee, common.RelayerFee)
	matchingFee = matchingFee.Add(matchingFee, common.RelayerFee)

	if common.EmptyHash(takerExOwner.Hash()) || common.EmptyHash(makerExOwner.Hash()) {
		return fmt.Errorf("Echange owner empty , Taker: %v , maker : %v ", takerExOwner, makerExOwner)
	}
	mapBalances := map[common.Address]map[common.Address]*big.Int{}
	//Checking balance
	newTakerInTotal, err := chancoinx_state.CheckAddTokenBalance(takerOrder.UserAddress, settleBalance.Taker.InTotal, settleBalance.Taker.InToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Taker.InToken] == nil {
		mapBalances[settleBalance.Taker.InToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Taker.InToken][takerOrder.UserAddress] = newTakerInTotal
	newTakerOutTotal, err := chancoinx_state.CheckSubTokenBalance(takerOrder.UserAddress, settleBalance.Taker.OutTotal, settleBalance.Taker.OutToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Taker.OutToken] == nil {
		mapBalances[settleBalance.Taker.OutToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Taker.OutToken][takerOrder.UserAddress] = newTakerOutTotal
	newMakerInTotal, err := chancoinx_state.CheckAddTokenBalance(makerOrder.UserAddress, settleBalance.Maker.InTotal, settleBalance.Maker.InToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Maker.InToken] == nil {
		mapBalances[settleBalance.Maker.InToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Maker.InToken][makerOrder.UserAddress] = newMakerInTotal
	newMakerOutTotal, err := chancoinx_state.CheckSubTokenBalance(makerOrder.UserAddress, settleBalance.Maker.OutTotal, settleBalance.Maker.OutToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[settleBalance.Maker.OutToken] == nil {
		mapBalances[settleBalance.Maker.OutToken] = map[common.Address]*big.Int{}
	}
	mapBalances[settleBalance.Maker.OutToken][makerOrder.UserAddress] = newMakerOutTotal
	newTakerFee, err := chancoinx_state.CheckAddTokenBalance(takerExOwner, settleBalance.Taker.Fee, makerOrder.QuoteToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	if mapBalances[makerOrder.QuoteToken] == nil {
		mapBalances[makerOrder.QuoteToken] = map[common.Address]*big.Int{}
	}
	mapBalances[makerOrder.QuoteToken][takerExOwner] = newTakerFee
	newMakerFee, err := chancoinx_state.CheckAddTokenBalance(makerExOwner, settleBalance.Maker.Fee, makerOrder.QuoteToken, statedb, mapBalances)
	if err != nil {
		return err
	}
	mapBalances[makerOrder.QuoteToken][makerExOwner] = newMakerFee

	mapRelayerFee := map[common.Address]*big.Int{}
	newRelayerTakerFee, err := chancoinx_state.CheckSubRelayerFee(takerOrder.ExchangeAddress, common.RelayerFee, statedb, mapRelayerFee)
	if err != nil {
		return err
	}
	mapRelayerFee[takerOrder.ExchangeAddress] = newRelayerTakerFee
	newRelayerMakerFee, err := chancoinx_state.CheckSubRelayerFee(makerOrder.ExchangeAddress, common.RelayerFee, statedb, mapRelayerFee)
	if err != nil {
		return err
	}
	mapRelayerFee[makerOrder.ExchangeAddress] = newRelayerMakerFee
	chancoinx_state.SetSubRelayerFee(takerOrder.ExchangeAddress, newRelayerTakerFee, common.RelayerFee, statedb)
	chancoinx_state.SetSubRelayerFee(makerOrder.ExchangeAddress, newRelayerMakerFee, common.RelayerFee, statedb)

	masternodeOwner := statedb.GetOwner(coinbase)
	statedb.AddBalance(masternodeOwner, matchingFee)

	chancoinx_state.SetTokenBalance(takerOrder.UserAddress, newTakerInTotal, settleBalance.Taker.InToken, statedb)
	chancoinx_state.SetTokenBalance(takerOrder.UserAddress, newTakerOutTotal, settleBalance.Taker.OutToken, statedb)

	chancoinx_state.SetTokenBalance(makerOrder.UserAddress, newMakerInTotal, settleBalance.Maker.InToken, statedb)
	chancoinx_state.SetTokenBalance(makerOrder.UserAddress, newMakerOutTotal, settleBalance.Maker.OutToken, statedb)

	// add balance for relayers
	//log.Debug("ApplyChancoinXMatchedTransaction settle fee for relayers",
	//	"takerRelayerOwner", takerExOwner,
	//	"takerFeeToken", quoteToken, "takerFee", settleBalanceResult[takerAddr][chancoinx.Fee].(*big.Int),
	//	"makerRelayerOwner", makerExOwner,
	//	"makerFeeToken", quoteToken, "makerFee", settleBalanceResult[makerAddr][chancoinx.Fee].(*big.Int))
	// takerFee
	chancoinx_state.SetTokenBalance(takerExOwner, newTakerFee, makerOrder.QuoteToken, statedb)
	chancoinx_state.SetTokenBalance(makerExOwner, newMakerFee, makerOrder.QuoteToken, statedb)
	return nil
}

func (chancoinx *ChancoinX) ProcessCancelOrder(chancoinXstatedb *chancoinx_state.ChancoinXStateDB, statedb *state.StateDB, chain consensus.ChainContext, coinbase common.Address, orderBook common.Hash, order *chancoinx_state.OrderItem) (error, bool) {
	if err := chancoinx_state.CheckRelayerFee(order.ExchangeAddress, common.RelayerCancelFee, statedb); err != nil {
		log.Debug("Relayer not enough fee when cancel order", "err", err)
		return nil, true
	}
	baseTokenDecimal, err := chancoinx.GetTokenDecimal(chain, statedb, coinbase, order.BaseToken)
	if err != nil || baseTokenDecimal.Sign() == 0 {
		log.Debug("Fail to get tokenDecimal ", "Token", order.BaseToken.String(), "err", err)
		return err, false
	}
	// order: basic order information (includes orderId, orderHash, baseToken, quoteToken) which user send to chancoinx to cancel order
	// originOrder: full order information getting from order trie
	originOrder := chancoinXstatedb.GetOrder(orderBook, common.BigToHash(new(big.Int).SetUint64(order.OrderID)))
	if originOrder == chancoinx_state.EmptyOrder {
		return fmt.Errorf("order not found. OrderId: %v. Base: %s. Quote: %s", order.OrderID, order.BaseToken.Hex(), order.QuoteToken.Hex()), false
	}
	var tokenBalance *big.Int
	switch originOrder.Side {
	case chancoinx_state.Ask:
		tokenBalance = chancoinx_state.GetTokenBalance(originOrder.UserAddress, originOrder.BaseToken, statedb)
	case chancoinx_state.Bid:
		tokenBalance = chancoinx_state.GetTokenBalance(originOrder.UserAddress, originOrder.QuoteToken, statedb)
	default:
		log.Debug("Not found order side", "Side", originOrder.Side)
		return nil, false
	}
	log.Debug("ProcessCancelOrder", "baseToken", originOrder.BaseToken, "quoteToken", originOrder.QuoteToken)
	feeRate := chancoinx_state.GetExRelayerFee(originOrder.ExchangeAddress, statedb)
	tokenCancelFee := getCancelFee(baseTokenDecimal, feeRate, &originOrder)
	if tokenBalance.Cmp(tokenCancelFee) < 0 {
		log.Debug("User not enough balance when cancel order", "Side", originOrder.Side, "balance", tokenBalance, "fee", tokenCancelFee)
		return nil, true
	}

	err = chancoinXstatedb.CancelOrder(orderBook, order)
	if err != nil {
		log.Debug("Error when cancel order", "order", order)
		return err, false
	}
	// relayers pay CHANCOIN for masternode
	chancoinx_state.SubRelayerFee(originOrder.ExchangeAddress, common.RelayerCancelFee, statedb)
	masternodeOwner := statedb.GetOwner(coinbase)
	// relayers pay CHANCOIN for masternode
	statedb.AddBalance(masternodeOwner, common.RelayerCancelFee)

	relayerOwner := chancoinx_state.GetRelayerOwner(originOrder.ExchangeAddress, statedb)
	switch originOrder.Side {
	case chancoinx_state.Ask:
		// users pay token (which they have) for relayer
		chancoinx_state.SubTokenBalance(originOrder.UserAddress, tokenCancelFee, originOrder.BaseToken, statedb)
		chancoinx_state.AddTokenBalance(relayerOwner, tokenCancelFee, originOrder.BaseToken, statedb)
	case chancoinx_state.Bid:
		// users pay token (which they have) for relayer
		chancoinx_state.SubTokenBalance(originOrder.UserAddress, tokenCancelFee, originOrder.QuoteToken, statedb)
		chancoinx_state.AddTokenBalance(relayerOwner, tokenCancelFee, originOrder.QuoteToken, statedb)
	default:
	}
	return nil, false
}

func getCancelFee(baseTokenDecimal *big.Int, feeRate *big.Int, order *chancoinx_state.OrderItem) *big.Int {
	cancelFee := big.NewInt(0)
	if order.Side == chancoinx_state.Ask {
		// SELL 1 BTC => CHANCOIN ,,
		// order.Quantity =1 && fee rate =2
		// ==> cancel fee = 2/10000
		// order.Quantity already included baseToken decimal
		cancelFee = new(big.Int).Mul(order.Quantity, feeRate)
		cancelFee = new(big.Int).Div(cancelFee, common.ChancoinXBaseCancelFee)
	} else {
		// BUY 1 BTC => CHANCOIN with Price : 10000
		// quoteTokenQuantity = 10000 && fee rate =2
		// => cancel fee =2
		quoteTokenQuantity := new(big.Int).Mul(order.Quantity, order.Price)
		quoteTokenQuantity = quoteTokenQuantity.Div(quoteTokenQuantity, baseTokenDecimal)
		// Fee
		// makerFee = quoteTokenQuantity * feeRate / baseFee = quantityToTrade * makerPrice / baseTokenDecimal * feeRate / baseFee
		cancelFee = new(big.Int).Mul(quoteTokenQuantity, feeRate)
		cancelFee = new(big.Int).Div(quoteTokenQuantity, common.ChancoinXBaseCancelFee)
	}
	return cancelFee
}
