package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/crons"
	"github.com/tomochain/tomox-sdk/daos"
	"github.com/tomochain/tomox-sdk/endpoints"
	"github.com/tomochain/tomox-sdk/engine"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/ethereum"
	"github.com/tomochain/tomox-sdk/operator"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/services"
	"github.com/tomochain/tomox-sdk/swap"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/ws"
)

const (
	swaggerUIDir = "/swaggerui/"
)

func Start() {
	env := os.Getenv("GO_ENV")

	if err := app.LoadConfig("./config", env); err != nil {
		panic(err)
	}

	if err := errors.LoadMessages(app.Config.ErrorFile); err != nil {
		panic(err)
	}

	// connect to the database
	_, err := daos.InitSession(nil)
	if err != nil {
		panic(err)
	}

	rabbitConn := rabbitmq.InitConnection(app.Config.RabbitMQURL)

	provider := ethereum.NewWebsocketProvider()

	router := NewRouter(provider, rabbitConn)
	// http.Handle("/", router)
	router.HandleFunc("/socket", ws.ConnectionEndpoint)

	// start the server
	address := fmt.Sprintf(":%v", app.Config.ServerPort)
	log.Printf("server %v is started at %v\n", app.Version, address)

	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Accept", "Authorization", "Access-Control-Allow-Origin"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	panic(http.ListenAndServe(address, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router)))
}

func NewRouter(
	provider *ethereum.EthereumProvider,
	rabbitConn *rabbitmq.Connection,
) *mux.Router {

	r := mux.NewRouter()

	// get daos for dependency injection
	orderDao := daos.NewOrderDao()
	stopOrderDao := daos.NewStopOrderDao()
	tokenDao := daos.NewTokenDao()
	pairDao := daos.NewPairDao()
	tradeDao := daos.NewTradeDao()
	accountDao := daos.NewAccountDao()
	walletDao := daos.NewWalletDao()
	configDao := daos.NewConfigDao()
	associationDao := daos.NewAssociationDao()
	fiatPriceDao := daos.NewFiatPriceDao()
	notificationDao := daos.NewNotificationDao()

	// instantiate engine
	eng := engine.NewEngine(rabbitConn, orderDao, stopOrderDao, tradeDao, pairDao, provider)
	swapEngine := swap.NewEngine(app.Config.Deposit)

	// get services for injection
	accountService := services.NewAccountService(accountDao, tokenDao)
	ohlcvService := services.NewOHLCVService(tradeDao)
	tokenService := services.NewTokenService(tokenDao)
	validatorService := services.NewValidatorService(provider, accountDao, orderDao, pairDao)
	pairService := services.NewPairService(pairDao, tokenDao, tradeDao, orderDao, eng, provider)
	orderService := services.NewOrderService(orderDao, stopOrderDao, pairDao, accountDao, tradeDao, notificationDao, eng, validatorService, rabbitConn)
	orderBookService := services.NewOrderBookService(pairDao, tokenDao, orderDao, eng)
	tradeService := services.NewTradeService(orderDao, tradeDao, rabbitConn)

	walletService := services.NewWalletService(walletDao)

	// txservice for deposit
	// wallet := &types.NewWalletFromPrivateKey(app.Config.Deposit.Tomochain.SignerPrivateKey)
	// we already have them so no need to re-calculate
	wallet := &types.Wallet{
		Address:    app.Config.Deposit.Tomochain.GetPublicKey(),
		PrivateKey: app.Config.Deposit.Tomochain.GetPrivateKey(),
	}
	txService := services.NewTxService(walletDao, wallet)
	depositService := services.NewDepositService(configDao, associationDao, pairDao, orderDao, swapEngine, eng, rabbitConn)
	priceBoardService := services.NewPriceBoardService(tokenDao, tradeDao)
	fiatPriceService := services.NewFiatPriceService(tokenDao, fiatPriceDao)
	marketsService := services.NewMarketsService(pairDao, orderDao, tradeDao, ohlcvService, fiatPriceService)
	notificationService := services.NewNotificationService(notificationDao)

	// start cron service
	cronService := crons.NewCronService(ohlcvService, priceBoardService, pairService, fiatPriceService, eng)

	// deploy operator
	op, err := operator.NewOperator(
		walletService,
		tradeService,
		orderService,
		provider,
		rabbitConn,
		accountService,
		tokenService,
	)

	if err != nil {
		panic(err)
	}

	// deploy http and ws endpoints
	endpoints.ServeInfoResource(r, walletService, tokenService)
	endpoints.ServeAccountResource(r, accountService)
	endpoints.ServeTokenResource(r, tokenService)
	endpoints.ServePairResource(r, pairService)
	endpoints.ServeOrderBookResource(r, orderBookService)
	endpoints.ServeOHLCVResource(r, ohlcvService)
	endpoints.ServeTradeResource(r, tradeService)
	endpoints.ServeOrderResource(r, orderService, accountService)

	endpoints.ServeDepositResource(r, depositService, walletService, txService)
	endpoints.ServePriceBoardResource(r, priceBoardService)
	endpoints.ServeMarketsResource(r, marketsService)
	endpoints.ServeNotificationResource(r, notificationService)

	// Swagger UI
	sh := http.StripPrefix(swaggerUIDir, http.FileServer(http.Dir("."+swaggerUIDir)))
	r.PathPrefix(swaggerUIDir).Handler(sh)

	//initialize rabbitmq subscriptions
	rabbitConn.SubscribeOrders(eng.HandleOrders)
	rabbitConn.SubscribeEngineResponses(orderService.HandleEngineResponse)
	rabbitConn.SubscribeTrades(op.HandleTrades)
	rabbitConn.SubscribeOperator(orderService.HandleOperatorMessages)

	// Initialize fiat price
	fiatPriceService.InitFiatPrice()

	// initialize MongoDB Change Streams
	go orderService.WatchChanges()
	go tradeService.WatchChanges()

	cronService.InitCrons()
	return r
}
