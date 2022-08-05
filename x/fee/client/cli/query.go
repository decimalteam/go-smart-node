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
		cmdQueryBaseDenomPrice(),
	)

	return queryCmd
}

// GetCmdQueryCollectionSupply queries the supply of a nft collection
func cmdQueryBaseDenomPrice() *cobra.Command {
	return &cobra.Command{
		Use:   "price",
		Short: "base denom price from oracle",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get current denom price from oracle.

Example:
$ %s query %s price
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

			res, err := queryClient.QueryBaseDenomPrice(context.Background(), &types.QueryBaseDenomPriceRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}
