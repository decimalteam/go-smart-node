package types

import (
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ DelegationI = Delegation{}

// NewDelegation creates a new delegation object.
func NewDelegation(delegator sdk.AccAddress, validator sdk.ValAddress, stake Stake) Delegation {
	return Delegation{
		Delegator: delegator.String(),
		Validator: validator.String(),
		Stake:     stake,
	}
}

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

// return the delegation
func UnmarshalDelegation(cdc codec.BinaryCodec, value []byte) (delegation Delegation, err error) {
	err = cdc.Unmarshal(value, &delegation)
	return delegation, err
}

func (d Delegation) GetDelegator() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(d.Delegator)
}

func (d Delegation) GetValidator() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(d.Validator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (d Delegation) GetStake() StakeI { return d.Stake }

// Delegations is a collection of delegations.
type Delegations []Delegation

func (d Delegations) String() (out string) {
	for _, del := range d {
		out += del.String() + "\n"
	}

	return strings.TrimSpace(out)
}

func NewRedelegationEntry(creationHeight int64, completionTime time.Time, stake Stake) RedelegationEntry {
	return RedelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		Stake:          stake,
	}
}

// IsMature - is the current entry mature
func (e RedelegationEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

//nolint:interfacer
func NewRedelegation(
	delegatorAddr sdk.AccAddress, validatorSrcAddr, validatorDstAddr sdk.ValAddress,
	creationHeight int64, minTime time.Time, stake Stake,
) Redelegation {
	return Redelegation{
		Delegator:    delegatorAddr.String(),
		ValidatorSrc: validatorSrcAddr.String(),
		ValidatorDst: validatorDstAddr.String(),
		Entries: []RedelegationEntry{
			NewRedelegationEntry(creationHeight, minTime, stake),
		},
	}
}

// AddEntry - append entry to the unbonding delegation
func (red *Redelegation) AddEntry(creationHeight int64, minTime time.Time, stake Stake) {
	entry := NewRedelegationEntry(creationHeight, minTime, stake)
	red.Entries = append(red.Entries, entry)
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (red *Redelegation) RemoveEntry(i int64) {
	red.Entries = append(red.Entries[:i], red.Entries[i+1:]...)
}

// MustMarshalRED returns the Redelegation bytes. Panics if fails.
func MustMarshalRED(cdc codec.BinaryCodec, red Redelegation) []byte {
	return cdc.MustMarshal(&red)
}

// MustUnmarshalRED unmarshals a redelegation from a store value. Panics if fails.
func MustUnmarshalRED(cdc codec.BinaryCodec, value []byte) Redelegation {
	red, err := UnmarshalRED(cdc, value)
	if err != nil {
		panic(err)
	}

	return red
}

// UnmarshalRED unmarshals a redelegation from a store value
func UnmarshalRED(cdc codec.BinaryCodec, value []byte) (red Redelegation, err error) {
	err = cdc.Unmarshal(value, &red)
	return red, err
}

// Redelegations are a collection of Redelegation
type Redelegations []Redelegation

func (d Redelegations) String() (out string) {
	for _, red := range d {
		out += red.String() + "\n"
	}

	return strings.TrimSpace(out)
}

func NewUndelegationEntry(creationHeight int64, completionTime time.Time, stake Stake) UndelegationEntry {
	return UndelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		Stake:          stake,
	}
}

// IsMature - is the current entry mature.
func (e UndelegationEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// NewUndelegation - create a new unbonding delegation object.
func NewUndelegation(
	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	creationHeight int64, minTime time.Time, stake Stake,
) Undelegation {
	return Undelegation{
		Delegator: delegatorAddr.String(),
		Validator: validatorAddr.String(),
		Entries: []UndelegationEntry{
			NewUndelegationEntry(creationHeight, minTime, stake),
		},
	}
}

// AddEntry - append entry to the unbonding delegation
func (ubd *Undelegation) AddEntry(creationHeight int64, minTime time.Time, stake Stake) {
	entry := NewUndelegationEntry(creationHeight, minTime, stake)
	ubd.Entries = append(ubd.Entries, entry)
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (ubd *Undelegation) RemoveEntry(i int64) {
	ubd.Entries = append(ubd.Entries[:i], ubd.Entries[i+1:]...)
}

// return the unbonding delegation
func MustMarshalUBD(cdc codec.BinaryCodec, ubd Undelegation) []byte {
	return cdc.MustMarshal(&ubd)
}

// unmarshal a unbonding delegation from a store value
func MustUnmarshalUBD(cdc codec.BinaryCodec, value []byte) Undelegation {
	ubd, err := UnmarshalUBD(cdc, value)
	if err != nil {
		panic(err)
	}

	return ubd
}

// unmarshal a unbonding delegation from a store value
func UnmarshalUBD(cdc codec.BinaryCodec, value []byte) (ubd Undelegation, err error) {
	err = cdc.Unmarshal(value, &ubd)
	return ubd, err
}

// Undelegations is a collection of Undelegation
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
