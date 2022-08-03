package tx

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	signingTypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authSigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
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
func BuildTransaction(acc *wallet.Account, msgs []sdk.Msg, memo string, feeDenom string) (*TxConstructor, error) {
	txc, err := newTxConstructor(msgs, memo)
	if err != nil {
		return nil, err
	}
	oldFee := sdk.ZeroInt()
	newFee := sdk.OneInt()
	for !oldFee.Equal(newFee) {
		oldFee = sdk.ZeroInt().Add(newFee) //=copy, sdk.Int is reference type
		newFee, err = calculateFee(acc, msgs, memo, feeDenom, oldFee)
		if err != nil {
			return nil, err
		}
	}
	txc.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, newFee)))
	// additional parameters may be usable
	// txBuilder.SetTimeoutHeight(f.TimeoutHeight())

	return txc, nil
}

func calculateFee(acc *wallet.Account, msgs []sdk.Msg, memo string, feeDenom string, fee sdk.Int) (sdk.Int, error) {
	txc, err := newTxConstructor(msgs, memo)
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
	newFee, err := appAnte.CalculateFee(msgs, int64(len(bz)), sdk.OneDec())
	if err != nil {
		// with zero fee, decimal node will calculate correct fee itself
		return sdk.ZeroInt(), err
	}
	return newFee, nil
}

func newTxConstructor(msgs []sdk.Msg, memo string) (*TxConstructor, error) {
	// 1. create TxBuilder
	interfaceRegistry := codecTypes.NewInterfaceRegistry()
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
	const signMode = signingTypes.SignMode_SIGN_MODE_DIRECT
	// Check chain ID, account number and sequence
	if acc.ChainID() == "" {
		return fmt.Errorf("chain ID is not set up")
	}
	// TODO
	//if acc.accountNumber == 0 || acc.sequence == 0 {
	//	return tx, errors.New("account number or sequence is not set up")
	//}

	// save signatures
	var prevSignatures []signing.SignatureV2
	prevSignatures, err := constructor.builder.GetTx().GetSignaturesV2()
	if err != nil {
		return err
	}

	// 3. signing
	// signerData need to get bytesToSign
	signerData := authSigning.SignerData{
		ChainID:       acc.ChainID(),
		AccountNumber: acc.AccountNumber(),
		Sequence:      acc.Sequence(),
	}
	// sig need for builder
	sig := signingTypes.SignatureV2{
		PubKey: acc.PubKey(),
		Data: &signingTypes.SingleSignatureData{
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
