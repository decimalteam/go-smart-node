package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/testutil"
	nftkeeper "bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	nfttestutil "bitbucket.org/decimalteam/go-smart-node/x/nft/testutil"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	//simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	"testing"
)

var (
//PKs          = simtestutil.CreateTestPubKeys(500) =
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	nftKeeper  nftkeeper.Keeper
	bankKeeper bankkeeper.Keeper

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
	bankKeeper := nfttestutil.NewMockKeeper(ctrl)
	//bankKeeper.EXPECT()
	// --

	// -- create nft keeper
	space, ok := paramsKeeper.GetSubspace(types.ModuleName)
	s.Require().True(ok)
	keeper := nftkeeper.NewKeeper(
		encCfg.Codec,
		key,
		space,
		bankKeeper,
	)
	keeper.SetParams(ctx, types.DefaultParams())
	// --

	s.ctx = ctx
	s.nftKeeper = *keeper
	s.bankKeeper = bankKeeper

	// -- register services
	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.nftKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
	s.msgServer = keeper
	//
}

func (s *KeeperTestSuite) TestParams() {
	ctx, keeper := s.ctx, s.nftKeeper
	require := s.Require()

	expParams := types.DefaultParams()
	expParams.MaxCollectionSize = 555
	expParams.MinReserveAmount = math.NewInt(111)
	keeper.SetParams(ctx, expParams)
	resParams := keeper.GetParams(ctx)
	require.True(expParams.Equal(resParams))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
