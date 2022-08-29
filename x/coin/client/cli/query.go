package cli

import (
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/btcutil/base58"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
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
		QueryCoinCmd(),
		QueryCoinsCmd(),
		QueryCheckCmd(),
		QueryChecksCmd(),
		QueryParamsCmd(),
	)
	return cmd
}

func QueryCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coin [symbol]",
		Short: "Query specific coin by symbol (denom)",
		Long: fmt.Sprintf(`Query coin full information 

Example: 	
$ %s query %s coin del`, config.AppBinName, types.ModuleName),
		Args: cobra.ExactArgs(1),
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

func QueryCoinsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coins",
		Short: "Query all existing coins",
		Long: fmt.Sprintf(`Query all coins full information 

Example: 	
$ %s query %s coins`, config.AppBinName, types.ModuleName),
		Args: cobra.NoArgs,
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

func QueryCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check [hash]",
		Short: "Query specific check by hash in hex format",
		Long: fmt.Sprintf(`Query check information from blockchain

Example: 	
$ %s query %s check 3YEtqixL7ccFTZJaMUHx3TgsQEqzrqoj...(result of command 'issue-check')`, config.AppBinName, types.ModuleName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// Decode provided check from base58 format to raw bytes
			checkBytes := base58.Decode(args[0])
			if len(checkBytes) == 0 {
				return errors.UnableDecodeCheckBase58
			}

			// Parse provided check from raw bytes to ensure it is valid
			check, err := types.ParseCheck(checkBytes)
			if err != nil {
				return err
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

func QueryChecksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checks",
		Short: "Query all existing checks",
		Long: fmt.Sprintf(`Query all checks information from blockchain

Example: 	
$ %s query %s checks`, config.AppBinName, types.ModuleName),
		Args: cobra.NoArgs,
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

func QueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: fmt.Sprintf("Query the current parameters of the module %s", types.ModuleName),
		Long: fmt.Sprintf(`Query module params from blockchain

Example: 	
$ %s query %s params`, config.AppBinName, types.ModuleName),
		Args: cobra.NoArgs,
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
