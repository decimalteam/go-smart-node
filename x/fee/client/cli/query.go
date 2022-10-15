package cli

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the Fee module",
		RunE:  client.ValidateCmd,
	}

	queryCmd.AddCommand(
		cmdQueryCoinPrice(),
		cmdQueryCoinPrices(),
	)

	return queryCmd
}

// cmdQueryCoinPrice queries price for piar denom/quote
func cmdQueryCoinPrice() *cobra.Command {
	return &cobra.Command{
		Use:   "price [denom] [quote]",
		Short: "current denom/quote price",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get current denom/quote price.

Example:
$ %s query %s price del usd
`, config.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CoinPrice(context.Background(), &types.QueryCoinPriceRequest{
				Denom: args[0],
				Quote: args[1],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// cmdQueryCoinPrice queries price for piar denom/quote
func cmdQueryCoinPrices() *cobra.Command {
	return &cobra.Command{
		Use:   "prices",
		Short: "all prices pairs",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get prices for all pairs denom/quote.

Example:
$ %s query %s prices
`, config.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CoinPrices(context.Background(), &types.QueryCoinPricesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}
