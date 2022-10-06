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
			err = k.iterateTokens(ctx, creator, collection.Denom, func(token *types.Token) bool {
				collection.Tokens = append(collection.Tokens, token)
				return false
			})
			if err != nil {
				return
			}

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

	err := k.iterateTokens(ctx, creator, collection.Denom, func(token *types.Token) bool {
		collection.Tokens = append(collection.Tokens, token)
		return false
	})
	if err != nil {
		return nil, err
	}

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
	err := k.iterateSubTokens(ctx, req.TokenId, func(subToken *types.SubToken) bool {
		if subToken.Reserve == nil {
			subToken.Reserve = &token.Reserve
		}

		token.SubTokens = append(token.SubTokens, subToken)

		return false
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryTokenResponse{
		Token: token,
	}, nil
}

func (k Keeper) SubToken(c context.Context, req *types.QuerySubTokenRequest) (*types.QuerySubTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	subTokenID, err := strconv.ParseInt(req.SubTokenId, 10, 32)
	if err != nil {
		return nil, err
	}

	// read NFT sub-token
	subToken, found := k.GetSubToken(ctx, req.TokenId, uint32(subTokenID))
	if !found {
		return nil, errors.UnknownNFT
	}

	if subToken.Reserve == nil {
		token, found := k.GetToken(ctx, req.TokenId)
		if !found {
			return nil, errors.InvalidTokenID
		}

		subToken.Reserve = &token.Reserve
	}

	return &types.QuerySubTokenResponse{
		SubToken: subToken,
	}, nil
}
