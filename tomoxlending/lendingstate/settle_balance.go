package lendingstate

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
	InQuantity  *big.Int
	InTotal     *big.Int
	OutToken    common.Address
	OutQuantity *big.Int
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

func GetSettleBalance(quoteInterest *big.Int, takerSide string, takerFeeRate *big.Int, baseToken, quoteToken common.Address, makerInterest *big.Int, makerFeeRate *big.Int, baseTokenDecimal *big.Int, quoteTokenDecimal *big.Int, quantityToTrade *big.Int) (*SettleBalance, error) {
	log.Debug("GetSettleBalance", "takerSide", takerSide, "takerFeeRate", takerFeeRate, "baseToken", baseToken, "quoteToken", quoteToken, "makerInterest", makerInterest, "makerFeeRate", makerFeeRate, "baseTokenDecimal", baseTokenDecimal, "quantityToTrade", quantityToTrade)
	var result *SettleBalance
	//result = map[common.Address]map[string]interface{}{}
	if takerSide == Bid {
		// maker InQuantity quoteTokenQuantity=(quantityToTrade*maker.Interest/baseTokenDecimal)
		quoteTokenQuantity := new(big.Int).Mul(quantityToTrade, makerInterest)
		quoteTokenQuantity = quoteTokenQuantity.Div(quoteTokenQuantity, baseTokenDecimal)
		// Fee
		// charge on the token he/she has before the trade, in this case: quoteToken
		// charge on the token he/she has before the trade, in this case: baseToken
		// takerFee = quoteTokenQuantity*takerFeeRate/baseFee=(quantityToTrade*maker.Interest/baseTokenDecimal) * makerFeeRate/baseFee
		takerFee := new(big.Int).Mul(quoteTokenQuantity, takerFeeRate)
		takerFee = new(big.Int).Div(takerFee, common.TomoXBaseFee)
		// charge on the token he/she has before the trade, in this case: baseToken
		makerFee := new(big.Int).Mul(quoteTokenQuantity, makerFeeRate)
		makerFee = new(big.Int).Div(makerFee, common.TomoXBaseFee)
		if quoteTokenQuantity.Cmp(makerFee) <= 0 {
			log.Debug("quantity trade too small", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee)
			return result, ErrQuantityTradeTooSmall
		}
		if baseToken.String() != common.TomoNativeAddress && quoteInterest != nil && quoteInterest.Cmp(common.Big0) > 0 {
			exMakerReceivedFee := new(big.Int).Mul(makerFee, quoteInterest)
			exMakerReceivedFee = exMakerReceivedFee.Div(exMakerReceivedFee, quoteTokenDecimal)
			log.Debug("exMakerReceivedFee", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "quoteInterest", quoteInterest)
			if exMakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("makerFee too small", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "quoteInterest", quoteInterest)
				return result, ErrQuantityTradeTooSmall
			}
			exTakerReceivedFee := new(big.Int).Mul(takerFee, quoteInterest)
			exTakerReceivedFee = exTakerReceivedFee.Div(exTakerReceivedFee, quoteTokenDecimal)
			log.Debug("exTakerReceivedFee", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "quoteInterest", quoteInterest)
			if exTakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("takerFee too small", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "quoteInterest", quoteInterest)
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
		//takerOutTotal= quoteTokenQuantity + takerFee =  quantityToTrade*maker.Interest/baseTokenDecimal + quantityToTrade*maker.Interest/baseTokenDecimal * takerFeeRate/baseFee
		// = quantityToTrade *  maker.Interest/baseTokenDecimal ( 1 +  takerFeeRate/baseFee)
		// = quantityToTrade * maker.Interest * (baseFee + takerFeeRate ) / ( baseTokenDecimal * baseFee)
		takerOutTotal := new(big.Int).Add(quoteTokenQuantity, takerFee)

		result = &SettleBalance{
			Taker: TradeResult{
				Fee:         takerFee,
				InToken:     baseToken,
				InQuantity:  quantityToTrade,
				InTotal:     quantityToTrade,
				OutToken:    quoteToken,
				OutQuantity: quoteTokenQuantity,
				OutTotal:    takerOutTotal,
			},
			Maker: TradeResult{
				Fee:         makerFee,
				InToken:     quoteToken,
				InQuantity:  quoteTokenQuantity,
				InTotal:     inTotal,
				OutToken:    baseToken,
				OutQuantity: quantityToTrade,
				OutTotal:    quantityToTrade,
			},
		}
	} else {
		// Taker InQuantity
		// quoteTokenQuantity = quantityToTrade * makerInterest / baseTokenDecimal
		quoteTokenQuantity := new(big.Int).Mul(quantityToTrade, makerInterest)
		quoteTokenQuantity = quoteTokenQuantity.Div(quoteTokenQuantity, baseTokenDecimal)
		// maker InQuantity

		// Fee
		// charge on the token he/she has before the trade, in this case: baseToken
		// makerFee = quoteTokenQuantity * makerFeeRate / baseFee = quantityToTrade * makerInterest / baseTokenDecimal * makerFeeRate / baseFee
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
		if baseToken.String() != common.TomoNativeAddress && quoteInterest != nil && quoteInterest.Cmp(common.Big0) > 0 {
			exMakerReceivedFee := new(big.Int).Mul(makerFee, quoteInterest)
			exMakerReceivedFee = exMakerReceivedFee.Div(exMakerReceivedFee, quoteTokenDecimal)
			log.Debug("exMakerReceivedFee", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "quoteInterest", quoteInterest)
			if exMakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("makerFee too small", "quoteTokenQuantity", quoteTokenQuantity, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "quoteInterest", quoteInterest)
				return result, ErrQuantityTradeTooSmall
			}
			exTakerReceivedFee := new(big.Int).Mul(takerFee, quoteInterest)
			exTakerReceivedFee = exTakerReceivedFee.Div(exTakerReceivedFee, quoteTokenDecimal)
			log.Debug("exTakerReceivedFee", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "quoteInterest", quoteInterest)
			if exTakerReceivedFee.Cmp(common.RelayerFee) <= 0 {
				log.Debug("takerFee too small", "quoteTokenQuantity", quoteTokenQuantity, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "quoteInterest", quoteInterest)
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
		// makerOutTotal = quoteTokenQuantity + makerFee  = quantityToTrade * makerInterest / baseTokenDecimal + quantityToTrade * makerInterest / baseTokenDecimal * makerFeeRate / baseFee
		// =  quantityToTrade * makerInterest / baseTokenDecimal * (1+makerFeeRate / baseFee)
		// = quantityToTrade  * makerInterest * (baseFee + makerFeeRate) / ( baseTokenDecimal * baseFee )
		makerOutTotal := new(big.Int).Add(quoteTokenQuantity, makerFee)
		// Fee
		result = &SettleBalance{
			Taker: TradeResult{
				Fee:         takerFee,
				InToken:     quoteToken,
				InQuantity:  quoteTokenQuantity,
				InTotal:     inTotal,
				OutToken:    baseToken,
				OutQuantity: quantityToTrade,
				OutTotal:    quantityToTrade,
			},
			Maker: TradeResult{
				Fee:         makerFee,
				InToken:     baseToken,
				InQuantity:  quantityToTrade,
				InTotal:     quantityToTrade,
				OutToken:    quoteToken,
				OutQuantity: quoteTokenQuantity,
				OutTotal:    makerOutTotal,
			},
		}
	}
	return result, nil
}
