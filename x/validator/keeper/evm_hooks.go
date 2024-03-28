// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	types2 "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	"encoding/json"
	"fmt"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var _ evmtypes.EvmHooks = Hooks{}

// Hooks wrapper struct for erc20 keeper
type Hooks struct {
	k Keeper
}

// Hooks Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// PostTxProcessing is a wrapper for calling the EVM PostTxProcessing hook on
// the module keeper
func (h Hooks) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	return h.k.PostTxProcessing(ctx, msg, receipt)
}

// PostTxProcessing implements EvmHooks.PostTxProcessing. The EVM hooks allows
// users to convert ERC20s to Cosmos Coins by sending an Ethereum tx transfer to
// the module account address. This hook applies to both token pairs that have
// been registered through a native Cosmos coin or an ERC20 token. If token pair
// has been registered with:
//   - coin -> burn tokens and transfer escrowed coins on module to sender
//   - token -> escrow tokens on module account and mint & transfer coins to sender
//
// Note that the PostTxProcessing hook is only called by sending an EVM
// transaction that triggers `ApplyTransaction`. A cosmos tx with a
// `ConvertERC20` msg does not trigger the hook as it only calls `ApplyMessage`.
func (k Keeper) PostTxProcessing(
	ctx sdk.Context,
	msg core.Message,
	recipient *ethtypes.Receipt,
) error {
	//params := k.GetParams(ctx)
	//if !params.EnableErc20 || !params.EnableEVMHook {
	//	// no error is returned to avoid reverting the tx and allow for other post
	//	// processing txs to pass and
	//	return nil
	//}

	validatorMaster, _ := contracts.MasterValidatorMetaData.GetAbi()
	delegatorCenter, _ := contracts.DelegationMetaData.GetAbi()

	// this var is only for new token create from token center
	var tokenDelegate contracts.ContractStaked
	var tokenUndelegate contracts.ContractsRequestWithdraw
	var tokenRedelegation contracts.ContractsRequestTransfer
	var newValidator contracts.MasterValidatorValidatorAdded
	var updateValidator contracts.MasterValidatorValidatorUpdated

	for _, log := range recipient.Logs {
		eventValidatorByID, errEvent := validatorMaster.EventByID(log.Topics[0])
		if errEvent == nil {
			if eventValidatorByID.Name == "ValidatorAdded" {
				_ = validatorMaster.UnpackIntoInterface(&newValidator, eventValidatorByID.Name, log.Data)
				newValidator.Validator = common.BytesToAddress(log.Topics[1].Bytes())
				var validatorInfo contracts.MasterValidatorValidatorAddedMeta
				_ = json.Unmarshal([]byte(newValidator.Meta), &validatorInfo)
				//if err != nil {
				//	return err
				//}
				fmt.Println(validatorInfo)
			}
			if eventValidatorByID.Name == "ValidatorUpdated" {
				cosmosAddressValidator, _ := types2.GetDecimalAddressFromHex(common.BytesToAddress(log.Topics[1].Bytes()).String())
				if updateValidator.Status == 2 {
					err := k.SetOnlineFromEvm(ctx, cosmosAddressValidator.String())
					if err != nil {
						return err
					}
				}
				if updateValidator.Status == 1 {
					err := k.SetOfflineFromEvm(ctx, cosmosAddressValidator.String())
					if err != nil {
						return err
					}
				}
			}
		}
		eventDelegationByID, errEvent := delegatorCenter.EventByID(log.Topics[0])
		if errEvent == nil {
			if eventDelegationByID.Name == "Staked" {
				_ = delegatorCenter.UnpackIntoInterface(&tokenDelegate, eventDelegationByID.Name, log.Data)
				fmt.Println(tokenDelegate)
				err := k.Staked(ctx, tokenDelegate)
				if err != nil {
					return err
				}
			}

			if eventDelegationByID.Name == "RequestWithdraw1" {
				_ = delegatorCenter.UnpackIntoInterface(&tokenUndelegate, eventDelegationByID.Name, log.Data)
				fmt.Println(tokenUndelegate)
				err := k.RequestWithdraw(ctx, tokenUndelegate)
				if err != nil {
					return err
				}
			}
			if eventDelegationByID.Name == "RequestTransfer1" {
				_ = delegatorCenter.UnpackIntoInterface(&tokenRedelegation, eventDelegationByID.Name, log.Data)
				fmt.Println(tokenRedelegation)
				err := k.RequestTransfer(ctx, tokenRedelegation)
				if err != nil {
					return err
				}
			}
		}
	}

	// Check if processed method
	//switch methodId.Name {
	//case types.ContractMethodCreateValidator:
	//
	//
	//default:
	//	return nil
	//}

	return nil
}

func (k Keeper) Staked(ctx sdk.Context, stakeData contracts.ContractStaked) error {
	coinStake, err := k.coinKeeper.GetCoinByDRC(ctx, stakeData.Stake.Token.String())
	if err != nil {
		return errors.CoinDoesNotExist
	}

	stake := types.NewStakeCoin(sdk.Coin{Denom: coinStake.Denom, Amount: math.NewIntFromBigInt(stakeData.Stake.Amount)})

	cosmosAddress, _ := types2.GetDecimalAddressFromHex(stakeData.Stake.Delegator.String())

	mintCoinForDelegation := sdk.NewCoins(sdk.NewCoin(coinStake.Denom, math.NewIntFromBigInt(stakeData.Stake.Amount)))
	err = k.bankKeeper.MintCoins(ctx, cointypes.ModuleName, mintCoinForDelegation)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, cointypes.ModuleName, cosmosAddress, mintCoinForDelegation)
	if err != nil {
		return err
	}

	//cosmosAddressValidator, _ := types2.GetDecimalAddressFromHex(tokenStaked.Stake.Validator.String())

	valAddr, err := sdk.ValAddressFromBech32("d0valoper1t4qx5x570wglgesc5g5gvf3a0n3jf9ngsn76pl")

	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return fmt.Errorf("not found validator %s", valAddr)
	}

	_ = k.Delegate(ctx, cosmosAddress, validator, stake)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) RequestWithdraw(ctx sdk.Context, tokenUndelegate contracts.ContractsRequestWithdraw) error {
	coinStake, err := k.coinKeeper.GetCoinByDRC(ctx, tokenUndelegate.FrozenStake.Stake.Token.String())
	if err != nil {
		return errors.CoinDoesNotExist
	}

	stake := types.NewStakeCoin(sdk.Coin{Denom: coinStake.Denom, Amount: math.NewIntFromBigInt(tokenUndelegate.FrozenStake.Stake.Amount)})

	cosmosAddress, _ := types2.GetDecimalAddressFromHex(tokenUndelegate.FrozenStake.Stake.Delegator.String())

	//cosmosAddressValidator, _ := types2.GetDecimalAddressFromHex(tokenStaked.Stake.Validator.String())

	valAddr, err := sdk.ValAddressFromBech32("d0valoper1t4qx5x570wglgesc5g5gvf3a0n3jf9ngsn76pl")

	delegation, found := k.GetDelegation(ctx, cosmosAddress, valAddr, stake.ID)
	if !found {
		return errors.DelegationNotFound
	}

	remainStake, err := k.CalculateRemainStake(ctx, delegation.Stake, stake)
	if err != nil {
		return err
	}

	_, err = k.Undelegate(ctx, cosmosAddress, valAddr, stake, remainStake)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) RequestTransfer(ctx sdk.Context, tokenRedelegation contracts.ContractsRequestTransfer) error {
	coinStake, err := k.coinKeeper.GetCoinByDRC(ctx, tokenRedelegation.FrozenStake.Stake.Token.String())
	if err != nil {
		return errors.CoinDoesNotExist
	}

	stake := types.NewStakeCoin(sdk.Coin{Denom: coinStake.Denom, Amount: math.NewIntFromBigInt(tokenRedelegation.FrozenStake.Stake.Amount)})

	cosmosAddress, _ := types2.GetDecimalAddressFromHex(tokenRedelegation.FrozenStake.Stake.Delegator.String())

	//cosmosAddressValidator, _ := types2.GetDecimalAddressFromHex(tokenStaked.Stake.Validator.String())

	valAddr, err := sdk.ValAddressFromBech32("d0valoper1t4qx5x570wglgesc5g5gvf3a0n3jf9ngsn76pl")

	delegation, found := k.GetDelegation(ctx, cosmosAddress, valAddr, stake.ID)
	if !found {
		return errors.DelegationNotFound
	}

	remainStake, err := k.CalculateRemainStake(ctx, delegation.Stake, stake)
	if err != nil {
		return err
	}

	_, err = k.BeginRedelegation(
		ctx, cosmosAddress, valAddr, valAddr, stake, remainStake,
	)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) CreateValidatorFromEVM(ctx sdk.Context, meta string) error {

	msg := types.MsgCreateValidator{
		OperatorAddress: DAOAddress2,
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return err
	}
	rewardAddr, err := sdk.AccAddressFromBech32(msg.RewardAddress)
	if err != nil {
		return err
	}

	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetValidator(ctx, valAddr); found {
		return errors.ValidatorAlreadyExists
	}

	_, err = k.coinKeeper.GetCoin(ctx, msg.Stake.Denom)
	if err != nil {
		return err
	}

	pk, ok := msg.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return errors.InvalidConsensusPubKey
	}

	if _, found := k.GetValidatorByConsAddrDecimal(ctx, sdk.GetConsAddress(pk)); found {
		return errors.ValidatorPublicKeyAlreadyExists
	}

	if _, err = msg.Description.EnsureLength(); err != nil {
		return err
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
			return errors.UnsupportedPubKeyType
		}
	}

	validator, err := types.NewValidator(valAddr, rewardAddr, pk, msg.Description, msg.Commission)
	if err != nil {
		return err
	}
	validator.Online = false
	validator.Jailed = false

	k.SetValidator(ctx, validator)
	k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)

	// call the after-creation hook
	if err = k.AfterValidatorCreated(ctx, validator.GetOperator()); err != nil {
		return err
	}

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	// NOTE source will always be from a wallet which are unbonded
	stake := types.NewStakeCoin(msg.Stake)
	err = k.Delegate(ctx, sdk.AccAddress(valAddr), validator, stake)
	if err != nil {
		return err
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
		return errors.Internal.Wrapf("err: %s", err.Error())
	}

	baseCoin := k.ToBaseCoin(ctx, msg.Stake)

	err = events.EmitTypedEvent(ctx, &types.EventDelegate{
		Delegator:  sdk.AccAddress(valAddr).String(),
		Validator:  valAddr.String(),
		Stake:      stake,
		AmountBase: baseCoin.Amount,
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}
	return nil
}

// SetOnlineFromEvm defines a method for turning on a validator into the blockchain consensus.
func (k Keeper) SetOnlineFromEvm(goCtx sdk.Context, validatorAddr string) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return err
	}
	// validator must already be registered
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return errors.ValidatorNotFound
	}

	if validator.Online {
		if !validator.Jailed {
			return nil
		}
	}

	// validator without delegations can't become online
	if !k.HasDelegations(ctx, valAddr) {
		return errors.ValidatorHasNoDelegations
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
		return err
	}

	stake := TokensToConsensusPower(totalStake)
	if stake == 0 {
		return errors.ValidatorStakeTooSmall
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
		return err
	}
	k.SetStartHeight(ctx, consAdr, ctx.BlockHeight())

	err = events.EmitTypedEvent(ctx, &types.EventSetOnline{
		Sender:    sdk.AccAddress(valAddr).String(),
		Validator: valAddr.String(),
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}

	return nil
}

// SetOfflineFromEvm defines a method for turning on a validator into the blockchain consensus.
func (k Keeper) SetOfflineFromEvm(goCtx sdk.Context, validatorAddrHex string) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(validatorAddrHex)
	if err != nil {
		return err
	}
	// validator must already be registered
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return errors.ValidatorNotFound
	}
	if !validator.Online {
		return errors.ValidatorAlreadyOffline
	}

	validator.Online = false
	// TODO: optimize
	k.SetValidator(ctx, validator)

	consAdr, err := validator.GetConsAddr()
	if err != nil {
		return err
	}
	k.DeleteStartHeight(ctx, consAdr)

	err = events.EmitTypedEvent(ctx, &types.EventSetOffline{
		Sender:    sdk.AccAddress(valAddr).String(),
		Validator: valAddr.String(),
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}

	return nil
}
