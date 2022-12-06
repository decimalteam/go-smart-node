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

	const validRecipient = "d01tlhpwr6t9nnq95xjet3ap2lc9zlxyw9dnyx3ya"

	sender, err := sdk.AccAddressFromBech32("d01lx4lvt8sjuxj8vw5dcf6knnq0pacre4w7swzpn")
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
		r                 string
		s                 string
		expectError       bool
	}{
		{
			tag:               "valid redeem",
			sender:            sender,
			recipient:         validRecipient,
			amount:            sdk.NewInt(1),
			tokenSymbol:       cmdcfg.BaseDenom,
			transactionNumber: "123",
			fromChain:         1,
			destChain:         3,
			r:                 "d8c0c8ff4a9b168be168f480bae61ead0a7f2b973f983a038f867621451fa553",
			s:                 "641ba9f5749afbb425e83b69ecacb3a0c6e32e2431609d474d4300a7cce5eb41",
			expectError:       false,
		},
		{
			tag:               "invalid recipient",
			sender:            sender,
			recipient:         validRecipient + "0",
			amount:            sdk.NewInt(1),
			tokenSymbol:       cmdcfg.BaseDenom,
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
			tokenSymbol:       cmdcfg.BaseDenom,
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
			tokenSymbol:       cmdcfg.BaseDenom,
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
			tokenSymbol:       cmdcfg.BaseDenom,
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
			tokenSymbol:       cmdcfg.BaseDenom,
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
			tokenSymbol:       cmdcfg.BaseDenom,
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
			tokenSymbol:       cmdcfg.BaseDenom,
			transactionNumber: "123",
			fromChain:         1,
			destChain:         3,
			expectError:       true,
		},
		{
			tag:               "invalid R 1",
			sender:            sender,
			recipient:         validRecipient,
			amount:            sdk.NewInt(1),
			tokenSymbol:       cmdcfg.BaseDenom,
			transactionNumber: "123",
			fromChain:         1,
			destChain:         3,
			r:                 "d8c0c8ff4a9b168be168f480bae61ead0a7f2b973f983a038f867621451fa55",
			s:                 "641ba9f5749afbb425e83b69ecacb3a0c6e32e2431609d474d4300a7cce5eb41",
			expectError:       true,
		},
		{
			tag:               "invalid R 2",
			sender:            sender,
			recipient:         validRecipient,
			amount:            sdk.NewInt(1),
			tokenSymbol:       cmdcfg.BaseDenom,
			transactionNumber: "123",
			fromChain:         1,
			destChain:         3,
			r:                 "d8c0c8ff4a9b168be168f480bae61ead0a7f2b973f983a038f867621451fa55z",
			s:                 "641ba9f5749afbb425e83b69ecacb3a0c6e32e2431609d474d4300a7cce5eb41",
			expectError:       true,
		},
		{
			tag:               "invalid S 1",
			sender:            sender,
			recipient:         validRecipient,
			amount:            sdk.NewInt(1),
			tokenSymbol:       cmdcfg.BaseDenom,
			transactionNumber: "123",
			fromChain:         1,
			destChain:         3,
			r:                 "d8c0c8ff4a9b168be168f480bae61ead0a7f2b973f983a038f867621451fa553",
			s:                 "641ba9f5749afbb425e83b69ecacb3a0c6e32e2431609d474d4300a7cce5eb4",
			expectError:       true,
		},
		{
			tag:               "invalid S 2",
			sender:            sender,
			recipient:         validRecipient,
			amount:            sdk.NewInt(1),
			tokenSymbol:       cmdcfg.BaseDenom,
			transactionNumber: "123",
			fromChain:         1,
			destChain:         3,
			r:                 "d8c0c8ff4a9b168be168f480bae61ead0a7f2b973f983a038f867621451fa553",
			s:                 "641ba9f5749afbb425e83b69ecacb3a0c6e32e2431609d474d4300a7cce5eb4z",
			expectError:       true,
		},
	}
	for _, tc := range testCases {
		msg := NewMsgRedeemSwap(
			tc.sender,
			tc.sender.String(),
			tc.recipient,
			tc.amount,
			tc.tokenSymbol,
			tc.transactionNumber,
			tc.fromChain,
			tc.destChain,
			0,
			tc.r,
			tc.s,
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

	sender, err := sdk.AccAddressFromBech32(DefaultSwapServiceAddress)
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
		msg := NewMsgActivateChain(
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

	sender, err := sdk.AccAddressFromBech32(DefaultSwapServiceAddress)
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
			"invalid chain number",
			sender,
			0,
			true,
		},
	}
	for _, tc := range testCases {
		msg := NewMsgDeactivateChain(
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
