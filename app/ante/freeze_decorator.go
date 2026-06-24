package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// FreezeKeeper is the subset of the validator keeper needed to read the
// emergency soft-freeze state.
type FreezeKeeper interface {
	IsFrozen(ctx sdk.Context) bool
}

// FreezeDecorator rejects every transaction while the chain has been soft-frozen
// by the emergency halt-admin, except the emergency control messages themselves
// (so the admin can still unfreeze, resume, or escalate to a hard halt). The
// chain keeps producing empty blocks while frozen.
type FreezeDecorator struct {
	vk FreezeKeeper
}

// NewFreezeDecorator creates a new FreezeDecorator.
func NewFreezeDecorator(vk FreezeKeeper) FreezeDecorator {
	return FreezeDecorator{vk: vk}
}

// AnteHandle implements sdk.AnteDecorator.
func (fd FreezeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if !fd.vk.IsFrozen(ctx) {
		return next(ctx, tx, simulate)
	}
	for _, msg := range tx.GetMsgs() {
		if !isEmergencyMsg(msg) {
			return ctx, ChainIsFrozen
		}
	}
	return next(ctx, tx, simulate)
}

// isEmergencyMsg reports whether the message is one of the emergency admin
// control messages that remain allowed while the chain is frozen.
func isEmergencyMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case *validatortypes.MsgHaltChain,
		*validatortypes.MsgResumeChain,
		*validatortypes.MsgFreezeChain,
		*validatortypes.MsgUnfreezeChain:
		return true
	default:
		return false
	}
}
