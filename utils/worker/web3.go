package worker

import (
	"fmt"
	"math/big"
	"time"

	web3types "github.com/ethereum/go-ethereum/core/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

const TxReceiptsBatchSize = 16
const RequestRetryDelay = 32 * time.Millisecond

func (w *Worker) fetchBlockWeb3(height int64, ch chan *web3types.Block) {

	// Request block by number
	result, err := w.web3Client.BlockByNumber(w.ctx, big.NewInt(height))
	w.panicError(err)

	// Send result to the channel
	ch <- result
}

func (w *Worker) fetchBlockTxReceiptsWeb3(block *web3types.Block, ch chan web3types.Receipts) {
	txCount := len(block.Transactions())
	results := make(web3types.Receipts, txCount)

	// NOTE: Try to retrieve results in the loop since it looks like there is some delay before results are ready to by retrieved
	for c := 0; true; c++ {
		if c > 0 {
			w.logger.Debug(fmt.Sprintf("%d attempt to fetch block height: %d, time %s", c, block.NumberU64(), time.Now().String()))
		}
		// Prepare batch requests to retrieve the receipt for each transaction in the block
		requests := make([]ethrpc.BatchElem, txCount)
		for i, tx := range block.Transactions() {
			requests[i] = ethrpc.BatchElem{
				Method: "eth_getTransactionReceipt",
				Args:   []interface{}{tx.Hash()},
				Result: &results[i],
			}
		}
		// Request transaction receipts with a batch
		err := w.ethRpcClient.BatchCall(requests)
		if err == nil {
			// Ensure all transaction receipts are retrieved
			for i := range requests {
				if requests[i].Error != nil {
					err = requests[i].Error
					w.logger.Error(fmt.Sprintf("Error: %v", err))
					// w.panicError(err)
				}
				if results[i] == nil {
					txHash := requests[i].Args[0].([]byte)
					err = fmt.Errorf("got null result for tx with hash %X", txHash)
					w.logger.Error(fmt.Sprintf("Error: %v", err))
					// w.panicError(err)
				}
			}
			if err == nil {
				break
			}
		}
		// Sleep some time before next try
		time.Sleep(RequestRetryDelay)
	}

	// Send results to the channel
	ch <- results
}
