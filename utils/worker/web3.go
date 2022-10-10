package worker

import (
	"fmt"
	"math/big"
	"time"

	web3common "github.com/ethereum/go-ethereum/common"
	web3types "github.com/ethereum/go-ethereum/core/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

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
	if txCount == 0 {
		ch <- results
		return
	}

	// NOTE: Try to retrieve tx receipts in the loop since it looks like there is some delay before receipts are ready to by retrieved
	for c := 1; true; c++ {
		if c > 5 {
			w.logger.Debug(fmt.Sprintf("%d attempt to fetch transaction receipts with height: %d, time %s", c, block.NumberU64(), time.Now().String()))
		}
		// Prepare batch requests to retrieve the receipt for each transaction in the block
		for i, tx := range block.Transactions() {
			results[i] = &web3types.Receipt{}
			requests[i] = ethrpc.BatchElem{
				Method: "eth_getTransactionReceipt",
				Args:   []interface{}{tx.Hash()},
				Result: results[i],
			}
		}
		// Request transaction receipts with a batch
		err := w.ethRpcClient.BatchCall(requests[:])
		if err == nil {
			// Ensure all transaction receipts are retrieved
			for i := range requests {
				if requests[i].Error != nil {
					err = requests[i].Error
					w.logger.Error(fmt.Sprintf("Error: %v", err))
				}
				if results[i].BlockNumber == nil || results[i].BlockNumber.Sign() == 0 {
					txHash := requests[i].Args[0].(web3common.Hash)
					err = fmt.Errorf("got null result for tx with hash %v", txHash)
					w.logger.Error(fmt.Sprintf("Error: %v", err))
				}
			}
		}
		if err == nil {
			break
		}
		// Sleep some time before next try
		time.Sleep(RequestRetryDelay)
	}

	// Send results to the channel
	ch <- results
}
