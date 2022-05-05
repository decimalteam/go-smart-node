package types

import "strings"

// Findable is an interface for iterable types that allows the FindUtil function to work
type Findable interface {
	ElAtIndex(index int) string
	Len() int
}

// FindUtil is a binary search funcion for types that support the Findable interface (elements must be sorted)
func FindUtil(group Findable, el string) int {
	if group.Len() == 0 {
		return -1
	}
	low := 0
	high := group.Len() - 1
	median := 0
	for low <= high {
		median = (low + high) / 2
		switch compare := strings.Compare(group.ElAtIndex(median), el); {
		case compare == 0:
			// if group[median].element == el
			return median
		case compare == -1:
			// if group[median].element < el
			low = median + 1
		default:
			// if group[median].element > el
			high = median - 1
		}
	}
	return -1
}

type FindableInt64 interface {
	ElAtIndex(index int) int64
	Len() int
}

func FindUtilInt64(group FindableInt64, el int64) int {
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
