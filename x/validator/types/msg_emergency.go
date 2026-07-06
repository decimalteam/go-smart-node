package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgHaltChain{}
	_ sdk.Msg = &MsgResumeChain{}
	_ sdk.Msg = &MsgFreezeChain{}
	_ sdk.Msg = &MsgUnfreezeChain{}
)

const (
	TypeMsgHaltChain     = "halt_chain"
	TypeMsgResumeChain   = "resume_chain"
	TypeMsgFreezeChain   = "freeze_chain"
	TypeMsgUnfreezeChain = "unfreeze_chain"
)

////////////////////////////////////////////////////////////////
// MsgHaltChain
////////////////////////////////////////////////////////////////

// NewMsgHaltChain creates a new instance of MsgHaltChain.
func NewMsgHaltChain(sender sdk.AccAddress, height int64, reason string) *MsgHaltChain {
	return &MsgHaltChain{
		Sender: sender.String(),
		Height: height,
		Reason: reason,
	}
}

// Route should return the name of the module.
func (msg *MsgHaltChain) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgHaltChain) Type() string { return TypeMsgHaltChain }

// GetSignBytes encodes the message for signing.
func (msg *MsgHaltChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgHaltChain) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgHaltChain) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if msg.Height < 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("halt height must not be negative")
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgResumeChain
////////////////////////////////////////////////////////////////

// NewMsgResumeChain creates a new instance of MsgResumeChain.
func NewMsgResumeChain(sender sdk.AccAddress) *MsgResumeChain {
	return &MsgResumeChain{
		Sender: sender.String(),
	}
}

// Route should return the name of the module.
func (msg *MsgResumeChain) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgResumeChain) Type() string { return TypeMsgResumeChain }

// GetSignBytes encodes the message for signing.
func (msg *MsgResumeChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgResumeChain) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgResumeChain) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgFreezeChain
////////////////////////////////////////////////////////////////

// NewMsgFreezeChain creates a new instance of MsgFreezeChain.
func NewMsgFreezeChain(sender sdk.AccAddress, reason string) *MsgFreezeChain {
	return &MsgFreezeChain{
		Sender: sender.String(),
		Reason: reason,
	}
}

// Route should return the name of the module.
func (msg *MsgFreezeChain) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgFreezeChain) Type() string { return TypeMsgFreezeChain }

// GetSignBytes encodes the message for signing.
func (msg *MsgFreezeChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgFreezeChain) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgFreezeChain) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgUnfreezeChain
////////////////////////////////////////////////////////////////

// NewMsgUnfreezeChain creates a new instance of MsgUnfreezeChain.
func NewMsgUnfreezeChain(sender sdk.AccAddress) *MsgUnfreezeChain {
	return &MsgUnfreezeChain{
		Sender: sender.String(),
	}
}

// Route should return the name of the module.
func (msg *MsgUnfreezeChain) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgUnfreezeChain) Type() string { return TypeMsgUnfreezeChain }

// GetSignBytes encodes the message for signing.
func (msg *MsgUnfreezeChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgUnfreezeChain) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgUnfreezeChain) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	return nil
}
