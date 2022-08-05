package api

import (
	"fmt"
	"log"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	resty "github.com/go-resty/resty/v2"
	"google.golang.org/grpc"

	config "bitbucket.org/decimalteam/go-smart-node/cmd/config"
)

// API struct accumulates all queries to blockchain node
// and makes broadcast of prepared transaction
type API struct {
	rpc        *resty.Client // this is interface
	rest       *resty.Client
	grpcClient *grpc.ClientConn

	useGRPC bool

	// network parameters from genesis
	chainID  string
	maxGas   uint64
	baseCoin string
}

type ConnectionOptions struct {
	EndpointHost   string // hostname or IP without any protocol description like "http://"
	TendermintPort int    // tendermint RPC port, default 26657
	RestPort       int    // REST server port, default 1317
	GRPCPort       int    // gRPC port, default 9090
	Timeout        uint   // timeout in seconds
	Debug          bool   //turn on debugging via stdlib log
	UseGRPC        bool
}

//resty logger implementation
type log2log struct{}

type Logger interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

func (l log2log) Errorf(format string, v ...interface{}) {
	log.Printf("L2LERR:"+format, v...)
}
func (l log2log) Warnf(format string, v ...interface{}) {
	log.Printf("L2LWRN:"+format, v...)
}
func (l log2log) Debugf(format string, v ...interface{}) {
	log.Printf("L2LDBG:"+format, v...)
}

func NewAPI(opts ConnectionOptions) (*API, error) {
	var err error
	api := &API{}
	if opts.TendermintPort == 0 {
		opts.TendermintPort = 26657
	}
	if opts.RestPort == 0 {
		opts.RestPort = 1317
	}
	if opts.GRPCPort == 0 {
		opts.GRPCPort = 9090
	}
	if opts.Timeout == 0 {
		opts.Timeout = 10
	}
	// rpc client
	api.rpc = resty.New().
		SetTimeout(time.Duration(opts.Timeout) * time.Second).
		SetBaseURL(fmt.Sprintf("http://%s:%d", opts.EndpointHost, opts.TendermintPort))
	// rest client
	api.rest = resty.New().
		SetTimeout(time.Duration(opts.Timeout) * time.Second).
		SetBaseURL(fmt.Sprintf("http://%s:%d", opts.EndpointHost, opts.RestPort))
	// gRPC client
	if opts.UseGRPC {
		api.useGRPC = true
		api.grpcClient, err = grpc.Dial(
			fmt.Sprintf("%s:%d", opts.EndpointHost, opts.GRPCPort), // your gRPC server address.
			grpc.WithInsecure(), // The Cosmos SDK doesn't support any transport security mechanism.
		)
		if err != nil {
			return nil, err
		}
	}
	//
	if opts.Debug {
		api.rest = api.rest.SetDebug(true).SetLogger(log2log{})
		api.rpc = api.rpc.SetDebug(true).SetLogger(log2log{})
	}
	// init global cosmos sdk prefixes
	initConfig()
	return api, nil
}

func (api *API) Close() error {
	return api.grpcClient.Close()
}

// ChainID() returns blockchain network chain id
func (api *API) ChainID() string {
	return api.chainID
}

// MaxGas() returns max gas from genesis. Need for correct transaction building
func (api *API) MaxGas() uint64 {
	return api.maxGas
}

// BaseCoin() returns base coin symbol from genesis. Need for correct transaction building
func (api *API) BaseCoin() string {
	return api.baseCoin
}

// GetParameters() get blockchain parameters
func (api *API) GetParameters() error {
	return api.restGetParameters()
}

func (api *API) restGetParameters() error {
	type respDirectGenesis struct {
		Result struct {
			Genesis struct {
				ChainID         string `json:"chain_id"`
				ConsensusParams struct {
					Block struct {
						MaxGas string `json:"max_gas"`
					} `json:"block"`
				} `json:"consensus_params"`
				AppState struct {
					Coin struct {
						Params struct {
							BaseSymbol string `json:"base_symbol"`
						} `json:"params"`
					} `json:"coin"`
				} `json:"app_state"`
			} `json:"genesis"`
		} `json:"result"`
	}
	// request
	res, err := api.rpc.R().Get("/genesis")
	if err = processConnectionError(res, err); err != nil {
		return err
	}
	// json decode
	respValue := respDirectGenesis{}
	err = universalJSONDecode(res.Body(), &respValue, nil, func() (bool, bool) {
		return respValue.Result.Genesis.ChainID > "", false
	})
	if err != nil {
		return err
	}
	// process results
	maxGas, err := strconv.ParseUint(respValue.Result.Genesis.ConsensusParams.Block.MaxGas, 10, 64)
	if err != nil {
		return err
	}
	api.chainID = respValue.Result.Genesis.ChainID
	api.baseCoin = respValue.Result.Genesis.AppState.Coin.Params.BaseSymbol
	api.maxGas = maxGas
	return nil
}

// Init global cosmos sdk config
// Do not seal config or rework to use sealed config
func initConfig() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(config.Bech32PrefixAccAddr, config.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(config.Bech32PrefixValAddr, config.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(config.Bech32PrefixConsAddr, config.Bech32PrefixConsPub)
}
