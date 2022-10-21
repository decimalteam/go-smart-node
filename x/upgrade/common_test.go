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
			Name: "https://devnet-dsc-repo.decimalchain.com/37201",
			//Name:   "http://127.0.0.1:8080/50",
			Height: 37201,
			Info:   `{"linux/centos/7":["40374b0402ddda0172c0deea7c570bfac13334f75be0cdaf4f165ece1bb4ee38"],"linux/centos/8":["41d21011fecda9f995bc7611d41f1932379f136791152cee7a9d2b23a73df275"],"linux/debian/10":["42c1ec1017097bb52d5b68390d0ff839ef57e24783730256cd315a23d33e443e"],"linux/debian/11":["d50513c0543b53c4e96c1a451a3fedb480f06d25a5570eda363978b94d033bec"],"linux/debian/9":["8631bb9f59d546e94806cfd9d27e0aa7e439ea737741663d4b6f0d21094ce50a"],"linux/ol/8.5":["9911d63e99dbf6d62aa30dec8945afb7212a0f285fcba58c3938f514b07a800f"],"linux/ubuntu/16.04":["c4bd4d02a27b2c76d3039f0076c41f63ca3fec44915a1172476c59a53fea8a03"],"linux/ubuntu/18.04":["4e57519cb75b78f7e53d2fd1ba34d8c6b7b819c87b4aa2a8da27560a69df332e"],"linux/ubuntu/20.04":["67eab5ca824ff92bf5729591a86cb3dd3c1e8d81b1d06e9c95eb6c1471e19cd8"],"linux/ubuntu/22.04":["a2660f3f84a81d93baa82592744cfa6dda287d0d9dac2d2c800cd49fd5c5f206"]}`,
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
