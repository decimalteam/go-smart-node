package errors

import (
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	"cosmossdk.io/errors"
)

var codespace = types.ModuleName

var (
	Internal                                       = errors.New(codespace, 1, "internal error")
	LegacyDoesNotExist                             = errors.New(codespace, 2, "legacy does not exist")
	InvalidPublicKeyLength                         = errors.New(codespace, 3, "invalid pub key length")
	CannnotGetLegacyAddressFromPublicKey           = errors.New(codespace, 4, "can not get legacy address from public key")
	NoMatchReceiverAndPKey                         = errors.New(codespace, 5, "receiver address and address from public key not match")
	NoLegacyBalance                                = errors.New(codespace, 6, "no legacy balance")
	NotFoundLegacyAddress                          = errors.New(codespace, 7, "not found legacy address in store")
	CannotGetActualAddressFromPublicKey            = errors.New(codespace, 8, "can not get actual address from public key")
	LegacyAddressesDuplicatedOnGenesis             = errors.New(codespace, 9, "legacy address duplicated on genesis")
	InvalidLegacyBech32Address                     = errors.New(codespace, 10, "address is not bech32 valid address")
	NoInfoForLegacyAddress                         = errors.New(codespace, 11, "no info for legacy address")
	OneOfLegacyAddresыСoinsBalanceIsNegativeOrZero = errors.New(codespace, 12, "one of address coin balance is negative or zero")
	WalletAddressIsNotValidBech32                  = errors.New(codespace, 13, "wallet address is not bech32 valid address")
)
