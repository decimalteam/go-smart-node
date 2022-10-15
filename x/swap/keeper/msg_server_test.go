package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/x/swap/types"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (s *KeeperTestSuite) TestMsgActivateChain() {
	ctx, _, msgServer := s.ctx, s.swapKeeper, s.msgServer
	require := s.Require()

	testCases := []struct {
		name   string
		input  *types.MsgActivateChain
		expErr bool
	}{
		{
			"valid request",
			types.NewMsgActivateChain(user1, defaultChainID, defaultChainName),
			false,
		},
		{
			"sender is not service",
			types.NewMsgActivateChain(user2, defaultChainID, defaultChainName),
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.ActivateChain(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgDeactivateChain() {
	ctx, k, msgServer := s.ctx, s.swapKeeper, s.msgServer
	require := s.Require()

	k.SetChain(ctx, &types.Chain{
		Id:     defaultChainID,
		Name:   defaultChainName,
		Active: true,
	})

	testCases := []struct {
		name   string
		input  *types.MsgDeactivateChain
		expErr bool
	}{
		{
			"valid request",
			types.NewMsgDeactivateChain(user1, defaultChainID),
			false,
		},
		{
			"sender is not a service",
			types.NewMsgDeactivateChain(user2, defaultChainID),
			true,
		},
		{
			"chain with given ID does not exists",
			types.NewMsgDeactivateChain(user1, 10),
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.DeactivateChain(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgInitializeSwap() {
	ctx, k, msgServer := s.ctx, s.swapKeeper, s.msgServer
	require := s.Require()

	k.SetChain(ctx, &types.Chain{
		Id:     defaultChainID,
		Name:   defaultChainName,
		Active: true,
	})

	k.SetChain(ctx, &types.Chain{
		Id:     destChainID,
		Name:   destChainName,
		Active: true,
	})

	testCases := []struct {
		name   string
		input  *types.MsgInitializeSwap
		expErr bool
	}{
		{
			"valid request",
			types.NewMsgInitializeSwap(user2, user3.String(), defaultAmount, defaultTokenSymbol, "0", defaultChainID, destChainID),
			false,
		},
		{
			"from chain does not exist",
			types.NewMsgInitializeSwap(user2, user3.String(), defaultAmount, defaultTokenSymbol, "0", 10, destChainID),
			true,
		},
		{
			"to chain does not exist",
			types.NewMsgInitializeSwap(user2, user3.String(), defaultAmount, defaultTokenSymbol, "0", defaultChainID, 10),
			true,
		},
		{
			"insuffucient funds",
			types.NewMsgInitializeSwap(user3, user2.String(), defaultAmount, defaultTokenSymbol, "0", defaultChainID, destChainID),
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.InitializeSwap(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgRedeem() {
	ctx, k, msgServer := s.ctx, s.swapKeeper, s.msgServer
	require := s.Require()

	swap := types.MsgRedeemSwap{
		Sender:            user2.String(),
		From:              user2.String(),
		Recipient:         user3.String(),
		Amount:            defaultAmount,
		TokenSymbol:       defaultTokenSymbol,
		TransactionNumber: "1",
		FromChain:         defaultChainID,
		DestChain:         destChainID,
		V:                 0,
		R:                 "",
		S:                 "",
	}

	transactionNumber, _ := sdk.NewIntFromString(swap.TransactionNumber)

	// get message hash
	hash, err := types.GetHash(transactionNumber, swap.TokenSymbol, swap.Amount, swap.Recipient, swap.FromChain, swap.DestChain)
	require.NoError(err)

	// sign message hash
	sig, err := privkey.Sign(hash[:])
	require.NoError(err)

	// get V R S
	v, r, sH := decodeSignature(sig)

	var _r types.Hash
	copy(_r[:], r.BigInt().Bytes())

	var _s types.Hash
	copy(_s[:], sH.BigInt().Bytes())

	swap.V = uint32(v.Uint64())
	swap.R = hex.EncodeToString(_r[:])
	swap.S = hex.EncodeToString(_s[:])

	k.SetChain(ctx, &types.Chain{
		Id:     defaultChainID,
		Name:   defaultChainName,
		Active: true,
	})

	k.SetChain(ctx, &types.Chain{
		Id:     destChainID,
		Name:   destChainName,
		Active: true,
	})

	testCases := []struct {
		name   string
		input  *types.MsgRedeemSwap
		expErr bool
	}{
		{
			"valid request",
			types.NewMsgRedeemSwap(user2, user2.String(), user3.String(), defaultAmount, defaultTokenSymbol, "1",
				defaultChainID, destChainID, swap.V, hex.EncodeToString(_r[:]), hex.EncodeToString(_s[:])),
			false,
		},
		{
			"swap is already redeemed",
			types.NewMsgRedeemSwap(user2, user2.String(), user3.String(), defaultAmount, defaultTokenSymbol, "1",
				defaultChainID, destChainID, swap.V, hex.EncodeToString(_r[:]), hex.EncodeToString(_s[:])),
			true,
		},
		{
			"invalid transaction number",
			types.NewMsgRedeemSwap(user2, user2.String(), user3.String(), defaultAmount, defaultTokenSymbol, "asf",
				defaultChainID, destChainID, 27, hex.EncodeToString(_r[:]), hex.EncodeToString(_s[:])),
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.RedeemSwap(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}

func decodeSignature(sig []byte) (v, r, s sdkmath.Int) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = sdk.NewIntFromBigInt(new(big.Int).SetBytes(sig[:32]))
	s = sdk.NewIntFromBigInt(new(big.Int).SetBytes(sig[32:64]))
	v = sdk.NewIntFromBigInt(new(big.Int).SetBytes([]byte{sig[64] + 27}))
	return
}
