package errors

import (
	"cosmossdk.io/errors"
)

var codespace = "swap"

var (
	Internal                 = errors.New(codespace, 1, "internal error")
	InvalidSenderAddress     = errors.New(codespace, 2, "invalid sender address")
	SenderIsNotSwapService   = errors.New(codespace, 3, "sender is not swap service address")
	InvalidAmount            = errors.New(codespace, 4, "amount should be greater than 0")
	InvalidChainNumber       = errors.New(codespace, 5, "chain number should be greater than 0")
	InvalidChainName         = errors.New(codespace, 6, "chain name should be not empty")
	ChainDoesNotExists       = errors.New(codespace, 7, "chain does not exist")
	InvalidRecipientAddress  = errors.New(codespace, 8, "invalid recipient address")
	InsufficientAccountFunds = errors.New(codespace, 9, "account does not have coins")
	InvalidTransactionNumber = errors.New(codespace, 10, "invalid tx number")
	AlreadyRedeemed          = errors.New(codespace, 11, "swap already redeemed")
	InvalidServiceAddress    = errors.New(codespace, 12, "invalid service address")
	InsufficientPoolFunds    = errors.New(codespace, 13, "insufficient pool funds")
	ChainNumbersAreSame      = errors.New(codespace, 14, "from chain and dest chain are same")
)
