package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/evmos/ethermint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// // BlockValidatorUpdates calculates the ValidatorUpdates for the current block
// // Called in each EndBlock
//
//	func (k Keeper) BlockValidatorUpdates(ctx sdk.Context) []abci.ValidatorUpdate {
//		// Calculate validator set changes.
//		//
//		// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
//		// UnbondAllMatureValidatorQueue.
//		// This fixes a bug when the unbonding period is instant (is the case in
//		// some of the tests). The test expected the validator to be completely
//		// unbonded after the Endblocker (go from Bonded -> Unbonding during
//		// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
//		// UnbondAllMatureValidatorQueue).
//		validatorUpdates, err := k.ApplyAndReturnValidatorSetUpdates(ctx)
//		if err != nil {
//			panic(err)
//		}
//
//		// unbond all mature validators from the unbonding queue
//		k.UnbondAllMatureValidators(ctx)
//
//		// Remove all mature unbonding delegations from the ubd queue.
//		matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
//		for _, dvPair := range matureUnbonds {
//			addr, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
//			if err != nil {
//				panic(err)
//			}
//			delegatorAddress := sdk.MustAccAddressFromBech32(dvPair.DelegatorAddress)
//
//			balances, err := k.CompleteUnbonding(ctx, delegatorAddress, addr)
//			if err != nil {
//				continue
//			}
//
//			ctx.EventManager().EmitEvent(
//				sdk.NewEvent(
//					types.EventTypeCompleteUnbonding,
//					sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
//					sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress),
//					sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress),
//				),
//			)
//		}
//
//		// Remove all mature redelegations from the red queue.
//		matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx)
//		for _, dvvTriplet := range matureRedelegations {
//			valSrcAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorSrcAddress)
//			if err != nil {
//				panic(err)
//			}
//			valDstAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorDstAddress)
//			if err != nil {
//				panic(err)
//			}
//			delegatorAddress := sdk.MustAccAddressFromBech32(dvvTriplet.DelegatorAddress)
//
//			balances, err := k.CompleteRedelegation(
//				ctx,
//				delegatorAddress,
//				valSrcAddr,
//				valDstAddr,
//			)
//			if err != nil {
//				continue
//			}
//
//			ctx.EventManager().EmitEvent(
//				sdk.NewEvent(
//					types.EventTypeCompleteRedelegation,
//					sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
//					sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress),
//					sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress),
//					sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress),
//				),
//			)
//		}
//
//		return validatorUpdates
//	}
//
// ApplyAndReturnValidatorSetUpdates applies and return accumulated updates to the bonded validator set. Also,
// * Updates the active valset as keyed by LastValidatorPowerKey.
// * Updates the total power as keyed by LastTotalPowerKey.
// * Updates validator status' according to updated powers.
// * Updates the fee pool bonded vs not-bonded tokens.
// * Updates relevant indices.
// It gets called once after genesis, another time maybe after genesis transactions,
// then once at every EndBlock.
//
// CONTRACT: Only validators with non-zero power or zero-power that were bonded
// at the previous block height or were removed from the validator set entirely
// are returned to Tendermint.
func (k Keeper) ApplyAndReturnValidatorSetUpdates(ctx sdk.Context) (updates []abci.ValidatorUpdate, err error) {
	vals := k.GetAllValidators(ctx)
	for _, val := range vals {
		if val.IsBonded() && !val.IsJailed() && val.Online {
			up := val.ABCIValidatorUpdate(ethtypes.PowerReduction)
			up.Power = k.GetLastValidatorPower(ctx, val.GetOperator())
			updates = append(updates, up)
		}
	}
	return
}

//
//// Validator state transitions
//
//func (k Keeper) bondedToUnbonding(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
//	if !validator.IsBonded() {
//		panic(fmt.Sprintf("bad state transition bondedToUnbonding, validator: %v\n", validator))
//	}
//
//	return k.beginUnbondingValidator(ctx, validator)
//}
//
//func (k Keeper) unbondingToBonded(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
//	if !validator.IsUnbonding() {
//		panic(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
//	}
//
//	return k.bondValidator(ctx, validator)
//}
//
//func (k Keeper) unbondedToBonded(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
//	if !validator.IsUnbonded() {
//		panic(fmt.Sprintf("bad state transition unbondedToBonded, validator: %v\n", validator))
//	}
//
//	return k.bondValidator(ctx, validator)
//}
//
//// UnbondingToUnbonded switches a validator from unbonding state to unbonded state
//func (k Keeper) UnbondingToUnbonded(ctx sdk.Context, validator types.Validator) types.Validator {
//	if !validator.IsUnbonding() {
//		panic(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
//	}
//
//	return k.completeUnbondingValidator(ctx, validator)
//}
//
//// send a validator to jail
//func (k Keeper) jailValidator(ctx sdk.Context, validator types.Validator) {
//	if validator.Jailed {
//		panic(fmt.Sprintf("cannot jail already jailed validator, validator: %v\n", validator))
//	}
//
//	validator.Jailed = true
//	k.SetValidator(ctx, validator)
//	//k.DeleteValidatorByPowerIndex(ctx, validator)
//}
//
//// remove a validator from jail
//func (k Keeper) unjailValidator(ctx sdk.Context, validator types.Validator) {
//	if !validator.Jailed {
//		panic(fmt.Sprintf("cannot unjail already unjailed validator, validator: %v\n", validator))
//	}
//
//	validator.Jailed = false
//	k.SetValidator(ctx, validator)
//	//k.SetValidatorByPowerIndex(ctx, validator)
//}
//
//// perform all the store operations for when a validator status becomes bonded
//func (k Keeper) bondValidator(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
//	// delete the validator by power index, as the key will change
//	//k.DeleteValidatorByPowerIndex(ctx, validator)
//
//	validator = validator.UpdateStatus(types.BondStatus_Bonded)
//
//	// save the now bonded validator record to the two referenced stores
//	k.SetValidator(ctx, validator)
//	//k.SetValidatorByPowerIndex(ctx, validator)
//
//	// delete from queue if present
//	k.DeleteValidatorQueue(ctx, validator)
//
//	// trigger hook
//	consAddr, err := validator.GetConsAddr()
//	if err != nil {
//		return validator, err
//	}
//	k.AfterValidatorBonded(ctx, consAddr, validator.GetOperator())
//
//	return validator, err
//}
//
//// perform all the store operations for when a validator begins unbonding
//func (k Keeper) beginUnbondingValidator(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
//	params := k.GetParams(ctx)
//
//	// delete the validator by power index, as the key will change
//	//k.DeleteValidatorByPowerIndex(ctx, validator)
//
//	// sanity check
//	if validator.Status != types.BondStatus_Bonded {
//		panic(fmt.Sprintf("should not already be unbonded or unbonding, validator: %v\n", validator))
//	}
//
//	validator = validator.UpdateStatus(types.BondStatus_Unbonding)
//
//	// set the unbonding completion time and completion height appropriately
//	validator.UnbondingTime = ctx.BlockHeader().Time.Add(params.UndelegationTime)
//	validator.UnbondingHeight = ctx.BlockHeader().Height
//
//	// save the now unbonded validator record and power index
//	k.SetValidator(ctx, validator)
//	//k.SetValidatorByPowerIndex(ctx, validator)
//
//	// Adds to unbonding validator queue
//	k.InsertUnbondingValidatorQueue(ctx, validator)
//
//	// trigger hook
//	consAddr, err := validator.GetConsAddr()
//	if err != nil {
//		return validator, err
//	}
//	k.AfterValidatorBeginUnbonding(ctx, consAddr, validator.GetOperator())
//
//	return validator, nil
//}
//
//// perform all the store operations for when a validator status becomes unbonded
//func (k Keeper) completeUnbondingValidator(ctx sdk.Context, validator types.Validator) types.Validator {
//	validator = validator.UpdateStatus(types.BondStatus_Unbonded)
//	k.SetValidator(ctx, validator)
//
//	return validator
//}
//
//// map of operator bech32-addresses to serialized power
//// We use bech32 strings here, because we can't have slices as keys: map[[]byte][]byte
//type validatorsByAddr map[string][]byte
//
//// get the last validator set
//func (k Keeper) getLastValidatorsByAddr(ctx sdk.Context) (validatorsByAddr, error) {
//	last := make(validatorsByAddr)
//
//	iterator := k.LastValidatorsIterator(ctx)
//	defer iterator.Close()
//
//	for ; iterator.Valid(); iterator.Next() {
//		// extract the validator address from the key (prefix is 1-byte, addrLen is 1-byte)
//		valAddr := types.AddressFromLastValidatorPowerKey(iterator.Key())
//		valAddrStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
//		if err != nil {
//			return nil, err
//		}
//
//		powerBytes := iterator.Value()
//		last[valAddrStr] = make([]byte, len(powerBytes))
//		copy(last[valAddrStr], powerBytes)
//	}
//
//	return last, nil
//}
//
//// given a map of remaining validators to previous bonded power
//// returns the list of validators to be unbonded, sorted by operator address
//func sortNoLongerBonded(last validatorsByAddr) ([][]byte, error) {
//	// sort the map keys for determinism
//	noLongerBonded := make([][]byte, len(last))
//	index := 0
//
//	for valAddrStr := range last {
//		valAddrBytes, err := sdk.ValAddressFromBech32(valAddrStr)
//		if err != nil {
//			return nil, err
//		}
//		noLongerBonded[index] = valAddrBytes
//		index++
//	}
//	// sorted by address - order doesn't matter
//	sort.SliceStable(noLongerBonded, func(i, j int) bool {
//		// -1 means strictly less than
//		return bytes.Compare(noLongerBonded[i], noLongerBonded[j]) == -1
//	})
//
//	return noLongerBonded, nil
//}