package worker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/status-im/keycard-go/hexutils"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

func (w *Worker) fetchBlock(height int64) *ctypes.ResultBlock {
	// Request until get block
	for first, start, deadline := true, time.Now(), time.Now().Add(RequestTimeout); true; first = false {
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
		// Stop trying when the deadline is reached
		if time.Now().After(deadline) {
			w.logger.Error("Failed to fetch block", "block", height, "error", err)
			return nil
		}
		// Sleep some time before next try
		time.Sleep(RequestRetryDelay)
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

func (w *Worker) fetchBlockResults(height int64, block ctypes.ResultBlock, ea *EventAccumulator, ch chan []Tx, brch chan *ctypes.ResultBlockResults) {
	var err error

	// Request block results from the node
	// NOTE: Try to retrieve results in the loop since it looks like there is some delay before results are ready to by retrieved
	var blockResults *ctypes.ResultBlockResults
	for c := 1; true; c++ {
		if c > 5 {
			w.logger.Debug(fmt.Sprintf("%d attempt to fetch block height: %d, time %s", c, height, time.Now().String()))
		}
		// Request block results
		blockResults, err = w.rpcClient.BlockResults(w.ctx, &height)
		if err == nil {
			break
		}
		// Sleep some time before next try
		time.Sleep(RequestRetryDelay)
	}

	// Prepare block results by overall processing
	var results []Tx
	for i, tx := range block.Block.Txs {
		var result Tx
		var txLog []interface{}
		txr := blockResults.TxsResults[i]

		recoveredTx, err := w.cdc.TxConfig.TxDecoder()(tx)
		w.panicError(err)

		// Parse transaction results logs
		err = json.Unmarshal([]byte(txr.Log), &txLog)
		if err != nil {
			result.Log = []interface{}{FailedTxLog{Log: txr.Log}}
		} else {
			result.Log = txLog
		}

		result.Info = w.parseTxInfo(recoveredTx)
		result.Data = txr.Data
		result.Hash = hexutils.BytesToHex(tx.Hash())
		result.Code = txr.Code
		result.Codespace = txr.Codespace
		result.GasUsed = txr.GasUsed
		result.GasWanted = txr.GasWanted

		results = append(results, result)

		// process events for transactions
		for _, event := range txr.Events {
			err := ea.AddEvent(event, hexutils.BytesToHex(tx.Hash()))
			if err != nil {
				fmt.Printf("error in event %v\n", event.Type)
				w.panicError(err)
			}
		}
	}

	// Send results to the channel
	ch <- results
	brch <- blockResults
}
