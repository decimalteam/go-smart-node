package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
)

var (
	_ sdk.Msg                            = &MsgCreateValidator{}
	_ codectypes.UnpackInterfacesMessage = (*MsgCreateValidator)(nil)
	_ sdk.Msg                            = &MsgEditValidator{}
	_ sdk.Msg                            = &MsgSetOnline{}
	_ sdk.Msg                            = &MsgSetOffline{}
	_ sdk.Msg                            = &MsgDelegate{}
	_ sdk.Msg                            = &MsgDelegateNFT{}
	_ sdk.Msg                            = &MsgRedelegate{}
	_ sdk.Msg                            = &MsgRedelegateNFT{}
	_ sdk.Msg                            = &MsgUndelegate{}
	_ sdk.Msg                            = &MsgUndelegateNFT{}
	_ sdk.Msg                            = &MsgCancelRedelegation{}
	_ sdk.Msg                            = &MsgCancelRedelegationNFT{}
	_ sdk.Msg                            = &MsgCancelUndelegation{}
	_ sdk.Msg                            = &MsgCancelUndelegationNFT{}
)

const (
	TypeMsgCreateValidator       = "create_validator"
	TypeMsgEditValidator         = "edit_validator"
	TypeMsgSetOnline             = "set_online"
	TypeMsgSetOffline            = "set_offline"
	TypeMsgDelegate              = "delegate"
	TypeMsgDelegateNFT           = "delegate_nft"
	TypeMsgRedelegate            = "redelegate"
	TypeMsgRedelegateNFT         = "redelegate_nft"
	TypeMsgUndelegate            = "undelegate"
	TypeMsgUndelegateNFT         = "undelegate_nft"
	TypeMsgCancelRedelegation    = "cancel_redelegation"
	TypeMsgCancelRedelegationNFT = "cancel_redelegation_nft"
	TypeMsgCancelUndelegation    = "cancel_undelegation"
	TypeMsgCancelUndelegationNFT = "cancel_undelegation_nft"
)

////////////////////////////////////////////////////////////////
// MsgCreateValidator
////////////////////////////////////////////////////////////////

// NewMsgCreateValidator creates a new instance of MsgCreateValidator.
func NewMsgCreateValidator(
	operatorAddr sdk.ValAddress,
	rewardAddr sdk.AccAddress,
	pubKey cryptotypes.PubKey,
	description Description,
	commission sdk.Dec,
	stake sdk.Coin,
) (*MsgCreateValidator, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}
	return &MsgCreateValidator{
		OperatorAddress: operatorAddr.String(),
		RewardAddress:   rewardAddr.String(),
		ConsensusPubkey: pkAny,
		Description:     description,
		Commission:      commission,
		Stake:           stake,
	}, nil
}

// Route should return the name of the module.
func (msg *MsgCreateValidator) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgCreateValidator) Type() string { return TypeMsgCreateValidator }

// GetSignBytes encodes the message for signing.
func (msg *MsgCreateValidator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgCreateValidator) GetSigners() []sdk.AccAddress {
	addr, err := sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{sdk.AccAddress(addr)}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgCreateValidator) ValidateBasic() error {
	_, err := sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.RewardAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid reward address: %s", err)
	}
	if msg.ConsensusPubkey == nil {
		return ErrEmptyValidatorPubKey
	}
	if len(msg.ConsensusPubkey.Value) == 0 {
		return ErrEmptyValidatorPubKey
	}
	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}
	if msg.Commission.IsNil() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty commission")
	}
	if msg.Commission.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "negative commission")
	}
	if !msg.Stake.IsValid() || !msg.Stake.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid initial delegation")
	}
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces.
func (msg MsgCreateValidator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.ConsensusPubkey, &pubKey)
}

////////////////////////////////////////////////////////////////
// MsgEditValidator
////////////////////////////////////////////////////////////////

// NewMsgEditValidator creates a new instance of MsgEditValidator.
func NewMsgEditValidator(
	operatorAddr sdk.ValAddress,
	rewardAddr sdk.AccAddress,
	description Description,
) *MsgEditValidator {
	return &MsgEditValidator{
		OperatorAddress: operatorAddr.String(),
		RewardAddress:   rewardAddr.String(),
		Description:     description,
	}
}

// Route should return the name of the module.
func (msg *MsgEditValidator) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgEditValidator) Type() string { return TypeMsgEditValidator }

// GetSignBytes encodes the message for signing.
func (msg *MsgEditValidator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgEditValidator) GetSigners() []sdk.AccAddress {
	addr, err := sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{sdk.AccAddress(addr)}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgEditValidator) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.OperatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.RewardAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", err)
	}
	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgSetOnline
////////////////////////////////////////////////////////////////

// NewMsgSetOnline creates a new instance of MsgSetOnline.
func NewMsgSetOnline(operatorAddr sdk.ValAddress) *MsgSetOnline {
	return &MsgSetOnline{
		Validator: operatorAddr.String(),
	}
}

// Route should return the name of the module.
func (msg *MsgSetOnline) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgSetOnline) Type() string { return TypeMsgSetOnline }

// GetSignBytes encodes the message for signing.
func (msg *MsgSetOnline) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgSetOnline) GetSigners() []sdk.AccAddress {
	addr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{sdk.AccAddress(addr)}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgSetOnline) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgSetOffline
////////////////////////////////////////////////////////////////

// NewMsgSetOffline creates a new instance of MsgSetOffline.
func NewMsgSetOffline(operatorAddr sdk.ValAddress) *MsgSetOffline {
	return &MsgSetOffline{
		Validator: operatorAddr.String(),
	}
}

// Route should return the name of the module.
func (msg *MsgSetOffline) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgSetOffline) Type() string { return TypeMsgSetOffline }

// GetSignBytes encodes the message for signing.
func (msg *MsgSetOffline) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgSetOffline) GetSigners() []sdk.AccAddress {
	addr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{sdk.AccAddress(addr)}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgSetOffline) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgDelegate
////////////////////////////////////////////////////////////////

// NewMsgDelegate creates a new instance of MsgDelegate.
func NewMsgDelegate(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, coin sdk.Coin) *MsgDelegate {
	return &MsgDelegate{
		Delegator: delegatorAddr.String(),
		Validator: validatorAddr.String(),
		Coin:      coin,
	}
}

// Route should return the name of the module.
func (msg *MsgDelegate) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgDelegate) Type() string { return TypeMsgDelegate }

// GetSignBytes encodes the message for signing.
func (msg MsgDelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgDelegate) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgDelegate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Delegator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	if !msg.Coin.IsValid() || !msg.Coin.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid delegation amount")
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgDelegateNFT
////////////////////////////////////////////////////////////////

// NewMsgDelegateNFT creates a new instance of MsgDelegateNFT.
func NewMsgDelegateNFT(
	delegatorAddr sdk.AccAddress,
	validatorAddr sdk.ValAddress,
	tokenID string,
	subTokenIDs []uint32,
) *MsgDelegateNFT {
	return &MsgDelegateNFT{
		Delegator:   delegatorAddr.String(),
		Validator:   validatorAddr.String(),
		TokenID:     tokenID,
		SubTokenIDs: subTokenIDs,
	}
}

// Route should return the name of the module.
func (msg *MsgDelegateNFT) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgDelegateNFT) Type() string { return TypeMsgDelegateNFT }

// GetSignBytes encodes the message for signing.
func (msg *MsgDelegateNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgDelegateNFT) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgDelegateNFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Delegator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	if len(msg.TokenID) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("empty token ID")
	}
	if len(msg.SubTokenIDs) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("empty sub-token IDs")
	}
	if !isUnique(msg.SubTokenIDs) {
		return errors.SubTokenIDsDublicates
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgRedelegate
////////////////////////////////////////////////////////////////

// NewMsgRedelegate creates a new instance of MsgRedelegate.
func NewMsgRedelegate(
	delegatorAddr sdk.AccAddress,
	validatorSrcAddr sdk.ValAddress,
	validatorDstAddr sdk.ValAddress,
	coin sdk.Coin,
) *MsgRedelegate {
	return &MsgRedelegate{
		Delegator:    delegatorAddr.String(),
		ValidatorSrc: validatorSrcAddr.String(),
		ValidatorDst: validatorDstAddr.String(),
		Coin:         coin,
	}
}

// Route should return the name of the module.
func (msg *MsgRedelegate) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgRedelegate) Type() string { return TypeMsgRedelegate }

// GetSignBytes encodes the message for signing.
func (msg *MsgRedelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgRedelegate) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgRedelegate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Delegator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorSrc); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid source validator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorDst); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid destination validator address: %s", err)
	}
	if msg.ValidatorSrc == msg.ValidatorDst {
		return errors.SelfRedelegation
	}
	if !msg.Coin.IsValid() || !msg.Coin.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid shares amount")
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgRedelegateNFT
////////////////////////////////////////////////////////////////

// NewMsgRedelegateNFT creates a new instance of MsgRedelegateNFT.
func NewMsgRedelegateNFT(
	delegatorAddr sdk.AccAddress,
	validatorSrcAddr sdk.ValAddress,
	validatorDstAddr sdk.ValAddress,
	tokenID string,
	subTokenIDs []uint32,
) *MsgRedelegateNFT {
	return &MsgRedelegateNFT{
		Delegator:    delegatorAddr.String(),
		ValidatorSrc: validatorSrcAddr.String(),
		ValidatorDst: validatorDstAddr.String(),
		TokenID:      tokenID,
		SubTokenIDs:  subTokenIDs,
	}
}

// Route should return the name of the module.
func (msg *MsgRedelegateNFT) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgRedelegateNFT) Type() string { return TypeMsgRedelegateNFT }

// GetSignBytes encodes the message for signing.
func (msg *MsgRedelegateNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgRedelegateNFT) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgRedelegateNFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Delegator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorSrc); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid source validator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorDst); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid destination validator address: %s", err)
	}
	if msg.ValidatorSrc == msg.ValidatorDst {
		return errors.SelfRedelegation
	}
	if len(msg.TokenID) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("empty token ID")
	}
	if len(msg.SubTokenIDs) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("empty sub-token IDs")
	}
	if !isUnique(msg.SubTokenIDs) {
		return errors.SubTokenIDsDublicates
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgUndelegate
////////////////////////////////////////////////////////////////

// NewMsgUndelegate creates a new instance of MsgUndelegate.
func NewMsgUndelegate(
	delegatorAddr sdk.AccAddress,
	validatorAddr sdk.ValAddress,
	coin sdk.Coin,
) *MsgUndelegate {
	return &MsgUndelegate{
		Delegator: delegatorAddr.String(),
		Validator: validatorAddr.String(),
		Coin:      coin,
	}
}

// Route should return the name of the module.
func (msg *MsgUndelegate) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgUndelegate) Type() string { return TypeMsgUndelegate }

// GetSignBytes encodes the message for signing.
func (msg *MsgUndelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgUndelegate) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgUndelegate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Delegator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	if !msg.Coin.IsValid() || !msg.Coin.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid shares amount")
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgUndelegateNFT
////////////////////////////////////////////////////////////////

// NewMsgUndelegateNFT creates a new instance of MsgUndelegateNFT.
func NewMsgUndelegateNFT(
	delegatorAddr sdk.AccAddress,
	validatorAddr sdk.ValAddress,
	tokenID string,
	subTokenIDs []uint32,
) *MsgUndelegateNFT {
	return &MsgUndelegateNFT{
		Delegator:   delegatorAddr.String(),
		Validator:   validatorAddr.String(),
		TokenID:     tokenID,
		SubTokenIDs: subTokenIDs,
	}
}

// Route should return the name of the module.
func (msg *MsgUndelegateNFT) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgUndelegateNFT) Type() string { return TypeMsgUndelegateNFT }

// GetSignBytes encodes the message for signing.
func (msg *MsgUndelegateNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgUndelegateNFT) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgUndelegateNFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Delegator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	if len(msg.TokenID) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("empty token ID")
	}
	if len(msg.SubTokenIDs) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("empty sub-token IDs")
	}
	if !isUnique(msg.SubTokenIDs) {
		return errors.SubTokenIDsDublicates
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgCancelRedelegation
////////////////////////////////////////////////////////////////

// NewMsgCancelRedelegation creates a new instance of MsgCancelRedelegation.
func NewMsgCancelRedelegation(
	delegatorAddr sdk.AccAddress,
	validatorSrcAddr sdk.ValAddress,
	validatorDstAddr sdk.ValAddress,
	creationHeight int64,
	coin sdk.Coin,
) *MsgCancelRedelegation {
	return &MsgCancelRedelegation{
		Delegator:      delegatorAddr.String(),
		ValidatorSrc:   validatorSrcAddr.String(),
		ValidatorDst:   validatorDstAddr.String(),
		CreationHeight: creationHeight,
		Coin:           coin,
	}
}

// Route should return the name of the module.
func (msg *MsgCancelRedelegation) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgCancelRedelegation) Type() string { return TypeMsgCancelRedelegation }

// GetSignBytes encodes the message for signing.
func (msg *MsgCancelRedelegation) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgCancelRedelegation) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgCancelRedelegation) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Delegator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorSrc); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid source validator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorDst); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid destination validator address: %s", err)
	}
	if msg.CreationHeight <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid creation height")
	}
	if !msg.Coin.IsValid() || !msg.Coin.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid shares amount")
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgCancelRedelegationNFT
////////////////////////////////////////////////////////////////

// NewMsgCancelRedelegationNFT creates a new instance of MsgCancelRedelegationNFT.
func NewMsgCancelRedelegationNFT(
	delegatorAddr sdk.AccAddress,
	validatorSrcAddr sdk.ValAddress,
	validatorDstAddr sdk.ValAddress,
	creationHeight int64,
	tokenID string,
	subTokenIDs []uint32,
) *MsgCancelRedelegationNFT {
	return &MsgCancelRedelegationNFT{
		Delegator:      delegatorAddr.String(),
		ValidatorSrc:   validatorSrcAddr.String(),
		ValidatorDst:   validatorDstAddr.String(),
		CreationHeight: creationHeight,
		TokenID:        tokenID,
		SubTokenIDs:    subTokenIDs,
	}
}

// Route should return the name of the module.
func (msg *MsgCancelRedelegationNFT) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgCancelRedelegationNFT) Type() string { return TypeMsgCancelRedelegationNFT }

// GetSignBytes encodes the message for signing.
func (msg *MsgCancelRedelegationNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgCancelRedelegationNFT) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgCancelRedelegationNFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Delegator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorSrc); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid source validator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorDst); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid destination validator address: %s", err)
	}
	if msg.CreationHeight <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid creation height")
	}
	if len(msg.TokenID) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("empty token ID")
	}
	if len(msg.SubTokenIDs) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("empty sub-token IDs")
	}
	if !isUnique(msg.SubTokenIDs) {
		return errors.SubTokenIDsDublicates
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgCancelUndelegation
////////////////////////////////////////////////////////////////

// NewMsgCancelUndelegation creates a new instance of MsgCancelUndelegation.
func NewMsgCancelUndelegation(
	delegatorAddr sdk.AccAddress,
	validatorAddr sdk.ValAddress,
	creationHeight int64,
	coin sdk.Coin,
) *MsgCancelUndelegation {
	return &MsgCancelUndelegation{
		Delegator:      delegatorAddr.String(),
		Validator:      validatorAddr.String(),
		CreationHeight: creationHeight,
		Coin:           coin,
	}
}

// Route should return the name of the module.
func (msg *MsgCancelUndelegation) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgCancelUndelegation) Type() string { return TypeMsgCancelUndelegation }

// GetSignBytes encodes the message for signing.
func (msg *MsgCancelUndelegation) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgCancelUndelegation) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgCancelUndelegation) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Delegator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	if msg.CreationHeight <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid creation height")
	}
	if !msg.Coin.IsValid() || !msg.Coin.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid shares amount")
	}
	return nil
}

////////////////////////////////////////////////////////////////
// MsgCancelUndelegationNFT
////////////////////////////////////////////////////////////////

// NewMsgCancelUndelegationNFT creates a new instance of MsgCancelUndelegationNFT.
func NewMsgCancelUndelegationNFT(
	delegatorAddr sdk.AccAddress,
	validatorAddr sdk.ValAddress,
	creationHeight int64,
	tokenID string,
	subTokenIDs []uint32,
) *MsgCancelUndelegationNFT {
	return &MsgCancelUndelegationNFT{
		Delegator:      delegatorAddr.String(),
		Validator:      validatorAddr.String(),
		CreationHeight: creationHeight,
		TokenID:        tokenID,
		SubTokenIDs:    subTokenIDs,
	}
}

// Route should return the name of the module.
func (msg *MsgCancelUndelegationNFT) Route() string { return RouterKey }

// Type should return the action.
func (msg *MsgCancelUndelegationNFT) Type() string { return TypeMsgCancelUndelegationNFT }

// GetSignBytes encodes the message for signing.
func (msg *MsgCancelUndelegationNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg *MsgCancelUndelegationNFT) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs stateless checks on the message.
func (msg *MsgCancelUndelegationNFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Delegator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	if msg.CreationHeight <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid creation height")
	}
	if len(msg.TokenID) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("empty token ID")
	}
	if len(msg.SubTokenIDs) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("empty sub-token IDs")
	}
	if !isUnique(msg.SubTokenIDs) {
		return errors.SubTokenIDsDublicates
	}
	return nil
}

// returns true if all list elements apperas only one time
func isUnique(list []uint32) bool {
	var looked = make(map[uint32]bool)
	for _, id := range list {
		if looked[id] {
			return false
		}
		looked[id] = true
	}
	return true
}
