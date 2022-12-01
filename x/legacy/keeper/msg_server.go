package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	commonTypes "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/errors"
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
	actualAddress, _ := bech32.ConvertAndEncode(config.Bech32Prefix, actualSdkAddress)

	// send coins to new address if we have some balance on old address
	coins := k.bankKeeper.GetAllBalances(ctx, legacySdkAddress)
	if !coins.IsZero() {
		err := k.bankKeeper.SendCoins(ctx, legacySdkAddress, actualSdkAddress, coins)
		if err != nil {
			return nil, err
		}
		// Emit send event
		err = events.EmitTypedEvent(ctx, &types.EventReturnLegacyCoins{
			LegacyOwner: legacyAddress,
			Owner:       actualAddress,
			Coins:       coins,
		})
		if err != nil {
			return nil, errors.Internal.Wrapf("err: %s", err.Error())
		}
	}

	return &types.MsgReturnLegacyResponse{}, nil
}
