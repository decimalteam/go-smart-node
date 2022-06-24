package types

import (
	"strings"
)

type Coins []Coin

func (c Coins) String() string {
	result := make([]string, len(c))
	for i, v := range c {
		result[i] = v.String()
	}

	return strings.Join(result, "\n")
}
