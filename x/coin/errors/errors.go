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
	maxCoinNameBytes   = 46
)

var (
	minCoinSupply  = helpers.EtherToWei(sdk.NewInt(1))
	maxCoinSupply  = helpers.EtherToWei(sdk.NewInt(1000000000000000))
	minCoinReserve = helpers.EtherToWei(sdk.NewInt(1000))
)

var (
	CoinDoesNotExist              = errors.New(codespace, 1, "coin does not exist")
	CoinAlreadyExists             = errors.New(codespace, 2, "coin already exist")
	InvalidCRR                    = errors.New(codespace, 3, "CRR must be between 10 and 100")
	InvalidCoinSymbol             = errors.New(codespace, 4, fmt.Sprintf("symbol must match this regular expression: %s", allowedCoinSymbols))
	InvalidCoinTitle              = errors.New(codespace, 5, fmt.Sprintf("invalid coin title. Allowed up to %d bytes", maxCoinNameBytes))
	InvalidCoinInitialVolume      = errors.New(codespace, 6, fmt.Sprintf("coin initial volume should be between %s and %s.", minCoinSupply.String(), maxCoinSupply.String()))
	InvalidCoinInitialReserve     = errors.New(codespace, 7, fmt.Sprintf("coin initial reserve should be greater than or equal to reserve %s", minCoinReserve))
	NewLimitVolumeLess            = errors.New(codespace, 8, "new limit volume should be grater than old limit volume")
	UpdateOnlyForCreator          = errors.New(codespace, 9, "updating allowed only for creator of coin")
	InvalidLimitVolume            = errors.New(codespace, 10, fmt.Sprintf("volume limit should be less or equal than %s", maxCoinSupply.String()))
	InvalidAmount                 = errors.New(codespace, 11, "amount should be greater than 0")
	SameCoin                      = errors.New(codespace, 12, "can't operating same coins")
	InvalidReceiverAddress        = errors.New(codespace, 13, "invalid receiver address:")
	InvalidSenderAddress          = errors.New(codespace, 14, "invalid sender address")
	CheckDoesNotExist             = errors.New(codespace, 15, "check does not exist")
	InvalidCheckSig               = errors.New(codespace, 16, "invalid transaction v, r, s values")
	InvalidProof                  = errors.New(codespace, 17, "provided proof is invalid")
	InvalidPassphrase             = errors.New(codespace, 18, "unable to create private key from passphrase")
	InvalidChainID                = errors.New(codespace, 19, "received invalid chain ID")
	InvalidNonce                  = errors.New(codespace, 20, "nonce is too big (should be up to 16 bytes)")
	CheckExpired                  = errors.New(codespace, 21, "check was expired")
	CheckRedeemed                 = errors.New(codespace, 22, "check was redeemed already")
	UnableDecodeCheckBase58       = errors.New(codespace, 23, "unable to decode check from base58")
	UnableRPLEncodeAddress        = errors.New(codespace, 24, "unable to encode address in rpl")
	UnableSignCheck               = errors.New(codespace, 25, "unable to sign check receiver address by private key generated from received passphrase")
	UnableDecodeProofBase64       = errors.New(codespace, 26, "unable to decode proof from base64")
	UnableRecoverAddressFromCheck = errors.New(codespace, 27, "unable to recover issuer address from check")
	UnableRecoverLockPkey         = errors.New(codespace, 28, "unable to recover lock public key from check")
	FailedToRecoverPKFromSig      = errors.New(codespace, 29, "failed to recover the pub key from the signature")
	InvalidPubKey                 = errors.New(codespace, 30, "pub key isn't valid")
	TxBreaksVolumeLimit           = errors.New(codespace, 31, "tx breaks coin LimitVolume rule: volume < limit volume")
	TxBreaksMinVolumeLimit        = errors.New(codespace, 32, "tx breaks min volume rule: volume > min volume")
	TxBreaksMinReserveRule        = errors.New(codespace, 33, fmt.Sprintf("tx breaks MinReserveLimit rule: reserve > %s", minCoinReserve))
	MaximumValueToSellReached     = errors.New(codespace, 34, "wanted limit amount of coins for sale is less than it actually took")
	MinimumValueToBuyReached      = errors.New(codespace, 35, "wanted minimum amount to buy is less than actually amount")
	InsufficientCoinReserve       = errors.New(codespace, 36, "coin reserve balance is not sufficient for transaction")
	InsufficientFunds             = errors.New(codespace, 37, "wallet not enough funds")
	DecodeRLP                     = errors.New(codespace, 38, "failed to parse rlp bytes")
	UnableRPLEncodeCheck          = errors.New(codespace, 39, "unable to RPL encode check")
	UnableRPLEncodeToBytesCheck   = errors.New(codespace, 40, "unable to RPL encode check to bytes")
	DuplicateCoinInGenesis        = errors.New(codespace, 41, "coin symbol duplicated on genesis")
	DuplicateCheckInGenesis       = errors.New(codespace, 42, "check hash duplicated on genesis")
	Internal                      = errors.New(codespace, 43, "internal error")
)
