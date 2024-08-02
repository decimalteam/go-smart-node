package types

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	lastBlock      = 46_656_000
	firstReward    = 500
	firstOldReward = 50

	firstIncrease          = 10
	firstOldIncrease       = 5
	blockStartCalcEmission = 21837188
	blockStartEmission     = 10214628091
)

func GetAllEmission(ctx sdk.Context) sdk.Int {
	allEmision := helpers.EtherToWei(sdk.NewInt(blockStartEmission))
	if !helpers.IsMainnet(ctx.ChainID()) {
		allEmision = helpers.EtherToWei(sdk.NewInt(1))
		for j := uint64(1); j < uint64(ctx.BlockHeight()); j++ {
			allEmision = allEmision.Add(GetRewardOldForBlock(j))
		}
	} else {
		for j := uint64(blockStartCalcEmission); j < uint64(ctx.BlockHeight()); j++ {
			allEmision = allEmision.Add(GetRewardOldForBlock(j))
		}
	}

	return allEmision
}

func GetRewardForBlock(blockHeight uint64) sdk.Int {
	if blockHeight >= lastBlock {
		return sdk.NewInt(0)
	}

	reward := sdk.NewInt(firstReward)
	rewardIncrease := sdk.NewInt(firstIncrease)

	reward = reward.Add(sdk.NewInt(int64(blockHeight / 432000)).Mul(rewardIncrease))
	return helpers.BipToPip(reward)
}

func GetRewardOldForBlock(blockHeight uint64) sdk.Int {
	if blockHeight >= lastBlock {
		return sdk.NewInt(0)
	}

	reward := sdk.NewInt(firstOldReward)
	rewardIncrease := sdk.NewInt(firstOldIncrease)

	if blockHeight/5184000 == 0 {
		reward = reward.Add(sdk.NewInt(int64(blockHeight / 432000)).Mul(rewardIncrease))
		return helpers.BipToPip(reward)
	}

	reward = reward.Add(sdk.NewInt(11).Mul(rewardIncrease))
	for i := uint64(1); i <= blockHeight/5184000; i++ {
		if blockHeight/5184000-i >= 1 {
			rewardIncrease = rewardIncrease.Add(sdk.NewInt(12))
			reward = reward.Add(sdk.NewInt(12).Mul(rewardIncrease))
			continue
		}
		rewardIncrease = rewardIncrease.Add(sdk.NewInt(12))
		//log.Println(reward, rewardIncrease.String(), ((blockHeight-i*5184000)/432000%12)+1)
		reward = reward.Add(sdk.NewInt(int64((blockHeight-i*5184000)/432000%12) + 1).Mul(rewardIncrease))
	}

	return helpers.BipToPip(reward)
}
