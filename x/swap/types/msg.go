package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgSwapInitialize{}
	_ sdk.Msg = &MsgSwapRedeem{}
	_ sdk.Msg = &MsgChainActivate{}
	_ sdk.Msg = &MsgChainDeactivate{}
)

const (
	TypeMsgSwapInitialize  = "swap_initialize"
	TypeMsgSwapRedeem      = "swap_redeem"
	TypeMsgChainActivate   = "chain_activate"
	TypeMsgChainDeactivate = "chain_deactivate"
)

////////////////////////////////////////////////////////////////
// MsgSwapInitialize
////////////////////////////////////////////////////////////////

// NewMsgSwapInitialize creates a new instance of MsgSwapInitialize.
func NewMsgSwapInitialize(
	sender sdk.AccAddress,
	recipient string,
	amount sdk.Int,
	tokenSymbol string,
	transactionNumber string,
	fromChain uint32,
	destChain uint32,
) *MsgSwapInitialize {
	return &MsgSwapInitialize{
		Sender:            sender.String(),
		Recipient:         recipient,
		Amount:            amount,
		TokenSymbol:       tokenSymbol,
		TransactionNumber: transactionNumber,
		FromChain:         fromChain,
		DestChain:         destChain,
	}
}

// Route should return the name of the module.
func (msg MsgSwapInitialize) Route() string { return RouterKey }

// Type should return the action.
func (msg MsgSwapInitialize) Type() string { return TypeMsgSwapInitialize }

// GetSignBytes encodes the message for signing.
func (msg *MsgSwapInitialize) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgSwapInitialize) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgSwapInitialize) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return ErrInvalidSenderAddress(msg.Sender)
	}
	if !msg.Amount.IsPositive() {
		return ErrInvalidAmount()
	}
	if msg.FromChain == 0 {
		return ErrInvalidChainNumber()
	}
	if msg.DestChain == 0 {
		return ErrInvalidChainNumber()
	}
	if _, ok := sdk.NewIntFromString(msg.TransactionNumber); !ok {
		return ErrInvalidTransactionNumber(msg.TransactionNumber)
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgSwapRedeem
////////////////////////////////////////////////////////////////

func NewMsgSwapRedeem(
	sender sdk.AccAddress,
	from string,
	recipient string,
	amount sdk.Int,
	tokenSymbol string,
	transactionNumber string,
	fromChain uint32,
	destChain uint32,
	v uint32,
	r *Hash,
	s *Hash,
) *MsgSwapRedeem {
	return &MsgSwapRedeem{
		Sender:            sender.String(),
		From:              from,
		Recipient:         recipient,
		Amount:            amount,
		TokenSymbol:       tokenSymbol,
		TransactionNumber: transactionNumber,
		FromChain:         fromChain,
		DestChain:         destChain,
		V:                 v,
		R:                 r.Copy(),
		S:                 s.Copy(),
	}
}

func (msg MsgSwapRedeem) Route() string { return RouterKey }

func (msg MsgSwapRedeem) Type() string { return TypeMsgSwapRedeem }

func (msg *MsgSwapRedeem) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSwapRedeem) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgSwapRedeem) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return ErrInvalidSenderAddress(msg.Sender)
	}
	if !msg.Amount.IsPositive() {
		return ErrInvalidAmount()
	}
	if msg.FromChain == 0 {
		return ErrInvalidChainNumber()
	}
	if msg.DestChain == 0 {
		return ErrInvalidChainNumber()
	}
	if _, ok := sdk.NewIntFromString(msg.TransactionNumber); !ok {
		return ErrInvalidTransactionNumber(msg.TransactionNumber)
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgChainActivate
////////////////////////////////////////////////////////////////

func NewMsgChainActivate(
	sender sdk.AccAddress,
	chainNumber uint32,
	chainName string,
) *MsgChainActivate {
	return &MsgChainActivate{
		Sender:      sender.String(),
		ChainNumber: chainNumber,
		ChainName:   chainName,
	}
}

func (msg MsgChainActivate) Route() string { return RouterKey }

func (msg MsgChainActivate) Type() string { return TypeMsgChainActivate }

func (msg *MsgChainActivate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgChainActivate) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgChainActivate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return ErrInvalidSenderAddress(msg.Sender)
	}

	if msg.Sender != SwapServiceAddress {
		return ErrSenderIsNotSwapService(msg.Sender)
	}

	if msg.ChainNumber == 0 {
		return ErrInvalidChainNumber()
	}

	return nil
}

////////////////////////////////////////////////////////////////
// MsgChainDeactivate
////////////////////////////////////////////////////////////////

func NewMsgChainDeactivate(
	sender sdk.AccAddress,
	chainNumber uint32,
) *MsgChainDeactivate {
	return &MsgChainDeactivate{
		Sender:      sender.String(),
		ChainNumber: chainNumber,
	}
}

func (msg MsgChainDeactivate) Route() string { return RouterKey }

func (msg MsgChainDeactivate) Type() string { return TypeMsgChainDeactivate }

func (msg *MsgChainDeactivate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgChainDeactivate) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgChainDeactivate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return ErrInvalidSenderAddress(msg.Sender)
	}

	if msg.Sender != SwapServiceAddress {
		return ErrSenderIsNotSwapService(msg.Sender)
	}

	if msg.ChainNumber == 0 {
		return ErrInvalidChainNumber()
	}

	return nil
}
