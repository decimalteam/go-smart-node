package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/* --------------------------------------------------------------------------- */
// MsgSaveBaseDenomPrice
/* --------------------------------------------------------------------------- */

// NewMsgSaveBaseDenomPrice is a constructor function for MsgSaveBaseDenomPrice
func NewMsgSaveBaseDenomPrice(
	sender string,
	denom string,
	price sdk.Dec,
) *MsgSaveBaseDenomPrice {
	return &MsgSaveBaseDenomPrice{
		Sender:    sender,
		BaseDenom: denom,
		Price:     price,
	}
}

// Route Implements Msg
func (m *MsgSaveBaseDenomPrice) Route() string { return RouterKey }

// Type Implements Msg
func (m *MsgSaveBaseDenomPrice) Type() string { return "save_base_denom_price" }

// ValidateBasic Implements Msg.
func (m *MsgSaveBaseDenomPrice) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if !m.Price.IsPositive() {
		return ErrWrongPrice(m.Price.String())
	}
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
