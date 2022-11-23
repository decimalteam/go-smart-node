package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	"github.com/spf13/viper"

	// Tendermint
	abci "github.com/tendermint/tendermint/abci/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmdb "github.com/tendermint/tm-db"

	// SDK
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"

	// SDK modules
	auth "github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	bank "github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capability "github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	feegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	genutil "github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	params "github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	// IBC
	ibc "github.com/cosmos/ibc-go/v5/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v5/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v5/modules/core/02-client/client"
	ibchost "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"

	// Ethermint
	ethencoding "github.com/evmos/ethermint/encoding"
	ethsrvflags "github.com/evmos/ethermint/server/flags"
	ethtypes "github.com/evmos/ethermint/types"

	// Ethermint modules
	evm "github.com/evmos/ethermint/x/evm"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	evmgeth "github.com/evmos/ethermint/x/evm/vm/geth"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	// Unnamed import of statik for swagger UI support
	// _ "bitbucket.org/decimalteam/go-smart-node/client/docs/statik"

	// Decimal
	ante "bitbucket.org/decimalteam/go-smart-node/app/ante"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"

	// Decimal modules
	coin "bitbucket.org/decimalteam/go-smart-node/x/coin"
	coinkeeper "bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	fee "bitbucket.org/decimalteam/go-smart-node/x/fee"
	feekeeper "bitbucket.org/decimalteam/go-smart-node/x/fee/keeper"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	legacy "bitbucket.org/decimalteam/go-smart-node/x/legacy"
	legacykeeper "bitbucket.org/decimalteam/go-smart-node/x/legacy/keeper"
	legacytypes "bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	multisig "bitbucket.org/decimalteam/go-smart-node/x/multisig"
	multisigkeeper "bitbucket.org/decimalteam/go-smart-node/x/multisig/keeper"
	multisigtypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nft "bitbucket.org/decimalteam/go-smart-node/x/nft"
	nftkeeper "bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	swap "bitbucket.org/decimalteam/go-smart-node/x/swap"
	swapkeeper "bitbucket.org/decimalteam/go-smart-node/x/swap/keeper"
	swaptypes "bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	upgrade "bitbucket.org/decimalteam/go-smart-node/x/upgrade"

	validator "bitbucket.org/decimalteam/go-smart-node/x/validator"
	validatorkeeper "bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// TODO: Move to some other place and use address with known private key!
const (
	UpgraderAddress = "dx1y7sex8yvrazyd8pljjxvnvpndaavn99tjd3ppm"
)

var (
	// DefaultNodeHome defines default home directory for the application daemon.
	DefaultNodeHome = os.ExpandEnv(fmt.Sprintf("$HOME/.%s/daemon", cmdcfg.AppName))
)

var (
	// MainnetChainIDPrefix defines the EVM EIP155 chain ID prefix for Decimal mainnet.
	MainnetChainIDPrefix = fmt.Sprintf("%s_%d", cmdcfg.AppName, cmdcfg.MainnetChainID)

	// TestnetChainIDPrefix defines the EVM EIP155 chain ID prefix for Decimal testnet.
	TestnetChainIDPrefix = fmt.Sprintf("%s_%d", cmdcfg.AppName, cmdcfg.TestnetChainID)
)

var (
	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration and genesis verification.
	ModuleBasics = module.NewBasicManager(
		// SDK
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		//staking.AppModuleBasic{},
		//distr.AppModuleBasic{},
		gov.NewAppModuleBasic([]govclient.ProposalHandler{
			// SDK proposal handlers
			paramsclient.ProposalHandler,
			distrclient.ProposalHandler,
			upgradeclient.LegacyProposalHandler,
			upgradeclient.LegacyCancelProposalHandler,
			// IBC proposal handlers
			ibcclientclient.UpdateClientProposalHandler,
			ibcclientclient.UpgradeProposalHandler,
			// Decimal proposal handlers
			// TODO: ...
		}),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		//slashing.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		//evidence.AppModuleBasic{},
		// IBC
		ibc.AppModuleBasic{},
		// Ethermint
		evm.AppModuleBasic{},
		// Decimal
		coin.AppModuleBasic{},
		legacy.AppModuleBasic{},
		nft.AppModuleBasic{},
		multisig.AppModuleBasic{},
		swap.AppModuleBasic{},
		fee.AppModuleBasic{},
		validator.AppModuleBasic{},
	)

	// Module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:       {authtypes.Burner, authtypes.Minter},
		distrtypes.ModuleName:            nil,
		validatortypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		validatortypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		validatortypes.ModuleName:        {authtypes.Burner, authtypes.Minter}, // used to store coins that will soon be paid out
		govtypes.ModuleName:              {authtypes.Burner},
		evmtypes.ModuleName:              {authtypes.Minter, authtypes.Burner}, // used for secure addition and subtraction of balance using module account
		cointypes.ModuleName:             {authtypes.Minter, authtypes.Burner},
		nfttypes.ReservedPool:            {authtypes.Minter, authtypes.Burner},
		legacytypes.LegacyCoinPool:       nil, // special account to hold legacy balances
		swaptypes.SwapPool:               nil, // special account to hold locked coins in swap process
		feetypes.BurningPool:             {authtypes.Burner},
	}

	// Module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName: true,
		cointypes.ModuleName:  true, // TODO: ?
		swaptypes.SwapPool:    true,
	}
)

var (
	_ servertypes.Application = (*DSC)(nil)
	_ simapp.App              = (*DSC)(nil)
	//_ ibctesting.TestingApp   = (*DSC)(nil)
)

// DSC implements an extended ABCI application. It is an application that may process
// transactions through Ethereum's EVM running atop of Tendermint consensus.
type DSC struct {
	*baseapp.BaseApp

	// Encoding
	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// Keys to access the substores
	keys    map[string]*store.KVStoreKey
	tkeys   map[string]*store.TransientStoreKey
	memKeys map[string]*store.MemoryStoreKey

	// SDK keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	//StakingKeeper    stakingkeeper.Keeper
	//SlashingKeeper   slashingkeeper.Keeper
	//DistrKeeper    distrkeeper.Keeper
	GovKeeper      govkeeper.Keeper
	CrisisKeeper   crisiskeeper.Keeper
	UpgradeKeeper  upgradekeeper.Keeper
	ParamsKeeper   paramskeeper.Keeper
	FeeGrantKeeper feegrantkeeper.Keeper
	AuthzKeeper    authzkeeper.Keeper
	//EvidenceKeeper   evidencekeeper.Keeper

	// IBC keepers
	IBCKeeper *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly

	// Make scoped keepers public for test purposes
	ScopedIBCKeeper capabilitykeeper.ScopedKeeper

	// Ethermint keepers
	EvmKeeper *evmkeeper.Keeper

	// Decimal keepers
	CoinKeeper      coinkeeper.Keeper
	SwapKeeper      swapkeeper.Keeper
	MultisigKeeper  multisigkeeper.Keeper
	NFTKeeper       nftkeeper.Keeper
	FeeKeeper       feekeeper.Keeper
	LegacyKeeper    legacykeeper.Keeper
	ValidatorKeeper validatorkeeper.Keeper

	// Module manager
	mm *module.Manager

	// Simulation manager
	sm *module.SimulationManager

	// Configurator
	configurator module.Configurator

	tpsCounter *tpsCounter

	// application options is *viper.Viper
	// need for proper telemetry initialization
	appOpts servertypes.AppOptions
}

// NewDSC returns a reference to a new initialized Decimal application.
func NewDSC(
	logger tmlog.Logger,
	db tmdb.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig simappparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *DSC {
	appCodec := encodingConfig.Codec
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	// Manually update the power reduction by replacing micro (u) -> atto (a) del
	sdk.DefaultPowerReduction = ethtypes.PowerReduction

	// NOTE we use custom transaction decoder that supports the sdk.Tx interface instead of sdk.StdTx
	bApp := baseapp.NewBaseApp(
		cmdcfg.AppBinName,
		logger,
		db,
		encodingConfig.TxConfig.TxDecoder(),
		baseAppOptions...,
	)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		// SDK keys
		authtypes.StoreKey,
		banktypes.StoreKey,
		//stakingtypes.StoreKey,
		//distrtypes.StoreKey,
		//slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		//evidencetypes.StoreKey,
		capabilitytypes.StoreKey,
		feegrant.StoreKey,
		authzkeeper.StoreKey,
		// IBC keys
		ibchost.StoreKey,
		// Ethermint keys
		evmtypes.StoreKey,
		// Decimal keys
		cointypes.StoreKey,
		multisigtypes.StoreKey,
		swaptypes.StoreKey,
		nfttypes.StoreKey,
		feetypes.StoreKey,
		legacytypes.StoreKey,
		validatortypes.StoreKey,
		upgradetypes.StoreKey,
	)

	// Add the EVM transient store key
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey, feemarkettypes.TransientKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	// Create app instance
	app := &DSC{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
		appOpts:           appOpts,
	}

	// Init params keeper and subspaces
	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()))

	// Add capability keeper and ScopeToModule for IBC module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	// Grant capabilities for the IBC and IBC-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)

	// Applications that wish to enforce statically created ScopedKeepers should call `Seal`
	// after creating their scoped modules in `NewApp` with `ScopeToModule`
	app.CapabilityKeeper.Seal()

	// Create SDK keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		ethtypes.ProtoAccount,
		maccPerms,
		cmdcfg.Bech32Prefix,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.GetSubspace(banktypes.ModuleName),
		app.BlockedAddrs(),
	)
	// app.StakingKeeper = stakingkeeper.NewKeeper(
	// 	appCodec,
	// 	keys[stakingtypes.StoreKey],
	// 	app.AccountKeeper,
	// 	app.BankKeeper,
	// 	app.GetSubspace(stakingtypes.ModuleName),
	// )
	// app.DistrKeeper = distrkeeper.NewKeeper(
	// 	appCodec,
	// 	keys[distrtypes.StoreKey],
	// 	app.GetSubspace(distrtypes.ModuleName),
	// 	app.AccountKeeper,
	// 	app.BankKeeper,
	// 	&app.StakingKeeper,
	// 	authtypes.FeeCollectorName,
	// )
	// app.SlashingKeeper = slashingkeeper.NewKeeper(
	// 	appCodec,
	// 	keys[slashingtypes.StoreKey],
	// 	&app.StakingKeeper,
	// 	app.GetSubspace(slashingtypes.ModuleName),
	// )
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)
	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)
	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		app.BaseApp.MsgServiceRouter(),
		app.AccountKeeper,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	// app.StakingKeeper.SetHooks(
	// 	stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	// )

	// Create decimal keeper because Validator = Staking+Slashing+Evidence
	app.CoinKeeper = *coinkeeper.NewKeeper(
		appCodec,
		keys[cointypes.StoreKey],
		app.GetSubspace(cointypes.ModuleName),
		app.AccountKeeper,
		&app.FeeKeeper,
		app.BankKeeper,
	)
	app.NFTKeeper = *nftkeeper.NewKeeper(
		appCodec,
		keys[nfttypes.StoreKey],
		app.GetSubspace(nfttypes.ModuleName),
		app.BankKeeper,
	)
	app.MultisigKeeper = *multisigkeeper.NewKeeper(
		appCodec,
		keys[multisigtypes.StoreKey],
		app.GetSubspace(multisigtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.MsgServiceRouter(),
	)
	app.FeeKeeper = *feekeeper.NewKeeper(
		appCodec,
		keys[feetypes.StoreKey],
		app.GetSubspace(feetypes.ModuleName),
		app.BankKeeper,
		&app.CoinKeeper,
		app.AccountKeeper,
		cmdcfg.BaseDenom,
		ante.CalculateFee,
	)
	app.ValidatorKeeper = validatorkeeper.NewKeeper(
		appCodec,
		keys[validatortypes.StoreKey],
		app.GetSubspace(validatortypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&app.NFTKeeper,
		&app.CoinKeeper,
		&app.MultisigKeeper,
	)

	// Create Ethermint keepers

	app.EvmKeeper = evmkeeper.NewKeeper(
		appCodec,
		keys[evmtypes.StoreKey],
		tkeys[evmtypes.TransientKey],
		app.GetSubspace(evmtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&app.ValidatorKeeper,
		app.FeeKeeper,
		nil,
		evmgeth.NewEVM,
		cast.ToString(appOpts.Get(ethsrvflags.EVMTracer)),
	)
	app.EvmKeeper = app.EvmKeeper.SetHooks(
		evmkeeper.NewMultiEvmHooks(),
	)

	// Create upgrade keeper
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
		UpgraderAddress,
	)

	// Create IBC keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibchost.StoreKey],
		app.GetSubspace(ibchost.ModuleName),
		app.ValidatorKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	// Register the proposal types
	govLegacyRouter := govtypesv1beta1.NewRouter()
	govLegacyRouter.AddRoute(govtypes.RouterKey, govtypesv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))
	govConfig := govtypes.DefaultConfig()

	app.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.GetSubspace(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.ValidatorKeeper,
		govLegacyRouter,
		app.BaseApp.MsgServiceRouter(),
		govConfig,
	)

	// If evidence needs to be handled for the app, set routes in router here and seal

	// Create Decimal keepers
	app.LegacyKeeper = *legacykeeper.NewKeeper(
		appCodec,
		keys[legacytypes.StoreKey],
		app.BankKeeper,
		&app.NFTKeeper,
		&app.MultisigKeeper,
		app.ValidatorKeeper,
	)
	app.SwapKeeper = *swapkeeper.NewKeeper(
		appCodec,
		keys[swaptypes.StoreKey],
		app.GetSubspace(swaptypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
	)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		// SDK app modules
		genutil.NewAppModule(app.AccountKeeper, app.ValidatorKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		//slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		//distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		//staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(appCodec, app.UpgradeKeeper),
		//evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		// IBC app modules
		ibc.NewAppModule(app.IBCKeeper),
		// Ethermint app modules
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper),

		// Decimal app modules
		coin.NewAppModule(appCodec, app.CoinKeeper, app.AccountKeeper, app.BankKeeper),
		multisig.NewAppModule(appCodec, app.MultisigKeeper, app.AccountKeeper, app.BankKeeper),
		swap.NewAppModule(appCodec, app.SwapKeeper, app.AccountKeeper, app.BankKeeper),
		nft.NewAppModule(appCodec, app.NFTKeeper),
		fee.NewAppModule(appCodec, app.FeeKeeper),
		legacy.NewAppModule(app.appCodec, app.LegacyKeeper),
		validator.NewAppModule(appCodec, app.ValidatorKeeper, app.AccountKeeper, app.BankKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: upgrade module must go first to handle software upgrades.
	// NOTE: staking module is required if HistoricalEntries param > 0.
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		evmtypes.ModuleName,
		validatortypes.ModuleName,
		//distrtypes.ModuleName,
		//slashingtypes.ModuleName,
		//evidencetypes.ModuleName,
		//stakingtypes.ModuleName,
		ibchost.ModuleName,
		// no-op modules
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		cointypes.ModuleName,
		multisigtypes.ModuleName,
		swaptypes.ModuleName,
		nfttypes.ModuleName,
		feetypes.ModuleName,
		legacytypes.ModuleName,
	)

	// NOTE: fee market module must go last in order to retrieve the block gas used.
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		validatortypes.ModuleName,
		//stakingtypes.ModuleName,
		evmtypes.ModuleName,
		//claimstypes.ModuleName,
		// no-op modules
		ibchost.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		//distrtypes.ModuleName,
		//slashingtypes.ModuleName,
		genutiltypes.ModuleName,
		//evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		cointypes.ModuleName,
		multisigtypes.ModuleName,
		swaptypes.ModuleName,
		nfttypes.ModuleName,
		legacytypes.ModuleName,
		feetypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		// SDK modules
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		//distrtypes.ModuleName,
		//stakingtypes.ModuleName,
		//slashingtypes.ModuleName,
		govtypes.ModuleName,
		ibchost.ModuleName,
		//evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		// fee decimal replacer
		feetypes.ModuleName,
		// Ethermint modules
		evmtypes.ModuleName,
		//recoverytypes.ModuleName,
		// Decimal modules
		cointypes.ModuleName,
		multisigtypes.ModuleName,
		swaptypes.ModuleName,
		nfttypes.ModuleName,
		legacytypes.ModuleName,
		validatortypes.ModuleName,
		genutiltypes.ModuleName,
		// NOTE: crisis module must go at the end to check for invariants on each module
		crisistypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	// NOTE: this is not required apps that don't use the simulator for fuzz testing transactions
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		//staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		//distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.ValidatorKeeper),
		//slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		params.NewAppModule(app.ParamsKeeper),
		//evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper),
		coin.NewAppModule(appCodec, app.CoinKeeper, app.AccountKeeper, app.BankKeeper),
		swap.NewAppModule(appCodec, app.SwapKeeper, app.AccountKeeper, app.BankKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// Initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// Initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	maxGasWanted := cast.ToUint64(appOpts.Get(ethsrvflags.EVMMaxTxGasWanted))
	options := ante.HandlerOptions{
		Cdc:             appCodec,
		AccountKeeper:   app.AccountKeeper,
		BankKeeper:      app.BankKeeper,
		EvmKeeper:       app.EvmKeeper,
		FeeMarketKeeper: app.FeeKeeper,
		FeegrantKeeper:  app.FeeGrantKeeper,
		IBCKeeper:       app.IBCKeeper,
		CoinKeeper:      &app.CoinKeeper,
		FeeKeeper:       &app.FeeKeeper,
		LegacyKeeper:    &app.LegacyKeeper,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:  SigVerificationGasConsumer,
		MaxTxGasWanted:  maxGasWanted,
		// TODO: ethermint stopped to use this option and checker
		//ExtensionOptionChecker: ethtypes.HasDynamicFeeExtensionOption,
		//TxFeeChecker:           ethante.NewDynamicFeeChecker(app.EvmKeeper),
	}

	if err := options.Validate(); err != nil {
		panic(err)
	}

	app.SetAnteHandler(ante.NewAnteHandler(options))
	app.SetEndBlocker(app.EndBlocker)
	app.setupUpgradeHandlers()

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper

	// Finally start the tpsCounter.
	app.tpsCounter = newTPSCounter(logger)
	go func() {
		// Unfortunately golangci-lint is so pedantic
		// so we have to ignore this error explicitly.
		_ = app.tpsCounter.start(context.Background())
	}()

	return app
}

// Name returns the name of the App.
func (app *DSC) Name() string { return app.BaseApp.Name() }

// BeginBlocker updates every begin block.
func (app *DSC) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker updates every end block.
func (app *DSC) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// We are intentionally decomposing the DeliverTx method so as to calculate the transactions per second.
func (app *DSC) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	defer func() {
		// TODO: Record the count along with the code and or reason so as to display
		// in the transactions per second live dashboards.
		if res.IsErr() {
			app.tpsCounter.incrementFailure()
		} else {
			app.tpsCounter.incrementSuccess()
		}
	}()

	return app.BaseApp.DeliverTx(req)
}

// InitChainer updates at chain initialization.
func (app *DSC) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var gs simapp.GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &gs); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.appCodec, gs)
}

// LoadHeight loads state at a particular height.
func (app *DSC) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *DSC) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *DSC) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// LegacyAmino returns DSC's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *DSC) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns DSC's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *DSC) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns DSC's InterfaceRegistry
func (app *DSC) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *DSC) GetKey(storeKey string) *store.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *DSC) GetTKey(storeKey string) *store.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *DSC) GetMemKey(storeKey string) *store.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *DSC) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *DSC) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *DSC) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}

	// TODO: remove after telemetry support in ethermint
	// this is workaround for ethermint startInProcess
	// ethermint.startInProcess does not initialize telemetry
	v, ok := app.appOpts.(*viper.Viper)
	if !ok {
		app.Logger().Error("can't convert appOpts to viper")
		return
	}
	cfg, err := config.ParseConfig(v)
	if err != nil {
		app.Logger().Error("can't parse config: ", err)
		return
	}
	if cfg.API.Enable && cfg.Telemetry.Enabled {
		metrics, err := telemetry.New(cfg.Telemetry)
		if err != nil {
			app.Logger().Error("can't create telemetry: ", err)
			return
		}
		apiSvr.SetTelemetry(metrics)
	}

}

func (app *DSC) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *DSC) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}

// IBC Go TestingApp functions

// GetBaseApp implements the TestingApp interface.
func (app *DSC) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetStakingKeeper implements the TestingApp interface.
// TODO: fix it?
func (app *DSC) GetStakingKeeper() validatorkeeper.Keeper {
	return app.ValidatorKeeper
}

// GetIBCKeeper implements the TestingApp interface.
func (app *DSC) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (app *DSC) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

// GetTxConfig implements the TestingApp interface.
func (app *DSC) GetTxConfig() client.TxConfig {
	cfg := ethencoding.MakeConfig(ModuleBasics)
	return cfg.TxConfig
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(_ client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}

// initParamsKeeper inits params keeper and its subspaces.
func initParamsKeeper(
	appCodec codec.BinaryCodec,
	legacyAmino *codec.LegacyAmino,
	key store.StoreKey,
	tkey store.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)
	// SDK subspaces
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	//paramsKeeper.Subspace(stakingtypes.ModuleName)
	//paramsKeeper.Subspace(distrtypes.ModuleName)
	//paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypesv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	// IBC subspaces
	paramsKeeper.Subspace(ibchost.ModuleName)
	// Ethermint subspaces
	paramsKeeper.Subspace(evmtypes.ModuleName)
	//paramsKeeper.Subspace(recoverytypes.ModuleName)
	// Decimal subspaces
	paramsKeeper.Subspace(cointypes.ModuleName)
	paramsKeeper.Subspace(feetypes.ModuleName)
	paramsKeeper.Subspace(multisigtypes.ModuleName)
	paramsKeeper.Subspace(nfttypes.ModuleName)
	paramsKeeper.Subspace(swaptypes.ModuleName)
	paramsKeeper.Subspace(validatortypes.ModuleName)
	return paramsKeeper
}

func (app *DSC) setupUpgradeHandlers() {
	// When a planned update height is reached, the old binary will panic
	// writing on disk the height and name of the update that triggered it
	// This will read that value, and execute the preparations for the upgrade.
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Errorf("failed to read upgrade info from disk: %w", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	var storeUpgrades *store.StoreUpgrades

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
