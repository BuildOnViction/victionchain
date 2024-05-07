package eth

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/params"
)

func TestRewardInflation(t *testing.T) {
	blocksPerYear := uint64(10)
	params.AllPosvProtocolChanges.SaigonBlock = big.NewInt(65)
	big250VIC := common.InitialSaigonRewardPerEpoch
	big125VIC := new(big.Int).Div(big250VIC, big.NewInt(2))
	big62Point5VIC := new(big.Int).Div(big250VIC, big.NewInt(4))

	/*
		| duration                                    | block range | pre-Saigon + post-Saigon epoch reward |
		|---------------------------------------------|-------------|---------------------------------------|
		| first 2 years                                | 0 -> 19     | 250 VIC + 0 VIC                       |
		| 2nd -> 5th year                             | 20 -> 49    | 125 VIC + 0 VIC                       |
		| 5th year -> Saigon HF                       | 50 -> 64    | 62.5 VIC + 0 VIC                      |
		| Saigon HF -> 8th year                       | 65 -> 79    | 62.5 VIC + 250 VIC                    |
		| 8th year -> Saigon HF + 4 years             | 80 -> 104   | 0 VIC + 250 VIC                       |
		| Saigon HF + 4 years -> Saigon HF + 8 years  | 105 -> 144  | 0 VIC + 125 VIC                       |
		| Saigon HF + 8 years -> Saigon HF + 12 years | 145 -> 184  | 0 VIC + 62.5 VIC                      |
	*/
	for i := uint64(0); i < 185; i++ {
		chainReward := new(big.Int).Add(
			preSaigonEpochReward(params.AllPosvProtocolChanges, i, blocksPerYear),
			postSaigonEpochReward(params.AllPosvProtocolChanges, new(big.Int).SetUint64(i), blocksPerYear),
		)
		switch i {
		case 0:
		case 19:
			assert.Equal(t, 0, chainReward.Cmp(big250VIC), "0 -> 2 years reward mismatch")
		case 20:
		case 49:
			assert.Equal(t, 0, chainReward.Cmp(big125VIC), "2 -> 5 years reward mismatch")
		case 50:
		case 64:
			assert.Equal(t, 0, chainReward.Cmp(big62Point5VIC), "5 years -> before SaigonBlock reward mismatch")
		case 65:
		case 79:
			assert.Equal(t, 0, chainReward.Cmp(new(big.Int).Add(big62Point5VIC, big250VIC)), "SaigonBlock -> 8 years reward mismatch")
		case 80:
		case 104:
			assert.Equal(t, 0, chainReward.Cmp(big250VIC), "8 years -> SaigonBlock + 4 years reward mismatch")
		case 105:
		case 144:
			assert.Equal(t, 0, chainReward.Cmp(big125VIC), "SaigonBlock + 4 years  -> SaigonBlock + 8 years reward mismatch")
		case 145:
		case 184:
			assert.Equal(t, 0, chainReward.Cmp(big62Point5VIC), "SaigonBlock + 8 years  -> SaigonBlock + 12 years reward mismatch")
		}
	}
}
