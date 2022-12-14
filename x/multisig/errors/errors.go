package errors

import (
	"fmt"

	"cosmossdk.io/errors"
)

var codespace = "multisig"

const (
	minOwnerCount = 2
	maxOwnerCount = 16
	minWeight     = 1
	maxWeight     = 1024
)

var (
	Internal                        = errors.New(codespace, 101, "internal error")
	InvalidSender                   = errors.New(codespace, 102, "invalid sender address")
	InvalidOwnerCount               = errors.New(codespace, 103, fmt.Sprintf("owner count must be in range [%d, %d]", minOwnerCount, maxOwnerCount))
	InvalidOwner                    = errors.New(codespace, 104, "invalid owner address")
	InvalidWeightCount              = errors.New(codespace, 105, "invalid weight count: weight count is not equal to owner count")
	InvalidWeight                   = errors.New(codespace, 106, fmt.Sprintf("weight must be less %d and greater than %d", maxWeight, minWeight))
	InvalidThreshold                = errors.New(codespace, 107, "sum of weights is less than threshold")
	WalletNotFound                  = errors.New(codespace, 108, "wallet not found")
	DuplicateOwner                  = errors.New(codespace, 109, "owners are duplicate")
	InvalidWallet                   = errors.New(codespace, 110, "invalid wallet address")
	UnableToCreateWallet            = errors.New(codespace, 111, "unable to create wallet")
	WalletAlreadyExists             = errors.New(codespace, 112, "wallet with address already exists")
	AccountAlreadyExists            = errors.New(codespace, 113, "account with address already exists")
	UnableToCreateTransaction       = errors.New(codespace, 114, "unable to create multi-signature transaction")
	InvalidTransactionIDPrefix      = errors.New(codespace, 115, "transaction ID have invalid prefix")
	InvalidTransactionIDError       = errors.New(codespace, 116, "transaction ID bech32 error")
	InvalidReceiver                 = errors.New(codespace, 117, "invalid receiver address")
	InvalidAmount                   = errors.New(codespace, 118, "invalid amount to send")
	InsufficientFunds               = errors.New(codespace, 119, "insufficient funds")
	NoCoinsToSend                   = errors.New(codespace, 120, "no coins to send")
	TransactionNotFound             = errors.New(codespace, 121, "tx not found")
	AlreadyEnoughSignatures         = errors.New(codespace, 122, "transaction already has enough signatures")
	TransactionAlreadySigned        = errors.New(codespace, 123, "tx already signed")
	SignerIsNotOwner                = errors.New(codespace, 124, "transaction signer is not owner")
	EmptyValueInKVStore             = errors.New(codespace, 125, "empty value in the key-value store")
	DuplicateWallet                 = errors.New(codespace, 126, "wallets are duplicate")
	DuplicateTxsID                  = errors.New(codespace, 127, "multisig transaction id duplicated on genesis")
	UnknownWalletInTx               = errors.New(codespace, 128, "multisig transaction have unknown wallet")
	TxSignersNotEqualToWalletOwners = errors.New(codespace, 129, "in multisig transaction signers count != wallet owners:")
	UnknownSignerInTx               = errors.New(codespace, 130, "multisig transaction have unknown signer")

	NoSignersInInternal         = errors.New(codespace, 131, "no signers in internal message")
	WalletIsNotSignerInInternal = errors.New(codespace, 132, "wallet isn't in signers in internal message")
	NoHandlerForInternal        = errors.New(codespace, 133, "no handlers to process internal message")
	EmptyInternalMessage        = errors.New(codespace, 134, "empty internal message")
)
