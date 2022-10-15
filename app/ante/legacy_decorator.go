package ante

import (
	commonTypes "bitbucket.org/decimalteam/go-smart-node/types"
	legacyTypes "bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type LegacyActualizerDecorator struct {
	legacyKeeper legacyTypes.LegacyKeeper
}

// NewFeeDecorator creates new FeeDecorator to deduct fee
func NewLegacyActualizerDecorator(lk legacyTypes.LegacyKeeper) LegacyActualizerDecorator {
	return LegacyActualizerDecorator{
		legacyKeeper: lk,
	}
}

// AnteHandle implements sdk.AnteHandler function.
func (lad LegacyActualizerDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx,
	simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.ErrTxDecode
	}
	pubkeys, err := sigTx.GetPubKeys()
	if err != nil {
		return ctx, err
	}
	for _, pubkey := range pubkeys {
		if pubkey == nil {
			continue
		}
		legacyAddress, err := commonTypes.GetLegacyAddressFromPubKey(pubkey.Bytes())
		if err != nil {
			return ctx, err
		}
		if lad.legacyKeeper.IsLegacyAddress(ctx, legacyAddress) {
			err = lad.legacyKeeper.ActualizeLegacy(ctx, pubkey.Bytes())
			if err != nil {
				return ctx, err
			}
		}
	}
	return next(ctx, tx, simulate)
}
