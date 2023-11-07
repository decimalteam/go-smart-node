package multisig_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	feemarkettypes "github.com/decimalteam/ethermint/x/feemarket/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
)

func TestDoubleWallet(t *testing.T) {
	const addrCount = 100

	_, dsc, ctx := getBaseAppWithCustomKeeper(t)
	addrs, _ := generateAddresses(dsc, ctx, addrCount,
		sdk.NewCoins(
			sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1000))),
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
			sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1000))),
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

func TestSignTransaction(t *testing.T) {
	const addrCount = 100

	_, dsc, ctx := getBaseAppWithCustomKeeper(t)
	addrs, _ := generateAddresses(dsc, ctx, addrCount,
		sdk.NewCoins(
			sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1000))),
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
	err = dsc.BankKeeper.SendCoins(ctx, sender, wAdr, sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(10)))))
	require.NoError(t, err)

	// create 'send 1 coin' to receiver
	msgTx, err := types.NewMsgCreateTransaction(addrs[1], walletResponse.Wallet,
		cointypes.NewMsgSendCoin(
			sdk.MustAccAddressFromBech32(walletResponse.Wallet),
			receiver,
			sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1))),
		),
	)
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
	walletBalance := dsc.BankKeeper.GetBalance(ctx, wAdr, cmdcfg.BaseDenom)
	require.True(t, walletBalance.Amount.Equal(helpers.EtherToWei(sdk.NewInt(9))), "wallet balance")
	// 1001 coins
	receiverBalance := dsc.BankKeeper.GetBalance(ctx, receiver, cmdcfg.BaseDenom)
	require.True(t, receiverBalance.Amount.Equal(helpers.EtherToWei(sdk.NewInt(1001))), "receiver balance")
}

func TestTryOverspend(t *testing.T) {
	const addrCount = 100

	_, dsc, ctx := getBaseAppWithCustomKeeper(t)
	addrs, _ := generateAddresses(dsc, ctx, addrCount,
		sdk.NewCoins(
			sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1000))),
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
	err = dsc.BankKeeper.SendCoins(ctx, sender, wAdr, sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(10)))))
	require.NoError(t, err)

	// create 'send 10 coin' to receiver twice
	txIDs := []string{}
	for i := 0; i < 2; i++ {
		ctx = ctx.WithTxBytes([]byte{byte(i)}) // for tx id salt
		msgTx, err := types.NewMsgCreateTransaction(addrs[1], walletResponse.Wallet,
			cointypes.NewMsgSendCoin(
				sdk.MustAccAddressFromBech32(walletResponse.Wallet),
				receiver,
				sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(10))),
			),
		)
		require.NoError(t, err)
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
	walletBalance := dsc.BankKeeper.GetBalance(ctx, wAdr, cmdcfg.BaseDenom)
	require.True(t, walletBalance.Amount.Equal(helpers.EtherToWei(sdk.NewInt(0))), "wallet balance")

	// second
	msgSign = types.NewMsgSignTransaction(addrs[3], txIDs[1])
	err = msgSign.ValidateBasic()
	require.NoError(t, err)
	goCtx = sdk.WrapSDKContext(ctx)
	_, err = dsc.MultisigKeeper.SignTransaction(goCtx, msgSign)
	require.Error(t, err)
}

func TestUniversalTx(t *testing.T) {

	const addrCount = 100

	_, dsc, ctx := getBaseAppWithCustomKeeper(t)
	addrs, _ := generateAddresses(dsc, ctx, addrCount,
		sdk.NewCoins(
			sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1000))),
		),
	)

	var sender = addrs[0]
	var receiver = addrs[10]
	var owners = []string{addrs[1].String(), addrs[2].String(), addrs[3].String(), addrs[4].String()}
	var weights = []uint32{1, 1, 1, 1}
	var threshold uint32 = 3

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
	err = dsc.BankKeeper.SendCoins(ctx, sender, wAdr, sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(10)))))
	require.NoError(t, err)

	// create universal tx
	msgU, err := types.NewMsgCreateTransaction(
		addrs[1], walletResponse.Wallet,
		cointypes.NewMsgSendCoin(wAdr, receiver, sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(10)))),
	)
	require.NoError(t, err)

	txres, err := dsc.MultisigKeeper.CreateTransaction(goCtx, msgU)
	require.NoError(t, err)

	// second owner sign
	msgS := types.NewMsgSignTransaction(addrs[2], txres.ID)
	_, err = dsc.MultisigKeeper.SignTransaction(goCtx, msgS)
	require.NoError(t, err)

	// check for double sign
	_, err = dsc.MultisigKeeper.SignTransaction(goCtx, msgS)
	require.Error(t, err)

	// third owner sign
	msgS = types.NewMsgSignTransaction(addrs[3], txres.ID)
	_, err = dsc.MultisigKeeper.SignTransaction(goCtx, msgS)
	require.NoError(t, err)

	// check internal tx result
	walletBalance := dsc.BankKeeper.GetBalance(ctx, wAdr, cmdcfg.BaseDenom)
	require.True(t, walletBalance.Amount.Equal(helpers.EtherToWei(sdk.NewInt(0))), "wallet balance")
	receiverBalance := dsc.BankKeeper.GetBalance(ctx, receiver, cmdcfg.BaseDenom)
	require.True(t, receiverBalance.Amount.Equal(helpers.EtherToWei(sdk.NewInt(1010))), "receiver balance", receiverBalance.String())
	require.True(t, dsc.MultisigKeeper.IsCompleted(ctx, txres.ID))

	// fourth owner, transaction already completed
	msgS = types.NewMsgSignTransaction(addrs[4], txres.ID)
	_, err = dsc.MultisigKeeper.SignTransaction(goCtx, msgS)
	require.Error(t, err)

}

func TestGenesisTransactions(t *testing.T) {
	_, dsc, ctx := getBaseAppWithCustomKeeper(t)
	walletAdr := sdk.AccAddress([]byte{1, 2, 3}).String()
	adr1 := sdk.AccAddress([]byte{1, 2, 3, 1}).String()
	adr2 := sdk.AccAddress([]byte{1, 2, 3, 2}).String()
	adr3 := sdk.AccAddress([]byte{1, 2, 3, 3}).String()
	txID, err := bech32.ConvertAndEncode(types.MultisigTransactionIDPrefix, []byte{0, 0, 0, 1})
	require.NoError(t, err)

	var genesis = types.GenesisState{
		Wallets: []types.Wallet{
			{
				Address: walletAdr,
				Owners: []string{
					adr1,
					adr2,
					adr3,
				},
				Weights:   []uint32{1, 1, 1},
				Threshold: 2,
			},
		},
		Transactions: []types.GenesisTransaction{
			{
				Id:        txID,
				Wallet:    walletAdr,
				Receiver:  adr3,
				Coins:     sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, sdkmath.NewInt(1))),
				Signers:   []string{adr1},
				CreatedAt: 1,
			},
		},
	}

	multisig.InitGenesis(ctx, dsc.MultisigKeeper, &genesis)

	_, err = dsc.MultisigKeeper.GetTransaction(ctx, txID)
	require.NoError(t, err)
	require.True(t, dsc.MultisigKeeper.IsSigned(ctx, txID, adr1))
	require.False(t, dsc.MultisigKeeper.IsSigned(ctx, txID, adr2))
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
		dsc.MsgServiceRouter(),
	)

	return codec.NewLegacyAmino(), dsc, ctx
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(dsc *app.DSC, ctx sdk.Context, numAddrs int, accCoins sdk.Coins) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := app.AddTestAddrsIncremental(dsc, ctx, numAddrs, accCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}
