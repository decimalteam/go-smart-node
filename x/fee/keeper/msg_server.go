package keeper

import (
	"context"

	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) UpdateCoinPrices(c context.Context, msg *types.MsgUpdateCoinPrices) (*types.MsgUpdateCoinPricesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	params := k.GetModuleParams(ctx)
	if msg.Oracle != params.Oracle {
		return nil, errors.UnknownOracle
	}

	for _, price := range msg.Prices {
		err := k.SavePrice(ctx, price)
		if err != nil {
			return nil, errors.SavingError
		}
	}
	err := events.EmitTypedEvent(ctx,
		&types.EventUpdateCoinPrices{
			Oracle: msg.Oracle,
			Prices: msg.Prices,
		},
	)
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgUpdateCoinPricesResponse{}, nil
}
