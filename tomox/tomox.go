package tomox

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/tomox/tomox_state"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/sync/syncmap"
)

const (
	ProtocolName       = "tomox"
	ProtocolVersion    = uint64(1)
	ProtocolVersionStr = "1.0"
	overflowIdx         // Indicator of message queue overflow
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

type TxDataMatch struct {
	Order  []byte // serialized data of order has been processed in this tx
	Trades []map[string]string
	RejectedOders []*tomox_state.OrderItem
}

type TxMatchBatch struct {
	Data      []TxDataMatch
	Timestamp int64
	TxHash    common.Hash
}

// DefaultConfig represents (shocker!) the default configuration.
var DefaultConfig = Config{
	DataDir: "",
}

type TomoX struct {
	// Order related
	db         OrderDao
	mongodb    OrderDao
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

func NewLDBEngine(cfg *Config) *BatchDatabase {
	datadir := cfg.DataDir
	batchDB := NewBatchDatabaseWithEncode(datadir, 0)
	return batchDB
}

func NewMongoDBEngine(cfg *Config) *MongoDatabase {
	mongoDB, err := NewMongoDatabase(nil, cfg.DBName, cfg.ConnectionUrl, cfg.ReplicaSetName, 0)

	if err != nil {
		log.Crit("Failed to init mongodb engine", "err", err)
	}

	return mongoDB
}

func New(cfg *Config) *TomoX {
	tokenDecimalCache, _ := lru.New(defaultCacheLimit)
	orderCache, _ := lru.New(defaultCacheLimit)
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

func (tomox *TomoX) GetDB() OrderDao {
	return tomox.db
}

func (tomox *TomoX) GetMongoDB() OrderDao {
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

func (tomox *TomoX) ProcessOrderPending(pending map[common.Address]types.OrderTransactions, statedb *state.StateDB, tomoXstatedb *tomox_state.TomoXStateDB) []TxDataMatch {
	txMatches := []TxDataMatch{}
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
		if order.Status == OrderStatusCancelled {
			cancel = true
		}

		log.Info("Process order pending", "orderPending", order, "BaseToken", order.BaseToken.Hex(), "QuoteToken", order.QuoteToken)
		originalOrder := &tomox_state.OrderItem{}
		*originalOrder = *order
		originalOrder.Quantity = CloneBigInt(order.Quantity)

		if cancel {
			order.Status = OrderStatusCancelled
		}
		trades, rejectes, err := ProcessOrder(statedb, tomoXstatedb, GetOrderBookHash(order.BaseToken, order.QuoteToken), order)

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
		originalOrderValue, err := EncodeBytesItem(originalOrder)
		if err != nil {
			log.Error("Can't encode", "order", originalOrder, "err", err)
			continue
		}
		txMatch := TxDataMatch{
			Order:  originalOrderValue,
			Trades: trades,
			RejectedOders: rejectes,
		}
		txMatches = append(txMatches, txMatch)

	}
	return txMatches
}

// there are 3 tasks need to complete to update data in SDK nodes after matching
// 1. txMatchData.Order: order has been processed. This order should be put to `orders` collection with status sdktypes.OrderStatusOpen
// 2. txMatchData.Trades: includes information of matched orders.
// 		a. PutObject them to `trades` collection
// 		b. Update status of regrading orders to sdktypes.OrderStatusFilled
func (tomox *TomoX) SyncDataToSDKNode(txDataMatch TxDataMatch, txHash common.Hash, txMatchTime time.Time, statedb *state.StateDB) error {
	var (
		// originTakerOrder: order get from db, nil if it doesn't exist
		// takerOrderInTx: order decoded from txdata
		// updatedTakerOrder: order with new status, filledAmount, CreatedAt, UpdatedAt. This will be inserted to db
		originTakerOrder, takerOrderInTx, updatedTakerOrder *tomox_state.OrderItem
		makerDirtyHashes                                    []string
		makerDirtyFilledAmount                              map[string]*big.Int
		err                                                 error
	)
	db := tomox.GetMongoDB()
	sc := db.InitBulk()

	// 1. put processed takerOrderInTx to db
	if takerOrderInTx, err = txDataMatch.DecodeOrder(); err != nil {
		log.Error("SDK node decode takerOrderInTx failed", "txDataMatch", txDataMatch)
		return fmt.Errorf("SDK node decode takerOrderInTx failed")
	}
	lastState := OrderHistoryItem{}
	val, err := db.GetObject(takerOrderInTx.Hash.Bytes(), &tomox_state.OrderItem{})
	if err == nil && val != nil {
		originTakerOrder = val.(*tomox_state.OrderItem)
		lastState = OrderHistoryItem{
			TxHash:       originTakerOrder.TxHash,
			FilledAmount: CloneBigInt(originTakerOrder.FilledAmount),
			Status:       originTakerOrder.Status,
		}
	}
	if originTakerOrder != nil {
		updatedTakerOrder = originTakerOrder
	} else {
		updatedTakerOrder = takerOrderInTx
	}

	if takerOrderInTx.Status != OrderStatusCancelled {
		updatedTakerOrder.Status = OrderStatusOpen
	} else {
		updatedTakerOrder.Status = OrderStatusCancelled
	}
	updatedTakerOrder.TxHash = txHash
	if updatedTakerOrder.CreatedAt.IsZero() {
		updatedTakerOrder.CreatedAt = txMatchTime
	}
	updatedTakerOrder.UpdatedAt = txMatchTime

	tomox.UpdateOrderCache(updatedTakerOrder.PairName, updatedTakerOrder.OrderID, txHash, lastState)

	// 2. put trades to db and update status to FILLED
	trades := txDataMatch.GetTrades()
	log.Debug("Got trades", "number", len(trades), "trades", trades)
	makerDirtyFilledAmount = make(map[string]*big.Int)
	for _, trade := range trades {
		// 2.a. put to trades
		tradeRecord := &Trade{}
		quantity := ToBigInt(trade[TradeQuantity])
		price := ToBigInt(trade[TradePrice])
		if price.Cmp(big.NewInt(0)) <= 0 || quantity.Cmp(big.NewInt(0)) <= 0 {
			return fmt.Errorf("trade misses important information. tradedPrice %v, tradedQuantity %v", price, quantity)
		}
		tradeRecord.Amount = quantity
		tradeRecord.PricePoint = price
		tradeRecord.PairName = updatedTakerOrder.PairName
		tradeRecord.BaseToken = updatedTakerOrder.BaseToken
		tradeRecord.QuoteToken = updatedTakerOrder.QuoteToken
		tradeRecord.Status = TradeStatusSuccess
		tradeRecord.Taker = updatedTakerOrder.UserAddress
		tradeRecord.Maker = common.HexToAddress(trade[TradeMaker])
		tradeRecord.TakerOrderHash = updatedTakerOrder.Hash
		tradeRecord.MakerOrderHash = common.HexToHash(trade[TradeMakerOrderHash])
		tradeRecord.TxHash = txHash
		tradeRecord.TakerOrderSide = updatedTakerOrder.Side
		tradeRecord.TakerExchange = updatedTakerOrder.ExchangeAddress
		tradeRecord.MakerExchange = common.HexToAddress(trade[TradeMakerExchange])

		// feeAmount: all fees are calculated in quoteToken
		quoteTokenQuantity := big.NewInt(0).Mul(quantity, price)
		quoteTokenQuantity = big.NewInt(0).Div(quoteTokenQuantity, common.BasePrice)
		takerFee := big.NewInt(0).Mul(quoteTokenQuantity, tomox_state.GetExRelayerFee(updatedTakerOrder.ExchangeAddress, statedb))
		takerFee = big.NewInt(0).Div(takerFee, common.TomoXBaseFee)
		tradeRecord.TakeFee = takerFee

		makerFee := big.NewInt(0).Mul(quoteTokenQuantity, tomox_state.GetExRelayerFee(common.HexToAddress(trade[TradeMakerExchange]), statedb))
		makerFee = big.NewInt(0).Div(makerFee, common.TomoXBaseFee)
		tradeRecord.MakeFee = makerFee
		if tradeRecord.CreatedAt.IsZero() {
			tradeRecord.CreatedAt = txMatchTime
		}
		tradeRecord.UpdatedAt = txMatchTime
		tradeRecord.Hash = tradeRecord.ComputeHash()

		log.Debug("TRADE history", "trade", tradeRecord)
		if err := db.PutObject(EmptyKey(), tradeRecord); err != nil {
			return fmt.Errorf("SDKNode: failed to store tradeRecord %s", err.Error())
		}

		// 2.b. update status and filledAmount
		filledAmount := quantity
		// maker dirty order
		makerFilledAmount := big.NewInt(0)
		if amount, ok := makerDirtyFilledAmount[trade[TradeMakerOrderHash]]; ok {
			makerFilledAmount = CloneBigInt(amount)
		}
		makerFilledAmount.Add(makerFilledAmount, filledAmount)
		makerDirtyFilledAmount[trade[TradeMakerOrderHash]] = makerFilledAmount
		makerDirtyHashes = append(makerDirtyHashes, trade[TradeMakerOrderHash])

		//updatedTakerOrder = tomox.updateMatchedOrder(updatedTakerOrder, filledAmount, txMatchTime, txHash)
		//  update filledAmount, status of takerOrder
		updatedTakerOrder.FilledAmount.Add(updatedTakerOrder.FilledAmount, filledAmount)
		if updatedTakerOrder.FilledAmount.Cmp(updatedTakerOrder.Quantity) < 0 {
			updatedTakerOrder.Status = OrderStatusPartialFilled
		} else {
			updatedTakerOrder.Status = OrderStatusFilled
		}
	}
	// 3. put rejected orders to db and update status REJECTED
	rejectedOders := txDataMatch.GetRejectedOrders()
	log.Debug("Got rejected orders", "number", len(rejectedOders), "trades", rejectedOders)

	if len(rejectedOders) > 0 {
		// updateRejectedOrders
		for _, rejectedOder := range rejectedOders {
			// update order status of relating orders
			if err := tomox.updateRejectedOrders(rejectedOder.Hash); err != nil {
				return fmt.Errorf("updateRejectedOrders %s", err.Error())
			}
		}
	}
	log.Debug("PutObject processed takerOrder", "takerOrder", ToJSON(updatedTakerOrder))
	if err := db.PutObject(updatedTakerOrder.Hash.Bytes(), updatedTakerOrder); err != nil {
		return fmt.Errorf("SDKNode: failed to put processed takerOrder. Hash: %s Error: %s", updatedTakerOrder.Hash.Hex(), err.Error())
	}
	makerOrders := db.GetListOrderByHashes(makerDirtyHashes)
	for _, o := range makerOrders {
		lastState = OrderHistoryItem{
			TxHash:       o.TxHash,
			FilledAmount: CloneBigInt(o.FilledAmount),
			Status:       o.Status,
		}
		tomox.UpdateOrderCache(o.PairName, o.OrderID, txHash, lastState)
		o.TxHash = txHash
		o.UpdatedAt = txMatchTime
		o.FilledAmount.Add(o.FilledAmount, makerDirtyFilledAmount[o.Hash.Hex()])
		if o.FilledAmount.Cmp(o.Quantity) < 0 {
			o.Status = OrderStatusPartialFilled
		} else {
			o.Status = OrderStatusFilled
		}
		log.Debug("PutObject processed makerOrder", "makerOrder", ToJSON(o))
		if err := db.PutObject(o.Hash.Bytes(), o); err != nil {
			return fmt.Errorf("SDKNode: failed to put processed makerOrder. Hash: %s Error: %s", o.Hash.Hex(), err.Error())
		}
	}
	if err := db.CommitBulk(sc); err != nil {
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

func (tomox *TomoX) UpdateOrderCache(pair string, orderId uint64, txhash common.Hash, lastState OrderHistoryItem) {
	var orderCacheAtTxHash map[common.Hash]OrderHistoryItem
	c, ok := tomox.orderCache.Get(txhash)
	if !ok || c == nil {
		orderCacheAtTxHash = make(map[common.Hash]OrderHistoryItem)
	} else {
		orderCacheAtTxHash = c.(map[common.Hash]OrderHistoryItem)
	}
	orderKey := GetOrderHistoryKey(pair, orderId)
	_, ok = orderCacheAtTxHash[orderKey]
	if !ok {
		orderCacheAtTxHash[orderKey] = lastState
	}
	tomox.orderCache.Add(txhash, orderCacheAtTxHash)
}

func (tomox *TomoX) RollbackReorgTxMatch(txhash common.Hash) {
	db := tomox.GetMongoDB()
	for _, order := range db.GetOrderByTxHash(txhash) {
		c, ok := tomox.orderCache.Get(txhash)
		log.Debug("Tomox reorg: rollback order", "txhash", txhash.Hex(), "order", ToJSON(order), "orderHistoryItem", c)
		if !ok {
			log.Debug("Tomox reorg: remove order due to no orderCache", "order", ToJSON(order))
			if err := db.DeleteObject(order.Hash.Bytes()); err != nil {
				log.Error("SDKNode: failed to remove reorg order", "err", err.Error(), "order", ToJSON(order))
			}
			continue
		}
		orderCacheAtTxHash := c.(map[common.Hash]OrderHistoryItem)
		orderHistoryItem, _ := orderCacheAtTxHash[GetOrderHistoryKey(order.PairName, order.OrderID)]
		if (orderHistoryItem == OrderHistoryItem{}) {
			log.Debug("Tomox reorg: remove order due to empty orderHistory", "order", ToJSON(order))
			if err := db.DeleteObject(order.Hash.Bytes()); err != nil {
				log.Error("SDKNode: failed to remove reorg order", "err", err.Error(), "order", ToJSON(order))
			}
			continue
		}
		order.TxHash = orderHistoryItem.TxHash
		order.Status = orderHistoryItem.Status
		order.FilledAmount = CloneBigInt(orderHistoryItem.FilledAmount)
		log.Debug("Tomox reorg: update order to the last orderHistoryItem", "order", ToJSON(order), "orderHistoryItem", orderHistoryItem)
		if err := db.PutObject(order.Hash.Bytes(), order); err != nil {
			log.Error("SDKNode: failed to update reorg order", "err", err.Error(), "order", ToJSON(order))
		}
	}
	log.Debug("Tomox reorg: DeleteTradeByTxHash", "txhash", txhash.Hex())
	db.DeleteTradeByTxHash(txhash)

}
func (tomox *TomoX) updateRejectedOrders(hash common.Hash) error {
	log.Debug("updateRejectedOrders", "hash", hash)
	db := tomox.GetMongoDB()
	val, err := db.GetObject(hash.Bytes(), &tomox_state.OrderItem{})
	if err != nil || val == nil {
		return fmt.Errorf("SDKNode: failed to get order. Key: %s", hash)
	}
	rejectedOder := val.(*tomox_state.OrderItem)
	rejectedOder.Status = OrderStatusRejected

	if err = db.PutObject(rejectedOder.Hash.Bytes(), rejectedOder); err != nil {
		return fmt.Errorf("SDKNode: failed to update rejectedOder to sdkNode %s", err.Error())
	}
	return nil
}
