package types

import "strings"

// NewOwnerCollection creates a new IDCollection instance
func NewOwnerCollection(denom string, ids []string) OwnerCollection {
	return OwnerCollection{
		Denom: strings.TrimSpace(denom),
		NFTs:  SortedStringArray(ids).Sort(),
	}
}

func (m OwnerCollection) AddID(id string) OwnerCollection {
	if len(m.NFTs) == 0 {
		m.NFTs = []string{id}
		return m
	}

	m.NFTs = append(m.NFTs, id).Sort()
	return m
}
