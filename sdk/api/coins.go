package api

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// see x/coin/types
type Coin struct {
	Title       string  `json:"title"`
	Symbol      string  `json:"symbol"`
	CRR         uint64  `json:"constant_reserve_ratio"`
	Reserve     sdk.Int `json:"reserve"`
	Volume      sdk.Int `json:"volume"`
	LimitVolume sdk.Int `json:"limit_volume"`
	Creator     string  `json:"creator"`
	Identity    string  `json:"identity"`
}

func (api *API) Coins() ([]Coin, error) {
	type directCoinsResult struct {
		Height string `json:"height"`
		Result []struct {
			Title       string `json:"title"`
			Symbol      string `json:"symbol"`
			CRR         string `json:"constant_reserve_ratio"`
			Reserve     string `json:"reserve"`
			Volume      string `json:"volume"`
			LimitVolume string `json:"limit_volume"`
			Creator     string `json:"creator"`
			Identity    string `json:"identity"`
		} `json:"result"`
	}
	// request
	res, err := api.rest.R().Get("/coin/all_coins")
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	// json decode
	respValue, respErr := directCoinsResult{}, Error{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.Height > "", respErr.StatusCode != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr)
	}
	// process result
	var result []Coin
	for _, rawCoin := range respValue.Result {
		var tmp uint64
		var ok bool
		coin := Coin{}
		coin.Title = rawCoin.Title
		coin.Symbol = rawCoin.Symbol
		tmp, err = strconv.ParseUint(rawCoin.CRR, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: field CRR='%s'", err, rawCoin.CRR)
		}
		coin.CRR = tmp
		coin.Reserve, ok = sdk.NewIntFromString(rawCoin.Reserve)
		if !ok {
			return nil, fmt.Errorf("not ok field Reserve='%s'", rawCoin.Reserve)
		}
		coin.Volume, ok = sdk.NewIntFromString(rawCoin.Volume)
		if !ok {
			return nil, fmt.Errorf("not ok field Volume='%s'", rawCoin.Volume)
		}
		coin.LimitVolume, ok = sdk.NewIntFromString(rawCoin.LimitVolume)
		if !ok {
			return nil, fmt.Errorf("not ok field Volume='%s'", rawCoin.LimitVolume)
		}
		coin.Creator = rawCoin.Creator
		coin.Identity = rawCoin.Identity
		result = append(result, coin)
	}
	return result, nil
}
