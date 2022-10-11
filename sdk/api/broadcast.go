package api

import (
	"context"

	cosmostx "github.com/cosmos/cosmos-sdk/types/tx"
	tmtypes "github.com/tendermint/tendermint/abci/types"
)

// Response of broadcast_tx_sync
type TxSyncResponse struct {
	// transaction hash
	Hash string
	// error info. Code = 0 mean no error
	Code      uint32
	Log       string
	Codespace string
	Events    []tmtypes.Event
}

type TxSimulateResponse struct {
	Log    string
	Events []tmtypes.Event
}

func (api *API) BroadcastTxSync(data []byte) (*TxSyncResponse, error) {
	return api.grpcBroadcastTx(data, false)
}

func (api *API) BroadcastTxCommit(data []byte) (*TxSyncResponse, error) {
	return api.grpcBroadcastTx(data, true)
}

func (api *API) grpcBroadcastTx(data []byte, commitMode bool) (*TxSyncResponse, error) {
	mode := cosmostx.BroadcastMode_BROADCAST_MODE_SYNC
	if commitMode {
		mode = cosmostx.BroadcastMode_BROADCAST_MODE_BLOCK
	}
	client := cosmostx.NewServiceClient(api.grpcClient)
	resp, err := client.BroadcastTx(context.Background(), &cosmostx.BroadcastTxRequest{
		TxBytes: data,
		Mode:    mode,
	})
	if err != nil {
		return nil, err
	}
	return &TxSyncResponse{
		Hash:      resp.TxResponse.TxHash,
		Code:      resp.TxResponse.Code,
		Log:       resp.TxResponse.RawLog,
		Codespace: resp.TxResponse.Codespace,
		Events:    resp.TxResponse.Events,
	}, nil
}

func (api *API) SimulateTx(data []byte) (*TxSimulateResponse, error) {
	client := cosmostx.NewServiceClient(api.grpcClient)
	resp, err := client.Simulate(context.Background(), &cosmostx.SimulateRequest{
		Tx:      nil,
		TxBytes: data,
	})
	if err != nil {
		return nil, err
	}
	return &TxSimulateResponse{
		Log:    resp.Result.Log,
		Events: resp.Result.Events,
	}, nil
}
