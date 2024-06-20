// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegation"
	"bitbucket.org/decimalteam/go-smart-node/contracts/token"
	"bitbucket.org/decimalteam/go-smart-node/contracts/validator"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// QueryIfNeedExecuteFinish returns the data of a deployed ERC20 contract
func (k *Keeper) QueryIfNeedExecuteFinish(
	ctx sdk.Context,
	contract common.Address,
) (bool, error) {

	contractDelegation, _ := delegation.DelegationMetaData.GetAbi()
	methodCall := "isQueueReady"
	// Address token center
	res, err := k.evmKeeper.CallEVM(ctx, *contractDelegation, common.Address(types.ModuleAddress), contract, false, methodCall)
	if err != nil {
		return false, err
	}
	data, err := contractDelegation.Unpack(methodCall, res.Ret)

	return data[0].(bool), err
}

// QuerySymbolToken returns the data of a deployed ERC20 contract
func (k *Keeper) QuerySymbolToken(
	ctx sdk.Context,
	contract common.Address,
) (string, error) {

	contractDelegation, _ := token.TokenMetaData.GetAbi()
	methodCall := "symbol"
	// Address token center
	res, err := k.evmKeeper.CallEVM(ctx, *contractDelegation, common.Address(types.ModuleAddress), contract, false, methodCall)
	if err != nil {
		return new(common.Address).Hex(), err
	}
	data, err := contractDelegation.Unpack(methodCall, res.Ret)

	return data[0].(string), err
}

// QueryOwnerDelegation returns the data of a deployed ERC20 contract
func (k *Keeper) QueryOwnerDelegation(
	ctx sdk.Context,
	contract common.Address,
) (string, error) {

	contractDelegation, _ := delegation.DelegationMetaData.GetAbi()
	methodCall := "owner"
	// Address token center
	res, err := k.evmKeeper.CallEVM(ctx, *contractDelegation, common.Address(types.ModuleAddress), contract, false, methodCall)
	if err != nil {
		return new(common.Address).Hex(), err
	}
	data, err := contractDelegation.Unpack(methodCall, res.Ret)

	return data[0].(common.Address).String(), err
}

// ExecuteQueueEVMAction returns the data of a deployed ERC20 contract
func (k *Keeper) ExecuteQueueEVMAction(
	ctx sdk.Context,
	contract common.Address,
) (bool, error) {

	contractDelegation, _ := delegation.DelegationMetaData.GetAbi()
	methodCall := "completeStake"
	// Address token center
	res, err := k.evmKeeper.CallEVM(ctx, *contractDelegation, common.Address(types.ModuleAddress), contract, true, methodCall)
	if err != nil {
		return false, err
	}
	data, err := contractDelegation.Unpack(methodCall, res.Ret)

	return data[0].(bool), err
}

// ExecuteAddPenalty returns the data of a deployed ERC20 contract
func (k *Keeper) ExecuteAddPenalty(
	ctx sdk.Context,
	contract common.Address,
	validatorAddress common.Address,
	penaltyPercent uint16,
) (bool, error) {

	contractDelegation, _ := validator.ValidatorMetaData.GetAbi()
	methodCall := "addPenalty"
	_, err := k.evmKeeper.CallEVM(ctx, *contractDelegation, common.Address(types.ModuleAddress), contract, true, methodCall, validatorAddress, penaltyPercent)
	if err != nil {
		return false, err
	}
	return true, err
}

// ExecuteBurnPenaltyTokens returns the data of a deployed ERC20 contract
func (k *Keeper) ExecuteBurnPenaltyTokens(
	ctx sdk.Context,
	contract common.Address,
	validatorAddress common.Address,
) (bool, error) {

	contractDelegation, _ := delegation.DelegationMetaData.GetAbi()
	methodCall := "burnPenaltyTokensValidator"
	_, err := k.evmKeeper.CallEVM(ctx, *contractDelegation, common.Address(types.ModuleAddress), contract, true, methodCall, validatorAddress)
	if err != nil {
		return false, err
	}
	return true, err
}
