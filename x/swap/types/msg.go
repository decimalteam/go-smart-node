package types

import (
	"bitbucket.org/decimalteam/go-smart-node/x/swap/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgInitializeSwap{}
	_ sdk.Msg = &MsgRedeemSwap{}
	_ sdk.Msg = &MsgActivateChain{}
	_ sdk.Msg = &MsgDeactivateChain{}
)

const (
	TypeMsgSwapInitialize  = "swap_initialize"
	TypeMsgSwapRedeem      = "swap_redeem"
	TypeMsgChainActivate   = "chain_activate"
	TypeMsgChainDeactivate = "chain_deactivate"
)

////////////////////////////////////////////////////////////////
// MsgInitializeSwap
////////////////////////////////////////////////////////////////

// NewMsgInitializeSwap creates a new instance of MsgInitializeSwap.
func NewMsgInitializeSwap(
	sender sdk.AccAddress,
	recipient string,
	amount sdk.Int,
	tokenSymbol string,
	transactionNumber string,
	fromChain uint32,
	destChain uint32,
) *MsgInitializeSwap {
	return &MsgInitializeSwap{
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
func (msg MsgInitializeSwap) Route() string { return RouterKey }

// Type should return the action.
func (msg MsgInitializeSwap) Type() string { return TypeMsgSwapInitialize }

// GetSignBytes encodes the message for signing.
func (msg *MsgInitializeSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgInitializeSwap) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgInitializeSwap) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}
	if !msg.Amount.IsPositive() {
		return errors.InvalidAmount
	}
	if msg.FromChain == 0 {
		return errors.InvalidChainNumber
	}
	if msg.DestChain == 0 {
		return errors.InvalidChainNumber
	}
	if msg.FromChain == msg.DestChain {
		return errors.ChainNumbersAreSame
	}
	if _, ok := sdk.NewIntFromString(msg.TransactionNumber); !ok {
		return errors.InvalidTransactionNumber
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgRedeemSwap
////////////////////////////////////////////////////////////////

func NewMsgRedeemSwap(
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
) *MsgRedeemSwap {
	return &MsgRedeemSwap{
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

func (msg MsgRedeemSwap) Route() string { return RouterKey }

func (msg MsgRedeemSwap) Type() string { return TypeMsgSwapRedeem }

func (msg *MsgRedeemSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRedeemSwap) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgRedeemSwap) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}
	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return errors.InvalidRecipientAddress
	}
	if !msg.Amount.IsPositive() {
		return errors.InvalidAmount
	}
	if msg.FromChain == 0 {
		return errors.InvalidChainNumber
	}
	if msg.DestChain == 0 {
		return errors.InvalidChainNumber
	}
	if msg.FromChain == msg.DestChain {
		return errors.ChainNumbersAreSame
	}
	if _, ok := sdk.NewIntFromString(msg.TransactionNumber); !ok {
		return errors.InvalidTransactionNumber
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgChainActivate
////////////////////////////////////////////////////////////////

func NewMsgActivateChain(
	sender sdk.AccAddress,
	id uint32,
	name string,
) *MsgActivateChain {
	return &MsgActivateChain{
		Sender: sender.String(),
		ID:     id,
		Name:   name,
	}
}

func (msg MsgActivateChain) Route() string { return RouterKey }

func (msg MsgActivateChain) Type() string { return TypeMsgChainActivate }

func (msg *MsgActivateChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgActivateChain) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgActivateChain) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}

	if msg.ID == 0 {
		return errors.InvalidChainNumber
	}

	if msg.Name == "" {
		return errors.InvalidChainName
	}

	return nil
}

////////////////////////////////////////////////////////////////
// MsgDeactivateChain
////////////////////////////////////////////////////////////////

func NewMsgDeactivateChain(
	sender sdk.AccAddress,
	id uint32,
) *MsgDeactivateChain {
	return &MsgDeactivateChain{
		Sender: sender.String(),
		ID:     id,
	}
}

func (msg MsgDeactivateChain) Route() string { return RouterKey }

func (msg MsgDeactivateChain) Type() string { return TypeMsgChainDeactivate }

func (msg *MsgDeactivateChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDeactivateChain) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgDeactivateChain) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}

	if msg.ID == 0 {
		return errors.InvalidChainNumber
	}

	return nil
}
