package types

import (
	fmt "fmt"

	"bitbucket.org/decimalteam/go-smart-node/utils/errors"
)

const (
	// Default codespace
	DefaultCodespace string = ModuleName

	CodeInternal uint32 = 100

	// CreateWallet
	CodeInvalidSender         uint32 = 101
	CodeInvalidOwnerCount     uint32 = 102
	CodeInvalidOwner          uint32 = 103
	CodeInvalidWeightCount    uint32 = 104
	CodeInvalidWeight         uint32 = 105
	CodeInvalidThreshold      uint32 = 106
	CodeWalletAccountNotFound uint32 = 107
	CodeDuplicateOwner        uint32 = 108
	CodeInvalidWallet         uint32 = 109
	CodeUnableToCreateWallet  uint32 = 110
	CodeWalletAlreadyExists   uint32 = 111
	CodeAccountAlreadyExists  uint32 = 112

	// CreateTransaction
	CodeUnableToCreateTransaction  uint32 = 201
	CodeInvalidTransactionIDError  uint32 = 202
	CodeInvalidTransactionIDPrefix uint32 = 203
	CodeInvalidAmountToSend        uint32 = 204
	CodeInsufficientFunds          uint32 = 205
	CodeInvalidReceiver            uint32 = 206
	CodeTransactionNotFound        uint32 = 207
	CodeNoCoinsToSend              uint32 = 208

	// SignTransaction
	CodeAlreadyEnoughSignatures  uint32 = 301
	CodeTransactionAlreadySigned uint32 = 302
	CodeSignerIsNotOwner         uint32 = 303
)

func ErrInternal(err string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInternal,
		fmt.Sprintf("Internal error: %s", err),
		errors.NewParam("error", err),
	)
}

func ErrInvalidSender(address string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidSender,
		fmt.Sprintf("invalid sender address: %s", address),
		errors.NewParam("address", address),
	)
}

func ErrInvalidOwnerCount(count, minCount, maxCount string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidOwnerCount,
		fmt.Sprintf("Invalid owner count: %s must in range [%s; %s]", count, minCount, maxCount),
		errors.NewParam("count", count),
		errors.NewParam("min_count", minCount),
		errors.NewParam("max_count", maxCount),
	)
}

func ErrInvalidOwner(address string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidOwner,
		fmt.Sprintf("invalid owner address: %s", address),
		errors.NewParam("address", address),
	)
}

func ErrInvalidWeightCount(lenMsgWeights string, lenMsgOwners string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidWeightCount,
		fmt.Sprintf("Invalid weight count: weight count (%s) is not equal to owner count (%s)", lenMsgWeights, lenMsgOwners),
		errors.NewParam("len_weight", lenMsgWeights),
		errors.NewParam("len_owners", lenMsgOwners),
	)
}

func ErrInvalidWeight(weight string, data string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidWeight,
		fmt.Sprintf("Invalid weight: weight cannot be %s than %s", data, weight),
	)
}

func ErrInvalidThreshold(sumOfWeights, threshold string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidThreshold,
		fmt.Sprintf("Sum of weights is less than threshold: %s < %s", sumOfWeights, threshold),
		errors.NewParam("sum_of_weights", sumOfWeights),
		errors.NewParam("threshold", threshold),
	)
}

func ErrWalletAccountNotFound(address string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeWalletAccountNotFound,
		fmt.Sprintf("wallet account %s not found", address),
		errors.NewParam("address", address),
	)
}

func ErrDuplicateOwner(address string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeDuplicateOwner,
		fmt.Sprintf("Invalid owners: owner with address %s is duplicated", address),
		errors.NewParam("address", address),
	)
}

func ErrInvalidWallet(address string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidWallet,
		fmt.Sprintf("invalid wallet address: %s", address),
		errors.NewParam("address", address),
	)
}

func ErrUnableToCreateWallet(err string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableToCreateWallet,
		fmt.Sprintf("unable to create wallet: %s", err),
		errors.NewParam("error", err),
	)
}

func ErrWalletAlreadyExists(address string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeWalletAlreadyExists,
		fmt.Sprintf("wallet with address %s already exists", address),
		errors.NewParam("address", address),
	)
}

func ErrAccountAlreadyExists(address string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeAccountAlreadyExists,
		fmt.Sprintf("account with address %s already exists", address),
		errors.NewParam("address", address),
	)
}

////////////////////////////////////////////////////////////////
// Transaction
////////////////////////////////////////////////////////////////

func ErrUnableToCreateTransaction(err string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableToCreateTransaction,
		fmt.Sprintf("unable to create multi-signature transaction: %s", err),
		errors.NewParam("error", err),
	)
}

func ErrInvalidTransactionIDPrefix(txID, prefixWant, prefixGot string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidTransactionIDPrefix,
		fmt.Sprintf("transaction ID %s must have prefix %s, got %s", txID, prefixWant, prefixGot),
		errors.NewParam("tx_id", txID),
		errors.NewParam("prefix_want", prefixWant),
		errors.NewParam("prefix_want", prefixGot),
	)
}

func ErrInvalidTransactionIDError(txID, err string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidTransactionIDError,
		fmt.Sprintf("transaction ID %s bech32 error: %s", txID, err),
		errors.NewParam("tx_id", txID),
		errors.NewParam("error", err),
	)
}

func ErrInvalidReceiver(address string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidReceiver,
		fmt.Sprintf("invalid receiver address: %s", address),
		errors.NewParam("address", address),
	)
}

func ErrInvalidAmount(coin, amount string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidAmountToSend,
		fmt.Sprintf("invalid amount for coin %s to send: %s", coin, amount),
		errors.NewParam("denom", coin),
		errors.NewParam("amount", amount),
	)
}

func ErrInsufficientFunds(funds, balance string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientFunds,
		fmt.Sprintf("Insufficient funds: wanted %s, but has %s", funds, balance),
		errors.NewParam("funds", funds),
		errors.NewParam("balance", balance),
	)
}

func ErrNoCoinsToSend() error {
	return errors.Encode(
		DefaultCodespace,
		CodeNoCoinsToSend,
		"No coins to send",
	)
}

func ErrTransactionNotFound(txID string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeTransactionNotFound,
		fmt.Sprintf("transaction with id %s not found", txID),
	)
}

func ErrAlreadyEnoughSignatures(confirmations, threshold string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeAlreadyEnoughSignatures,
		fmt.Sprintf("transaction already has enough signatures (%s >= %s)", confirmations, threshold),
		errors.NewParam("confirmations", confirmations),
		errors.NewParam("threshold", threshold),
	)
}

func ErrTransactionAlreadySigned(signer string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeTransactionAlreadySigned,
		fmt.Sprintf("transaction already signed by %s", signer),
		errors.NewParam("signer", signer),
	)
}

func ErrSignerIsNotOwner(signer, wallet string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeSignerIsNotOwner,
		fmt.Sprintf("transaction signer %s is not owner of %s", signer, wallet),
		errors.NewParam("signer", signer),
		errors.NewParam("wallet", wallet),
	)
}

func ErrUnablePreformTransaction(txID, err string) error {
	return errors.Encode(
		DefaultCodespace,
		CodeSignerIsNotOwner,
		fmt.Sprintf("unable to perform multi-signature transaction %s: %s", txID, err),
		errors.NewParam("tx_id", txID),
		errors.NewParam("error", err),
	)
}
