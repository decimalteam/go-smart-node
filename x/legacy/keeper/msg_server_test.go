package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/decimalteam/ethermint/crypto/ethsecp256k1"
	feemarkettypes "github.com/decimalteam/ethermint/x/feemarket/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	commonTypes "bitbucket.org/decimalteam/go-smart-node/types"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy"
	"bitbucket.org/decimalteam/go-smart-node/x/legacy/types"
	multisigtypes "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

var (
	invalidAccPk         = secp256k1.GenPrivKey().PubKey()
	invalidOldAddress, _ = commonTypes.GetLegacyAddressFromPubKey(invalidAccPk.Bytes())
	invalidActualAddress = sdk.AccAddress(ethsecp256k1.PubKey{Key: invalidAccPk.Bytes()}.Address())
	defaultTokenID       = "test_record_token_ID"
	defaultRecord        = types.Record{
		LegacyAddress: oldAddress,
		Coins:         sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(100)), sdk.NewCoin("test2", sdk.NewInt(102))),
		Wallets:       []string{defaultMultisigWalletBefore.Address},
		NFTs:          []string{defaultTokenID},
		Validators:    []string{defaultOperatorAddressSdk.String()},
	}
	subTokenReserve        = sdk.NewCoin(baseDenom, sdk.NewInt(10))
	defaultSubTokensBefore = []nfttypes.SubToken{
		{
			ID:      1,
			Owner:   oldAddress,
			Reserve: &subTokenReserve,
		},
	}
	defaultSubTokensAfter = []nfttypes.SubToken{
		{
			ID:      1,
			Owner:   actualAddress,
			Reserve: &subTokenReserve,
		},
	}
	defaultMultisigWalletBefore = multisigtypes.Wallet{
		Address: "wallet",
		Owners:  []string{oldAddress, invalidActualAddress.String()},
	}
	defaultMultisigWalletAfter = multisigtypes.Wallet{
		Address: "wallet",
		Owners:  []string{actualAddress, invalidActualAddress.String()},
	}

	defaultOperatorAddressSdk = sdk.ValAddress([]byte{1, 2, 3})
	defaultOperatorAddress    = "d0valoper1aaa"

	defaultValidatorBefore = validatortypes.Validator{
		OperatorAddress: defaultOperatorAddress,
		RewardAddress:   oldAddress,
	}
	defaultValidatorAfter = validatortypes.Validator{
		OperatorAddress: defaultOperatorAddress,
		RewardAddress:   actualAddress,
	}
	defaultRewards = sdk.NewCoins(sdk.NewCoin("test1", sdk.NewInt(100)))
)

func (s *KeeperTestSuite) TestMsgReturnLegacy() {
	ctx, k, msgServer := s.ctx, s.legacyKeeper, s.msgServer
	require := s.Require()

	k.SetLegacyRecord(ctx, defaultRecord)

	testCases := []struct {
		name   string
		input  *types.MsgReturnLegacy
		expErr bool
	}{
		{
			"valid request",
			types.NewMsgReturnLegacy(newAddress, publicKey),
			false,
		},
		{
			"not have record",
			types.NewMsgReturnLegacy(invalidActualAddress, invalidAccPk.Bytes()),
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.ReturnLegacy(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}

func TestLegacyReturnForValidator(t *testing.T) {
	const validatorAddress = "d0valoper14elhyzmq95f98wrkvujtsr5cyudffp6qwyerml"
	var publickey = []byte{0x3, 0x44, 0x8e, 0x6b, 0x3d, 0x50, 0xd6, 0xa3, 0x9c, 0xab, 0x3b, 0xab, 0xaa, 0x4a, 0xa2, 0xb0, 0x88, 0x5f, 0x55, 0x6f, 0xe0, 0x5d, 0x71, 0x49, 0x88, 0x5a, 0x5, 0xa0, 0xe7, 0x94, 0xa, 0x7e, 0x4f}
	const legacyRewardAddress = "dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy"
	const actualRewardAddress = "d01xp6aqad49te7vsfga6str8hrdeh24r9jhplgxv"

	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	dsc.ValidatorKeeper.SetValidator(ctx, validatortypes.Validator{
		OperatorAddress: validatorAddress,
		RewardAddress:   legacyRewardAddress,
	})

	genesisState := types.GenesisState{
		Records: []types.Record{
			{
				LegacyAddress: legacyRewardAddress,
				Validators:    []string{validatorAddress},
			},
		},
	}
	legacy.InitGenesis(ctx, dsc.LegacyKeeper, &genesisState)
	dsc.LegacyKeeper.ActualizeLegacy(ctx, publickey)
}
