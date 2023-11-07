package keeper_test

import (
	"encoding/hex"
	"testing"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/testutil"
	"bitbucket.org/decimalteam/go-smart-node/x/swap/keeper"
	swaptestutil "bitbucket.org/decimalteam/go-smart-node/x/swap/testutil"
	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/decimalteam/ethermint/crypto/ethsecp256k1"
	feemarkettypes "github.com/decimalteam/ethermint/x/feemarket/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

var (
	privkey, _ = ethsecp256k1.GenerateKey()

	pk1   = ed25519.GenPrivKey().PubKey()
	user1 = sdk.AccAddress(pk1.Address())
	pk2   = ed25519.GenPrivKey().PubKey()
	user2 = sdk.AccAddress(pk2.Address())
	pk3   = ed25519.GenPrivKey().PubKey()
	user3 = sdk.AccAddress(pk3.Address())

	defaultTokenSymbol string = cmdcfg.BaseDenom
	defaultAmount             = sdk.NewInt(10)
	defaultCoins              = sdk.NewCoins(sdk.NewCoin(defaultTokenSymbol, defaultAmount))

	defaultChainID   uint32 = 1
	defaultChainName        = "test_chain1"
	destChainID      uint32 = 2
	destChainName           = "test_chain2"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	swapKeeper *keeper.Keeper
	bankKeeper bankkeeper.Keeper

	queryClient   types.QueryClient
	fmQueryClient feemarkettypes.QueryClient
	msgServer     types.MsgServer
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
	poolAddr := cosmosAuthTypes.NewModuleAddress(types.SwapPool)
	ctrl := gomock.NewController(s.T())
	bankKeeper := swaptestutil.NewMockBankKeeper(ctrl)
	bankKeeper.EXPECT().GetAllBalances(ctx, user2).AnyTimes().Return(defaultCoins)
	bankKeeper.EXPECT().GetAllBalances(ctx, user3).AnyTimes().Return(sdk.NewCoins())
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, user2, types.SwapPool, defaultCoins).AnyTimes().Return(nil)
	bankKeeper.EXPECT().GetAllBalances(ctx, poolAddr).AnyTimes().Return(defaultCoins)
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, types.SwapPool, user3, defaultCoins).AnyTimes().Return(nil)
	accountKeeper := swaptestutil.NewMockAccountKeeperI(ctrl)
	// --

	// -- create nft keeper
	space, ok := paramsKeeper.GetSubspace(types.ModuleName)
	s.Require().True(ok)
	k := keeper.NewKeeper(
		encCfg.Codec,
		key,
		space,
		accountKeeper,
		bankKeeper,
	)
	dp := types.DefaultParams()
	dp.ServiceAddress = user1.String()
	dp.CheckingAddress = hex.EncodeToString(privkey.PubKey().Address())
	k.SetParams(ctx, dp)
	// --

	s.ctx = ctx
	s.swapKeeper = k
	s.bankKeeper = bankKeeper

	// -- register services
	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.swapKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
	s.msgServer = k
	//
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
