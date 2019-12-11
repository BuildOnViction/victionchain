package tomoxlending

import (
	"errors"
	"fmt"
	"github.com/tomochain/tomochain/consensus"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/p2p"
	"github.com/tomochain/tomochain/tomox"
	"github.com/tomochain/tomochain/tomoxDAO"
	"github.com/tomochain/tomochain/tomoxlending/lendingstate"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
	"math/big"
	"strconv"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/rpc"
	"golang.org/x/sync/syncmap"
)

const (
	ProtocolName       = "tomoxlending"
	ProtocolVersion    = uint64(1)
	ProtocolVersionStr = "1.0"
	overflowIdx        // Indicator of message queue overflow
)

var (
	ErrNonceTooHigh = errors.New("nonce too high")
	ErrNonceTooLow  = errors.New("nonce too low")
)


type Lending struct {
	// Order related
	leveldb    tomoxDAO.TomoXDAO
	mongodb    tomoxDAO.TomoXDAO
	Triegc     *prque.Prque          // Priority queue mapping block numbers to tries to gc
	StateCache lendingstate.Database // State database to reuse between imports (contains state cache)    *lendingstate.TomoXStateDB

	orderNonce map[common.Address]*big.Int

	tomox *tomox.TomoX
	settings           syncmap.Map // holds configuration settings that can be dynamically changed
	tokenDecimalCache  *lru.Cache
	lendingItemHistory *lru.Cache
}

func (l *Lending) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

func (l *Lending) Start(server *p2p.Server) error {
	return nil
}

func (l *Lending) Stop() error {
	return nil
}


func New(tomox *tomox.TomoX) *Lending {
	tokenDecimalCache, _ := lru.New(defaultCacheLimit)
	orderCache, _ := lru.New(lendingstate.OrderCacheLimit)
	lending := &Lending{
		orderNonce:         make(map[common.Address]*big.Int),
		Triegc:             prque.New(),
		tokenDecimalCache:  tokenDecimalCache,
		lendingItemHistory: orderCache,
	}

	lending.leveldb = tomox.GetDB()

	if tomox.IsSDKNode() { // this is an add-on DBEngine for SDK nodes
		lending.mongodb = tomox.GetMongoDB()
	}

	lending.StateCache = lendingstate.NewDatabase(lending.leveldb)
	lending.settings.Store(overflowIdx, false)

	return lending
}

// Overflow returns an indication if the message queue is full.
func (l *Lending) Overflow() bool {
	val, _ := l.settings.Load(overflowIdx)
	return val.(bool)
}

func (l *Lending) GetDB() tomoxDAO.TomoXDAO {
	return l.leveldb
}

func (l *Lending) GetMongoDB() tomoxDAO.TomoXDAO {
	return l.mongodb
}

// APIs returns the RPC descriptors the Lending implementation offers
func (l *Lending) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: ProtocolName,
			Version:   ProtocolVersionStr,
			Service:   NewPublicTomoXLendingAPI(l),
			Public:    true,
		},
	}
}

// Version returns the Lending sub-protocols version number.
func (l *Lending) Version() uint64 {
	return ProtocolVersion
}

func (l *Lending) ProcessOrderPending(coinbase common.Address, chain consensus.ChainContext, pending map[common.Address]types.OrderTransactions, statedb *state.StateDB, tomoXstatedb *lendingstate.TomoXStateDB) ([]*lendingstate.LendingItem, map[common.Hash]lendingstate.MatchingResult) {
	lendingItems := []*lendingstate.LendingItem{}
	matchingResults := map[common.Hash]lendingstate.MatchingResult{}

	txs := types.NewOrderTransactionByNonce(types.OrderTxSigner{}, pending)
	for {
		tx := txs.Peek()
		if tx == nil {
			break
		}
		log.Debug("ProcessOrderPending start", "len", len(pending))
		log.Debug("Get pending orders to process", "address", tx.UserAddress(), "nonce", tx.Nonce())
		V, R, S := tx.Signature()

		bigstr := V.String()
		n, e := strconv.ParseInt(bigstr, 10, 8)
		if e != nil {
			continue
		}

		order := &lendingstate.LendingItem{
			Nonce:           big.NewInt(int64(tx.Nonce())),
			Quantity:        tx.Quantity(),
			Interest:        tx.Price(),
			Relayer:         tx.ExchangeAddress(),
			UserAddress:     tx.UserAddress(),
			LendingToken:    tx.BaseToken(),
			CollateralToken: tx.QuoteToken(),
			Status:          tx.Status(),
			Side:            tx.Side(),
			Type:            tx.Type(),
			Hash:            tx.OrderHash(),
			LendingId:       tx.OrderID(),
			Signature: &lendingstate.Signature{
				V: byte(n),
				R: common.BigToHash(R),
				S: common.BigToHash(S),
			},
		}
		cancel := false
		if order.Status == lendingstate.LendingStatusCancelled {
			cancel = true
		}

		log.Info("Process order pending", "orderPending", order, "LendingToken", order.LendingToken.Hex(), "CollateralToken", order.CollateralToken)
		originalOrder := &lendingstate.LendingItem{}
		*originalOrder = *order
		originalOrder.Quantity = lendingstate.CloneBigInt(order.Quantity)

		if cancel {
			order.Status = lendingstate.LendingStatusCancelled
		}

		_, newRejectedOrders, err := l.CommitOrder(coinbase, chain, statedb, tomoXstatedb, lendingstate.GetOrderBookHash(order.LendingToken, order.CollateralToken), order)

		for _, reject := range newRejectedOrders {
			log.Debug("Reject order", "reject", *reject)
		}

		switch err {
		case ErrNonceTooLow:
			// New head notification data race between the transaction pool and miner, shift
			log.Debug("Skipping order with low nonce", "sender", tx.UserAddress(), "nonce", tx.Nonce())
			txs.Shift()
			continue

		case ErrNonceTooHigh:
			// Reorg notification data race between the transaction pool and miner, skip account =
			log.Debug("Skipping order account with high nonce", "sender", tx.UserAddress(), "nonce", tx.Nonce())
			txs.Pop()
			continue

		case nil:
			// everything ok
			txs.Shift()

		default:
			// Strange error, discard the transaction and get the next in line (note, the
			// nonce-too-high clause will prevent us from executing in vain).
			log.Debug("Transaction failed, account skipped", "hash", tx.Hash(), "err", err)
			txs.Shift()
			continue
		}

		// orderID has been updated
		originalOrder.LendingId = order.LendingId
		lendingItems = append(lendingItems, originalOrder)
		matchingResults[order.Hash] = lendingstate.MatchingResult{
			//Trades:  newTrades,
			Rejects: newRejectedOrders,
		}
	}
	return lendingItems, matchingResults
}

// there are 3 tasks need to complete to update data in SDK nodes after matching
// 1. txMatchData.Order: order has been processed. This order should be put to `orders` collection with status sdktypes.OrderStatusOpen
// 2. txMatchData.Trades: includes information of matched orders.
// 		a. PutObject them to `trades` collection
// 		b. Update status of regrading orders to sdktypes.OrderStatusFilled
func (l *Lending) SyncDataToSDKNode(takerOrderInTx *lendingstate.LendingItem, txHash common.Hash, txMatchTime time.Time, statedb *state.StateDB, trades []*lendingstate.LendingTrade, rejectedOrders []*lendingstate.LendingItem, dirtyOrderCount *uint64) error {
	var (
		// originTakerOrder: order get from leveldb, nil if it doesn't exist
		// takerOrderInTx: order decoded from txdata
		// updatedTakerOrder: order with new status, filledAmount, CreatedAt, UpdatedAt. This will be inserted to leveldb
		originTakerOrder, updatedTakerOrder *lendingstate.LendingItem
		makerDirtyHashes                    []string
		makerDirtyFilledAmount              map[string]*big.Int
		err                                 error
	)
	db := l.GetMongoDB()
	sc := db.InitBulk()
	defer sc.Close()
	// 1. put processed takerOrderInTx to leveldb
	lastState := lendingstate.LendingItemHistoryItem{}
	val, err := db.GetObject(takerOrderInTx.Hash, &lendingstate.LendingItem{})
	if err == nil && val != nil {
		originTakerOrder = val.(*lendingstate.LendingItem)
		lastState = lendingstate.LendingItemHistoryItem{
			TxHash:       originTakerOrder.TxHash,
			FilledAmount: lendingstate.CloneBigInt(originTakerOrder.FilledAmount),
			Status:       originTakerOrder.Status,
			UpdatedAt:    originTakerOrder.UpdatedAt,
		}
	}
	if originTakerOrder != nil {
		updatedTakerOrder = originTakerOrder
	} else {
		updatedTakerOrder = takerOrderInTx
	}

	if takerOrderInTx.Status != lendingstate.LendingStatusCancelled {
		updatedTakerOrder.Status = lendingstate.LendingStatusOpen
	} else {
		updatedTakerOrder.Status = lendingstate.LendingStatusCancelled
	}
	updatedTakerOrder.TxHash = txHash
	if updatedTakerOrder.CreatedAt.IsZero() {
		updatedTakerOrder.CreatedAt = txMatchTime
	}
	if txMatchTime.Before(updatedTakerOrder.UpdatedAt) || (txMatchTime.Equal(updatedTakerOrder.UpdatedAt) && *dirtyOrderCount == 0) {
		log.Debug("Ignore old orders/trades taker", "txHash", txHash.Hex(), "txTime", txMatchTime.UnixNano(), "updatedAt", updatedTakerOrder.UpdatedAt.UnixNano())
		return nil
	}
	*dirtyOrderCount++

	l.UpdateOrderCache(updatedTakerOrder.LendingToken, updatedTakerOrder.CollateralToken, updatedTakerOrder.Hash, txHash, lastState)
	updatedTakerOrder.UpdatedAt = txMatchTime

	// 2. put trades to leveldb and update status to FILLED
	log.Debug("Got trades", "number", len(trades), "txhash", txHash.Hex())
	makerDirtyFilledAmount = make(map[string]*big.Int)
	for _, trade := range trades {
		// 2.a. put to trades
		tradeRecord := trade

		if tradeRecord.CreatedAt.IsZero() {
			tradeRecord.CreatedAt = txMatchTime
		}
		tradeRecord.UpdatedAt = txMatchTime

		log.Debug("TRADE history", "Term", tradeRecord.Term, "amount", tradeRecord.Amount, "Interest", tradeRecord.Interest,
			"borrower", tradeRecord.Borrower.Hex(), "investor", tradeRecord.Investor.Hex(), "TakerOrderHash", tradeRecord.TakerOrderHash.Hex(), "MakerOrderHash", tradeRecord.MakerOrderHash.Hex(),
			"borrowing", tradeRecord.BorrowingFee.String(), "investingFee", tradeRecord.InvestingFee.String())
		if err := db.PutObject(tradeRecord.Hash, tradeRecord); err != nil {
			return fmt.Errorf("SDKNode: failed to store lendingTrade %s", err.Error())
		}

		// 2.b. update status and filledAmount
		filledAmount := trade.Amount
		// maker dirty order
		makerFilledAmount := big.NewInt(0)
		if amount, ok := makerDirtyFilledAmount[trade.MakerOrderHash.Hex()]; ok {
			makerFilledAmount = lendingstate.CloneBigInt(amount)
		}
		makerFilledAmount.Add(makerFilledAmount, filledAmount)
		makerDirtyFilledAmount[trade.MakerOrderHash.Hex()] = makerFilledAmount
		makerDirtyHashes = append(makerDirtyHashes, trade.MakerOrderHash.Hex())

		//updatedTakerOrder = l.updateMatchedOrder(updatedTakerOrder, filledAmount, txMatchTime, txHash)
		//  update filledAmount, status of takerOrder
		updatedTakerOrder.FilledAmount.Add(updatedTakerOrder.FilledAmount, filledAmount)
		if updatedTakerOrder.FilledAmount.Cmp(updatedTakerOrder.Quantity) < 0 && updatedTakerOrder.Type == lendingstate.Limit {
			updatedTakerOrder.Status = lendingstate.LendingStatusPartialFilled
		} else {
			updatedTakerOrder.Status = lendingstate.LendingStatusFilled
		}
	}

	// update status for Market orders
	if updatedTakerOrder.Type == lendingstate.Market {
		if updatedTakerOrder.FilledAmount.Cmp(big.NewInt(0)) > 0 {
			updatedTakerOrder.Status = lendingstate.LendingStatusFilled
		} else {
			updatedTakerOrder.Status = lendingstate.LendingStatusReject
		}
	}
	log.Debug("PutObject processed takerOrder",
		"term", updatedTakerOrder.Term, "userAddr", updatedTakerOrder.UserAddress.Hex(), "side", updatedTakerOrder.Side,
		"Interest", updatedTakerOrder.Interest, "quantity", updatedTakerOrder.Quantity, "filledAmount", updatedTakerOrder.FilledAmount, "status", updatedTakerOrder.Status,
		"hash", updatedTakerOrder.Hash.Hex(), "txHash", updatedTakerOrder.TxHash.Hex())
	if err := db.PutObject(updatedTakerOrder.Hash, updatedTakerOrder); err != nil {
		return fmt.Errorf("SDKNode: failed to put processed takerOrder. Hash: %s Error: %s", updatedTakerOrder.Hash.Hex(), err.Error())
	}
	makerOrders := db.GetListLendingItemByHashes(makerDirtyHashes)
	log.Debug("Maker dirty orders", "len", len(makerOrders), "txhash", txHash.Hex())
	for _, o := range makerOrders {
		if txMatchTime.Before(o.UpdatedAt) {
			log.Debug("Ignore old orders/trades maker", "txHash", txHash.Hex(), "txTime", txMatchTime.UnixNano(), "updatedAt", updatedTakerOrder.UpdatedAt.UnixNano())
			continue
		}
		lastState = lendingstate.LendingItemHistoryItem{
			TxHash:       o.TxHash,
			FilledAmount: lendingstate.CloneBigInt(o.FilledAmount),
			Status:       o.Status,
			UpdatedAt:    o.UpdatedAt,
		}
		l.UpdateOrderCache(o.LendingToken, o.CollateralToken, o.Hash, txHash, lastState)
		o.TxHash = txHash
		o.UpdatedAt = txMatchTime
		o.FilledAmount.Add(o.FilledAmount, makerDirtyFilledAmount[o.Hash.Hex()])
		if o.FilledAmount.Cmp(o.Quantity) < 0 {
			o.Status = lendingstate.LendingStatusPartialFilled
		} else {
			o.Status = lendingstate.LendingStatusFilled
		}
		log.Debug("PutObject processed makerOrder",
			"term", o.Term, "userAddr", o.UserAddress.Hex(), "side", o.Side,
			"Interest", o.Interest, "quantity", o.Quantity, "filledAmount", o.FilledAmount, "status", o.Status,
			"hash", o.Hash.Hex(), "txHash", o.TxHash.Hex())
		if err := db.PutObject(o.Hash, o); err != nil {
			return fmt.Errorf("SDKNode: failed to put processed makerOrder. Hash: %s Error: %s", o.Hash.Hex(), err.Error())
		}
	}

	// 3. put rejected orders to leveldb and update status REJECTED
	log.Debug("Got rejected orders", "number", len(rejectedOrders), "rejectedOrders", rejectedOrders)

	if len(rejectedOrders) > 0 {
		var rejectedHashes []string
		// updateRejectedOrders
		for _, rejectedOrder := range rejectedOrders {
			rejectedHashes = append(rejectedHashes, rejectedOrder.Hash.Hex())
			if updatedTakerOrder.Hash == rejectedOrder.Hash && !txMatchTime.Before(updatedTakerOrder.UpdatedAt) {
				// cache order history for handling reorg
				orderHistoryRecord := lendingstate.LendingItemHistoryItem{
					TxHash:       updatedTakerOrder.TxHash,
					FilledAmount: lendingstate.CloneBigInt(updatedTakerOrder.FilledAmount),
					Status:       updatedTakerOrder.Status,
					UpdatedAt:    updatedTakerOrder.UpdatedAt,
				}
				l.UpdateOrderCache(updatedTakerOrder.LendingToken, updatedTakerOrder.CollateralToken, updatedTakerOrder.Hash, txHash, orderHistoryRecord)

				updatedTakerOrder.Status = lendingstate.LendingStatusReject
				updatedTakerOrder.TxHash = txHash
				updatedTakerOrder.UpdatedAt = txMatchTime
				if err := db.PutObject(updatedTakerOrder.Hash, updatedTakerOrder); err != nil {
					return fmt.Errorf("SDKNode: failed to reject takerOrder. Hash: %s Error: %s", updatedTakerOrder.Hash.Hex(), err.Error())
				}
			}
		}
		dirtyRejectedOrders := db.GetListLendingItemByHashes(rejectedHashes)
		for _, order := range dirtyRejectedOrders {
			if txMatchTime.Before(order.UpdatedAt) {
				log.Debug("Ignore old orders/trades reject", "txHash", txHash.Hex(), "txTime", txMatchTime.UnixNano(), "updatedAt", updatedTakerOrder.UpdatedAt.UnixNano())
				continue
			}
			// cache order history for handling reorg
			orderHistoryRecord := lendingstate.LendingItemHistoryItem{
				TxHash:       order.TxHash,
				FilledAmount: lendingstate.CloneBigInt(order.FilledAmount),
				Status:       order.Status,
				UpdatedAt:    order.UpdatedAt,
			}
			l.UpdateOrderCache(order.LendingToken, order.CollateralToken, order.Hash, txHash, orderHistoryRecord)
			dirtyFilledAmount, ok := makerDirtyFilledAmount[order.Hash.Hex()]
			if ok && dirtyFilledAmount != nil {
				order.FilledAmount.Add(order.FilledAmount, dirtyFilledAmount)
			}
			order.Status = lendingstate.LendingStatusReject
			order.TxHash = txHash
			order.UpdatedAt = txMatchTime
			if err = db.PutObject(order.Hash, order); err != nil {
				return fmt.Errorf("SDKNode: failed to update rejectedOder to sdkNode %s", err.Error())
			}
		}
	}

	if err := db.CommitBulk(); err != nil {
		return fmt.Errorf("SDKNode fail to commit bulk update orders, trades at txhash %s . Error: %s", txHash.Hex(), err.Error())
	}
	return nil
}

func (l *Lending) GetTomoxState(block *types.Block) (*lendingstate.TomoXStateDB, error) {
	root, err := l.GetTomoxStateRoot(block)
	if err != nil {
		return nil, err
	}
	if l.StateCache == nil {
		return nil, errors.New("Not initialized l")
	}
	return lendingstate.New(root, l.StateCache)
}

func (l *Lending) GetStateCache() lendingstate.Database {
	return l.StateCache
}

func (l *Lending) GetTriegc() *prque.Prque {
	return l.Triegc
}

func (l *Lending) GetTomoxStateRoot(block *types.Block) (common.Hash, error) {
	for _, tx := range block.Transactions() {
		if tx.To() != nil && tx.To().Hex() == common.TomoXStateAddr {
			if len(tx.Data()) > 0 {
				return common.BytesToHash(tx.Data()), nil
			}
		}
	}
	return lendingstate.EmptyRoot, nil
}

func (l *Lending) UpdateOrderCache(LendingToken, CollateralToken common.Address, orderHash common.Hash, txhash common.Hash, lastState lendingstate.LendingItemHistoryItem) {
	var orderCacheAtTxHash map[common.Hash]lendingstate.LendingItemHistoryItem
	c, ok := l.lendingItemHistory.Get(txhash)
	if !ok || c == nil {
		orderCacheAtTxHash = make(map[common.Hash]lendingstate.LendingItemHistoryItem)
	} else {
		orderCacheAtTxHash = c.(map[common.Hash]lendingstate.LendingItemHistoryItem)
	}
	orderKey := lendingstate.GetLendingItemHistoryKey(LendingToken, CollateralToken, orderHash)
	_, ok = orderCacheAtTxHash[orderKey]
	if !ok {
		orderCacheAtTxHash[orderKey] = lastState
	}
	l.lendingItemHistory.Add(txhash, orderCacheAtTxHash)
}

func (l *Lending) RollbackReorgTxMatch(txhash common.Hash) {
	db := l.GetMongoDB()
	defer l.lendingItemHistory.Remove(txhash)

	for _, order := range db.GetLendingItemByTxHash(txhash) {
		c, ok := l.lendingItemHistory.Get(txhash)
		log.Debug("Tomox reorg: rollback order", "txhash", txhash.Hex(), "order", lendingstate.ToJSON(order), "orderHistoryItem", c)
		if !ok {
			log.Debug("Tomox reorg: remove order due to no lendingItemHistory", "order", lendingstate.ToJSON(order))
			if err := db.DeleteObject(order.Hash, &lendingstate.LendingItem{}); err != nil {
				log.Error("SDKNode: failed to remove reorg order", "err", err.Error(), "order", lendingstate.ToJSON(order))
			}
			continue
		}
		orderCacheAtTxHash := c.(map[common.Hash]lendingstate.LendingItemHistoryItem)
		orderHistoryItem, _ := orderCacheAtTxHash[lendingstate.GetLendingItemHistoryKey(order.LendingToken, order.CollateralToken, order.Hash)]
		if (orderHistoryItem == lendingstate.LendingItemHistoryItem{}) {
			log.Debug("Tomox reorg: remove order due to empty orderHistory", "order", lendingstate.ToJSON(order))
			if err := db.DeleteObject(order.Hash, &lendingstate.LendingItem{}); err != nil {
				log.Error("SDKNode: failed to remove reorg order", "err", err.Error(), "order", lendingstate.ToJSON(order))
			}
			continue
		}
		order.TxHash = orderHistoryItem.TxHash
		order.Status = orderHistoryItem.Status
		order.FilledAmount = lendingstate.CloneBigInt(orderHistoryItem.FilledAmount)
		order.UpdatedAt = orderHistoryItem.UpdatedAt
		log.Debug("Tomox reorg: update order to the last orderHistoryItem", "order", lendingstate.ToJSON(order), "orderHistoryItem", orderHistoryItem)
		if err := db.PutObject(order.Hash, order); err != nil {
			log.Error("SDKNode: failed to update reorg order", "err", err.Error(), "order", lendingstate.ToJSON(order))
		}
	}
	log.Debug("Tomox reorg: DeleteTradeByTxHash", "txhash", txhash.Hex())
	db.DeleteTradeByTxHash(txhash)

}
