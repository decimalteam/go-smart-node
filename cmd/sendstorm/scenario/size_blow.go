package scenario

import (
	"fmt"
	"math/rand"
	"time"

	"bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/actions"
	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NFTBlowScenario struct {
	accs []*stormTypes.StormAccount
	api  *dscApi.API
	nfts []dscApi.NFTToken
	rnd  *rand.Rand
}

// 1. create nft
// 2. transfer nft
// 3. query nft, go to 2

func NewNFTBlowScenario(api *dscApi.API, accs []*stormTypes.StormAccount) *NFTBlowScenario {
	return &NFTBlowScenario{
		accs: accs,
		api:  api,
		rnd:  rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (sc *NFTBlowScenario) CreateNFTs(subtokensCount uint32) error {
	for _, acc := range sc.accs {
		acc.UpdateBalance()
		acc.UpdateNumberSequence()
		id := acc.Address()
		msg := dscTx.NewMsgMintToken(
			acc.Account().SdkAddress(),
			"nft_blow_collection",
			id,
			id,
			false,
			acc.Account().SdkAddress(),
			subtokensCount,
			sdk.NewCoin("del", helpers.EtherToWei(sdkmath.NewInt(1))),
		)
		tx, err := dscTx.BuildTransaction(acc.Account(), []sdk.Msg{msg}, "-", "del", sc.api.GetFeeCalculationOptions())
		if err != nil {
			fmt.Printf("CreateNFTs-BuildTransaction err: %s\n", err.Error())
			continue
		}
		err = tx.SignTransaction(acc.Account())
		if err != nil {
			fmt.Printf("CreateNFTs-SignTransaction err: %s\n", err.Error())
			continue
		}
		bz, err := tx.BytesToSend()
		if err != nil {
			fmt.Printf("CreateNFTs-BytesToSend err: %s\n", err.Error())
			continue
		}
		resp, err := sc.api.BroadcastTxSync(bz)
		if err != nil {
			fmt.Printf("CreateNFTs-BytesToSend err: %s\n", err.Error())
			continue
		}
		if resp.Code != 0 {
			fmt.Printf("TxResult code=%d, codespace=%s, msg=%s\n", resp.Code, resp.Codespace, resp.Log)
		}
	}
	return nil
}

func (sc *NFTBlowScenario) UpdateNFT() {
	sc.nfts = make([]dscApi.NFTToken, 0)
	for _, acc := range sc.accs {
		nft, err := sc.api.NFTToken(acc.Address()) // id = acc.Address()
		if err != nil {
			fmt.Printf("get NFT err: %s\n", err.Error())
			continue
		}
		sc.nfts = append(sc.nfts, nft)
	}
}

func (sc *NFTBlowScenario) SendNFT() error {
	sc.UpdateNFT()
	for j, acc := range sc.accs {
		acc.UpdateBalance()
		acc.UpdateNumberSequence()
		var receiver sdk.AccAddress
		if j == 0 {
			receiver = sc.accs[1].Account().SdkAddress()
		} else {
			receiver = sc.accs[j-1].Account().SdkAddress()
		}
		nftsToSend := make([]dscApi.NFTToken, 0)
		for _, nft := range sc.nfts {
			for _, sub := range nft.SubTokens {
				if sub.Owner == acc.Address() {
					nftsToSend = append(nftsToSend, nft)
					break
				}
			}
		}
		if len(nftsToSend) == 0 {
			continue
		}
		nft := actions.RandomChoice(sc.rnd, nftsToSend)
		subTokens := []uint32{}
		for _, sub := range nft.SubTokens {
			if sub.Owner == acc.Address() {
				subTokens = append(subTokens, sub.ID)
			}
		}
		////
		msg := dscTx.NewMsgSendToken(
			acc.Account().SdkAddress(),
			receiver,
			nft.ID,
			subTokens,
		)
		tx, err := dscTx.BuildTransaction(acc.Account(), []sdk.Msg{msg}, "-", "del", sc.api.GetFeeCalculationOptions())
		if err != nil {
			fmt.Printf("SendNFT-BuildTransaction err: %s\n", err.Error())
			continue
		}
		err = tx.SignTransaction(acc.Account())
		if err != nil {
			fmt.Printf("SendNFT-SignTransaction err: %s\n", err.Error())
			continue
		}
		bz, err := tx.BytesToSend()
		if err != nil {
			fmt.Printf("SendNFT-BytesToSend err: %s\n", err.Error())
			continue
		}
		resp, err := sc.api.BroadcastTxSync(bz)
		if err != nil {
			fmt.Printf("SendNFT-BytesToSend err: %s\n", err.Error())
			continue
		}
		if resp.Code != 0 {
			fmt.Printf("TxResult code=%d, codespace=%s, msg=%s\n", resp.Code, resp.Codespace, resp.Log)
		}
	}
	return nil
}
