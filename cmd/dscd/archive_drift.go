package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	sdkmath "cosmossdk.io/math"
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
	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/remotestore"
)

// ArchiveDriftCmd takes today's repair plan from an offline node home and, for every
// key it would change, reads the node delegation and the key's contract stake amounts
// at a historical height from an archive node (versioned per-key reads only), to tell
// incident damage from drift that already existed then.
func ArchiveDriftCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive-drift",
		Short: "Classify today's repair keys by whether they already differed at a historical height (read-only)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			homePath := cast.ToString(appOpts.Get(flags.FlagHome))
			rpc, _ := cmd.Flags().GetString("archive-rpc")
			height, _ := cmd.Flags().GetInt64("height")
			workers, _ := cmd.Flags().GetInt("workers")
			outPath, _ := cmd.Flags().GetString("out")
			out := cmd.OutOrStdout()

			db, err := openRedenomAppDB(homePath, appOpts)
			if err != nil {
				return fmt.Errorf("open application db: %w", err)
			}
			defer db.Close()
			dscApp := app.NewDSC(serverCtx.Logger, db, nil, true, map[int64]bool{}, homePath, uint(1), encCfg, appOpts)
			local := dscApp.LastBlockHeight()
			unix := time.Now().Unix()
			if blockDB, err := cometbftdb.NewGoLevelDB("blockstore", filepath.Join(homePath, "data")); err == nil {
				bs := cmtstore.NewBlockStore(blockDB)
				if meta := bs.LoadBlockMeta(bs.Height()); meta != nil {
					unix = meta.Header.Time.Unix()
				}
				blockDB.Close()
			}
			header := tmproto.Header{Height: local + 1, ChainID: readChainID(homePath), Time: time.Unix(unix, 0).UTC()}
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

			client := &remotestore.Client{RPC: rpc, Height: height, HTTP: &http.Client{Timeout: 2 * time.Minute}, Progress: os.Stderr}
			rows, err := app.ArchiveDrift(ctx, dscApp, client.KeysBatch, workers)
			if err != nil {
				return err
			}
			type agg struct {
				n    int
				base sdkmath.Int
				dels map[string]bool
			}
			by := map[string]*agg{}
			for _, r := range rows {
				a, ok := by[r.Class]
				if !ok {
					a = &agg{base: sdkmath.ZeroInt(), dels: map[string]bool{}}
					by[r.Class] = a
				}
				a.n++
				b, _ := sdkmath.NewIntFromString(r.DriftNowBase)
				a.base = a.base.Add(b)
				a.dels[r.Delegator] = true
			}
			fmt.Fprintf(out, "local height %d, historical height %d, %d keys, %d archive reads\n", local, height, len(rows), client.Calls)
			for _, c := range []string{app.DriftIncident, app.DriftOldUnchanged, app.DriftOldNodeMoved, app.DriftOldContractMove} {
				if a := by[c]; a != nil {
					fmt.Fprintf(out, "  %-19s keys=%-6d delegators=%-5d |drift| %s base coin\n", c, a.n, len(a.dels), sdkmath.LegacyNewDecFromIntWithPrec(a.base, 18).String())
				}
			}
			if outPath != "" {
				bz, _ := json.MarshalIndent(rows, "", " ")
				if err := os.WriteFile(outPath, bz, 0o644); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "offline node home with TODAY's state")
	cmd.Flags().String("archive-rpc", "https://dscd-node-archive-2.decimalchain.com/rpc", "archive node Tendermint RPC")
	cmd.Flags().Int64("height", 32811495, "historical height (last block before the 32811496 upgrade)")
	cmd.Flags().Int("workers", 8, "parallel archive batch requests (100 keys each)")
	cmd.Flags().String("out", "", "write every key as JSON")
	return cmd
}
