package errors

import (
	"cosmossdk.io/errors"
)

var codespace = "legacy"

var (
	Internal                                       = errors.New(codespace, 1, "internal error")
	CannotGetLegacyAddressFromPublicKey            = errors.New(codespace, 2, "can not get legacy address from public key")
	NotFoundLegacyAddress                          = errors.New(codespace, 3, "not found legacy address in store")
	CannotGetActualAddressFromPublicKey            = errors.New(codespace, 4, "can not get actual address from public key")
	LegacyAddressesDuplicatedOnGenesis             = errors.New(codespace, 5, "legacy address duplicated on genesis")
	InvalidLegacyBech32Address                     = errors.New(codespace, 6, "address is not bech32 valid address")
	NoInfoForLegacyAddress                         = errors.New(codespace, 7, "no info for legacy address")
	OneOfLegacyAddressCoinsBalanceIsNegativeOrZero = errors.New(codespace, 8, "one of address coin balance is negative or zero")
	WalletAddressIsNotValidBech32                  = errors.New(codespace, 9, "wallet address is not bech32 valid address")
)
