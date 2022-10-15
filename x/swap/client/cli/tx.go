package cli

import (
	"fmt"
	"strconv"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
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
		NewSwapInitializeCmd(),
		NewSwapRedeemCmd(),
		NewChainActivateCmd(),
		NewChainDeactivateCmd(),
	)

	return coinCmd
}

func NewSwapInitializeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [recipient] [amount] [token_symbol] [tx_number] [from_chain] [dest_chain]",
		Short: "Swap initialize",
		Long: fmt.Sprintf(`Init swap from our blockchain to other blockchain.

Example: 	
$ %s tx %s init 0x12345 1000 bnb 128943 1 3 --from mykey`, config.AppBinName, types.ModuleName),
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			recipient := args[0]
			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid amount")
			}
			tokenSymbol := args[2]
			txNumber := args[3]
			fromChain, err := strconv.ParseUint(args[4], 10, 32)
			if err != nil {
				return err
			}
			destChain, err := strconv.ParseUint(args[5], 10, 32)
			if err != nil {
				return err
			}

			msg := types.NewMsgInitializeSwap(from, recipient, amount, tokenSymbol, txNumber, uint32(fromChain), uint32(destChain))
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

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}

func NewSwapRedeemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem [from] [recipient] [amount] [token_symbol] [tx_number] [from_chain] [dest_chain] [v] [r] [s]",
		Short: "Swap redeem",
		Long: fmt.Sprintf(`Redeem swap from other blockchain to our blockchain.
'from' - address of account from other blockchain
'recipient' - address of account in our blockchain
v,r,s - hex encoded proof

Example: 	
$ %s tx %s redeem 0x12345 dx1..addr 1000 del 128943 3 1 0 ae45.. df350.. --from mykey`, config.AppBinName, types.ModuleName),
		Args: cobra.ExactArgs(10),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress()

			from := args[0]
			recipient := args[1]
			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid amount")
			}
			tokenSymbol := args[3]
			txNumber := args[4]
			fromChain, err := strconv.ParseUint(args[5], 10, 32)
			if err != nil {
				return err
			}
			destChain, err := strconv.ParseUint(args[6], 10, 32)
			if err != nil {
				return err
			}

			v, err := strconv.Atoi(args[7])
			if err != nil {
				return err
			}

			msg := types.NewMsgRedeemSwap(
				sender, from, recipient, amount, tokenSymbol, txNumber, uint32(fromChain), uint32(destChain), uint32(v), args[8], args[9])
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

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}

func NewChainActivateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain-activate [number] [name]",
		Short: "Activate chain",
		Long: fmt.Sprintf(`Activate/create blockchain record for swap.

Example: 	
$ %s tx %s chain-activate 2 "Ethereum" --from mykey`, config.AppBinName, types.ModuleName),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress()

			number, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}

			name := args[1]

			msg := types.NewMsgActivateChain(sender, uint32(number), name)
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

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}

func NewChainDeactivateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain-deactivate [number]",
		Short: "Deactivate chain",
		Long: fmt.Sprintf(`Deactivate blockchain record for swap.

Example: 	
$ %s tx %s chain-deactivate 2 --from mykey`, config.AppBinName, types.ModuleName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress()

			number, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeactivateChain(sender, uint32(number))
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

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}
