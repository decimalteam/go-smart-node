package api

import (
	"context"

	tmtypes "github.com/cometbft/cometbft/abci/types"
	cosmostx "github.com/cosmos/cosmos-sdk/types/tx"
)

type TxResult struct {
	Height    int64
	Code      uint32
	GasWanted int64
	GasUsed   int64
	Codespace string
	RawLog    string
	Log       []TxLog
	Events    []tmtypes.Event
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
	return api.grpcTransaction(hash)
}

func (api *API) grpcTransaction(hash string) (*TxResult, error) {
	client := cosmostx.NewServiceClient(api.grpcClient)
	resp, err := client.GetTx(context.Background(), &cosmostx.GetTxRequest{Hash: hash})
	if err != nil {
		return nil, err
	}
	resLogs := make([]TxLog, len(resp.TxResponse.Logs))
	for i, log := range resp.TxResponse.Logs {
		resEvents := make([]TxEvent, len(log.Events))
		for j, event := range log.Events {
			resAttrs := make([]TxAttribute, len(event.Attributes))
			for k, attr := range event.Attributes {
				resAttrs[k] = TxAttribute{
					Key:   attr.Key,
					Value: attr.Value,
				}
			}
			resEvents[j] = TxEvent{
				Type:       event.Type,
				Attributes: resAttrs,
			}
		}
		resLogs[i] = TxLog{resEvents}
	}
	return &TxResult{
		Height:    resp.TxResponse.Height,
		Code:      resp.TxResponse.Code,
		GasWanted: resp.TxResponse.GasWanted,
		GasUsed:   resp.TxResponse.GasUsed,
		Codespace: resp.TxResponse.Codespace,
		RawLog:    resp.TxResponse.RawLog,
		Log:       resLogs,
		Events:    resp.TxResponse.Events,
	}, nil
}

// TODO
func (api *API) TransactionsByBlock(height uint64) ([]string, error) {
	return nil, nil
}
