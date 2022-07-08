package api

import (
	"encoding/hex"
)

// Response of broadcast_tx_sync
type TxSyncResponse struct {
	// transaction hash
	Hash string
	// error info. Code = 0 mean no error
	Code      int
	Log       string
	Codespace string
}

// Send transaction data in sync mode. NOTE: marked by tendermint as deprecated
func (api *API) BroadcastTxSync(data []byte) (*TxSyncResponse, error) {
	type directSyncResponse struct {
		Result struct {
			Code      int    `json:"code"`
			Hash      string `json:"hash"`
			Log       string `json:"log"`
			Codespace string `json:"codespace"`
		} `json:"result"`
	}
	// request
	res, err := api.rpc.R().Get("/broadcast_tx_sync?tx=0x" + hex.EncodeToString(data))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	// json decode
	respValue := directSyncResponse{}
	err = universalJSONDecode(res.Body(), &respValue, nil, func() (bool, bool) {
		return respValue.Result.Hash > "", false
	})
	if err != nil {
		return nil, err
	}
	// process result
	return &TxSyncResponse{
		Code:      respValue.Result.Code,
		Hash:      respValue.Result.Hash,
		Log:       respValue.Result.Log,
		Codespace: respValue.Result.Codespace,
	}, nil
}

// Send transaction data in commit mode. NOTE: marked by tendermint as deprecated
func (api *API) BroadcastTxCommit(data []byte) (*TxSyncResponse, error) {
	type directSyncResponse struct {
		Result struct {
			CheckTx struct {
				Code      int    `json:"code"`
				Log       string `json:"log"`
				Codespace string `json:"codespace"`
			} `json:"check_tx"`
			DeliverTx struct {
				Code      int    `json:"code"`
				Log       string `json:"log"`
				Codespace string `json:"codespace"`
			} `json:"deliver_tx"`
			Hash string `json:"hash"`
		} `json:"result"`
	}
	// request
	res, err := api.rpc.R().Get("/broadcast_tx_commit?tx=0x" + hex.EncodeToString(data))
	if err = processConnectionError(res, err); err != nil {
		return nil, err
	}
	// json decode
	respValue := directSyncResponse{}
	err = universalJSONDecode(res.Body(), &respValue, nil, func() (bool, bool) {
		return respValue.Result.Hash > "", false
	})
	if err != nil {
		return nil, err
	}
	if respValue.Result.CheckTx.Code != 0 {
		return &TxSyncResponse{
			Code:      respValue.Result.CheckTx.Code,
			Hash:      respValue.Result.Hash,
			Log:       respValue.Result.CheckTx.Log,
			Codespace: respValue.Result.CheckTx.Codespace,
		}, nil
	}
	// process result
	return &TxSyncResponse{
		Code:      respValue.Result.DeliverTx.Code,
		Hash:      respValue.Result.Hash,
		Log:       respValue.Result.DeliverTx.Log,
		Codespace: respValue.Result.DeliverTx.Codespace,
	}, nil
}
