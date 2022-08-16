package types

import (
	fmt "fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	ParamStoreKeyLockedTimeOut = []byte("LockedTimeOut")
	ParamStoreKeyLockedTimeIn  = []byte("LockedTimeIn")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(lockedTimeOut, lockedTimeIn time.Duration) Params {
	return Params{
		LockedTimeOut: lockedTimeOut,
		LockedTimeIn:  lockedTimeIn,
	}
}

func DefaultParams() Params {
	return NewParams(DefaultLockedTimeOut, DefaultLockedTimeIn)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyLockedTimeOut, &p.LockedTimeOut, validateLockedTime),
		paramtypes.NewParamSetPair(ParamStoreKeyLockedTimeIn, &p.LockedTimeIn, validateLockedTime),
	}
}

func validateLockedTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("locked time must be positive: %d", v)
	}

	return nil
}

func (p Params) Validate() error {
	return nil
}
