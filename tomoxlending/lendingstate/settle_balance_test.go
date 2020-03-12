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
