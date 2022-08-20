package errors

import (
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	"cosmossdk.io/errors"
	"fmt"
)

var codespace = types.ModuleName

const (
	minOwnerCount = 2
	maxOwnerCount = 16
	minWeight     = 1
	maxWeight     = 1024
)

var (
	Internal                        = errors.New(codespace, 1, "internal error")
	InvalidSender                   = errors.New(codespace, 2, "invalid sender address")
	InvalidOwnerCount               = errors.New(codespace, 3, fmt.Sprintf("owner count must be in range [%d, %d]", minOwnerCount, maxOwnerCount))
	InvalidOwner                    = errors.New(codespace, 4, "invalid owner address")
	InvalidWeightCount              = errors.New(codespace, 5, "invalid weight count: weight count is not equal to owner count")
	InvalidWeight                   = errors.New(codespace, 6, fmt.Sprintf("weight must be less %d and greater than %d", maxWeight, minWeight))
	InvalidThreshold                = errors.New(codespace, 7, "sum of weights is less than threshold")
	WalletNotFound                  = errors.New(codespace, 8, "wallet not found")
	DuplicateOwner                  = errors.New(codespace, 9, "owners are duplicate")
	InvalidWallet                   = errors.New(codespace, 10, "invalid wallet address")
	UnableToCreateWallet            = errors.New(codespace, 11, "unable to create wallet")
	WalletAlreadyExists             = errors.New(codespace, 12, "wallet with address already exists")
	AccountAlreadyExists            = errors.New(codespace, 13, "account with address already exists")
	UnableToCreateTransaction       = errors.New(codespace, 14, "unable to create multi-signature transaction")
	InvalidTransactionIDPrefix      = errors.New(codespace, 15, "transaction ID have invalid prefix")
	InvalidTransactionIDError       = errors.New(codespace, 16, "transaction ID bech32 error")
	InvalidReceiver                 = errors.New(codespace, 17, "invalid receiver address")
	InvalidAmount                   = errors.New(codespace, 18, "invalid amount to send")
	InsufficientFunds               = errors.New(codespace, 19, "insufficient funds")
	NoCoinsToSend                   = errors.New(codespace, 20, "no coins to send")
	TransactionNotFound             = errors.New(codespace, 21, "tx not found")
	AlreadyEnoughSignatures         = errors.New(codespace, 22, "transaction already has enough signatures")
	TransactionAlreadySigned        = errors.New(codespace, 23, "tx already signed")
	SignerIsNotOwner                = errors.New(codespace, 24, "transaction signer is not owner")
	EmptyValueInKVStore             = errors.New(codespace, 25, "empty value in the key-value store")
	DuplicateWallet                 = errors.New(codespace, 26, "wallets are duplicate")
	DuplicateTxsID                  = errors.New(codespace, 27, "multisig transaction id duplicated on genesis")
	UnknownWalletInTx               = errors.New(codespace, 28, "multisig transaction have unknown wallet")
	TxSignersNotEqualToWalletOwners = errors.New(codespace, 29, "in multisig transaction signers count != wallet owners:")
	UnknownSignerInTx               = errors.New(codespace, 30, "multisig transaction have unknown signer")
)
