package types

import (
	"fmt"
)

func NewTokenOwner(address string, subTokenIDs []uint64) *TokenOwner {
	return &TokenOwner{
		Address:     address,
		SubTokenIDs: SortedUintArray(subTokenIDs).Sort(),
	}
}

func (t *TokenOwner) GetAddress() string {
	return t.Address
}

func (t *TokenOwner) GetSubTokenIDs() SortedUintArray {
	return t.SubTokenIDs
}

func (t *TokenOwner) SetSubTokenID(subTokenID uint64) *TokenOwner {
	index := t.SubTokenIDs.Find(subTokenID)
	if index == -1 {
		t.SubTokenIDs = append(t.SubTokenIDs, subTokenID).Sort()
	} else {
		t.SubTokenIDs[index] = subTokenID
	}
	return t
}

func (t *TokenOwner) RemoveSubTokenID(subTokenID uint64) *TokenOwner {
	index := t.SubTokenIDs.Find(subTokenID)
	if index != -1 {
		t.SubTokenIDs = append(t.SubTokenIDs[:index], t.SubTokenIDs[index+1:]...)
	}
	return t
}

type TokenOwners []TokenOwner

func (t TokenOwners) SetOwner(owner *TokenOwner) TokenOwners {
	for i, o := range t {
		if o.GetAddress() == owner.GetAddress() {
			t[i] = *owner
			return t
		}
	}

	t = append(t, *owner)

	return t
}

func (t TokenOwners) GetOwner(address string) *TokenOwner {
	for _, owner := range t {
		ownerAddr := owner.GetAddress()
		if ownerAddr == address {
			return &owner
		}
	}
	return nil
}

func (t TokenOwners) String() string {
	if len(t) == 0 {
		return ""
	}

	out := ""
	for _, owner := range t {
		out += fmt.Sprintf("%v\n", owner)
	}
	return out[:len(out)-1]
}
