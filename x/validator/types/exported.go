package types

import (
	tmprotocrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"

	sdkmath "cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakeI defines interface for a delegation stake.
type StakeI interface {
	GetType() StakeType
	GetID() string
	GetStake() sdk.Coin
	GetSubTokenIDs() []int64
}

// DelegationI defines interface for a delegation bonded to a validator.
type DelegationI interface {
	GetDelegator() sdk.AccAddress
	GetValidator() sdk.ValAddress
	GetStake() StakeI
}

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
	GetConsensusPower(sdkmath.Int) int64               // validation power in tendermint
}
