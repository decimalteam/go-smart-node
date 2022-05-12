package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (m *GenesisState) GetOwners() []Owner {
	if m != nil {
		return m.Owners
	}
	return nil
}

func (m *GenesisState) GetCollections() []Collection {
	if m != nil {
		return m.Collections
	}
	return nil
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(owners []Owner, collections Collections) GenesisState {
	return GenesisState{
		Owners:      owners,
		Collections: collections,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState([]Owner{}, NewCollections())
}

// ValidateGenesis performs basic validation of nfts genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	for _, Owner := range data.Owners {
		addr, err := sdk.AccAddressFromBech32(Owner.Address)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
		}
		if addr.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "address cannot be empty")
		}
	}
	return nil
}
