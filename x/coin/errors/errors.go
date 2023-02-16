package errors

import (
	"fmt"

	"cosmossdk.io/errors"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/config"
)

var codespace = "coin"

var (
	CoinDoesNotExist              = errors.New(codespace, 101, "coin does not exist")
	CoinAlreadyExists             = errors.New(codespace, 102, "coin already exist")
	InvalidCRR                    = errors.New(codespace, 103, "CRR must be between 10 and 100")
	InvalidCoinDenom              = errors.New(codespace, 104, fmt.Sprintf("symbol must match this regular expression: %s", config.CoinDenomRegExp))
	InvalidCoinTitle              = errors.New(codespace, 105, fmt.Sprintf("invalid coin title. Allowed up to %d bytes", config.MaxCoinTitleLength))
	InvalidCoinInitialVolume      = errors.New(codespace, 106, fmt.Sprintf("coin initial volume should be between %s and %s.", config.MinCoinSupply, config.MaxCoinSupply))
	InvalidCoinInitialReserve     = errors.New(codespace, 107, fmt.Sprintf("coin initial reserve should be greater than or equal to reserve %s", config.MinCoinReserve))
	NewLimitVolumeLess            = errors.New(codespace, 108, "new limit volume should be grater than volume")
	UpdateOnlyForCreator          = errors.New(codespace, 109, "updating allowed only for creator of coin")
	InvalidLimitVolume            = errors.New(codespace, 110, fmt.Sprintf("volume limit should be greater than initial volume and less or equal than %s", config.MaxCoinSupply))
	InvalidAmount                 = errors.New(codespace, 111, "amount should be greater than 0")
	SameCoin                      = errors.New(codespace, 112, "can't operating same coins")
	InvalidRecipientAddress       = errors.New(codespace, 113, "invalid recipient address:")
	InvalidSenderAddress          = errors.New(codespace, 114, "invalid sender address")
	CheckDoesNotExist             = errors.New(codespace, 115, "check does not exist")
	InvalidCheckSig               = errors.New(codespace, 116, "invalid transaction v, r, s values")
	InvalidProof                  = errors.New(codespace, 117, "provided proof is invalid")
	InvalidPassphrase             = errors.New(codespace, 118, "unable to create private key from passphrase")
	InvalidChainID                = errors.New(codespace, 119, "received invalid chain ID")
	InvalidNonce                  = errors.New(codespace, 120, "nonce is too big (should be up to 16 bytes)")
	CheckExpired                  = errors.New(codespace, 121, "check was expired")
	CheckRedeemed                 = errors.New(codespace, 122, "check was redeemed already")
	UnableDecodeCheckBase58       = errors.New(codespace, 123, "unable to decode check from base58")
	UnableRLPEncodeAddress        = errors.New(codespace, 124, "unable to encode address in rpl")
	UnableSignCheck               = errors.New(codespace, 125, "unable to sign check receiver address by private key generated from received passphrase")
	UnableDecodeProofBase64       = errors.New(codespace, 126, "unable to decode proof from base64")
	UnableRecoverAddressFromCheck = errors.New(codespace, 127, "unable to recover issuer address from check")
	UnableRecoverLockPkey         = errors.New(codespace, 128, "unable to recover lock public key from check")
	FailedToRecoverPKFromSig      = errors.New(codespace, 129, "failed to recover the pub key from the signature")
	InvalidPubKey                 = errors.New(codespace, 130, "pub key isn't valid")
	TxBreaksVolumeLimit           = errors.New(codespace, 131, "tx breaks coin LimitVolume rule: volume < limit volume")
	TxBreaksMinVolumeLimit        = errors.New(codespace, 132, "tx breaks min volume rule: volume > min volume")
	TxBreaksMinReserveRule        = errors.New(codespace, 133, fmt.Sprintf("tx breaks MinReserveLimit rule: reserve > %s", config.MinCoinReserve))
	MaximumValueToSellReached     = errors.New(codespace, 134, "wanted limit amount of coins for sale is less than it actually took")
	MinimumValueToBuyReached      = errors.New(codespace, 135, "wanted minimum amount to buy is less than actually amount")
	InsufficientCoinReserve       = errors.New(codespace, 136, "coin reserve balance is not sufficient for transaction")
	InsufficientFunds             = errors.New(codespace, 137, "wallet not enough funds")
	DecodeRLP                     = errors.New(codespace, 138, "failed to parse RPL bytes")
	UnableRLPEncodeCheck          = errors.New(codespace, 139, "unable to RPL encode check")
	UnableRLPEncodeToBytesCheck   = errors.New(codespace, 140, "unable to RPL encode check to bytes")
	DuplicateCoinInGenesis        = errors.New(codespace, 141, "coin symbol duplicated on genesis")
	DuplicateCheckInGenesis       = errors.New(codespace, 142, "check hash duplicated on genesis")
	Internal                      = errors.New(codespace, 143, "internal error")
	InvalidCoinMinEmission        = errors.New(codespace, 144, fmt.Sprintf("coin min emission should be between %s and %s.", config.MinCoinSupply, config.MaxCoinSupply))
	UneditableCoinMinEmission     = errors.New(codespace, 145, "coin min emission cannot be enabled or disabled after coin creation")
	TxBreaksMinEmission           = errors.New(codespace, 146, "tx breaks coin min emission rule: volume > min emission")
	TooBigCoinMinEmission         = errors.New(codespace, 147, fmt.Sprintf("coin min emission should be between %s and %s.", config.MinCoinSupply, "volume limit"))
)
