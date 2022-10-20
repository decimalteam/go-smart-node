package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// TODO: add delegations nft, add checks of nft owners
func TestStateOnlineOffline(t *testing.T) {
	_, dsc, ctx := createTestInput(t)
	msgsrv := keeper.NewMsgServerImpl(dsc.ValidatorKeeper)

	// 0. genesis
	balance := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100000)))
	accs, vals := generateAddresses(dsc, ctx, 10, sdk.NewCoins(balance))

	genesisVals := dsc.ValidatorKeeper.GetValidators(ctx, 10)
	require.Len(t, genesisVals, 1)
	genesisVal := genesisVals[0]
	require.True(t, genesisVal.ConsensusPower() > 0)
	require.Equal(t, genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
	balanceNB := dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress())
	require.True(t, balanceNB.IsZero())
	startBalanceB := dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetBondedPool(ctx).GetAddress())
	balanceB := dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetBondedPool(ctx).GetAddress())

	////////////////////////////////////////////////
	// 1. create second validator
	creatorStake := sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(100)))
	msgCreate, err := types.NewMsgCreateValidator(vals[0], accs[0], PKs[0], types.Description{Moniker: "monik"},
		sdk.ZeroDec(), creatorStake)
	require.NoError(t, err)
	goCtx := sdk.WrapSDKContext(ctx)
	_, err = msgsrv.CreateValidator(goCtx, msgCreate)
	require.NoError(t, err)

	updates := keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})
	// new validator is not online, there is not changes in tendermint validators and powers
	require.Len(t, updates, 0)
	require.Equal(t, genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)

	balanceNB = dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress())
	require.True(t, balanceNB.IsEqual(sdk.NewCoins(creatorStake)))
	balanceB = dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetBondedPool(ctx).GetAddress())
	require.True(t, balanceB.IsEqual(startBalanceB))

	////////////////////////////////////////////////
	// 2. increment block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))

	////////////////////////////////////////////////
	// 3. set second validator online
	msgOnline := types.NewMsgSetOnline(vals[0])
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = msgsrv.SetOnline(goCtx, msgOnline)
	require.NoError(t, err)
	// last validators must be changes after ApplyAndReturnValidatorSetUpdates in EndBlocker
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	require.Len(t, updates, 1)
	require.Equal(t, updates[0].Power, int64(100)) // see MsgCreateValidator stake
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 2)
	newValidator, found := dsc.ValidatorKeeper.GetValidator(ctx, vals[0])
	require.True(t, found)
	require.Equal(t, newValidator.ConsensusPower()+genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
	balanceNB = dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress())
	require.True(t, balanceNB.IsZero())
	balanceB = dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetBondedPool(ctx).GetAddress())
	require.True(t, balanceB.IsEqual(startBalanceB.Add(creatorStake)))

	////////////////////////////////////////////////
	// 4. increment block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(time.Second * 5))

	////////////////////////////////////////////////
	// 5. set second validator offline
	msgOffline := types.NewMsgSetOffline(vals[0])
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = msgsrv.SetOffline(goCtx, msgOffline)
	require.NoError(t, err)
	// last validator must be changes after ApplyAndReturnValidatorSetUpdates in EndBlocker
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 2)

	updates = keeper.EndBlocker(ctx, dsc.ValidatorKeeper, abci.RequestEndBlock{})

	require.Len(t, updates, 1)
	require.Equal(t, updates[0].Power, int64(0)) // 0 mean 'remove from validators'
	require.Len(t, dsc.ValidatorKeeper.GetLastValidators(ctx), 1)
	require.Equal(t, genesisVal.ConsensusPower(), dsc.ValidatorKeeper.GetLastTotalPower(ctx).Int64())
	balanceNB = dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress())
	require.True(t, balanceNB.IsEqual(sdk.NewCoins(creatorStake)))
	balanceB = dsc.BankKeeper.GetAllBalances(ctx, dsc.ValidatorKeeper.GetBondedPool(ctx).GetAddress())
	require.True(t, balanceB.IsEqual(startBalanceB))
}
