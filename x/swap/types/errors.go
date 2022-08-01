package types

import (
	fmt "fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-smart-node/utils/errors"
)

const (
	DefaultCodespace string = ModuleName

	CodeInvalidSenderAddress   uint32 = 100
	CodeSenderIsNotSwapService uint32 = 101
	CodeInvalidAmount          uint32 = 102
	CodeInvalidChainNumber     uint32 = 103
	CodeInvalidChainName       uint32 = 104

	// SwapInitialize
	CodeChainDoesNotExists       uint32 = 201
	CodeInsufficientAccountFunds uint32 = 202

	// SwapRedeem
	CodeInvalidRecipientAddress  uint32 = 301
	CodeInvalidTransactionNumber uint32 = 302
	CodeAlreadyRedeemed          uint32 = 303
	CodeInvalidServiceAddress    uint32 = 304
	CodeInsufficientPoolFunds    uint32 = 305
	CodeChainNumbersAreSame      uint32 = 306
)

func ErrInvalidSenderAddress(address string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidSenderAddress,
		fmt.Sprintf("invalid sender address: %s", address),
		errors.NewParam("address", address),
	)
}

func ErrSenderIsNotSwapService(address string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidSenderAddress,
		fmt.Sprintf("sender is not swap service address: %s", address),
		errors.NewParam("address", address),
	)
}

func ErrInvalidAmount() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidAmount,
		"amount should be greater than 0",
	)
}

func ErrInvalidChainNumber() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidChainNumber,
		"chain number should be greater than 0",
	)
}

func ErrInvalidChainName() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidChainName,
		"chain name should be not empty",
	)
}

// SwapInitialize
func ErrChainDoesNotExists(chainNumber string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeChainDoesNotExists,
		fmt.Sprintf("chain number %s dose not exists", chainNumber),
		errors.NewParam("chain_number", chainNumber),
	)
}

func ErrInsufficientAccountFunds(address string, coins string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientAccountFunds,
		fmt.Sprintf("account %s dose not have %s", address, coins),
		errors.NewParam("address", address),
		errors.NewParam("coins", coins),
	)
}

// SwapRedeem
func ErrInvalidRecipientAddress(address string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidRecipientAddress,
		fmt.Sprintf("invalid recipient address: %s", address),
		errors.NewParam("address", address),
	)
}

func ErrInvalidTransactionNumber(transactionNumber string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidTransactionNumber,
		fmt.Sprintf("invalid transaction number: %s", transactionNumber),
		errors.NewParam("transaction_number", transactionNumber),
	)
}

func ErrAlreadyRedeemed(hash string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeAlreadyRedeemed,
		fmt.Sprintf("swap already redeemed: %s", hash),
		errors.NewParam("hash", hash),
	)
}

func ErrInvalidServiceAddress(requiredAddress, givenAddress string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidServiceAddress,
		fmt.Sprintf("invalid service address: wait for %s, but get %s", requiredAddress, givenAddress),
		errors.NewParam("required_address", requiredAddress),
		errors.NewParam("given_address", givenAddress),
	)
}

func ErrInsufficientPoolFunds(want string, exists string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientPoolFunds,
		fmt.Sprintf("insufficient pool funds: want = %s, exists = %s", want, exists),
	)
}

func ErrChainNumbersAreSame() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeChainNumbersAreSame,
		"from chain and dest chain are same",
	)
}
