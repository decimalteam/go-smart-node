package fee

import (
	"context"
	"encoding/json"
	"fmt"

	"bitbucket.org/decimalteam/go-smart-node/x/fee/client/cli"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
	_ module.AppModule      = AppModule{}
)

// AppModuleBasic app module basics object
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {
}

func (AppModuleBasic) RegisterInterfaces(interfaceRegistry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(interfaceRegistry)
}

// DefaultGenesis default genesis state
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the module.
func (b AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genesisState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genesisState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genesisState.Validate()
}

// RegisterRESTRoutes registers rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	//rest.RegisterRoutes(ctx, rtr, types.RouterKey)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (b AppModuleBasic) RegisterGRPCGatewayRoutes(c client.Context, serveMux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(c)); err != nil {
		panic(err)
	}
}

// GetTxCmd gets the root tx command of this module
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd gets the root query command of this module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// AppModule supply app module
type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// Name defines module name
func (AppModule) Name() string {
	return types.ModuleName
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(&am.keeper)
}

// RegisterInvariants registers the fee module invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {

}

// Route returns the module's message routing key.
func (am AppModule) Route() sdk.Route {
	return sdk.Route{}
}

// QuerierRoute module querier route name
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// LegacyQuerierHandler returns the module's Querier.
func (am AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return nil
}

// RegisterServices registers a GRPC query service to respond to the module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
	feemarkettypes.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// ConsensusVersion implements ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 {
	return 1
}

// InitGenesis performs the module's genesis initialization.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var gs types.GenesisState
	cdc.MustUnmarshalJSON(data, &gs)
	InitGenesis(ctx, am.keeper, gs)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock module begin-block
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
}

// EndBlock module end-block
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
