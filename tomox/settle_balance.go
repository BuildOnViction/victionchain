package tomox

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	Fee         = "Fee"         // Fee is calculated in quoteToken
	InToken     = "InToken"     // type of token which user is received as the result of the trade
	InQuantity  = "InQuantity"  // amount of token which user is received as the result of the trade, not include fee
	InTotal     = "InTotal"     // amount of token which user is received as the result of the trade, include fee
	OutToken    = "OutToken"    // type of token which user sends out to the partner as the result of the trade
	OutQuantity = "OutQuantity" // amount of token which user sends out to the partner as the result of the trade, not include fee
	OutTotal    = "OutTotal"    // amount of token which user sends out to the partner as the result of the trade, include fee
)

func SettleBalance(
	maker, taker common.Address,
	baseToken, quoteToken common.Address,
	isTakerBuy bool,
	makerFeeRate, takerFeeRate, baseFee *big.Int,
	quantity *big.Int,
	price *big.Int,
) map[common.Address]map[string]interface{} {
	result := map[common.Address]map[string]interface{}{}
	//
	//// pair: BASE_TOKEN / QUOTE_TOKEN
	//// Volume is calculated by quote token
	//// Therefore, baseTokenQuantity = price * quantity
	////				quoteQuantity = quantity
	//// Fee by quoteToken
	//
	if isTakerBuy {
		// taker InQuantity
		baseTokenQuantity := quantity

		// maker InQuantity
		quoteTokenQuantity := new(big.Int).SetUint64(0).Mul(quantity, price)

		// Fee
		// charge on the token he/she has before the trade, in this case: quoteToken
		takerFee := new(big.Int).SetUint64(0).Mul(quoteTokenQuantity, takerFeeRate)
		takerFee = takerFee.Div(takerFee, baseFee)
		// charge on the token he/she has before the trade, in this case: baseToken
		makerFee := new(big.Int).SetUint64(0).Mul(quoteTokenQuantity, makerFeeRate)
		makerFee = makerFee.Div(makerFee, baseFee)

		result[taker] = map[string]interface{}{
			Fee:         takerFee,
			InToken:     baseToken,
			InQuantity:  baseTokenQuantity,
			InTotal:     baseTokenQuantity,
			OutToken:    quoteToken,
			OutQuantity: quoteTokenQuantity,
			OutTotal:    new(big.Int).SetUint64(0).Add(quoteTokenQuantity, takerFee),
		}

		result[maker] = map[string]interface{}{
			Fee:         makerFee,
			InToken:     quoteToken,
			InQuantity:  quoteTokenQuantity,
			InTotal:     new(big.Int).SetUint64(0).Sub(quoteTokenQuantity, makerFee),
			OutToken:    baseToken,
			OutQuantity: baseTokenQuantity,
			OutTotal:    baseTokenQuantity,
		}
	} else {
		// taker InQuantity
		quoteTokenQuantity := new(big.Int).SetUint64(0).Mul(quantity, price)
		// maker InQuantity
		baseTokenQuantity := quantity

		// Fee
		// charge on the token he/she has before the trade, in this case: baseToken
		takerFee := new(big.Int).SetUint64(0).Mul(quoteTokenQuantity, takerFeeRate)
		takerFee = takerFee.Div(takerFee, baseFee)
		// charge on the token he/she has before the trade, in this case: quoteToken
		makerFee := new(big.Int).SetUint64(0).Mul(quoteTokenQuantity, makerFeeRate)
		makerFee = makerFee.Div(makerFee, baseFee)

		result[taker] = map[string]interface{}{
			Fee:         takerFee,
			InToken:     quoteToken,
			InQuantity:  quoteTokenQuantity,
			InTotal:     new(big.Int).SetUint64(0).Sub(quoteTokenQuantity, takerFee),
			OutToken:    baseToken,
			OutQuantity: baseTokenQuantity,
			OutTotal:    baseTokenQuantity,
		}

		result[maker] = map[string]interface{}{
			Fee:         makerFee,
			InToken:     baseToken,
			InQuantity:  baseTokenQuantity,
			InTotal:     baseTokenQuantity,
			OutToken:    quoteToken,
			OutQuantity: quoteTokenQuantity,
			OutTotal:    new(big.Int).SetUint64(0).Add(quoteTokenQuantity, makerFee),
		}
	}
	return result
}
