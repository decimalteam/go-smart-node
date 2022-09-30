package tmservice

import (
	"context"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

func getBlock(ctx context.Context, clientCtx client.Context, height *int64) (*coretypes.ResultBlock, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	return node.Block(ctx, height)
}

func getBlockResults(ctx context.Context, clientCtx client.Context, height int64) (*BlockResults, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}
	res, err := node.BlockResults(ctx, &height)

	beginBlockEvents := make([]*abci.Event, len(res.BeginBlockEvents))
	for i, v := range res.BeginBlockEvents {
		beginBlockEvents[i] = &v
	}
	endBlockEvents := make([]*abci.Event, len(res.EndBlockEvents))
	for i, v := range res.EndBlockEvents {
		endBlockEvents[i] = &v
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

func GetProtoBlock(ctx context.Context, clientCtx client.Context, height *int64) (tmproto.BlockID, *tmproto.Block, error) {
	block, err := getBlock(ctx, clientCtx, height)
	if err != nil {
		return tmproto.BlockID{}, nil, err
	}
	protoBlock, err := block.Block.ToProto()
	if err != nil {
		return tmproto.BlockID{}, nil, err
	}
	protoBlockID := block.BlockID.ToProto()

	return protoBlockID, protoBlock, nil
}
