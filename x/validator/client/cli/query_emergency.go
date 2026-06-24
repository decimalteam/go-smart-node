package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// GetCmdQueryEmergencyStatus implements the emergency-status query command.
func GetCmdQueryEmergencyStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "emergency-status",
		Args:  cobra.NoArgs,
		Short: "Query the current emergency hard-halt and soft-freeze state",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query whether a hard halt is scheduled (and at what height) and whether
the chain is currently soft-frozen.

Example:
$ %s query %s emergency-status
`, version.AppName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.EmergencyStatus(cmd.Context(), &types.QueryEmergencyStatusRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
