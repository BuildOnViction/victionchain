package lendingstate

import (
	"github.com/tomochain/tomochain/common"
	"math/big"
	"reflect"
	"testing"
)

func TestCalculateInterestRate(t *testing.T) {
	type args struct {
		repayTime       uint64
		liquidationTime uint64
		term            uint64
		apr             uint64
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		// apr = 10% per year
		// term 365 days
		// repay after one day
		// have to pay interest for a half of year
		// I = APR *(T + T1) / 2 / 365 = 10% * (365 + 1) / 2 /365 = 5,01369863 %
		// 1e8 is decimal of interestRate
		{
			"term 365 days: early repay",
			args{
				repayTime:       86400,
				liquidationTime: common.OneYear,
				term:            common.OneYear,
				apr:             10 * 1e8,
			},
			new(big.Int).SetUint64(501369863),
		},

		// apr = 10% per year (365 days)
		// term: 365 days
		// repay at the end
		// pay full interestRate 10%
		// I = APR *(T + T1) / 2 / 365 = 10% * (365 + 365) / 2 /365 = 10 %
		// 1e8 is decimal of interestRate
		{
			"term 365 days: repay at the end",
			args{
				repayTime:       common.OneYear,
				liquidationTime: common.OneYear,
				term:            common.OneYear,
				apr:             10 * 1e8,
			},
			new(big.Int).SetUint64(10 * 1e8),
		},

		// apr = 10% per year
		// term 30 days
		// repay after one day
		// have to pay interest for a half of year
		// I = APR *(T + T1) / 2 / 365 = 10% * (30 + 1) / 2 /365 = 0,424657534 %
		// 1e8 is decimal of interestRate
		{
			"term 30 days: early repay",
			args{
				repayTime:       86400,
				liquidationTime: 30 * 86400,
				term:            30 * 86400,
				apr:             10 * 1e8,
			},
			new(big.Int).SetUint64(42465753),
		},

		// apr = 10% per year (365 days)
		// term: 30 days
		// repay at the end
		// pay full interestRate 10%
		// I = APR *(T + T1) / 2 / 365 = 10% * (30 + 30) / 2 /365 = 0,821917808 %
		// 1e8 is decimal of interestRate
		{
			"term 30 days: repay at the end",
			args{
				repayTime:       30 * 86400,
				liquidationTime: 30 * 86400,
				term:            30 * 86400,
				apr:             10 * 1e8,
			},
			new(big.Int).SetUint64(82191780),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateInterestRate(tt.args.repayTime, tt.args.liquidationTime, tt.args.term, tt.args.apr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateInterestRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSettleBalance(t *testing.T) {
	lendQuantity, _ := new(big.Int).SetString("1000000000000000000000", 10)        // 1000
	fee, _ := new(big.Int).SetString("10000000000000000000", 10)                   // 10
	lendQuantityExcluded, _ := new(big.Int).SetString("990000000000000000000", 10) // 990
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000022")
	collateralLocked, _ := new(big.Int).SetString("1000000000000000000000", 10) // 1000
	collateralLocked = new(big.Int).Mul(big.NewInt(150), collateralLocked)
	collateralLocked = new(big.Int).Div(collateralLocked, big.NewInt(100))

	type GetSettleBalanceArg struct {
		isTomoXLendingFork     bool
		takerSide              string
		lendTokenTOMOPrice     *big.Int
		collateralPrice        *big.Int
		depositRate            *big.Int
		borrowFeeRate          *big.Int
		lendingToken           common.Address
		collateralToken        common.Address
		lendTokenDecimal       *big.Int
		collateralTokenDecimal *big.Int
		quantityToLend         *big.Int
	}
	tests := []struct {
		name    string
		args    GetSettleBalanceArg
		want    *LendingSettleBalance
		wantErr bool
	}{
		{
			"quantityToLend = borrowFee",
			GetSettleBalanceArg{
				true,
				Borrowing,
				common.BasePrice,
				common.BasePrice,
				big.NewInt(150),
				big.NewInt(10000), // 100%
				common.Address{},
				common.Address{},
				common.BasePrice,
				common.BasePrice,
				lendQuantity,
			},
			nil,
			true,
		},

		{
			"LendToken is TOMO, quantity too small",
			GetSettleBalanceArg{
				true,
				Borrowing,
				common.BasePrice,
				common.BasePrice,
				big.NewInt(150),
				big.NewInt(100), // 1%
				common.HexToAddress(common.TomoNativeAddress),
				common.Address{},
				common.BasePrice,
				common.BasePrice,
				common.BasePrice,
			},
			nil,
			true,
		},
		{
			"LendToken is not TOMO, quantity too small",
			GetSettleBalanceArg{
				true,
				Borrowing,
				common.BasePrice,
				common.BasePrice,
				big.NewInt(150),
				big.NewInt(100), // 1%
				common.Address{},
				common.Address{},
				common.BasePrice,
				common.BasePrice,
				common.BasePrice,
			},
			nil,
			true,
		},

		{
			"LendToken is not TOMO, no error",
			GetSettleBalanceArg{
				true,
				Borrowing,
				common.BasePrice,
				common.BasePrice,
				big.NewInt(150),
				big.NewInt(100), // 1%
				common.Address{},
				common.Address{},
				common.BasePrice,
				common.BasePrice,
				common.BasePrice,
			},
			nil,
			true,
		},

		{
			"LendToken is TOMO, no error, invest",
			GetSettleBalanceArg{
				true,
				Investing,
				common.BasePrice,
				common.BasePrice,
				big.NewInt(150),
				big.NewInt(100), // 1%
				common.HexToAddress(common.TomoNativeAddress),
				collateral,
				common.BasePrice,
				common.BasePrice,
				lendQuantity,
			},
			&LendingSettleBalance{
				Taker: TradeResult{
					Fee:      common.Big0,
					InToken:  common.Address{},
					InTotal:  common.Big0,
					OutToken: common.HexToAddress(common.TomoNativeAddress),
					OutTotal: lendQuantity,
				},
				Maker: TradeResult{
					Fee:      fee,
					InToken:  common.HexToAddress(common.TomoNativeAddress),
					InTotal:  lendQuantityExcluded,
					OutToken: collateral,
					OutTotal: collateralLocked,
				},
				CollateralLockedAmount: collateralLocked,
			},
			false,
		},

		{
			"LendToken is TOMO, no error, Borrow",
			GetSettleBalanceArg{
				true,
				Borrowing,
				common.BasePrice,
				common.BasePrice,
				big.NewInt(150),
				big.NewInt(100), // 1%
				common.HexToAddress(common.TomoNativeAddress),
				collateral,
				common.BasePrice,
				common.BasePrice,
				lendQuantity,
			},
			&LendingSettleBalance{
				Maker: TradeResult{
					Fee:      common.Big0,
					InToken:  common.Address{},
					InTotal:  common.Big0,
					OutToken: common.HexToAddress(common.TomoNativeAddress),
					OutTotal: lendQuantity,
				},
				Taker: TradeResult{
					Fee:      fee,
					InToken:  common.HexToAddress(common.TomoNativeAddress),
					InTotal:  lendQuantityExcluded,
					OutToken: collateral,
					OutTotal: collateralLocked,
				},
				CollateralLockedAmount: collateralLocked,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSettleBalance(tt.args.isTomoXLendingFork, tt.args.takerSide, tt.args.lendTokenTOMOPrice, tt.args.collateralPrice, tt.args.depositRate, tt.args.borrowFeeRate, tt.args.lendingToken, tt.args.collateralToken, tt.args.lendTokenDecimal, tt.args.collateralTokenDecimal, tt.args.quantityToLend)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSettleBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				t.Log(tt.want.String())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSettleBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}
