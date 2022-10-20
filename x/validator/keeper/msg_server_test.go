package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func TestCreateValidator(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	stake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(stake))

	goCtx := sdk.WrapSDKContext(ctx)

	// 0.
	startBalance := dsc.BankKeeper.GetBalance(ctx, dsc.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress(), stake.Denom)
	// 1. create regular validator
	msg, err := types.NewMsgCreateValidator(
		vals[0],
		accs[0],
		PKs[0],
		types.Description{
			Moniker: "monik",
		},
		sdk.OneDec(), stake)
	require.NoError(t, err)
	require.NoError(t, msg.ValidateBasic())
	_, err = msgsrv.CreateValidator(goCtx, msg)
	require.NoError(t, err)

	// check validator
	validator, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
	require.True(t, found)
	require.Equal(t, vals[0].String(), validator.OperatorAddress)
	require.Equal(t, accs[0].String(), validator.RewardAddress)
	pk, err := validator.ConsPubKey()
	require.NoError(t, err)
	require.Equal(t, PKs[0], pk)
	require.Equal(t, types.Description{Moniker: "monik"}, validator.Description)
	require.Equal(t, sdk.OneDec(), validator.Commission)

	// check delegations
	delegations := dsc.ValidatorKeeper.GetValidatorDelegations(ctx, vals[0])
	require.Len(t, delegations, 1)
	require.Equal(t, stake, delegations[0].Stake.Stake)

	// check pool changes
	balance := dsc.BankKeeper.GetBalance(ctx, dsc.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress(), stake.Denom)
	require.Equal(t, stake.Amount, balance.Amount.Sub(startBalance.Amount))

	// 2. create with same public key
	msg, err = types.NewMsgCreateValidator(
		vals[1],
		accs[1],
		PKs[0],
		types.Description{
			Moniker: "monik2",
		},
		sdk.OneDec(), stake)
	require.NoError(t, err)
	require.NoError(t, msg.ValidateBasic())
	_, err = msgsrv.CreateValidator(goCtx, msg)
	require.Error(t, err)
}

func TestEditValidator(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	stake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(stake))

	goCtx := sdk.WrapSDKContext(ctx)
	// 1. create validator
	msg, err := types.NewMsgCreateValidator(
		vals[0],
		accs[0],
		PKs[0],
		types.Description{
			Moniker:  "monik",
			Identity: "somesome",
		},
		sdk.OneDec(), stake)
	require.NoError(t, err)
	require.NoError(t, msg.ValidateBasic())
	_, err = msgsrv.CreateValidator(goCtx, msg)
	require.NoError(t, err)

	// 2. edit
	msgEdit := types.NewMsgEditValidator(
		vals[0],
		accs[1],
		types.Description{
			Moniker:  "monik2",
			Identity: types.DoNotModifyDesc, //
		},
	)
	require.NoError(t, msgEdit.ValidateBasic())
	_, err = msgsrv.EditValidator(goCtx, msgEdit)
	require.NoError(t, err)

	validator, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
	require.True(t, found)
	require.Equal(t, accs[1].String(), validator.RewardAddress)
	require.Equal(t, types.Description{
		Moniker:  "monik2",
		Identity: "somesome",
	}, validator.Description)
}

//
//func TestCancelUnbondingDelegation(t *testing.T) {
//	// setup the app
//	app := simapp.Setup(t, false)
//	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
//	msgServer := keeper.NewMsgServerImpl(app.StakingKeeper)
//	bondDenom := app.StakingKeeper.BondDenom(ctx)
//
//	// set the not bonded pool module account
//	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
//	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 5)
//
//	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens))))
//	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	moduleBalance := app.BankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), app.StakingKeeper.BondDenom(ctx))
//	require.Equal(t, sdk.NewInt64Coin(bondDenom, startTokens.Int64()), moduleBalance)
//
//	// accounts
//	delAddrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.NewInt(10000))
//	validators := app.StakingKeeper.GetValidators(ctx, 10)
//	require.Equal(t, len(validators), 1)
//
//	validatorAddr, err := sdk.ValAddressFromBech32(validators[0].OperatorAddress)
//	require.NoError(t, err)
//	delegatorAddr := delAddrs[0]
//
//	// setting the ubd entry
//	unbondingAmount := sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 5)
//	ubd := types.NewUnbondingDelegation(
//		delegatorAddr, validatorAddr, 10,
//		ctx.BlockTime().Add(time.Minute*10),
//		unbondingAmount.Amount,
//	)
//
//	// set and retrieve a record
//	app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
//	resUnbond, found := app.StakingKeeper.GetUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
//	require.True(t, found)
//	require.Equal(t, ubd, resUnbond)
//
//	testCases := []struct {
//		Name      string
//		ExceptErr bool
//		req       types.MsgCancelUnbondingDelegation
//	}{
//		{
//			Name:      "invalid height",
//			ExceptErr: true,
//			req: types.MsgCancelUnbondingDelegation{
//				DelegatorAddress: resUnbond.DelegatorAddress,
//				ValidatorAddress: resUnbond.ValidatorAddress,
//				Amount:           sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), sdk.NewInt(4)),
//				CreationHeight:   0,
//			},
//		},
//		{
//			Name:      "invalid coin",
//			ExceptErr: true,
//			req: types.MsgCancelUnbondingDelegation{
//				DelegatorAddress: resUnbond.DelegatorAddress,
//				ValidatorAddress: resUnbond.ValidatorAddress,
//				Amount:           sdk.NewCoin("dump_coin", sdk.NewInt(4)),
//				CreationHeight:   0,
//			},
//		},
//		{
//			Name:      "validator not exists",
//			ExceptErr: true,
//			req: types.MsgCancelUnbondingDelegation{
//				DelegatorAddress: resUnbond.DelegatorAddress,
//				ValidatorAddress: sdk.ValAddress(sdk.AccAddress("asdsad")).String(),
//				Amount:           unbondingAmount,
//				CreationHeight:   0,
//			},
//		},
//		{
//			Name:      "invalid delegator address",
//			ExceptErr: true,
//			req: types.MsgCancelUnbondingDelegation{
//				DelegatorAddress: "invalid_delegator_addrtess",
//				ValidatorAddress: resUnbond.ValidatorAddress,
//				Amount:           unbondingAmount,
//				CreationHeight:   0,
//			},
//		},
//		{
//			Name:      "invalid amount",
//			ExceptErr: true,
//			req: types.MsgCancelUnbondingDelegation{
//				DelegatorAddress: resUnbond.DelegatorAddress,
//				ValidatorAddress: resUnbond.ValidatorAddress,
//				Amount:           unbondingAmount.Add(sdk.NewInt64Coin(bondDenom, 10)),
//				CreationHeight:   10,
//			},
//		},
//		{
//			Name:      "success",
//			ExceptErr: false,
//			req: types.MsgCancelUnbondingDelegation{
//				DelegatorAddress: resUnbond.DelegatorAddress,
//				ValidatorAddress: resUnbond.ValidatorAddress,
//				Amount:           unbondingAmount.Sub(sdk.NewInt64Coin(bondDenom, 1)),
//				CreationHeight:   10,
//			},
//		},
//		{
//			Name:      "success",
//			ExceptErr: false,
//			req: types.MsgCancelUnbondingDelegation{
//				DelegatorAddress: resUnbond.DelegatorAddress,
//				ValidatorAddress: resUnbond.ValidatorAddress,
//				Amount:           unbondingAmount.Sub(unbondingAmount.Sub(sdk.NewInt64Coin(bondDenom, 1))),
//				CreationHeight:   10,
//			},
//		},
//	}
//
//	for _, testCase := range testCases {
//		t.Run(testCase.Name, func(t *testing.T) {
//			_, err := msgServer.CancelUnbondingDelegation(ctx, &testCase.req)
//			if testCase.ExceptErr {
//				require.Error(t, err)
//			} else {
//				require.NoError(t, err)
//				balanceForNotBondedPool := app.BankKeeper.GetBalance(ctx, sdk.AccAddress(notBondedPool.GetAddress()), bondDenom)
//				require.Equal(t, balanceForNotBondedPool, moduleBalance.Sub(testCase.req.Amount))
//				moduleBalance = moduleBalance.Sub(testCase.req.Amount)
//			}
//		})
//	}
//}
