package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// BeginBlocker will persist the current header and validator set as a historical entry
// and prune the oldest entry based on the HistoricalEntries parameter
func BeginBlocker(ctx sdk.Context, k Keeper, req abci.RequestBeginBlock) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	/*
		params := k.GetParams(ctx)

		// Iterate over all the validators which *should* have signed this block
		// store whether or not they have actually signed it and slash/unbond any
		// which have missed too many blocks in a row (downtime slashing)
		for _, voteInfo := range req.LastCommitInfo.GetVotes() {
			k.HandleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, voteInfo.SignedLastBlock, params)
		}

		// Iterate through any newly discovered evidence of infraction
		// Slash any validators (and since-unbonded stake within the unbonding period)
		// who contributed to valid infractions
		for _, evidence := range req.ByzantineValidators {
			switch evidence.Type {
			case abci.EvidenceType_DUPLICATE_VOTE:
				k.HandleDoubleSign(ctx, evidence.Validator.Address, evidence.Height, evidence.Time, evidence.Validator.Power, params)
			default:
				k.Logger(ctx).Error(fmt.Sprintf("ignored unknown evidence type: %s", evidence.Type))
			}
		}
	*/
	k.TrackHistoricalInfo(ctx)
}

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k Keeper, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	height := ctx.BlockHeight()

	/*
			validatorUpdates, err := k.ApplyAndReturnValidatorSetUpdates(ctx)
			if err != nil {
				panic(err)
			}
			return []abci.ValidatorUpdate{} // k.BlockValidatorUpdates(ctx)



		// Unbond all mature validators from the unbonding queue.
		k.UnbondAllMatureValidatorQueue(ctx)

		//Remove all mature unbonding delegations from the ubd queue.
		matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
		for _, dvPair := range matureUnbonds {
			delAddr := sdk.MustAccAddressFromBech32(dvPair.DelegatorAddress)
			valAddr, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
			if err != nil {
				panic(err)
			}

			_, found := k.GetUndelegation(ctx, delAddr, valAddr)
			if !found {
				continue
			}

			err = k.CompleteUnbonding(ctx, dvPair.DelegatorAddress, dvPair.ValidatorAddress)
			if err != nil {
				continue
			}

			//ctxTime := ctx.BlockHeader().Time

			//ctx.EventManager().EmitEvents(delegation.GetEvents(ctxTime))
		}
	*/
	// calculate emmission
	rewards := types.GetRewardForBlock(uint64(height))

	err := ctx.EventManager().EmitTypedEvents(&types.EventEmission{
		Amount: rewards,
	})
	if err != nil {
		panic(err)
	}

	// calculate rewards
	baseDenom := k.BaseDenom(ctx)
	baseCoin, err := k.coinKeeper.GetCoin(ctx, baseDenom)
	if err != nil {
		panic(err)
	}

	feeCollector := k.authKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	feesCollectedCoins := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	for _, fee := range feesCollectedCoins {
		feeInBaseCoin := k.ToBaseCoin(ctx, fee)

		rewards.Add(feeInBaseCoin.Amount)
	}
	err = k.coinKeeper.BurnPoolCoins(ctx, authtypes.FeeCollectorName, feesCollectedCoins)
	if err != nil {
		panic(err)
	}

	// create coins for delegators
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(baseDenom, rewards)))
	if err != nil {
		panic(err)
	}
	err = k.coinKeeper.UpdateCoinVR(ctx, baseDenom, baseCoin.Reserve, baseCoin.Volume.Add(rewards))
	if err != nil {
		panic(err)
	}

	// pay rewards to validators
	remainder := sdk.NewIntFromBigInt(rewards.BigInt())

	vals, powers, totalPower := k.GetAllValidatorsByPowerIndex(ctx)

	for i, val := range vals {
		if !val.Online {
			continue
		}
		power := powers[i]

		r := sdk.ZeroInt()
		r = rewards.Mul(sdk.NewInt(power)).Quo(totalPower)
		remainder = remainder.Sub(r)
		err = k.AddAccumRewards(ctx, val.GetOperator(), r)
		if err != nil {
			panic(err)
		}
	}

	err = k.bankKeeper.MintCoins(ctx, authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(k.BaseDenom(ctx), remainder)))
	if err != nil {
		panic(err)
	}

	//if height%120 == 0 {
	//	err = k.PayRewards(ctx, totalPower)
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	return []abci.ValidatorUpdate{}
}
