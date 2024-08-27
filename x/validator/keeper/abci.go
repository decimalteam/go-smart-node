package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// BeginBlocker will persist the current header and validator set as a historical entry
// and prune the oldest entry based on the HistoricalEntries parameter
func BeginBlocker(ctx sdk.Context, k Keeper, req abci.RequestBeginBlock) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

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

	k.TrackHistoricalInfo(ctx)
}

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k Keeper, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	start := time.Now()
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	height := ctx.BlockHeight()

	updates := k.BlockValidatorUpdates(ctx)

	// send coin reserve
	fmt.Println("coins")
	nbPool := k.GetNotBondedPool(ctx).GetAddress()
	coins := k.bankKeeper.GetAllBalances(ctx, nbPool)
	fmt.Println(coins)

	bPool := k.GetBondedPool(ctx).GetAddress()
	coins = k.bankKeeper.GetAllBalances(ctx, bPool)
	fmt.Println(coins.String())

	k.PayValidators(ctx)

	if height%120 == 0 {
		ctx.Logger().Debug(
			fmt.Sprintf("Duration 120 block (%s)", helpers.DurationToString(time.Since(start))),
		)
		err := k.PayRewards(ctx)
		if err != nil {
			panic(err)
		}

		dataAddress, err := contracts.GetAddressFromContractCenter(ctx, k.evmKeeper, contracts.NameOfSlugForGetAddressDelegation)
		if err == nil {
			params := k.GetParams(ctx)
			undelegationTime, err := contracts.GetTimeUndelegate(ctx, k.evmKeeper, common.HexToAddress(dataAddress))
			if err == nil {
				params.UndelegationTime = time.Second * time.Duration(undelegationTime.Int64())
			}
			redelegationTime, err := contracts.GetTimeRedelegation(ctx, k.evmKeeper, common.HexToAddress(dataAddress))
			if err == nil {
				params.RedelegationTime = time.Second * time.Duration(redelegationTime.Int64())
			}
			k.SetParams(ctx, params)
		}
		k.DeleteHoldMature(ctx)
	} else {
		ctx.Logger().Debug(
			fmt.Sprintf("Duration simple block (%s)", helpers.DurationToString(time.Since(start))),
		)
	}
	return updates
}

func (k Keeper) PayValidators(ctx sdk.Context) {
	height := ctx.BlockHeight()

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

	if baseCoin.LimitVolume.IsZero() {
		baseCoin.LimitVolume = types.GetAllEmission(ctx)
	} else {
		baseCoin.LimitVolume = baseCoin.LimitVolume.Add(rewards)
	}

	k.coinKeeper.SetCoin(ctx, baseCoin)

	feeCollector := k.authKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	feesCollectedCoins := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	for _, fee := range feesCollectedCoins {
		if fee.Denom == baseDenom {
			rewards = rewards.Add(fee.Amount)
			continue
		}

		coin, err := k.coinKeeper.GetCoin(ctx, fee.Denom)
		if err != nil {
			panic(err)
		}

		baseAmount := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, uint(coin.CRR), fee.Amount)
		rewards = rewards.Add(baseAmount)
	}
	err = k.coinKeeper.BurnPoolCoins(ctx, authtypes.FeeCollectorName, feesCollectedCoins)
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

	// create coins for delegators
	// remainder to FeeCollector
	err = k.bankKeeper.MintCoins(ctx, authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(k.BaseDenom(ctx), remainder)))
	// distributed to validator module for delegators
	distributed := rewards.Sub(remainder)
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(baseDenom, distributed)))
	if err != nil {
		panic(err)
	}
	err = k.coinKeeper.UpdateCoinVR(ctx, baseDenom, baseCoin.Volume.Add(distributed), baseCoin.Reserve)
	if err != nil {
		panic(err)
	}
}
