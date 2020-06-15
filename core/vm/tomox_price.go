package vm

import (
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/log"
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
	if t.tomoxState != nil && len(input) == 64 {
		base := common.BytesToAddress(input[12:32]) // 20 bytes from 13-32
		quote := common.BytesToAddress(input[44:]) // 20 bytes from 45-64
		price := t.tomoxState.GetLastPrice(tradingstate.GetTradingOrderBookHash(base, quote))
		if price != nil {
			log.Debug("Run GetLastPrice", "base", base.Hex(), "quote", quote.Hex(), "price", price)
			return price.Bytes(), nil
		}
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
	if t.tomoxState != nil && len(input) == 64 {
		base := common.BytesToAddress(input[12:32]) // 20 bytes from 13-32
		quote := common.BytesToAddress(input[44:]) // 20 bytes from 45-64
		price := t.tomoxState.GetMediumPriceBeforeEpoch(tradingstate.GetTradingOrderBookHash(base, quote))
		if price != nil {
			log.Debug("Run GetEpochPrice", "base", base.Hex(), "quote", quote.Hex(), "price", price)
			return price.Bytes(), nil
		}
	}
	return []byte{}, nil
}

func (t *tomoxEpochPrice) SetTradingState(tomoxState *tradingstate.TradingStateDB) {
	if tomoxState != nil {
		t.tomoxState = tomoxState.Copy()
	}
}


