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

	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/sync/syncmap"
)

const (
	ProtocolName         = "tomox"
	ProtocolVersion      = uint64(1)
	ProtocolVersionStr   = "1.0"
	overflowIdx                // Indicator of message queue overflow
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
}

type TxMatchBatch struct {
	Data      []TxDataMatch
	Timestamp uint64
	TxHash    common.Hash
}

// DefaultConfig represents (shocker!) the default configuration.
var DefaultConfig = Config{
	DataDir: "",
}

type TomoX struct {
	// Order related
	db         OrderDao
	Triegc     *prque.Prque         // Priority queue mapping block numbers to tries to gc
	StateCache tomox_state.Database // State database to reuse between imports (contains state cache)    *tomox_state.TomoXStateDB

	orderNonce map[common.Address]*big.Int

	sdkNode     bool
	settings syncmap.Map // holds configuration settings that can be dynamically changed
	tokenDecimalCache *lru.Cache
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
	tomoX := &TomoX{
		orderNonce:          make(map[common.Address]*big.Int),
		Triegc:              prque.New(),
		tokenDecimalCache:   tokenDecimalCache,
	}
	switch cfg.DBEngine {
	case "leveldb":
		tomoX.db = NewLDBEngine(cfg)
		tomoX.sdkNode = false
	case "mongodb":
		tomoX.db = NewMongoDBEngine(cfg)
		tomoX.sdkNode = true
	default:
		log.Crit("wrong database engine, only accept either leveldb or mongodb")
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


		log.Info("Process order pending", "orderPending", order)
		originalOrder := &tomox_state.OrderItem{}
		*originalOrder = *order
		originalOrder.Quantity = CloneBigInt(order.Quantity)

		if cancel {
			order.Status = OrderStatusCancelled
		}
		trades, _, err := ProcessOrder(statedb, tomoXstatedb, common.StringToHash(order.PairName), order)

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
func (tomox *TomoX) SyncDataToSDKNode(txDataMatch TxDataMatch, txHash common.Hash, statedb *state.StateDB) error {
	var (
		order *tomox_state.OrderItem
		err   error
	)
	db := tomox.GetDB()

	// 1. put processed order to db
	if order, err = txDataMatch.DecodeOrder(); err != nil {
		log.Error("SDK node decode order failed", "txDataMatch", txDataMatch)
		return fmt.Errorf("SDK node decode order failed")
	}

	if order.Status != OrderStatusCancelled {
		order.Status = OrderStatusOpen
	}
	order.TxHash = txHash

	log.Debug("PutObject processed order", "order", order)
	if err := db.PutObject(order.Hash.Bytes(), order, false, common.Hash{}); err != nil {
		return fmt.Errorf("SDKNode: failed to put processed order. Error: %s", err.Error())
	}
	if order.Status == OrderStatusCancelled {
		return nil
	}
	order.TxHash = txHash
	// 2. put trades to db and update status to FILLED
	trades := txDataMatch.GetTrades()
	log.Debug("Got trades", "number", len(trades), "trades", trades)
	for _, trade := range trades {
		// 2.a. put to trades
		tradeSDK := &Trade{}
		quantity := ToBigInt(trade[TradeQuantity])
		price := ToBigInt(trade[TradePrice])
		if price.Cmp(big.NewInt(0)) <= 0 || quantity.Cmp(big.NewInt(0)) <= 0 {
			return fmt.Errorf("trade misses important information. tradedPrice %v, tradedQuantity %v", price, quantity)
		}
		tradeSDK.Amount = quantity
		tradeSDK.PricePoint = price
		tradeSDK.PairName = order.PairName
		tradeSDK.BaseToken = order.BaseToken
		tradeSDK.QuoteToken = order.QuoteToken
		tradeSDK.Status = TradeStatusSuccess
		tradeSDK.Taker = order.UserAddress
		tradeSDK.Maker = common.HexToAddress(trade[TradeMaker])
		tradeSDK.TakerOrderHash = order.Hash
		tradeSDK.MakerOrderHash = common.HexToHash(trade[TradeMakerOrderHash])
		tradeSDK.TxHash = txHash
		tradeSDK.TakerOrderSide = order.Side
		tradeSDK.TakerExchange = order.ExchangeAddress
		tradeSDK.MakerExchange = common.HexToAddress(trade[TradeMakerExchange])

		// feeAmount: all fees are calculated in quoteToken
		quoteTokenQuantity := big.NewInt(0).Mul(quantity, price)
		quoteTokenQuantity = big.NewInt(0).Div(quoteTokenQuantity, common.BasePrice)
		takerFee := big.NewInt(0).Mul(quoteTokenQuantity, tomox_state.GetExRelayerFee(order.ExchangeAddress, statedb))
		takerFee = big.NewInt(0).Div(takerFee, common.TomoXBaseFee)
		tradeSDK.TakeFee = takerFee

		makerFee := big.NewInt(0).Mul(quoteTokenQuantity, tomox_state.GetExRelayerFee(common.HexToAddress(trade[TradeMakerExchange]), statedb))
		makerFee = big.NewInt(0).Div(makerFee, common.TomoXBaseFee)
		tradeSDK.MakeFee = makerFee

		tradeSDK.Hash = tradeSDK.ComputeHash()
		log.Debug("TRADE history", "order", order, "trade", tradeSDK)
		if err := db.PutObject(EmptyKey(), tradeSDK, false, common.Hash{}); err != nil {
			return fmt.Errorf("SDKNode: failed to store tradeSDK %s", err.Error())
		}

		// 2.b. update status and filledAmount
		filledAmount := quantity
		// update order status of relating orders
		if err := tomox.updateMatchedOrder(trade[TradeMakerOrderHash], filledAmount); err != nil {
			return err
		}
		if err := tomox.updateMatchedOrder(trade[TradeTakerOrderHash], filledAmount); err != nil {
			return err
		}
	}
	return nil
}

func (tomox *TomoX) updateMatchedOrder(hashString string, filledAmount *big.Int) error {
	log.Debug("updateMatchedOrder", "hash", hashString, "filledAmount", filledAmount)
	db := tomox.GetDB()
	orderHashBytes, err := hex.DecodeString(hashString)
	if err != nil {
		return fmt.Errorf("SDKNode: failed to decode orderKey. Key: %s", hashString)
	}
	val, err := db.GetObject(orderHashBytes, &tomox_state.OrderItem{}, false, common.Hash{})
	if err != nil || val == nil {
		return fmt.Errorf("SDKNode: failed to get order. Key: %s", hashString)
	}
	matchedOrder := val.(*tomox_state.OrderItem)
	updatedFillAmount := new(big.Int)
	updatedFillAmount.Add(matchedOrder.FilledAmount, filledAmount)
	matchedOrder.FilledAmount = updatedFillAmount
	if matchedOrder.FilledAmount.Cmp(matchedOrder.Quantity) < 0 {
		matchedOrder.Status = OrderStatusPartialFilled
	} else {
		matchedOrder.Status = OrderStatusFilled
	}
	if err = db.PutObject(matchedOrder.Hash.Bytes(), matchedOrder, false, common.Hash{}); err != nil {
		return fmt.Errorf("SDKNode: failed to update matchedOrder to sdkNode %s", err.Error())
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
