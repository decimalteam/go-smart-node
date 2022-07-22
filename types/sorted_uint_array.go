package types

import (
	"sort"
	"strconv"
	"strings"
)

type FindableUint64 interface {
	ElAtIndex(index int) uint64
	Len() int
}

func FindUtilUint64(group FindableUint64, el uint64) int {
	if group.Len() == 0 {
		return -1
	}
	low := 0
	high := group.Len() - 1
	median := 0
	for low <= high {
		median = (low + high) / 2
		if group.ElAtIndex(median) == el {
			return median
		} else if group.ElAtIndex(median) < el {
			low = median + 1
		} else if group.ElAtIndex(median) > el {
			high = median - 1
		}
	}
	return -1
}

//-----------------------------------------------------------------------------

type SortedUintArray []uint64

func (sa SortedUintArray) Max() uint64 {
	if sa.Len() == 0 {
		return 0
	}

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
