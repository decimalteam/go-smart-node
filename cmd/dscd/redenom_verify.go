package main

import (
	"fmt"
	"math/big"
	"text/tabwriter"

	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/app/redenom"
	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// RedenomVerifyCmd is the EXHAUSTIVE counterpart to redenom-dryrun. Where the dry-run
// samples a few thousand accounts and checks high-level invariants, this command decodes
// EVERY DEL-denominated value in the offline snapshot (Cosmos module state + EVM contract
// storage), runs the real redenomination in a never-committed CacheMultiStore, then
// re-decodes every value and asserts new == floor(old/1000) for each one individually.
//
// It also does a full-storage diff of the three most-rewritten EVM contracts (delegation,
// WDEL, checks): every storage slot that changed must be a clean ÷1000 floor or a void
// (deletion), proving the upgrade neither corrupted an unrelated slot nor missed a DEL slot.
//
// The node MUST be stopped (LevelDB is single-writer); nothing is committed.
func RedenomVerifyCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redenom-verify",
		Short: "Exhaustively verify the DEL ÷1000 redenomination against an offline snapshot (every value, writes nothing)",
		Long: `Decodes every DEL-denominated value in the application database (all bank balances,
delegations/undelegations/redelegations, validator rewards, coin reserves, NFT and
sub-token reserves, legacy records, threshold params, and all EVM delegation/WDEL/NFT-
collection/checks DEL slots), runs the full redenomination in a cache-only context, and
verifies each value individually became floor(old/1000). Prints per-category counts of
values checked and any mismatches. The node must be stopped; nothing is committed.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			homePath := cast.ToString(appOpts.Get(flags.FlagHome))
			if homePath == "" {
				return fmt.Errorf("--home is required")
			}

			db, err := openRedenomAppDB(homePath, appOpts)
			if err != nil {
				return fmt.Errorf("open application db: %w", err)
			}
			defer db.Close()

			dscApp := app.NewDSC(serverCtx.Logger, db, nil, true, map[int64]bool{}, homePath, uint(1), encCfg, appOpts)
			height := dscApp.LastBlockHeight()
			chainID := readChainID(homePath)

			cms := dscApp.CommitMultiStore().CacheMultiStore()
			ctx := sdk.NewContext(cms, tmproto.Header{Height: height, ChainID: chainID}, false, serverCtx.Logger)

			base := dscApp.ValidatorKeeper.BaseDenom(ctx)
			div := redenom.DefaultDivisor
			out := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintf(out, "DEL redenomination EXHAUSTIVE verify\tchain-id %s\theight %d\tbase %q\tdivisor %s\n\n",
				chainID, height, base, div)

			res := newCheckResult()

			// 1. Capture every DEL value BEFORE.
			before := captureAll(ctx, dscApp, base, chainID, div)
			fmt.Fprintf(out, "decoded BEFORE: %s\n", before.summary())

			// 2. Run the real redenomination (same call the upgrade handler makes).
			keepers := redenom.Keepers{
				Bank: dscApp.BankKeeper, Coin: &dscApp.CoinKeeper, Validator: dscApp.ValidatorKeeper,
				NFT: &dscApp.NFTKeeper, Legacy: &dscApp.LegacyKeeper, Gov: dscApp.GovKeeper,
				Account: dscApp.AccountKeeper, EVM: &dscApp.EvmKeeper, Fee: &dscApp.FeeKeeper,
			}
			storeKeys := redenom.StoreKeys{Bank: dscApp.GetKey(banktypes.StoreKey), EVM: dscApp.GetKey(evmtypes.StoreKey)}
			report, err := redenom.Redenominate(ctx, keepers, storeKeys, div)
			if err != nil {
				return fmt.Errorf("redenominate: %w", err)
			}

			// 3. Verify every captured value individually AFTER.
			verifyAll(ctx, dscApp, base, div, before, report, res)

			printReport(out, report)
			res.print(out)
			out.Flush()
			if res.failed > 0 {
				return fmt.Errorf("verify FAILED: %d checks failed", res.failed)
			}
			return nil
		},
	}
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory")
	return cmd
}

// floorI returns floor(x/div) for a sdkmath.Int.
func floorI(x, div sdkmath.Int) sdkmath.Int {
	if x.IsNil() {
		return sdkmath.ZeroInt()
	}
	return x.Quo(div)
}

// evmSlotVal is one captured EVM storage slot with its pre-scale value.
type evmSlotVal struct {
	addr common.Address
	slot common.Hash
	old  *big.Int
}

// fullState is every DEL value decoded from the snapshot, keyed for after-lookup.
type fullState struct {
	// Cosmos
	bankDel    map[string]sdkmath.Int // addr.String() -> old del balance
	oldSupply  sdkmath.Int
	delegs     []amt    // delegation del stakes (positional; no deletions)
	undelegs   []amt    // undelegation entry del stakes
	redelegs   []amt    // redelegation entry del stakes
	delReserve sdkmath.Int
	delLimit   sdkmath.Int            // base coin LimitVolume (emission cap / percentForHold denominator; floored)
	custVol    map[string]sdkmath.Int // custom denom -> volume (must be unchanged)
	custRes    map[string]sdkmath.Int // custom denom -> reserve (floored)
	nftRes     []amt                  // nft token + sub-token del reserves
	legacy     []amt                  // legacy record del coin amounts
	valRewards []valRS                // per-validator rewards/totalRewards/power
	thrCoinVol sdkmath.Int            // coin params BaseVolume
	thrNftMin  sdkmath.Int            // nft params MinReserveAmount
	thrGovDep  []amt                  // gov MinDeposit del entries
	feePrice   sdk.Dec                // base-denom fiat oracle price (x/fee) — fee-layer root
	minGas     sdk.Dec                // derived EVM min gas price / base fee

	// EVM (independently re-identified DEL slots -> old value)
	evmDeleg    []evmSlotVal // delegation contract DEL slots (coin/nft/frozen/autounbond)
	evmWdel     []evmSlotVal // WDEL ledger DEL slots
	evmNft      []evmSlotVal // NFT-collection DEL reserve slots whose native survives scaling
	evmNftDrop  []evmSlotVal // reserve slots in collections whose native del floors to 0 (redenom drops these)
	nftDropColl int          // count of dropped collections
	nftDropDel  sdkmath.Int  // total (old) reserve dust left unscaled in dropped collections
	// Checks: slots that must be voided (==0) after; refund accounting.
	checksAddr     common.Address
	checkVoid      []common.Hash          // slots that should be zero after
	checkRefund    sdkmath.Int            // total expected refund (sum of floored DEL check amounts)
	checkRefundN   int                    // number of outstanding DEL checks
	refundByAddr   map[string]sdkmath.Int // creator AccAddress.String() -> floored refund credited

	// Full storage snapshots for corruption diff (delegation, wdel, checks).
	delegAddr  common.Address
	wdelAddr   common.Address
	delegStore stakescan.Storage
	wdelStore  stakescan.Storage
	checksStore stakescan.Storage
}

type amt struct {
	id  string
	old sdkmath.Int
}
type valRS struct {
	op            string
	rewards       sdkmath.Int
	totalRewards  sdkmath.Int
	power         int64
	bonded        bool
}

func (s *fullState) summary() string {
	return fmt.Sprintf("bankDel=%d delegs=%d undelegs=%d redelegs=%d custCoins=%d nftReserves=%d legacy=%d validators=%d | evmDeleg=%d evmWdel=%d evmNft=%d checksVoidSlots=%d",
		len(s.bankDel), len(s.delegs), len(s.undelegs), len(s.redelegs), len(s.custRes),
		len(s.nftRes), len(s.legacy), len(s.valRewards), len(s.evmDeleg), len(s.evmWdel), len(s.evmNft), len(s.checkVoid))
}

func captureAll(ctx sdk.Context, a *app.DSC, base, chainID string, div sdkmath.Int) *fullState {
	s := &fullState{
		bankDel: map[string]sdkmath.Int{},
		custVol: map[string]sdkmath.Int{},
		custRes: map[string]sdkmath.Int{},
	}

	// --- Bank: every del balance (deletions possible -> map keyed by address) ---
	s.oldSupply = a.BankKeeper.GetSupply(ctx, base).Amount
	a.BankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, c sdk.Coin) bool {
		if c.Denom == base {
			s.bankDel[addr.String()] = c.Amount
		}
		return false
	})

	// --- Validator stakes: delegations / undelegations / redelegations ---
	for _, d := range a.ValidatorKeeper.GetAllDelegations(ctx) {
		if d.Stake.Stake.Denom == base {
			id := fmt.Sprintf("%s|%s|%s", d.Delegator, d.Validator, d.Stake.ID)
			s.delegs = append(s.delegs, amt{id, d.Stake.Stake.Amount})
		}
	}
	a.ValidatorKeeper.IterateUndelegations(ctx, func(_ int64, ubd validatortypes.Undelegation) bool {
		for i, e := range ubd.Entries {
			if e.Stake.Stake.Denom == base {
				id := fmt.Sprintf("%s|%s|%d", ubd.Delegator, ubd.Validator, i)
				s.undelegs = append(s.undelegs, amt{id, e.Stake.Stake.Amount})
			}
		}
		return false
	})
	a.ValidatorKeeper.IterateRedelegations(ctx, func(_ int64, red validatortypes.Redelegation) bool {
		for i, e := range red.Entries {
			if e.Stake.Stake.Denom == base {
				id := fmt.Sprintf("%s|%s|%s|%d", red.Delegator, red.ValidatorSrc, red.ValidatorDst, i)
				s.redelegs = append(s.redelegs, amt{id, e.Stake.Stake.Amount})
			}
		}
		return false
	})

	// --- Coins: del reserve, custom volumes (unchanged) + reserves (floored) ---
	for _, c := range a.CoinKeeper.GetCoins(ctx) {
		if c.Denom == base {
			s.delReserve = c.Reserve
			s.delLimit = c.LimitVolume
			continue
		}
		s.custVol[c.Denom] = c.Volume
		s.custRes[c.Denom] = c.Reserve
	}

	// --- NFT token + sub-token DEL reserves ---
	for _, col := range a.NFTKeeper.GetCollections(ctx) {
		creator, err := sdk.AccAddressFromBech32(col.Creator)
		if err != nil {
			continue
		}
		for _, token := range a.NFTKeeper.GetTokens(ctx, creator, col.Denom) {
			if token.Reserve.Denom == base {
				s.nftRes = append(s.nftRes, amt{"tok|" + token.ID, token.Reserve.Amount})
			}
			for _, st := range a.NFTKeeper.GetSubTokens(ctx, token.ID) {
				if st.Reserve != nil && st.Reserve.Denom == base {
					s.nftRes = append(s.nftRes, amt{fmt.Sprintf("sub|%s|%d", token.ID, st.ID), st.Reserve.Amount})
				}
			}
		}
	}

	// --- Legacy records ---
	for _, rec := range a.LegacyKeeper.GetLegacyRecords(ctx) {
		for i, c := range rec.Coins {
			if c.Denom == base {
				s.legacy = append(s.legacy, amt{fmt.Sprintf("%s|%d", rec.LegacyAddress, i), c.Amount})
			}
		}
	}

	// --- Validator rewards + power ---
	for _, v := range a.ValidatorKeeper.GetAllValidators(ctx) {
		rs, err := a.ValidatorKeeper.GetValidatorRS(ctx, v.GetOperator())
		if err != nil {
			continue
		}
		s.valRewards = append(s.valRewards, valRS{
			op: v.OperatorAddress, rewards: rs.Rewards, totalRewards: rs.TotalRewards,
			power: rs.Stake, bonded: v.Status == validatortypes.BondStatus_Bonded,
		})
	}

	// --- Threshold params ---
	s.thrCoinVol = a.CoinKeeper.GetParams(ctx).BaseVolume
	s.thrNftMin = a.NFTKeeper.GetParams(ctx).MinReserveAmount
	for _, c := range a.GovKeeper.GetDepositParams(ctx).MinDeposit {
		if c.Denom == base {
			s.thrGovDep = append(s.thrGovDep, amt{c.Denom, c.Amount})
		}
	}
	if p, err := a.FeeKeeper.GetPrice(ctx, base, "usd"); err == nil {
		s.feePrice = p.Price
		s.minGas = a.FeeKeeper.GetMinGasPrice(ctx)
	}

	// --- EVM ---
	captureEVM(ctx, a, base, chainID, div, s)
	return s
}

// captureEVM re-identifies every DEL-denominated EVM storage slot independently (mirroring
// app/redenom/evm.go) and records its pre-scale value, plus full-storage snapshots of the
// delegation/WDEL/checks contracts for the corruption diff.
func captureEVM(ctx sdk.Context, a *app.DSC, base, chainID string, div sdkmath.Int, s *fullState) {
	cc, ok := redenom.ContractCenterFor(chainID)
	if !ok {
		return
	}
	evmKey := a.GetKey(evmtypes.StoreKey)
	kv := ctx.KVStore(evmKey)
	ccStore := scanStorage(kv, cc)
	delegAddr := stakescan.ResolveAddressBySymbol(ccStore, contracts.NameOfSlugForGetAddressDelegation)
	wdelAddr := stakescan.ResolveAddressBySymbol(ccStore, contracts.NameOfSlugForGetAddressWDEL)
	tokenCenter := stakescan.ResolveAddressBySymbol(ccStore, contracts.NameOfSlugForGetAddressTokenCenter)
	s.delegAddr, s.wdelAddr = delegAddr, wdelAddr

	// 1. Delegation contract DEL slots.
	if delegAddr != (common.Address{}) {
		st := scanStorage(kv, delegAddr)
		s.delegStore = st
		r := stakescan.Reconstruct(st, stakescan.DelegationBase)
		read := func(slot common.Hash) *big.Int { return new(big.Int).SetBytes(st[slot].Bytes()) }
		for _, cs := range r.CoinStakes {
			if cs.TokenType != 4 {
				continue
			}
			slot := stakescan.StakeAmountSlot(cs.BaseSlot)
			s.evmDeleg = append(s.evmDeleg, evmSlotVal{delegAddr, slot, read(slot)})
		}
		for _, ns := range r.NFTStakes {
			// Mirror the migration: only DEL-denominated NFT reserves (reserveToken == wdel)
			// are scaled; custom-coin reserves are left untouched and must not be checked.
			if ns.ReserveToken != wdelAddr {
				continue
			}
			slot := stakescan.NFTStakeReserveAmountSlot(ns.BaseSlot)
			s.evmDeleg = append(s.evmDeleg, evmSlotVal{delegAddr, slot, read(slot)})
		}
		for _, fz := range r.FrozenLive {
			if fz.Stake.TokenType != 4 {
				continue
			}
			slot := stakescan.FrozenStakeAmountSlot(stakescan.DelegationBase, fz.Index, false)
			s.evmDeleg = append(s.evmDeleg, evmSlotVal{delegAddr, slot, read(slot)})
		}
		for _, fz := range r.FrozenDeprecated {
			if fz.Stake.TokenType != 4 {
				continue
			}
			slot := stakescan.FrozenStakeAmountSlot(stakescan.DelegationBase, fz.Index, true)
			s.evmDeleg = append(s.evmDeleg, evmSlotVal{delegAddr, slot, read(slot)})
		}
		for _, q := range stakescan.FindAutoUnbondQueues(st, 1.0) {
			for _, e := range q.Entries {
				isDEL := e.Token == (common.Address{}) || (wdelAddr != (common.Address{}) && e.Token == wdelAddr)
				if !isDEL {
					continue
				}
				slot := stakescan.AutoUnbondAmountSlot(q.FieldSlot, e.Index)
				s.evmDeleg = append(s.evmDeleg, evmSlotVal{delegAddr, slot, read(slot)})
			}
		}
	}

	// 2. WDEL ledger: every non-meta, non-sentinel, non-zero slot.
	if wdelAddr != (common.Address{}) {
		st := scanStorage(kv, wdelAddr)
		s.wdelStore = st
		meta := map[common.Hash]bool{
			common.BigToHash(big.NewInt(0)): true,
			common.BigToHash(big.NewInt(1)): true,
			common.BigToHash(big.NewInt(2)): true,
		}
		inf := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
		for slot, val := range st {
			if meta[slot] {
				continue
			}
			v := new(big.Int).SetBytes(val.Bytes())
			if v.Sign() == 0 || v.Cmp(inf) == 0 {
				continue
			}
			s.evmWdel = append(s.evmWdel, evmSlotVal{wdelAddr, slot, v})
		}
	}

	// 3. EVM NFT-collection DEL reserve amount slots (mirror scaleNFTCollectionReserves).
	// IMPORTANT: scaleNFTCollectionReserves runs AFTER scaleBankDel, so it only sees
	// collections whose native del SURVIVED scaling (floor(native/div) > 0). A collection
	// whose native floors to 0 has its bank balance deleted and drops out of redenom's
	// iteration entirely — its reserve is left UNSCALED. We capture pre-scale (native still
	// present), so we split the two buckets to mirror redenom exactly and surface the drop.
	s.nftDropDel = sdkmath.ZeroInt()
	beaconSlot := common.HexToHash("0xa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d50")
	delType := big.NewInt(int64(1)) // ReserveType.DEL
	a.BankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, c sdk.Coin) bool {
		if c.Denom != base || len(addr.Bytes()) != 20 {
			return false
		}
		eth := common.BytesToAddress(addr.Bytes())
		beacon := a.EvmKeeper.GetState(ctx, eth, beaconSlot)
		if beacon == (common.Hash{}) || common.BytesToAddress(beacon.Bytes()) == tokenCenter {
			return false
		}
		survives := floorI(c.Amount, div).IsPositive() // does redenom still see this collection?
		st := scanStorage(kv, eth)
		dropped := false
		for slotHash, v := range st {
			if new(big.Int).SetBytes(v.Bytes()).Cmp(delType) != 0 {
				continue
			}
			sBig := new(big.Int).SetBytes(slotHash.Bytes())
			amountSlot := common.BigToHash(new(big.Int).Sub(sBig, big.NewInt(1)))
			tokenSlot := common.BigToHash(new(big.Int).Sub(sBig, big.NewInt(2)))
			if st[tokenSlot] != (common.Hash{}) {
				continue
			}
			av := new(big.Int).SetBytes(st[amountSlot].Bytes())
			if av.Sign() == 0 {
				continue
			}
			if survives {
				s.evmNft = append(s.evmNft, evmSlotVal{eth, amountSlot, av})
			} else {
				s.evmNftDrop = append(s.evmNftDrop, evmSlotVal{eth, amountSlot, av})
				s.nftDropDel = s.nftDropDel.Add(sdkmath.NewIntFromBigInt(av))
				dropped = true
			}
		}
		if dropped {
			s.nftDropColl++
		}
		return false
	})

	// 4. Checks: outstanding DEL checks to be voided + refunded.
	checksAddr := stakescan.ResolveAddressBySymbol(ccStore, "checks")
	var checks []stakescan.Check
	if checksAddr != (common.Address{}) {
		checks = stakescan.ReconstructChecks(scanStorage(kv, checksAddr), stakescan.ChecksBase)
	} else {
		checksAddr, checks = locateChecksContractMain(kv)
	}
	s.checksAddr = checksAddr
	s.refundByAddr = map[string]sdkmath.Int{}
	if checksAddr != (common.Address{}) {
		s.checksStore = scanStorage(kv, checksAddr)
		s.checkRefund = sdkmath.ZeroInt()
		for _, c := range checks {
			if c.Status != 0 || c.TypeChecks != 0 {
				continue
			}
			s.checkRefundN++
			refund := floorI(sdkmath.NewIntFromBigInt(c.Amount), div)
			s.checkRefund = s.checkRefund.Add(refund)
			creatorKey := sdk.AccAddress(c.Creator.Bytes()).String()
			if prev, ok := s.refundByAddr[creatorKey]; ok {
				s.refundByAddr[creatorKey] = prev.Add(refund)
			} else {
				s.refundByAddr[creatorKey] = refund
			}
			s.checkVoid = append(s.checkVoid, c.BaseSlot, stakescan.SlotAdd(c.BaseSlot, 1))
			detailsBase := stakescan.CheckDetailsBase(c.CheckDetailsHash)
			for i := uint64(0); i < 5; i++ {
				s.checkVoid = append(s.checkVoid, stakescan.SlotAdd(detailsBase, i))
			}
		}
	}
}

func verifyAll(ctx sdk.Context, a *app.DSC, base string, div sdkmath.Int, b *fullState, report redenom.Report, res *checkResult) {
	// --- Bank: every del balance floored; supply == sum == report ---
	// Expected = floor(old/div) for every account, PLUS the check-refund transfer applied by
	// processChecks AFTER scaleBankDel: each creator is credited its floored refund, and the
	// DecimalChecks contract is debited the total. Model that here so refunds aren't flagged.
	checksAccStr := ""
	if b.checksAddr != (common.Address{}) {
		checksAccStr = sdk.AccAddress(b.checksAddr.Bytes()).String()
	}
	bankMiss, bankFirst, refundAccounted := 0, "", 0
	for addrStr, old := range b.bankDel {
		addr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			continue
		}
		want := floorI(old, div)
		if r, ok := b.refundByAddr[addrStr]; ok {
			want = want.Add(r)
			refundAccounted++
		}
		if addrStr == checksAccStr {
			want = want.Sub(b.checkRefund)
			refundAccounted++
		}
		got := a.BankKeeper.GetBalance(ctx, addr, base).Amount
		if !got.Equal(want) {
			bankMiss++
			if bankFirst == "" {
				bankFirst = fmt.Sprintf("%s old=%s got=%s want=%s", addrStr, old, got, want)
			}
		}
	}
	res.check(fmt.Sprintf("ALL bank del balances floored (%d checked, %d incl. check-refund adj)", len(b.bankDel), refundAccounted), bankMiss == 0,
		fmt.Sprintf("%d mismatches; first: %s", bankMiss, bankFirst))

	newSupply := a.BankKeeper.GetSupply(ctx, base).Amount
	sum := sdkmath.ZeroInt()
	a.BankKeeper.IterateAllBalances(ctx, func(_ sdk.AccAddress, c sdk.Coin) bool {
		if c.Denom == base {
			sum = sum.Add(c.Amount)
		}
		return false
	})
	res.check("supply == report.NewSupplyDel", newSupply.Equal(report.NewSupplyDel), fmt.Sprintf("supply=%s report=%s", newSupply, report.NewSupplyDel))
	res.check("supply == sum(del balances)", newSupply.Equal(sum), fmt.Sprintf("supply=%s sum=%s", newSupply, sum))
	res.check("new supply <= floor(old supply)", newSupply.LTE(floorI(b.oldSupply, div)), fmt.Sprintf("new=%s floor(old)=%s", newSupply, floorI(b.oldSupply, div)))

	// --- Validator stakes (positional re-iteration; no deletions) ---
	verifyCosmosAmts(res, "delegation del stakes", b.delegs, collectDelegs(ctx, a, base))
	verifyCosmosAmts(res, "undelegation del entries", b.undelegs, collectUndelegs(ctx, a, base))
	verifyCosmosAmts(res, "redelegation del entries", b.redelegs, collectRedelegs(ctx, a, base))

	// --- Coins ---
	delVolOK, custResMiss, custVolMiss := true, 0, 0
	var newDelReserve sdkmath.Int = sdkmath.ZeroInt()
	var newDelLimit sdkmath.Int = sdkmath.ZeroInt()
	for _, c := range a.CoinKeeper.GetCoins(ctx) {
		if c.Denom == base {
			delVolOK = c.Volume.Equal(newSupply) // del coin volume is set to the new supply
			newDelReserve = c.Reserve
			newDelLimit = c.LimitVolume
			continue
		}
		if ov, ok := b.custVol[c.Denom]; ok && !c.Volume.Equal(ov) {
			custVolMiss++
		}
		if or, ok := b.custRes[c.Denom]; ok && !c.Reserve.Equal(floorI(or, div)) {
			custResMiss++
		}
	}
	res.check("del coin volume == new supply", delVolOK, "")
	res.check("del coin reserve floored", newDelReserve.Equal(floorI(b.delReserve, div)), fmt.Sprintf("old=%s new=%s", b.delReserve, newDelReserve))
	res.check("del coin LimitVolume floored", newDelLimit.Equal(floorI(b.delLimit, div)), fmt.Sprintf("old=%s new=%s", b.delLimit, newDelLimit))
	res.check(fmt.Sprintf("ALL custom coin reserves floored (%d)", len(b.custRes)), custResMiss == 0, fmt.Sprintf("%d mismatches", custResMiss))
	res.check(fmt.Sprintf("ALL custom coin volumes unchanged (%d)", len(b.custVol)), custVolMiss == 0, fmt.Sprintf("%d changed", custVolMiss))

	// --- NFT reserves ---
	verifyCosmosAmts(res, "NFT token+subtoken del reserves", b.nftRes, collectNFTRes(ctx, a, base))

	// --- Legacy ---
	verifyCosmosAmts(res, "legacy del coin amounts", b.legacy, collectLegacy(ctx, a, base))

	// --- Validator rewards + power ---
	rwMiss, zeroedBonded := 0, 0
	for _, v := range b.valRewards {
		op, err := sdk.ValAddressFromBech32(v.op)
		if err != nil {
			continue
		}
		rs, err := a.ValidatorKeeper.GetValidatorRS(ctx, op)
		if err != nil {
			continue
		}
		if !rs.Rewards.Equal(floorI(v.rewards, div)) || !rs.TotalRewards.Equal(floorI(v.totalRewards, div)) {
			rwMiss++
		}
		if v.bonded && rs.Stake == 0 {
			zeroedBonded++
		}
	}
	res.check(fmt.Sprintf("ALL validator rewards/totalRewards floored (%d)", len(b.valRewards)), rwMiss == 0, fmt.Sprintf("%d mismatches", rwMiss))
	res.check("no bonded validator floored to 0 power", zeroedBonded == 0 && len(report.ValidatorsZeroed) == 0,
		fmt.Sprintf("independent=%d report=%d", zeroedBonded, len(report.ValidatorsZeroed)))

	// --- Threshold params ---
	cp := a.CoinKeeper.GetParams(ctx)
	np := a.NFTKeeper.GetParams(ctx)
	res.check("coin BaseVolume floored", cp.BaseVolume.Equal(floorI(b.thrCoinVol, div)), fmt.Sprintf("old=%s new=%s", b.thrCoinVol, cp.BaseVolume))
	res.check("nft MinReserveAmount floored", np.MinReserveAmount.Equal(floorI(b.thrNftMin, div)), fmt.Sprintf("old=%s new=%s", b.thrNftMin, np.MinReserveAmount))
	govMiss := 0
	govDep := a.GovKeeper.GetDepositParams(ctx).MinDeposit
	for _, e := range b.thrGovDep {
		for _, c := range govDep {
			if c.Denom == e.id && !c.Amount.Equal(floorI(e.old, div)) {
				govMiss++
			}
		}
	}
	res.check(fmt.Sprintf("gov MinDeposit del floored (%d)", len(b.thrGovDep)), govMiss == 0, fmt.Sprintf("%d mismatches", govMiss))

	// Fee layer (F-9/F-10, fork rehearsal 2026-07-06): the base coin's fiat oracle
	// price ×div, and the derived EVM min gas price / base fee therefore ÷div.
	if b.feePrice.IsNil() {
		res.check("fee oracle base price ×div (fee layer)", false, "no base-denom fiat price found before the upgrade")
	} else {
		newPrice, ferr := a.FeeKeeper.GetPrice(ctx, base, "usd")
		res.check("fee oracle base price ×div (fee layer)",
			ferr == nil && newPrice.Price.Equal(b.feePrice.MulInt(div)),
			fmt.Sprintf("old=%s new=%v err=%v", b.feePrice, newPrice.Price, ferr))
		newMinGas := a.FeeKeeper.GetMinGasPrice(ctx)
		res.check("EVM min gas price / base fee ÷div", newMinGas.MulInt(div).Equal(b.minGas),
			fmt.Sprintf("old=%s new=%s", b.minGas, newMinGas))
	}

	// --- EVM: every re-identified DEL slot floored ---
	verifyEVMSlots(ctx, a, div, res, "EVM delegation DEL slots", b.evmDeleg)
	verifyEVMSlots(ctx, a, div, res, "EVM WDEL ledger slots", b.evmWdel)
	verifyEVMSlots(ctx, a, div, res, "EVM NFT-collection reserve slots (native survives)", b.evmNft)

	// Dropped collections: native del floored to 0 -> redenom never scales their reserve.
	// Confirm they are indeed left UNSCALED (matches redenom) and surface them as a finding:
	// each is now a tiny unbacked reserve (native 0, reserve.amount unscaled). Dust-level.
	if len(b.evmNftDrop) > 0 {
		stillUnscaled := 0
		for _, sv := range b.evmNftDrop {
			got := new(big.Int).SetBytes(a.EvmKeeper.GetState(ctx, sv.addr, sv.slot).Bytes())
			if got.Cmp(sv.old) == 0 {
				stillUnscaled++
			}
		}
		res.check(fmt.Sprintf("dropped-collection reserves left unscaled as expected (%d slots, %d colls)", len(b.evmNftDrop), b.nftDropColl),
			stillUnscaled == len(b.evmNftDrop), fmt.Sprintf("%d of %d unexpectedly changed", len(b.evmNftDrop)-stillUnscaled, len(b.evmNftDrop)))
		res.note(fmt.Sprintf("FINDING: %d EVM NFT collection(s) whose native del floored to 0 were dropped from reserve scaling — %d reserve slot(s), %s base-unit DEL total left unscaled+now-unbacked (dust; redeem of these would revert). Compare: redenom scaled %d slots, we identified %d total.",
			b.nftDropColl, len(b.evmNftDrop), b.nftDropDel, report.EVMNFTReserveSlots, len(b.evmNft)+len(b.evmNftDrop)))
	}

	// --- EVM checks: every voided slot is now empty; refund total matches ---
	if b.checksAddr != (common.Address{}) {
		voidMiss := 0
		for _, slot := range b.checkVoid {
			if a.EvmKeeper.GetState(ctx, b.checksAddr, slot) != (common.Hash{}) {
				voidMiss++
			}
		}
		res.check(fmt.Sprintf("ALL outstanding DEL check slots voided (%d checks, %d slots)", b.checkRefundN, len(b.checkVoid)), voidMiss == 0, fmt.Sprintf("%d slots still set", voidMiss))
		res.check("checks refund total == report", b.checkRefund.Equal(report.ChecksRefundedDel), fmt.Sprintf("independent=%s report=%s", b.checkRefund, report.ChecksRefundedDel))
		res.check("checks voided count == report", b.checkRefundN == report.ChecksVoided, fmt.Sprintf("independent=%d report=%d", b.checkRefundN, report.ChecksVoided))
	}

	// --- EVM corruption diff: every changed slot in delegation/WDEL/checks is a clean floor or void ---
	diffEVMContract(ctx, a, div, res, "delegation", b.delegAddr, b.delegStore, false)
	diffEVMContract(ctx, a, div, res, "WDEL", b.wdelAddr, b.wdelStore, false)
	diffEVMContract(ctx, a, div, res, "checks", b.checksAddr, b.checksStore, true)
}

// verifyEVMSlots asserts each captured slot's post-scale value equals floor(old/div).
func verifyEVMSlots(ctx sdk.Context, a *app.DSC, div sdkmath.Int, res *checkResult, label string, slots []evmSlotVal) {
	if len(slots) == 0 {
		res.check(fmt.Sprintf("ALL %s floored (0)", label), true, "none present")
		return
	}
	miss, first := 0, ""
	dv := div.BigInt()
	for _, sv := range slots {
		got := new(big.Int).SetBytes(a.EvmKeeper.GetState(ctx, sv.addr, sv.slot).Bytes())
		want := new(big.Int).Quo(sv.old, dv)
		if got.Cmp(want) != 0 {
			miss++
			if first == "" {
				first = fmt.Sprintf("%s/%s old=%s got=%s want=%s", sv.addr.Hex(), sv.slot.Hex(), sv.old, got, want)
			}
		}
	}
	res.check(fmt.Sprintf("ALL %s floored (%d checked)", label, len(slots)), miss == 0, fmt.Sprintf("%d mismatches; first: %s", miss, first))
}

// diffEVMContract re-reads a contract's full storage after scaling and asserts every slot
// that changed is either a clean ÷div floor of its old value, or a deletion (new==0:
// floored-to-zero for value slots, or a void for the checks contract). A change that is
// neither (corruption) or a slot created from zero fails the check.
func diffEVMContract(ctx sdk.Context, a *app.DSC, div sdkmath.Int, res *checkResult, label string, addr common.Address, beforeStore stakescan.Storage, allowVoid bool) {
	if addr == (common.Address{}) || beforeStore == nil {
		return
	}
	dv := div.BigInt()
	changed, floored, voided, bad := 0, 0, 0, 0
	first := ""
	for slot, oldHash := range beforeStore {
		old := new(big.Int).SetBytes(oldHash.Bytes())
		newV := new(big.Int).SetBytes(a.EvmKeeper.GetState(ctx, addr, slot).Bytes())
		if old.Cmp(newV) == 0 {
			continue
		}
		changed++
		want := new(big.Int).Quo(old, dv)
		switch {
		case newV.Cmp(want) == 0 && newV.Sign() != 0:
			floored++
		case newV.Sign() == 0 && (want.Sign() == 0 || allowVoid):
			voided++ // floored-to-zero (deleted) or an intended void
		default:
			bad++
			if first == "" {
				first = fmt.Sprintf("%s old=%s new=%s want=%s", slot.Hex(), old, newV, want)
			}
		}
	}
	res.check(fmt.Sprintf("EVM %s diff: every changed slot is a clean floor/void (%d changed: %d floored, %d voided)", label, changed, floored, voided),
		bad == 0, fmt.Sprintf("%d unexplained; first: %s", bad, first))
}

// --- after-state collectors (positional, same iteration order as capture) ---

func collectDelegs(ctx sdk.Context, a *app.DSC, base string) map[string]sdkmath.Int {
	m := map[string]sdkmath.Int{}
	for _, d := range a.ValidatorKeeper.GetAllDelegations(ctx) {
		if d.Stake.Stake.Denom == base {
			m[fmt.Sprintf("%s|%s|%s", d.Delegator, d.Validator, d.Stake.ID)] = d.Stake.Stake.Amount
		}
	}
	return m
}
func collectUndelegs(ctx sdk.Context, a *app.DSC, base string) map[string]sdkmath.Int {
	m := map[string]sdkmath.Int{}
	a.ValidatorKeeper.IterateUndelegations(ctx, func(_ int64, ubd validatortypes.Undelegation) bool {
		for i, e := range ubd.Entries {
			if e.Stake.Stake.Denom == base {
				m[fmt.Sprintf("%s|%s|%d", ubd.Delegator, ubd.Validator, i)] = e.Stake.Stake.Amount
			}
		}
		return false
	})
	return m
}
func collectRedelegs(ctx sdk.Context, a *app.DSC, base string) map[string]sdkmath.Int {
	m := map[string]sdkmath.Int{}
	a.ValidatorKeeper.IterateRedelegations(ctx, func(_ int64, red validatortypes.Redelegation) bool {
		for i, e := range red.Entries {
			if e.Stake.Stake.Denom == base {
				m[fmt.Sprintf("%s|%s|%s|%d", red.Delegator, red.ValidatorSrc, red.ValidatorDst, i)] = e.Stake.Stake.Amount
			}
		}
		return false
	})
	return m
}
func collectNFTRes(ctx sdk.Context, a *app.DSC, base string) map[string]sdkmath.Int {
	m := map[string]sdkmath.Int{}
	for _, col := range a.NFTKeeper.GetCollections(ctx) {
		creator, err := sdk.AccAddressFromBech32(col.Creator)
		if err != nil {
			continue
		}
		for _, token := range a.NFTKeeper.GetTokens(ctx, creator, col.Denom) {
			if token.Reserve.Denom == base {
				m["tok|"+token.ID] = token.Reserve.Amount
			}
			for _, st := range a.NFTKeeper.GetSubTokens(ctx, token.ID) {
				if st.Reserve != nil && st.Reserve.Denom == base {
					m[fmt.Sprintf("sub|%s|%d", token.ID, st.ID)] = st.Reserve.Amount
				}
			}
		}
	}
	return m
}
func collectLegacy(ctx sdk.Context, a *app.DSC, base string) map[string]sdkmath.Int {
	m := map[string]sdkmath.Int{}
	for _, rec := range a.LegacyKeeper.GetLegacyRecords(ctx) {
		for i, c := range rec.Coins {
			if c.Denom == base {
				m[fmt.Sprintf("%s|%d", rec.LegacyAddress, i)] = c.Amount
			}
		}
	}
	return m
}

// verifyCosmosAmts checks each captured old value's post-scale counterpart (looked up by
// id) equals floor(old/div), and that the count is unchanged.
func verifyCosmosAmts(res *checkResult, label string, old []amt, after map[string]sdkmath.Int) {
	div := redenom.DefaultDivisor
	miss, first := 0, ""
	for _, o := range old {
		got, ok := after[o.id]
		if !ok {
			miss++
			if first == "" {
				first = fmt.Sprintf("%s MISSING after", o.id)
			}
			continue
		}
		want := floorI(o.old, div)
		if !got.Equal(want) {
			miss++
			if first == "" {
				first = fmt.Sprintf("%s old=%s got=%s want=%s", o.id, o.old, got, want)
			}
		}
	}
	res.check(fmt.Sprintf("ALL %s floored (%d checked)", label, len(old)), miss == 0, fmt.Sprintf("%d mismatches; first: %s", miss, first))
}

// locateChecksContractMain mirrors app/redenom.locateChecksContract (which is unexported):
// it scans the EVM storage namespace and returns the first contract whose storage
// reconstructs verified checks. DecimalChecks does not self-register in the ContractCenter,
// so it must be detected by storage signature (the reconstruction self-verifies, so there
// are no false positives).
func locateChecksContractMain(kv sdk.KVStore) (common.Address, []stakescan.Check) {
	ps := prefix.NewStore(kv, evmtypes.KeyPrefixStorage)
	it := ps.Iterator(nil, nil)
	defer it.Close()

	var (
		cur     common.Address
		curHas  bool
		curStor stakescan.Storage
	)
	flush := func() (common.Address, []stakescan.Check, bool) {
		if !curHas {
			return common.Address{}, nil, false
		}
		checks := stakescan.ReconstructChecks(curStor, stakescan.ChecksBase)
		for _, c := range checks {
			if c.DetailsFound {
				return cur, checks, true
			}
		}
		return common.Address{}, nil, false
	}
	for ; it.Valid(); it.Next() {
		key := it.Key()
		if len(key) != 20+32 {
			continue
		}
		addr := common.BytesToAddress(key[:20])
		if !curHas || addr != cur {
			if a, checks, found := flush(); found {
				return a, checks
			}
			cur, curHas, curStor = addr, true, stakescan.Storage{}
		}
		curStor[common.BytesToHash(key[20:])] = common.BytesToHash(it.Value())
	}
	if a, checks, found := flush(); found {
		return a, checks
	}
	return common.Address{}, nil
}
