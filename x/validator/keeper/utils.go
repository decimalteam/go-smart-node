package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) AddAccumRewards(ctx sdk.Context, valAddr sdk.ValAddress, r sdkmath.Int) error {
	rewards, err := k.GetValidatorRS(ctx, valAddr)
	if err != nil {
		return err
	}
	rewards.Rewards = rewards.Rewards.Add(r)
	rewards.TotalRewards = rewards.TotalRewards.Add(r)
	k.SetValidatorRS(ctx, valAddr, rewards)

	return nil
}

func (k *Keeper) ToBaseCoin(ctx sdk.Context, c sdk.Coin) sdk.Coin {
	baseDenom := k.BaseDenom(ctx)
	if baseDenom == c.Denom {
		return c
	}

	customCoin, err := k.coinKeeper.GetCoin(ctx, c.Denom)
	if err != nil {
		panic(err)
	}
	amountInBaseCoin := formulas.CalculateSaleReturn(customCoin.Volume, customCoin.Reserve, uint(customCoin.CRR), c.Amount)

	return sdk.NewCoin(baseDenom, amountInBaseCoin)
}
