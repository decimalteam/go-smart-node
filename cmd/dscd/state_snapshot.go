package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	cometbftdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	cmtstore "github.com/tendermint/tendermint/store"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
)

// StateSnapshotCmd writes a full JSON snapshot of node delegation state and the
// delegation contract (stakes, frozen queue, auto-unbond queue, token balances) of an
// offline node home, for before/after comparisons. Read-only.
func StateSnapshotCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state-snapshot",
		Short: "Write a JSON snapshot of node + delegation-contract state (offline, read-only)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			homePath := cast.ToString(appOpts.Get(flags.FlagHome))
			outPath, _ := cmd.Flags().GetString("out")
			db, err := openRedenomAppDB(homePath, appOpts)
			if err != nil {
				return err
			}
			defer db.Close()
			dscApp := app.NewDSC(serverCtx.Logger, db, nil, true, map[int64]bool{}, homePath, uint(1), encCfg, appOpts)
			height := dscApp.LastBlockHeight()
			unix := time.Now().Unix()
			if blockDB, err := cometbftdb.NewGoLevelDB("blockstore", filepath.Join(homePath, "data")); err == nil {
				bs := cmtstore.NewBlockStore(blockDB)
				if meta := bs.LoadBlockMeta(bs.Height()); meta != nil {
					unix = meta.Header.Time.Unix()
				}
				blockDB.Close()
			}
			header := tmproto.Header{Height: height, ChainID: readChainID(homePath), Time: time.Unix(unix, 0).UTC()}
			ctx := sdk.NewContext(dscApp.CommitMultiStore().CacheMultiStore(), header, false, serverCtx.Logger)
			for _, v := range dscApp.ValidatorKeeper.GetAllValidators(ctx) {
				if v.IsBonded() {
					if cons, err := v.GetConsAddr(); err == nil {
						header.ProposerAddress = cons
						break
					}
				}
			}
			ctx = ctx.WithBlockHeader(header)
			snap, err := app.TakeStateSnapshot(ctx, dscApp)
			if err != nil {
				return err
			}
			bz, err := json.Marshal(snap)
			if err != nil {
				return err
			}
			if err := os.WriteFile(outPath, bz, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "snapshot height %d: %d validators, %d delegations, %d contract stakes, queue %d -> %s\n",
				height, len(snap.Validators), len(snap.Delegations), len(snap.Contract.Stakes), snap.Contract.QueueLen, outPath)
			return nil
		},
	}
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "offline node home")
	cmd.Flags().String("out", "snapshot.json", "output file")
	return cmd
}
