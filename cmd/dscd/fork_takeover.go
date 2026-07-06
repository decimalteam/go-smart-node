package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// ForkTakeoverCmd rewrites an offline node snapshot so that OUR priv_validator_key
// becomes the SOLE consensus validator, consistently across (a) the x/validator
// application state, (b) the CometBFT state.db validator sets, and (c) the block
// store. After it runs, `dscd start` boots into single-validator consensus and
// produces blocks. The node MUST be stopped. Work only in a disposable COPY.
//
// Mechanism:
//  1. Pick a bonded+online validator V (highest power). Swap its consensus pubkey to
//     OUR ed25519 pubkey, keep it bonded/online, and make it the only entry in the
//     power index and last-validator-powers. Commit -> app H -> H+1, new hash A'.
//  2. Read CometBFT State_H, fabricate an empty block H+1 whose validators are OUR
//     validator with a seen-commit signed by OUR key, and Bootstrap the CometBFT
//     state to H+1 (LastValidators=Validators=NextValidators=OUR set, AppHash=A').
//     On restart store==state==app==H+1 -> no replay, consensus resumes solo at H+2.
//  3. Reset the consensus WAL and priv_validator_state.
func ForkTakeoverCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fork-takeover",
		Short: "Rewrite an offline node COPY so our priv_validator_key is the sole validator",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			home := cast.ToString(appOpts.Get(flags.FlagHome))
			if home == "" {
				return fmt.Errorf("--home is required")
			}
			out := cmd.OutOrStdout()
			chainID := readChainID(home)
			if chainID == "" {
				return fmt.Errorf("could not read chain-id from genesis")
			}

			// --- OUR consensus key ---
			ourPriv, ourPubTm, err := loadForkPrivKey(filepath.Join(home, "config", "priv_validator_key.json"))
			if err != nil {
				return fmt.Errorf("load priv_validator_key: %w", err)
			}
			ourAddr := ourPubTm.Address()
			fmt.Fprintf(out, "fork validator addr %X chain-id %s\n", ourAddr, chainID)

			// --- read CometBFT state + blockstore ---
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
			H := stateH.LastBlockHeight
			if blockStore.Height() != H {
				return fmt.Errorf("blockstore height %d != state height %d", blockStore.Height(), H)
			}
			commitH := blockStore.LoadSeenCommit(H)
			if commitH == nil {
				commitH = blockStore.LoadBlockCommit(H)
			}
			if commitH == nil {
				return fmt.Errorf("no commit found for height %d", H)
			}
			fmt.Fprintf(out, "cometbft H=%d appHash(A)=%X lastVals=%d\n", H, stateH.AppHash, stateH.LastValidators.Size())

			// --- mutate app state and commit -> H+1 ---
			appHashNew, power, takenOver, err := forkMutateAppState(serverCtx.Logger, home, appOpts, encCfg, H, chainID, ourPubTm.Bytes())
			if err != nil {
				return fmt.Errorf("mutate app state: %w", err)
			}
			fmt.Fprintf(out, "took over %s power=%d newAppHash(A')=%X\n", takenOver, power, appHashNew)

			// --- OUR validator set ---
			ourSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{tmtypes.NewValidator(ourPubTm, power)})

			// --- fabricate empty block H+1 signed by OUR key ---
			makeState := stateH.Copy()
			makeState.Validators = ourSet.Copy()
			makeState.NextValidators = ourSet.Copy()
			block, partSet := makeState.MakeBlock(H+1, nil, commitH, nil, ourAddr)
			blockID := tmtypes.BlockID{Hash: block.Hash(), PartSetHeader: partSet.Header()}
			// The precommit timestamp becomes the BFT median time of the LastCommit used
			// to build H+2, which must be STRICTLY greater than H+1's block time.
			ts := block.Time.Add(time.Second)

			vote := &tmtypes.Vote{
				Type:             tmproto.PrecommitType,
				Height:           H + 1,
				Round:            0,
				BlockID:          blockID,
				Timestamp:        ts,
				ValidatorAddress: ourAddr,
				ValidatorIndex:   0,
			}
			sig, err := ourPriv.Sign(tmtypes.VoteSignBytes(chainID, vote.ToProto()))
			if err != nil {
				return fmt.Errorf("sign H+1 precommit: %w", err)
			}
			seenCommit := &tmtypes.Commit{
				Height:     H + 1,
				Round:      0,
				BlockID:    blockID,
				Signatures: []tmtypes.CommitSig{tmtypes.NewCommitSigForBlock(sig, ourAddr, ts)},
			}
			blockStore.SaveBlock(block, partSet, seenCommit)
			fmt.Fprintf(out, "saved fabricated block H+1=%d id=%X\n", H+1, blockID.Hash)

			// --- bootstrap CometBFT state to H+1 ---
			s1 := stateH.Copy()
			s1.LastBlockHeight = H + 1
			s1.LastBlockID = blockID
			s1.LastBlockTime = block.Time
			s1.LastValidators = ourSet.Copy()
			s1.Validators = ourSet.Copy()
			s1.NextValidators = ourSet.Copy()
			s1.LastHeightValidatorsChanged = H + 1
			s1.AppHash = appHashNew
			s1.LastResultsHash = tmtypes.NewResults(nil).Hash()
			if err := stateStore.Bootstrap(s1); err != nil {
				return fmt.Errorf("bootstrap cometbft state: %w", err)
			}
			fmt.Fprintf(out, "bootstrapped cometbft state to H+1=%d\n", H+1)

			// --- reset WAL + priv validator state ---
			if err := os.RemoveAll(filepath.Join(dataDir, "cs.wal")); err != nil {
				return fmt.Errorf("remove cs.wal: %w", err)
			}
			_ = os.WriteFile(filepath.Join(dataDir, "priv_validator_state.json"),
				[]byte("{\n  \"height\": \"0\",\n  \"round\": 0,\n  \"step\": 0\n}\n"), 0600)

			fmt.Fprintf(out, "DONE: fork ready at height %d. start: dscd start --home %s\n", H+1, home)
			return nil
		},
	}
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory (COPY only!)")
	return cmd
}

// forkMutateAppState swaps a bonded+online validator's consensus pubkey to ourPub,
// makes it the sole entry in the power index / last-validator-powers, commits, and
// returns the new app hash, that validator's consensus power, and its operator addr.
func forkMutateAppState(logger log.Logger, home string, appOpts servertypes.AppOptions, encCfg params.EncodingConfig, H int64, chainID string, ourPub []byte) ([]byte, int64, string, error) {
	db, err := openRedenomAppDB(home, appOpts)
	if err != nil {
		return nil, 0, "", err
	}
	defer db.Close()

	dscApp := app.NewDSC(logger, db, nil, true, map[int64]bool{}, home, uint(1), encCfg, appOpts)
	if got := dscApp.LastBlockHeight(); got != H {
		return nil, 0, "", fmt.Errorf("app height %d != cometbft height %d", got, H)
	}

	cms := dscApp.CommitMultiStore()
	ctx := sdk.NewContext(cms, tmproto.Header{Height: H, ChainID: chainID}, false, logger)
	k := dscApp.ValidatorKeeper
	store := ctx.KVStore(dscApp.GetKey(validatortypes.StoreKey))

	// pick V: bonded + online + highest power
	allVals := k.GetAllValidators(ctx)
	var V validatortypes.Validator
	found := false
	for _, v := range allVals {
		if v.Status == validatortypes.BondStatus_Bonded && v.Online && v.Stake > 0 {
			if !found || v.Stake > V.Stake {
				V = v
				found = true
			}
		}
	}
	if !found {
		return nil, 0, "", fmt.Errorf("no bonded+online validator with power > 0")
	}
	oldConsAddr, err := V.GetConsAddr()
	if err != nil {
		return nil, 0, "", fmt.Errorf("old cons addr: %w", err)
	}

	// wipe the entire power index and last-validator-powers, then re-add only V
	wipePrefix(store, validatortypes.GetValidatorsByPowerIndexKey())
	wipePrefix(store, validatortypes.GetLastValidatorPowersKey())
	store.Delete(validatortypes.GetValidatorByConsAddrIndexKey(oldConsAddr))

	pkAny, err := codectypes.NewAnyWithValue(&sdked25519.PubKey{Key: ourPub})
	if err != nil {
		return nil, 0, "", fmt.Errorf("pack pubkey: %w", err)
	}
	power := V.Stake
	V.ConsensusPubkey = pkAny
	V.Online = true
	V.Jailed = false
	V.Status = validatortypes.BondStatus_Bonded
	V.Stake = power

	k.SetValidator(ctx, V)
	k.SetValidatorRS(ctx, V.GetOperator(), validatortypes.ValidatorRS{
		Rewards: V.Rewards, TotalRewards: V.TotalRewards, Stake: power,
	})
	if err := k.SetValidatorByConsAddr(ctx, V); err != nil {
		return nil, 0, "", fmt.Errorf("set cons addr index: %w", err)
	}
	k.SetValidatorByPowerIndex(ctx, V)
	k.SetLastValidatorPower(ctx, V.GetOperator(), power)
	k.SetLastTotalPower(ctx, sdk.NewInt(power))

	newConsAddr, err := V.GetConsAddr()
	if err != nil {
		return nil, 0, "", fmt.Errorf("new cons addr: %w", err)
	}
	k.SetStartHeight(ctx, newConsAddr, H)

	// Force every OTHER validator offline. They keep their Bonded status + delegations
	// (so no pool-invariant breakage) but Online=false guarantees that even if a later
	// PayRewards/price recalc re-inserts them into the power index, the EndBlocker will
	// only ever transition them Bonded->Unbonded (zero power, no ABCI update) and never
	// emit a positive-power update that CometBFT would add to the (unsignable) valset.
	for _, v := range allVals {
		if v.OperatorAddress == V.OperatorAddress || !v.Online {
			continue
		}
		v.Online = false
		k.SetValidator(ctx, v)
	}

	commitID := cms.Commit()
	return commitID.Hash, power, V.OperatorAddress, nil
}

func wipePrefix(store sdk.KVStore, prefix []byte) {
	it := sdk.KVStorePrefixIterator(store, prefix)
	var keys [][]byte
	for ; it.Valid(); it.Next() {
		keys = append(keys, append([]byte(nil), it.Key()...))
	}
	it.Close()
	for _, key := range keys {
		store.Delete(key)
	}
}

// loadForkPrivKey reads a CometBFT priv_validator_key.json.
func loadForkPrivKey(path string) (tmed25519.PrivKey, tmed25519.PubKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var doc struct {
		PrivKey struct {
			Value string `json:"value"`
		} `json:"priv_key"`
		PubKey struct {
			Value string `json:"value"`
		} `json:"pub_key"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, nil, err
	}
	privBz, err := base64.StdEncoding.DecodeString(doc.PrivKey.Value)
	if err != nil {
		return nil, nil, fmt.Errorf("decode priv: %w", err)
	}
	pubBz, err := base64.StdEncoding.DecodeString(doc.PubKey.Value)
	if err != nil {
		return nil, nil, fmt.Errorf("decode pub: %w", err)
	}
	if len(privBz) != 64 || len(pubBz) != 32 {
		return nil, nil, fmt.Errorf("unexpected key sizes priv=%d pub=%d", len(privBz), len(pubBz))
	}
	return tmed25519.PrivKey(privBz), tmed25519.PubKey(pubBz), nil
}
