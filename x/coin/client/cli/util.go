package cli

import (
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// Return coin instance from State
func existCoinSymbol(ctx client.Context, symbol string) error {
	queryClient := types.NewQueryClient(ctx)
	res, err := queryClient.Coin(context.Background(), types.NewQueryCoinRequest(symbol))
	switch {
	case err != nil:
		return err
	case res == nil:
		return types.ErrCoinDoesNotExist(symbol)
	default:
		return nil
	}
}

func getCoin(ctx client.Context, symbol string) (*types.QueryCoinResponse, error) {
	queryClient := types.NewQueryClient(ctx)

	return queryClient.Coin(context.Background(), types.NewQueryCoinRequest(symbol))
}

func getBalanceWithSymbol(ctx client.Context, address sdk.AccAddress, symbol string) (*bankTypes.QueryBalanceResponse, error) {
	queryClient := bankTypes.NewQueryClient(ctx)

	return queryClient.Balance(context.Background(), bankTypes.NewQueryBalanceRequest(address, symbol))
}

func getBalances(ctx client.Context, address sdk.AccAddress, req *query.PageRequest) (*bankTypes.QueryAllBalancesResponse, error) {
	queryClient := bankTypes.NewQueryClient(ctx)

	return queryClient.AllBalances(context.Background(), bankTypes.NewQueryAllBalancesRequest(address, req))
}
