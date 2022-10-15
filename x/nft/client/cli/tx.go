package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

const (
	flagTokenURI = "tokenURI"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		RunE:  client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdMintToken(),
		cmdUpdateToken(),
		cmdUpdateReserve(),
		cmdSendToken(),
		cmdBurnToken(),
	)

	return cmd
}

func cmdMintToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [denom] [token_id] [recipient] [quantity] [reserve] [allow_mint]",
		Short: "mint an NFT and set the owner to the recipient",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Mint an NFT from a given collection that has a 
specific id (SHA-256 hex hash) and set the ownership to a specific address.

Example:
$ %s tx %s mint crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa \
dx1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p 10 100000000000000000000del t --from mykey --tokenURI myuri
`,
				cmdcfg.AppBinName, types.ModuleName,
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
				return errors.EmptyTokenURI
			}

			quantity, err := strconv.ParseUint(args[3], 10, 32)
			if err != nil {
				return errors.InvalidQuantity
			}

			reserve, err := sdk.ParseCoinNormalized(args[4])
			if err != nil {
				return errors.InvalidReserve
			}

			var allowMint bool
			if args[5] == "t" {
				allowMint = true
			}

			msg := types.NewMsgMintToken(clientCtx.GetFromAddress(), denom, tokenID, tokenURI, allowMint, recipient, uint32(quantity), reserve)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagTokenURI, "", "URI for supplemental off-chain metadata (should return a JSON object)")

	return cmd
}

func cmdUpdateToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update_uri [token_id]",
		Short: "update the NFT token URI",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update the NFT token URI.

Example:
$ %s tx %s update_uri d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa \
--tokenURI path_to_token_URI_JSON --from mykey
`,
				cmdcfg.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokenID := args[0]
			tokenURI := viper.GetString(flagTokenURI)

			msg := types.NewMsgUpdateToken(clientCtx.GetFromAddress(), tokenID, tokenURI)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagTokenURI, "", "Extra properties available for querying")
	return cmd
}

func cmdUpdateReserve() *cobra.Command {
	return &cobra.Command{
		Use:   "update_reserve [token_id] [sub_token_ids] [new_reserve]",
		Short: "update reserve NFT",
		Long: strings.TrimSpace(
			fmt.Sprintf(`NFT Reserve Update an NFT from a given collection that has a 
			specific id (SHA-256 hex hash).

Example:
$ %s tx %s update_reserve d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa 1,2 1000del \
--from mykey
`,
				cmdcfg.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokenID := args[0]
			subTokenIDsStr := strings.Split(args[1], ",")
			subTokenIDs := make([]uint32, len(subTokenIDsStr))
			for i, d := range subTokenIDsStr {
				subTokenID, err := strconv.ParseUint(d, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid subTokenID")
				}
				subTokenIDs[i] = uint32(subTokenID)
			}

			newCoinReserve, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return errors.InvalidReserve
			}

			msg := types.NewMsgUpdateReserve(clientCtx.GetFromAddress(), tokenID, subTokenIDs, newCoinReserve)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

func cmdSendToken() *cobra.Command {
	return &cobra.Command{
		Use:   "send [sender] [recipient] [token_id] [sub_token_ids]",
		Short: "send a NFT to a recipient",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Send a NFT from a given collection that has a 
specific id (SHA-256 hex hash) to a specific recipient.

Example:
$ %s tx %s send 
dx1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p dx1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm \
crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa 1,2 \
--from mykey
`,
				cmdcfg.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(4),
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

			tokenID := args[2]

			subTokenIDsStr := strings.Split(args[3], ",")
			subTokenIDs := make([]uint32, len(subTokenIDsStr))
			for i, d := range subTokenIDsStr {
				subTokenID, err := strconv.ParseUint(d, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid subTokenID")
				}
				subTokenIDs[i] = uint32(subTokenID)
			}

			msg := types.NewMsgSendToken(sender, recipient, tokenID, subTokenIDs)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

func cmdBurnToken() *cobra.Command {
	return &cobra.Command{
		Use:   "burn [token_id] [sub_token_ids]",
		Short: "burn an NFT",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Burn (i.e permanently delete) an NFT from a given collection that has a 
			specific id (SHA-256 hex hash).

Example:
$ %s tx %s burn crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa 2,3\
--from mykey
`,
				cmdcfg.AppBinName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokenID := args[0]

			subTokenIDsStr := strings.Split(args[1], ",")
			subTokenIDs := make([]uint32, len(subTokenIDsStr))
			for i, d := range subTokenIDsStr {
				subTokenID, err := strconv.ParseUint(d, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid subTokenID")
				}
				subTokenIDs[i] = uint32(subTokenID)
			}

			msg := types.NewMsgBurnToken(clientCtx.GetFromAddress(), tokenID, subTokenIDs)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}
