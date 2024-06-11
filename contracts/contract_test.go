package contracts

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts/center"
	"bitbucket.org/decimalteam/go-smart-node/contracts/nft721"
	"bitbucket.org/decimalteam/go-smart-node/contracts/nftCenter"
	"bitbucket.org/decimalteam/go-smart-node/contracts/token"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"os"
	"testing"
)

func TestInitCmd(t *testing.T) {

	web3Client, err := ethclient.Dial("https://devnet-val.decimalchain.com/web3/")
	if err != nil {
		os.Exit(1)
	}

	contractCenter, err := center.NewCenterCaller(
		common.HexToAddress("0xa052da26a526e251db6390834009464ab0398ddc"),
		web3Client,
	)

	address, err := contractCenter.GetAddress(&bind.CallOpts{}, NameOfSlugForGetAddressNftCenter)
	if err != nil {
		return
	}
	fmt.Println(address)
	nftCenterService, err := nftCenter.NewNftCenter(address, web3Client)
	if err != nil {
		fmt.Println(err)
	}
	recipient, err := web3Client.TransactionReceipt(context.Background(), common.HexToHash("0xbd6a379f8cf44bc558db3cd0046860a7c674a6375d8c17d84f54a4857363c854"))
	if err != nil {
		return
	}

	nftContractCenter, _ := nftCenter.NftCenterMetaData.GetAbi()
	nftContract721, _ := nft721.Nft721MetaData.GetAbi()
	tokenContract, _ := token.TokenMetaData.GetAbi()

	var tokenAddress nftCenter.NftCenterNFTCreated
	var tokenReserve token.TokenTransfer
	var nft721Mint nft721.Nft721Transfer

	for _, log := range recipient.Logs {
		eventCenterByID, errEvent := nftContractCenter.EventByID(log.Topics[0])
		if errEvent == nil {
			if eventCenterByID.Name == "NFTCreated" {
				_ = nftContractCenter.UnpackIntoInterface(&tokenAddress, eventCenterByID.Name, log.Data)
				//tokenAddress.TokenAddress = common.HexToAddress(log.Topics[1].Hex())
				fmt.Println(tokenAddress)
				//err = k.CreateCoinEvent(ctx, tokenUpdated.NewReserve, tokenAddress.Meta, tokenAddress.TokenAddress.String())
				//if err != nil {
				//	return status.Error(codes.Internal, err.Error())
				//}
			}
		}
		eventTokenByID, errEvent := tokenContract.EventByID(log.Topics[0])
		if errEvent == nil {
			if eventTokenByID.Name == "Transfer" {

				_ = UnpackLog(tokenContract, &tokenReserve, eventTokenByID.Name, log)
				fmt.Println(tokenReserve)
				//err = k.CreateCoinEvent(ctx, tokenUpdated.NewReserve, tokenAddress.Meta, tokenAddress.TokenAddress.String())
				//if err != nil {
				//	return status.Error(codes.Internal, err.Error())
				//}
			}
		}
		event721ByID, errEvent := nftContract721.EventByID(log.Topics[0])
		if errEvent == nil {
			if event721ByID.Name == "Transfer" {
				_ = UnpackLog(nftContract721, &nft721Mint, event721ByID.Name, log)
				//nft721Mint.From = common.HexToAddress(log.Topics[1].Hex())
				//nft721Mint.To = common.HexToAddress(log.Topics[2].Hex())
				//
				//_ = abi.ParseTopics(nft721Mint, indexed, log.Topics[1:])

				fmt.Println(log.Topics[0].Hex())
				fmt.Println(log.Topics[1].Hex())
				fmt.Println(log.Topics[2].Hex())
				fmt.Println(nft721Mint)
				fmt.Println(nft721Mint.From)
				fmt.Println(nft721Mint.To)
				fmt.Println(nft721Mint.TokenId.String())
				//err = k.CreateCoinEvent(ctx, tokenUpdated.NewReserve, tokenAddress.Meta, tokenAddress.TokenAddress.String())
				//if err != nil {
				//	return status.Error(codes.Internal, err.Error())
				//}
			}
		}
	}
	//fmt.Println(hash)
	//nft721.P
	fmt.Println(nftCenterService)
}
