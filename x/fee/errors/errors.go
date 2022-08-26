package errors

import (
	"cosmossdk.io/errors"
)

var codespace = "fee"

var (
	WrongPrice    = errors.New(codespace, 1, "wrong price")
	UnknownOracle = errors.New(codespace, 2, "unknown oracle address")
	SavingError   = errors.New(codespace, 3, "price saving error")
)
