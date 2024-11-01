package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/contracts/center"
	"bitbucket.org/decimalteam/go-smart-node/contracts/tokenCenter"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
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
	res, err := k.evmKeeper.CallEVM(ctx, *contractCenter, common.Address(types.ModuleAddress), contract, false, methodCall, contracts.NameOfSlugForGetAddressWDEL)
	if err != nil {
		return new(common.Address).Hex(), err
	}
	data, err := contractCenter.Unpack(methodCall, res.Ret)
	if len(data) == 0 {
		return new(common.Address).Hex(), err
	}
	return data[0].(common.Address).String(), err
}

// QueryAddressTokenCenter returns the data of a deployed ERC20 contract
func (k *Keeper) QueryAddressTokenCenter(
	ctx sdk.Context,
	contract common.Address,
) (string, error) {

	contractCenter, _ := tokenCenter.TokenCenterMetaData.GetAbi()
	methodCall := "getAddress"
	// Address token center
	res, err := k.evmKeeper.CallEVM(ctx, *contractCenter, common.Address(types.ModuleAddress), contract, false, methodCall, contracts.NameOfSlugForGetAddressTokenCenter)
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

// BalanceOf queries an account's balance for a given ERC20 contract
func (k *Keeper) BalanceOf(
	ctx sdk.Context,
	abi abi.ABI,
	contract, account common.Address,
) *big.Int {
	res, err := k.evmKeeper.CallEVM(ctx, abi, common.Address(types.ModuleAddress), contract, false, "balanceOf", account)
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
