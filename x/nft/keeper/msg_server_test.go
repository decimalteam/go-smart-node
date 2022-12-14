package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

func (s *KeeperTestSuite) TestMsgMintToken() {
	ctx, _, msgServer := s.ctx, s.nftKeeper, s.msgServer
	require := s.Require()

	var (
		denom        = "test_msg_mint_1"
		ID           = "mint_token_1"
		ID2          = "mint_token_2"
		ID3          = "mint_token_3"
		pk           = ed25519.GenPrivKey().PubKey()
		invalidOwner = sdk.AccAddress(pk.Address())
	)

	testCases := []struct {
		name   string
		input  *types.MsgMintToken
		expErr bool
	}{
		{
			"Valid request",
			types.NewMsgMintToken(addr, denom, ID, ID, false, addr, 1, defaultCoin),
			false,
		},
		{
			"Valid request - second token mint",
			types.NewMsgMintToken(addr, denom, ID2, ID2, true, addr, 1, defaultCoin),
			false,
		},
		{
			"two collection - one token",
			types.NewMsgMintToken(addr, "invalid denom", ID2, ID2, true, addr, 1, defaultCoin),
			true,
		},
		{
			"Creator not owner",
			types.NewMsgMintToken(invalidOwner, denom, ID2, ID2, true, invalidOwner, 1, defaultCoin),
			true,
		},
		{
			"Not unique URI if token exists",
			types.NewMsgMintToken(addr, denom, ID3, ID, true, addr, 1, defaultCoin),
			true,
		},
		{
			"Not allowed to mint",
			types.NewMsgMintToken(addr, denom, ID, ID, false, addr, 1, defaultCoin),
			true,
		},
		{
			"New token reserve is less than min",
			types.NewMsgMintToken(addr, denom, ID3, ID3, true, addr, 1, sdk.NewCoin(cmdcfg.BaseDenom, sdk.NewInt(1000))),
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.MintToken(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurnToken() {
	ctx, k, msgServer := s.ctx, s.nftKeeper, s.msgServer
	require := s.Require()

	var (
		denom        = "test_msg_burn"
		ID           = "msg_burn_1"
		pk           = ed25519.GenPrivKey().PubKey()
		invalidOwner = sdk.AccAddress(pk.Address())
	)

	collection := types.Collection{
		Creator: addr.String(),
		Denom:   denom,
	}
	token := types.Token{
		Creator:   addr.String(),
		Denom:     denom,
		ID:        ID,
		URI:       ID,
		Reserve:   defaultCoin,
		AllowMint: true,
		Minted:    0,
		Burnt:     0,
	}
	k.CreateToken(ctx, collection, token)
	for i := 1; i < 10; i++ {
		k.SetSubToken(ctx, ID, types.SubToken{
			ID:      uint32(i),
			Owner:   addr.String(),
			Reserve: &defaultCoin,
		})
	}
	k.SetSubToken(ctx, ID, types.SubToken{
		ID:      uint32(11),
		Owner:   invalidOwner.String(),
		Reserve: &defaultCoin,
	})
	testCases := []struct {
		name   string
		input  *types.MsgBurnToken
		expErr bool
	}{
		{
			"valid request",
			types.NewMsgBurnToken(addr, ID, []uint32{1}),
			false,
		},
		{
			"not exists sub token",
			types.NewMsgBurnToken(addr, ID, []uint32{1000}),
			true,
		},
		{
			"not exists token",
			types.NewMsgBurnToken(addr, "invalid ID", []uint32{1}),
			true,
		},
		{
			"not token owner",
			types.NewMsgBurnToken(invalidOwner, ID, []uint32{1}),
			true,
		},
		{
			"owned token, but not owned subtoken",
			types.NewMsgBurnToken(addr, ID, []uint32{11}),
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.BurnToken(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgSendToken() {
	ctx, k, msgServer := s.ctx, s.nftKeeper, s.msgServer
	require := s.Require()

	var (
		denom     = "test_msg_send"
		ID        = "msg_send_1"
		pk        = ed25519.GenPrivKey().PubKey()
		recipient = sdk.AccAddress(pk.Address())
	)

	collection := types.Collection{
		Creator: addr.String(),
		Denom:   denom,
	}
	token := types.Token{
		Creator:   addr.String(),
		Denom:     denom,
		ID:        ID,
		URI:       ID,
		Reserve:   defaultCoin,
		AllowMint: true,
		Minted:    0,
		Burnt:     0,
	}
	k.CreateToken(ctx, collection, token)
	for i := 1; i < 10; i++ {
		k.SetSubToken(ctx, ID, types.SubToken{
			ID:      uint32(i),
			Owner:   addr.String(),
			Reserve: &defaultCoin,
		})
	}
	k.SetSubToken(ctx, ID, types.SubToken{
		ID:      11,
		Owner:   recipient.String(),
		Reserve: &defaultCoin,
	})

	testCases := []struct {
		name   string
		input  *types.MsgSendToken
		expErr bool
	}{
		{
			"valid request",
			types.NewMsgSendToken(addr, recipient, ID, []uint32{1, 2, 3}),
			false,
		},
		{
			"not exists token",
			types.NewMsgSendToken(addr, recipient, "invalid token", []uint32{1, 2, 3}),
			true,
		},
		{
			"not exists subtokens",
			types.NewMsgSendToken(addr, recipient, ID, []uint32{4, 5, 133}),
			true,
		},
		{
			"not owned subtoken",
			types.NewMsgSendToken(addr, recipient, ID, []uint32{11}),
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.SendToken(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgUpdateToken() {
	ctx, k, msgServer := s.ctx, s.nftKeeper, s.msgServer
	require := s.Require()

	var (
		denom        = "test_msg_update_token"
		ID           = "msg_update_token_1"
		pk           = ed25519.GenPrivKey().PubKey()
		invalidOwner = sdk.AccAddress(pk.Address())
	)

	collection := types.Collection{
		Creator: addr.String(),
		Denom:   denom,
	}
	token := types.Token{
		Creator:   addr.String(),
		Denom:     denom,
		ID:        ID,
		URI:       ID,
		Reserve:   defaultCoin,
		AllowMint: true,
	}
	k.CreateToken(ctx, collection, token)

	testCases := []struct {
		name   string
		input  *types.MsgUpdateToken
		expErr bool
	}{
		{
			"valid request",
			types.NewMsgUpdateToken(addr, ID, "new_uri"),
			false,
		},
		{
			"new uri is identical to the old one",
			types.NewMsgUpdateToken(addr, ID, "new_uri"),
			true,
		},
		{
			"not owned token",
			types.NewMsgUpdateToken(invalidOwner, ID, "uri"),
			true,
		},
		{
			"not exists token",
			types.NewMsgUpdateToken(addr, "invalid token", "uri"),
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.UpdateToken(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgUpdateReserve() {
	ctx, k, msgServer := s.ctx, s.nftKeeper, s.msgServer
	require := s.Require()

	var (
		denom        = "test_msg_update_reserve"
		ID           = "msg_update_reserve_1"
		pk           = ed25519.GenPrivKey().PubKey()
		invalidOwner = sdk.AccAddress(pk.Address())
	)

	collection := types.Collection{
		Creator: addr.String(),
		Denom:   denom,
	}
	token := types.Token{
		Creator:   addr.String(),
		Denom:     denom,
		ID:        ID,
		URI:       ID,
		Reserve:   defaultCoin,
		AllowMint: true,
	}
	k.CreateToken(ctx, collection, token)
	for i := 1; i < 10; i++ {
		k.SetSubToken(ctx, ID, types.SubToken{
			ID:      uint32(i),
			Owner:   addr.String(),
			Reserve: &defaultCoin,
		})
	}
	k.SetSubToken(ctx, ID, types.SubToken{
		ID:      11,
		Owner:   invalidOwner.String(),
		Reserve: &defaultCoin,
	})

	testCases := []struct {
		name   string
		input  *types.MsgUpdateReserve
		expErr bool
	}{
		{
			"valid request",
			types.NewMsgUpdateReserve(addr, ID, []uint32{1, 2}, sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(2)))),
			false,
		},
		{
			"not exists token",
			types.NewMsgUpdateReserve(addr, "not exists", []uint32{1, 2}, sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(2)))),
			true,
		},
		{
			"not owned token",
			types.NewMsgUpdateReserve(invalidOwner, ID, []uint32{1, 2}, sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(2)))),
			true,
		},
		{
			"new reserve has another denom",
			types.NewMsgUpdateReserve(addr, ID, []uint32{1, 2}, sdk.NewCoin("test", helpers.EtherToWei(sdkmath.NewInt(2)))),
			true,
		},
		{
			"not exists subtoken",
			types.NewMsgUpdateReserve(addr, ID, []uint32{133}, sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(2)))),
			true,
		},
		{
			"owned token, but now owned subtoken",
			types.NewMsgUpdateReserve(addr, ID, []uint32{11}, sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(2)))),
			true,
		},
		{
			"new reserve less than actual",
			types.NewMsgUpdateReserve(addr, ID, []uint32{8}, sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdkmath.NewInt(1)))),
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.UpdateReserve(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}
