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

func (k Keeper) CoinPrices(c context.Context, req *types.QueryCoinPricesRequest) (*types.QueryCoinPricesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	prices, err := k.GetPrices(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCoinPricesResponse{Prices: prices}, nil
}

func (k Keeper) CoinPrice(c context.Context, req *types.QueryCoinPriceRequest) (*types.QueryCoinPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	price, err := k.GetPrice(ctx, req.Denom, req.Quote)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCoinPriceResponse{Price: &price}, nil
}

func (k Keeper) ModuleParams(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
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
	// TODO: rework when ethermint starts use BlockGas from EVM
	return &feemarkettypes.QueryBlockGasResponse{
		Gas: 0,
	}, nil
}
