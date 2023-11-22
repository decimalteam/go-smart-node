package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/decimalteam/ethermint/crypto/ethsecp256k1"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/testutil"
	commonTypes "bitbucket.org/decimalteam/go-smart-node/types"
	legacykeeper "bitbucket.org/decimalteam/go-smart-node/x/legacy/keeper"
	legacytestutil "bitbucket.org/decimalteam/go-smart-node/x/legacy/testutil"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
)

var (
	baseDenom = config.BaseDenom
	publicKey = []byte{0x3, 0x44, 0x8e, 0x6b, 0x3d, 0x50, 0xd6, 0xa3, 0x9c, 0xab, 0x3b, 0xab, 0xaa,
		0x4a, 0xa2, 0xb0, 0x88, 0x5f, 0x55, 0x6f, 0xe0, 0x5d, 0x71, 0x49, 0x88, 0x5a, 0x5, 0xa0, 0xe7, 0x94, 0xa, 0x7e, 0x4f}
	pk               = ethsecp256k1.PubKey{Key: publicKey}
	pkOld            = secp256k1.PubKey{Key: publicKey}
	oldAddress, _    = commonTypes.GetLegacyAddressFromPubKey(publicKey)
	oldSdkAddress    = pkOld.Address()
	newAddress       = sdk.AccAddress(pk.Address())
	actualAddress, _ = bech32.ConvertAndEncode(config.Bech32Prefix, newAddress)
)

// Init global cosmos sdk cmdcfg
func initConfig() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(cmdcfg.Bech32PrefixAccAddr, cmdcfg.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(cmdcfg.Bech32PrefixValAddr, cmdcfg.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(cmdcfg.Bech32PrefixConsAddr, cmdcfg.Bech32PrefixConsPub)
}

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	legacyKeeper   *legacykeeper.Keeper
	bankKeeper     types.BankKeeper
	nftKeeper      types.NftKeeper
	multisigKeeper types.MultisigKeeper

	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func (s *KeeperTestSuite) SetupTest() {
	initConfig()

	key := sdk.NewKVStoreKey(types.StoreKey)
	keys := []storetypes.StoreKey{
		key,
	}

	tkey := sdk.NewTransientStoreKey("transient_test")
	tkeys := []storetypes.StoreKey{
		tkey,
	}

	testCtx := testutil.DefaultContextWithDB(s.T(), keys, tkeys)
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})
	encCfg := testutil.MakeTestEncodingConfig()

	// -- create mock controller
	ctrl := gomock.NewController(s.T())
	bankKeeper := legacytestutil.NewMockBankKeeper(ctrl)
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, types.LegacyCoinPool, newAddress, defaultRecord.Coins).AnyTimes().Return(nil)
	bankKeeper.EXPECT().GetAllBalances(ctx, sdk.AccAddress(oldSdkAddress)).AnyTimes().Return(defaultRewards)
	bankKeeper.EXPECT().SendCoins(ctx, sdk.AccAddress(oldSdkAddress), newAddress, defaultRewards).AnyTimes().Return(nil)
	bankKeeper.EXPECT().GetAllBalances(ctx, gomock.Any()).AnyTimes().Return(sdk.NewCoins())
	nftKeeper := legacytestutil.NewMockNftKeeper(ctrl)
	nftKeeper.EXPECT().GetToken(ctx, defaultTokenID).AnyTimes().Return(nfttypes.Token{Denom: defaultTokenID}, true)
	nftKeeper.EXPECT().GetSubTokens(ctx, defaultTokenID).AnyTimes().Return(defaultSubTokensBefore)
	nftKeeper.EXPECT().ReplaceSubTokenOwner(ctx, defaultTokenID, defaultSubTokensAfter[0].ID, actualAddress).AnyTimes().Return(nil)
	multisigKeeper := legacytestutil.NewMockMultisigKeeper(ctrl)
	multisigKeeper.EXPECT().GetWallet(ctx, defaultMultisigWalletBefore.Address).AnyTimes().Return(defaultMultisigWalletBefore, nil)
	multisigKeeper.EXPECT().SetWallet(ctx, defaultMultisigWalletAfter).AnyTimes()
	validatorKeeper := legacytestutil.NewMockValidatorKeeper(ctrl)
	validatorKeeper.EXPECT().GetValidator(ctx, defaultOperatorAddressSdk).AnyTimes().Return(defaultValidatorBefore, true)
	validatorKeeper.EXPECT().SetValidator(ctx, defaultValidatorAfter).AnyTimes()

	// --

	// -- create nft keeper
	k := legacykeeper.NewKeeper(
		encCfg.Codec,
		key,
		bankKeeper,
		nftKeeper,
		multisigKeeper,
		validatorKeeper,
	)
	// --

	s.ctx = ctx
	s.legacyKeeper = k
	s.bankKeeper = bankKeeper

	// -- register services
	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.legacyKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
	s.msgServer = k
	//
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
