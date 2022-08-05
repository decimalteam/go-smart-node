package types

import (
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/utils/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type CodeType = uint32

const (
	// Default coin codespace
	DefaultCodespace string = ModuleName

	CodeInvaliPrice   CodeType = 101
	CodeUnknownOracle CodeType = 102
	CodeSavingError   CodeType = 103
)

func ErrWrongPrice(price string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvaliPrice,
		fmt.Sprintf("wrong price: %s", price),
		errors.NewParam("price", price),
	)
}

func ErrUnknownOracle(address string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnknownOracle,
		fmt.Sprintf("unknown oracle address: %s", address),
		errors.NewParam("address", address),
	)
}

func ErrSavingError(err string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeSavingError,
		fmt.Sprintf("price saving error: %s", err),
		errors.NewParam("error", err),
	)
}
