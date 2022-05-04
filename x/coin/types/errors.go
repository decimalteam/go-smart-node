package types

// DONTCOVER

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-smart-node/utils/errors"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

const (
	// Default codespace
	DefaultCodespace string = ModuleName

	// Create coin
	CodeInvalidCRR                      uint32 = 100
	CodeCoinDoesNotExist                uint32 = 101
	CodeInvalidCoinSymbol               uint32 = 102
	CodeForbiddenCoinSymbol             uint32 = 103
	CodeRetrievedAnotherCoin            uint32 = 104
	CodeCoinAlreadyExists               uint32 = 105
	CodeInvalidCoinTitle                uint32 = 106
	CodeInvalidCoinInitialVolume        uint32 = 107
	CodeInvalidCoinInitialReserve       uint32 = 108
	CodeInternal                        uint32 = 109
	CodeInsufficientCoinReserve         uint32 = 110
	CodeInsufficientCoinToPayCommission uint32 = 111
	CodeInsufficientFunds               uint32 = 112
	CodeCalculateCommission             uint32 = 113
	CodeForbiddenUpdate                 uint32 = 114

	// Buy/Sell coin
	CodeSameCoins                  uint32 = 200
	CodeInsufficientFundsToSellAll uint32 = 201
	CodeTxBreaksVolumeLimit        uint32 = 202
	CodeTxBreaksMinReserveLimit    uint32 = 203
	CodeMaximumValueToSellReached  uint32 = 204
	CodeMinimumValueToBuyReached   uint32 = 205
	CodeUpdateBalance              uint32 = 206
	CodeLimitVolumeBroken          uint32 = 207

	// Send coin
	CodeInvalidAmount          uint32 = 300
	CodeInvalidReceiverAddress uint32 = 301
	CodeInvalidSenderAddress   uint32 = 302

	// Redeem check
	CodeInvalidCheck          uint32 = 400
	CodeInvalidProof          uint32 = 401
	CodeInvalidPassphrase     uint32 = 402
	CodeInvalidChainID        uint32 = 403
	CodeInvalidNonce          uint32 = 404
	CodeCheckExpired          uint32 = 405
	CodeCheckRedeemed         uint32 = 406
	CodeUnableDecodeCheck     uint32 = 407
	CodeUnableRPLEncodeCheck  uint32 = 408
	CodeUnableSignCheck       uint32 = 409
	CodeUnableDecodeProof     uint32 = 410
	CodeUnableRecoverAddress  uint32 = 411
	CodeUnableRecoverLockPkey uint32 = 412

	// AccountKeys
	CodeInvalidPkey              uint32 = 500
	CodeUnableRetriveArmoredPkey uint32 = 501
	CodeUnableRetrivePkey        uint32 = 502
	CodeUnableRetriveSECPPkey    uint32 = 503
)

const maxCoinNameBytes = 64
const allowedCoinSymbols = "^[a-zA-Z][a-zA-Z0-9]{2,9}$"

var minCoinSupply = sdk.NewInt(1)
var maxCoinSupply = helpers.BipToPip(sdk.NewInt(1000000000000000))
var MinCoinReserve = helpers.BipToPip(sdk.NewInt(1000))

func ErrInvalidCRR(crr string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCRR,
		fmt.Sprintf("coin CRR must be between 10 and 100, crr is: %s", crr),
		errors.NewParam("crr", crr),
	)
}

func ErrCoinDoesNotExist(symbol string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCoinDoesNotExist,
		fmt.Sprintf("coin %s does not exist", symbol),
		errors.NewParam("symbol", symbol),
	)
}

func ErrInvalidCoinSymbol(symbol string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCoinSymbol,
		fmt.Sprintf("invalid coin symbol %s. Symbol must match this regular expression: %s", symbol, allowedCoinSymbols),
		errors.NewParam("symbol", symbol),
	)
}

func ErrForbiddenCoinSymbol(symbol string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeForbiddenCoinSymbol,
		fmt.Sprintf("forbidden coin symbol %s", symbol),
		errors.NewParam("symbol", symbol),
	)
}

func ErrRetrievedAnotherCoin(symbolWant string, symbolRetrieved string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeRetrievedAnotherCoin,
		fmt.Sprintf("retrieved coin %s instead %s", symbolRetrieved, symbolWant),
		errors.NewParam("symbol_want", symbolWant),
		errors.NewParam("symbol_retrieved", symbolRetrieved),
	)
}

func ErrCoinAlreadyExists(coin string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCoinAlreadyExists,
		fmt.Sprintf("coin %s already exist", coin),
		errors.NewParam("coin", coin),
	)
}

func ErrInvalidCoinTitle(title string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCoinTitle,
		fmt.Sprintf("invalid coin title: %s. Allowed up to %d bytes", title, maxCoinNameBytes),
		errors.NewParam("title", title),
	)
}

func ErrInvalidCoinInitialVolume(initialVolume string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCoinInitialVolume,
		fmt.Sprintf("coin initial volume should be between %s and %s. Given %s", minCoinSupply.String(), maxCoinSupply.String(), initialVolume),
		errors.NewParam("min_coin_supply", minCoinSupply.String()),
		errors.NewParam("max_coin_supply", maxCoinSupply.String()),
		errors.NewParam("initial_volume", initialVolume),
	)
}

func ErrInvalidCoinInitialReserve(reserve string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCoinInitialReserve,
		fmt.Sprintf("coin initial reserve should be greater than or equal to %s", reserve),
		errors.NewParam("reserve", reserve),
	)
}

func ErrInternal(err string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInternal,
		fmt.Sprintf("Internal error: %s", err),
		errors.NewParam("err", err),
	)
}

func ErrInsufficientCoinReserve() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientCoinReserve,
		"not enough coin to reserve",
	)
}

func ErrInsufficientFundsToPayCommission(commission string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientCoinToPayCommission,
		fmt.Sprintf("insufficient funds to pay commission: wanted = %s", commission),
		errors.NewParam("commission", commission),
	)
}

func ErrInsufficientFunds(fundsWant string, fundsExist string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientFunds,
		fmt.Sprintf("insufficient account funds; %s < %s", fundsExist, fundsWant),
		errors.NewParam("funds_want", fundsWant),
		errors.NewParam("funds_exist", fundsExist),
	)
}

func ErrCalculateCommission(err string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCalculateCommission,
		err,
	)
}

func ErrUpdateOnlyForCreator() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeForbiddenUpdate,
		"updating allowed only for creator of coin",
	)
}

func ErrSameCoin() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeSameCoins,
		"can't buy same coins",
	)
}

func ErrInsufficientFundsToSellAll() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientFundsToSellAll,
		"not enough coin to sell",
	)
}

func ErrTxBreaksVolumeLimit(volume string, limitVolume string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeTxBreaksVolumeLimit,
		fmt.Sprintf("tx breaks LimitVolume rule: %s > %s", volume, limitVolume),
		errors.NewParam("volume", volume),
		errors.NewParam("limit_volume", limitVolume),
	)
}

func ErrTxBreaksMinReserveRule(minCoinReserve string, reserve string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeTxBreaksMinReserveLimit,
		fmt.Sprintf("tx breaks MinReserveLimit rule: %s < %s", reserve, minCoinReserve),
		errors.NewParam("reserve", reserve),
		errors.NewParam("min_coin_reserve", minCoinReserve),
	)
}

func ErrMaximumValueToSellReached(amount string, max string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeMaximumValueToSellReached,
		fmt.Sprintf("wanted to sell maximum %s, but required to spend %s at the moment", max, amount),
		errors.NewParam("max", max),
		errors.NewParam("amount", amount),
	)
}

func ErrMinimumValueToBuyReached(amount string, min string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeMinimumValueToBuyReached,
		fmt.Sprintf("wanted to buy minimum %s, but expected to receive %s at the moment", min, amount),
		errors.NewParam("min", min),
		errors.NewParam("amount", amount),
	)
}

func ErrUpdateBalance(account string, err string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUpdateBalance,
		fmt.Sprintf("unable to update balance of account %s: %s", account, err),
		errors.NewParam("account", account),
		errors.NewParam("err", err),
	)
}

func ErrLimitVolumeBroken(volume string, limit string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeLimitVolumeBroken,
		fmt.Sprintf("volume should be less than or equal the volume limit: %s > %s", volume, limit),
		errors.NewParam("volume", volume),
		errors.NewParam("limit", limit),
	)
}

func ErrInvalidAmount() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidAmount,
		"amount should be greater than 0",
	)
}

func ErrInvalidReceiverAddress(address string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidReceiverAddress,
		fmt.Sprintf("invalid receiver address: %s", address),
	)
}

func ErrInvalidSenderAddress(address string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidSenderAddress,
		fmt.Sprintf("invalid sender address: %s", address),
	)
}

// Redeem check

func ErrInvalidCheck(data string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCheck,
		fmt.Sprintf("unable to parse check: %s", data),
		errors.NewParam("data", data),
	)
}

func ErrInvalidProof(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidProof,
		fmt.Sprintf("provided proof is invalid %s", error),
		errors.NewParam("error", error),
	)
}

func ErrInvalidPassphrase(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidPassphrase,
		fmt.Sprintf("unable to create private key from passphrase: %s", error),
		errors.NewParam("error", error),
	)
}

func ErrInvalidChainID(wanted string, issued string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidChainID,
		fmt.Sprintf("wanted chain ID %s, but check is issued for chain with ID %s", wanted, issued),
		errors.NewParam("wanted", wanted),
		errors.NewParam("issued", issued),
	)
}

func ErrInvalidNonce() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidNonce,
		"nonce is too big (should be up to 16 bytes)",
	)
}

func ErrCheckExpired(block string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCheckExpired,
		fmt.Sprintf("check was expired at block %s", block),
		errors.NewParam("block", block),
	)
}

func ErrCheckRedeemed() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCheckRedeemed,
		"check was redeemed already",
	)
}

func ErrUnableDecodeCheck(check string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableDecodeCheck,
		fmt.Sprintf("unable to decode check from base58 %s", check),
		errors.NewParam("check", check),
	)
}

func ErrUnableRPLEncodeCheck(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableRPLEncodeCheck,
		fmt.Sprintf("unable to RLP encode check receiver address: %s", error),
		errors.NewParam("error", error),
	)
}

func ErrUnableSignCheck(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableSignCheck,
		fmt.Sprintf("unable to sign check receiver address by private key generated from passphrase: %s", error),
		errors.NewParam("error", error),
	)
}

func ErrUnableDecodeProof() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableDecodeProof,
		"unable to decode proof from base64",
	)
}

func ErrUnableRecoverAddress(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableRecoverAddress,
		fmt.Sprintf("unable to recover check issuer address: %s", error),
		errors.NewParam("error", error),
	)
}

func ErrUnableRecoverLockPkey(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableRecoverLockPkey,
		fmt.Sprintf("unable to recover lock public key from check: %s", error),
		errors.NewParam("error", error),
	)
}

// AccountKeys Errors

func ErrInvalidPkey() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidPkey,
		"invalid private key",
	)
}

func ErrUnableRetrieveArmoredPkey(name string, error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableRetriveArmoredPkey,
		fmt.Sprintf("unable to retrieve armored private key for account %s: %s", name, error),
		errors.NewParam("name", name),
		errors.NewParam("error", error),
	)
}

func ErrUnableRetrievePkey(name string, error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableRetrivePkey,
		fmt.Sprintf("unable to retrieve private key for account %s: %s", name, error),
		errors.NewParam("name", name),
		errors.NewParam("error", error),
	)
}

func ErrUnableRetrieveSECPPkey(name string, algo string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableRetriveSECPPkey,
		fmt.Sprintf("unable to retrieve secp256k1 private key for account %s: %s private key retrieved instead", name, algo),
		errors.NewParam("name", name),
		errors.NewParam("algo", algo),
	)
}
