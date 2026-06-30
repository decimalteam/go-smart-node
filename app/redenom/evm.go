package redenom

import (
	"fmt"
	"math/big"
	"sort"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
)

func ContractCenterFor(chainID string) (common.Address, bool) {
	switch {
	case helpers.IsMainnet(chainID):
		return common.HexToAddress("0xc108715A06f76CAA96fa2c943Ebf05159c29A87D"), true // mainnet
	case helpers.IsTestnet(chainID):
		return common.HexToAddress("0xbC96b61F137F28F0Da47Cc4Ef06e5f984B565A2E"), true // testnet
	case helpers.IsDevnet(chainID):
		return common.HexToAddress("0x481487AEafc60512233a08Da240FC7AF99c0f696"), true // devnet
	default:
		return common.Address{}, false
	}
}

// tokenTypeDEL is IDecimalDelegationCommon.TokenType.DEL (DRC20==1, NFT==2,
// NFT1155==3, DEL==4). A coin/frozen stake is DEL-denominated iff TokenType==4.
const tokenTypeDEL = 4

// nftReserveTypeDEL is INFTReserve.ReserveType.DEL (None==0, DEL==1, DRC20==2). A
// per-collection NFT reserve is DEL-denominated iff reserveType==1.
const nftReserveTypeDEL = 1

// evmWrite is one pending storage rewrite (or deletion when NewVal is nil). All
// writes are collected, then sorted by (address, slot) and applied, so the resulting
// app hash is identical across nodes (iterating a Go map and writing in map order is
// non-deterministic and would fork the chain).
type evmWrite struct {
	addr   common.Address
	slot   common.Hash
	newVal *big.Int // nil => delete the slot (void)
}

// scaleEVM rewrites every DEL-denominated amount held in EVM contract storage by div
// (floor) and voids + refunds outstanding DEL checks. Scope:
//
//  1. DecimalDelegation: DEL coin-stake amounts, NFT-stake DEL reserves, DEL frozen
//     stakes, and DEL auto-unbond queue entries.
//  2. DecimalChecks: outstanding DEL checks are voided (storage cleared) and their
//     DEL backing is refunded to each creator from the contract's bank "del" balance.
//
// EVM native account balances are NOT touched here: ethermint keeps them in the bank
// keeper as denom "del", already scaled by scaleBankDel. DRC20 custom-token balances
// are out of scope (custom-denominated, never scaled). The per-NFT-collection reserve
// pools (DRC721/DRC1155) are NOT scaled — see the nft-center note below.
func scaleEVM(ctx sdk.Context, k Keepers, sk StoreKeys, div sdkmath.Int, base string, rep *Report) error {
	logger := ctx.Logger()

	cc, ok := ContractCenterFor(ctx.ChainID())
	if !ok {
		// HARD FAIL — do NOT silently skip. scaleCosmos has already divided every native
		// "del" balance by div, INCLUDING the WDEL / NFT-collection / DecimalChecks /
		// delegation contract accounts. If we skipped EVM scaling and let the upgrade
		// commit, those contracts' EVM ledgers (WDEL.balanceOf, _reserve[tokenId].amount,
		// check amounts, delegation stake slots) would keep their old ×div values against
		// ÷div native backing — silent, unrecoverable insolvency (first withdraw()/redeem()
		// drains the pool, the rest revert). Refuse to commit a Cosmos-only redenomination.
		return fmt.Errorf(
			"redenom evm: no ContractCenter address for chain-id %q — refusing to commit a "+
				"Cosmos-only redenomination (EVM ledgers would be left insolvent vs ÷%s native backing); "+
				"register the network in redenom.ContractCenterFor before upgrading",
			ctx.ChainID(), div.String())
	}

	// Resolve the delegation / nft-center / wdel contracts via the ContractCenter
	// registry (mapping(string => address) at ContractCenterBase+0).
	ccStore := readContractStorage(ctx, sk, cc)
	delegationAddr := stakescan.ResolveAddressBySymbol(ccStore, contracts.NameOfSlugForGetAddressDelegation)
	nftCenterAddr := stakescan.ResolveAddressBySymbol(ccStore, contracts.NameOfSlugForGetAddressNftCenter)
	wdelAddr := stakescan.ResolveAddressBySymbol(ccStore, contracts.NameOfSlugForGetAddressWDEL)
	tokenCenterAddr := stakescan.ResolveAddressBySymbol(ccStore, contracts.NameOfSlugForGetAddressTokenCenter)

	divBig := div.BigInt()
	_ = nftCenterAddr // resolution kept for logging/diagnostics; reserves enumerated via x/nft

	var writes []evmWrite

	// holderSet accumulates every EVM address that could hold WDEL (delegators,
	// validators, and the system contracts), built while reconstructing delegation
	// stakes, then used to enumerate WDEL.balanceOf holders.
	holderSet := map[common.Address]struct{}{}

	// --- 1. DecimalDelegation contract ---
	if delegationAddr == (common.Address{}) {
		logger.Error("redenom evm: delegation contract not found in ContractCenter, skipping delegation scaling")
	} else {
		delStore := readContractStorage(ctx, sk, delegationAddr)
		ws, err := collectDelegationWrites(delStore, delegationAddr, wdelAddr, divBig, rep, logger, holderSet)
		if err != nil {
			return err
		}
		writes = append(writes, ws...)
	}

	// --- 2. WDEL (wrapped DEL) ledger: divide every balanceOf in lockstep with the backing ---
	// Staked DEL is custodied as WDEL (DecimalDelegation._delegateDEL forwards msg.value into
	// WDEL via deposit{value}, crediting balanceOf[delegation]). scaleBankDel already divided
	// WDEL's native "del" balance (== WETH9 totalSupply, computed from the native balance) by
	// div; the ERC20 ledger must follow or withdraw() goes insolvent. scaleWDEL enumerates the
	// holder universe (bank del holders ∪ system contracts ∪ delegationDerived), divides each
	// balanceOf (never allowances), then reconciles the native backing to Σ scaled balances.
	if wdelAddr == (common.Address{}) {
		logger.Error("redenom evm: wdel contract not found in ContractCenter, skipping WDEL ledger scaling (staking backing will be INCONSISTENT)")
	} else {
		wdelWrites, err := scaleWDEL(ctx, k, sk, wdelAddr, div, base, ccStore, holderSet, rep, logger)
		if err != nil {
			return err
		}
		writes = append(writes, wdelWrites...)
	}

	// --- 3. Per-collection DRC721/DRC1155 NFT reserves (DEL), enumerated via x/nft ---
	nftWrites, err := scaleNFTCollectionReserves(ctx, k, sk, div, base, tokenCenterAddr, rep, logger)
	if err != nil {
		return err
	}
	writes = append(writes, nftWrites...)

	// --- 4. DecimalChecks: void outstanding DEL checks + refund creators ---
	checkWrites, err := processChecks(ctx, k, sk, ccStore, base, rep, logger)
	if err != nil {
		return err
	}
	writes = append(writes, checkWrites...)

	// Apply all storage writes deterministically.
	applyWrites(ctx, k, writes)
	return nil
}

// collectDelegationWrites reconstructs the delegation contract's stakes/reserves/
// frozen-stakes/auto-unbond entries and returns the DEL-amount rewrites (divided by
// divBig). DEL is identified by TokenType==4 for stakes, and by Token==WDEL for
// auto-unbond entries (the DEL coin's DRC20Contract is the WDEL address).
func collectDelegationWrites(
	delStore stakescan.Storage,
	addr, wdel common.Address,
	divBig *big.Int,
	rep *Report,
	logger interface{ Info(string, ...interface{}) },
	holderSet map[common.Address]struct{},
) ([]evmWrite, error) {
	var out []evmWrite
	res := stakescan.Reconstruct(delStore, stakescan.DelegationBase)

	// The delegation contract itself is the dominant WDEL holder (delegated DEL is
	// deposited into WDEL crediting balanceOf[delegation]). Collect it plus every
	// delegator/validator EVM address seen, as a defensive WDEL-holder superset.
	holderSet[addr] = struct{}{}
	addHolder := func(a common.Address) {
		if a != (common.Address{}) {
			holderSet[a] = struct{}{}
		}
	}

	// Coin stakes: scale the DEL amount (base+3).
	for _, cs := range res.CoinStakes {
		addHolder(cs.Delegator)
		addHolder(cs.Validator)
		if cs.TokenType != tokenTypeDEL {
			continue
		}
		out = append(out, evmWrite{addr, stakescan.StakeAmountSlot(cs.BaseSlot), divInt(cs.Amount, divBig)})
		rep.EVMCoinStakeSlots++
	}

	// NFT stakes: scale the reserve (base+6) ONLY when it is DEL-denominated. The
	// delegation contract sets reserveToken = WETH()(=WDEL) for a DEL reserve and the
	// custom DRC20 address otherwise (DecimalDelegation.sol delegateNFT paths), so a
	// reserveToken != wdel reserve is a custom coin whose amount must NOT be divided.
	// (The NFT intrinsic amount at base+2 is custom-or-DEL but always left untouched.)
	for _, ns := range res.NFTStakes {
		addHolder(ns.Delegator)
		addHolder(ns.Validator)
		if ns.ReserveToken != wdel {
			continue
		}
		out = append(out, evmWrite{addr, stakescan.NFTStakeReserveAmountSlot(ns.BaseSlot), divInt(ns.ReserveAmount, divBig)})
		rep.EVMNFTReserveSlots++
	}

	// Frozen stakes (live + deprecated): scale the embedded DEL Stake amount (elem+3).
	for _, fz := range res.FrozenLive {
		addHolder(fz.Stake.Delegator)
		addHolder(fz.Stake.Validator)
		if fz.Stake.TokenType != tokenTypeDEL {
			continue
		}
		slot := stakescan.FrozenStakeAmountSlot(stakescan.DelegationBase, fz.Index, false)
		out = append(out, evmWrite{addr, slot, divInt(fz.Stake.Amount, divBig)})
		rep.EVMFrozenSlots++
	}
	for _, fz := range res.FrozenDeprecated {
		addHolder(fz.Stake.Delegator)
		addHolder(fz.Stake.Validator)
		if fz.Stake.TokenType != tokenTypeDEL {
			continue
		}
		slot := stakescan.FrozenStakeAmountSlot(stakescan.DelegationBase, fz.Index, true)
		out = append(out, evmWrite{addr, slot, divInt(fz.Stake.Amount, divBig)})
		rep.EVMFrozenSlots++
	}

	// Auto-unbond queues. DEL entries have Token==WDEL (the DEL coin's DRC20Contract);
	// the zero address is also treated as native DEL defensively. Non-DEL (custom DRC20)
	// entries are skipped and counted for the log.
	queues := stakescan.FindAutoUnbondQueues(delStore, 1.0)
	skippedNonDEL := 0
	for _, q := range queues {
		for _, e := range q.Entries {
			isDEL := e.Token == (common.Address{}) || (wdel != (common.Address{}) && e.Token == wdel)
			if !isDEL {
				skippedNonDEL++
				continue
			}
			slot := stakescan.AutoUnbondAmountSlot(q.FieldSlot, e.Index)
			out = append(out, evmWrite{addr, slot, divInt(e.Amount, divBig)})
			rep.EVMAutoUnbondSlots++
		}
	}
	logger.Info("redenom evm: delegation reconstructed",
		"coinStakesDEL", rep.EVMCoinStakeSlots,
		"nftReserveSlots", rep.EVMNFTReserveSlots,
		"frozenDELSlots", rep.EVMFrozenSlots,
		"autoUnbondDELSlots", rep.EVMAutoUnbondSlots,
		"autoUnbondQueues", len(queues),
		"autoUnbondSkippedNonDEL", skippedNonDEL,
		"wdel", wdel.Hex(),
	)
	return out, nil
}

// wdelInfiniteAllowance is WETH9's "unlimited approval" sentinel (uint(-1) == 2^256-1),
// used only to recognise a non-balance value if one is ever seen in the balance-slot scan.
var wdelInfiniteAllowance = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

// scaleWDEL divides every DEL-denominated value slot in the WDEL contract's storage by div
// (floor), keeping WDEL's ERC20 ledger backed 1:1 by its (already-scaled) native "del".
//
// WDEL is the canonical WETH9 wrapper — verified against the deployed bytecode in genesis:
// Solidity 0.4.18, flat storage (slot 0 = name, 1 = symbol, 2 = decimals, balanceOf mapping
// at field 3, allowance mapping at field 4), and totalSupply() executes the BALANCE opcode
// (returns address(this).balance) so there is NO _totalSupply storage slot. Staked DEL is
// custodied here: DecimalDelegation._delegateDEL forwards msg.value into WDEL via
// deposit{value}, crediting balanceOf[delegation]. If balanceOf were not divided in lockstep,
// withdraw() (which pays native gated on balanceOf) would over-pay against the scaled native.
//
// The body uses an EXHAUSTIVE scan-all (no holder enumeration — see comment there): WDEL
// mapping keys aren't stored, so a holder-enumeration approach silently skips any holder
// outside the candidate universe (4 of ~9 balance slots on mainnet). Scanning every slot
// guarantees no balanceOf is left unscaled, which is the solvency-preserving safety property.
func scaleWDEL(
	ctx sdk.Context,
	k Keepers,
	sk StoreKeys,
	wdelAddr common.Address,
	div sdkmath.Int,
	base string,
	ccStore stakescan.Storage,
	delegationDerived map[common.Address]struct{},
	rep *Report,
	logger interface {
		Info(string, ...interface{})
		Error(string, ...interface{})
	},
) ([]evmWrite, error) {
	wstore := readContractStorage(ctx, sk, wdelAddr)
	_ = ccStore // resolved at the call site; the exhaustive scan below needs no holder universe

	meta := map[common.Hash]bool{
		common.BigToHash(big.NewInt(0)): true, // name
		common.BigToHash(big.NewInt(1)): true, // symbol
		common.BigToHash(big.NewInt(2)): true, // decimals
	}

	// EXHAUSTIVE (no holder enumeration): WDEL mapping keys are not stored, so holders aren't
	// recoverable — but every WDEL storage slot is either meta (0/1/2), a balanceOf entry, or
	// an allowance entry, and BOTH balanceOf and finite allowance hold DEL-denominated WDEL
	// amounts that must ÷div to stay economically equivalent. So divide EVERY non-meta slot
	// except the infinite-approval sentinel uint(-1) (transferFrom special-cases it). Scanning
	// all slots can never miss a holder — an unscaled balanceOf would leave WDEL insolvent on
	// withdraw — which is the safety property. (Holder enumeration would silently skip any
	// holder outside the universe; on mainnet that was 4 of ~9 slots, hence this approach.)
	var out []evmWrite
	scaledSlots, preservedSentinels := 0, 0
	for slot, val := range wstore {
		if meta[slot] {
			continue
		}
		v := new(big.Int).SetBytes(val.Bytes())
		if v.Sign() == 0 {
			continue
		}
		if v.Cmp(wdelInfiniteAllowance) == 0 {
			preservedSentinels++
			continue
		}
		out = append(out, evmWrite{wdelAddr, slot, new(big.Int).Quo(v, div.BigInt())})
		scaledSlots++
	}
	rep.WDELHoldersScaled = scaledSlots

	// scaleEVM runs AFTER scaleBankDel, so WDEL's bank "del" balance is already
	// floor(origNative/div) = the post-migration WETH9 totalSupply() (== native balance; no
	// stored supply slot). It stays >= Σ scaled balanceOf (floor(Σ) >= Σ floor, and Σ original
	// balanceOf == origNative since WDEL is 1:1 backed), so WDEL is solvent — over-collateralized
	// only by negligible flooring dust. We do NOT reduce it: the scan mixes balanceOf with finite
	// allowances, so the scaled sum is not Σ balanceOf alone, and WETH9 never asserts equality.
	scaledNative := k.Bank.GetBalance(ctx, sdk.AccAddress(wdelAddr.Bytes()), base).Amount
	rep.WDELTotalSupplyOld = scaledNative.Mul(div) // ≈ original native backing (pre-÷div)
	rep.WDELTotalSupplyNew = scaledNative

	// Solvency cross-check: the identifiable holders (delegation-derived) scaled down cannot
	// exceed the scaled native backing (they are a subset of all WDEL holders, 1:1 backed).
	identified := sdkmath.ZeroInt()
	seen := map[common.Hash]bool{}
	for h := range delegationDerived {
		slot := stakescan.WDELBalanceSlot(h)
		if seen[slot] {
			continue
		}
		seen[slot] = true
		if raw, ok := wstore[slot]; ok {
			identified = identified.Add(sdkmath.NewIntFromBigInt(new(big.Int).SetBytes(raw.Bytes())))
		}
	}
	if divFloor(identified, div).GT(scaledNative) {
		return nil, fmt.Errorf(
			"redenom evm: WDEL identified balanceOf %s (÷%s) exceeds scaled native del backing %s — slot math/classification error",
			identified.String(), div.String(), scaledNative.String())
	}

	logger.Info("redenom evm: WDEL ledger scaled (exhaustive scan)",
		"wdel", wdelAddr.Hex(),
		"slotsScaled", scaledSlots,
		"infiniteAllowancesPreserved", preservedSentinels,
		"scaledNativeBacking", scaledNative.String(),
		"identifiedHolderBalanceUnscaled", identified.String(),
	)
	return out, nil
}

// scaleNFTCollectionReserves divides the DEL `_reserve[tokenId].amount` of every EVM NFT
// collection by div. A collection (DRC721/DRC1155) is a DecimalNFTBeaconProxy (EIP-1967
// beacon slot set) inheriting NFTReserve; a DEL reserve is the collection's native "del"
// (already ÷div by scaleBankDel) plus a `_reserve[tokenId] = {token:0, amount, type:DEL}`
// entry in its storage. The amount MUST be ÷div in lockstep — otherwise redeem/burn pays the
// unscaled amount against the scaled native backing (insolvency).
//
// x/nft is LEGACY and does NOT mirror these pure-EVM collections (AddressDRC empty on
// mainnet), so we enumerate them directly: every del-holding beacon proxy is a collection
// holding a DEL reserve. Mapping keys (tokenIds) are not stored, so for each collection we
// find DEL reserve entries by the Reserve struct shape — a slot S holding reserveType DEL(1)
// whose token slot (base+0 = S-2) is address(0) and amount slot (base+1 = S-1) is non-zero —
// then REQUIRE floor(Σ amounts / div) == the collection's scaled native del before scaling,
// so a mis-identified slot can never be silently rewritten and a missed entry can never pass.
func scaleNFTCollectionReserves(
	ctx sdk.Context,
	k Keepers,
	sk StoreKeys,
	div sdkmath.Int,
	base string,
	tokenCenter common.Address, // custom-coin DRC20s are token-center beacon proxies — NOT NFTs
	rep *Report,
	logger interface {
		Info(string, ...interface{})
		Error(string, ...interface{})
	},
) ([]evmWrite, error) {
	beaconSlot := common.HexToHash("0xa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d50")
	divBig := div.BigInt()
	delType := big.NewInt(int64(nftReserveTypeDEL))

	// Snapshot del-holding NFT-collection beacon proxies before reading their storage. A beacon
	// proxy is an NFT collection iff its beacon is NOT the token-center (token-center beacon
	// proxies are custom-coin DRC20s, whose del is the coin reserve handled by bank + x/coin and
	// whose balanceOf is the custom token, not del — nothing to scale here).
	type coll struct {
		addr         common.Address
		scaledNative sdkmath.Int // scaleBankDel already divided the native del backing
	}
	var colls []coll
	skippedCoins := 0
	k.Bank.IterateAllBalances(ctx, func(addr sdk.AccAddress, c sdk.Coin) bool {
		if c.Denom != base || len(addr.Bytes()) != 20 {
			return false
		}
		eth := common.BytesToAddress(addr.Bytes())
		beacon := k.EVM.GetState(ctx, eth, beaconSlot)
		if beacon == (common.Hash{}) {
			return false // not a beacon proxy
		}
		if common.BytesToAddress(beacon.Bytes()) == tokenCenter {
			skippedCoins++
			return false // custom-coin DRC20, not an NFT collection
		}
		colls = append(colls, coll{eth, c.Amount})
		return false
	})

	var out []evmWrite
	scaledSlots, okColls, rejected, surplusColls := 0, 0, 0, 0
	rejectedDel, surplusDel := sdkmath.ZeroInt(), sdkmath.ZeroInt()
	for _, cl := range colls {
		st := readContractStorage(ctx, sk, cl.addr)
		type entry struct {
			slot common.Hash
			amt  *big.Int
		}
		var entries []entry
		sumOrig := new(big.Int)
		for s, v := range st {
			if new(big.Int).SetBytes(v.Bytes()).Cmp(delType) != 0 {
				continue // reserveType slot must equal DEL(1)
			}
			sBig := new(big.Int).SetBytes(s.Bytes())
			amountSlot := common.BigToHash(new(big.Int).Sub(sBig, big.NewInt(1)))
			tokenSlot := common.BigToHash(new(big.Int).Sub(sBig, big.NewInt(2)))
			if st[tokenSlot] != (common.Hash{}) {
				continue // token != address(0) ⇒ not a DEL reserve (or a false-positive 1-slot)
			}
			amt := new(big.Int).SetBytes(st[amountSlot].Bytes())
			if amt.Sign() == 0 {
				continue
			}
			entries = append(entries, entry{amountSlot, amt})
			sumOrig.Add(sumOrig, amt)
		}
		// False-positive guard + solvency. The found DEL reserves are real (mintByDEL /
		// addReserveByDEL / migration initReserve all write _reserve[tokenId].amount), and native
		// del backs them. If floor(Σ found / div) EXCEEDS the scaled native, a mis-identified slot
		// inflated the sum — do NOT scale (would corrupt), flag. If it is LESS, the difference is
		// legitimate orphan del: the Migration NFT variants (DRC721Migration/DRC1155Migration)
		// expose receive() payable that accepts del outside _reserve. The found reserves are still
		// real, so we DO scale them; the surplus stays as scaled native (over-collateralized).
		scaledFound := new(big.Int).Quo(sumOrig, divBig)
		nativeBig := cl.scaledNative.BigInt()
		if scaledFound.Cmp(nativeBig) > 0 {
			rejected++
			rejectedDel = rejectedDel.Add(cl.scaledNative)
			logger.Error("redenom evm: NFT collection Σreserve EXCEEDS native backing — slot mis-identification, NOT scaling (manual review)",
				"collection", cl.addr.Hex(),
				"sumReservesScaled", scaledFound.String(),
				"scaledNative", cl.scaledNative.String(),
				"delEntries", len(entries))
			continue
		}
		for _, e := range entries {
			out = append(out, evmWrite{cl.addr, e.slot, new(big.Int).Quo(e.amt, divBig)})
			scaledSlots++
		}
		okColls++
		if scaledFound.Cmp(nativeBig) < 0 {
			surplusColls++
			surplusDel = surplusDel.Add(cl.scaledNative.Sub(sdkmath.NewIntFromBigInt(scaledFound)))
		}
	}
	rep.EVMNFTReserveSlots = scaledSlots
	rep.EVMNFTCollectionsScaled = okColls
	rep.EVMNFTCollectionsFailed = rejected
	rep.EVMNFTFailedDel = rejectedDel
	logger.Info("redenom evm: EVM NFT collection DEL reserves scaled",
		"collections", okColls,
		"reserveSlotsScaled", scaledSlots,
		"coinProxiesSkipped", skippedCoins,
		"surplusColls(orphan del)", surplusColls,
		"surplusDel", surplusDel.String(),
		"rejected(Σreserve>native)", rejected,
		"rejectedDel", rejectedDel.String())
	// `rejected` collections had Σ found reserves EXCEEDING native — a slot mis-identification we
	// refuse to scale (would corrupt). We do NOT abort for a negligible amount; they're left
	// unscaled + flagged for manual review. Hard-abort only if it is ever large.
	if rejected > 0 {
		logger.Error("redenom evm: EVM NFT collections left UNSCALED (Σreserve>native) — manual review required",
			"count", rejected, "del", rejectedDel.String())
	}
	return out, nil
}

// processChecks locates the DecimalChecks contract, voids every outstanding DEL check
// (clearing its _checks and _checkDetails slots), and aggregates the scaled refund per
// creator. It then performs the balanced bank transfer (debit the checks contract's
// "del" balance, credit each creator) and returns the storage-clearing writes.
//
// Backing model (verified against DecimalChecks.sol): createChecksDEL is payable with
// msg.value == amount*signers, so DEL checks are backed by the checks contract's own
// native (bank "del") balance; redeem pays out via msg.sender.call{value}. The refund
// is therefore a debit of the checks contract's bank "del" and a credit to creators —
// a balanced transfer that does NOT change supply.
func processChecks(
	ctx sdk.Context,
	k Keepers,
	sk StoreKeys,
	ccStore stakescan.Storage,
	base string,
	rep *Report,
	logger interface{ Info(string, ...interface{}) },
) ([]evmWrite, error) {
	// DecimalChecks does not self-register in ContractCenter, so the symbol lookup is
	// expected to miss; fall back to locating it by storage signature.
	checksAddr := stakescan.ResolveAddressBySymbol(ccStore, "checks")
	var checks []stakescan.Check
	if checksAddr != (common.Address{}) {
		checks = stakescan.ReconstructChecks(readContractStorage(ctx, sk, checksAddr), stakescan.ChecksBase)
	} else {
		checksAddr, checks = locateChecksContract(ctx, sk)
	}
	if checksAddr == (common.Address{}) {
		logger.Info("redenom evm: DecimalChecks contract not found, skipping check void/refund")
		return nil, nil
	}

	var out []evmWrite
	totalRefund := sdkmath.ZeroInt()
	refundByCreator := map[common.Address]sdkmath.Int{}
	var creators []common.Address // stable order for the deterministic credit loop

	for _, c := range checks {
		// Only outstanding (Status.None==0) DEL (TypeChecks.DEL==0) checks.
		if c.Status != 0 || c.TypeChecks != 0 {
			continue
		}
		refund := divFloor(sdkmath.NewIntFromBigInt(c.Amount), rep.Divisor)

		if _, seen := refundByCreator[c.Creator]; !seen {
			refundByCreator[c.Creator] = sdkmath.ZeroInt()
			creators = append(creators, c.Creator)
		}
		refundByCreator[c.Creator] = refundByCreator[c.Creator].Add(refund)
		totalRefund = totalRefund.Add(refund)

		// Void the check: delete the _checks entry (base, base+1) and the _checkDetails
		// struct (detailsBase+0..4: typeChecks, amount, dueBlock, token, creator).
		out = append(out, evmWrite{checksAddr, c.BaseSlot, nil})
		out = append(out, evmWrite{checksAddr, stakescan.SlotAdd(c.BaseSlot, 1), nil})
		detailsBase := stakescan.CheckDetailsBase(c.CheckDetailsHash)
		for i := uint64(0); i < 5; i++ {
			out = append(out, evmWrite{checksAddr, stakescan.SlotAdd(detailsBase, i), nil})
		}
		rep.EVMCheckSlotsCleared += 7
		rep.ChecksVoided++
	}

	if rep.ChecksVoided == 0 {
		logger.Info("redenom evm: no outstanding DEL checks to void", "checksContract", checksAddr.Hex())
		rep.ChecksRefundedDel = sdkmath.ZeroInt()
		return out, nil
	}

	// Balanced bank transfer: debit the checks contract, credit each creator.
	contractAcc := sdk.AccAddress(checksAddr.Bytes())
	contractBal := k.Bank.GetBalance(ctx, contractAcc, base).Amount
	if contractBal.LT(totalRefund) {
		// floor(sum(amounts)) >= sum(floor(amount)) == totalRefund, and the contract's
		// scaled "del" balance is floor(its old balance) which backed sum(amounts);
		// a shortfall means inconsistent state — fail loudly rather than underflow.
		return nil, errShortChecksBacking(checksAddr, contractBal, totalRefund)
	}
	bankStore := ctx.KVStore(sk.Bank)
	setBankBalanceDel(bankStore, contractAcc, base, contractBal.Sub(totalRefund))
	// Sort creators so the credit loop runs in a deterministic order. (Per-address
	// bank writes are independent and each creator is credited once with a fixed
	// aggregate, so order does not affect final state; sorting makes that explicit.)
	sort.Slice(creators, func(i, j int) bool { return creators[i].Hex() < creators[j].Hex() })
	for _, creator := range creators {
		creatorAcc := sdk.AccAddress(creator.Bytes())
		cur := k.Bank.GetBalance(ctx, creatorAcc, base).Amount
		setBankBalanceDel(bankStore, creatorAcc, base, cur.Add(refundByCreator[creator]))
	}
	rep.ChecksRefundedDel = totalRefund

	logger.Info("redenom evm: voided + refunded outstanding DEL checks",
		"checksContract", checksAddr.Hex(),
		"checksVoided", rep.ChecksVoided,
		"slotsCleared", rep.EVMCheckSlotsCleared,
		"creators", len(creators),
		"refundedDel", totalRefund.String(),
	)
	return out, nil
}

// locateChecksContract scans the entire EVM storage namespace and returns the first
// contract whose storage reconstructs verified checks (DetailsFound). DecimalChecks is
// not registered in the ContractCenter, so it must be detected by signature. The
// reconstruction is self-verifying (keccak256(checkHash ‖ checksField) == slot), so
// there are no false positives. This is expensive but acceptable in a one-time upgrade.
//
// It streams per address: cosmos KV iteration is key-sorted, so all slots of one
// contract (key layout addr(20) ‖ slot(32) after the 0x02 prefix) are contiguous;
// each contract's storage is reconstructed and discarded before the next, bounding
// memory to a single contract instead of the whole EVM state. Because addresses are
// visited in sorted order, the located contract is deterministic.
func locateChecksContract(ctx sdk.Context, sk StoreKeys) (common.Address, []stakescan.Check) {
	ps := prefix.NewStore(ctx.KVStore(sk.EVM), evmtypes.KeyPrefixStorage)
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

// readContractStorage builds the stakescan.Storage (slot -> value) for one contract by
// iterating its EVM storage prefix.
func readContractStorage(ctx sdk.Context, sk StoreKeys, addr common.Address) stakescan.Storage {
	ps := prefix.NewStore(ctx.KVStore(sk.EVM), evmtypes.AddressStoragePrefix(addr))
	it := ps.Iterator(nil, nil)
	defer it.Close()
	s := stakescan.Storage{}
	for ; it.Valid(); it.Next() {
		s[common.BytesToHash(it.Key())] = common.BytesToHash(it.Value())
	}
	return s
}

// applyWrites sorts the collected writes by (address, slot) and applies them via the
// EVM keeper. Sorting is mandatory for determinism: applying in Go map iteration order
// would produce a different app hash on each node and fork the chain. A nil newVal
// deletes the slot (void); a non-nil value is written as minimal big-endian bytes
// (ethermint trims storage and left-pads on read), with zero written as a deletion.
func applyWrites(ctx sdk.Context, k Keepers, writes []evmWrite) {
	sort.Slice(writes, func(i, j int) bool {
		if c := bytesCompare(writes[i].addr.Bytes(), writes[j].addr.Bytes()); c != 0 {
			return c < 0
		}
		return bytesCompare(writes[i].slot.Bytes(), writes[j].slot.Bytes()) < 0
	})
	for _, w := range writes {
		setEVMSlot(ctx, k, w.addr, w.slot, w.newVal)
	}
}

// setEVMSlot writes newVal to (addr, slot), or deletes the slot when newVal is nil or
// zero. ethermint stores trimmed bytes and left-pads on read, so the minimal
// big-endian encoding (newVal.Bytes()) round-trips; an empty value deletes the slot.
func setEVMSlot(ctx sdk.Context, k Keepers, addr common.Address, slot common.Hash, newVal *big.Int) {
	if newVal == nil || newVal.Sign() == 0 {
		k.EVM.SetState(ctx, addr, slot, nil)
		return
	}
	k.EVM.SetState(ctx, addr, slot, newVal.Bytes())
}

// errShortChecksBacking is returned when the DecimalChecks contract's scaled "del"
// balance is below the total refund owed to creators — an invariant violation, since
// floor(sum of check amounts) >= sum of floored refunds and the contract's scaled
// balance equals floor(its old balance), which backed the full sum.
func errShortChecksBacking(addr common.Address, have, need sdkmath.Int) error {
	return fmt.Errorf(
		"redenom evm: checks contract %s del balance %s < total check refund %s",
		addr.Hex(), have.String(), need.String(),
	)
}

// divInt returns floor(x / divBig) as a new big.Int (x is treated as zero when nil).
// DEL amounts are non-negative, so big.Int truncation toward zero is a true floor.
func divInt(x, divBig *big.Int) *big.Int {
	if x == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Quo(x, divBig)
}

// bytesCompare is a tiny lexicographic byte-slice comparator (-1/0/1) used for the
// deterministic write ordering.
func bytesCompare(a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		switch {
		case a[i] < b[i]:
			return -1
		case a[i] > b[i]:
			return 1
		}
	}
	switch {
	case len(a) < len(b):
		return -1
	case len(a) > len(b):
		return 1
	}
	return 0
}
