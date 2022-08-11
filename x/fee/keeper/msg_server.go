package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"context"
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

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBaseDenomSaved,
			sdk.NewAttribute(types.AttributeKeyPrice, msg.Price.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.BaseDenom),
		),
	})

	return &types.MsgSaveBaseDenomPriceResponse{}, nil
}
