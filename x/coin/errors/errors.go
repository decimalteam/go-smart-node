package errors

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var codespace = "coin"

const (
	allowedCoinSymbols = "^[a-zA-Z][a-zA-Z0-9]{2,9}$"
	maxCoinNameBytes   = 64
)

var (
	minCoinSupply  = helpers.EtherToWei(sdk.NewInt(1))
	maxCoinSupply  = helpers.EtherToWei(sdk.NewInt(1000000000000000))
	minCoinReserve = helpers.EtherToWei(sdk.NewInt(1000))
)

var (
	InvalidCRR                = errors.New(codespace, 1, "CRR must be between 10 and 100")
	CoinDoesNotExist          = errors.New(codespace, 2, "coin does not exist")
	InvalidCoinSymbol         = errors.New(codespace, 3, fmt.Sprintf("symbol must match this regular expression: %s", allowedCoinSymbols))
	CoinAlreadyExists         = errors.New(codespace, 4, "coin already exist")
	InvalidCoinTitle          = errors.New(codespace, 5, fmt.Sprintf("invalid coin title. Allowed up to %d bytes", maxCoinNameBytes))
	InvalidCoinInitialVolume  = errors.New(codespace, 6, fmt.Sprintf("coin initial volume should be between %s and %s.", minCoinSupply.String(), maxCoinSupply.String()))
	InvalidCoinInitialReserve = errors.New(codespace, 7, fmt.Sprintf("coin initial reserve should be greater than or equal to reserve %s", minCoinReserve))
	Internal                  = errors.New(codespace, 8, "internal error")
	InsufficientCoinReserve   = errors.New(codespace, 9, "not enough coin to coin reserve")
	UpdateOnlyForCreator      = errors.New(codespace, 11, "updating allowed only for creator of coin")
	SameCoin                  = errors.New(codespace, 12, "can't operating same coins")
	TxBreaksMinVolumeLimit    = errors.New(codespace, 13, "tx breaks min volume rule: volume > min volume")
	TxBreaksMinReserveRule    = errors.New(codespace, 14, fmt.Sprintf("tx breaks MinReserveLimit rule: reserve > %s", minCoinReserve))
	MaximumValueToSellReached = errors.New(codespace, 15, "wanted limit amount of coins for sale is less than it actually took")
	MinimumValueToBuyReached  = errors.New(codespace, 16, "wanted minimum amount to buy is less than actually amount")
	LimitVolumeBroken         = errors.New(codespace, 17, "volume should be less than or equal the volume limit")
	InvalidAmount             = errors.New(codespace, 18, "amount should be greater than 0")
	InvalidReceiverAddress    = errors.New(codespace, 19, "invalid receiver address:")
	InvalidSenderAddress      = errors.New(codespace, 20, "invalid sender address")
	TxBreaksVolumeLimit       = errors.New(codespace, 21, "tx breaks coin LimitVolume rule: volume < limit volume")
	InvalidCheck              = errors.New(codespace, 22, "unable to parse check")
	InvalidProof              = errors.New(codespace, 23, "provided proof is invalid")
	InvalidPassphrase         = errors.New(codespace, 24, "unable to create private key from passphrase")
	InvalidChainID            = errors.New(codespace, 25, "received invalid chain ID")
	InvalidNonce              = errors.New(codespace, 26, "nonce is too big (should be up to 16 bytes)")
	CheckExpired              = errors.New(codespace, 27, "check was expired")
	CheckRedeemed             = errors.New(codespace, 28, "check was redeemed already")
	UnableDecodeCheck         = errors.New(codespace, 29, "unable to decode check from base58")
	UnableRPLEncodeCheck      = errors.New(codespace, 30, "unable to RLP encode check receiver address")
	UnableSignCheck           = errors.New(codespace, 31, "unable to sign check receiver address by private key generated from received passphrase")
	UnableDecodeProof         = errors.New(codespace, 32, "unable to decode proof from base64")
	UnableRecoverAddress      = errors.New(codespace, 33, "unable to recover check from issuer address")
	UnableRecoverLockPkey     = errors.New(codespace, 34, "unable to recover lock public key from check")

	FailedSend          = errors.New(codespace, 35, "an error occurred while sending tokens")
	InsufficientFunds   = errors.New(codespace, 36, "wallet not enough funds")
	CalculateCommission = errors.New(codespace, 37, "")
	FailedBurnCoins     = errors.New(codespace, 38, "failed to burn coins")
	FailedMintCoins     = errors.New(codespace, 39, "failed to mint coins")
)
