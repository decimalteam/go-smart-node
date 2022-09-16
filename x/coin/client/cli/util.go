package cli

import (
	"context"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// Return coin instance from State
func existCoinDenom(ctx client.Context, denom string) error {
	queryClient := types.NewQueryClient(ctx)
	res, err := queryClient.Coin(context.Background(), types.NewQueryCoinRequest(denom))
	switch {
	case err != nil:
		return err
	case res == nil:
		return errors.CoinDoesNotExist
	default:
		return nil
	}
}

func getCoin(ctx client.Context, denom string) (*types.QueryCoinResponse, error) {
	queryClient := types.NewQueryClient(ctx)

	return queryClient.Coin(context.Background(), types.NewQueryCoinRequest(denom))
}

func getBalanceWithDenom(ctx client.Context, address sdk.AccAddress, denom string) (*bankTypes.QueryBalanceResponse, error) {
	queryClient := bankTypes.NewQueryClient(ctx)

	return queryClient.Balance(context.Background(), bankTypes.NewQueryBalanceRequest(address, denom))
}

func getBalances(ctx client.Context, address sdk.AccAddress, req *query.PageRequest) (*bankTypes.QueryAllBalancesResponse, error) {
	queryClient := bankTypes.NewQueryClient(ctx)

	return queryClient.AllBalances(context.Background(), bankTypes.NewQueryAllBalancesRequest(address, req))
}

func parseCoin(clientCtx client.Context, amount string) (sdk.Coin, error) {
	var (
		coin sdk.Coin
		err  error
	)
	coin, err = sdk.ParseCoinNormalized(amount)
	if err != nil {
		return coin, err
	}

	resp, err := getCoin(clientCtx, coin.Denom)
	switch {
	case err != nil:
		return coin, err
	case resp == nil:
		return coin, errors.CoinDoesNotExist
	default:
		return coin, nil
	}
}

func checkBalance(clientCtx client.Context, address sdk.AccAddress, needAmount sdkmath.Int, denom string) error {
	balance, err := getBalanceWithDenom(clientCtx, address, denom)
	if err != nil {
		return err
	}

	amountBalance := balance.Balance.Amount

	if amountBalance.LT(needAmount) {
		return errors.InsufficientFunds
	}

	return nil
}
