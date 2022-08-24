package main

import (
	"bitbucket.org/decimalteam/go-smart-node/encoding"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cnfcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/worker"
)

var (
	cdc    encoding.EncodingConfig
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

	// Load .env file from current directory
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
		os.Exit(1)
	}

	// Prepare worker configuration
	config := &worker.Config{
		IndexerEndpoint: os.Getenv("INDEXER_URL"),
		RpcEndpoint:     os.Getenv("RPC_URL"),
		Web3Endpoint:    os.Getenv("WEB3_URL"),
		WorkersCount:    1,
	}
	if len(os.Getenv("WORKERS")) > 0 {
		workersCount, err := strconv.Atoi(os.Getenv("WORKERS"))
		if err == nil {
			config.WorkersCount = workersCount
		}
	}

	// Create worker
	w, err := worker.NewWorker(cdc, logger, config)
	if err != nil {
		panic(err)
	}

	// Start worker
	w.Start()
}
