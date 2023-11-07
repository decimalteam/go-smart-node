package types

import (
	commonTypes "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/decimalteam/ethermint/crypto/ethsecp256k1"
)

var (
	_ sdk.Msg = &MsgReturnLegacy{}
)

const (
	TypeMsgReturnLegacy = "return_legacy"
)

////////////////////////////////////////////////////////////////
// MsgReturnLegacy
////////////////////////////////////////////////////////////////

// NewMsgReturnLegacy creates a new instance of MsgReturnLegacy.

func NewMsgReturnLegacy(
	sender sdk.AccAddress,
	publicKeyBytes []byte,
) *MsgReturnLegacy {
	return &MsgReturnLegacy{
		Sender:    sender.String(),
		PublicKey: publicKeyBytes,
	}

}

// Route should return the name of the module.
func (msg MsgReturnLegacy) Route() string { return RouterKey }

// Type should return the action.
func (msg MsgReturnLegacy) Type() string { return TypeMsgReturnLegacy }

// GetSignBytes encodes the message for signing.
func (msg *MsgReturnLegacy) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgReturnLegacy) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.

func (msg MsgReturnLegacy) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}
	// Validate public key
	if len(msg.PublicKey) != ethsecp256k1.PubKeySize {
		return errors.InvalidPublicKeyLength
	}
	// Validate address from public key
	if _, err := commonTypes.GetLegacyAddressFromPubKey(msg.PublicKey); err != nil {
		return errors.CannotGetLegacyAddressFromPublicKey
	}
	return nil
}
