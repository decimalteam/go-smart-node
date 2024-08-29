package keeper_test

import (
	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/testutil"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/keeper"
	mstestutil "bitbucket.org/decimalteam/go-smart-node/x/multisig/testutil"
	"bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	feemarkettypes "github.com/decimalteam/ethermint/x/feemarket/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

var (
	pk1   = ed25519.GenPrivKey().PubKey()
	user1 = sdk.AccAddress(pk1.Address())
	pk2   = ed25519.GenPrivKey().PubKey()
	user2 = sdk.AccAddress(pk2.Address())
	pk3   = ed25519.GenPrivKey().PubKey()
	user3 = sdk.AccAddress(pk3.Address())
	pk4   = ed25519.GenPrivKey().PubKey()
	user4 = sdk.AccAddress(pk4.Address())

	defaultOwners, defaultWeights        = []string{user1.String(), user2.String(), user3.String()}, []uint32{1, 1, 1}
	defaultThreeshold             uint32 = 2
	defaultWallet, _                     = types.NewWallet(defaultOwners, defaultWeights, defaultThreeshold, []byte{})
	defaultWalletAddress, _              = sdk.AccAddressFromBech32(defaultWallet.Address)
	existsWallet, _                      = types.NewWallet(defaultOwners, defaultWeights, defaultThreeshold, []byte{112})
	existsWalletAddress, _               = sdk.AccAddressFromBech32(existsWallet.Address)
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	msKeeper   *keeper.Keeper
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
	paramsKeeper.Subspace(types.ModuleName)
	// --

	// -- create mock controller
	var emptyAccount authtypes.AccountI
	ctrl := gomock.NewController(s.T())
	bankKeeper := mstestutil.NewMockKeeper(ctrl)
	bankKeeper.EXPECT().GetAllBalances(ctx, defaultWalletAddress).AnyTimes().Return(sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, sdk.NewInt(1))))
	bankKeeper.EXPECT().GetAllBalances(ctx, existsWalletAddress).AnyTimes().Return(sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, sdk.NewInt(0))))
	bankKeeper.EXPECT().SendCoins(ctx, defaultWalletAddress, user1, sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, sdk.NewInt(1)))).AnyTimes().Return(nil)
	accountKeeper := mstestutil.NewMockAccountKeeperI(ctrl)
	accountKeeper.EXPECT().GetAccount(ctx, defaultWalletAddress).AnyTimes().Return(nil)
	accountKeeper.EXPECT().GetAccount(ctx.WithTxBytes([]byte{1}), existsWalletAddress).AnyTimes().Return(emptyAccount)
	// --

	// --
	router := baseapp.NewMsgServiceRouter()
	// --

	// -- create multisig keeper
	space, ok := paramsKeeper.GetSubspace(types.ModuleName)
	s.Require().True(ok)
	k := keeper.NewKeeper(
		encCfg.Codec,
		key,
		space,
		accountKeeper,
		bankKeeper,
		router,
	)
	// --

	s.ctx = ctx
	s.msKeeper = k
	s.bankKeeper = bankKeeper

	// -- register services
	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.msKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
	s.msgServer = k
	//
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
