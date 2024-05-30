// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts/center"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegation"
	"bitbucket.org/decimalteam/go-smart-node/contracts/validator"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// QueryAddressWDEL returns the data of a deployed ERC20 contract
func (k *Keeper) QueryAddressWDEL(
	ctx sdk.Context,
	contract common.Address,
) (string, error) {

	contractCenter, _ := center.CenterMetaData.GetAbi()
	methodCall := "getAddress"
	// Address token center
	res, err := k.evmKeeper.CallEVM(ctx, *contractCenter, common.Address(types.ModuleAddress), contract, false, methodCall, types.NameOfSlugForGetAddressWDEL)
	if err != nil {
		return new(common.Address).Hex(), err
	}
	data, err := contractCenter.Unpack(methodCall, res.Ret)
	if len(data) == 0 {
		return new(common.Address).Hex(), err
	}
	return data[0].(common.Address).String(), err
}

// QueryAddressDelegation returns the data of a deployed ERC20 contract
func (k *Keeper) QueryAddressDelegation(
	ctx sdk.Context,
	contract common.Address,
) (string, error) {

	contractCenter, _ := center.CenterMetaData.GetAbi()
	methodCall := "getAddress"
	// Address token center
	res, err := k.evmKeeper.CallEVM(ctx, *contractCenter, common.Address(types.ModuleAddress), contract, false, methodCall, types.NameOfSlugForGetAddressDelegation)
	if err != nil {
		return new(common.Address).Hex(), err
	}
	data, err := contractCenter.Unpack(methodCall, res.Ret)
	fmt.Println(data)
	if len(data) == 0 {
		return new(common.Address).Hex(), err
	}
	return data[0].(common.Address).String(), err
}

// QueryAddressMasterValidator returns the data of a deployed ERC20 contract
func (k *Keeper) QueryAddressMasterValidator(
	ctx sdk.Context,
	contract common.Address,
) (string, error) {

	contractCenter, _ := center.CenterMetaData.GetAbi()
	methodCall := "getAddress"
	// Address token center
	res, err := k.evmKeeper.CallEVM(ctx, *contractCenter, common.Address(types.ModuleAddress), contract, false, methodCall, types.NameOfSlugForGetAddressMasterValidator)
	if err != nil {
		return new(common.Address).Hex(), err
	}
	data, err := contractCenter.Unpack(methodCall, res.Ret)
	if len(data) == 0 {
		return new(common.Address).Hex(), err
	}
	return data[0].(common.Address).String(), err
}

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
