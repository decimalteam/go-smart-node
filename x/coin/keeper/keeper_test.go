package keeper_test

import (
	"context"
	"fmt"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/testcoin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func bootstrapKeeperTest(t *testing.T, numAddrs int, accCoins sdk.Coins) (*app.DSC, sdk.Context, []sdk.AccAddress, []sdk.ValAddress) {
	_, dsc, ctx := testkeeper.GetTestAppWithCoinKeeper(t)

	addrDels, addrVals := testkeeper.GenerateAddresses(dsc, ctx, numAddrs, accCoins)
	require.NotNil(t, addrDels)
	require.NotNil(t, addrVals)

	return dsc, ctx, addrDels, addrVals
}

var (
	baseDenom  = "del"
	baseAmount = helpers.EtherToWei(sdk.NewInt(1000000000000))
)

func TestKeeper_Coin(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapKeeperTest(t, 1, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})

	symbol := "testcoin"
	newCoin := types.Coin{
		Title:       "test keeper coin functions coin",
		Symbol:      symbol,
		CRR:         50,
		Reserve:     helpers.EtherToWei(sdk.NewInt(5000)),
		Volume:      helpers.EtherToWei(sdk.NewInt(10000)),
		LimitVolume: helpers.EtherToWei(sdk.NewInt(1000000000)),
		Creator:     addrs[0].String(),
		Identity:    "",
	}

	// check set coin
	dsc.CoinKeeper.SetCoin(ctx, newCoin)

	// check get exist coin
	getCoin, err := dsc.CoinKeeper.GetCoin(ctx, symbol)
	require.NoError(t, err)
	require.True(t, getCoin.Equal(newCoin))
	// check get not exist coin
	_, err = dsc.CoinKeeper.GetCoin(ctx, "not exist symbol")
	require.Error(t, err)
	// check get coins
	coins := dsc.CoinKeeper.GetCoins(ctx)
	require.Equal(t, 2, len(coins))

	// edit coin
	dsc.CoinKeeper.EditCoin(ctx, getCoin, helpers.EtherToWei(sdk.NewInt(1000000001)), helpers.EtherToWei(sdk.NewInt(10002)))
}

func TestKeeper_Check(t *testing.T) {
	dsc, ctx, _, _ := bootstrapKeeperTest(t, 1, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})

	newCheck := testcoin.CreateNewCheck(ctx.ChainID(), fmt.Sprintf("10000%s", baseDenom), "9", "", 123)

	// verify new check is not redeemed
	ok := dsc.CoinKeeper.IsCheckRedeemed(ctx, &newCheck)
	require.False(t, ok)
	// set new check
	dsc.CoinKeeper.SetCheck(ctx, &newCheck)
	// get check
	newCheckHash := newCheck.HashFull()
	getCheck, err := dsc.CoinKeeper.GetCheck(ctx, newCheckHash[:])
	require.NoError(t, err)
	require.True(t, getCheck.Equal(newCheck))
	// get checks
	checks := dsc.CoinKeeper.GetChecks(ctx)
	require.Equal(t, 1, len(checks))
	//  verify new check is redeemed
	ok = dsc.CoinKeeper.IsCheckRedeemed(ctx, &newCheck)
	require.True(t, ok)
}

func TestKeeper_Params(t *testing.T) {
	dsc, ctx, _, _ := bootstrapKeeperTest(t, 1, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})

	defParams := types.DefaultParams()
	// set params
	dsc.CoinKeeper.SetParams(ctx, defParams)
	// get params
	getParams := dsc.CoinKeeper.GetParams(ctx)
	require.True(t, defParams.Equal(getParams))
}

func TestKeeper_Helpers(t *testing.T) {
	custCoinDenom := "custcoin"
	//custCoinAmount := helpers.EtherToWei(sdk.NewInt(10000))

	dsc, ctx, addrs, _ := bootstrapKeeperTest(t, 1, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})

	newCoin := types.Coin{
		Title:       "test keeper coin functions coin",
		Symbol:      custCoinDenom,
		CRR:         50,
		Reserve:     helpers.EtherToWei(sdk.NewInt(5000)),
		Volume:      helpers.EtherToWei(sdk.NewInt(10000)),
		LimitVolume: helpers.EtherToWei(sdk.NewInt(1000000000)),
		Creator:     addrs[0].String(),
		Identity:    "",
	}
	dsc.CoinKeeper.SetCoin(ctx, newCoin)

	// get base denom
	denom := dsc.CoinKeeper.GetBaseDenom(ctx)
	require.Equal(t, baseDenom, denom)

	// check the input denom equal to base denom
	ok := dsc.CoinKeeper.IsCoinBase(ctx, baseDenom)
	require.True(t, ok)

	// commission calculate ----
	// fee with base coin
	_, _, err := dsc.CoinKeeper.GetCommission(ctx, helpers.EtherToWei(sdk.NewInt(10)))
	require.NoError(t, err)

	// fee with custom coin
	ctxWithFee := ctx
	ctxWithFee = ctx.WithContext(context.WithValue(ctx.Context(), "fee", sdk.Coins{
		{
			Denom:  custCoinDenom,
			Amount: helpers.EtherToWei(sdk.NewInt(100)),
		},
	}))
	require.NotNil(t, ctxWithFee.Context())

	_, _, err = dsc.CoinKeeper.GetCommission(ctxWithFee, helpers.EtherToWei(sdk.NewInt(10)))
	require.NoError(t, err)

	// fee custom coin not exist
	ctxWithNotExistCoinFee := ctx
	ctxWithNotExistCoinFee = ctx.WithContext(context.WithValue(ctx.Context(), "fee", sdk.Coins{
		{
			Denom:  "notexistcoin",
			Amount: helpers.EtherToWei(sdk.NewInt(100)),
		},
	}))
	require.NotNil(t, ctxWithNotExistCoinFee.Context())

	_, _, err = dsc.CoinKeeper.GetCommission(ctxWithNotExistCoinFee, helpers.EtherToWei(sdk.NewInt(10)))
	require.Error(t, err)

	// fee custom coin reserve less than need fee base coin amount
	ctxWithLessReserveFee := ctx
	ctxWithLessReserveFee = ctx.WithContext(context.WithValue(ctx.Context(), "fee", sdk.Coins{
		{
			Denom:  custCoinDenom,
			Amount: helpers.EtherToWei(sdk.NewInt(1000)),
		},
	}))
	require.NotNil(t, ctxWithLessReserveFee.Context())

	_, _, err = dsc.CoinKeeper.GetCommission(ctxWithLessReserveFee, helpers.EtherToWei(sdk.NewInt(1000000000)))
	require.Error(t, err)
}

func TestKeeper_CoinCache(t *testing.T) {
	dsc, ctx, _, _ := bootstrapKeeperTest(t, 1, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})

	// set coin cache
	dsc.CoinKeeper.SetCachedCoin(ctx, baseDenom)
	// get coin cache
	ok := dsc.CoinKeeper.GetCoinCache(baseDenom)
	require.True(t, ok)
	// clear coin cache
	dsc.CoinKeeper.ClearCoinCache(ctx)
	ok = dsc.CoinKeeper.GetCoinCache(baseDenom)
	require.False(t, ok)
}
