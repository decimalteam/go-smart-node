package keeper

import (
	"time"

	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegation"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// ProcessAutoUnbond checks all offline validators and enqueues their
// delegations for auto-unbond in the EVM delegation contract.
func (k Keeper) ProcessAutoUnbond(ctx sdk.Context) {
	timeout := k.AutoUnbondTimeout(ctx)

	// Feature disabled
	if timeout == 0 {
		return
	}

	cutoff := ctx.BlockTime().Add(-timeout)

	k.IterateValidatorOfflineSince(ctx, func(valAddr sdk.ValAddress, offlineSince time.Time) bool {
		if offlineSince.After(cutoff) {
			return false // not yet expired
		}

		k.enqueueAutoUnbond(ctx, valAddr)
		k.DeleteValidatorOfflineSince(ctx, valAddr)
		return false
	})
}

// enqueueAutoUnbond writes all delegations for a validator into the
// EVM contract's auto-unbond queue.
func (k Keeper) enqueueAutoUnbond(ctx sdk.Context, valAddr sdk.ValAddress) {
	delegations := k.GetValidatorDelegations(ctx, valAddr)

	for _, del := range delegations {
		err := k.ExecuteAutoUnbondEnqueue(ctx, del)
		if err != nil {
			ctx.Logger().Error("auto-unbond: enqueue failed",
				"validator", valAddr,
				"delegator", del.Delegator,
				"err", err)
			continue
		}
	}

	ctx.Logger().Info("auto-unbond: enqueued delegations for withdrawal",
		"validator", valAddr,
		"delegations_count", len(delegations),
	)
}

// ExecuteAutoUnbondEnqueue calls the EVM delegation contract to enqueue
// an auto-unbond entry. Follows the ExecuteForceWithdrawal pattern.
func (k *Keeper) ExecuteAutoUnbondEnqueue(ctx sdk.Context, del types.Delegation) error {
	delegationAddress, err := contracts.GetAddressFromContractCenter(
		ctx, k.evmKeeper, contracts.NameOfSlugForGetAddressDelegation)
	if err != nil {
		return err
	}

	valAddr, err := sdk.ValAddressFromBech32(del.Validator)
	if err != nil {
		return err
	}
	delAddr, err := sdk.AccAddressFromBech32(del.Delegator)
	if err != nil {
		return err
	}
	amount := del.GetStake().GetStake().Amount.BigInt()
	coin, err := k.coinKeeper.GetCoin(ctx, del.GetStake().GetStake().Denom)
	if err != nil {
		return err
	}

	contractDelegation, _ := delegation.DelegationMetaData.GetAbi()
	_, err = k.evmKeeper.CallEVM(ctx, *contractDelegation,
		common.Address{}, // msg.sender == address(0)
		common.HexToAddress(delegationAddress),
		true, // commit
		"autoUnbondEnqueue",
		common.BytesToAddress(valAddr.Bytes()),
		common.BytesToAddress(delAddr.Bytes()),
		amount,
		common.HexToAddress(coin.DRC20Contract),
	)
	return err
}
