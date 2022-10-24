package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/nft/config"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/errors"
)

var (
	_ sdk.Msg = &MsgMintToken{}
	_ sdk.Msg = &MsgUpdateToken{}
	_ sdk.Msg = &MsgUpdateReserve{}
	_ sdk.Msg = &MsgSendToken{}
	_ sdk.Msg = &MsgBurnToken{}
)

const (
	TypeMsgMintToken     = "mint_token"
	TypeMsgUpdateToken   = "update_token"
	TypeMsgUpdateReserve = "update_reserve"
	TypeMsgSendToken     = "send_token"
	TypeMsgBurnToken     = "burn_token"
)

////////////////////////////////////////////////////////////////
// MsgMintToken
////////////////////////////////////////////////////////////////

// NewMsgMintToken creates a new instance of MsgMintToken.
func NewMsgMintToken(
	sender sdk.AccAddress,
	denom string,
	tokenID string,
	tokenURI string,
	allowMint bool,
	recipient sdk.AccAddress,
	quantity uint32,
	reserve sdk.Coin,
) *MsgMintToken {
	return &MsgMintToken{
		Sender:    sender.String(),
		Denom:     strings.TrimSpace(denom),
		TokenID:   strings.TrimSpace(tokenID),
		TokenURI:  strings.TrimSpace(tokenURI),
		AllowMint: allowMint,
		Recipient: recipient.String(),
		Quantity:  quantity,
		Reserve:   reserve,
	}
}

// Route should return the name of the module.
func (msg *MsgMintToken) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgMintToken) Type() string { return TypeMsgMintToken }

// GetSignBytes encodes the message for signing.
func (msg *MsgMintToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgMintToken) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgMintToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.InvalidSender
	}
	_, err = sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return errors.InvalidRecipientAddress
	}
	if strings.TrimSpace(msg.Denom) == "" {
		return errors.InvalidDenom
	}
	if strings.TrimSpace(msg.TokenID) == "" {
		return errors.InvalidNFT
	}
	if msg.Quantity < 1 {
		return errors.InvalidQuantity
	}
	if !msg.Reserve.IsPositive() {
		return errors.InvalidReserve
	}
	if !config.CollectionDenomValidator.MatchString(msg.Denom) {
		return errors.InvalidDenom
	}
	if !config.CollectionDenomValidator.MatchString(msg.TokenID) {
		return errors.InvalidTokenID
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgUpdateToken
////////////////////////////////////////////////////////////////

// NewMsgUpdateToken creates a new instance of MsgUpdateToken.
func NewMsgUpdateToken(sender sdk.AccAddress, tokenID string, tokenURI string) *MsgUpdateToken {
	return &MsgUpdateToken{
		Sender:   sender.String(),
		TokenID:  strings.TrimSpace(tokenID),
		TokenURI: strings.TrimSpace(tokenURI),
	}
}

// Route should return the name of the module.
func (msg *MsgUpdateToken) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgUpdateToken) Type() string { return TypeMsgUpdateToken }

// GetSignBytes encodes the message for signing.
func (msg *MsgUpdateToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgUpdateToken) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgUpdateToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.InvalidSender
	}
	if strings.TrimSpace(msg.TokenID) == "" {
		return errors.InvalidNFT
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgUpdateReserve
////////////////////////////////////////////////////////////////

// NewMsgUpdateReserve creates a new instance of MsgUpdateReserve.
func NewMsgUpdateReserve(sender sdk.AccAddress, tokenID string, subTokenIDs []uint32, reserve sdk.Coin) *MsgUpdateReserve {
	return &MsgUpdateReserve{
		Sender:      sender.String(),
		TokenID:     strings.TrimSpace(tokenID),
		SubTokenIDs: subTokenIDs,
		Reserve:     reserve,
	}
}

// Route should return the name of the module.
func (msg *MsgUpdateReserve) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgUpdateReserve) Type() string { return TypeMsgUpdateReserve }

// GetSignBytes encodes the message for signing.
func (msg *MsgUpdateReserve) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgUpdateReserve) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgUpdateReserve) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.InvalidSender
	}
	if strings.TrimSpace(msg.TokenID) == "" {
		return errors.InvalidNFT
	}
	if !CheckUnique(msg.SubTokenIDs) {
		return errors.NotUniqueSubTokenIDs
	}
	if msg.Reserve.IsZero() {
		return errors.InvalidReserve
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgSendToken
////////////////////////////////////////////////////////////////

// NewMsgSendToken creates a new instance of MsgSendToken.
func NewMsgSendToken(sender sdk.AccAddress, recipient sdk.AccAddress, tokenID string, subTokenIDs []uint32) *MsgSendToken {
	return &MsgSendToken{
		Sender:      sender.String(),
		TokenID:     strings.TrimSpace(tokenID),
		SubTokenIDs: subTokenIDs,
		Recipient:   recipient.String(),
	}
}

// Route should return the name of the module.
func (msg *MsgSendToken) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgSendToken) Type() string { return TypeMsgSendToken }

// GetSignBytes encodes the message for signing.
func (msg *MsgSendToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgSendToken) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgSendToken) ValidateBasic() error {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.InvalidSender
	}
	if strings.TrimSpace(msg.TokenID) == "" {
		return errors.InvalidCollection
	}
	if !CheckUnique(msg.SubTokenIDs) {
		return errors.NotUniqueSubTokenIDs
	}
	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return errors.InvalidRecipientAddress
	}
	if sender.Equals(recipient) {
		return errors.ForbiddenToTransferToYourself
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgBurnToken
////////////////////////////////////////////////////////////////

// NewMsgBurnToken creates a new instance of MsgBurnToken.
func NewMsgBurnToken(sender sdk.AccAddress, tokenID string, subTokenIDs []uint32) *MsgBurnToken {
	return &MsgBurnToken{
		Sender:      sender.String(),
		TokenID:     strings.TrimSpace(tokenID),
		SubTokenIDs: subTokenIDs,
	}
}

// Route should return the name of the module.
func (msg *MsgBurnToken) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgBurnToken) Type() string { return TypeMsgBurnToken }

// GetSignBytes encodes the message for signing.
func (msg *MsgBurnToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgBurnToken) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgBurnToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.InvalidSender
	}
	if strings.TrimSpace(msg.TokenID) == "" {
		return errors.InvalidNFT
	}
	if !CheckUnique(msg.SubTokenIDs) {
		return errors.NotUniqueSubTokenIDs
	}

	return nil
}

////////////////////////////////////////////////////////////////

func CheckUnique(arr []uint32) bool {
	for i, el := range arr {
		for j, el2 := range arr {
			if i != j && el == el2 {
				return false
			}
		}
	}
	return true
}
