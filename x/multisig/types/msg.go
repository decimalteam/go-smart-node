package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
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

// Multisignature wallet limitations.
const (
	MinOwnerCount = 2
	MaxOwnerCount = 16
	MinWeight     = 1
	MaxWeight     = 1024
)

////////////////////////////////////////////////////////////////
// MsgCreateWallet
////////////////////////////////////////////////////////////////

func NewMsgCreateWallet(
	sender sdk.AccAddress,
	owners []string,
	weights []uint64,
	threshold uint64,
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
		return ErrInvalidSender(msg.Sender)
	}
	// Validate owner count
	if len(msg.Owners) < MinOwnerCount || len(msg.Owners) > MaxOwnerCount {
		return ErrInvalidOwnerCount(strconv.Itoa(len(msg.Owners)), strconv.Itoa(MinOwnerCount), strconv.Itoa(MaxOwnerCount))
	}
	// Validate weight count
	if len(msg.Owners) != len(msg.Weights) {
		return ErrInvalidWeightCount(strconv.Itoa(len(msg.Weights)), strconv.Itoa(len(msg.Owners)))
	}
	// Validate owners (ensure there are no duplicates)
	owners := make(map[string]bool, len(msg.Owners))
	for i, c := 0, len(msg.Owners); i < c; i++ {
		if _, err := sdk.AccAddressFromBech32(msg.Owners[i]); err != nil {
			return ErrInvalidOwner(msg.Owners[i])
		}
		if owners[msg.Owners[i]] {
			return ErrDuplicateOwner(msg.Owners[i])
		}
		owners[msg.Owners[i]] = true
	}
	// Validate weights
	for i, c := 0, len(msg.Weights); i < c; i++ {
		if msg.Weights[i] < MinWeight {
			return ErrInvalidWeight(strconv.Itoa(MinWeight), "less")
		}
		if msg.Weights[i] > MaxWeight {
			return ErrInvalidWeight(strconv.Itoa(MaxWeight), "greater")
		}
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
	receiver string,
	coins sdk.Coins,
) *MsgCreateTransaction {
	return &MsgCreateTransaction{
		Sender:   sender.String(),
		Wallet:   wallet,
		Receiver: receiver,
		Coins:    coins,
	}
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
		return ErrInvalidSender(msg.Sender)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Wallet); err != nil {
		return ErrInvalidWallet(msg.Wallet)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Receiver); err != nil {
		return ErrInvalidReceiver(msg.Receiver)
	}
	// Amount should be positive
	for _, coin := range msg.Coins {
		if coin.Amount.LT(sdk.NewInt(0)) {
			return ErrInvalidAmount(coin.Denom, coin.Amount.String())
		}
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
		TxID:   txID,
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
		return ErrInvalidSender(msg.Sender)
	}
	prefix, _, err := bech32.DecodeAndConvert(msg.TxID)
	if err != nil {
		return ErrInvalidTransactionIDError(msg.TxID, err.Error())
	}
	if prefix != MultisigTransactionIDPrefix {
		return ErrInvalidTransactionIDPrefix(msg.TxID, MultisigTransactionIDPrefix, prefix)
	}
	// TODO: TxID length
	return nil
}
