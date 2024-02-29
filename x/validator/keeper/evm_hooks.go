// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
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

	// this var is only for new token create from token center
	var tokenAddress contracts.TokenCenterDeployed
	var tokenUpdata contracts.TokenReserveUpdated

	for _, log := range recipient.Logs {
		eventCenterByID, errEvent := validatorMaster.EventByID(log.Topics[0])
		if errEvent == nil {
			if eventCenterByID.Name == "TokenDeployed" {
				_ = validatorMaster.UnpackIntoInterface(&tokenAddress, eventCenterByID.Name, log.Data)
			}
		}
		eventCoinByID, errEvent := validatorMaster.EventByID(log.Topics[0])
		if errEvent == nil {
			if eventCoinByID.Name == "ReserveUpdated" {
				_ = contracts.UnpackInputsData(&tokenUpdata, eventCoinByID.Inputs, log.Data)
			}
		}
	}

	methodId, err := validatorMaster.MethodById(msg.Data)
	if err != nil {
		return nil
	}
	// Check if processed method
	switch methodId.Name {
	case types.ContractMethodCreateValidator:

		//var tokenNew NewToken
		//err = contracts.UnpackInputsData(&tokenNew, methodId.Inputs, msg.Data[4:])
		//
		//err = k.CreateCoinEvent(ctx, msg.Value, tokenNew.TokenData, tokenAddress.TokenAddress.String())
		//if err != nil {
		//	return status.Error(codes.Internal, err.Error())
		//}
	default:
		return nil
	}

	return nil
}
