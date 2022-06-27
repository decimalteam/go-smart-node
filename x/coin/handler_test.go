package coin_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/testcoin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"github.com/cosmos/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethereumCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"golang.org/x/crypto/sha3"
)

func bootstrapHandlerGenesisTest(t *testing.T, numAddrs int, accCoins sdk.Coins) (*app.DSC, sdk.Context, []sdk.AccAddress, []sdk.ValAddress) {
	_, dsc, ctx := getBaseAppWithCustomKeeper()

	addrDels, addrVals := generateAddresses(dsc, ctx, numAddrs, accCoins)
	require.NotNil(t, addrDels)
	require.NotNil(t, addrVals)

	return dsc, ctx, addrDels, addrVals
}

var (
	baseDenom  = "del"
	baseAmount = helpers.EtherToWei(sdk.NewInt(1000000000000))

	// valid test coin params
	title              = "Its Test Coin"
	symbol             = "tstcoin"
	crr         uint64 = 50
	initVolume         = helpers.EtherToWei(sdk.NewInt(1000))
	initReserve        = helpers.EtherToWei(sdk.NewInt(10000))
	limitVolume        = helpers.EtherToWei(sdk.NewInt(10000000000000000))
)

func TestCreateCoinHandler(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapHandlerGenesisTest(t, 2, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})
	tscoin := testcoin.NewHelper(t, ctx, dsc.CoinKeeper)
	addr1 := addrs[0]
	addr2 := addrs[1]

	// create coin
	coin := tscoin.CreateCoin(addr1, title, symbol, crr, initVolume, initReserve, limitVolume, "", true)

	// check store with coin equals
	storeCoin, err := dsc.CoinKeeper.GetCoin(ctx, coin.Symbol)
	require.NoError(t, err)
	require.True(t, coin.Equal(storeCoin))

	// create coin with exist symbol
	_ = tscoin.CreateCoin(addr1, title, symbol, crr, initVolume, initReserve, limitVolume, "", false)

	// create coin with custom fee coin
	ctxWithFee := tscoin.Ctx
	ctxWithFee = ctxWithFee.WithContext(context.WithValue(ctxWithFee.Context(), "fee", sdk.Coins{
		{
			Denom:  symbol,
			Amount: helpers.EtherToWei(sdk.NewInt(100)),
		},
	}))

	require.NotNil(t, ctxWithFee.Context())
	_ = tscoin.CreateCoinWithContext(ctxWithFee, addr1, title, "customFeeTest", crr, initVolume, initReserve, limitVolume, "", true)

	// create coin without balance at user
	addrWithouBalance := app.CreateAddr()
	_ = tscoin.CreateCoin(addrWithouBalance, title, "withoutBaseBalance", crr, initVolume, initReserve, limitVolume, "", false)

	// create coin without balance at user with custom fee coin
	_ = tscoin.CreateCoinWithContext(ctxWithFee, addrWithouBalance, title, "customFeeTestWithoutBaseBalance", crr, initVolume, initReserve, limitVolume, "", false)

	// create coin without balance at user with custom fee coin
	_ = tscoin.CreateCoinWithContext(ctxWithFee, addr2, title, "customFeeTestWithoutCustomBalance", crr, initVolume, initReserve, limitVolume, "", false)

}

func TestUpdateCoinHandler(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapHandlerGenesisTest(t, 2, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})
	tscoin := testcoin.NewHelper(t, ctx, dsc.CoinKeeper)
	addr1 := addrs[0]
	addr2 := addrs[1]

	// create coin
	_ = tscoin.CreateCoin(addr1, title, symbol, crr, initVolume, initReserve, limitVolume, "", true)

	// update coin
	newLimitVolume := limitVolume.Add(helpers.EtherToWei(sdk.NewInt(1009000)))
	tscoin.UpdateCoin(addr1, symbol, newLimitVolume, "", true)

	storeCoin, err := dsc.CoinKeeper.GetCoin(ctx, symbol)
	require.NoError(t, err)
	require.True(t, storeCoin.LimitVolume.Equal(newLimitVolume))

	// update coin not from creator
	tscoin.UpdateCoin(addr2, symbol, newLimitVolume, "", false)

	// update not exist coin
	tscoin.UpdateCoin(addr1, "notExistCoin", newLimitVolume, "", false)

	// update with less limit volume
	lessLimitVolume := newLimitVolume.Sub(helpers.EtherToWei(sdk.NewInt(10000)))
	tscoin.UpdateCoin(addr1, symbol, lessLimitVolume, "", false)
}

func TestSendCoinHandler(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapHandlerGenesisTest(t, 2, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})

	tscoin := testcoin.NewHelper(t, ctx, dsc.CoinKeeper)
	addr1 := addrs[0]
	addr2 := addrs[1]
	emptyBalanceAddr := app.CreateAddr()

	var (
		validSendCoin   = validCoin(baseDenom, 100)
		invalidSendCoin = invalidCoin()
	)

	tscoin.SendCoin(addr1, addr2, validSendCoin, true)
	tscoin.SendCoin(addr1, addr2, invalidSendCoin, false)
	tscoin.SendCoin(emptyBalanceAddr, addr2, validSendCoin, false)
}

func TestMultiSendCoinHandler(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapHandlerGenesisTest(t, 4, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})

	tscoin := testcoin.NewHelper(t, ctx, dsc.CoinKeeper)
	addr1 := addrs[0]
	addr2 := addrs[1]
	addr3 := addrs[2]
	addr4 := addrs[3]
	emptyBalanceAddr := app.CreateAddr()

	var (
		coin1 = validCoin(baseDenom, 100)

		send1 = types.Send{
			Coin:     coin1,
			Receiver: addr2.String(),
		}
		send2 = types.Send{
			Coin:     coin1,
			Receiver: addr3.String(),
		}
		invalidSend1 = types.Send{
			Coin:     invalidCoin(),
			Receiver: addr4.String(),
		}
	)

	tscoin.MultiSendCoin(addr1, []types.Send{send1, send2}, true)
	tscoin.MultiSendCoin(addr1, []types.Send{invalidSend1}, false)
	tscoin.MultiSendCoin(emptyBalanceAddr, []types.Send{send1, send2}, false)
}

func TestBuyHandler(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapHandlerGenesisTest(t, 2, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})

	tscoin := testcoin.NewHelper(t, ctx, dsc.CoinKeeper)
	addr1 := addrs[0]
	addr2 := addrs[1]

	var (
		coinToBuy     = validCoin(baseDenom, 10)
		maxCoinToSell = validCoin(symbol, 10000)

		coinToBuy2     = validCoin(symbol, 100)
		maxCoinToSell2 = validCoin(baseDenom, 10)

		secondTestCoin = "buycointest"
	)

	_ = tscoin.CreateCoin(addr1, title, symbol, crr, helpers.EtherToWei(sdk.NewInt(10000000)), initReserve, limitVolume, "", true)
	_ = tscoin.CreateCoin(addr1, "Its second test coin", secondTestCoin, crr, helpers.EtherToWei(sdk.NewInt(10000)), helpers.EtherToWei(sdk.NewInt(10000000)), limitVolume, "", true)

	// valid requests
	tscoin.BuyCoin(addr1, coinToBuy, maxCoinToSell, true)
	tscoin.BuyCoin(addr1, coinToBuy2, maxCoinToSell2, true)

	// overflow limit volume
	tscoin.BuyCoin(addr1, validCoin(symbol, 10000000000000000), validCoin(baseDenom, 10), false)
	// coin to buy does not exist
	tscoin.BuyCoin(addr1, invalidCoin(), maxCoinToSell, false)
	// coin to sell does not exist
	tscoin.BuyCoin(addr1, coinToBuy, invalidCoin(), false)
	// base coin reserve in custom coin less than amount to buy
	tscoin.BuyCoin(addr1, validCoin(baseDenom, 100000), validCoin(symbol, 100000000), false)
	// custom coin reserve is less than amount to buy
	tscoin.BuyCoin(addr1, validCoin(secondTestCoin, 1000000000000), validCoin(symbol, 100000000), false)
	// maxAmountToSell is less than real amount to sell
	tscoin.BuyCoin(addr1, validCoin(baseDenom, 1000), validCoin(symbol, 100), false)
	// reserve after sell is less than minReserve for coin
	tscoin.BuyCoin(addr1, validCoin(baseDenom, 9001), validCoin(symbol, 1000000000), false)
	// addr dont have tokens for sell
	tscoin.BuyCoin(addr2, coinToBuy, maxCoinToSell, false)
}

func TestSellHadnler(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapHandlerGenesisTest(t, 2, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: helpers.EtherToWei(sdk.NewInt(100000000000000000)),
		},
	})

	tscoin := testcoin.NewHelper(t, ctx, dsc.CoinKeeper)
	addr1 := addrs[0]
	addr2 := addrs[1]

	var (
		secondCustomCoin = "secondcustomcoin"
	)

	_ = tscoin.CreateCoin(addr1, title, symbol, crr, helpers.EtherToWei(sdk.NewInt(10000000)), initReserve, limitVolume, "", true)
	_ = tscoin.CreateCoin(addr1, "Its second custom coin", secondCustomCoin, crr, helpers.EtherToWei(sdk.NewInt(10000000)), helpers.EtherToWei(sdk.NewInt(5000)), helpers.EtherToWei(sdk.NewInt(100000000000)), "", true)

	tscoin.SellCoin(addr1, validCoin(baseDenom, 100), validCoin(symbol, 10000), true)
	tscoin.SellCoin(addr1, validCoin(secondCustomCoin, 100000), validCoin(symbol, 1000), true)
	tscoin.SellCoin(addr1, validCoin(symbol, 100000), validCoin(symbol, 1000), true)

	// coin to buy does not exist
	tscoin.SellCoin(addr1, invalidCoin(), validCoin(baseDenom, 10000), false)
	// coin to sell does not exist
	tscoin.SellCoin(addr1, validCoin(baseDenom, 10000), invalidCoin(), false)
	// addr not have tokenst to sell
	tscoin.SellCoin(addr2, validCoin(baseDenom, 100), validCoin(symbol, 100000), false)
	// custom coin reserve less than minCoinReserve
	tscoin.SellCoin(addr1, validCoin(symbol, 10000000000), validCoin(baseDenom, 10000), false)
	// custom coin to sell reserve is less than minCoinReserve
	tscoin.SellCoin(addr1, validCoin(symbol, 10000000), validCoin(baseDenom, 9001), false)
	// custom coin to buy supply less than this limit volume
	tscoin.SellCoin(addr1, validCoin(baseDenom, 10000000000000000), validCoin(secondCustomCoin, 100000000000), false)
}

func TestSellAllHandler(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapHandlerGenesisTest(t, 1, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: helpers.EtherToWei(sdk.NewInt(10000000000000)),
		},
	})

	tscoin := testcoin.NewHelper(t, ctx, dsc.CoinKeeper)
	addr1 := addrs[0]

	_ = tscoin.CreateCoin(addr1, title, symbol, crr, helpers.EtherToWei(sdk.NewInt(10000000)), initReserve, limitVolume, "", true)

	tscoin.SellAllCoin(addr1, validCoin(baseDenom, 10000000000000), validCoin(symbol, 5000), true)
}

func TestRedeemHandler(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapHandlerGenesisTest(t, 2, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})
	tscoin := testcoin.NewHelper(t, ctx, dsc.CoinKeeper)
	addr1 := addrs[0]

	var ()

	check1, priv, proof := createNewCheck(t, ctx.ChainID(), "1000del", "9", "", 10)
	addCoinToAddr(t, ctx, dsc, sdk.AccAddress(priv.PubKey().Address()), sdk.Coins{validCoin(baseDenom, 1000000000000)})
	customCoin := tscoin.CreateCoin(addr1, title, symbol, crr, helpers.EtherToWei(sdk.NewInt(10000000)), initReserve, limitVolume, "", true)

	tscoin.CheckRedeem(sdk.AccAddress(priv.PubKey().Address()), checkToRlpString(t, check1), proof, true)

	// invalid check base58
	tscoin.CheckRedeem(sdk.AccAddress(priv.PubKey().Address()), "invalidCheck", proof, false)
	// invalid rlp check bytes
	tscoin.CheckRedeem(sdk.AccAddress(priv.PubKey().Address()), base58.Encode(priv.PubKey().Address()), proof, false)
	// invalid proof base64
	tscoin.CheckRedeem(sdk.AccAddress(priv.PubKey().Address()), checkToRlpString(t, check1), "invalidProof-%^", false)
	// invalid sender check
	invalidCheck1 := check1
	invR := sdk.NewInt(7)
	invalidCheck1.R = &invR
	tscoin.CheckRedeem(sdk.AccAddress(priv.PubKey().Address()), checkToRlpString(t, invalidCheck1), proof, false)
	// invalid coin Denom
	invalidCheck1.R = check1.R
	invalidCheck1.Coin = "invalidCoin"
	tscoin.CheckRedeem(sdk.AccAddress(priv.PubKey().Address()), checkToRlpString(t, invalidCheck1), proof, false)
	// addr balance less than checkCoinAmount
	invalidCheck1.Coin = baseDenom
	invalidCheck1.Amount = helpers.EtherToWei(sdk.NewInt(10000000000000))
	tscoin.CheckRedeem(sdk.AccAddress(priv.PubKey().Address()), checkToRlpString(t, invalidCheck1), proof, false)
	// addr custom coin balance less than checkCreateFeeAmount
	invalidCheck1.Coin = customCoin.Symbol
	invalidCheck1.Amount = helpers.EtherToWei(sdk.NewInt(100))
	tscoin.CheckRedeem(sdk.AccAddress(priv.PubKey().Address()), checkToRlpString(t, invalidCheck1), proof, false)
	// if custom coin enough, then baseCoin balance less than FeeAmount
	check2, priv2, proof2 := createNewCheck(t, ctx.ChainID(), fmt.Sprintf("1000%s", customCoin.Symbol), "9", "", 10)
	addCoinToAddr(t, ctx, dsc, sdk.AccAddress(priv2.PubKey().Address()), sdk.Coins{validCoin(customCoin.Symbol, 1000000000000)})
	tscoin.CheckRedeem(sdk.AccAddress(priv2.PubKey().Address()), checkToRlpString(t, check2), proof2, false)
	// addr base coin check amount enough, then balance less than checkAmount+feeAmount
	check3, priv3, proof3 := createNewCheck(t, ctx.ChainID(), fmt.Sprintf("1000000000000000000000000000000%s", baseDenom), "9", "", 10)
	addCoinToAddr(t, ctx, dsc, sdk.AccAddress(priv3.PubKey().Address()), sdk.Coins{validCoin(baseDenom, 1000000000000)})
	tscoin.CheckRedeem(sdk.AccAddress(priv3.PubKey().Address()), checkToRlpString(t, check3), proof3, false)
	// invalid chain_id
	check4, priv4, proof4 := createNewCheck(t, "notValidChainId-9000-1", fmt.Sprintf("100%s", baseDenom), "9", "", 10)
	addCoinToAddr(t, ctx, dsc, sdk.AccAddress(priv4.PubKey().Address()), sdk.Coins{validCoin(baseDenom, 1000000000000)})
	tscoin.CheckRedeem(sdk.AccAddress(priv4.PubKey().Address()), checkToRlpString(t, check4), proof4, false)
	// nonce len > 16
	check5, priv5, proof5 := createNewCheck(t, ctx.ChainID(), fmt.Sprintf("100%s", baseDenom), "189247582944712043914891349311313875902143479", "", 10)
	addCoinToAddr(t, ctx, dsc, sdk.AccAddress(priv5.PubKey().Address()), sdk.Coins{validCoin(baseDenom, 1000000000000)})
	tscoin.CheckRedeem(sdk.AccAddress(priv5.PubKey().Address()), checkToRlpString(t, check5), proof5, false)
	// this check has already been redeemed
	tscoin.CheckRedeem(sdk.AccAddress(priv.PubKey().Address()), checkToRlpString(t, check1), proof, false)
	// sender pubkey and checkAuthor pubkey not equal
	check6, priv6, proof6 := createNewCheck(t, ctx.ChainID(), "1000del", "10", "", 11)
	addCoinToAddr(t, ctx, dsc, sdk.AccAddress(priv6.PubKey().Address()), sdk.Coins{validCoin(baseDenom, 1000000000000)})
	tscoin.CheckRedeem(app.CreateAddr(), checkToRlpString(t, check6), proof6, false)
}

func TestLegacyReturn(t *testing.T) {
	app, ctx := bootstrapGenesisTest()

	// write genesis
	params := app.CoinKeeper.GetParams(ctx)

	coins := []types.Coin{
		{
			Title:  "del",
			Symbol: "del",
		},
		{
			Title:       "Test coin",
			Symbol:      "foo",
			CRR:         50,
			Reserve:     sdk.NewInt(1_000_000_000),
			Volume:      sdk.NewInt(1_000_000_000_0),
			LimitVolume: sdk.NewInt(1_000_000_000_000_000_000),
			Creator:     "uatom",
			Identity:    "dx1hs2wdrm87c92rzhq0vgmgrxr6u57xpr2lcygc2",
		},
	}

	var legacyAccounts = []struct {
		pubKeyBytes      []byte
		oldAddress       string
		newAddress       string
		fooLegacyBalance sdk.Int
		delLegacyBalance sdk.Int
		fooActualBalance sdk.Int
		delActualBalance sdk.Int
	}{
		{
			// account without new balance at moment of return
			// gasp history river forget aware wide dance velvet weather rain rail dry cliff assault coach jelly choose spirit shoulder isolate kidney outer trust message
			[]byte{0x3, 0x44, 0x8e, 0x6b, 0x3d, 0x50, 0xd6, 0xa3, 0x9c, 0xab, 0x3b, 0xab, 0xaa, 0x4a, 0xa2, 0xb0, 0x88, 0x5f, 0x55, 0x6f, 0xe0, 0x5d, 0x71, 0x49, 0x88, 0x5a, 0x5, 0xa0, 0xe7, 0x94, 0xa, 0x7e, 0x4f},
			"dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
			"dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
			sdk.NewInt(100),
			sdk.NewInt(200),
			sdk.NewInt(0),
			sdk.NewInt(0),
		},
		{
			// account with new balance at moment of return
			// section jeans evoke hockey result spell dish zero merge actress pink resource loan afford fitness install purity duck cannon ugly session stereo pattern spawn
			[]byte{0x3, 0x16, 0x18, 0x96, 0x7d, 0x77, 0xf1, 0xe4, 0x90, 0xd4, 0x1f, 0xc0, 0xe0, 0xc0, 0xc8, 0xb4, 0xb0, 0x47, 0x99, 0xe1, 0x16, 0x82, 0x4b, 0xea, 0x8, 0x47, 0x33, 0xe4, 0x63, 0x43, 0x20, 0xca, 0x3},
			"dx1m3eg7v6pu0dga2knj9zm4683dk9c8800j9nfw0",
			"dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f",
			sdk.NewInt(100),
			sdk.NewInt(200),
			sdk.NewInt(10),
			sdk.NewInt(20),
		},
	}
	stub, err := sdk.Bech32ifyAddressBytes(config.Bech32Prefix, types.StubCoinAddress)
	require.NoError(t, err, "StubCoinAddress to bech32")

	bankGenesisState := &cosmosBankTypes.GenesisState{
		Params: cosmosBankTypes.DefaultParams(),
		Balances: []cosmosBankTypes.Balance{
			{
				Address: stub,
				Coins: sdk.Coins{
					{
						Denom:  params.BaseSymbol,
						Amount: sdk.NewInt(400), // 200 + 200 from accounts above
					},
					{
						Denom:  "foo",
						Amount: sdk.NewInt(200), // 100 + 100 from account above
					},
				},
			},
		},
	}
	// add actual non-zero balances to new accounts
	for _, la := range legacyAccounts {
		if la.delActualBalance.IsZero() {
			continue
		}
		bankGenesisState.Balances = append(bankGenesisState.Balances,
			cosmosBankTypes.Balance{
				Address: la.newAddress,
				Coins: sdk.Coins{
					{
						Denom:  "del",
						Amount: la.delActualBalance,
					},
					{
						Denom:  "foo",
						Amount: la.fooActualBalance,
					},
				},
			},
		)
	}
	require.NoError(t, bankGenesisState.Validate(), "bankGenesisState")
	app.BankKeeper.InitGenesis(ctx, bankGenesisState)

	coinGenesisState := types.NewGenesisState(params, coins, []types.Check{}, []types.LegacyBalance{})
	for _, la := range legacyAccounts {
		coinGenesisState.LegacyBalances = append(coinGenesisState.LegacyBalances, types.LegacyBalance{
			OldAddress: la.oldAddress,
			Entries: []types.LegacyBalanceEntry{
				{
					CoinDenom: "del",
					Balance:   la.delLegacyBalance,
				},
				{
					CoinDenom: "foo",
					Balance:   la.fooLegacyBalance,
				},
			},
		})
	}

	require.NoError(t, coinGenesisState.Validate(), "coinGenesisState")
	coin.InitGenesis(ctx, app.CoinKeeper, coinGenesisState)

	sender, err := sdk.Bech32ifyAddressBytes(config.Bech32Prefix, []byte("sender"))
	require.NoError(t, err, "sender address")

	msgHandler := coin.NewHandler(app.CoinKeeper)
	//1. invalid public key
	msg := &types.MsgReturnLegacyBalance{
		Sender:         sender,
		Receiver:       legacyAccounts[0].newAddress,
		PublicKeyBytes: []byte{0},
	}

	require.Error(t, msg.ValidateBasic(), "for invalid public key must be error")

	// 2. return coins
	for _, la := range legacyAccounts {
		msg := &types.MsgReturnLegacyBalance{
			Sender:         sender,
			Receiver:       la.newAddress,
			PublicKeyBytes: la.pubKeyBytes,
		}
		_, err = msgHandler(ctx, msg)
		require.NoError(t, err, "error on return for '%s'", la.newAddress)
		// there must be nothing left in legacy
		lb, err := app.CoinKeeper.GetLegacyBalance(ctx, la.oldAddress)
		require.Error(t, err, "something leave for '%s': %+v", la.newAddress, lb)
		// check balances in bank
		receiver, err := sdk.AccAddressFromBech32(la.newAddress)
		require.NoError(t, err, "receiver address '%s'", la.newAddress)
		newCoins := app.BankKeeper.GetAllBalances(ctx, receiver)
		for _, coin := range newCoins {
			if coin.Denom == "del" {
				expectAmount := la.delActualBalance.Add(la.delLegacyBalance)
				require.True(t, expectAmount.Equal(coin.Amount),
					"expect amount for '%s' and coin 'del': %s != %s", la.newAddress, expectAmount.Sign(), coin.Amount.String())
			}
			if coin.Denom == "foo" {
				expectAmount := la.fooActualBalance.Add(la.fooLegacyBalance)
				require.True(t, expectAmount.Equal(coin.Amount),
					"expect amount for '%s' and coin 'foo': %s != %s", la.newAddress, expectAmount.Sign(), coin.Amount.String())
			}
		}
	}
}

///////////////////
// helper functions
///////////////////

func createNewCheck(t *testing.T, chainID, coinAmountStr, nonceStr, password string, dueBlock uint64) (types.Check, ethsecp256k1.PrivKey, string) {
	var (
		coinAmount, _ = sdk.ParseCoinNormalized(coinAmountStr)
		nonce, _      = sdk.NewIntFromString(nonceStr)
	)

	priv, _ := ethsecp256k1.GenerateKey()

	passphraseHash := sha256.Sum256([]byte(password))
	passphrasePrivKey, err := ethereumCrypto.ToECDSA(passphraseHash[:])
	require.NoError(t, err)

	check := &types.Check{
		ChainID:  chainID,
		Coin:     coinAmount.Denom,
		Amount:   coinAmount.Amount,
		Nonce:    nonce.BigInt().Bytes(),
		DueBlock: dueBlock,
	}

	checkHash := check.HashWithoutLock()
	lock, _ := ethereumCrypto.Sign(checkHash[:], passphrasePrivKey)
	check.Lock = lock

	// un armor key
	key, _ := priv.ToECDSA()

	checkHash = check.Hash()
	signature, _ := ethereumCrypto.Sign(checkHash[:], key)

	check.SetSignature(signature)

	// Prepare bytes to sign by private key generated from passphrase
	receiverAddressHash := make([]byte, 32)
	hw := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hw, []interface{}{
		sdk.AccAddress(priv.PubKey().Address()),
	})
	require.NoError(t, err)
	hw.Sum(receiverAddressHash[:0])

	// Sign receiver address by private key generated from passphrase
	proofSignature, err := ethereumCrypto.Sign(receiverAddressHash[:], passphrasePrivKey)
	require.NoError(t, err)
	proofBase64 := base64.StdEncoding.EncodeToString(proofSignature)

	return *check, *priv, proofBase64
}

func checkToRlpString(t *testing.T, check types.Check) string {
	checkBytes, err := rlp.EncodeToBytes(&check)
	require.NoError(t, err)

	return base58.Encode(checkBytes)
}

func addCoinToAddr(t *testing.T, ctx sdk.Context, dsc *app.DSC, addr sdk.AccAddress, coins sdk.Coins) {
	err := dsc.BankKeeper.MintCoins(ctx, types.ModuleName, coins)
	require.NoError(t, err)

	err = dsc.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
	require.NoError(t, err)
}
