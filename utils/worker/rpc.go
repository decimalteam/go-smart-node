package worker

import (
	"encoding/json"
	"fmt"
	"time"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

func (w *Worker) fetchBlock(height int64) *ctypes.ResultBlock {
	start := time.Now()

	// Request until get block
	for first := true; true; first = false {
		// Request block
		result, err := w.rpcClient.Block(w.ctx, &height)
		if err == nil {
			if !first {
				w.logger.Info(
					fmt.Sprintf("Fetched block (after %s)", helpers.DurationToString(time.Since(start))),
					"block", height,
				)
			} else {
				w.logger.Info(
					fmt.Sprintf("Fetched block (%s)", helpers.DurationToString(time.Since(start))),
					"block", height,
				)
			}
			return result
		}
	}

	return nil
}

func (w *Worker) fetchBlockSize(height int64, ch chan int) {

	// Request blockchain info
	result, err := w.rpcClient.BlockchainInfo(w.ctx, height, height)
	w.panicError(err)

	// Send result to the channel
	ch <- result.BlockMetas[0].BlockSize
}

func (w *Worker) fetchBlockTxs(height int64, total int, ea *EventAccumulator, ch chan []Tx) {
	query := fmt.Sprintf("tx.height=%d", height)
	page, perPage := 1, 100

	var results []Tx
	for len(results) < total {

		// Request transactions
		result, err := w.rpcClient.TxSearch(w.ctx, query, true, &page, &perPage, "")
		w.panicError(err)

		for _, tx := range result.Txs {
			var result Tx
			var txLog []interface{}

			// Recover messages from raw transaction bytes
			recoveredTx, err := w.cdc.TxConfig.TxDecoder()(tx.Tx)
			w.panicError(err)

			// Parse transaction results logs
			err = json.Unmarshal([]byte(tx.TxResult.Log), &txLog)
			if err != nil {
				result.Log = []interface{}{FailedTxLog{Log: tx.TxResult.Log}}
			} else {
				result.Log = txLog
			}

			result.Info = w.parseTxInfo(recoveredTx)
			result.Data = tx.TxResult.Data
			result.Hash = tx.Hash.String()
			result.Code = tx.TxResult.Code
			result.GasUsed = tx.TxResult.GasUsed
			result.GasWanted = tx.TxResult.GasWanted

			results = append(results, result)

			// process events for successful transactions
			if tx.TxResult.Code == 0 {
				for _, event := range tx.TxResult.Events {
					err := ea.AddEvent(event, tx.Hash.String(), tx.Height)
					if err != nil {
						w.panicError(err)
					}
				}
			}
		}

		if len(result.Txs) > 0 {
			page++
		}
	}

	// Send results to the channel
	ch <- results
}

func (w *Worker) fetchBlockTxResults(height int64, ch chan *ctypes.ResultBlockResults) {

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
