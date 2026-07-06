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
	"github.com/ethereum/go-ethereum/common"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	cmtstate "github.com/tendermint/tendermint/state"
	cmtstore "github.com/tendermint/tendermint/store"
	tmtypes "github.com/tendermint/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
)

// OpenZeppelin v5 OwnableUpgradeable ERC-7201 storage location:
//   keccak256(abi.encode(uint256(keccak256("openzeppelin.storage.Ownable")) - 1)) & ~0xff
// _owner is the first (and only) field of the struct, so it sits at exactly this slot.
const ozOwnableSlot = "0x9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300"

// ForkSetOwnerCmd overwrites the OpenZeppelin-v5 Ownable `_owner` storage slot of a
// contract (e.g. the delegation proxy) directly in an offline fork COPY's EVM state,
// then advances BOTH the app and CometBFT one height so the node boots consistently
// (app==store==state) — the exact fork-fund recipe, only the app mutation differs.
//
// Use it so a throwaway fork key can run owner-gated flows (pause/upgrade/...) on a
// mainnet fork WITHOUT ever touching the real owner's private key.
//
// The node MUST be stopped (LevelDB single-writer). Work only in a disposable copy.
func ForkSetOwnerCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fork-setowner",
		Short: "Overwrite a contract's OZ-Ownable owner slot in an offline fork COPY (app + CometBFT -> H+1)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			logger := serverCtx.Logger
			home := cast.ToString(appOpts.Get(flags.FlagHome))
			if home == "" {
				return fmt.Errorf("--home is required")
			}
			out := cmd.OutOrStdout()
			contractStr, _ := cmd.Flags().GetString("contract")
			ownerStr, _ := cmd.Flags().GetString("owner")
			slotStr, _ := cmd.Flags().GetString("slot")
			if !common.IsHexAddress(contractStr) {
				return fmt.Errorf("invalid --contract %q (expected 0x-hex eth address)", contractStr)
			}
			if !common.IsHexAddress(ownerStr) {
				return fmt.Errorf("invalid --owner %q (expected 0x-hex eth address)", ownerStr)
			}
			contractAddr := common.HexToAddress(contractStr)
			newOwner := common.HexToAddress(ownerStr)
			slot := common.HexToHash(slotStr)
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
				// normal path: overwrite the owner slot then commit -> Hc+1
				cms := dscApp.CommitMultiStore()
				ctx := sdk.NewContext(cms, tmproto.Header{Height: appH, ChainID: chainID}, false, logger)
				before := dscApp.EvmKeeper.GetState(ctx, contractAddr, slot)
				// 32-byte, left-padded value (EVM word holding the 20-byte address)
				value := common.BytesToHash(newOwner.Bytes())
				dscApp.EvmKeeper.SetState(ctx, contractAddr, slot, value.Bytes())
				after := dscApp.EvmKeeper.GetState(ctx, contractAddr, slot)
				commitID := cms.Commit()
				appHashNew = commitID.Hash
				fmt.Fprintf(out, "contract %s slot %s\n", contractAddr.Hex(), slot.Hex())
				fmt.Fprintf(out, "  owner %s -> %s\n", common.BytesToAddress(before.Bytes()).Hex(), common.BytesToAddress(after.Bytes()).Hex())
				fmt.Fprintf(out, "  app committed %d -> %d appHash %X\n", appH, appH+1, appHashNew)
				if common.BytesToAddress(after.Bytes()) != newOwner {
					return fmt.Errorf("post-write readback mismatch: got %s want %s", common.BytesToAddress(after.Bytes()).Hex(), newOwner.Hex())
				}
			case appH == Hc+1:
				// recovery: app already committed one ahead; only advance comet
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
	cmd.Flags().String("contract", "", "0x eth address of the contract (proxy) whose owner slot to overwrite")
	cmd.Flags().String("owner", "", "0x eth address to install as the new owner")
	cmd.Flags().String("slot", ozOwnableSlot, "storage slot holding _owner (default: OZ v5 Ownable ERC-7201 slot)")
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory (COPY only!)")
	return cmd
}
