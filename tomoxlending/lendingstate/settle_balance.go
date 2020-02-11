package lendingstate

import (
	"encoding/json"
	"errors"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/log"
	"math/big"
)

var (
	ErrQuantityTradeTooSmall = errors.New("quantity trade too small")
	ErrInvalidCollateralPrice = errors.New("unable to retrieve price of this collateral. Please try another collateral")
)

type TradeResult struct {
	Fee      *big.Int
	InToken  common.Address
	InTotal  *big.Int
	OutToken common.Address
	OutTotal *big.Int
}
type LendingSettleBalance struct {
	Taker                  TradeResult
	Maker                  TradeResult
	CollateralLockedAmount *big.Int
}

func (settleBalance *LendingSettleBalance) String() string {
	jsonData, _ := json.Marshal(settleBalance)
	return string(jsonData)
}

func GetSettleBalance(takerSide string,
	lendTokenTOMOPrice,
	collateralPrice,
	depositRate,
	borrowFee *big.Int,
	lendingToken,
	collateralToken common.Address,
	lendTokenDecimal,
	collateralTokenDecimal *big.Int,
	quantityToLend *big.Int) (*LendingSettleBalance, error) {
	log.Debug("GetSettleBalance", "takerSide", takerSide, "borrowFee", borrowFee, "lendingToken", lendingToken, "collateralToken", collateralToken, "quantityToLend", quantityToLend)
	if collateralPrice == nil || collateralPrice.Sign() <= 0 {
		return nil, ErrInvalidCollateralPrice
	}

	var result *LendingSettleBalance
	//result = map[common.Address]map[string]interface{}{}
	if takerSide == Borrowing {
		// taker = Borrower : takerOutTotal = CollateralLockedAmount = quantityToLend * collateral Token Decimal/ CollateralPrice  * deposit rate
		takerOutTotal := new(big.Int).Mul(quantityToLend, collateralTokenDecimal)
		takerOutTotal = new(big.Int).Mul(takerOutTotal, depositRate) // eg: depositRate = 150%
		takerOutTotal = new(big.Int).Div(takerOutTotal, big.NewInt(100))
		takerOutTotal = new(big.Int).Div(takerOutTotal, collateralPrice)
		// Fee
		// takerFee = quantityToLend*borrowFee/baseFee
		takerFee := new(big.Int).Mul(quantityToLend, borrowFee)
		takerFee = new(big.Int).Div(takerFee, common.TomoXBaseFee)
		if quantityToLend.Cmp(takerFee) <= 0 {
			log.Debug("quantity lending too small", "quantityToLend", quantityToLend, "takerFee", takerFee)
			return result, ErrQuantityTradeTooSmall
		}
		if lendingToken.String() != common.TomoNativeAddress && lendTokenTOMOPrice != nil && lendTokenTOMOPrice.Cmp(common.Big0) > 0 {
			exTakerReceivedFee := new(big.Int).Mul(takerFee, lendTokenTOMOPrice)
			exTakerReceivedFee = new(big.Int).Div(exTakerReceivedFee, lendTokenDecimal)
			log.Debug("exTakerReceivedFee", "quantityToLend", quantityToLend, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "borrowFee", borrowFee)
			if exTakerReceivedFee.Cmp(common.RelayerLendingFee) <= 0 {
				log.Debug("takerFee too small", "quantityToLend", quantityToLend, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "borrowFee", borrowFee)
				return result, ErrQuantityTradeTooSmall
			}
		} else if lendingToken.String() == common.TomoNativeAddress {
			exTakerReceivedFee := takerFee
			log.Debug("exTakerReceivedFee", "quantityToLend", quantityToLend, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "borrowFee", borrowFee)
			if exTakerReceivedFee.Cmp(common.RelayerLendingFee) <= 0 {
				log.Debug("takerFee too small", "quantityToLend", quantityToLend, "takerFee", takerFee, "exTakerReceivedFee", exTakerReceivedFee, "borrowFee", borrowFee)
				return result, ErrQuantityTradeTooSmall
			}
		}
		result = &LendingSettleBalance{
			//Borrower
			Taker: TradeResult{
				Fee:      takerFee,
				InToken:  lendingToken,
				InTotal:  new(big.Int).Sub(quantityToLend, takerFee),
				OutToken: collateralToken,
				OutTotal: takerOutTotal,
			},
			// Investor : makerOutTotal = quantityToLend
			Maker: TradeResult{
				Fee:      common.Big0,
				InToken:  common.Address{},
				InTotal:  common.Big0,
				OutToken: lendingToken,
				OutTotal: quantityToLend,
			},
			CollateralLockedAmount: takerOutTotal,
		}
	} else {
		// maker =  Borrower : makerOutTotal = CollateralLockedAmount = quantityToLend * collateral Token Decimal / CollateralPrice  * deposit rate
		makerOutTotal := new(big.Int).Mul(quantityToLend, collateralTokenDecimal)
		makerOutTotal = new(big.Int).Mul(makerOutTotal, depositRate) // eg: depositRate = 150%
		makerOutTotal = new(big.Int).Div(makerOutTotal, big.NewInt(100))
		makerOutTotal = new(big.Int).Div(makerOutTotal, collateralPrice)
		// Fee
		makerFee := new(big.Int).Mul(quantityToLend, borrowFee)
		makerFee = new(big.Int).Div(makerFee, common.TomoXBaseFee)
		if quantityToLend.Cmp(makerFee) <= 0 {
			log.Debug("quantity lending too small", "quantityToLend", quantityToLend, "makerFee", makerFee)
			return result, ErrQuantityTradeTooSmall
		}
		if lendingToken.String() != common.TomoNativeAddress && lendTokenTOMOPrice != nil && lendTokenTOMOPrice.Cmp(common.Big0) > 0 {
			exMakerReceivedFee := new(big.Int).Mul(makerFee, lendTokenTOMOPrice)
			exMakerReceivedFee = new(big.Int).Div(exMakerReceivedFee, lendTokenDecimal)
			log.Debug("exMakerReceivedFee", "quantityToLend", quantityToLend, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "borrowFee", borrowFee)
			if exMakerReceivedFee.Cmp(common.RelayerLendingFee) <= 0 {
				log.Debug("makerFee too small", "quantityToLend", quantityToLend, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "borrowFee", borrowFee)
				return result, ErrQuantityTradeTooSmall
			}
		} else if lendingToken.String() == common.TomoNativeAddress {
			exMakerReceivedFee := makerFee
			log.Debug("exMakerReceivedFee", "quantityToLend", quantityToLend, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "borrowFee", borrowFee)
			if exMakerReceivedFee.Cmp(common.RelayerLendingFee) <= 0 {
				log.Debug("makerFee too small", "quantityToLend", quantityToLend, "makerFee", makerFee, "exMakerReceivedFee", exMakerReceivedFee, "borrowFee", borrowFee)
				return result, ErrQuantityTradeTooSmall
			}
		}
		result = &LendingSettleBalance{
			Taker: TradeResult{
				Fee:      common.Big0,
				InToken:  common.Address{},
				InTotal:  common.Big0,
				OutToken: lendingToken,
				OutTotal: quantityToLend,
			},
			Maker: TradeResult{
				Fee:      makerFee,
				InToken:  lendingToken,
				InTotal:  new(big.Int).Add(quantityToLend, makerFee),
				OutToken: collateralToken,
				OutTotal: makerOutTotal,
			},
			CollateralLockedAmount: makerOutTotal,
		}
	}
	return result, nil
}
