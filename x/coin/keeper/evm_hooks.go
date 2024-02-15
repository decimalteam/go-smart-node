// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

var _ evmtypes.EvmHooks = Hooks{}

// Hooks wrapper struct for erc20 keeper
type Hooks struct {
	k Keeper
}

// Return the wrapper struct
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

	//coinCenter, _ := contracts.TokenCenterMetaData.GetAbi()
	//
	//methodId, err := coinCenter.MethodById(msg.Data)
	//if err != nil {
	//	return nil
	//}
	//
	//type NewToken struct {
	//	TokenData contracts.DecimalTokenCenterToken `abi:"tokenData"`
	//}
	//// Check if processed method
	//switch methodId.Name {
	//case types.DRC20MethodCreateToken:
	//
	//	var tokenAddress contracts.TokenCenterDeployed
	//	for _, log := range recipient.Logs {
	//		eventByID, errEvent := coinCenter.EventByID(log.Topics[0])
	//		if errEvent == nil {
	//			if eventByID.Name == "TokenDeployed" {
	//				_ = coinCenter.UnpackIntoInterface(&tokenAddress, eventByID.Name, log.Data)
	//			}
	//		}
	//	}
	//
	//	var tokenNew NewToken
	//	err = contracts.UnpackInputsData(&tokenNew, methodId.Inputs, msg.Data[4:])
	//	err = k.CreateCoinEvent(ctx, msg.Value, tokenNew.TokenData, tokenAddress.TokenAddress.String())
	//	if err != nil {
	//		return status.Error(codes.Internal, err.Error())
	//	}
	//default:
	//	return nil
	//}

	//for i, log := range receipt.Logs {
	//	// Note: the `Transfer` event contains 3 topics (id, from, to)
	//	if len(log.Topics) != 3 {
	//		continue
	//	}
	//
	//	// Check if event is included in ERC20
	//	eventID := log.Topics[0]
	//	methodId, err := coinCenter.MethodById(eventID)
	//	if err != nil {
	//		continue
	//	}
	//	// Check if event .
	//	switch event.Name {
	//	case types.DRC20EventCreateToken:
	//		_, err = k.CallEVM(ctx, erc20, types.ModuleAddress, contractAddr, true, "burn", tokens)
	//	case types.OWNER_EXTERNAL:
	//		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	//	default:
	//		continue
	//	}
	//
	//	// Check if event is a `Transfer` event.
	//	if event.Name != types.ERC20EventTransfer {
	//		k.Logger(ctx).Info("emitted event", "name", event.Name, "signature", event.Sig)
	//		continue
	//	}
	//
	//	transferEvent, err := erc20.Unpack(event.Name, log.Data)
	//	if err != nil {
	//		k.Logger(ctx).Error("failed to unpack transfer event", "error", err.Error())
	//		continue
	//	}
	//
	//	if len(transferEvent) == 0 {
	//		continue
	//	}
	//
	//	tokens, ok := transferEvent[0].(*big.Int)
	//	// safety check and ignore if amount not positive
	//	if !ok || tokens == nil || tokens.Sign() != 1 {
	//		continue
	//	}
	//
	//	// Check that the contract is a registered token pair
	//	contractAddr := log.Address
	//	id := k.GetERC20Map(ctx, contractAddr)
	//	if len(id) == 0 {
	//		continue
	//	}
	//
	//	pair, found := k.GetTokenPair(ctx, id)
	//	if !found {
	//		continue
	//	}
	//
	//	// Check if tokens are sent to module address
	//	to := common.BytesToAddress(log.Topics[2].Bytes())
	//	if !bytes.Equal(to.Bytes(), types.ModuleAddress.Bytes()) {
	//		continue
	//	}
	//
	//	// Check that conversion for the pair is enabled. Fail
	//	if !pair.Enabled {
	//		// continue to allow transfers for the ERC20 in case the token pair is
	//		// disabled
	//		k.Logger(ctx).Debug(
	//			"ERC20 token -> Cosmos coin conversion is disabled for pair",
	//			"coin", pair.Denom, "contract", pair.Erc20Address,
	//		)
	//		continue
	//	}
	//
	//	// create the corresponding sdk.Coin that is paired with ERC20
	//	coins := sdk.Coins{{Denom: pair.Denom, Amount: sdk.NewIntFromBigInt(tokens)}}
	//
	//	// Perform token conversion. We can now assume that the sender of a
	//	// registered token wants to mint a Cosmos coin.
	//	switch pair.ContractOwner {
	//	case types.OWNER_MODULE:
	//		_, err = k.CallEVM(ctx, erc20, types.ModuleAddress, contractAddr, true, "burn", tokens)
	//	case types.OWNER_EXTERNAL:
	//		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	//	default:
	//		err = types.ErrUndefinedOwner
	//	}
	//
	//	if err != nil {
	//		k.Logger(ctx).Debug(
	//			"failed to process EVM hook for ER20 -> coin conversion",
	//			"coin", pair.Denom, "contract", pair.Erc20Address, "error", err.Error(),
	//		)
	//		continue
	//	}
	//
	//	// Only need last 20 bytes from log.topics
	//	from := common.BytesToAddress(log.Topics[1].Bytes())
	//	recipient := sdk.AccAddress(from.Bytes())
	//
	//	// transfer the tokens from ModuleAccount to sender address
	//	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins); err != nil {
	//		k.Logger(ctx).Debug(
	//			"failed to process EVM hook for ER20 -> coin conversion",
	//			"tx-hash", receipt.TxHash.Hex(), "log-idx", i,
	//			"coin", pair.Denom, "contract", pair.Erc20Address, "error", err.Error(),
	//		)
	//		continue
	//	}
	//}

	return nil
}

// CreateCoinEvent returns the coin if exists in KVStore.
func (k *Keeper) CreateCoinEvent(ctx sdk.Context, reserve *big.Int, token contracts.DecimalTokenCenterToken, tokenAddress string) error {
	k.Logger(ctx).Info("emitted event", token)

	coinDenom := token.Symbol

	var coin = types.Coin{}

	// Ensure coin does not exist
	coinExist, err := k.GetCoin(ctx, coinDenom)
	if err == nil {
		coin = coinExist
	} else {
		// get authority address
		authAddr := authtypes.NewModuleAddress(cointypes.ModuleName)

		// Create new coin instance
		coin = types.Coin{
			Title:         token.Name,
			Denom:         coinDenom,
			CRR:           uint32(token.Crr),
			Reserve:       math.NewIntFromBigInt(reserve),
			Volume:        math.NewIntFromBigInt(token.InitialMint),
			LimitVolume:   math.NewIntFromBigInt(token.MaxTotalSupply),
			MinVolume:     math.NewIntFromBigInt(token.MinTotalSupply),
			Creator:       authAddr.String(),
			Identity:      token.Identity,
			DRC20Contract: tokenAddress,
		}
	}

	// Save coin to the storage
	k.SetCoin(ctx, coin)

	// Emit transaction events
	_ = events.EmitTypedEvent(ctx, &types.EventCreateCoin{
		Sender:               coin.Creator,
		Denom:                coinDenom,
		Title:                coin.Title,
		CRR:                  coin.CRR,
		InitialVolume:        coin.Volume.String(),
		InitialReserve:       coin.Reserve.String(),
		LimitVolume:          coin.LimitVolume.String(),
		MinVolume:            coin.MinVolume.String(),
		Identity:             coin.Identity,
		CommissionCreateCoin: "0",
	})

	return nil
}
