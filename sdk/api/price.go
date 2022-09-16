package api

import (
	"context"

	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetFeeParams returns all for ccorrect transaction fee calculation
func (api *API) GetFeeParams() (sdk.Dec, feetypes.Params, error) {
	client := feetypes.NewQueryClient(api.grpcClient)
	// 1. price
	resp, err := client.QueryBaseDenomPrice(
		context.Background(),
		&feetypes.QueryBaseDenomPriceRequest{},
	)
	if err != nil {
		return sdk.ZeroDec(), feetypes.DefaultParams(), err
	}
	price, err := sdk.NewDecFromStr(resp.Price)
	if err != nil {
		return sdk.ZeroDec(), feetypes.DefaultParams(), err
	}
	// 2. params
	respP, err := client.QueryModuleParams(
		context.Background(),
		&feetypes.QueryParamsRequest{},
	)
	if err != nil {
		return sdk.ZeroDec(), feetypes.DefaultParams(), err
	}
	return price, respP.Params, nil
}
