package tomox

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	trc21 "github.com/ethereum/go-ethereum/contracts/trc21issuer/contract"
	"github.com/ethereum/go-ethereum/ethclient"
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

func (tomox *TomoX) SettleBalance(
	ipcEndpoint string,
	maker, taker common.Address,
	baseToken, quoteToken common.Address,
	isTakerBuy bool,
	makerFeeRate, takerFeeRate, baseFee *big.Int,
	quantity *big.Int,
	price *big.Int,
) (map[common.Address]map[string]interface{}, error) {
	result := map[common.Address]map[string]interface{}{}
	//
	//// pair: BASE_TOKEN / QUOTE_TOKEN
	//// Volume is calculated by quote token
	//// Therefore, baseTokenQuantity = price * quantity
	////				quoteQuantity = quantity
	//// Fee by quoteToken
	//
	baseTokenDecimal, err := tomox.GetTokenDecimal(ipcEndpoint, baseToken)
	if err != nil {
		return nil, fmt.Errorf("Fail to get tokenDecimal. Token: %v . Err: %v", baseToken.String(), err)
	}
	if isTakerBuy {
		// taker InQuantity
		baseTokenQuantity := quantity

		// maker InQuantity
		quoteTokenQuantity := big.NewInt(0).Mul(quantity, price)
		quoteTokenQuantity = big.NewInt(0).Div(quoteTokenQuantity, baseTokenDecimal)

		// Fee
		// charge on the token he/she has before the trade, in this case: quoteToken
		takerFee := big.NewInt(0).Mul(quoteTokenQuantity, takerFeeRate)
		takerFee = big.NewInt(0).Div(takerFee, baseFee)
		// charge on the token he/she has before the trade, in this case: baseToken
		makerFee := big.NewInt(0).Mul(quoteTokenQuantity, makerFeeRate)
		makerFee = big.NewInt(0).Div(makerFee, baseFee)

		result[taker] = map[string]interface{}{
			Fee:         takerFee,
			InToken:     baseToken,
			InQuantity:  baseTokenQuantity,
			InTotal:     baseTokenQuantity,
			OutToken:    quoteToken,
			OutQuantity: quoteTokenQuantity,
			OutTotal:    big.NewInt(0).Add(quoteTokenQuantity, takerFee),
		}

		result[maker] = map[string]interface{}{
			Fee:         makerFee,
			InToken:     quoteToken,
			InQuantity:  quoteTokenQuantity,
			InTotal:     big.NewInt(0).Sub(quoteTokenQuantity, makerFee),
			OutToken:    baseToken,
			OutQuantity: baseTokenQuantity,
			OutTotal:    baseTokenQuantity,
		}
	} else {
		// taker InQuantity
		quoteTokenQuantity := big.NewInt(0).Mul(quantity, price)
		quoteTokenQuantity = big.NewInt(0).Div(quoteTokenQuantity, baseTokenDecimal)
		// maker InQuantity
		baseTokenQuantity := quantity

		// Fee
		// charge on the token he/she has before the trade, in this case: baseToken
		takerFee := big.NewInt(0).Mul(quoteTokenQuantity, takerFeeRate)
		takerFee = big.NewInt(0).Div(takerFee, baseFee)
		// charge on the token he/she has before the trade, in this case: quoteToken
		makerFee := big.NewInt(0).Mul(quoteTokenQuantity, makerFeeRate)
		makerFee = big.NewInt(0).Div(makerFee, baseFee)

		result[taker] = map[string]interface{}{
			Fee:         takerFee,
			InToken:     quoteToken,
			InQuantity:  quoteTokenQuantity,
			InTotal:     big.NewInt(0).Sub(quoteTokenQuantity, takerFee),
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
			OutTotal:    big.NewInt(0).Add(quoteTokenQuantity, makerFee),
		}
	}
	return result, nil
}

func (tomox *TomoX) GetTokenDecimal(ipcEndpoint string,tokenAddr common.Address) (*big.Int, error) {
	if tokenDecimal, ok := tomox.tokenDecimalCache.Get(tokenAddr); ok {
		return tokenDecimal.(*big.Int), nil
	}
	if tokenAddr.String() == common.TomoNativeAddress {
		tomox.tokenDecimalCache.Add(tokenAddr, common.BasePrice)
		return common.BasePrice, nil
	}

	client, err := ethclient.Dial(ipcEndpoint)
	if err != nil {
		return nil, err
	}
	opts := new(bind.CallOpts)
	trc21Contract, err := trc21.NewMyTRC21(tokenAddr, client)
	decimal, err := trc21Contract.Decimals(opts)
	if err != nil {
		return nil, err
	}
	tokenDecimal := new(big.Int).SetUint64(0).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	tomox.tokenDecimalCache.Add(tokenAddr, tokenDecimal)
	return tokenDecimal, nil
}
