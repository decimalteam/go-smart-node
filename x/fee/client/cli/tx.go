package cli

import (
	"bitbucket.org/decimalteam/go-smart-node/x/fee/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		Short: "save price to fee module, sender must be oracle",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			price, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return errors.WrongPrice
			}

			msg := types.NewMsgSaveBaseDenomPrice(clientCtx.GetFromAddress().String(), args[0], price)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
