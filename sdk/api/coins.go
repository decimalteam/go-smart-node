package api

import (
	"context"

	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
)

type Coin = coinTypes.Coin

func (api *API) Coins() ([]Coin, error) {
	client := coinTypes.NewQueryClient(api.grpcClient)
	coins := make([]coinTypes.Coin, 0)
	req := &coinTypes.QueryCoinsRequest{
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
