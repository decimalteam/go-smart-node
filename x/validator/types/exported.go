package types

import (
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorI defines interface for a validator.
type ValidatorI interface {
	GetOperator() sdk.ValAddress                       // operator address to receive/return validators coins
	GetConsAddr() (sdk.ConsAddress, error)             // validation consensus address
	GetMoniker() string                                // moniker of the validator
	GetCommission() sdk.Dec                            // validator commission rate
	GetStatus() BondStatus                             // status of the validator
	IsJailed() bool                                    // whether the validator is jailed
	IsBonded() bool                                    // check if has a bonded status
	IsUnbonded() bool                                  // check if has status unbonded
	IsUnbonding() bool                                 // check if has status unbonding
	ConsPubKey() (cryptotypes.PubKey, error)           // validation consensus pubkey (cryptotypes.PubKey)
	TmConsPublicKey() (tmprotocrypto.PublicKey, error) // validation consensus pubkey (Tendermint)
	ConsensusPower() int64                             // validation power in tendermint
}

// DelegationI defines interface for a delegation bonded to a validator.
type DelegationI interface {
	GetDelegator() sdk.AccAddress
	GetValidator() sdk.ValAddress
	GetStake() StakeI
}

// RedelegationI defines interface for a redelegation from one validator to another.
type RedelegationI interface {
	GetDelegator() sdk.AccAddress
	GetValidatorSrc() sdk.ValAddress
	GetValidatorDst() sdk.ValAddress
	GetEntries() []RedelegationEntry
}

// UndelegationI defines interface for a undelegation from a validator.
type UndelegationI interface {
	GetDelegator() sdk.AccAddress
	GetValidator() sdk.ValAddress
	GetEntries() []UndelegationEntry
}

// StakeI defines interface for a delegation stake.
type StakeI interface {
	GetType() StakeType
	GetID() string
	GetStake() sdk.Coin
	GetSubTokenIDs() []uint32
}
