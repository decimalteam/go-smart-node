package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/exported"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func NewTokenOwner(address string, subTokenIDs []int64) exported.TokenOwner {
	return &TokenOwner{
		Address:     address,
		SubTokenIDs: subTokenIDs,
	}
}

func (t *TokenOwner) GetAddress() string {
	return t.Address
}

func (t *TokenOwner) GetSubTokenIDs() []int64 {
	return t.SubTokenIDs
}

func (t *TokenOwner) SetSubTokenID(subTokenID int64) exported.TokenOwner {
	index := t.SubTokenIDs.Find(subTokenID)
	if index == -1 {
		t.SubTokenIDs = append(t.SubTokenIDs, subTokenID).Sort()
	} else {
		t.SubTokenIDs[index] = subTokenID
	}
	return t
}

func (t *TokenOwner) RemoveSubTokenID(subTokenID int64) exported.TokenOwner {
	index := t.SubTokenIDs.Find(subTokenID)
	if index != -1 {
		t.SubTokenIDs = append(t.SubTokenIDs[:index], t.SubTokenIDs[index+1:]...)
	}
	return t
}

func (t *TokenOwner) SortSubTokensFix() exported.TokenOwner {
	t.SubTokenIDs = t.SubTokenIDs.Sort()
	return t
}

type TokenOwners []TokenOwner

func (t TokenOwners) GetOwners() []exported.TokenOwner {
	var result = make([]exported.TokenOwner, len(t))
	for i, v := range t {
		result[i] = &v
	}

	return result
}

func (t TokenOwners) SetOwner(owner exported.TokenOwner) (exported.TokenOwners, error) {
	for i, o := range t {
		if o.GetAddress() == owner.GetAddress() {
			t[i] = *(owner.(*TokenOwner))
			return t, nil
		}
	}

	t = append(t, TokenOwner{
		Address:     owner.GetAddress(),
		SubTokenIDs: owner.GetSubTokenIDs(),
	})

	return t, nil
}

func (t TokenOwners) GetOwner(address string) exported.TokenOwner {
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

type IDCollections []IDCollection

// String follows stringer interface
func (idCollections IDCollections) String() string {
	if len(idCollections) == 0 {
		return ""
	}

	out := ""
	for _, idCollection := range idCollections {
		out += fmt.Sprintf("%v\n", idCollection.String())
	}
	return out[:len(out)-1]
}

// Append appends IDCollections to IDCollections
func (idCollections IDCollections) Append(idCollections2 ...IDCollection) IDCollections {
	return append(idCollections, idCollections2...).Sort()
}
func (idCollections IDCollections) find(denom string) int {
	return FindUtil(idCollections, denom)
}

// Sort and Findable interface for IDCollections

func (idCollections IDCollections) ElAtIndex(index int) string { return idCollections[index].Denom }
func (idCollections IDCollections) Len() int                   { return len(idCollections) }
func (idCollections IDCollections) Less(i, j int) bool {
	return strings.Compare(idCollections[i].Denom, idCollections[j].Denom) == -1
}
func (idCollections IDCollections) Swap(i, j int) {
	idCollections[i].Denom, idCollections[j].Denom = idCollections[j].Denom, idCollections[i].Denom
}

var _ sort.Interface = IDCollections{}

// Sort is a helper function to sort the set of strings in place
func (idCollections IDCollections) Sort() IDCollections {
	sort.Sort(idCollections)
	return idCollections
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

//-----------------------------------------------------------------------------

type SortedIntArray []int64

// Sort and Findable interface for SortedIntArray

func (sa SortedIntArray) ElAtIndex(index int) int64 { return sa[index] }
func (sa SortedIntArray) Len() int                  { return len(sa) }
func (sa SortedIntArray) Less(i, j int) bool {
	return sa[i] < sa[j]
}
func (sa SortedIntArray) Swap(i, j int) {
	sa[i], sa[j] = sa[j], sa[i]
}

var _ sort.Interface = SortedStringArray{}

// Sort is a helper function to sort the set of strings in place
func (sa SortedIntArray) Sort() SortedIntArray {
	sort.Sort(sa)
	return sa
}

func (sa SortedIntArray) Find(el int64) (idx int) {
	return FindUtilInt64(sa, el)
}

// String is the string representation
func (sa SortedIntArray) String() string {
	str := make([]string, sa.Len())
	for i, v := range sa {
		str[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(str[:], ",")
}
