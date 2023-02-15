package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func TestPayValidators(t *testing.T) {
	const legacyRewardAddress = "dx14elhyzmq95f98wrkvujtsr5cyudffp6q2hfkhs"

	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)
	nbPool := dsc.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress()
	//bPool := dsc.ValidatorKeeper.GetBondedPool(ctx).GetAddress()
	valK := dsc.ValidatorKeeper

	// 0. genesis
	balance := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100000000)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(balance))

	genesisVals := dsc.ValidatorKeeper.GetValidators(ctx, 10)
	require.Len(t, genesisVals, 1)
	genesisVal := genesisVals[0]
	require.True(t, genesisVal.ConsensusPower() > 0)
	require.Equal(t, genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
	balanceNB := dsc.BankKeeper.GetAllBalances(ctx, nbPool)
	require.True(t, balanceNB.IsZero())
	//startBalanceB := dsc.BankKeeper.GetAllBalances(ctx, bPool)
	//balanceB := dsc.BankKeeper.GetAllBalances(ctx, bPool)

	//
	ccDenom1 := "custom"
	initVolume1 := keeper.TokensFromConsensusPower(100000000)
	initReserve1 := keeper.TokensFromConsensusPower(1000)
	limitVolume1 := keeper.TokensFromConsensusPower(1000000000000000000)
	crr := uint64(50)

	_, err := dsc.CoinKeeper.CreateCoin(ctx, cointypes.NewMsgCreateCoin(accs[0], ccDenom1, "d", crr, initVolume1, initReserve1, limitVolume1, sdkmath.ZeroInt(), ""))
	require.NoError(t, err)

	ccDenom2 := "custom2"
	initVolume2 := keeper.TokensFromConsensusPower(100000000)
	initReserve2 := keeper.TokensFromConsensusPower(3000)
	limitVolume2 := keeper.TokensFromConsensusPower(1000000000000000)

	_, err = dsc.CoinKeeper.CreateCoin(ctx, cointypes.NewMsgCreateCoin(accs[0], ccDenom2, "dd", crr, initVolume2, initReserve2, limitVolume2, sdkmath.ZeroInt(), ""))
	require.NoError(t, err)

	goCtx := sdk.WrapSDKContext(ctx)

	// first val create
	firstStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(333)))
	{
		creatorStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(333)))
		msgCreate, err := types.NewMsgCreateValidator(vals[0], accs[0], PKs[0], types.Description{Moniker: "monik"},
			sdk.ZeroDec(), creatorStake)
		require.NoError(t, err)

		_, err = msgsrv.CreateValidator(goCtx, msgCreate)
		require.NoError(t, err)

		msgOnline := types.NewMsgSetOnline(vals[0])
		_, err = msgsrv.SetOnline(goCtx, msgOnline)
		require.NoError(t, err)
	}
	// second val create
	secondStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(666)))
	{
		creatorStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(666)))
		msgCreate, err := types.NewMsgCreateValidator(vals[1], accs[1], PKs[1], types.Description{Moniker: "monik1"},
			sdk.ZeroDec(), creatorStake)
		require.NoError(t, err)

		_, err = msgsrv.CreateValidator(goCtx, msgCreate)
		require.NoError(t, err)

		msgOnline := types.NewMsgSetOnline(vals[1])
		_, err = msgsrv.SetOnline(goCtx, msgOnline)
		require.NoError(t, err)

		//set legacy address for second validator
		v, found := valK.GetValidator(ctx, vals[1])
		require.True(t, found)
		v.RewardAddress = legacyRewardAddress
		valK.SetValidator(ctx, v)
	}

	// check validators created and online
	{
		updates := valK.BlockValidatorUpdates(ctx)
		// new validator is not online, there is not changes in tendermint validators and powers
		require.Len(t, updates, 2)
		require.Equal(t, genesisVal.ConsensusPower()+keeper.TokensToConsensusPower(firstStake.Amount)+keeper.TokensToConsensusPower(secondStake.Amount), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
		require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 3)
	}

	// ---------------------------- PAY VALIDATORS

	feeCollector := dsc.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	validatorPool := dsc.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)

	validatorPoolBalance := dsc.BankKeeper.GetAllBalances(ctx, validatorPool.GetAddress())
	feeCollectorBalance := dsc.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	valRs1, _ := valK.GetValidatorRS(ctx, genesisVal.GetOperator())
	valRs2, _ := valK.GetValidatorRS(ctx, genesisVal.GetOperator())
	valRs3, _ := valK.GetValidatorRS(ctx, genesisVal.GetOperator())
	// first pay validators
	{
		height := ctx.BlockHeight()
		// add constant rewards
		rewards := types.GetRewardForBlock(uint64(height))
		// send custom coin for distribute
		amount1 := initVolume1.Quo(sdk.NewInt(10))
		coin1 := sdk.NewCoin(ccDenom1, amount1)

		err = dsc.BankKeeper.SendCoinsFromAccountToModule(ctx, accs[0], authtypes.FeeCollectorName, sdk.NewCoins(coin1))
		require.NoError(t, err)

		baseAmount1 := formulas.CalculateSaleReturn(initVolume1, initReserve1, uint(crr), amount1)
		rewards = rewards.Add(baseAmount1)
		initVolume1 = initVolume1.Sub(amount1)
		initReserve1 = initReserve1.Sub(baseAmount1)

		// pay rewards
		valK.PayValidators(ctx)

		// check custom coin was updated
		updatedCCoin, err := dsc.CoinKeeper.GetCoin(ctx, ccDenom1)
		require.NoError(t, err)
		require.Equal(t, initVolume1, updatedCCoin.Volume)
		require.Equal(t, initReserve1, updatedCCoin.Reserve)

		// check correct rewards distribute
		totalPower := dsc.ValidatorKeeper.GetLastTotalPower(ctx)
		remainder := rewards
		{
			beforeRewards := valRs1.Rewards
			beforeTotalRewards := valRs1.TotalRewards
			valRs1, err = valK.GetValidatorRS(ctx, genesisVal.GetOperator())
			require.NoError(t, err)

			r := sdk.ZeroInt()
			r = rewards.Mul(sdk.NewInt(valRs1.Stake)).Quo(totalPower) //Quo(keeper.TokensFromConsensusPower(totalPower))
			require.Equal(t, valRs1.Rewards, beforeRewards.Add(r))
			require.Equal(t, valRs1.TotalRewards, beforeTotalRewards.Add(r))
			remainder = remainder.Sub(r)
		}

		{
			beforeRewards := valRs2.Rewards
			beforeTotalRewards := valRs2.TotalRewards
			valRs2, err = valK.GetValidatorRS(ctx, vals[0])
			require.NoError(t, err)

			r := sdk.ZeroInt()
			r = rewards.Mul(sdk.NewInt(keeper.TokensToConsensusPower(firstStake.Amount))).Quo(totalPower)
			require.Equal(t, valRs2.Rewards, beforeRewards.Add(r))
			require.Equal(t, valRs2.TotalRewards, beforeTotalRewards.Add(r))
			remainder = remainder.Sub(r)
		}

		{
			beforeRewards := valRs3.Rewards
			beforeTotalRewards := valRs3.TotalRewards
			valRs3, err = valK.GetValidatorRS(ctx, vals[1])
			require.NoError(t, err)

			r := sdk.ZeroInt()
			r = rewards.Mul(sdk.NewInt(keeper.TokensToConsensusPower(secondStake.Amount))).Quo(totalPower)
			require.Equal(t, valRs3.Rewards, beforeRewards.Add(r))
			require.Equal(t, valRs3.TotalRewards, beforeTotalRewards.Add(r))
			remainder = remainder.Sub(r)
		}
		// check validator pool for delegator distribute
		validatorPoolBalance = dsc.BankKeeper.GetAllBalances(ctx, validatorPool.GetAddress())
		require.Equal(t, sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, rewards.Sub(remainder))), validatorPoolBalance)

		// check feeCollectorPool
		feeCollectorBalance = dsc.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
		require.Equal(t, sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, remainder)), feeCollectorBalance)
	}
	// second pay validators
	{
		height := ctx.BlockHeight()
		// add constant rewards
		rewards := types.GetRewardForBlock(uint64(height))
		// add remainder from last distribute
		feeCollectorBase := dsc.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()).AmountOf(cmdcfg.BaseDenom)
		rewards = rewards.Add(feeCollectorBase)
		// send custom coin for distribute
		amount2 := initVolume2.Quo(sdk.NewInt(10))
		coin2 := sdk.NewCoin(ccDenom2, amount2)
		err = dsc.BankKeeper.SendCoinsFromAccountToModule(ctx, accs[0], authtypes.FeeCollectorName, sdk.NewCoins(coin2))
		require.NoError(t, err)
		baseAmount2 := formulas.CalculateSaleReturn(initVolume2, initReserve2, uint(crr), amount2)
		rewards = rewards.Add(baseAmount2)
		initVolume2 = initVolume2.Sub(amount2)
		initReserve2 = initReserve2.Sub(baseAmount2)

		// pay rewards
		beforeValidatorBalance := validatorPoolBalance
		valK.PayValidators(ctx)

		// check custom coin was updated
		updatedCCoin, err := dsc.CoinKeeper.GetCoin(ctx, ccDenom2)
		require.NoError(t, err)
		require.Equal(t, initVolume2, updatedCCoin.Volume)
		require.Equal(t, initReserve2, updatedCCoin.Reserve)

		// check correct rewards distribute
		totalPower := dsc.ValidatorKeeper.GetLastTotalPower(ctx)
		remainder := rewards
		{
			beforeRewards := valRs1.Rewards
			beforeTotalRewards := valRs1.TotalRewards
			valRs1, err = valK.GetValidatorRS(ctx, genesisVal.GetOperator())
			require.NoError(t, err)
			r := sdk.ZeroInt()
			r = rewards.Mul(sdk.NewInt(valRs1.Stake)).Quo(totalPower) //Quo(keeper.TokensFromConsensusPower(totalPower))
			require.Equal(t, beforeRewards.Add(r), valRs1.Rewards)
			require.Equal(t, beforeTotalRewards.Add(r), valRs1.TotalRewards)
			remainder = remainder.Sub(r)
		}

		{
			beforeRewards := valRs2.Rewards
			beforeTotalRewards := valRs2.TotalRewards
			valRs2, err = valK.GetValidatorRS(ctx, vals[0])
			require.NoError(t, err)

			r := sdk.ZeroInt()
			r = rewards.Mul(sdk.NewInt(keeper.TokensToConsensusPower(firstStake.Amount))).Quo(totalPower)
			require.Equal(t, valRs2.Rewards, beforeRewards.Add(r))
			require.Equal(t, valRs2.TotalRewards, beforeTotalRewards.Add(r))
			remainder = remainder.Sub(r)
		}

		{
			beforeRewards := valRs3.Rewards
			beforeTotalRewards := valRs3.TotalRewards
			valRs3, err = valK.GetValidatorRS(ctx, vals[1])
			require.NoError(t, err)

			r := sdk.ZeroInt()
			r = rewards.Mul(sdk.NewInt(keeper.TokensToConsensusPower(secondStake.Amount))).Quo(totalPower)
			require.Equal(t, valRs3.Rewards, beforeRewards.Add(r))
			require.Equal(t, valRs3.TotalRewards, beforeTotalRewards.Add(r))
			remainder = remainder.Sub(r)
		}
		// check validator pool for delegator distribute
		validatorPoolBalance = dsc.BankKeeper.GetAllBalances(ctx, validatorPool.GetAddress())
		require.Equal(t, beforeValidatorBalance.Add(sdk.NewCoin(cmdcfg.BaseDenom, rewards.Sub(remainder))), validatorPoolBalance)

		// check feeCollectorPool
		feeCollectorBalance = dsc.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
		require.Equal(t, sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, remainder)), feeCollectorBalance)
	}

	// ---------------------------- PAY DELEGATORS
	val2, _ := valK.GetValidator(ctx, vals[0])
	val3, _ := valK.GetValidator(ctx, vals[1])
	// send custom coins for delegators
	_, err = dsc.CoinKeeper.SendCoin(ctx, cointypes.NewMsgSendCoin(accs[0], accs[3], sdk.NewCoin(ccDenom1, initVolume1.Quo(sdk.NewInt(2)))))
	require.NoError(t, err)
	_, err = dsc.CoinKeeper.SendCoin(ctx, cointypes.NewMsgSendCoin(accs[0], accs[4], sdk.NewCoin(ccDenom2, initVolume2.Quo(sdk.NewInt(2)))))
	require.NoError(t, err)
	// delegates
	stake32 := types.NewStakeCoin(sdk.NewCoin(ccDenom1, helpers.EtherToWei(sdkmath.NewInt(1000))))
	stake42 := types.NewStakeCoin(sdk.NewCoin(ccDenom2, helpers.EtherToWei(sdkmath.NewInt(1200))))
	stake33 := types.NewStakeCoin(sdk.NewCoin(ccDenom1, helpers.EtherToWei(sdkmath.NewInt(1000))))
	stake34 := types.NewStakeCoin(sdk.NewCoin(ccDenom2, helpers.EtherToWei(sdkmath.NewInt(1200))))
	{
		err = valK.Delegate(ctx, accs[3], val2, stake32)
		require.NoError(t, err)
		err = valK.Delegate(ctx, accs[4], val2, stake42)
		require.NoError(t, err)
		err = valK.Delegate(ctx, accs[3], val3, stake33)
		require.NoError(t, err)
		err = valK.Delegate(ctx, accs[4], val3, stake34)
		require.NoError(t, err)
	}
	// clear delegators balances
	{
		balances := dsc.BankKeeper.GetAllBalances(ctx, accs[3])
		for _, v := range balances {
			_, err = dsc.CoinKeeper.SendCoin(ctx, cointypes.NewMsgSendCoin(accs[3], accs[6], v))
			require.NoError(t, err)
		}
		balances = dsc.BankKeeper.GetAllBalances(ctx, accs[4])
		for _, v := range balances {
			_, err = dsc.CoinKeeper.SendCoin(ctx, cointypes.NewMsgSendCoin(accs[4], accs[6], v))
			require.NoError(t, err)
		}
	}
	// check pay
	customCoinStaked := valK.GetAllCustomCoinsStaked(ctx)
	require.Equal(t, helpers.EtherToWei(sdkmath.NewInt(2000)), customCoinStaked[ccDenom1])
	require.Equal(t, helpers.EtherToWei(sdkmath.NewInt(2400)), customCoinStaked[ccDenom2])

	customCoinPrices := valK.CalculateCustomCoinPrices(ctx, customCoinStaked)
	require.Equal(t, sdk.NewDecFromInt(formulas.CalculateSaleReturn(initVolume1, initReserve1, uint(crr), helpers.EtherToWei(sdkmath.NewInt(2000)))).Quo(sdk.NewDecFromInt(helpers.EtherToWei(sdk.NewInt(2000)))), customCoinPrices[ccDenom1])
	require.Equal(t, sdk.NewDecFromInt(formulas.CalculateSaleReturn(initVolume2, initReserve2, uint(crr), helpers.EtherToWei(sdkmath.NewInt(2400)))).Quo(sdk.NewDecFromInt(helpers.EtherToWei(sdk.NewInt(2400)))), customCoinPrices[ccDenom2])

	totalComission := sdk.ZeroDec().Add(keeper.DAOCommission).Add(keeper.DevelopCommission)
	acc3Stake := sdk.NewDecFromInt(helpers.EtherToWei(sdkmath.NewInt(1000))).Mul(customCoinPrices[ccDenom1]).TruncateInt()
	acc4Stake := sdk.NewDecFromInt(helpers.EtherToWei(sdkmath.NewInt(1200))).Mul(customCoinPrices[ccDenom2]).TruncateInt()

	// calculate expected rewards
	delByValidator := valK.GetAllDelegationsByValidator(ctx)
	acc3TotalBalance := sdk.NewCoins()
	acc4TotalBalance := sdk.NewCoins()
	{
		totalStake, err := valK.CalculateTotalPowerWithDelegationsAndPrices(ctx, val2.GetOperator(), delByValidator[val2.GetOperator().String()], customCoinPrices)
		require.NoError(t, err)
		comission := sdk.NewDecFromInt(val2.Rewards).Mul(totalComission.Add(val2.Commission)).TruncateInt()
		rewards := val2.Rewards.Sub(comission)

		acc3TotalBalance = acc3TotalBalance.Add(sdk.NewCoin(cmdcfg.BaseDenom, rewards.Mul(acc3Stake).Quo(totalStake)))
		acc4TotalBalance = acc4TotalBalance.Add(sdk.NewCoin(cmdcfg.BaseDenom, rewards.Mul(acc4Stake).Quo(totalStake)))
	}
	{
		totalStake, err := valK.CalculateTotalPowerWithDelegationsAndPrices(ctx, val3.GetOperator(), delByValidator[val3.GetOperator().String()], customCoinPrices)
		require.NoError(t, err)
		comission := sdk.NewDecFromInt(val3.Rewards).Mul(totalComission.Add(val3.Commission)).TruncateInt()
		rewards := val3.Rewards.Sub(comission)

		acc3TotalBalance = acc3TotalBalance.Add(sdk.NewCoin(cmdcfg.BaseDenom, rewards.Mul(acc3Stake).Quo(totalStake)))
		acc4TotalBalance = acc4TotalBalance.Add(sdk.NewCoin(cmdcfg.BaseDenom, rewards.Mul(acc4Stake).Quo(totalStake)))
	}

	// pay rewards
	err = valK.PayRewards(ctx)
	require.NoError(t, err)

	// check balances after pay
	balances := dsc.BankKeeper.GetAllBalances(ctx, accs[3])
	require.Equal(t, acc3TotalBalance, balances)
	balances = dsc.BankKeeper.GetAllBalances(ctx, accs[4])
	require.Equal(t, acc4TotalBalance, balances)
}
