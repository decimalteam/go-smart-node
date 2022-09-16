package worker

import (
	"encoding/json"
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type EventMintNFT struct {
	NftID        string      `json:"nftId"`
	Denom        string      `json:"nftCollection"`
	Creator      string      `json:"creator"`
	Owner        string      `json:"owner"`
	Quantity     sdkmath.Int `json:"quantity"`
	StartReserve sdk.Coin    `json:"startReserve"`
	TotalReserve sdk.Coin    `json:"totalReserve"`
	TokenURI     string      `json:"tokenUri"`
	AllowMint    bool        `json:"allowMint"`
	SubTokenIDs  []uint64    `json:"subTokenIds"`
	// from tx
	TxHash  string `json:"txHash"`
	BlockID int64  `json:"blockId"`

	// ??? nonFungible: boolean;
}

type EventTransferNFT struct {
	Sender      string   `json:"sender"`
	Recipient   string   `json:"recipient"`
	NftID       string   `json:"nftId"`
	Denom       string   `json:"nftCollection"`
	SubTokenIDs []uint64 `json:"subTokenIds"`
	TxHash      string   `json:"txHash"`
	BlockID     int64    `json:"blockId"`
}

// decimal.nft.v1.EventMintNFT
func processEventMintNFT(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
			string sender = 1;
		    string recipient = 2;
		    // aka collection
		    string denom = 3;
		    // aka id, token_id
		    string nft_id = 4 [
		        (gogoproto.customname) = "NFTID"
		    ];
		    string token_uri = 5 [
		        (gogoproto.customname) = "TokenURI"
		    ];
		    bool allow_mint = 6;
		    string reserve = 7;
		    repeated uint64 sub_token_ids = 8 [
		        (gogoproto.customname) = "SubTokenIDs"
		    ];
	*/
	var err error
	var emn EventMintNFT
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			emn.Creator = string(attr.Value)
		case "recipient":
			emn.Owner = string(attr.Value)
		case "denom":
			emn.Denom = string(attr.Value)
		case "nft_id":
			emn.NftID = string(attr.Value)
		case "token_uri":
			emn.TokenURI = string(attr.Value)
		case "allow_mint":
			emn.AllowMint = false
			if string(attr.Value) == "true" {
				emn.AllowMint = true
			}
		case "reserve":
			emn.StartReserve, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse reserve '%s': %s", string(attr.Value), err.Error())
			}
		case "sub_token_ids":
			var subIds []string
			err := json.Unmarshal(attr.Value, &subIds)
			if err != nil {
				return fmt.Errorf("can't unmarshal sub_token_ids: %s", err.Error())
			}
			emn.Quantity = sdk.NewInt(int64(len(subIds)))
			for _, s := range subIds {
				v, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return fmt.Errorf("can't parse sub token id '%s': %s", s, err.Error())
				}
				emn.SubTokenIDs = append(emn.SubTokenIDs, v)
			}
		}

	}
	emn.TxHash = txHash
	emn.BlockID = blockId
	emn.TotalReserve = sdk.NewCoin(emn.StartReserve.Denom, emn.StartReserve.Amount.Mul(emn.Quantity))
	ea.addBalanceChange(emn.Creator, emn.TotalReserve.Denom, emn.TotalReserve.Amount.Neg())

	ea.NFTMints = append(ea.NFTMints, emn)
	return nil
}

// decimal.nft.v1.EventTransferNFT
func processEventTransferNFT(ea *EventAccumulator, event abci.Event, txHash string, blockId int64) error {
	/*
		string sender = 1;
		string recipient = 2;
		string denom = 3;
		string nft_id = 4 [
		   (gogoproto.customname) = "NFTID"
		];
		repeated uint64 sub_token_ids = 8 [
		    (gogoproto.customname) = "SubTokenIDs"
		];
	*/
	var etn EventTransferNFT
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			etn.Sender = string(attr.Value)
		case "recipient":
			etn.Recipient = string(attr.Value)
		case "denom":
			etn.Denom = string(attr.Value)
		case "nft_id":
			etn.NftID = string(attr.Value)
		case "sub_token_ids":
			var subIds []string
			err := json.Unmarshal(attr.Value, &subIds)
			if err != nil {
				return fmt.Errorf("can't unmarshal sub_token_ids: %s", err.Error())
			}
			for _, s := range subIds {
				v, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return fmt.Errorf("can't parse sub token id '%s': %s", s, err.Error())
				}
				etn.SubTokenIDs = append(etn.SubTokenIDs, v)
			}
		}
	}
	etn.TxHash = txHash
	etn.BlockID = blockId

	return nil
}
