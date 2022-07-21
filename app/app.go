package app

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	// Tendermint
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	// SDK
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
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
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidence "github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	feegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	genutil "github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	params "github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	// IBC
	ibctransfer "github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v3/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v3/modules/core/02-client/client"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"

	// Ethermint
	"github.com/tharsis/ethermint/encoding"
	srvflags "github.com/tharsis/ethermint/server/flags"
	ethermint "github.com/tharsis/ethermint/types"

	// Ethermint modules
	evm "github.com/tharsis/ethermint/x/evm"
	evmrest "github.com/tharsis/ethermint/x/evm/client/rest"
	evmkeeper "github.com/tharsis/ethermint/x/evm/keeper"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	feemarket "github.com/tharsis/ethermint/x/feemarket"
	feemarketkeeper "github.com/tharsis/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"

	// Evmos modules
	claims "github.com/tharsis/evmos/v3/x/claims"
	claimskeeper "github.com/tharsis/evmos/v3/x/claims/keeper"
	claimstypes "github.com/tharsis/evmos/v3/x/claims/types"
	epochs "github.com/tharsis/evmos/v3/x/epochs"
	epochskeeper "github.com/tharsis/evmos/v3/x/epochs/keeper"
	epochstypes "github.com/tharsis/evmos/v3/x/epochs/types"
	erc20 "github.com/tharsis/evmos/v3/x/erc20"
	erc20client "github.com/tharsis/evmos/v3/x/erc20/client"
	erc20keeper "github.com/tharsis/evmos/v3/x/erc20/keeper"
	erc20types "github.com/tharsis/evmos/v3/x/erc20/types"
	incentives "github.com/tharsis/evmos/v3/x/incentives"
	incentivesclient "github.com/tharsis/evmos/v3/x/incentives/client"
	incentiveskeeper "github.com/tharsis/evmos/v3/x/incentives/keeper"
	incentivestypes "github.com/tharsis/evmos/v3/x/incentives/types"
	inflation "github.com/tharsis/evmos/v3/x/inflation"
	inflationkeeper "github.com/tharsis/evmos/v3/x/inflation/keeper"
	inflationtypes "github.com/tharsis/evmos/v3/x/inflation/types"
	recovery "github.com/tharsis/evmos/v3/x/recovery"
	recoverykeeper "github.com/tharsis/evmos/v3/x/recovery/keeper"
	recoverytypes "github.com/tharsis/evmos/v3/x/recovery/types"
	vesting "github.com/tharsis/evmos/v3/x/vesting"
	vestingkeeper "github.com/tharsis/evmos/v3/x/vesting/keeper"
	vestingtypes "github.com/tharsis/evmos/v3/x/vesting/types"

	// Unnamed import of statik for swagger UI spport
	// _ "bitbucket.org/decimalteam/go-smart-node/client/docs/statik"

	// Decimal
	"bitbucket.org/decimalteam/go-smart-node/app/ante"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"

	// Decimal modules
	coin "bitbucket.org/decimalteam/go-smart-node/x/coin"
	coinkeeper "bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	nftkeeper "bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
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
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			// SDK proposal handlers
			paramsclient.ProposalHandler,
			distrclient.ProposalHandler,
			upgradeclient.ProposalHandler,
			upgradeclient.CancelProposalHandler,
			// IBC proposal handlers
			ibcclientclient.UpdateClientProposalHandler,
			ibcclientclient.UpgradeProposalHandler,
			// Evmos proposal handlers
			erc20client.RegisterCoinProposalHandler,
			erc20client.RegisterERC20ProposalHandler,
			erc20client.ToggleTokenRelayProposalHandler,
			erc20client.UpdateTokenPairERC20ProposalHandler,
			incentivesclient.RegisterIncentiveProposalHandler,
			incentivesclient.CancelIncentiveProposalHandler,
			// Decimal proposal handlers
			// TODO: ?
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		// IBC
		ibc.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		// Ethermint
		evm.AppModuleBasic{},
		feemarket.AppModuleBasic{},
		// Evmos
		vesting.AppModuleBasic{},
		inflation.AppModuleBasic{},
		erc20.AppModuleBasic{},
		incentives.AppModuleBasic{},
		epochs.AppModuleBasic{},
		claims.AppModuleBasic{},
		recovery.AppModuleBasic{},
		// Decimal
		coin.AppModuleBasic{},
		nft.AppModuleBasic{},
	)

	// Module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     {authtypes.Burner},
		distrtypes.ModuleName:          nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		evmtypes.ModuleName:            {authtypes.Minter, authtypes.Burner}, // used for secure addition and subtraction of balance using module account
		inflationtypes.ModuleName:      {authtypes.Minter},
		erc20types.ModuleName:          {authtypes.Minter, authtypes.Burner},
		claimstypes.ModuleName:         nil,
		incentivestypes.ModuleName:     {authtypes.Minter, authtypes.Burner},
		cointypes.ModuleName:           {authtypes.Minter, authtypes.Burner},
		nfttypes.ReservedPool:          {authtypes.Minter, authtypes.Burner},
		// special account to hold legacy balances
		cointypes.LegacyCoinPool: nil,
	}

	// Module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName:      true,
		incentivestypes.ModuleName: true,
		cointypes.ModuleName:       true, // TODO: ?
	}
)

var (
	_ servertypes.Application = (*DSC)(nil)
	_ simapp.App              = (*DSC)(nil)
	_ ibctesting.TestingApp   = (*DSC)(nil)
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
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// SDK keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper
	EvidenceKeeper   evidencekeeper.Keeper

	// IBC keepers
	IBCKeeper         *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCTransferKeeper ibctransferkeeper.Keeper

	// Make scoped keepers public for test purposes
	ScopedIBCKeeper         capabilitykeeper.ScopedKeeper
	ScopedIBCTransferKeeper capabilitykeeper.ScopedKeeper

	// Ethermint keepers
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper

	// Evmos keepers
	InflationKeeper  inflationkeeper.Keeper
	ClaimsKeeper     *claimskeeper.Keeper
	Erc20Keeper      erc20keeper.Keeper
	IncentivesKeeper incentiveskeeper.Keeper
	EpochsKeeper     epochskeeper.Keeper
	VestingKeeper    vestingkeeper.Keeper
	RecoveryKeeper   *recoverykeeper.Keeper

	// Decimal keepers
	CoinKeeper coinkeeper.Keeper
	NFTKeeper  nftkeeper.Keeper

	// Module manager
	mm *module.Manager

	// Simulation manager
	sm *module.SimulationManager

	// Configurator
	configurator module.Configurator

	tpsCounter *tpsCounter
}

// NewDSC returns a reference to a new initialized Ethermint application.
func NewDSC(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig simappparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *DSC {
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	// Manually update the power reduction by replacing micro (u) -> atto (a) del
	sdk.DefaultPowerReduction = ethermint.PowerReduction

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
		stakingtypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		capabilitytypes.StoreKey,
		feegrant.StoreKey,
		authzkeeper.StoreKey,
		// IBC keys
		ibchost.StoreKey,
		ibctransfertypes.StoreKey,
		// Ethermint keys
		evmtypes.StoreKey,
		feemarkettypes.StoreKey,
		// Evmos keys
		inflationtypes.StoreKey,
		erc20types.StoreKey,
		incentivestypes.StoreKey,
		epochstypes.StoreKey,
		claimstypes.StoreKey,
		vestingtypes.StoreKey,
		// Decimal keys
		cointypes.StoreKey,
		nfttypes.StoreKey,
	)

	// Add the EVM transient store key
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey)
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
	}

	// Init params keeper and subspaces
	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// Add capability keeper and ScopeToModule for IBC module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	// Grant capabilities for the IBC and IBC-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedIBCTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	// Applications that wish to enforce statically created ScopedKeepers should call `Seal`
	// after creating their scoped modules in `NewApp` with `ScopeToModule`
	app.CapabilityKeeper.Seal()

	// Create SDK keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		ethermint.ProtoAccount,
		maccPerms,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.GetSubspace(banktypes.ModuleName),
		app.BlockedAddrs(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(stakingtypes.ModuleName),
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		app.GetSubspace(distrtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		authtypes.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		keys[slashingtypes.StoreKey],
		&stakingKeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)
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
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
	)
	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		app.BaseApp.MsgServiceRouter(),
	)

	// Create Ethermint keepers
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec,
		keys[feemarkettypes.StoreKey],
		app.GetSubspace(feemarkettypes.ModuleName),
	)
	app.EvmKeeper = evmkeeper.NewKeeper(
		appCodec,
		keys[evmtypes.StoreKey],
		tkeys[evmtypes.TransientKey],
		app.GetSubspace(evmtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		app.FeeMarketKeeper,
		cast.ToString(appOpts.Get(srvflags.EVMTracer)),
	)

	// Create IBC keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), &stakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	)

	// Register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper)).
		AddRoute(erc20types.RouterKey, erc20.NewErc20ProposalHandler(&app.Erc20Keeper)).
		AddRoute(incentivestypes.RouterKey, incentives.NewIncentivesProposalHandler(&app.IncentivesKeeper))
	govKeeper := govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.GetSubspace(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		govRouter,
	)

	// Create Evmos keepers
	app.InflationKeeper = inflationkeeper.NewKeeper(
		keys[inflationtypes.StoreKey],
		appCodec,
		app.GetSubspace(inflationtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		&stakingKeeper,
		authtypes.FeeCollectorName,
	)
	app.ClaimsKeeper = claimskeeper.NewKeeper(
		appCodec,
		keys[claimstypes.StoreKey],
		app.GetSubspace(claimstypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		app.DistrKeeper,
	)
	// Register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	// NOTE: Distr, Slashing and Claim must be created before calling the Hooks method to avoid returning a Keeper without its table generated
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks(),
			app.ClaimsKeeper.Hooks(),
		),
	)
	app.VestingKeeper = vestingkeeper.NewKeeper(
		keys[vestingtypes.StoreKey],
		appCodec,
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
	)
	app.Erc20Keeper = erc20keeper.NewKeeper(
		keys[erc20types.StoreKey],
		appCodec,
		app.GetSubspace(erc20types.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.EvmKeeper,
	)
	app.IncentivesKeeper = incentiveskeeper.NewKeeper(
		keys[incentivestypes.StoreKey],
		appCodec,
		app.GetSubspace(incentivestypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.InflationKeeper,
		app.StakingKeeper,
		app.EvmKeeper,
	)
	epochsKeeper := epochskeeper.NewKeeper(
		appCodec,
		keys[epochstypes.StoreKey],
	)
	app.EpochsKeeper = *epochsKeeper.SetHooks(
		epochskeeper.NewMultiEpochHooks(
			// Insert epoch hooks receivers here
			app.IncentivesKeeper.Hooks(),
			app.InflationKeeper.Hooks(),
		),
	)
	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
			app.ClaimsKeeper.Hooks(),
		),
	)
	app.EvmKeeper = app.EvmKeeper.SetHooks(
		evmkeeper.NewMultiEvmHooks(
			app.Erc20Keeper.Hooks(),
			app.IncentivesKeeper.Hooks(),
			app.ClaimsKeeper.Hooks(),
		),
	)

	// Create Transfer Stack

	// SendPacket, since it is originating from the application to core IBC:
	// ibctransferkeeper.SendPacket -> claim.SendPacket -> recovery.SendPacket -> channel.SendPacket

	// RecvPacket, message that originates from core IBC and goes down to app, the flow is the otherway
	// channel.RecvPacket -> recovery.OnRecvPacket -> claim.OnRecvPacket -> ibctransfer.OnRecvPacket

	app.IBCTransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.ClaimsKeeper, // ICS4 Wrapper: claims IBC middleware
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedIBCTransferKeeper,
	)
	app.RecoveryKeeper = recoverykeeper.NewKeeper(
		app.GetSubspace(recoverytypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.IBCTransferKeeper,
		app.ClaimsKeeper,
	)

	// Set the ICS4 wrappers for claims and recovery middlewares
	app.RecoveryKeeper.SetICS4Wrapper(app.IBCKeeper.ChannelKeeper)
	app.ClaimsKeeper.SetICS4Wrapper(app.RecoveryKeeper)
	// NOTE: ICS4 wrapper for Transfer Keeper already set

	ibctransferModule := ibctransfer.NewAppModule(app.IBCTransferKeeper)

	// Transfer stack contains (from top to bottom):
	// - Recovery Middleware
	// - Airdrop Claims Middleware
	// - Transfer

	// Create IBC module from bottom to top of stack
	var ibctransferStack porttypes.IBCModule

	ibctransferStack = ibctransfer.NewIBCModule(app.IBCTransferKeeper)
	ibctransferStack = claims.NewIBCMiddleware(*app.ClaimsKeeper, ibctransferStack)
	ibctransferStack = recovery.NewIBCMiddleware(*app.RecoveryKeeper, ibctransferStack)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, ibctransferStack)
	app.IBCKeeper.SetRouter(ibcRouter)

	// Create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		keys[evidencetypes.StoreKey],
		&app.StakingKeeper,
		app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	// Create Decimal keepers
	coinKeeper := coinkeeper.NewKeeper(
		appCodec,
		keys[cointypes.StoreKey],
		app.GetSubspace(cointypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
	)
	app.CoinKeeper = *coinKeeper
	nftKeeper := nftkeeper.NewKeeper(
		appCodec,
		keys[nfttypes.StoreKey],
		app.BankKeeper,
		cmdcfg.BaseDenom,
	)
	app.NFTKeeper = *nftKeeper

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		// SDK app modules
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		// IBC app modules
		ibc.NewAppModule(app.IBCKeeper),
		ibctransferModule,
		// Ethermint app modules
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		// Evmos app modules
		inflation.NewAppModule(app.InflationKeeper, app.AccountKeeper, app.StakingKeeper),
		erc20.NewAppModule(app.Erc20Keeper, app.AccountKeeper),
		incentives.NewAppModule(app.IncentivesKeeper, app.AccountKeeper),
		epochs.NewAppModule(appCodec, app.EpochsKeeper),
		claims.NewAppModule(appCodec, *app.ClaimsKeeper),
		vesting.NewAppModule(app.VestingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		recovery.NewAppModule(*app.RecoveryKeeper),
		// Decimal app modules
		coin.NewAppModule(appCodec, app.CoinKeeper, app.AccountKeeper, app.BankKeeper),
		nft.NewAppModule(app.NFTKeeper, app.AccountKeeper),
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
		// Note: epochs' begin should be "real" start of epochs, we keep epochs beginblock at the beginning
		epochstypes.ModuleName,
		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibchost.ModuleName,
		// no-op modules
		ibctransfertypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		inflationtypes.ModuleName,
		erc20types.ModuleName,
		claimstypes.ModuleName,
		incentivestypes.ModuleName,
		recoverytypes.ModuleName,
		cointypes.ModuleName,
		nfttypes.ModuleName,
	)

	// NOTE: fee market module must go last in order to retrieve the block gas used.
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		// Note: epochs' endblock should be "real" end of epochs, we keep epochs endblock at the end
		epochstypes.ModuleName,
		claimstypes.ModuleName,
		// no-op modules
		ibchost.ModuleName,
		ibctransfertypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		inflationtypes.ModuleName,
		erc20types.ModuleName,
		incentivestypes.ModuleName,
		recoverytypes.ModuleName,
		cointypes.ModuleName,
		nfttypes.ModuleName,
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
		distrtypes.ModuleName,
		// NOTE: staking requires the claiming hook
		claimstypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		ibchost.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		// Ethermint modules
		evmtypes.ModuleName, feemarkettypes.ModuleName,
		// Evmos modules
		vestingtypes.ModuleName,
		inflationtypes.ModuleName,
		erc20types.ModuleName,
		incentivestypes.ModuleName,
		epochstypes.ModuleName,
		recoverytypes.ModuleName,
		// Decimal modules
		cointypes.ModuleName,
		nfttypes.ModuleName,
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
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		params.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		ibctransferModule,
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper),
		epochs.NewAppModule(appCodec, app.EpochsKeeper),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		coin.NewAppModule(appCodec, app.CoinKeeper, app.AccountKeeper, app.BankKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// Initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// Initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	maxGasWanted := cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted))
	options := ante.HandlerOptions{
		AccountKeeper:   app.AccountKeeper,
		BankKeeper:      app.BankKeeper,
		EvmKeeper:       app.EvmKeeper,
		FeegrantKeeper:  app.FeeGrantKeeper,
		IBCKeeper:       app.IBCKeeper,
		FeeMarketKeeper: app.FeeMarketKeeper,
		CoinKeeper:      &app.CoinKeeper,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:  SigVerificationGasConsumer,
		Cdc:             appCodec,
		MaxTxGasWanted:  maxGasWanted,
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
	app.ScopedIBCTransferKeeper = scopedIBCTransferKeeper

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
func (app *DSC) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *DSC) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *DSC) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
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
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)

	evmrest.RegisterTxRoutes(clientCtx, apiSvr.Router)

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

func (app *DSC) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *DSC) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// IBC Go TestingApp functions

// GetBaseApp implements the TestingApp interface.
func (app *DSC) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetStakingKeeper implements the TestingApp interface.
func (app *DSC) GetStakingKeeper() stakingkeeper.Keeper {
	return app.StakingKeeper
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
	cfg := encoding.MakeConfig(ModuleBasics)
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
	key sdk.StoreKey,
	tkey sdk.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)
	// SDK subspaces
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	// Ethermint subspaces
	paramsKeeper.Subspace(evmtypes.ModuleName)
	paramsKeeper.Subspace(feemarkettypes.ModuleName)
	// Evmos subspaces
	paramsKeeper.Subspace(inflationtypes.ModuleName)
	paramsKeeper.Subspace(erc20types.ModuleName)
	paramsKeeper.Subspace(claimstypes.ModuleName)
	paramsKeeper.Subspace(incentivestypes.ModuleName)
	paramsKeeper.Subspace(recoverytypes.ModuleName)
	// Decimal subspaces
	paramsKeeper.Subspace(cointypes.ModuleName)
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

	var storeUpgrades *storetypes.StoreUpgrades

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
