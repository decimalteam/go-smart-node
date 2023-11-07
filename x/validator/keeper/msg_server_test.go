package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
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

func TestCancelUndelegation(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(1000))),
	))

	// 1. create validator, delegate coins and nfts
	goCtx := sdk.WrapSDKContext(ctx)
	msg, err := types.NewMsgCreateValidator(
		vals[0],
		accs[0],
		PKs[0],
		types.Description{
			Moniker:  "monik",
			Identity: "somesome",
		},
		sdk.OneDec(),
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100))),
	)
	require.NoError(t, err)
	require.NoError(t, msg.ValidateBasic())
	_, err = msgsrv.CreateValidator(goCtx, msg)
	require.NoError(t, err)

	_, err = dsc.NFTKeeper.MintToken(goCtx, nfttypes.NewMsgMintToken(
		accs[0],
		"nfts",
		"token1",
		"http://localhost",
		false,
		accs[0],
		3,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100))),
	))
	require.NoError(t, err)

	_, err = msgsrv.DelegateNFT(goCtx, types.NewMsgDelegateNFT(
		accs[0],
		vals[0],
		"token1",
		[]uint32{1, 2, 3},
	))
	require.NoError(t, err)

	_, err = dsc.NFTKeeper.MintToken(goCtx, nfttypes.NewMsgMintToken(
		accs[0],
		"nfts",
		"token2",
		"http://localghost",
		false,
		accs[0],
		1,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100))),
	))
	require.NoError(t, err)

	require.NoError(t, err)
	_, err = msgsrv.DelegateNFT(goCtx, types.NewMsgDelegateNFT(
		accs[0],
		vals[0],
		"token2",
		[]uint32{1},
	))
	require.NoError(t, err)

	// 2. create partial undelegations
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	goCtx = sdk.WrapSDKContext(ctx)
	heightOfUndelegation := ctx.BlockHeight()

	_, err = msgsrv.Undelegate(goCtx, types.NewMsgUndelegate(
		accs[0],
		vals[0],
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(50))),
	))
	require.NoError(t, err)
	_, err = msgsrv.UndelegateNFT(goCtx, types.NewMsgUndelegateNFT(
		accs[0],
		vals[0],
		"token1",
		[]uint32{1, 3},
	))
	require.NoError(t, err)
	_, err = msgsrv.UndelegateNFT(goCtx, types.NewMsgUndelegateNFT(
		accs[0],
		vals[0],
		"token2",
		[]uint32{1},
	))
	require.NoError(t, err)

	// 3. cancel partial undelegations
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	goCtx = sdk.WrapSDKContext(ctx)

	_, err = msgsrv.CancelUndelegation(goCtx, types.NewMsgCancelUndelegation(
		accs[0],
		vals[0],
		heightOfUndelegation,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(25))),
	))
	require.NoError(t, err)
	_, err = msgsrv.CancelUndelegationNFT(goCtx, types.NewMsgCancelUndelegationNFT(
		accs[0],
		vals[0],
		heightOfUndelegation,
		"token1",
		[]uint32{3},
	))
	require.NoError(t, err)

	undel, found := dsc.ValidatorKeeper.GetUndelegation(ctx, accs[0], vals[0])
	require.True(t, found)
	require.Len(t, undel.Entries, 3)

	// 4. cancel full nft undelegation
	_, err = msgsrv.CancelUndelegationNFT(goCtx, types.NewMsgCancelUndelegationNFT(
		accs[0],
		vals[0],
		heightOfUndelegation,
		"token2",
		[]uint32{1},
	))
	require.NoError(t, err)

	undel, found = dsc.ValidatorKeeper.GetUndelegation(ctx, accs[0], vals[0])
	require.True(t, found)
	require.Len(t, undel.Entries, 2)

	// 5. checks
	// check coin delegations
	del, found := dsc.ValidatorKeeper.GetDelegation(ctx, accs[0], vals[0], cmdcfg.BaseDenom)
	require.True(t, found)
	require.True(t, del.Stake.Stake.Equal(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(75)))),
		"delegation stake is %s", del.Stake.Stake)
	// check nft delegations
	del, found = dsc.ValidatorKeeper.GetDelegation(ctx, accs[0], vals[0], "token1")
	require.True(t, found)
	require.True(t, types.SetHasSubset(del.Stake.SubTokenIDs, []uint32{2, 3}) && types.SetHasSubset([]uint32{2, 3}, del.Stake.SubTokenIDs))
	del, found = dsc.ValidatorKeeper.GetDelegation(ctx, accs[0], vals[0], "token2")
	require.True(t, found)
	require.True(t, types.SetHasSubset(del.Stake.SubTokenIDs, []uint32{1}) && types.SetHasSubset([]uint32{1}, del.Stake.SubTokenIDs))

	// check undelegation
	for _, ent := range undel.Entries {
		switch ent.Stake.ID {
		case cmdcfg.BaseDenom:
			require.True(t, ent.Stake.Stake.Equal(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(25)))))
		case "token1":
			require.True(t, types.SetHasSubset(ent.Stake.SubTokenIDs, []uint32{1}) && types.SetHasSubset([]uint32{1}, ent.Stake.SubTokenIDs))
		case "token2":
			require.Fail(t, "token2 must be total delegated")
		}
	}

	// 6. cancel too much, trying to get errors
	_, err = msgsrv.CancelUndelegation(goCtx, types.NewMsgCancelUndelegation(
		accs[0],
		vals[0],
		heightOfUndelegation,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100))),
	))
	require.Error(t, err)
	_, err = msgsrv.CancelUndelegationNFT(goCtx, types.NewMsgCancelUndelegationNFT(
		accs[0],
		vals[0],
		heightOfUndelegation,
		"token1",
		[]uint32{3},
	))
	require.Error(t, err)
	_, err = msgsrv.CancelUndelegation(goCtx, types.NewMsgCancelUndelegation(
		accs[0],
		vals[0],
		heightOfUndelegation-1, // no undelegation for this height
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(25))),
	))
	require.Error(t, err)

	// 7. cancel all remain undelegations
	_, err = msgsrv.CancelUndelegation(goCtx, types.NewMsgCancelUndelegation(
		accs[0],
		vals[0],
		heightOfUndelegation,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(25))),
	))
	require.NoError(t, err)
	_, err = msgsrv.CancelUndelegationNFT(goCtx, types.NewMsgCancelUndelegationNFT(
		accs[0],
		vals[0],
		heightOfUndelegation,
		"token1",
		[]uint32{1},
	))
	require.NoError(t, err)

	undel, found = dsc.ValidatorKeeper.GetUndelegation(ctx, accs[0], vals[0])
	require.False(t, found)
}

func TestCancelRedelegation(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(1000))),
	))

	// 1. create validators, delegate coins and nfts
	goCtx := sdk.WrapSDKContext(ctx)
	msg, err := types.NewMsgCreateValidator(
		vals[0],
		accs[0],
		PKs[0],
		types.Description{
			Moniker:  "monik",
			Identity: "somesome",
		},
		sdk.OneDec(),
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100))),
	)
	require.NoError(t, err)
	require.NoError(t, msg.ValidateBasic())
	_, err = msgsrv.CreateValidator(goCtx, msg)
	require.NoError(t, err)

	msg, err = types.NewMsgCreateValidator(
		vals[1],
		accs[1],
		PKs[1],
		types.Description{
			Moniker:  "monik2",
			Identity: "somesome",
		},
		sdk.OneDec(),
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100))),
	)
	require.NoError(t, err)
	require.NoError(t, msg.ValidateBasic())
	_, err = msgsrv.CreateValidator(goCtx, msg)
	require.NoError(t, err)
	// need set first validator online for 'bonded' status
	_, err = msgsrv.SetOnline(goCtx, types.NewMsgSetOnline(vals[0]))
	require.NoError(t, err)

	_, err = dsc.NFTKeeper.MintToken(goCtx, nfttypes.NewMsgMintToken(
		accs[0],
		"nfts",
		"token1",
		"http://localhost",
		false,
		accs[0],
		3,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100))),
	))
	require.NoError(t, err)

	_, err = msgsrv.DelegateNFT(goCtx, types.NewMsgDelegateNFT(
		accs[0],
		vals[0],
		"token1",
		[]uint32{1, 2, 3},
	))
	require.NoError(t, err)

	_, err = dsc.NFTKeeper.MintToken(goCtx, nfttypes.NewMsgMintToken(
		accs[0],
		"nfts",
		"token2",
		"http://localghost",
		false,
		accs[0],
		1,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100))),
	))
	require.NoError(t, err)

	require.NoError(t, err)
	_, err = msgsrv.DelegateNFT(goCtx, types.NewMsgDelegateNFT(
		accs[0],
		vals[0],
		"token2",
		[]uint32{1},
	))
	require.NoError(t, err)

	keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	// 2. create partial redelegations
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	goCtx = sdk.WrapSDKContext(ctx)
	heightOfRedelegation := ctx.BlockHeight()

	_, err = msgsrv.Redelegate(goCtx, types.NewMsgRedelegate(
		accs[0],
		vals[0],
		vals[1],
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(50))),
	))
	require.NoError(t, err)
	_, err = msgsrv.RedelegateNFT(goCtx, types.NewMsgRedelegateNFT(
		accs[0],
		vals[0],
		vals[1],
		"token1",
		[]uint32{1, 3},
	))
	require.NoError(t, err)
	_, err = msgsrv.RedelegateNFT(goCtx, types.NewMsgRedelegateNFT(
		accs[0],
		vals[0],
		vals[1],
		"token2",
		[]uint32{1},
	))
	require.NoError(t, err)

	// 3. cancel partial redelegations
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	goCtx = sdk.WrapSDKContext(ctx)

	_, err = msgsrv.CancelRedelegation(goCtx, types.NewMsgCancelRedelegation(
		accs[0],
		vals[0],
		vals[1],
		heightOfRedelegation,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(25))),
	))
	require.NoError(t, err)
	_, err = msgsrv.CancelRedelegationNFT(goCtx, types.NewMsgCancelRedelegationNFT(
		accs[0],
		vals[0],
		vals[1],
		heightOfRedelegation,
		"token1",
		[]uint32{3},
	))
	require.NoError(t, err)

	redel, found := dsc.ValidatorKeeper.GetRedelegation(ctx, accs[0], vals[0], vals[1])
	require.True(t, found)
	require.Len(t, redel.Entries, 3)

	// 4. cancel full nft undelegation
	_, err = msgsrv.CancelRedelegationNFT(goCtx, types.NewMsgCancelRedelegationNFT(
		accs[0],
		vals[0],
		vals[1],
		heightOfRedelegation,
		"token2",
		[]uint32{1},
	))
	require.NoError(t, err)

	redel, found = dsc.ValidatorKeeper.GetRedelegation(ctx, accs[0], vals[0], vals[1])
	require.True(t, found)
	require.Len(t, redel.Entries, 2)

	// 5. checks
	// check coin delegations
	del, found := dsc.ValidatorKeeper.GetDelegation(ctx, accs[0], vals[0], cmdcfg.BaseDenom)
	require.True(t, found)
	require.True(t, del.Stake.Stake.Equal(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(75)))),
		"delegation stake is %s", del.Stake.Stake)
	// check nft delegations
	del, found = dsc.ValidatorKeeper.GetDelegation(ctx, accs[0], vals[0], "token1")
	require.True(t, found)
	require.True(t, types.SetHasSubset(del.Stake.SubTokenIDs, []uint32{2, 3}) && types.SetHasSubset([]uint32{2, 3}, del.Stake.SubTokenIDs))
	del, found = dsc.ValidatorKeeper.GetDelegation(ctx, accs[0], vals[0], "token2")
	require.True(t, found)
	require.True(t, types.SetHasSubset(del.Stake.SubTokenIDs, []uint32{1}) && types.SetHasSubset([]uint32{1}, del.Stake.SubTokenIDs))

	// check undelegation
	for _, ent := range redel.Entries {
		switch ent.Stake.ID {
		case cmdcfg.BaseDenom:
			require.True(t, ent.Stake.Stake.Equal(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(25)))))
		case "token1":
			require.True(t, types.SetHasSubset(ent.Stake.SubTokenIDs, []uint32{1}) && types.SetHasSubset([]uint32{1}, ent.Stake.SubTokenIDs))
		case "token2":
			require.Fail(t, "token2 must be total delegated")
		}
	}

	// 6. cancel too much, trying to get errors
	_, err = msgsrv.CancelRedelegation(goCtx, types.NewMsgCancelRedelegation(
		accs[0],
		vals[0],
		vals[1],
		heightOfRedelegation,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100))),
	))
	require.Error(t, err)
	_, err = msgsrv.CancelRedelegationNFT(goCtx, types.NewMsgCancelRedelegationNFT(
		accs[0],
		vals[0],
		vals[1],
		heightOfRedelegation,
		"token1",
		[]uint32{3},
	))
	require.Error(t, err)
	_, err = msgsrv.CancelRedelegation(goCtx, types.NewMsgCancelRedelegation(
		accs[0],
		vals[0],
		vals[1],
		heightOfRedelegation-1, // no undelegation for this height
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(25))),
	))
	require.Error(t, err)

	// 7. cancel all remain undelegations
	_, err = msgsrv.CancelRedelegation(goCtx, types.NewMsgCancelRedelegation(
		accs[0],
		vals[0],
		vals[1],
		heightOfRedelegation,
		sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(25))),
	))
	require.NoError(t, err)
	_, err = msgsrv.CancelRedelegationNFT(goCtx, types.NewMsgCancelRedelegationNFT(
		accs[0],
		vals[0],
		vals[1],
		heightOfRedelegation,
		"token1",
		[]uint32{1},
	))
	require.NoError(t, err)

	redel, found = dsc.ValidatorKeeper.GetRedelegation(ctx, accs[0], vals[0], vals[1])
	require.False(t, found)
}
