package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// Return all validators that a delegator is bonded to. If maxRetrieve is supplied, the respective amount will be returned.
func (k Keeper) GetDelegatorValidators(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint32) types.Validators {
	validators := make([]types.Validator, maxRetrieve)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegatorDelegationsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())

		validator, found := k.GetValidator(ctx, delegation.GetValidator())
		if !found {
			panic(types.ErrNoValidatorFound)
		}

		validators[i] = validator
		i++
	}

	return validators[:i] // trim
}

// return a validator that a delegator is bonded to
func (k Keeper) GetDelegatorValidator(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress) (v types.Validator, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetDelegationsKey(delegator, validator))
	defer iterator.Close()

	var found = false // any delegation found

	for ; iterator.Valid(); iterator.Next() {
		found = true
		break
	}
	if !found {
		return types.Validator{}, errors.ValidatorNotFound
	}

	v, found = k.GetValidator(ctx, validator)
	if !found {
		return types.Validator{}, errors.ValidatorNotFound
	}

	return v, nil
}

// return all delegations for a delegator
func (k Keeper) GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []types.Delegation {
	delegations := make([]types.Delegation, 0)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegatorDelegationsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	i := 0

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
		i++
	}

	return delegations
}

// return all unbonding-delegations for a delegator
func (k Keeper) GetAllUndelegations(ctx sdk.Context, delegator sdk.AccAddress) []types.Undelegation {
	undelegations := make([]types.Undelegation, 0)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetUBDsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	for i := 0; iterator.Valid(); iterator.Next() {
		undelegation := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		undelegations = append(undelegations, undelegation)
		i++
	}

	return undelegations
}

// return all redelegations for a delegator
func (k Keeper) GetAllRedelegations(
	ctx sdk.Context, delegator sdk.AccAddress, srcValAddress, dstValAddress sdk.ValAddress,
) []types.Redelegation {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetREDsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	srcValFilter := !(srcValAddress.Empty())
	dstValFilter := !(dstValAddress.Empty())

	redelegations := []types.Redelegation{}

	for ; iterator.Valid(); iterator.Next() {
		redelegation := types.MustUnmarshalRED(k.cdc, iterator.Value())
		valSrcAddr, err := sdk.ValAddressFromBech32(redelegation.ValidatorSrc)
		if err != nil {
			panic(err)
		}
		valDstAddr, err := sdk.ValAddressFromBech32(redelegation.ValidatorDst)
		if err != nil {
			panic(err)
		}
		if srcValFilter && !(srcValAddress.Equals(valSrcAddr)) {
			continue
		}

		if dstValFilter && !(dstValAddress.Equals(valDstAddr)) {
			continue
		}

		redelegations = append(redelegations, redelegation)
	}

	return redelegations
}
