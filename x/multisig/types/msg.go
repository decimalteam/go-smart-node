package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
)

var (
	_ sdk.Msg = &MsgCreateWallet{}
	_ sdk.Msg = &MsgCreateTransaction{}
	_ sdk.Msg = &MsgSignTransaction{}
)

const (
	TypeMsgCreateWallet      = "create_wallet"
	TypeMsgCreateTransaction = "create_transaction"
	TypeMsgSignTransaction   = "sign_transaction"
)

////////////////////////////////////////////////////////////////
// MsgCreateWallet
////////////////////////////////////////////////////////////////

func NewMsgCreateWallet(
	sender sdk.AccAddress,
	owners []string,
	weights []uint32,
	threshold uint32,
) *MsgCreateWallet {
	return &MsgCreateWallet{
		Sender:    sender.String(),
		Owners:    owners,
		Weights:   weights,
		Threshold: threshold,
	}
}

// Route should return the name of the module.
func (msg MsgCreateWallet) Route() string { return RouterKey }

// Type should return the action.
func (msg MsgCreateWallet) Type() string { return TypeMsgCreateWallet }

// GetSignBytes encodes the message for signing.
func (msg *MsgCreateWallet) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgCreateWallet) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgCreateWallet) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSender
	}
	// Validate owner count
	if len(msg.Owners) < MinOwnerCount || len(msg.Owners) > MaxOwnerCount {
		return errors.InvalidOwnerCount
	}
	// Validate weight count
	if len(msg.Owners) != len(msg.Weights) {
		return errors.InvalidWeightCount
	}
	// Validate owners (ensure there are no duplicates)
	owners := make(map[string]bool, len(msg.Owners))
	for i := 0; i < len(msg.Owners); i++ {
		if _, err := sdk.AccAddressFromBech32(msg.Owners[i]); err != nil {
			return errors.InvalidOwner
		}
		if owners[msg.Owners[i]] {
			return errors.DuplicateOwner
		}
		owners[msg.Owners[i]] = true
	}
	// Validate weights
	var sumOfWeights uint32
	for i := 0; i < len(msg.Weights); i++ {
		if msg.Weights[i] < MinWeight || msg.Weights[i] > MaxWeight {
			return errors.InvalidWeight
		}
		sumOfWeights += msg.Weights[i]
	}
	if sumOfWeights < msg.Threshold {
		return errors.InvalidThreshold
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgCreateTransaction
////////////////////////////////////////////////////////////////

// NewMsgCreateTransaction creates a new MsgCreateTransaction instance.
func NewMsgCreateTransaction(
	sender sdk.AccAddress,
	wallet string,
	content sdk.Msg,
) (*MsgCreateTransaction, error) {
	anys, err := sdktx.SetMsgs([]sdk.Msg{content})
	if err != nil {
		return nil, err
	}
	return &MsgCreateTransaction{
		Sender:  sender.String(),
		Wallet:  wallet,
		Content: anys[0],
	}, nil
}

// Route returns name of the route for the message.
func (msg *MsgCreateTransaction) Route() string { return RouterKey }

// Type returns the name of the type for the message.
func (msg *MsgCreateTransaction) Type() string { return TypeMsgCreateTransaction }

// GetSignBytes returns the canonical byte representation of the message used to generate a signature.
func (msg *MsgCreateTransaction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners returns the list of signers required to sign the message.
func (msg *MsgCreateTransaction) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic performs basic validation of the message.
func (msg *MsgCreateTransaction) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSender
	}
	if _, err := sdk.AccAddressFromBech32(msg.Wallet); err != nil {
		return errors.InvalidWallet
	}

	return nil
}

////////////////////////////////////////////////////////////////
// MsgSignTransaction
////////////////////////////////////////////////////////////////

func NewMsgSignTransaction(
	sender sdk.AccAddress,
	txID string,
) *MsgSignTransaction {
	return &MsgSignTransaction{
		Sender: sender.String(),
		ID:     txID,
	}
}

// Route returns name of the route for the message.
func (msg *MsgSignTransaction) Route() string { return RouterKey }

// Type returns the name of the type for the message.
func (msg *MsgSignTransaction) Type() string { return TypeMsgSignTransaction }

// GetSignBytes returns the canonical byte representation of the message used to generate a signature.
func (msg *MsgSignTransaction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners returns the list of signers required to sign the message.
func (msg *MsgSignTransaction) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic performs basic validation of the message.
func (msg *MsgSignTransaction) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSender
	}
	prefix, _, err := bech32.DecodeAndConvert(msg.ID)
	if err != nil {
		return errors.InvalidTransactionIDError
	}
	if prefix != MultisigTransactionIDPrefix {
		return errors.InvalidTransactionIDPrefix
	}
	// TODO: TxID length
	return nil
}
