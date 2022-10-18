package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
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

func (k Keeper) TotalStakeInBaseCoin(ctx sdk.Context, valAddress sdk.ValAddress) (sdkmath.Int, error) {
	delegations := k.GetValidatorDelegations(ctx, valAddress)

	return k.delegationsTotalStake(ctx, delegations)
}

func (k Keeper) delegationsTotalStake(ctx sdk.Context, delegations []types.Delegation) (sdkmath.Int, error) {
	totalStakeInBaseCoin := sdk.ZeroInt()
	for _, del := range delegations {
		delStake := del.Stake.Stake

		if del.Stake.SubTokenIDs != nil && len(del.Stake.SubTokenIDs) != 0 {
			delStake = k.getSumSubTokensReserve(ctx, del.GetStake().GetID(), del.GetStake().GetSubTokenIDs())
		}

		baseCoin := k.ToBaseCoin(ctx, delStake)
		totalStakeInBaseCoin.Add(baseCoin.Amount)
	}

	return totalStakeInBaseCoin, nil
}

func (k Keeper) getSumSubTokensReserve(ctx sdk.Context, id string, subToken []uint32) sdk.Coin {
	sum := sdk.Coin{Amount: sdk.ZeroInt()}

	if len(subToken) != 0 {
		for _, v := range subToken {
			subtoken, found := k.nftKeeper.GetSubToken(ctx, id, v)
			if !found {
				panic("not found subtoken")
			}
			sum.Denom = subtoken.Reserve.Denom
			sum.Amount.Add(subtoken.Reserve.Amount)
		}
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
