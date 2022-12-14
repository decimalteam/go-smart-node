package keeper

import (
	"context"
	"time"

	sdkmath "cosmossdk.io/math"
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
	k.SetNewValidatorByPowerIndex(ctx, validator)

	// call the after-creation hook
	if err := k.AfterValidatorCreated(ctx, validator.GetOperator()); err != nil {
		return nil, err
	}

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	// NOTE source will always be from a wallet which are unbonded
	stake := types.NewStakeCoin(msg.Stake)
	err = k.Keeper.Delegate(ctx, sdk.AccAddress(valAddr), validator, stake)
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

	baseCoin := k.Keeper.ToBaseCoin(ctx, msg.Stake)

	err = events.EmitTypedEvent(ctx, &types.EventDelegate{
		Delegator:  sdk.AccAddress(valAddr).String(),
		Validator:  valAddr.String(),
		Stake:      stake,
		AmountBase: baseCoin.Amount,
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
		return nil, errors.ValidatorNotFound
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

	// validator without delegations can't become online
	if !k.HasDelegations(ctx, valAddr) {
		return nil, errors.ValidatorHasNoDelegations
	}

	k.DeleteValidatorByPowerIndex(ctx, validator)

	// TODO: move Online and Jailed to store keys?
	validator.Online = true
	validator.Jailed = false

	delByValidator := k.GetAllDelegationsByValidator(ctx)
	customCoinStaked := k.GetAllCustomCoinsStaked(ctx)
	customCoinPrices := k.CalculateCustomCoinPrices(ctx, customCoinStaked)
	totalStake, err := k.CalculateTotalPowerWithDelegationsAndPrices(ctx, validator.GetOperator(), delByValidator[validator.OperatorAddress], customCoinPrices)
	if err != nil {
		return nil, err
	}

	stake := TokensToConsensusPower(totalStake)
	if stake == 0 {
		return nil, errors.ValidatorStakeTooSmall
	}

	validator.Stake = stake

	rs, err := k.GetValidatorRS(ctx, valAddr)
	if err != nil {
		rs = types.ValidatorRS{
			Rewards:      sdkmath.ZeroInt(),
			TotalRewards: sdkmath.ZeroInt(),
		}
	}
	rs.Stake = stake
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	k.SetValidatorRS(ctx, valAddr, rs)

	// StartHeight need for correct calculation of missing blocks
	consAdr, err := validator.GetConsAddr()
	if err != nil {
		return nil, err
	}
	k.SetStartHeight(ctx, consAdr, ctx.BlockHeight())

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

	consAdr, err := validator.GetConsAddr()
	if err != nil {
		return nil, err
	}
	k.DeleteStartHeight(ctx, consAdr)

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

	if !k.coinKeeper.IsCoinExists(ctx, msg.Coin.Denom) {
		return nil, errors.CoinDoesNotExist
	}

	stake := types.NewStakeCoin(msg.Coin)

	err := k._delegate(ctx, msg.Delegator, msg.Validator, stake)
	if err != nil {
		return nil, err
	}

	return &types.MsgDelegateResponse{}, nil
}

// DelegateNFT defines a method for performing a delegation of NFTs from a delegator to a validator.
func (k msgServer) DelegateNFT(goCtx context.Context, msg *types.MsgDelegateNFT) (*types.MsgDelegateNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	subtokens, err := k.prepareSubTokens(ctx, msg.TokenID, msg.SubTokenIDs)
	if err != nil {
		return nil, err
	}
	// check ownership
	for _, sub := range subtokens {
		if sub.Owner != msg.Delegator {
			return nil, errors.DelegatorIsNotOwnerOfSubtoken
		}
	}

	stake := types.NewStakeNFT(msg.TokenID, msg.SubTokenIDs, sumSubTokens(subtokens))

	err = k._delegate(ctx, msg.Delegator, msg.Validator, stake)
	if err != nil {
		return nil, err
	}

	return &types.MsgDelegateNFTResponse{}, nil
}

func (k msgServer) _delegate(ctx sdk.Context, msgDelegator string, msgValidator string, stake types.Stake) error {
	valAddr, err := sdk.ValAddressFromBech32(msgValidator)
	if err != nil {
		return err
	}
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return errors.ValidatorNotFound
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msgDelegator)
	if err != nil {
		return err
	}

	err = k.Keeper.Delegate(ctx, delegatorAddress, validator, stake)
	if err != nil {
		return err
	}

	baseCoin := k.Keeper.ToBaseCoin(ctx, stake.Stake)

	err = events.EmitTypedEvent(ctx, &types.EventDelegate{
		Delegator:  msgDelegator,
		Validator:  msgValidator,
		Stake:      stake,
		AmountBase: baseCoin.Amount,
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}
	return nil
}

// Redelegate defines a method for performing a redelegation of coins from a source validator to destination one.
func (k msgServer) Redelegate(goCtx context.Context, msg *types.MsgRedelegate) (*types.MsgRedelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	stake := types.NewStakeCoin(msg.Coin)

	completionTime, err := k._redelegate(ctx, msg.Delegator, msg.ValidatorSrc, msg.ValidatorDst, stake)
	if err != nil {
		return nil, err
	}

	return &types.MsgRedelegateResponse{
		CompletionTime: completionTime,
	}, nil
}

// RedelegateNFT defines a method for performing a redelegation of NFTs from a source validator to destination one.
func (k msgServer) RedelegateNFT(goCtx context.Context, msg *types.MsgRedelegateNFT) (*types.MsgRedelegateNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	subtokens, err := k.prepareSubTokens(ctx, msg.TokenID, msg.SubTokenIDs)
	if err != nil {
		return nil, err
	}
	stake := types.NewStakeNFT(msg.TokenID, msg.SubTokenIDs, sumSubTokens(subtokens))

	completionTime, err := k._redelegate(ctx, msg.Delegator, msg.ValidatorSrc, msg.ValidatorDst, stake)
	if err != nil {
		return nil, err
	}

	return &types.MsgRedelegateNFTResponse{
		CompletionTime: completionTime,
	}, nil
}

func (k msgServer) _redelegate(ctx sdk.Context, msgDelegator, msgValidatorSrc, msgValidatorDst string, stake types.Stake) (time.Time, error) {
	valSrcAddr, err := sdk.ValAddressFromBech32(msgValidatorSrc)
	if err != nil {
		return time.Time{}, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msgDelegator)
	if err != nil {
		return time.Time{}, err
	}
	valDstAddr, err := sdk.ValAddressFromBech32(msgValidatorDst)
	if err != nil {
		return time.Time{}, err
	}

	delegation, found := k.GetDelegation(ctx, delegatorAddress, valSrcAddr, stake.ID)
	if !found {
		return time.Time{}, errors.DelegationNotFound
	}

	remainStake, err := k.CalculateRemainStake(ctx, delegation.Stake, stake)
	if err != nil {
		return time.Time{}, err
	}

	completionTime, err := k.BeginRedelegation(
		ctx, delegatorAddress, valSrcAddr, valDstAddr, stake, remainStake,
	)
	if err != nil {
		return time.Time{}, err
	}

	baseCoin := k.ToBaseCoin(ctx, stake.Stake)

	err = events.EmitTypedEvent(ctx, &types.EventRedelegate{
		Delegator:    msgDelegator,
		ValidatorSrc: msgValidatorSrc,
		ValidatorDst: msgValidatorDst,
		Stake:        stake,
		AmountBase:   baseCoin.Amount,
		CompleteAt:   completionTime,
	})
	if err != nil {
		return time.Time{}, errors.Internal.Wrapf("err: %s", err.Error())
	}
	return completionTime, nil
}

// Undelegate defines a method for performing an undelegation of coins from a validator.
func (k msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	stake := types.NewStakeCoin(msg.Coin)

	completionTime, err := k._undelegate(ctx, msg.Delegator, msg.Validator, stake)
	if err != nil {
		return nil, err
	}

	return &types.MsgUndelegateResponse{
		CompletionTime: completionTime,
	}, nil
}

// UndelegateNFT defines a method for performing an undelegation of NFTs from a validator.
func (k msgServer) UndelegateNFT(goCtx context.Context, msg *types.MsgUndelegateNFT) (*types.MsgUndelegateNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	subtokens, err := k.prepareSubTokens(ctx, msg.TokenID, msg.SubTokenIDs)
	if err != nil {
		return nil, err
	}
	stake := types.NewStakeNFT(msg.TokenID, msg.SubTokenIDs, sumSubTokens(subtokens))

	completionTime, err := k._undelegate(ctx, msg.Delegator, msg.Validator, stake)
	if err != nil {
		return nil, err
	}

	return &types.MsgUndelegateNFTResponse{
		CompletionTime: completionTime,
	}, nil
}

func (k msgServer) _undelegate(ctx sdk.Context, msgDelegator string, msgValidator string, stake types.Stake) (time.Time, error) {

	validatorAddr, err := sdk.ValAddressFromBech32(msgValidator)
	if err != nil {
		return time.Time{}, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msgDelegator)
	if err != nil {
		return time.Time{}, err
	}

	delegation, found := k.GetDelegation(ctx, delegatorAddress, validatorAddr, stake.ID)
	if !found {
		return time.Time{}, errors.DelegationNotFound
	}

	remainStake, err := k.CalculateRemainStake(ctx, delegation.Stake, stake)
	if err != nil {
		return time.Time{}, err
	}

	completionTime, err := k.Keeper.Undelegate(ctx, delegatorAddress, validatorAddr, stake, remainStake)
	if err != nil {
		return time.Time{}, err
	}

	baseCoin := k.ToBaseCoin(ctx, stake.Stake)

	err = events.EmitTypedEvent(ctx, &types.EventUndelegate{
		Delegator:  msgDelegator,
		Validator:  msgValidator,
		Stake:      stake,
		AmountBase: baseCoin.Amount,
		CompleteAt: completionTime,
	})
	if err != nil {
		return time.Time{}, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return completionTime, nil
}

// CancelUndelegation defines a method for canceling the undelegation and delegate back the validator.
func (k msgServer) CancelUndelegation(goCtx context.Context, msg *types.MsgCancelUndelegation) (*types.MsgCancelUndelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	stake := types.NewStakeCoin(msg.Coin)

	err := k._cancelUndelegation(ctx, msg.Delegator, msg.Validator, msg.CreationHeight, stake)
	if err != nil {
		return nil, err
	}

	return &types.MsgCancelUndelegationResponse{}, nil
}

// CancelUndelegationNFT defines a method for canceling the undelegation and delegate back the validator.
func (k msgServer) CancelUndelegationNFT(goCtx context.Context, msg *types.MsgCancelUndelegationNFT) (*types.MsgCancelUndelegationNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	subtokens, err := k.prepareSubTokens(ctx, msg.TokenID, msg.SubTokenIDs)
	if err != nil {
		return nil, err
	}
	stake := types.NewStakeNFT(msg.TokenID, msg.SubTokenIDs, sumSubTokens(subtokens))

	err = k._cancelUndelegation(ctx, msg.Delegator, msg.Validator, msg.CreationHeight, stake)
	if err != nil {
		return nil, err
	}

	return &types.MsgCancelUndelegationNFTResponse{}, nil
}

func (k msgServer) _cancelUndelegation(ctx sdk.Context, msgDelegator string, msgValidator string, msgCreationHeight int64, stake types.Stake) error {
	valAddr, err := sdk.ValAddressFromBech32(msgValidator)
	if err != nil {
		return err
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msgDelegator)
	if err != nil {
		return err
	}

	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return errors.ValidatorNotFound
	}

	if validator.IsJailed() {
		return errors.ValidatorJailed
	}

	ubd, found := k.GetUndelegation(ctx, delegatorAddress, valAddr)
	if !found {
		return errors.UBDNotFound
	}

	var (
		unbondEntry      types.UndelegationEntry
		unbondEntryIndex int64 = -1
	)

	for i, entry := range ubd.Entries {
		if entry.CreationHeight == msgCreationHeight && entry.Stake.ID == stake.ID {
			unbondEntry = entry
			unbondEntryIndex = int64(i)
			break
		}
	}
	if unbondEntryIndex == -1 {
		return errors.UBDEntryNotFound
	}
	if unbondEntry.CompletionTime.Before(ctx.BlockTime()) {
		return errors.UBDAlreadyProcessed
	}

	remainStake, err := k.CalculateRemainStake(ctx, unbondEntry.Stake, stake)
	if err != nil {
		return err
	}

	// delegate back the unbonding delegation amount to the validator
	delegation, found := k.GetDelegation(ctx, delegatorAddress, valAddr, stake.ID)
	if !found {
		delegation = types.NewDelegation(delegatorAddress, valAddr, stake)
		k.IncrementDelegationsCount(ctx, valAddr)
	} else {
		delegation.Stake, err = delegation.Stake.Add(stake)
		if err != nil {
			return err
		}
	}
	k.SetDelegation(ctx, delegation)

	err = k.TransferStakeBetweenPools(ctx, types.BondStatus_Unbonded, validator.GetStatus(), stake)
	if err != nil {
		return err
	}

	if remainStake.IsEmpty() {
		ubd.RemoveEntry(unbondEntryIndex)
	} else {
		// update the undelegationEntryBalance and InitialBalance for ubd entry
		unbondEntry.Stake = remainStake
		ubd.Entries[unbondEntryIndex] = unbondEntry
	}

	// set the unbonding delegation or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUndelegation(ctx, ubd)
	} else {
		k.SetUndelegation(ctx, ubd)
	}

	err = events.EmitTypedEvent(ctx, &types.EventCancelUndelegation{
		Delegator:      msgDelegator,
		Validator:      msgValidator,
		CreationHeight: msgCreationHeight,
		Stake:          stake,
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}

	return nil
}

// CancelUndelegation defines a method for canceling the undelegation and delegate back to the validator.
func (k msgServer) CancelRedelegation(goCtx context.Context, msg *types.MsgCancelRedelegation) (*types.MsgCancelRedelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	stake := types.NewStakeCoin(msg.Coin)

	err := k._cancelRedelegation(ctx, msg.Delegator, msg.ValidatorSrc, msg.ValidatorDst, msg.CreationHeight, stake)
	if err != nil {
		return nil, err
	}

	return &types.MsgCancelRedelegationResponse{}, nil
}

// CancelUndelegationNFT defines a method for canceling the undelegation and delegate back to the validator.
func (k msgServer) CancelRedelegationNFT(goCtx context.Context, msg *types.MsgCancelRedelegationNFT) (*types.MsgCancelRedelegationNFTResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	subtokens, err := k.prepareSubTokens(ctx, msg.TokenID, msg.SubTokenIDs)
	if err != nil {
		return nil, err
	}
	stake := types.NewStakeNFT(msg.TokenID, msg.SubTokenIDs, sumSubTokens(subtokens))

	err = k._cancelRedelegation(ctx, msg.Delegator, msg.ValidatorSrc, msg.ValidatorDst, msg.CreationHeight, stake)
	if err != nil {
		return nil, err
	}

	return &types.MsgCancelRedelegationNFTResponse{}, nil
}

func (k msgServer) _cancelRedelegation(ctx sdk.Context, msgDelegator, msgValidatorSrc, msgValidatorDst string, msgCreationHeight int64, stake types.Stake) error {
	valSrcAddr, err := sdk.ValAddressFromBech32(msgValidatorSrc)
	if err != nil {
		return err
	}
	valDstAddr, err := sdk.ValAddressFromBech32(msgValidatorDst)
	if err != nil {
		return err
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msgDelegator)
	if err != nil {
		return err
	}

	validatorSrc, found := k.GetValidator(ctx, valSrcAddr)
	if !found {
		return errors.ValidatorNotFound
	}
	_, found = k.GetValidator(ctx, valDstAddr)
	if !found {
		return errors.ValidatorNotFound
	}

	if validatorSrc.IsJailed() {
		return errors.ValidatorJailed
	}

	red, found := k.GetRedelegation(ctx, delegatorAddress, valSrcAddr, valDstAddr)
	if !found {
		return errors.UBDNotFound
	}

	var (
		redEntry      types.RedelegationEntry
		redEntryIndex int64 = -1
	)

	for i, entry := range red.Entries {
		if entry.CreationHeight == msgCreationHeight && entry.Stake.ID == stake.ID {
			redEntry = entry
			redEntryIndex = int64(i)
			break
		}
	}
	if redEntryIndex == -1 {
		return errors.REDEntryNotFound
	}
	if redEntry.CompletionTime.Before(ctx.BlockTime()) {
		return errors.REDAlreadyProcessed
	}

	remainStake, err := k.CalculateRemainStake(ctx, redEntry.Stake, stake)
	if err != nil {
		return err
	}

	// delegate back the unbonding delegation amount to the validator
	delegation, found := k.GetDelegation(ctx, delegatorAddress, valSrcAddr, stake.ID)
	if !found {
		delegation = types.NewDelegation(delegatorAddress, valSrcAddr, stake)
		k.IncrementDelegationsCount(ctx, valSrcAddr)
	} else {
		delegation.Stake, err = delegation.Stake.Add(stake)
		if err != nil {
			return err
		}
	}
	k.SetDelegation(ctx, delegation)

	err = k.TransferStakeBetweenPools(ctx, types.BondStatus_Unbonded, validatorSrc.GetStatus(), stake)
	if err != nil {
		return err
	}

	if remainStake.IsEmpty() {
		red.RemoveEntry(redEntryIndex)
	} else {
		// update the undelegationEntryBalance and InitialBalance for ubd entry
		redEntry.Stake = remainStake
		red.Entries[redEntryIndex] = redEntry
	}

	// set the unbonding delegation or remove it if there are no more entries
	if len(red.Entries) == 0 {
		k.RemoveRedelegation(ctx, red)
	} else {
		k.SetRedelegation(ctx, red)
	}

	err = events.EmitTypedEvent(ctx, &types.EventCancelRedelegation{
		Delegator:      msgDelegator,
		ValidatorSrc:   msgValidatorSrc,
		ValidatorDst:   msgValidatorDst,
		CreationHeight: msgCreationHeight,
		Stake:          stake,
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}

	return nil
}
