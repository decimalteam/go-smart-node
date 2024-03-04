// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	types2 "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
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
	var tokenStaked contracts.ContractsStaked
	//var tokenUpdata contracts.TokenReserveUpdated

	for _, log := range recipient.Logs {
		eventValidatorByID, errEvent := validatorMaster.EventByID(log.Topics[0])
		if errEvent == nil {
			if eventValidatorByID.Name == "TokenDeployed" {
				_ = validatorMaster.UnpackIntoInterface(&tokenStaked, eventValidatorByID.Name, log.Data)
			}
		}
		eventDelegationByID, errEvent := delegatorCenter.EventByID(log.Topics[0])
		if errEvent == nil {
			if eventDelegationByID.Name == "Staked1" {
				_ = validatorMaster.UnpackIntoInterface(&tokenStaked, eventDelegationByID.Name, log.Data)

				coinStake, err := k.coinKeeper.GetCoinByDRC(ctx, tokenStaked.Stake.Token.String())
				if err != nil {
					return errors.CoinDoesNotExist
				}

				stake := types.NewStakeCoin(sdk.Coin{Denom: coinStake.Denom, Amount: math.NewIntFromBigInt(tokenStaked.Stake.Amount)})

				cosmosAddress, _ := types2.GetDecimalAddressFromHex(tokenStaked.Stake.Delegator.String())
				cosmosAddressValidator, _ := types2.GetDecimalAddressFromHex(tokenStaked.Stake.Validator.String())

				valAddr, err := sdk.ValAddressFromBech32(cosmosAddressValidator.String())

				validator, found := k.GetValidator(ctx, valAddr)
				if !found {
					return fmt.Errorf("not found validator %s", valAddr)
				}

				_ = k.Delegate(ctx, cosmosAddress, validator, stake)
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
