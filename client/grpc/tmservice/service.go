package tmservice

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
)

var (
	_ ServiceServer = queryServer{}
)

type (
	abciQueryFn = func(abci.RequestQuery) abci.ResponseQuery

	queryServer struct {
		clientCtx         client.Context
		interfaceRegistry codectypes.InterfaceRegistry
		queryFn           abciQueryFn
	}
)

// NewQueryServer creates a new tendermint query server.
func NewQueryServer(
	clientCtx client.Context,
	interfaceRegistry codectypes.InterfaceRegistry,
	queryFn abciQueryFn,
) ServiceServer {
	return queryServer{
		clientCtx:         clientCtx,
		interfaceRegistry: interfaceRegistry,
		queryFn:           queryFn,
	}
}

func (s queryServer) GetBlockchainInfo(ctx context.Context, req *GetBlockchainInfoRequest) (*GetBlockchainInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	info, err := getBlockchainInfo(ctx, s.clientCtx, req.GetMinHeight(), req.GetMaxHeight())
	if err != nil {
		return nil, err
	}

	protoMetas := make([]*types.BlockMeta, len(info.BlockMetas))
	for i, v := range info.BlockMetas {
		protoMetas[i] = v.ToProto()
	}

	return &GetBlockchainInfoResponse{BlockMetas: protoMetas}, err
}

func (s queryServer) GetTxSearch(ctx context.Context, req *GetTxSearchRequest) (*GetTxSearchResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	res, totalCount, err := getTxSearch(ctx, s.clientCtx, req.Query, req.Prove, req.Page, req.PerPage, req.OrderBy)
	if err != nil {
		return nil, err
	}

	return &GetTxSearchResponse{
		Txs:        res,
		TotalCount: int64(totalCount),
	}, nil
}

func (s queryServer) GetBlockResults(ctx context.Context, req *GetBlockResultsRequest) (*GetBlockResultsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	res, err := getBlockResults(ctx, s.clientCtx, req.Height)
	if err != nil {
		return nil, err
	}

	return &GetBlockResultsResponse{
		BlockResults: res,
	}, nil
}

// RegisterTendermintService registers the tendermint queries on the gRPC router.
func RegisterTendermintService(
	clientCtx client.Context,
	server gogogrpc.Server,
	iRegistry codectypes.InterfaceRegistry,
	queryFn abciQueryFn,
) {
	RegisterServiceServer(server, NewQueryServer(clientCtx, iRegistry, queryFn))
}

// RegisterGRPCGatewayRoutes mounts the tendermint service's GRPC-gateway routes on the
// given Mux.
func RegisterGRPCGatewayRoutes(clientConn gogogrpc.ClientConn, mux *runtime.ServeMux) {
	_ = RegisterServiceHandlerClient(context.Background(), mux, NewServiceClient(clientConn))
}
