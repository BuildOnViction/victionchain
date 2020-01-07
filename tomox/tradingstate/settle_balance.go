package tradingstate

import (
	"encoding/json"
	"errors"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/log"
	"math/big"
)

var ErrQuantityTradeTooSmall = errors.New("quantity trade too small")

type TradeResult struct {
	Fee         *big.Int
	InToken     common.Address
	InTotal     *big.Int
	OutToken    common.Address
	OutTotal    *big.Int
}
type SettleBalance struct {
	Taker TradeResult
	Maker TradeResult
}

func (settleBalance *SettleBalance) String() string {
	jsonData, _ := json.Marshal(settleBalance)
	return string(jsonData)
}

func GetSettleBalance(quotePrice *big.Int, takerSide string, takerFeeRate *big.Int, baseToken, quoteToken common.Address, makerPrice *big.Int, makerFeeRate *big.Int, baseTokenDecimal *big.Int, quoteTokenDecimal *big.Int, quantityToTrade *big.Int) (*SettleBalance, error) {
	log.Debug("GetSettleBalance", "takerSide", takerSide, "takerFeeRate", takerFeeRate, "baseToken", baseToken, "quoteToken", quoteToken, "makerPrice", makerPrice, "makerFeeRate", makerFeeRate, "baseTokenDecimal", baseTokenDecimal, "quantityToTrade", quantityToTrade)
	var result *SettleBalance
	//result = map[common.Address]map[string]interface{}{}
	if takerSide == Bid {
		// maker InQuantity quoteTokenQuantity=(quantityToTrade*maker.Price/baseTokenDecimal)
		quoteTokenQuantity := new(big.Int).Mul(quantityToTrade, makerPrice)
		quoteTokenQuantity = quoteTokenQuantity.Div(quoteTokenQuantity, baseTokenDecimal)
		// Fee
		// charge on the token he/she has before the trade, in this case: quoteToken
		// charge on the token he/she has before the trade, in this case: baseToken
		// takerFee = quoteTokenQuantity*takerFeeRate/baseFee=(quantityToTrade*maker.Price/baseTokenDecimal) * makerFeeRate/baseFee
		takerFee := new(big.Int).Mul(quoteTokenQuantity, takerFeeRate)
		takerFee = new(big.Int).Div(takerFee, common.TomoXBaseFee)
		// charge on the token he/she has before the trade, in this case: baseToken
		makerFee := new(big.Int).Mul(quoteTokenQuantity, makerFeeRate)
		makerFee = new(big.Int).Div(makerFee, common.TomoXBaseFee)
		if quoteTokenQuantity.Cmp(makerFee) <= 0 {
			log.Debug("quantity trade too small", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee)
			return result, ErrQuantityTradeTooSmall
		}
		if baseToken.String() != common.TomoNativeAddress && quotePrice != nil && quotePrice.Cmp(common.Big0) > 0 {
			exMakerReceivedFee := new(big.Int).Mul(makerFee, quotePrice)
			exMakerReceivedFee = exMakerReceivedFee.Div(exMakerReceivedFee, quoteTokenDecimal)
			log.Debug("exMakerReceivedFee", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "quotePrice", quotePrice)
			if exMakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("makerFee too small", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "quotePrice", quotePrice)
				return result, ErrQuantityTradeTooSmall
			}
			exTakerReceivedFee := new(big.Int).Mul(takerFee, quotePrice)
			exTakerReceivedFee = exTakerReceivedFee.Div(exTakerReceivedFee, quoteTokenDecimal)
			log.Debug("exTakerReceivedFee", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "quotePrice", quotePrice)
			if exTakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("takerFee too small", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "quotePrice", quotePrice)
				return result, ErrQuantityTradeTooSmall
			}
		} else if baseToken.String() == common.TomoNativeAddress {
			exMakerReceivedFee := new(big.Int).Mul(quantityToTrade, makerFeeRate)
			exMakerReceivedFee = exMakerReceivedFee.Div(exMakerReceivedFee, common.TomoXBaseFee)
			log.Debug("exMakerReceivedFee", "quantityToTrade", quantityToTrade, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "makerFeeRate", makerFeeRate)
			if exMakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("makerFee too small", "quantityToTrade", quantityToTrade, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "makerFeeRate", makerFeeRate)
				return result, ErrQuantityTradeTooSmall
			}
			exTakerReceivedFee := new(big.Int).Mul(quantityToTrade, takerFeeRate)
			exTakerReceivedFee = exTakerReceivedFee.Div(exTakerReceivedFee, common.TomoXBaseFee)
			log.Debug("exTakerReceivedFee", "quantityToTrade", quantityToTrade, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "takerFeeRate", takerFeeRate)
			if exTakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("takerFee too small", "quantityToTrade", quantityToTrade, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "takerFeeRate", takerFeeRate)
				return result, ErrQuantityTradeTooSmall
			}
		}
		inTotal := new(big.Int).Sub(quoteTokenQuantity, makerFee)
		//takerOutTotal= quoteTokenQuantity + takerFee =  quantityToTrade*maker.Price/baseTokenDecimal + quantityToTrade*maker.Price/baseTokenDecimal * takerFeeRate/baseFee
		// = quantityToTrade *  maker.Price/baseTokenDecimal ( 1 +  takerFeeRate/baseFee)
		// = quantityToTrade * maker.Price * (baseFee + takerFeeRate ) / ( baseTokenDecimal * baseFee)
		takerOutTotal := new(big.Int).Add(quoteTokenQuantity, takerFee)

		result = &SettleBalance{
			Taker: TradeResult{
				Fee:         takerFee,
				InToken:     baseToken,
				InTotal:     quantityToTrade,
				OutToken:    quoteToken,
				OutTotal:    takerOutTotal,
			},
			Maker: TradeResult{
				Fee:         makerFee,
				InToken:     quoteToken,
				InTotal:     inTotal,
				OutToken:    baseToken,
				OutTotal:    quantityToTrade,
			},
		}
	} else {
		// Taker InQuantity
		// quoteTokenQuantity = quantityToTrade * makerPrice / baseTokenDecimal
		quoteTokenQuantity := new(big.Int).Mul(quantityToTrade, makerPrice)
		quoteTokenQuantity = quoteTokenQuantity.Div(quoteTokenQuantity, baseTokenDecimal)
		// maker InQuantity

		// Fee
		// charge on the token he/she has before the trade, in this case: baseToken
		// makerFee = quoteTokenQuantity * makerFeeRate / baseFee = quantityToTrade * makerPrice / baseTokenDecimal * makerFeeRate / baseFee
		// charge on the token he/she has before the trade, in this case: quoteToken
		makerFee := new(big.Int).Mul(quoteTokenQuantity, makerFeeRate)
		makerFee = makerFee.Div(makerFee, common.TomoXBaseFee)

		// charge on the token he/she has before the trade, in this case: baseToken
		takerFee := new(big.Int).Mul(quoteTokenQuantity, takerFeeRate)
		takerFee = new(big.Int).Div(takerFee, common.TomoXBaseFee)
		if quoteTokenQuantity.Cmp(takerFee) <= 0 {
			log.Debug("quantity trade too small", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee)
			return result, ErrQuantityTradeTooSmall
		}
		if baseToken.String() != common.TomoNativeAddress && quotePrice != nil && quotePrice.Cmp(common.Big0) > 0 {
			exMakerReceivedFee := new(big.Int).Mul(makerFee, quotePrice)
			exMakerReceivedFee = exMakerReceivedFee.Div(exMakerReceivedFee, quoteTokenDecimal)
			log.Debug("exMakerReceivedFee", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "quotePrice", quotePrice)
			if exMakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("makerFee too small", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "quotePrice", quotePrice)
				return result, ErrQuantityTradeTooSmall
			}
			exTakerReceivedFee := new(big.Int).Mul(takerFee, quotePrice)
			exTakerReceivedFee = exTakerReceivedFee.Div(exTakerReceivedFee, quoteTokenDecimal)
			log.Debug("exTakerReceivedFee", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "quotePrice", quotePrice)
			if exTakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("takerFee too small", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "quotePrice", quotePrice)
				return result, ErrQuantityTradeTooSmall
			}
		} else if baseToken.String() == common.TomoNativeAddress {
			exMakerReceivedFee := new(big.Int).Mul(quantityToTrade, makerFeeRate)
			exMakerReceivedFee = exMakerReceivedFee.Div(exMakerReceivedFee, common.TomoXBaseFee)
			log.Debug("exMakerReceivedFee", "quantityToTrade", quantityToTrade, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "makerFeeRate", makerFeeRate)
			if exMakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("makerFee too small", "quantityToTrade", quantityToTrade, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "makerFeeRate", makerFeeRate)
				return result, ErrQuantityTradeTooSmall
			}
			exTakerReceivedFee := new(big.Int).Mul(quantityToTrade, takerFeeRate)
			exTakerReceivedFee = exTakerReceivedFee.Div(exTakerReceivedFee, common.TomoXBaseFee)
			log.Debug("exTakerReceivedFee", "quantityToTrade", quantityToTrade, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "takerFeeRate", takerFeeRate)
			if exTakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("takerFee too small", "quantityToTrade", quantityToTrade, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "takerFeeRate", takerFeeRate)
				return result, ErrQuantityTradeTooSmall
			}
		}
		inTotal := new(big.Int).Sub(quoteTokenQuantity, takerFee)
		// makerOutTotal = quoteTokenQuantity + makerFee  = quantityToTrade * makerPrice / baseTokenDecimal + quantityToTrade * makerPrice / baseTokenDecimal * makerFeeRate / baseFee
		// =  quantityToTrade * makerPrice / baseTokenDecimal * (1+makerFeeRate / baseFee)
		// = quantityToTrade  * makerPrice * (baseFee + makerFeeRate) / ( baseTokenDecimal * baseFee )
		makerOutTotal := new(big.Int).Add(quoteTokenQuantity, makerFee)
		// Fee
		result = &SettleBalance{
			Taker: TradeResult{
				Fee:         takerFee,
				InToken:     quoteToken,
				InTotal:     inTotal,
				OutToken:    baseToken,
				OutTotal:    quantityToTrade,
			},
			Maker: TradeResult{
				Fee:         makerFee,
				InToken:     baseToken,
				InTotal:     quantityToTrade,
				OutToken:    quoteToken,
				OutTotal:    makerOutTotal,
			},
		}
	}
	return result, nil
}
