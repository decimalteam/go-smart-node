package contracts

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts/center"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegation"
	"bitbucket.org/decimalteam/go-smart-node/contracts/nftCenter"
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
		common.HexToAddress("0xe5268fd6a4d041f20cbb92c662ceff1efe4c861e"),
		web3Client,
	)

	address, err := contractCenter.GetContractAddress(&bind.CallOpts{}, NameOfSlugForGetAddressNftCenter)
	if err != nil {
		return
	}
	fmt.Println(address)
	nftCenterService, err := nftCenter.NewNftCenter(address, web3Client)
	if err != nil {
		fmt.Println(err)
	}
	recipient, err := web3Client.TransactionReceipt(context.Background(), common.HexToHash("0x60819d8d2ec33e8de31613c0a46cb6c13990eb56b083e5b29eb17bebc1c05c4e"))
	if err != nil {
		return
	}

	delegatorCenter, _ := delegation.DelegationMetaData.GetAbi()
	//nftContractCenter, _ := nftCenter.NftCenterMetaData.GetAbi()
	//nftContract721, _ := nft721.Nft721MetaData.GetAbi()
	//tokenContract, _ := token.TokenMetaData.GetAbi()
	//validatorContract, _ := validator.ValidatorMetaData.GetAbi()

	var tokenDelegate delegation.DelegationStakeUpdated
	//var tokenAddress nftCenter.NftCenterNFTCreated
	//var tokenReserve token.TokenTransfer
	//var nft721Mint nft721.Nft721Transfer
	//var tokenDelegationUser validator.ValidatorValidatorMetaUpdated
	//var validatorMeta validator.ValidatorValidatorMetaUpdated

	for _, log := range recipient.Logs {
		eventDelegationByID, errEvent := delegatorCenter.EventByID(log.Topics[0])
		if errEvent == nil {
			fmt.Println(eventDelegationByID.Name)
			//if eventDelegationByID.Name == "StakeUpdated" && stakeUpdate == 0 {
			//	if tokenDelegationAmount.ChangedAmount == nil {
			//		return errors.DelegationSumIsNotSet
			//	}
			_ = UnpackLog(delegatorCenter, &tokenDelegate, eventDelegationByID.Name, log)
			//	stakeUpdate = stakeUpdate + 1
			//	_, err := k.coinKeeper.GetCoinByDRC(ctx, tokenDelegate.Stake.Token.String())
			//	if err != nil {
			//		symbolToken, _ := k.QuerySymbolToken(ctx, tokenDelegate.Stake.Token)
			//		coinUpdate, err := k.coinKeeper.GetCoin(ctx, symbolToken)
			//		if err == nil {
			//			_ = k.coinKeeper.UpdateCoinDRC(ctx, symbolToken, tokenDelegate.Stake.Token.String())
			//			coinUpdate.DRC20Contract = tokenDelegate.Stake.Token.String()
			//			k.coinKeeper.SetCoin(ctx, coinUpdate)
			//		}
			//	}
			//	tokenDelegate.Stake.Amount = tokenDelegationAmount.ChangedAmount
			//	err = k.Staked(ctx, tokenDelegate)
			//	if err != nil {
			//		return err
			//	}
			//}
			//
			//if eventDelegationByID.Name == "RequestWithdraw" {
			//	_ = delegatorCenter.UnpackIntoInterface(&tokenUndelegate, eventDelegationByID.Name, log.Data)
			//	_, err := k.coinKeeper.GetCoinByDRC(ctx, tokenDelegate.Stake.Token.String())
			//	if err != nil {
			//		symbolToken, _ := k.QuerySymbolToken(ctx, tokenDelegate.Stake.Token)
			//		coinUpdate, err := k.coinKeeper.GetCoin(ctx, symbolToken)
			//		if err == nil {
			//			_ = k.coinKeeper.UpdateCoinDRC(ctx, symbolToken, tokenDelegate.Stake.Token.String())
			//			coinUpdate.DRC20Contract = tokenDelegate.Stake.Token.String()
			//			k.coinKeeper.SetCoin(ctx, coinUpdate)
			//		}
			//	}
			//	err = k.RequestWithdraw(ctx, tokenUndelegate)
			//	if err != nil {
			//		return err
			//	}
			//}
			//if eventDelegationByID.Name == "RequestTransfer" {
			//	_ = delegatorCenter.UnpackIntoInterface(&tokenRedelegation, eventDelegationByID.Name, log.Data)
			//	_, err := k.coinKeeper.GetCoinByDRC(ctx, tokenDelegate.Stake.Token.String())
			//	if err != nil {
			//		symbolToken, _ := k.QuerySymbolToken(ctx, tokenDelegate.Stake.Token)
			//		coinUpdate, err := k.coinKeeper.GetCoin(ctx, symbolToken)
			//		if err == nil {
			//			_ = k.coinKeeper.UpdateCoinDRC(ctx, symbolToken, tokenDelegate.Stake.Token.String())
			//			coinUpdate.DRC20Contract = tokenDelegate.Stake.Token.String()
			//			k.coinKeeper.SetCoin(ctx, coinUpdate)
			//		}
			//	}
			//	err = k.RequestTransfer(ctx, tokenRedelegation)
			//	if err != nil {
			//		return err
			//	}
			//}
		}
		//eventValidatorByID, errEvent := validatorContract.EventByID(log.Topics[0])
		//if errEvent == nil {
		//	if eventValidatorByID.Name == "ValidatorMetaUpdated" {
		//		_ = UnpackLog(validatorContract, &validatorMeta, eventValidatorByID.Name, log)
		//		//tokenAddress.TokenAddress = common.HexToAddress(log.Topics[1].Hex())
		//		fmt.Println(validatorMeta)
		//		var validatorInfo MasterValidatorValidatorAddedMeta
		//		_ = json.Unmarshal([]byte(validatorMeta.Meta), &validatorInfo)
		//		fmt.Println(validatorMeta)
		//		//err = k.CreateCoinEvent(ctx, tokenUpdated.NewReserve, tokenAddress.Meta, tokenAddress.TokenAddress.String())
		//		//if err != nil {
		//		//	return status.Error(codes.Internal, err.Error())
		//		//}
		//	}
		//}
		//eventDelegationByID, errEvent := tokenDelegation.EventByID(log.Topics[0])
		//if errEvent == nil {
		//	if eventDelegationByID.Name == "StakeUpdated" {
		//		_ = UnpackLog(tokenDelegation, &tokenDelegationUser, eventDelegationByID.Name, log)
		//		//tokenAddress.TokenAddress = common.HexToAddress(log.Topics[1].Hex())
		//		fmt.Println(tokenDelegationUser)
		//		fmt.Println(tokenDelegationUser)
		//		//err = k.CreateCoinEvent(ctx, tokenUpdated.NewReserve, tokenAddress.Meta, tokenAddress.TokenAddress.String())
		//		//if err != nil {
		//		//	return status.Error(codes.Internal, err.Error())
		//		//}
		//	}
		//}
		//eventCenterByID, errEvent := nftContractCenter.EventByID(log.Topics[0])
		//if errEvent == nil {
		//	if eventCenterByID.Name == "NFTCreated" {
		//		_ = nftContractCenter.UnpackIntoInterface(&tokenAddress, eventCenterByID.Name, log.Data)
		//		//tokenAddress.TokenAddress = common.HexToAddress(log.Topics[1].Hex())
		//		fmt.Println(tokenAddress)
		//		//err = k.CreateCoinEvent(ctx, tokenUpdated.NewReserve, tokenAddress.Meta, tokenAddress.TokenAddress.String())
		//		//if err != nil {
		//		//	return status.Error(codes.Internal, err.Error())
		//		//}
		//	}
		//}
		//eventTokenByID, errEvent := tokenContract.EventByID(log.Topics[0])
		//if errEvent == nil {
		//	if eventTokenByID.Name == "Transfer" {
		//
		//		_ = UnpackLog(tokenContract, &tokenReserve, eventTokenByID.Name, log)
		//		fmt.Println(tokenReserve)
		//		//err = k.CreateCoinEvent(ctx, tokenUpdated.NewReserve, tokenAddress.Meta, tokenAddress.TokenAddress.String())
		//		//if err != nil {
		//		//	return status.Error(codes.Internal, err.Error())
		//		//}
		//	}
		//}
		//event721ByID, errEvent := nftContract721.EventByID(log.Topics[0])
		//if errEvent == nil {
		//	if event721ByID.Name == "Transfer" {
		//		_ = UnpackLog(nftContract721, &nft721Mint, event721ByID.Name, log)
		//		//nft721Mint.From = common.HexToAddress(log.Topics[1].Hex())
		//		//nft721Mint.To = common.HexToAddress(log.Topics[2].Hex())
		//		//
		//		//_ = abi.ParseTopics(nft721Mint, indexed, log.Topics[1:])
		//
		//		fmt.Println(log.Topics[0].Hex())
		//		fmt.Println(log.Topics[1].Hex())
		//		fmt.Println(log.Topics[2].Hex())
		//		fmt.Println(nft721Mint)
		//		fmt.Println(nft721Mint.From)
		//		fmt.Println(nft721Mint.To)
		//		fmt.Println(nft721Mint.TokenId.String())
		//		//err = k.CreateCoinEvent(ctx, tokenUpdated.NewReserve, tokenAddress.Meta, tokenAddress.TokenAddress.String())
		//		//if err != nil {
		//		//	return status.Error(codes.Internal, err.Error())
		//		//}
		//	}
		//}
	}
	//fmt.Println(hash)
	//nft721.P
	fmt.Println(nftCenterService)
}
