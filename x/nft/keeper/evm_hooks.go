// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/contracts/nftCenter"
	"bitbucket.org/decimalteam/go-smart-node/contracts/tokenCenter"
	"bitbucket.org/decimalteam/go-smart-node/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
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
	//params := k.GetParams(ctx)
	//if params.TokenCenter == "" {
	//	// no error is returned to avoid reverting the tx and allow for other post
	//	// processing txs to pass and
	//	fmt.Print(params)
	//}

	contractNftCenter, err := contracts.GetAddressFromContractCenter(ctx, k.evmKeeper, contracts.NameOfSlugForGetAddressNftCenter)
	//
	//tokenCenter := center{}
	//fmt.Print(err)
	fmt.Println("nft hooks")
	fmt.Println(contracts.GetContractCenter(ctx.ChainID()))
	fmt.Println(contractNftCenter)
	fmt.Println(err)
	//fmt.Print(tokenCenter)

	// parser for create new or update nft collection
	nftContractCenter, _ := nftCenter.NftCenterMetaData.GetAbi()
	var nftCreated nftCenter.NftCenterNFTCreated

	for _, log := range recipient.Logs {
		eventCenterByID, errEvent := nftContractCenter.EventByID(log.Topics[0])
		if errEvent == nil {
			if eventCenterByID.Name == "NFTCreated" && log.Address.Hex() == contractNftCenter {
				_ = nftContractCenter.UnpackIntoInterface(&nftCreated, eventCenterByID.Name, log.Data)
				creatorAddress, _ := types.GetDecimalAddressFromHex(nftCreated.TokenAddress.Hex())

				// retrieve NFT collection
				collection, collectionExists := k.GetCollection(ctx, creatorAddress, nftCreated.Nft.Symbol)
				if !collectionExists {
					// create NFT collection
					collection = nfttypes.Collection{
						Creator:    creatorAddress.String(),
						Denom:      nftCreated.Nft.Symbol,
						Supply:     0,
						Tokens:     nil,
						TypeNft:    nfttypes.NftType_Unspecified,
						AddressDRC: nftCreated.TokenAddress.String(),
					}
				} else {
					collection.TypeNft = nfttypes.NftType_Unspecified
					collection.AddressDRC = nftCreated.TokenAddress.String()
				}
				if nftCreated.NftType == 0 {
					collection.TypeNft = nfttypes.NftType_NFT721
				}
				if nftCreated.NftType == 1 {
					collection.TypeNft = nfttypes.NftType_NFT1155
				}

				// write collection with it's counter
				k.SetCollection(ctx, collection)
			}
		}
	}
	return nil
}
