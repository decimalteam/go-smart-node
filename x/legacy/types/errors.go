package types

import (
	"fmt"
	"strconv"

	"bitbucket.org/decimalteam/go-smart-node/utils/errors"
)

const (
	// Default codespace
	DefaultCodespace string = ModuleName

	CodeInternal                       uint32 = 100
	CodeLegacyDoesNotExist             uint32 = 101
	CodeInvalidPublicKeyLength         uint32 = 102
	CodeCannnotGetAddressFromPublicKey uint32 = 103
	CodeNoMatchReceiverAndPKey         uint32 = 104
	CodeNoLegacyBalance                uint32 = 105
)

func ErrInternal(err string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInternal,
		fmt.Sprintf("Internal error: %s", err),
		errors.NewParam("err", err),
	)
}

func ErrLegacyDoesNotExist(address string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeLegacyDoesNotExist,
		fmt.Sprintf("legacy for address '%s' does not exist", address),
		errors.NewParam("address", address),
	)
}

// Legacy return errors
func ErrInvalidPublicKeyLength(publicKeyLength int) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidPublicKeyLength,
		fmt.Sprintf("invalid public key length %d", publicKeyLength),
		errors.NewParam("public_key_length", strconv.Itoa(publicKeyLength)),
	)
}

func ErrCannnotGetAddressFromPublicKey(err string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeCannnotGetAddressFromPublicKey,
		fmt.Sprintf("can not get address from public key: %s", err),
		errors.NewParam("error", err),
	)
}

func ErrNoMatchReceiverAndPKey(expectAddress, gotAddress string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeNoMatchReceiverAndPKey,
		fmt.Sprintf("receiver address and address from public key not match: expect %s, but got %s",
			expectAddress, gotAddress),
		errors.NewParam("receiver_address", expectAddress),
		errors.NewParam("pubkey_address", gotAddress),
	)
}

func ErrNoLegacyBalance(receiver, legacyAddress string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeNoLegacyBalance,
		fmt.Sprintf("no legacy balance for receiver %s and legacy address %s", receiver, legacyAddress),
		errors.NewParam("receiver_address", receiver),
		errors.NewParam("legacy_address", legacyAddress),
	)
}
