package tomox

import (
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/tomox/tradingstate"
	"math/big"
	"reflect"
	"testing"
)

func Test_getCancelFee(t *testing.T) {
	type CancelFeeArg struct {
		baseTokenDecimal *big.Int
		feeRate          *big.Int
		order            *tradingstate.OrderItem
	}
	tests := []struct {
		name string
		args CancelFeeArg
		want *big.Int
	}{
		// zero fee test: SELL
		{
			"zero fee test: SELL",
			CancelFeeArg{
				baseTokenDecimal: common.Big1,
				feeRate:          common.Big0,
				order: &tradingstate.OrderItem{
					Quantity: new(big.Int).SetUint64(10000),
					Side:     tradingstate.Ask,
				},
			},
			common.Big0,
		},

		// zero fee test: BUY
		{
			"zero fee test: BUY",
			CancelFeeArg{
				baseTokenDecimal: common.Big1,
				feeRate:          common.Big0,
				order: &tradingstate.OrderItem{
					Quantity: new(big.Int).SetUint64(10000),
					Price:    new(big.Int).SetUint64(1),
					Side:     tradingstate.Bid,
				},
			},
			common.Big0,
		},

		// test getCancelFee: SELL
		{
			"test getCancelFee:: SELL",
			CancelFeeArg{
				baseTokenDecimal: common.Big1,
				feeRate:          new(big.Int).SetUint64(10), // 10/1000 = 0.1%
				order: &tradingstate.OrderItem{
					Quantity: new(big.Int).SetUint64(10000),
					Side:     tradingstate.Ask,
				},
			},
			common.Big1,
		},

		// test getCancelFee:: BUY
		{
			"test getCancelFee:: BUY",
			CancelFeeArg{
				baseTokenDecimal: common.Big1,
				feeRate:          new(big.Int).SetUint64(10), // 10/1000 = 0.1%
				order: &tradingstate.OrderItem{
					Quantity: new(big.Int).SetUint64(10000),
					Price:    new(big.Int).SetUint64(1),
					Side:     tradingstate.Bid,
				},
			},
			common.Big1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCancelFee(tt.args.baseTokenDecimal, tt.args.feeRate, tt.args.order); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCancelFee() = %v, want %v", got, tt.want)
			}
		})
	}
}
