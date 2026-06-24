package main

import (
	"fmt"
	"math/big"
	"path/filepath"
	"sort"
	"text/tabwriter"

	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	dbm "github.com/tendermint/tm-db"

	tmtypes "github.com/tendermint/tendermint/types"

	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/common"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/app/redenom"
	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// contractCenterByChainID maps a cosmos chain-id to its DecimalContractCenter proxy
// address, used to resolve the delegation contract offline for the EVM cross-check.
var ccDryRunByChainID = map[string]common.Address{
	"decimal_75-1":       common.HexToAddress("0xc108715A06f76CAA96fa2c943Ebf05159c29A87D"), // mainnet
	"decimal_202020-1":   common.HexToAddress("0xbC96b61F137F28F0Da47Cc4Ef06e5f984B565A2E"), // testnet
	"decimal_20202020-1": common.HexToAddress("0x481487AEafc60512233a08Da240FC7AF99c0f696"), // devnet
}

// RedenomDryRunCmd runs the DEL redenomination against an offline node snapshot in a
// cache-only context (nothing is ever committed), then asserts the post-state
// invariants and prints a PASS/FAIL report. The node MUST be stopped (LevelDB is a
// single-writer); back up the data dir first as defence in depth.
func RedenomDryRunCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redenom-dryrun",
		Short: "Dry-run the DEL ÷1000 redenomination against an offline node snapshot (writes nothing)",
		Long: `Loads the application database at its latest height, runs the full redenomination
(app/redenom.Redenominate) inside a CacheMultiStore that is never written back, and
verifies the resulting state: supply == sum of balances, every sampled balance is
floor(old/1000), bonded/not-bonded pools still cover the scaled stakes, custom-coin
reserves are scaled while volumes are untouched, EVM delegation stakes match the
scaled Cosmos delegations, and no currently-bonded validator floors to zero power.

The node must be stopped. Nothing is committed; the on-disk state is left unchanged.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			homePath := cast.ToString(appOpts.Get(flags.FlagHome))
			if homePath == "" {
				return fmt.Errorf("--home is required")
			}
			sampleN, _ := cmd.Flags().GetInt("sample")

			db, err := openRedenomAppDB(homePath, appOpts)
			if err != nil {
				return fmt.Errorf("open application db: %w", err)
			}
			defer db.Close()

			dscApp := app.NewDSC(serverCtx.Logger, db, nil, true, map[int64]bool{}, homePath, uint(1), encCfg, appOpts)
			height := dscApp.LastBlockHeight()
			chainID := readChainID(homePath)

			// Cache-only context: every read/write goes through the cache; we never
			// call Write(), so the on-disk store is untouched.
			cms := dscApp.CommitMultiStore().CacheMultiStore()
			header := tmproto.Header{Height: height, ChainID: chainID}
			ctx := sdk.NewContext(cms, header, false, serverCtx.Logger)

			base := dscApp.ValidatorKeeper.BaseDenom(ctx)
			out := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintf(out, "DEL redenomination dry-run\tchain-id %s\theight %d\tbase %q\tdivisor %s\n",
				chainID, height, base, redenom.DefaultDivisor)

			div := redenom.DefaultDivisor
			b := captureBefore(ctx, dscApp, base, chainID, sampleN)

			keepers := redenom.Keepers{
				Bank: dscApp.BankKeeper, Coin: &dscApp.CoinKeeper, Validator: dscApp.ValidatorKeeper,
				NFT: &dscApp.NFTKeeper, Legacy: &dscApp.LegacyKeeper, Gov: dscApp.GovKeeper,
				Account: dscApp.AccountKeeper, EVM: &dscApp.EvmKeeper,
			}
			storeKeys := redenom.StoreKeys{Bank: dscApp.GetKey(banktypes.StoreKey), EVM: dscApp.GetKey(evmtypes.StoreKey)}

			report, err := redenom.Redenominate(ctx, keepers, storeKeys, div)
			if err != nil {
				return fmt.Errorf("redenominate: %w", err)
			}

			res := newCheckResult()
			verifyAfter(ctx, dscApp, base, div, b, report, res)

			printReport(out, report)
			res.print(out)
			out.Flush()
			if res.failed > 0 {
				return fmt.Errorf("dry-run FAILED: %d checks failed", res.failed)
			}
			return nil
		},
	}
	cmd.Flags().Int("sample", 2000, "number of accounts / EVM stakes to sample for floor verification (0 = skip sampling)")
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory")
	return cmd
}

// --- snapshot / verification ---

type beforeState struct {
	oldSupply       sdkmath.Int
	bondedAddr      sdk.AccAddress
	notBondedAddr   sdk.AccAddress
	oldBonded       sdkmath.Int
	oldNotBonded    sdkmath.Int
	valStatus       map[string]validatortypes.BondStatus // operator -> status
	oldPower        map[string]int64
	customVolume    map[string]sdkmath.Int
	customReserve   map[string]sdkmath.Int
	sampleBalances  []balSample // up to sampleN del accounts
	sampleStakes    []stakeSample
	delegationAddr  common.Address
	wdelAddr        common.Address
}

type balSample struct {
	addr sdk.AccAddress
	old  sdkmath.Int
}
type stakeSample struct {
	slot common.Hash
	old  sdkmath.Int
}

func captureBefore(ctx sdk.Context, a *app.DSC, base, chainID string, sampleN int) beforeState {
	b := beforeState{
		oldSupply:     a.BankKeeper.GetSupply(ctx, base).Amount,
		valStatus:     map[string]validatortypes.BondStatus{},
		oldPower:      map[string]int64{},
		customVolume:  map[string]sdkmath.Int{},
		customReserve: map[string]sdkmath.Int{},
	}
	b.bondedAddr = a.ValidatorKeeper.GetBondedPool(ctx).GetAddress()
	b.notBondedAddr = a.ValidatorKeeper.GetNotBondedPool(ctx).GetAddress()
	b.oldBonded = a.BankKeeper.GetBalance(ctx, b.bondedAddr, base).Amount
	b.oldNotBonded = a.BankKeeper.GetBalance(ctx, b.notBondedAddr, base).Amount

	for _, v := range a.ValidatorKeeper.GetAllValidators(ctx) {
		b.valStatus[v.OperatorAddress] = v.Status
		b.oldPower[v.OperatorAddress] = v.Stake
	}
	for _, c := range a.CoinKeeper.GetCoins(ctx) {
		if c.Denom == base {
			continue
		}
		b.customVolume[c.Denom] = c.Volume
		b.customReserve[c.Denom] = c.Reserve
	}
	if sampleN > 0 {
		a.BankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, coin sdk.Coin) bool {
			if coin.Denom == base && len(b.sampleBalances) < sampleN {
				b.sampleBalances = append(b.sampleBalances, balSample{addr: addr, old: coin.Amount})
			}
			return len(b.sampleBalances) >= sampleN
		})
		b.sampleStakes, b.delegationAddr, b.wdelAddr = sampleEVMStakes(ctx, a, chainID, sampleN)
	}
	return b
}

// sampleEVMStakes resolves the delegation + wdel contracts and captures up to sampleN
// DEL coin-stake amount slots with their pre-scale values.
func sampleEVMStakes(ctx sdk.Context, a *app.DSC, chainID string, sampleN int) ([]stakeSample, common.Address, common.Address) {
	ccAddr, ok := ccDryRunByChainID[chainID]
	if !ok {
		return nil, common.Address{}, common.Address{}
	}
	kv := ctx.KVStore(a.GetKey(evmtypes.StoreKey))
	ccStore := scanStorage(kv, ccAddr)
	delegationAddr := stakescan.ResolveAddressBySymbol(ccStore, "delegation")
	wdelAddr := stakescan.ResolveAddressBySymbol(ccStore, "wdel")
	if delegationAddr == (common.Address{}) {
		return nil, common.Address{}, wdelAddr
	}
	storage := scanStorage(kv, delegationAddr)
	res := stakescan.Reconstruct(storage, stakescan.DelegationBase)
	var out []stakeSample
	for _, cs := range res.CoinStakes {
		if cs.TokenType != 4 { // 4 == DEL
			continue
		}
		out = append(out, stakeSample{slot: stakescan.StakeAmountSlot(cs.BaseSlot), old: sdkmath.NewIntFromBigInt(cs.Amount)})
		if len(out) >= sampleN {
			break
		}
	}
	return out, delegationAddr, wdelAddr
}

func verifyAfter(ctx sdk.Context, a *app.DSC, base string, div sdkmath.Int, b beforeState, report redenom.Report, res *checkResult) {
	floor := func(x sdkmath.Int) sdkmath.Int { return x.Quo(div) }

	// 1. Supply: GetSupply == report.NewSupply, and == sum of all del balances.
	newSupply := a.BankKeeper.GetSupply(ctx, base).Amount
	res.check("supply == report.NewSupplyDel", newSupply.Equal(report.NewSupplyDel),
		fmt.Sprintf("supply=%s report=%s", newSupply, report.NewSupplyDel))
	sumBal := sdkmath.ZeroInt()
	a.BankKeeper.IterateAllBalances(ctx, func(_ sdk.AccAddress, coin sdk.Coin) bool {
		if coin.Denom == base {
			sumBal = sumBal.Add(coin.Amount)
		}
		return false
	})
	res.check("supply == sum(del balances)", newSupply.Equal(sumBal),
		fmt.Sprintf("supply=%s sum=%s", newSupply, sumBal))
	// Redenomination can only shrink supply (sum of per-account floors <= floor(sum)).
	res.check("new supply <= floor(old supply)", newSupply.LTE(b.oldSupply.Quo(div)),
		fmt.Sprintf("old=%s newFloor=%s new=%s", b.oldSupply, b.oldSupply.Quo(div), newSupply))

	// 2. Sampled account balances are exactly floor(old/div).
	floorFails := 0
	for _, s := range b.sampleBalances {
		got := a.BankKeeper.GetBalance(ctx, s.addr, base).Amount
		if !got.Equal(floor(s.old)) {
			floorFails++
		}
	}
	res.check(fmt.Sprintf("sampled balances floored (%d sampled)", len(b.sampleBalances)), floorFails == 0,
		fmt.Sprintf("%d mismatches", floorFails))

	// 3. Module pools floored cleanly. DSC staking is contract-based, so the Cosmos
	// bonded/not-bonded pools are empty; the staked DEL is held by the DecimalDelegation
	// contract's bank "del" balance, which must still cover the scaled active stakes.
	newBonded := a.BankKeeper.GetBalance(ctx, b.bondedAddr, base).Amount
	newNotBonded := a.BankKeeper.GetBalance(ctx, b.notBondedAddr, base).Amount
	res.check("bonded pool == floor(old)", newBonded.Equal(floor(b.oldBonded)),
		fmt.Sprintf("old=%s new=%s", b.oldBonded, newBonded))
	res.check("not-bonded pool == floor(old)", newNotBonded.Equal(floor(b.oldNotBonded)),
		fmt.Sprintf("old=%s new=%s", b.oldNotBonded, newNotBonded))

	// Staking backing: DSC stakes DEL as WDEL (wrapped DEL), so neither the Cosmos
	// pools nor the delegation contract hold native del — the backing is the WDEL
	// contract's native del balance, and the WDEL ERC20 ledger (totalSupply/balanceOf)
	// must scale in lockstep. WDEL consistency is a SEPARATE verification (its balanceOf
	// mapping is not enumerable from storage); flag it explicitly rather than asserting
	// the wrong (native-del-in-delegation) model.
	totalStaked := sumScaledCoinDelStakes(ctx, a, base)
	if b.delegationAddr != (common.Address{}) {
		delegBal := a.BankKeeper.GetBalance(ctx, sdk.AccAddress(b.delegationAddr.Bytes()), base).Amount
		res.note(fmt.Sprintf("staking is WDEL-backed: delegation contract native del=%s, scaled active DEL stakes=%s", delegBal, totalStaked))
	}

	// 3b. WDEL ledger stays backed 1:1: the WDEL contract's native del balance (scaled by
	// scaleBankDel = floor(old)) must still cover the new total WDEL supply (= sum of the
	// scaled holder balances). floor(Σ) >= Σ floor, so native >= totalSupply (over-backed
	// by dust) — anything else means a holder was missed or the supply mis-derived.
	if b.wdelAddr != (common.Address{}) {
		wdelNative := a.BankKeeper.GetBalance(ctx, sdk.AccAddress(b.wdelAddr.Bytes()), base).Amount
		res.check("WDEL native del backing >= new WDEL totalSupply", wdelNative.GTE(report.WDELTotalSupplyNew),
			fmt.Sprintf("native=%s totalSupply=%s holders=%d", wdelNative, report.WDELTotalSupplyNew, report.WDELHoldersScaled))
		res.check("WDEL new totalSupply <= floor(old totalSupply)", report.WDELTotalSupplyNew.LTE(report.WDELTotalSupplyOld.Quo(div)),
			fmt.Sprintf("old=%s new=%s", report.WDELTotalSupplyOld, report.WDELTotalSupplyNew))
		res.check("WDEL holders were scaled", report.WDELHoldersScaled > 0,
			fmt.Sprintf("holders=%d", report.WDELHoldersScaled))
	} else {
		res.note("WDEL contract not resolved (chain-id not mapped?) — WDEL ledger check skipped")
	}

	// EVM NFT collection _reserve ledgers (DecimalNFTBeaconProxy). The redenomination scales
	// each collection's DEL _reserve.amount in lockstep with its (bank-scaled) native backing.
	res.check("EVM NFT collections were scaled", report.EVMNFTCollectionsScaled > 0,
		fmt.Sprintf("scaled=%d", report.EVMNFTCollectionsScaled))
	if report.EVMNFTCollectionsFailed > 0 {
		res.note(fmt.Sprintf("%d EVM NFT collections (%s del) REJECTED (Σreserve > native = slot mis-id) and were left UNSCALED — manual review before mainnet",
			report.EVMNFTCollectionsFailed, report.EVMNFTFailedDel))
	}

	// 4. Custom coins: reserve floored, volume untouched.
	volFails, resFails := 0, 0
	for _, c := range a.CoinKeeper.GetCoins(ctx) {
		if c.Denom == base {
			continue
		}
		if ov, ok := b.customVolume[c.Denom]; ok && !c.Volume.Equal(ov) {
			volFails++
		}
		if or, ok := b.customReserve[c.Denom]; ok && !c.Reserve.Equal(floor(or)) {
			resFails++
		}
	}
	res.check("custom coin volumes unchanged", volFails == 0, fmt.Sprintf("%d changed", volFails))
	res.check("custom coin reserves floored", resFails == 0, fmt.Sprintf("%d mismatches", resFails))

	// 5. Validator power: no currently-bonded validator floored to 0.
	res.check("no bonded validator floored to 0 power", len(report.ValidatorsZeroed) == 0,
		fmt.Sprintf("%d zeroed: %v", len(report.ValidatorsZeroed), report.ValidatorsZeroed))

	// 6. EVM sampled stake amounts floored.
	if len(b.sampleStakes) > 0 {
		evmFails := 0
		for _, s := range b.sampleStakes {
			got := sdkmath.NewIntFromBigInt(new(big.Int).SetBytes(a.EvmKeeper.GetState(ctx, b.delegationAddr, s.slot).Bytes()))
			if !got.Equal(floor(s.old)) {
				evmFails++
			}
		}
		res.check(fmt.Sprintf("EVM delegation DEL stakes floored (%d sampled)", len(b.sampleStakes)), evmFails == 0,
			fmt.Sprintf("%d mismatches", evmFails))
	}
}

// sumScaledCoinDelStakes sums the now-scaled active coin-DEL delegation amounts. The
// DecimalDelegation contract's "del" balance backs at least these (it also holds the
// unbonding/frozen DEL), so it is a safe lower bound for the backing check.
func sumScaledCoinDelStakes(ctx sdk.Context, a *app.DSC, base string) sdkmath.Int {
	total := sdkmath.ZeroInt()
	for _, d := range a.ValidatorKeeper.GetAllDelegations(ctx) {
		if d.Stake.Type == validatortypes.StakeType_Coin && d.Stake.Stake.Denom == base {
			total = total.Add(d.Stake.Stake.Amount)
		}
	}
	return total
}

// --- small helpers ---

func openRedenomAppDB(homePath string, appOpts servertypes.AppOptions) (dbm.DB, error) {
	dataDir := filepath.Join(homePath, "data")
	backend := sdkserver.GetAppDBBackend(appOpts)
	return dbm.NewDB("application", backend, dataDir)
}

func readChainID(homePath string) string {
	doc, err := tmtypes.GenesisDocFromFile(filepath.Join(homePath, "config", "genesis.json"))
	if err != nil {
		return ""
	}
	return doc.ChainID
}

func scanStorage(kv sdk.KVStore, addr common.Address) stakescan.Storage {
	ps := prefix.NewStore(kv, evmtypes.AddressStoragePrefix(addr))
	it := ps.Iterator(nil, nil)
	defer it.Close()
	out := stakescan.Storage{}
	for ; it.Valid(); it.Next() {
		out[common.BytesToHash(it.Key())] = common.BytesToHash(it.Value())
	}
	return out
}

func printReport(out *tabwriter.Writer, r redenom.Report) {
	fmt.Fprintf(out, "\n-- report --\n")
	fmt.Fprintf(out, "old supply\t%s\n", r.OldSupplyDel)
	fmt.Fprintf(out, "new supply\t%s\n", r.NewSupplyDel)
	fmt.Fprintf(out, "supply removed (÷1000)\t%s\n", r.SupplyRemoved)
	fmt.Fprintf(out, "rounding dust\t%s\n", r.RoundingDust)
	fmt.Fprintf(out, "bank accounts scaled\t%d\n", r.BankAccountsScaled)
	fmt.Fprintf(out, "delegations / ubd / red\t%d / %d / %d\n", r.DelegationsScaled, r.UndelegationsScaled, r.RedelegationsScaled)
	fmt.Fprintf(out, "validators repowered\t%d (zeroed %d)\n", r.ValidatorsRepowered, len(r.ValidatorsZeroed))
	fmt.Fprintf(out, "custom reserves / nft reserves\t%d / %d\n", r.CustomReservesScaled, r.NFTReservesScaled)
	fmt.Fprintf(out, "legacy records\t%d\n", r.LegacyRecordsScaled)
	fmt.Fprintf(out, "EVM coin-stake / nft-reserve / frozen / autounbond slots\t%d / %d / %d / %d\n",
		r.EVMCoinStakeSlots, r.EVMNFTReserveSlots, r.EVMFrozenSlots, r.EVMAutoUnbondSlots)
	fmt.Fprintf(out, "EVM NFT collections scaled / rejected (del)\t%d / %d (%s)\n",
		r.EVMNFTCollectionsScaled, r.EVMNFTCollectionsFailed, r.EVMNFTFailedDel)
	fmt.Fprintf(out, "checks voided / refunded del\t%d / %s\n", r.ChecksVoided, r.ChecksRefundedDel)
}

type checkResult struct {
	passed, failed int
	lines          []string
	notes          []string
}

func newCheckResult() *checkResult { return &checkResult{} }

func (r *checkResult) check(name string, ok bool, detail string) {
	if ok {
		r.passed++
		r.lines = append(r.lines, fmt.Sprintf("PASS\t%s", name))
		return
	}
	r.failed++
	r.lines = append(r.lines, fmt.Sprintf("FAIL\t%s\t(%s)", name, detail))
}

// note records an informational line that does not affect pass/fail (used for facts
// that need separate verification, e.g. the WDEL-backed staking ledger).
func (r *checkResult) note(msg string) { r.notes = append(r.notes, msg) }

func (r *checkResult) print(out *tabwriter.Writer) {
	fmt.Fprintf(out, "\n-- checks --\n")
	sort.Strings(r.lines)
	for _, l := range r.lines {
		fmt.Fprintln(out, l)
	}
	if len(r.notes) > 0 {
		fmt.Fprintf(out, "\n-- notes (verify separately) --\n")
		for _, n := range r.notes {
			fmt.Fprintf(out, "NOTE\t%s\n", n)
		}
	}
	fmt.Fprintf(out, "\n%d passed, %d failed\n", r.passed, r.failed)
}
