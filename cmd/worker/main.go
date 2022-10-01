package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/ethermint/encoding"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/worker"
)

var (
	cdc    params.EncodingConfig
	logger log.Logger
)

func init() {

	// Setup config
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	cmdcfg.SetBip44CoinType(config)
	cmdcfg.RegisterBaseDenom()
	config.Seal()

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

	grpcPort, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		logger.Error("Error convert GRPC port from str to int")
		panic(err)
	}
	// Prepare worker configuration
	config := &worker.Config{
		IndexerEndpoint: os.Getenv("INDEXER_URL"),
		RpcEndpoint:     os.Getenv("RPC_URL"),
		Web3Endpoint:    os.Getenv("WEB3_URL"),
		GRPCHost:        os.Getenv("GRPC_HOST"),
		GRPCPort:        grpcPort,
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
