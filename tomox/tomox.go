package tomox

import (
	"errors"
	"fmt"
	"github.com/tomochain/tomochain/consensus"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/p2p"
	"github.com/tomochain/tomochain/tomox/tomox_state"
	"github.com/tomochain/tomochain/tomoxDAO"
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
	ProtocolName       = "tomox"
	ProtocolVersion    = uint64(1)
	ProtocolVersionStr = "1.0"
	overflowIdx        // Indicator of message queue overflow
	defaultCacheLimit = 1024
)

var (
	ErrNonceTooHigh = errors.New("nonce too high")
	ErrNonceTooLow  = errors.New("nonce too low")
)

type Config struct {
	DataDir        string `toml:",omitempty"`
	DBEngine       string `toml:",omitempty"`
	DBName         string `toml:",omitempty"`
	ConnectionUrl  string `toml:",omitempty"`
	ReplicaSetName string `toml:",omitempty"`
}

// DefaultConfig represents (shocker!) the default configuration.
var DefaultConfig = Config{
	DataDir: "",
}

type TomoX struct {
	// Order related
	db         tomoxDAO.TomoXDAO
	mongodb    tomoxDAO.TomoXDAO
	Triegc     *prque.Prque         // Priority queue mapping block numbers to tries to gc
	StateCache tomox_state.Database // State database to reuse between imports (contains state cache)    *tomox_state.TomoXStateDB

	orderNonce map[common.Address]*big.Int

	sdkNode           bool
	settings          syncmap.Map // holds configuration settings that can be dynamically changed
	tokenDecimalCache *lru.Cache
	orderCache        *lru.Cache
}

func (tomox *TomoX) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

func (tomox *TomoX) Start(server *p2p.Server) error {
	return nil
}

func (tomox *TomoX) Stop() error {
	return nil
}

func NewLDBEngine(cfg *Config) *tomoxDAO.BatchDatabase {
	datadir := cfg.DataDir
	batchDB := tomoxDAO.NewBatchDatabaseWithEncode(datadir, 0)
	return batchDB
}

func NewMongoDBEngine(cfg *Config) *tomoxDAO.MongoDatabase {
	mongoDB, err := tomoxDAO.NewMongoDatabase(nil, cfg.DBName, cfg.ConnectionUrl, cfg.ReplicaSetName, 0)

	if err != nil {
		log.Crit("Failed to init mongodb engine", "err", err)
	}

	return mongoDB
}

func New(cfg *Config) *TomoX {
	tokenDecimalCache, _ := lru.New(defaultCacheLimit)
	orderCache, _ := lru.New(tomox_state.OrderCacheLimit)
	tomoX := &TomoX{
		orderNonce:        make(map[common.Address]*big.Int),
		Triegc:            prque.New(),
		tokenDecimalCache: tokenDecimalCache,
		orderCache:        orderCache,
	}

	// default DBEngine: levelDB
	tomoX.db = NewLDBEngine(cfg)
	tomoX.sdkNode = false

	if cfg.DBEngine == "mongodb" { // this is an add-on DBEngine for SDK nodes
		tomoX.mongodb = NewMongoDBEngine(cfg)
		tomoX.sdkNode = true
	}

	tomoX.StateCache = tomox_state.NewDatabase(tomoX.db)
	tomoX.settings.Store(overflowIdx, false)

	return tomoX
}

// Overflow returns an indication if the message queue is full.
func (tomox *TomoX) Overflow() bool {
	val, _ := tomox.settings.Load(overflowIdx)
	return val.(bool)
}

func (tomox *TomoX) IsSDKNode() bool {
	return tomox.sdkNode
}

func (tomox *TomoX) GetDB() tomoxDAO.TomoXDAO {
	return tomox.db
}

func (tomox *TomoX) GetMongoDB() tomoxDAO.TomoXDAO {
	return tomox.mongodb
}

// APIs returns the RPC descriptors the TomoX implementation offers
func (tomox *TomoX) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: ProtocolName,
			Version:   ProtocolVersionStr,
			Service:   NewPublicTomoXAPI(tomox),
			Public:    true,
		},
	}
}

// Version returns the TomoX sub-protocols version number.
func (tomox *TomoX) Version() uint64 {
	return ProtocolVersion
}

func (tomox *TomoX) ProcessOrderPending(coinbase common.Address, chain consensus.ChainContext, pending map[common.Address]types.OrderTransactions, statedb *state.StateDB, tomoXstatedb *tomox_state.TomoXStateDB) ([]tomox_state.TxDataMatch, map[common.Hash]tomox_state.MatchingResult) {
	txMatches := []tomox_state.TxDataMatch{}
	matchingResults := map[common.Hash]tomox_state.MatchingResult{}

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

		order := &tomox_state.OrderItem{
			Nonce:           big.NewInt(int64(tx.Nonce())),
			Quantity:        tx.Quantity(),
			Price:           tx.Price(),
			ExchangeAddress: tx.ExchangeAddress(),
			UserAddress:     tx.UserAddress(),
			BaseToken:       tx.BaseToken(),
			QuoteToken:      tx.QuoteToken(),
			Status:          tx.Status(),
			Side:            tx.Side(),
			Type:            tx.Type(),
			Hash:            tx.OrderHash(),
			OrderID:         tx.OrderID(),
			Signature: &tomox_state.Signature{
				V: byte(n),
				R: common.BigToHash(R),
				S: common.BigToHash(S),
			},
			PairName: tx.PairName(),
		}
		cancel := false
		if order.Status == tomox_state.OrderStatusCancelled {
			cancel = true
		}

		log.Info("Process order pending", "orderPending", order, "BaseToken", order.BaseToken.Hex(), "QuoteToken", order.QuoteToken)
		originalOrder := &tomox_state.OrderItem{}
		*originalOrder = *order
		originalOrder.Quantity = tomox_state.CloneBigInt(order.Quantity)

		if cancel {
			order.Status = tomox_state.OrderStatusCancelled
		}

		newTrades, newRejectedOrders, err := tomox.CommitOrder(coinbase, chain, statedb, tomoXstatedb, tomox_state.GetOrderBookHash(order.BaseToken, order.QuoteToken), order)

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
		originalOrder.OrderID = order.OrderID
		originalOrderValue, err := tomox_state.EncodeBytesItem(originalOrder)
		if err != nil {
			log.Error("Can't encode", "order", originalOrder, "err", err)
			continue
		}
		txMatch := tomox_state.TxDataMatch{
			Order: originalOrderValue,
		}
		txMatches = append(txMatches, txMatch)
		matchingResults[order.Hash] = tomox_state.MatchingResult{
			Trades:  newTrades,
			Rejects: newRejectedOrders,
		}
	}
	return txMatches, matchingResults
}

// there are 3 tasks need to complete to update data in SDK nodes after matching
// 1. txMatchData.Order: order has been processed. This order should be put to `orders` collection with status sdktypes.OrderStatusOpen
// 2. txMatchData.Trades: includes information of matched orders.
// 		a. PutObject them to `trades` collection
// 		b. Update status of regrading orders to sdktypes.OrderStatusFilled
func (tomox *TomoX) SyncDataToSDKNode(takerOrderInTx *tomox_state.OrderItem, txHash common.Hash, txMatchTime time.Time, statedb *state.StateDB, trades []map[string]string, rejectedOrders []*tomox_state.OrderItem, dirtyOrderCount *uint64) error {
	var (
		// originTakerOrder: order get from db, nil if it doesn't exist
		// takerOrderInTx: order decoded from txdata
		// updatedTakerOrder: order with new status, filledAmount, CreatedAt, UpdatedAt. This will be inserted to db
		originTakerOrder, updatedTakerOrder *tomox_state.OrderItem
		makerDirtyHashes                    []string
		makerDirtyFilledAmount              map[string]*big.Int
		err                                 error
	)
	db := tomox.GetMongoDB()
	sc := db.InitBulk()
	defer sc.Close()
	// 1. put processed takerOrderInTx to db
	lastState := tomox_state.OrderHistoryItem{}
	val, err := db.GetObject(takerOrderInTx.Hash, &tomox_state.OrderItem{})
	if err == nil && val != nil {
		originTakerOrder = val.(*tomox_state.OrderItem)
		lastState = tomox_state.OrderHistoryItem{
			TxHash:       originTakerOrder.TxHash,
			FilledAmount: tomox_state.CloneBigInt(originTakerOrder.FilledAmount),
			Status:       originTakerOrder.Status,
			UpdatedAt:    originTakerOrder.UpdatedAt,
		}
	}
	if originTakerOrder != nil {
		updatedTakerOrder = originTakerOrder
	} else {
		updatedTakerOrder = takerOrderInTx
	}

	if takerOrderInTx.Status != tomox_state.OrderStatusCancelled {
		updatedTakerOrder.Status = tomox_state.OrderStatusOpen
	} else {
		updatedTakerOrder.Status = tomox_state.OrderStatusCancelled
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

	tomox.UpdateOrderCache(updatedTakerOrder.BaseToken, updatedTakerOrder.QuoteToken, updatedTakerOrder.Hash, txHash, lastState)
	updatedTakerOrder.UpdatedAt = txMatchTime

	// 2. put trades to db and update status to FILLED
	log.Debug("Got trades", "number", len(trades), "txhash", txHash.Hex())
	makerDirtyFilledAmount = make(map[string]*big.Int)
	for _, trade := range trades {
		// 2.a. put to trades
		tradeRecord := &tomox_state.Trade{}
		quantity := tomox_state.ToBigInt(trade[tomox_state.TradeQuantity])
		price := tomox_state.ToBigInt(trade[tomox_state.TradePrice])
		if price.Cmp(big.NewInt(0)) <= 0 || quantity.Cmp(big.NewInt(0)) <= 0 {
			return fmt.Errorf("trade misses important information. tradedPrice %v, tradedQuantity %v", price, quantity)
		}
		tradeRecord.Amount = quantity
		tradeRecord.PricePoint = price
		tradeRecord.PairName = updatedTakerOrder.PairName
		tradeRecord.BaseToken = updatedTakerOrder.BaseToken
		tradeRecord.QuoteToken = updatedTakerOrder.QuoteToken
		tradeRecord.Status = tomox_state.TradeStatusSuccess
		tradeRecord.Taker = updatedTakerOrder.UserAddress
		tradeRecord.Maker = common.HexToAddress(trade[tomox_state.TradeMaker])
		tradeRecord.TakerOrderHash = updatedTakerOrder.Hash
		tradeRecord.MakerOrderHash = common.HexToHash(trade[tomox_state.TradeMakerOrderHash])
		tradeRecord.TxHash = txHash
		tradeRecord.TakerOrderSide = updatedTakerOrder.Side
		tradeRecord.TakerExchange = updatedTakerOrder.ExchangeAddress
		tradeRecord.MakerExchange = common.HexToAddress(trade[tomox_state.TradeMakerExchange])

		// feeAmount: all fees are calculated in quoteToken
		quoteTokenQuantity := big.NewInt(0).Mul(quantity, price)
		quoteTokenQuantity = big.NewInt(0).Div(quoteTokenQuantity, common.BasePrice)
		takerFee := big.NewInt(0).Mul(quoteTokenQuantity, tomox_state.GetExRelayerFee(updatedTakerOrder.ExchangeAddress, statedb))
		takerFee = big.NewInt(0).Div(takerFee, common.TomoXBaseFee)
		tradeRecord.TakeFee = takerFee

		makerFee := big.NewInt(0).Mul(quoteTokenQuantity, tomox_state.GetExRelayerFee(common.HexToAddress(trade[tomox_state.TradeMakerExchange]), statedb))
		makerFee = big.NewInt(0).Div(makerFee, common.TomoXBaseFee)
		tradeRecord.MakeFee = makerFee

		// set makerOrderType, takerOrderType
		tradeRecord.MakerOrderType = trade[tomox_state.MakerOrderType]
		tradeRecord.TakerOrderType = updatedTakerOrder.Type

		if tradeRecord.CreatedAt.IsZero() {
			tradeRecord.CreatedAt = txMatchTime
		}
		tradeRecord.UpdatedAt = txMatchTime
		tradeRecord.Hash = tradeRecord.ComputeHash()

		log.Debug("TRADE history", "pairName", tradeRecord.PairName, "amount", tradeRecord.Amount, "pricepoint", tradeRecord.PricePoint,
			"taker", tradeRecord.Taker.Hex(), "maker", tradeRecord.Maker.Hex(), "takerOrder", tradeRecord.TakerOrderHash.Hex(), "makerOrder", tradeRecord.MakerOrderHash.Hex(),
			"takerFee", tradeRecord.TakeFee, "makerFee", tradeRecord.MakeFee)
		if err := db.PutObject(tradeRecord.Hash, tradeRecord); err != nil {
			return fmt.Errorf("SDKNode: failed to store tradeRecord %s", err.Error())
		}

		// 2.b. update status and filledAmount
		filledAmount := quantity
		// maker dirty order
		makerFilledAmount := big.NewInt(0)
		if amount, ok := makerDirtyFilledAmount[trade[tomox_state.TradeMakerOrderHash]]; ok {
			makerFilledAmount = tomox_state.CloneBigInt(amount)
		}
		makerFilledAmount.Add(makerFilledAmount, filledAmount)
		makerDirtyFilledAmount[trade[tomox_state.TradeMakerOrderHash]] = makerFilledAmount
		makerDirtyHashes = append(makerDirtyHashes, trade[tomox_state.TradeMakerOrderHash])

		//updatedTakerOrder = tomox.updateMatchedOrder(updatedTakerOrder, filledAmount, txMatchTime, txHash)
		//  update filledAmount, status of takerOrder
		updatedTakerOrder.FilledAmount.Add(updatedTakerOrder.FilledAmount, filledAmount)
		if updatedTakerOrder.FilledAmount.Cmp(updatedTakerOrder.Quantity) < 0 && updatedTakerOrder.Type == tomox_state.Limit {
			updatedTakerOrder.Status = tomox_state.OrderStatusPartialFilled
		} else {
			updatedTakerOrder.Status = tomox_state.OrderStatusFilled
		}
	}

	// update status for Market orders
	if updatedTakerOrder.Type == tomox_state.Market {
		if updatedTakerOrder.FilledAmount.Cmp(big.NewInt(0)) > 0 {
			updatedTakerOrder.Status = tomox_state.OrderStatusFilled
		} else {
			updatedTakerOrder.Status = tomox_state.OrderStatusRejected
		}
	}
	log.Debug("PutObject processed takerOrder",
		"pairName", updatedTakerOrder.PairName, "userAddr", updatedTakerOrder.UserAddress.Hex(), "side", updatedTakerOrder.Side,
		"price", updatedTakerOrder.Price, "quantity", updatedTakerOrder.Quantity, "filledAmount", updatedTakerOrder.FilledAmount, "status", updatedTakerOrder.Status,
		"hash", updatedTakerOrder.Hash.Hex(), "txHash", updatedTakerOrder.TxHash.Hex())
	if err := db.PutObject(updatedTakerOrder.Hash, updatedTakerOrder); err != nil {
		return fmt.Errorf("SDKNode: failed to put processed takerOrder. Hash: %s Error: %s", updatedTakerOrder.Hash.Hex(), err.Error())
	}
	makerOrders := db.GetListOrderByHashes(makerDirtyHashes)
	log.Debug("Maker dirty orders", "len", len(makerOrders), "txhash", txHash.Hex())
	for _, o := range makerOrders {
		if txMatchTime.Before(o.UpdatedAt) {
			log.Debug("Ignore old orders/trades maker", "txHash", txHash.Hex(), "txTime", txMatchTime.UnixNano(), "updatedAt", updatedTakerOrder.UpdatedAt.UnixNano())
			continue
		}
		lastState = tomox_state.OrderHistoryItem{
			TxHash:       o.TxHash,
			FilledAmount: tomox_state.CloneBigInt(o.FilledAmount),
			Status:       o.Status,
			UpdatedAt:    o.UpdatedAt,
		}
		tomox.UpdateOrderCache(o.BaseToken, o.QuoteToken, o.Hash, txHash, lastState)
		o.TxHash = txHash
		o.UpdatedAt = txMatchTime
		o.FilledAmount.Add(o.FilledAmount, makerDirtyFilledAmount[o.Hash.Hex()])
		if o.FilledAmount.Cmp(o.Quantity) < 0 {
			o.Status = tomox_state.OrderStatusPartialFilled
		} else {
			o.Status = tomox_state.OrderStatusFilled
		}
		log.Debug("PutObject processed makerOrder",
			"pairName", o.PairName, "userAddr", o.UserAddress.Hex(), "side", o.Side,
			"price", o.Price, "quantity", o.Quantity, "filledAmount", o.FilledAmount, "status", o.Status,
			"hash", o.Hash.Hex(), "txHash", o.TxHash.Hex())
		if err := db.PutObject(o.Hash, o); err != nil {
			return fmt.Errorf("SDKNode: failed to put processed makerOrder. Hash: %s Error: %s", o.Hash.Hex(), err.Error())
		}
	}

	// 3. put rejected orders to db and update status REJECTED
	log.Debug("Got rejected orders", "number", len(rejectedOrders), "rejectedOrders", rejectedOrders)

	if len(rejectedOrders) > 0 {
		var rejectedHashes []string
		// updateRejectedOrders
		for _, rejectedOrder := range rejectedOrders {
			rejectedHashes = append(rejectedHashes, rejectedOrder.Hash.Hex())
			if updatedTakerOrder.Hash == rejectedOrder.Hash && !txMatchTime.Before(updatedTakerOrder.UpdatedAt) {
				// cache order history for handling reorg
				orderHistoryRecord := tomox_state.OrderHistoryItem{
					TxHash:       updatedTakerOrder.TxHash,
					FilledAmount: tomox_state.CloneBigInt(updatedTakerOrder.FilledAmount),
					Status:       updatedTakerOrder.Status,
					UpdatedAt:    updatedTakerOrder.UpdatedAt,
				}
				tomox.UpdateOrderCache(updatedTakerOrder.BaseToken, updatedTakerOrder.QuoteToken, updatedTakerOrder.Hash, txHash, orderHistoryRecord)

				updatedTakerOrder.Status = tomox_state.OrderStatusRejected
				updatedTakerOrder.TxHash = txHash
				updatedTakerOrder.UpdatedAt = txMatchTime
				if err := db.PutObject(updatedTakerOrder.Hash, updatedTakerOrder); err != nil {
					return fmt.Errorf("SDKNode: failed to reject takerOrder. Hash: %s Error: %s", updatedTakerOrder.Hash.Hex(), err.Error())
				}
			}
		}
		dirtyRejectedOrders := db.GetListOrderByHashes(rejectedHashes)
		for _, order := range dirtyRejectedOrders {
			if txMatchTime.Before(order.UpdatedAt) {
				log.Debug("Ignore old orders/trades reject", "txHash", txHash.Hex(), "txTime", txMatchTime.UnixNano(), "updatedAt", updatedTakerOrder.UpdatedAt.UnixNano())
				continue
			}
			// cache order history for handling reorg
			orderHistoryRecord := tomox_state.OrderHistoryItem{
				TxHash:       order.TxHash,
				FilledAmount: tomox_state.CloneBigInt(order.FilledAmount),
				Status:       order.Status,
				UpdatedAt:    order.UpdatedAt,
			}
			tomox.UpdateOrderCache(order.BaseToken, order.QuoteToken, order.Hash, txHash, orderHistoryRecord)
			dirtyFilledAmount, ok := makerDirtyFilledAmount[order.Hash.Hex()]
			if ok && dirtyFilledAmount != nil {
				order.FilledAmount.Add(order.FilledAmount, dirtyFilledAmount)
			}
			order.Status = tomox_state.OrderStatusRejected
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

func (tomox *TomoX) GetTomoxState(block *types.Block) (*tomox_state.TomoXStateDB, error) {
	root, err := tomox.GetTomoxStateRoot(block)
	if err != nil {
		return nil, err
	}
	if tomox.StateCache == nil {
		return nil, errors.New("Not initialized tomox")
	}
	return tomox_state.New(root, tomox.StateCache)
}

func (tomox *TomoX) GetStateCache() tomox_state.Database {
	return tomox.StateCache
}

func (tomox *TomoX) GetTriegc() *prque.Prque {
	return tomox.Triegc
}

func (tomox *TomoX) GetTomoxStateRoot(block *types.Block) (common.Hash, error) {
	for _, tx := range block.Transactions() {
		if tx.To() != nil && tx.To().Hex() == common.TomoXStateAddr {
			if len(tx.Data()) > 0 {
				return common.BytesToHash(tx.Data()), nil
			}
		}
	}
	return tomox_state.EmptyRoot, nil
}

func (tomox *TomoX) UpdateOrderCache(baseToken, quoteToken common.Address, orderHash common.Hash, txhash common.Hash, lastState tomox_state.OrderHistoryItem) {
	var orderCacheAtTxHash map[common.Hash]tomox_state.OrderHistoryItem
	c, ok := tomox.orderCache.Get(txhash)
	if !ok || c == nil {
		orderCacheAtTxHash = make(map[common.Hash]tomox_state.OrderHistoryItem)
	} else {
		orderCacheAtTxHash = c.(map[common.Hash]tomox_state.OrderHistoryItem)
	}
	orderKey := tomox_state.GetOrderHistoryKey(baseToken, quoteToken, orderHash)
	_, ok = orderCacheAtTxHash[orderKey]
	if !ok {
		orderCacheAtTxHash[orderKey] = lastState
	}
	tomox.orderCache.Add(txhash, orderCacheAtTxHash)
}

func (tomox *TomoX) RollbackReorgTxMatch(txhash common.Hash) {
	db := tomox.GetMongoDB()
	defer tomox.orderCache.Remove(txhash)

	for _, order := range db.GetOrderByTxHash(txhash) {
		c, ok := tomox.orderCache.Get(txhash)
		log.Debug("Tomox reorg: rollback order", "txhash", txhash.Hex(), "order", tomox_state.ToJSON(order), "orderHistoryItem", c)
		if !ok {
			log.Debug("Tomox reorg: remove order due to no orderCache", "order", tomox_state.ToJSON(order))
			if err := db.DeleteObject(order.Hash, &tomox_state.OrderItem{}); err != nil {
				log.Error("SDKNode: failed to remove reorg order", "err", err.Error(), "order", tomox_state.ToJSON(order))
			}
			continue
		}
		orderCacheAtTxHash := c.(map[common.Hash]tomox_state.OrderHistoryItem)
		orderHistoryItem, _ := orderCacheAtTxHash[tomox_state.GetOrderHistoryKey(order.BaseToken, order.QuoteToken, order.Hash)]
		if (orderHistoryItem == tomox_state.OrderHistoryItem{}) {
			log.Debug("Tomox reorg: remove order due to empty orderHistory", "order", tomox_state.ToJSON(order))
			if err := db.DeleteObject(order.Hash, &tomox_state.OrderItem{}); err != nil {
				log.Error("SDKNode: failed to remove reorg order", "err", err.Error(), "order", tomox_state.ToJSON(order))
			}
			continue
		}
		order.TxHash = orderHistoryItem.TxHash
		order.Status = orderHistoryItem.Status
		order.FilledAmount = tomox_state.CloneBigInt(orderHistoryItem.FilledAmount)
		order.UpdatedAt = orderHistoryItem.UpdatedAt
		log.Debug("Tomox reorg: update order to the last orderHistoryItem", "order", tomox_state.ToJSON(order), "orderHistoryItem", orderHistoryItem)
		if err := db.PutObject(order.Hash, order); err != nil {
			log.Error("SDKNode: failed to update reorg order", "err", err.Error(), "order", tomox_state.ToJSON(order))
		}
	}
	log.Debug("Tomox reorg: DeleteTradeByTxHash", "txhash", txhash.Hex())
	db.DeleteTradeByTxHash(txhash)

}
