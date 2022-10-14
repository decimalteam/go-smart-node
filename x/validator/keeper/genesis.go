package keeper

import (
	"fmt"

	ethtypes "github.com/evmos/ethermint/types"
	abci "github.com/tendermint/tendermint/abci/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// InitGenesis sets the pool and parameters for the provided keeper.
// For each validator in data, it sets that validator in the keeper along with manually setting the indexes.
// In addition, it also sets any delegations found in data. Finally, it updates the bonded validators.
// Returns final validator set after applying all declaration and delegations
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) (res []abci.ValidatorUpdate) {
	if err := data.UnpackInterfaces(k.cdc); err != nil {
		panic(err)
	}

	bondedCoins := sdk.NewCoins()
	notBondedCoins := sdk.NewCoins()

	// We need to pretend to be "n blocks before genesis", where "n" is the validator update delay,
	// so that e.g. slashing periods are correctly initialized for the validator set e.g. with a one-block
	// offset - the first TM block is at height 1, so state updates applied from genesis.json are in block 0.
	ctx = ctx.WithBlockHeight(1 - sdk.ValidatorUpdateDelay)

	k.SetParams(ctx, data.Params)
	k.SetLastTotalPower(ctx, sdkmath.NewInt(data.LastTotalPower))

	var valStatus = make(map[string]types.BondStatus)

	for _, validator := range data.Validators {
		k.SetValidator(ctx, validator)

		// Manually set indices for the first time
		k.SetValidatorByConsAddr(ctx, validator)
		//k.SetValidatorByPowerIndex(ctx, validator.GetOperator(), validator.Stake)
		k.SetValidatorRS(ctx, validator.GetOperator(), types.ValidatorRS{
			Rewards:      validator.Rewards,
			TotalRewards: validator.Rewards,
			Stake:        validator.Stake,
		})
		var hasPower bool
		for _, lp := range data.LastValidatorPowers {
			if lp.Address == validator.OperatorAddress {
				hasPower = true
				break
			}
		}
		if hasPower {
			k.SetValidatorByPowerIndex(ctx, validator.GetOperator(), validator.Stake)
		}

		// Call the creation hook if not exported
		if !data.Exported {
			if err := k.AfterValidatorCreated(ctx, validator.GetOperator()); err != nil {
				panic(err)
			}
		}

		// update timeslice if necessary
		if validator.IsUnbonding() {
			k.InsertUnbondingValidatorQueue(ctx, validator)
		}

		valStatus[validator.OperatorAddress] = validator.GetStatus()
	}

	coinMap := make(map[string]bool)
	coinMap[k.BaseDenom(ctx)] = true
	for _, coinInfo := range k.coinKeeper.GetCoins(ctx) {
		coinMap[coinInfo.Denom] = true
	}

	for _, delegation := range data.Delegations {
		delegatorAddress := sdk.MustAccAddressFromBech32(delegation.Delegator)

		// Call the before-creation hook if not exported
		if !data.Exported {
			if err := k.BeforeDelegationCreated(ctx, delegatorAddress, delegation.GetValidator()); err != nil {
				panic(err)
			}
		}

		if err := k.checkStake(ctx, delegation.Stake, coinMap); err != nil {
			panic(err)
		}

		k.SetDelegation(ctx, delegation)

		// Call the after-modification hook if not exported
		if !data.Exported {
			if err := k.AfterDelegationModified(ctx, delegatorAddress, delegation.GetValidator()); err != nil {
				panic(err)
			}
		}

		if delegation.Stake.Type == types.StakeType_Coin {
			switch valStatus[delegation.Validator] {
			case types.BondStatus_Bonded:
				bondedCoins = bondedCoins.Add(delegation.Stake.Stake)
			case types.BondStatus_Unbonding, types.BondStatus_Unbonded:
				notBondedCoins = notBondedCoins.Add(delegation.Stake.Stake)
			default:
				panic(fmt.Errorf("invalid validator %s status", delegation.Validator))
			}
		}
	}

	for _, ubd := range data.Undelegations {
		for _, entry := range ubd.Entries {
			if err := k.checkStake(ctx, entry.Stake, coinMap); err != nil {
				panic(err)
			}
		}
		k.SetUndelegation(ctx, ubd)

		for _, entry := range ubd.Entries {
			k.InsertUBDQueue(ctx, ubd, entry.CompletionTime)
			notBondedCoins = notBondedCoins.Add(entry.Stake.Stake)
		}
	}

	for _, red := range data.Redelegations {
		for _, entry := range red.Entries {
			if err := k.checkStake(ctx, entry.Stake, coinMap); err != nil {
				panic(err)
			}
		}
		k.SetRedelegation(ctx, red)

		for _, entry := range red.Entries {
			k.InsertRedelegationQueue(ctx, red, entry.CompletionTime)
		}
	}

	// check if the unbonded and bonded pools accounts exists
	bondedPool := k.GetBondedPool(ctx)
	if bondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	// TODO: remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	bondedBalance := k.bankKeeper.GetAllBalances(ctx, bondedPool.GetAddress())
	if bondedBalance.IsZero() {
		k.authKeeper.SetModuleAccount(ctx, bondedPool)
	}

	// if balance is different from bonded coins panic because genesis is most likely malformed
	if !bondedBalance.IsEqual(bondedCoins) {
		panic(fmt.Sprintf("bonded pool balance is different from bonded coins: %s <-> %s", bondedBalance, bondedCoins))
	}

	notBondedPool := k.GetNotBondedPool(ctx)
	if notBondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	notBondedBalance := k.bankKeeper.GetAllBalances(ctx, notBondedPool.GetAddress())
	if notBondedBalance.IsZero() {
		k.authKeeper.SetModuleAccount(ctx, notBondedPool)
	}

	// If balance is different from non bonded coins panic because genesis is most
	// likely malformed.
	if !notBondedBalance.IsEqual(notBondedCoins) {
		panic(fmt.Sprintf("not bonded pool balance is different from not bonded coins: %s <-> %s", notBondedBalance, notBondedCoins))
	}

	// don't need to run Tendermint updates if we exported
	//if data.Exported {
	for _, lv := range data.LastValidatorPowers {
		valAddr, err := sdk.ValAddressFromBech32(lv.Address)
		if err != nil {
			panic(err)
		}

		k.SetLastValidatorPower(ctx, valAddr, lv.Power)
		validator, found := k.GetValidator(ctx, valAddr)

		if !found {
			panic(fmt.Sprintf("validator %s not found", lv.Address))
		}

		update := validator.ABCIValidatorUpdate(ethtypes.PowerReduction)
		update.Power = lv.Power // keep the next-val-set offset, use the last power for the first block

		res = append(res, update)
	}
	//} else {
	//	var err error
	// res, err = k.ApplyAndReturnValidatorSetUpdates(ctx)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	return res
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, params, validators, and bonds found in
// the keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	var undelegations []types.Undelegation

	k.IterateUndelegations(ctx, func(_ int64, ubd types.Undelegation) (stop bool) {
		undelegations = append(undelegations, ubd)
		return false
	})

	var redelegations []types.Redelegation

	k.IterateRedelegations(ctx, func(_ int64, red types.Redelegation) (stop bool) {
		redelegations = append(redelegations, red)
		return false
	})

	var lastValidatorPowers []types.LastValidatorPower

	k.IterateLastValidatorPowers(ctx, func(addr sdk.ValAddress, power int64) (stop bool) {
		lastValidatorPowers = append(lastValidatorPowers, types.LastValidatorPower{Address: addr.String(), Power: power})
		return false
	})

	return &types.GenesisState{
		Params:              k.GetParams(ctx),
		LastTotalPower:      k.GetLastTotalPower(ctx).Int64(),
		LastValidatorPowers: lastValidatorPowers,
		Validators:          k.GetAllValidators(ctx),
		Delegations:         k.GetAllDelegations(ctx),
		Undelegations:       undelegations,
		Redelegations:       redelegations,
		Exported:            true,
	}
}

func (k Keeper) checkStake(ctx sdk.Context, stake types.Stake, coinMap map[string]bool) error {
	switch stake.Type {
	case types.StakeType_Coin:
		if !coinMap[stake.ID] {
			return fmt.Errorf("coin '%s' does not exists", stake.ID)
		}
	case types.StakeType_NFT:
		for _, subID := range stake.SubTokenIDs {
			_, found := k.nftKeeper.GetSubToken(ctx, stake.ID, subID)
			if !found {
				return fmt.Errorf("token '%s' subtoken '%d' does not exists", stake.ID, subID)
			}
		}
	}
	return nil
}
