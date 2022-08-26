package errors

import (
	"cosmossdk.io/errors"
)

var codespace = "fee"

var (
	Internal      = errors.New(codespace, 101, "internal error")
	WrongPrice    = errors.New(codespace, 102, "wrong price")
	UnknownOracle = errors.New(codespace, 103, "unknown oracle address")
	SavingError   = errors.New(codespace, 104, "price saving error")
)
