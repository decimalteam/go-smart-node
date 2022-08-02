package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) SaveBaseDenomPrice(c context.Context, msg *types.MsgSaveBaseDenomPrice) (*types.MsgSaveBaseDenomPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_ = k.SavePrice(ctx, msg.Price)

	return &types.MsgSaveBaseDenomPriceResponse{}, nil
}
