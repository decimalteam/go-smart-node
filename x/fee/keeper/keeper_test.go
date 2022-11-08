package keeper_test

import (
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/testutil"
	feekeeper "bitbucket.org/decimalteam/go-smart-node/x/fee/keeper"
	feetestutil "bitbucket.org/decimalteam/go-smart-node/x/fee/testutil"
	"bitbucket.org/decimalteam/go-smart-node/x/fee/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

var (
	baseDenom = "del"
	pk        = ed25519.GenPrivKey().PubKey()
	oracle    = sdk.AccAddress(pk.Address()).String()
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	feeKeeper  feekeeper.Keeper
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
	ctrl := gomock.NewController(s.T())
	bankKeeper := feetestutil.NewMockKeeper(ctrl)
	coinKeeper := feetestutil.NewMockCoinKeeper(ctrl)
	authKeeper := feetestutil.NewMockAccountKeeper(ctrl)
	// --

	// -- create nft keeper
	space, ok := paramsKeeper.GetSubspace(types.ModuleName)
	s.Require().True(ok)
	k := feekeeper.NewKeeper(
		encCfg.Codec,
		key,
		space,
		bankKeeper,
		coinKeeper,
		authKeeper,
		baseDenom,
	)
	dp := types.DefaultParams()
	dp.Oracle = oracle
	k.SetModuleParams(ctx, dp)
	// --

	s.ctx = ctx
	s.feeKeeper = *k
	s.bankKeeper = bankKeeper

	// -- register services
	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.feeKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
	feemarkettypes.RegisterQueryServer(queryHelper, k)
	s.fmQueryClient = feemarkettypes.NewQueryClient(queryHelper)
	s.msgServer = k
	//
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
