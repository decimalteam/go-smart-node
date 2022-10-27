package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/ethermint/encoding"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cnfcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/worker"
)

var (
	cdc    params.EncodingConfig
	logger log.Logger
)

func init() {

	// Setup config
	cfg := sdk.GetConfig()
	cnfcfg.SetBech32Prefixes(cfg)
	cnfcfg.SetBip44CoinType(cfg)
	cnfcfg.RegisterBaseDenom()
	cfg.Seal()

	// Register interfaces and create codec
	cdc = encoding.MakeConfig(app.ModuleBasics)

	// Logger
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

func main() {
	runtime.SetCPUProfileRate(100000)
	f, err := os.Create("worker.pprof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = pprof.StartCPUProfile(f)
	if err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	// Prepare worker configuration
	config := &worker.Config{
		IndexerEndpoint: "",
		RpcEndpoint:     "http://localhost:26657",
		Web3Endpoint:    "http://localhost:8545",
		WorkersCount:    1,
	}

	// Create worker
	w, err := worker.NewWorker(cdc, logger, config)
	if err != nil {
		panic(err)
	}

	for block := int64(916); block < 917; block++ {
		fmt.Printf("process block=%d\n", block)
		w.GetBlockResult(block, -1)
	}
}
