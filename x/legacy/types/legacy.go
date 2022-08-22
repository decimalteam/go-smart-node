package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (rec *LegacyRecord) Validate() error {
	// 'dx' is prefix from old Decimal
	_, err := sdk.GetFromBech32(rec.Address, DecimalPrefix)
	if err != nil {
		return fmt.Errorf("address '%s' is not bech32 valid address: %w", rec.Address, err)
	}
	// record must be not empty
	if len(rec.Coins) == 0 && len(rec.Nfts) == 0 && len(rec.Wallets) == 0 {
		return fmt.Errorf("no info for legacy address '%s'", rec.Address)
	}

	for _, coin := range rec.Coins {
		if coin.Amount.IsZero() || coin.Amount.IsNegative() {
			return fmt.Errorf("for address '%s', coin '%s' balance '%s' must be > 0",
				rec.Address, coin.Denom, coin.Amount.String())
		}
	}
	// wallets addresses must be valid bech32 addresses
	for _, w := range rec.Wallets {
		_, err := sdk.GetFromBech32(w, DecimalPrefix)
		if err != nil {
			return fmt.Errorf("for owner '%s' address '%s' is not bech32 valid address: %w", rec.Address, w, err)
		}
	}
	return nil
}
