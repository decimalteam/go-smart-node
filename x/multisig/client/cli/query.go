package cli

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
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
		cmdQueryWallet(),
		cmdQueryWallets(),
		cmdQueryTransaction(),
		cmdQueryTransactions(),
	)
	return cmd
}

func cmdQueryWallet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallet [address]",
		Short: "Query multisig wallet by address",
		Long: fmt.Sprintf(`Query full information about multisig wallet

Example: 	
$ %s query multisig wallet dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f`, config.AppBinName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryWalletRequest{
				Wallet: strings.ToLower(args[0]),
			}

			res, err := queryClient.Wallet(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryWallets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallets [owner]",
		Short: "Query all multisig wallets by owner",
		Long: fmt.Sprintf(`Query full information about multisig wallets with owner

Example: 	
$ %s query multisig wallets dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f`, config.AppBinName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryWalletsRequest{
				Owner:      strings.ToLower(args[0]),
				Pagination: pageReq,
			}

			res, err := queryClient.Wallets(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all wallets")

	return cmd
}

func cmdQueryTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transaction [tx_id]",
		Short: "Query multisig transaction by id",
		Long: fmt.Sprintf(`Query information about multisig transaction

Example: 	
$ %s query multisig transaction dx18...`, config.AppBinName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryTransactionRequest{
				Id: strings.ToLower(args[0]),
			}

			res, err := queryClient.Transaction(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryTransactions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transactions [wallet]",
		Short: "Query all multisig transactions by wallet address",
		Long: fmt.Sprintf(`Query information about multisig transactions by wallet address

Example: 	
$ %s query multisig transactions dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f`, config.AppBinName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryTransactionsRequest{
				Wallet:     strings.ToLower(args[0]),
				Pagination: pageReq,
			}

			res, err := queryClient.Transactions(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all transactions")

	return cmd
}
