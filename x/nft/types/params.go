package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

// NFT params default values.
var (
	DefaultMaxCollectionSize uint32      = 10000
	DefaultMaxTokenQuantity  uint32      = 10000
	DefaultMinReserveAmount  sdkmath.Int = helpers.EtherToWei(sdkmath.NewInt(1))
)

// Parameter store keys.
var (
	KeyMaxCollectionSize = []byte("MaxCollectionSize")
	KeyMaxTokenQuantity  = []byte("MaxTokenQuantity")
	KeyMinReserveAmount  = []byte("MinReserveAmount")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MaxCollectionSize: DefaultMaxCollectionSize,
		MaxTokenQuantity:  DefaultMaxTokenQuantity,
		MinReserveAmount:  DefaultMinReserveAmount,
	}
}

// NewParams creates a new Params instance.
func NewParams(maxCollectionSize uint32, maxTokenQuantity uint32, minReserveAmount sdkmath.Int) Params {
	return Params{
		MaxCollectionSize: maxCollectionSize,
		MaxTokenQuantity:  maxTokenQuantity,
		MinReserveAmount:  minReserveAmount,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxCollectionSize, &p.MaxCollectionSize, validateMaxCollectionSize),
		paramtypes.NewParamSetPair(KeyMaxTokenQuantity, &p.MaxTokenQuantity, validateMaxTokenQuantity),
		paramtypes.NewParamSetPair(KeyMinReserveAmount, &p.MinReserveAmount, validateMinReserveAmount),
	}
}

// Validate validates the set of params.
func (p *Params) Validate() (err error) {
	if err = validateMaxCollectionSize(p.MaxCollectionSize); err != nil {
		return
	}
	if err = validateMaxTokenQuantity(p.MaxTokenQuantity); err != nil {
		return
	}
	if err = validateMinReserveAmount(p.MinReserveAmount); err != nil {
		return
	}
	return
}

func validateMaxCollectionSize(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("max collection size must be positive: %d", v)
	}
	return nil
}

func validateMaxTokenQuantity(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("max token quantity must be positive: %d", v)
	}
	return nil
}

func validateMinReserveAmount(i interface{}) error {
	_, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	// TODO
	return nil
}
