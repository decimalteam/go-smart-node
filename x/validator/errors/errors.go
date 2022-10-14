package errors

import (
	"cosmossdk.io/errors"
)

var codespace = "validator"

var (
	Internal                        = errors.New(codespace, 101, "internal error")
	ValidatorAlreadyExists          = errors.New(codespace, 102, "validator already exists")
	InvalidConsensusPubKey          = errors.New(codespace, 103, "invalid consensus public key")
	ValidatorPublicKeyAlreadyExists = errors.New(codespace, 104, "validator public key already exists")
	UnsupportedPubKeyType           = errors.New(codespace, 105, "unsupported public key type")
	ValidatorNotFound               = errors.New(codespace, 106, "validator not found")
	ValidatorAlreadyOnline          = errors.New(codespace, 107, "validator already online")
	ValidatorAlreadyOffline         = errors.New(codespace, 108, "validator already offline")
	NFTSubTokenNotFound             = errors.New(codespace, 109, "NFT subtoken does not exists")
	DelegatorIsNotOwnerOfSubtoken   = errors.New(codespace, 109, "delegator is not owner of NFT subtoken")
)
