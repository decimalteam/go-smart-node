package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	cometbftdb "github.com/cometbft/cometbft-db"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	cmtstate "github.com/tendermint/tendermint/state"
	cmtstore "github.com/tendermint/tendermint/store"
	tmtypes "github.com/tendermint/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
)

// ForkSchedulePlanCmd schedules an x/upgrade plan directly in an offline fork COPY's
// app state so the fork can rehearse an in-place upgrade (e.g. the DEL
// redenomination) without going through governance. It mirrors fork-fund's
// "mutate app state, commit -> H+1, fabricate+bootstrap CometBFT to H+1" recipe.
//
// The plan is scheduled at EXACTLY comet H+2 — the first block the restarted node
// produces — because the x/upgrade BeginBlocker panics ("BINARY UPDATED BEFORE
// TRIGGER") on any block where a handler is already registered for a plan that is
// not yet due. Scheduling at the first post-restart height leaves no such block.
//
// The plan name is "<prefix>/<exec-height>", which must have a handler registered in
// app/upgradeslist.go for this chain-id; the command refuses to touch state when the
// handler is missing, so an upgradeslist/plan-height mismatch is caught before any
// mutation instead of halting the node at the upgrade height.
//
// The node MUST be stopped (LevelDB single-writer). Work only in a disposable copy.
func ForkSchedulePlanCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fork-scheduleplan",
		Short: "Schedule an x/upgrade plan at the first post-restart height in an offline fork COPY (app + CometBFT -> H+1)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			logger := serverCtx.Logger
			home := cast.ToString(appOpts.Get(flags.FlagHome))
			if home == "" {
				return fmt.Errorf("--home is required")
			}
			out := cmd.OutOrStdout()
			prefix, _ := cmd.Flags().GetString("name-prefix")
			info, _ := cmd.Flags().GetString("info")
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

			// Consensus resumes at Hc+2 after the H+1 jump below; the plan must be due on
			// exactly that first block (see the pre-trigger panic note in the doc comment).
			execHeight := Hc + 2
			name := fmt.Sprintf("%s/%d", prefix, execHeight)

			var appHashNew []byte
			switch {
			case appH == Hc:
				if !dscApp.UpgradeKeeper.HasHandler(name) {
					return fmt.Errorf("no upgrade handler registered for %q (chain-id %s) — add it to app/upgradeslist.go at exactly this height and rebuild before scheduling", name, chainID)
				}
				cms := dscApp.CommitMultiStore()
				ctx := sdk.NewContext(cms, tmproto.Header{Height: appH, ChainID: chainID, Time: stateH.LastBlockTime}, false, logger)
				plan := upgradetypes.Plan{Name: name, Height: execHeight, Info: info}
				if err := dscApp.UpgradeKeeper.ScheduleUpgrade(ctx, plan); err != nil {
					return fmt.Errorf("schedule upgrade: %w", err)
				}
				commitID := cms.Commit()
				appHashNew = commitID.Hash
				fmt.Fprintf(out, "scheduled plan %q at height %d\n", name, execHeight)
				fmt.Fprintf(out, "  app committed %d -> %d appHash %X\n", appH, appH+1, appHashNew)
			case appH == Hc+1:
				// recovery: app already scheduled/committed one ahead; only advance comet
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

			fmt.Fprintf(out, "DONE: store==state==app==%d, plan %q due at %d (first post-restart block). Start the node.\n", targetH, name, targetH+1)
			return nil
		},
	}
	cmd.Flags().String("name-prefix", "https://repo.decimalchain.com", "plan name prefix; full name is <prefix>/<exec-height>")
	cmd.Flags().String("info", "{}", "plan info payload (JSON)")
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory (COPY only!)")
	return cmd
}
