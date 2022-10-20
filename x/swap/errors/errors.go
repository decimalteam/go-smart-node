package errors

import (
	"cosmossdk.io/errors"
)

var codespace = "swap"

var (
	Internal                 = errors.New(codespace, 101, "internal error")
	InvalidSenderAddress     = errors.New(codespace, 102, "invalid sender address")
	SenderIsNotSwapService   = errors.New(codespace, 103, "sender is not swap service address")
	InvalidAmount            = errors.New(codespace, 104, "amount should be greater than 0")
	InvalidChainNumber       = errors.New(codespace, 105, "chain number should be greater than 0")
	InvalidChainName         = errors.New(codespace, 106, "chain name should be not empty")
	ChainDoesNotExists       = errors.New(codespace, 107, "chain does not exist")
	InvalidRecipientAddress  = errors.New(codespace, 108, "invalid recipient address")
	InsufficientAccountFunds = errors.New(codespace, 109, "account does not have coins")
	InvalidTransactionNumber = errors.New(codespace, 110, "invalid tx number")
	AlreadyRedeemed          = errors.New(codespace, 111, "swap already redeemed")
	InvalidServiceAddress    = errors.New(codespace, 112, "invalid service address")
	InsufficientPoolFunds    = errors.New(codespace, 113, "insufficient pool funds")
	ChainNumbersAreSame      = errors.New(codespace, 114, "from chain and dest chain are same")
	InvalidHexStringR        = errors.New(codespace, 115, "invalid hex representation of R")
	InvalidHexStringS        = errors.New(codespace, 116, "invalid hex representation of S")
	InvalidLengthR           = errors.New(codespace, 117, "length R must be 32 bytes")
	InvalidLengthS           = errors.New(codespace, 118, "length R must be 32 bytes")
)
