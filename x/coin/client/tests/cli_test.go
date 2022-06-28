//go:build norace

package cli_tests

import (
	"bitbucket.org/decimalteam/go-smart-node/testutil/network"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)

type CliIntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *CliIntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	if testing.Short() {
		s.T().Skip("skipping test in unit-tests mode.")
	}
	cfg := network.DefaultConfig()

	s.cfg = cfg

	baseDir, err := ioutil.TempDir(s.T().TempDir(), cfg.ChainID)
	require.NoError(s.T(), err)
	s.T().Logf("created temporary directory: %s", baseDir)

	s.network, err = network.New(s.T(), baseDir, cfg)
	s.Require().NoError(err)
}

func (s *CliIntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestCliIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CliIntegrationTestSuite))
}

func (s *CliIntegrationTestSuite) TestCreateCoinCmd() {
	require := s.Require()
	val := s.network.Validators[0]

	info, _, err := val.ClientCtx.Keyring.NewMnemonic("NewValidator", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	require.NoError(err)

	newAddr := sdk.AccAddress(info.GetPubKey().Address())
	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, helpers.EtherToWei(sdk.NewInt(10000000000000)))).String()),
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
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr),
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
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr),
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
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr),
			},
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewCreateValidatorCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				require.Error(err)
				require.Nil(out)
			} else {
				require.NoError(err, "test: %s\noutput: %s", tc.name, out.String())
				//err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType)
				//require.NoError(err, out.String(), "test: %s, output\n:", tc.name, out.String())

				//txResp := tc.respType.(*sdk.TxResponse)
				//require.Equal(tc.expectedCode, txResp.Code,
				//	"test: %s, output\n:", tc.name, out.String())

				//events := txResp.Logs[0].GetEvents()
				//for i := 0; i < len(events); i++ {
				//	if events[i].GetType() == "create_validator" {
				//		attributes := events[i].GetAttributes()
				//		require.Equal(attributes[1].Value, "100stake")
				//		break
				//	}
				//}
			}
		})
	}
}

//func (s *CliIntegrationTestSuite) TestCreateCoinCmd() {
//	val := s.network.Validators[0]
//
//}
//func (s *CliIntegrationTestSuite) TestCreateCoinCmd() {
//	val := s.network.Validators[0]
//
//}
