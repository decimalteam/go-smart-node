package keeper

import (
	"context"

	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Swap(c context.Context, req *types.QuerySwapRequest) (*types.QuerySwapResponse, error) {
	return &types.QuerySwapResponse{}, nil
}

func (k Keeper) ActiveSwaps(c context.Context, req *types.QueryActiveSwapsRequest) (*types.QueryActiveSwapsResponse, error) {
	return &types.QueryActiveSwapsResponse{}, nil
}

func (k Keeper) Pool(c context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	coins := k.GetLockedFunds(ctx)

	return &types.QueryPoolResponse{Amount: coins}, nil
}
