package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
	"strings"
)

/* --------------------------------------------------------------------------- */
// MsgMintNFT
/* --------------------------------------------------------------------------- */

const regName = "^[a-zA-Z0-9_-]{1,255}$"

var MinReserve = sdk.NewInt(100)

//var NewMinReserve = helpers.BipToPip(sdk.NewInt(100))
//var NewMinReserve2 = helpers.BipToPip(sdk.NewInt(1))

// Route Implements Msg
func (m *MsgMintNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg *MsgMintNFT) Type() string { return "mint_nft" }

// ValidateBasic Implements Msg.
func (msg *MsgMintNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return ErrInvalidSenderAddress(msg.Sender)
	}
	_, err = sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return ErrInvalidRecipientAddress(msg.Recipient)
	}

	if strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidDenom(msg.Denom)
	}
	if strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidNFT(msg.ID)
	}
	if !msg.Quantity.IsPositive() {
		return ErrInvalidQuantity(msg.Quantity.String())
	}

	if !msg.Reserve.IsPositive() || msg.Reserve.LT(MinReserve) {
		return ErrInvalidReserve(msg.Reserve.String())
	}
	if match, _ := regexp.MatchString(regName, msg.Denom); !match {
		return ErrInvalidDenom(msg.Denom)
	}
	if match, _ := regexp.MatchString(regName, msg.ID); !match {
		return ErrInvalidTokenID(msg.ID)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg *MsgMintNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg *MsgMintNFT) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender}
}

/* --------------------------------------------------------------------------- */
// MsgBurnNFT
/* --------------------------------------------------------------------------- */

// NewMsgBurnNFT is a constructor function for MsgBurnNFT
func NewMsgBurnNFT(sender sdk.AccAddress, id string, denom string, subTokenIDs []int64) MsgBurnNFT {
	return MsgBurnNFT{
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
		return ErrInvalidSenderAddress(msg.Sender)
	}
	if strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidDenom(msg.Denom)
	}
	if strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidNFT(msg.ID)
	}
	if !CheckUnique(msg.SubTokenIDs) {
		return ErrNotUniqueSubTokenIDs()
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
func NewMsgUpdateReserveNFT(sender sdk.AccAddress, id string, denom string, subTokenIDs []int64, newReserveNFT sdk.Int) MsgUpdateReserveNFT {
	return MsgUpdateReserveNFT{
		Sender:        sender.String(),
		ID:            strings.TrimSpace(id),
		Denom:         strings.TrimSpace(denom),
		SubTokenIDs:   subTokenIDs,
		NewReserveNFT: newReserveNFT,
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
		return ErrInvalidSenderAddress(msg.Sender)
	}

	if strings.TrimSpace(msg.Denom) == "" {

		return ErrInvalidDenom(msg.Denom)
	}
	if strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidNFT(msg.ID)
	}
	if !CheckUnique(msg.SubTokenIDs) {
		return ErrNotUniqueSubTokenIDs()
	}

	if msg.NewReserveNFT.IsZero() {
		return ErrInvalidReserve("Reserv can not be equal to zero")
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
func NewMsgTransferNFT(sender, recipient sdk.AccAddress, denom, id string, subTokenIDs []int64) MsgTransferNFT {
	return MsgTransferNFT{
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
		return ErrInvalidCollection(msg.Denom)
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return ErrInvalidSenderAddress(msg.Sender)
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return ErrInvalidRecipientAddress(msg.Recipient)
	}

	if sender.Equals(recipient) {
		return ErrForbiddenToTransferToYourself()
	}

	if strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidCollection(msg.ID)
	}
	if !CheckUnique(msg.SubTokenIDs) {
		return ErrNotUniqueSubTokenIDs()
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
func NewMsgEditNFTMetadata(sender sdk.AccAddress, id,
	denom, tokenURI string,
) MsgEditNFTMetadata {
	return MsgEditNFTMetadata{
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
		return ErrInvalidSenderAddress(msg.Sender)
	}

	if strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidDenom(msg.Denom)
	}
	if strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidNFT(msg.ID)
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
func CheckUnique(arr []int64) bool {
	for i, el := range arr {
		for j, el2 := range arr {
			if i != j && el == el2 {
				return false
			}
		}
	}
	return true
}
