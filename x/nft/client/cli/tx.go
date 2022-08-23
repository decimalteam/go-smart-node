package cli

import (
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Edit metadata flags
const (
	flagTokenURI = "tokenURI"
)

func GetTxCmd() *cobra.Command {
	nftTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "NFT transactions subcommands",
		RunE:  client.ValidateCmd,
	}

	nftTxCmd.AddCommand(
		GetCmdMintNFT(),
		GetCmdTransferNFT(),
		GetCmdEditNFTMetadata(),
		GetCmdBurnNFT(),
		GetCmdUpdateReserveNFT(),
	)

	return nftTxCmd
}

// GetCmdMintNFT is the CLI command for a MintNFT transaction
func GetCmdMintNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [denom] [token_id] [recipient] [quantity] [reserve] [allow_mint]",
		Short: "mint an NFT and set the owner to the recipient",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Mint an NFT from a given collection that has a 
specific id (SHA-256 hex hash) and set the ownership to a specific address.

Example:
$ %s tx %s mint crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa \
dx1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --from mykey
`,
				config.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denom := args[0]
			tokenID := args[1]

			recipient, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			tokenURI := viper.GetString(flagTokenURI)
			if tokenURI == "" {
				return types.ErrEmptyTokenURI()
			}

			quantity, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.ErrInvalidQuantity(args[2])
			}

			coinReserve, err := sdk.ParseCoinNormalized(args[4])
			if err != nil {
				return types.ErrInvalidReserve(args[4])
			}

			var allowMint bool
			if args[5] == "t" {
				allowMint = true
			}

			msg := types.NewMsgMintNFT(clientCtx.GetFromAddress(), recipient, tokenID, denom, tokenURI, quantity, coinReserve, allowMint)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagTokenURI, "", "URI for supplemental off-chain metadata (should return a JSON object)")

	return cmd
}

// GetCmdTransferNFT is the CLI command for sending a TransferNFT transaction
func GetCmdTransferNFT() *cobra.Command {
	return &cobra.Command{
		Use:   "transfer [sender] [recipient] [denom] [token_id] [sub_token_ids]",
		Short: "transfer a NFT to a recipient",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Transfer a NFT from a given collection that has a 
specific id (SHA-256 hex hash) to a specific recipient.

Example:
$ %s tx %s transfer 
dx1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p dx1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm \
crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa 1,2 \
--from mykey
`,
				config.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			recipient, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			denom := args[2]
			tokenID := args[3]

			subTokenIDsStr := strings.Split(args[4], ",")
			subTokenIDs := make([]uint64, len(subTokenIDsStr))
			for i, d := range subTokenIDsStr {
				subTokenID, err := strconv.ParseUint(d, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid subTokenID")
				}
				subTokenIDs[i] = subTokenID
			}

			msg := types.NewMsgTransferNFT(sender, recipient, denom, tokenID, subTokenIDs)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

// GetCmdEditNFTMetadata is the CLI command for sending an EditMetadata transaction
func GetCmdEditNFTMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit_metadata [denom] [token_id]",
		Short: "edit the metadata of an NFT",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Edit the metadata of an NFT from a given collection that has a 
specific id (SHA-256 hex hash).

Example:
$ %s tx %s edit-metadata crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa \
--tokenURI path_to_token_URI_JSON --from mykey
`,
				config.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denom := args[0]
			tokenID := args[1]
			tokenURI := viper.GetString(flagTokenURI)

			msg := types.NewMsgEditNFTMetadata(clientCtx.GetFromAddress(), tokenID, denom, tokenURI)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagTokenURI, "", "Extra properties available for querying")
	return cmd
}

// GetCmdBurnNFT is the CLI command for sending a BurnNFT transaction
func GetCmdBurnNFT() *cobra.Command {
	return &cobra.Command{
		Use:   "burn [denom] [token_id] [sub_token_ids]",
		Short: "burn an NFT",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Burn (i.e permanently delete) an NFT from a given collection that has a 
			specific id (SHA-256 hex hash).

Example:
$ %s tx %s burn crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa 2,3\
--from mykey
`,
				config.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

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

			msg := types.NewMsgBurnNFT(clientCtx.GetFromAddress(), tokenID, denom, subTokenIDs)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

// GetCmdUpdateReserveNFT is the CLI command for sending a BurnNFT transaction
func GetCmdUpdateReserveNFT() *cobra.Command {
	return &cobra.Command{
		Use:   "update_reserve [denom] [token_id] [sub_token_ids] [new_reserve]",
		Short: "update reserve NFT",
		Long: strings.TrimSpace(
			fmt.Sprintf(`NFT Reserve Update an NFT from a given collection that has a 
			specific id (SHA-256 hex hash).

Example:
$ %s tx %s update_reserve crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa 1,2 1000del \
--from mykey
`,
				config.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

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

			newCoinReserve, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return types.ErrInvalidReserve(args[3])
			}

			msg := types.NewMsgUpdateReserveNFT(clientCtx.GetFromAddress(), tokenID, denom, subTokenIDs, newCoinReserve)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}
