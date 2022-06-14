package cli

import (
	"context"
	"fmt"
	"github.com/cosmos/btcutil/base58"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// GetQueryCmd returns the parent command for all the module's CLI query commands.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		cmdQueryCoin(),
		cmdQueryCoins(),
		cmdQueryCheck(),
		cmdQueryChecks(),
		cmdQueryParams(),
	)
	return cmd
}

func cmdQueryCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coin [symbol]",
		Short: "Query specific coin by symbol (denom)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryCoinRequest{
				Symbol: strings.ToLower(args[0]),
			}

			res, err := queryClient.Coin(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryCoins() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coins",
		Short: "Query all existing coins",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryCoinsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Coins(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all coins")

	return cmd
}

func cmdQueryCheck() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check [hash]",
		Short: "Query specific check by hash in hex format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// Decode provided check from base58 format to raw bytes
			checkBytes := base58.Decode(args[0])
			if len(checkBytes) == 0 {
				return types.ErrUnableDecodeCheck(args[0])
			}

			// Parse provided check from raw bytes to ensure it is valid
			check, err := types.ParseCheck(checkBytes)
			if err != nil {
				return types.ErrInvalidCheck(err.Error())
			}

			hash := check.HashFull()
			req := &types.QueryCheckRequest{
				Hash: hash[:],
			}

			res, err := queryClient.Check(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryChecks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checks",
		Short: "Query all existing checks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			req := &types.QueryChecksRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Checks(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all checks")

	return cmd
}

func cmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: fmt.Sprintf("Query the current parameters of the module %s", types.ModuleName),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryParamsRequest{}

			res, err := queryClient.Params(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
