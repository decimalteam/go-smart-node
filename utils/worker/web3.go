package worker

import (
	"fmt"
	"math/big"
	"sync"
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
	requests := make([]ethrpc.BatchElem, txCount)

	// Prepare batch requests to retrieve the receipt for each transaction in the block
	for i, tx := range block.Transactions() {
		requests[i] = ethrpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{tx.Hash()},
			Result: &results[i],
		}
	}

	// Request transaction receipts with batches in parallel
	wg := &sync.WaitGroup{}
	for i, s := 0, TxReceiptsBatchSize; i < txCount; i += s {
		end := i + s
		if end > txCount {
			end = txCount
		}
		wg.Add(1)
		go func(requests []ethrpc.BatchElem) {
			defer wg.Done()
			// NOTE: Try to retrieve receipts in the loop since it looks like there is some delay before receipts are ready to by retrieved
			for c := 0; true; c++ {
				if c > 0 {
					w.logger.Debug(fmt.Sprintf("%d attempt to fetch tx receipts height: %d, time %s", c, block.NumberU64(), time.Now().String()))
				}
				// Request transaction receipts with the batch
				err := w.ethRpcClient.BatchCall(requests)
				if err == nil {
					break
				}
				// Sleep some time before next try
				time.Sleep(RequestRetryDelay)
			}
		}(requests[i:end])
	}
	wg.Wait()

	// Ensure all transaction receipts are retrieved
	for i := range requests {
		if requests[i].Error != nil {
			w.panicError(requests[i].Error)
		}
		if results[i] == nil {
			txHash := requests[i].Args[0].([]byte)
			err := fmt.Errorf("got null result for tx with hash %X", txHash)
			w.panicError(err)
		}
	}

	// Send results to the channel
	ch <- results
}
