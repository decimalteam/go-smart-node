package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/nft/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) QueryCollectionSupply(c context.Context, req *types.QueryCollectionSupplyRequest) (*types.QueryCollectionSupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	collection, found := k.GetCollection(ctx, req.Denom)
	if !found {
		return nil, errors.UnknownCollection
	}

	return &types.QueryCollectionSupplyResponse{Supply: int64(collection.Supply())}, nil
}

func (k Keeper) QueryOwnerCollections(c context.Context, req *types.QueryOwnerCollectionsRequest) (*types.QueryOwnerCollectionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ownerAddress, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	owner := types.Owner{
		Address: req.GetOwner(),
	}

	owner.Collections, err = k.GetOwnerCollections(ctx, ownerAddress)
	if err != nil {
		return nil, err
	}

	return &types.QueryOwnerCollectionsResponse{Owner: owner}, nil
}

func (k Keeper) QueryCollection(c context.Context, req *types.QueryCollectionRequest) (*types.QueryCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	collection, found := k.GetCollection(ctx, req.Denom)
	if !found {
		return nil, errors.UnknownCollection
	}

	return &types.QueryCollectionResponse{
		Collection: collection,
	}, nil
}

func (k Keeper) QueryDenoms(c context.Context, _ *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	denoms, err := k.GetDenoms(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryDenomsResponse{
		Denoms: denoms,
	}, nil
}

func (k Keeper) QueryNFT(c context.Context, req *types.QueryNFTRequest) (*types.QueryNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	nft, err := k.GetNFT(ctx, req.Denom, req.TokenId)
	if err != nil {
		return nil, errors.UnknownNFT
	}

	return &types.QueryNFTResponse{
		NFT: nft,
	}, nil
}

func (k Keeper) QuerySubTokens(c context.Context, req *types.QuerySubTokensRequest) (*types.QuerySubTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var subTokens []types.SubToken
	subTokens, err := k.GetSubTokens(ctx, req.TokenID)
	if err != nil {
		return nil, err
	}

	return &types.QuerySubTokensResponse{
		SubTokens: subTokens,
	}, nil
}
