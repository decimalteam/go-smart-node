// One-off audit probe (fork rehearsal 2026-07-06): measure _nftBackedAmounts
// occupancy in the DecimalDelegation contract on an offline snapshot — the R1
// finding (docs/redenom-full-audit-2026-07-02.md) is only live where a coin
// stake has NFTBacked > 0. Read-only.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/decimalteam/ethermint/encoding"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/spf13/viper"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
)

func main() {
	home := os.Args[1]
	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, filepath.Join(home, "data"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	doc, err := tmtypes.GenesisDocFromFile(filepath.Join(home, "config", "genesis.json"))
	if err != nil {
		panic(err)
	}

	encCfg := encoding.MakeConfig(app.ModuleBasics)
	v := viper.New()
	v.Set("home", home)
	dscApp := app.NewDSC(log.NewNopLogger(), db, nil, true, map[int64]bool{}, home, uint(1), encCfg, v)
	height := dscApp.LastBlockHeight()
	ctx := sdk.NewContext(dscApp.CommitMultiStore().CacheMultiStore(),
		tmproto.Header{Height: height, ChainID: doc.ChainID}, false, log.NewNopLogger())

	kv := ctx.KVStore(dscApp.GetKey(evmtypes.StoreKey))
	readStorage := func(addr common.Address) stakescan.Storage {
		ps := prefix.NewStore(kv, evmtypes.AddressStoragePrefix(addr))
		it := ps.Iterator(nil, nil)
		defer it.Close()
		s := stakescan.Storage{}
		for ; it.Valid(); it.Next() {
			s[common.BytesToHash(it.Key())] = common.BytesToHash(it.Value())
		}
		return s
	}

	ccAddr := common.HexToAddress("0xc108715A06f76CAA96fa2c943Ebf05159c29A87D") // mainnet ContractCenter
	delegationAddr := stakescan.ResolveAddressBySymbol(readStorage(ccAddr), "delegation")
	fmt.Printf("height=%d delegation=%s\n", height, delegationAddr.Hex())

	res := stakescan.Reconstruct(readStorage(delegationAddr), stakescan.DelegationBase)
	total, nonzero := 0, 0
	for _, cs := range res.CoinStakes {
		total++
		if cs.NFTBacked != nil && cs.NFTBacked.Sign() > 0 {
			nonzero++
			fmt.Printf("NFTBacked>0: stakeID=%x tokenType=%d amount=%s nftBacked=%s delegator=%s validator=%s\n",
				cs.StakeID, cs.TokenType, cs.Amount, cs.NFTBacked, cs.Delegator.Hex(), cs.Validator.Hex())
		}
	}
	fmt.Printf("coinStakes=%d nftBackedNonZero=%d\n", total, nonzero)
}
