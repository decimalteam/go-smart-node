package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	stakingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the staking module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingQueryCmd.AddCommand(
		GetCmdQueryValidator(),
		GetCmdQueryValidators(),
		GetCmdQueryValidatorDelegations(),
		GetCmdQueryValidatorUndelegations(),
		GetCmdQueryValidatorRedelegations(),
		GetCmdQueryDelegations(),
		GetCmdQueryUndelegation(),
		GetCmdQueryRedelegations(),
		GetCmdQueryDelegatorDelegations(),
		GetCmdQueryDelegatorUndelegations(),
		GetCmdQueryDelegatorRedelegations(),
		GetCmdQueryDelegatorValidators(),
		GetCmdQueryDelegatorValidator(),
		GetCmdQueryHistoricalInfo(),
		GetCmdQueryParams(),
		GetCmdQueryPool(),
		GetCmdQueryCustomCoinPrice(),
		GetCmdQueryTotalCustomCoin(),
	)

	return stakingQueryCmd
}

// GetCmdQueryValidator implements the validator query command.
func GetCmdQueryValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator [validator-addr]",
		Short: "Query a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about an individual validator.

Example:
$ %s query %s validator dx1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryValidatorRequest{Validator: addr.String()}
			res, err := queryClient.Validator(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Validator)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryValidators implements the query all validators command.
func GetCmdQueryValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Short: "Query for all validators",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all validators on a network.

Example:
$ %s query %s validators
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			result, err := queryClient.Validators(cmd.Context(), &types.QueryValidatorsRequest{
				// Leaving status empty on purpose to query all validators.
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "validators")

	return cmd
}

// GetCmdQueryValidatorDelegations implements the query all delegatations from a validator command.
func GetCmdQueryValidatorDelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegations-from [validator-addr]",
		Short: "Query all delegatations from a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegations _from_ a validator.

Example:
$ %s query %s delegations-from dxvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryValidatorDelegationsRequest{
				Validator:  valAddr.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.ValidatorDelegations(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "validator delegations")

	return cmd
}

// GetCmdQueryValidatorUndelegations implements the query all unbonding delegatations from a validator command.
func GetCmdQueryValidatorUndelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "undelegations-from [validator-addr]",
		Short: "Query all unbonding delegatations from a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegations that are unbonding _from_ a validator.

Example:
$ %s query %s undelegations-from dxvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryValidatorUndelegationsRequest{
				Validator:  valAddr.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.ValidatorUndelegations(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "unbonding delegations")

	return cmd
}

// GetCmdQueryValidatorRedelegations implements the query all redelegatations
// from a validator command.
func GetCmdQueryValidatorRedelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redelegations-from [validator-addr]",
		Short: "Query all outgoing redelegatations from a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegations that are redelegating _from_ a validator.

Example:
$ %s query %s redelegations-from dxvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			valSrcAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryValidatorRedelegationsRequest{
				Validator:  valSrcAddr.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.ValidatorRedelegations(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "validator redelegations")

	return cmd
}

// GetCmdQueryDelegations implements the query all delegations from a validator-delegator pair command.
func GetCmdQueryDelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-delegator-delegations [validator-addr] [delegator-addr]",
		Short: "Query all delegations from a delegator-validator pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegations _from_ a delegator-validator pair.

Example:
$ %s query %s validator-delegator-delegations dxvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj dx1mq0fvnxlkqh4xv5sxgrsf3lrcfvnrs0e9ajma5
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			delAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			params := &types.QueryDelegationsRequest{
				Validator: valAddr.String(),
				Delegator: delAddr.String(),
			}

			res, err := queryClient.Delegations(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "validator-delegator pair delegations")

	return cmd
}

// GetCmdQueryUndelegation implements the query unbonding delegation from a validator-delegator pair command.
func GetCmdQueryUndelegation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-delegator-undelegation [validator-addr] [delegator-addr]",
		Short: "Query undelegation from a delegator-validator pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query undelegation _from_ a delegator-validator pair.

Example:
$ %s query %s validator-delegator-undelegation dxvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj dx1mq0fvnxlkqh4xv5sxgrsf3lrcfvnrs0e9ajma5
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			delAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			params := &types.QueryUndelegationRequest{
				Validator: valAddr.String(),
				Delegator: delAddr.String(),
			}

			res, err := queryClient.Undelegation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "validator-delegator pair undelegation")

	return cmd
}

// GetCmdQueryRedelegations implements the query all redelegations from a validator-delegator pair command.
func GetCmdQueryRedelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-delegator-redelegations [validator-addr] [delegator-addr]",
		Short: "Query all redelegations from a delegator-validator pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all redelegations _from_ a delegator-validator pair.

Example:
$ %s query %s validator-delegator-redelegations dxvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj dx1mq0fvnxlkqh4xv5sxgrsf3lrcfvnrs0e9ajma5
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			delAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			params := &types.QueryRedelegationsRequest{
				Validator: valAddr.String(),
				Delegator: delAddr.String(),
			}

			res, err := queryClient.Redelegations(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "validator-delegator pair undelegation")

	return cmd
}

// GetCmdQueryDelegatorDelegations implements the query all delegations from a delegator command.
func GetCmdQueryDelegatorDelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator-delegations [delegator-addr]",
		Short: "Query all delegations from a delegator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegations _from_ a delegator.

Example:
$ %s query %s delegator-delegations dx1mq0fvnxlkqh4xv5sxgrsf3lrcfvnrs0e9ajma5
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryDelegatorDelegationsRequest{
				Delegator: delAddr.String(),
			}

			res, err := queryClient.DelegatorDelegations(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator delegations")

	return cmd
}

// GetCmdQueryDelegatorUndelegations implements the query all unbonding delegations from a delegator command.
func GetCmdQueryDelegatorUndelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator-undelegations [delegator-addr]",
		Short: "Query all undelegations from a delegator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query unbonding delegations _from_ a delegator.

Example:
$ %s query %s delegator-undelegations dx1mq0fvnxlkqh4xv5sxgrsf3lrcfvnrs0e9ajma5
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryDelegatorUndelegationsRequest{
				Delegator: delAddr.String(),
			}

			res, err := queryClient.DelegatorUndelegations(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator undelegations")

	return cmd
}

// GetCmdQueryDelegatorRedelegations implements the query all redelegations from a delegator command.
func GetCmdQueryDelegatorRedelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator-redelegations [delegator-addr]",
		Short: "Query all redelegations from a delegator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query redelegations _from_ a delegator.

Example:
$ %s query %s delegator-redelegations dx1mq0fvnxlkqh4xv5sxgrsf3lrcfvnrs0e9ajma5
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryDelegatorRedelegationsRequest{
				Delegator: delAddr.String(),
			}

			res, err := queryClient.DelegatorRedelegations(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator delegations")

	return cmd
}

// GetCmdQueryDelegatorValidators implements the query all validators delegated by the delegator command.
func GetCmdQueryDelegatorValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator-validators [delegator-addr]",
		Short: "Query all validators delegated by the delegator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query validators delegated by the delegator.

Example:
$ %s query %s delegator-validators dx1mq0fvnxlkqh4xv5sxgrsf3lrcfvnrs0e9ajma5
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryDelegatorValidatorsRequest{
				Delegator: delAddr.String(),
			}

			res, err := queryClient.DelegatorValidators(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator validators")

	return cmd
}

// GetCmdQueryDelegatorValidator implements the query validator info for given delegator validator pair command.
func GetCmdQueryDelegatorValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator-validator [delegator-addr] [validator-addr]",
		Short: "Query validator info for given delegator validator pair.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query validator info for given delegator validator pair..

Example:
$ %s query %s delegator-validator dx1mq0fvnxlkqh4xv5sxgrsf3lrcfvnrs0e9ajma5 dxvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			valAddr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			params := &types.QueryDelegatorValidatorRequest{
				Delegator: delAddr.String(),
				Validator: valAddr.String(),
			}

			res, err := queryClient.DelegatorValidator(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator validator")

	return cmd
}

// GetCmdQueryHistoricalInfo implements the historical info query command
func GetCmdQueryHistoricalInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "historical-info [height]",
		Args:  cobra.ExactArgs(1),
		Short: "Query historical info at given height",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query historical info at given height.

Example:
$ %s query %s historical-info 5
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			height, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil || height < 0 {
				return fmt.Errorf("height argument provided must be a non-negative-integer: %v", err)
			}

			params := &types.QueryHistoricalInfoRequest{Height: height}
			res, err := queryClient.HistoricalInfo(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Hist)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryPool implements the pool query command.
func GetCmdQueryPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool",
		Args:  cobra.NoArgs,
		Short: "Query the current staking pool values",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values for amounts stored in the staking pool.

Example:
$ %s query %s pool
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Pool(cmd.Context(), &types.QueryPoolRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Pool)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current staking parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as staking parameters.

Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryCustomCoinPrice query custom coin price for delegation command
func GetCmdQueryCustomCoinPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom-coin-price [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query custom coin price for delegation ",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query custom coin price for delegation.

Example:
$ %s query %s custom-coin-price satoshicoin
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			denom := args[0]

			params := &types.QueryCustomCoinPriceRequest{Denom: denom}

			res, err := queryClient.CustomCoinPrice(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTotalCustomCoin query bonded custom coin amount info command
func GetCmdQueryTotalCustomCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom-coin-bonded [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query bonded custom coin amount info",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query bonded custom coin amount info .

Example:
$ %s query %s custom-coin-bonded satoshicoin
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			denom := args[0]

			params := &types.QueryTotalCustomCoinRequest{Denom: denom}

			res, err := queryClient.TotalCustomCoin(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
