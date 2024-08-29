package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/app"
	testkeeper "bitbucket.org/decimalteam/go-smart-node/testutil/keeper"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	cmdconfig "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/testutil"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	nftkeeper "bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	nfttestutil "bitbucket.org/decimalteam/go-smart-node/x/nft/testutil"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	sdkmath "cosmossdk.io/math"
)

var (
	pk          = ed25519.GenPrivKey().PubKey()
	addr        = sdk.AccAddress(pk.Address())
	defaultCoin = sdk.NewCoin(cmdcfg.BaseDenom, types.DefaultMinReserveAmount)
	baseDenom   = cmdconfig.BaseDenom
	baseAmount  = helpers.EtherToWei(sdkmath.NewInt(1000000000000))
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
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, addr, types.ReservedPool, sdk.NewCoins(defaultCoin)).AnyTimes().Return(nil)
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, addr, types.ReservedPool, sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(2))))).AnyTimes().Return(nil)
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ReservedPool, addr, sdk.NewCoins(defaultCoin)).AnyTimes().Return(nil)
	// --

	// -- create nft keeper
	space, ok := paramsKeeper.GetSubspace(types.ModuleName)
	s.Require().True(ok)
	keeper := nftkeeper.NewKeeper(
		encCfg.Codec,
		key,
		space,
		bankKeeper,
		nil,
		nil,
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

}

func bootstrapKeeperTest(t *testing.T, numAddrs int, accCoins sdk.Coins) (*app.DSC, sdk.Context, []sdk.AccAddress, []sdk.ValAddress) {
	_, dsc, ctx := testkeeper.GetBaseAppWithCustomKeeper(t)

	addrDels, addrVals := testkeeper.GenerateAddresses(dsc, ctx, numAddrs, accCoins)
	require.NotNil(t, addrDels)
	require.NotNil(t, addrVals)

	return dsc, ctx, addrDels, addrVals
}

func TestKeeper_Collection_Create(t *testing.T) {
	dsc, ctx, addrs, _ := bootstrapKeeperTest(t, 1, sdk.Coins{
		{
			Denom:  baseDenom,
			Amount: baseAmount,
		},
	})

	//creatorAddress, _ := types2.GetDecimalAddressFromHex("0xED5bd94A68B4EeB55F0516d1a2b51cfa6c2EE829")
	collection := types.Collection{
		Creator:    addrs[0].String(),
		Denom:      "symbol",
		Supply:     0,
		Tokens:     nil,
		TypeNft:    types.NftType_Unspecified,
		AddressDRC: "d01jnxuf8xk0gj20v595vchv40965exq4dx99x7de",
	}
	// write collection with it's counter
	dsc.NFTKeeper.SetCollection(ctx, collection)
	fmt.Println(collection)
	creator, err := dsc.NFTKeeper.CollectionsByCreator(ctx, &types.QueryCollectionsByCreatorRequest{
		Creator:    addrs[0].String(),
		Pagination: nil,
	})
	if err != nil {
		return
	}
	fmt.Println(creator)
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
