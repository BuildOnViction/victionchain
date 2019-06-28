package ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/contracts/contractsinterfaces"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/utils"
)

type EthereumProvider struct {
	Client interfaces.EthereumClient
	Config interfaces.EthereumConfig
}

func NewEthereumProvider(c interfaces.EthereumClient) *EthereumProvider {
	url := app.Config.Ethereum["http_url"]
	exchange := common.HexToAddress(app.Config.Ethereum["exchange_address"])
	config := NewEthereumConfig(url, exchange)

	return &EthereumProvider{
		Client: c,
		Config: config,
	}
}

func NewDefaultEthereumProvider() *EthereumProvider {
	url := app.Config.Ethereum["http_url"]
	exchange := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	conn, err := rpc.DialHTTP(app.Config.Ethereum["http_url"])
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(conn)
	config := NewEthereumConfig(url, exchange)

	return &EthereumProvider{
		Client: client,
		Config: config,
	}
}

func NewWebsocketProvider() *EthereumProvider {
	url := app.Config.Ethereum["ws_url"]
	exchange := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	conn, err := rpc.DialWebsocket(context.Background(), url, "")
	if err != nil {
		panic(err)
	}

	ethClient := ethclient.NewClient(conn)
	config := NewEthereumConfig(url, exchange)

	return &EthereumProvider{
		Client: ethClient,
		Config: config,
	}
}

func NewSimulatedEthereumProvider(accs []common.Address) *EthereumProvider {
	url := app.Config.Ethereum["http_url"]
	exchange := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	config := NewEthereumConfig(url, exchange)
	client := NewSimulatedClient(accs)

	return &EthereumProvider{
		Client: client,
		Config: config,
	}
}

func (e *EthereumProvider) WaitMined(hash common.Hash) (*eth.Receipt, error) {
	ctx := context.Background()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		receipt, _ := e.Client.TransactionReceipt(ctx, hash)
		if receipt != nil {
			return receipt, nil
		}

		// if err != nil {
		// 	logger.Error(err)
		// 	// return nil, err
		// }

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}

func (e *EthereumProvider) GetBalanceAt(a common.Address) (*big.Int, error) {
	ctx := context.Background()
	nonce, err := e.Client.BalanceAt(ctx, a, nil)
	if err != nil {
		logger.Error(err)
		return big.NewInt(0), err
	}

	return nonce, nil
}

func (e *EthereumProvider) GetPendingNonceAt(a common.Address) (uint64, error) {
	ctx := context.Background()
	nonce, err := e.Client.PendingNonceAt(ctx, a)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return nonce, nil
}

func (e *EthereumProvider) Decimals(token common.Address) (uint8, error) {
	var tokenInterface *contractsinterfaces.ERC20
	var err error

	// retry in case the connection with the ethereum client is asleep
	err = utils.Retry(3, func() error {
		tokenInterface, err = contractsinterfaces.NewERC20(token, e.Client)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return 0, err
	}

	opts := &bind.CallOpts{Pending: true}
	decimals, err := tokenInterface.Decimals(opts)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return decimals, nil
}

func (e *EthereumProvider) Symbol(token common.Address) (string, error) {
	// retry in case the connection with the ethereum client is asleep
	var tokenInterface *contractsinterfaces.ERC20
	var err error

	err = utils.Retry(3, func() error {
		tokenInterface, err = contractsinterfaces.NewERC20(token, e.Client)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return "", err
	}

	opts := &bind.CallOpts{Pending: true}
	symbol, err := tokenInterface.Symbol(opts)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	return symbol, nil
}

func (e *EthereumProvider) BalanceOf(owner common.Address, token common.Address) (*big.Int, error) {
	tokenInterface, err := contractsinterfaces.NewToken(token, e.Client)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	opts := &bind.CallOpts{Pending: true}
	b, err := tokenInterface.BalanceOf(opts, owner)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return b, nil
}
