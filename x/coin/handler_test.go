package coin_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"

	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/testcoin"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"github.com/cosmos/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereumCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
)

func bootstrapHandlerTest(t *testing.T, numAddrs int, accCoins sdk.Coins) (*app.DSC, sdk.Context, []sdk.AccAddress, []sdk.ValAddress) {
	_, dsc, ctx := testkeeper.GetTestAppWithCoinKeeper(t)

	addrDels, addrVals := testkeeper.GenerateAddresses(dsc, ctx, numAddrs, accCoins)
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
	dsc, ctx, addrs, _ := bootstrapHandlerTest(t, 2, sdk.Coins{
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
	dsc, ctx, addrs, _ := bootstrapHandlerTest(t, 2, sdk.Coins{
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
	dsc, ctx, addrs, _ := bootstrapHandlerTest(t, 2, sdk.Coins{
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
	dsc, ctx, addrs, _ := bootstrapHandlerTest(t, 4, sdk.Coins{
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
	dsc, ctx, addrs, _ := bootstrapHandlerTest(t, 2, sdk.Coins{
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

func TestSellHandler(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapHandlerTest(t, 2, sdk.Coins{
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
	dsc, ctx, addrs, _ := bootstrapHandlerTest(t, 1, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: helpers.EtherToWei(sdk.NewInt(10000000000000)),
		},
	})

	tscoin := testcoin.NewHelper(t, ctx, dsc.CoinKeeper)
	addr1 := addrs[0]

	_ = tscoin.CreateCoin(addr1, title, symbol, crr, helpers.EtherToWei(sdk.NewInt(10000000)), initReserve, limitVolume, "", true)

	tscoin.SellAllCoin(addr1, baseDenom, validCoin(symbol, 5000), true)
}

func TestBurnCoinHandler(t *testing.T) {
	const customSymbol = "somecoin"
	var customVolume = helpers.EtherToWei(sdk.NewInt(2000))
	var customReserve = helpers.EtherToWei(sdk.NewInt(1000))

	dsc, ctx, addrs, _ := bootstrapHandlerTest(t, 2, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})

	handler := coin.NewHandler(dsc.CoinKeeper)
	_, err := handler(ctx, types.NewMsgCreateCoin(
		addrs[0],
		"somecoin",
		customSymbol,
		10,
		customVolume.Mul(sdk.NewInt(10)),
		customReserve,
		customVolume.Mul(sdk.NewInt(100)),
		"",
	))
	require.NoError(t, err, "create coin")
	balance := dsc.BankKeeper.GetBalance(ctx, addrs[0], customSymbol)
	require.True(t, balance.Amount.Equal(customVolume.Mul(sdk.NewInt(10))), "balance: %s", balance.String())

	_, err = handler(ctx, types.NewMsgBurnCoin(
		addrs[0],
		sdk.NewCoin(customSymbol, customVolume),
	))
	require.NoError(t, err, "burn coin")
	balance = dsc.BankKeeper.GetBalance(ctx, addrs[0], customSymbol)
	require.True(t, balance.Amount.Equal(customVolume.Mul(sdk.NewInt(9))), "balance: %s", balance.String())
	inf, err := dsc.CoinKeeper.GetCoin(ctx, customSymbol)
	require.NoError(t, err, "coin info")
	require.True(t, inf.Reserve.Equal(customReserve), "check reserve")

	//try to burn to break limits
	_, err = handler(ctx, types.NewMsgBurnCoin(
		addrs[0],
		sdk.NewCoin(customSymbol, customVolume.Mul(sdk.NewInt(9))),
	))
	require.Error(t, err, "overburn coin")
	// balance must be same
	balance = dsc.BankKeeper.GetBalance(ctx, addrs[0], customSymbol)
	require.True(t, balance.Amount.Equal(customVolume.Mul(sdk.NewInt(9))), "balance: %s", balance.String())

	// burn to minimal volume
	balance = dsc.BankKeeper.GetBalance(ctx, addrs[0], customSymbol)
	volumeToBurn := balance.Amount.Sub(types.MinCoinSupply)
	_, err = handler(ctx, types.NewMsgBurnCoin(
		addrs[0],
		sdk.NewCoin(customSymbol, volumeToBurn),
	))
	require.NoError(t, err, "burn coin to minimum")
	inf, err = dsc.CoinKeeper.GetCoin(ctx, customSymbol)
	require.NoError(t, err, "coin info")

	// this call check MinCoinSupply after burn
	// If MinCoinSupply is too small, there will be panic
	formulas.CalculatePurchaseAmount(inf.Volume, inf.Reserve, uint(inf.CRR), helpers.EtherToWei(sdk.NewInt(1)))
	formulas.CalculatePurchaseAmount(inf.Volume, inf.Reserve, uint(inf.CRR), helpers.FinneyToWei(sdk.NewInt(1)))

	////////
	// check base coin burning
	_, err = handler(ctx, types.NewMsgBurnCoin(
		addrs[1],
		sdk.NewCoin(baseDenom, helpers.EtherToWei(sdk.NewInt(1))),
	))
	require.NoError(t, err, "burn base coin")
}

func TestRedeemHandler(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapHandlerTest(t, 2, sdk.Coins{
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
