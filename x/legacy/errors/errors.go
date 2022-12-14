package errors

import (
	"cosmossdk.io/errors"
)

var codespace = "legacy"

var (
	Internal                                       = errors.New(codespace, 101, "internal error")
	CannotGetLegacyAddressFromPublicKey            = errors.New(codespace, 102, "can not get legacy address from public key")
	NotFoundLegacyAddress                          = errors.New(codespace, 103, "not found legacy address in store")
	CannotGetActualAddressFromPublicKey            = errors.New(codespace, 104, "can not get actual address from public key")
	LegacyAddressesDuplicatedOnGenesis             = errors.New(codespace, 105, "legacy address duplicated on genesis")
	InvalidLegacyBech32Address                     = errors.New(codespace, 106, "address is not bech32 valid address")
	NoInfoForLegacyAddress                         = errors.New(codespace, 107, "no info for legacy address")
	OneOfLegacyAddressCoinsBalanceIsNegativeOrZero = errors.New(codespace, 108, "one of address coin balance is negative or zero")
	WalletAddressIsNotValidBech32                  = errors.New(codespace, 109, "wallet address is not bech32 valid address")
	InvalidPublicKeyLength                         = errors.New(codespace, 110, "invalid public key length")
	InvalidSenderAddress                           = errors.New(codespace, 111, "invalid sender address")
	ValidatorAddressIsNotValidBech32               = errors.New(codespace, 112, "validator address is not bech32 valid address")
)
