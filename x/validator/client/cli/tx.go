package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	ethtypes "github.com/decimalteam/ethermint/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// default values
var (
	DefaultTokens         = sdk.TokensFromConsensusPower(1, ethtypes.PowerReduction)
	defaultAmount         = DefaultTokens.String() + cointypes.DefaultBaseDenom
	defaultCommissionRate = "0.1"
)

// GetTxCmd returns a root CLI command handler for all x/staking transaction commands.
func GetTxCmd() *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Staking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingTxCmd.AddCommand(
		NewCreateValidatorCmd(),
		NewEditValidatorCmd(),
		NewDelegateCmd(),
		NewDelegateNFTCmd(),
		NewUndelegateCmd(),
		NewUndelegateNFTCmd(),
		NewRedelegateCmd(),
		NewRedelegateNFTCmd(),
		NewSetOnlineCmd(),
		NewSetOfflineCmd(),
		NewCancelUndelegateCmd(),
		NewCancelUndelegateNFTCmd(),
		NewCancelRedelegateCmd(),
		NewCancelRedelegateNFTCmd(),
	)

	return stakingTxCmd
}

// NewCreateValidatorCmd returns a CLI command handler for creating a MsgCreateValidator transaction.
func NewCreateValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "create new validator initialized with a self-delegation to it",
		Long: fmt.Sprintf(`
Create new validator command example:

%s tx %s create-validator 
--moniker validator
--details='Hey, I am a new validator. Please delegate!'
--identity=AFAF00C4
--website=https://newvalidator.io
--security-contact=contact@newvalidator.io
--amount 100000del
--from mykey
--pubkey '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"cdafs02U0NcdgX1PigeBmxMNleH+kUCr+eEdnZnNSag="}' 
--commission-rate="0.10"
`, cmdcfg.AppBinName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildCreateValidatorMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetPublicKey())
	cmd.Flags().AddFlagSet(FlagSetAmount())
	cmd.Flags().AddFlagSet(flagSetDescriptionCreate())
	cmd.Flags().AddFlagSet(FlagSetCommissionCreate())

	cmd.Flags().String(FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", flags.FlagGenerateOnly))
	cmd.Flags().String(FlagNodeID, "", "The node's ID")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagPubKey)
	_ = cmd.MarkFlagRequired(FlagMoniker)

	return cmd
}

// NewEditValidatorCmd returns a CLI command handler for creating a MsgEditValidator transaction.
func NewEditValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-validator [validator-address] [reward-address]",
		Short: "edit an existing validator account",
		Args:  cobra.ExactArgs(2),
		Long: fmt.Sprintf(`Edit validator  

Example: 	
$ %s tx %s dxvaloper1w4m22z5nuvyaphcyccwcqx4vpm7f5q496yu9m5 dx1w4m22z5nuvyaphcyccwcqx4vpm7f5q49xkmgwl`,
			cmdcfg.AppBinName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			rewardAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			moniker, _ := cmd.Flags().GetString(FlagEditMoniker)
			identity, _ := cmd.Flags().GetString(FlagIdentity)
			website, _ := cmd.Flags().GetString(FlagWebsite)
			security, _ := cmd.Flags().GetString(FlagSecurityContact)
			details, _ := cmd.Flags().GetString(FlagDetails)
			description := types.NewDescription(moniker, identity, website, security, details)

			msg := types.NewMsgEditValidator(valAddress, rewardAddress, description)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(flagSetDescriptionEdit())
	cmd.Flags().AddFlagSet(flagSetCommissionUpdate())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewDelegateCmd returns a CLI command handler for creating a MsgDelegate transaction.
func NewDelegateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate [validator-addr] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Delegate liquid tokens to a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Delegate an amount of liquid coins to a validator from your wallet.

Example:
$ %s tx %s delegate dxvaloper1q3pjfs20lcakezjmd3l8a4fcnzq9p69hcc28zh 1000del --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDelegate(delAddr, valAddr, amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewDelegateNFTCmd returns a CLI command handler for creating a MsgDelegateNFT transaction.
func NewDelegateNFTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate-nft [validator-addr] [tokenID] [sub_token_ids]",
		Args:  cobra.ExactArgs(3),
		Short: "Delegate nft tokens to a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Delegate nft subtokens to a validator from your wallet.

Example:
$ %s tx %s delegate-nft dxvaloper1q3pjfs20lcakezjmd3l8a4fcnzq9p69hcc28zh shinigami 10,12 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			subTokenIDsStr := strings.Split(args[2], ",")
			subTokenIDs := make([]uint32, len(subTokenIDsStr))
			for i, d := range subTokenIDsStr {
				subTokenID, err := strconv.ParseUint(d, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid quantity")
				}
				subTokenIDs[i] = uint32(subTokenID)
			}

			tokenID := args[1]

			msg := types.NewMsgDelegateNFT(delAddr, valAddr, tokenID, subTokenIDs)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUndelegateCmd returns a CLI command handler for creating a MsgUndelegate transaction.
func NewUndelegateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "undelegate [validator-addr] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Undelegate liquid tokens to a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Undelegate an amount of liquid coins from a validator to your wallet.

Example:
$ %s tx %s undelegate dxvaloper1q3pjfs20lcakezjmd3l8a4fcnzq9p69hcc28zh 999del --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgUndelegate(delAddr, valAddr, amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUndelegateNFTCmd returns a CLI command handler for creating a MsgDelegateNFT transaction.
func NewUndelegateNFTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "undelegate-nft [validator-addr] [tokenID] [sub_token_ids]",
		Args:  cobra.ExactArgs(3),
		Short: "Undelegate nft tokens from validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Undelegate nft subtokens from validator to your wallet.

Example:
$ %s tx %s delegate-nft dxvaloper1q3pjfs20lcakezjmd3l8a4fcnzq9p69hcc28zh shinigami 10,12 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			subTokenIDsStr := strings.Split(args[2], ",")
			subTokenIDs := make([]uint32, len(subTokenIDsStr))
			for i, d := range subTokenIDsStr {
				subTokenID, err := strconv.ParseUint(d, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid quantity")
				}
				subTokenIDs[i] = uint32(subTokenID)
			}

			tokenID := args[1]

			msg := types.NewMsgUndelegateNFT(delAddr, valAddr, tokenID, subTokenIDs)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewRedelegateCmd returns a CLI command handler for creating a MsgRedelegate transaction.
func NewRedelegateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redelegate [src-validator-addr] [dst-validator-addr] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Redelegate liquid tokens from validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Redelegate an amount of liquid coins from one validator to another validator.

Example:
$ %s tx %s redelegate dxvaloper1qqwy22055u6yrem8s8gv9j2ndv2hv2z9magtat dxvaloper1qqwy22055u6yrem8s8gv9j2ndv2hv2z9magtat 999del --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			srcValAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			dstValAddr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgRedelegate(delAddr, srcValAddr, dstValAddr, amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewRedelegateNFTCmd returns a CLI command handler for creating a MsgRedelegateNFT transaction.
func NewRedelegateNFTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redelegate-nft [src-validator-addr] [dst-validator-addr] [tokenID] [sub_token_ids]",
		Args:  cobra.ExactArgs(4),
		Short: "Redelegate nft tokens from validator to another validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Redelegate nft subtokens from validator to another validator.

Example:
$ %s tx %s delegate-nft dxvaloper1qqwy22055u6yrem8s8gv9j2ndv2hv2z9magtat dxvaloper1qqwy22055u6yrem8s8gv9j2ndv2hv2z9magtat shinigami 10,12 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			srcValAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			dstValAddr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			subTokenIDsStr := strings.Split(args[3], ",")
			subTokenIDs := make([]uint32, len(subTokenIDsStr))
			for i, d := range subTokenIDsStr {
				subTokenID, err := strconv.ParseUint(d, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid quantity")
				}
				subTokenIDs[i] = uint32(subTokenID)
			}

			tokenID := args[2]

			msg := types.NewMsgRedelegateNFT(delAddr, srcValAddr, dstValAddr, tokenID, subTokenIDs)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewSetOnlineCmd returns a CLI command handler for creating a MsgSetOnline transaction.
func NewSetOnlineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-online [validator-addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Set online validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Validator status update.

Example:
$ %s tx %s set-online --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetOnline(valAddr)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewSetOfflineCmd returns a CLI command handler for creating a MsgSetOnline transaction.
func NewSetOfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-offline [validator-addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Set offline validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Validator status.

Example:
$ %s tx %s set-offline --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetOffline(valAddr)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewCancelUndelegateCmd returns a CLI command handler for creating a MsgCancelUndelegate transaction.
func NewCancelUndelegateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-undelegate [validator-addr] [amount] [create-height]",
		Args:  cobra.ExactArgs(3),
		Short: "Cancel undelegate liquid tokens to a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel undelegate an amount of liquid coins from a validator.

Example:
$ %s tx %s cancel-undelegate dxvaloper1q3pjfs20lcakezjmd3l8a4fcnzq9p69hcc28zh 999del 23 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			creationHeight, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid height")
			}

			msg := types.NewMsgCancelUndelegation(delAddr, valAddr, creationHeight, amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewCancelUndelegateNFTCmd returns a CLI command handler for creating a MsgCancelUndelegateNFT transaction.
func NewCancelUndelegateNFTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-undelegate-nft [validator-addr] [tokenID] [sub_token_ids] [creation height]",
		Args:  cobra.ExactArgs(4),
		Short: "Cancel undelegate nft tokens from validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel undelegate nft tokens from a validator.

Example:
$ %s tx %s cancel-undelegate-nft dxvaloper1q3pjfs20lcakezjmd3l8a4fcnzq9p69hcc28zh shinigami 10,12 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			subTokenIDsStr := strings.Split(args[2], ",")
			subTokenIDs := make([]uint32, len(subTokenIDsStr))
			for i, d := range subTokenIDsStr {
				subTokenID, err := strconv.ParseUint(d, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid quantity")
				}
				subTokenIDs[i] = uint32(subTokenID)
			}

			tokenID := args[1]

			creationHeight, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid height")
			}

			msg := types.NewMsgCancelUndelegationNFT(delAddr, valAddr, creationHeight, tokenID, subTokenIDs)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewCancelRedelegateCmd returns a CLI command handler for creating a MsgCancelRedelegate transaction.
func NewCancelRedelegateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-redelegate [src-validator-addr] [dst-validator-addr] [amount] [create-height]",
		Args:  cobra.ExactArgs(4),
		Short: "Cancel redelegate liquid tokens to a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel redelegate an amount of liquid coins from a validator.

Example:
$ %s tx %s cancel-redelegate dxvaloper1q3pjfs20lcakezjmd3l8a4fcnzq9p69hcc28zh 999del 23 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			srcValAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			dstValAddr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			creationHeight, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid height")
			}

			msg := types.NewMsgCancelRedelegation(delAddr, srcValAddr, dstValAddr, creationHeight, amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewCancelRedelegateNFTCmd returns a CLI command handler for creating a MsgCancelRedelegateNFT transaction.
func NewCancelRedelegateNFTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-redelegate-nft [src-validator-addr] [dst-validator-addr] [tokenID] [sub_token_ids] [creation height]",
		Args:  cobra.ExactArgs(5),
		Short: "Cancel redelegate nft tokens from validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel undelegate nft tokens from a validator.

Example:
$ %s tx %s cancel-undelegate-nft dxvaloper1q3pjfs20lcakezjmd3l8a4fcnzq9p69hcc28zh dxvaloper1qqwy22055u6yrem8s8gv9j2ndv2hv2z9magtat shinigami 10,12 32 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			srcValAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			dstValAddr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			subTokenIDsStr := strings.Split(args[3], ",")
			subTokenIDs := make([]uint32, len(subTokenIDsStr))
			for i, d := range subTokenIDsStr {
				subTokenID, err := strconv.ParseUint(d, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid quantity")
				}
				subTokenIDs[i] = uint32(subTokenID)
			}

			tokenID := args[2]

			creationHeight, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid height")
			}

			msg := types.NewMsgCancelRedelegationNFT(delAddr, srcValAddr, dstValAddr, creationHeight, tokenID, subTokenIDs)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func newBuildCreateValidatorMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgCreateValidator, error) {
	fAmount, _ := fs.GetString(FlagAmount)
	amount, err := sdk.ParseCoinNormalized(fAmount)
	if err != nil {
		return txf, nil, err
	}

	valAddr := clientCtx.GetFromAddress()
	pkStr, err := fs.GetString(FlagPubKey)
	if err != nil {
		return txf, nil, err
	}

	var pk cryptotypes.PubKey
	if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(pkStr), &pk); err != nil {
		return txf, nil, err
	}

	moniker, _ := fs.GetString(FlagMoniker)
	identity, _ := fs.GetString(FlagIdentity)
	website, _ := fs.GetString(FlagWebsite)
	security, _ := fs.GetString(FlagSecurityContact)
	details, _ := fs.GetString(FlagDetails)
	description := types.NewDescription(
		moniker,
		identity,
		website,
		security,
		details,
	)
	// get the initial validator commission parameters
	fcomissionRates, err := fs.GetString(FlagCommissionRate)
	if err != nil {
		return txf, nil, err
	}

	commissionRates, err := sdk.NewDecFromStr(fcomissionRates)
	if err != nil {
		return txf, nil, err
	}

	msg, err := types.NewMsgCreateValidator(
		sdk.ValAddress(valAddr), clientCtx.FromAddress, pk, description, commissionRates, amount,
	)

	if err != nil {
		return txf, nil, err
	}
	if err := msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	genOnly, _ := fs.GetBool(flags.FlagGenerateOnly)
	if genOnly {
		ip, _ := fs.GetString(FlagIP)
		nodeID, _ := fs.GetString(FlagNodeID)

		if nodeID != "" && ip != "" {
			txf = txf.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
		}
	}

	return txf, msg, nil
}
