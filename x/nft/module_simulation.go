package nft

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

var _ module.AppModuleSimulation = AppModule{}

// RegisterStoreDecoder registers a decoder for nft module's types
func (AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {
	//sdr[types.StoreKey] = simulation.DecodeStore
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// GenerateGenesisState creates a randomized GenState of the nft module.
func (AppModule) GenerateGenesisState(_ *module.SimulationState) {
	//simulation.RandomizedGenState(simState)
}

// RandomizedParams doesn't create randomized nft param changes for the simulator.
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {
	return nil
}

// WeightedOperations doesn't return any operation for the nft module.
func (am AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}
