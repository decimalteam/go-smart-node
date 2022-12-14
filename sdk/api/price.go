package api

import (
	"context"

	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const QuotePair = "usd"

// GetFeeParams returns all for ccorrect transaction fee calculation
func (api *API) GetFeeParams(baseDenom, quoteDenom string) (sdk.Dec, feetypes.Params, error) {
	client := feetypes.NewQueryClient(api.grpcClient)
	// 1. price
	resp, err := client.CoinPrice(
		context.Background(),
		&feetypes.QueryCoinPriceRequest{
			Denom: baseDenom,
			Quote: quoteDenom,
		},
	)
	if err != nil {
		return sdk.ZeroDec(), feetypes.DefaultParams(), err
	}
	// 2. params
	respP, err := client.ModuleParams(
		context.Background(),
		&feetypes.QueryModuleParamsRequest{},
	)
	if err != nil {
		return sdk.ZeroDec(), feetypes.DefaultParams(), err
	}
	return resp.Price.Price, respP.Params, nil
}
