package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/precompile/drc20cosmos"
	"context"
	"fmt"
	"github.com/decimalteam/ethermint/crypto/ethsecp256k1"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/testutil"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	coinkeeper "bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/testcoin"
	cointestutil "bitbucket.org/decimalteam/go-smart-node/x/coin/testutil"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cmdconfig "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

var (
	pk    = ed25519.GenPrivKey().PubKey()
	addr  = sdk.AccAddress(pk.Address())
	pk1   = ed25519.GenPrivKey().PubKey()
	addr1 = sdk.AccAddress(pk.Address())

	baseDenom  = cmdconfig.BaseDenom
	baseAmount = helpers.EtherToWei(sdkmath.NewInt(1000000000000))
	//baseInitialReserve = sdk.NewCoin(baseDenom, helpers.EtherToWei(10000))
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	coinKeeper keeper.Keeper
	bankKeeper bankkeeper.Keeper
	acKeeper   authkeeper.AccountKeeperI

	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func (s *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(types.StoreKey)
	paramkey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	keys := []storetypes.StoreKey{
		key, paramkey,
	}

	tkey := sdk.NewTransientStoreKey("transient_test")
	tparamskey := sdk.NewTransientStoreKey("transient_param_test")
	tkeys := []storetypes.StoreKey{
		tkey, tparamskey,
	}

	testCtx := testutil.DefaultContextWithDB(s.T(), keys, tkeys)
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})
	encCfg := testutil.MakeTestEncodingConfig()

	// -- create params keeper
	paramsKeeper := paramskeeper.NewKeeper(
		encCfg.Codec,
		encCfg.Amino,
		paramkey,
		tparamskey,
	)
	paramsKeeper.Subspace(types.ModuleName).WithKeyTable(types.ParamKeyTable())
	// --

	// -- create mock controller
	ctrl := gomock.NewController(s.T())
	bankKeeper := cointestutil.NewMockKeeper(ctrl)
	acKeeper := cointestutil.NewMockAccountKeeperI(ctrl)
	fkKeeper := cointestutil.NewMockFeeMarketKeeper(ctrl)
	// --

	// -- create nft keeper
	space, ok := paramsKeeper.GetSubspace(types.ModuleName)
	s.Require().True(ok)
	k := coinkeeper.NewKeeper(
		encCfg.Codec,
		key,
		space,
		acKeeper,
		fkKeeper,
		bankKeeper,
		nil,
	)
	k.SetParams(ctx, types.DefaultParams())
	// --

	s.ctx = ctx
	s.coinKeeper = *k
	s.acKeeper = acKeeper
	s.bankKeeper = bankKeeper

	// -- register services
	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.coinKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
	s.msgServer = k
	//
}

func bootstrapKeeperTest(t *testing.T, numAddrs int, accCoins sdk.Coins) (*app.DSC, sdk.Context, []sdk.AccAddress, []sdk.ValAddress) {
	_, dsc, ctx := testkeeper.GetTestAppWithCoinKeeper(t)

	addrDels, addrVals := testkeeper.GenerateAddresses(dsc, ctx, numAddrs, accCoins)
	require.NotNil(t, addrDels)
	require.NotNil(t, addrVals)

	return dsc, ctx, addrDels, addrVals
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func TestKeeper_Coin(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapKeeperTest(t, 1, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})

	//paramsCtx := dsc.EvmKeeper.GetParams(ctx)
	//ethCfg := paramsCtx.ChainConfig.EthereumConfig(dsc.EvmKeeper.ChainID())

	//signer1 := ethtypes.MakeSigner(ethCfg, big.NewInt(ctx.BlockHeight()))

	denom := "testcoin"
	newCoin := types.Coin{
		Denom:        denom,
		Title:        "test keeper coin functions coin",
		CRR:          50,
		Reserve:      helpers.EtherToWei(sdkmath.NewInt(5000)),
		Volume:       helpers.EtherToWei(sdkmath.NewInt(10000)),
		LimitVolume:  helpers.EtherToWei(sdkmath.NewInt(1000000000)),
		MinVolume:    sdk.ZeroInt(),
		Creator:      addrs[0].String(),
		Identity:     "",
		Drc20Address: "",
	}

	// check set coin
	dsc.CoinKeeper.SetCoin(ctx, newCoin)

	// account key, use a constant account to keep unit test deterministic.
	ecdsaPriv, err := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	require.NoError(t, err)
	priv := &ethsecp256k1.PrivKey{
		Key: crypto.FromECDSA(ecdsaPriv),
	}
	address := common.BytesToAddress(priv.PubKey().Address().Bytes())
	//signer2 := tests.NewSigner(priv)

	contractCreateTx := &ethtypes.AccessListTx{
		GasPrice: nil,
		Gas:      53000,
		To:       nil,
		Data:     []byte("contract_data"),
		Nonce:    1,
	}
	ethTx := ethtypes.NewTx(contractCreateTx)
	ethMsg := &evmtypes.MsgEthereumTx{}
	ethMsg.FromEthereumTx(ethTx)
	ethMsg.From = address.Hex()
	fmt.Print(address.Hex())
	//ethMsg.Sign(signer1, nil)

	dsc.EvmKeeper.SetBalance(ctx, common.HexToAddress(drc20cosmos.AddressForContractOwner), big.NewInt(100000000000000))

	drc20, err1 := drc20cosmos.NewDrc20Cosmos(ctx, dsc.EvmKeeper, dsc.BankKeeper, newCoin)
	if err1 != nil {
		ctx.Logger().Info(err1.Error())
	}

	_, err2 := drc20.CreateContractIfNotSet()
	if err2 != nil {
		ctx.Logger().Info(err.Error())
	}

	newCoin = drc20.Coin
	dsc.CoinKeeper.SetCoin(ctx, newCoin)

	// check get exist coin
	getCoin, err := dsc.CoinKeeper.GetCoin(ctx, denom)
	require.NoError(t, err)
	require.True(t, getCoin.Equal(newCoin))
	// check get not exist coin
	_, err = dsc.CoinKeeper.GetCoin(ctx, "not exist coin")
	require.Error(t, err)
	// check get coins
	coins := dsc.CoinKeeper.GetCoins(ctx)
	require.Equal(t, 2, len(coins))

	// update coin volume and reserve
	dsc.CoinKeeper.UpdateCoinVR(ctx, getCoin.Denom, helpers.EtherToWei(sdkmath.NewInt(10002)), helpers.EtherToWei(sdkmath.NewInt(1000000001)))
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
	//custCoinAmount := helpers.EtherToWei(sdkmath.NewInt(10000))

	dsc, ctx, addrs, _ := bootstrapKeeperTest(t, 1, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})

	newCoin := types.Coin{
		Denom:       custCoinDenom,
		Title:       "test keeper coin functions coin",
		CRR:         50,
		Reserve:     helpers.EtherToWei(sdkmath.NewInt(5000)),
		Volume:      helpers.EtherToWei(sdkmath.NewInt(10000)),
		LimitVolume: helpers.EtherToWei(sdkmath.NewInt(1000000000)),
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
	_, _, err := dsc.CoinKeeper.GetCommission(ctx, helpers.EtherToWei(sdkmath.NewInt(10)))
	require.NoError(t, err)

	// fee with custom coin
	ctxWithFee := ctx
	ctxWithFee = ctx.WithContext(context.WithValue(ctx.Context(), types.ContextFeeKey{}, sdk.Coins{
		{
			Denom:  custCoinDenom,
			Amount: helpers.EtherToWei(sdkmath.NewInt(100)),
		},
	}))
	require.NotNil(t, ctxWithFee.Context())

	_, _, err = dsc.CoinKeeper.GetCommission(ctxWithFee, helpers.EtherToWei(sdkmath.NewInt(10)))
	require.NoError(t, err)

	// fee custom coin not exist
	ctxWithNotExistCoinFee := ctx
	ctxWithNotExistCoinFee = ctx.WithContext(context.WithValue(ctx.Context(), types.ContextFeeKey{}, sdk.Coins{
		{
			Denom:  "notexistcoin",
			Amount: helpers.EtherToWei(sdkmath.NewInt(100)),
		},
	}))
	require.NotNil(t, ctxWithNotExistCoinFee.Context())

	_, _, err = dsc.CoinKeeper.GetCommission(ctxWithNotExistCoinFee, helpers.EtherToWei(sdkmath.NewInt(10)))
	require.Error(t, err)

	// fee custom coin reserve less than need fee base coin amount
	ctxWithLessReserveFee := ctx
	ctxWithLessReserveFee = ctx.WithContext(context.WithValue(ctx.Context(), types.ContextFeeKey{}, sdk.Coins{
		{
			Denom:  custCoinDenom,
			Amount: helpers.EtherToWei(sdkmath.NewInt(1000)),
		},
	}))
	require.NotNil(t, ctxWithLessReserveFee.Context())

	_, _, err = dsc.CoinKeeper.GetCommission(ctxWithLessReserveFee, helpers.EtherToWei(sdkmath.NewInt(1000000000)))
	require.Error(t, err)
}
