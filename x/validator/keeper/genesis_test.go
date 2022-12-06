package keeper_test

import (
	"fmt"
	"testing"
	"time"

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
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
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

// test helper indexes after init of genesis
func TestDelegationsAfterGenesis(t *testing.T) {
	var completionTime1 = time.Now().Add(time.Second)
	var completionTime2 = time.Now().Add(time.Second * 2)

	dsc, ctx, addrs := bootstrapGenesisTest(t, 200)

	genesisVal := dsc.ValidatorKeeper.GetValidators(ctx, 10)[0]

	dsc.CoinKeeper.SetCoin(ctx, cointypes.Coin{
		Denom:       "custom1",
		Title:       "custom1",
		Creator:     addrs[0].String(),
		CRR:         10,
		LimitVolume: helpers.EtherToWei(sdkmath.NewInt(1_000_000)),
		Identity:    "",
		Volume:      helpers.EtherToWei(sdkmath.NewInt(2_000)),
		Reserve:     helpers.EtherToWei(sdkmath.NewInt(1_000)),
	})

	dsc.CoinKeeper.SetCoin(ctx, cointypes.Coin{
		Denom:       "custom2",
		Title:       "custom2",
		Creator:     addrs[1].String(),
		CRR:         20,
		LimitVolume: helpers.EtherToWei(sdkmath.NewInt(1_000_000)),
		Identity:    "",
		Volume:      helpers.EtherToWei(sdkmath.NewInt(3_000)),
		Reserve:     helpers.EtherToWei(sdkmath.NewInt(1_000)),
	})

	pk0, err := codectypes.NewAnyWithValue(PKs[0])
	require.NoError(t, err)
	pk1, err := codectypes.NewAnyWithValue(PKs[1])
	require.NoError(t, err)

	// validator 0 bonded, validator 1 unbonded
	genesisState := &types.GenesisState{
		Params: types.DefaultParams(),
		Validators: []types.Validator{
			{
				OperatorAddress: sdk.ValAddress(addrs[0]).String(),
				RewardAddress:   addrs[0].String(),
				ConsensusPubkey: pk0,
				Description:     types.NewDescription("#0", "", "", "", ""),
				Commission:      sdk.ZeroDec(),
				Status:          types.BondStatus_Bonded,
				Online:          true,
				Stake:           1000,
			},
			{
				OperatorAddress: sdk.ValAddress(addrs[1]).String(),
				RewardAddress:   addrs[1].String(),
				ConsensusPubkey: pk1,
				Description:     types.NewDescription("#1", "", "", "", ""),
				Commission:      sdk.ZeroDec(),
				Status:          types.BondStatus_Unbonded,
				Online:          false,
				Stake:           0,
			},
		},
		Delegations: []types.Delegation{
			// validator 0, delegator 0
			{
				Delegator: addrs[0].String(),
				Validator: sdk.ValAddress(addrs[0]).String(),
				Stake:     types.NewStakeCoin(sdk.NewCoin("custom1", helpers.EtherToWei(sdkmath.NewInt(10)))), // 10custom1
			},
			{
				Delegator: addrs[0].String(),
				Validator: sdk.ValAddress(addrs[0]).String(),
				Stake:     types.NewStakeCoin(sdk.NewCoin("custom2", helpers.EtherToWei(sdkmath.NewInt(100)))), // 100custom2
			},
			// validator 0, delegator 1
			{
				Delegator: addrs[1].String(),
				Validator: sdk.ValAddress(addrs[0]).String(),
				Stake:     types.NewStakeCoin(sdk.NewCoin("custom1", helpers.EtherToWei(sdkmath.NewInt(20)))), // 20custom1
			},
			{
				Delegator: addrs[1].String(),
				Validator: sdk.ValAddress(addrs[0]).String(),
				Stake:     types.NewStakeCoin(sdk.NewCoin("custom2", helpers.EtherToWei(sdkmath.NewInt(200)))), // 200custom2
			},
			// validator 1, delegator 0
			{
				Delegator: addrs[0].String(),
				Validator: sdk.ValAddress(addrs[1]).String(),
				Stake:     types.NewStakeCoin(sdk.NewCoin("custom1", helpers.EtherToWei(sdkmath.NewInt(11)))), // 11custom1
			},
			{
				Delegator: addrs[0].String(),
				Validator: sdk.ValAddress(addrs[1]).String(),
				Stake:     types.NewStakeCoin(sdk.NewCoin("custom2", helpers.EtherToWei(sdkmath.NewInt(111)))), // 111custom2
			},
			// validator 1, delegator 1
			{
				Delegator: addrs[1].String(),
				Validator: sdk.ValAddress(addrs[1]).String(),
				Stake:     types.NewStakeCoin(sdk.NewCoin("custom1", helpers.EtherToWei(sdkmath.NewInt(22)))), // 22custom1
			},
			{
				Delegator: addrs[1].String(),
				Validator: sdk.ValAddress(addrs[1]).String(),
				Stake:     types.NewStakeCoin(sdk.NewCoin("custom2", helpers.EtherToWei(sdkmath.NewInt(222)))), // 222custom2
			},
		},
		Undelegations: []types.Undelegation{
			{
				Delegator: addrs[0].String(),
				Validator: sdk.ValAddress(addrs[0]).String(),
				Entries: []types.UndelegationEntry{
					{
						CreationHeight: 0,
						CompletionTime: completionTime1,
						Stake:          types.NewStakeCoin(sdk.NewCoin("custom1", helpers.EtherToWei(sdkmath.NewInt(30)))),
					},
					{
						CreationHeight: 0,
						CompletionTime: completionTime2,
						Stake:          types.NewStakeCoin(sdk.NewCoin("custom2", helpers.EtherToWei(sdkmath.NewInt(300)))),
					},
				},
			},
			{
				Delegator: addrs[0].String(),
				Validator: sdk.ValAddress(addrs[1]).String(),
				Entries: []types.UndelegationEntry{
					{
						CreationHeight: 0,
						CompletionTime: completionTime1,
						Stake:          types.NewStakeCoin(sdk.NewCoin("custom1", helpers.EtherToWei(sdkmath.NewInt(31)))),
					},
					{
						CreationHeight: 0,
						CompletionTime: completionTime2,
						Stake:          types.NewStakeCoin(sdk.NewCoin("custom2", helpers.EtherToWei(sdkmath.NewInt(311)))),
					},
				},
			},
			{
				Delegator: addrs[1].String(),
				Validator: sdk.ValAddress(addrs[0]).String(),
				Entries: []types.UndelegationEntry{
					{
						CreationHeight: 0,
						CompletionTime: completionTime1,
						Stake:          types.NewStakeCoin(sdk.NewCoin("custom1", helpers.EtherToWei(sdkmath.NewInt(33)))),
					},
					{
						CreationHeight: 0,
						CompletionTime: completionTime2,
						Stake:          types.NewStakeCoin(sdk.NewCoin("custom2", helpers.EtherToWei(sdkmath.NewInt(333)))),
					},
				},
			},
		},
		Redelegations: []types.Redelegation{
			{
				Delegator:    addrs[0].String(),
				ValidatorSrc: sdk.ValAddress(addrs[0]).String(),
				ValidatorDst: sdk.ValAddress(addrs[1]).String(),
				Entries: []types.RedelegationEntry{
					{
						CreationHeight: 0,
						CompletionTime: completionTime1,
						Stake:          types.NewStakeCoin(sdk.NewCoin("custom1", helpers.EtherToWei(sdkmath.NewInt(40)))),
					},
					{
						CreationHeight: 0,
						CompletionTime: completionTime2,
						Stake:          types.NewStakeCoin(sdk.NewCoin("custom2", helpers.EtherToWei(sdkmath.NewInt(400)))),
					},
				},
			},
			{
				Delegator:    addrs[0].String(),
				ValidatorSrc: sdk.ValAddress(addrs[1]).String(),
				ValidatorDst: sdk.ValAddress(addrs[0]).String(),
				Entries: []types.RedelegationEntry{
					{
						CreationHeight: 0,
						CompletionTime: completionTime1,
						Stake:          types.NewStakeCoin(sdk.NewCoin("custom1", helpers.EtherToWei(sdkmath.NewInt(41)))),
					},
					{
						CreationHeight: 0,
						CompletionTime: completionTime2,
						Stake:          types.NewStakeCoin(sdk.NewCoin("custom2", helpers.EtherToWei(sdkmath.NewInt(411)))),
					},
				},
			},
			{
				Delegator:    addrs[1].String(),
				ValidatorSrc: sdk.ValAddress(addrs[0]).String(),
				ValidatorDst: sdk.ValAddress(addrs[1]).String(),
				Entries: []types.RedelegationEntry{
					{
						CreationHeight: 0,
						CompletionTime: completionTime1,
						Stake:          types.NewStakeCoin(sdk.NewCoin("custom1", helpers.EtherToWei(sdkmath.NewInt(44)))),
					},
					{
						CreationHeight: 0,
						CompletionTime: completionTime2,
						Stake:          types.NewStakeCoin(sdk.NewCoin("custom2", helpers.EtherToWei(sdkmath.NewInt(444)))),
					},
				},
			},
		},
	}

	genesisState.Delegations = append(genesisState.Delegations, dsc.ValidatorKeeper.GetAllDelegations(ctx)...)
	genesisState.Validators = append(genesisState.Validators, dsc.ValidatorKeeper.GetAllValidators(ctx)...)

	bondedPool := sdk.NewCoins(
		sdk.NewCoin("custom1", helpers.EtherToWei(sdkmath.NewInt(10+20))),
		sdk.NewCoin("custom2", helpers.EtherToWei(sdkmath.NewInt(100+200))),
	)

	notBondedPool := sdk.NewCoins(
		sdk.NewCoin("custom1", helpers.EtherToWei(sdkmath.NewInt(11+22+30+31+33+40+41+44))),
		sdk.NewCoin("custom2", helpers.EtherToWei(sdkmath.NewInt(111+222+300+311+333+400+411+444))),
	)

	testvalidator.FundModuleAccount(
		dsc.BankKeeper,
		ctx,
		types.BondedPoolName,
		bondedPool,
	)
	testvalidator.FundModuleAccount(
		dsc.BankKeeper,
		ctx,
		types.NotBondedPoolName,
		notBondedPool,
	)

	dsc.ValidatorKeeper.InitGenesis(ctx, genesisState)

	// check custom coin pool
	ccs := dsc.ValidatorKeeper.GetAllCustomCoinsStaked(ctx)
	for denom, amount := range ccs {
		switch denom {
		case "custom1":
			require.True(t, amount.Equal(helpers.EtherToWei(sdkmath.NewInt(10+20+11+22+30+31+33+40+41+44))))
		case "custom2":
			require.True(t, amount.Equal(helpers.EtherToWei(sdkmath.NewInt(100+200+111+222+300+311+333+400+411+444))))
		}
	}

	// check delegations
	require.Len(t, dsc.ValidatorKeeper.GetAllDelegations(ctx), 9)

	// check delegations counts
	counts := dsc.ValidatorKeeper.GetAllDelegationsCount(ctx)
	for val, count := range counts {
		switch val {
		case sdk.ValAddress(addrs[0]).String():
			require.Equal(t, uint32(4), count)
		case sdk.ValAddress(addrs[1]).String():
			require.Equal(t, uint32(4), count)
		}
	}

	// check delegations indexes
	// GetValidatorDelegations uses 'DelegationsByVal:        0x37<validator><delegator><stake_id>'
	require.Len(t, dsc.ValidatorKeeper.GetValidatorDelegations(ctx, genesisVal.GetOperator()), 1)
	require.Len(t, dsc.ValidatorKeeper.GetValidatorDelegations(ctx, sdk.ValAddress(addrs[0])), 4)
	require.Len(t, dsc.ValidatorKeeper.GetValidatorDelegations(ctx, sdk.ValAddress(addrs[1])), 4)

	// check undelegation indexes
	require.Len(t, dsc.ValidatorKeeper.GetAllUndelegations(ctx, addrs[0]), 2)
	require.Len(t, dsc.ValidatorKeeper.GetAllUndelegations(ctx, addrs[1]), 1)
	require.Len(t, dsc.ValidatorKeeper.GetUndelegationsFromValidator(ctx, sdk.ValAddress(addrs[0])), 2)
	require.Len(t, dsc.ValidatorKeeper.GetUndelegationsFromValidator(ctx, sdk.ValAddress(addrs[1])), 1)
	require.Len(t, dsc.ValidatorKeeper.GetUBDQueueTimeSlice(ctx, completionTime1), 3)
	require.Len(t, dsc.ValidatorKeeper.GetUBDQueueTimeSlice(ctx, completionTime2), 3)

	// check redelegation indexes
	require.Len(t, dsc.ValidatorKeeper.GetRedelegations(ctx, addrs[0], 100), 2)
	require.Len(t, dsc.ValidatorKeeper.GetRedelegations(ctx, addrs[1], 100), 1)
	require.Len(t, dsc.ValidatorKeeper.GetRedelegationsFromSrcValidator(ctx, sdk.ValAddress(addrs[0])), 2)
	require.Len(t, dsc.ValidatorKeeper.GetRedelegationsFromSrcValidator(ctx, sdk.ValAddress(addrs[1])), 1)
	require.Len(t, dsc.ValidatorKeeper.GetRedelegationQueueTimeSlice(ctx, completionTime1), 3)
	require.Len(t, dsc.ValidatorKeeper.GetRedelegationQueueTimeSlice(ctx, completionTime2), 3)
}
