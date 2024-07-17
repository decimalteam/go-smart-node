// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegation"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegationNft"
	"bitbucket.org/decimalteam/go-smart-node/contracts/validator"
	types2 "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	"encoding/json"
	"fmt"
	typescodec "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"strings"
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
	fmt.Println("validator hook")
	addressDelegation, _ := contracts.GetAddressFromContractCenter(ctx, k.evmKeeper, contracts.NameOfSlugForGetAddressDelegation)
	addressDelegationNft, _ := contracts.GetAddressFromContractCenter(ctx, k.evmKeeper, contracts.NameOfSlugForGetAddressDelegationNft)
	addressDelegation = strings.ToLower(addressDelegation)
	validatorMaster, _ := validator.ValidatorMetaData.GetAbi()
	delegatorCenter, _ := delegation.DelegationMetaData.GetAbi()
	delegatorNftCenter, _ := delegationNft.DelegationNftMetaData.GetAbi()

	// this var is only for new token create from token center
	var tokenDelegate delegation.DelegationStakeUpdated
	var tokenUndelegate delegation.DelegationWithdrawRequest
	var tokenRedelegation delegation.DelegationTransferRequest
	var newValidator validator.ValidatorValidatorMetaUpdated
	var updateValidator validator.ValidatorValidatorUpdated

	for _, log := range recipient.Logs {
		eventValidatorByID, errEvent := validatorMaster.EventByID(log.Topics[0])
		if errEvent == nil {
			if eventValidatorByID.Name == "ValidatorMetaUpdated" {
				_ = validatorMaster.UnpackIntoInterface(&newValidator, eventValidatorByID.Name, log.Data)
				var validatorInfo contracts.MasterValidatorValidatorAddedMeta
				_ = json.Unmarshal([]byte(newValidator.Meta), &validatorInfo)
				valAddr, _ := sdk.ValAddressFromHex(msg.From.String()[2:])
				validatorInfo.OperatorAddress = valAddr.String()

				err := k.CreateValidatorFromEVM(ctx, validatorInfo)
				fmt.Println(validatorInfo)
				fmt.Println(err)
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
		if errEvent == nil && strings.ToLower(log.Address.String()) == addressDelegation {
			if eventDelegationByID.Name == "StakeUpdated" {
				_ = contracts.UnpackLog(delegatorCenter, &tokenDelegate, eventDelegationByID.Name, log)
				_, err := k.coinKeeper.GetCoinByDRC(ctx, tokenDelegate.Stake.Token.String())
				if err != nil {
					symbolToken, _ := k.QuerySymbolToken(ctx, tokenDelegate.Stake.Token)
					coinUpdate, err := k.coinKeeper.GetCoin(ctx, symbolToken)
					if err == nil {
						_ = k.coinKeeper.UpdateCoinDRC(ctx, symbolToken, tokenDelegate.Stake.Token.String())
						coinUpdate.DRC20Contract = tokenDelegate.Stake.Token.String()
						k.coinKeeper.SetCoin(ctx, coinUpdate)
					}
				}
				err = k.Staked(ctx, tokenDelegate)
				if err != nil {
					return err
				}
			}

			if eventDelegationByID.Name == "RequestWithdraw" {
				_ = delegatorCenter.UnpackIntoInterface(&tokenUndelegate, eventDelegationByID.Name, log.Data)
				fmt.Println(tokenUndelegate)
				err := k.RequestWithdraw(ctx, tokenUndelegate)
				if err != nil {
					return err
				}
			}
			if eventDelegationByID.Name == "RequestTransfer" {
				_ = delegatorCenter.UnpackIntoInterface(&tokenRedelegation, eventDelegationByID.Name, log.Data)
				fmt.Println(tokenRedelegation)
				err := k.RequestTransfer(ctx, tokenRedelegation)
				if err != nil {
					return err
				}
			}
		}
		eventDelegationNftByID, errEvent := delegatorNftCenter.EventByID(log.Topics[0])
		if errEvent == nil && log.Address.String() == addressDelegationNft {
			if eventDelegationNftByID.Name == "StakedUpdated" {
				_ = delegatorCenter.UnpackIntoInterface(&tokenDelegate, eventDelegationNftByID.Name, log.Data)
				fmt.Println(tokenDelegate)
				err := k.Staked(ctx, tokenDelegate)
				if err != nil {
					return err
				}
			}

			if eventDelegationNftByID.Name == "RequestWithdraw" {
				_ = delegatorCenter.UnpackIntoInterface(&tokenUndelegate, eventDelegationNftByID.Name, log.Data)
				fmt.Println(tokenUndelegate)
				err := k.RequestWithdraw(ctx, tokenUndelegate)
				if err != nil {
					return err
				}
			}
			if eventDelegationNftByID.Name == "RequestTransfer" {
				_ = delegatorCenter.UnpackIntoInterface(&tokenRedelegation, eventDelegationNftByID.Name, log.Data)
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

func (k Keeper) Staked(ctx sdk.Context, stakeData delegation.DelegationStakeUpdated) error {

	coinStake, err := k.coinKeeper.GetCoinByDRC(ctx, stakeData.Stake.Token.String())
	if err != nil {
		return errors.CoinDoesNotExist
	}

	stake := types.NewStakeCoin(sdk.Coin{Denom: coinStake.Denom, Amount: math.NewIntFromBigInt(stakeData.Stake.Amount)})

	//if stakeData.Stake.HoldTimestamp != nil {
	//	stake.HoldStartTime = time.Now().Unix()
	//	stake.HoldEndTime = stakeData.Stake.HoldTimestamp.Int64()
	//}

	delegatorAddress, _ := types2.GetDecimalAddressFromHex(stakeData.Stake.Delegator.String())

	mintCoinForDelegation := sdk.NewCoins(sdk.NewCoin(coinStake.Denom, math.NewIntFromBigInt(stakeData.Stake.Amount)))
	err = k.bankKeeper.MintCoins(ctx, cointypes.ModuleName, mintCoinForDelegation)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, cointypes.ModuleName, delegatorAddress, mintCoinForDelegation)
	if err != nil {
		return err
	}

	valAddr, err := sdk.ValAddressFromHex(stakeData.Stake.Validator.String()[2:])

	validatorCosmos, found := k.GetValidator(ctx, valAddr)
	if !found {
		return fmt.Errorf("not found validator %s", valAddr)
	}

	_ = k.Delegate(ctx, delegatorAddress, validatorCosmos, stake)
	if err != nil {
		return err
	}
	fmt.Println("finish delegation")
	fmt.Println(stake)
	fmt.Println(delegatorAddress)
	fmt.Println(validatorCosmos)
	fmt.Println("finish delegation")
	return nil
}

func (k Keeper) RequestWithdraw(ctx sdk.Context, tokenUndelegate delegation.DelegationWithdrawRequest) error {
	coinStake, err := k.coinKeeper.GetCoinByDRC(ctx, tokenUndelegate.FrozenStake.Stake.Token.String())
	if err != nil {
		return errors.CoinDoesNotExist
	}

	stake := types.NewStakeCoin(sdk.Coin{Denom: coinStake.Denom, Amount: math.NewIntFromBigInt(tokenUndelegate.FrozenStake.Stake.Amount)})

	delegatorAddress, _ := types2.GetDecimalAddressFromHex(tokenUndelegate.FrozenStake.Stake.Delegator.String())

	valAddr, err := sdk.ValAddressFromBech32(tokenUndelegate.FrozenStake.Stake.Validator.String()[2:])

	delegationCosmos, found := k.GetDelegation(ctx, delegatorAddress, valAddr, stake.ID)
	if !found {
		return errors.DelegationNotFound
	}

	remainStake, err := k.CalculateRemainStake(ctx, delegationCosmos.Stake, stake)
	if err != nil {
		return err
	}

	_, err = k.Undelegate(ctx, delegatorAddress, valAddr, stake, remainStake)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) RequestTransfer(ctx sdk.Context, tokenRedelegation delegation.DelegationTransferRequest) error {
	coinStake, err := k.coinKeeper.GetCoinByDRC(ctx, tokenRedelegation.FrozenStake.Stake.Token.String())
	if err != nil {
		return errors.CoinDoesNotExist
	}

	stake := types.NewStakeCoin(sdk.Coin{Denom: coinStake.Denom, Amount: math.NewIntFromBigInt(tokenRedelegation.FrozenStake.Stake.Amount)})

	delegatorAddress, _ := types2.GetDecimalAddressFromHex(tokenRedelegation.FrozenStake.Stake.Delegator.String())

	valAddr, err := sdk.ValAddressFromHex(tokenRedelegation.FrozenStake.Stake.Validator.String()[2:])

	delegationCosmos, found := k.GetDelegation(ctx, delegatorAddress, valAddr, stake.ID)
	if !found {
		return errors.DelegationNotFound
	}

	remainStake, err := k.CalculateRemainStake(ctx, delegationCosmos.Stake, stake)
	if err != nil {
		return err
	}

	_, err = k.BeginRedelegation(
		ctx, delegatorAddress, valAddr, valAddr, stake, remainStake,
	)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) CreateValidatorFromEVM(ctx sdk.Context, validatorMeta contracts.MasterValidatorValidatorAddedMeta) error {

	commission, _ := sdkmath.NewIntFromString(validatorMeta.Commission)
	stakeSum, _ := sdkmath.NewIntFromString(validatorMeta.Stake)

	msg := types.MsgCreateValidator{
		OperatorAddress: validatorMeta.OperatorAddress,
		RewardAddress:   validatorMeta.RewardAddress,
		ConsensusPubkey: typescodec.UnsafePackAny(validatorMeta.ConsensusPubkey),
		Description: types.Description{
			Moniker:         validatorMeta.Description.Moniker,
			Identity:        validatorMeta.Description.Identity,
			Website:         validatorMeta.Description.Website,
			SecurityContact: validatorMeta.Description.SecurityContact,
			Details:         validatorMeta.Description.Details,
		},
		Commission: sdk.NewDecFromInt(commission),
		Stake: sdk.Coin{
			Denom:  validatorMeta.Coin,
			Amount: helpers.EtherToWei(stakeSum),
		},
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
		return nil
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

	validatorCosmos, err := types.NewValidator(valAddr, rewardAddr, pk, msg.Description, msg.Commission)
	if err != nil {
		return err
	}
	validatorCosmos.Online = false
	validatorCosmos.Jailed = false

	k.SetValidator(ctx, validatorCosmos)
	k.SetValidatorByConsAddr(ctx, validatorCosmos)
	k.SetNewValidatorByPowerIndex(ctx, validatorCosmos)

	// call the after-creation hook
	if err = k.AfterValidatorCreated(ctx, validatorCosmos.GetOperator()); err != nil {
		return err
	}

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	// NOTE source will always be from a wallet which are unbonded
	stake := types.NewStakeCoin(msg.Stake)
	//err = k.Delegate(ctx, sdk.AccAddress(valAddr), validator, stake)
	//if err != nil {
	//	return err
	//}

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
	validatorCosmos, found := k.GetValidator(ctx, valAddr)
	if !found {
		return errors.ValidatorNotFound
	}

	if validatorCosmos.Online {
		if !validatorCosmos.Jailed {
			return nil
		}
	}

	// validator without delegations can't become online
	if !k.HasDelegations(ctx, valAddr) {
		return errors.ValidatorHasNoDelegations
	}

	k.DeleteValidatorByPowerIndex(ctx, validatorCosmos)

	// TODO: move Online and Jailed to store keys?
	validatorCosmos.Online = true
	validatorCosmos.Jailed = false

	delByValidator := k.GetAllDelegationsByValidator(ctx)
	customCoinStaked := k.GetAllCustomCoinsStaked(ctx)
	customCoinPrices := k.CalculateCustomCoinPrices(ctx, customCoinStaked)
	totalStake, err := k.CalculateTotalPowerWithDelegationsAndPrices(ctx, validatorCosmos.GetOperator(), delByValidator[validatorCosmos.OperatorAddress], customCoinPrices)
	if err != nil {
		return err
	}

	stake := TokensToConsensusPower(totalStake)
	if stake == 0 {
		return errors.ValidatorStakeTooSmall
	}

	validatorCosmos.Stake = stake

	rs, err := k.GetValidatorRS(ctx, valAddr)
	if err != nil {
		rs = types.ValidatorRS{
			Rewards:      sdkmath.ZeroInt(),
			TotalRewards: sdkmath.ZeroInt(),
		}
	}
	rs.Stake = stake
	k.SetValidator(ctx, validatorCosmos)
	k.SetValidatorByPowerIndex(ctx, validatorCosmos)
	k.SetValidatorRS(ctx, valAddr, rs)

	// StartHeight need for correct calculation of missing blocks
	consAdr, err := validatorCosmos.GetConsAddr()
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
	validatorCosmos, found := k.GetValidator(ctx, valAddr)
	if !found {
		return errors.ValidatorNotFound
	}
	if !validatorCosmos.Online {
		return errors.ValidatorAlreadyOffline
	}

	validatorCosmos.Online = false
	// TODO: optimize
	k.SetValidator(ctx, validatorCosmos)

	consAdr, err := validatorCosmos.GetConsAddr()
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
