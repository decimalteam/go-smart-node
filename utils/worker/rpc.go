package worker

import (
	"bitbucket.org/decimalteam/go-smart-node/client/grpc/tmservice"
	"encoding/json"
	"fmt"
	cosmostmservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/proto/tendermint/types"
	"time"

	"github.com/status-im/keycard-go/hexutils"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

func (w *Worker) fetchBlock(height int64) (*cosmostmservice.Block, *types.BlockID) {
	start := time.Now()

	// Request until get block
	for first := true; true; first = false {
		// Request block
		result, err := w.cTmClient.GetBlockByHeight(w.ctx, &cosmostmservice.GetBlockByHeightRequest{Height: height})
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
			return result.SdkBlock, result.BlockId
		}
		// Sleep some time before next try
		time.Sleep(RequestRetryDelay)
	}

	return nil, nil
}

func (w *Worker) fetchBlockSize(height int64, ch chan int) {

	// Request blockchain info
	result, err := w.tmClient.GetBlockchainInfo(w.ctx, &tmservice.GetBlockchainInfoRequest{
		MinHeight: height,
		MaxHeight: height,
	},
	)
	w.panicError(err)

	// Send result to the channel
	ch <- int(result.BlockMetas[0].BlockSize)
}

func (w *Worker) fetchBlockResults(height int64, block *cosmostmservice.Block, ea *EventAccumulator, ch chan []Tx, brch chan *ctypes.ResultBlockResults) {
	// Request block results from the node
	// NOTE: Try to retrieve results in the loop since it looks like there is some delay before results are ready to by retrieved
	var blockResults *ctypes.ResultBlockResults
	for c := 0; true; c++ {
		if c > 0 {
			w.logger.Debug(fmt.Sprintf("%d attempt to fetch block height: %d, time %s", c, height, time.Now().String()))
		}
		// Request block results
		resp, err := w.tmClient.GetBlockResults(w.ctx, &tmservice.GetBlockResultsRequest{
			Height: height,
		})
		if err == nil {
			blockResults = resp.BlockResults.ToCoreTypes()
			break
		}
		// Sleep some time before next try
		time.Sleep(RequestRetryDelay)
	}

	// Prepare block results by overall processing
	var results []Tx
	for i, tx := range block.Data.Txs {
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
		result.Hash = hexutils.BytesToHex(tmhash.Sum(tx))
		result.Code = txr.Code
		result.GasUsed = txr.GasUsed
		result.GasWanted = txr.GasWanted

		results = append(results, result)

		// process events for successful transactions
		if txr.Code == 0 {
			for _, event := range txr.Events {
				err := ea.AddEvent(event, hexutils.BytesToHex(tmhash.Sum(tx)))
				if err != nil {
					fmt.Printf("error in event %v\n", event.Type)
					w.panicError(err)
				}
			}
		}
	}

	// Send results to the channel
	ch <- results
	brch <- blockResults
}
