package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the NFT Querier
const (
	QuerySupply       = "supply"
	QueryOwner        = "owner"
	QueryOwnerByDenom = "ownerByDenom"
	QueryCollection   = "collection"
	QueryDenoms       = "denoms"
	QueryNFT          = "nft"
	QuerySubTokens    = "sub_tokens"
)

// NewQuerier is the module level router for state queries
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QuerySupply:
			return querySupply(ctx, path[1:], req, k, legacyQuerierCdc)
		case QueryOwner:
			return queryOwner(ctx, path[1:], req, k, legacyQuerierCdc)
		case QueryOwnerByDenom:
			return queryOwnerByDenom(ctx, path[1:], req, k, legacyQuerierCdc)
		case QueryCollection:
			return queryCollection(ctx, path[1:], req, k, legacyQuerierCdc)
		case QueryDenoms:
			return queryDenoms(ctx, path[1:], req, k, legacyQuerierCdc)
		case QueryNFT:
			return queryNFT(ctx, path[1:], req, k, legacyQuerierCdc)
		case QuerySubTokens:
			return querySubTokens(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown nft query endpoint")
		}
	}
}

func querySupply(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryCollectionParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("incorrectly formatted request data %v", err.Error()))
	}

	collection, found := k.GetCollection(ctx, params.Denom)
	if !found {
		return nil, types.ErrUnknownCollection(params.Denom)
	}

	bz, err := legacyQuerierCdc.MarshalJSON(strconv.Itoa(collection.Supply()))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryOwner(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryBalanceParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	owner := types.Owner{
		Address:     params.Owner.String(),
		Collections: k.GetOwnerCollections(ctx, params.Owner),
	}

	bz, err := legacyQuerierCdc.MarshalJSON(owner)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryOwnerByDenom(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryBalanceParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	ownerCollection, found := k.GetOwnerCollectionByDenom(ctx, params.Owner, params.Denom)
	if !found {
		ownerCollection = types.NewOwnerCollection(params.Denom, []string{})
	}

	owner := types.Owner{
		Address: params.Owner.String(),
		Collections: []types.OwnerCollection{
			ownerCollection,
		},
	}

	bz, err := legacyQuerierCdc.MarshalJSON(owner)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryCollection(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryCollectionParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	collection, found := k.GetCollection(ctx, params.Denom)
	if !found {
		return nil, types.ErrUnknownCollection(params.Denom)
	}

	bz, err := legacyQuerierCdc.MarshalJSON(collection)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDenoms(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	denoms := k.GetDenoms(ctx)

	bz, err := legacyQuerierCdc.MarshalJSON(denoms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryNFT(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryNFTParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	nft, err := k.GetNFT(ctx, params.Denom, params.TokenID)
	if err != nil {
		return nil, types.ErrUnknownNFT(params.Denom, params.TokenID)
	}

	bz, err := legacyQuerierCdc.MarshalJSON(nft)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func querySubTokens(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QuerySubTokensParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	var response types.ResponseSubTokens

	for _, id := range params.SubTokenIDs {
		subToken, ok := k.GetSubToken(ctx, params.TokenID, id)
		if !ok {
			continue
		}
		response = append(response, types.ResponseSubToken{
			ID:      subToken.ID,
			Reserve: subToken.Reserve,
		})
	}

	bz, err := legacyQuerierCdc.MarshalJSON(response)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
