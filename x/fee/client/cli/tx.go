package cli

import (
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

func GetTxCmd() *cobra.Command {
	nftTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Fee transactions subcommands",
		RunE:  client.ValidateCmd,
	}

	nftTxCmd.AddCommand(
		GetCmdSaveBaseDenomPrice(),
	)

	return nftTxCmd
}

// GetCmdSaveBaseDenomPrice is the CLI command for a MintNFT transaction
func GetCmdSaveBaseDenomPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save-price [denom] [price]",
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
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			price, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return types.ErrWrongPrice(args[1])
			}

			msg := types.NewMsgSaveBaseDenomPrice(args[0], price, clientCtx.GetFromAddress().String())
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
