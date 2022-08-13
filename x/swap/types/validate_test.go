package types

import (
	"testing"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestSwapRedeem(t *testing.T) {
	cfg := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(cfg)
	cmdcfg.SetBip44CoinType(cfg)

	const validRecipient = "dx1tlhpwr6t9nnq95xjet3ap2lc9zlxyw9dhr9y0z"
	zeroHash := Hash([32]byte{0})

	sender, err := sdk.AccAddressFromBech32("dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v")
	require.NoError(t, err)

	var testCases = []struct {
		tag               string
		sender            sdk.AccAddress
		recipient         string
		amount            sdk.Int
		tokenSymbol       string
		transactionNumber string
		fromChain         uint32
		destChain         uint32
		expectError       bool
	}{
		{
			tag:               "valid redeem",
			sender:            sender,
			recipient:         validRecipient,
			amount:            sdk.NewInt(1),
			tokenSymbol:       "del",
			transactionNumber: "123",
			fromChain:         1,
			destChain:         3,
			expectError:       false,
		},
		{
			tag:               "invalid recipient",
			sender:            sender,
			recipient:         validRecipient + "0",
			amount:            sdk.NewInt(1),
			tokenSymbol:       "del",
			transactionNumber: "123",
			fromChain:         1,
			destChain:         3,
			expectError:       true,
		},
		{
			tag:               "invalid chain 1",
			sender:            sender,
			recipient:         validRecipient,
			amount:            sdk.NewInt(1),
			tokenSymbol:       "del",
			transactionNumber: "123",
			fromChain:         0,
			destChain:         3,
			expectError:       true,
		},
		{
			tag:               "invalid chain 2",
			sender:            sender,
			recipient:         validRecipient,
			amount:            sdk.NewInt(1),
			tokenSymbol:       "del",
			transactionNumber: "123",
			fromChain:         1,
			destChain:         0,
			expectError:       true,
		},
		{
			tag:               "invalid chain 3",
			sender:            sender,
			recipient:         validRecipient,
			amount:            sdk.NewInt(1),
			tokenSymbol:       "del",
			transactionNumber: "123",
			fromChain:         1,
			destChain:         1,
			expectError:       true,
		},
		{
			tag:               "invalid transaction number",
			sender:            sender,
			recipient:         validRecipient,
			amount:            sdk.NewInt(1),
			tokenSymbol:       "del",
			transactionNumber: "a123",
			fromChain:         1,
			destChain:         3,
			expectError:       true,
		},
		{
			tag:               "invalid amount",
			sender:            sender,
			recipient:         validRecipient,
			amount:            sdk.NewInt(-1),
			tokenSymbol:       "del",
			transactionNumber: "123",
			fromChain:         1,
			destChain:         3,
			expectError:       true,
		},
		{
			tag:               "invalid amount2",
			sender:            sender,
			recipient:         validRecipient,
			amount:            sdk.NewInt(0),
			tokenSymbol:       "del",
			transactionNumber: "123",
			fromChain:         1,
			destChain:         3,
			expectError:       true,
		},
	}
	for _, tc := range testCases {
		msg := NewMsgSwapRedeem(
			tc.sender,
			tc.sender.String(),
			tc.recipient,
			tc.amount,
			tc.tokenSymbol,
			tc.transactionNumber,
			tc.fromChain,
			tc.destChain,
			0,
			&zeroHash,
			&zeroHash,
		)
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}

}

func TestChainActivate(t *testing.T) {
	cfg := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(cfg)
	cmdcfg.SetBip44CoinType(cfg)

	sender, err := sdk.AccAddressFromBech32(ChainActivatorAddress)
	require.NoError(t, err)
	invalidSender, err := sdk.AccAddressFromBech32("dx1p844kydt9eljvuef4nk52dm6lcgj5c42q4zmvd")
	require.NoError(t, err)

	var testCases = []struct {
		tag         string
		sender      sdk.AccAddress
		chainNumber uint32
		chainName   string
		expectError bool
	}{
		{
			"valid chain activate",
			sender,
			1,
			"some blockchain",
			false,
		},
		{
			"invalid sender",
			invalidSender,
			1,
			"some blockchain",
			true,
		},
		{
			"invalid chain number",
			sender,
			0,
			"some blockchain",
			true,
		},
		{
			"invalid chain name",
			sender,
			1,
			"",
			true,
		},
	}
	for _, tc := range testCases {
		msg := NewMsgChainActivate(
			tc.sender,
			tc.chainNumber,
			tc.chainName,
		)
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}
}

func TestChainDeactivate(t *testing.T) {
	cfg := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(cfg)
	cmdcfg.SetBip44CoinType(cfg)

	sender, err := sdk.AccAddressFromBech32(ChainActivatorAddress)
	require.NoError(t, err)
	invalidSender, err := sdk.AccAddressFromBech32("dx1p844kydt9eljvuef4nk52dm6lcgj5c42q4zmvd")
	require.NoError(t, err)

	var testCases = []struct {
		tag         string
		sender      sdk.AccAddress
		chainNumber uint32
		expectError bool
	}{
		{
			"valid chain activate",
			sender,
			1,
			false,
		},
		{
			"invalid sender",
			invalidSender,
			1,
			true,
		},
		{
			"invalid chain number",
			sender,
			0,
			true,
		},
	}
	for _, tc := range testCases {
		msg := NewMsgChainDeactivate(
			tc.sender,
			tc.chainNumber,
		)
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}
}
