package bitcoin

import (
	"math/big"

	"github.com/tomochain/tomox-sdk/swap/config"
	"github.com/tomochain/tomox-sdk/utils/binance"
)

// ValueToTomo need to convert BTC to ETH, because Tomo using ETH as unit
func (t *Transaction) ValueToTomo() string {
	valueSat := new(big.Int).SetInt64(t.ValueSat)
	valueBtc := new(big.Rat).Quo(new(big.Rat).SetInt(valueSat), satInBtc)
	return valueBtc.FloatString(config.TomoAmountPrecision)
}

func (t *Transaction) ValueToWei() string {

	lastPrice, err := binance.GetLastPrice("ETH", "BTC")
	if err != nil {
		logger.Error(err)
		return ""
	}
	multiplier := new(big.Rat).SetInt64(1e10) // decimals of eth is 10^18, btc is 10^8
	valueSat := new(big.Rat).SetInt64(t.ValueSat)

	exchangeRate, ok := new(big.Rat).SetString(lastPrice)
	if !ok {
		return ""
	}
	logger.Infof("Last price: %s", lastPrice)
	valueWei := new(big.Rat).Quo(valueSat, exchangeRate)
	valueWei = valueWei.Mul(valueWei, multiplier)
	return valueWei.FloatString(0)
}
