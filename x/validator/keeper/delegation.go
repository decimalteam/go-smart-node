package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

////////////////////////////////////////////////////////////////
// Delegations
////////////////////////////////////////////////////////////////

// GetDelegation returns specific delegation by the given delegator address, validator and staked coin's denom.
func (k Keeper) GetDelegation(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress, denom string) (delegation types.Delegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDelegationKey(delegator, validator, denom)
	value := store.Get(key)
	if value == nil {
		return delegation, false
	}
	delegation = types.MustUnmarshalDelegation(k.cdc, value)
	return delegation, true
}

// GetAllDelegations returns all delegations (used during genesis dump).
func (k Keeper) GetAllDelegations(ctx sdk.Context) (delegations []types.Delegation) {
	k.IterateAllDelegations(ctx, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})
	return delegations
}

// GetAllDelegationsByValidator returns all delegations by validator stored in the application state.
func (k Keeper) GetAllDelegationsByValidator(ctx sdk.Context) (delegations map[string][]types.Delegation) {
	delegations = make(map[string][]types.Delegation)
	k.IterateAllDelegations(ctx, func(delegation types.Delegation) bool {
		valAddress := delegation.GetValidator().String()
		delegations[valAddress] = append(delegations[valAddress], delegation)
		return false
	})
	return
}

// GetValidatorDelegations returns all delegations to a specific validator. Useful for querier.
func (k Keeper) GetValidatorDelegations(ctx sdk.Context, validator sdk.ValAddress) (delegations []types.Delegation) { //nolint:interfacer
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetValidatorDelegationsKey(validator))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := types.GetDelegationKeyFromValIndexKey(iterator.Key())
		value := store.Get(key)
		delegation := types.MustUnmarshalDelegation(k.cdc, value)
		delegations = append(delegations, delegation)
	}
	return delegations
}

func (k Keeper) HasDelegations(ctx sdk.Context, validator sdk.ValAddress) bool { //nolint:interfacer
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetValidatorDelegationsKey(validator))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		return true
	}
	return false
}

func (k Keeper) HasUndelegations(ctx sdk.Context, validator sdk.ValAddress) bool { //nolint:interfacer
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetUBDsByValIndexKey(validator))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		return true
	}
	return false
}

func (k Keeper) HasRedelegations(ctx sdk.Context, validator sdk.ValAddress) bool { //nolint:interfacer
	store := ctx.KVStore(k.storeKey)
	iteratorSrc := sdk.KVStorePrefixIterator(store, types.GetREDsFromValSrcIndexKey(validator))
	defer iteratorSrc.Close()
	for ; iteratorSrc.Valid(); iteratorSrc.Next() {
		return true
	}
	iteratorDst := sdk.KVStorePrefixIterator(store, types.GetREDsToValDstIndexKey(validator))
	defer iteratorDst.Close()
	for ; iteratorDst.Valid(); iteratorDst.Next() {
		return true
	}
	return false
}

// GetDelegatorDelegations returns a given amount of all the delegations from a delegator.
func (k Keeper) GetDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (delegations []types.Delegation) {
	delegations = make([]types.Delegation, 0, maxRetrieve)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetDelegatorDelegationsKey(delegator))
	defer iterator.Close()
	for i := 0; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
		i++
	}
	return delegations
}

// GetDelegatorValidatorDelegations returns a given amount of all the delegations between the validator and the delegator.
func (k Keeper) GetDelegatorValidatorDelegations(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress, maxRetrieve uint16) (delegations []types.Delegation) {
	delegations = make([]types.Delegation, 0, maxRetrieve)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetDelegationsKey(delegator, validator))
	defer iterator.Close()
	for i := 0; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
		i++
	}
	return delegations
}

// IterateAllDelegations iterates through all of the delegations.
func (k Keeper) IterateAllDelegations(ctx sdk.Context, cb func(delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetAllDelegationsKey())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// SetDelegation sets a delegation.
func (k Keeper) SetDelegation(ctx sdk.Context, delegation types.Delegation) {
	delegator := delegation.GetDelegator()
	validator := delegation.GetValidator()
	denom := delegation.GetStake().GetID()
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(types.GetDelegationKey(delegator, validator, denom), b)
	store.Set(types.GetValidatorDelegatorDelegationKey(validator, delegator, denom), []byte{}) //index
}

// RemoveDelegation removes a delegation
func (k Keeper) RemoveDelegation(ctx sdk.Context, delegation types.Delegation) error {
	delegator := delegation.GetDelegator()
	validator := delegation.GetValidator()
	denom := delegation.GetStake().GetID()
	// TODO: Consider calling hooks outside of the store wrapper functions, it's unobvious.
	if err := k.BeforeDelegationRemoved(ctx, delegator, validator); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegationKey(delegator, validator, denom))
	store.Delete(types.GetValidatorDelegatorDelegationKey(validator, delegator, denom))
	return nil
}

////////////////////////////////////////////////////////////////
// Redelegations
////////////////////////////////////////////////////////////////

// GetRedelegations returns a given amount of all the delegator redelegations.
func (k Keeper) GetRedelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (redelegations []types.Redelegation) {
	redelegations = make([]types.Redelegation, 0, maxRetrieve)
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetREDsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close()
	for i := 0; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		redelegation := types.MustUnmarshalRED(k.cdc, iterator.Value())
		redelegations = append(redelegations, redelegation)
		i++
	}
	return redelegations
}

// GetRedelegation returns a redelegation.
func (k Keeper) GetRedelegation(ctx sdk.Context, delegator sdk.AccAddress, validatorSrc, validatorDst sdk.ValAddress) (red types.Redelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetREDKey(delegator, validatorSrc, validatorDst)
	value := store.Get(key)
	if value == nil {
		return red, false
	}
	red = types.MustUnmarshalRED(k.cdc, value)
	return red, true
}

// GetRedelegationsFromSrcValidator returns all redelegations from a particular validator.
func (k Keeper) GetRedelegationsFromSrcValidator(ctx sdk.Context, validatorSrc sdk.ValAddress) (reds []types.Redelegation) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetREDsFromValSrcIndexKey(validatorSrc))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := types.GetREDKeyFromValSrcIndexKey(iterator.Key())
		value := store.Get(key)
		red := types.MustUnmarshalRED(k.cdc, value)
		reds = append(reds, red)
	}
	return reds
}

// HasReceivingRedelegation checks if validator is receiving a redelegation.
func (k Keeper) HasReceivingRedelegation(ctx sdk.Context, delegator sdk.AccAddress, validatorDst sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	prefix := types.GetREDsByDelToValDstIndexKey(delegator, validatorDst)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	return iterator.Valid()
}

// HasMaxRedelegationEntries checks if redelegation has maximum number of entries.
func (k Keeper) HasMaxRedelegationEntries(ctx sdk.Context, delegator sdk.AccAddress, validatorSrc, validatorDst sdk.ValAddress) bool {
	red, found := k.GetRedelegation(ctx, delegator, validatorSrc, validatorDst)
	if !found {
		return false
	}
	return len(red.Entries) >= int(k.MaxEntries(ctx))
}

// SetRedelegation set a redelegation and associated index.
func (k Keeper) SetRedelegation(ctx sdk.Context, red types.Redelegation) {
	delegator := red.GetDelegator()
	validatorSrc := red.GetValidatorSrc()
	validatorDst := red.GetValidatorDst()
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalRED(k.cdc, red)
	key := types.GetREDKey(delegator, validatorSrc, validatorDst)
	store.Set(key, bz)
	store.Set(types.GetREDByValSrcIndexKey(delegator, validatorSrc, validatorDst), []byte{})
	store.Set(types.GetREDByValDstIndexKey(delegator, validatorSrc, validatorDst), []byte{})
}

// SetRedelegationEntry adds an entry to the unbonding delegation at the given addresses.
// It creates the undelegation if it does not exist.
func (k Keeper) SetRedelegationEntry(
	ctx sdk.Context,
	delegator sdk.AccAddress,
	validatorSrc sdk.ValAddress,
	validatorDst sdk.ValAddress,
	creationHeight int64,
	minTime time.Time,
	stake types.Stake,
) types.Redelegation {
	red, found := k.GetRedelegation(ctx, delegator, validatorSrc, validatorDst)
	if found {
		red.AddEntry(creationHeight, minTime, stake)
	} else {
		red = types.NewRedelegation(delegator, validatorSrc,
			validatorDst, creationHeight, minTime, stake)
	}

	k.SetRedelegation(ctx, red)

	return red
}

// IterateRedelegations iterates through all redelegations.
func (k Keeper) IterateRedelegations(ctx sdk.Context, fn func(index int64, red types.Redelegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetAllREDsKey())
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		red := types.MustUnmarshalRED(k.cdc, iterator.Value())
		if stop := fn(i, red); stop {
			break
		}
		i++
	}
}

// RemoveRedelegation removes a redelegation object and associated index.
func (k Keeper) RemoveRedelegation(ctx sdk.Context, red types.Redelegation) {
	delegator := red.GetDelegator()
	validatorSrc := red.GetValidatorSrc()
	validatorDst := red.GetValidatorDst()
	store := ctx.KVStore(k.storeKey)
	redKey := types.GetREDKey(delegator, validatorSrc, validatorDst)
	store.Delete(redKey)
	store.Delete(types.GetREDByValSrcIndexKey(delegator, validatorSrc, validatorDst))
	store.Delete(types.GetREDByValDstIndexKey(delegator, validatorSrc, validatorDst))
}

// GetRedelegationQueueTimeSlice gets a specific redelegation queue timeslice. A
// timeslice is a slice of DVVTriplets corresponding to redelegations that
// expire at a certain time.
func (k Keeper) GetRedelegationQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvvTriplets []stakingtypes.DVVTriplet) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetRedelegationsTimeKey(timestamp))
	if bz == nil {
		return []stakingtypes.DVVTriplet{}
	}

	triplets := stakingtypes.DVVTriplets{}
	k.cdc.MustUnmarshal(bz, &triplets)

	return triplets.Triplets
}

// SetRedelegationQueueTimeSlice sets a specific redelegation queue timeslice.
func (k Keeper) SetRedelegationQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []stakingtypes.DVVTriplet) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&stakingtypes.DVVTriplets{Triplets: keys})
	store.Set(types.GetRedelegationsTimeKey(timestamp), bz)
}

// InsertRedelegationQueue insert an redelegation delegation to the appropriate
// timeslice in the redelegation queue.
func (k Keeper) InsertRedelegationQueue(ctx sdk.Context, red types.Redelegation, completionTime time.Time) {
	timeSlice := k.GetRedelegationQueueTimeSlice(ctx, completionTime)
	dvvTriplet := stakingtypes.DVVTriplet{
		DelegatorAddress:    red.Delegator,
		ValidatorSrcAddress: red.ValidatorSrc,
		ValidatorDstAddress: red.ValidatorDst,
	}

	if len(timeSlice) == 0 {
		k.SetRedelegationQueueTimeSlice(ctx, completionTime, []stakingtypes.DVVTriplet{dvvTriplet})
	} else {
		timeSlice = append(timeSlice, dvvTriplet)
		k.SetRedelegationQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// RedelegationQueueIterator returns all the redelegation queue timeslices from
// time 0 until endTime.
func (k Keeper) RedelegationQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.GetAllRedelegationsTimeKey(), sdk.InclusiveEndBytes(types.GetRedelegationsTimeKey(endTime)))
}

// DequeueAllMatureRedelegationQueue returns a concatenated list of all the
// timeslices inclusively previous to currTime, and deletes the timeslices from
// the queue.
func (k Keeper) DequeueAllMatureRedelegationQueue(ctx sdk.Context) (matureRedelegations []stakingtypes.DVVTriplet) {
	store := ctx.KVStore(k.storeKey)

	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	redelegationTimesliceIterator := k.RedelegationQueueIterator(ctx, ctx.BlockHeader().Time)
	defer redelegationTimesliceIterator.Close()

	for ; redelegationTimesliceIterator.Valid(); redelegationTimesliceIterator.Next() {
		timeslice := stakingtypes.DVVTriplets{}
		value := redelegationTimesliceIterator.Value()
		k.cdc.MustUnmarshal(value, &timeslice)

		matureRedelegations = append(matureRedelegations, timeslice.Triplets...)

		store.Delete(redelegationTimesliceIterator.Key())
	}

	return matureRedelegations
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

// GetUndelegations returns a given amount of all the delegator unbonding-delegations.
func (k Keeper) GetUndelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (undelegations []types.Undelegation) {
	undelegations = make([]types.Undelegation, maxRetrieve)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetUBDsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		undelegation := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		undelegations[i] = undelegation
		i++
	}

	return undelegations[:i] // trim if the array length < maxRetrieve
}

// GetUndelegation returns a unbonding delegation.
func (k Keeper) GetUndelegation(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress) (ubd types.Undelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDKey(delegator, validator)
	value := store.Get(key)

	if value == nil {
		return ubd, false
	}

	ubd = types.MustUnmarshalUBD(k.cdc, value)

	return ubd, true
}

// GetUndelegationsFromValidator returns all unbonding delegations from a particular validator.
func (k Keeper) GetUndelegationsFromValidator(ctx sdk.Context, validator sdk.ValAddress) (ubds []types.Undelegation) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetUBDsByValIndexKey(validator))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := types.GetUBDKeyFromValIndexKey(iterator.Key())
		value := store.Get(key)
		ubd := types.MustUnmarshalUBD(k.cdc, value)
		ubds = append(ubds, ubd)
	}

	return ubds
}

// IterateUndelegations iterates through all of the unbonding delegations.
func (k Keeper) IterateUndelegations(ctx sdk.Context, fn func(index int64, ubd types.Undelegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetAllUBDsKey())
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		ubd := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		if stop := fn(i, ubd); stop {
			break
		}
		i++
	}
}

// GetDelegatorUnbonding returns the total amount a delegator has unbonding.
func (k Keeper) GetDelegatorUnbonding(ctx sdk.Context, delegator sdk.AccAddress) sdkmath.Int {
	unbonding := sdk.ZeroInt()

	k.IterateDelegatorUndelegations(ctx, delegator, func(ubd types.Undelegation) bool {
		for _, entry := range ubd.Entries {
			unbonding = unbonding.Add(entry.Stake.Stake.Amount)
		}
		return false
	})
	return unbonding
}

// IterateDelegatorUndelegations iterates through a delegator's unbonding delegations.
func (k Keeper) IterateDelegatorUndelegations(ctx sdk.Context, delegator sdk.AccAddress, cb func(ubd types.Undelegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetUBDsKey(delegator))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		ubd := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		if cb(ubd) {
			break
		}
	}
}

// GetDelegatorBonded returs the total amount a delegator has bonded.
func (k Keeper) GetDelegatorBonded(ctx sdk.Context, delegator sdk.AccAddress) sdkmath.Int {
	bonded := sdk.ZeroInt()

	k.IterateDelegatorDelegations(ctx, delegator, func(delegation types.Delegation) bool {
		amount := delegation.Stake.GetStake().Amount
		bonded = bonded.Add(amount)

		return false
	})
	return bonded
}

// IterateDelegatorDelegations iterates through one delegator's delegations.
func (k Keeper) IterateDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress, cb func(delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegatorDelegationsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// IterateDelegatorRedelegations iterates through one delegator's redelegations.
func (k Keeper) IterateDelegatorRedelegations(ctx sdk.Context, delegator sdk.AccAddress, cb func(red types.Redelegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetREDsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		red := types.MustUnmarshalRED(k.cdc, iterator.Value())
		if cb(red) {
			break
		}
	}
}

// HasMaxUndelegationEntries - check if unbonding delegation has maximum number of entries.
func (k Keeper) HasMaxUndelegationEntries(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress) bool {
	ubd, found := k.GetUndelegation(ctx, delegator, validator)
	if !found {
		return false
	}

	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
}

// SetUndelegation sets the unbonding delegation and associated index.
func (k Keeper) SetUndelegation(ctx sdk.Context, ubd types.Undelegation) {
	delegator := sdk.MustAccAddressFromBech32(ubd.Delegator)

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalUBD(k.cdc, ubd)
	addr, err := sdk.ValAddressFromBech32(ubd.Validator)
	if err != nil {
		panic(err)
	}
	key := types.GetUBDKey(delegator, addr)
	store.Set(key, bz)
	store.Set(types.GetUBDByValIndexKey(delegator, addr), []byte{}) // index, store empty bytes
}

// RemoveUndelegation removes the unbonding delegation object and associated index.
func (k Keeper) RemoveUndelegation(ctx sdk.Context, ubd types.Undelegation) {
	delegator := sdk.MustAccAddressFromBech32(ubd.Delegator)

	store := ctx.KVStore(k.storeKey)
	addr, err := sdk.ValAddressFromBech32(ubd.Validator)
	if err != nil {
		panic(err)
	}
	key := types.GetUBDKey(delegator, addr)
	store.Delete(key)
	store.Delete(types.GetUBDByValIndexKey(delegator, addr))
}

// SetUndelegationEntry adds an entry to the unbonding delegation at
// the given addresses. It creates the unbonding delegation if it does not exist.
func (k Keeper) SetUndelegationEntry(
	ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress,
	creationHeight int64, minTime time.Time, stake types.Stake,
) types.Undelegation {
	ubd, found := k.GetUndelegation(ctx, delegator, validator)
	if found {
		ubd.AddEntry(creationHeight, minTime, stake)
	} else {
		ubd = types.NewUndelegation(delegator, validator, creationHeight, minTime, stake)
	}

	k.SetUndelegation(ctx, ubd)

	return ubd
}

// unbonding delegation queue timeslice operations

// GetUBDQueueTimeSlice gets a specific unbonding queue timeslice. A timeslice
// is a slice of DVPairs corresponding to unbonding delegations that expire at a
// certain time.
func (k Keeper) GetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvPairs []stakingtypes.DVPair) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetUndelegationsTimeKey(timestamp))
	if bz == nil {
		return []stakingtypes.DVPair{}
	}

	pairs := stakingtypes.DVPairs{}
	k.cdc.MustUnmarshal(bz, &pairs)

	return pairs.Pairs
}

// SetUBDQueueTimeSlice sets a specific unbonding queue timeslice.
func (k Keeper) SetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []stakingtypes.DVPair) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&stakingtypes.DVPairs{Pairs: keys})
	store.Set(types.GetUndelegationsTimeKey(timestamp), bz)
}

// InsertUBDQueue inserts an unbonding delegation to the appropriate timeslice
// in the unbonding queue.
func (k Keeper) InsertUBDQueue(ctx sdk.Context, ubd types.Undelegation, completionTime time.Time) {
	dvPair := stakingtypes.DVPair{DelegatorAddress: ubd.Delegator, ValidatorAddress: ubd.Validator}

	timeSlice := k.GetUBDQueueTimeSlice(ctx, completionTime)
	if len(timeSlice) == 0 {
		k.SetUBDQueueTimeSlice(ctx, completionTime, []stakingtypes.DVPair{dvPair})
	} else {
		timeSlice = append(timeSlice, dvPair)
		k.SetUBDQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// UBDQueueIterator returns all the unbonding queue timeslices from time 0 until endTime.
func (k Keeper) UBDQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.GetAllUndelegationsTimeKey(),
		sdk.InclusiveEndBytes(types.GetUndelegationsTimeKey(endTime)))
}

// DequeueAllMatureUBDQueue returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue.
func (k Keeper) DequeueAllMatureUBDQueue(ctx sdk.Context, currTime time.Time) (matureUnbonds []stakingtypes.DVPair) {
	store := ctx.KVStore(k.storeKey)

	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator := k.UBDQueueIterator(ctx, currTime)
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := stakingtypes.DVPairs{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshal(value, &timeslice)

		matureUnbonds = append(matureUnbonds, timeslice.Pairs...)

		store.Delete(unbondingTimesliceIterator.Key())
	}

	return matureUnbonds
}

////////////////////////////////////////////////////////////////
// CustomCoinStaked
////////////////////////////////////////////////////////////////

func (k Keeper) SetCustomCoinStaked(ctx sdk.Context, denom string, amount sdkmath.Int) {
	store := ctx.KVStore(k.storeKey)

	if amount.IsZero() {
		store.Delete(types.GetCustomCoinStaked(denom))
		return
	}

	bz, err := amount.Marshal()
	if err != nil {
		panic(err)
	}

	store.Set(types.GetCustomCoinStaked(denom), bz)
}

func (k Keeper) GetCustomCoinStaked(ctx sdk.Context, denom string) sdkmath.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetCustomCoinStaked(denom))
	amount := sdk.ZeroInt()
	if bz == nil {
		return amount
	}
	err := amount.Unmarshal(bz)
	if err != nil {
		panic(err)
	}

	return amount
}

func (k Keeper) GetAllCustomCoinsStaked(ctx sdk.Context) map[string]sdkmath.Int {
	result := make(map[string]sdkmath.Int)

	k.IterateAllCustomCoinStaked(ctx, func(denom string, amount sdkmath.Int) bool {
		result[denom] = amount

		return false
	})

	return result
}

func (k Keeper) IterateAllCustomCoinStaked(ctx sdk.Context, cb func(denom string, amount sdkmath.Int) bool) {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetAllCustomCoinsStaked()

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		amount := sdk.ZeroInt()
		err := amount.Unmarshal(iterator.Value())
		if err != nil {
			panic(err)
		}

		denom := string(iterator.Key()[1:])

		if cb(denom, amount) {
			break
		}
	}
}

func (k Keeper) AddCustomCoinStaked(ctx sdk.Context, coin sdk.Coin) {
	if coin.Denom == k.BaseDenom(ctx) {
		return
	}
	amount := k.GetCustomCoinStaked(ctx, coin.Denom)
	amount = amount.Add(coin.Amount)
	// emit event
	k.SetCustomCoinStaked(ctx, coin.Denom, amount)
	err := events.EmitTypedEvent(ctx, &types.EventUpdateCoinsStaked{
		Denom:       coin.Denom,
		TotalAmount: amount,
	})
	if err != nil {
		panic(err)
	}
}

func (k Keeper) SubCustomCoinStaked(ctx sdk.Context, coin sdk.Coin) {
	if coin.Denom == k.BaseDenom(ctx) {
		return
	}
	amount := k.GetCustomCoinStaked(ctx, coin.Denom)
	amount = amount.Sub(coin.Amount)
	if amount.IsNegative() {
		panic(fmt.Errorf("amount of staked custom coin '%s' become negative", coin.Denom))
	}
	k.SetCustomCoinStaked(ctx, coin.Denom, amount)
	// emit event
	err := events.EmitTypedEvent(ctx, &types.EventUpdateCoinsStaked{
		Denom:       coin.Denom,
		TotalAmount: amount,
	})
	if err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

// Delegate performs a delegation, set/update everything necessary within the store.
// tokenSrc indicates the bond status of the incoming funds.
// NFT subtoken ownership MUST BE checked before
func (k Keeper) Delegate(
	ctx sdk.Context, delegator sdk.AccAddress, validator types.Validator, stake types.Stake,
) error {
	var err error
	// 1. Get or create the delegation object
	delegation, delegationFound := k.GetDelegation(ctx, delegator, validator.GetOperator(), stake.ID)
	if !delegationFound {
		delegation = types.NewDelegation(delegator, validator.GetOperator(), stake)
		k.IncrementDelegationsCount(ctx, validator.GetOperator())
	} else {
		delegation.Stake, err = delegation.Stake.Add(stake)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	// 3. transfer coin/nft
	notBondedPool := k.GetNotBondedPool(ctx).GetAddress()
	bondedPool := k.GetBondedPool(ctx).GetAddress()
	var receiverName string
	var receiverPool sdk.AccAddress
	switch {
	case validator.IsBonded():
		receiverName = types.BondedPoolName
		receiverPool = bondedPool
	case validator.IsUnbonding(), validator.IsUnbonded():
		receiverName = types.NotBondedPoolName
		receiverPool = notBondedPool
	default:
		return errors.ValidatorStatusUnknown
	}

	// the transfer of user assets is carried out in coins or nft
	switch stake.Type {
	case types.StakeType_Coin:
		if err := k.bankKeeper.DelegateCoinsFromAccountToModule(ctx, delegator, receiverName, sdk.NewCoins(stake.Stake)); err != nil {
			return err
		}
	case types.StakeType_NFT:
		if err := k.nftKeeper.TransferSubTokens(ctx, delegator, receiverPool, stake.ID, stake.SubTokenIDs); err != nil {
			return err
		}
	}

	// Update delegation
	k.SetDelegation(ctx, delegation)

	// update validator info
	valAddress := validator.GetOperator()

	// clean index
	k.DeleteValidatorByPowerIndex(ctx, validator)

	rs, err := k.GetValidatorRS(ctx, valAddress)
	if err != nil {
		if err == errors.RewardsNotFound {
			rs = types.ValidatorRS{
				Rewards:      sdkmath.ZeroInt(),
				TotalRewards: sdkmath.ZeroInt(),
				Stake:        0,
			}
		} else {
			return err
		}
	}

	// calculate validator new stake
	// change stake/validator power only if validator is online
	if validator.Online {
		delStake := k.ToBaseCoin(ctx, stake.GetStake())
		power := TokensToConsensusPower(delStake.Amount)
		rs.Stake += power
		validator.Stake += power
	} else {
		if rs.Stake > 0 {
			rs.Stake = 0
			validator.Stake = 0
		}
	}
	// write index
	k.SetValidatorRS(ctx, valAddress, rs)
	k.SetValidatorByPowerIndex(ctx, validator)

	k.AddCustomCoinStaked(ctx, stake.GetStake())
	return nil
}

func (k Keeper) TransferStakeBetweenPools(ctx sdk.Context, statusSrc types.BondStatus, statusDst types.BondStatus, stake types.Stake) error {
	switch stake.Type {
	case types.StakeType_Coin:
		return k.transferBetweenPools(ctx, statusSrc, statusDst, sdk.NewCoins(stake.Stake), nil)
	case types.StakeType_NFT:
		return k.transferBetweenPools(ctx, statusSrc, statusDst, nil, []nftTransferRecord{
			{
				tokenID:     stake.ID,
				subTokenIDs: stake.SubTokenIDs,
			},
		})
	}
	return nil
}

// universal transfer between pools
type nftTransferRecord struct {
	tokenID     string
	subTokenIDs []uint32
}

func (k Keeper) transferBetweenPools(ctx sdk.Context, statusSrc types.BondStatus, statusDst types.BondStatus, coins sdk.Coins, nfts []nftTransferRecord) error {
	notBondedPool := k.GetNotBondedPool(ctx).GetAddress()
	bondedPool := k.GetBondedPool(ctx).GetAddress()

	switch {
	case statusSrc == types.BondStatus_Bonded && statusDst == types.BondStatus_Bonded:
		// do nothing
	case (statusSrc == types.BondStatus_Unbonded || statusSrc == types.BondStatus_Unbonding) &&
		(statusDst == types.BondStatus_Unbonded || statusDst == types.BondStatus_Unbonding):
		// do nothing
	case statusSrc == types.BondStatus_Bonded && (statusDst == types.BondStatus_Unbonded || statusDst == types.BondStatus_Unbonding):
		// transfer pools bond->not bond
		if !coins.IsZero() {
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.BondedPoolName, types.NotBondedPoolName, coins); err != nil {
				return err
			}
		}
		if len(nfts) > 0 {
			for _, rec := range nfts {
				if err := k.nftKeeper.TransferSubTokens(ctx, bondedPool, notBondedPool, rec.tokenID, rec.subTokenIDs); err != nil {
					return err
				}
			}
		}
	case (statusSrc == types.BondStatus_Unbonded || statusSrc == types.BondStatus_Unbonding) && statusDst == types.BondStatus_Bonded:
		// transfer pools not bond->bond
		if !coins.IsZero() {
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.NotBondedPoolName, types.BondedPoolName, coins); err != nil {
				return err
			}
		}
		if len(nfts) > 0 {
			for _, rec := range nfts {
				if err := k.nftKeeper.TransferSubTokens(ctx, notBondedPool, bondedPool, rec.tokenID, rec.subTokenIDs); err != nil {
					return err
				}
			}
		}
	default:
		return errors.IncompatibleBondStatuses
	}
	return nil
}

// Undelegate unbonds a stake from a given validator.
// stake and remainStake must be calculated before by CalculateUnbondStake
// It create an unbonding object and insert into the unbonding queue which will be
// processed during EndBlocker.
func (k Keeper) Undelegate(
	ctx sdk.Context, delegator sdk.AccAddress, valAddress sdk.ValAddress,
	stake types.Stake, remainStake types.Stake,
) (time.Time, error) {
	validator, found := k.GetValidator(ctx, valAddress)
	if !found {
		return time.Time{}, errors.ValidatorNotFound
	}

	if k.HasMaxUndelegationEntries(ctx, delegator, valAddress) {
		return time.Time{}, errors.MaxUndelegationEntries
	}

	err := k.Unbond(ctx, delegator, valAddress, stake, remainStake)
	if err != nil {
		return time.Time{}, err
	}

	// transfer in pool
	err = k.TransferStakeBetweenPools(ctx, validator.GetStatus(), types.BondStatus_Unbonding, stake)
	if err != nil {
		return time.Time{}, err
	}

	completionTime := ctx.BlockHeader().Time.Add(k.UndelegationTime(ctx))
	ubd := k.SetUndelegationEntry(ctx, delegator, valAddress, ctx.BlockHeight(), completionTime, stake)
	k.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, nil
}

// CompleteUnbonding completes the unbonding of all mature entries in the
// retrieved unbonding delegation object and returns the total unbonding balance
// or an error upon failure.
func (k Keeper) CompleteUnbonding(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress) (err error) {
	ubd, found := k.GetUndelegation(ctx, delegator, validator)
	if !found {
		return types.ErrNoUndelegation
	}

	ctxTime := ctx.BlockHeader().Time

	delegator, err = sdk.AccAddressFromBech32(ubd.Delegator)
	if err != nil {
		return err
	}

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) {
			ubd.RemoveEntry(int64(i))
			i--
			stake := entry.Stake
			switch stake.Type {
			case types.StakeType_Coin:
				amt := stake.Stake

				if err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(
					ctx, types.NotBondedPoolName, delegator, sdk.NewCoins(amt),
				); err != nil {
					return err
				}
			case types.StakeType_NFT:
				if err := k.nftKeeper.TransferSubTokens(ctx, k.GetNotBondedPool(ctx).GetAddress(), delegator, stake.GetID(), stake.GetSubTokenIDs()); err != nil {
					return err
				}
			}
			k.SubCustomCoinStaked(ctx, entry.Stake.Stake)

			err = events.EmitTypedEvent(ctx, &types.EventUndelegateComplete{
				Delegator: delegator.String(),
				Validator: validator.String(),
				Stake:     stake,
			})
			if err != nil {
				return err
			}
		}
	}

	// set the unbonding delegation or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUndelegation(ctx, ubd)
	} else {
		k.SetUndelegation(ctx, ubd)
	}

	return nil
}

// BeginRedelegation begins unbonding / redelegation and creates a redelegation
// record.
// stake and remainStake MUST BE calculated before by ValidateUnbondStake
func (k Keeper) BeginRedelegation(
	ctx sdk.Context, delegator sdk.AccAddress, validatorSrc, validatorDst sdk.ValAddress, stake types.Stake, remainStake types.Stake,
) (completionTime time.Time, err error) {
	// 1. preparations, checks
	if bytes.Equal(validatorSrc, validatorDst) {
		return time.Time{}, errors.SelfRedelegation
	}

	_, found := k.GetValidator(ctx, validatorDst)
	if !found {
		return time.Time{}, errors.BadRedelegationDst
	}

	srcValidator, found := k.GetValidator(ctx, validatorSrc)
	if !found {
		return time.Time{}, errors.BadRedelegationSrc
	}

	if k.HasMaxRedelegationEntries(ctx, delegator, validatorSrc, validatorDst) {
		return time.Time{}, errors.MaxRedelegationEntries
	}

	// 2. modify or remove current delegation
	// call the before-delegation-modified hook
	err = k.Unbond(ctx, delegator, validatorSrc, stake, remainStake)
	if err != nil {
		return time.Time{}, err
	}

	// 3. transfer in pool
	err = k.TransferStakeBetweenPools(ctx, srcValidator.GetStatus(), types.BondStatus_Unbonded, stake)
	if err != nil {
		return time.Time{}, err
	}

	// create the unbonding delegation
	completionTime, height := k.getBeginInfo(ctx, validatorSrc)

	red := k.SetRedelegationEntry(
		ctx, delegator, validatorSrc, validatorDst,
		height, completionTime, stake,
	)
	k.InsertRedelegationQueue(ctx, red, completionTime)

	return completionTime, nil
}

// CompleteRedelegation completes the redelegations of all mature entries in the
// retrieved redelegation object and returns the total redelegation (initial)
// balance or an error upon failure.
func (k Keeper) CompleteRedelegation(
	ctx sdk.Context, delegator sdk.AccAddress, validatorSrc, validatorDst sdk.ValAddress,
) (err error) {
	red, found := k.GetRedelegation(ctx, delegator, validatorSrc, validatorDst)
	if !found {
		return types.ErrNoRedelegation
	}

	ctxTime := ctx.BlockHeader().Time

	// loop through all the entries and complete mature redelegation entries
	for i := 0; i < len(red.Entries); i++ {
		entry := red.Entries[i]
		if entry.IsMature(ctxTime) {
			red.RemoveEntry(int64(i))
			i--

			stake := entry.Stake
			// return coins
			switch stake.Type {
			case types.StakeType_Coin:
				amt := stake.Stake

				if err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(
					ctx, types.NotBondedPoolName, delegator, sdk.NewCoins(amt),
				); err != nil {
					return err
				}

			case types.StakeType_NFT:
				if err := k.nftKeeper.TransferSubTokens(ctx, k.GetNotBondedPool(ctx).GetAddress(), delegator, stake.GetID(), stake.GetSubTokenIDs()); err != nil {
					return err
				}
			}
			k.SubCustomCoinStaked(ctx, stake.Stake)

			// delegate
			validator, found := k.GetValidator(ctx, validatorDst)
			if !found {
				return fmt.Errorf("not found validator %s", validatorDst)
			}
			err := k.Delegate(ctx, delegator, validator, entry.Stake)
			if err != nil {
				return err
			}

			err = events.EmitTypedEvent(ctx, &types.EventRedelegateComplete{
				Delegator:    delegator.String(),
				ValidatorSrc: validatorSrc.String(),
				ValidatorDst: validatorDst.String(),
				Stake:        stake,
			})
			if err != nil {
				return err
			}
		}
	}

	// set the redelegation or remove it if there are no more entries
	if len(red.Entries) == 0 {
		k.RemoveRedelegation(ctx, red)
	} else {
		k.SetRedelegation(ctx, red)
	}

	return nil
}

func (k Keeper) Unbond(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, stake, remainStake types.Stake) (err error) {
	// check if a delegation object exists in the store
	delegation, found := k.GetDelegation(ctx, delAddr, valAddr, stake.ID)
	if !found {
		return types.ErrNoDelegatorForAddress
	}

	// get validator
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return types.ErrNoValidatorFound
	}

	// set delegation new stake
	delegation.Stake = remainStake
	if k.ToBaseCoin(ctx, remainStake.GetStake()).IsZero() {
		err = k.RemoveDelegation(ctx, delegation)
		k.DecrementDelegationsCount(ctx, validator.GetOperator())
	} else {
		k.SetDelegation(ctx, delegation)
	}

	stakeBaseCoin := k.ToBaseCoin(ctx, stake.GetStake())
	stakePower := TokensToConsensusPower(stakeBaseCoin.Amount)

	// clean index
	k.DeleteValidatorByPowerIndex(ctx, validator)
	rs, err := k.GetValidatorRS(ctx, valAddr)
	if err != nil {
		panic(err)
	}

	// calculate validator new stake
	if validator.Online {
		validator.Stake -= stakePower
		if validator.Stake < 0 {
			validator.Stake = 0
		}
		rs.Stake -= stakePower
		if rs.Stake < 0 {
			rs.Stake = 0
		}
	} else {
		if rs.Stake > 0 {
			rs.Stake = 0
			validator.Stake = 0
		}
	}
	// write index
	k.SetValidatorRS(ctx, valAddr, rs)
	k.SetValidatorByPowerIndex(ctx, validator)

	return
}

// CalculateRemainStake validates that a given stake can be
// substracted from source. If the stake is valid, the remain stake is returned,
// otherwise an error is returned.
func (k Keeper) CalculateRemainStake(
	ctx sdk.Context, source, stake types.Stake,
) (types.Stake, error) {
	if source.Type != stake.Type {
		return types.Stake{}, errors.WrongStakeType
	}

	var remainStake types.Stake

	switch stake.Type {
	case types.StakeType_Coin:
		if !source.Stake.IsGTE(stake.Stake) {
			return types.Stake{}, errors.StakeTooSmall
		}
		remainStake = types.NewStakeCoin(source.Stake.Sub(stake.Stake))
	case types.StakeType_NFT:
		source.Stake = k.getSumSubTokensReserve(ctx, source.ID, source.GetSubTokenIDs())
		if !types.SetHasSubset(source.SubTokenIDs, stake.SubTokenIDs) {
			return types.Stake{}, errors.StakeDoesNotHaveSubTokenID
		}
		remainIDs := types.SetSubstract(source.SubTokenIDs, stake.SubTokenIDs)
		subtokens, err := k.prepareSubTokens(ctx, stake.ID, stake.SubTokenIDs)
		if err != nil {
			return types.Stake{}, errors.NFTSubTokenNotFound
		}
		amount := source.Stake
		for _, sub := range subtokens {
			amount = amount.Sub(*sub.Reserve)
		}
		remainStake = types.NewStakeNFT(stake.ID, remainIDs, amount)
	}

	return remainStake, nil
}

// getBeginInfo returns the completion time and height of a redelegation, along with
// a boolean signaling if the redelegation is complete based on the source validator.
func (k Keeper) getBeginInfo(ctx sdk.Context, validatorSrc sdk.ValAddress) (completionTime time.Time, height int64) {
	validator, found := k.GetValidator(ctx, validatorSrc)

	// TODO: When would the validator not be found?
	switch {
	case !found || validator.IsBonded() || validator.IsUnbonded():
		// the longest wait - just unbonding period from now
		completionTime = ctx.BlockHeader().Time.Add(k.RedelegationTime(ctx))
		height = ctx.BlockHeight()
		return completionTime, height

	default:
		panic(fmt.Sprintf("unknown validator status for redelegation: %s", validator.Status))
	}
}

func (k Keeper) IncrementDelegationsCount(ctx sdk.Context, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorDelegationsCount(valAddr)
	counter := uint32(0)
	value := store.Get(key)
	if len(value) > 0 {
		counter = binary.BigEndian.Uint32(value)
	}
	counter++
	bz := make([]byte, 4)
	binary.BigEndian.PutUint32(bz, counter)
	store.Set(key, bz)
}

func (k Keeper) DecrementDelegationsCount(ctx sdk.Context, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorDelegationsCount(valAddr)
	counter := uint32(0)
	value := store.Get(key)
	if len(value) > 0 {
		counter = binary.BigEndian.Uint32(value)
	}
	counter--
	bz := make([]byte, 4)
	binary.BigEndian.PutUint32(bz, counter)
	store.Set(key, bz)
}

func (k Keeper) GetDelegationsCount(ctx sdk.Context, valAddr sdk.ValAddress) uint32 {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorDelegationsCount(valAddr)
	counter := uint32(0)
	value := store.Get(key)
	if len(value) > 0 {
		counter = binary.BigEndian.Uint32(value)
	}
	return counter
}

func (k Keeper) GetAllDelegationsCount(ctx sdk.Context) map[string]uint32 {
	result := make(map[string]uint32)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetAllDelegationsCount())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		counter := uint32(0)
		if len(iterator.Value()) > 0 {
			counter = binary.BigEndian.Uint32(iterator.Value())
		}
		valAddr := types.ParseValidatorDelegationsCountKey(iterator.Key())
		result[valAddr.String()] = counter
	}
	return result
}
