package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

const flagHaltHeight = "height"

// NewHaltChainCmd returns a CLI command handler for creating a MsgHaltChain transaction.
func NewHaltChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "halt-chain [reason]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Emergency: schedule a hard consensus halt of the chain (halt-admin only)",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Schedule a hard consensus halt of the chain. The chain stops producing
blocks at the target height and stays halted until operators restart their nodes
with --unsafe-skip-halt and the halt-admin broadcasts "resume-chain".

Example:
$ %s tx %s halt-chain "security incident" --from haltadmin
$ %s tx %s halt-chain "scheduled maintenance" --height 1250000 --from haltadmin
`, version.AppName, types.ModuleName, version.AppName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			height, err := cmd.Flags().GetInt64(flagHaltHeight)
			if err != nil {
				return err
			}

			var reason string
			if len(args) > 0 {
				reason = args[0]
			}

			msg := types.NewMsgHaltChain(clientCtx.GetFromAddress(), height, reason)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Int64(flagHaltHeight, 0, "Block height at which to halt (0 = next block)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewResumeChainCmd returns a CLI command handler for creating a MsgResumeChain transaction.
func NewResumeChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume-chain",
		Args:  cobra.NoArgs,
		Short: "Emergency: clear a scheduled or active hard halt (halt-admin only)",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Clear a scheduled or active hard halt. After a chain has actually halted,
this transaction can only be included once operators restart with --unsafe-skip-halt.

Example:
$ %s tx %s resume-chain --from haltadmin
`, version.AppName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgResumeChain(clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewFreezeChainCmd returns a CLI command handler for creating a MsgFreezeChain transaction.
func NewFreezeChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "freeze-chain [reason]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Emergency: soft-freeze the chain, rejecting all non-emergency transactions (halt-admin only)",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Soft-freeze the chain. The chain keeps producing empty blocks but rejects
all transactions except emergency admin messages. Reversible with "unfreeze-chain"
without a node restart.

Example:
$ %s tx %s freeze-chain "investigating anomaly" --from haltadmin
`, version.AppName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var reason string
			if len(args) > 0 {
				reason = args[0]
			}

			msg := types.NewMsgFreezeChain(clientCtx.GetFromAddress(), reason)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUnfreezeChainCmd returns a CLI command handler for creating a MsgUnfreezeChain transaction.
func NewUnfreezeChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unfreeze-chain",
		Args:  cobra.NoArgs,
		Short: "Emergency: lift a soft freeze (halt-admin only)",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Lift a soft freeze so the chain accepts transactions again.

Example:
$ %s tx %s unfreeze-chain --from haltadmin
`, version.AppName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUnfreezeChain(clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
