package upgrade_test

import (
	"fmt"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/testutil"
	"bitbucket.org/decimalteam/go-smart-node/x/upgrade"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	msg := types.MsgSoftwareUpgrade{
		Authority: "d01y7sex8yvrazyd8pljjxvnvpndaavn99tk2j52y",
		Plan: types.Plan{
			Name: "https://devnet-dsc-repo.decimalchain.com/63261",
			//Name:   "http://127.0.0.1:8080/50",
			Height: 63261,
			Info:   `{"linux/centos/7":["89440dac66fcb9e590a9d7bb8752b61f006d46c21a47adc24d67711b8e09b53b"],"linux/centos/8":["9b340bbcc38571d69c9b8c1a37b296c337e2cdf69ce5429f5d4d8a700b63e0c8"],"linux/debian/10":["7292b25d82e0e938eab84a0d6d5b3415cebf87568bc73d87f99bc0c996f1b950"],"linux/debian/11":["3ff725aeee6c59c39954e8b9a391546cae2c57d35b17c10145b094f16c1dab47"],"linux/debian/9":["0755d0bbd9dfa77d7b34181e4d977bee5ea0430c07a60dee9970a403ff9cd93b"],"linux/ol/8.5":["0dca000b55988d73044a3d887f6153ac16268b4ad4863820d0bf8f4a0255c94f"],"linux/ubuntu/16.04":["7fa6659194a04416e2f0d8a644ae859453f68a8800556dc3db36ca02e98baa66"],"linux/ubuntu/18.04":["0f3c0dbc0398381ff6b12c00102967edc5c049bbdeda249f2af551a9da3c4ac4"],"linux/ubuntu/20.04":["814ee010d89f202a3d32279c1cdaf3173699f3227ef8b9d23d4ba421be4396c7"],"linux/ubuntu/22.04":["ae0f16711451fba5201e3af9a394b109ce7447d1ced2ae5be550147399267132"]}`,
			//Info:   `{"darwin":["8f02709b53f523c413391f4970595c098bba62e04b00c5d6c86c6e25b0d1d851"]}`,
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
