package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type TxResult struct {
	Height    uint64
	Code      int
	GasWanted uint64
	GasUsed   uint64
	Codespace string
	RawLog    string
	Log       []TxLog
}

// TxLog contains API response fields.
type TxLog struct {
	Events []TxEvent `json:"events"`
}

// TxEvent contains API response fields.
type TxEvent struct {
	Type       string        `json:"type"`
	Attributes []TxAttribute `json:"attributes"`
}

// TxAttribute contains API response fields.
type TxAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (api *API) Transaction(hash string) (*TxResult, error) {
	type directTxResult struct {
		Result struct {
			Hash     string `json:"hash"`
			Height   string `json:"height"`
			TxResult struct {
				Code      int    `json:"code"`
				GasUsed   string `json:"gas_used"`
				GasWanted string `json:"gas_wanted"`
				Log       string `json:"log"`
				Codespace string `json:"codespace"`
			} `json:"tx_result"`
		} `json:"result"`
	}
	type directError struct {
		err RPCError `json:"error"`
	}
	// request
	res, err := api.rpc.R().Get("/tx?hash=0x" + hash)
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	// json decode
	respValue, respErr := directTxResult{}, directError{}
	err = universalJSONDecode(res.Body(), &respValue, &respErr, func() (bool, bool) {
		return respValue.Result.Hash > "", respErr.err.Code != 0
	})
	if err != nil {
		return nil, joinErrors(err, respErr.err)
	}
	// process results
	result := &TxResult{}
	// TODO: case of transaction error
	height, err := strconv.ParseUint(respValue.Result.Height, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: field Height='%s'", err, respValue.Result.Height)
	}
	gasWanted, err := strconv.ParseUint(respValue.Result.TxResult.GasWanted, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: field GasWanted='%s'", err, respValue.Result.TxResult.GasWanted)
	}
	gasUsed, err := strconv.ParseUint(respValue.Result.TxResult.GasUsed, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: field GasUsed='%s'", err, respValue.Result.TxResult.GasUsed)
	}
	//
	result.Height = height
	result.GasUsed = gasUsed
	result.GasWanted = gasWanted
	result.Code = respValue.Result.TxResult.Code
	result.Codespace = respValue.Result.TxResult.Codespace
	result.RawLog = respValue.Result.TxResult.Log

	err = json.Unmarshal([]byte(respValue.Result.TxResult.Log), &result.Log)
	if err != nil {
		return nil, fmt.Errorf("%w: field Log='%s'", err, respValue.Result.TxResult.Log)
	}

	return result, nil
}

func (api *API) TransactionsByBlock(height uint64) ([]string, error) {
	return nil, nil
}
