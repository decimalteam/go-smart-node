package testcoin

import (
	"bitbucket.org/decimalteam/go-smart-node/x/coin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

// Helper is a structure which wraps the staking handler
// and provides methods useful in tests
type Helper struct {
	t *testing.T
	h sdk.Handler
	k keeper.Keeper

	Ctx sdk.Context
	//Commission cointypes.CommissionRates
	// Coin Denomination
	//Denom string
}

// NewHelper creates staking Handler wrapper for tests
func NewHelper(t *testing.T, ctx sdk.Context, k keeper.Keeper) *Helper {
	return &Helper{
		t:   t,
		h:   coin.NewHandler(k),
		k:   k,
		Ctx: ctx,
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
func (sh *Helper) CreateCoin(sender sdk.AccAddress, title, symbol string, crr uint64, initVolume, initReserve, limitVolume sdk.Int, identity string, ok bool) types.Coin {
	msg := types.NewMsgCreateCoin(sender, title, symbol, crr, initVolume, initReserve, limitVolume, identity)
	sh.Handle(sh.Ctx, msg, ok)
	return types.Coin{
		Title:       title,
		Symbol:      symbol,
		CRR:         crr,
		Reserve:     initReserve,
		Volume:      initVolume,
		LimitVolume: limitVolume,
		Creator:     sender.String(),
		Identity:    identity,
	}
}

// CreateCoinWithContext create msg and handle create coin with custom context
func (sh *Helper) CreateCoinWithContext(ctx sdk.Context, sender sdk.AccAddress, title, symbol string, crr uint64, initVolume, initReserve, limitVolume sdk.Int, identity string, ok bool) types.Coin {
	msg := types.NewMsgCreateCoin(sender, title, symbol, crr, initVolume, initReserve, limitVolume, identity)
	sh.Handle(ctx, msg, ok)
	return types.Coin{
		Title:       title,
		Symbol:      symbol,
		CRR:         crr,
		Reserve:     initReserve,
		Volume:      initVolume,
		LimitVolume: limitVolume,
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
func (sh *Helper) UpdateCoin(sender sdk.AccAddress, symbol string, limitVolume sdk.Int, identity string, ok bool) {
	msg := types.NewMsgUpdateCoin(sender, symbol, limitVolume, identity)
	sh.Handle(sh.Ctx, msg, ok)
}

// SendCoin create msg and handler send coin to other address
func (sh *Helper) SendCoin(sender, receiver sdk.AccAddress, coin sdk.Coin, ok bool) {
	msg := types.NewMsgSendCoin(sender, coin, receiver)
	sh.Handle(sh.Ctx, msg, ok)
}

// MultiSendCoin create msg and handler send coin to other addresses
func (sh *Helper) MultiSendCoin(sender sdk.AccAddress, sends []types.Send, ok bool) {
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
func (sh *Helper) SellAllCoin(sender sdk.AccAddress, coinToSell, minCoinToBuy sdk.Coin, ok bool) {
	msg := types.NewMsgSellAllCoin(sender, coinToSell, minCoinToBuy)
	sh.Handle(sh.Ctx, msg, ok)
}

func (sh *Helper) GetCoin(symbol string, ok bool) {
	_, err := sh.k.GetCoin(sh.Ctx, symbol)
	if ok {
		require.NoError(sh.t, err)
	} else {
		require.Error(sh.t, err)
	}
}
