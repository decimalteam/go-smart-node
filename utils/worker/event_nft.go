package worker

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"encoding/json"
	"fmt"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var reservedPool = authtypes.NewModuleAddress(types.ReservedPool)

// CreateOrUpdate in postgres
type EventUpdateCollection struct {
	Creator string `json:"creator"`
	Denom   string `json:"nftCollection"`
	Supply  uint32 `json:"supply"`
	// from tx
	TxHash string `json:"txHash"`
}

type EventCreateToken struct {
	NftID         string   `json:"nftId"`
	NftCollection string   `json:"nftCollection"`
	TokenURI      string   `json:"tokenUri"`
	Creator       string   `json:"creator"`
	StartReserve  sdk.Coin `json:"startReserve"`
	TotalReserve  sdk.Coin `json:"totalReserve"`
	AllowMint     bool     `json:"allowMint"`
	Recipient     string   `json:"recipient"`
	Quantity      uint32   `json:"quantity"`
	SubTokenIDs   []uint32 `json:"subTokenIds"`
	// from tx
	TxHash string `json:"txHash"`
}

type EventMintToken struct {
	Creator       string   `json:"creator"`
	Recipient     string   `json:"recipient"`
	NftCollection string   `json:"nftCollection"`
	NftID         string   `json:"nftId"`
	StartReserve  sdk.Coin `json:"startReserve"`
	SubTokenIDs   []uint32 `json:"subTokenIds"`
	// from tx
	TxHash string `json:"txHash"`
}

type EventBurnToken struct {
	Sender      string   ` json:"sender"`
	NftID       string   `json:"nftId"`
	SubTokenIDs []uint32 `json:"subTokenIds"`
	// from tx
	TxHash string `json:"txHash"`
}

type EventUpdateToken struct {
	//Sender   string ` json:"sender"`
	NftID    string `json:"nftId"`
	TokenURI string `json:"tokenUri"`
	// from tx
	TxHash string `json:"txHash"`
}

type EventUpdateReserve struct {
	//Sender      string   ` json:"sender"`
	NftID       string   `json:"nftId"`
	Reserve     sdk.Coin `json:"reserve"`
	Refill      sdk.Coin `json:"refill"`
	SubTokenIDs []uint32 `json:"subTokenIds"`
	// from tx
	TxHash string `json:"txHash"`
}

type EventSendToken struct {
	Sender      string   ` json:"sender"`
	NftID       string   `json:"nftId"`
	Recipient   string   `json:"recipient"`
	SubTokenIDs []uint32 `json:"subTokenIds"`
	// from tx
	TxHash string `json:"txHash"`
}

func processEventCreateCollection(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	  string creator = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string denom = 2;
	  uint32 supply = 3;
	*/
	//var err error
	//var e EventUpdateCollection
	//for _, attr := range event.Attributes {
	//	switch string(attr.Key) {
	//	case "creator":
	//		e.Creator = string(attr.Value)
	//	case "denom":
	//		e.Denom = string(attr.Value)
	//	case "supply":
	//		e.Supply = binary.LittleEndian.Uint32(attr.Value)
	//	}
	//}
	//e.TxHash = txHash
	//
	//ea.Collection = append(ea.Collection, e)
	return nil
}

func processEventCreateToken(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
		string creator = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
		string denom = 2;
		string id = 3 [ (gogoproto.customname) = "ID" ];
		string uri = 4 [ (gogoproto.customname) = "URI" ];
		bool allowMint = 5;
		string reserve = 6;
		string recipient = 7;
		repeated uint32 subTokenIds = 8 [ (gogoproto.customname) = "SubTokenIDs" ];
	*/
	var err error
	var e EventCreateToken
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "creator":
			e.Creator = string(attr.Value)
		case "denom":
			e.NftCollection = string(attr.Value)
		case "id":
			e.NftID = string(attr.Value)
		case "uri":
			e.TokenURI = string(attr.Value)
		case "allowMint":
			e.AllowMint = false
			if string(attr.Value) == "true" {
				e.AllowMint = true
			}
		case "reserve":
			e.StartReserve, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse reserve '%s': %s", string(attr.Value), err.Error())
			}
		case "subTokenIds":
			var subIds []string
			err := json.Unmarshal(attr.Value, &subIds)
			if err != nil {
				return fmt.Errorf("can't unmarshal subTokenIds: %s", err.Error())
			}
			for _, s := range subIds {
				v, err := strconv.ParseUint(s, 10, 32)
				if err != nil {
					return fmt.Errorf("can't parse sub token id '%s': %s", s, err.Error())
				}
				e.SubTokenIDs = append(e.SubTokenIDs, uint32(v))
			}
		}
	}

	// TODO возможно стоит убрать поля которые есть в mint из create token
	e.TxHash = txHash

	e.TotalReserve = sdk.NewCoin(e.StartReserve.Denom, e.StartReserve.Amount.Mul(sdk.NewInt(int64(len(e.SubTokenIDs)))))
	e.Quantity = uint32(len(e.SubTokenIDs))
	ea.addBalanceChange(e.Creator, e.TotalReserve.Denom, e.TotalReserve.Amount.Neg())
	//ea.addBalanceChange(reservedPool.String(), e.TotalReserve.Denom, e.TotalReserve.Amount)
	ea.addMintSubTokens(EventMintToken{
		Creator:       e.Creator,
		Recipient:     e.Recipient,
		NftCollection: e.NftCollection,
		NftID:         e.NftID,
		StartReserve:  e.StartReserve,
		SubTokenIDs:   e.SubTokenIDs,
		TxHash:        txHash,
	})

	ea.CreateToken = append(ea.CreateToken, e)
	return nil
}

func processEventMintNFT(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	  string creator = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string denom = 2;
	  string id = 3 [ (gogoproto.customname) = "ID" ];
	  string reserve = 4;
	  string recipient = 5;
	  repeated uint32 sub_token_ids = 6 [ (gogoproto.customname) = "SubTokenIDs" ];
	*/
	var err error
	var e EventMintToken
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "creator":
			e.Creator = string(attr.Value)
		case "denom":
			e.NftCollection = string(attr.Value)
		case "id":
			e.NftID = string(attr.Value)
		case "reserve":
			e.StartReserve, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse reserve '%s': %s", string(attr.Value), err.Error())
			}
		case "sub_token_ids":
			var subIds []string
			err := json.Unmarshal(attr.Value, &subIds)
			if err != nil {
				return fmt.Errorf("can't unmarshal subTokenIds: %s", err.Error())
			}
			for _, s := range subIds {
				v, err := strconv.ParseUint(s, 10, 32)
				if err != nil {
					return fmt.Errorf("can't parse sub token id '%s': %s", s, err.Error())
				}
				e.SubTokenIDs = append(e.SubTokenIDs, uint32(v))
			}
		}
	}
	e.TxHash = txHash

	totalReserve := sdk.NewCoin(e.StartReserve.Denom, e.StartReserve.Amount.Mul(sdk.NewInt(int64(len(e.SubTokenIDs)))))
	ea.addBalanceChange(e.Creator, totalReserve.Denom, totalReserve.Amount.Neg())
	//ea.addBalanceChange(reservedPool.String(), totalReserve.Denom, totalReserve.Amount)

	ea.addMintSubTokens(e)

	return nil
}

func processEventBurnNFT(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	  string sender = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string id = 2 [ (gogoproto.customname) = "ID" ];
	  string return = 3;  // coin that was returned in total per transaction for all NFT sub-tokens
	  repeated uint32 sub_token_ids = 4 [ (gogoproto.customname) = "SubTokenIDs" ];
	*/
	var (
		err         error
		returnCoins sdk.Coin
		e           EventBurnToken
	)
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			e.Sender = string(attr.Value)
		case "id":
			e.NftID = string(attr.Value)
		case "return":
			returnCoins, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse reserve '%s': %s", string(attr.Value), err.Error())
			}
		case "sub_token_ids":
			var subIds []string
			err := json.Unmarshal(attr.Value, &subIds)
			if err != nil {
				return fmt.Errorf("can't unmarshal subTokenIds: %s", err.Error())
			}
			for _, s := range subIds {
				v, err := strconv.ParseUint(s, 10, 32)
				if err != nil {
					return fmt.Errorf("can't parse sub token id '%s': %s", s, err.Error())
				}
				e.SubTokenIDs = append(e.SubTokenIDs, uint32(v))
			}
		}
	}
	e.TxHash = txHash

	ea.addBalanceChange(e.Sender, returnCoins.Denom, returnCoins.Amount)
	//ea.addBalanceChange(reservedPool.String(), returnCoins.Denom, returnCoins.Amount.Neg())
	ea.addBurnSubTokens(e)

	return nil
}

func processEventUpdateToken(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	  string sender = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string id = 2 [ (gogoproto.customname) = "ID" ];
	  string uri = 3 [ (gogoproto.customname) = "URI" ];
	*/

	//var err error
	var e EventUpdateToken
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "id":
			e.NftID = string(attr.Value)
		case "uri":
			e.TokenURI = string(attr.Value)
		}
	}

	ea.UpdateToken = append(ea.UpdateToken, e)
	return nil
}

func processEventUpdateReserve(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	  string sender = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string id = 2 [ (gogoproto.customname) = "ID" ];
	  string reserve = 3; // coin that defines new reserve for all updating NFT-subtokens
	  string refill = 4;  // coin that was added in total per transaction for all NFT sub-tokens
	  repeated uint32 sub_token_ids = 5 [ (gogoproto.customname) = "SubTokenIDs" ];
	*/
	var (
		sender string
		err    error
		e      EventUpdateReserve
	)
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			sender = string(attr.Value)
		case "id":
			e.NftID = string(attr.Value)
		case "reserve":
			e.Reserve, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse reserve '%s': %s", string(attr.Value), err.Error())
			}
		case "refill":
			e.Refill, err = sdk.ParseCoinNormalized(string(attr.Value))
			if err != nil {
				return fmt.Errorf("can't parse reserve '%s': %s", string(attr.Value), err.Error())
			}
		case "sub_token_ids":
			var subIds []string
			err := json.Unmarshal(attr.Value, &subIds)
			if err != nil {
				return fmt.Errorf("can't unmarshal subTokenIds: %s", err.Error())
			}
			for _, s := range subIds {
				v, err := strconv.ParseUint(s, 10, 32)
				if err != nil {
					return fmt.Errorf("can't parse sub token id '%s': %s", s, err.Error())
				}
				e.SubTokenIDs = append(e.SubTokenIDs, uint32(v))
			}
		}
	}

	ea.addBalanceChange(sender, e.Refill.Denom, e.Refill.Amount.Neg())
	//ea.addBalanceChange(reservedPool.String(), e.Refill.Denom, e.Refill.Amount)
	ea.UpdateReserve = append(ea.UpdateReserve, e)

	return nil
}

func processEventSendNFT(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	  string sender = 1 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  string id = 2 [ (gogoproto.customname) = "ID" ];
	  string recipient = 3 [ (cosmos_proto.scalar) = "cosmos.AddressString" ];
	  repeated uint32 sub_token_ids = 4 [ (gogoproto.customname) = "SubTokenIDs" ];
	*/

	var e EventSendToken
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			e.Sender = string(attr.Value)
		case "id":
			e.NftID = string(attr.Value)
		case "recipient":
			e.Recipient = string(attr.Value)
		case "sub_token_ids":
			var subIds []string
			err := json.Unmarshal(attr.Value, &subIds)
			if err != nil {
				return fmt.Errorf("can't unmarshal subTokenIds: %s", err.Error())
			}
			for _, s := range subIds {
				v, err := strconv.ParseUint(s, 10, 32)
				if err != nil {
					return fmt.Errorf("can't parse sub token id '%s': %s", s, err.Error())
				}
				e.SubTokenIDs = append(e.SubTokenIDs, uint32(v))
			}
		}
	}

	ea.SendNFTs = append(ea.SendNFTs, e)
	return nil
}
