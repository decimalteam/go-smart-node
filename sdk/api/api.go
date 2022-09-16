package api

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	coinTypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	tmservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
)

// this is default limit for queries with pagination
const queryLimit = 10

// API struct accumulates all queries to blockchain node
// and makes broadcast of prepared transaction
type API struct {
	grpcClient *grpc.ClientConn

	// network parameters from genesis
	chainID  string
	baseCoin string
}

type ConnectionOptions struct {
	EndpointHost string // hostname or IP without any protocol description like "http://"
	GRPCPort     int    // gRPC port, default 9090
	Timeout      uint   // timeout in seconds
}

func NewAPI(opts ConnectionOptions) (*API, error) {
	var err error

	// init global cosmos sdk prefixes
	initConfig()

	api := &API{}
	if opts.GRPCPort == 0 {
		opts.GRPCPort = 9090
	}
	if opts.Timeout == 0 {
		opts.Timeout = 10
	}
	// gRPC client

	api.grpcClient, err = grpc.Dial(
		fmt.Sprintf("%s:%d", opts.EndpointHost, opts.GRPCPort), // your gRPC server address.
		grpc.WithInsecure(), // The Cosmos SDK doesn't support any transport security mechanism.
	)
	if err != nil {
		return nil, err
	}

	return api, nil
}

func (api *API) Close() error {
	return api.grpcClient.Close()
}

// ChainID() returns blockchain network chain id
func (api *API) ChainID() string {
	return api.chainID
}

// BaseCoin() returns base coin symbol from genesis. Need for correct transaction building
func (api *API) BaseCoin() string {
	return api.baseCoin
}

// GetParameters() get blockchain parameters
func (api *API) GetParameters() error {
	return api.grpcGetParameters()
}

func (api *API) grpcGetParameters() error {
	// chain id
	{
		client := tmservice.NewServiceClient(api.grpcClient)
		resp, err := client.GetLatestBlock(context.Background(), &tmservice.GetLatestBlockRequest{})
		if err != nil {
			return err
		}
		api.chainID = resp.Block.Header.ChainID
	}
	// base coin
	{
		client := coinTypes.NewQueryClient(api.grpcClient)
		resp, err := client.Params(context.Background(), &coinTypes.QueryParamsRequest{})
		if err != nil {
			return err
		}
		api.baseCoin = resp.Params.BaseSymbol
	}
	return nil
}

// Init global cosmos sdk config
// Do not seal config or rework to use sealed config
func initConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(cmdcfg.Bech32PrefixAccAddr, cmdcfg.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(cmdcfg.Bech32PrefixValAddr, cmdcfg.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(cmdcfg.Bech32PrefixConsAddr, cmdcfg.Bech32PrefixConsPub)
}
