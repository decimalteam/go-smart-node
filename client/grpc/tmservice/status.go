package tmservice

import (
	"context"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/client"
)

func getNodeStatus(ctx context.Context, clientCtx client.Context) (*coretypes.ResultStatus, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return &coretypes.ResultStatus{}, err
	}
	return node.Status(ctx)
}

func getBlockchainInfo(ctx context.Context, clientCtx client.Context, minHeight, maxHeight int64) (*coretypes.ResultBlockchainInfo, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return &coretypes.ResultBlockchainInfo{}, err
	}
	return node.BlockchainInfo(ctx, minHeight, maxHeight)
}
