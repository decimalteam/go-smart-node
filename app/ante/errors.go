package ante

import (
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/utils/errors"
)

// Local code type
type CodeType = uint32

const (
	// Default ante codespace
	DefaultRootCodespace string = "ante"

	CodeFeePayerAddressDoesNotExist    CodeType = 101
	CodeFeeLessThanCommission          CodeType = 102
	CodeFailedToSendCoins              CodeType = 103
	CodeInsufficientFundsToPayFee      CodeType = 104
	CodeInvalidFeeAmount               CodeType = 105
	CodeCoinDoesNotExist               CodeType = 106
	CodeNotStdTxType                   CodeType = 107
	CodeNotFeeTxType                   CodeType = 108
	CodeOutOfGas                       CodeType = 109
	CodeNotGasTxType                   CodeType = 110
	CodeInvalidAddressOfCreatedAccount CodeType = 111
	CodeUnableToFindCreatedAccount     CodeType = 112
	CodeUnknownTransaction             CodeType = 113
	CodeCoinReserveInsufficient        CodeType = 114
	CodeCoinReserveBecomeInsufficient  CodeType = 115
	CodeCoinVolumeBecomeInsufficient   CodeType = 116
)

func ErrFeePayerAddressDoesNotExist(feePayer string) error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeFeePayerAddressDoesNotExist,
		fmt.Sprintf("fee payer address does not exist: %s", feePayer),
		errors.NewParam("fee_payer", feePayer),
	)
}

func ErrFeeLessThanCommission(feeInBaseCoin, commissionInBaseCoin string) error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeFeeLessThanCommission,
		fmt.Sprintf("insufficient funds to pay for fees; %s < %s", feeInBaseCoin, commissionInBaseCoin),
		errors.NewParam("fee", feeInBaseCoin),
		errors.NewParam("commission", commissionInBaseCoin),
	)
}

func ErrFailedToSendCoins(err string) error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeFailedToSendCoins,
		fmt.Sprintf("failed to send coins: %s", err),
		errors.NewParam("error", err),
	)
}

func ErrInsufficientFundsToPayFee(coins, fee string) error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeInsufficientFundsToPayFee,
		fmt.Sprintf("insufficient funds to pay for fee; %s < %s", coins, fee),
		errors.NewParam("coins", coins),
		errors.NewParam("fee", fee),
	)
}

func ErrInvalidFeeAmount(fee string) error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeInvalidFeeAmount,
		"invalid fee amount",
		errors.NewParam("fee", fee),
	)
}

func ErrCoinDoesNotExist(feeDenom string) error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeCoinDoesNotExist,
		fmt.Sprintf("coin not exist: %s", feeDenom),
		errors.NewParam("fee_denom", feeDenom),
	)
}

func ErrNotStdTxType() error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeNotStdTxType,
		"Tx must be StdTx",
	)
}

func ErrNotFeeTxType() error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeNotFeeTxType,
		"x must be a FeeTx",
	)
}

func ErrNotGasTxType() error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeNotGasTxType,
		"Tx must be GasTx",
	)
}

func ErrOutOfGas(location, gasWanted, gasUsed string) error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeOutOfGas,
		fmt.Sprintf("out of gas in location: %s; gasWanted: %s, gasUsed: %s", location, gasWanted, gasUsed),
		errors.NewParam("location", location),
		errors.NewParam("gas_wanted", gasWanted),
		errors.NewParam("gas_used", gasUsed),
	)
}

func ErrInvalidAddressOfCreatedAccount() error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeInvalidAddressOfCreatedAccount,
		"invalid address of created account",
	)
}

func ErrUnableToFindCreatedAccount() error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeUnableToFindCreatedAccount,
		"unable to find created account",
	)
}

func ErrUnknownTransaction(txType string) error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeUnknownTransaction,
		fmt.Sprintf("unknown transaction type: %s", txType),
		errors.NewParam("tx_type", txType),
	)
}

func ErrCoinReserveInsufficient(reserve, commission string) error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeCoinReserveInsufficient,
		fmt.Sprintf("coin reserve balance is not sufficient for transaction. Has: %s, required %s", reserve, commission),
		errors.NewParam("reserve", reserve),
		errors.NewParam("commission", commission),
	)
}

func ErrCoinReserveBecomeInsufficient(reserve, decreasing, minReserve string) error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeCoinReserveBecomeInsufficient,
		fmt.Sprintf("coin reserve will become lower than minimal reserve. Has: %s, decreasing: %s, minimal reserve: %s", reserve, decreasing, minReserve),
		errors.NewParam("reserve", reserve),
		errors.NewParam("decreasing", decreasing),
		errors.NewParam("min_reserve", minReserve),
	)
}

func ErrCoinVolumeBecomeInsufficient(volume, decreasing, minVolume string) error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeCoinVolumeBecomeInsufficient,
		fmt.Sprintf("coin volume will become lower than minimal volume. Has: %s, decreasing: %s, minimal volume: %s", volume, decreasing, minVolume),
		errors.NewParam("volume", volume),
		errors.NewParam("decreasing", decreasing),
		errors.NewParam("min_volume", minVolume),
	)
}
