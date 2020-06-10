package vm

import (
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/params"
	"github.com/tomochain/tomochain/tomox/tradingstate"
)

// tomoxPrice implements a pre-compile contract to get token price in tomox

type tomoxLastPrice struct {
	tomoxState *tradingstate.TradingStateDB
}
type tomoxEpochPrice struct {
	tomoxState *tradingstate.TradingStateDB
}

func (t *tomoxLastPrice) RequiredGas(input []byte) uint64 {
	return params.TomoXPriceGas
}

func (t *tomoxLastPrice) Run(input []byte) ([]byte, error) {
	// input includes baseTokenAddress, quoteTokenAddress
	if t.tomoxState != nil && len(input) == 40 {
		base := common.BytesToAddress(input[:20])
		quote := common.BytesToAddress(input[20:])
		return t.tomoxState.GetLastPrice(tradingstate.GetTradingOrderBookHash(base, quote)).Bytes(), nil
	}
	return []byte{}, nil
}

func (t *tomoxLastPrice) SetTradingState(tomoxState *tradingstate.TradingStateDB) {
	if tomoxState != nil {
		t.tomoxState = tomoxState.Copy()
	}
}

func (t *tomoxEpochPrice) RequiredGas(input []byte) uint64 {
	return params.TomoXPriceGas
}

func (t *tomoxEpochPrice) Run(input []byte) ([]byte, error) {
	// input includes baseTokenAddress, quoteTokenAddress
	if t.tomoxState != nil && len(input) == 40 {
		base := common.BytesToAddress(input[:20])
		quote := common.BytesToAddress(input[20:])
		return t.tomoxState.GetMediumPriceBeforeEpoch(tradingstate.GetTradingOrderBookHash(base, quote)).Bytes(), nil
	}
	return []byte{}, nil
}

func (t *tomoxEpochPrice) SetTradingState(tomoxState *tradingstate.TradingStateDB) {
	if tomoxState != nil {
		t.tomoxState = tomoxState.Copy()
	}
}


