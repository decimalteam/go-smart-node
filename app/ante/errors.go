package ante

import (
	"cosmossdk.io/errors"
)

var codespace = "ante"

var (
	FeePayerAddressDoesNotExist    = errors.New(codespace, 1, "fee payer address does not exist in network")
	FeeLessThanCommission          = errors.New(codespace, 2, "insufficient funds to pay for fees")
	FailedToSendCoins              = errors.New(codespace, 3, "failed to send coins")
	InsufficientFundsToPayFee      = errors.New(codespace, 4, "insufficient funds to pay for fee")
	InvalidFeeAmount               = errors.New(codespace, 5, "invalid fee amount")
	UnknownTransaction             = errors.New(codespace, 6, "unknown transaction type")
	CoinReserveInsufficient        = errors.New(codespace, 7, "coin reserve balance is not sufficient for transaction")
	CoinReserveBecomeInsufficient  = errors.New(codespace, 8, "coin reserve will become lower than minimal reserve.")
	NotFeeTxType                   = errors.New(codespace, 9, "x must be a FeeTx")
	CountOfMsgsMustBeOne           = errors.New(codespace, 10, "count of messages must be 1")
	InvalidAddressOfCreatedAccount = errors.New(codespace, 11, "invalid address of created account")
	UnableToFindCreatedAccount     = errors.New(codespace, 12, "unable to find created account")
)
