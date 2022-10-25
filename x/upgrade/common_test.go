package upgrade_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/testutil"
	"bitbucket.org/decimalteam/go-smart-node/x/upgrade"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestName(t *testing.T) {
	msg := types.MsgSoftwareUpgrade{
		Authority: "dx1y7sex8yvrazyd8pljjxvnvpndaavn99tjd3ppm",
		Plan: types.Plan{
			//Name: "https://devnet-dsc-repo.decimalchain.com/50701",
			Name:   "http://127.0.0.1:8080/50",
			Height: 50,
			Info:   `{"darwin":["e277359d00df49bd9284c896f4de8826c4e298296cdc9623705bffc35186ab6b"]}`,
		},
	}

	encCfg := testutil.MakeTestEncodingConfig(upgrade.AppModule{})

	kb, err := keyring.New(t.Name(), "test", t.TempDir(), nil, encCfg.Codec)
	require.NoError(t, err)

	path := hd.CreateHDPath(60, 0, 0).String()

	_, _, err = kb.NewMnemonic("test_key1", keyring.English, path, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	require.NoError(t, err)

	txf := tx.Factory{}.
		WithTxConfig(encCfg.TxConfig).
		//WithAccountNumber(8).
		//WithSequence(0).
		//WithFees("26320000000000000000000del").
		//WithMemo("memo").
		WithChainID("decimal_2020-22100500").
		WithKeybase(kb)

	tx, err := txf.BuildUnsignedTx(&msg)
	require.NoError(t, err)
	require.NotNil(t, tx)

	json, err := encCfg.TxConfig.TxJSONEncoder()(tx.GetTx())
	require.NoError(t, err)

	t.Log(fmt.Sprintf("%s\n", json))
}

// getBaseAppWithCustomKeeper Returns a simapp with custom CoinKeeper
// to avoid messing with the hooks.
//func getBaseAppWithCustomKeeper(skip map[int64]bool) (*codec.LegacyAmino, *app.DSC, sdk.Context) {
//	dsc := app.Setup(false, feemarkettypes.DefaultGenesisState())
//	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})
//
//	appCodec := dsc.AppCodec()
//
//	if skip == nil {
//		skip = make(map[int64]bool)
//	}
//	dsc.UpgradeKeeper = keeper.NewKeeper(
//		skip,
//		dsc.GetKey(types.StoreKey),
//		appCodec,
//		app.DefaultNodeHome,
//		dsc.BaseApp,
//		app.UpgraderAddress,
//	)
//
//	return codec.NewLegacyAmino(), dsc, ctx
//}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(dsc *app.DSC, ctx sdk.Context, numAddrs int, accCoins sdk.Coins) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := app.AddTestAddrsIncremental(dsc, ctx, numAddrs, accCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}
