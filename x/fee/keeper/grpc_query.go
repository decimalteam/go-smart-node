package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) QueryBaseDenomPrice(c context.Context, req *types.QueryBaseDenomPriceRequest) (*types.QueryBaseDenomPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryBaseDenomPriceResponse{Price: 100.}, nil
}
