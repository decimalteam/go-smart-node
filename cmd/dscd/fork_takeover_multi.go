package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdked25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	cometbftdb "github.com/cometbft/cometbft-db"
	tmed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	cmtstate "github.com/tendermint/tendermint/state"
	cmtstore "github.com/tendermint/tendermint/store"
	tmtypes "github.com/tendermint/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// ForkTakeoverMultiCmd is fork-takeover for N validators: the N highest-power
// bonded+online validators get the consensus keys from --keys (keeping their real
// stake), become the whole CometBFT set, and sign the fabricated block H+1 together.
// Clone the result into N homes, give home i priv_validator_key i, and peer them: the
// fork then runs real multi-validator consensus. Work only in a disposable COPY.
func ForkTakeoverMultiCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fork-takeover-multi",
		Short: "Rewrite an offline node COPY so N local keys become the validator set",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			home := cast.ToString(appOpts.Get(flags.FlagHome))
			keysFlag, _ := cmd.Flags().GetString("keys")
			out := cmd.OutOrStdout()
			chainID := readChainID(home)
			if chainID == "" {
				return fmt.Errorf("could not read chain-id from genesis")
			}
			var privs []tmed25519.PrivKey
			var pubs []tmed25519.PubKey
			for _, p := range strings.Split(keysFlag, ",") {
				priv, pub, err := loadForkPrivKey(strings.TrimSpace(p))
				if err != nil {
					return fmt.Errorf("load %s: %w", p, err)
				}
				privs = append(privs, priv)
				pubs = append(pubs, pub)
			}
			if len(privs) < 2 {
				return fmt.Errorf("need at least 2 keys")
			}

			dataDir := filepath.Join(home, "data")
			stateDB, err := cometbftdb.NewGoLevelDB("state", dataDir)
			if err != nil {
				return err
			}
			defer stateDB.Close()
			blockDB, err := cometbftdb.NewGoLevelDB("blockstore", dataDir)
			if err != nil {
				return err
			}
			defer blockDB.Close()
			stateStore := cmtstate.NewStore(stateDB, cmtstate.StoreOptions{})
			blockStore := cmtstore.NewBlockStore(blockDB)
			stateH, err := stateStore.Load()
			if err != nil || stateH.IsEmpty() {
				return fmt.Errorf("load cometbft state: %v", err)
			}
			H := stateH.LastBlockHeight
			if blockStore.Height() != H {
				return fmt.Errorf("blockstore height %d != state height %d", blockStore.Height(), H)
			}
			commitH := blockStore.LoadSeenCommit(H)
			if commitH == nil {
				commitH = blockStore.LoadBlockCommit(H)
			}
			if commitH == nil {
				return fmt.Errorf("no commit for %d", H)
			}

			pubBytes := make([][]byte, len(pubs))
			for i, p := range pubs {
				pubBytes[i] = p.Bytes()
			}
			planOffset, _ := cmd.Flags().GetInt64("plan-offset")
			planPrefix, _ := cmd.Flags().GetString("plan-prefix")
			var plan *upgradetypes.Plan
			if planOffset > 0 {
				exec := H + 1 + planOffset
				plan = &upgradetypes.Plan{Name: fmt.Sprintf("%s/%d", planPrefix, exec), Height: exec, Info: "{}"}
				fmt.Fprintf(out, "scheduling plan %q at height %d\n", plan.Name, exec)
			}
			appHashNew, powers, ops, err := forkMutateAppStateMulti(serverCtx.Logger, home, appOpts, encCfg, H, chainID, pubBytes, plan)
			if err != nil {
				return fmt.Errorf("mutate app state: %w", err)
			}
			vals := make([]*tmtypes.Validator, len(pubs))
			for i := range pubs {
				vals[i] = tmtypes.NewValidator(pubs[i], powers[i])
				fmt.Fprintf(out, "validator %d: %s power=%d consaddr=%X\n", i, ops[i], powers[i], pubs[i].Address())
			}
			set := tmtypes.NewValidatorSet(vals)

			makeState := stateH.Copy()
			makeState.Validators = set.Copy()
			makeState.NextValidators = set.Copy()
			proposer := set.GetProposer().Address
			block, partSet := makeState.MakeBlock(H+1, nil, commitH, nil, proposer)
			blockID := tmtypes.BlockID{Hash: block.Hash(), PartSetHeader: partSet.Header()}
			ts := block.Time.Add(time.Second)

			privByAddr := map[string]tmed25519.PrivKey{}
			for i, p := range pubs {
				privByAddr[string(p.Address())] = privs[i]
			}
			sigs := make([]tmtypes.CommitSig, set.Size())
			for idx, v := range set.Validators {
				vote := &tmtypes.Vote{Type: tmproto.PrecommitType, Height: H + 1, Round: 0, BlockID: blockID,
					Timestamp: ts, ValidatorAddress: v.Address, ValidatorIndex: int32(idx)}
				sig, err := privByAddr[string(v.Address)].Sign(tmtypes.VoteSignBytes(chainID, vote.ToProto()))
				if err != nil {
					return err
				}
				sigs[idx] = tmtypes.NewCommitSigForBlock(sig, v.Address, ts)
			}
			seen := &tmtypes.Commit{Height: H + 1, Round: 0, BlockID: blockID, Signatures: sigs}
			blockStore.SaveBlock(block, partSet, seen)

			s1 := stateH.Copy()
			s1.LastBlockHeight = H + 1
			s1.LastBlockID = blockID
			s1.LastBlockTime = block.Time
			s1.LastValidators = set.Copy()
			s1.Validators = set.Copy()
			s1.NextValidators = set.Copy()
			s1.LastHeightValidatorsChanged = H + 1
			s1.AppHash = appHashNew
			s1.LastResultsHash = tmtypes.NewResults(nil).Hash()
			if err := stateStore.Bootstrap(s1); err != nil {
				return err
			}
			if err := os.RemoveAll(filepath.Join(dataDir, "cs.wal")); err != nil {
				return err
			}
			_ = os.WriteFile(filepath.Join(dataDir, "priv_validator_state.json"),
				[]byte("{\n  \"height\": \"0\",\n  \"round\": 0,\n  \"step\": 0\n}\n"), 0600)
			fmt.Fprintf(out, "DONE: %d-validator fork ready at height %d appHash %X\n", set.Size(), H+1, appHashNew)
			return nil
		},
	}
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory (COPY only!)")
	cmd.Flags().String("keys", "", "comma-separated priv_validator_key.json paths, one per validator")
	cmd.Flags().Int64("plan-offset", 0, "also schedule an x/upgrade plan this many blocks after the fork height (0 = none)")
	cmd.Flags().String("plan-prefix", "https://repo.decimalchain.com", "plan name prefix; full name is <prefix>/<exec-height>")
	return cmd
}

func forkMutateAppStateMulti(logger log.Logger, home string, appOpts servertypes.AppOptions, encCfg params.EncodingConfig, H int64, chainID string, pubs [][]byte, plan *upgradetypes.Plan) ([]byte, []int64, []string, error) {
	db, err := openRedenomAppDB(home, appOpts)
	if err != nil {
		return nil, nil, nil, err
	}
	defer db.Close()
	dscApp := app.NewDSC(logger, db, nil, true, map[int64]bool{}, home, uint(1), encCfg, appOpts)
	if got := dscApp.LastBlockHeight(); got != H {
		return nil, nil, nil, fmt.Errorf("app height %d != cometbft height %d", got, H)
	}
	cms := dscApp.CommitMultiStore()
	ctx := sdk.NewContext(cms, tmproto.Header{Height: H, ChainID: chainID}, false, logger)
	k := dscApp.ValidatorKeeper
	store := ctx.KVStore(dscApp.GetKey(validatortypes.StoreKey))

	all := k.GetAllValidators(ctx)
	var cand []validatortypes.Validator
	for _, v := range all {
		// Online is not required: a prior single fork-takeover (needed by fork-setowner /
		// fork-fund) marks every other validator offline. Chosen ones are set online.
		if v.Status == validatortypes.BondStatus_Bonded && !v.Jailed && v.Stake > 0 {
			cand = append(cand, v)
		}
	}
	sort.Slice(cand, func(i, j int) bool {
		if cand[i].Stake != cand[j].Stake {
			return cand[i].Stake > cand[j].Stake
		}
		return cand[i].OperatorAddress < cand[j].OperatorAddress
	})
	if len(cand) < len(pubs) {
		return nil, nil, nil, fmt.Errorf("only %d bonded+online validators for %d keys", len(cand), len(pubs))
	}
	chosen := cand[:len(pubs)]

	wipePrefix(store, validatortypes.GetValidatorsByPowerIndexKey())
	wipePrefix(store, validatortypes.GetLastValidatorPowersKey())
	powers := make([]int64, len(pubs))
	ops := make([]string, len(pubs))
	total := int64(0)
	keep := map[string]bool{}
	for i, V := range chosen {
		oldCons, err := V.GetConsAddr()
		if err != nil {
			return nil, nil, nil, err
		}
		store.Delete(validatortypes.GetValidatorByConsAddrIndexKey(oldCons))
		pkAny, err := codectypes.NewAnyWithValue(&sdked25519.PubKey{Key: pubs[i]})
		if err != nil {
			return nil, nil, nil, err
		}
		power := V.Stake
		V.ConsensusPubkey = pkAny
		V.Online = true
		V.Jailed = false
		V.Status = validatortypes.BondStatus_Bonded
		k.SetValidator(ctx, V)
		k.SetValidatorRS(ctx, V.GetOperator(), validatortypes.ValidatorRS{Rewards: V.Rewards, TotalRewards: V.TotalRewards, Stake: power})
		if err := k.SetValidatorByConsAddr(ctx, V); err != nil {
			return nil, nil, nil, err
		}
		k.SetValidatorByPowerIndex(ctx, V)
		k.SetLastValidatorPower(ctx, V.GetOperator(), power)
		newCons, err := V.GetConsAddr()
		if err != nil {
			return nil, nil, nil, err
		}
		k.SetStartHeight(ctx, newCons, H)
		powers[i], ops[i] = power, V.OperatorAddress
		total += power
		keep[V.OperatorAddress] = true
	}
	k.SetLastTotalPower(ctx, sdk.NewInt(total))
	// Every other validator offline (see forkMutateAppState).
	for _, v := range all {
		if keep[v.OperatorAddress] || !v.Online {
			continue
		}
		v.Online = false
		k.SetValidator(ctx, v)
	}
	if plan != nil {
		if err := dscApp.UpgradeKeeper.ScheduleUpgrade(ctx, *plan); err != nil {
			return nil, nil, nil, fmt.Errorf("schedule upgrade: %w", err)
		}
	}
	commitID := cms.Commit()
	return commitID.Hash, powers, ops, nil
}
