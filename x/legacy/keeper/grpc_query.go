package keeper

import (
	"context"

	commonTypes "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Records(c context.Context, req *types.QueryRecordsRequest) (*types.QueryRecordsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)

	records := make([]types.Record, 0)
	pageRes, err := query.Paginate(
		store,
		req.Pagination,
		func(_, value []byte) (err error) {
			var rec types.Record
			err = k.cdc.UnmarshalLengthPrefixed(value, &rec)
			if err != nil {
				return
			}
			records = append(records, rec)
			return
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRecordsResponse{
		Records:    records,
		Pagination: pageRes,
	}, nil

}

func (k Keeper) Record(c context.Context, req *types.QueryRecordRequest) (*types.QueryRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	rec, err := k.GetLegacyRecord(ctx, req.LegacyAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRecordResponse{
		Record: rec,
	}, nil
}

func (k Keeper) Check(c context.Context, req *types.QueryCheckRequest) (*types.QueryCheckResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	legacyAddress, err := commonTypes.GetLegacyAddressFromPubKey(req.Pubkey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	rec, err := k.GetLegacyRecord(ctx, legacyAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCheckResponse{
		Record: rec,
	}, nil
}
