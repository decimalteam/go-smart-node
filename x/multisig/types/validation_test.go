package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

func TestValidateWallet(t *testing.T) {
	const addrCount = 100
	addrs := generateAddresses(addrCount)

	overMaxOwners := make([]string, MaxOwnerCount+1)
	overMaxWeights := make([]uint32, MaxOwnerCount+1)
	for i := 0; i < MaxOwnerCount+1; i++ {
		overMaxOwners[i] = addrs[i+10].String()
		overMaxWeights[i] = 1
	}

	var testCases = []struct {
		tag         string
		sender      sdk.AccAddress
		owners      []string
		weights     []uint32
		threshold   uint32
		expectError bool
	}{
		{
			tag:         "valid wallet",
			sender:      addrs[0],
			owners:      []string{addrs[1].String(), addrs[2].String(), addrs[3].String()},
			weights:     []uint32{1, 1, 1},
			threshold:   2,
			expectError: false,
		},
		{
			tag:         "1 owner",
			sender:      addrs[0],
			owners:      []string{addrs[1].String()},
			weights:     []uint32{1},
			threshold:   2,
			expectError: true,
		},
		{
			tag:         "over max owners",
			sender:      addrs[0],
			owners:      overMaxOwners,
			weights:     overMaxWeights,
			threshold:   2,
			expectError: true,
		},
		{
			tag:         "double owner",
			sender:      addrs[0],
			owners:      []string{addrs[1].String(), addrs[1].String(), addrs[2].String()},
			weights:     []uint32{1, 1, 1},
			threshold:   2,
			expectError: true,
		},
		{
			tag:         "owner count != weight count",
			sender:      addrs[0],
			owners:      []string{addrs[1].String(), addrs[2].String()},
			weights:     []uint32{1},
			threshold:   2,
			expectError: true,
		},
		{
			tag:         "invalid weight 1",
			sender:      addrs[0],
			owners:      []string{addrs[1].String(), addrs[2].String()},
			weights:     []uint32{MinWeight - 1, 1}, // < 1
			threshold:   2,
			expectError: true,
		},
		{
			tag:         "invalid weight 2",
			sender:      addrs[0],
			owners:      []string{addrs[1].String(), addrs[2].String()},
			weights:     []uint32{1, MaxWeight + 1}, // > 1024
			threshold:   2,
			expectError: true,
		},
		{
			tag:         "threshold over sum of weights",
			sender:      addrs[0],
			owners:      []string{addrs[1].String(), addrs[2].String()},
			weights:     []uint32{1, 1},
			threshold:   3,
			expectError: true,
		},
		{
			tag:         "invalid owner address",
			sender:      addrs[0],
			owners:      []string{addrs[1].String() + "0", addrs[2].String()},
			weights:     []uint32{1, 2},
			threshold:   2,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		msg := NewMsgCreateWallet(tc.sender, tc.owners, tc.weights, tc.threshold)
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}
}

func TestValidateTransaction(t *testing.T) {
	const addrCount = 100
	addrs := generateAddresses(addrCount)

	wallet, err := NewWallet(
		[]string{addrs[0].String(), addrs[1].String(), addrs[2].String()},
		[]uint32{1, 1, 1},
		2,
		[]byte{1},
	)
	require.NoError(t, err)

	var testCases = []struct {
		tag         string
		sender      sdk.AccAddress
		wallet      string
		receiver    string
		coins       sdk.Coins
		expectError bool
	}{
		{
			tag:         "valid transaction",
			sender:      addrs[0],
			wallet:      wallet.Address,
			receiver:    addrs[1].String(),
			coins:       sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1)))),
			expectError: false,
		},
		{
			tag:         "invalid wallet",
			sender:      addrs[0],
			wallet:      wallet.Address + "0",
			receiver:    addrs[1].String(),
			coins:       sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1)))),
			expectError: true,
		},
		{
			tag:         "invalid receiver",
			sender:      addrs[0],
			wallet:      wallet.Address,
			receiver:    addrs[1].String() + "0",
			coins:       sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1)))),
			expectError: true,
		},
		{
			tag:         "invalid coins",
			sender:      addrs[0],
			wallet:      wallet.Address,
			receiver:    addrs[1].String(),
			coins:       sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, sdk.ZeroInt())), // it will be empty
			expectError: true,
		},
	}

	for _, tc := range testCases {
		msg := NewMsgCreateTransaction(tc.sender, tc.wallet, tc.receiver, tc.coins)
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}
}

func TestValidateSignTransaction(t *testing.T) {
	const addrCount = 100
	addrs := generateAddresses(addrCount)

	wallet, err := NewWallet(
		[]string{addrs[0].String(), addrs[1].String(), addrs[2].String()},
		[]uint32{1, 1, 1},
		2,
		[]byte{1},
	)
	require.NoError(t, err)

	tx, err := NewTransaction(wallet.Address, addrs[1].String(),
		sdk.NewCoins(sdk.NewCoin(cmdcfg.BaseDenom, helpers.EtherToWei(sdk.NewInt(1)))),
		3, 100, []byte{1})
	require.NoError(t, err)

	var testCases = []struct {
		tag         string
		sender      sdk.AccAddress
		txID        string
		expectError bool
	}{
		{
			tag:         "valid sign",
			sender:      addrs[0],
			txID:        tx.Id,
			expectError: false,
		},
		{
			tag:         "invalid tx id",
			sender:      addrs[0],
			txID:        tx.Id + "0",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		msg := NewMsgSignTransaction(tc.sender, tc.txID)
		err := msg.ValidateBasic()
		if tc.expectError {
			require.Error(t, err, tc.tag)
		} else {
			require.NoError(t, err, tc.tag)
		}
	}
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(numAddrs int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, numAddrs)
	for i := 0; i < numAddrs; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}
	return testAddrs
}
