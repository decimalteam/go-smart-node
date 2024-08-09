package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"context"
	"encoding/hex"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"

	feemarkettypes "github.com/decimalteam/ethermint/x/feemarket/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	feeconfig "bitbucket.org/decimalteam/go-smart-node/x/fee/config"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
)

var _ types.QueryServer = Keeper{}
var _ feemarkettypes.QueryServer = Keeper{}

/////////////
// Fee Keeper
/////////////

func (k Keeper) CoinPrices(c context.Context, req *types.QueryCoinPricesRequest) (*types.QueryCoinPricesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	prices, err := k.GetPrices(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCoinPricesResponse{Prices: prices}, nil
}

func (k Keeper) CoinPrice(c context.Context, req *types.QueryCoinPriceRequest) (*types.QueryCoinPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	price, err := k.GetPrice(ctx, req.Denom, req.Quote)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCoinPriceResponse{Price: &price}, nil
}

func (k Keeper) ModuleParams(c context.Context, req *types.QueryModuleParamsRequest) (*types.QueryModuleParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	params := k.GetModuleParams(ctx)

	return &types.QueryModuleParamsResponse{Params: params}, nil
}

/////////////
// Fee Market Keeper
/////////////

// BaseFee implements the Query/BaseFee gRPC method
func (k Keeper) BaseFee(c context.Context, _ *feemarkettypes.QueryBaseFeeRequest) (*feemarkettypes.QueryBaseFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	res := &feemarkettypes.QueryBaseFeeResponse{}
	baseFee := k.GetBaseFee(ctx)

	if baseFee != nil {
		aux := sdkmath.NewIntFromBigInt(baseFee)
		res.BaseFee = &aux
	}

	return res, nil
}

// Params implements the Query/Params gRPC method
func (k Keeper) Params(c context.Context, _ *feemarkettypes.QueryParamsRequest) (*feemarkettypes.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &feemarkettypes.QueryParamsResponse{
		Params: params,
	}, nil
}

// BlockGas implements the Query/BlockGas gRPC method
func (k Keeper) BlockGas(c context.Context, _ *feemarkettypes.QueryBlockGasRequest) (*feemarkettypes.QueryBlockGasResponse, error) {
	// TODO: rework when ethermint starts use BlockGas from EVM
	return &feemarkettypes.QueryBlockGasResponse{
		Gas: 0,
	}, nil
}

func (k Keeper) CalculateCommission(c context.Context, req *types.QueryCalculateCommissionRequest) (*types.QueryCalculateCommissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	baseDenom := helpers.GetBaseDenom(ctx.ChainID())

	bz, err := hex.DecodeString(req.TxBytes)
	if err != nil {
		return nil, err
	}

	msg, err := txDecode(k.cdc, bz)
	if err != nil {
		return nil, err
	}
	params := k.GetModuleParams(ctx)
	delPrice, err := k.GetPrice(ctx, baseDenom, feeconfig.DefaultQuote)
	if err != nil {
		return nil, err
	}

	commission, err := k.calcFunc(k.cdc, []sdk.Msg{msg}, int64(len(bz)), delPrice.Price, params)
	if req.Denom == baseDenom {
		return &types.QueryCalculateCommissionResponse{
			Commission: commission,
		}, nil
	}

	coinInfo, err := k.coinKeeper.GetCoin(ctx, req.Denom)
	if err != nil {
		return nil, err
	}
	commission = formulas.CalculateSaleAmount(coinInfo.Volume, coinInfo.Reserve, uint(coinInfo.CRR), commission)
	return &types.QueryCalculateCommissionResponse{
		Commission: commission,
	}, nil
}

// DefaultTxDecoder returns a default protobuf TxDecoder using the provided Marshaler.
func txDecode(cdc codec.BinaryCodec, txBytes []byte) (sdk.Msg, error) {
	var raw sdktx.TxRaw

	err := cdc.Unmarshal(txBytes, &raw)
	if err != nil {
		return nil, err
	}

	var body sdktx.TxBody

	err = cdc.Unmarshal(raw.BodyBytes, &body)
	if err != nil {
		return nil, err
	}

	if len(body.Messages) != 1 {
		return nil, errors.MustBeOneMessage
	}

	var msg sdk.Msg
	err = cdc.UnpackAny(body.Messages[0], &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
