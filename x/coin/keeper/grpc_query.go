package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Coin(c context.Context, req *types.QueryCoinRequest) (*types.QueryCoinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	coin, err := k.GetCoin(ctx, req.Symbol)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCoinResponse{Coin: coin}, nil
}

func (k Keeper) Coins(c context.Context, req *types.QueryCoinsRequest) (*types.QueryCoinsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCoin)

	coins := []types.Coin{}

	pageRes, err := query.Paginate(
		store,
		req.Pagination,
		func(key, value []byte) error {
			var coin types.Coin
			if err := k.cdc.Unmarshal(value, &coin); err != nil {
				return err
			}
			coins = append(coins, coin)
			return nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCoinsResponse{
		Coins:      coins,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) Check(c context.Context, req *types.QueryCheckRequest) (*types.QueryCheckResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	check, err := k.GetCheck(ctx, req.Hash)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCheckResponse{Check: *check}, nil
}

func (k Keeper) Checks(c context.Context, req *types.QueryChecksRequest) (*types.QueryChecksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCheck)

	checks := []types.Check{}

	pageRes, err := query.Paginate(
		store,
		req.Pagination,
		func(key, value []byte) error {
			var check types.Check
			if err := k.cdc.Unmarshal(value, &check); err != nil {
				return err
			}
			checks = append(checks, check)
			return nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryChecksResponse{
		Checks:     checks,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}
