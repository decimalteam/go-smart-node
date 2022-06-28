package worker

import (
	"math/big"
	"sync"

	web3types "github.com/ethereum/go-ethereum/core/types"
)

func (w *Worker) fetchBlockWeb3(height int64, ch chan *web3types.Block) {

	// Request block by number
	result, err := w.web3Client.BlockByNumber(w.ctx, big.NewInt(height))
	w.panicError(err)

	// Send result to the channel
	ch <- result
}

func (w *Worker) fetchBlockTxReceiptsWeb3(block *web3types.Block, ch chan web3types.Receipts) {

	// Request transaction receipts by hashes in parallel
	results := make(web3types.Receipts, len(block.Transactions()))
	wg := &sync.WaitGroup{}
	wg.Add(len(block.Transactions()))
	for i := range block.Transactions() {
		go func(i int) {
			defer wg.Done()
			result, err := w.web3Client.TransactionReceipt(w.ctx, block.Transactions()[i].Hash())
			w.panicError(err)
			results[i] = result
		}(i)
	}
	wg.Wait()

	// Send results to the channel
	ch <- results
}
