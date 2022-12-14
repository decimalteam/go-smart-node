package errors

import (
	"cosmossdk.io/errors"
)

var codespace = "fee"

var (
	Internal           = errors.New(codespace, 101, "internal error")
	WrongPrice         = errors.New(codespace, 102, "wrong price")
	UnknownOracle      = errors.New(codespace, 103, "unknown oracle address")
	SavingError        = errors.New(codespace, 104, "price saving error")
	DuplicateCoinPrice = errors.New(codespace, 105, "duplicate coin price by denom and quote")
	PriceNotFound      = errors.New(codespace, 106, "price is not found in the key-value store")
	MustBeOneMessage   = errors.New(codespace, 107, "transaction must contain one message")
)
