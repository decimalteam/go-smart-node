package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	coinconfig "bitbucket.org/decimalteam/go-smart-node/x/coin/config"
)

// Validator params default values.
const (
	// Default maximum number of bonded validators.
	DefaultMaxValidators uint32 = 256

	// Default maximum number of bonded delegations per validator.
	DefaultMaxDelegations uint32 = 1000

	// Default maximum entries in a UBD/RED pair.
	DefaultMaxEntries uint32 = 7

	// DefaultHistorical entries is 10000. Apps that don't use IBC can ignore this value by not
	// adding the validator module to the application module manager's SetOrderBeginBlockers.
	DefaultHistoricalEntries uint32 = 10000

	// DefaultRedelegationTime reflects a week in nanoseconds as the default redelegating time.
	DefaultRedelegationTime time.Duration = time.Hour * 24 * 7

	// DefaultUndelegationTime reflects a month in nanoseconds as the default unbonding time.
	DefaultUndelegationTime time.Duration = time.Hour * 24 * 7 * 4

	// DefaultSignedBlocksWindow 24 * ~5 sec = ~120 sec window
	DefaultSignedBlocksWindow int64 = 24

	DefaultBaseDenom = cmdcfg.BaseDenom
)

var (
	// DefaultMinSignedPerWindow 0.5 of 24 blocks = ~ 60 sec
	DefaultMinSignedPerWindow = sdk.NewDec(1).Quo(sdk.NewDec(2))
	// DefaultSlashFractionDowntime 1% of stake
	DefaultSlashFractionDowntime = sdk.NewDec(1).Quo(sdk.NewDec(100))
	// DefaultSlashFractionDoubleSign 5% of stake
	DefaultSlashFractionDoubleSign = sdk.NewDec(1).Quo(sdk.NewDec(20))
)

var (
	KeyMaxValidators     = []byte("MaxValidators")
	KeyMaxDelegations    = []byte("MaxDelegations")
	KeyMaxEntries        = []byte("MaxEntries")
	KeyHistoricalEntries = []byte("HistoricalEntries")
	KeyRedelegationTime  = []byte("RedelegationTime")
	KeyUndelegationTime  = []byte("UndelegationTime")
	KeyBaseDenom         = []byte("BaseDenom")

	KeySignedBlocksWindow      = []byte("SignedBlocksWindow")
	KeyMinSignedPerWindow      = []byte("MinSignedPerWindow")
	KeySlashFractionDowntime   = []byte("SlashFractionDowntime")
	KeySlashFractionDoubleSign = []byte("SlashFractionDoubleSign")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MaxValidators:           DefaultMaxValidators,
		MaxDelegations:          DefaultMaxDelegations,
		MaxEntries:              DefaultMaxEntries,
		HistoricalEntries:       DefaultHistoricalEntries,
		RedelegationTime:        DefaultRedelegationTime,
		UndelegationTime:        DefaultUndelegationTime,
		BaseDenom:               DefaultBaseDenom,
		SignedBlocksWindow:      DefaultSignedBlocksWindow,
		MinSignedPerWindow:      DefaultMinSignedPerWindow,
		SlashFractionDowntime:   DefaultSlashFractionDowntime,
		SlashFractionDoubleSign: DefaultSlashFractionDoubleSign,
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
	baseDenom string,
	signedBlockWindow int64,
	minSignedPerWindow sdk.Dec,
	slashFractionDowntime sdk.Dec,
	slashFractionDoubleSign sdk.Dec,
) Params {
	return Params{
		MaxValidators:           maxValidators,
		MaxDelegations:          maxDelegations,
		MaxEntries:              maxEntries,
		HistoricalEntries:       historicalEntries,
		RedelegationTime:        redelegationTime,
		UndelegationTime:        undelegationTime,
		BaseDenom:               baseDenom,
		SignedBlocksWindow:      signedBlockWindow,
		MinSignedPerWindow:      minSignedPerWindow,
		SlashFractionDowntime:   slashFractionDowntime,
		SlashFractionDoubleSign: slashFractionDoubleSign,
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
		paramtypes.NewParamSetPair(KeyBaseDenom, &p.BaseDenom, validateBaseDenom),
		paramtypes.NewParamSetPair(KeySignedBlocksWindow, &p.SignedBlocksWindow, validateSignedBlockWindow),
		paramtypes.NewParamSetPair(KeyMinSignedPerWindow, &p.MinSignedPerWindow, validateDec),
		paramtypes.NewParamSetPair(KeySlashFractionDowntime, &p.SlashFractionDowntime, validateDec),
		paramtypes.NewParamSetPair(KeySlashFractionDoubleSign, &p.SlashFractionDoubleSign, validateDec),
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
	if err = validateBaseDenom(p.BaseDenom); err != nil {
		return
	}
	if err = validateSignedBlockWindow(p.SignedBlocksWindow); err != nil {
		return
	}
	if err = validateDec(p.MinSignedPerWindow); err != nil {
		return
	}
	if err = validateDec(p.SlashFractionDowntime); err != nil {
		return
	}
	if err = validateDec(p.SlashFractionDoubleSign); err != nil {
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

func validateSignedBlockWindow(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < DefaultSignedBlocksWindow {
		return fmt.Errorf("signed block window too small: %d < %d", v, DefaultSignedBlocksWindow)
	}
	return nil
}

func validateDec(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() || !v.IsPositive() {
		return fmt.Errorf("wrong sdk.Dec value")
	}
	return nil
}

func validateBaseDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == "" || !coinconfig.CoinDenomValidator.MatchString(v) {
		return fmt.Errorf("wrong denom value")
	}

	return nil
}
