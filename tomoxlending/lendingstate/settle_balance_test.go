package lendingstate

import (
	"github.com/tomochain/tomochain/common"
	"math/big"
	"reflect"
	"testing"
)

func TestCalculateInterestRate(t *testing.T) {
	type args struct {
		finalizeTime    uint64
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
		// finalize after one day
		// interestRate = 1 * 10 * interestDecimal / 365 = 10 * 1e8 / 365 = 2739726
		// early finalize in the first half => interestRate = interestRate / 2 = 1369863
		// mean 1369863 / 1e8 =  0,01369863 %
		{
			"early finalize in the first half",
			args{
				finalizeTime:    86400,
				liquidationTime: common.OneYear,
				term:            common.OneYear,
				apr:             10 * 1e8,
			},
			new(big.Int).SetUint64(1369863),
		},

		// apr = 10% per year (365 days)
		// term: 365 days
		// finalize at the end
		// pay full interestRate 10%
		{
			"finalize at the end, term : 365 days",
			args{
				finalizeTime:    common.OneYear,
				liquidationTime: common.OneYear,
				term:            common.OneYear,
				apr:             10 * 1e8,
			},
			new(big.Int).SetUint64(10 * 1e8),
		},

		// apr = 10% per year (365 days)
		// term: 30 days
		// finalize after 15 days
		// pay a half of interestRate 10% for 15 days / 365 days
		// interestRate = 10% * 15 /365 / 2 = 0,41095890 % / 2 = 0,20547945 %
		{
			"term: 30 days, finalize after 15 days",
			args{
				finalizeTime:    15 * 86400,
				liquidationTime: 30 * 86400,
				term:            30 * 86400,
				apr:             10 * 1e8,
			},
			new(big.Int).SetUint64(20547945),
		},

		// apr = 10% per year (365 days)
		// term: 30 days
		// finalize at the end
		// pay full interestRate 10% for 30 days / 365 days
		// interestRate = 10% * 30 /365 = 0,821917808 %
		{
			"finalize at the end, term: 30 days",
			args{
				finalizeTime:    30 * 86400,
				liquidationTime: 30 * 86400,
				term:            30 * 86400,
				apr:             10 * 1e8,
			},
			new(big.Int).SetUint64(82191780),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateInterestRate(tt.args.finalizeTime, tt.args.liquidationTime, tt.args.term, tt.args.apr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateInterestRate() = %v, want %v", got, tt.want)
			}
		})
	}
}
