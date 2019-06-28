package ethereum

import (
	"context"
	"math/big"

	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/utils/math"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
)

type SimulatedClient struct {
	*backends.SimulatedBackend
}

func (b *SimulatedClient) PendingBalanceAt(ctx context.Context, acc common.Address) (*big.Int, error) {
	return nil, errors.New("PendingBalanceAt is not implemented on the simulated backend")
}

func NewSimulatedClientWithGasLimit(accs []common.Address, gasLimit uint64) *SimulatedClient {
	weiBalance := &big.Int{}
	ether := big.NewInt(1e18)
	etherBalance := big.NewInt(1000)
	firstEtherBalance := big.NewInt(1e8)

	alloc := make(core.GenesisAlloc)
	weiBalance.Mul(etherBalance, ether)
	firstWeiBalance := math.Mul(firstEtherBalance, ether)

	for index, a := range accs {
		if index == 0 {
			alloc[a] = core.GenesisAccount{Balance: firstWeiBalance}
		} else {
			alloc[a] = core.GenesisAccount{Balance: weiBalance}
		}
	}

	client := backends.NewSimulatedBackend(alloc, gasLimit)

	return &SimulatedClient{client}
}

func NewSimulatedClient(accs []common.Address) *SimulatedClient {
	return NewSimulatedClientWithGasLimit(accs, 5e6)
}
