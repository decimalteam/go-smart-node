package types

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// NewOwnerCollection creates a new IDCollection instance
func NewOwnerCollection(denom string, ids []string) OwnerCollection {
	return OwnerCollection{
		Denom: strings.TrimSpace(denom),
		NFTs:  SortedStringArray(ids).Sort(),
	}
}

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

func (t TokenOwners) GetOwners() []TokenOwner {
	return t
}

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

func (m OwnerCollection) AddID(id string) OwnerCollection {
	if len(m.NFTs) == 0 {
		m.NFTs = []string{id}
		return m
	}

	m.NFTs = append(m.NFTs, id).Sort()
	return m
}

//-----------------------------------------------------------------------------

// SortedStringArray is an array of strings whose sole purpose is to help with find
type SortedStringArray []string

// String is the string representation
func (sa SortedStringArray) String() string { return strings.Join(sa[:], ",") }

// Sort and Findable interface for SortedStringArray

func (sa SortedStringArray) ElAtIndex(index int) string { return sa[index] }
func (sa SortedStringArray) Len() int                   { return len(sa) }
func (sa SortedStringArray) Less(i, j int) bool {
	return strings.Compare(sa[i], sa[j]) == -1
}
func (sa SortedStringArray) Swap(i, j int) {
	sa[i], sa[j] = sa[j], sa[i]
}

var _ sort.Interface = SortedStringArray{}

// Sort is a helper function to sort the set of strings in place
func (sa SortedStringArray) Sort() SortedStringArray {
	sort.Sort(sa)
	return sa
}

func (sa SortedStringArray) Has(item string) bool {
	for _, saItem := range sa {
		if saItem == item {
			return true
		}
	}

	return false
}

//-----------------------------------------------------------------------------

type SortedUintArray []uint64

func (sa SortedUintArray) Max() uint64 {
	return sa[sa.Len()-1]
}

// Sort and Findable interface for SortedIntArray

func (sa SortedUintArray) ElAtIndex(index int) uint64 { return sa[index] }
func (sa SortedUintArray) Len() int                   { return len(sa) }
func (sa SortedUintArray) Less(i, j int) bool {
	return sa[i] < sa[j]
}
func (sa SortedUintArray) Swap(i, j int) {
	sa[i], sa[j] = sa[j], sa[i]
}

var _ sort.Interface = SortedStringArray{}

// Sort is a helper function to sort the set of strings in place
func (sa SortedUintArray) Sort() SortedUintArray {
	sort.Sort(sa)
	return sa
}

func (sa SortedUintArray) Find(el uint64) (idx int) {
	return FindUtilUint64(sa, el)
}

// String is the string representation
func (sa SortedUintArray) String() string {
	str := make([]string, sa.Len())
	for i, v := range sa {
		str[i] = strconv.FormatUint(v, 10)
	}
	return strings.Join(str[:], ",")
}
