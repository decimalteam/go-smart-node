package keeper_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
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

	_, err = dsc.CoinKeeper.CreateCoin(ctx, cointypes.NewMsgCreateCoin(addrDels[0], ccDenom, "d", crr, initVolume, initReserve, limitVolume, sdkmath.ZeroInt(), ""))
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
	reserveForStake := subTokenReserve
	reserveForStake.Amount = reserveForStake.Amount.Mul(sdk.NewInt(3))
	defaultNftStake := types.NewStakeNFT(nftDenom, []uint32{1, 2, 3}, reserveForStake)

	ccStake := types.NewStakeCoin(sdk.NewCoin(ccDenom, keeper.TokensFromConsensusPower(400000000)))
	reserveForStake = subTokenCCReserve
	reserveForStake.Amount = reserveForStake.Amount.Mul(sdk.NewInt(3))
	ccNftStake := types.NewStakeNFT(nftCCDenom, []uint32{1, 2, 3}, reserveForStake)

	// construct the validators
	amts := []sdkmath.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(7)}
	var validators [3]types.Validator
	for i := range amts {
		validators[i], err = types.NewValidator(valAddrs[i], addrDels[i], PKs[i], types.Description{}, sdk.ZeroDec())
		require.NoError(t, err)
		// Bonded must be Online
		validators[i].Status = types.BondStatus_Bonded
		validators[i].Online = true
		validators[i].Stake = 10
		valK.CreateValidator(ctx, validators[i])
	}

	valAddr := validators[0].GetOperator()

	// delegate validator base coin
	{
		val, found := valK.GetValidator(ctx, valAddr)
		require.True(t, found)
		del := addrDels[0]

		err = valK.Delegate(ctx, del, val, defaultStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		// validator power updated
		require.Equal(t, val.Stake+keeper.TokensToConsensusPower(defaultStake.Stake.Amount), rs.Stake)
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
		require.Equal(t, val.Stake+keeper.TokensToConsensusPower(subTokenReserve.Amount.Mul(sdk.NewInt(3))), rs.Stake)
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
		require.Equal(t, val.Stake+keeper.TokensToConsensusPower((expectedStakeInBaseCoins)), rs.Stake)
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
		require.Equal(t, val.Stake+keeper.TokensToConsensusPower(expectedStakeInBaseCoins), rs.Stake)
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
		require.Equal(t, val.Stake+keeper.TokensToConsensusPower(defaultStake.Stake.Amount), rs.Stake)
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
		require.Equal(t, val.Stake+keeper.TokensToConsensusPower(expectedStakeInBaseCoins), rs.Stake)
		ccStaked := valK.GetCustomCoinStaked(ctx, ccDenom)
		// custom coin staked updated
		require.True(t, oldCcs.Add(ccStake.Stake.Amount).Equal(ccStaked))
	}

	// delegate again nfts with base coin reserve
	{
		stake := defaultNftStake
		stake.SubTokenIDs = []uint32{4, 5}
		stake.Stake.Amount = subTokenReserve.Amount.Mul(sdk.NewInt(2))
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
		require.Equal(t, val.Stake+keeper.TokensToConsensusPower(subTokenReserve.Amount.Mul(sdk.NewInt(2))), rs.Stake)
	}

	// delegate again nfts with custom coin reserve
	{
		coin, _ := dsc.CoinKeeper.GetCoin(ctx, ccDenom)
		stake := ccNftStake
		stake.SubTokenIDs = []uint32{4, 5}
		stake.Stake.Amount = subTokenCCReserve.Amount.Mul(sdk.NewInt(2))
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
		require.Equal(t, val.Stake+keeper.TokensToConsensusPower(expectedStakeInBaseCoins), rs.Stake)
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

func TestUndelegation(t *testing.T) {
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

	_, err = dsc.CoinKeeper.CreateCoin(ctx, cointypes.NewMsgCreateCoin(addrDels[0], ccDenom, "d", crr, initVolume, initReserve, limitVolume, sdkmath.ZeroInt(), ""))
	require.NoError(t, err)
	// ----------------------------

	// create nfts
	nftDenom := "nft_denom"
	subTokenReserve := sdk.NewCoin(cmdcfg.BaseDenom, keeper.TokensFromConsensusPower(100))

	_, err = dsc.NFTKeeper.MintToken(ctx, nfttypes.NewMsgMintToken(addrDels[0], "collection", nftDenom, "uri", true, addrDels[0], 7, subTokenReserve))
	require.NoError(t, err)
	// ----------------------------

	// create nfts with custom coin reserve
	nftCCDenom := "nft_cc_denom"
	subTokenCCReserve := sdk.NewCoin(ccDenom, keeper.TokensFromConsensusPower(100000000))

	_, err = dsc.NFTKeeper.MintToken(ctx, nfttypes.NewMsgMintToken(addrDels[0], "collection", nftCCDenom, "uri2", true, addrDels[0], 7, subTokenCCReserve))
	require.NoError(t, err)
	// ----------------------------

	defaultStake := types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, keeper.TokensFromConsensusPower(2)))
	reserveForStake := subTokenReserve
	reserveForStake.Amount = reserveForStake.Amount.Mul(sdk.NewInt(3))
	defaultNftStake := types.NewStakeNFT(nftDenom, []uint32{1, 2, 3}, reserveForStake)

	ccStake := types.NewStakeCoin(sdk.NewCoin(ccDenom, keeper.TokensFromConsensusPower(400000000)))
	reserveForStake = subTokenCCReserve
	reserveForStake.Amount = subTokenCCReserve.Amount.Mul(sdk.NewInt(3))
	ccNftStake := types.NewStakeNFT(nftCCDenom, []uint32{1, 2, 3}, reserveForStake)

	// construct the validators
	amts := []sdkmath.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(7)}
	var validators [3]types.Validator
	for i := range amts {
		validators[i], err = types.NewValidator(valAddrs[i], addrDels[i], PKs[i], types.Description{}, sdk.ZeroDec())
		require.NoError(t, err)
		// Bonded = Online
		validators[i].Status = types.BondStatus_Bonded
		validators[i].Online = true
		valK.CreateValidator(ctx, validators[i])
	}

	// delegates
	delAddr := addrDels[0]
	valAddr := validators[0].GetOperator()

	err = valK.Delegate(ctx, delAddr, validators[0], defaultStake)
	require.NoError(t, err)
	err = valK.Delegate(ctx, delAddr, validators[0], ccStake)
	require.NoError(t, err)
	err = valK.Delegate(ctx, delAddr, validators[0], defaultNftStake)
	require.NoError(t, err)
	err = valK.Delegate(ctx, delAddr, validators[0], ccNftStake)
	require.NoError(t, err)

	coins := sdk.NewCoins()
	defaultNfts := make(map[uint32]bool)
	customNfts := make(map[uint32]bool)
	// undelegate base coin del
	{
		val, _ := valK.GetValidator(ctx, valAddr)
		unStake := defaultStake
		unStake.Stake.Amount = unStake.Stake.Amount.Quo(sdk.NewInt(2))
		remainStake, err := valK.CalculateRemainStake(ctx, defaultStake, unStake)
		require.NoError(t, err)

		_, err = valK.Undelegate(ctx, delAddr, valAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		delegation, found := valK.GetDelegation(ctx, delAddr, valAddr, unStake.GetID())
		require.True(t, found)

		// tokens are sub from delegation
		require.True(t, delegation.Stake.Equal(&remainStake))
		// validator power updated
		require.Equal(t, val.Stake-keeper.TokensToConsensusPower(unStake.Stake.Amount), rs.Stake)

		coins = coins.Add(unStake.Stake)
		defaultStake = remainStake
	}
	// undelegate custom coin del
	{
		coin, _ := dsc.CoinKeeper.GetCoin(ctx, ccDenom)
		val, _ := valK.GetValidator(ctx, valAddr)
		unStake := ccStake
		ccs := valK.GetCustomCoinStaked(ctx, ccDenom)
		unStake.Stake.Amount = unStake.Stake.Amount.Quo(sdk.NewInt(2))
		unStakeInBaseCoin := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, uint(crr), unStake.Stake.Amount)
		remainStake, err := valK.CalculateRemainStake(ctx, ccStake, unStake)
		require.NoError(t, err)

		_, err = valK.Undelegate(ctx, delAddr, valAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		delegation, found := valK.GetDelegation(ctx, delAddr, valAddr, unStake.GetID())
		require.True(t, found)

		// custom coin staked are not sub
		require.True(t, ccs.Equal(valK.GetCustomCoinStaked(ctx, ccDenom)))
		// tokens are sub from delegation
		require.True(t, delegation.Stake.Equal(&remainStake))
		// validator power updated
		require.Equal(t, val.Stake-keeper.TokensToConsensusPower(unStakeInBaseCoin), rs.Stake)

		coins = coins.Add(unStake.Stake)
		ccStake = remainStake
	}
	// undelegate nfts with base coin reserve
	{
		val, _ := valK.GetValidator(ctx, valAddr)
		unStake := defaultNftStake
		unStake.SubTokenIDs = []uint32{3}
		unStake.Stake.Amount = subTokenReserve.Amount.Mul(sdk.NewInt(1))
		remainStake, err := valK.CalculateRemainStake(ctx, defaultNftStake, unStake)
		require.NoError(t, err)
		reserveTotalStake := subTokenReserve.Amount.Mul(sdk.NewInt(1))

		_, err = valK.Undelegate(ctx, delAddr, valAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		delegation, found := valK.GetDelegation(ctx, delAddr, valAddr, unStake.GetID())
		require.True(t, found)

		// tokens are sub from delegation
		require.True(t, delegation.Stake.Equal(&remainStake))
		// validator power updated
		require.Equal(t, val.Stake-keeper.TokensToConsensusPower(reserveTotalStake), rs.Stake)

		for _, v := range unStake.SubTokenIDs {
			defaultNfts[v] = false
		}

		defaultNftStake = remainStake
	}
	// undelegate nfts with custom coin reserve
	{
		coin, _ := dsc.CoinKeeper.GetCoin(ctx, ccDenom)
		val, _ := valK.GetValidator(ctx, valAddr)
		unStake := ccNftStake
		unStake.SubTokenIDs = []uint32{2, 3}
		unStake.Stake.Amount = subTokenCCReserve.Amount.Mul(sdk.NewInt(2))
		ccs := valK.GetCustomCoinStaked(ctx, ccDenom)
		remainStake, err := valK.CalculateRemainStake(ctx, ccNftStake, unStake)
		require.NoError(t, err)
		reserveTotalStake := subTokenCCReserve.Amount.Mul(sdk.NewInt(2))
		unStakeInBaseCoin := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, uint(crr), reserveTotalStake)

		_, err = valK.Undelegate(ctx, delAddr, valAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		delegation, found := valK.GetDelegation(ctx, delAddr, valAddr, unStake.GetID())
		require.True(t, found)

		// custom coin staked are not sub
		require.True(t, ccs.Equal(valK.GetCustomCoinStaked(ctx, ccDenom)))
		// tokens are sub from delegation
		require.True(t, delegation.Stake.Equal(&remainStake))
		// validator power updated
		require.Equal(t, val.Stake-keeper.TokensToConsensusPower(unStakeInBaseCoin), rs.Stake)

		for _, v := range unStake.SubTokenIDs {
			customNfts[v] = false
		}
		ccNftStake = remainStake
	}
	// undelegate again base coin del
	{
		val, _ := valK.GetValidator(ctx, valAddr)
		unStake := defaultStake
		unStake.Stake.Amount = unStake.Stake.Amount.Quo(sdk.NewInt(2))
		remainStake, err := valK.CalculateRemainStake(ctx, defaultStake, unStake)
		require.NoError(t, err)

		_, err = valK.Undelegate(ctx, delAddr, valAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		delegation, found := valK.GetDelegation(ctx, delAddr, valAddr, unStake.GetID())
		require.True(t, found)

		// tokens are sub from delegation
		require.True(t, delegation.Stake.Equal(&remainStake))
		// validator power updated
		require.Equal(t, val.Stake-keeper.TokensToConsensusPower(unStake.Stake.Amount), rs.Stake)

		coins = coins.Add(unStake.Stake)
		defaultStake = remainStake
	}
	// undelegate again custom coin del
	{
		coin, _ := dsc.CoinKeeper.GetCoin(ctx, ccDenom)
		val, _ := valK.GetValidator(ctx, valAddr)
		unStake := ccStake
		ccs := valK.GetCustomCoinStaked(ctx, ccDenom)
		//unStake.Stake.Amount = unStake.Stake.Amount.Quo(sdk.NewInt())
		unStakeInBaseCoin := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, uint(crr), unStake.Stake.Amount)
		remainStake, err := valK.CalculateRemainStake(ctx, ccStake, unStake)
		require.NoError(t, err)

		_, err = valK.Undelegate(ctx, delAddr, valAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		_, found := valK.GetDelegation(ctx, delAddr, valAddr, unStake.GetID())
		require.False(t, found)

		// custom coin staked are not sub
		require.True(t, ccs.Equal(valK.GetCustomCoinStaked(ctx, ccDenom)))
		// validator power updated
		require.Equal(t, val.Stake-keeper.TokensToConsensusPower(unStakeInBaseCoin), rs.Stake)

		coins = coins.Add(unStake.Stake)
		ccStake = remainStake
	}
	// undelegate again nfts with base coin reserve
	{
		val, _ := valK.GetValidator(ctx, valAddr)
		unStake := defaultNftStake
		unStake.SubTokenIDs = []uint32{2}
		unStake.Stake.Amount = subTokenReserve.Amount.Mul(sdk.NewInt(1))
		remainStake, err := valK.CalculateRemainStake(ctx, defaultNftStake, unStake)
		require.NoError(t, err)
		reserveTotalStake := subTokenReserve.Amount.Mul(sdk.NewInt(1))

		_, err = valK.Undelegate(ctx, delAddr, valAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		delegation, found := valK.GetDelegation(ctx, delAddr, valAddr, unStake.GetID())
		require.True(t, found)

		// tokens are sub from delegation
		require.True(t, delegation.Stake.Equal(&remainStake))
		// validator power updated
		require.Equal(t, val.Stake-keeper.TokensToConsensusPower(reserveTotalStake), rs.Stake)

		for _, v := range unStake.SubTokenIDs {
			defaultNfts[v] = false
		}
	}

	balances := dsc.BankKeeper.GetAllBalances(ctx, valK.GetNotBondedPool(ctx).GetAddress())
	require.True(t, balances.IsEqual(coins))
	expectBalances := balances.Sub(coins...)

	blockHeader := ctx.BlockHeader()
	blockHeader.Time = blockHeader.Time.Add(time.Hour * 100000)
	ctx = ctx.WithBlockHeader(blockHeader)
	err = valK.CompleteUnbonding(ctx, delAddr, valAddr)
	require.NoError(t, err)

	balances = dsc.BankKeeper.GetAllBalances(ctx, valK.GetNotBondedPool(ctx).GetAddress())
	require.True(t, expectBalances.IsEqual(balances))

	for i := range defaultNfts {
		st, ok := dsc.NFTKeeper.GetSubToken(ctx, nftDenom, i)
		require.True(t, ok)
		require.Equal(t, st.Owner, delAddr.String())
	}
	for i := range customNfts {
		st, ok := dsc.NFTKeeper.GetSubToken(ctx, nftCCDenom, i)
		require.True(t, ok)
		require.Equal(t, st.Owner, delAddr.String())
	}
}

func TestRedelegation(t *testing.T) {
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

	_, err = dsc.CoinKeeper.CreateCoin(ctx, cointypes.NewMsgCreateCoin(addrDels[0], ccDenom, "d", crr, initVolume, initReserve, limitVolume, sdkmath.ZeroInt(), ""))
	require.NoError(t, err)
	// ----------------------------

	// create nfts
	nftDenom := "nft_denom"
	subTokenReserve := sdk.NewCoin(cmdcfg.BaseDenom, keeper.TokensFromConsensusPower(100))

	_, err = dsc.NFTKeeper.MintToken(ctx, nfttypes.NewMsgMintToken(addrDels[0], "collection", nftDenom, "uri", true, addrDels[0], 7, subTokenReserve))
	require.NoError(t, err)
	// ----------------------------

	// create nfts with custom coin reserve
	nftCCDenom := "nft_cc_denom"
	subTokenCCReserve := sdk.NewCoin(ccDenom, keeper.TokensFromConsensusPower(100000000))

	_, err = dsc.NFTKeeper.MintToken(ctx, nfttypes.NewMsgMintToken(addrDels[0], "collection", nftCCDenom, "uri2", true, addrDels[0], 7, subTokenCCReserve))
	require.NoError(t, err)
	// ----------------------------

	defaultStake := types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, keeper.TokensFromConsensusPower(2)))
	reserveForStake := subTokenReserve
	reserveForStake.Amount = reserveForStake.Amount.Mul(sdk.NewInt(3))
	defaultNftStake := types.NewStakeNFT(nftDenom, []uint32{1, 2, 3}, reserveForStake)

	ccStake := types.NewStakeCoin(sdk.NewCoin(ccDenom, keeper.TokensFromConsensusPower(400000000)))
	reserveForStake = subTokenCCReserve
	reserveForStake.Amount = reserveForStake.Amount.Mul(sdk.NewInt(3))
	ccNftStake := types.NewStakeNFT(nftCCDenom, []uint32{1, 2, 3}, reserveForStake)

	// construct the validators
	amts := []sdkmath.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(7)}
	var validators [3]types.Validator
	for i := range amts {
		validators[i], err = types.NewValidator(valAddrs[i], addrDels[i], PKs[i], types.Description{}, sdk.ZeroDec())
		require.NoError(t, err)
		// Bonded = Online
		validators[i].Status = types.BondStatus_Bonded
		validators[i].Online = true
		valK.CreateValidator(ctx, validators[i])
	}

	// delegates
	delAddr := addrDels[0]
	valSrcAddr := validators[0].GetOperator()
	valDstAddr := validators[1].GetOperator()

	err = valK.Delegate(ctx, delAddr, validators[0], defaultStake)
	require.NoError(t, err)
	err = valK.Delegate(ctx, delAddr, validators[0], ccStake)
	require.NoError(t, err)
	err = valK.Delegate(ctx, delAddr, validators[0], defaultNftStake)
	require.NoError(t, err)
	err = valK.Delegate(ctx, delAddr, validators[0], ccNftStake)
	require.NoError(t, err)

	expectDelegations := make(map[string]types.Delegation)

	// redelegate base coin
	{
		valSrc, _ := valK.GetValidator(ctx, valSrcAddr)
		//valDst, _ := valK.GetValidator(ctx, valDstAddr)
		unStake := defaultStake
		unStake.Stake.Amount = unStake.Stake.Amount.Quo(sdk.NewInt(2))
		remainStake, err := valK.CalculateRemainStake(ctx, defaultStake, unStake)
		require.NoError(t, err)

		_, err = valK.BeginRedelegation(ctx, delAddr, valSrcAddr, valDstAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, valSrcAddr)
		require.NoError(t, err)

		delegation, found := valK.GetDelegation(ctx, delAddr, valSrcAddr, unStake.GetID())
		require.True(t, found)

		// tokens are sub from delegation
		require.True(t, delegation.Stake.Equal(&remainStake))
		// validator power updated
		require.Equal(t, valSrc.Stake-keeper.TokensToConsensusPower(unStake.Stake.Amount), rs.Stake)

		expectDelegations[unStake.ID] = types.Delegation{
			Delegator: delAddr.String(),
			Validator: valDstAddr.String(),
			Stake:     unStake,
		}
		defaultStake = remainStake
	}
	// redelegate custom coin
	{
		coin, _ := dsc.CoinKeeper.GetCoin(ctx, ccDenom)
		valSrc, _ := valK.GetValidator(ctx, valSrcAddr)
		unStake := ccStake
		ccs := valK.GetCustomCoinStaked(ctx, ccDenom)
		unStake.Stake.Amount = unStake.Stake.Amount.Quo(sdk.NewInt(2))
		unStakeInBaseCoin := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, uint(crr), unStake.Stake.Amount)
		remainStake, err := valK.CalculateRemainStake(ctx, ccStake, unStake)
		require.NoError(t, err)

		_, err = valK.BeginRedelegation(ctx, delAddr, valSrcAddr, valDstAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, valSrcAddr)
		require.NoError(t, err)

		delegation, found := valK.GetDelegation(ctx, delAddr, valSrcAddr, unStake.GetID())
		require.True(t, found)

		// custom coin staked are not sub
		require.True(t, ccs.Equal(valK.GetCustomCoinStaked(ctx, ccDenom)))
		// tokens are sub from delegation
		require.True(t, delegation.Stake.Equal(&remainStake))
		// validator power updated
		require.Equal(t, valSrc.Stake-keeper.TokensToConsensusPower(unStakeInBaseCoin), rs.Stake)

		expectDelegations[unStake.ID] = types.Delegation{
			Delegator: delAddr.String(),
			Validator: valDstAddr.String(),
			Stake:     unStake,
		}
		ccStake = remainStake
	}
	// redelegate nfts with base coin reserve
	{
		valSrc, _ := valK.GetValidator(ctx, valSrcAddr)
		unStake := defaultNftStake
		unStake.SubTokenIDs = []uint32{3}
		unStake.Stake.Amount = subTokenReserve.Amount.Mul(sdk.NewInt(1))
		remainStake, err := valK.CalculateRemainStake(ctx, defaultNftStake, unStake)
		require.NoError(t, err)
		reserveTotalStake := subTokenReserve.Amount.Mul(sdk.NewInt(1))

		_, err = valK.BeginRedelegation(ctx, delAddr, valSrcAddr, valDstAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, valSrcAddr)
		require.NoError(t, err)

		delegation, found := valK.GetDelegation(ctx, delAddr, valSrcAddr, unStake.GetID())
		require.True(t, found)

		// tokens are sub from delegation
		require.True(t, delegation.Stake.Equal(&remainStake))
		// validator power updated
		require.Equal(t, valSrc.Stake-keeper.TokensToConsensusPower(reserveTotalStake), rs.Stake)

		expectDelegations[unStake.ID] = types.Delegation{
			Delegator: delAddr.String(),
			Validator: valDstAddr.String(),
			Stake:     unStake,
		}
		defaultNftStake = remainStake
	}
	// redelegate nfts with custom coin reserve
	{
		coin, _ := dsc.CoinKeeper.GetCoin(ctx, ccDenom)
		valSrc, _ := valK.GetValidator(ctx, valSrcAddr)
		unStake := ccNftStake
		unStake.SubTokenIDs = []uint32{2, 3}
		unStake.Stake.Amount = subTokenCCReserve.Amount.Mul(sdk.NewInt(2))
		ccs := valK.GetCustomCoinStaked(ctx, ccDenom)
		remainStake, err := valK.CalculateRemainStake(ctx, ccNftStake, unStake)
		require.NoError(t, err)
		reserveTotalStake := subTokenCCReserve.Amount.Mul(sdk.NewInt(2))
		unStakeInBaseCoin := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, uint(crr), reserveTotalStake)

		_, err = valK.BeginRedelegation(ctx, delAddr, valSrcAddr, valDstAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, valSrc.GetOperator())
		require.NoError(t, err)

		delegation, found := valK.GetDelegation(ctx, delAddr, valSrcAddr, unStake.GetID())
		require.True(t, found)

		// custom coin staked are not sub
		require.True(t, ccs.Equal(valK.GetCustomCoinStaked(ctx, ccDenom)))
		// tokens are sub from delegation
		require.True(t, delegation.Stake.Equal(&remainStake))
		// validator power updated
		require.Equal(t, valSrc.Stake-keeper.TokensToConsensusPower(unStakeInBaseCoin), rs.Stake)

		expectDelegations[unStake.ID] = types.Delegation{
			Delegator: delAddr.String(),
			Validator: valDstAddr.String(),
			Stake:     unStake,
		}
		ccNftStake = remainStake
	}
	// redelegate again base coin del
	{
		valSrc, _ := valK.GetValidator(ctx, valSrcAddr)
		unStake := defaultStake
		unStake.Stake.Amount = unStake.Stake.Amount.Quo(sdk.NewInt(2))
		remainStake, err := valK.CalculateRemainStake(ctx, defaultStake, unStake)
		require.NoError(t, err)

		_, err = valK.BeginRedelegation(ctx, delAddr, valSrcAddr, valDstAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, valSrcAddr)
		require.NoError(t, err)

		delegation, found := valK.GetDelegation(ctx, delAddr, valSrcAddr, unStake.GetID())
		require.True(t, found)

		// tokens are sub from delegation
		require.True(t, delegation.Stake.Equal(&remainStake))
		// validator power updated
		require.Equal(t, valSrc.Stake-keeper.TokensToConsensusPower(unStake.Stake.Amount), rs.Stake)

		oldExDel := expectDelegations[unStake.ID]
		oldExDel.Stake, err = oldExDel.Stake.Add(unStake)
		require.NoError(t, err)
		expectDelegations[unStake.ID] = oldExDel
		defaultStake = remainStake
	}
	// redelegate again custom coin del
	{
		coin, _ := dsc.CoinKeeper.GetCoin(ctx, ccDenom)
		valSrc, _ := valK.GetValidator(ctx, valSrcAddr)
		unStake := ccStake
		ccs := valK.GetCustomCoinStaked(ctx, ccDenom)
		//unStake.Stake.Amount = unStake.Stake.Amount.Quo(sdk.NewInt())
		unStakeInBaseCoin := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, uint(crr), unStake.Stake.Amount)
		remainStake, err := valK.CalculateRemainStake(ctx, ccStake, unStake)
		require.NoError(t, err)

		_, err = valK.BeginRedelegation(ctx, delAddr, valSrcAddr, valDstAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, valSrcAddr)
		require.NoError(t, err)

		_, found := valK.GetDelegation(ctx, delAddr, valSrcAddr, unStake.GetID())
		require.False(t, found)

		// custom coin staked are not sub
		require.True(t, ccs.Equal(valK.GetCustomCoinStaked(ctx, ccDenom)))
		// validator power updated
		require.Equal(t, valSrc.Stake-keeper.TokensToConsensusPower(unStakeInBaseCoin), rs.Stake)

		oldExDel := expectDelegations[unStake.ID]
		oldExDel.Stake, err = oldExDel.Stake.Add(unStake)
		require.NoError(t, err)
		expectDelegations[unStake.ID] = oldExDel
		ccStake = remainStake
	}
	// redelegate again nfts with base coin reserve
	{
		val, _ := valK.GetValidator(ctx, valSrcAddr)
		unStake := defaultNftStake
		unStake.SubTokenIDs = []uint32{2}
		unStake.Stake.Amount = subTokenReserve.Amount.Mul(sdk.NewInt(1))
		remainStake, err := valK.CalculateRemainStake(ctx, defaultNftStake, unStake)
		require.NoError(t, err)
		reserveTotalStake := subTokenReserve.Amount.Mul(sdk.NewInt(1))

		_, err = valK.BeginRedelegation(ctx, delAddr, valSrcAddr, valDstAddr, unStake, remainStake)
		require.NoError(t, err)

		rs, err := valK.GetValidatorRS(ctx, val.GetOperator())
		require.NoError(t, err)

		delegation, found := valK.GetDelegation(ctx, delAddr, valSrcAddr, unStake.GetID())
		require.True(t, found)

		// tokens are sub from delegation
		require.True(t, delegation.Stake.Equal(&remainStake))
		// validator power updated
		require.Equal(t, val.Stake-keeper.TokensToConsensusPower(reserveTotalStake), rs.Stake)

		oldExDel := expectDelegations[unStake.ID]
		oldExDel.Stake, err = oldExDel.Stake.Add(unStake)
		require.NoError(t, err)
		expectDelegations[unStake.ID] = oldExDel
	}

	blockHeader := ctx.BlockHeader()
	blockHeader.Time = blockHeader.Time.Add(time.Hour * 1000000)
	ctx = ctx.WithBlockHeader(blockHeader)
	err = valK.CompleteRedelegation(ctx, delAddr, valSrcAddr, valDstAddr)
	require.NoError(t, err)

	delegations = valK.GetValidatorDelegations(ctx, valDstAddr)
	require.NotNil(t, delegations)
	for _, v := range delegations {
		exDel, ok := expectDelegations[v.Stake.ID]
		require.True(t, ok)
		require.True(t, v.Equal(exDel))
	}
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
func TestRedelegationInStore(t *testing.T) {
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
