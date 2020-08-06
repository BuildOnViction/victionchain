package vm

import (
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/params"
	"github.com/tomochain/tomochain/tomox/tradingstate"
)
const TomoXPriceNumberOfBytesReturn = 32
// tomoxPrice implements a pre-compile contract to get token price in tomox

type tomoxLastPrice struct {
	tradingStateDB *tradingstate.TradingStateDB
}
type tomoxEpochPrice struct {
	tradingStateDB *tradingstate.TradingStateDB
}

func (t *tomoxLastPrice) RequiredGas(input []byte) uint64 {
	return params.TomoXPriceGas
}

func (t *tomoxLastPrice) Run(input []byte) ([]byte, error) {
	// input includes baseTokenAddress, quoteTokenAddress
	if t.tradingStateDB != nil && len(input) == 64 {
		base := common.BytesToAddress(input[12:32]) // 20 bytes from 13-32
		quote := common.BytesToAddress(input[44:]) // 20 bytes from 45-64
		price := t.tradingStateDB.GetLastPrice(tradingstate.GetTradingOrderBookHash(base, quote))
		if price != nil {
			log.Debug("Run GetLastPrice", "base", base.Hex(), "quote", quote.Hex(), "price", price)
			return common.LeftPadBytes(price.Bytes(), TomoXPriceNumberOfBytesReturn), nil
		}
	}
	return common.LeftPadBytes([]byte{}, TomoXPriceNumberOfBytesReturn), nil
}

func (t *tomoxLastPrice) SetTradingState(tradingStateDB *tradingstate.TradingStateDB) {
	if tradingStateDB != nil {
		t.tradingStateDB = tradingStateDB.Copy()
	} else {
		t.tradingStateDB = nil
	}
}

func (t *tomoxEpochPrice) RequiredGas(input []byte) uint64 {
	return params.TomoXPriceGas
}

func (t *tomoxEpochPrice) Run(input []byte) ([]byte, error) {
	// input includes baseTokenAddress, quoteTokenAddress
	if t.tradingStateDB != nil && len(input) == 64 {
		base := common.BytesToAddress(input[12:32]) // 20 bytes from 13-32
		quote := common.BytesToAddress(input[44:]) // 20 bytes from 45-64
		price := t.tradingStateDB.GetMediumPriceBeforeEpoch(tradingstate.GetTradingOrderBookHash(base, quote))
		if price != nil {
			log.Debug("Run GetEpochPrice", "base", base.Hex(), "quote", quote.Hex(), "price", price)
			return common.LeftPadBytes(price.Bytes(), TomoXPriceNumberOfBytesReturn), nil
		}
	}
	return common.LeftPadBytes([]byte{}, TomoXPriceNumberOfBytesReturn), nil
}

func (t *tomoxEpochPrice) SetTradingState(tradingStateDB *tradingstate.TradingStateDB) {
	if tradingStateDB != nil {
		t.tradingStateDB = tradingStateDB.Copy()
	} else {
		t.tradingStateDB = nil
	}
}


