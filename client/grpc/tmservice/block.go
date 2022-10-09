package tmservice

import (
	"context"

	abci "github.com/tendermint/tendermint/abci/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/client"
)

func getBlockResults(ctx context.Context, clientCtx client.Context, height int64) (*BlockResults, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}
	res, err := node.BlockResults(ctx, &height)

	beginBlockEvents := make([]*abci.Event, len(res.BeginBlockEvents))
	for i, v := range res.BeginBlockEvents {
		event := v
		beginBlockEvents[i] = &event
	}
	endBlockEvents := make([]*abci.Event, len(res.EndBlockEvents))
	for i, v := range res.EndBlockEvents {
		event := v
		endBlockEvents[i] = &event
	}
	validatorUpdates := make([]*abci.ValidatorUpdate, len(res.ValidatorUpdates))
	for i, v := range validatorUpdates {
		validatorUpdates[i] = v
	}
	return &BlockResults{
		Height:                res.Height,
		TxsResults:            res.TxsResults,
		BeginBlockEvents:      beginBlockEvents,
		EndBlockEvents:        endBlockEvents,
		ValidatorUpdates:      validatorUpdates,
		ConsensusParamUpdates: res.ConsensusParamUpdates,
	}, nil
}

func getBlockchainInfo(ctx context.Context, clientCtx client.Context, minHeight, maxHeight int64) (*coretypes.ResultBlockchainInfo, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return &coretypes.ResultBlockchainInfo{}, err
	}
	return node.BlockchainInfo(ctx, minHeight, maxHeight)
}
