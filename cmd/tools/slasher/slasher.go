package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Event struct {
	Type       string       `json:"type"`
	Attributes []*Attribute `json:"attributes"`
}

type Result struct {
	BeginBlockEvents []*Event `json:"begin_block_events"`
	EndBlockEvents   []*Event `json:"end_block_events"`
}

type Response struct {
	JsonRPC string  `json:"jsonrpc"`
	ID      int     `json:"id"`
	Result  *Result `json:"result"`
}

func decodeBase64(encoded string) string {
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encoded))
	buffer := make([]byte, 64*1024)
	read, err := decoder.Read(buffer)
	if err != nil {
		panic(err)
	}
	return string(buffer[:read])
}

func findAttribute(key string, attributes []*Attribute) string {
	for _, a := range attributes {
		if a.Key == key {
			return a.Value
		}
	}
	return ""
}

func main() {
	// url := "http://kihot.crypton.studio/rpc/block_results?height=9288729"
	url := "http://46.101.127.241/rpc/block_results?height=2211"

	client := http.Client{
		Timeout: time.Second * 5, // Timeout after 5 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "unfair-slashes-counter")

	// fmt.Println("Retrieving data from the node's RPC...")

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// fmt.Println("Parsing the response...")

	resp := Response{}
	jsonErr := json.Unmarshal(body, &resp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	// fmt.Println("Decoding base64 attribute's keys/values and splitting by event type...")

	eventsByTypes := make(map[string][]*Event)

	for _, event := range resp.Result.BeginBlockEvents {
		for _, attribute := range event.Attributes {
			attribute.Key = decodeBase64(attribute.Key)
			attribute.Value = decodeBase64(attribute.Value)
		}
		eventsByTypes[event.Type] = append(eventsByTypes[event.Type], event)
	}

	// fmt.Println("Printing results...")

	// result, _ := json.MarshalIndent(resp.Result, "", "  ")
	// fmt.Println(string(result))

	fmt.Println("Validator;Delegator;Type;Slash Amount;Slash Denom;NFT Denom;NFT ID;NFT SubTokenID;NFT SubToken Reserve")
	for t, events := range eventsByTypes {
		switch t {
		case "liveness":
			for _, event := range events {
				validator := findAttribute("validator", event.Attributes)
				delegator := findAttribute("delegator", event.Attributes)
				amount := findAttribute("slash_amount", event.Attributes)
				if len(amount) == 0 {
					continue
				}
				coin, err := sdk.ParseCoinNormalized(amount)
				if err != nil {
					panic(err)
				}
				fmt.Println(validator, delegator, "coin", coin.Amount, coin.Denom)
			}
		case "liveness_nft":
			for _, event := range events {
				validator := findAttribute("validator", event.Attributes)
				delegator := findAttribute("delegator", event.Attributes)
				amount := findAttribute("slash_amount", event.Attributes)
				if len(amount) == 0 {
					continue
				}
				coin, err := sdk.ParseCoinNormalized(amount)
				if err != nil {
					panic(err)
				}
				denom := findAttribute("denom", event.Attributes)
				id := findAttribute("id", event.Attributes)
				subTokenID := findAttribute("sub_token_id", event.Attributes)
				subTokenReserve := findAttribute("sub_token_id_reserve", event.Attributes)
				fmt.Println(validator, delegator, "nft", coin.Amount, coin.Denom, denom, id, subTokenID, subTokenReserve)
			}
		case "slash":
			// for _, event := range events {
			// 	address := findAttribute("address", event.Attributes)
			// 	amount := findAttribute("slash_amount", event.Attributes)
			// 	fmt.Println("SLASH:", address, amount)
			// }
		case "update_coin":
			for _, event := range events {
				symbol := findAttribute("symbol", event.Attributes)
				volume := findAttribute("volume", event.Attributes)
				reserve := findAttribute("reserve", event.Attributes)
				v, _ := sdk.NewIntFromString(volume)
				r, _ := sdk.NewIntFromString(reserve)
				fmt.Printf("% 12s %s\n", symbol, sdk.NewDecFromInt(v).QuoInt(r).String())
			}
		}
	}
}
