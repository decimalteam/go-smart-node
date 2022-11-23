package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
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
	if c.Denom == baseDenom {
		return c
	}

	customCoinStaked := k.GetCustomCoinStaked(ctx, c.Denom)
	if customCoinStaked.Equal(sdk.ZeroInt()) {
		customCoin, err := k.coinKeeper.GetCoin(ctx, c.Denom)
		if err != nil {
			panic(err)
		}
		amountInBaseCoin := formulas.CalculateSaleReturn(customCoin.Volume, customCoin.Reserve, uint(customCoin.CRR), c.Amount)

		return sdk.NewCoin(baseDenom, amountInBaseCoin)
	}

	customCoinPrice := k.calculateCustomCoinPrice(ctx, c.Denom, customCoinStaked)

	baseAmount := sdk.NewDecFromInt(c.Amount).Mul(customCoinPrice).TruncateInt()

	return sdk.NewCoin(baseDenom, baseAmount)
}

func (k Keeper) getSumSubTokensReserve(ctx sdk.Context, id string, subToken []uint32) sdk.Coin {
	st, err := k.prepareSubTokens(ctx, id, subToken)
	if err != nil {
		panic(err)
	}

	sum := sdk.Coin{Amount: sdk.ZeroInt()}
	for _, subtoken := range st {
		sum.Denom = subtoken.Reserve.Denom
		sum.Amount = sum.Amount.Add(subtoken.Reserve.Amount)
	}

	return sum
}

// reads subtokens and fill reserve if it's nil
func (k Keeper) prepareSubTokens(ctx sdk.Context, tokenID string, subTokenIDs []uint32) ([]nfttypes.SubToken, error) {
	var result []nfttypes.SubToken
	nft, found := k.nftKeeper.GetToken(ctx, tokenID)
	if !found {
		return []nfttypes.SubToken{}, errors.NFTTokenNotFound
	}
	for _, subID := range subTokenIDs {
		subtoken, found := k.nftKeeper.GetSubToken(ctx, tokenID, subID)
		if !found {
			return []nfttypes.SubToken{}, errors.NFTSubTokenNotFound
		}
		if subtoken.Reserve == nil {
			subtoken.Reserve = &nft.Reserve
		}
		result = append(result, subtoken)
	}
	return result, nil
}

func sumSubTokens(subtokens []nfttypes.SubToken) sdk.Coin {
	sum := sdk.Coin{Amount: sdk.ZeroInt()}
	for _, sub := range subtokens {
		sum.Denom = sub.Reserve.Denom
		sum.Amount = sum.Amount.Add(sub.Reserve.Amount)
	}
	return sum
}
