package multisig_test

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestDoubleWallet(t *testing.T) {
	const addrCount = 100

	_, dsc, ctx := getBaseAppWithCustomKeeper(t)
	addrs, _ := generateAddresses(dsc, ctx, addrCount,
		sdk.NewCoins(
			sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(1000))),
		),
	)

	var sender = addrs[0]
	var owners = []string{addrs[1].String(), addrs[2].String(), addrs[3].String()}
	var weights = []uint32{1, 1, 1}
	var threshold uint32 = 2

	msg := types.NewMsgCreateWallet(sender, owners, weights, threshold)
	err := msg.ValidateBasic()
	require.NoError(t, err)

	ctx = ctx.WithTxBytes([]byte{byte(1)}) // for wallet salt
	goCtx := sdk.WrapSDKContext(ctx)
	_, err = dsc.MultisigKeeper.CreateWallet(goCtx, msg)
	require.NoError(t, err)

	// try to create wallet again
	_, err = dsc.MultisigKeeper.CreateWallet(goCtx, msg)
	require.Error(t, err)
}

func TestAccountWithSameAddress(t *testing.T) {
	const addrCount = 100

	_, dsc, ctx := getBaseAppWithCustomKeeper(t)
	addrs, _ := generateAddresses(dsc, ctx, addrCount,
		sdk.NewCoins(
			sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(1000))),
		),
	)

	var sender = addrs[0]
	var owners = []string{addrs[1].String(), addrs[2].String(), addrs[3].String()}
	var weights = []uint32{1, 1, 1}
	var threshold uint32 = 2

	wallet, err := types.NewWallet(owners, weights, threshold, []byte{1})
	require.NoError(t, err)

	// some cheating - createing accoutn with address of future wallet
	addr, err := sdk.AccAddressFromBech32(wallet.Address)
	require.NoError(t, err)
	acc := authtypes.NewBaseAccountWithAddress(addr)
	dsc.AccountKeeper.SetAccount(ctx, acc)

	// create wallet
	msg := types.NewMsgCreateWallet(sender, owners, weights, threshold)
	err = msg.ValidateBasic()
	require.NoError(t, err)
	ctx = ctx.WithTxBytes([]byte{byte(1)}) // for wallet salt
	goCtx := sdk.WrapSDKContext(ctx)
	_, err = dsc.MultisigKeeper.CreateWallet(goCtx, msg)
	// address already exists
	require.Error(t, err)
}

func TestLowBalance(t *testing.T) {
	const addrCount = 100

	_, dsc, ctx := getBaseAppWithCustomKeeper(t)
	addrs, _ := generateAddresses(dsc, ctx, addrCount,
		sdk.NewCoins(
			sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(1000))),
		),
	)

	var sender = addrs[0]
	var owners = []string{addrs[1].String(), addrs[2].String(), addrs[3].String()}
	var weights = []uint32{1, 1, 1}
	var threshold uint32 = 2

	// create wallet with empty balance
	msg := types.NewMsgCreateWallet(sender, owners, weights, threshold)
	err := msg.ValidateBasic()
	require.NoError(t, err)
	ctx = ctx.WithTxBytes([]byte{byte(1)}) // for wallet salt
	goCtx := sdk.WrapSDKContext(ctx)
	walletResponse, err := dsc.MultisigKeeper.CreateWallet(goCtx, msg)
	require.NoError(t, err)

	msgTx := types.NewMsgCreateTransaction(addrs[1], walletResponse.Wallet, addrs[10].String(), sdk.NewCoins(sdk.NewCoin("del", sdk.NewInt(1))))
	err = msgTx.ValidateBasic()
	require.NoError(t, err)
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = dsc.MultisigKeeper.CreateTransaction(goCtx, msgTx)
	require.Error(t, err)
}

func TestSenderNotOwner(t *testing.T) {
	const addrCount = 100

	_, dsc, ctx := getBaseAppWithCustomKeeper(t)
	addrs, _ := generateAddresses(dsc, ctx, addrCount,
		sdk.NewCoins(
			sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(1000))),
		),
	)

	var sender = addrs[0]
	var owners = []string{addrs[1].String(), addrs[2].String(), addrs[3].String()}
	var weights = []uint32{1, 1, 1}
	var threshold uint32 = 2

	// create wallet with empty balance
	msg := types.NewMsgCreateWallet(sender, owners, weights, threshold)
	err := msg.ValidateBasic()
	require.NoError(t, err)
	ctx = ctx.WithTxBytes([]byte{byte(1)}) // for wallet salt
	goCtx := sdk.WrapSDKContext(ctx)
	walletResponse, err := dsc.MultisigKeeper.CreateWallet(goCtx, msg)
	require.NoError(t, err)

	// send 10 coins to wallet
	wAdr, err := sdk.AccAddressFromBech32(walletResponse.Wallet)
	require.NoError(t, err)
	err = dsc.BankKeeper.SendCoins(ctx, sender, wAdr, sdk.NewCoins(sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(10)))))
	require.NoError(t, err)

	msgTx := types.NewMsgCreateTransaction(sender, walletResponse.Wallet, addrs[10].String(), sdk.NewCoins(sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(1)))))
	err = msgTx.ValidateBasic()
	require.NoError(t, err)
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = dsc.MultisigKeeper.CreateTransaction(goCtx, msgTx)
	require.Error(t, err)
}

func TestSignTransaction(t *testing.T) {
	const addrCount = 100

	_, dsc, ctx := getBaseAppWithCustomKeeper(t)
	addrs, _ := generateAddresses(dsc, ctx, addrCount,
		sdk.NewCoins(
			sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(1000))),
		),
	)

	var sender = addrs[0]
	var receiver = addrs[10]
	var owners = []string{addrs[1].String(), addrs[2].String(), addrs[3].String()}
	var weights = []uint32{1, 1, 1}
	var threshold uint32 = 2

	// create wallet with empty balance
	msg := types.NewMsgCreateWallet(sender, owners, weights, threshold)
	err := msg.ValidateBasic()
	require.NoError(t, err)
	ctx = ctx.WithTxBytes([]byte{byte(1)}) // for wallet salt
	goCtx := sdk.WrapSDKContext(ctx)
	walletResponse, err := dsc.MultisigKeeper.CreateWallet(goCtx, msg)
	require.NoError(t, err)

	// send 10 coins to wallet
	wAdr, err := sdk.AccAddressFromBech32(walletResponse.Wallet)
	require.NoError(t, err)
	err = dsc.BankKeeper.SendCoins(ctx, sender, wAdr, sdk.NewCoins(sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(10)))))
	require.NoError(t, err)

	// create 'send 1 coin' to receiver
	msgTx := types.NewMsgCreateTransaction(addrs[1], walletResponse.Wallet, receiver.String(),
		sdk.NewCoins(sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(1)))))
	err = msgTx.ValidateBasic()
	require.NoError(t, err)
	goCtx = sdk.WrapSDKContext(ctx)
	txResponse, err := dsc.MultisigKeeper.CreateTransaction(goCtx, msgTx)
	require.NoError(t, err)

	// tx already signed by first owner, time to sign by third
	msgSign := types.NewMsgSignTransaction(addrs[3], txResponse.ID)
	err = msgSign.ValidateBasic()
	require.NoError(t, err)
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = dsc.MultisigKeeper.SignTransaction(goCtx, msgSign)
	require.NoError(t, err)

	// check balance
	// 9 coins
	walletBalance := dsc.BankKeeper.GetBalance(ctx, wAdr, "del")
	require.True(t, walletBalance.Amount.Equal(helpers.EtherToWei(sdk.NewInt(9))), "wallet balance")
	// 1001 coins
	receiverBalance := dsc.BankKeeper.GetBalance(ctx, receiver, "del")
	require.True(t, receiverBalance.Amount.Equal(helpers.EtherToWei(sdk.NewInt(1001))), "receiver balance")
}

func TestTryOverspend(t *testing.T) {
	const addrCount = 100

	_, dsc, ctx := getBaseAppWithCustomKeeper(t)
	addrs, _ := generateAddresses(dsc, ctx, addrCount,
		sdk.NewCoins(
			sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(1000))),
		),
	)

	var sender = addrs[0]
	var receiver = addrs[10]
	var owners = []string{addrs[1].String(), addrs[2].String(), addrs[3].String()}
	var weights = []uint32{1, 1, 1}
	var threshold uint32 = 2

	// create wallet with empty balance
	msg := types.NewMsgCreateWallet(sender, owners, weights, threshold)
	err := msg.ValidateBasic()
	require.NoError(t, err)
	ctx = ctx.WithTxBytes([]byte{byte(1)}) // for wallet salt
	goCtx := sdk.WrapSDKContext(ctx)
	walletResponse, err := dsc.MultisigKeeper.CreateWallet(goCtx, msg)
	require.NoError(t, err)

	// send 10 coins to wallet
	wAdr, err := sdk.AccAddressFromBech32(walletResponse.Wallet)
	require.NoError(t, err)
	err = dsc.BankKeeper.SendCoins(ctx, sender, wAdr, sdk.NewCoins(sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(10)))))
	require.NoError(t, err)

	// create 'send 10 coin' to receiver twice
	txIDs := []string{}
	for i := 0; i < 2; i++ {
		ctx = ctx.WithTxBytes([]byte{byte(i)}) // for tx id salt
		msgTx := types.NewMsgCreateTransaction(addrs[1], walletResponse.Wallet, receiver.String(),
			sdk.NewCoins(sdk.NewCoin("del", helpers.EtherToWei(sdk.NewInt(10)))))
		err = msgTx.ValidateBasic()
		require.NoError(t, err)
		goCtx = sdk.WrapSDKContext(ctx)
		txResponse, err := dsc.MultisigKeeper.CreateTransaction(goCtx, msgTx)
		require.NoError(t, err)
		txIDs = append(txIDs, txResponse.ID)
	}

	// now we have 2 transactions
	// second tx will overspend

	// first
	msgSign := types.NewMsgSignTransaction(addrs[3], txIDs[0])
	err = msgSign.ValidateBasic()
	require.NoError(t, err)
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = dsc.MultisigKeeper.SignTransaction(goCtx, msgSign)
	require.NoError(t, err)
	// 0 coins
	walletBalance := dsc.BankKeeper.GetBalance(ctx, wAdr, "del")
	require.True(t, walletBalance.Amount.Equal(helpers.EtherToWei(sdk.NewInt(0))), "wallet balance")

	// second
	msgSign = types.NewMsgSignTransaction(addrs[3], txIDs[1])
	err = msgSign.ValidateBasic()
	require.NoError(t, err)
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = dsc.MultisigKeeper.SignTransaction(goCtx, msgSign)
	require.Error(t, err)
}

// getBaseAppWithCustomKeeper Returns a simapp with custom keepers
// to avoid messing with the hooks.
func getBaseAppWithCustomKeeper(t *testing.T) (*codec.LegacyAmino, *app.DSC, sdk.Context) {
	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	appCodec := dsc.AppCodec()

	dsc.MultisigKeeper = *keeper.NewKeeper(
		appCodec,
		dsc.GetKey(types.StoreKey),
		dsc.GetSubspace(types.ModuleName),
		dsc.AccountKeeper,
		dsc.BankKeeper,
	)
	//dsc.MultisigKeeper.SetParams(ctx, types.DefaultParams())

	return codec.NewLegacyAmino(), dsc, ctx
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(dsc *app.DSC, ctx sdk.Context, numAddrs int, accCoins sdk.Coins) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := app.AddTestAddrsIncremental(dsc, ctx, numAddrs, accCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}
