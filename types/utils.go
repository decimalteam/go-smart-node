package types

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

// IsSupportedKey returns true if the pubkey type is supported by the chain.
func IsSupportedKey(pk cryptotypes.PubKey) bool {
	switch pk := pk.(type) {
	case *ed25519.PubKey, *ethsecp256k1.PubKey:
		return true
	case multisig.PubKey:
		if len(pk.GetPubKeys()) == 0 {
			return false
		}
		for _, pk := range pk.GetPubKeys() {
			switch pk.(type) {
			case *ed25519.PubKey, *ethsecp256k1.PubKey:
				continue
			default:
				return false
			}
		}
		return true
	default:
		return false
	}
}

// GetDecimalAddressFromBech32 returns the sdk.Account address of given address, while
// also changing bech32 human readable prefix (HRP) to the value set on the global sdk.Config (eg: `dx`).
// The function fails if the provided bech32 address is invalid.
func GetDecimalAddressFromBech32(address string) (sdk.AccAddress, error) {
	bech32Prefix := strings.SplitN(address, "1", 2)[0]
	if bech32Prefix == address {
		return nil, sdkerrors.ErrInvalidAddress
	}

	addressBz, err := sdk.GetFromBech32(address, bech32Prefix)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress
	}

	// safety check: shouldn't happen
	if err := sdk.VerifyAddressFormat(addressBz); err != nil {
		return nil, err
	}

	return sdk.AccAddress(addressBz), nil
}

// GetLegacyAddressFromPubKey returns wallets address in old blockchain for given public key bytes
func GetLegacyAddressFromPubKey(pubKeyBytes []byte) (string, error) {
	oldPubKey := secp256k1.PubKey{Key: pubKeyBytes}
	return bech32.ConvertAndEncode("dx", oldPubKey.Address())
}
