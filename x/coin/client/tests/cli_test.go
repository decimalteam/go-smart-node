//go:build norace
// +build norace

package tests

import (
	"bitbucket.org/decimalteam/go-smart-node/testutil/network"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/client/cli"
	"fmt"
	"github.com/cosmos/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	etherminthd "github.com/evmos/ethermint/crypto/hd"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"io/ioutil"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	if testing.Short() {
		s.T().Skip("skipping test in unit-tests mode.")
	}

	cfg := network.DefaultConfig()
	cfg.NumValidators = 1

	s.cfg = cfg

	baseDir, err := ioutil.TempDir(s.T().TempDir(), s.cfg.ChainID)
	require.NoError(s.T(), err)
	s.T().Logf("created temporary directory: %s", baseDir)

	s.network, err = network.New(s.T(), baseDir, s.cfg)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestCreateCoinCmd() {
	require := s.Require()
	val := s.network.Validators[0]

	info, _, err := val.ClientCtx.Keyring.NewMnemonic("createCoin", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, etherminthd.EthSecp256k1)
	require.NoError(err)

	pk, err := info.GetPubKey()
	require.NoError(err)

	newAddr := sdk.AccAddress(pk.Address())
	_, err = MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, helpers.EtherToWei(sdk.NewInt(100000)))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, helpers.EtherToWei(sdk.NewInt(1)))).String()),
	)
	require.NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "create coin valid command",
			args: []string{
				"valid coin title",          //title
				"TEST",                      //denom
				"80",                        //crr
				"1000000000000000000000000", //initReserve
				"9990000000",                // initVolume
				"91000000000000000000000",   //limitVolume
				"",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: false,
		},
		{
			name: "invalid crr",
			args: []string{
				"valid coin title",          //title
				"TEST",                      //denom
				"800",                       //crr
				"1000000000000000000000000", //initReserve
				"9990000000",                // initVolume
				"91000000000000000000000",   //limitVolume
				"",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "coin exist",
			args: []string{
				"valid coin title",          //title
				"TEST",                      //denom
				"800",                       //crr
				"1000000000000000000000000", //initReserve
				"9990000000",                // initVolume
				"91000000000000000000000",   //limitVolume
				"",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewCreateCoinCmd()
			clientCtx := val.ClientCtx
			args := addBasicFlagsForTxCmd(tc.args, s.cfg.BondDenom)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestUpdateCoinCmd() {
	require := s.Require()
	val := s.network.Validators[0]

	info, _, err := val.ClientCtx.Keyring.NewMnemonic("updateCoin", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, etherminthd.EthSecp256k1)
	require.NoError(err)

	pk, err := info.GetPubKey()
	require.NoError(err)

	newAddr := sdk.AccAddress(pk.Address())
	_, err = MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, helpers.EtherToWei(sdk.NewInt(100000)))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, helpers.EtherToWei(sdk.NewInt(1)))).String()),
	)
	require.NoError(err)

	s.T().Log(val.Address.String())
	symbol := "CUSTUPDATE"
	_, err = MsgCreateCoinExec(
		val.ClientCtx,
		"Custom coin for update",
		symbol,
		"50",
		helpers.EtherToWei(sdk.NewInt(2000)).String(),
		helpers.EtherToWei(sdk.NewInt(50)).String(),
		helpers.EtherToWei(sdk.NewInt(100)).String(),
		"",
		val.Address.String(),
		addBasicFlagsForTxCmd([]string{}, s.cfg.BondDenom)...,
	)
	require.NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "valid update coin request",
			args: []string{
				symbol, // custom coin symbol
				helpers.EtherToWei(sdk.NewInt(101)).String(), // limit volume
				"", // identity
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: false,
		},
		{
			name: "not exist coin symbol request",
			args: []string{
				"invalid symbol", // custom coin symbol
				helpers.EtherToWei(sdk.NewInt(100)).String(), // limit volume
				"", // identity
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "invalid limit volume request",
			args: []string{
				symbol,                 // custom coin symbol
				"invalid limit volume", // limit volume
				"",                     // identity
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "author and from address not equal request",
			args: []string{
				symbol, // custom coin symbol
				helpers.EtherToWei(sdk.NewInt(100)).String(), // limit volume
				"", // identity
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewUpdateCoinCmd()
			clientCtx := val.ClientCtx
			args := addBasicFlagsForTxCmd(tc.args, s.cfg.BondDenom)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestBuyCoinCmd() {
	require := s.Require()
	val := s.network.Validators[0]

	info, _, err := val.ClientCtx.Keyring.NewMnemonic("newAcc", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, etherminthd.EthSecp256k1)
	require.NoError(err)

	pk, err := info.GetPubKey()
	require.NoError(err)

	newAddr := sdk.AccAddress(pk.Address())
	symbol := "buycoin"
	_, err = MsgCreateCoinExec(
		val.ClientCtx,
		"Custom coin for update",
		symbol,
		"50",
		helpers.EtherToWei(sdk.NewInt(2000)).String(),
		helpers.EtherToWei(sdk.NewInt(500000)).String(),
		helpers.EtherToWei(sdk.NewInt(10000000)).String(),
		"",
		val.Address.String(),
		addBasicFlagsForTxCmd([]string{}, s.cfg.BondDenom)...,
	)
	require.NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "valid buy coin request",
			args: []string{
				fmt.Sprintf("100%s", symbol),                 //amountToBuy
				fmt.Sprintf("1000000000%s", s.cfg.BondDenom), //amountToSell
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: false,
		},
		{
			name: "buy coin not exist request",
			args: []string{
				fmt.Sprintf("100%s", "notexistcoin"),         //amountToBuy
				fmt.Sprintf("1000000000%s", s.cfg.BondDenom), //amountToSell
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "sell coin not exist request",
			args: []string{
				fmt.Sprintf("100%s", symbol),                //amountToBuy
				fmt.Sprintf("1000000000%s", "notexistcoin"), //amountToSell
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "invalid buy coin request",
			args: []string{
				"invalidcoindenom",                           //amountToBuy
				fmt.Sprintf("1000000000%s", s.cfg.BondDenom), //amountToSell
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "invalid sell coin request",
			args: []string{
				fmt.Sprintf("100%s", symbol), //amountToBuy
				"invalidcoindenom",           //amountToSell
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "addr not have buy coins request",
			args: []string{
				fmt.Sprintf("100%s", symbol),                 //amountToBuy
				fmt.Sprintf("1000000000%s", s.cfg.BondDenom), //amountToSell
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewBuyCoinCmd()
			clientCtx := val.ClientCtx
			args := addBasicFlagsForTxCmd(tc.args, s.cfg.BondDenom)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestSellCoinCmd() {
	require := s.Require()
	val := s.network.Validators[0]

	info, _, err := val.ClientCtx.Keyring.NewMnemonic("sellCoin", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, etherminthd.EthSecp256k1)
	require.NoError(err)

	pk, err := info.GetPubKey()
	require.NoError(err)

	newAddr := sdk.AccAddress(pk.Address())

	symbol := "sellcoin"
	_, err = MsgCreateCoinExec(
		val.ClientCtx,
		"Custom coin for update",
		symbol,
		"50",
		helpers.EtherToWei(sdk.NewInt(2000)).String(),
		helpers.EtherToWei(sdk.NewInt(500000)).String(),
		helpers.EtherToWei(sdk.NewInt(10000000)).String(),
		"",
		val.Address.String(),
		addBasicFlagsForTxCmd([]string{}, s.cfg.BondDenom)...,
	)
	require.NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "valid sell coin request",
			args: []string{
				fmt.Sprintf("1000000000%s", s.cfg.BondDenom), //amountToSell
				fmt.Sprintf("10000%s", symbol),               //minAmountToBuy
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: false,
		},
		{
			name: "not exist sell coin request",
			args: []string{
				fmt.Sprintf("1000000000%s", "notexistcoin"), //amountToSell
				fmt.Sprintf("10000%s", symbol),              //minAmountToBuy
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "not exist buy coin request",
			args: []string{
				fmt.Sprintf("1000000000%s", s.cfg.BondDenom), //amountToSell
				fmt.Sprintf("10000%s", "not exist coin"),     //minAmountToBuy
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "invalid send coin request",
			args: []string{
				"invalidcoinamount",                      //amountToSell
				fmt.Sprintf("10000%s", "not exist coin"), //minAmountToBuy
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "invalid buy coin request",
			args: []string{
				fmt.Sprintf("1000000000%s", s.cfg.BondDenom), //amountToSell
				"invalidcoinamount",                          //minAmountToBuy
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "addr not have sell coins request",
			args: []string{
				fmt.Sprintf("100%s", symbol),                 //amountToSell
				fmt.Sprintf("1000000000%s", s.cfg.BondDenom), //minAmountToBuy
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewSellCoinCmd()
			clientCtx := val.ClientCtx
			args := addBasicFlagsForTxCmd(tc.args, s.cfg.BondDenom)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestSendCoinCmd() {
	require := s.Require()
	val := s.network.Validators[0]

	info, _, err := val.ClientCtx.Keyring.NewMnemonic("sendCoin", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, etherminthd.EthSecp256k1)
	require.NoError(err)

	pk, err := info.GetPubKey()
	require.NoError(err)
	newAddr := sdk.AccAddress(pk.Address())

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "valid send coin request",
			args: []string{
				newAddr.String(), // recipient
				fmt.Sprintf("1000000000%s", s.cfg.BondDenom), // amount
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: false,
		},
		{
			name: "not exist coin to send request",
			args: []string{
				newAddr.String(), // recipient
				fmt.Sprintf("1000000000%s", "notexistcoin"), // amount
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "not valid coin to send request",
			args: []string{
				newAddr.String(), // recipient
				"notvalidamount", // amount
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "not valid recipient bech32  request",
			args: []string{
				"notvalidbech32", // recipient
				fmt.Sprintf("1000000000%s", s.cfg.BondDenom), // amount
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "valid send coin request",
			args: []string{
				val.Address.String(),                            // recipient
				fmt.Sprintf("1000000000000%s", s.cfg.BondDenom), // amount
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr.String()),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewSendCoinCmd()
			clientCtx := val.ClientCtx
			args := addBasicFlagsForTxCmd(tc.args, s.cfg.BondDenom)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestMultisendCmd() {
	val := s.network.Validators[0]
	require := s.Require()

	info, _, err := val.ClientCtx.Keyring.NewMnemonic("multisend1", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, etherminthd.EthSecp256k1)
	require.NoError(err)
	info1, _, err := val.ClientCtx.Keyring.NewMnemonic("multisend2", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, etherminthd.EthSecp256k1)
	require.NoError(err)
	info2, _, err := val.ClientCtx.Keyring.NewMnemonic("multisend3", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, etherminthd.EthSecp256k1)
	require.NoError(err)

	pk, err := info.GetPubKey()
	require.NoError(err)

	pk1, err := info1.GetPubKey()
	require.NoError(err)

	pk2, err := info2.GetPubKey()
	require.NoError(err)

	newAddr := sdk.AccAddress(pk.Address())
	newAddr1 := sdk.AccAddress(pk1.Address())
	newAddr2 := sdk.AccAddress(pk2.Address())

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "valid multi-send coin request",
			args: []string{
				newAddr.String(),
				fmt.Sprintf("1000000000000000%s", s.cfg.BondDenom),
				newAddr1.String(),
				fmt.Sprintf("1000000000000000%s", s.cfg.BondDenom),
				newAddr2.String(),
				fmt.Sprintf("1000000%s", s.cfg.BondDenom),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: false,
		},
		{
			name: "args len is not even request",
			args: []string{
				newAddr.String(),
				fmt.Sprintf("1000000000000000%s", s.cfg.BondDenom),
				newAddr1.String(),
				fmt.Sprintf("1000000000000000%s", s.cfg.BondDenom),
				newAddr2.String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "invalid reciever bech32 address request",
			args: []string{
				"invalidBech32Addr",
				fmt.Sprintf("1000000000000000%s", s.cfg.BondDenom),
				newAddr1.String(),
				fmt.Sprintf("1000000000000000%s", s.cfg.BondDenom),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "multi-send coin not exist request",
			args: []string{
				newAddr.String(),
				fmt.Sprintf("1000000000000000%s", "notexistcoin"),
				newAddr1.String(),
				fmt.Sprintf("1000000000000000%s", s.cfg.BondDenom),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "multi-send coin is invalid request",
			args: []string{
				newAddr.String(),
				"invalidCoinAmount",
				newAddr1.String(),
				fmt.Sprintf("1000000000000000%s", s.cfg.BondDenom),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			},
			expectErr: true,
		},
		{
			name: "addr not have balance request",
			args: []string{
				newAddr.String(),
				fmt.Sprintf("1000000000000000%s", s.cfg.BondDenom),
				newAddr1.String(),
				fmt.Sprintf("1000000000000000%s", s.cfg.BondDenom),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr2.String()),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewMultiSendCoinCmd()
			clientCtx := val.ClientCtx
			args := addBasicFlagsForTxCmd(tc.args, s.cfg.BondDenom)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestSellAllCmd() {
	val := s.network.Validators[0]
	require := s.Require()

	symbol := "sellall"
	_, err := MsgCreateCoinExec(
		val.ClientCtx,
		"Custom coin for update",
		symbol,
		"50",
		helpers.EtherToWei(sdk.NewInt(20000)).String(),
		helpers.EtherToWei(sdk.NewInt(500000)).String(),
		helpers.EtherToWei(sdk.NewInt(100000000000000)).String(),
		"",
		val.Moniker,
		addBasicFlagsForTxCmd([]string{}, s.cfg.BondDenom)...,
	)
	require.NoError(err)

	info, _, err := val.ClientCtx.Keyring.NewMnemonic("sellall", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, etherminthd.EthSecp256k1)
	require.NoError(err)

	newAddr, err := info.GetAddress()
	require.NoError(err)
	_, err = MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, helpers.EtherToWei(sdk.NewInt(10000)))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, helpers.EtherToWei(sdk.NewInt(1)))).String()),
	)
	require.NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "valid sell all coins request",
			args: []string{
				s.cfg.BondDenom,                  // symCoinToSell
				fmt.Sprintf("1000000%s", symbol), // minAmountToBuyCoins
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr.String()),
			},
			expectErr: false,
		},
		{
			name: "invalid symbol sell coin request",
			args: []string{
				"invalidcoinsymbol",                      // symCoinToSell
				fmt.Sprintf("100000000000000%s", symbol), // minAmountTIBuyCoins
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr.String()),
			},
			expectErr: true,
		},
		{
			name: "buy coin not exist request",
			args: []string{
				s.cfg.BondDenom, // symCoinToSell
				fmt.Sprintf("100000000000000%s", "notexistcoin"), // minAmountTIBuyCoins
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr.String()),
			},
			expectErr: true,
		},
		{
			name: "invalid amount buy coin request",
			args: []string{
				s.cfg.BondDenom, // symCoinToSell
				"invalidamount", // minAmountTIBuyCoins
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr.String()),
			},
			expectErr: true,
		},
		{
			name: "balance less than minAmountToBuy request",
			args: []string{
				s.cfg.BondDenom,                         // symCoinToSell
				fmt.Sprintf("10000000000000%s", symbol), // minAmountTIBuyCoins
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr.String()),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewSellAllCoinCmd()
			clientCtx := val.ClientCtx
			args := addBasicFlagsForTxCmd(tc.args, s.cfg.BondDenom)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestIssueCheckCmd() {
	val := s.network.Validators[0]
	require := s.Require()

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "valid issue-check request",
			args: []string{
				fmt.Sprintf("10000%s", s.cfg.BondDenom), // amount
				"9",                                     // nonce
				"123213",                                // dueBlock
				"",                                      // password
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Moniker),
			},
			expectErr: false,
		},
		{
			name: "invalid dueBlock request",
			args: []string{
				fmt.Sprintf("10000%s", s.cfg.BondDenom), // amount
				"9",                                     // nonce
				"182913074327431298523432412341234489712390", // dueBlock
				"", // password
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Moniker),
			},
			expectErr: true,
		},
		{
			name: "not exist coin request",
			args: []string{
				fmt.Sprintf("10000%s", "notexistcoin"), // amount
				"9",                                    // nonce
				"123213",                               // dueBlock
				"",                                     // password
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Moniker),
			},
			expectErr: true,
		},
		{
			name: "not exist coin request",
			args: []string{
				"invalidcoinamount", // amount
				"9",                 // nonce
				"123213",            // dueBlock
				"",                  // password
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Moniker),
			},
			expectErr: true,
		},
		{
			name: "not exist key request",
			args: []string{
				fmt.Sprintf("10000%s", s.cfg.BondDenom), // amount
				"9",                                     // nonce
				"123213",                                // dieBlock
				"",                                      // password
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "notexistkey"),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewIssueCheckCmd()
			clientCtx := val.ClientCtx
			args := addBasicFlagsForTxCmd(tc.args, s.cfg.BondDenom)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCheckRedeemCmd() {
	val := s.network.Validators[0]
	require := s.Require()

	pass := ""
	hash, err := MsgIssueCheckExec(val.ClientCtx, val.Address.String(), fmt.Sprintf("10000%s", s.cfg.BondDenom), "9", "123", pass)
	require.NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "valid redeem check request",
			args: []string{
				hash.String(), // check
				pass,          // pass
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Moniker),
			},
			expectErr: false,
		},
		{
			name: "invalid base58 hash request",
			args: []string{
				"invalidBase58hash", // check
				"",                  // pass
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Moniker),
			},
			expectErr: true,
		},
		{
			name: "invalid check hash request",
			args: []string{
				base58.Encode([]byte("invalidcheckstructure")), // check
				"invalidpass", // pass
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Moniker),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewRedeemCheckCmd()
			clientCtx := val.ClientCtx
			args := addBasicFlagsForTxCmd(tc.args, s.cfg.BondDenom)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryCoinCmd() {
	val := s.network.Validators[0]
	require := s.Require()

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name:      "valid query get coin request",
			args:      []string{s.cfg.BondDenom},
			expectErr: false,
		},
		{
			name:      "not exist coin request",
			args:      []string{"invalid coin symbol"},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.QueryCoinCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryCoinsCmd() {
	val := s.network.Validators[0]
	require := s.Require()

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "valid get coins request",
			args: []string{
				fmt.Sprintf("--%s=%d", flags.FlagPage, 1),
				fmt.Sprintf("--%s=%d", flags.FlagOffset, 0),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			expectErr: false,
		},
		{
			name: "invalid page params in get coins request",
			args: []string{
				fmt.Sprintf("--%s=%d", flags.FlagPage, 3),
				fmt.Sprintf("--%s=%d", flags.FlagOffset, 1),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.QueryCoinsCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryCheckCmd() {
	val := s.network.Validators[0]
	require := s.Require()

	pass := ""
	hash, err := MsgIssueCheckExec(val.ClientCtx, val.Address.String(), fmt.Sprintf("10000%s", s.cfg.BondDenom), "9", "123", pass)
	require.NoError(err)

	_, err = MsgCheckRedeemExec(val.ClientCtx, hash.String(), pass, val.Address.String(), addBasicFlagsForTxCmd([]string{}, s.cfg.BondDenom)...)
	require.NoError(err)

	_, err = s.network.WaitForHeight(10)
	require.NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "valid get check request",
			args: []string{
				hash.String(), // check
			},
			expectErr: false,
		},
		{
			name: "invalid base58 check hash request",
			args: []string{
				"invalidBase58", // check
			},
			expectErr: true,
		},
		{
			name: "invalid check hash request",
			args: []string{
				base58.Encode([]byte("invalidcheckstructure")), // check
			},
			expectErr: true,
		},
		{
			name: "not exist check request",
			args: []string{
				"RdZpEdZnxmSEiH2eQYiXspVVbFSVVqr5DcAxg2ubgeMburp3FUUCD49xySyedDVWbMyLLD7y3mk6CsLS4v92nX9EGfbUc6R65vnYEX7DNbDCv7dwUmiPNQj3Ssi57h6D7eGXeA9586xobVqr4v7NviGNJnX6tdtxFthSA9ukGqUu12ZoQceRe8Qh8BCKsMCGVhKrjemnbwvqUSJfTtpTxTKvCAvqWbRCZU6g4Ce",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.QueryCheckCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryChecksCmd() {
	val := s.network.Validators[0]
	require := s.Require()

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name: "valid get checks request",
			args: []string{
				fmt.Sprintf("--%s=%d", flags.FlagPage, 1),
				fmt.Sprintf("--%s=%d", flags.FlagOffset, 0),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			expectErr: false,
		},
		{
			name: "invalid page params for checks request",
			args: []string{
				fmt.Sprintf("--%s=%d", flags.FlagPage, 3),
				fmt.Sprintf("--%s=%d", flags.FlagOffset, 1),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.QueryChecksCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				require.Error(err)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
			}
		})
	}
}

func addBasicFlagsForTxCmd(args []string, bondDenom string) []string {
	return append(args,
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(bondDenom, helpers.EtherToWei(sdk.NewInt(1)))).String()),
	)
}
