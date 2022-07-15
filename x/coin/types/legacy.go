package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (lb *LegacyBalance) Validate() error {
	// 'dx' is prefix from old Decimal
	_, err := sdk.GetFromBech32(lb.LegacyAddress, "dx")
	if err != nil {
		return fmt.Errorf("address '%s' is not bech32 valid address: %w", lb.LegacyAddress, err)
	}
	for _, coin := range lb.Coins {
		if coin.Amount.IsZero() || coin.Amount.IsNegative() {
			return fmt.Errorf("for address '%s', coin '%s' balance '%s' must be > 0",
				lb.LegacyAddress, coin.Denom, coin.Amount.String())
		}
	}
	return nil
}
