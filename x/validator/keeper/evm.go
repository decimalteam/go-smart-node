// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// QueryAddressDelegation returns the data of a deployed ERC20 contract
func (k *Keeper) QueryAddressDelegation(
	ctx sdk.Context,
	contract common.Address,
) (string, error) {

	contractCenter, _ := contracts.ContractCenterMetaData.GetAbi()
	methodCall := "getAddress"
	// Address token center
	res, err := k.evmKeeper.CallEVM(ctx, *contractCenter, common.Address(types.ModuleAddress), contract, false, methodCall, types.NameOfSlugForGetAddressDelegation)
	if err != nil {
		return new(common.Address).Hex(), err
	}
	data, err := contractCenter.Unpack(methodCall, res.Ret)

	return data[0].(common.Address).String(), err
}

// QueryIfNeedExecuteFinish returns the data of a deployed ERC20 contract
func (k *Keeper) QueryIfNeedExecuteFinish(
	ctx sdk.Context,
	contract common.Address,
) (bool, error) {

	contractDelegation, _ := contracts.DelegationMetaData.GetAbi()
	methodCall := "isFrozenStakesQueueReady"
	// Address token center
	res, err := k.evmKeeper.CallEVM(ctx, *contractDelegation, common.Address(types.ModuleAddress), contract, false, methodCall)
	if err != nil {
		return false, err
	}
	data, err := contractDelegation.Unpack(methodCall, res.Ret)

	return data[0].(bool), err
}
