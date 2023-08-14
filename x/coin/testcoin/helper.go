package testcoin

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"bitbucket.org/decimalteam/go-smart-node/x/coin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// Helper is a structure which wraps the staking handler
// and provides methods useful in tests
type Helper struct {
	t *testing.T
	h sdk.Handler
	k keeper.Keeper

	msgServer   types.MsgServer
	queryServer types.QueryServer
	Ctx         sdk.Context
}

// NewHelper creates staking Handler wrapper for tests
func NewHelper(t *testing.T, ctx sdk.Context, k keeper.Keeper) *Helper {
	return &Helper{
		t: t,
		h: coin.NewHandler(k),
		k: k,

		msgServer:   &k,
		queryServer: &k,
		Ctx:         ctx,
	}
}

// Handle calls staking handler on a given message
func (sh *Helper) Handle(ctx sdk.Context, msg sdk.Msg, ok bool) *sdk.Result {
	res, err := sh.h(ctx, msg)
	if ok {
		require.NoError(sh.t, err)
		require.NotNil(sh.t, res)
	} else {
		require.Error(sh.t, err)
		require.Nil(sh.t, res)
	}
	return res
}

// CreateCoin create msg and handle create coin
func (sh *Helper) CreateCoin(sender sdk.AccAddress, title, denom string, crr uint64, initVolume, initReserve, limitVolume sdkmath.Int, identity string, ok bool) types.Coin {
	msg := types.NewMsgCreateCoin(sender, denom, title, crr, initVolume, initReserve, limitVolume, sdkmath.ZeroInt(), identity)
	sh.Handle(sh.Ctx, msg, ok)
	return types.Coin{
		Title:        title,
		Denom:        denom,
		CRR:          uint32(crr),
		Reserve:      initReserve,
		Volume:       initVolume,
		LimitVolume:  limitVolume,
		MinVolume:    sdkmath.ZeroInt(),
		Creator:      sender.String(),
		Identity:     identity,
		Drc20Address: "0x1a7e5e7e6c9f33b7d34fd76eeffbcee6a006f700",
	}
}

// CreateCoinWithContext create msg and handle create coin with custom context
func (sh *Helper) CreateCoinWithContext(ctx sdk.Context, sender sdk.AccAddress, title, denom string, crr uint64, initVolume, initReserve, limitVolume sdkmath.Int, identity string, ok bool) types.Coin {
	msg := types.NewMsgCreateCoin(sender, denom, title, crr, initVolume, initReserve, limitVolume, sdkmath.ZeroInt(), identity)
	sh.Handle(ctx, msg, ok)
	return types.Coin{
		Title:       title,
		Denom:       denom,
		CRR:         uint32(crr),
		Reserve:     initReserve,
		Volume:      initVolume,
		LimitVolume: limitVolume,
		MinVolume:   sdkmath.ZeroInt(),
		Creator:     sender.String(),
		Identity:    identity,
	}
}

// CheckRedeem create msg and handle redeem check
func (sh *Helper) CheckRedeem(sender sdk.AccAddress, check, proof string, ok bool) {
	msg := types.NewMsgRedeemCheck(sender, check, proof)
	sh.Handle(sh.Ctx, msg, ok)
}

// UpdateCoin create msg and handle update coin
func (sh *Helper) UpdateCoin(sender sdk.AccAddress, denom string, limitVolume sdkmath.Int, identity string, ok bool) {
	msg := types.NewMsgUpdateCoin(sender, denom, limitVolume, sdkmath.ZeroInt(), identity)
	sh.Handle(sh.Ctx, msg, ok)
}

// SendCoin create msg and handler send coin to other address
func (sh *Helper) SendCoin(sender, recipient sdk.AccAddress, coin sdk.Coin, ok bool) {
	msg := types.NewMsgSendCoin(sender, recipient, coin)
	sh.Handle(sh.Ctx, msg, ok)
}

// MultiSendCoin create msg and handler send coin to other addresses
func (sh *Helper) MultiSendCoin(sender sdk.AccAddress, sends []types.MultiSendEntry, ok bool) {
	msg := types.NewMsgMultiSendCoin(sender, sends)
	sh.Handle(sh.Ctx, msg, ok)
}

// BuyCoin create msg and handler buy coins request
func (sh *Helper) BuyCoin(sender sdk.AccAddress, coinToBuy, maxCoinToSell sdk.Coin, ok bool) {
	msg := types.NewMsgBuyCoin(sender, coinToBuy, maxCoinToSell)
	sh.Handle(sh.Ctx, msg, ok)
}

// SellCoin create msg and handler sell coins request
func (sh *Helper) SellCoin(sender sdk.AccAddress, coinToSell, maxCoinToBuy sdk.Coin, ok bool) {
	msg := types.NewMsgSellCoin(sender, coinToSell, maxCoinToBuy)
	sh.Handle(sh.Ctx, msg, ok)
}

// SellAllCoin create msg and handler sell all coins request
func (sh *Helper) SellAllCoin(sender sdk.AccAddress, coinDenomToSell string, minCoinToBuy sdk.Coin, ok bool) {
	msg := types.NewMsgSellAllCoin(sender, coinDenomToSell, minCoinToBuy)
	sh.Handle(sh.Ctx, msg, ok)
}

func (sh *Helper) GetCoin(denom string, ok bool) {
	_, err := sh.k.GetCoin(sh.Ctx, denom)
	if ok {
		require.NoError(sh.t, err)
	} else {
		require.Error(sh.t, err)
	}
}

func (sh *Helper) QueryCoins() []types.Coin {
	resp, err := sh.queryServer.Coins(sh.Ctx, types.NewQueryCoinsRequest(&query.PageRequest{Limit: 1000}))
	require.NoError(sh.t, err)

	return resp.Coins
}

func (sh *Helper) QueryCoin(denom string) types.Coin {
	resp, err := sh.queryServer.Coin(sh.Ctx, types.NewQueryCoinRequest(denom))
	require.NoError(sh.t, err)

	return resp.Coin
}

func (sh *Helper) QueryChecks() []types.Check {
	resp, err := sh.queryServer.Checks(sh.Ctx, types.NewQueryChecksRequest(&query.PageRequest{Limit: 1000}))
	require.NoError(sh.t, err)

	return resp.Checks
}

func (sh *Helper) QueryCheck(hash []byte) types.Check {
	resp, err := sh.queryServer.Check(sh.Ctx, types.NewQueryCheckRequest(hash))
	require.NoError(sh.t, err)

	return resp.Check
}

func (sh *Helper) QueryParams() types.Params {
	resp, err := sh.queryServer.Params(sh.Ctx, types.NewQueryParamsRequest())
	require.NoError(sh.t, err)

	return resp.Params
}
