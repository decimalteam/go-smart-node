package tests

import (
	"context"
	"fmt"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	etherminthd "github.com/decimalteam/ethermint/crypto/hd"
	"github.com/decimalteam/ethermint/encoding"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// helper set up functions
func setUpCliTest(t *testing.T, accCount int) (client.Context, []keyring.Record, *cliTestResulter) {
	var err error

	// make codec and context
	encCfg := encoding.MakeConfig(app.ModuleBasics)
	kb, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendMemory, "",
		nil, encCfg.Codec, []keyring.Option{etherminthd.EthSecp256k1Option()}...)
	require.NoError(t, err)
	clientCtx := client.Context{}.WithKeyring(kb).WithCodec(encCfg.Codec)
	clientCtx = clientCtx.WithGenerateOnly(true)
	clientCtx = clientCtx.WithSignModeStr(flags.SignModeDirect)
	result := &cliTestResulter{}
	clientCtx = clientCtx.WithTxConfig(result)

	// add accounts
	accs := make([]keyring.Record, 0)
	for i := 0; i < accCount; i++ {
		info, _, err := kb.NewMnemonic(fmt.Sprintf("acc%d", i), keyring.English, sdk.FullFundraiserPath,
			keyring.DefaultBIP39Passphrase, etherminthd.EthSecp256k1)
		require.NoError(t, err)
		accs = append(accs, *info)
	}

	return clientCtx, accs, result
}

func setUpCmd(t *testing.T, cmd *cobra.Command, clientCtx client.Context, from string) context.Context {
	var err error

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)
	//err = cmd.Flags().Set(flags.FlagChainID, "decimal_202020-1")
	//require.NoError(t, err)
	err = cmd.Flags().Set(flags.FlagOffline, "true")
	require.NoError(t, err)
	err = cmd.Flags().Set(flags.FlagFrom, from)
	require.NoError(t, err)
	err = cmd.Flags().Set(flags.FlagAccountNumber, "1")
	require.NoError(t, err)
	err = cmd.Flags().Set(flags.FlagSequence, "1")
	require.NoError(t, err)

	return ctx
}

// Helper works as TxConfig+TxBuilder and
// accumulates transaction messages
type cliTestResulter struct {
	msgs []sdk.Msg
}

// TxConfig
func (c *cliTestResulter) TxEncoder() sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return nil, nil
	}
}
func (c *cliTestResulter) TxDecoder() sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		return &sdkTx.Tx{}, nil
	}
}
func (c *cliTestResulter) TxJSONEncoder() sdk.TxEncoder {
	return c.TxEncoder()
}
func (c *cliTestResulter) TxJSONDecoder() sdk.TxDecoder {
	return c.TxDecoder()
}
func (c *cliTestResulter) MarshalSignatureJSON(sign []signingtypes.SignatureV2) ([]byte, error) {
	return nil, nil
}
func (c *cliTestResulter) UnmarshalSignatureJSON([]byte) ([]signingtypes.SignatureV2, error) {
	return []signingtypes.SignatureV2{}, nil
}
func (c *cliTestResulter) NewTxBuilder() client.TxBuilder {
	return c
}
func (c *cliTestResulter) WrapTxBuilder(sdk.Tx) (client.TxBuilder, error) {
	return c, nil
}
func (c *cliTestResulter) SignModeHandler() signing.SignModeHandler {
	return nil
}

// TxBuilder
func (c *cliTestResulter) GetTx() signing.Tx {
	return nil
}
func (c *cliTestResulter) SetMsgs(msgs ...sdk.Msg) error {
	c.msgs = append(c.msgs, msgs...)
	return nil
}
func (c *cliTestResulter) SetSignatures(signatures ...signingtypes.SignatureV2) error {
	return nil
}
func (c *cliTestResulter) AddAuxSignerData(data sdkTx.AuxSignerData) error {
	return nil
}
func (c *cliTestResulter) SetFeePayer(feePayer sdk.AccAddress) {
}
func (c *cliTestResulter) SetTip(tip *sdkTx.Tip) {
}

func (c *cliTestResulter) SetMemo(memo string)                     {}
func (c *cliTestResulter) SetFeeAmount(amount sdk.Coins)           {}
func (c *cliTestResulter) SetGasLimit(limit uint64)                {}
func (c *cliTestResulter) SetTimeoutHeight(height uint64)          {}
func (c *cliTestResulter) SetFeeGranter(feeGranter sdk.AccAddress) {}
