package worker

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"

	abci "github.com/cometbft/cometbft/abci/types"
)

// var swapPool = authtypes.NewModuleAddress(swaptypes.SwapPool)

type EventActivateChain struct {
	Sender string `json:"sender"`
	ID     uint32 `json:"id"`
	Name   string `json:"name"`
	TxHash string `json:"txHash"`
}

type EventDeactivateChain struct {
	Sender string `json:"sender"`
	ID     uint32 `json:"id"`
	TxHash string `json:"txHash"`
}

type EventSwapInitialize struct {
	Sender            string      `json:"sender"`
	DestChain         uint32      `json:"destChain"`
	FromChain         uint32      `json:"fromChain"`
	Recipient         string      `json:"recipient"`
	Amount            sdkmath.Int `json:"amount"`
	TokenDenom        string      `json:"tokenDenom"`
	TransactionNumber string      `json:"transactionNumber"`
	TxHash            string      `json:"txHash"`
}

type EventSwapRedeem struct {
	Sender            string      `json:"sender"`
	From              string      `json:"from"`
	Recipient         string      `json:"recipient"`
	Amount            sdkmath.Int `json:"amount"`
	TokenDenom        string      `json:"tokenDenom"`
	TransactionNumber string      `json:"transactionNumber"`
	DestChain         uint32      `json:"destChain"`
	FromChain         uint32      `json:"fromChain"`
	HashReedem        string      `json:"hashReedem"`
	V                 string      `json:"v"`
	R                 string      `json:"r"`
	S                 string      `json:"s"`
	TxHash            string      `json:"txHash"`
}

func processEventActivateChain(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
		string sender
		uint32 id = 1 [ (gogoproto.customname) = "ID" ];
		string name = 2;
	*/

	var e EventActivateChain
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			e.Sender = string(attr.Value)
		case "id":
			strId := string(attr.Value)
			id, err := strconv.ParseUint(strId, 10, 32)
			if err != nil {
				return fmt.Errorf("can't parse chain id '%s': %s", strId, err.Error())
			}
			e.ID = uint32(id)
		case "name":
			e.Name = string(attr.Value)
		}
	}

	e.TxHash = txHash
	ea.ActivateChain = append(ea.ActivateChain, e)
	return nil
}

func processEventDeactivateChain(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
		string sender;
		uint32 id = 1 [ (gogoproto.customname) = "ID" ];
	*/

	var e EventDeactivateChain
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			e.Sender = string(attr.Value)
		case "id":
			strId := string(attr.Value)
			id, err := strconv.ParseUint(strId, 10, 32)
			if err != nil {
				return fmt.Errorf("can't parse chain id '%s': %s", strId, err.Error())
			}
			e.ID = uint32(id)
		}
	}

	e.TxHash = txHash
	ea.DeactivateChain = append(ea.DeactivateChain, e)
	return nil
}

func processEventSwapInitialize(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	   string sender = 1;
	   string recipient = 5;
	   string amount = 6;
	   string token_symbol = 8;
	   string transaction_number = 7;
	   uint32 from_chain = 3;
	   uint32 dest_chain = 4;
	*/

	var (
		e EventSwapInitialize
	)
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			e.Sender = string(attr.Value)
		case "recipient":
			e.Recipient = string(attr.Value)
		case "amount":
			amountStr := string(attr.Value)
			amount, ok := sdkmath.NewIntFromString(amountStr)
			if !ok {
				return fmt.Errorf("failed to parse amount: %s", amountStr)
			}
			e.Amount = amount
		case "token_symbol":
			e.Sender = string(attr.Value)
		case "transaction_number":
			e.TransactionNumber = string(attr.Value)
		case "dest_chain":
			strId := string(attr.Value)
			id, err := strconv.ParseUint(strId, 10, 32)
			if err != nil {
				return fmt.Errorf("can't parse chain id '%s': %s", strId, err.Error())
			}
			e.DestChain = uint32(id)
		case "from_chain":
			strId := string(attr.Value)
			id, err := strconv.ParseUint(strId, 10, 32)
			if err != nil {
				return fmt.Errorf("can't parse chain id '%s': %s", strId, err.Error())
			}
			e.FromChain = uint32(id)
		}
	}

	e.TxHash = txHash
	//ea.addBalanceChange(e.Sender, e.TokenDenom, e.Amount.Neg())
	//ea.addBalanceChange(swaptypes.SwapPool, e.TokenDenom, e.Amount)
	ea.SwapInitialize = append(ea.SwapInitialize, e)
	return nil
}

func processEventSwapRedeem(ea *EventAccumulator, event abci.Event, txHash string) error {
	/*
	   string sender = 1;
	   string from = 2;
	   string recipient =3;
	   string amount = 4;
	   string token_symbol = 5;
	   string transaction_number = 6;
	   uint32 from_chain = 7;
	   uint32 dest_chain = 8;
	   string hashRedeem = 9;
	   string v = 10;
	   string r = 11;
	   string s = 12;
	*/
	var (
		e EventSwapRedeem
	)
	for _, attr := range event.Attributes {
		switch string(attr.Key) {
		case "sender":
			e.Sender = string(attr.Value)
		case "from":
			e.From = string(attr.Value)
		case "recipient":
			e.Recipient = string(attr.Value)
		case "amount":
			amountStr := string(attr.Value)
			amount, ok := sdkmath.NewIntFromString(amountStr)
			if !ok {
				return fmt.Errorf("failed to parse amount: %s", amountStr)
			}
			e.Amount = amount
		case "token_symbol":
			e.Sender = string(attr.Value)
		case "transaction_number":
			e.TransactionNumber = string(attr.Value)
		case "dest_chain":
			strId := string(attr.Value)
			id, err := strconv.ParseUint(strId, 10, 32)
			if err != nil {
				return fmt.Errorf("can't parse chain id '%s': %s", strId, err.Error())
			}
			e.DestChain = uint32(id)
		case "src_chain":
			strId := string(attr.Value)
			id, err := strconv.ParseUint(strId, 10, 32)
			if err != nil {
				return fmt.Errorf("can't parse chain id '%s': %s", strId, err.Error())
			}
			e.FromChain = uint32(id)
		case "v":
			e.V = string(attr.Value)
		case "r":
			e.R = string(attr.Value)
		case "s":
			e.S = string(attr.Value)
		case "hash_redeem":
			e.HashReedem = string(attr.Value)
		}
	}

	e.TxHash = txHash
	//ea.addBalanceChange(e.Recipient, e.TokenDenom, e.Amount)
	//ea.addBalanceChange(swaptypes.SwapPool, e.TokenDenom, e.Amount.Neg())
	ea.SwapRedeem = append(ea.SwapRedeem, e)
	return nil
}
