package api

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	tmservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/decimalteam/ethermint/encoding"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
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

	// codec for fee calculation
	appCodec  codec.BinaryCodec
	delPrice  sdk.Dec
	feeParams feetypes.Params
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
		fmt.Sprintf("%s:%d", opts.EndpointHost, opts.GRPCPort),   // your gRPC server address.
		grpc.WithTransportCredentials(insecure.NewCredentials()), // The Cosmos SDK doesn't support any transport security mechanism.
	)
	if err != nil {
		return nil, err
	}

	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	api.appCodec = encodingConfig.Codec

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
		api.chainID = resp.SdkBlock.Header.ChainID
	}
	// base coin
	{
		client := cointypes.NewQueryClient(api.grpcClient)
		resp, err := client.Params(context.Background(), &cointypes.QueryParamsRequest{})
		if err != nil {
			return err
		}
		api.baseCoin = resp.Params.BaseDenom
	}
	// price and fee params
	{
		var err error
		// TODO: parametrize quote
		api.delPrice, api.feeParams, err = api.GetFeeParams(api.baseCoin, "usd")
		if err != nil {
			return err
		}
	}
	return nil
}

func (api *API) GetFeeCalculationOptions() *tx.FeeCalculationOptions {
	return &tx.FeeCalculationOptions{
		DelPrice:  api.delPrice,
		FeeParams: api.feeParams,
		AppCodec:  api.appCodec,
	}
}

// GetParameters() get blockchain parameters
func (api *API) GetLastHeight() int64 {
	client := tmservice.NewServiceClient(api.grpcClient)
	resp, err := client.GetLatestBlock(context.Background(), &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return 0
	}
	return resp.SdkBlock.Header.Height
}

func (api *API) GetSupply(denom string) (sdk.Int, error) {
	bankClient := bankTypes.NewQueryClient(api.grpcClient)
	req := &bankTypes.QuerySupplyOfRequest{
		Denom: denom,
	}
	res, err := bankClient.SupplyOf(context.Background(), req)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return res.Amount.Amount, nil
}

// Init global cosmos sdk config
// Do not seal config or rework to use sealed config
func initConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(cmdcfg.Bech32PrefixAccAddr, cmdcfg.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(cmdcfg.Bech32PrefixValAddr, cmdcfg.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(cmdcfg.Bech32PrefixConsAddr, cmdcfg.Bech32PrefixConsPub)
}
