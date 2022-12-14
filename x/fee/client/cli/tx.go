package cli

import (
	"time"

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
		GetCmdUpdatePrice(),
	)

	return nftTxCmd
}

// GetCmdSaveBaseDenomPrice is the CLI command for a MintNFT transaction
func GetCmdUpdatePrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-price [denom] [quote] [price]",
		Short: "update price for pair denom/quote in fee module, sender must be oracle",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			price, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return errors.WrongPrice
			}

			msg := types.NewMsgUpdateCoinPrices(clientCtx.GetFromAddress().String(),
				[]types.CoinPrice{
					{
						Denom:     args[0],
						Quote:     args[1],
						Price:     price,
						UpdatedAt: time.Now(),
					},
				})

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
