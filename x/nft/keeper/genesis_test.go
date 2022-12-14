package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

func (s *KeeperTestSuite) TestInitGenesis() {
	dgs := types.DefaultGenesisState()
	require := s.Require()

	denom := "Test_Collection"
	pk := ed25519.GenPrivKey().PubKey()
	owner := sdk.AccAddress(pk.Address())
	tokenID := "token-1"
	reserve := sdk.NewCoin(cmdcfg.BaseDenom, sdk.NewInt(100000))

	dgs.Collections = []types.Collection{
		{
			Denom:   "Test_Collection",
			Creator: owner.String(),
			Supply:  1,
			Tokens: []*types.Token{
				{
					Creator:   owner.String(),
					Denom:     denom,
					ID:        tokenID,
					URI:       tokenID,
					Reserve:   reserve,
					AllowMint: true,
					Minted:    10,
					Burnt:     5,
					SubTokens: types.SubTokens{
						{
							ID:      1,
							Owner:   owner.String(),
							Reserve: &reserve,
						},
						{
							ID:      2,
							Owner:   owner.String(),
							Reserve: &reserve,
						},
						{
							ID:      3,
							Owner:   owner.String(),
							Reserve: &reserve,
						},
						{
							ID:      4,
							Owner:   owner.String(),
							Reserve: &reserve,
						},
						{
							ID:      5,
							Owner:   owner.String(),
							Reserve: &reserve,
						},
					},
				},
			},
		},
	}

	keeper.InitGenesis(s.ctx, s.nftKeeper, dgs)

	// check store collection
	collection, found := s.nftKeeper.GetCollection(s.ctx, owner, denom)
	require.True(found)

	// check store tokens
	tokens := s.nftKeeper.GetTokens(s.ctx, owner, denom)
	if len(tokens) < 1 {
		s.T().Fatal("tokens less than needed")
	}

	// check store subtokens
	subtokens := s.nftKeeper.GetSubTokens(s.ctx, tokenID)
	if len(subtokens) != len(dgs.Collections[0].Tokens[0].SubTokens) {
		s.T().Fatal("the number of tokens in genesis and storage is different")
	}

	for i, subtoken := range subtokens {
		require.True(dgs.Collections[0].Tokens[0].SubTokens[i].Equal(subtoken))
	}

	dgs.Collections[0].Tokens[0].SubTokens = types.SubTokens{}

	for i, token := range tokens {
		require.True(dgs.Collections[0].Tokens[i].Equal(token))
	}

	dgs.Collections[0].Tokens = types.Tokens{}

	require.True(dgs.Collections[0].Equal(collection))

	state := keeper.ExportGenesis(s.ctx, s.nftKeeper)

	require.True(state.Params.Equal(dgs.Params))
	for i, v := range dgs.Collections {
		v.Equal(state.Collections[i])
	}
}
