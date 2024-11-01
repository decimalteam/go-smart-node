// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/contracts/token"
	"bitbucket.org/decimalteam/go-smart-node/contracts/tokenCenter"
	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
	"strings"
)

var _ evmtypes.EvmHooks = Hooks{}

// Hooks wrapper struct for erc20 keeper
type Hooks struct {
	k Keeper
}

type NewToken struct {
	TokenData tokenCenter.DecimalTokenCenterToken `abi:"tokenData"`
}

type ContractCenter struct {
	Address string
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
func (k *Keeper) PostTxProcessing(
	ctx sdk.Context,
	msg core.Message,
	recipient *ethtypes.Receipt,
) error {

	contractTokenCenter, err := contracts.GetAddressFromContractCenter(ctx, k.evmKeeper, contracts.NameOfSlugForGetAddressTokenCenter)

	//fmt.Print(contractTokenCenter)
	tokenContractCenter, _ := tokenCenter.TokenCenterMetaData.GetAbi()
	coinContract, _ := token.TokenMetaData.GetAbi()
	addressWDEL, _ := contracts.GetAddressFromContractCenter(ctx, k.evmKeeper, contracts.NameOfSlugForGetAddressWDEL)
	coinDel, _ := k.GetCoin(ctx, k.GetBaseDenom(ctx))
	if coinDel.DRC20Contract == "" || coinDel.DRC20Contract != addressWDEL {
		_ = k.UpdateCoinDRC(ctx, coinDel.Denom, addressWDEL)
		coinDel.DRC20Contract = addressWDEL
		k.SetCoin(ctx, coinDel)
	}

	// this var is only for new token create from token center
	var tokenAddress tokenCenter.TokenCenterTokenDeployed
	var tokenUpdated token.TokenReserveUpdated

	for _, log := range recipient.Logs {
		eventCoinByID, errEvent := coinContract.EventByID(log.Topics[0])
		if errEvent == nil {
			fmt.Println(eventCoinByID.Name)
			if eventCoinByID.Name == "ReserveUpdated" {
				_ = contracts.UnpackLog(coinContract, &tokenUpdated, eventCoinByID.Name, log)
				_, err = k.GetCoinByDRC(ctx, log.Address.String())
				if err != nil {
					continue
				}
				_ = k.UpdateCoinFromEvent(ctx, tokenUpdated, log.Address.String())
			}
		}
	}

	for _, log := range recipient.Logs {
		if log.Address.String() != contractTokenCenter {
			continue
		}
		eventCenterByID, errEvent := tokenContractCenter.EventByID(log.Topics[0])
		if errEvent == nil {
			fmt.Println(eventCenterByID.Name)
			if eventCenterByID.Name == "TokenDeployed" {
				_ = contracts.UnpackLog(tokenContractCenter, &tokenAddress, eventCenterByID.Name, log)
				fmt.Println(tokenUpdated)
				fmt.Println(tokenAddress)
				if tokenUpdated.NewReserve == nil {
					return fmt.Errorf("reserve is nil")
				}
				fmt.Println(tokenUpdated)
				err = k.CreateCoinEvent(ctx, tokenUpdated.NewReserve, tokenAddress.Meta, tokenAddress.TokenAddress.String())
				if err != nil {
					return status.Error(codes.Internal, err.Error())
				}
			}
		}
	}

	return nil
}

// UpdateCoinFromEvent update reserve and volume by event
func (k *Keeper) UpdateCoinFromEvent(ctx sdk.Context, dataUpdate token.TokenReserveUpdated, tokenAddress string) error {

	// Ensure coin does not exist
	coinExist, err := k.GetCoinByDRC(ctx, tokenAddress)
	if err != nil {
		return nil
	}

	_ = k.UpdateCoinVR(ctx, coinExist.Denom, math.NewIntFromBigInt(dataUpdate.NewSupply), math.NewIntFromBigInt(dataUpdate.NewReserve))

	// Emit transaction events
	_ = events.EmitTypedEvent(ctx, &types.EventUpdateCoinVR{
		Denom:   coinExist.Denom,
		Volume:  math.NewIntFromBigInt(dataUpdate.NewSupply).String(),
		Reserve: math.NewIntFromBigInt(dataUpdate.NewReserve).String(),
	})

	return nil
}

// CreateCoinEvent returns the coin if exists in KVStore.
func (k *Keeper) CreateCoinEvent(ctx sdk.Context, reserve *big.Int, token tokenCenter.DecimalTokenCenterToken, tokenAddress string) error {

	tokenAddress = strings.ToLower(tokenAddress)
	coinDenom := token.Symbol

	if reserve == nil {
		return fmt.Errorf("reserve is nil")
	}

	// Ensure coin does not exist
	coinExist, err := k.GetCoin(ctx, coinDenom)
	if err == nil {
		_ = k.UpdateCoinDRC(ctx, coinDenom, tokenAddress)
		coinExist.DRC20Contract = tokenAddress
		k.SetCoin(ctx, coinExist)
		return nil
	}
	// get authority address
	authAddr := authtypes.NewModuleAddress(cointypes.ModuleName)

	// Create new coin instance
	var coin = types.Coin{
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
