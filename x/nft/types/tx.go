package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/errors"
	sdkmath "cosmossdk.io/math"

	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/* --------------------------------------------------------------------------- */
// MsgMintNFT
/* --------------------------------------------------------------------------- */

// NewMsgMintNFT is a constructor function for MsgMintNFT
func NewMsgMintNFT(
	sender, recipient sdk.AccAddress,
	id, denom, tokenURI string,
	quantity sdkmath.Int,
	reserve sdk.Coin,
	allowMint bool,
) *MsgMintNFT {
	return &MsgMintNFT{
		Sender:    sender.String(),
		Recipient: recipient.String(),
		ID:        strings.TrimSpace(id),
		Denom:     strings.TrimSpace(denom),
		TokenURI:  strings.TrimSpace(tokenURI),
		Quantity:  quantity,
		Reserve:   reserve,
		AllowMint: allowMint,
	}
}

// Route Implements Msg
func (m *MsgMintNFT) Route() string { return RouterKey }

// Type Implements Msg
func (m *MsgMintNFT) Type() string { return "mint_nft" }

// ValidateBasic Implements Msg.
func (m *MsgMintNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return errors.InvalidSender
	}
	_, err = sdk.AccAddressFromBech32(m.Recipient)
	if err != nil {
		return errors.InvalidRecipientAddress
	}

	if strings.TrimSpace(m.Denom) == "" {
		return errors.InvalidDenom
	}
	if strings.TrimSpace(m.ID) == "" {
		return errors.InvalidNFT
	}
	if !m.Quantity.IsPositive() {
		return errors.InvalidQuantity
	}

	if !m.Reserve.IsPositive() || m.Reserve.Amount.LT(MinReserve) {
		return errors.InvalidReserve
	}
	if match, _ := regexp.MatchString(regName, m.Denom); !match {
		return errors.InvalidDenom
	}
	if match, _ := regexp.MatchString(regName, m.ID); !match {
		return errors.InvalidTokenID
	}

	return nil
}

// GetSignBytes Implements Msg.
func (m *MsgMintNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (m *MsgMintNFT) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender}
}

/* --------------------------------------------------------------------------- */
// MsgBurnNFT
/* --------------------------------------------------------------------------- */

// NewMsgBurnNFT is a constructor function for MsgBurnNFT
func NewMsgBurnNFT(sender sdk.AccAddress, id string, denom string, subTokenIDs []uint64) *MsgBurnNFT {
	return &MsgBurnNFT{
		Sender:      sender.String(),
		ID:          strings.TrimSpace(id),
		Denom:       strings.TrimSpace(denom),
		SubTokenIDs: subTokenIDs,
	}
}

// Route Implements Msg
func (msg *MsgBurnNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg *MsgBurnNFT) Type() string { return "burn_nft" }

// ValidateBasic Implements Msg.
func (msg *MsgBurnNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.InvalidSender
	}
	if strings.TrimSpace(msg.Denom) == "" {
		return errors.InvalidDenom
	}
	if strings.TrimSpace(msg.ID) == "" {
		return errors.InvalidNFT
	}
	if !CheckUnique(msg.SubTokenIDs) {
		return errors.NotUniqueSubTokenIDs
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg *MsgBurnNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg *MsgBurnNFT) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender}
}

/* --------------------------------------------------------------------------- */
// MsgUpdateReservNFT
/* --------------------------------------------------------------------------- */

// NewUpdateReservNFT is a constructor function for MsgUpdateReservNFT
func NewMsgUpdateReserveNFT(sender sdk.AccAddress, id string, denom string, subTokenIDs []uint64, newReserve sdk.Coin) *MsgUpdateReserveNFT {
	return &MsgUpdateReserveNFT{
		Sender:      sender.String(),
		ID:          strings.TrimSpace(id),
		Denom:       strings.TrimSpace(denom),
		SubTokenIDs: subTokenIDs,
		NewReserve:  newReserve,
	}
}

// Route Implements Msg
func (msg *MsgUpdateReserveNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg *MsgUpdateReserveNFT) Type() string { return "update_nft_reserve" }

// ValidateBasic Implements Msg.
func (msg *MsgUpdateReserveNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.InvalidSender
	}

	if strings.TrimSpace(msg.Denom) == "" {

		return errors.InvalidDenom
	}
	if strings.TrimSpace(msg.ID) == "" {
		return errors.InvalidNFT
	}
	if !CheckUnique(msg.SubTokenIDs) {
		return errors.NotUniqueSubTokenIDs
	}

	if msg.NewReserve.IsZero() {
		return errors.InvalidReserve
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg *MsgUpdateReserveNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg *MsgUpdateReserveNFT) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender}
}

/* --------------------------------------------------------------------------- */
// MsgTransferNFT
/* --------------------------------------------------------------------------- */

// NewMsgTransferNFT is a constructor function for MsgSetName
func NewMsgTransferNFT(sender, recipient sdk.AccAddress, denom, id string, subTokenIDs []uint64) *MsgTransferNFT {
	return &MsgTransferNFT{
		Sender:      sender.String(),
		Recipient:   recipient.String(),
		Denom:       strings.TrimSpace(denom),
		ID:          strings.TrimSpace(id),
		SubTokenIDs: subTokenIDs,
	}
}

// Route Implements Msg
func (msg *MsgTransferNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg *MsgTransferNFT) Type() string { return "transfer_nft" }

// ValidateBasic Implements Msg.
func (msg *MsgTransferNFT) ValidateBasic() error {
	if strings.TrimSpace(msg.Denom) == "" {
		return errors.InvalidCollection
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.InvalidSender
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return errors.InvalidRecipientAddress
	}

	if sender.Equals(recipient) {
		return errors.ForbiddenToTransferToYourself
	}

	if strings.TrimSpace(msg.ID) == "" {
		return errors.InvalidCollection
	}
	if !CheckUnique(msg.SubTokenIDs) {
		return errors.NotUniqueSubTokenIDs
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg *MsgTransferNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg *MsgTransferNFT) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender}
}

/* --------------------------------------------------------------------------- */
// MsgEditNFTMetadata
/* --------------------------------------------------------------------------- */

// NewMsgEditNFTMetadata is a constructor function for MsgSetName
func NewMsgEditNFTMetadata(sender sdk.AccAddress, id, denom, tokenURI string) *MsgEditNFTMetadata {
	return &MsgEditNFTMetadata{
		Sender:   sender.String(),
		ID:       strings.TrimSpace(id),
		Denom:    strings.TrimSpace(denom),
		TokenURI: strings.TrimSpace(tokenURI),
	}
}

// Route Implements Msg
func (msg *MsgEditNFTMetadata) Route() string { return RouterKey }

// Type Implements Msg
func (msg *MsgEditNFTMetadata) Type() string { return "edit_nft_metadata" }

// ValidateBasic Implements Msg.
func (msg *MsgEditNFTMetadata) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.InvalidSender
	}

	if strings.TrimSpace(msg.Denom) == "" {
		return errors.InvalidDenom
	}
	if strings.TrimSpace(msg.ID) == "" {
		return errors.InvalidNFT
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg *MsgEditNFTMetadata) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg *MsgEditNFTMetadata) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender}
}

/* --------------------------------------------------------------------------- */

func CheckUnique(arr []uint64) bool {
	for i, el := range arr {
		for j, el2 := range arr {
			if i != j && el == el2 {
				return false
			}
		}
	}
	return true
}
