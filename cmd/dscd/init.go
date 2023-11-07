package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	cfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/libs/cli"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	pvm "github.com/cometbft/cometbft/privval"
	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/cosmos/go-bip39"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	validatorkeeper "bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) printInfo {
	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(info printInfo) error {
	out, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(os.Stderr, "%s\n", string(sdk.MustSortJSON(out))); err != nil {
		return err
	}

	return nil
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(mbm module.BasicManager, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			// Initialize config with default values
			config := serverCtx.Config
			config.Moniker = args[0]
			config.P2P.MaxNumInboundPeers = 100
			config.P2P.MaxNumOutboundPeers = 30
			config.Mempool.Size = 10000
			config.StateSync.TrustPeriod = 112 * time.Hour
			config.SetRoot(clientCtx.HomeDir)

			// Get chain id from flags
			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if len(chainID) == 0 {
				chainID = fmt.Sprintf("decimal_202020-%v", tmrand.Str(6))
			}

			// Get bip39 mnemonic
			var mnemonic string
			recover, _ := cmd.Flags().GetBool(genutilcli.FlagRecover)
			if recover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				value, err := input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				mnemonic = value
				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			// Initialize validator and node files
			nodeID, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			genFile := config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(genutilcli.FlagOverwrite)

			if !overwrite && tmos.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}

			appState, err := json.MarshalIndent(mbm.DefaultGenesis(clientCtx.Codec), "", " ")
			if err != nil {
				return errors.Wrap(err, "Failed to marshall default genesis state")
			}

			genDoc := &tmtypes.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				genDoc, err = tmtypes.GenesisDocFromFile(genFile)
				if err != nil {
					return errors.Wrap(err, "Failed to read genesis doc from file")
				}
			}

			genDoc.ChainID = chainID
			genDoc.Validators = nil
			genDoc.AppState = appState

			if err := genutil.ExportGenesisFile(genDoc, genFile); err != nil {
				return errors.Wrap(err, "Failed to export gensis file")
			}

			toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return displayInfo(toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(genutilcli.FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().Bool(genutilcli.FlagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")

	return cmd
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func SelfDelegationCmd(mbm module.BasicManager, txEncCfg client.TxEncodingConfig, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "selfdelegation [stake]",
		Short: "Add private validator to genesis file as validator with stake in base coin from specified account",
		Long: fmt.Sprintf(`Add private validator to genesis file as validator with stake from specified account.
Like gentx, but just add validator without transaction.

Example:
$ %s selfdelegation 100000000del --home=/path/to/home/dir --from keyname
`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			cdc := clientCtx.Codec

			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			privValidator := pvm.LoadFilePV(serverCtx.Config.PrivValidatorKeyFile(), serverCtx.Config.PrivValidatorStateFile())

			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis doc file %s", config.GenesisFile())
			}

			var genesisState map[string]json.RawMessage
			if err = json.Unmarshal(genDoc.AppState, &genesisState); err != nil {
				return errors.Wrap(err, "failed to unmarshal genesis state")
			}

			if err = mbm.ValidateGenesis(cdc, txEncCfg, genesisState); err != nil {
				return errors.Wrap(err, "failed to validate genesis state")
			}

			// create validator
			tmpubkey, err := privValidator.GetPubKey()
			if err != nil {
				return err
			}
			pubkey, err := cryptocodec.FromTmPubKeyInterface(tmpubkey)
			if err != nil {
				return err
			}
			valAdr := sdk.ValAddress(pubkey.Address())
			validator, err := validatortypes.NewValidator(
				valAdr,
				clientCtx.FromAddress,
				pubkey,
				validatortypes.NewDescription(
					config.Moniker,
					config.Moniker,
					"http://example.org",
					"example@exxample.org",
					"Details",
				),
				sdk.MustNewDecFromStr("0.10"),
			)
			validator.Status = validatortypes.BondStatus_Bonded
			validator.Online = true
			// create delegation
			stakeCoin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}
			delegation := validatortypes.NewDelegation(clientCtx.FromAddress, valAdr, validatortypes.NewStakeCoin(stakeCoin))
			// stake power
			power := validatortypes.LastValidatorPower{
				Address: valAdr.String(),
				Power:   validatorkeeper.TokensToConsensusPower(stakeCoin.Amount),
			}
			// record for tendermint validators
			// tmVal := tmtypes.GenesisValidator{
			// 	Address: sdk.ConsAddress(pubkey.Address()).Bytes(),
			// 	Name:    config.Moniker,
			// 	Power:   validatorkeeper.TokensToConsensusPower(stakeCoin.Amount),
			// 	PubKey:  tmpubkey,
			// }

			// insert into validator state
			var vgs validatortypes.GenesisState
			cdc.MustUnmarshalJSON(genesisState["validator"], &vgs)
			vgs.Validators = append(vgs.Validators, validator)
			vgs.Delegations = append(vgs.Delegations, delegation)
			vgs.LastValidatorPowers = append(vgs.LastValidatorPowers, power)
			vgs.LastTotalPower += power.Power
			genesisState["validator"] = cdc.MustMarshalJSON(&vgs)
			// insert into bond pool
			poolAddress, _ := bech32.ConvertAndEncode(cmdcfg.Bech32Prefix, authtypes.NewModuleAddress("bonded_tokens_pool"))
			var bgs banktypes.GenesisState
			cdc.MustUnmarshalJSON(genesisState["bank"], &bgs)
			var added = false
			for i := range bgs.Balances {
				if bgs.Balances[i].Address == poolAddress {
					bgs.Balances[i].Coins = bgs.Balances[i].Coins.Add(stakeCoin)
					added = true
				}
			}
			if !added {
				bgs.Balances = append(bgs.Balances, banktypes.Balance{
					Address: poolAddress,
					Coins:   sdk.NewCoins(stakeCoin),
				})
			}
			genesisState["bank"] = cdc.MustMarshalJSON(&bgs)
			// insert into tendermint validators
			//genDoc.Validators = append(genDoc.Validators, tmVal)

			genDoc.AppState, err = json.Marshal(genesisState)
			if err != nil {
				return errors.Wrap(err, "Failed to marshal app state")
			}

			if err := genutil.ExportGenesisFile(genDoc, config.GenesisFile()); err != nil {
				return errors.Wrap(err, "Failed to export gensis file")
			}

			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flags.FlagFrom, "", "key of staker")

	return cmd
}
