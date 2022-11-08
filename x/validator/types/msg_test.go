package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

var (
	coinPos  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)
	coinZero = sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)
)

// test ValidateBasic for MsgCreateValidator
func TestMsgCreateValidator(t *testing.T) {
	commission1 := sdk.ZeroDec()
	commission2 := sdk.MustNewDecFromStr("0.1")

	tests := []struct {
		name, moniker, identity, website, securityContact, details string
		commission                                                 sdk.Dec
		validatorAddr                                              sdk.ValAddress
		rewardAddr                                                 sdk.AccAddress
		pubkey                                                     cryptotypes.PubKey
		bond                                                       sdk.Coin
		expectPass                                                 bool
	}{
		{"basic good", "a", "b", "c", "d", "e", commission1, valAddr1, rewardAddr1, pk1, coinPos, true},
		{"partial description", "", "", "c", "", "", commission1, valAddr1, rewardAddr1, pk1, coinPos, true},
		{"empty description", "", "", "", "", "", commission2, valAddr1, rewardAddr1, pk1, coinPos, false},
		{"empty address", "a", "b", "c", "d", "e", commission2, emptyAddr, rewardAddr1, pk1, coinPos, false},
		{"empty pubkey", "a", "b", "c", "d", "e", commission1, valAddr1, rewardAddr1, emptyPubkey, coinPos, false},
		{"empty pubkey2", "a", "b", "c", "d", "e", commission1, valAddr1, rewardAddr1, &pkNil, coinPos, false},
		{"empty bond", "a", "b", "c", "d", "e", commission2, valAddr1, rewardAddr1, pk1, coinZero, false},
		{"nil bond", "a", "b", "c", "d", "e", commission2, valAddr1, rewardAddr1, pk1, sdk.Coin{}, false},
		{"negative commission", "a", "b", "c", "d", "e", sdk.MustNewDecFromStr("-0.1"), valAddr1, rewardAddr1, pk1, coinPos, false},
	}

	for _, tc := range tests {
		description := types.NewDescription(tc.moniker, tc.identity, tc.website, tc.securityContact, tc.details)
		msg, err := types.NewMsgCreateValidator(tc.validatorAddr, tc.rewardAddr, tc.pubkey, description,
			tc.commission, tc.bond)
		require.NoError(t, err)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgEditValidator
func TestMsgEditValidator(t *testing.T) {
	tests := []struct {
		name, moniker, identity, website, securityContact, details string
		validatorAddr                                              sdk.ValAddress
		rewardAddr                                                 sdk.AccAddress
		expectPass                                                 bool
	}{
		{"basic good", "a", "b", "c", "d", "e", valAddr1, rewardAddr1, true},
		{"partial description", "", "", "c", "", "", valAddr1, rewardAddr1, true},
		{"empty description", "", "", "", "", "", valAddr1, rewardAddr1, false},
		{"empty address", "a", "b", "c", "d", "e", emptyAddr, rewardAddr1, false},
	}

	for _, tc := range tests {
		description := types.NewDescription(tc.moniker, tc.identity, tc.website, tc.securityContact, tc.details)

		msg := types.NewMsgEditValidator(tc.validatorAddr, tc.rewardAddr, description)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgDelegate
func TestMsgDelegate(t *testing.T) {
	tests := []struct {
		name          string
		delegatorAddr sdk.AccAddress
		validatorAddr sdk.ValAddress
		bond          sdk.Coin
		expectPass    bool
	}{
		{"basic good", sdk.AccAddress(valAddr1), valAddr2, coinPos, true},
		{"self bond", sdk.AccAddress(valAddr1), valAddr1, coinPos, true},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, coinPos, false},
		{"empty validator", sdk.AccAddress(valAddr1), emptyAddr, coinPos, false},
		{"empty bond", sdk.AccAddress(valAddr1), valAddr2, coinZero, false},
		{"nil bold", sdk.AccAddress(valAddr1), valAddr2, sdk.Coin{}, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgDelegate(tc.delegatorAddr, tc.validatorAddr, tc.bond)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgDelegateNFT
func TestMsgDelegateNFT(t *testing.T) {
	tests := []struct {
		name          string
		delegatorAddr sdk.AccAddress
		validatorAddr sdk.ValAddress
		tokenID       string
		subTokens     []uint32
		expectPass    bool
	}{
		{"basic good", sdk.AccAddress(valAddr1), valAddr2, "abcdef", []uint32{1}, true},
		{"self bond", sdk.AccAddress(valAddr1), valAddr1, "abcdef", []uint32{1}, true},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, "abcdef", []uint32{1}, false},
		{"empty validator", sdk.AccAddress(valAddr1), emptyAddr, "abcdef", []uint32{1}, false},
		{"empty token id", sdk.AccAddress(valAddr1), valAddr2, "", []uint32{1}, false},
		{"empty subtokens", sdk.AccAddress(valAddr1), valAddr2, "abcdef", []uint32{}, false},
		{"non unique subtokens", sdk.AccAddress(valAddr1), valAddr2, "abcdef", []uint32{1, 2, 1}, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgDelegateNFT(tc.delegatorAddr, tc.validatorAddr, tc.tokenID, tc.subTokens)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgRedelegate
func TestMsgRedelegate(t *testing.T) {
	tests := []struct {
		name             string
		delegatorAddr    sdk.AccAddress
		validatorSrcAddr sdk.ValAddress
		validatorDstAddr sdk.ValAddress
		amount           sdk.Coin
		expectPass       bool
	}{
		{"regular", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), true},
		{"zero amount", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), false},
		{"nil amount", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.Coin{}, false},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
		{"empty source validator", sdk.AccAddress(valAddr1), emptyAddr, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
		{"empty destination validator", sdk.AccAddress(valAddr1), valAddr2, emptyAddr, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
		{"same validator", sdk.AccAddress(valAddr1), valAddr2, valAddr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
	}

	for _, tc := range tests {
		msg := types.NewMsgRedelegate(tc.delegatorAddr, tc.validatorSrcAddr, tc.validatorDstAddr, tc.amount)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

func TestMsgRedelegateNFT(t *testing.T) {
	tests := []struct {
		name             string
		delegatorAddr    sdk.AccAddress
		validatorSrcAddr sdk.ValAddress
		validatorDstAddr sdk.ValAddress
		tokenID          string
		subTokens        []uint32
		expectPass       bool
	}{
		{"regular", sdk.AccAddress(valAddr1), valAddr2, valAddr3, "abcdef", []uint32{1}, true},
		{"empty token id", sdk.AccAddress(valAddr1), valAddr2, valAddr3, "", []uint32{1}, false},
		{"empty subtoken", sdk.AccAddress(valAddr1), valAddr2, valAddr3, "abcdef", []uint32{}, false},
		{"dublicate subtoken", sdk.AccAddress(valAddr1), valAddr2, valAddr3, "abcdef", []uint32{1, 2, 1}, false},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, valAddr3, "abcdef", []uint32{1}, false},
		{"empty source validator", sdk.AccAddress(valAddr1), emptyAddr, valAddr3, "abcdef", []uint32{1}, false},
		{"empty destination validator", sdk.AccAddress(valAddr1), valAddr2, emptyAddr, "abcdef", []uint32{1}, false},
		{"same validator", sdk.AccAddress(valAddr1), valAddr2, valAddr2, "abcdef", []uint32{1}, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgRedelegateNFT(tc.delegatorAddr, tc.validatorSrcAddr, tc.validatorDstAddr, tc.tokenID, tc.subTokens)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgUndelegate
func TestMsgUndelegate(t *testing.T) {
	tests := []struct {
		name          string
		delegatorAddr sdk.AccAddress
		validatorAddr sdk.ValAddress
		amount        sdk.Coin
		expectPass    bool
	}{
		{"regular", sdk.AccAddress(valAddr1), valAddr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), true},
		{"zero amount", sdk.AccAddress(valAddr1), valAddr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), false},
		{"nil amount", sdk.AccAddress(valAddr1), valAddr2, sdk.Coin{}, false},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
		{"empty validator", sdk.AccAddress(valAddr1), emptyAddr, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
	}

	for _, tc := range tests {
		msg := types.NewMsgUndelegate(tc.delegatorAddr, tc.validatorAddr, tc.amount)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgUndelegateNFT
func TestMsgUndelegateNFT(t *testing.T) {
	tests := []struct {
		name          string
		delegatorAddr sdk.AccAddress
		validatorAddr sdk.ValAddress
		tokenID       string
		subTokens     []uint32
		expectPass    bool
	}{
		{"regular", sdk.AccAddress(valAddr1), valAddr2, "abcdef", []uint32{1, 2}, true},
		{"empty token id", sdk.AccAddress(valAddr1), valAddr2, "", []uint32{1, 2}, false},
		{"empty subtokens", sdk.AccAddress(valAddr1), valAddr2, "abcdef", []uint32{}, false},
		{"repeated subtokens", sdk.AccAddress(valAddr1), valAddr2, "abcdef", []uint32{1, 2, 1}, false},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, "abcdef", []uint32{1, 2}, false},
		{"empty validator", sdk.AccAddress(valAddr1), emptyAddr, "abcdef", []uint32{1, 2}, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgUndelegateNFT(tc.delegatorAddr, tc.validatorAddr, tc.tokenID, tc.subTokens)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgCancelRedelegation
func TestMsgCancelRedelegate(t *testing.T) {
	tests := []struct {
		name             string
		delegatorAddr    sdk.AccAddress
		validatorAddrSrc sdk.ValAddress
		validatorAddrDst sdk.ValAddress
		amount           sdk.Coin
		height           int64
		expectPass       bool
	}{
		{"regular", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), 1, true},
		{"zero amount", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), 1, false},
		{"nil amount", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.Coin{}, 1, false},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), 1, false},
		{"empty validator src", sdk.AccAddress(valAddr1), emptyAddr, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), 1, false},
		{"empty validator dst", sdk.AccAddress(valAddr1), valAddr1, emptyAddr, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), 1, false},
		{"zero height", sdk.AccAddress(valAddr1), valAddr2, valAddr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), 0, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgCancelRedelegation(tc.delegatorAddr, tc.validatorAddrSrc, tc.validatorAddrDst, tc.height, tc.amount)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgCancelRedelegation
func TestMsgCancelRedelegateNFT(t *testing.T) {
	tests := []struct {
		name             string
		delegatorAddr    sdk.AccAddress
		validatorAddrSrc sdk.ValAddress
		validatorAddrDst sdk.ValAddress
		tokenID          string
		subTokens        []uint32
		height           int64
		expectPass       bool
	}{
		{"regular", sdk.AccAddress(valAddr1), valAddr2, valAddr3, "abcdef", []uint32{1}, 1, true},
		{"empty token id", sdk.AccAddress(valAddr1), valAddr2, valAddr3, "", []uint32{1}, 1, false},
		{"empty subtokens", sdk.AccAddress(valAddr1), valAddr2, valAddr3, "abcdef", []uint32{}, 1, false},
		{"repeated subtokens", sdk.AccAddress(valAddr1), valAddr2, valAddr3, "abcdef", []uint32{1, 2, 1}, 1, false},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, valAddr3, "abcdef", []uint32{1}, 1, false},
		{"empty validator src", sdk.AccAddress(valAddr1), emptyAddr, valAddr3, "abcdef", []uint32{1}, 1, false},
		{"empty validator dst", sdk.AccAddress(valAddr1), valAddr1, emptyAddr, "abcdef", []uint32{1}, 1, false},
		{"zero height", sdk.AccAddress(valAddr1), valAddr2, valAddr3, "abcdef", []uint32{1}, 0, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgCancelRedelegationNFT(tc.delegatorAddr, tc.validatorAddrSrc, tc.validatorAddrDst, tc.height, tc.tokenID, tc.subTokens)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgCancelUndelegation
func TestMsgCancelUndelegate(t *testing.T) {
	tests := []struct {
		name          string
		delegatorAddr sdk.AccAddress
		validatorAddr sdk.ValAddress
		amount        sdk.Coin
		height        int64
		expectPass    bool
	}{
		{"regular", sdk.AccAddress(valAddr1), valAddr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), 1, true},
		{"zero amount", sdk.AccAddress(valAddr1), valAddr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), 1, false},
		{"nil amount", sdk.AccAddress(valAddr1), valAddr2, sdk.Coin{}, 1, false},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), 1, false},
		{"empty validator", sdk.AccAddress(valAddr1), emptyAddr, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), 1, false},
		{"zero height", sdk.AccAddress(valAddr1), valAddr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), 0, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgCancelUndelegation(tc.delegatorAddr, tc.validatorAddr, tc.height, tc.amount)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for MsgCancelUndelegationNFT
func TestMsgCancelUndelegateNFT(t *testing.T) {
	tests := []struct {
		name          string
		delegatorAddr sdk.AccAddress
		validatorAddr sdk.ValAddress
		tokenID       string
		subTokens     []uint32
		height        int64
		expectPass    bool
	}{
		{"regular", sdk.AccAddress(valAddr1), valAddr2, "abcdef", []uint32{1}, 1, true},
		{"empty token id", sdk.AccAddress(valAddr1), valAddr2, "", []uint32{1}, 1, false},
		{"empty subtokens", sdk.AccAddress(valAddr1), valAddr2, "abcdef", []uint32{}, 1, false},
		{"repeated subtokens", sdk.AccAddress(valAddr1), valAddr2, "abcdef", []uint32{1, 2, 1}, 1, false},
		{"empty delegator", sdk.AccAddress(emptyAddr), valAddr1, "abcdef", []uint32{1}, 1, false},
		{"empty validator", sdk.AccAddress(valAddr1), emptyAddr, "abcdef", []uint32{1}, 1, false},
		{"zero height", sdk.AccAddress(valAddr1), valAddr2, "abcdef", []uint32{1}, 0, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgCancelUndelegationNFT(tc.delegatorAddr, tc.validatorAddr, tc.height, tc.tokenID, tc.subTokens)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}
