package cli

import (
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	nftQueryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the Fee module",
		RunE:  client.ValidateCmd,
	}

	nftQueryCmd.AddCommand(
		cmdQueryBaseDenomPrice(),
	)

	return nftQueryCmd
}

// GetCmdQueryCollectionSupply queries the supply of a nft collection
func cmdQueryBaseDenomPrice() *cobra.Command {
	return &cobra.Command{
		Use:   "price",
		Short: "base denom price from oracle",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the total count of NFTs that match a certain denomination.

Example:
$ %s query %s supply crypto-kitties
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

			req := types.QueryBaseDenomPriceRequest{}

			res, err := queryClient.QueryBaseDenomPrice(context.Background(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}
