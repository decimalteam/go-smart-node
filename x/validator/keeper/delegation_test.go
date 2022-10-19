package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	sdkmath "cosmossdk.io/math"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// tests GetDelegation, GetDelegatorDelegations, SetDelegation, RemoveDelegation, GetDelegatorDelegations
func TestDelegation(t *testing.T) {
	var err error
	_, dsc, ctx := createTestInput(t)

	valK := dsc.ValidatorKeeper

	// remove genesis validator delegations
	delegations := valK.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)

	valK.RemoveDelegation(ctx, delegations[0])

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 3, defaultCoins)
	valAddrs := app.ConvertAddrsToValAddrs(addrDels)

	// create custom coin
	ccDenom := "custom"
	initVolume := keeper.TokensFromConsensusPower(100000000000)
	initReserve := keeper.TokensFromConsensusPower(1000)
	limitVolume := keeper.TokensFromConsensusPower(100000000000000000)
	crr := uint64(50)

	_, err = dsc.CoinKeeper.CreateCoin(ctx, cointypes.NewMsgCreateCoin(addrDels[0], ccDenom, "d", crr, initVolume, initReserve, limitVolume, ""))
	require.NoError(t, err)
	// ----------------------------

	// create nfts
	nftDenom := "nft_denom"
	subTokenReserve := sdk.NewCoin(cmdcfg.BaseDenom, keeper.TokensFromConsensusPower(100))

	_, err = dsc.NFTKeeper.MintToken(ctx, nfttypes.NewMsgMintToken(addrDels[0], "collection", nftDenom, "uri", true, addrDels[0], 5, subTokenReserve))
	require.NoError(t, err)
	// ----------------------------

	// create nfts with custom coin reserve
	nftCCDenom := "nft_cc_denom"
	subTokenCCReserve := sdk.NewCoin(ccDenom, keeper.TokensFromConsensusPower(100000000))

	_, err = dsc.NFTKeeper.MintToken(ctx, nfttypes.NewMsgMintToken(addrDels[0], "collection", nftCCDenom, "uri2", true, addrDels[0], 5, subTokenCCReserve))
	require.NoError(t, err)
	// ----------------------------

	defaultStake := types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, keeper.TokensFromConsensusPower(2)))
	defaultNftStake := types.NewStakeNFT(nftDenom, []uint32{1, 2, 3}, subTokenReserve)
	ccStake := types.NewStakeCoin(sdk.NewCoin(ccDenom, keeper.TokensFromConsensusPower(400000000)))
	ccNftStake := types.NewStakeNFT(nftCCDenom, []uint32{1, 2, 3}, subTokenCCReserve)

	// construct the validators
	amts := []sdkmath.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(7)}
	var validators [3]types.Validator
	for i := range amts {
		validators[i], err = types.NewValidator(valAddrs[i], addrDels[i], PKs[i], types.Description{}, sdk.ZeroDec())
		require.NoError(t, err)
		validators[i].Status = types.BondStatus_Bonded
		valK.CreateValidator(ctx, validators[i])
	}

	valAddr := validators[0].GetOperator()

	// delegate validator base coin
	{
		val, _ := valK.GetValidator(ctx, valAddr)
		del := addrDels[0]

		err = valK.Delegate(ctx, del, val, defaultStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		// validator power updated
		require.Equal(t, keeper.TokensToConsensusPower(val.Stake.Add(defaultStake.Stake.Amount)), keeper.TokensToConsensusPower(rs.Stake))
	}

	// delegate validator nfts with base reserve
	{
		val, _ := valK.GetValidator(ctx, valAddr)
		del := addrDels[0]

		err = valK.Delegate(ctx, del, val, defaultNftStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		// validator power updated
		require.Equal(t, keeper.TokensToConsensusPower(val.Stake.Add(subTokenReserve.Amount.Mul(sdk.NewInt(3)))), keeper.TokensToConsensusPower(rs.Stake))
	}

	//custom coin delegate
	valK.SetCustomCoinStaked(ctx, ccDenom, sdk.ZeroInt())
	{
		val, _ := valK.GetValidator(ctx, valAddr)
		del := addrDels[0]
		expectedStakeInBaseCoins := formulas.CalculateSaleReturn(initVolume, initReserve, uint(crr), ccStake.Stake.Amount)

		err = valK.Delegate(ctx, del, val, ccStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		// validator power updated
		require.Equal(t, keeper.TokensToConsensusPower(val.Stake.Add(expectedStakeInBaseCoins)), keeper.TokensToConsensusPower(rs.Stake))
		ccStaked := valK.GetCustomCoinStaked(ctx, ccDenom)
		// custom coin staked updated
		require.True(t, ccStake.Stake.Amount.Equal(ccStaked))
	}

	//  delegate validator nfts with custom coin reserve
	{
		val, _ := valK.GetValidator(ctx, valAddr)
		del := addrDels[0]
		ccs := valK.GetCustomCoinStaked(ctx, ccDenom)
		totalReserve := subTokenCCReserve.Amount.Mul(sdk.NewInt(3))
		expectedStakeInBaseCoins := formulas.CalculateSaleReturn(initVolume, initReserve, uint(crr), totalReserve)

		err = valK.Delegate(ctx, del, val, ccNftStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		// validator power updated
		require.Equal(t, keeper.TokensToConsensusPower(val.Stake.Add(expectedStakeInBaseCoins)), keeper.TokensToConsensusPower(rs.Stake))
		ccStaked := valK.GetCustomCoinStaked(ctx, ccDenom)
		// custom coin staked updated
		require.True(t, ccs.Add(totalReserve).Equal(ccStaked))
	}

	// delegate again base coin
	{
		val, _ := valK.GetValidator(ctx, valAddr)
		del := addrDels[0]
		oldDelegation, found := valK.GetDelegation(ctx, del, val.GetOperator(), defaultStake.Stake.Denom)
		require.True(t, found)

		err = valK.Delegate(ctx, del, val, defaultStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		newDelegation, found := valK.GetDelegation(ctx, del, val.GetOperator(), defaultStake.Stake.Denom)
		require.True(t, found)

		oldDelegation.Stake, err = oldDelegation.Stake.Add(defaultStake)
		require.NoError(t, err)

		// new amount added to current delegation
		require.True(t, newDelegation.Equal(oldDelegation))
		// validator power updated
		require.Equal(t, keeper.TokensToConsensusPower(val.Stake.Add(defaultStake.Stake.Amount)), keeper.TokensToConsensusPower(rs.Stake))
	}

	// delegate again custom coin
	{
		val, _ := valK.GetValidator(ctx, valAddr)
		del := addrDels[0]
		oldCcs := valK.GetCustomCoinStaked(ctx, ccDenom)
		expectedStakeInBaseCoins := formulas.CalculateSaleReturn(initVolume, initReserve, uint(crr), ccStake.Stake.Amount)

		oldDelegation, found := valK.GetDelegation(ctx, del, val.GetOperator(), ccStake.Stake.Denom)
		require.True(t, found)

		err = valK.Delegate(ctx, del, val, ccStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		newDelegation, found := valK.GetDelegation(ctx, del, val.GetOperator(), ccStake.Stake.Denom)
		require.True(t, found)

		oldDelegation.Stake, err = oldDelegation.Stake.Add(ccStake)
		require.NoError(t, err)

		// new amount added to current delegation
		require.True(t, newDelegation.Equal(oldDelegation))
		// validator power updated
		require.Equal(t, keeper.TokensToConsensusPower(val.Stake.Add(expectedStakeInBaseCoins)), keeper.TokensToConsensusPower(rs.Stake))
		ccStaked := valK.GetCustomCoinStaked(ctx, ccDenom)
		// custom coin staked updated
		require.True(t, oldCcs.Add(ccStake.Stake.Amount).Equal(ccStaked))
	}

	// delegate again nfts with base coin reserve
	{
		stake := defaultNftStake
		stake.SubTokenIDs = []uint32{4, 5}
		val, _ := valK.GetValidator(ctx, valAddr)
		del := addrDels[0]
		oldDelegation, found := valK.GetDelegation(ctx, del, val.GetOperator(), stake.ID)
		require.True(t, found)

		err = valK.Delegate(ctx, del, val, stake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		newDelegation, found := valK.GetDelegation(ctx, del, val.GetOperator(), stake.ID)
		require.True(t, found)

		oldDelegation.Stake, err = oldDelegation.Stake.Add(stake)
		require.NoError(t, err)

		// new amount added to current delegation
		require.True(t, newDelegation.Equal(oldDelegation))
		// validator power updated
		require.Equal(t, keeper.TokensToConsensusPower(val.Stake.Add(subTokenReserve.Amount.Mul(sdk.NewInt(2)))), keeper.TokensToConsensusPower(rs.Stake))
	}

	// delegate again nfts with custom coin reserve
	{
		coin, _ := dsc.CoinKeeper.GetCoin(ctx, ccDenom)
		stake := ccNftStake
		stake.SubTokenIDs = []uint32{4, 5}
		val, _ := valK.GetValidator(ctx, valAddr)
		del := addrDels[0]
		ccs := valK.GetCustomCoinStaked(ctx, ccDenom)
		totalReserve := subTokenCCReserve.Amount.Mul(sdk.NewInt(2))
		expectedStakeInBaseCoins := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, uint(crr), totalReserve)

		oldDelegation, found := valK.GetDelegation(ctx, del, val.GetOperator(), stake.ID)
		require.True(t, found)

		err = valK.Delegate(ctx, del, val, stake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		newDelegation, found := valK.GetDelegation(ctx, del, val.GetOperator(), stake.ID)
		require.True(t, found)

		oldDelegation.Stake, err = oldDelegation.Stake.Add(stake)
		require.NoError(t, err)

		// new amount added to current delegation
		require.True(t, newDelegation.Equal(oldDelegation))

		// validator power updated
		require.Equal(t, keeper.TokensToConsensusPower(val.Stake.Add(expectedStakeInBaseCoins)), keeper.TokensToConsensusPower(rs.Stake))
		ccStaked := valK.GetCustomCoinStaked(ctx, ccDenom)
		// custom coin staked updated
		require.True(t, ccs.Add(totalReserve).Equal(ccStaked))
	}
	delegations = valK.GetDelegatorDelegations(ctx, addrDels[0], 5)
	require.Len(t, delegations, 4)
}

func TestDelegationsInStore(t *testing.T) {
	var err error
	_, dsc, ctx := createTestInput(t)

	valK := dsc.ValidatorKeeper

	// remove genesis validator delegations
	delegations := valK.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)

	valK.RemoveDelegation(ctx, delegations[0])

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 3, defaultCoins)
	valAddrs := app.ConvertAddrsToValAddrs(addrDels)

	amts := []sdkmath.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(7)}
	var validators [3]types.Validator
	for i := range amts {
		validators[i], err = types.NewValidator(valAddrs[i], addrDels[i], PKs[i], types.Description{}, sdk.ZeroDec())
		require.NoError(t, err)
		validators[i].Status = types.BondStatus_Bonded
		valK.CreateValidator(ctx, validators[i])
	}

	defaultStake := types.NewStakeCoin(sdk.NewInt64Coin(cmdcfg.BaseDenom, 9))
	// first add a validators[0] to delegate too
	bond1to1 := types.NewDelegation(addrDels[0], valAddrs[0], defaultStake)

	//check the empty keeper first
	_, found := dsc.ValidatorKeeper.GetDelegation(ctx, addrDels[0], valAddrs[0], bond1to1.Stake.ID)
	require.False(t, found)

	//set and retrieve a record
	dsc.ValidatorKeeper.SetDelegation(ctx, bond1to1)

	resBond, found := dsc.ValidatorKeeper.GetDelegation(ctx, addrDels[0], valAddrs[0], bond1to1.Stake.ID)
	require.True(t, found)
	require.Equal(t, bond1to1, resBond)

	// modify a records, save, and retrieve
	bond1to1.Stake.Stake.Amount = sdk.NewInt(99)
	dsc.ValidatorKeeper.SetDelegation(ctx, bond1to1)
	resBond, found = dsc.ValidatorKeeper.GetDelegation(ctx, addrDels[0], valAddrs[0], bond1to1.Stake.ID)
	require.True(t, found)
	require.Equal(t, bond1to1, resBond)

	// add some more records
	bond1to2 := types.NewDelegation(addrDels[0], valAddrs[1], defaultStake)
	bond1to3 := types.NewDelegation(addrDels[0], valAddrs[2], defaultStake)
	bond2to1 := types.NewDelegation(addrDels[1], valAddrs[0], defaultStake)
	bond2to2 := types.NewDelegation(addrDels[1], valAddrs[1], defaultStake)
	bond2to3 := types.NewDelegation(addrDels[1], valAddrs[2], defaultStake)
	dsc.ValidatorKeeper.SetDelegation(ctx, bond1to2)
	dsc.ValidatorKeeper.SetDelegation(ctx, bond1to3)
	dsc.ValidatorKeeper.SetDelegation(ctx, bond2to1)
	dsc.ValidatorKeeper.SetDelegation(ctx, bond2to2)
	dsc.ValidatorKeeper.SetDelegation(ctx, bond2to3)

	// test all bond retrieve capabilities
	resBonds := dsc.ValidatorKeeper.GetDelegatorDelegations(ctx, addrDels[0], 5)
	require.Equal(t, 3, len(resBonds))
	require.Equal(t, bond1to1, resBonds[0])
	require.Equal(t, bond1to2, resBonds[1])
	require.Equal(t, bond1to3, resBonds[2])
	resBonds = dsc.ValidatorKeeper.GetAllDelegatorDelegations(ctx, addrDels[0])
	require.Equal(t, 3, len(resBonds))
	resBonds = dsc.ValidatorKeeper.GetDelegatorDelegations(ctx, addrDels[0], 2)
	require.Equal(t, 2, len(resBonds))
	resBonds = dsc.ValidatorKeeper.GetDelegatorDelegations(ctx, addrDels[1], 5)
	require.Equal(t, 3, len(resBonds))
	require.Equal(t, bond2to1, resBonds[0])
	require.Equal(t, bond2to2, resBonds[1])
	require.Equal(t, bond2to3, resBonds[2])
	allBonds := dsc.ValidatorKeeper.GetAllDelegations(ctx)
	require.Equal(t, 6, len(allBonds))
	require.Equal(t, bond1to1, allBonds[0])
	require.Equal(t, bond1to2, allBonds[1])
	require.Equal(t, bond1to3, allBonds[2])
	require.Equal(t, bond2to1, allBonds[3])
	require.Equal(t, bond2to2, allBonds[4])
	require.Equal(t, bond2to3, allBonds[5])

	resVals := dsc.ValidatorKeeper.GetDelegatorValidators(ctx, addrDels[0], 3)
	require.Equal(t, 3, len(resVals))
	resVals = dsc.ValidatorKeeper.GetDelegatorValidators(ctx, addrDels[1], 4)
	require.Equal(t, 3, len(resVals))

	for i := 0; i < 3; i++ {
		resVal, err := dsc.ValidatorKeeper.GetDelegatorValidator(ctx, addrDels[0], valAddrs[i])
		require.Nil(t, err)
		require.Equal(t, valAddrs[i], resVal.GetOperator())

		resVal, err = dsc.ValidatorKeeper.GetDelegatorValidator(ctx, addrDels[1], valAddrs[i])
		require.Nil(t, err)
		require.Equal(t, valAddrs[i], resVal.GetOperator())

		resDels := dsc.ValidatorKeeper.GetValidatorDelegations(ctx, valAddrs[i])
		require.Len(t, resDels, 2)
	}

	// test total bonded for single delegator
	expBonded := bond1to1.Stake.Stake.Add(bond2to1.Stake.Stake).Add(bond1to3.Stake.Stake)
	resDelBond := dsc.ValidatorKeeper.GetDelegatorBonded(ctx, addrDels[0])
	require.Equal(t, expBonded.Amount, resDelBond)

	// delete a record
	dsc.ValidatorKeeper.RemoveDelegation(ctx, bond2to3)
	_, found = dsc.ValidatorKeeper.GetDelegation(ctx, addrDels[1], valAddrs[2], bond2to3.Stake.ID)
	require.False(t, found)
	resBonds = dsc.ValidatorKeeper.GetDelegatorDelegations(ctx, addrDels[1], 5)
	require.Equal(t, 2, len(resBonds))
	require.Equal(t, bond2to1, resBonds[0])
	require.Equal(t, bond2to2, resBonds[1])

	resBonds = dsc.ValidatorKeeper.GetAllDelegatorDelegations(ctx, addrDels[1])
	require.Equal(t, 2, len(resBonds))

	// delete all the records from delegator 2
	dsc.ValidatorKeeper.RemoveDelegation(ctx, bond2to1)
	dsc.ValidatorKeeper.RemoveDelegation(ctx, bond2to2)
	_, found = dsc.ValidatorKeeper.GetDelegation(ctx, addrDels[1], valAddrs[0], bond2to1.Stake.ID)
	require.False(t, found)
	_, found = dsc.ValidatorKeeper.GetDelegation(ctx, addrDels[1], valAddrs[1], bond2to2.Stake.ID)
	require.False(t, found)
	resBonds = dsc.ValidatorKeeper.GetDelegatorDelegations(ctx, addrDels[1], 5)
	require.Equal(t, 0, len(resBonds))
}

// tests Get/Set/Remove UnbondingDelegation
func TestUnbondingDelegation(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	delAddrs := app.AddTestAddrsIncremental(dsc, ctx, 2, defaultCoins)
	valAddrs := app.ConvertAddrsToValAddrs(delAddrs)

	defaultStake := types.NewStakeCoin(sdk.NewInt64Coin(cmdcfg.BaseDenom, 9))

	ubd := types.NewUndelegation(
		delAddrs[0],
		valAddrs[0],
		0,
		time.Unix(0, 0).UTC(),
		defaultStake,
	)

	// set and retrieve a record
	dsc.ValidatorKeeper.SetUndelegation(ctx, ubd)
	resUnbond, found := dsc.ValidatorKeeper.GetUndelegation(ctx, delAddrs[0], valAddrs[0])
	require.True(t, found)
	require.Equal(t, ubd, resUnbond)

	// modify a records, save, and retrieve
	expUnbond := sdk.NewInt(21)
	ubd.Entries[0].Stake.Stake.Amount = expUnbond
	dsc.ValidatorKeeper.SetUndelegation(ctx, ubd)

	resUnbonds := dsc.ValidatorKeeper.GetUndelegations(ctx, delAddrs[0], 5)
	require.Equal(t, 1, len(resUnbonds))

	resUnbonds = dsc.ValidatorKeeper.GetAllUndelegations(ctx, delAddrs[0])
	require.Equal(t, 1, len(resUnbonds))

	resUnbond, found = dsc.ValidatorKeeper.GetUndelegation(ctx, delAddrs[0], valAddrs[0])
	require.True(t, found)
	require.Equal(t, ubd, resUnbond)

	resDelUnbond := dsc.ValidatorKeeper.GetDelegatorUnbonding(ctx, delAddrs[0])
	require.Equal(t, expUnbond, resDelUnbond)

	// delete a record
	dsc.ValidatorKeeper.RemoveUndelegation(ctx, ubd)
	_, found = dsc.ValidatorKeeper.GetUndelegation(ctx, delAddrs[0], valAddrs[0])
	require.False(t, found)

	resUnbonds = dsc.ValidatorKeeper.GetUndelegations(ctx, delAddrs[0], 5)
	require.Equal(t, 0, len(resUnbonds))

	resUnbonds = dsc.ValidatorKeeper.GetAllUndelegations(ctx, delAddrs[0])
	require.Equal(t, 0, len(resUnbonds))
}

// func TestUnbondDelegation(t *testing.T) {
// 	_, app, ctx := createTestInput(t)

// 	delAddrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.NewInt(10000))
// 	valAddrs := simapp.ConvertAddrsToValAddrs(delAddrs)

// 	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
// 	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)

// 	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens))))
// 	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

// 	// create a validator and a delegator to that validator
// 	// note this validator starts not-bonded
// 	validator := testvalidator.NewValidator(t, valAddrs[0], PKs[0])

// 	validator, issuedShares := validator.AddTokensFromDel(startTokens)
// 	require.Equal(t, startTokens, issuedShares.RoundInt())

// 	validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)

// 	delegation := types.NewDelegation(delAddrs[0], valAddrs[0], issuedShares)
// 	app.StakingKeeper.SetDelegation(ctx, delegation)

// 	bondTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 6)
// 	amount, err := app.StakingKeeper.Unbond(ctx, delAddrs[0], valAddrs[0], sdk.NewDecFromInt(bondTokens))
// 	require.NoError(t, err)
// 	require.Equal(t, bondTokens, amount) // shares to be added to an unbonding delegation

// 	delegation, found := app.StakingKeeper.GetDelegation(ctx, delAddrs[0], valAddrs[0])
// 	require.True(t, found)
// 	validator, found = app.StakingKeeper.GetValidator(ctx, valAddrs[0])
// 	require.True(t, found)

// 	remainingTokens := startTokens.Sub(bondTokens)
// 	require.Equal(t, remainingTokens, delegation.Shares.RoundInt())
// 	require.Equal(t, remainingTokens, validator.BondedTokens())
// }

//	func TestUnbondingDelegationsMaxEntries(t *testing.T) {
//		_, app, ctx := createTestInput(t)
//
//		addrDels := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.NewInt(10000))
//		addrVals := simapp.ConvertAddrsToValAddrs(addrDels)
//
//		startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//
//		bondDenom := app.StakingKeeper.BondDenom(ctx)
//		notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
//
//		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(bondDenom, startTokens))))
//		app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
//
//		// create a validator and a delegator to that validator
//		validator := testvalidator.NewValidator(t, addrVals[0], PKs[0])
//
//		validator, issuedShares := validator.AddTokensFromDel(startTokens)
//		require.Equal(t, startTokens, issuedShares.RoundInt())
//
//		validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//		require.True(sdk.IntEq(t, startTokens, validator.BondedTokens()))
//		require.True(t, validator.IsBonded())
//
//		delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
//		app.StakingKeeper.SetDelegation(ctx, delegation)
//
//		maxEntries := app.StakingKeeper.MaxEntries(ctx)
//
//		oldBonded := app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
//		oldNotBonded := app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//
//		// should all pass
//		var completionTime time.Time
//		for i := uint32(0); i < maxEntries; i++ {
//			var err error
//			completionTime, err = app.StakingKeeper.Undelegate(ctx, addrDels[0], addrVals[0], sdk.NewDec(1))
//			require.NoError(t, err)
//		}
//
//		newBonded := app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
//		newNotBonded := app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//		require.True(sdk.IntEq(t, newBonded, oldBonded.SubRaw(int64(maxEntries))))
//		require.True(sdk.IntEq(t, newNotBonded, oldNotBonded.AddRaw(int64(maxEntries))))
//
//		oldBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
//		oldNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//
//		// an additional unbond should fail due to max entries
//		_, err := app.StakingKeeper.Undelegate(ctx, addrDels[0], addrVals[0], sdk.NewDec(1))
//		require.Error(t, err)
//
//		newBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
//		newNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//
//		require.True(sdk.IntEq(t, newBonded, oldBonded))
//		require.True(sdk.IntEq(t, newNotBonded, oldNotBonded))
//
//		// mature unbonding delegations
//		ctx = ctx.WithBlockTime(completionTime)
//		_, err = app.StakingKeeper.CompleteUnbonding(ctx, addrDels[0], addrVals[0])
//		require.NoError(t, err)
//
//		newBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
//		newNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//		require.True(sdk.IntEq(t, newBonded, oldBonded))
//		require.True(sdk.IntEq(t, newNotBonded, oldNotBonded.SubRaw(int64(maxEntries))))
//
//		oldNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//
//		// unbonding  should work again
//		_, err = app.StakingKeeper.Undelegate(ctx, addrDels[0], addrVals[0], sdk.NewDec(1))
//		require.NoError(t, err)
//
//		newBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
//		newNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//		require.True(sdk.IntEq(t, newBonded, oldBonded.SubRaw(1)))
//		require.True(sdk.IntEq(t, newNotBonded, oldNotBonded.AddRaw(1)))
//	}
//
// // // test undelegating self delegation from a validator pushing it below MinSelfDelegation
// // // shift it from the bonded to unbonding state and jailed
//
//	func TestUndelegateSelfDelegationBelowMinSelfDelegation(t *testing.T) {
//		_, app, ctx := createTestInput(t)
//
//		addrDels := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.NewInt(10000))
//		addrVals := simapp.ConvertAddrsToValAddrs(addrDels)
//		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//		delCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), delTokens))
//
//		// create a validator with a self-delegation
//		validator := testvalidator.NewValidator(t, addrVals[0], PKs[0])
//
//		validator.MinSelfDelegation = delTokens
//		validator, issuedShares := validator.AddTokensFromDel(delTokens)
//		require.Equal(t, delTokens, issuedShares.RoundInt())
//
//		// add bonded tokens to pool for delegations
//		notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
//		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), delCoins))
//		app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
//
//		validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//		app.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
//		require.True(t, validator.IsBonded())
//
//		selfDelegation := types.NewDelegation(sdk.AccAddress(addrVals[0].Bytes()), addrVals[0], issuedShares)
//		app.StakingKeeper.SetDelegation(ctx, selfDelegation)
//
//		// add bonded tokens to pool for delegations
//		bondedPool := app.StakingKeeper.GetBondedPool(ctx)
//		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
//		app.AccountKeeper.SetModuleAccount(ctx, bondedPool)
//
//		// create a second delegation to this validator
//		app.StakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
//		validator, issuedShares = validator.AddTokensFromDel(delTokens)
//		require.True(t, validator.IsBonded())
//		require.Equal(t, delTokens, issuedShares.RoundInt())
//
//		// add bonded tokens to pool for delegations
//		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
//		app.AccountKeeper.SetModuleAccount(ctx, bondedPool)
//
//		validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//		delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
//		app.StakingKeeper.SetDelegation(ctx, delegation)
//
//		val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//		_, err := app.StakingKeeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewDecFromInt(app.StakingKeeper.TokensFromConsensusPower(ctx, 6)))
//		require.NoError(t, err)
//
//		// end block
//		applyValidatorSetUpdates(t, ctx, app.StakingKeeper, 1)
//
//		validator, found := app.StakingKeeper.GetValidator(ctx, addrVals[0])
//		require.True(t, found)
//		require.Equal(t, app.StakingKeeper.TokensFromConsensusPower(ctx, 14), validator.Tokens)
//		require.Equal(t, types.Unbonding, validator.Status)
//		require.True(t, validator.Jailed)
//	}
//
//	func TestUndelegateFromUnbondingValidator(t *testing.T) {
//		_, app, ctx := createTestInput(t)
//		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//		delCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), delTokens))
//
//		addrDels := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.NewInt(0))
//		addrVals := simapp.ConvertAddrsToValAddrs(addrDels)
//
//		// create a validator with a self-delegation
//		validator := testvalidator.NewValidator(t, addrVals[0], PKs[0])
//		app.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
//
//		validator, issuedShares := validator.AddTokensFromDel(delTokens)
//		require.Equal(t, delTokens, issuedShares.RoundInt())
//
//		// add bonded tokens to pool for delegations
//		notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
//		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), delCoins))
//		app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
//
//		validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//		require.True(t, validator.IsBonded())
//
//		selfDelegation := types.NewDelegation(addrVals[0].Bytes(), addrVals[0], issuedShares)
//		app.StakingKeeper.SetDelegation(ctx, selfDelegation)
//
//		// add bonded tokens to pool for delegations
//		bondedPool := app.StakingKeeper.GetBondedPool(ctx)
//		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
//		app.AccountKeeper.SetModuleAccount(ctx, bondedPool)
//
//		// create a second delegation to this validator
//		app.StakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
//
//		validator, issuedShares = validator.AddTokensFromDel(delTokens)
//		require.Equal(t, delTokens, issuedShares.RoundInt())
//
//		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
//		app.AccountKeeper.SetModuleAccount(ctx, bondedPool)
//
//		validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//		delegation := types.NewDelegation(addrDels[1], addrVals[0], issuedShares)
//		app.StakingKeeper.SetDelegation(ctx, delegation)
//
//		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
//		app.AccountKeeper.SetModuleAccount(ctx, bondedPool)
//
//		header := ctx.BlockHeader()
//		blockHeight := int64(10)
//		header.Height = blockHeight
//		blockTime := time.Unix(333, 0)
//		header.Time = blockTime
//		ctx = ctx.WithBlockHeader(header)
//
//		// unbond the all self-delegation to put validator in unbonding state
//		val0AccAddr := sdk.AccAddress(addrVals[0])
//		_, err := app.StakingKeeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewDecFromInt(delTokens))
//		require.NoError(t, err)
//
//		// end block
//		applyValidatorSetUpdates(t, ctx, app.StakingKeeper, 1)
//
//		validator, found := app.StakingKeeper.GetValidator(ctx, addrVals[0])
//		require.True(t, found)
//		require.Equal(t, blockHeight, validator.UnbondingHeight)
//		params := app.StakingKeeper.GetParams(ctx)
//		require.True(t, blockTime.Add(params.UnbondingTime).Equal(validator.UnbondingTime))
//
//		blockHeight2 := int64(20)
//		blockTime2 := time.Unix(444, 0).UTC()
//		ctx = ctx.WithBlockHeight(blockHeight2)
//		ctx = ctx.WithBlockTime(blockTime2)
//
//		// unbond some of the other delegation's shares
//		_, err = app.StakingKeeper.Undelegate(ctx, addrDels[1], addrVals[0], sdk.NewDec(6))
//		require.NoError(t, err)
//
//		// retrieve the unbonding delegation
//		ubd, found := app.StakingKeeper.GetUnbondingDelegation(ctx, addrDels[1], addrVals[0])
//		require.True(t, found)
//		require.Len(t, ubd.Entries, 1)
//		require.True(t, ubd.Entries[0].Balance.Equal(sdk.NewInt(6)))
//		assert.Equal(t, blockHeight2, ubd.Entries[0].CreationHeight)
//		assert.True(t, blockTime2.Add(params.UnbondingTime).Equal(ubd.Entries[0].CompletionTime))
//	}
//
//	func TestUndelegateFromUnbondedValidator(t *testing.T) {
//		_, app, ctx := createTestInput(t)
//		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//		delCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), delTokens))
//
//		addrDels := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.NewInt(0))
//		addrVals := simapp.ConvertAddrsToValAddrs(addrDels)
//
//		// add bonded tokens to pool for delegations
//		notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
//		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), delCoins))
//		app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
//
//		// create a validator with a self-delegation
//		validator := testvalidator.NewValidator(t, addrVals[0], PKs[0])
//		app.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
//
//		valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//		validator, issuedShares := validator.AddTokensFromDel(valTokens)
//		require.Equal(t, valTokens, issuedShares.RoundInt())
//		validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//		require.True(t, validator.IsBonded())
//
//		val0AccAddr := sdk.AccAddress(addrVals[0])
//		selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//		app.StakingKeeper.SetDelegation(ctx, selfDelegation)
//
//		// add bonded tokens to pool for delegations
//		bondedPool := app.StakingKeeper.GetBondedPool(ctx)
//		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
//		app.AccountKeeper.SetModuleAccount(ctx, bondedPool)
//
//		// create a second delegation to this validator
//		app.StakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
//		validator, issuedShares = validator.AddTokensFromDel(delTokens)
//		require.Equal(t, delTokens, issuedShares.RoundInt())
//		validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//		require.True(t, validator.IsBonded())
//		delegation := types.NewDelegation(addrDels[1], addrVals[0], issuedShares)
//		app.StakingKeeper.SetDelegation(ctx, delegation)
//
//		ctx = ctx.WithBlockHeight(10)
//		ctx = ctx.WithBlockTime(time.Unix(333, 0))
//
//		// unbond the all self-delegation to put validator in unbonding state
//		_, err := app.StakingKeeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewDecFromInt(valTokens))
//		require.NoError(t, err)
//
//		// end block
//		applyValidatorSetUpdates(t, ctx, app.StakingKeeper, 1)
//
//		validator, found := app.StakingKeeper.GetValidator(ctx, addrVals[0])
//		require.True(t, found)
//		require.Equal(t, ctx.BlockHeight(), validator.UnbondingHeight)
//		params := app.StakingKeeper.GetParams(ctx)
//		require.True(t, ctx.BlockHeader().Time.Add(params.UnbondingTime).Equal(validator.UnbondingTime))
//
//		// unbond the validator
//		ctx = ctx.WithBlockTime(validator.UnbondingTime)
//		app.StakingKeeper.UnbondAllMatureValidators(ctx)
//
//		// Make sure validator is still in state because there is still an outstanding delegation
//		validator, found = app.StakingKeeper.GetValidator(ctx, addrVals[0])
//		require.True(t, found)
//		require.Equal(t, validator.Status, types.Unbonded)
//
//		// unbond some of the other delegation's shares
//		unbondTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 6)
//		_, err = app.StakingKeeper.Undelegate(ctx, addrDels[1], addrVals[0], sdk.NewDecFromInt(unbondTokens))
//		require.NoError(t, err)
//
//		// unbond rest of the other delegation's shares
//		remainingTokens := delTokens.Sub(unbondTokens)
//		_, err = app.StakingKeeper.Undelegate(ctx, addrDels[1], addrVals[0], sdk.NewDecFromInt(remainingTokens))
//		require.NoError(t, err)
//
//		//  now validator should be deleted from state
//		validator, found = app.StakingKeeper.GetValidator(ctx, addrVals[0])
//		require.False(t, found, "%v", validator)
//	}
//
//	func TestUnbondingAllDelegationFromValidator(t *testing.T) {
//		_, app, ctx := createTestInput(t)
//		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//		delCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), delTokens))
//
//		addrDels := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.NewInt(0))
//		addrVals := simapp.ConvertAddrsToValAddrs(addrDels)
//
//		// add bonded tokens to pool for delegations
//		notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
//		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), delCoins))
//		app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
//
//		// create a validator with a self-delegation
//		validator := testvalidator.NewValidator(t, addrVals[0], PKs[0])
//		app.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
//
//		valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//		validator, issuedShares := validator.AddTokensFromDel(valTokens)
//		require.Equal(t, valTokens, issuedShares.RoundInt())
//
//		validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//		require.True(t, validator.IsBonded())
//		val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//
//		selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//		app.StakingKeeper.SetDelegation(ctx, selfDelegation)
//
//		// create a second delegation to this validator
//		app.StakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
//		validator, issuedShares = validator.AddTokensFromDel(delTokens)
//		require.Equal(t, delTokens, issuedShares.RoundInt())
//
//		// add bonded tokens to pool for delegations
//		bondedPool := app.StakingKeeper.GetBondedPool(ctx)
//		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
//		app.AccountKeeper.SetModuleAccount(ctx, bondedPool)
//
//		validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//		require.True(t, validator.IsBonded())
//
//		delegation := types.NewDelegation(addrDels[1], addrVals[0], issuedShares)
//		app.StakingKeeper.SetDelegation(ctx, delegation)
//
//		ctx = ctx.WithBlockHeight(10)
//		ctx = ctx.WithBlockTime(time.Unix(333, 0))
//
//		// unbond the all self-delegation to put validator in unbonding state
//		_, err := app.StakingKeeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewDecFromInt(valTokens))
//		require.NoError(t, err)
//
//		// end block
//		applyValidatorSetUpdates(t, ctx, app.StakingKeeper, 1)
//
//		// unbond all the remaining delegation
//		_, err = app.StakingKeeper.Undelegate(ctx, addrDels[1], addrVals[0], sdk.NewDecFromInt(delTokens))
//		require.NoError(t, err)
//
//		// validator should still be in state and still be in unbonding state
//		validator, found := app.StakingKeeper.GetValidator(ctx, addrVals[0])
//		require.True(t, found)
//		require.Equal(t, validator.Status, types.Unbonding)
//
//		// unbond the validator
//		ctx = ctx.WithBlockTime(validator.UnbondingTime)
//		app.StakingKeeper.UnbondAllMatureValidators(ctx)
//
//		// validator should now be deleted from state
//		_, found = app.StakingKeeper.GetValidator(ctx, addrVals[0])
//		require.False(t, found)
//	}
//
// Make sure that that the retrieving the delegations doesn't affect the state
func TestGetRedelegationsFromSrcValidator(t *testing.T) {
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 2, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	defaultStake := types.NewStakeCoin(sdk.NewInt64Coin(cmdcfg.BaseDenom, 5))

	rd := types.NewRedelegation(
		addrDels[0],
		addrVals[0],
		addrVals[1],
		0,
		time.Unix(0, 0),
		defaultStake,
	)

	// set and retrieve a record
	dsc.ValidatorKeeper.SetRedelegation(ctx, rd)
	resBond, found := dsc.ValidatorKeeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
	require.True(t, found)

	// get the redelegations one time
	redelegations := dsc.ValidatorKeeper.GetRedelegationsFromSrcValidator(ctx, addrVals[0])
	require.Equal(t, 1, len(redelegations))
	require.Equal(t, redelegations[0], resBond)

	// get the redelegations a second time, should be exactly the same
	redelegations = dsc.ValidatorKeeper.GetRedelegationsFromSrcValidator(ctx, addrVals[0])
	require.Equal(t, 1, len(redelegations))
	require.Equal(t, redelegations[0], resBond)
}

// tests Get/Set/Remove/Has UnbondingDelegation
func TestRedelegation(t *testing.T) {
	cmpRedelegations := func(rd1, rd2 types.Redelegation) bool {
		if rd1.Delegator != rd2.Delegator {
			return false
		}
		if rd1.ValidatorSrc != rd2.ValidatorSrc {
			return false
		}
		if rd1.ValidatorDst != rd2.ValidatorDst {
			return false
		}
		if len(rd1.Entries) != len(rd2.Entries) {
			return false
		}
		for i := range rd1.Entries {
			if !rd1.Entries[i].CompletionTime.Equal(rd2.Entries[i].CompletionTime) {
				return false
			}
			if rd1.Entries[i].CreationHeight != rd2.Entries[i].CreationHeight {
				return false
			}
			if !reflect.DeepEqual(rd1.Entries[i].Stake, rd2.Entries[i].Stake) {
				return false
			}
		}
		return true
	}
	_, dsc, ctx := createTestInput(t)

	addrDels := app.AddTestAddrsIncremental(dsc, ctx, 2, defaultCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	defaultStake := types.NewStakeCoin(sdk.NewInt64Coin(cmdcfg.BaseDenom, 5))

	rd := types.NewRedelegation(
		addrDels[0],
		addrVals[0],
		addrVals[1],
		0,
		time.Unix(0, 0),
		defaultStake,
	)

	// test shouldn't have and redelegations
	has := dsc.ValidatorKeeper.HasReceivingRedelegation(ctx, addrDels[0], addrVals[1])
	require.False(t, has)

	// set and retrieve a record
	dsc.ValidatorKeeper.SetRedelegation(ctx, rd)
	resRed, found := dsc.ValidatorKeeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
	require.True(t, found)

	redelegations := dsc.ValidatorKeeper.GetRedelegationsFromSrcValidator(ctx, addrVals[0])
	require.Equal(t, 1, len(redelegations))
	require.True(t, cmpRedelegations(redelegations[0], resRed))

	redelegations = dsc.ValidatorKeeper.GetRedelegations(ctx, addrDels[0], 5)
	require.Equal(t, 1, len(redelegations))
	require.True(t, cmpRedelegations(redelegations[0], resRed))

	redelegations = dsc.ValidatorKeeper.GetAllRedelegations(ctx, addrDels[0], nil, nil)
	require.Equal(t, 1, len(redelegations))
	require.True(t, cmpRedelegations(redelegations[0], resRed))

	// check if has the redelegation
	has = dsc.ValidatorKeeper.HasReceivingRedelegation(ctx, addrDels[0], addrVals[1])
	require.True(t, has)

	// modify a records, save, and retrieve
	rd.Entries[0].Stake.Stake.Amount = sdk.NewInt(21)
	dsc.ValidatorKeeper.SetRedelegation(ctx, rd)

	resRed, found = dsc.ValidatorKeeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
	require.True(t, found)
	require.True(t, cmpRedelegations(rd, resRed))

	redelegations = dsc.ValidatorKeeper.GetRedelegationsFromSrcValidator(ctx, addrVals[0])
	require.Equal(t, 1, len(redelegations))
	require.True(t, cmpRedelegations(redelegations[0], resRed))

	redelegations = dsc.ValidatorKeeper.GetRedelegations(ctx, addrDels[0], 5)
	require.Equal(t, 1, len(redelegations))
	require.True(t, cmpRedelegations(redelegations[0], resRed))

	// delete a record
	dsc.ValidatorKeeper.RemoveRedelegation(ctx, rd)
	_, found = dsc.ValidatorKeeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
	require.False(t, found)

	redelegations = dsc.ValidatorKeeper.GetRedelegations(ctx, addrDels[0], 5)
	require.Equal(t, 0, len(redelegations))

	redelegations = dsc.ValidatorKeeper.GetAllRedelegations(ctx, addrDels[0], nil, nil)
	require.Equal(t, 0, len(redelegations))
}

//func TestRedelegateToSameValidator(t *testing.T) {
//	_, app, ctx := createTestInput(t)
//
//	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.NewInt(0))
//	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)
//
//	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//	startCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), valTokens))
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
//	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), startCoins))
//	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	// create a validator with a self-delegation
//	validator := testvalidator.NewValidator(t, addrVals[0], PKs[0])
//	validator, issuedShares := validator.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//	require.True(t, validator.IsBonded())
//
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//	app.StakingKeeper.SetDelegation(ctx, selfDelegation)
//
//	_, err := app.StakingKeeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[0], sdk.NewDec(5))
//	require.Error(t, err)
//}
//
//func TestRedelegationMaxEntries(t *testing.T) {
//	_, app, ctx := createTestInput(t)
//
//	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.NewInt(0))
//	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)
//
//	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 20)
//	startCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens))
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
//	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), startCoins))
//	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	// create a validator with a self-delegation
//	validator := testvalidator.NewValidator(t, addrVals[0], PKs[0])
//	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//	validator, issuedShares := validator.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//	app.StakingKeeper.SetDelegation(ctx, selfDelegation)
//
//	// create a second validator
//	validator2 := testvalidator.NewValidator(t, addrVals[1], PKs[1])
//	validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//
//	validator2 = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator2, true)
//	require.Equal(t, types.Bonded, validator2.Status)
//
//	maxEntries := app.StakingKeeper.MaxEntries(ctx)
//
//	// redelegations should pass
//	var completionTime time.Time
//	for i := uint32(0); i < maxEntries; i++ {
//		var err error
//		completionTime, err = app.StakingKeeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1], sdk.NewDec(1))
//		require.NoError(t, err)
//	}
//
//	// an additional redelegation should fail due to max entries
//	_, err := app.StakingKeeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1], sdk.NewDec(1))
//	require.Error(t, err)
//
//	// mature redelegations
//	ctx = ctx.WithBlockTime(completionTime)
//	_, err = app.StakingKeeper.CompleteRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1])
//	require.NoError(t, err)
//
//	// redelegation should work again
//	_, err = app.StakingKeeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1], sdk.NewDec(1))
//	require.NoError(t, err)
//}
//
//func TestRedelegateSelfDelegation(t *testing.T) {
//	_, app, ctx := createTestInput(t)
//
//	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.NewInt(0))
//	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)
//
//	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 30)
//	startCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens))
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
//	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), startCoins))
//	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	// create a validator with a self-delegation
//	validator := testvalidator.NewValidator(t, addrVals[0], PKs[0])
//	app.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
//
//	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//	validator, issuedShares := validator.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//
//	validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//
//	val0AccAddr := sdk.AccAddress(addrVals[0])
//	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//	app.StakingKeeper.SetDelegation(ctx, selfDelegation)
//
//	// create a second validator
//	validator2 := testvalidator.NewValidator(t, addrVals[1], PKs[1])
//	validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator2 = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator2, true)
//	require.Equal(t, types.Bonded, validator2.Status)
//
//	// create a second delegation to validator 1
//	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//	validator, issuedShares = validator.AddTokensFromDel(delTokens)
//	require.Equal(t, delTokens, issuedShares.RoundInt())
//	validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//
//	delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
//	app.StakingKeeper.SetDelegation(ctx, delegation)
//
//	_, err := app.StakingKeeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1], sdk.NewDecFromInt(delTokens))
//	require.NoError(t, err)
//
//	// end block
//	applyValidatorSetUpdates(t, ctx, app.StakingKeeper, 2)
//
//	validator, found := app.StakingKeeper.GetValidator(ctx, addrVals[0])
//	require.True(t, found)
//	require.Equal(t, valTokens, validator.Tokens)
//	require.Equal(t, types.Unbonding, validator.Status)
//}
//
//func TestRedelegateFromUnbondingValidator(t *testing.T) {
//	_, app, ctx := createTestInput(t)
//
//	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.NewInt(0))
//	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)
//
//	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 30)
//	startCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens))
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
//	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), startCoins))
//	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	// create a validator with a self-delegation
//	validator := testvalidator.NewValidator(t, addrVals[0], PKs[0])
//	app.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
//
//	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//	validator, issuedShares := validator.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//	app.StakingKeeper.SetDelegation(ctx, selfDelegation)
//
//	// create a second delegation to this validator
//	app.StakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
//	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//	validator, issuedShares = validator.AddTokensFromDel(delTokens)
//	require.Equal(t, delTokens, issuedShares.RoundInt())
//	validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//	delegation := types.NewDelegation(addrDels[1], addrVals[0], issuedShares)
//	app.StakingKeeper.SetDelegation(ctx, delegation)
//
//	// create a second validator
//	validator2 := testvalidator.NewValidator(t, addrVals[1], PKs[1])
//	validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator2 = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator2, true)
//
//	header := ctx.BlockHeader()
//	blockHeight := int64(10)
//	header.Height = blockHeight
//	blockTime := time.Unix(333, 0)
//	header.Time = blockTime
//	ctx = ctx.WithBlockHeader(header)
//
//	// unbond the all self-delegation to put validator in unbonding state
//	_, err := app.StakingKeeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewDecFromInt(delTokens))
//	require.NoError(t, err)
//
//	// end block
//	applyValidatorSetUpdates(t, ctx, app.StakingKeeper, 1)
//
//	validator, found := app.StakingKeeper.GetValidator(ctx, addrVals[0])
//	require.True(t, found)
//	require.Equal(t, blockHeight, validator.UnbondingHeight)
//	params := app.StakingKeeper.GetParams(ctx)
//	require.True(t, blockTime.Add(params.UnbondingTime).Equal(validator.UnbondingTime))
//
//	// change the context
//	header = ctx.BlockHeader()
//	blockHeight2 := int64(20)
//	header.Height = blockHeight2
//	blockTime2 := time.Unix(444, 0)
//	header.Time = blockTime2
//	ctx = ctx.WithBlockHeader(header)
//
//	// unbond some of the other delegation's shares
//	redelegateTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 6)
//	_, err = app.StakingKeeper.BeginRedelegation(ctx, addrDels[1], addrVals[0], addrVals[1], sdk.NewDecFromInt(redelegateTokens))
//	require.NoError(t, err)
//
//	// retrieve the unbonding delegation
//	ubd, found := app.StakingKeeper.GetRedelegation(ctx, addrDels[1], addrVals[0], addrVals[1])
//	require.True(t, found)
//	require.Len(t, ubd.Entries, 1)
//	assert.Equal(t, blockHeight, ubd.Entries[0].CreationHeight)
//	assert.True(t, blockTime.Add(params.UnbondingTime).Equal(ubd.Entries[0].CompletionTime))
//}
//
//func TestRedelegateFromUnbondedValidator(t *testing.T) {
//	_, app, ctx := createTestInput(t)
//
//	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.NewInt(0))
//	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)
//
//	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 30)
//	startCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens))
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
//	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), startCoins))
//	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	// create a validator with a self-delegation
//	validator := testvalidator.NewValidator(t, addrVals[0], PKs[0])
//	app.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
//
//	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//	validator, issuedShares := validator.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//	app.StakingKeeper.SetDelegation(ctx, selfDelegation)
//
//	// create a second delegation to this validator
//	app.StakingKeeper.DeleteValidatorByPowerIndex(ctx, validator)
//	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
//	validator, issuedShares = validator.AddTokensFromDel(delTokens)
//	require.Equal(t, delTokens, issuedShares.RoundInt())
//	validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//	delegation := types.NewDelegation(addrDels[1], addrVals[0], issuedShares)
//	app.StakingKeeper.SetDelegation(ctx, delegation)
//
//	// create a second validator
//	validator2 := testvalidator.NewValidator(t, addrVals[1], PKs[1])
//	validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator2 = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator2, true)
//	require.Equal(t, types.Bonded, validator2.Status)
//
//	ctx = ctx.WithBlockHeight(10)
//	ctx = ctx.WithBlockTime(time.Unix(333, 0))
//
//	// unbond the all self-delegation to put validator in unbonding state
//	_, err := app.StakingKeeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewDecFromInt(delTokens))
//	require.NoError(t, err)
//
//	// end block
//	applyValidatorSetUpdates(t, ctx, app.StakingKeeper, 1)
//
//	validator, found := app.StakingKeeper.GetValidator(ctx, addrVals[0])
//	require.True(t, found)
//	require.Equal(t, ctx.BlockHeight(), validator.UnbondingHeight)
//	params := app.StakingKeeper.GetParams(ctx)
//	require.True(t, ctx.BlockHeader().Time.Add(params.UnbondingTime).Equal(validator.UnbondingTime))
//
//	// unbond the validator
//	app.StakingKeeper.UnbondingToUnbonded(ctx, validator)
//
//	// redelegate some of the delegation's shares
//	redelegationTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 6)
//	_, err = app.StakingKeeper.BeginRedelegation(ctx, addrDels[1], addrVals[0], addrVals[1], sdk.NewDecFromInt(redelegationTokens))
//	require.NoError(t, err)
//
//	// no red should have been found
//	red, found := app.StakingKeeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
//	require.False(t, found, "%v", red)
//}
