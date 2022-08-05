package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
)

// GetTxCmd returns the transaction commands for the module.
func GetTxCmd() *cobra.Command {
	coinCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	coinCmd.AddCommand(
		NewCreateWalletCmd(),
		NewCreateTransactionCmd(),
	)

	return coinCmd
}

func NewCreateWalletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-wallet [owners] [weights] [threshold]",
		Short: "Creates new multi signature wallet. Owners must be list of addresses splitted by comma; weights must be list of ints,",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from            = clientCtx.GetFromAddress()
				ownersStrings   = strings.Split(args[0], ",")
				weightsStrings  = strings.Split(args[1], ",")
				thresholdString = args[2]
			)

			if len(ownersStrings) != len(weightsStrings) {
				return fmt.Errorf("count of owners and weights must be same, but %d != %d", len(ownersStrings), len(weightsStrings))
			}

			// check owners for valid addresses, for duplicates
			ownersDups := make(map[string]bool)
			for i, address := range ownersStrings {
				_, err = sdk.AccAddressFromBech32(address)
				if err != nil {
					return fmt.Errorf("owner address %s at pos %d: %s", address, i+1, err.Error())
				}
				if ownersDups[address] {
					return fmt.Errorf("owner address %s duplicates", address)
				}
				ownersDups[address] = true
			}

			weights := make([]uint64, len(weightsStrings))
			for i, weightString := range weightsStrings {
				weight, err := strconv.ParseUint(weightString, 10, 64)
				if err != nil {
					return fmt.Errorf("weight %s and pos %d: %s", weightString, i+1, err)
				}
				weights[i] = weight
			}

			threshold, err := strconv.ParseUint(thresholdString, 10, 64)
			if err != nil {
				return fmt.Errorf("threshold %s: %s", thresholdString, err.Error())
			}

			msg := types.NewMsgCreateWallet(from, ownersStrings, weights, threshold)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	// workaround for cosmos
	cmd.Flags().String(flags.FlagChainID, "", "network chain id")

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewCreateTransactionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-transaction [wallet] [receiver] [coins]",
		Short: "Creates transaction for multisig wallet to send coins. Transaction sender must be in wallet owners",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from        = clientCtx.GetFromAddress()
				wallet      = args[0]
				receiver    = args[1]
				coinsString = args[2]
			)

			if _, err := sdk.AccAddressFromBech32(wallet); err != nil {
				return fmt.Errorf("invalid wallet address %s: %s", wallet, err.Error())
			}

			if _, err := sdk.AccAddressFromBech32(receiver); err != nil {
				return fmt.Errorf("invalid receiver address %s: %s", receiver, err.Error())
			}

			coins, err := sdk.ParseCoinsNormalized(coinsString)
			if err != nil {
				return fmt.Errorf("error whil parse coins %s: %s", coinsString, err.Error())
			}

			msg := types.NewMsgCreateTransaction(from, wallet, receiver, coins)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	// workaround for cosmos
	cmd.Flags().String(flags.FlagChainID, "", "network chain id")

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewSignTransactionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign-transaction [tx_id]",
		Short: "Sign transaction for multisig wallet to send coins",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from = clientCtx.GetFromAddress()
				txID = args[0]
			)

			msg := types.NewMsgSignTransaction(from, txID)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	// workaround for cosmos
	cmd.Flags().String(flags.FlagChainID, "", "network chain id")

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
