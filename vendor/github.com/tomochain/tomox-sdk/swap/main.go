package swap

import (
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/gommon/log"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/swap/bitcoin"
	"github.com/tomochain/tomox-sdk/swap/config"
	"github.com/tomochain/tomox-sdk/swap/ethereum"
	"github.com/tomochain/tomox-sdk/swap/queue"
	"github.com/tomochain/tomox-sdk/swap/storage"
	"github.com/tomochain/tomox-sdk/swap/tomochain"
	"github.com/tomochain/tomox-sdk/utils"
)

// swap is engine
var logger = utils.Logger

// JS SDK use to communicate.
const ProtocolVersion int = 2

type Engine struct {
	Config *config.Config `inject:""`

	bitcoinListener         *bitcoin.Listener         `inject:""`
	bitcoinAddressGenerator *bitcoin.AddressGenerator `inject:""`

	ethereumListener         *ethereum.Listener         `inject:""`
	ethereumAddressGenerator *ethereum.AddressGenerator `inject:""`

	tomochainAccountConfigurator *tomochain.AccountConfigurator `inject:""`
	transactionsQueue            queue.Queue                    `inject:""`

	minimumValueEth string
	minimumValueBtc string

	// decimals in bitcoin is small
	minimumValueSat int64
	minimumValueWei *big.Int

	signerPublicKey common.Address
}

func NewEngine(cfg *config.Config) *Engine {
	engine := &Engine{
		Config: cfg,
	}

	// config blockchains
	engine.configEthereum()
	engine.configBitcoin()

	engine.configTomochain()

	return engine
}

func (engine *Engine) configEthereum() {
	if engine.Config.Ethereum != nil {
		if engine.Config.Ethereum.MasterPublicKey == "" {
			logger.Error("Error: Ethereum master public key is not set")
			return
		}

		ethereumListener := &ethereum.Listener{}
		ethereumClient, err := ethclient.Dial(fmt.Sprintf("%s", engine.Config.Ethereum.RpcServer))
		if err != nil {
			logger.Error("Error connecting to geth")
			os.Exit(-1)
		}

		// config ethereum listener
		ethereumListener.Enabled = true
		ethereumListener.NetworkID = engine.Config.Ethereum.NetworkID
		ethereumListener.ConfirmedBlockNumber = engine.Config.Ethereum.ConfirmedBlockNumber
		ethereumListener.Client = ethereumClient

		engine.minimumValueEth = engine.Config.Ethereum.MinimumValueEth

		ethereumAddressGenerator, err := ethereum.NewAddressGenerator(engine.Config.Ethereum.MasterPublicKey)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}

		engine.ethereumAddressGenerator = ethereumAddressGenerator
		engine.ethereumListener = ethereumListener
	}

}

func (engine *Engine) configBitcoin() {
	if engine.Config.Bitcoin != nil {
		if engine.Config.Bitcoin.MasterPublicKey == "" {
			logger.Error("Error: Bitcoin master public key is not set")
			return
		}

		bitcoinListener := &bitcoin.Listener{}
		connConfig := &rpcclient.ConnConfig{
			Host:         engine.Config.Bitcoin.RpcServer,
			User:         engine.Config.Bitcoin.RpcUser,
			Pass:         engine.Config.Bitcoin.RpcPass,
			HTTPPostMode: true,
			DisableTLS:   true,
		}
		// do not receive notifications
		bitcoinClient, err := rpcclient.New(connConfig, nil)
		if err != nil {
			logger.Error("Error connecting to bitcoin-core")
			os.Exit(-1)
		}

		// config bitcoin listener
		bitcoinListener.Enabled = true
		bitcoinListener.Testnet = engine.Config.Bitcoin.Testnet
		bitcoinListener.ConfirmedBlockNumber = engine.Config.Bitcoin.ConfirmedBlockNumber
		bitcoinListener.Client = bitcoinClient

		engine.minimumValueBtc = engine.Config.Bitcoin.MinimumValueBtc

		var chainParams *chaincfg.Params

		if bitcoinListener.Testnet {
			chainParams = &chaincfg.TestNet3Params
		} else {
			chainParams = &chaincfg.MainNetParams
		}

		bitcoinAddressGenerator, err := bitcoin.NewAddressGenerator(
			engine.Config.Bitcoin.MasterPublicKey, chainParams)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}

		engine.bitcoinAddressGenerator = bitcoinAddressGenerator
		engine.bitcoinListener = bitcoinListener

	}
}

func (engine *Engine) configTomochain() {
	if engine.Config.Tomochain != nil {
		// get signer public key
		engine.signerPublicKey = engine.Config.Tomochain.GetPublicKey()
		// tomochain account configurator
		tomochainAccountConfigurator := tomochain.NewAccountConfigurator(engine.Config.Tomochain)
		tomochainAccountConfigurator.Enabled = true

		if engine.Config.Tomochain.StartingBalance == "" {
			tomochainAccountConfigurator.StartingBalance = "100.00"
		}

		if engine.Config.Ethereum != nil {
			tomochainAccountConfigurator.TokenPriceETH = engine.Config.Ethereum.TokenPrice
		}

		if engine.Config.Bitcoin != nil {
			tomochainAccountConfigurator.TokenPriceBTC = engine.Config.Bitcoin.TokenPrice
		}

		engine.tomochainAccountConfigurator = tomochainAccountConfigurator
	}
}

// SetStorage : update storage mechanism
func (engine *Engine) SetStorage(storage storage.Storage) {

	// set storage for both bitcoin and ethereum
	if engine.ethereumListener != nil {
		engine.ethereumListener.Storage = storage
	}

	if engine.bitcoinListener != nil {
		engine.bitcoinListener.Storage = storage
	}
}

// SetQueue : update queue mechanism, may be rabbitmq implementation
func (engine *Engine) SetQueue(queue queue.Queue) {
	engine.transactionsQueue = queue
}

func (engine *Engine) SetDelegate(handler interfaces.SwapEngineHandler) {
	// delegate some handlers
	if engine.ethereumListener != nil && handler.OnNewEthereumTransaction != nil {
		engine.ethereumListener.TransactionHandler = handler.OnNewEthereumTransaction
	}

	if engine.bitcoinListener != nil && handler.OnNewBitcoinTransaction != nil {
		engine.bitcoinListener.TransactionHandler = handler.OnNewBitcoinTransaction
	}

	engine.tomochainAccountConfigurator.OnSubmitTransaction = handler.OnSubmitTransaction
	engine.tomochainAccountConfigurator.OnAccountCreated = handler.OnTomochainAccountCreated
	engine.tomochainAccountConfigurator.OnExchanged = handler.OnExchanged
	engine.tomochainAccountConfigurator.OnExchangedTimelocked = handler.OnExchangedTimelocked
	engine.tomochainAccountConfigurator.LoadAccountHandler = handler.LoadAccountHandler
}

func (engine *Engine) Start() error {

	ethereumEnabled := engine.ethereumListener != nil && engine.ethereumListener.Enabled
	bitcoinEnabled := engine.bitcoinListener != nil && engine.bitcoinListener.Enabled

	if !ethereumEnabled && !bitcoinEnabled {
		return errors.New("At least one listener (BitcoinListener or EthereumListener) must be enabled")
	}

	var err error

	if ethereumEnabled {
		engine.minimumValueWei, err = ethereum.EthToWei(engine.minimumValueEth)
		if err != nil {
			return errors.Wrapf(err, "Invalid minimum accepted Ethereum transaction value: %s", engine.minimumValueEth)
		}

		if engine.minimumValueWei.Cmp(new(big.Int)) == 0 {
			return errors.New("Minimum accepted Ethereum transaction value must be larger than 0")
		}

		err = engine.ethereumListener.Start()
		if err != nil {
			return errors.Wrap(err, "Error starting EthereumListener")
		}
	}

	if bitcoinEnabled {
		engine.minimumValueSat, err = bitcoin.BtcToSat(engine.minimumValueBtc)

		if err != nil {
			return errors.Wrapf(err, "Invalid minimum accepted Bitcoin transaction value: %s"+engine.minimumValueBtc)
		}

		if engine.minimumValueSat == 0 {
			return errors.New("Minimum accepted Bitcoin transaction value must be larger than 0")
		}

		err = engine.bitcoinListener.Start()
		if err != nil {
			return errors.Wrap(err, "Error starting BitcoinListener")
		}
	}

	err = engine.tomochainAccountConfigurator.Start()
	if err != nil {
		return errors.Wrap(err, "Error starting TomochainAccountConfigurator")
	}

	go engine.poolTransactionsQueue()

	return nil
}

func (engine *Engine) TransactionsQueue() queue.Queue {
	return engine.transactionsQueue
}

// public method to access private properties, this avoids setting props directly cause mistmatch from config
func (engine *Engine) EthereumAddressGenerator() *ethereum.AddressGenerator {
	return engine.ethereumAddressGenerator
}

func (engine *Engine) BitcoinAddressGenerator() *bitcoin.AddressGenerator {
	return engine.bitcoinAddressGenerator
}

func (engine *Engine) TomochainAccountConfigurator() *tomochain.AccountConfigurator {
	return engine.tomochainAccountConfigurator
}

func (engine *Engine) SignerPublicKey() common.Address {
	return engine.signerPublicKey
}

func (engine *Engine) MinimumValueEth() string {
	return engine.minimumValueEth
}

func (engine *Engine) MinimumValueBtc() string {
	return engine.minimumValueBtc
}

func (engine *Engine) MinimumValueSat() int64 {
	return engine.minimumValueSat
}

func (engine *Engine) MinimumValueWei() *big.Int {
	return engine.minimumValueWei
}

// poolTransactionsQueue pools transactions queue which contains only processed and
// validated transactions and sends it to TomochainAccountConfigurator for account configuration.
func (engine *Engine) poolTransactionsQueue() {
	logger.Infof("Started pooling transactions queue")
	msgs, err := engine.transactionsQueue.QueuePool()

	if err != nil {
		logger.Infof("Error pooling transactions queue")
		time.Sleep(5 * time.Second)
		engine.shutdown()
		return
	}

	signalInterrupt := make(chan os.Signal, 1)
	signal.Notify(signalInterrupt, os.Interrupt)

	var endWaiter sync.WaitGroup
	endWaiter.Add(1)

	// eating messages from the read-only channel
	go func() {
		for {
			select {
			case transaction := <-msgs:
				if transaction == nil {
					time.Sleep(time.Second)
					continue
				}

				logger.Infof("Received transaction from transactions queue: %v", transaction)
				go engine.tomochainAccountConfigurator.ConfigureAccount(transaction)
			case <-signalInterrupt:
				// wait for interrupt
				endWaiter.Done()
			default:
				time.Sleep(time.Second)
			}
		}
	}()

	endWaiter.Wait()

	logger.Infof("Ending transaction queue")
	engine.shutdown()

	os.Exit(0)
}

func (engine *Engine) shutdown() {
	// do something
	engine.ethereumListener.Stop()
	engine.tomochainAccountConfigurator.Stop()
}
