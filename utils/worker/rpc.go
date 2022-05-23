package worker

import (
	"encoding/json"
	"fmt"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (w *Worker) fetchBlockInfo(height int64) *ctypes.ResultBlock {
	// Request until get block
	for {
		// Request block
		result, err := w.rpcClient.Block(w.ctx, &height)
		if err == nil {
			return result
		}
	}
}

func (w *Worker) fetchAllTxs(height int64, total int, ch chan *[]Tx) {
	query := fmt.Sprintf("tx.height=%d", height)
	page, perPage := 1, 100

	var txs []Tx
	for len(txs) < total {

		// Request transactions
		result, err := w.rpcClient.TxSearch(w.ctx, query, true, &page, &perPage, "")
		w.panicError(err)

		for _, tx := range result.Txs {
			var parsedTx Tx
			var txLog []interface{}

			// Recover messages from raw transaction bytes
			recoveredTx, err := w.cdc.TxConfig.TxDecoder()(tx.Tx)
			w.panicError(err)

			// Parse transaction results logs
			err = json.Unmarshal([]byte(tx.TxResult.Log), &txLog)
			if err != nil {
				parsedTx.Log = []interface{}{FailedTxLog{Log: tx.TxResult.Log}}
			} else {
				parsedTx.Log = txLog
			}

			parsedTx.Tx = w.parseTxFromStd(recoveredTx)
			parsedTx.Data = tx.TxResult.Data
			parsedTx.Hash = tx.Hash.String()
			parsedTx.Code = tx.TxResult.Code
			parsedTx.GasUsed = tx.TxResult.GasUsed
			parsedTx.GasWanted = tx.TxResult.GasWanted

			txs = append(txs, parsedTx)
		}

		if len(result.Txs) > 0 {
			page++
		}
	}

	// Send result to the channel
	ch <- &txs
}

func (w *Worker) fetchBlockSize(height int64, ch chan int) {

	// Request blockchain info
	result, err := w.rpcClient.BlockchainInfo(w.ctx, height, height)
	w.panicError(err)

	// Send result to the channel
	ch <- result.BlockMetas[0].BlockSize
}

func (w *Worker) fetchBlockResults(height int64, ch chan *ctypes.ResultBlockResults) {

	// Request until get block results
	for {
		// Request block results
		result, err := w.rpcClient.BlockResults(w.ctx, &height)
		if err == nil { // len(result.EndBlockEvents) != 0
			ch <- result
			break
		}
	}
}
