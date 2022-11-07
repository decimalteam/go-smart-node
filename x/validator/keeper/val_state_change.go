package keeper

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"sort"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/evmos/ethermint/types"
	gogotypes "github.com/gogo/protobuf/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// BlockValidatorUpdates calculates the ValidatorUpdates for the current block
// Called in each EndBlock
func (k Keeper) BlockValidatorUpdates(ctx sdk.Context) []abci.ValidatorUpdate {
	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	validatorUpdates, err := k.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		panic(err)
	}

	// unbond all mature validators from the unbonding queue
	k.UnbondAllMatureValidators(ctx)

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		addr, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress := sdk.MustAccAddressFromBech32(dvPair.DelegatorAddress)

		err = k.CompleteUnbonding(ctx, delegatorAddress, addr)
		if err != nil {
			continue
		}
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx)
	for _, dvvTriplet := range matureRedelegations {
		valSrcAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorSrcAddress)
		if err != nil {
			panic(err)
		}
		valDstAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress := sdk.MustAccAddressFromBech32(dvvTriplet.DelegatorAddress)

		err = k.CompleteRedelegation(
			ctx,
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
		)
		if err != nil {
			continue
		}
	}

	return validatorUpdates
}

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
	params := k.GetParams(ctx)
	maxValidators := params.MaxValidators
	totalPower := sdk.ZeroInt()
	amtFromBondedToNotBonded, amtFromNotBondedToBonded := sdk.NewCoins(), sdk.NewCoins()
	nftsFromBondedToNotBonded, nftsFromNotBondedToBonded := []nftTransferRecord{}, []nftTransferRecord{}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("stacktrace from panic: %s \n%s\n", r, string(debug.Stack()))
		}
	}()

	// Retrieve the last validator set.
	// The persistent set is updated later in this function.
	// (see LastValidatorPowerKey).
	last, err := k.getLastValidatorsByAddr(ctx)
	if err != nil {
		return nil, err
	}

	validators := k.GetAllValidatorsByPowerIndexReversed(ctx)
	delegations := k.GetAllDelegationsByValidator(ctx)
	for _, validator := range validators {
		if validator.Jailed {
			continue
		}
		k.CheckDelegations(ctx, validator, delegations[validator.GetOperator().String()])
	}

	// Iterate over validators, highest power to lowest.
	iterator := k.ValidatorsPowerStoreIterator(ctx)
	defer iterator.Close()
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
		// everything that is iterated in this loop is becoming or already a
		// part of the bonded validator set
		valAddr := sdk.ValAddress(iterator.Value()).String()
		validator := k.mustGetValidator(ctx, sdk.ValAddress(iterator.Value()))
		if validator.Jailed {
			return nil, fmt.Errorf("ApplyAndReturnValidatorSetUpdates: should never retrieve a jailed validator from the power store")
		}

		// Unbonding case: no delegations, but status != Unbonding
		// stake and consensus power = 0
		if len(delegations[valAddr]) == 0 && validator.Status != types.BondStatus_Unbonding {
			// force to set offline
			validator.Online = false
			validator, err = k.beginUnbondingValidator(ctx, validator)
			if err != nil {
				return nil, err
			}
			// continue to add validator to noLongerBonded, and remove from power index
			continue
		}

		// if we get to a zero-power validator (which we don't bond),
		// there are no more possible bonded validators
		if validator.PotentialConsensusPower() == 0 {
			// continue for potential unbonding validators
			continue
		}

		// apply the appropriate state change if necessary
		switch {
		case validator.IsUnbonded():
			if !validator.Online {
				break // break switch
			}

			validator, err = k.unbondedToBonded(ctx, validator)
			if err != nil {
				return
			}
			for _, delegation := range delegations[valAddr] {
				switch delegation.Stake.Type {
				case types.StakeType_Coin:
					amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(delegation.GetStake().GetStake())
				case types.StakeType_NFT:
					nftsFromNotBondedToBonded = append(nftsFromNotBondedToBonded, nftTransferRecord{
						tokenID:     delegation.GetStake().GetID(),
						subTokenIDs: delegation.GetStake().GetSubTokenIDs(),
					})
				}
			}
		case validator.IsUnbonding():
			if !validator.Online {
				break // break switch
			}

			validator, err = k.unbondingToBonded(ctx, validator)
			if err != nil {
				return
			}
			for _, delegation := range delegations[valAddr] {
				switch delegation.Stake.Type {
				case types.StakeType_Coin:
					amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(delegation.GetStake().GetStake())
				case types.StakeType_NFT:
					nftsFromNotBondedToBonded = append(nftsFromNotBondedToBonded, nftTransferRecord{
						tokenID:     delegation.GetStake().GetID(),
						subTokenIDs: delegation.GetStake().GetSubTokenIDs(),
					})
				}
			}
		case validator.IsBonded():
			if validator.Online {
				break // break switch
			}

			validator, err = k.bondedToUnbonded(ctx, validator)
			if err != nil {
				return
			}
			for _, delegation := range delegations[valAddr] {
				switch delegation.Stake.Type {
				case types.StakeType_Coin:
					amtFromBondedToNotBonded = amtFromBondedToNotBonded.Add(delegation.GetStake().GetStake())
				case types.StakeType_NFT:
					nftsFromBondedToNotBonded = append(nftsFromBondedToNotBonded, nftTransferRecord{
						tokenID:     delegation.GetStake().GetID(),
						subTokenIDs: delegation.GetStake().GetSubTokenIDs(),
					})
				}
			}
		default:
			panic("unexpected validator status")
		}

		// fetch the old power bytes
		oldPowerBytes, found := last[valAddr]

		newPower := validator.ConsensusPower()
		newPowerBytes := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: newPower})

		// update the validator set if power has changed
		if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
			if newPower > 0 {
				k.SetLastValidatorPower(ctx, validator.GetOperator(), newPower)
				updates = append(updates, validator.ABCIValidatorUpdate(ethtypes.PowerReduction))
			} else {
				k.DeleteLastValidatorPower(ctx, validator.GetOperator())
				updates = append(updates, validator.ABCIValidatorUpdateZero())
			}
		}

		delete(last, valAddr)
		count++

		totalPower = totalPower.Add(sdk.NewInt(newPower))
	}

	noLongerBonded, err := sortNoLongerBonded(last)
	if err != nil {
		return nil, err
	}

	for _, valAddrBytes := range noLongerBonded {
		validator := k.mustGetValidator(ctx, sdk.ValAddress(valAddrBytes))
		k.DeleteLastValidatorPower(ctx, validator.GetOperator())
		updates = append(updates, validator.ABCIValidatorUpdateZero())
	}

	// Update the pools based on the recent updates in the validator set:
	// - The tokens from the non-bonded candidates that enter the new validator set need to be transferred
	// to the Bonded pool.
	// - The tokens from the bonded validators that are being kicked out from the validator set
	// need to be transferred to the NotBonded pool.

	err = k.transferBetweenPools(ctx, types.BondStatus_Bonded, types.BondStatus_Unbonded, amtFromBondedToNotBonded, nftsFromBondedToNotBonded)
	if err != nil {
		return nil, err
	}
	err = k.transferBetweenPools(ctx, types.BondStatus_Unbonded, types.BondStatus_Bonded, amtFromNotBondedToBonded, nftsFromNotBondedToBonded)
	if err != nil {
		return nil, err
	}

	// set total power on lookup index if there are any updates
	if len(updates) > 0 {
		k.SetLastTotalPower(ctx, totalPower)
	}

	return updates, err
}

// Validator state transitions

func (k Keeper) bondedToUnbonding(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsBonded() {
		panic(fmt.Sprintf("bad state transition bondedToUnbonding, validator: %v\n", validator))
	}

	return k.beginUnbondingValidator(ctx, validator)
}

func (k Keeper) unbondingToBonded(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsUnbonding() {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
	}

	return k.bondValidator(ctx, validator)
}

func (k Keeper) unbondedToBonded(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsUnbonded() {
		panic(fmt.Sprintf("bad state transition unbondedToBonded, validator: %v\n", validator))
	}

	return k.bondValidator(ctx, validator)
}

func (k Keeper) bondedToUnbonded(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsBonded() {
		panic(fmt.Sprintf("bad state transition bondedToUnbonded, validator: %v\n", validator))
	}

	return k.unbondValidator(ctx, validator)
}

// UnbondingToUnbonded switches a validator from unbonding state to unbonded state

func (k Keeper) UnbondingToUnbonded(ctx sdk.Context, validator types.Validator) types.Validator {
	if !validator.IsUnbonding() {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
	}

	return k.completeUnbondingValidator(ctx, validator)
}

// perform all the store operations for when a validator status becomes bonded

func (k Keeper) bondValidator(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	// delete the validator by power index, as the key will change
	k.DeleteValidatorByPowerIndex(ctx, validator)

	validator = validator.UpdateStatus(types.BondStatus_Bonded)

	// save the now bonded validator record to the two referenced stores
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)

	// delete from queue if present
	k.DeleteValidatorQueue(ctx, validator)

	// trigger hook //TODO needed hooks?
	//consAddr, err := validator.GetConsAddr()
	//if err != nil {
	//	return validator, err
	//}
	//k.AfterValidatorBonded(ctx, consAddr, validator.GetOperator())

	return validator, nil
}

// force unbond status when validator become offline
func (k Keeper) unbondValidator(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	// delete the validator by power index, as the key will change
	k.DeleteValidatorByPowerIndex(ctx, validator)

	valAdr := validator.GetOperator()
	validator = validator.UpdateStatus(types.BondStatus_Unbonded)
	validator.Stake = 0
	rs, err := k.GetValidatorRS(ctx, valAdr)
	if err != nil {
		return types.Validator{}, err
	}
	rs.Stake = 0
	k.SetValidatorRS(ctx, valAdr, rs)

	// save the now bonded validator record to the two referenced stores
	k.SetValidator(ctx, validator)

	// delete from unbonding queue if present
	k.DeleteValidatorQueue(ctx, validator)

	return validator, nil
}

// perform all the store operations for when a validator begins unbonding

func (k Keeper) beginUnbondingValidator(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	params := k.GetParams(ctx)

	// delete the validator by power index, as the key will change
	k.DeleteValidatorByPowerIndex(ctx, validator)

	validator = validator.UpdateStatus(types.BondStatus_Unbonding)
	valAdr := validator.GetOperator()
	validator.Stake = 0
	rs, err := k.GetValidatorRS(ctx, valAdr)
	if err != nil {
		return types.Validator{}, err
	}
	rs.Stake = 0
	k.SetValidatorRS(ctx, valAdr, rs)

	// set the unbonding completion time and completion height appropriately
	validator.UnbondingTime = ctx.BlockHeader().Time.Add(params.UndelegationTime)
	validator.UnbondingHeight = ctx.BlockHeader().Height

	// save the now unbonded validator record and power index
	k.SetValidator(ctx, validator)
	//k.SetValidatorByPowerIndex(ctx, validator)

	// Adds to unbonding validator queue
	k.InsertUnbondingValidatorQueue(ctx, validator)

	// trigger hook
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return validator, err
	}
	k.AfterValidatorBeginUnbonding(ctx, consAddr, validator.GetOperator())

	return validator, nil
}

// perform all the store operations for when a validator status becomes unbonded

func (k Keeper) completeUnbondingValidator(ctx sdk.Context, validator types.Validator) types.Validator {
	validator = validator.UpdateStatus(types.BondStatus_Unbonded)
	k.SetValidator(ctx, validator)

	return validator
}

// map of operator bech32-addresses to serialized power
// We use bech32 strings here, because we can't have slices as keys: map[[]byte][]byte
type validatorsByAddr map[string][]byte

// get the last validator set
func (k Keeper) getLastValidatorsByAddr(ctx sdk.Context) (validatorsByAddr, error) {
	last := make(validatorsByAddr)

	iterator := k.LastValidatorsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// extract the validator address from the key (prefix is 1-byte, addrLen is 1-byte)
		valAddr := types.AddressFromLastValidatorPowerKey(iterator.Key())
		valAddrStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
		if err != nil {
			return nil, err
		}

		powerBytes := iterator.Value()
		last[valAddrStr] = make([]byte, len(powerBytes))
		copy(last[valAddrStr], powerBytes)
	}

	return last, nil
}

// given a map of remaining validators to previous bonded power
// returns the list of validators to be unbonded, sorted by operator address
func sortNoLongerBonded(last validatorsByAddr) ([][]byte, error) {
	// sort the map keys for determinism
	noLongerBonded := make([][]byte, len(last))
	index := 0

	for valAddrStr := range last {
		valAddrBytes, err := sdk.ValAddressFromBech32(valAddrStr)
		if err != nil {
			return nil, err
		}
		noLongerBonded[index] = valAddrBytes
		index++
	}
	// sorted by address - order doesn't matter
	sort.SliceStable(noLongerBonded, func(i, j int) bool {
		// -1 means strictly less than
		return bytes.Compare(noLongerBonded[i], noLongerBonded[j]) == -1
	})

	return noLongerBonded, nil
}

func (k Keeper) CheckDelegations(ctx sdk.Context, validator types.Validator, delegations []types.Delegation) {
	if len(delegations) <= int(k.MaxDelegations(ctx)) {
		return
	}

	baseAmounts := make(map[int]sdkmath.Int)
	for i := range delegations {
		baseAmounts[i] = k.ToBaseCoin(ctx, delegations[i].Stake.GetStake()).Amount
	}

	sort.SliceStable(delegations, func(i, j int) bool {
		amountI := baseAmounts[i]
		amountJ := baseAmounts[j]
		return amountI.GT(amountJ)
	})

	for i := int(k.MaxDelegations(ctx)); i < len(delegations); i++ {
		stake := delegations[i].Stake
		delegator := delegations[i].GetDelegator()
		switch validator.Status {
		case types.BondStatus_Bonded:
			switch stake.Type {
			case types.StakeType_Coin:
				amt := stake.Stake

				if err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(
					ctx, types.BondedPoolName, delegator, sdk.NewCoins(amt),
				); err != nil {
					panic(err)
				}

			case types.StakeType_NFT:
				if err := k.nftKeeper.TransferSubTokens(ctx, k.GetBondedPool(ctx).GetAddress(), delegator, stake.GetID(), stake.GetSubTokenIDs()); err != nil {
					panic(err)
				}
			}
		case types.BondStatus_Unbonded:
			switch stake.Type {
			case types.StakeType_Coin:
				amt := stake.Stake

				if err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(
					ctx, types.NotBondedPoolName, delegator, sdk.NewCoins(amt),
				); err != nil {
					panic(err)
				}

			case types.StakeType_NFT:
				if err := k.nftKeeper.TransferSubTokens(ctx, k.GetNotBondedPool(ctx).GetAddress(), delegator, stake.GetID(), stake.GetSubTokenIDs()); err != nil {
					panic(err)
				}
			}
		}

		remainStake, err := k.CalculateRemainStake(ctx, delegations[i].Stake, delegations[i].Stake)
		if err != nil {
			panic(err)
		}

		err = k.Unbond(ctx, delegator, validator.GetOperator(), delegations[i].Stake, remainStake)
		if err != nil {
			panic(err)
		}
	}
}
