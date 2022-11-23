package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

var _ types.StakingHooks = Keeper{}

// AfterValidatorCreated - call hook if registered
func (k Keeper) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.AfterValidatorCreated(ctx, valAddr)
	}
	return nil
}

// BeforeValidatorModified - call hook if registered
func (k Keeper) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.BeforeValidatorModified(ctx, valAddr)
	}
	return nil
}

// AfterValidatorRemoved - call hook if registered
func (k Keeper) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.AfterValidatorRemoved(ctx, consAddr, valAddr)
	}
	return nil
}

// AfterValidatorBonded - call hook if registered
func (k Keeper) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.AfterValidatorBonded(ctx, consAddr, valAddr)
	}
	return nil
}

// AfterValidatorBeginUnbonding - call hook if registered
func (k Keeper) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
	}
	return nil
}

// BeforeDelegationCreated - call hook if registered
func (k Keeper) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.BeforeDelegationCreated(ctx, delAddr, valAddr)
	}
	return nil
}

// BeforeDelegationSharesModified - call hook if registered
func (k Keeper) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	}
	return nil
}

// BeforeDelegationRemoved - call hook if registered
func (k Keeper) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		k.hooks.BeforeDelegationRemoved(ctx, delAddr, valAddr)
	}
	return nil
}

// AfterDelegationModified - call hook if registered
func (k Keeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.AfterDelegationModified(ctx, delAddr, valAddr)
	}
	return nil
}

// BeforeValidatorSlashed - call hook if registered
func (k Keeper) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	if k.hooks != nil {
		return k.hooks.BeforeValidatorSlashed(ctx, valAddr, fraction)
	}
	return nil
}

////////////////////////////////////////////////////////
// Validator Module Hooks //////////////////////////////
////////////////////////////////////////////////////////

// BeforeUpdateDelegation before update, subtruct all delegation  staked custom coins
func (k Keeper) BeforeUpdateDelegation(ctx sdk.Context, del types.Delegation, denom string) {
	switch del.GetStake().GetType() {
	case types.StakeType_Coin:
		if denom == k.BaseDenom(ctx) {
			return
		}

		ccs := k.GetCustomCoinStaked(ctx, denom)
		ccs = ccs.Sub(del.GetStake().GetStake().Amount)
		k.SetCustomCoinStaked(ctx, denom, ccs)
	case types.StakeType_NFT:
		reserve := del.GetStake().GetStake()
		if reserve.Denom == k.BaseDenom(ctx) {
			return
		}
		ccs := k.GetCustomCoinStaked(ctx, reserve.Denom)
		ccs = ccs.Sub(reserve.Amount)
		k.SetCustomCoinStaked(ctx, reserve.Denom, ccs)
	}
}

// AfterUpdateDelegation after update sum delegation staked custom coin
func (k Keeper) AfterUpdateDelegation(ctx sdk.Context, denom string, amount sdkmath.Int) {
	if denom == k.BaseDenom(ctx) {
		return
	}

	ccs := k.GetCustomCoinStaked(ctx, denom)
	ccs = ccs.Add(amount)
	k.SetCustomCoinStaked(ctx, denom, ccs)
	err := events.EmitTypedEvent(ctx, &types.EventUpdateCoinsStaked{
		Denom:       denom,
		TotalAmount: ccs,
	})
	if err != nil {
		panic(err)
	}

	return
}
