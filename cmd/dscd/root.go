package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	tmcfg "github.com/tendermint/tendermint/config"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/snapshots"
	snapshottypes "github.com/cosmos/cosmos-sdk/snapshots/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	ethermintclient "github.com/decimalteam/ethermint/client"
	"github.com/decimalteam/ethermint/client/debug"
	"github.com/decimalteam/ethermint/encoding"
	ethermintserver "github.com/decimalteam/ethermint/server"
	servercfg "github.com/decimalteam/ethermint/server/config"
	srvflags "github.com/decimalteam/ethermint/server/flags"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	dsckr "bitbucket.org/decimalteam/go-smart-node/crypto/keyring"
)

const (
	EnvPrefix = "DECIMAL"
)

// NewRootCmd creates a new root command for dscd. It is called once in the main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {

	// Initialize client context
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome).
		WithKeyringOptions(dsckr.Option()).
		WithViper(EnvPrefix)

	// Initialize root command
	rootCmd := &cobra.Command{
		Use:   cmdcfg.AppBinName,
		Short: "Decimal Smart Chain daemon",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {

			// Set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			// Disable ledger temporarily
			//useLedger, _ := cmd.Flags().GetBool(flags.FlagUseLedger)
			//if useLedger {
			//	return fmt.Errorf("--%s flag passed: Ledger device is currently not supported", flags.FlagUseLedger)
			//}

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			// TODO: define our own token
			customAppTemplate, customAppConfig := initAppConfig()

			return sdkserver.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, tmcfg.DefaultConfig())
		},
	}

	// TODO: double-check
	// authclient.Codec = encodingConfig.Marshaler

	cfg := sdk.GetConfig()
	cfg.Seal()

	rootCmd.AddCommand(
		ethermintclient.ValidateChainID(
			InitCmd(app.ModuleBasics, app.DefaultNodeHome),
		),
		SelfDelegationCmd(app.ModuleBasics, encodingConfig.TxConfig, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		NewTestnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		debug.Cmd(),
		config.Cmd(),
	)

	a := appCreator{encodingConfig}
	ethermintserver.AddCommands(
		rootCmd,
		ethermintserver.NewDefaultStartOptions(a.newApp, app.DefaultNodeHome),
		a.appExport,
		addModuleInitFlags,
	)

	// Add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		ethermintclient.KeyCommands(app.DefaultNodeHome),
	)
	rootCmd, err := srvflags.AddTxFlags(rootCmd)
	if err != nil {
		panic(err)
	}

	// Add rosetta
	rootCmd.AddCommand(sdkserver.RosettaCommand(encodingConfig.InterfaceRegistry, encodingConfig.Codec))

	return rootCmd, encodingConfig
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	app.ModuleBasics.AddQueryCommands(cmd)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
	)

	app.ModuleBasics.AddTxCommands(cmd)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	customAppTemplate, customAppConfig := servercfg.AppConfig(cmdcfg.BaseDenom)

	srvCfg, ok := customAppConfig.(servercfg.Config)
	if !ok {
		panic(fmt.Errorf("unknown app config type %T", customAppConfig))
	}

	// TODO: Is it necessary?
	// srvCfg.StateSync.SnapshotInterval = 1500
	// srvCfg.StateSync.SnapshotKeepRecent = 2

	return customAppTemplate, srvCfg
}

type appCreator struct {
	encCfg params.EncodingConfig
}

// newApp is an appCreator
func (a appCreator) newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(sdkserver.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(sdkserver.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := sdkserver.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	if cast.ToBool(appOpts.Get(cmdcfg.FlagPruningForValidatorOnly)) {

	}

	snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := sdk.NewLevelDB("metadata", snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	dscApp := app.NewDSC(
		logger,
		db,
		traceStore,
		true,
		skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(sdkserver.FlagInvCheckPeriod)),
		a.encCfg,
		appOpts,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(sdkserver.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(sdkserver.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(sdkserver.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(sdkserver.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(sdkserver.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(sdkserver.FlagIndexEvents))),
		baseapp.SetSnapshot(snapshotStore, snapshottypes.SnapshotOptions{
			Interval:   cast.ToUint64(appOpts.Get(sdkserver.FlagStateSyncSnapshotInterval)),
			KeepRecent: cast.ToUint32(appOpts.Get(sdkserver.FlagStateSyncSnapshotKeepRecent)),
		}),
	)

	return dscApp
}

// appExport creates a new simapp (optionally at a given height) and exports state.
func (a appCreator) appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
) (servertypes.ExportedApp, error) {
	var dscApp *app.DSC
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || len(homePath) == 0 {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	loadLatest := height < 0
	dscApp = app.NewDSC(logger, db, traceStore, loadLatest, map[int64]bool{}, "", uint(1), a.encCfg, appOpts)
	if !loadLatest {
		if err := dscApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return dscApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
}
