package keeper

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

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
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
	}
	return delegations
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
	store.Set(types.GetValidatorDelegatorDelegationKey(validator, delegator, denom), b)
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
	return store.Iterator(types.GetRedelegationsTimeKey(endTime), sdk.InclusiveEndBytes(types.GetRedelegationsTimeKey(endTime)))
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

	iterator := sdk.KVStorePrefixIterator(store, types.GetAllDelegationsKey())
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
	return store.Iterator(types.GetUndelegationsTimeKey(endTime),
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

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

// Delegate performs a delegation, set/update everything necessary within the store.
// tokenSrc indicates the bond status of the incoming funds.
func (k Keeper) Delegate(
	ctx sdk.Context, delegator sdk.AccAddress, denom string, coinAmount *sdkmath.Int, subTokenIDs []uint32, tokenSrc types.BondStatus,
	validator types.Validator, subtractAccount bool,
) (totalStake sdkmath.Int, err error) {
	// create stake entity
	var stake types.Stake
	switch {
	case coinAmount != nil && subTokenIDs == nil: // if stake is coin
		stake = types.NewStakeCoin(sdk.NewCoin(denom, *coinAmount))
	case subTokenIDs != nil && coinAmount == nil: // if stake is nft
		if len(subTokenIDs) == 0 {
			return sdk.ZeroInt(), fmt.Errorf("not coin delegate and not nft delegate, who is you") //TODO Error
		}

		var reserve *sdk.Coin
		for _, v := range subTokenIDs {
			st, ok := k.nftKeeper.GetSubToken(ctx, denom, v)
			if !ok {
				return sdk.ZeroInt(), fmt.Errorf("not found sub token") // TODO error
			}
			if reserve == nil {
				reserve = st.Reserve
			}
			reserve.Add(*st.Reserve)
		}
		stake = types.NewStakeNFT(denom, subTokenIDs, *reserve)
	default:
		return sdkmath.Int{}, fmt.Errorf("not coin delegate and not nft delegate, who is you") //TODO error
	}

	// Get or create the delegation object
	delegation, found := k.GetDelegation(ctx, delegator, validator.GetOperator(), denom)
	if !found {
		delegation = types.NewDelegation(delegator, validator.GetOperator(), stake)
	}

	// call the appropriate hook if present
	if found {
		k.BeforeUpdateDelegation(ctx, delegation, denom)
		//err = k.BeforeDelegationSharesModified(ctx, delegator, validator.GetOperator())
	} else {
		// nothing now
		//err = k.BeforeDelegationCreated(ctx, delegator, validator.GetOperator())
	}

	if err != nil {
		return sdk.ZeroInt(), err
	}

	// if subtractAccount is true then we are
	// performing a delegation and not a redelegation, thus the source tokens are
	// all non bonded
	notBondedPool := k.GetNotBondedPool(ctx).GetAddress()
	bondedPool := k.GetBondedPool(ctx).GetAddress()
	if subtractAccount {
		if tokenSrc == types.BondStatus_Bonded {
			panic("delegation token source cannot be bonded")
		}

		var sendName string
		var sendPool sdk.AccAddress
		switch {
		case validator.IsBonded():
			sendName = types.BondedPoolName
			sendPool = bondedPool
		case validator.IsUnbonding(), validator.IsUnbonded():
			sendName = types.NotBondedPoolName
			sendPool = notBondedPool
		default:
			panic("invalid validator status")
		}

		// the transfer of user assets is carried out in coins or nft
		switch {
		case coinAmount != nil && subTokenIDs == nil:
			coins := sdk.NewCoins(stake.Stake)
			if err := k.bankKeeper.DelegateCoinsFromAccountToModule(ctx, delegator, sendName, coins); err != nil {
				return sdk.ZeroInt(), err
			}
		case subTokenIDs != nil && coinAmount == nil:
			if err := k.nftKeeper.TransferSubTokens(ctx, delegator, sendPool, denom, subTokenIDs); err != nil {
				return sdk.ZeroInt(), err
			}
		}
	} else {
		// potentially transfer tokens between pools, if
		switch {
		case tokenSrc == types.BondStatus_Bonded && validator.IsBonded():
			// do nothing
		case (tokenSrc == types.BondStatus_Unbonded || tokenSrc == types.BondStatus_Unbonding) && !validator.IsBonded():
			// do nothing
		case (tokenSrc == types.BondStatus_Unbonded || tokenSrc == types.BondStatus_Unbonding) && validator.IsBonded():
			// transfer pools
			// stake is coin or nft
			switch {
			case coinAmount != nil && subTokenIDs == nil:
				coins := sdk.NewCoins(stake.GetStake())
				k.sendCoinsToBonded(ctx, coins)
			case subTokenIDs != nil && coinAmount == nil:
				if err = k.nftKeeper.TransferSubTokens(ctx, notBondedPool, bondedPool, denom, subTokenIDs); err != nil {
					return sdk.ZeroInt(), err
				}
			}
		case tokenSrc == types.BondStatus_Bonded && !validator.IsBonded():
			// transfer pools
			// stake is coin or nft
			switch {
			case coinAmount != nil && subTokenIDs == nil:
				coins := sdk.NewCoins(stake.GetStake())
				k.sendCoinsToNotBonded(ctx, coins)
			case subTokenIDs != nil && coinAmount == nil:
				if err = k.nftKeeper.TransferSubTokens(ctx, bondedPool, notBondedPool, denom, subTokenIDs); err != nil {
					return sdk.ZeroInt(), err
				}
			}
		default:
			panic("unknown token source bond status")
		}
	}

	// Update delegation

	if found {
		delegation.Stake.Stake = delegation.Stake.Stake.Add(stake.GetStake())

		switch stake.Type {
		case types.StakeType_Coin:
		case types.StakeType_NFT:
			delegation.Stake.SubTokenIDs, err = delegation.Stake.AddSubTokens(stake.SubTokenIDs)
			if err != nil {
				return sdkmath.Int{}, err
			}
		}
	}

	// Update delegation

	k.SetDelegation(ctx, delegation)

	// update validator info
	valAddress := validator.GetOperator()
	totalStake, err = k.TotalStakeInBaseCoin(ctx, valAddress)

	k.DeleteValidatorByPowerIndex(ctx, valAddress, validator.Stake)
	rs, err := k.GetValidatorRS(ctx, valAddress)
	if err != nil {
		return sdkmath.Int{}, err
	}
	rs.Stake = totalStake.Int64()
	k.SetValidatorRS(ctx, valAddress, rs)
	k.SetValidatorByPowerIndex(ctx, valAddress, totalStake.Int64())

	if err = k.AfterUpdateDelegation(ctx, delegation.GetStake().GetStake().Denom, delegation.GetStake().GetStake().Amount); err != nil {
		return sdk.ZeroInt(), err
	}

	return totalStake, nil
}

// Unbond unbonds a particular delegation and perform associated store operations.
//func (k Keeper) Unbond(
//	ctx sdk.Context, delegator sdk.AccAddress, validator types.Validator,
//	denom string, coinAmount *sdkmath.Int, subTokenIDs []uint32,
//) (amount sdkmath.Int, err error) {
//	// check if a delegation object exists in the store
//	delegation, found := k.GetDelegation(ctx, delegator, validator.GetOperator(), denom)
//	if !found {
//		return amount, types.ErrNoDelegatorForAddress
//	}
//
//	// call the before-delegation-modified hook
//	if err := k.BeforeDelegationSharesModified(ctx, delegator, validator.GetOperator()); err != nil {
//		return amount, err
//	}
//
//	//// ensure that we have enough shares to remove
//	//if delegation.Shares.LT(shares) {
//	//	return amount, sdkerrors.Wrap(types.ErrNotEnoughDelegationShares, delegation.Shares.String())
//	//}
//
//	// subtract shares from delegation
//	//delegation.Shares = delegation.Shares.Sub(shares)
//
//	isValidatorOperator := delegator.Equals(validator.GetOperator())
//
//	// If the delegation is the operator of the validator and undelegating will decrease the validator's
//	// self-delegation below their minimum, we jail the validator.
//	if isValidatorOperator && !validator.Jailed &&
//		validator.TokensFromShares(delegation.Shares).TruncateInt().LT(validator.MinSelfDelegation) {
//		k.jailValidator(ctx, validator)
//		validator = k.mustGetValidator(ctx, validator.GetOperator())
//	}
//
//	if delegation.Shares.IsZero() {
//		err = k.RemoveDelegation(ctx, delegation)
//	} else {
//		k.SetDelegation(ctx, delegation)
//		// call the after delegation modification hook
//		err = k.AfterDelegationModified(ctx, delegator, delegation.GetValidatorAddr())
//	}
//
//	if err != nil {
//		return amount, err
//	}
//
//	// remove the shares and coins from the validator
//	// NOTE that the amount is later (in keeper.Delegation) moved between staking module pools
//	validator, amount = k.RemoveValidatorTokensAndShares(ctx, validator, shares)
//	if validator.DelegatorShares.IsZero() && validator.IsUnbonded() {
//		// if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
//		k.RemoveValidator(ctx, validator.GetOperator())
//	}
//
//	return amount, nil
//}
//
//// Undelegate unbonds an amount of delegator shares from a given validator. It
//// will verify that the unbonding entries between the delegator and validator
//// are not exceeded and unbond the staked tokens (based on shares) by creating
//// an unbonding object and inserting it into the unbonding queue which will be
//// processed during the staking EndBlocker.
//func (k Keeper) Undelegate(
//	ctx sdk.Context, delegator sdk.AccAddress, valAddress sdk.ValAddress,
//	denom string, coinAmount *sdkmath.Int, subTokenIDs []uint32,
//) (time.Time, error) {
//	validator, found := k.GetValidator(ctx, valAddress)
//	if !found {
//		return time.Time{}, types.ErrNoDelegatorForAddress
//	}
//
//	if k.HasMaxUndelegationEntries(ctx, delegator, valAddress) {
//		return time.Time{}, types.ErrMaxUndelegationEntries
//	}
//
//	returnAmount, err := k.Unbond(ctx, delegator, validator, sharesAmount)
//	if err != nil {
//		return time.Time{}, err
//	}
//
//	// transfer the validator tokens to the not bonded pool
//	if validator.IsBonded() {
//		k.sendCoinsToNotBonded(ctx, returnAmount)
//	}
//
//	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
//	ubd := k.SetUndelegationEntry(ctx, delegator, validator, ctx.BlockHeight(), completionTime, returnAmount)
//	k.InsertUBDQueue(ctx, ubd, completionTime)
//
//	return completionTime, nil
//}
//
//// CompleteUnbonding completes the unbonding of all mature entries in the
//// retrieved unbonding delegation object and returns the total unbonding balance
//// or an error upon failure.
//func (k Keeper) CompleteUnbonding(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress) (sdk.Coins, error) {
//	ubd, found := k.GetUndelegation(ctx, delegator, validator)
//	if !found {
//		return nil, types.ErrNoUndelegation
//	}
//
//	bondDenom := k.GetParams(ctx).BondDenom
//	balances := sdk.NewCoins()
//	ctxTime := ctx.BlockHeader().Time
//
//	delegator, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
//	if err != nil {
//		return nil, err
//	}
//
//	// loop through all the entries and complete unbonding mature entries
//	for i := 0; i < len(ubd.Entries); i++ {
//		entry := ubd.Entries[i]
//		if entry.IsMature(ctxTime) {
//			ubd.RemoveEntry(int64(i))
//			i--
//
//			// track undelegation only when remaining or truncated shares are non-zero
//			if !entry.Balance.IsZero() {
//				amt := sdk.NewCoin(bondDenom, entry.Balance)
//				if err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(
//					ctx, types.NotBondedPoolName, delegator, sdk.NewCoins(amt),
//				); err != nil {
//					return nil, err
//				}
//
//				balances = balances.Add(amt)
//			}
//		}
//	}
//
//	// set the unbonding delegation or remove it if there are no more entries
//	if len(ubd.Entries) == 0 {
//		k.RemoveUndelegation(ctx, ubd)
//	} else {
//		k.SetUndelegation(ctx, ubd)
//	}
//
//	return balances, nil
//}
//
//// BeginRedelegation begins unbonding / redelegation and creates a redelegation
//// record.
//func (k Keeper) BeginRedelegation(
//	ctx sdk.Context, delegator sdk.AccAddress, validatorSrc, validatorDst sdk.ValAddress, sharesAmount sdk.Dec,
//) (completionTime time.Time, err error) {
//	if bytes.Equal(validatorSrc, validatorDst) {
//		return time.Time{}, types.ErrSelfRedelegation
//	}
//
//	dstValidator, found := k.GetValidator(ctx, validatorDst)
//	if !found {
//		return time.Time{}, types.ErrBadRedelegationDst
//	}
//
//	srcValidator, found := k.GetValidator(ctx, validatorSrc)
//	if !found {
//		return time.Time{}, types.ErrBadRedelegationDst
//	}
//
//	// check if this is a transitive redelegation
//	if k.HasReceivingRedelegation(ctx, delegator, validatorSrc) {
//		return time.Time{}, types.ErrTransitiveRedelegation
//	}
//
//	if k.HasMaxRedelegationEntries(ctx, delegator, validatorSrc, validatorDst) {
//		return time.Time{}, types.ErrMaxRedelegationEntries
//	}
//
//	returnAmount, err := k.Unbond(ctx, delegator, validatorSrc, sharesAmount)
//	if err != nil {
//		return time.Time{}, err
//	}
//
//	if returnAmount.IsZero() {
//		return time.Time{}, types.ErrTinyRedelegationAmount
//	}
//
//	sharesCreated, err := k.Delegate(ctx, delegator, returnAmount, srcValidator.GetStatus(), dstValidator, false)
//	if err != nil {
//		return time.Time{}, err
//	}
//
//	// create the unbonding delegation
//	completionTime, height, completeNow := k.getBeginInfo(ctx, validatorSrc)
//
//	if completeNow { // no need to create the redelegation object
//		return completionTime, nil
//	}
//
//	red := k.SetRedelegationEntry(
//		ctx, delegator, validatorSrc, validatorDst,
//		height, completionTime, returnAmount, sharesAmount, sharesCreated,
//	)
//	k.InsertRedelegationQueue(ctx, red, completionTime)
//
//	return completionTime, nil
//}
//
//// CompleteRedelegation completes the redelegations of all mature entries in the
//// retrieved redelegation object and returns the total redelegation (initial)
//// balance or an error upon failure.
//func (k Keeper) CompleteRedelegation(
//	ctx sdk.Context, delegator sdk.AccAddress, validatorSrc, validatorDst sdk.ValAddress,
//) (sdk.Coins, error) {
//	red, found := k.GetRedelegation(ctx, delegator, validatorSrc, validatorDst)
//	if !found {
//		return nil, types.ErrNoRedelegation
//	}
//
//	bondDenom := k.GetParams(ctx).BondDenom
//	balances := sdk.NewCoins()
//	ctxTime := ctx.BlockHeader().Time
//
//	// loop through all the entries and complete mature redelegation entries
//	for i := 0; i < len(red.Entries); i++ {
//		entry := red.Entries[i]
//		if entry.IsMature(ctxTime) {
//			red.RemoveEntry(int64(i))
//			i--
//
//			if !entry.InitialBalance.IsZero() {
//				balances = balances.Add(sdk.NewCoin(bondDenom, entry.InitialBalance))
//			}
//		}
//	}
//
//	// set the redelegation or remove it if there are no more entries
//	if len(red.Entries) == 0 {
//		k.RemoveRedelegation(ctx, red)
//	} else {
//		k.SetRedelegation(ctx, red)
//	}
//
//	return balances, nil
//}
//
//// ValidateUnbondAmount validates that a given unbond or redelegation amount is
//// valied based on upon the converted shares. If the amount is valid, the total
//// amount of respective shares is returned, otherwise an error is returned.
//func (k Keeper) ValidateUnbondAmount(
//	ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress, amt sdkmath.Int,
//) (shares sdk.Dec, err error) {
//	validator, found := k.GetValidator(ctx, validator)
//	if !found {
//		return shares, types.ErrNoValidatorFound
//	}
//
//	del, found := k.GetDelegation(ctx, delegator, validator)
//	if !found {
//		return shares, types.ErrNoDelegation
//	}
//
//	shares, err = validator.SharesFromTokens(amt)
//	if err != nil {
//		return shares, err
//	}
//
//	sharesTruncated, err := validator.SharesFromTokensTruncated(amt)
//	if err != nil {
//		return shares, err
//	}
//
//	delShares := del.GetShares()
//	if sharesTruncated.GT(delShares) {
//		return shares, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid shares amount")
//	}
//
//	// Cap the shares at the delegation's shares. Shares being greater could occur
//	// due to rounding, however we don't want to truncate the shares or take the
//	// minimum because we want to allow for the full withdraw of shares from a
//	// delegation.
//	if shares.GT(delShares) {
//		shares = delShares
//	}
//
//	return shares, nil
//}
//
//// getBeginInfo returns the completion time and height of a redelegation, along with
//// a boolean signaling if the redelegation is complete based on the source validator.
//func (k Keeper) getBeginInfo(ctx sdk.Context, validatorSrc sdk.ValAddress) (completionTime time.Time, height int64, completeNow bool) {
//	validator, found := k.GetValidator(ctx, validatorSrc)
//
//	// TODO: When would the validator not be found?
//	switch {
//	case !found || validator.IsBonded():
//		// the longest wait - just unbonding period from now
//		completionTime = ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
//		height = ctx.BlockHeight()
//
//		return completionTime, height, false
//
//	case validator.IsUnbonded():
//		return completionTime, height, true
//
//	case validator.IsUnbonding():
//		return validator.UnbondingTime, validator.UnbondingHeight, false
//
//	default:
//		panic(fmt.Sprintf("unknown validator status: %s", validator.Status))
//	}
//}

func (k Keeper) TotalStakeInBaseCoin(ctx sdk.Context, valAddress sdk.ValAddress) (sdkmath.Int, error) {
	delegations := k.GetValidatorDelegations(ctx, valAddress)
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
