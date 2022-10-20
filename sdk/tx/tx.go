package tx

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"

	appAnte "bitbucket.org/decimalteam/go-smart-node/app/ante"
	"bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
)

// TxConstruct is used in process of building, signing and sending transactions
type TxConstructor struct {
	config  client.TxConfig
	builder client.TxBuilder
}

// BuildTransaction creates transaction builder with automatic fee calculation
// if delPrice is zero, fee amount will be set to zero - this mean that
// DSC node will calculate fee during transaction execution
func BuildTransaction(acc *wallet.Account, msgs []sdk.Msg, memo string, feeDenom string, opts *FeeCalculationOptions) (*TxConstructor, error) {
	if opts == nil {
		return nil, fmt.Errorf("nil opts")
	}
	txc, err := newTxConstructor(msgs, memo)
	if err != nil {
		return nil, err
	}
	oldFee := sdk.ZeroInt()
	newFee := sdk.OneInt()
	if opts.DelPrice.IsZero() {
		newFee = sdk.ZeroInt()
	} else {
		for !oldFee.Equal(newFee) {
			oldFee = sdk.ZeroInt().Add(newFee) // = copy, sdkmath.Int is reference type
			newFee, err = calculateFee(acc, msgs, memo, feeDenom, oldFee, opts)
			if err != nil {
				return nil, err
			}
		}
	}
	txc.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, newFee)))
	// additional parameters may be usable
	// txBuilder.SetTimeoutHeight(f.TimeoutHeight())

	return txc, nil
}

func calculateFee(acc *wallet.Account, msgs []sdk.Msg, memo string, feeDenom string, fee sdkmath.Int, opts *FeeCalculationOptions) (sdkmath.Int, error) {
	txc, err := newTxConstructor(msgs, memo)
	if err != nil {
		// with zero fee, decimal node will calculate correct fee itself
		return sdk.ZeroInt(), err
	}
	txc.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, fee)))
	err = txc.SignTransaction(acc)
	if err != nil {
		// with zero fee, decimal node will calculate correct fee itself
		return sdk.ZeroInt(), err
	}
	bz, err := txc.BytesToSend()
	if err != nil {
		// with zero fee, decimal node will calculate correct fee itself
		return sdk.ZeroInt(), err
	}
	// TODO: in future need to get feetypes.Param by api query
	newFee, err := appAnte.CalculateFee(opts.AppCodec, msgs, int64(len(bz)), opts.DelPrice, opts.FeeParams)
	if err != nil {
		// with zero fee, decimal node will calculate correct fee itself
		return sdk.ZeroInt(), err
	}
	return newFee, nil
}

func newTxConstructor(msgs []sdk.Msg, memo string) (*TxConstructor, error) {
	// 1. create TxBuilder
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authTx.NewTxConfig(marshaler, authTx.DefaultSignModes)
	txBuilder := txConfig.NewTxBuilder()
	// 2. set transaction info
	if err := txBuilder.SetMsgs(msgs...); err != nil {
		return nil, err
	}
	txBuilder.SetMemo(memo)
	return &TxConstructor{txConfig, txBuilder}, nil
}

// SignTransaction signs transaction and appends signature to transaction signatures.
func (constructor *TxConstructor) SetFeeAmount(coins sdk.Coins) {
	constructor.builder.SetFeeAmount(coins)
}

// SignTransaction signs transaction and appends signature to transaction signatures.
func (constructor *TxConstructor) SignTransaction(acc *wallet.Account) error {
	const signMode = signing.SignMode_SIGN_MODE_DIRECT
	// Check chain ID, account number and sequence
	if acc.ChainID() == "" {
		return fmt.Errorf("chain ID is not set up")
	}
	// TODO
	// if acc.accountNumber == 0 || acc.sequence == 0 {
	//	return tx, errors.New("account number or sequence is not set up")
	// }

	// save signatures
	var prevSignatures []signing.SignatureV2
	prevSignatures, err := constructor.builder.GetTx().GetSignaturesV2()
	if err != nil {
		return err
	}

	// 3. signing
	// signerData need to get bytesToSign
	signerData := authsigning.SignerData{
		ChainID:       acc.ChainID(),
		AccountNumber: acc.AccountNumber(),
		Sequence:      acc.Sequence(),
	}
	// sig need for builder
	sig := signing.SignatureV2{
		PubKey: acc.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signMode,
			Signature: nil,
		},
		Sequence: acc.Sequence(),
	}

	if err = constructor.builder.SetSignatures(sig); err != nil {
		return err
	}

	// Generate the bytes to be signed.
	bytesToSign, err := constructor.config.SignModeHandler().GetSignBytes(signMode, signerData, constructor.builder.GetTx())
	if err != nil {
		return err
	}

	// Sign those bytes
	sigBytes, err := acc.Sign(bytesToSign)
	if err != nil {
		return err
	}

	// Construct final SignatureV2 struct
	sig = signing.SignatureV2{
		PubKey: acc.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signMode,
			Signature: sigBytes,
		},
		Sequence: acc.Sequence(),
	}

	prevSignatures = append(prevSignatures, sig)
	if err = constructor.builder.SetSignatures(prevSignatures...); err != nil {
		return err
	}

	return nil
}

// BytesToSend return binary encoded transaction
func (constructor *TxConstructor) BytesToSend() ([]byte, error) {
	return constructor.config.TxEncoder()(constructor.builder.GetTx())
}
