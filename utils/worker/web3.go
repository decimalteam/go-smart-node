package worker

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/status-im/keycard-go/hexutils"

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
	wg := &sync.WaitGroup{}

	// Request transaction receipts by hashes in parallel
	results := make(web3types.Receipts, len(block.Transactions()))
	batchElems := make([]ethrpc.BatchElem, len(block.Transactions()))

	for i, tx := range block.Transactions() {
		batchElems[i] = ethrpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{tx.Hash()},
			Result: &results[i],
		}
	}

	quantity := 15 // requests count in batch

	end := quantity
	if len(batchElems) < end {
		end = len(batchElems)
	}
	for i := 0; i < len(batchElems); {
		elems := batchElems[i:end]
		wg.Add(1)

		go func(requests []ethrpc.BatchElem) {
			defer wg.Done()
			err := w.ethRpcClient.BatchCall(requests)
			w.panicError(err)
		}(elems)

		i = end
		if end+quantity > len(batchElems) {
			end = len(batchElems)
			continue
		}
		end += quantity
	}

	wg.Wait()

	for i := range batchElems {
		if batchElems[i].Error != nil {
			w.panicError(batchElems[i].Error)
		}
		if results[i] == nil {
			w.panicError(fmt.Errorf("got null result for tx with hash %s", hexutils.BytesToHex(batchElems[i].Args[0].([]byte))))
		}
	}

	ch <- results
}
