package keeper

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the validator MsgServer interface for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateValidator defines a method for creating a new validator.
func (k msgServer) CreateValidator(goCtx context.Context, msg *types.MsgCreateValidator) (*types.MsgCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}
	rewardAddr, err := sdk.AccAddressFromBech32(msg.RewardAddress)
	if err != nil {
		return nil, err
	}

	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetValidator(ctx, valAddr); found {
		return nil, errors.ValidatorAlreadyExists
	}

	_, err = k.coinKeeper.GetCoin(ctx, msg.Stake.Denom)
	if err != nil {
		return nil, err
	}

	pk, ok := msg.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, errors.InvalidConsensusPubKey
	}

	if _, found := k.GetValidatorByConsAddrDecimal(ctx, sdk.GetConsAddress(pk)); found {
		return nil, errors.ValidatorPublicKeyAlreadyExists
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return nil, err
	}

	cp := ctx.ConsensusParams()
	if cp != nil && cp.Validator != nil {
		pkType := pk.Type()
		hasKeyType := false
		for _, keyType := range cp.Validator.PubKeyTypes {
			if pkType == keyType {
				hasKeyType = true
				break
			}
		}
		if !hasKeyType {
			return nil, errors.UnsupportedPubKeyType
		}
	}

	validator, err := types.NewValidator(valAddr, rewardAddr, pk, msg.Description, msg.Commission)
	if err != nil {
		return nil, err
	}
	validator.Online = false
	validator.Jailed = false

	k.SetValidator(ctx, validator)
	k.SetValidatorByConsAddr(ctx, validator)
	// TODO: calculate power
	k.SetNewValidatorByPowerIndex(ctx, validator.GetOperator(), 0)

	// call the after-creation hook
	if err := k.AfterValidatorCreated(ctx, validator.GetOperator()); err != nil {
		return nil, err
	}

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	// NOTE source will always be from a wallet which are unbonded
	_, _, err = k.Keeper.Delegate(ctx, sdk.AccAddress(valAddr), msg.Stake.Denom,
		&msg.Stake.Amount, nil, validator)
	if err != nil {
		return nil, err
	}

	err = events.EmitTypedEvent(ctx, &types.EventCreateValidator{
		Sender:          sdk.AccAddress(valAddr).String(),
		Validator:       valAddr.String(),
		RewardAddress:   rewardAddr.String(),
		ConsensusPubkey: pk.String(),
		Description:     msg.Description,
		Commission:      msg.Commission,
		Stake:           msg.Stake,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgCreateValidatorResponse{}, nil
}

// EditValidator defines a method for editing an existing validator.
func (k msgServer) EditValidator(goCtx context.Context, msg *types.MsgEditValidator) (*types.MsgEditValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}
	// validator must already be registered
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, types.ErrNoValidatorFound
	}

	// replace all editable fields (clients should autofill existing values)
	description, err := validator.Description.UpdateDescription(msg.Description)
	if err != nil {
		return nil, err
	}

	validator.Description = description
	validator.RewardAddress = msg.RewardAddress

	k.SetValidator(ctx, validator)

	err = events.EmitTypedEvent(ctx, &types.EventEditValidator{
		Sender:        sdk.AccAddress(valAddr).String(),
		Validator:     valAddr.String(),
		RewardAddress: msg.RewardAddress,
		Description:   description,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgEditValidatorResponse{}, nil
}

// SetOnline defines a method for turning on a validator into the blockchain consensus.
func (k msgServer) SetOnline(goCtx context.Context, msg *types.MsgSetOnline) (*types.MsgSetOnlineResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	// validator must already be registered
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, errors.ValidatorNotFound
	}

	if validator.Online {
		if !validator.Jailed {
			return nil, errors.ValidatorAlreadyOnline
		}
	}

	validator.Online = true
	validator.Jailed = false

	// TODO: optimize
	k.SetValidator(ctx, validator)

	err = events.EmitTypedEvent(ctx, &types.EventSetOnline{
		Sender:    sdk.AccAddress(valAddr).String(),
		Validator: valAddr.String(),
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgSetOnlineResponse{}, nil
}

// SetOffline defines a method for turning off a validator from the blockchain consensus.
func (k msgServer) SetOffline(goCtx context.Context, msg *types.MsgSetOffline) (*types.MsgSetOfflineResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	// validator must already be registered
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, errors.ValidatorNotFound
	}
	if !validator.Online {
		return nil, errors.ValidatorAlreadyOffline
	}

	validator.Online = false
	// TODO: optimize
	k.SetValidator(ctx, validator)

	err = events.EmitTypedEvent(ctx, &types.EventSetOffline{
		Sender:    sdk.AccAddress(valAddr).String(),
		Validator: valAddr.String(),
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgSetOfflineResponse{}, nil
}

// Delegate defines a method for performing a delegation of coins from a delegator to a validator.
func (k msgServer) Delegate(goCtx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	// validator must already be registered
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, errors.ValidatorNotFound
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil, err
	}

	if !k.coinKeeper.IsCoinExists(ctx, msg.Coin.Denom) {
		return nil, err
	}

	_, stake, err := k.Keeper.Delegate(ctx, delegatorAddress, msg.Coin.Denom, &msg.Coin.Amount, nil, validator)
	if err != nil {
		return nil, err
	}

	baseCoin := k.Keeper.ToBaseCoin(ctx, msg.Coin)

	err = events.EmitTypedEvent(ctx, &types.EventDelegate{
		Delegator:  msg.Delegator,
		Validator:  msg.Validator,
		Stake:      stake,
		AmountBase: baseCoin.Amount,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgDelegateResponse{}, nil
}

// DelegateNFT defines a method for performing a delegation of NFTs from a delegator to a validator.
func (k msgServer) DelegateNFT(goCtx context.Context, msg *types.MsgDelegateNFT) (*types.MsgDelegateNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	// validator must already be registered
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, errors.ValidatorNotFound
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil, err
	}

	for _, subID := range msg.SubTokenIDs {
		subToken, found := k.nftKeeper.GetSubToken(ctx, msg.TokenID, subID)
		if !found {
			return nil, errors.NFTSubTokenNotFound
		}
		if subToken.Owner != msg.Delegator {
			return nil, errors.DelegatorIsNotOwnerOfSubtoken
		}
	}

	_, stake, err := k.Keeper.Delegate(ctx, delegatorAddress, msg.TokenID, nil, msg.SubTokenIDs, validator)
	if err != nil {
		return nil, err
	}

	baseCoin := k.Keeper.ToBaseCoin(ctx, stake.Stake)

	err = events.EmitTypedEvent(ctx, &types.EventDelegate{
		Delegator:  msg.Delegator,
		Validator:  msg.Validator,
		Stake:      stake,
		AmountBase: baseCoin.Amount,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgDelegateNFTResponse{}, nil
}

// Redelegate defines a method for performing a redelegation of coins from a source validator to destination one.
func (k msgServer) Redelegate(goCtx context.Context, msg *types.MsgRedelegate) (*types.MsgRedelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valSrcAddr, err := sdk.ValAddressFromBech32(msg.ValidatorSrc)
	if err != nil {
		return nil, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil, err
	}
	valDstAddr, err := sdk.ValAddressFromBech32(msg.ValidatorDst)
	if err != nil {
		return nil, err
	}

	stake := types.NewStakeCoin(msg.Coin)
	remainStake, err := k.CalculateUnbondStake(ctx, delegatorAddress, valSrcAddr, stake)
	if err != nil {
		return nil, err
	}

	completionTime, err := k.BeginRedelegation(
		ctx, delegatorAddress, valSrcAddr, valDstAddr, stake, remainStake,
	)
	if err != nil {
		return nil, err
	}

	baseCoin := k.ToBaseCoin(ctx, stake.Stake)

	err = events.EmitTypedEvent(ctx, &types.EventRedelegate{
		Delegator:    msg.Delegator,
		ValidatorSrc: msg.ValidatorSrc,
		ValidatorDst: msg.ValidatorDst,
		Stake:        stake,
		AmountBase:   baseCoin.Amount,
		CompleteAt:   completionTime,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgRedelegateResponse{
		CompletionTime: completionTime,
	}, nil
}

// RedelegateNFT defines a method for performing a redelegation of NFTs from a source validator to destination one.
func (k msgServer) RedelegateNFT(goCtx context.Context, msg *types.MsgRedelegateNFT) (*types.MsgRedelegateNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valSrcAddr, err := sdk.ValAddressFromBech32(msg.ValidatorSrc)
	if err != nil {
		return nil, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil, err
	}
	valDstAddr, err := sdk.ValAddressFromBech32(msg.ValidatorDst)
	if err != nil {
		return nil, err
	}

	subtokens, err := k.prepareSubTokens(ctx, msg.TokenID, msg.SubTokenIDs)
	if err != nil {
		return nil, err
	}
	stake := types.NewStakeNFT(msg.TokenID, msg.SubTokenIDs, sumSubTokens(subtokens))
	remainStake, err := k.CalculateUnbondStake(ctx, delegatorAddress, valSrcAddr, stake)
	if err != nil {
		return nil, err
	}

	completionTime, err := k.BeginRedelegation(
		ctx, delegatorAddress, valSrcAddr, valDstAddr, stake, remainStake,
	)
	if err != nil {
		return nil, err
	}

	baseCoin := k.ToBaseCoin(ctx, stake.Stake)

	err = events.EmitTypedEvent(ctx, &types.EventRedelegate{
		Delegator:    msg.Delegator,
		ValidatorSrc: msg.ValidatorSrc,
		ValidatorDst: msg.ValidatorDst,
		Stake:        stake,
		AmountBase:   baseCoin.Amount,
		CompleteAt:   completionTime,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgRedelegateNFTResponse{
		CompletionTime: completionTime,
	}, nil
}

// Undelegate defines a method for performing an undelegation of coins from a validator.
func (k msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validatorAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil, err
	}

	stake := types.NewStakeCoin(msg.Coin)
	remainStake, err := k.CalculateUnbondStake(ctx, delegatorAddress, validatorAddr, stake)
	if err != nil {
		return nil, err
	}

	completionTime, err := k.Keeper.Undelegate(ctx, delegatorAddress, validatorAddr, stake, remainStake)
	if err != nil {
		return nil, err
	}

	baseCoin := k.ToBaseCoin(ctx, stake.Stake)

	err = events.EmitTypedEvent(ctx, &types.EventUndelegate{
		Delegator:  msg.Delegator,
		Validator:  msg.Validator,
		Stake:      stake,
		AmountBase: baseCoin.Amount,
		CompleteAt: completionTime,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgUndelegateResponse{
		CompletionTime: completionTime,
	}, nil
}

// UndelegateNFT defines a method for performing an undelegation of NFTs from a validator.
func (k msgServer) UndelegateNFT(goCtx context.Context, msg *types.MsgUndelegateNFT) (*types.MsgUndelegateNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validatorAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil, err
	}

	subtokens, err := k.prepareSubTokens(ctx, msg.TokenID, msg.SubTokenIDs)
	if err != nil {
		return nil, err
	}
	stake := types.NewStakeNFT(msg.TokenID, msg.SubTokenIDs, sumSubTokens(subtokens))
	remainStake, err := k.CalculateUnbondStake(ctx, delegatorAddress, validatorAddr, stake)
	if err != nil {
		return nil, err
	}

	completionTime, err := k.Keeper.Undelegate(ctx, delegatorAddress, validatorAddr, stake, remainStake)
	if err != nil {
		return nil, err
	}

	baseCoin := k.ToBaseCoin(ctx, stake.Stake)

	err = events.EmitTypedEvent(ctx, &types.EventUndelegate{
		Delegator:  msg.Delegator,
		Validator:  msg.Validator,
		Stake:      stake,
		AmountBase: baseCoin.Amount,
		CompleteAt: completionTime,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgUndelegateNFTResponse{
		CompletionTime: completionTime,
	}, nil
}

// CancelRedelegation defines a method for canceling the redelegation and delegate back the validator.
func (k msgServer) CancelUndelegation(goCtx context.Context, msg *types.MsgCancelUndelegation) (*types.MsgCancelUndelegationResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	//
	//valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	//if err != nil {
	//	return nil, err
	//}
	//
	//delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	//if err != nil {
	//	return nil, err
	//}
	//
	//bondDenom := k.BondDenom(ctx)
	//if msg.Amount.Denom != bondDenom {
	//	return nil, sdkerrors.Wrapf(
	//		sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
	//	)
	//}
	//
	//validator, found := k.GetValidator(ctx, valAddr)
	//if !found {
	//	return nil, types.ErrNoValidatorFound
	//}
	//
	//if validator.IsJailed() {
	//	return nil, types.ErrValidatorJailed
	//}
	//
	//ubd, found := k.GetUndelegation(ctx, delegatorAddress, valAddr)
	//if !found {
	//	return nil, status.Errorf(
	//		codes.NotFound,
	//		"unbonding delegation with delegator %s not found for validator %s",
	//		msg.DelegatorAddress, msg.ValidatorAddress,
	//	)
	//}
	//
	//var (
	//	unbondEntry      types.UndelegationEntry
	//	unbondEntryIndex int64 = -1
	//)
	//
	//for i, entry := range ubd.Entries {
	//	if entry.CreationHeight == msg.CreationHeight {
	//		unbondEntry = entry
	//		unbondEntryIndex = int64(i)
	//		break
	//	}
	//}
	//if unbondEntryIndex == -1 {
	//	return nil, sdkerrors.ErrNotFound.Wrapf("unbonding delegation entry is not found at block height %d", msg.CreationHeight)
	//}
	//
	//if unbondEntry.Balance.LT(msg.Amount.Amount) {
	//	return nil, sdkerrors.ErrInvalidRequest.Wrap("amount is greater than the unbonding delegation entry balance")
	//}
	//
	//if unbondEntry.CompletionTime.Before(ctx.BlockTime()) {
	//	return nil, sdkerrors.ErrInvalidRequest.Wrap("unbonding delegation is already processed")
	//}
	//
	//// delegate back the unbonding delegation amount to the validator
	//_, err = k.Keeper.Delegate(ctx, delegatorAddress, msg.Amount.Amount, types.Unbonding, validator, false)
	//if err != nil {
	//	return nil, err
	//}
	//
	//amount := unbondEntry.Balance.Sub(msg.Amount.Amount)
	//if amount.IsZero() {
	//	ubd.RemoveEntry(unbondEntryIndex)
	//} else {
	//	// update the undelegationEntryBalance and InitialBalance for ubd entry
	//	unbondEntry.Balance = amount
	//	unbondEntry.InitialBalance = unbondEntry.InitialBalance.Sub(msg.Amount.Amount)
	//	ubd.Entries[unbondEntryIndex] = unbondEntry
	//}
	//
	//// set the unbonding delegation or remove it if there are no more entries
	//if len(ubd.Entries) == 0 {
	//	k.RemoveUndelegation(ctx, ubd)
	//} else {
	//	k.SetUndelegation(ctx, ubd)
	//}
	//
	//ctx.EventManager().EmitEvent(
	//	sdk.NewEvent(
	//		types.EventTypeCancelUndelegation,
	//		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
	//		sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
	//		sdk.NewAttribute(types.AttributeKeyDelegator, msg.DelegatorAddress),
	//		sdk.NewAttribute(types.AttributeKeyCreationHeight, strconv.FormatInt(msg.CreationHeight, 10)),
	//	),
	//)

	return &types.MsgCancelUndelegationResponse{}, nil
}

// CancelRedelegationNFT defines a method for canceling the redelegation and delegate back the validator.
func (k msgServer) CancelUndelegationNFT(goCtx context.Context, msg *types.MsgCancelUndelegationNFT) (*types.MsgCancelUndelegationNFTResponse, error) {

	// TODO: Implement!

	return &types.MsgCancelUndelegationNFTResponse{}, nil
}

// CancelUndelegation defines a method for canceling the undelegation and delegate back to the validator.
func (k msgServer) CancelRedelegation(goCtx context.Context, msg *types.MsgCancelRedelegation) (*types.MsgCancelRedelegationResponse, error) {

	// TODO: Implement!

	return &types.MsgCancelRedelegationResponse{}, nil
}

// CancelUndelegationNFT defines a method for canceling the undelegation and delegate back to the validator.
func (k msgServer) CancelRedelegationNFT(goCtx context.Context, msg *types.MsgCancelRedelegationNFT) (*types.MsgCancelRedelegationNFTResponse, error) {

	// TODO: Implement!

	return &types.MsgCancelRedelegationNFTResponse{}, nil
}
