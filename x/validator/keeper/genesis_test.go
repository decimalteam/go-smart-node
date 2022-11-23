package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/validator"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/testvalidator"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func bootstrapGenesisTest(t *testing.T, numAddrs int) (*app.DSC, sdk.Context, []sdk.AccAddress) {
	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	addrDels, _ := generateAddresses(dsc, ctx, numAddrs,
		sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(10))),
	)
	return dsc, ctx, addrDels
}

// compare simple, without nft
func isEqualDelegations(dels1, dels2 []types.Delegation) bool {
	if len(dels1) != len(dels2) {
		return false
	}
	for i := range dels1 {
		isHere := false
		for j := range dels2 {
			if dels1[i].Delegator == dels2[j].Delegator &&
				dels1[i].Validator == dels2[j].Validator &&
				dels1[i].Stake.Equal(&dels2[j].Stake) {
				isHere = true
				break
			}
		}
		if !isHere {
			return false
		}
	}
	return true
}

func TestInitGenesis(t *testing.T) {
	dsc, ctx, addrs := bootstrapGenesisTest(t, 10)

	valTokens := int64(1000)

	params := dsc.ValidatorKeeper.GetParams(ctx)
	validators := dsc.ValidatorKeeper.GetAllValidators(ctx)
	require.Len(t, validators, 1)
	var delegations []types.Delegation

	pk0, err := codectypes.NewAnyWithValue(PKs[0])
	require.NoError(t, err)

	pk1, err := codectypes.NewAnyWithValue(PKs[1])
	require.NoError(t, err)

	// initialize the validators
	bondedVal1 := types.Validator{
		OperatorAddress: sdk.ValAddress(addrs[0]).String(),
		ConsensusPubkey: pk0,
		Status:          types.BondStatus_Bonded,
		Stake:           valTokens,
		Description:     types.NewDescription("hoop", "", "", "", ""),
	}
	bondedVal2 := types.Validator{
		OperatorAddress: sdk.ValAddress(addrs[1]).String(),
		ConsensusPubkey: pk1,
		Status:          types.BondStatus_Bonded,
		Stake:           valTokens,
		Description:     types.NewDescription("bloop", "", "", "", ""),
	}

	// append new bonded validators to the list
	validators = append(validators, bondedVal1, bondedVal2)

	// mint coins in the bonded pool representing the validators coins
	i2 := int64(len(validators) - 1) // -1 to exclude genesis validator
	require.NoError(t,
		testvalidator.FundModuleAccount(
			dsc.BankKeeper,
			ctx,
			types.BondedPoolName,
			sdk.NewCoins(
				sdk.NewCoin(params.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(i2*valTokens))),
			),
		),
	)

	genesisDelegations := dsc.ValidatorKeeper.GetAllDelegations(ctx)
	delegations = append(delegations, genesisDelegations...)
	delegations = append(delegations, types.NewDelegation(
		addrs[0],
		sdk.ValAddress(addrs[0]),
		types.NewStakeCoin(sdk.NewCoin(params.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(valTokens)))),
	))
	delegations = append(delegations, types.NewDelegation(
		addrs[1],
		sdk.ValAddress(addrs[1]),
		types.NewStakeCoin(sdk.NewCoin(params.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(valTokens)))),
	))

	genesisState := types.NewGenesisState(params, validators, delegations,
		[]types.LastValidatorPower{
			{bondedVal1.OperatorAddress, valTokens},
			{bondedVal2.OperatorAddress, valTokens},
		},
	)
	vals := dsc.ValidatorKeeper.InitGenesis(ctx, genesisState)

	actualGenesis := dsc.ValidatorKeeper.ExportGenesis(ctx)
	require.Equal(t, genesisState.Params, actualGenesis.Params)
	require.True(t, isEqualDelegations(genesisState.Delegations, actualGenesis.Delegations))
	require.EqualValues(t, dsc.ValidatorKeeper.GetAllValidators(ctx), actualGenesis.Validators)

	// Ensure validators have addresses.
	vals2, err := validator.WriteValidators(ctx, dsc.ValidatorKeeper)
	require.NoError(t, err)

	for _, val := range vals2 {
		require.NotEmpty(t, val.Address)
	}

	// now make sure the validators are bonded and intra-tx counters are correct
	resVal, found := dsc.ValidatorKeeper.GetValidator(ctx, sdk.ValAddress(addrs[0]))
	require.True(t, found)
	require.Equal(t, types.BondStatus_Bonded, resVal.Status)

	resVal, found = dsc.ValidatorKeeper.GetValidator(ctx, sdk.ValAddress(addrs[1]))
	require.True(t, found)
	require.Equal(t, types.BondStatus_Bonded, resVal.Status)

	abcivals := make([]abci.ValidatorUpdate, len(vals))

	validators = validators[1:] // remove genesis validator
	for i, val := range validators {
		abcivals[i] = val.ABCIValidatorUpdate(sdkmath.NewInt(0)) // TODO: refactor ABCIValidatorUpdate - r is not needed
	}

	require.Equal(t, abcivals, vals)
}

func TestInitGenesis_PoolsBalanceMismatch(t *testing.T) {
	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	consPub, err := codectypes.NewAnyWithValue(PKs[0])
	require.NoError(t, err)

	validator := types.Validator{
		OperatorAddress: sdk.ValAddress("12345678901234567890").String(),
		ConsensusPubkey: consPub,
		Jailed:          false,
		Stake:           10,
		Description:     types.NewDescription("bloop", "", "", "", ""),
	}

	params := types.Params{
		UndelegationTime: 10000,
		MaxValidators:    1,
		MaxEntries:       10,
		BaseDenom:        "stake",
	}

	require.Panics(t, func() {
		// setting validator status to bonded so the balance counts towards bonded pool
		validator.Status = types.BondStatus_Bonded
		dsc.ValidatorKeeper.InitGenesis(ctx, &types.GenesisState{
			Params:     params,
			Validators: []types.Validator{validator},
		})
	},
		"should panic because bonded pool balance is different from bonded pool coins",
	)

	require.Panics(t, func() {
		// setting validator status to unbonded so the balance counts towards not bonded pool
		validator.Status = types.BondStatus_Unbonded
		dsc.ValidatorKeeper.InitGenesis(ctx, &types.GenesisState{
			Params:     params,
			Validators: []types.Validator{validator},
		})
	},
		"should panic because not bonded pool balance is different from not bonded pool coins",
	)
}

func TestInitGenesisLargeValidatorSet(t *testing.T) {
	size := 200
	require.True(t, size > 100)

	dsc, ctx, addrs := bootstrapGenesisTest(t, 200)
	genesisValidators := dsc.ValidatorKeeper.GetAllValidators(ctx)

	params := dsc.ValidatorKeeper.GetParams(ctx)
	params.MaxValidators = 100
	dsc.ValidatorKeeper.SetParams(ctx, params)

	delegations := dsc.ValidatorKeeper.GetAllDelegations(ctx)
	validators := make([]types.Validator, size)
	powers := []types.LastValidatorPower{}

	var err error

	bondedPoolAmt := sdk.ZeroInt()
	for i := range validators {
		validators[i], err = types.NewValidator(
			sdk.ValAddress(addrs[i]),
			addrs[i],
			PKs[i],
			types.NewDescription(fmt.Sprintf("#%d", i), "", "", "", ""),
			sdk.ZeroDec(),
		)
		require.NoError(t, err)
		validators[i].Status = types.BondStatus_Bonded

		validators[i].Stake = 1
		if i < 100 {
			validators[i].Stake = 2
		}
		tokens := helpers.EtherToWei(sdkmath.NewInt(validators[i].Stake))
		powers = append(powers, types.LastValidatorPower{
			Address: sdk.ValAddress(addrs[i]).String(),
			Power:   validators[i].Stake,
		})
		delegations = append(delegations, types.NewDelegation(
			addrs[i],
			sdk.ValAddress(addrs[i]),
			types.NewStakeCoin(sdk.NewCoin(cmdcfg.BaseDenom, tokens)),
		))
		// add bonded coins
		bondedPoolAmt = bondedPoolAmt.Add(tokens)
	}

	validators = append(validators, genesisValidators...)
	genesisState := types.NewGenesisState(params, validators, delegations, powers)

	// mint coins in the bonded pool representing the validators coins
	require.NoError(t,
		testvalidator.FundModuleAccount(
			dsc.BankKeeper,
			ctx,
			types.BondedPoolName,
			sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, bondedPoolAmt)),
		),
	)

	vals := dsc.ValidatorKeeper.InitGenesis(ctx, genesisState)

	abcivals := make([]abci.ValidatorUpdate, 100)
	for i, val := range validators[:100] {
		abcivals[i] = val.ABCIValidatorUpdate(sdkmath.ZeroInt())
	}

	// remove genesis validator
	vals = vals[:100]
	require.Equal(t, abcivals, vals)
}
