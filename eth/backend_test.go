package eth

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/params"
)

func TestRewardInflation(t *testing.T) {
	params.AllPosvProtocolChanges.SaigonBlock = big.NewInt(60)
	baseSaigonBlockRewardPerEpoch := new(big.Int).Div(common.InitialSaigonBlockRewardPerYear, big.NewInt(int64(common.EpochPerYear)))
	// the first 2 years
	initialBlockRewardPerEpoch := new(big.Int).Mul(new(big.Int).SetUint64(params.AllPosvProtocolChanges.Posv.Reward), new(big.Int).SetUint64(params.Ether))
	// 3rd year, 4th year, 5th year
	halfReward := new(big.Int).Mul(new(big.Int).SetUint64(125), new(big.Int).SetUint64(params.Ether))
	// after 5 years and before SaigonBlock
	quarterReward := new(big.Int).Mul(new(big.Int).SetUint64(62.5*1000), new(big.Int).SetUint64(params.Finney))
	// first 4 years after SaigonBlock
	thirdHalvingReward := new(big.Int).Div(baseSaigonBlockRewardPerEpoch, big.NewInt(1))
	// next 4 years
	fourthHalvingReward := new(big.Int).Div(baseSaigonBlockRewardPerEpoch, big.NewInt(2))
	// next 4 years
	fifthHalvingReward := new(big.Int).Div(baseSaigonBlockRewardPerEpoch, big.NewInt(4))
	// next 4 years
	sixthHalvingReward := new(big.Int).Div(baseSaigonBlockRewardPerEpoch, big.NewInt(8))

	for i := int64(0); i < 200; i++ {
		chainReward := rewardInflation(params.AllPosvProtocolChanges, currentChainReward(params.AllPosvProtocolChanges, big.NewInt(i)), uint64(i), 10)
		switch i {
		case 0:
		case 19:
			assert.Equal(t, 0, chainReward.Cmp(initialBlockRewardPerEpoch), "0 -> 2 years reward mismatch",
				"chainReward", chainReward, "initialBlockRewardPerEpoch", initialBlockRewardPerEpoch)
		case 20:
		case 49:
			assert.Equal(t, 0, chainReward.Cmp(halfReward), "2 -> 5 years reward mismatch",
				"chainReward", chainReward, "halfReward", halfReward)
		case 50:
		case 59:
			assert.Equal(t, 0, chainReward.Cmp(quarterReward), "5 years -> before SaigonBlock reward mismatch",
				"chainReward", chainReward, "quarterReward", quarterReward)
		case 60:
		case 99:
			assert.Equal(t, 0, chainReward.Cmp(thirdHalvingReward), "SaigonBlock -> next 4 years reward mismatch",
				"chainReward", chainReward, "thirdHalvingReward", thirdHalvingReward)
		case 100:
		case 139:
			assert.Equal(t, 0, chainReward.Cmp(fourthHalvingReward), "SaigonBlock -> next 8 years reward mismatch",
				"chainReward", chainReward, "fourthHalvingReward", fourthHalvingReward)
		case 140:
		case 179:
			assert.Equal(t, 0, chainReward.Cmp(fifthHalvingReward), "SaigonBlock -> next 12 years reward mismatch",
				"chainReward", chainReward, "fifthHalvingReward", fifthHalvingReward)
		case 180:
		case 199:
			assert.Equal(t, 0, chainReward.Cmp(sixthHalvingReward), "SaigonBlock -> next 16 years reward mismatch",
				"chainReward", chainReward, "sixthHalvingReward", sixthHalvingReward)
		}
	}
}
