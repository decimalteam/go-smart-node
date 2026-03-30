package keeper

import (
	"math/big"
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
// auto-unbond entries for async processing by the external bot.
// It enqueues the base (non-held) stake with holdTimestamp=0 and each
// hold separately with its own holdTimestamp.
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
	coin, err := k.coinKeeper.GetCoin(ctx, del.GetStake().GetStake().Denom)
	if err != nil {
		return err
	}

	contractDelegation, _ := delegation.DelegationMetaData.GetAbi()
	evmVal := common.BytesToAddress(valAddr.Bytes())
	evmDel := common.BytesToAddress(delAddr.Bytes())
	evmToken := common.HexToAddress(coin.DRC20Contract)
	evmDelegationAddr := common.HexToAddress(delegationAddress)

	totalAmount := del.GetStake().GetStake().Amount
	holds := del.GetStake().GetHolds()

	// Enqueue each hold with its own holdTimestamp
	for _, hold := range holds {
		holdTS := big.NewInt(hold.HoldEndTime)
		_, err = k.evmKeeper.CallEVM(ctx, *contractDelegation,
			common.Address{},
			evmDelegationAddr,
			true,
			"autoUnbondEnqueue",
			evmVal, evmDel, hold.Amount.BigInt(), evmToken, holdTS,
		)
		if err != nil {
			return err
		}
		totalAmount = totalAmount.Sub(hold.Amount)
	}

	// Enqueue the base (non-held) portion with holdTimestamp=0
	if totalAmount.IsPositive() {
		_, err = k.evmKeeper.CallEVM(ctx, *contractDelegation,
			common.Address{},
			evmDelegationAddr,
			true,
			"autoUnbondEnqueue",
			evmVal, evmDel, totalAmount.BigInt(), evmToken, big.NewInt(0),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
