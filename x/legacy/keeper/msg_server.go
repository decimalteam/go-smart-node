package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	commonTypes "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) ReturnLegacy(goCtx context.Context, msg *types.MsgReturnLegacy) (*types.MsgReturnLegacyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.ActualizeLegacy(ctx, msg.PublicKey)
	if err != nil {
		return nil, err
	}
	// all errors already checked in ActualizeLegacy
	legacyAddress, _ := commonTypes.GetLegacyAddressFromPubKey(msg.PublicKey)
	legacySdkAddress := sdk.MustAccAddressFromBech32(legacyAddress)
	actualSdkAddress := sdk.AccAddress(ethsecp256k1.PubKey{Key: msg.PublicKey}.Address())

	// send coins to new address if we have some balance on old address
	coins := k.bankKeeper.GetAllBalances(ctx, legacySdkAddress)
	if !coins.IsZero() {
		// this send emit event 'transfer', which processed by worker
		err := k.bankKeeper.SendCoins(ctx, legacySdkAddress, actualSdkAddress, coins)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgReturnLegacyResponse{}, nil
}
