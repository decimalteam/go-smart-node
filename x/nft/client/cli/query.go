package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

// GetQueryCmd returns the parent command for all the module's CLI query commands.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		RunE:  client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdQueryCollections(),
		cmdQueryCollection(),
		cmdQueryToken(),
		cmdQuerySubToken(),
	)

	return cmd
}

func cmdQueryCollections() *cobra.Command {
	return &cobra.Command{
		Use:   "collections [creator]",
		Short: "get the complete list of existing NFT collections",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the complete list of existing NFT collections.

Example:
$ %s query %s collections
$ %s query %s collections dx1c6r9smnxccxlnf7rmxze9pajs0l2d3sftdvr32
`, cmdcfg.AppBinName, types.ModuleName, cmdcfg.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if len(args) == 0 {
				req := &types.QueryCollectionsRequest{}
				res, err := queryClient.Collections(context.Background(), req)
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res)
			} else {
				req := &types.QueryCollectionsByCreatorRequest{}
				res, err := queryClient.CollectionsByCreator(context.Background(), req)
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res)
			}
		},
	}
}

func cmdQueryCollection() *cobra.Command {
	return &cobra.Command{
		Use:   "collection [creator] [denom]",
		Short: "get existing NFT collections with specified creator address and collection denom",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get existing NFT collections with specified creator address and collection denom.

Example:
$ %s query %s collection dx1c6r9smnxccxlnf7rmxze9pajs0l2d3sftdvr32 crypto-kitties
`, cmdcfg.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryCollectionRequest{
				Creator: args[0],
				Denom:   args[1],
			}
			res, err := queryClient.Collection(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
}

func cmdQueryToken() *cobra.Command {
	return &cobra.Command{
		Use:   "token [token_id]",
		Short: "get existing NFT token with specified ID",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get existing NFT token with specified ID.

Example:
$ %s query %s token d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa
`, cmdcfg.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryTokenRequest{
				TokenId: args[0],
			}
			res, err := queryClient.Token(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
}

func cmdQuerySubToken() *cobra.Command {
	return &cobra.Command{
		Use:   "sub-token [token_id] [sub_token_id]",
		Short: "get existing NFT sub-token with specified token ID and sub-token ID",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get existing NFT sub-token with specified token ID and sub-token ID.

Example:
$ %s query %s sub-token d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa 1
`, cmdcfg.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QuerySubTokenRequest{
				TokenId:    args[0],
				SubTokenId: args[1],
			}
			res, err := queryClient.SubToken(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
}
