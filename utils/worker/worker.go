package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	web3 "github.com/ethereum/go-ethereum/ethclient"

	"github.com/tendermint/tendermint/libs/log"
	rpc "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/simapp/params"
)

type Worker struct {
	ctx        context.Context
	httpClient *http.Client
	cdc        params.EncodingConfig
	logger     log.Logger
	config     *Config
	hostname   string
	rpcClient  *rpc.HTTP
	web3Client *web3.Client
	query      chan *ParseTask
}

type Config struct {
	IndexerEndpoint string
	RpcEndpoint     string
	Web3Endpoint    string
	WorkersCount    int
}

func NewWorker(cdc params.EncodingConfig, logger log.Logger, config *Config) (*Worker, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{}
	rpcClient, err := rpc.NewWithClient(config.RpcEndpoint, config.RpcEndpoint, httpClient)
	if err != nil {
		return nil, err
	}
	web3Client, err := web3.Dial(config.Web3Endpoint)
	if err != nil {
		return nil, err
	}
	worker := &Worker{
		ctx:        context.Background(),
		httpClient: httpClient,
		cdc:        cdc,
		logger:     logger,
		config:     config,
		hostname:   hostname,
		rpcClient:  rpcClient,
		web3Client: web3Client,
		query:      make(chan *ParseTask, 1000),
	}
	return worker, nil
}

func (w *Worker) Start() {
	wg := &sync.WaitGroup{}
	wg.Add(w.config.WorkersCount)
	for i := 0; i < w.config.WorkersCount; i++ {
		go w.executeFromQuery(wg)
	}
	wg.Wait()
}

func (w *Worker) executeFromQuery(wg *sync.WaitGroup) {
	defer wg.Done()
	w.getWork()
	for {
		task := <-w.query
		w.getBlockResultAndSend(task.height, task.txNum)
		w.getWork()
	}
}

func (w *Worker) getBlockResultAndSend(height int64, txNum int) {

	// Fetch requested block from Tendermint RPC
	block := w.fetchBlockInfo(height)

	// Fetch everything needed from Tendermint RPC
	txsChan := make(chan *[]Tx)
	resultsChan := make(chan *ctypes.ResultBlockResults)
	sizeChan := make(chan int)
	var parseTxNum int
	if txNum == -1 {
		parseTxNum = len(block.Block.Data.Txs)
	} else {
		parseTxNum = txNum
	}
	go w.fetchAllTxs(height, parseTxNum, txsChan)
	go w.fetchBlockResults(height, resultsChan)
	go w.fetchBlockSize(height, sizeChan)
	txs := <-txsChan
	results := <-resultsChan
	size := <-sizeChan

	w.logger.Info(fmt.Sprintf(
		"Found %d txs in %d block. Events on begin/end block: %d/%d.",
		len(*txs), height, len(results.BeginBlockEvents), len(results.EndBlockEvents),
	))

	// Retrieve emission and rewards
	var emission string
	var rewards []ProposerReward
	var commissionRewards []CommissionReward
	for _, event := range results.EndBlockEvents {
		switch event.Type {
		case "emission":
			// Parse emission
			emission = string(event.Attributes[0].Value)

		case "proposer_reward":
			// Parse proposer rewards
			var reward ProposerReward
			for _, attr := range event.Attributes {
				switch string(attr.Key) {
				case "amount", "accum_rewards":
					reward.Reward = string(attr.Value)
				case "validator", "accum_rewards_validator":
					reward.Address = string(attr.Value)
				case "delegator":
					reward.Delegator = string(attr.Value)
				}
			}
			rewards = append(rewards, reward)

		case "commission_reward":
			// Parser commission reward
			var reward CommissionReward
			for _, attr := range event.Attributes {
				switch string(attr.Key) {
				case "amount":
					reward.Amount = string(attr.Value)
				case "validator":
					reward.Validator = string(attr.Value)
				case "reward_address":
					reward.RewardAddress = string(attr.Value)
				}
			}
			commissionRewards = append(commissionRewards, reward)
		}
	}

	// Create and fill Block object and then marshal to JSON
	var b Block
	b.Evidence = block.Block.Evidence
	b.Header = block.Block.Header
	b.LastCommit = block.Block.LastCommit
	b.Data = BlockData{Txs: *txs}
	b.Emission = emission
	b.Rewards = rewards
	b.CommissionRewards = commissionRewards
	b.EndBlockEvents = w.parseEvents(results.EndBlockEvents)
	b.BeginBlockEvents = w.parseEvents(results.BeginBlockEvents)
	b.Size = size
	data, err := json.Marshal(b)
	w.panicError(err)

	// Send
	w.sendBlock(data)
}

func (w *Worker) panicError(err error) {
	if err != nil {
		w.logger.Error("Error: ", err)
		panic(err)
	}
}
