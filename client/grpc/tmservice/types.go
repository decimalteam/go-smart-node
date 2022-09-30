package tmservice

import (
	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// ToABCIRequestQuery converts a gRPC ABCIQueryRequest type to an ABCI
// RequestQuery type.
func (req *ABCIQueryRequest) ToABCIRequestQuery() abci.RequestQuery {
	return abci.RequestQuery{
		Data:   req.Data,
		Path:   req.Path,
		Height: req.Height,
		Prove:  req.Prove,
	}
}

// FromABCIResponseQuery converts an ABCI ResponseQuery type to a gRPC
// ABCIQueryResponse type.
func FromABCIResponseQuery(res abci.ResponseQuery) *ABCIQueryResponse {
	var proofOps *ProofOps

	if res.ProofOps != nil {
		proofOps = &ProofOps{
			Ops: make([]ProofOp, len(res.ProofOps.Ops)),
		}
		for i, proofOp := range res.ProofOps.Ops {
			proofOps.Ops[i] = ProofOp{
				Type: proofOp.Type,
				Key:  proofOp.Key,
				Data: proofOp.Data,
			}
		}
	}

	return &ABCIQueryResponse{
		Code:      res.Code,
		Log:       res.Log,
		Info:      res.Info,
		Index:     res.Index,
		Key:       res.Key,
		Value:     res.Value,
		ProofOps:  proofOps,
		Height:    res.Height,
		Codespace: res.Codespace,
	}
}

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
