package types

import (
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ DelegationI = Delegation{}
var _ RedelegationI = Redelegation{}
var _ UndelegationI = Undelegation{}

////////////////////////////////////////////////////////////////
// Delegation
////////////////////////////////////////////////////////////////

// NewDelegation creates a new delegation object.
func NewDelegation(delegator sdk.AccAddress, validator sdk.ValAddress, stake Stake) Delegation {
	return Delegation{
		Delegator: delegator.String(),
		Validator: validator.String(),
		Stake:     stake,
	}
}

// GetDelegator returns delegator address.
func (d Delegation) GetDelegator() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(d.Delegator)
}

// GetValidator returns validator address.
func (d Delegation) GetValidator() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(d.Validator)
	if err != nil {
		panic(err)
	}
	return addr
}

// GetStake returns the stake delegated.
func (d Delegation) GetStake() StakeI { return d.Stake }

////////////////////////////////////////////////////////////////

// MustMarshalDelegation returns the delegation bytes. Panics if fails.
func MustMarshalDelegation(cdc codec.BinaryCodec, delegation Delegation) []byte {
	return cdc.MustMarshal(&delegation)
}

// MustUnmarshalDelegation returns the unmarshaled delegation from bytes. Panics if fails.
func MustUnmarshalDelegation(cdc codec.BinaryCodec, value []byte) Delegation {
	delegation, err := UnmarshalDelegation(cdc, value)
	if err != nil {
		panic(err)
	}
	return delegation
}

// UnmarshalDelegation returns the unmarshaled delegation from bytes.
func UnmarshalDelegation(cdc codec.BinaryCodec, value []byte) (delegation Delegation, err error) {
	err = cdc.Unmarshal(value, &delegation)
	return delegation, err
}

////////////////////////////////////////////////////////////////

// Delegations is a collection of delegations.
type Delegations []Delegation

func (d Delegations) String() (out string) {
	for _, del := range d {
		out += del.String() + "\n"
	}
	return strings.TrimSpace(out)
}

////////////////////////////////////////////////////////////////
// Redelegation
////////////////////////////////////////////////////////////////

// NewRedelegation creates a new redelegation object.
func NewRedelegation(
	delegator sdk.AccAddress,
	validatorSrc sdk.ValAddress,
	validatorDst sdk.ValAddress,
	creationHeight int64,
	minTime time.Time,
	stake Stake,
) Redelegation {
	return Redelegation{
		Delegator:    delegator.String(),
		ValidatorSrc: validatorSrc.String(),
		ValidatorDst: validatorDst.String(),
		Entries: []RedelegationEntry{
			NewRedelegationEntry(creationHeight, minTime, stake),
		},
	}
}

// GetDelegator returns delegator address.
func (red Redelegation) GetDelegator() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(red.Delegator)
}

// GetValidatorSrc returns source validator address.
func (red Redelegation) GetValidatorSrc() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(red.ValidatorSrc)
	if err != nil {
		panic(err)
	}
	return addr
}

// GetValidatorDst returns destination validator address.
func (red Redelegation) GetValidatorDst() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(red.ValidatorDst)
	if err != nil {
		panic(err)
	}
	return addr
}

// GetEntries returns redelegation entries.
func (red Redelegation) GetEntries() []RedelegationEntry { return red.Entries }

// AddEntry appends new entry to the redelegation.
func (red *Redelegation) AddEntry(creationHeight int64, minTime time.Time, stake Stake) {
	entry := NewRedelegationEntry(creationHeight, minTime, stake)
	red.Entries = append(red.Entries, entry)
}

// RemoveEntry removes existing entry at index i from the redelegation.
func (red *Redelegation) RemoveEntry(i int64) {
	red.Entries = append(red.Entries[:i], red.Entries[i+1:]...)
}

////////////////////////////////////////////////////////////////

// MustMarshalRED returns the redelegation bytes. Panics if fails.
func MustMarshalRED(cdc codec.BinaryCodec, red Redelegation) []byte {
	return cdc.MustMarshal(&red)
}

// MustUnmarshalRED returns the unmarshaled redelegation from bytes. Panics if fails.
func MustUnmarshalRED(cdc codec.BinaryCodec, value []byte) Redelegation {
	red, err := UnmarshalRED(cdc, value)
	if err != nil {
		panic(err)
	}
	return red
}

// UnmarshalRED returns the unmarshaled redelegation from bytes.
func UnmarshalRED(cdc codec.BinaryCodec, value []byte) (red Redelegation, err error) {
	err = cdc.Unmarshal(value, &red)
	return red, err
}

////////////////////////////////////////////////////////////////

// NewRedelegationEntry creates a new redelegation entry object.
func NewRedelegationEntry(creationHeight int64, completionTime time.Time, stake Stake) RedelegationEntry {
	return RedelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		Stake:          stake,
	}
}

// IsMature returns true if the entry is mature currently (redelegation is ready to be completed).
func (e RedelegationEntry) IsMature(now time.Time) bool {
	return !e.CompletionTime.After(now)
}

////////////////////////////////////////////////////////////////

// Redelegations is a collection of redelegations.
type Redelegations []Redelegation

func (d Redelegations) String() (out string) {
	for _, red := range d {
		out += red.String() + "\n"
	}
	return strings.TrimSpace(out)
}

////////////////////////////////////////////////////////////////
// Undelegation
////////////////////////////////////////////////////////////////

// NewUndelegation creates a new undelegation object.
func NewUndelegation(
	delegator sdk.AccAddress,
	validator sdk.ValAddress,
	creationHeight int64,
	minTime time.Time,
	stake Stake,
) Undelegation {
	return Undelegation{
		Delegator: delegator.String(),
		Validator: validator.String(),
		Entries: []UndelegationEntry{
			NewUndelegationEntry(creationHeight, minTime, stake),
		},
	}
}

// GetDelegator returns delegator address.
func (ubd Undelegation) GetDelegator() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(ubd.Delegator)
}

// GetValidator returns validator address.
func (ubd Undelegation) GetValidator() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(ubd.Validator)
	if err != nil {
		panic(err)
	}
	return addr
}

// GetEntries returns undelegation entries.
func (ubd Undelegation) GetEntries() []UndelegationEntry { return ubd.Entries }

// AddEntry appends new entry to the undelegation.
func (ubd *Undelegation) AddEntry(creationHeight int64, minTime time.Time, stake Stake) {
	entry := NewUndelegationEntry(creationHeight, minTime, stake)
	ubd.Entries = append(ubd.Entries, entry)
}

// RemoveEntry removes existing entry at index i from the undelegation.
func (ubd *Undelegation) RemoveEntry(i int64) {
	ubd.Entries = append(ubd.Entries[:i], ubd.Entries[i+1:]...)
}

////////////////////////////////////////////////////////////////

// MustMarshalUBD returns the undelegation bytes. Panics if fails.
func MustMarshalUBD(cdc codec.BinaryCodec, ubd Undelegation) []byte {
	return cdc.MustMarshal(&ubd)
}

// MustUnmarshalUBD returns the unmarshaled undelegation from bytes. Panics if fails.
func MustUnmarshalUBD(cdc codec.BinaryCodec, value []byte) Undelegation {
	ubd, err := UnmarshalUBD(cdc, value)
	if err != nil {
		panic(err)
	}
	return ubd
}

// UnmarshalUBD returns the unmarshaled undelegation from bytes.
func UnmarshalUBD(cdc codec.BinaryCodec, value []byte) (ubd Undelegation, err error) {
	err = cdc.Unmarshal(value, &ubd)
	return ubd, err
}

////////////////////////////////////////////////////////////////

// NewUndelegationEntry creates a new undelegation entry object.
func NewUndelegationEntry(creationHeight int64, completionTime time.Time, stake Stake) UndelegationEntry {
	return UndelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		Stake:          stake,
	}
}

// IsMature returns true if the entry is mature currently (undelegation is ready to be completed).
func (e UndelegationEntry) IsMature(now time.Time) bool {
	return !e.CompletionTime.After(now)
}

////////////////////////////////////////////////////////////////

// Undelegations is a collection of undelegations.
type Undelegations []Undelegation

func (ubds Undelegations) String() (out string) {
	for _, u := range ubds {
		out += u.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// // ----------------------------------------------------------------------------
// // Client Types

// // NewDelegationResp creates a new DelegationResponse instance
// func NewDelegationResp(
// 	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, shares sdk.Dec, balance sdk.Coin,
// ) DelegationResponse {
// 	return DelegationResponse{
// 		Delegation: NewDelegation(delegatorAddr, validatorAddr, shares),
// 		Balance:    balance,
// 	}
// }

// // String implements the Stringer interface for DelegationResponse.
// func (d DelegationResponse) String() string {
// 	return fmt.Sprintf("%s\n  Balance:   %s", d.Delegation.String(), d.Balance)
// }

// type delegationRespAlias DelegationResponse

// // MarshalJSON implements the json.Marshaler interface. This is so we can
// // achieve a flattened structure while embedding other types.
// func (d DelegationResponse) MarshalJSON() ([]byte, error) {
// 	return json.Marshal((delegationRespAlias)(d))
// }

// // UnmarshalJSON implements the json.Unmarshaler interface. This is so we can
// // achieve a flattened structure while embedding other types.
// func (d *DelegationResponse) UnmarshalJSON(bz []byte) error {
// 	return json.Unmarshal(bz, (*delegationRespAlias)(d))
// }

// // DelegationResponses is a collection of DelegationResp
// type DelegationResponses []DelegationResponse

// // String implements the Stringer interface for DelegationResponses.
// func (d DelegationResponses) String() (out string) {
// 	for _, del := range d {
// 		out += del.String() + "\n"
// 	}

// 	return strings.TrimSpace(out)
// }

// // NewRedelegationResponse crates a new RedelegationEntryResponse instance.
// //
// //nolint:interfacer
// func NewRedelegationResponse(
// 	delegatorAddr sdk.AccAddress, validatorSrc, validatorDst sdk.ValAddress, entries []RedelegationEntryResponse,
// ) RedelegationResponse {
// 	return RedelegationResponse{
// 		Redelegation: Redelegation{
// 			Delegator:    delegatorAddr.String(),
// 			ValidatorSrc: validatorSrc.String(),
// 			ValidatorDst: validatorDst.String(),
// 		},
// 		Entries: entries,
// 	}
// }

// // NewRedelegationEntryResponse creates a new RedelegationEntryResponse instance.
// func NewRedelegationEntryResponse(
// 	creationHeight int64, completionTime time.Time, sharesDst sdk.Dec, initialBalance, balance sdk.Int,
// ) RedelegationEntryResponse {
// 	return RedelegationEntryResponse{
// 		RedelegationEntry: NewRedelegationEntry(creationHeight, completionTime, initialBalance, sharesDst),
// 		Balance:           balance,
// 	}
// }

// type redelegationRespAlias RedelegationResponse

// // MarshalJSON implements the json.Marshaler interface. This is so we can
// // achieve a flattened structure while embedding other types.
// func (r RedelegationResponse) MarshalJSON() ([]byte, error) {
// 	return json.Marshal((redelegationRespAlias)(r))
// }

// // UnmarshalJSON implements the json.Unmarshaler interface. This is so we can
// // achieve a flattened structure while embedding other types.
// func (r *RedelegationResponse) UnmarshalJSON(bz []byte) error {
// 	return json.Unmarshal(bz, (*redelegationRespAlias)(r))
// }

// // RedelegationResponses are a collection of RedelegationResp
// type RedelegationResponses []RedelegationResponse

// func (r RedelegationResponses) String() (out string) {
// 	for _, red := range r {
// 		out += red.String() + "\n"
// 	}

// 	return strings.TrimSpace(out)
// }
