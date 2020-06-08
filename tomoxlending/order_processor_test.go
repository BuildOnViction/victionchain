package tomoxlending

import (
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/tomox/tradingstate"
	"github.com/tomochain/tomochain/tomoxlending/lendingstate"
	"math/big"
	"reflect"
	"testing"
)

func Test_getCancelFee(t *testing.T) {
	type CancelFeeArg struct {
		collateralTokenDecimal *big.Int
		collateralPrice        *big.Int
		borrowFeeRate          *big.Int
		order                  *lendingstate.LendingItem
	}
	tests := []struct {
		name string
		args CancelFeeArg
		want *big.Int
	}{
		// zero fee test: LEND
		{
			"zero fee test: LEND",
			CancelFeeArg{
				collateralTokenDecimal: common.Big1,
				collateralPrice:        common.Big1,
				borrowFeeRate:          common.Big0,
				order: &lendingstate.LendingItem{
					Quantity: new(big.Int).SetUint64(10000),
					Side:     tradingstate.Ask,
				},
			},
			common.Big0,
		},

		// zero fee test: BORROW
		{
			"zero fee test: BORROW",
			CancelFeeArg{
				collateralTokenDecimal: common.Big1,
				collateralPrice:        common.Big1,
				borrowFeeRate:          common.Big0,
				order: &lendingstate.LendingItem{
					Quantity: new(big.Int).SetUint64(10000),
					Side:     tradingstate.Bid,
				},
			},
			common.Big0,
		},

		// test getCancelFee: LEND
		{
			"test getCancelFee:: LEND",
			CancelFeeArg{
				collateralTokenDecimal: common.Big1,
				collateralPrice:        common.Big1,
				borrowFeeRate:          new(big.Int).SetUint64(30), // 30/1000 = 0.3%
				order: &lendingstate.LendingItem{
					Quantity: new(big.Int).SetUint64(10000),
					Side:     tradingstate.Ask,
				},
			},
			common.Big3,
		},

		// test getCancelFee:: BORROW
		{
			"test getCancelFee:: BORROW",
			CancelFeeArg{
				collateralTokenDecimal: common.Big1,
				collateralPrice:        common.Big1,
				borrowFeeRate:          new(big.Int).SetUint64(30), // 30/1000 = 0.3%
				order: &lendingstate.LendingItem{
					Quantity: new(big.Int).SetUint64(10000),
					Side:     tradingstate.Bid,
				},
			},
			common.Big3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCancelFee(tt.args.collateralTokenDecimal, tt.args.collateralPrice, tt.args.borrowFeeRate, tt.args.order); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCancelFee() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLendQuantity(t *testing.T) {
	depositRate := big.NewInt(150)
	lendQuantity := new(big.Int).Mul(big.NewInt(1000), common.BasePrice)
	collateralLocked, _ := new(big.Int).SetString("1000000000000000000000", 10) // 1000
	collateralLocked = new(big.Int).Mul(big.NewInt(150), collateralLocked)
	collateralLocked = new(big.Int).Div(collateralLocked, big.NewInt(100))
	type GetLendQuantityArg struct {
		takerSide              string
		collateralTokenDecimal *big.Int
		depositRate            *big.Int
		collateralPrice        *big.Int
		takerBalance           *big.Int
		makerBalance           *big.Int
		quantityToLend         *big.Int
	}
	tests := []struct {
		name  string
		args  GetLendQuantityArg
		lendQuantity  *big.Int
		rejectMaker bool
	}{
		{
			"taker: BORROW, takerBalance = 0, reject taker",
			GetLendQuantityArg{
				lendingstate.Borrowing,
				common.BasePrice,
				depositRate,
				common.BasePrice,
				common.Big0,
				common.Big0,
				lendQuantity,
			},
			common.Big0,
			false,
		},
		{
			"taker: BORROW, takerBalance not enough, reject partial of taker",
			GetLendQuantityArg{
				lendingstate.Borrowing,
				common.BasePrice,
				depositRate,
				common.BasePrice,
				new(big.Int).Div(collateralLocked, big.NewInt(2)), // 1/2
				lendQuantity,
				lendQuantity,
			},
			new(big.Int).Div(lendQuantity, big.NewInt(2)),
			false,
		},
		{
			"taker: BORROW, makerBalance = 0, reject maker",
			GetLendQuantityArg{
				lendingstate.Borrowing,
				common.BasePrice,
				depositRate,
				common.BasePrice,
				new(big.Int).Div(collateralLocked, big.NewInt(2)),
				common.Big0,
				lendQuantity,
			},
			common.Big0,
			true,
		},
		{
			"taker: BORROW, makerBalance not enough, reject partial of maker",
			GetLendQuantityArg{
				lendingstate.Borrowing,
				common.BasePrice,
				depositRate,
				common.BasePrice,
				collateralLocked,
				new(big.Int).Div(lendQuantity, big.NewInt(2)),
				lendQuantity,
			},
			new(big.Int).Div(lendQuantity, big.NewInt(2)),
			true,
		},
		{
			"taker: BORROW, don't reject",
			GetLendQuantityArg{
				lendingstate.Borrowing,
				common.BasePrice,
				depositRate,
				common.BasePrice,
				collateralLocked,
				lendQuantity,
				lendQuantity,
			},
			lendQuantity,
			false,
		},

		{
			"taker: INVEST, makerBalance = 0, reject maker",
			GetLendQuantityArg{
				lendingstate.Investing,
				common.BasePrice,
				depositRate,
				common.BasePrice,
				new(big.Int).Div(collateralLocked, big.NewInt(2)),
				common.Big0,
				lendQuantity,
			},
			common.Big0,
			true,
		},
		{
			"taker: INVEST, takerBalance not enough, reject partial of taker",
			GetLendQuantityArg{
				lendingstate.Investing,
				common.BasePrice,
				depositRate,
				common.BasePrice,
				new(big.Int).Div(lendQuantity, big.NewInt(2)), // 1/2
				collateralLocked,
				lendQuantity,
			},
			new(big.Int).Div(lendQuantity, big.NewInt(2)),
			false,
		},
		{
			"taker: INVEST, makerBalance = 0, reject maker",
			GetLendQuantityArg{
				lendingstate.Investing,
				common.BasePrice,
				depositRate,
				common.BasePrice,
				common.Big0,
				new(big.Int).Div(collateralLocked, big.NewInt(2)),
				lendQuantity,
			},
			common.Big0,
			false,
		},
		{
			"taker: INVEST, makerBalance not enough, reject partial of maker",
			GetLendQuantityArg{
				lendingstate.Investing,
				common.BasePrice,
				depositRate,
				common.BasePrice,
				collateralLocked,
				new(big.Int).Div(collateralLocked, big.NewInt(2)),
				lendQuantity,
			},
			new(big.Int).Div(lendQuantity, big.NewInt(2)),
			true,
		},
		{
			"taker: INVEST, don't reject",
			GetLendQuantityArg{
				lendingstate.Investing,
				common.BasePrice,
				depositRate,
				common.BasePrice,
				lendQuantity,
				collateralLocked,
				lendQuantity,
			},
			lendQuantity,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetLendQuantity(tt.args.takerSide, tt.args.collateralTokenDecimal, tt.args.depositRate, tt.args.collateralPrice, tt.args.takerBalance, tt.args.makerBalance, tt.args.quantityToLend)
			if !reflect.DeepEqual(got, tt.lendQuantity) {
				t.Errorf("GetLendQuantity() got = %v, want %v", got, tt.lendQuantity)
			}
			if got1 != tt.rejectMaker {
				t.Errorf("GetLendQuantity() got1 = %v, want %v", got1, tt.rejectMaker)
			}
		})
	}
}
