package interfaces

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	swapBitcoin "github.com/tomochain/tomox-sdk/swap/bitcoin"
	swapEthereum "github.com/tomochain/tomox-sdk/swap/ethereum"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/ws"
	"math/big"
)

type OrderDao interface {
	GetCollection() *mgo.Collection
	Create(o *types.Order) error
	Update(id bson.ObjectId, o *types.Order) error
	Upsert(id bson.ObjectId, o *types.Order) error
	Delete(orders ...*types.Order) error
	DeleteByHashes(hashes ...common.Hash) error
	UpdateAllByHash(h common.Hash, o *types.Order) error
	UpdateByHash(h common.Hash, o *types.Order) error
	UpsertByHash(h common.Hash, o *types.Order) error
	GetOrderCountByUserAddress(addr common.Address) (int, error)
	GetByID(id bson.ObjectId) (*types.Order, error)
	GetByHash(h common.Hash) (*types.Order, error)
	GetByHashes(hashes []common.Hash) ([]*types.Order, error)
	GetByUserAddress(addr, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error)
	GetOpenOrdersByUserAddress(addr common.Address) ([]*types.Order, error)
	GetCurrentByUserAddress(a common.Address, limit ...int) ([]*types.Order, error)
	GetHistoryByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error)
	GetMatchingBuyOrders(o *types.Order) ([]*types.Order, error)
	GetMatchingSellOrders(o *types.Order) ([]*types.Order, error)
	UpdateOrderFilledAmount(h common.Hash, value *big.Int) error
	UpdateOrderFilledAmounts(h []common.Hash, values []*big.Int) ([]*types.Order, error)
	UpdateOrderStatusesByHashes(status string, hashes ...common.Hash) ([]*types.Order, error)
	GetUserLockedBalance(account common.Address, token common.Address, p *types.Pair) (*big.Int, error)
	UpdateOrderStatus(h common.Hash, status string) error
	GetRawOrderBook(*types.Pair) ([]*types.Order, error)
	GetOrderBook(*types.Pair) ([]map[string]string, []map[string]string, error)
	GetSideOrderBook(p *types.Pair, side string, sort int, limit ...int) ([]map[string]string, error)
	GetOrderBookPricePoint(p *types.Pair, pp *big.Int, side string) (*big.Int, error)
	FindAndModify(h common.Hash, o *types.Order) (*types.Order, error)
	Drop() error
	Aggregate(q []bson.M) ([]*types.OrderData, error)
	AddNewOrder(o *types.Order, topic string) error
	CancelOrder(o *types.Order, topic string) error
	AddTopic(t []string) (string, error)
	DeleteTopic(t string) error
}

type StopOrderDao interface {
	Create(so *types.StopOrder) error
	Update(id bson.ObjectId, so *types.StopOrder) error
	UpdateByHash(h common.Hash, so *types.StopOrder) error
	Upsert(id bson.ObjectId, so *types.StopOrder) error
	UpsertByHash(h common.Hash, so *types.StopOrder) error
	UpdateAllByHash(h common.Hash, so *types.StopOrder) error
	GetByHash(h common.Hash) (*types.StopOrder, error)
	FindAndModify(h common.Hash, so *types.StopOrder) (*types.StopOrder, error)
	GetTriggeredStopOrders(baseToken, quoteToken common.Address, lastPrice *big.Int) ([]*types.StopOrder, error)
	Drop() error
}

type AccountDao interface {
	Create(account *types.Account) (err error)
	GetAll() (res []types.Account, err error)
	GetByID(id bson.ObjectId) (*types.Account, error)
	GetByAddress(owner common.Address) (response *types.Account, err error)
	GetTokenBalances(owner common.Address) (map[common.Address]*types.TokenBalance, error)
	GetTokenBalance(owner common.Address, token common.Address) (*types.TokenBalance, error)
	UpdateTokenBalance(owner common.Address, token common.Address, tokenBalance *types.TokenBalance) (err error)
	UpdateBalance(owner common.Address, token common.Address, balance *big.Int) (err error)
	FindOrCreate(addr common.Address) (*types.Account, error)
	Transfer(token common.Address, fromAddress common.Address, toAddress common.Address, amount *big.Int) error
	Drop()
	GetFavoriteTokens(owner common.Address) (map[common.Address]bool, error)
	AddFavoriteToken(owner, token common.Address) error
	DeleteFavoriteToken(owner, token common.Address) error
}

type ConfigDao interface {
	GetSchemaVersion() uint64
	GetAddressIndex(chain types.Chain) (uint64, error)
	IncrementAddressIndex(chain types.Chain) error
	ResetBlockCounters() error
	GetBlockToProcess(chain types.Chain) (uint64, error)
	SaveLastProcessedBlock(chain types.Chain, block uint64) error
	Drop()
}

type AssociationDao interface {
	GetAssociationByChainAddress(chain types.Chain, address common.Address) (*types.AddressAssociationRecord, error)
	GetAssociationByChainAssociatedAddress(chain types.Chain, associatedAddress common.Address) (*types.AddressAssociationRecord, error)

	// save mean if there is no item then insert, otherwise update
	SaveAssociation(record *types.AddressAssociationRecord) error
	SaveDepositTransaction(chain types.Chain, sourceAccount common.Address, txEnvelope string) error
	SaveAssociationStatus(chain types.Chain, sourceAccount common.Address, status string) error
}

type WalletDao interface {
	Create(wallet *types.Wallet) error
	GetAll() ([]types.Wallet, error)
	GetByID(id bson.ObjectId) (*types.Wallet, error)
	GetByAddress(addr common.Address) (*types.Wallet, error)
	GetDefaultAdminWallet() (*types.Wallet, error)
	GetOperatorWallets() ([]*types.Wallet, error)
}

type PairDao interface {
	Create(o *types.Pair) error
	GetAll() ([]types.Pair, error)
	GetActivePairs() ([]*types.Pair, error)
	GetByID(id bson.ObjectId) (*types.Pair, error)
	GetByName(name string) (*types.Pair, error)
	GetByTokenSymbols(baseTokenSymbol, quoteTokenSymbol string) (*types.Pair, error)
	GetByTokenAddress(baseToken, quoteToken common.Address) (*types.Pair, error)
	GetListedPairs() ([]types.Pair, error)
	GetUnlistedPairs() ([]types.Pair, error)
}

type TradeDao interface {
	GetCollection() *mgo.Collection
	Create(o ...*types.Trade) error
	Update(t *types.Trade) error
	UpdateByHash(h common.Hash, t *types.Trade) error
	GetAll() ([]types.Trade, error)
	Aggregate(q []bson.M) ([]*types.Tick, error)
	GetByPairName(name string) ([]*types.Trade, error)
	GetByHash(h common.Hash) (*types.Trade, error)
	GetByMakerOrderHash(h common.Hash) ([]*types.Trade, error)
	GetByTakerOrderHash(h common.Hash) ([]*types.Trade, error)
	GetByOrderHashes(hashes []common.Hash) ([]*types.Trade, error)
	GetSortedTrades(bt, qt common.Address, from, to int64, n int) ([]*types.Trade, error)
	GetSortedTradesByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Trade, error)
	GetNTradesByPairAddress(bt, qt common.Address, n int) ([]*types.Trade, error)
	GetTradesByPairAddress(bt, qt common.Address, n int) ([]*types.Trade, error)
	GetAllTradesByPairAddress(bt, qt common.Address) ([]*types.Trade, error)
	FindAndModify(h common.Hash, t *types.Trade) (*types.Trade, error)
	GetByUserAddress(a common.Address) ([]*types.Trade, error)
	GetLatestTrade(bt, qt common.Address) (*types.Trade, error)
	UpdateTradeStatus(h common.Hash, status string) error
	UpdateTradeStatuses(status string, hashes ...common.Hash) ([]*types.Trade, error)
	UpdateTradeStatusesByOrderHashes(status string, hashes ...common.Hash) ([]*types.Trade, error)
	Drop()
}

type TokenDao interface {
	Create(token *types.Token) error
	GetAll() ([]types.Token, error)
	GetByID(id bson.ObjectId) (*types.Token, error)
	GetByAddress(addr common.Address) (*types.Token, error)
	GetQuoteTokens() ([]types.Token, error)
	GetBaseTokens() ([]types.Token, error)
	UpdateFiatPriceBySymbol(symbol string, price float64) error
	Drop() error
}

type FiatPriceDao interface {
	GetLatestQuotes() (map[string]float64, error)
	GetCoinMarketChart(id string, vsCurrency string, days string) (*types.CoinsIDMarketChart, error)
	GetCoinMarketChartRange(id string, vsCurrency string, from int64, to int64) (*types.CoinsIDMarketChart, error)
	Get24hChart(symbol, fiatCurrency string) ([]*types.FiatPriceItem, error)
	Create(items ...*types.FiatPriceItem) error
	FindAndModify(symbol, fiatCurrency, timestamp string, i *types.FiatPriceItem) (*types.FiatPriceItem, error)
	Upsert(symbol, fiatCurrency, timestamp string, i *types.FiatPriceItem) error
}

type NotificationDao interface {
	Create(notifications ...*types.Notification) ([]*types.Notification, error)
	GetAll() ([]types.Notification, error)
	GetByUserAddress(addr common.Address, limit int, offset int) ([]*types.Notification, error)
	GetByID(id bson.ObjectId) (*types.Notification, error)
	FindAndModify(id bson.ObjectId, n *types.Notification) (*types.Notification, error)
	Update(n *types.Notification) error
	Upsert(id bson.ObjectId, n *types.Notification) error
	Delete(notifications ...*types.Notification) error
	DeleteByIds(ids ...bson.ObjectId) error
	Aggregate(q []bson.M) ([]*types.Notification, error)
	Drop()
}

type Engine interface {
	HandleOrders(msg *rabbitmq.Message) error
	// RecoverOrders(matches types.Matches) error
	// CancelOrder(order *types.Order) (*types.EngineResponse, error)
	// DeleteOrder(o *types.Order) error
	Provider() EthereumProvider
}

type WalletService interface {
	CreateAdminWallet(a common.Address) (*types.Wallet, error)
	GetDefaultAdminWallet() (*types.Wallet, error)
	GetOperatorWallets() ([]*types.Wallet, error)
	GetOperatorAddresses() ([]common.Address, error)
	GetAll() ([]types.Wallet, error)
	GetByAddress(addr common.Address) (*types.Wallet, error)
}

type OHLCVService interface {
	Unsubscribe(c *ws.Client)
	UnsubscribeChannel(c *ws.Client, p *types.SubscriptionPayload)
	Subscribe(c *ws.Client, p *types.SubscriptionPayload)
	GetOHLCV(p []types.PairAddresses, duration int64, unit string, timeInterval ...int64) ([]*types.Tick, error)
}

type EthereumService interface {
	WaitMined(hash common.Hash) (*eth.Receipt, error)
	GetPendingNonceAt(a common.Address) (uint64, error)
	GetBalanceAt(a common.Address) (*big.Int, error)
}

type OrderService interface {
	GetOrderCountByUserAddress(addr common.Address) (int, error)
	GetByID(id bson.ObjectId) (*types.Order, error)
	GetByHash(h common.Hash) (*types.Order, error)
	GetByHashes(hashes []common.Hash) ([]*types.Order, error)
	// GetTokenByAddress(a common.Address) (*types.Token, error)
	GetByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error)
	GetCurrentByUserAddress(a common.Address, limit ...int) ([]*types.Order, error)
	GetHistoryByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error)
	NewOrder(o *types.Order) error
	NewStopOrder(so *types.StopOrder) error
	CancelOrder(oc *types.OrderCancel) error
	CancelAllOrder(a common.Address) error
	CancelStopOrder(oc *types.OrderCancel) error
	HandleEngineResponse(res *types.EngineResponse) error
	GetTriggeredStopOrders(baseToken, quoteToken common.Address, lastPrice *big.Int) ([]*types.StopOrder, error)
	UpdateStopOrder(h common.Hash, so *types.StopOrder) error
}

type OrderBookService interface {
	GetOrderBook(bt, qt common.Address) (*types.OrderBook, error)
	GetRawOrderBook(bt, qt common.Address) (*types.RawOrderBook, error)
	SubscribeOrderBook(c *ws.Client, bt, qt common.Address)
	UnsubscribeOrderBook(c *ws.Client)
	UnsubscribeOrderBookChannel(c *ws.Client, bt, qt common.Address)
	SubscribeRawOrderBook(c *ws.Client, bt, qt common.Address)
	UnsubscribeRawOrderBook(c *ws.Client)
	UnsubscribeRawOrderBookChannel(c *ws.Client, bt, qt common.Address)
}

type PairService interface {
	Create(pair *types.Pair) error
	CreatePairs(token common.Address) ([]*types.Pair, error)
	GetByID(id bson.ObjectId) (*types.Pair, error)
	GetByTokenAddress(bt, qt common.Address) (*types.Pair, error)
	GetTokenPairData(bt, qt common.Address) ([]*types.Tick, error)
	GetAllTokenPairData() ([]*types.PairData, error)
	GetAll() ([]types.Pair, error)
	GetListedPairs() ([]types.Pair, error)
	GetUnlistedPairs() ([]types.Pair, error)
}

type TokenService interface {
	Create(token *types.Token) error
	GetByID(id bson.ObjectId) (*types.Token, error)
	GetByAddress(a common.Address) (*types.Token, error)
	GetAll() ([]types.Token, error)
	GetQuoteTokens() ([]types.Token, error)
	GetBaseTokens() ([]types.Token, error)
}

type TradeService interface {
	GetByPairName(p string) ([]*types.Trade, error)
	GetAllTradesByPairAddress(bt, qt common.Address) ([]*types.Trade, error)
	GetSortedTrades(bt, qt common.Address, from, to int64, n int) ([]*types.Trade, error)
	GetSortedTradesByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Trade, error)
	GetByUserAddress(a common.Address) ([]*types.Trade, error)
	GetByHash(h common.Hash) (*types.Trade, error)
	GetByOrderHashes(h []common.Hash) ([]*types.Trade, error)
	GetByMakerOrderHash(h common.Hash) ([]*types.Trade, error)
	GetByTakerOrderHash(h common.Hash) ([]*types.Trade, error)
	Subscribe(c *ws.Client, bt, qt common.Address)
	UnsubscribeChannel(c *ws.Client, bt, qt common.Address)
	Unsubscribe(c *ws.Client)
}

type PriceBoardService interface {
	Subscribe(c *ws.Client, bt, qt common.Address)
	UnsubscribeChannel(c *ws.Client, bt, qt common.Address)
	Unsubscribe(c *ws.Client)
}

type MarketsService interface {
	Subscribe(c *ws.Client)
	UnsubscribeChannel(c *ws.Client)
	Unsubscribe(c *ws.Client)
}

type FiatPriceService interface {
	InitFiatPrice()
	UpdateFiatPrice()
	SyncFiatPrice() error
	GetFiatPriceChart() (map[string][]*types.FiatPriceItem, error)
}

type NotificationService interface {
	Create(n *types.Notification) ([]*types.Notification, error)
	GetAll() ([]types.Notification, error)
	GetByUserAddress(a common.Address, limit int, offset int) ([]*types.Notification, error)
	GetByID(id bson.ObjectId) (*types.Notification, error)
	Update(n *types.Notification) (*types.Notification, error)
}

type TxService interface {
	GetTxCallOptions() *bind.CallOpts
	GetTxSendOptions() (*bind.TransactOpts, error)
	GetTxDefaultSendOptions() (*bind.TransactOpts, error)
	SetTxSender(w *types.Wallet)
	GetCustomTxSendOptions(w *types.Wallet) *bind.TransactOpts
}

type AccountService interface {
	GetAll() ([]types.Account, error)
	Create(account *types.Account) error
	GetByID(id bson.ObjectId) (*types.Account, error)
	GetByAddress(a common.Address) (*types.Account, error)
	FindOrCreate(a common.Address) (*types.Account, error)
	GetTokenBalance(owner common.Address, token common.Address) (*types.TokenBalance, error)
	GetTokenBalances(owner common.Address) (map[common.Address]*types.TokenBalance, error)
	Transfer(token common.Address, fromAddress common.Address, toAddress common.Address, amount *big.Int) error
	GetFavoriteTokens(account common.Address) (map[common.Address]bool, error)
	AddFavoriteToken(account, token common.Address) error
	DeleteFavoriteToken(account, token common.Address) error
}

type DepositService interface {
	SignerPublicKey() common.Address
	GenerateAddress(chain types.Chain) (common.Address, uint64, error)
	GetSchemaVersion() uint64
	RecoveryTransaction(chain types.Chain, address common.Address) error

	// one for wallet, one for relayer
	GetAssociationByChainAddress(chain types.Chain, userAddress common.Address) (*types.AddressAssociationRecord, error)
	GetAssociationByChainAssociatedAddress(chain types.Chain, associatedAddress common.Address) (*types.AddressAssociationRecord, error)

	SaveAssociationByChainAddress(chain types.Chain, address common.Address, index uint64, associatedAddress common.Address, pairAddreses *types.PairAddresses) error
	SaveAssociationStatusByChainAddress(addressAssociation *types.AddressAssociationRecord, status string) error
	SaveDepositTransaction(chain types.Chain, sourceAccount common.Address, txEnvelope string) error
	// SetDelegate to endpoint
	MinimumValueWei() *big.Int
	MinimumValueSat() int64

	SetDelegate(handler SwapEngineHandler)

	// Queue implementation
	QueueAdd(queueTx *types.DepositTransaction) error
	QueuePool() (<-chan *types.DepositTransaction, error)

	// help creating token
	EthereumClient() EthereumClient
}

type SwapEngineHandler interface {
	OnNewEthereumTransaction(transaction swapEthereum.Transaction) error
	OnNewBitcoinTransaction(transaction swapBitcoin.Transaction) error
	OnSubmitTransaction(chain types.Chain, destination string, transaction *types.AssociationTransaction) error
	OnTomochainAccountCreated(chain types.Chain, destination string)
	OnExchanged(chain types.Chain, destination string)
	OnExchangedTimelocked(chain types.Chain, destination string, transaction *types.AssociationTransaction)

	LoadAccountHandler(chain types.Chain, publicKey string) (*types.AddressAssociation, error)
}

type ValidatorService interface {
	ValidateBalance(o *types.Order) error
	ValidateAvailableBalance(o *types.Order) error
}

type EthereumConfig interface {
	GetURL() string
	ExchangeAddress() common.Address
}

type EthereumClient interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*eth.Receipt, error)
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error)
	SendTransaction(ctx context.Context, tx *eth.Transaction) error
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	BalanceAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (*big.Int, error)
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]eth.Log, error)
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- eth.Log) (ethereum.Subscription, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}

type EthereumProvider interface {
	WaitMined(h common.Hash) (*eth.Receipt, error)
	GetBalanceAt(a common.Address) (*big.Int, error)
	GetPendingNonceAt(a common.Address) (uint64, error)
	BalanceOf(owner common.Address, token common.Address) (*big.Int, error)
	Decimals(token common.Address) (uint8, error)
	Symbol(token common.Address) (string, error)
}
