package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (rec *LegacyRecord) Validate() error {
	// 'dx' is prefix from old Decimal
	_, err := sdk.GetFromBech32(rec.Address, DecimalPrefix)
	if err != nil {
		return errors.InvalidLegacyBech32Address
	}
	// record must be not empty
	if len(rec.Coins) == 0 && len(rec.Nfts) == 0 && len(rec.Wallets) == 0 {
		return errors.NoInfoForLegacyAddress
	}

	for _, coin := range rec.Coins {
		if coin.Amount.IsZero() || coin.Amount.IsNegative() {
			return errors.OneOfLegacyAddressCoinsBalanceIsNegativeOrZero
		}
	}
	// wallets addresses must be valid bech32 addresses
	for _, w := range rec.Wallets {
		_, err := sdk.GetFromBech32(w, DecimalPrefix)
		if err != nil {
			return errors.WalletAddressIsNotValidBech32
		}
	}
	return nil
}
