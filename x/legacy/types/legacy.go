package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (rec *Record) Validate() error {
	// 'dx' is prefix from old Decimal
	_, err := sdk.GetFromBech32(rec.LegacyAddress, DecimalPrefix)
	if err != nil {
		return errors.InvalidLegacyBech32Address
	}
	// record must be not empty
	if len(rec.Coins) == 0 && len(rec.NFTs) == 0 && len(rec.Wallets) == 0 && len(rec.Validators) == 0 {
		return errors.NoInfoForLegacyAddress
	}

	for _, coin := range rec.Coins {
		if coin.Amount.IsZero() || coin.Amount.IsNegative() {
			return errors.OneOfLegacyAddressCoinsBalanceIsNegativeOrZero
		}
	}
	// wallets addresses must be valid bech32 addresses
	for _, w := range rec.Wallets {
		_, err := sdk.AccAddressFromBech32(w)
		if err != nil {
			return errors.WalletAddressIsNotValidBech32
		}
	}
	// validators addresses must be valid bech32
	for _, v := range rec.Validators {
		_, err := sdk.ValAddressFromBech32(v)
		if err != nil {
			return errors.ValidatorAddressIsNotValidBech32
		}
	}
	return nil
}
