package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// requireHaltAdmin verifies that the given sender is the dedicated emergency
// halt-admin for the current chain.
func (k msgServer) requireHaltAdmin(ctx sdk.Context, sender string) error {
	admin := types.GetHaltAdmin(ctx.ChainID())
	if admin == "" {
		return errors.HaltAdminNotConfigured
	}
	if sender != admin {
		return errors.UnauthorizedHaltAdmin
	}
	return nil
}

// HaltChain schedules a hard consensus halt of the chain. The actual halt is
// performed by the validator BeginBlocker, which panics at and after the
// recorded height. Only the dedicated halt-admin may call this.
func (k msgServer) HaltChain(goCtx context.Context, msg *types.MsgHaltChain) (*types.MsgHaltChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.requireHaltAdmin(ctx, msg.Sender); err != nil {
		return nil, err
	}

	// Determine the effective halt height. 0 means "halt at the next block".
	haltHeight := msg.Height
	if haltHeight == 0 {
		haltHeight = ctx.BlockHeight() + 1
	} else if haltHeight <= ctx.BlockHeight() {
		return nil, errors.InvalidHaltHeight
	}

	k.SetHaltInfo(ctx, types.HaltInfo{
		Active: true,
		Height: haltHeight,
		Reason: msg.Reason,
		Sender: msg.Sender,
	})

	k.Logger(ctx).Error("EMERGENCY: hard chain halt scheduled",
		"height", haltHeight, "reason", msg.Reason, "admin", msg.Sender)

	err := events.EmitTypedEvent(ctx, &types.EventChainHaltScheduled{
		Sender: msg.Sender,
		Height: haltHeight,
		Reason: msg.Reason,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgHaltChainResponse{}, nil
}

// ResumeChain clears a scheduled or active hard halt. Note: once the chain has
// actually halted (BeginBlocker panicking), this transaction can only be
// included after operators restart their nodes with --unsafe-skip-halt.
func (k msgServer) ResumeChain(goCtx context.Context, msg *types.MsgResumeChain) (*types.MsgResumeChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.requireHaltAdmin(ctx, msg.Sender); err != nil {
		return nil, err
	}

	info, _ := k.GetHaltInfo(ctx)
	info.Active = false
	k.SetHaltInfo(ctx, info)

	k.Logger(ctx).Error("EMERGENCY: hard chain halt cleared", "admin", msg.Sender)

	err := events.EmitTypedEvent(ctx, &types.EventChainResumed{
		Sender: msg.Sender,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgResumeChainResponse{}, nil
}

// FreezeChain soft-freezes the chain so the ante handler rejects all
// transactions except emergency admin messages. The chain keeps producing
// empty blocks. Only the dedicated halt-admin may call this.
func (k msgServer) FreezeChain(goCtx context.Context, msg *types.MsgFreezeChain) (*types.MsgFreezeChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.requireHaltAdmin(ctx, msg.Sender); err != nil {
		return nil, err
	}

	k.SetFreezeInfo(ctx, types.FreezeInfo{
		Frozen: true,
		Reason: msg.Reason,
		Sender: msg.Sender,
	})

	k.Logger(ctx).Error("EMERGENCY: chain soft-frozen", "reason", msg.Reason, "admin", msg.Sender)

	err := events.EmitTypedEvent(ctx, &types.EventChainFrozen{
		Sender: msg.Sender,
		Reason: msg.Reason,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgFreezeChainResponse{}, nil
}

// UnfreezeChain lifts a soft freeze. Only the dedicated halt-admin may call this.
func (k msgServer) UnfreezeChain(goCtx context.Context, msg *types.MsgUnfreezeChain) (*types.MsgUnfreezeChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.requireHaltAdmin(ctx, msg.Sender); err != nil {
		return nil, err
	}

	info, _ := k.GetFreezeInfo(ctx)
	info.Frozen = false
	k.SetFreezeInfo(ctx, info)

	k.Logger(ctx).Error("EMERGENCY: chain soft-freeze lifted", "admin", msg.Sender)

	err := events.EmitTypedEvent(ctx, &types.EventChainUnfrozen{
		Sender: msg.Sender,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgUnfreezeChainResponse{}, nil
}
