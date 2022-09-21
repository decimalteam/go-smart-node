package keeper

import (
	"context"

	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) ReturnLegacy(goCtx context.Context, msg *types.MsgReturnLegacy) (*types.MsgReturnLegacyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.ActualizeLegacy(ctx, msg.PublicKey)
	if err != nil {
		return nil, err
	}
	return &types.MsgReturnLegacyResponse{}, nil
}
