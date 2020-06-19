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
				t.Errorf("getCancelFee() = %v, quantity %v", got, tt.want)
			}
		})
	}
}

func TestGetTradeQuantity(t *testing.T) {
	type GetTradeQuantityArg struct {
		takerSide        string
		takerFeeRate     *big.Int
		takerBalance     *big.Int
		makerPrice       *big.Int
		makerFeeRate     *big.Int
		makerBalance     *big.Int
		baseTokenDecimal *big.Int
		quantityToTrade  *big.Int
	}
	tests := []struct {
		name        string
		args        GetTradeQuantityArg
		quantity    *big.Int
		rejectMaker bool
	}{
		{
			"BUY: feeRate = 0, price 1, quantity 1000, taker balance 1000, maker balance 1000",
			GetTradeQuantityArg{
				takerSide: tradingstate.Bid,
				takerFeeRate: common.Big0,
				takerBalance: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			false,
		},
		{
			"BUY: feeRate = 0, price 1, quantity 1000, taker balance 1000, maker balance 900 -> reject maker",
			GetTradeQuantityArg{
				takerSide: tradingstate.Bid,
				takerFeeRate: common.Big0,
				takerBalance: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: new(big.Int).Mul(big.NewInt(900), common.BasePrice),
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			new(big.Int).Mul(big.NewInt(900), common.BasePrice),
			true,
		},
		{
			"BUY: feeRate = 0, price 1, quantity 1000, taker balance 900, maker balance 1000 -> reject taker",
			GetTradeQuantityArg{
				takerSide: tradingstate.Bid,
				takerFeeRate: common.Big0,
				takerBalance: new(big.Int).Mul(big.NewInt(900), common.BasePrice),
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			new(big.Int).Mul(big.NewInt(900), common.BasePrice),
			false,
		},
		{
			"BUY: feeRate = 0, price 1, quantity 1000, taker balance 0, maker balance 1000 -> reject taker",
			GetTradeQuantityArg{
				takerSide: tradingstate.Bid,
				takerFeeRate: common.Big0,
				takerBalance: common.Big0,
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			common.Big0,
			false,
		},
		{
			"BUY: feeRate = 0, price 1, quantity 1000, taker balance 0, maker balance 0 -> reject both taker",
			GetTradeQuantityArg{
				takerSide: tradingstate.Bid,
				takerFeeRate: common.Big0,
				takerBalance: common.Big0,
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: common.Big0,
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			common.Big0,
			false,
		},
		{
			"BUY: feeRate = 0, price 1, quantity 1000, taker balance 500, maker balance 100 -> reject both taker, maker",
			GetTradeQuantityArg{
				takerSide: tradingstate.Bid,
				takerFeeRate: common.Big0,
				takerBalance: new(big.Int).Mul(big.NewInt(500), common.BasePrice),
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: new(big.Int).Mul(big.NewInt(100), common.BasePrice),
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			new(big.Int).Mul(big.NewInt(100), common.BasePrice),
			true,
		},




		{
			"SELL: feeRate = 0, price 1, quantity 1000, taker balance 1000, maker balance 1000",
			GetTradeQuantityArg{
				takerSide: tradingstate.Ask,
				takerFeeRate: common.Big0,
				takerBalance: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			false,
		},
		{
			"SELL: feeRate = 0, price 1, quantity 1000, taker balance 1000, maker balance 900 -> reject maker",
			GetTradeQuantityArg{
				takerSide: tradingstate.Ask,
				takerFeeRate: common.Big0,
				takerBalance: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: new(big.Int).Mul(big.NewInt(900), common.BasePrice),
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			new(big.Int).Mul(big.NewInt(900), common.BasePrice),
			true,
		},
		{
			"SELL: feeRate = 0, price 1, quantity 1000, taker balance 900, maker balance 1000 -> reject taker",
			GetTradeQuantityArg{
				takerSide: tradingstate.Ask,
				takerFeeRate: common.Big0,
				takerBalance: new(big.Int).Mul(big.NewInt(900), common.BasePrice),
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			new(big.Int).Mul(big.NewInt(900), common.BasePrice),
			false,
		},
		{
			"SELL: feeRate = 0, price 1, quantity 1000, taker balance 0, maker balance 1000 -> reject taker",
			GetTradeQuantityArg{
				takerSide: tradingstate.Ask,
				takerFeeRate: common.Big0,
				takerBalance: common.Big0,
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			common.Big0,
			false,
		},
		{
			"SELL: feeRate = 0, price 1, quantity 1000, taker balance 0, maker balance 0 -> reject maker",
			GetTradeQuantityArg{
				takerSide: tradingstate.Ask,
				takerFeeRate: common.Big0,
				takerBalance: common.Big0,
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: common.Big0,
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			common.Big0,
			true,
		},
		{
			"SELL: feeRate = 0, price 1, quantity 1000, taker balance 500, maker balance 100 -> reject both taker, maker",
			GetTradeQuantityArg{
				takerSide: tradingstate.Ask,
				takerFeeRate: common.Big0,
				takerBalance: new(big.Int).Mul(big.NewInt(500), common.BasePrice),
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: new(big.Int).Mul(big.NewInt(100), common.BasePrice),
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			new(big.Int).Mul(big.NewInt(100), common.BasePrice),
			true,
		},
		{
			"SELL: feeRate = 0, price 1, quantity 1000, taker balance 0, maker balance 100 -> reject both taker, maker",
			GetTradeQuantityArg{
				takerSide: tradingstate.Ask,
				takerFeeRate: common.Big0,
				takerBalance: common.Big0,
				makerPrice: common.BasePrice,
				makerFeeRate: common.Big0,
				makerBalance: new(big.Int).Mul(big.NewInt(100), common.BasePrice),
				baseTokenDecimal: common.BasePrice,
				quantityToTrade: new(big.Int).Mul(big.NewInt(1000), common.BasePrice),
			},
			common.Big0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetTradeQuantity(tt.args.takerSide, tt.args.takerFeeRate, tt.args.takerBalance, tt.args.makerPrice, tt.args.makerFeeRate, tt.args.makerBalance, tt.args.baseTokenDecimal, tt.args.quantityToTrade)
			if !reflect.DeepEqual(got, tt.quantity) {
				t.Errorf("GetTradeQuantity() got = %v, quantity %v", got, tt.quantity)
			}
			if got1 != tt.rejectMaker {
				t.Errorf("GetTradeQuantity() got1 = %v, quantity %v", got1, tt.rejectMaker)
			}
		})
	}
}
