package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/* --------------------------------------------------------------------------- */
// MsgSaveBaseDenomPrice
/* --------------------------------------------------------------------------- */

// NewMsgSaveBaseDenomPrice is a constructor function for MsgMintNFT
func NewMsgSaveBaseDenomPrice(
	denom string,
	price float64,
	sender string,
) *MsgSaveBaseDenomPrice {
	return &MsgSaveBaseDenomPrice{
		BaseDenom: denom,
		Price:     price,
		Sender:    sender,
	}
}

const regName = "^[a-zA-Z0-9_-]{1,255}$"

// Route Implements Msg
func (m *MsgSaveBaseDenomPrice) Route() string { return RouterKey }

// Type Implements Msg
func (m *MsgSaveBaseDenomPrice) Type() string { return "mint_nft" }

// ValidateBasic Implements Msg.
func (m *MsgSaveBaseDenomPrice) ValidateBasic() error {
	// TODO implement
	return nil
}

// GetSignBytes Implements Msg.
func (m *MsgSaveBaseDenomPrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (m *MsgSaveBaseDenomPrice) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender}
}
