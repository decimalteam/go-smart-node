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

	"bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
)

// TxConstruct is used in process of building, signing and sending transactions
type TxConstructor struct {
	config  client.TxConfig
	builder client.TxBuilder
}

// BuildTransaction creates transaction builder without
func BuildTransaction(msgs []sdk.Msg, memo string, feeDenom string, maxGas uint64) (*TxConstructor, error) {
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
	fee := calculateFee(msgs, memo, maxGas)
	txBuilder.SetGasLimit(fee.Uint64())
	// additional parameters may be usable
	// txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, fee)))
	// txBuilder.SetTimeoutHeight(f.TimeoutHeight())

	return &TxConstructor{txConfig, txBuilder}, nil
}

func calculateFee(msgs []sdk.Msg, memo string, maxGas uint64) sdk.Int {
	return sdk.NewInt(int64(maxGas))
}

// SignTransaction signs transaction and appends signature to transaction signatures.
func (constructor *TxConstructor) SignTransaction(acc *wallet.Account) (*TxConstructor, error) {
	const signMode = signingTypes.SignMode_SIGN_MODE_DIRECT
	// Check chain ID, account number and sequence
	if acc.ChainID() == "" {
		return nil, fmt.Errorf("chain ID is not set up")
	}
	// TODO
	//if acc.accountNumber == 0 || acc.sequence == 0 {
	//	return tx, errors.New("account number or sequence is not set up")
	//}

	// save signatures
	var prevSignatures []signing.SignatureV2
	prevSignatures, err := constructor.builder.GetTx().GetSignaturesV2()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// Generate the bytes to be signed.
	bytesToSign, err := constructor.config.SignModeHandler().GetSignBytes(signMode, signerData, constructor.builder.GetTx())
	if err != nil {
		return nil, err
	}

	// Sign those bytes
	sigBytes, err := acc.Sign(bytesToSign)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return constructor, nil
}

// BytesToSend return binary encoded transaction
func (constructor *TxConstructor) BytesToSend() ([]byte, error) {
	return constructor.config.TxEncoder()(constructor.builder.GetTx())
}
