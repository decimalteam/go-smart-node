package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cometbftdb "github.com/cometbft/cometbft-db"
	"github.com/ethereum/go-ethereum/common"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	cmtstate "github.com/tendermint/tendermint/state"
	cmtstore "github.com/tendermint/tendermint/store"
	tmtypes "github.com/tendermint/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// ForkFundCmd credits an eth address a large base-denom balance directly in an
// offline fork COPY, then advances BOTH the app and CometBFT state one height so the
// node boots consistently (app==store==state) — mirroring fork-takeover's H+1 jump.
//
// It mints via the coin module (Minter perms), sends to sdk.AccAddress(ethAddr), and
// commits the app CommitMultiStore -> H+1 with new app hash A'. It then fabricates an
// empty block H+1 signed by OUR priv_validator key (the sole validator, unchanged
// since fork-takeover), SaveBlock, and Bootstraps CometBFT state to H+1 with AppHash=A'.
// On restart store==state==app==H+1 -> no replay, solo consensus resumes at H+2, and the
// funded balance is visible to the EVM (eth_getBalance reads the bank balance).
//
// Recovery: if the app is already committed one height ahead of CometBFT (e.g. a prior
// run committed the fund but died before advancing CometBFT), it skips funding and only
// advances CometBFT to match, using the app's existing LastCommitID hash.
//
// The node MUST be stopped (LevelDB single-writer). Work only in a disposable copy.
func ForkFundCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fork-fund",
		Short: "Credit an eth address a large base-denom balance in an offline fork COPY (app + CometBFT -> H+1)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			logger := serverCtx.Logger
			home := cast.ToString(appOpts.Get(flags.FlagHome))
			if home == "" {
				return fmt.Errorf("--home is required")
			}
			out := cmd.OutOrStdout()
			addrStr, _ := cmd.Flags().GetString("addr")
			amountStr, _ := cmd.Flags().GetString("amount")
			if !common.IsHexAddress(addrStr) {
				return fmt.Errorf("invalid --addr %q (expected 0x-hex eth address)", addrStr)
			}
			ethAddr := common.HexToAddress(addrStr)
			amount, ok := sdkmath.NewIntFromString(amountStr)
			if !ok || !amount.IsPositive() {
				return fmt.Errorf("invalid --amount %q (expected positive integer in base-denom units)", amountStr)
			}
			chainID := readChainID(home)
			if chainID == "" {
				return fmt.Errorf("could not read chain-id from genesis")
			}

			// --- OUR consensus key (unchanged since fork-takeover) ---
			ourPriv, ourPubTm, err := loadForkPrivKey(filepath.Join(home, "config", "priv_validator_key.json"))
			if err != nil {
				return fmt.Errorf("load priv_validator_key: %w", err)
			}
			ourAddr := ourPubTm.Address()

			// --- open CometBFT state + blockstore ---
			dataDir := filepath.Join(home, "data")
			stateDB, err := cometbftdb.NewGoLevelDB("state", dataDir)
			if err != nil {
				return fmt.Errorf("open state.db: %w", err)
			}
			defer stateDB.Close()
			blockDB, err := cometbftdb.NewGoLevelDB("blockstore", dataDir)
			if err != nil {
				return fmt.Errorf("open blockstore.db: %w", err)
			}
			defer blockDB.Close()
			stateStore := cmtstate.NewStore(stateDB, cmtstate.StoreOptions{})
			blockStore := cmtstore.NewBlockStore(blockDB)
			stateH, err := stateStore.Load()
			if err != nil {
				return fmt.Errorf("load cometbft state: %w", err)
			}
			if stateH.IsEmpty() {
				return fmt.Errorf("cometbft state is empty")
			}
			Hc := stateH.LastBlockHeight
			commitHc := blockStore.LoadSeenCommit(Hc)
			if commitHc == nil {
				commitHc = blockStore.LoadBlockCommit(Hc)
			}
			if commitHc == nil {
				return fmt.Errorf("no commit for comet height %d", Hc)
			}

			// --- OUR validator set / power (reuse the existing comet valset) ---
			_, ourVal := stateH.Validators.GetByAddress(ourAddr)
			if ourVal == nil {
				return fmt.Errorf("our validator %X not in comet valset at H=%d", ourAddr, Hc)
			}
			power := ourVal.VotingPower
			ourSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{tmtypes.NewValidator(ourPubTm, power)})

			// --- open app ---
			db, err := openRedenomAppDB(home, appOpts)
			if err != nil {
				return fmt.Errorf("open application db: %w", err)
			}
			defer db.Close()
			dscApp := app.NewDSC(logger, db, nil, true, map[int64]bool{}, home, uint(1), encCfg, appOpts)
			appH := dscApp.LastBlockHeight()

			var appHashNew []byte
			switch {
			case appH == Hc:
				// normal path: fund then commit -> Hc+1
				cms := dscApp.CommitMultiStore()
				ctx := sdk.NewContext(cms, tmproto.Header{Height: appH, ChainID: chainID}, false, logger)
				base := dscApp.ValidatorKeeper.BaseDenom(ctx)
				coins := sdk.NewCoins(sdk.NewCoin(base, amount))
				acc := sdk.AccAddress(ethAddr.Bytes())
				if dscApp.AccountKeeper.GetAccount(ctx, acc) == nil {
					dscApp.AccountKeeper.SetAccount(ctx, dscApp.AccountKeeper.NewAccountWithAddress(ctx, acc))
				}
				if err := dscApp.BankKeeper.MintCoins(ctx, cointypes.ModuleName, coins); err != nil {
					return fmt.Errorf("mint %s: %w", coins, err)
				}
				if err := dscApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, cointypes.ModuleName, acc, coins); err != nil {
					return fmt.Errorf("send to account: %w", err)
				}
				after := dscApp.BankKeeper.GetBalance(ctx, acc, base)
				commitID := cms.Commit()
				appHashNew = commitID.Hash
				fmt.Fprintf(out, "funded eth %s (cosmos %s)\n", ethAddr.Hex(), acc.String())
				fmt.Fprintf(out, "  minted %s%s; balance -> %s\n", amount.String(), base, after.Amount.String())
				fmt.Fprintf(out, "  app committed %d -> %d appHash %X\n", appH, appH+1, appHashNew)
			case appH == Hc+1:
				// recovery: app already funded/committed one ahead; only advance comet
				appHashNew = dscApp.LastCommitID().Hash
				fmt.Fprintf(out, "recovery: app already at %d (comet %d); advancing comet to match; appHash %X\n", appH, Hc, appHashNew)
			default:
				return fmt.Errorf("unexpected heights app=%d comet=%d (want app==comet or app==comet+1)", appH, Hc)
			}

			targetH := Hc + 1

			// --- fabricate empty block targetH signed by OUR key ---
			makeState := stateH.Copy()
			makeState.Validators = ourSet.Copy()
			makeState.NextValidators = ourSet.Copy()
			block, partSet := makeState.MakeBlock(targetH, nil, commitHc, nil, ourAddr)
			blockID := tmtypes.BlockID{Hash: block.Hash(), PartSetHeader: partSet.Header()}
			// precommit ts becomes the BFT median time of the LastCommit for targetH+1;
			// must be strictly greater than targetH's block time.
			ts := block.Time.Add(time.Second)
			vote := &tmtypes.Vote{
				Type: tmproto.PrecommitType, Height: targetH, Round: 0,
				BlockID: blockID, Timestamp: ts, ValidatorAddress: ourAddr, ValidatorIndex: 0,
			}
			sig, err := ourPriv.Sign(tmtypes.VoteSignBytes(chainID, vote.ToProto()))
			if err != nil {
				return fmt.Errorf("sign precommit: %w", err)
			}
			seenCommit := &tmtypes.Commit{
				Height: targetH, Round: 0, BlockID: blockID,
				Signatures: []tmtypes.CommitSig{tmtypes.NewCommitSigForBlock(sig, ourAddr, ts)},
			}
			blockStore.SaveBlock(block, partSet, seenCommit)

			// --- bootstrap CometBFT state to targetH ---
			s1 := stateH.Copy()
			s1.LastBlockHeight = targetH
			s1.LastBlockID = blockID
			s1.LastBlockTime = block.Time
			s1.LastValidators = ourSet.Copy()
			s1.Validators = ourSet.Copy()
			s1.NextValidators = ourSet.Copy()
			s1.LastHeightValidatorsChanged = targetH
			s1.AppHash = appHashNew
			s1.LastResultsHash = tmtypes.NewResults(nil).Hash()
			if err := stateStore.Bootstrap(s1); err != nil {
				return fmt.Errorf("bootstrap cometbft state: %w", err)
			}

			// --- reset WAL + priv validator state ---
			_ = os.RemoveAll(filepath.Join(dataDir, "cs.wal"))
			_ = os.WriteFile(filepath.Join(dataDir, "priv_validator_state.json"),
				[]byte("{\n  \"height\": \"0\",\n  \"round\": 0,\n  \"step\": 0\n}\n"), 0600)

			fmt.Fprintf(out, "DONE: store==state==app==%d appHash %X. Start the node.\n", targetH, appHashNew)
			return nil
		},
	}
	cmd.Flags().String("addr", "", "0x eth address to fund")
	cmd.Flags().String("amount", "", "base-denom amount as an integer (18-decimal base units, e.g. 1e21 = 1000 coins)")
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory (COPY only!)")
	return cmd
}
