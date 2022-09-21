package types

import (
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	feeconfig "bitbucket.org/decimalteam/go-smart-node/x/fee/config"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, initialPrices []CoinPrice) GenesisState {
	return GenesisState{
		Params: params,
		Prices: initialPrices,
	}
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		Prices: []CoinPrice{
			{
				Denom: config.BaseDenom,
				Quote: feeconfig.DefaultQuote,
				Price: sdk.OneDec(),
			},
		},
	}
}

func (gs *GenesisState) Validate() error {
	type pair struct {
		denom string
		quote string
	}
	if len(gs.Prices) == 0 {
		return errors.WrongPrice
	}
	knownPairs := make(map[pair]bool)
	for _, price := range gs.Prices {
		if !price.Price.IsPositive() {
			return errors.WrongPrice
		}
		key := pair{
			denom: price.Denom,
			quote: price.Quote,
		}
		if knownPairs[key] {
			return errors.DuplicateCoinPrice
		}
		knownPairs[key] = true
	}

	return gs.Params.Validate()
}
