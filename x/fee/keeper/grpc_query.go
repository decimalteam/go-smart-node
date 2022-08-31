package keeper

import (
	"context"

	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}
var _ feemarkettypes.QueryServer = Keeper{}

/////////////
// Fee Keeper
/////////////

func (k Keeper) QueryBaseDenomPrice(c context.Context, req *types.QueryBaseDenomPriceRequest) (*types.QueryBaseDenomPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	price, err := k.GetPrice(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBaseDenomPriceResponse{Price: price.String()}, nil
}

func (k Keeper) QueryModuleParams(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	params := k.GetModuleParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

/////////////
// Fee Market Keeper
/////////////

// BaseFee implements the Query/BaseFee gRPC method
func (k Keeper) BaseFee(c context.Context, _ *feemarkettypes.QueryBaseFeeRequest) (*feemarkettypes.QueryBaseFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	res := &feemarkettypes.QueryBaseFeeResponse{}
	baseFee := k.GetBaseFee(ctx)

	if baseFee != nil {
		aux := sdkmath.NewIntFromBigInt(baseFee)
		res.BaseFee = &aux
	}

	return res, nil
}

// Params implements the Query/Params gRPC method
func (k Keeper) Params(c context.Context, _ *feemarkettypes.QueryParamsRequest) (*feemarkettypes.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &feemarkettypes.QueryParamsResponse{
		Params: params,
	}, nil
}

// BlockGas implements the Query/BlockGas gRPC method
func (k Keeper) BlockGas(c context.Context, _ *feemarkettypes.QueryBlockGasRequest) (*feemarkettypes.QueryBlockGasResponse, error) {
	// TODO: do we need something

	return &feemarkettypes.QueryBlockGasResponse{
		Gas: 0,
	}, nil
}
