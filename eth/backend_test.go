package eth

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tomochain/tomochain/params"
)

func TestCalcInitialReward(t *testing.T) {
	for i := 0; i < 100; i++ {
		// the first 2 years
		rewardPerEpoch := new(big.Int).Mul(new(big.Int).SetUint64(250), new(big.Int).SetUint64(params.Ether))
		chainReward := calcInitialReward(rewardPerEpoch, uint64(i), 10)
		if 0 <= i && i < 20 && chainReward.Cmp(rewardPerEpoch) != 0 {
			t.Error("Fail tor calculate reward inflation for 0 -> 2 years", "chainReward", chainReward)
		}

		// 3rd year, 4th year, 5th year
		halfReward := new(big.Int).Mul(new(big.Int).SetUint64(125), new(big.Int).SetUint64(params.Ether))
		if 20 <= i && i < 50 && chainReward.Cmp(halfReward) != 0 {
			t.Error("Fail tor calculate reward inflation for 2 -> 5 years", "chainReward", chainReward)
		}

		// 6th year, 7th year, 8th year
		quarterReward := new(big.Int).Mul(new(big.Int).SetUint64(62.5*1000), new(big.Int).SetUint64(params.Finney))
		if 50 <= i && i < 80 && chainReward.Cmp(quarterReward) != 0 {
			t.Error("Fail tor calculate reward inflation for 5 -> 8 years", "chainReward", chainReward)
		}

		// 8th onwards
		zeroReward := big.NewInt(0)
		if 80 <= i && chainReward.Cmp(zeroReward) != 0 {
			t.Error("Fail tor calculate reward inflation above 8 years", "chainReward", chainReward)
		}
	}
}

func TestCalcSaigonReward(t *testing.T) {
	initialReward := new(big.Int).Mul(new(big.Int).SetUint64(250), new(big.Int).SetUint64(params.Ether))
	firstHalvingReward := new(big.Int).Div(initialReward, big.NewInt(2))
	secondHalvingReward := new(big.Int).Div(initialReward, big.NewInt(4))
	thirdHalvingReward := new(big.Int).Div(initialReward, big.NewInt(8))
	zeroReward := big.NewInt(0)

	for i := 0; i < 150; i++ {
		reward := calcSaigonReward(initialReward, big.NewInt(10), uint64(i), 7)
		if i < 10 {
			assert.True(t, reward.Cmp(zeroReward) == 0, "pre-Saigon Upgrade reward mismatch (%s) at epoch %s", reward, i)
		}
		if i >= 10 && i < 38 {
			assert.True(t, reward.Cmp(initialReward) == 0, "1st cycle reward mismatch (%s) at epoch %s", reward, i)
		}
		if i >= 38 && i < 66 {
			assert.True(t, reward.Cmp(firstHalvingReward) == 0, "2nd cycle reward mismatch (%s) at epoch %s", reward, i)
		}
		if i >= 66 && i < 94 {
			assert.True(t, reward.Cmp(secondHalvingReward) == 0, "3rd cycle reward mismatch (%s) at epoch %s", reward, i)
		}
		if i >= 94 && i < 122 {
			assert.True(t, reward.Cmp(thirdHalvingReward) == 0, "4rd cycle reward mismatch (%s) at epoch %s", reward, i)
		}
		if i >= 122 {
			assert.True(t, reward.Cmp(zeroReward) == 0, "post-Saigon Upgrade reward mismatch (%s) at epoch %s", reward, i)
		}
	}
}
