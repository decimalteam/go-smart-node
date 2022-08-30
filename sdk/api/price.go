package api

import (
	"context"

	feeTypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetFeeParams returns all for ccorrect transaction fee calculation
func (api *API) GetFeeParams() (sdk.Dec, feeTypes.Params, error) {
	client := feeTypes.NewQueryClient(api.grpcClient)
	// 1. price
	resp, err := client.QueryBaseDenomPrice(
		context.Background(),
		&feeTypes.QueryBaseDenomPriceRequest{},
	)
	if err != nil {
		return sdk.ZeroDec(), feeTypes.DefaultParams(), err
	}
	price, err := sdk.NewDecFromStr(resp.Price)
	if err != nil {
		return sdk.ZeroDec(), feeTypes.DefaultParams(), err
	}
	// 2. params
	respP, err := client.QueryParams(
		context.Background(),
		&feeTypes.QueryParamsRequest{},
	)
	if err != nil {
		return sdk.ZeroDec(), feeTypes.DefaultParams(), err
	}
	return price, respP.Params, nil
}
