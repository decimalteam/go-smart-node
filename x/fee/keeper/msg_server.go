package keeper

import (
	"context"

	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) SaveBaseDenomPrice(c context.Context, msg *types.MsgSaveBaseDenomPrice) (*types.MsgSaveBaseDenomPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	params := k.GetModuleParams(ctx)
	if msg.Sender != params.OracleAddress {
		return nil, types.ErrUnknownOracle(msg.Sender)
	}

	err := k.SavePrice(ctx, msg.Price)
	if err != nil {
		return nil, types.ErrSavingError(err.Error())
	}

	return &types.MsgSaveBaseDenomPriceResponse{}, nil
}
