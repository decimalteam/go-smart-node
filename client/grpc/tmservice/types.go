package tmservice

import (
	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (b *BlockResults) ToCoreTypes() *ctypes.ResultBlockResults {
	beginBlockEvents := make([]abci.Event, len(b.BeginBlockEvents))
	for i, v := range b.BeginBlockEvents {
		beginBlockEvents[i] = *v
	}
	endBlockEvents := make([]abci.Event, len(b.EndBlockEvents))
	for i, v := range b.EndBlockEvents {
		endBlockEvents[i] = *v
	}
	validatorUpdates := make([]abci.ValidatorUpdate, len(b.ValidatorUpdates))
	for i, v := range validatorUpdates {
		validatorUpdates[i] = v
	}
	return &ctypes.ResultBlockResults{
		Height:                b.Height,
		TxsResults:            b.TxsResults,
		BeginBlockEvents:      beginBlockEvents,
		EndBlockEvents:        endBlockEvents,
		ValidatorUpdates:      validatorUpdates,
		ConsensusParamUpdates: b.ConsensusParamUpdates,
	}
}
