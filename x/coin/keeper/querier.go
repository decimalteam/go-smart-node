package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, path[1:], req, k, legacyQuerierCdc)
		case types.QueryCoin:
			return queryCoin(ctx, path[1:], req, k, legacyQuerierCdc)
		case types.QueryCoins:
			return queryCoins(ctx, path[1:], req, k, legacyQuerierCdc)
		case types.QueryCheck:
			return queryCheck(ctx, path[1:], req, k, legacyQuerierCdc)
		case types.QueryChecks:
			return queryChecks(ctx, path[1:], req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, _ []string, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, k.GetParams(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryCoin(ctx sdk.Context, path []string, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	coinDenom := path[0]

	coin, err := k.GetCoin(ctx, coinDenom)
	if err != nil {
		return nil, errors.CoinDoesNotExist
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, coin)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryCoins(ctx sdk.Context, _ []string, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	coins := k.GetCoins(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, coins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryCheck(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := &types.QueryCheckParams{}

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	check, err := k.GetCheck(ctx, params.Hash)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrNotFound, err.Error()) //TODO: how err write this?
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, check)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryChecks(ctx sdk.Context, _ []string, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	checks := k.GetChecks(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, checks)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
