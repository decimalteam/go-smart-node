package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	nftQueryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the NFT module",
		RunE:  client.ValidateCmd,
	}

	nftQueryCmd.AddCommand(
		cmdQueryCollectionSupply(),
		cmdQueryOwner(),
		cmdQueryCollections(),
		cmdQueryDenoms(),
		cmdQueryNFT(),
		cmdQuerySubTokens(),
	)

	return nftQueryCmd
}

// GetCmdQueryCollectionSupply queries the supply of a nft collection
func cmdQueryCollectionSupply() *cobra.Command {
	return &cobra.Command{
		Use:   "supply [denom]",
		Short: "total supply of a collection of NFTs",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the total count of NFTs that match a certain denomination.

Example:
$ %s query %s supply crypto-kitties
`, config.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryCollectionSupplyRequest{
				Denom: args[0],
			}

			res, err := queryClient.QueryCollectionSupply(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
}

// GetCmdQueryOwner queries all the NFTs owned by an account
func cmdQueryOwner() *cobra.Command {
	return &cobra.Command{
		Use:   "owner [account_address] [denom]",
		Short: "get the NFTs owned by an account address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the NFTs owned by an account address optionally filtered by the denom of the NFTs.

Example:
$ %s query %s owner cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
$ %s query %s owner cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p crypto-kitties
`, config.AppBinName, types.ModuleName, config.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryOwnerCollectionsRequest{
				Owner: args[0],
			}
			if len(args) > 1 {
				req.Denom = args[1]
			}

			res, err := queryClient.QueryOwnerCollections(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// GetCmdQueryCollection queries all the NFTs from a collection
func cmdQueryCollections() *cobra.Command {
	return &cobra.Command{
		Use:   "collection [denom]",
		Short: "get all the NFTs from a given collection",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get a list of all NFTs from a given collection.

Example:
$ %s query %s collection crypto-kitties
`, config.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryCollectionRequest{
				Denom: args[0],
			}

			res, err := queryClient.QueryCollection(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// GetCmdQueryDenoms queries all denoms
func cmdQueryDenoms() *cobra.Command {
	return &cobra.Command{
		Use:   "denoms",
		Short: "queries all denominations of all collections of NFTs",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Gets all denominations of all the available collections of NFTs that
are stored on the chain.

Example:
$ %s query %s denoms
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

			req := &types.QueryDenomsRequest{}

			res, err := queryClient.QueryDenoms(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// GetCmdQueryNFT queries a single NFTs from a collection
func cmdQueryNFT() *cobra.Command {
	return &cobra.Command{
		Use:   "token [denom] [token_id]",
		Short: "query a single NFT from a collection",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get an NFT from a collection that has the given ID (SHA-256 hex hash).

Example:
$ %s query %s token crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa
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

			denom := args[0]
			tokenID := args[1]
			req := &types.QueryNFTRequest{
				Denom:   denom,
				TokenId: tokenID,
			}

			res, err := queryClient.QueryNFT(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

func cmdQuerySubTokens() *cobra.Command {
	return &cobra.Command{
		Use:   "sub_tokens [denom] [token_id] [sub_token_ids]",
		Short: "query a sub tokens information of single NFT from a collection",
		Long: fmt.Sprintf(`Get sub tokens of NFT from a collection that has the given ID (SHA-256 hex hash).

Example:
$ %s query %s sub_tokens crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa 1,2,3
`, config.AppBinName, types.ModuleName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			denom := args[0]
			tokenID := args[1]

			subTokenIDsStr := strings.Split(args[2], ",")
			subTokenIDs := make([]uint64, len(subTokenIDsStr))
			for i, d := range subTokenIDsStr {
				subTokenID, err := strconv.ParseUint(d, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid subTokenID")
				}
				subTokenIDs[i] = subTokenID
			}

			req := &types.QuerySubTokensRequest{
				Denom:       denom,
				TokenID:     tokenID,
				SubTokenIDs: subTokenIDs,
			}

			res, err := queryClient.QuerySubTokens(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}
