package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/config"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
)

var (
	_ sdk.Msg = &MsgCreateCoin{}
	_ sdk.Msg = &MsgUpdateCoin{}
	_ sdk.Msg = &MsgSendCoin{}
	_ sdk.Msg = &MsgMultiSendCoin{}
	_ sdk.Msg = &MsgBuyCoin{}
	_ sdk.Msg = &MsgSellCoin{}
	_ sdk.Msg = &MsgSellAllCoin{}
	_ sdk.Msg = &MsgBurnCoin{}
	_ sdk.Msg = &MsgRedeemCheck{}
)

const (
	TypeMsgCreateCoin    = "create_coin"
	TypeMsgUpdateCoin    = "update_coin"
	TypeMsgSendCoin      = "send_coin"
	TypeMsgMultiSendCoin = "multi_send_coin"
	TypeMsgBuyCoin       = "buy_coin"
	TypeMsgSellCoin      = "sell_coin"
	TypeMsgSellAllCoin   = "sell_all_coin"
	TypeMsgBurnCoin      = "burn_coin"
	TypeMsgRedeemCheck   = "redeem_check"
)

////////////////////////////////////////////////////////////////
// MsgCreateCoin
////////////////////////////////////////////////////////////////

// NewMsgCreateCoin creates a new instance of MsgCreateCoin.
func NewMsgCreateCoin(
	sender sdk.AccAddress,
	denom string,
	title string,
	crr uint64,
	initVolume sdkmath.Int,
	initReserve sdkmath.Int,
	limitVolume sdkmath.Int,
	minVolume sdkmath.Int,
	identity string,
) *MsgCreateCoin {
	return &MsgCreateCoin{
		Sender:         sender.String(),
		Denom:          denom,
		Title:          title,
		CRR:            uint32(crr),
		InitialVolume:  initVolume,
		InitialReserve: initReserve,
		LimitVolume:    limitVolume,
		MinVolume:      minVolume,
		Identity:       identity,
	}
}

// Route should return the name of the module.
func (msg *MsgCreateCoin) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgCreateCoin) Type() string { return TypeMsgCreateCoin }

// GetSignBytes encodes the message for signing.
func (msg *MsgCreateCoin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgCreateCoin) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgCreateCoin) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}
	// Validate coin title
	if len(msg.Title) > config.MaxCoinTitleLength {
		return errors.InvalidCoinTitle
	}
	// Validate coin denom
	if !config.CoinDenomValidator.MatchString(msg.Denom) {
		return errors.InvalidCoinDenom
	}
	// TODO: Looks like no need since should be no more possible to create such coins anyway
	// // Forbid creating coin with symbol DEL in testnet
	// if strings.HasPrefix(config.ChainID, "decimal-testnet") {
	// 	if strings.ToLower(msg.Denom) == config.DenomBaseCoin {
	// 		return ErrForbiddenCoinDenom(msg.Denom)
	// 	}
	// }
	// Validate coin CRR
	if msg.CRR < 10 || msg.CRR > 100 {
		return errors.InvalidCRR
	}
	// Check coin initial volume to be correct
	if msg.InitialVolume.LT(config.MinCoinSupply) || msg.InitialVolume.GT(config.MaxCoinSupply) {
		return errors.InvalidCoinInitialVolume
	}
	if msg.InitialVolume.GT(msg.LimitVolume) {
		return errors.InvalidLimitVolume
	}
	// Check coin initial reserve to be correct
	if msg.InitialReserve.LT(config.MinCoinReserve) {
		return errors.InvalidCoinInitialReserve
	}
	// Check limit volume
	if msg.LimitVolume.GT(config.MaxCoinSupply) {
		return errors.InvalidLimitVolume
	}
	// Check coin min emission to be disabled or enabled correctly
	if !msg.MinVolume.IsNil() && !msg.MinVolume.IsZero() {
		if msg.MinVolume.LT(config.MinCoinSupply) || msg.MinVolume.GT(config.MaxCoinSupply) {
			return errors.InvalidCoinMinEmission
		}
		if msg.MinVolume.GT(msg.LimitVolume) {
			return errors.TooBigCoinMinEmission
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgUpdateCoin
////////////////////////////////////////////////////////////////

// NewMsgUpdateCoin creates a new instance of MsgUpdateCoin.
func NewMsgUpdateCoin(
	sender sdk.AccAddress,
	denom string,
	limitVolume sdkmath.Int,
	minVolume sdkmath.Int,
	identity string,
) *MsgUpdateCoin {
	return &MsgUpdateCoin{
		Sender:      sender.String(),
		Denom:       denom,
		LimitVolume: limitVolume,
		MinVolume:   minVolume,
		Identity:    identity,
	}
}

// Route should return the name of the module.
func (msg *MsgUpdateCoin) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgUpdateCoin) Type() string { return TypeMsgUpdateCoin }

// GetSignBytes encodes the message for signing.
func (msg *MsgUpdateCoin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgUpdateCoin) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgUpdateCoin) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}
	// Validate coin denom
	if !config.CoinDenomValidator.MatchString(msg.Denom) {
		return errors.InvalidCoinDenom
	}
	// Check coin limit volume to be less than max coin supply
	if msg.LimitVolume.GT(config.MaxCoinSupply) {
		return errors.InvalidLimitVolume
	}
	// Check coin min emission to be disabled or enabled correctly
	if !msg.MinVolume.IsNil() && !msg.MinVolume.IsZero() {
		if msg.MinVolume.LT(config.MinCoinSupply) || msg.MinVolume.GT(config.MaxCoinSupply) {
			return errors.InvalidCoinMinEmission
		}
		if msg.MinVolume.GT(msg.LimitVolume) {
			return errors.TooBigCoinMinEmission
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgSendCoin
////////////////////////////////////////////////////////////////

// NewMsgSendCoin creates a new instance of MsgSendCoin.
func NewMsgSendCoin(
	sender sdk.AccAddress,
	recipient sdk.AccAddress,
	coin sdk.Coin,
) *MsgSendCoin {
	return &MsgSendCoin{
		Sender:    sender.String(),
		Recipient: recipient.String(),
		Coin:      coin,
	}
}

// Route should return the name of the module.
func (msg *MsgSendCoin) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgSendCoin) Type() string { return TypeMsgSendCoin }

// GetSignBytes encodes the message for signing.
func (msg *MsgSendCoin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgSendCoin) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgSendCoin) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}
	// Validate recipient
	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return errors.InvalidRecipientAddress
	}
	if msg.Sender == msg.Recipient {
		return errors.InvalidRecipientAddress
	}
	// Validate coin denom
	if !config.CoinDenomValidator.MatchString(msg.Coin.Denom) {
		return errors.InvalidCoinDenom
	}
	// Validate coin amount
	if !msg.Coin.Amount.IsPositive() {
		return errors.InvalidAmount
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgMultiSendCoin
////////////////////////////////////////////////////////////////

// NewMsgMultiSendCoin creates a new instance of MsgMultiSendCoin.
func NewMsgMultiSendCoin(
	sender sdk.AccAddress,
	sends []MultiSendEntry,
) *MsgMultiSendCoin {
	return &MsgMultiSendCoin{
		Sender: sender.String(),
		Sends:  sends,
	}
}

// Route should return the name of the module.
func (msg *MsgMultiSendCoin) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgMultiSendCoin) Type() string { return TypeMsgMultiSendCoin }

// GetSignBytes encodes the message for signing.
func (msg *MsgMultiSendCoin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgMultiSendCoin) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgMultiSendCoin) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}
	if len(msg.Sends) == 0 {
		return errors.InvalidAmount
	}
	for i := range msg.Sends {
		// Validate recipient
		if _, err := sdk.AccAddressFromBech32(msg.Sends[i].Recipient); err != nil {
			return errors.InvalidRecipientAddress
		}
		if msg.Sender == msg.Sends[i].Recipient {
			return errors.InvalidRecipientAddress
		}
		// Validate coin denom
		if !config.CoinDenomValidator.MatchString(msg.Sends[i].Coin.Denom) {
			return errors.InvalidCoinDenom
		}
		// Validate coin amount
		if !msg.Sends[i].Coin.Amount.IsPositive() {
			return errors.InvalidAmount
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgBuyCoin
////////////////////////////////////////////////////////////////

// NewMsgBuyCoin creates a new instance of MsgBuyCoin.
func NewMsgBuyCoin(
	sender sdk.AccAddress,
	coinToBuy sdk.Coin,
	maxCoinToSell sdk.Coin,
) *MsgBuyCoin {
	return &MsgBuyCoin{
		Sender:        sender.String(),
		CoinToBuy:     coinToBuy,
		MaxCoinToSell: maxCoinToSell,
	}
}

// Route should return the name of the module.
func (msg *MsgBuyCoin) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgBuyCoin) Type() string { return TypeMsgBuyCoin }

// GetSignBytes encodes the message for signing.
func (msg *MsgBuyCoin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgBuyCoin) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgBuyCoin) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}
	// Denoms of specified coins cannot be the same
	if msg.CoinToBuy.Denom == msg.MaxCoinToSell.Denom {
		return errors.SameCoin
	}
	// Validate coin denom
	if !config.CoinDenomValidator.MatchString(msg.CoinToBuy.Denom) {
		return errors.InvalidCoinDenom
	}
	if !config.CoinDenomValidator.MatchString(msg.MaxCoinToSell.Denom) {
		return errors.InvalidCoinDenom
	}
	// Check amount
	if !msg.CoinToBuy.Amount.IsPositive() {
		return errors.InvalidAmount
	}
	if !msg.MaxCoinToSell.Amount.IsPositive() {
		return errors.InvalidAmount
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgSellCoin
////////////////////////////////////////////////////////////////

// NewMsgSellCoin creates a new instance of MsgSellCoin.
func NewMsgSellCoin(
	sender sdk.AccAddress,
	coinToSell sdk.Coin,
	minCoinToBuy sdk.Coin,
) *MsgSellCoin {
	return &MsgSellCoin{
		Sender:       sender.String(),
		CoinToSell:   coinToSell,
		MinCoinToBuy: minCoinToBuy,
	}
}

// Route should return the name of the module.
func (msg *MsgSellCoin) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgSellCoin) Type() string { return TypeMsgSellCoin }

// GetSignBytes encodes the message for signing.
func (msg *MsgSellCoin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgSellCoin) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgSellCoin) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}
	// Denoms of specified coins cannot be the same
	if msg.CoinToSell.Denom == msg.MinCoinToBuy.Denom {
		return errors.SameCoin
	}
	// Validate coin denom
	if !config.CoinDenomValidator.MatchString(msg.CoinToSell.Denom) {
		return errors.InvalidCoinDenom
	}
	if !config.CoinDenomValidator.MatchString(msg.MinCoinToBuy.Denom) {
		return errors.InvalidCoinDenom
	}
	// Check amount
	if !msg.CoinToSell.Amount.IsPositive() {
		return errors.InvalidAmount
	}
	// sdk.Coin amount can not be negative
	// if msg.MinCoinToBuy.Amount.IsZero() {
	//	return ErrInvalidAmount()
	// }

	return nil
}

////////////////////////////////////////////////////////////////
// MsgSellAllCoin
////////////////////////////////////////////////////////////////

// NewMsgSellAllCoin creates a new instance of MsgSellAllCoin.
func NewMsgSellAllCoin(
	sender sdk.AccAddress,
	coinDenomToSell string,
	minCoinToBuy sdk.Coin,
) *MsgSellAllCoin {
	return &MsgSellAllCoin{
		Sender:          sender.String(),
		CoinDenomToSell: coinDenomToSell,
		MinCoinToBuy:    minCoinToBuy,
	}
}

// Route should return the name of the module.
func (msg *MsgSellAllCoin) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgSellAllCoin) Type() string { return TypeMsgSellAllCoin }

// GetSignBytes encodes the message for signing.
func (msg *MsgSellAllCoin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgSellAllCoin) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgSellAllCoin) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}
	// Validate coin denom
	if !config.CoinDenomValidator.MatchString(msg.CoinDenomToSell) {
		return errors.InvalidCoinDenom
	}
	if !config.CoinDenomValidator.MatchString(msg.MinCoinToBuy.Denom) {
		return errors.InvalidCoinDenom
	}
	// Denoms of specified coins cannot be the same
	if msg.CoinDenomToSell == msg.MinCoinToBuy.Denom {
		return errors.SameCoin
	}
	// Check amount
	if !msg.MinCoinToBuy.Amount.IsPositive() {
		return errors.InvalidAmount
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgBurnCoin
////////////////////////////////////////////////////////////////

// NewMsgSendCoin creates a new instance of MsgSendCoin.
func NewMsgBurnCoin(
	sender sdk.AccAddress,
	coin sdk.Coin,
) *MsgBurnCoin {
	return &MsgBurnCoin{
		Sender: sender.String(),
		Coin:   coin,
	}
}

// Route should return the name of the module.
func (msg *MsgBurnCoin) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgBurnCoin) Type() string { return TypeMsgSendCoin }

// GetSignBytes encodes the message for signing.
func (msg *MsgBurnCoin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgBurnCoin) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgBurnCoin) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}
	// Validate coin denom
	if !config.CoinDenomValidator.MatchString(msg.Coin.Denom) {
		return errors.InvalidCoinDenom
	}
	// Amount should be positive
	if !msg.Coin.Amount.IsPositive() {
		return errors.InvalidAmount
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgRedeemCheck
////////////////////////////////////////////////////////////////

// NewMsgRedeemCheck creates a new instance of MsgRedeemCheck.
func NewMsgRedeemCheck(
	sender sdk.AccAddress,
	check string,
	proof string,
) *MsgRedeemCheck {
	return &MsgRedeemCheck{
		Sender: sender.String(),
		Check:  check,
		Proof:  proof,
	}
}

// Route should return the name of the module.
func (msg *MsgRedeemCheck) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgRedeemCheck) Type() string { return TypeMsgRedeemCheck }

// GetSignBytes encodes the message for signing.
func (msg *MsgRedeemCheck) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgRedeemCheck) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgRedeemCheck) ValidateBasic() error {
	// Validate sender
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errors.InvalidSenderAddress
	}
	// TODO
	return nil
}
