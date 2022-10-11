package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Validator params default values.
const (
	// Default maximum number of bonded validators.
	DefaultMaxValidators uint32 = 256

	// Default maximum number of bonded delegations per validator.
	DefaultMaxDelegations uint16 = 1000

	// Default maximum entries in a UBD/RED pair.
	DefaultMaxEntries uint32 = 7

	// DefaultHistorical entries is 10000. Apps that don't use IBC can ignore this value by not
	// adding the validator module to the application module manager's SetOrderBeginBlockers.
	DefaultHistoricalEntries uint32 = 10000

	// DefaultRedelegationTime reflects a week in nanoseconds as the default redelegating time.
	DefaultRedelegationTime time.Duration = time.Hour * 24 * 7

	// DefaultUndelegationTime reflects a month in nanoseconds as the default unbonding time.
	DefaultUndelegationTime time.Duration = time.Hour * 24 * 7 * 4
)

var (
	KeyMaxValidators     = []byte("MaxValidators")
	KeyMaxDelegations    = []byte("MaxDelegations")
	KeyMaxEntries        = []byte("MaxEntries")
	KeyHistoricalEntries = []byte("HistoricalEntries")
	KeyRedelegationTime  = []byte("RedelegationTime")
	KeyUndelegationTime  = []byte("UndelegationTime")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MaxValidators:     DefaultMaxValidators,
		MaxDelegations:    0,
		MaxEntries:        DefaultMaxEntries,
		HistoricalEntries: DefaultHistoricalEntries,
		RedelegationTime:  DefaultRedelegationTime,
		UndelegationTime:  DefaultUndelegationTime,
	}
}

// NewParams creates a new Params instance.
func NewParams(
	maxValidators uint32,
	maxDelegations uint32,
	maxEntries uint32,
	historicalEntries uint32,
	redelegationTime time.Duration,
	undelegationTime time.Duration,
) Params {
	return Params{
		MaxValidators:     maxValidators,
		MaxDelegations:    maxDelegations,
		MaxEntries:        maxEntries,
		HistoricalEntries: historicalEntries,
		RedelegationTime:  redelegationTime,
		UndelegationTime:  undelegationTime,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxValidators, &p.MaxValidators, validateMaxValidators),
		paramtypes.NewParamSetPair(KeyMaxDelegations, &p.MaxDelegations, validateMaxDelegations),
		paramtypes.NewParamSetPair(KeyMaxEntries, &p.MaxEntries, validateMaxEntries),
		paramtypes.NewParamSetPair(KeyHistoricalEntries, &p.HistoricalEntries, validateHistoricalEntries),
		paramtypes.NewParamSetPair(KeyRedelegationTime, &p.RedelegationTime, validateRedelegationTime),
		paramtypes.NewParamSetPair(KeyUndelegationTime, &p.UndelegationTime, validateUndelegationTime),
	}
}

// Validate validates the set of params.
func (p Params) Validate() (err error) {
	if err = validateMaxValidators(p.MaxValidators); err != nil {
		return
	}
	if err = validateMaxDelegations(p.MaxDelegations); err != nil {
		return
	}
	if err = validateMaxEntries(p.MaxEntries); err != nil {
		return
	}
	if err = validateHistoricalEntries(p.HistoricalEntries); err != nil {
		return
	}
	if err = validateRedelegationTime(p.RedelegationTime); err != nil {
		return
	}
	if err = validateUndelegationTime(p.UndelegationTime); err != nil {
		return
	}
	return
}

func validateMaxValidators(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("max validators must be positive: %d", v)
	}
	return nil
}

func validateMaxDelegations(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("max delegations must be positive: %d", v)
	}
	return nil
}

func validateMaxEntries(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("max entries must be positive: %d", v)
	}
	return nil
}

func validateHistoricalEntries(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateRedelegationTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("redelegation time must be positive: %d", v)
	}
	return nil
}

func validateUndelegationTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("undelegation time must be positive: %d", v)
	}
	return nil
}
