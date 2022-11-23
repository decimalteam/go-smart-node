package api

import (
	"context"

	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
)

type Coin = cointypes.Coin

func (api *API) Coins() ([]Coin, error) {
	client := cointypes.NewQueryClient(api.grpcClient)
	coins := make([]cointypes.Coin, 0)
	req := &cointypes.QueryCoinsRequest{
		Pagination: &query.PageRequest{Limit: queryLimit},
	}
	for {
		res, err := client.Coins(
			context.Background(),
			req,
		)
		if err != nil {
			return []Coin{}, err
		}
		if len(res.Coins) == 0 {
			break
		}
		coins = append(coins, res.Coins...)
		if len(res.Pagination.NextKey) == 0 {
			break
		}
		req.Pagination.Key = res.Pagination.NextKey
	}
	return coins, nil
}

func (api *API) Coin(denom string) (Coin, error) {
	client := cointypes.NewQueryClient(api.grpcClient)
	req := &cointypes.QueryCoinRequest{
		Denom: denom,
	}
	res, err := client.Coin(
		context.Background(),
		req,
	)
	if err != nil {
		return Coin{}, err
	}
	return res.Coin, nil
}
