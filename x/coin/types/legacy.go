package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (lb *LegacyBalance) Validate() error {
	// 'dx' is prefix from old Decimal
	_, err := sdk.GetFromBech32(lb.OldAddress, "dx")
	if err != nil {
		return fmt.Errorf("address '%s' is not bech32 valid address: %w", lb.OldAddress, err)
	}
	for _, entry := range lb.Entries {
		if entry.Balance.IsZero() || entry.Balance.IsNegative() {
			return fmt.Errorf("for address '%s', coin '%s' balance '%s' must be > 0",
				lb.OldAddress, entry.CoinDenom, entry.Balance.String())
		}
	}
	return nil
}
