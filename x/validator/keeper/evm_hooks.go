// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegation"
	"bitbucket.org/decimalteam/go-smart-node/contracts/validator"
	"bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
	validatorType "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	typescodec "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/tendermint/tendermint/crypto"
	cmtjson "github.com/tendermint/tendermint/libs/json"
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
	//params.UndelegationTime =
	//k.SetParams(ctx)

	addressValidator, _ := contracts.GetAddressFromContractCenter(ctx, k.evmKeeper, contracts.NameOfSlugForGetAddressMasterValidator)
	addressDelegation, _ := contracts.GetAddressFromContractCenter(ctx, k.evmKeeper, contracts.NameOfSlugForGetAddressDelegation)
	addressValidator = strings.ToLower(addressValidator)
	addressDelegation = strings.ToLower(addressDelegation)
	validatorMaster, _ := validator.ValidatorMetaData.GetAbi()
	delegatorCenter, _ := delegation.DelegationMetaData.GetAbi()

	// this var is only for new token create from token center
	var tokenDelegate delegation.DelegationStakeUpdated
	var transferExistStake delegation.DelegationStakeUpdated
	var tokenUndelegate delegation.DelegationWithdrawRequest
	var tokenRedelegation delegation.DelegationTransferRequest
	var tokenDelegationAmount delegation.DelegationStakeAmountUpdated
	var newValidator validator.ValidatorValidatorMetaUpdated
	var updateValidator validator.ValidatorValidatorUpdated

	//var stakeUpdate []delegation.DelegationStakeUpdated

	redelegation := false
	undelegate := false
	stakedHold := false
	transferCompleted := false
	withdrawCompleted := false
	createValidator := false
	stakeReset := false
	stakeUpdate := 0

	for _, log := range recipient.Logs {
		eventValidatorByID, errEvent := validatorMaster.EventByID(log.Topics[0])
		if errEvent == nil && addressValidator == strings.ToLower(log.Address.String()) {
			fmt.Println(eventValidatorByID.Name)
			if eventValidatorByID.Name == "ValidatorMetaUpdated" {
				createValidator = true
				_ = validatorMaster.UnpackIntoInterface(&newValidator, eventValidatorByID.Name, log.Data)
				var validatorInfo contracts.MasterValidatorValidatorAddedMeta
				fmt.Println(newValidator.Meta)
				_ = json.Unmarshal([]byte(newValidator.Meta), &validatorInfo)
				valAddr, _ := sdk.ValAddressFromHex(msg.From.String()[2:])
				validatorInfo.OperatorAddress = valAddr.String()
				fmt.Println(validatorInfo)
				err := k.CreateValidatorFromEVM(ctx, validatorInfo)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, log := range recipient.Logs {
		eventValidatorByID, errEvent := validatorMaster.EventByID(log.Topics[0])
		if errEvent == nil && addressValidator == strings.ToLower(log.Address.String()) {
			fmt.Println(eventValidatorByID.Name)
			if eventValidatorByID.Name == "ValidatorUpdated" && !createValidator {
				_ = contracts.UnpackLog(validatorMaster, &updateValidator, eventValidatorByID.Name, log)
				fmt.Println(updateValidator)
				cosmosAddressValidator, _ := sdk.ValAddressFromHex(updateValidator.Validator.String()[2:])
				if updateValidator.Paused == false {
					err := k.SetOnlineFromEvm(ctx, cosmosAddressValidator.String())
					if err != nil {
						return err
					}
				}
				if updateValidator.Paused == true {
					err := k.SetOfflineFromEvm(ctx, cosmosAddressValidator.String())
					if err != nil {
						return err
					}
				}
			}
			// Contract-set per-block reward override. The node clamps the value
			// to the built-in schedule (reduce-only) when consuming it in
			// GetBlockReward, so this only ever lowers emission.
			if eventValidatorByID.Name == "RewardPerBlockUpdated" {
				var rewardPerBlockUpdated validator.ValidatorRewardPerBlockUpdated
				_ = contracts.UnpackLog(validatorMaster, &rewardPerBlockUpdated, eventValidatorByID.Name, log)
				if rewardPerBlockUpdated.Enabled {
					k.SetRewardPerBlockOverride(ctx, sdkmath.NewIntFromBigInt(rewardPerBlockUpdated.RewardPerBlock))
				} else {
					k.ClearRewardPerBlockOverride(ctx)
				}
			}
		}
	}

	for _, log := range recipient.Logs {
		eventDelegationByID, errEvent := delegatorCenter.EventByID(log.Topics[0])
		if errEvent == nil {
			if eventDelegationByID.Name == "WithdrawRequest" {
				undelegate = true
			}
			if eventDelegationByID.Name == "TransferRequest" {
				redelegation = true
			}
			if eventDelegationByID.Name == "StakeHolded" {
				stakedHold = true
			}
			if eventDelegationByID.Name == "TransferCompleted" {
				transferCompleted = true
			}
			if eventDelegationByID.Name == "WithdrawCompleted" {
				withdrawCompleted = true
			}
			if eventDelegationByID.Name == "StakeReset" {
				stakeReset = true
			}
		}
	}

	srcValidatorRedelegation := ""

	// Index the coin-side events (StakeUpdated / StakeAmountUpdated) by stakeId so an NFT-typed
	// WithdrawRequest/TransferRequest can be rewritten to coin-shaped before dispatch. NFT
	// withdraw/transfer emits both under the same stakeId as the frozen-stake request.
	coinByStakeId := make(map[[32]byte]delegation.IDecimalDelegationCommonStake)
	amountByStakeId := make(map[[32]byte]*big.Int)
	for _, log := range recipient.Logs {
		if len(log.Topics) == 0 {
			continue
		}
		ev, evErr := delegatorCenter.EventByID(log.Topics[0])
		if evErr != nil || strings.ToLower(log.Address.String()) != addressDelegation {
			continue
		}
		switch ev.Name {
		case "StakeUpdated":
			var su delegation.DelegationStakeUpdated
			if contracts.UnpackLog(delegatorCenter, &su, ev.Name, log) == nil {
				coinByStakeId[su.StakeId] = su.Stake
			}
		case "StakeAmountUpdated":
			var sa delegation.DelegationStakeAmountUpdated
			if contracts.UnpackLog(delegatorCenter, &sa, ev.Name, log) == nil {
				amountByStakeId[sa.StakeId] = sa.ChangedAmount
			}
		}
	}

	for _, log := range recipient.Logs {
		eventDelegationByID, errEvent := delegatorCenter.EventByID(log.Topics[0])
		if errEvent == nil && strings.ToLower(log.Address.String()) == addressDelegation {
			fmt.Println(eventDelegationByID.Name)
			if eventDelegationByID.Name == "StakeAmountUpdated" {
				_ = contracts.UnpackLog(delegatorCenter, &tokenDelegationAmount, eventDelegationByID.Name, log)
			}
			if eventDelegationByID.Name == "StakeUpdated" && redelegation && !undelegate && !stakedHold && !transferCompleted && !withdrawCompleted && !stakeReset {
				_ = contracts.UnpackLog(delegatorCenter, &tokenDelegate, eventDelegationByID.Name, log)
				srcValidatorRedelegation = tokenDelegate.Stake.Validator.String()
			}
			if eventDelegationByID.Name == "StakeUpdated" && !redelegation && !undelegate && !stakedHold && !transferCompleted && !withdrawCompleted && !stakeReset && stakeUpdate == 0 {
				if tokenDelegationAmount.ChangedAmount == nil {
					return errors.DelegationSumIsNotSet
				}
				stakeUpdate = stakeUpdate + 1
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
				tokenDelegate.Stake.Amount = tokenDelegationAmount.ChangedAmount
				err = k.Staked(ctx, tokenDelegate, true)
				if err != nil {
					return err
				}
			}

			if eventDelegationByID.Name == "StakeHolded" {
				_ = contracts.UnpackLog(delegatorCenter, &transferExistStake, eventDelegationByID.Name, log)
				_, err := k.coinKeeper.GetCoinByDRC(ctx, transferExistStake.Stake.Token.String())
				if err != nil {
					symbolToken, _ := k.QuerySymbolToken(ctx, transferExistStake.Stake.Token)
					coinUpdate, err := k.coinKeeper.GetCoin(ctx, symbolToken)
					if err == nil {
						_ = k.coinKeeper.UpdateCoinDRC(ctx, symbolToken, transferExistStake.Stake.Token.String())
						coinUpdate.DRC20Contract = transferExistStake.Stake.Token.String()
						k.coinKeeper.SetCoin(ctx, coinUpdate)
					}
				}
				transferExistStake.Stake.Amount = tokenDelegationAmount.ChangedAmount
				err = k.Staked(ctx, transferExistStake, false)
				if err != nil {
					return err
				}
			}

			if eventDelegationByID.Name == "WithdrawRequest" {
				_ = delegatorCenter.UnpackIntoInterface(&tokenUndelegate, eventDelegationByID.Name, log.Data)
				srcVal, applied, nftErr := rewriteFrozenStakeForNFT(&tokenUndelegate.FrozenStake.Stake, tokenUndelegate.StakeId, coinByStakeId, amountByStakeId)
				if nftErr != nil {
					return nftErr
				}
				_, err := k.coinKeeper.GetCoinByDRC(ctx, tokenUndelegate.FrozenStake.Stake.Token.String())
				if err != nil {
					symbolToken, _ := k.QuerySymbolToken(ctx, tokenUndelegate.FrozenStake.Stake.Token)
					coinUpdate, err := k.coinKeeper.GetCoin(ctx, symbolToken)
					if err == nil {
						_ = k.coinKeeper.UpdateCoinDRC(ctx, symbolToken, tokenUndelegate.FrozenStake.Stake.Token.String())
						coinUpdate.DRC20Contract = tokenUndelegate.FrozenStake.Stake.Token.String()
						k.coinKeeper.SetCoin(ctx, coinUpdate)
					}
				}
				//tokenUndelegate.FrozenStake.Stake.Amount = tokenDelegationAmount.ChangedAmount
				fmt.Println(tokenUndelegate)
				// Deployed-contract compat: a Transfer-typed NFT request arriving as WithdrawRequest
				// (some deployed contract versions emit one event for both unbond and redelegate) is a
				// redelegate, not an unbond — route it to RequestTransfer with the source validator from
				// the sibling coin event and the destination from the frozen stake, so the reserve
				// delegation moves instead of being stranded. A correctly-emitted TransferRequest is
				// handled in the branch below and never reaches here.
				if applied && tokenUndelegate.FrozenStake.FreezeType == freezeTypeTransfer {
					err = k.RequestTransfer(ctx, delegation.DelegationTransferRequest{
						StakeId:     tokenUndelegate.StakeId,
						StakeIndex:  tokenUndelegate.StakeIndex,
						FrozenStake: tokenUndelegate.FrozenStake,
					}, srcVal.String())
				} else {
					err = k.RequestWithdraw(ctx, tokenUndelegate)
				}
				if err != nil {
					return err
				}
			}
			if eventDelegationByID.Name == "TransferRequest" {
				_ = delegatorCenter.UnpackIntoInterface(&tokenRedelegation, eventDelegationByID.Name, log.Data)
				srcValidator := srcValidatorRedelegation
				if src, applied, nftErr := rewriteFrozenStakeForNFT(&tokenRedelegation.FrozenStake.Stake, tokenRedelegation.StakeId, coinByStakeId, amountByStakeId); nftErr != nil {
					return nftErr
				} else if applied {
					srcValidator = src.String()
				}
				_, err := k.coinKeeper.GetCoinByDRC(ctx, tokenRedelegation.FrozenStake.Stake.Token.String())
				if err != nil {
					symbolToken, _ := k.QuerySymbolToken(ctx, tokenRedelegation.FrozenStake.Stake.Token)
					coinUpdate, err := k.coinKeeper.GetCoin(ctx, symbolToken)
					if err == nil {
						_ = k.coinKeeper.UpdateCoinDRC(ctx, symbolToken, tokenDelegate.Stake.Token.String())
						coinUpdate.DRC20Contract = tokenRedelegation.FrozenStake.Stake.Token.String()
						k.coinKeeper.SetCoin(ctx, coinUpdate)
					}
				}
				fmt.Println(tokenRedelegation)
				fmt.Println(srcValidator)
				err = k.RequestTransfer(ctx, tokenRedelegation, srcValidator)
				if err != nil {
					return err
				}
			}
		}
	}

	//if len(stakeUpdate) == 0 {
	//
	//}
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

// NFT frozen-stake token types, mirroring Solidity enum TokenType
// (None=0, DRC20=1, DRC721=2, DRC1155=3).
const (
	tokenTypeDRC721  uint8 = 2
	tokenTypeDRC1155 uint8 = 3
)

// FrozenStake.FreezeType, mirroring the Solidity enum (Withdraw=1, Transfer=2).
// Some deployed DecimalDelegation versions emit a single WithdrawRequest event for
// both NFT unbond and NFT redelegate, distinguished only by this field.
const freezeTypeTransfer uint8 = 2

// isNFTStake reports whether a frozen-stake token type is an NFT type (DRC721 or DRC1155).
func isNFTStake(tokenType uint8) bool {
	return tokenType == tokenTypeDRC721 || tokenType == tokenTypeDRC1155
}

// rewriteFrozenStakeForNFT converts an NFT-typed frozen stake into a coin-shaped one so the
// existing coin undelegate/redelegate path can process an NFT withdraw/redelegate.
//
// NFT withdraw/transfer events carry frozenStake.Stake.Token = the NFT contract (which the
// contract's completion path needs) and Amount = 1/nftCount — the node cannot resolve that to a
// coin. But the same tx also emits the coin-side StakeUpdated (Stake.Token = reserveToken, a real
// coin; Stake.Validator = source validator) and StakeAmountUpdated (-reserveAmount), both under
// the same stakeId. This rewrites the frozen stake's Token/Amount to the reserve coin + reserve
// amount from those sibling events and returns the source validator for the redelegation path.
//
// applied=false (err=nil) for non-NFT stakes — the caller proceeds unchanged. err!=nil only when
// the stake IS NFT-typed but the sibling coin events are missing (malformed tx / event desync),
// so the failure is diagnosable in node logs instead of a silent revert.
func rewriteFrozenStakeForNFT(
	stake *delegation.IDecimalDelegationCommonStake,
	stakeId [32]byte,
	coinByStakeId map[[32]byte]delegation.IDecimalDelegationCommonStake,
	amountByStakeId map[[32]byte]*big.Int,
) (srcValidator common.Address, applied bool, err error) {
	if stake == nil || !isNFTStake(stake.TokenType) {
		return common.Address{}, false, nil
	}
	coin, hasCoin := coinByStakeId[stakeId]
	changed, hasAmount := amountByStakeId[stakeId]
	if !hasCoin || !hasAmount || changed == nil {
		return common.Address{}, true, fmt.Errorf(
			"nft exit stakeId %x: missing sibling coin event (haveStakeUpdated=%v haveStakeAmountUpdated=%v)",
			stakeId, hasCoin, hasAmount && changed != nil)
	}
	stake.Token = coin.Token
	stake.Amount = math.NewIntFromBigInt(changed).Abs().BigInt()
	return coin.Validator, true, nil
}

// resolveHoldStartTime picks a hold's start time. For a brand-new hold bucket
// (contract isNew flag) it trusts the contract-supplied start, falling back to the
// current block time when the event carries 0 (pre-backfill stakes / defense-in-depth
// against a 0 that would read as >=1yr old). For a top-up to an existing bucket it
// always uses the current block time, preserving the per-segment anti-gaming behaviour.
func resolveHoldStartTime(blockTime time.Time, eventStart *big.Int, isNewBucket bool) int64 {
	if isNewBucket && eventStart != nil && eventStart.Sign() > 0 {
		return eventStart.Int64()
	}
	return blockTime.Unix()
}

func (k Keeper) Staked(ctx sdk.Context, stakeData delegation.DelegationStakeUpdated, newStake bool) error {

	coinStake, err := k.coinKeeper.GetCoinByDRC(ctx, stakeData.Stake.Token.String())
	if err != nil {
		return errors.CoinDoesNotExist
	}

	if coinStake.Denom == "" {
		return errors.CoinDoesNotExist
	}

	stake := validatorType.NewStakeCoin(sdk.Coin{Denom: coinStake.Denom, Amount: math.NewIntFromBigInt(stakeData.Stake.Amount)})

	if stakeData.Stake.HoldTimestamp.Int64() != 0 {
		var newHold validatorType.StakeHold
		newHold.Amount = math.NewIntFromBigInt(stakeData.Stake.Amount)
		newHold.HoldStartTime = resolveHoldStartTime(ctx.BlockTime(), stakeData.Stake.HoldStartTime, stakeData.IsNew)
		newHold.HoldEndTime = stakeData.Stake.HoldTimestamp.Int64()
		stake.Holds = append(stake.Holds, &newHold)
	}

	delegatorAddress, _ := types.GetDecimalAddressFromHex(stakeData.Stake.Delegator.String())

	//mintCoinForDelegation := sdk.NewCoins(sdk.NewCoin(coinStake.Denom, math.NewIntFromBigInt(stakeData.Stake.Amount)))
	//err = k.bankKeeper.MintCoins(ctx, cointypes.ModuleName, mintCoinForDelegation)
	//if err != nil {
	//	return err
	//}
	//err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, cointypes.ModuleName, delegatorAddress, mintCoinForDelegation)
	//if err != nil {
	//	return err
	//}

	valAddr, err := sdk.ValAddressFromHex(stakeData.Stake.Validator.String()[2:])

	validatorCosmos, found := k.GetValidator(ctx, valAddr)
	if !found {
		return fmt.Errorf("not found validator %s", valAddr)
	}

	if newStake {
		_ = k.Delegate(ctx, delegatorAddress, validatorCosmos, stake)
		if err != nil {
			return err
		}
	} else {
		_ = k.TransferToHold(ctx, delegatorAddress, validatorCosmos, stake)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) RequestWithdraw(ctx sdk.Context, tokenUndelegate delegation.DelegationWithdrawRequest) error {

	coinStake, err := k.coinKeeper.GetCoinByDRC(ctx, tokenUndelegate.FrozenStake.Stake.Token.String())
	if err != nil {
		return errors.CoinDoesNotExist
	}

	if coinStake.Denom == "" {
		return errors.CoinDoesNotExist
	}

	stake := validatorType.NewStakeCoin(sdk.Coin{Denom: coinStake.Denom, Amount: math.NewIntFromBigInt(tokenUndelegate.FrozenStake.Stake.Amount)})

	if tokenUndelegate.FrozenStake.Stake.HoldTimestamp.Int64() != 0 {
		var newHold validatorType.StakeHold
		newHold.Amount = math.NewIntFromBigInt(tokenUndelegate.FrozenStake.Stake.Amount)
		newHold.HoldStartTime = resolveHoldStartTime(ctx.BlockTime(), tokenUndelegate.FrozenStake.Stake.HoldStartTime, true)
		newHold.HoldEndTime = tokenUndelegate.FrozenStake.Stake.HoldTimestamp.Int64()
		stake.Holds = append(stake.Holds, &newHold)
	}

	delegatorAddress, _ := types.GetDecimalAddressFromHex(tokenUndelegate.FrozenStake.Stake.Delegator.String())

	valAddr, err := sdk.ValAddressFromHex(tokenUndelegate.FrozenStake.Stake.Validator.String()[2:])

	delegationCosmos, found := k.GetDelegation(ctx, delegatorAddress, valAddr, stake.ID)
	if !found {
		// Delegation already removed by CheckDelegations (force-undelegate).
		// EVM-side withdrawal completes normally; Cosmos side is already done.
		ctx.Logger().Info("WithdrawRequest: delegation not found, skipping",
			"delegator", delegatorAddress,
			"validator", valAddr,
		)
		return nil
	}

	remainStake, err := k.CalculateRemainStake(ctx, delegationCosmos.Stake, stake)
	if err != nil {
		return err
	}

	if len(stake.Holds) != 0 {
		moved := stake.Holds[0]
		// Subtract the withdrawn held amount from the source delegation's remaining holds so
		// sum(holds) does not exceed the remaining stake (otherwise a partial held withdraw
		// over-pays the >=1yr long-hold reward bonus and over-enqueues auto-unbond). Withdraw
		// has no destination, so the rebuilt moved sub-holds are discarded.
		var leftover math.Int
		remainStake.Holds, _, leftover = applyTransferredHold(remainStake.Holds, moved)
		if leftover.IsPositive() {
			// Hold already (partly) pruned by DeleteHoldMature after expiry; clamp at zero.
			ctx.Logger().Debug("WithdrawRequest: withdrawn hold exceeds source hold records",
				"delegator", delegatorAddress.String(),
				"hold_end", moved.HoldEndTime,
				"leftover", leftover.String(),
			)
		}
	}

	_, err = k.Undelegate(ctx, delegatorAddress, valAddr, stake, remainStake, nil)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) RequestTransfer(ctx sdk.Context, tokenRedelegation delegation.DelegationTransferRequest, srcValidator string) error {
	coinStake, err := k.coinKeeper.GetCoinByDRC(ctx, tokenRedelegation.FrozenStake.Stake.Token.String())
	if err != nil {
		return errors.CoinDoesNotExist
	}

	stake := validatorType.NewStakeCoin(sdk.Coin{Denom: coinStake.Denom, Amount: math.NewIntFromBigInt(tokenRedelegation.FrozenStake.Stake.Amount)})

	if tokenRedelegation.FrozenStake.Stake.HoldTimestamp.Int64() != 0 {
		var newHold validatorType.StakeHold
		newHold.Amount = math.NewIntFromBigInt(tokenRedelegation.FrozenStake.Stake.Amount)
		// The start read from the FrozenStake event is the stored original hold start.
		// applyTransferredHold below rebuilds the moved hold(s) from source delegation
		// hold segments (each with its own real HoldStartTime), overwriting this value
		// for the redelegated (destination) hold. Pass isNewBucket=true so the event
		// value is trusted with the 0-guard fallback.
		newHold.HoldStartTime = resolveHoldStartTime(ctx.BlockTime(), tokenRedelegation.FrozenStake.Stake.HoldStartTime, true)
		newHold.HoldEndTime = tokenRedelegation.FrozenStake.Stake.HoldTimestamp.Int64()
		stake.Holds = append(stake.Holds, &newHold)
	}

	delegatorAddress, _ := types.GetDecimalAddressFromHex(tokenRedelegation.FrozenStake.Stake.Delegator.String())

	srcValAddr, err := sdk.ValAddressFromHex(srcValidator[2:])

	delegationCosmos, found := k.GetDelegation(ctx, delegatorAddress, srcValAddr, stake.ID)
	if !found {
		return errors.DelegationNotFound
	}

	valAddr, err := sdk.ValAddressFromHex(tokenRedelegation.FrozenStake.Stake.Validator.String()[2:])

	remainStake, err := k.CalculateRemainStake(ctx, delegationCosmos.Stake, stake)
	if err != nil {
		return err
	}

	if len(stake.Holds) != 0 {
		moved := stake.Holds[0]
		// Subtract the moved held amount from the source's holds and rebuild the moved
		// hold(s) as per-source-start segments, which become the destination's held credit.
		// The destination is credited exactly what was sourced; any leftover (a hold no
		// longer recorded on the node, e.g. matured and pruned) lands as ordinary unheld
		// stake on the destination instead of inheriting a >=1yr reward window.
		var leftover math.Int
		remainStake.Holds, stake.Holds, leftover = applyTransferredHold(remainStake.Holds, moved)
		if leftover.IsPositive() {
			ctx.Logger().Debug("RequestTransfer: transferred hold exceeds source hold records",
				"delegator", delegatorAddress.String(),
				"hold_end", moved.HoldEndTime,
				"leftover", leftover.String(),
			)
		}
	}

	_, err = k.BeginRedelegation(
		ctx, delegatorAddress, srcValAddr, valAddr, stake, remainStake, nil,
	)
	if err != nil {
		return err
	}

	return nil
}

// applyTransferredHold reconciles a redelegated/withdrawn held amount (moved) against a
// source delegation's holds at moved.HoldEndTime. It returns the source's remaining holds,
// the moved amount split into per-source-start sub-holds, and any amount that could not be
// sourced (leftover).
//
// The EVM merges all holds with the same (validator, delegator, token, holdTimestamp) into
// one stake, whereas the node stores holds as an un-merged list — so the node may carry
// several entries for one HoldEndTime, each with its own real HoldStartTime. The moved
// amount is drawn from those entries in order (FIFO), clamped so no hold goes negative, and
// each drawn segment becomes a moved sub-hold carrying the REAL start of the source hold it
// came from. This (a) prevents a young same-end hold from inheriting an older hold's >=1yr
// long-hold reward window, and (b) credits the destination a held amount equal to exactly
// what was sourced — never more. Fully consumed and any pre-existing non-positive source
// holds are dropped.
//
// leftover is the moved amount with no matching source hold (>0 only on EVM/node desync,
// e.g. a hold that already matured and was pruned by DeleteHoldMature but is still movable
// on the EVM). In that case sum(movedHolds) = moved.Amount - leftover, and the caller leaves
// the remaining moved principal as ordinary (unheld) stake on the destination.
func applyTransferredHold(
	remainHolds []*validatorType.StakeHold, moved *validatorType.StakeHold,
) (kept, movedHolds []*validatorType.StakeHold, leftover math.Int) {
	remaining := moved.Amount
	kept = make([]*validatorType.StakeHold, 0, len(remainHolds))
	for _, hold := range remainHolds {
		if hold.HoldEndTime == moved.HoldEndTime && remaining.IsPositive() && hold.Amount.IsPositive() {
			sub := hold.Amount
			if sub.GT(remaining) {
				sub = remaining
			}
			movedHolds = append(movedHolds, &validatorType.StakeHold{
				Amount:        sub,
				HoldStartTime: hold.HoldStartTime,
				HoldEndTime:   moved.HoldEndTime,
			})
			hold.Amount = hold.Amount.Sub(sub)
			remaining = remaining.Sub(sub)
		}
		if hold.Amount.IsPositive() {
			kept = append(kept, hold)
		}
	}
	return kept, movedHolds, remaining
}

func (k Keeper) CreateValidatorFromEVM(ctx sdk.Context, validatorMeta contracts.MasterValidatorValidatorAddedMeta) error {

	// Commission arrives in the validator meta as a percentage in [0, 100]
	// (e.g. "20" or "20.000000000000000000" for 20%). Validators persist it as
	// a fraction in [0, 1] to match genesis and the reward math
	// (reward.go: sdk.NewDecFromInt(rewards).Mul(val.Commission)), so divide by 100.
	commissionPct, err := sdk.NewDecFromStr(strings.TrimSpace(validatorMeta.Commission.String()))
	if err != nil {
		return errors.Internal.Wrapf("invalid validator commission %q: %s", validatorMeta.Commission.String(), err.Error())
	}
	if commissionPct.GT(sdk.NewDec(100)) {
		return errors.ValidatorCommissionIsTooBig
	}
	if commissionPct.IsNegative() {
		return errors.ValidatorCommissionIsTooSmall
	}
	commissionFraction := commissionPct.QuoInt64(100)

	rewardAddress, _ := types.GetDecimalAddressFromHex(validatorMeta.RewardAddress)

	var pubKey crypto.PubKey
	_ = cmtjson.Unmarshal(
		[]byte(fmt.Sprintf("{\"type\":\"tendermint/PubKeyEd25519\",\"value\":\"%s\"}", validatorMeta.ConsensusPubkey)),
		&pubKey)

	valPubKey, err := cryptocodec.FromTmPubKeyInterface(pubKey)
	if err != nil {
		return err
	}
	fmt.Println("MsgCreateValidator")
	msg := validatorType.MsgCreateValidator{
		OperatorAddress: rewardAddress.String(),
		RewardAddress:   rewardAddress.String(),
		ConsensusPubkey: typescodec.UnsafePackAny(valPubKey),
		Description: validatorType.Description{
			Moniker:         validatorMeta.Description.Moniker,
			Identity:        validatorMeta.Description.Identity,
			Website:         validatorMeta.Description.Website,
			SecurityContact: validatorMeta.Description.SecurityContact,
			Details:         validatorMeta.Description.Details,
		},
		Commission: commissionFraction,
		Stake:      sdk.Coin{},
	}
	fmt.Println("ValAddressFromBech32")
	valAddr, err := sdk.ValAddressFromBech32(validatorMeta.OperatorAddress)
	if err != nil {
		return err
	}
	rewardAddr, err := sdk.AccAddressFromBech32(msg.RewardAddress)
	if err != nil {
		return err
	}

	// check to see if the pubkey or sender has been registered before
	if valEdit, found := k.GetValidator(ctx, valAddr); found {
		// validator must already be registered
		// replace all editable fields (clients should autofill existing values)
		description, err := valEdit.Description.UpdateDescription(msg.Description)
		if err != nil {
			return err
		}

		valEdit.Description = description
		valEdit.RewardAddress = msg.RewardAddress
		valEdit.Commission = msg.Commission

		k.SetValidator(ctx, valEdit)

		err = events.EmitTypedEvent(ctx, &validatorType.EventEditValidator{
			Sender:        sdk.AccAddress(valAddr).String(),
			Validator:     valAddr.String(),
			RewardAddress: msg.RewardAddress,
			Description:   description,
		})
		return err
	}

	pk, ok := msg.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return errors.InvalidConsensusPubKey
	}
	fmt.Println("ValAddressFromBech32")
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

	validatorCosmos, err := validatorType.NewValidator(valAddr, rewardAddr, pk, msg.Description, msg.Commission)
	if err != nil {
		return err
	}
	validatorCosmos.Online = false
	validatorCosmos.Jailed = false
	fmt.Println("ValAddressFromBech32")
	k.SetValidator(ctx, validatorCosmos)
	k.SetValidatorByConsAddr(ctx, validatorCosmos)
	k.SetNewValidatorByPowerIndex(ctx, validatorCosmos)

	// call the after-creation hook
	if err = k.AfterValidatorCreated(ctx, validatorCosmos.GetOperator()); err != nil {
		return err
	}
	fmt.Println("ValAddressFromBech32")
	err = events.EmitTypedEvent(ctx, &validatorType.EventCreateValidator{
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
	fmt.Println("finish create")
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

	// Clear auto-unbond timer when validator comes back online
	k.DeleteValidatorOfflineSince(ctx, valAddr)

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
		rs = validatorType.ValidatorRS{
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

	err = events.EmitTypedEvent(ctx, &validatorType.EventSetOnline{
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

	// Start auto-unbond timer
	k.SetValidatorOfflineSince(ctx, valAddr, ctx.BlockTime())

	consAdr, err := validatorCosmos.GetConsAddr()
	if err != nil {
		return err
	}
	k.DeleteStartHeight(ctx, consAdr)

	err = events.EmitTypedEvent(ctx, &validatorType.EventSetOffline{
		Sender:    sdk.AccAddress(valAddr).String(),
		Validator: valAddr.String(),
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}

	return nil
}
