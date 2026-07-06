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

	cometbftdb "github.com/cometbft/cometbft-db"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	cmtstate "github.com/tendermint/tendermint/state"
	cmtstore "github.com/tendermint/tendermint/store"
	tmtypes "github.com/tendermint/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
)

// ForkTimewarpCmd shrinks the x/validator AutoUnbondTimeout param directly in an
// offline fork COPY's app state, so ProcessAutoUnbond's offline-duration cutoff
// (normally 30 days) can be crossed with a short, realistic wall-clock wait after
// a live MsgSetOffline tx instead of waiting 30 real days. Mirrors fork-fund's
// "mutate app state, commit -> H+1, fabricate+bootstrap CometBFT to H+1" recipe so
// the node reboots consistently (app==store==state) with no replay.
//
// The node MUST be stopped (LevelDB single-writer). Work only in a disposable copy.
func ForkTimewarpCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fork-timewarp",
		Short: "Shrink x/validator AutoUnbondTimeout in an offline fork COPY (app + CometBFT -> H+1)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			logger := serverCtx.Logger
			home := cast.ToString(appOpts.Get(flags.FlagHome))
			if home == "" {
				return fmt.Errorf("--home is required")
			}
			out := cmd.OutOrStdout()
			seconds, _ := cmd.Flags().GetInt64("seconds")
			if seconds <= 0 {
				return fmt.Errorf("--seconds must be positive")
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

			offlineValoper, _ := cmd.Flags().GetString("offline-validator")

			var appHashNew []byte
			switch {
			case appH == Hc:
				// normal path: shrink AutoUnbondTimeout (and optionally force a validator
				// offline + stamp OfflineSince=now) then commit -> Hc+1
				cms := dscApp.CommitMultiStore()
				ctx := sdk.NewContext(cms, tmproto.Header{Height: appH, ChainID: chainID, Time: stateH.LastBlockTime}, false, logger)
				k := dscApp.ValidatorKeeper
				p := k.GetParams(ctx)
				before := p.AutoUnbondTimeout
				p.AutoUnbondTimeout = time.Duration(seconds) * time.Second
				if md, _ := cmd.Flags().GetUint32("max-delegations"); md > 0 {
					beforeMD := p.MaxDelegations
					p.MaxDelegations = md
					fmt.Fprintf(out, "MaxDelegations %d -> %d\n", beforeMD, md)
				}
				k.SetParams(ctx, p)
				after := k.AutoUnbondTimeout(ctx)
				fmt.Fprintf(out, "AutoUnbondTimeout %s -> %s\n", before, after)

				if offlineValoper != "" {
					valAddr, verr := sdk.ValAddressFromBech32(offlineValoper)
					if verr != nil {
						return fmt.Errorf("bad --offline-validator: %w", verr)
					}
					v, found := k.GetValidator(ctx, valAddr)
					if !found {
						return fmt.Errorf("validator %s not found", offlineValoper)
					}
					v.Online = false
					k.SetValidator(ctx, v)
					k.SetValidatorOfflineSince(ctx, valAddr, ctx.BlockTime())
					fmt.Fprintf(out, "validator %s Online=false OfflineSince=%s\n", offlineValoper, ctx.BlockTime())
				}

				commitID := cms.Commit()
				appHashNew = commitID.Hash
				fmt.Fprintf(out, "  app committed %d -> %d appHash %X\n", appH, appH+1, appHashNew)
			case appH == Hc+1:
				// recovery: app already mutated/committed one ahead; only advance comet
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
	cmd.Flags().Int64("seconds", 10, "new AutoUnbondTimeout in seconds")
	cmd.Flags().Uint32("max-delegations", 0, "if >0, set x/validator MaxDelegations param (for CheckDelegations trim test)")
	cmd.Flags().String("offline-validator", "", "optional valoper address to force Online=false + OfflineSince=now")
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory (COPY only!)")
	return cmd
}
