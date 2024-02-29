// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"encoding/json"
	"fmt"
	"github.com/decimalteam/ethermint/server/config"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
)

// QueryAddressTokenCenter returns the data of a deployed ERC20 contract
func (k Keeper) QueryAddressTokenCenter(
	ctx sdk.Context,
	contract common.Address,
) (string, error) {

	contractCenter, _ := contracts.ContractCenterMetaData.GetAbi()

	// Name
	_, _ = k.CallEVM(ctx, *contractCenter, common.Address(types.ModuleAddress), contract, false, "getAddress", types.NameOfSlugForGetAddressTokenCenter)
	//if err != nil {
	//	return types.ERC20Data{}, err
	//}
	//
	//if err := erc20.UnpackIntoInterface(&nameRes, "name", res.Ret); err != nil {
	//	return types.ERC20Data{}, errorsmod.Wrapf(
	//		types.ErrABIUnpack, "failed to unpack name: %s", err.Error(),
	//	)
	//}

	return "", nil
}

// BalanceOf queries an account's balance for a given ERC20 contract
func (k Keeper) BalanceOf(
	ctx sdk.Context,
	abi abi.ABI,
	contract, account common.Address,
) *big.Int {
	res, err := k.CallEVM(ctx, abi, common.Address(types.ModuleAddress), contract, false, "balanceOf", account)
	if err != nil {
		return nil
	}

	unpacked, err := abi.Unpack("balanceOf", res.Ret)
	if err != nil || len(unpacked) == 0 {
		return nil
	}

	balance, ok := unpacked[0].(*big.Int)
	if !ok {
		return nil
	}

	return balance
}

// CallEVM performs a smart contract method call using given args
func (k Keeper) CallEVM(
	ctx sdk.Context,
	abi abi.ABI,
	from, contract common.Address,
	commit bool,
	method string,
	args ...interface{},
) (*evmtypes.MsgEthereumTxResponse, error) {
	data, _ := abi.Pack(method, args...)
	//if err != nil {
	//	return nil, errorsmod.Wrap(
	//		types.ErrABIPack,
	//		errorsmod.Wrap(err, "failed to create transaction data").Error(),
	//	)
	//}

	argsCall, err := json.Marshal(evmtypes.TransactionArgs{
		From: nil,
		To:   &contract,
		Data: (*hexutil.Bytes)(&data),
	})
	if err != nil {
		return nil, errorsmod.Wrapf(errortypes.ErrJSONMarshal, "failed to marshal tx args: %s", err.Error())
	}
	resp, err := k.evmKeeper.EthCall(sdk.WrapSDKContext(ctx), &evmtypes.EthCallRequest{
		Args:   argsCall,
		GasCap: config.DefaultGasCap,
	})
	fmt.Print(err)
	//resp, _ := k.CallEVMWithData(ctx, from, &contract, data, commit)
	//if err != nil {
	//	return nil, errorsmod.Wrapf(err, "contract call failed: method '%s', contract '%s'", method, contract)
	//}
	return resp, nil
}

// CallEVMWithData performs a smart contract method call using contract data
func (k Keeper) CallEVMWithData(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	data []byte,
	commit bool,
) (*evmtypes.MsgEthereumTxResponse, error) {
	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return nil, err
	}

	gasCap := config.DefaultGasCap
	if commit {
		args, err := json.Marshal(evmtypes.TransactionArgs{
			From: &from,
			To:   contract,
			Data: (*hexutil.Bytes)(&data),
		})
		if err != nil {
			return nil, errorsmod.Wrapf(errortypes.ErrJSONMarshal, "failed to marshal tx args: %s", err.Error())
		}

		gasRes, err := k.evmKeeper.EstimateGas(sdk.WrapSDKContext(ctx), &evmtypes.EthCallRequest{
			Args:   args,
			GasCap: config.DefaultGasCap,
		})
		if err != nil {
			return nil, err
		}
		gasCap = gasRes.Gas
	}

	msg := core.Message{
		From:              from,
		Nonce:             nonce,
		GasLimit:          gasCap,
		GasPrice:          big.NewInt(0),
		GasFeeCap:         big.NewInt(0),
		GasTipCap:         big.NewInt(0),
		To:                contract,
		Value:             big.NewInt(0),
		Data:              data,
		AccessList:        ethtypes.AccessList{},
		SkipAccountChecks: true,
	}

	res, err := k.evmKeeper.ApplyMessage(ctx, msg, evmtypes.NewNoOpTracer(), commit)
	fmt.Print(err)
	//if err != nil {
	//	return nil, err
	//}

	//if res.Failed() {
	//	return nil, errorsmod.Wrap(evmtypes.ErrVMExecution, res.VmError)
	//}

	return res, nil
}
