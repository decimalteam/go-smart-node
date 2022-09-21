package keeper

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"bitbucket.org/decimalteam/go-smart-node/x/nft/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Collections(c context.Context, req *types.QueryCollectionsRequest) (*types.QueryCollectionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	storePrefixed := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetCollectionsKey())

	// iterate over NFT collections with pagination
	collections := make([]types.Collection, 0)
	pageRes, err := query.Paginate(
		storePrefixed,
		req.Pagination,
		func(_, value []byte) (err error) {

			// parse NFT collection since it is already read
			var collection types.Collection
			err = k.cdc.UnmarshalLengthPrefixed(value, &collection)
			if err != nil {
				return
			}

			// read NFT collection counter separately
			creator := sdk.MustAccAddressFromBech32(collection.Creator)
			counter := k.getCollectionCounter(ctx, creator, collection.Denom)
			collection.Supply = counter.Supply

			collections = append(collections, collection)
			return
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCollectionsResponse{
		Collections: collections,
		Pagination:  pageRes,
	}, nil
}

func (k Keeper) CollectionsByCreator(c context.Context, req *types.QueryCollectionsByCreatorRequest) (*types.QueryCollectionsByCreatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	creator := sdk.MustAccAddressFromBech32(req.Creator)
	storePrefixed := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetCollectionsByCreatorKey(creator))

	// iterate over NFT collections with pagination
	collections := make([]types.Collection, 0)
	pageRes, err := query.Paginate(
		storePrefixed,
		req.Pagination,
		func(_, value []byte) (err error) {

			// parse NFT collection since it is already read
			var collection types.Collection
			err = k.cdc.UnmarshalLengthPrefixed(value, &collection)
			if err != nil {
				return
			}

			// read NFT collection counter separately
			creator := sdk.MustAccAddressFromBech32(collection.Creator)
			counter := k.getCollectionCounter(ctx, creator, collection.Denom)
			collection.Supply = counter.Supply

			// read NFT tokens within the collection
			k.iterateTokens(ctx, creator, collection.Denom, func(token *types.Token) bool {
				collection.Tokens = append(collection.Tokens, token)
				return false
			})

			collections = append(collections, collection)
			return
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCollectionsByCreatorResponse{
		Collections: collections,
		Pagination:  pageRes,
	}, nil
}

func (k Keeper) Collection(c context.Context, req *types.QueryCollectionRequest) (*types.QueryCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	creator := sdk.MustAccAddressFromBech32(req.Creator)

	// read NFT collection
	collection, found := k.GetCollection(ctx, creator, req.Denom)
	if !found {
		return nil, errors.UnknownCollection
	}
	// read NFT tokens within the collection

	k.iterateTokens(ctx, creator, collection.Denom, func(token *types.Token) bool {
		collection.Tokens = append(collection.Tokens, token)
		return false
	})

	return &types.QueryCollectionResponse{
		Collection: collection,
	}, nil
}

func (k Keeper) Token(c context.Context, req *types.QueryTokenRequest) (*types.QueryTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// read NFT token
	token, found := k.GetToken(ctx, req.TokenId)
	if !found {
		return nil, errors.UnknownNFT
	}

	// read NFT sub-tokens within the token
	k.iterateSubTokens(ctx, req.TokenId, func(subToken *types.SubToken) bool {
		token.SubTokens = append(token.SubTokens, subToken)
		return false
	})

	return &types.QueryTokenResponse{
		Token: token,
	}, nil
}

func (k Keeper) SubToken(c context.Context, req *types.QuerySubTokenRequest) (*types.QuerySubTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	subTokenId, err := strconv.ParseInt(req.SubTokenId, 10, 32)
	if err != nil {
		return nil, err
	}

	// read NFT sub-token
	subToken, found := k.GetSubToken(ctx, req.TokenId, uint32(subTokenId))
	if !found {
		return nil, errors.UnknownNFT
	}

	return &types.QuerySubTokenResponse{
		SubToken: subToken,
	}, nil
}
