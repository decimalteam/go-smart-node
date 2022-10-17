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
	NFTTokenNotFound                = errors.New(codespace, 110, "NFT does not exists")
	DelegationNotFound              = errors.New(codespace, 111, "delegation not found")
	DelegationWrongType             = errors.New(codespace, 112, "delegation has wrong type")
	DelegationTooSmall              = errors.New(codespace, 113, "delegation too small for undelegation/redelegation")
	SubTokenIDsDublicates           = errors.New(codespace, 114, "subtokes ID set has dublicates")
	StakeDoesNotHaveSubTokenID      = errors.New(codespace, 114, "stake does not have subtoken id")
	DelegateSubTokenTwice           = errors.New(codespace, 115, "trying to delegate subtoken id twice")
	BadRedelegationDst              = errors.New(codespace, 116, "redelegation destination validator not found")
	BadRedelegationSrc              = errors.New(codespace, 117, "redelegation source validator not found")
	SelfRedelegation                = errors.New(codespace, 118, "cannot redelegate to the same validator")
	TransitiveRedelegation          = errors.New(codespace, 119, "redelegation to this validator already in progress; first redelegation to this validator must complete before next redelegation")
	MaxRedelegationEntries          = errors.New(codespace, 120, "too many redelegation entries for (delegator, src-validator, dst-validator) tuple")
	IncompatibleBondStatuses        = errors.New(codespace, 121, "incompatible bond statuses")
	ValidatorStatusUnknown          = errors.New(codespace, 122, "validator status unknown")
	WrongStakeType                  = errors.New(codespace, 123, "wrong stake type")
	SubTokenExistsInStake           = errors.New(codespace, 124, "subtoken exists in stake")
	MaxUndelegationEntries          = errors.New(codespace, 125, "too many unbonding delegation entries for (delegator, validator) tuple")
)
