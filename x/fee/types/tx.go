package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/fee/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/* --------------------------------------------------------------------------- */
// MsgUpdateCoinPrices
/* --------------------------------------------------------------------------- */

// NewMsgUpdateCoinPrices is a constructor function for MsgUpdateCoinPrices
func NewMsgUpdateCoinPrices(
	sender string,
	prices []CoinPrice,
) *MsgUpdateCoinPrices {
	return &MsgUpdateCoinPrices{
		Oracle: sender,
		Prices: prices,
	}
}

// Route Implements Msg
func (m *MsgUpdateCoinPrices) Route() string { return RouterKey }

// Type Implements Msg
func (m *MsgUpdateCoinPrices) Type() string { return "save_base_denom_price" }

// ValidateBasic Implements Msg.
func (m *MsgUpdateCoinPrices) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Oracle); err != nil {
		return err
	}
	if len(m.Prices) == 0 {
		return errors.WrongPrice
	}
	type pair struct {
		denom string
		quote string
	}
	knownPairs := make(map[pair]bool)
	for _, price := range m.Prices {
		if !price.Price.IsPositive() {
			return errors.WrongPrice
		}
		key := pair{
			denom: price.Denom,
			quote: price.Quote,
		}
		if knownPairs[key] {
			return errors.DuplicateCoinPrice
		}
		knownPairs[key] = true
	}
	return nil
}

// GetSignBytes Implements Msg.
func (m *MsgUpdateCoinPrices) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (m *MsgUpdateCoinPrices) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Oracle)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender}
}
