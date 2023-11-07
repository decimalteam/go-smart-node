package types

import (
	"bytes"
	"sort"
	"strings"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	MaxMonikerLength         = 70
	MaxIdentityLength        = 3000
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
)

var (
	BondStatusUnspecified = BondStatus_name[int32(BondStatus_Unspecified)]
	BondStatusUnbonded    = BondStatus_name[int32(BondStatus_Unbonded)]
	BondStatusUnbonding   = BondStatus_name[int32(BondStatus_Unbonding)]
	BondStatusBonded      = BondStatus_name[int32(BondStatus_Bonded)]
)

var _ ValidatorI = Validator{}

// NewValidator creates new Validator instance.
func NewValidator(
	operatorAddress sdk.ValAddress,
	rewardAddress sdk.AccAddress,
	pubKey cryptotypes.PubKey,
	description Description,
	commission sdk.Dec,
) (Validator, error) {
	pkAny, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return Validator{}, err
	}
	return Validator{
		OperatorAddress: operatorAddress.String(),
		RewardAddress:   rewardAddress.String(),
		ConsensusPubkey: pkAny,
		Description:     description,
		Commission:      commission,
		Status:          BondStatus_Unbonded,
		Online:          false,
		Jailed:          false,
		UnbondingHeight: 0,
		UnbondingTime:   time.Unix(0, 0).UTC(),
		Rewards:         sdkmath.ZeroInt(),
		TotalRewards:    sdkmath.ZeroInt(),
		Stake:           0,
	}, nil
}

// Validators is a collection of Validator
type Validators []Validator

func (v Validators) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}

	return strings.TrimSpace(out)
}

// ToSDKValidators -  convenience function convert []Validator to []sdk.ValidatorI
func (v Validators) ToSDKValidators() (validators []ValidatorI) {
	for _, val := range v {
		validators = append(validators, val)
	}

	return validators
}

// Sort Validators sorts validator array in ascending operator address order
func (v Validators) Sort() {
	sort.Sort(v)
}

// Implements sort interface
func (v Validators) Len() int {
	return len(v)
}

// Implements sort interface
func (v Validators) Less(i, j int) bool {
	return bytes.Compare(v[i].GetOperator().Bytes(), v[j].GetOperator().Bytes()) == -1
}

// Implements sort interface
func (v Validators) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// ValidatorsByVotingPower implements sort.Interface for []Validator based on
// the VotingPower and Address fields.
// The validators are sorted first by their voting power (descending). Secondary index - Address (ascending).
// Copied from tendermint/types/validator_set.go
type ValidatorsByVotingPower []Validator

func (valz ValidatorsByVotingPower) Len() int { return len(valz) }

func (valz ValidatorsByVotingPower) Less(i, j int) bool {
	if valz[i].ConsensusPower() == valz[j].ConsensusPower() {
		addrI, errI := valz[i].GetConsAddr()
		addrJ, errJ := valz[j].GetConsAddr()
		// If either returns error, then return false
		if errI != nil || errJ != nil {
			return false
		}
		return bytes.Compare(addrI, addrJ) == -1
	}
	return valz[i].ConsensusPower() > valz[j].ConsensusPower()
}

func (valz ValidatorsByVotingPower) Swap(i, j int) {
	valz[i], valz[j] = valz[j], valz[i]
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v Validators) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range v {
		if err := v[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}

func MustMarshalValidator(cdc codec.BinaryCodec, validator *Validator) []byte {
	return cdc.MustMarshal(validator)
}

func MustUnmarshalValidator(cdc codec.BinaryCodec, value []byte) Validator {
	validator, err := UnmarshalValidator(cdc, value)
	if err != nil {
		panic(err)
	}

	return validator
}

// unmarshal a redelegation from a store value
func UnmarshalValidator(cdc codec.BinaryCodec, value []byte) (v Validator, err error) {
	err = cdc.Unmarshal(value, &v)
	return v, err
}

func MustMarshalValidatorRewards(cdc codec.BinaryCodec, rewards *ValidatorRS) []byte {
	return cdc.MustMarshal(rewards)
}

func MustUnmarshalValidatorRewards(cdc codec.BinaryCodec, value []byte) ValidatorRS {
	validator, err := UnmarshalValidatorRewards(cdc, value)
	if err != nil {
		panic(err)
	}

	return validator
}

// unmarshal a redelegation from a store value
func UnmarshalValidatorRewards(cdc codec.BinaryCodec, value []byte) (v ValidatorRS, err error) {
	err = cdc.Unmarshal(value, &v)
	return v, err
}

// IsBonded checks if the validator status equals Bonded
func (v Validator) IsBonded() bool {
	return v.GetStatus() == BondStatus_Bonded
}

// IsUnbonded checks if the validator status equals Unbonded
func (v Validator) IsUnbonded() bool {
	return v.GetStatus() == BondStatus_Unbonded
}

// IsUnbonding checks if the validator status equals Unbonding
func (v Validator) IsUnbonding() bool {
	return v.GetStatus() == BondStatus_Unbonding
}

func (v Validator) GetBondedTokens() sdk.Int {
	// TODO: implement
	return sdk.ZeroInt()
}

func (v Validator) GetConsensusPower(sdk.Int) int64 {
	// TODO: implement
	return 0
}

// constant used in flags to indicate that description field should not be updated
const DoNotModifyDesc = "[do-not-modify]"

func NewDescription(moniker, identity, website, securityContact, details string) Description {
	return Description{
		Moniker:         moniker,
		Identity:        identity,
		Website:         website,
		SecurityContact: securityContact,
		Details:         details,
	}
}

// // String implements the Stringer interface for a Description object.
// func (d Description) String() string {
// 	out, _ := yaml.Marshal(d)
// 	return string(out)
// }

// UpdateDescription updates the fields of a given description. An error is
// returned if the resulting description contains an invalid length.
func (d Description) UpdateDescription(d2 Description) (Description, error) {
	if d2.Moniker == DoNotModifyDesc {
		d2.Moniker = d.Moniker
	}

	if d2.Identity == DoNotModifyDesc {
		d2.Identity = d.Identity
	}

	if d2.Website == DoNotModifyDesc {
		d2.Website = d.Website
	}

	if d2.SecurityContact == DoNotModifyDesc {
		d2.SecurityContact = d.SecurityContact
	}

	if d2.Details == DoNotModifyDesc {
		d2.Details = d.Details
	}

	return NewDescription(
		d2.Moniker,
		d2.Identity,
		d2.Website,
		d2.SecurityContact,
		d2.Details,
	).EnsureLength()
}

// EnsureLength ensures the length of a validator's description.
func (d Description) EnsureLength() (Description, error) {
	if len(d.Moniker) > MaxMonikerLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid moniker length; got: %d, max: %d", len(d.Moniker), MaxMonikerLength)
	}

	if len(d.Identity) > MaxIdentityLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid identity length; got: %d, max: %d", len(d.Identity), MaxIdentityLength)
	}

	if len(d.Website) > MaxWebsiteLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid website length; got: %d, max: %d", len(d.Website), MaxWebsiteLength)
	}

	if len(d.SecurityContact) > MaxSecurityContactLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid security contact length; got: %d, max: %d", len(d.SecurityContact), MaxSecurityContactLength)
	}

	if len(d.Details) > MaxDetailsLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid details length; got: %d, max: %d", len(d.Details), MaxDetailsLength)
	}

	return d, nil
}

// ABCIValidatorUpdate returns an abci.ValidatorUpdate from a staking validator type with the full validator power.
func (v Validator) ABCIValidatorUpdate(r sdkmath.Int) abci.ValidatorUpdate {
	tmProtoPk, err := v.TmConsPublicKey()
	if err != nil {
		panic(err)
	}

	return abci.ValidatorUpdate{
		PubKey: tmProtoPk,
		Power:  v.ConsensusPower(),
	}
}

// ABCIValidatorUpdateZero returns an abci.ValidatorUpdate from a staking validator type with zero power used for validator updates.
func (v Validator) ABCIValidatorUpdateZero() abci.ValidatorUpdate {
	tmProtoPk, err := v.TmConsPublicKey()
	if err != nil {
		panic(err)
	}

	return abci.ValidatorUpdate{
		PubKey: tmProtoPk,
		Power:  0,
	}
}

// ConsensusPower gets the consensus-engine power. Aa reduction of 10^6 from validator tokens is applied.
func (v Validator) ConsensusPower() int64 {
	if v.IsBonded() {
		return v.PotentialConsensusPower()
	}

	return 0
}

// PotentialConsensusPower returns the potential consensus-engine power.
func (v Validator) PotentialConsensusPower() int64 {
	return v.Stake
}

// UpdateStatus updates the location of the shares within a validator to reflect the new status.
func (v Validator) UpdateStatus(newStatus BondStatus) Validator {
	v.Status = newStatus
	return v
}

// MinEqual defines a more minimum set of equality conditions when comparing two validators.
func (v *Validator) MinEqual(other *Validator) bool {
	return v.OperatorAddress == other.OperatorAddress &&
		v.RewardAddress == other.RewardAddress &&
		v.ConsensusPubkey.Equal(other.ConsensusPubkey) &&
		v.Description.Equal(other.Description) &&
		v.Commission.Equal(other.Commission) &&
		v.Status == other.Status &&
		v.Jailed == other.Jailed
}

// Equal checks if the receiver equals the parameter.
func (v *Validator) Equal(v2 *Validator) bool {
	return v.MinEqual(v2) &&
		v.UnbondingHeight == v2.UnbondingHeight &&
		v.UnbondingTime.Equal(v2.UnbondingTime)
}

func (v Validator) IsJailed() bool        { return v.Jailed }
func (v Validator) GetMoniker() string    { return v.Description.Moniker }
func (v Validator) GetStatus() BondStatus { return v.Status }
func (v Validator) GetOperator() sdk.ValAddress {
	if len(v.OperatorAddress) == 0 {
		return nil
	}
	addr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// ConsPubKey returns the validator PubKey as a cryptotypes.PubKey.
func (v Validator) ConsPubKey() (cryptotypes.PubKey, error) {
	pk, ok := v.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "expecting cryptotypes.PubKey, got %T", pk)
	}

	return pk, nil
}

// TmConsPublicKey casts Validator.ConsensusPubkey to tmprotocrypto.PubKey.
func (v Validator) TmConsPublicKey() (tmprotocrypto.PublicKey, error) {
	pk, err := v.ConsPubKey()
	if err != nil {
		return tmprotocrypto.PublicKey{}, err
	}
	tmPk, err := cryptocodec.ToTmProtoPublicKey(pk)
	if err != nil {
		return tmprotocrypto.PublicKey{}, err
	}
	return tmPk, nil
}

// GetConsAddr extracts Consensus key address
func (v Validator) GetConsAddr() (sdk.ConsAddress, error) {
	pk, ok := v.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "expecting cryptotypes.PubKey, got %T", pk)
	}
	return sdk.ConsAddress(pk.Address()), nil
}

func (v Validator) GetCommission() sdk.Dec { return v.Commission }

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces.
func (v Validator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(v.ConsensusPubkey, &pk)
}
