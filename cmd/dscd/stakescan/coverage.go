package stakescan

import (
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Additional DelegationStorage field offsets used only by the coverage audit.
const (
	fFreezeTime      = 4 // mapping(FreezeType => uint256)
	fValidatorTokens = 5 // mapping(address => mapping(bytes32 => ValidatorReserve))
)

// knownFixedSlots are well-known proxy (ERC-1967) and OpenZeppelin v5 (ERC-7201)
// storage locations that live outside the DelegationStorage namespace. Each base holds
// a small struct/value (or a mapping whose entries are scattered, e.g. AccessControl).
var knownFixedSlots = map[common.Hash]string{
	common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"): "erc1967.implementation",
	common.HexToHash("0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103"): "erc1967.admin",
	common.HexToHash("0xa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d50"): "erc1967.beacon",
	common.HexToHash("0xf0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00"): "oz.Initializable",
	common.HexToHash("0x9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300"): "oz.Ownable",
	common.HexToHash("0xcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300"): "oz.Pausable",
	common.HexToHash("0x9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00"): "oz.ReentrancyGuard",
	common.HexToHash("0x237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00"): "oz.AccessControl",
	// Slot 0 holds a leftover value 1 from the pre-ERC-7201 OZ v4 Initializable
	// (_initialized) layout, retained after the contract moved to namespaced storage.
	common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"): "slot0.legacyInitialized",
}

// SlotValue is a raw storage entry.
type SlotValue struct {
	Slot  common.Hash `json:"slot"`
	Value common.Hash `json:"value"`
}

// AuditReport accounts for every stored slot of the delegation contract.
type AuditReport struct {
	TotalSlots  int            `json:"totalSlots"`
	Explained   int            `json:"explained"`
	Unexplained int            `json:"unexplained"`
	ByCategory  map[string]int `json:"byCategory"`
	// Sample of unexplained slots (sorted), capped for display; Unexplained has the full count.
	UnexplainedSample []SlotValue `json:"unexplainedSample"`
}

// Audit classifies every stored slot against the structures the reconstruction explains
// and returns the slots that remain unaccounted for. sampleLimit caps UnexplainedSample
// (0 = include all).
func Audit(s Storage, base common.Hash, r Result, au []AutoUnbondQueue, sampleLimit int) AuditReport {
	cov := make(map[common.Hash]string, len(s))
	mark := func(slot common.Hash, cat string) {
		if _, ok := cov[slot]; !ok {
			cov[slot] = cat
		}
	}

	stakesField := slotAdd(base, fStakes)
	penaltyField := slotAdd(base, fStakePenalty)
	nftBackedField := slotAdd(base, fNFTBacked)
	nftField := slotAdd(base, fNFTStakes)
	nftPenaltyField := slotAdd(base, fNFTPenalty)
	valTokField := slotAdd(base, fValidatorTokens)
	freezeTimeField := slotAdd(base, fFreezeTime)

	// Base value field and array length words.
	mark(base, "contractCenter")                        // B+0
	mark(slotAdd(base, fFrozenDep), "frozenDeprecated") // B+3 length
	mark(slotAdd(base, fFrozen), "frozenLive")          // B+9 length

	// _freezeTime[FreezeType] for the three enum keys.
	for ft := int64(0); ft <= 2; ft++ {
		mark(mappingSlotBytes32(common.BigToHash(big.NewInt(ft)), freezeTimeField), "freezeTime")
	}

	// Per-stake mapping slots (from coin stakes and from stakes referenced by frozen entries).
	seenValTok := make(map[common.Hash]bool)
	markStake := func(st Stake) {
		bsl := mappingSlotBytes32(st.StakeID, stakesField)
		for i := uint64(0); i < 8; i++ {
			mark(slotAdd(bsl, i), "coinStake")
		}
		mark(mappingSlotBytes32(st.StakeID, penaltyField), "stakePenaltyIndex")
		mark(mappingSlotBytes32(st.StakeID, nftBackedField), "nftBackedAmount")

		// _validatorTokens[validator][keccak(token,tokenId)] = ValidatorReserve{penaltyIndex, reserve}.
		htid := hashedTokenID(st.Token, st.TokenID)
		dedup := crypto.Keccak256Hash(st.Validator.Bytes(), htid.Bytes())
		if !seenValTok[dedup] {
			seenValTok[dedup] = true
			mid := mappingSlotBytes32(common.BytesToHash(st.Validator.Bytes()), valTokField)
			rbase := mappingSlotBytes32(htid, mid)
			mark(rbase, "validatorTokens")
			mark(slotAdd(rbase, 1), "validatorTokens")
		}
	}
	for _, cs := range r.CoinStakes {
		markStake(cs.Stake)
	}
	for _, f := range r.FrozenLive {
		markStake(f.Stake)
	}
	for _, f := range r.FrozenDeprecated {
		markStake(f.Stake)
	}

	// NFT stake mapping slots.
	for _, n := range r.NFTStakes {
		bsl := mappingSlotBytes32(n.NFTStakeID, nftField)
		for i := uint64(0); i < 9; i++ {
			mark(slotAdd(bsl, i), "nftStake")
		}
		mark(mappingSlotBytes32(n.NFTStakeID, nftPenaltyField), "nftStakePenaltyIndex")
	}

	// Frozen array element slots (their own keccak region, distinct from _stakes).
	markArray(s, mark, slotAdd(base, fFrozenDep), legacyFrozenStride, "frozenDeprecated")
	markArray(s, mark, slotAdd(base, fFrozen), frozenStride, "frozenLive")

	// Auto-unbond queue element slots.
	for _, q := range au {
		mark(q.FieldSlot, "autoUnbond")
		start := arrayElemStart(q.FieldSlot)
		for i := uint64(0); i < q.Length; i++ {
			elem := slotAdd(start, i*autoUnbondStride)
			for j := uint64(0); j < autoUnbondStride; j++ {
				mark(slotAdd(elem, j), "autoUnbond")
			}
		}
	}

	// _frozenStakeMigrated[frozenStakeHash] = true: migration dedup flags (own namespace).
	// Flags persist after a migrated frozen stake is completed and popped from _frozenStakes.
	migNS := ERC7201("DECIMAL_DELEGATION_MIGRATION_STORAGE_LOCATION")
	for _, f := range r.FrozenDeprecated {
		mark(mappingSlotBytes32(FrozenStakeHash(f.Stake, f.UnfreezeTimestamp), migNS), "frozenStakeMigrated")
	}
	for _, f := range r.FrozenLive {
		mark(mappingSlotBytes32(FrozenStakeHash(f.Stake, f.UnfreezeTimestamp), migNS), "frozenStakeMigrated")
	}

	// Proxy / OpenZeppelin fixed slots.
	for slot, name := range knownFixedSlots {
		mark(slot, name)
	}

	// Tally storage against coverage.
	byCat := make(map[string]int)
	var unexplained []SlotValue
	for k, v := range s {
		if cat, ok := cov[k]; ok {
			byCat[cat]++
		} else {
			unexplained = append(unexplained, SlotValue{Slot: k, Value: v})
		}
	}
	sort.Slice(unexplained, func(i, j int) bool {
		return unexplained[i].Slot.Hex() < unexplained[j].Slot.Hex()
	})

	sample := unexplained
	if sampleLimit > 0 && len(sample) > sampleLimit {
		sample = sample[:sampleLimit]
	}
	return AuditReport{
		TotalSlots:        len(s),
		Explained:         len(s) - len(unexplained),
		Unexplained:       len(unexplained),
		ByCategory:        byCat,
		UnexplainedSample: sample,
	}
}

func markArray(s Storage, mark func(common.Hash, string), fieldSlot common.Hash, stride uint64, cat string) {
	n, ok := arrayLen(s, fieldSlot)
	if !ok {
		return
	}
	start := arrayElemStart(fieldSlot)
	for i := uint64(0); i < n; i++ {
		elem := slotAdd(start, i*stride)
		for j := uint64(0); j < stride; j++ {
			mark(slotAdd(elem, j), cat)
		}
	}
}

// ClassMatch reports how many of a candidate mapping's computed slots fall in the
// unexplained set (used to identify leftover storage).
type ClassMatch struct {
	Name      string `json:"name"`
	FieldSlot string `json:"fieldSlot"`
	Matched   int    `json:"matched"`
	Probed    int    `json:"probed"`
}

// ERC7201 computes a namespaced storage base:
// keccak256(abi.encode(uint256(keccak256(name)) - 1)) & ~bytes32(uint256(0xff)).
func ERC7201(name string) common.Hash {
	inner := new(big.Int).SetBytes(crypto.Keccak256([]byte(name)))
	inner.Sub(inner, big.NewInt(1))
	h := crypto.Keccak256(common.BigToHash(inner).Bytes())
	h[31] = 0 // & ~0xff
	return common.BytesToHash(h)
}

// FrozenStakeHash = keccak256(abi.encode(validator, delegator, token, amount, tokenId,
// tokenType, unfreezeTimestamp)) — the migration dedup key (_getFrozenStakeHash). abi.encode
// pads every field to a 32-byte word.
func FrozenStakeHash(st Stake, unfreeze *big.Int) common.Hash {
	buf := make([]byte, 0, 7*32)
	buf = append(buf, common.BytesToHash(st.Validator.Bytes()).Bytes()...)
	buf = append(buf, common.BytesToHash(st.Delegator.Bytes()).Bytes()...)
	buf = append(buf, common.BytesToHash(st.Token.Bytes()).Bytes()...)
	buf = append(buf, leftPad32(st.Amount)...)
	buf = append(buf, leftPad32(st.TokenID)...)
	buf = append(buf, leftPad32(big.NewInt(int64(st.TokenType)))...)
	buf = append(buf, leftPad32(unfreeze)...)
	return common.BytesToHash(crypto.Keccak256(buf))
}

// ClassifyUnexplained tests candidate mappings against the unexplained slot set to identify
// what the leftover storage is (e.g. the migration _frozenStakeMigrated bool mapping).
func ClassifyUnexplained(unexplained []SlotValue, r Result) []ClassMatch {
	set := make(map[common.Hash]bool, len(unexplained))
	for _, sv := range unexplained {
		set[sv.Slot] = true
	}

	frozen := append(append([]FrozenStake{}, r.FrozenDeprecated...), r.FrozenLive...)
	probeFrozen := func(name string, field common.Hash) ClassMatch {
		matched := 0
		for _, f := range frozen {
			if set[mappingSlotBytes32(FrozenStakeHash(f.Stake, f.UnfreezeTimestamp), field)] {
				matched++
			}
		}
		return ClassMatch{Name: name, FieldSlot: field.Hex(), Matched: matched, Probed: len(frozen)}
	}

	return []ClassMatch{
		probeFrozen("_frozenStakeMigrated[frozenStakeHash] (computed ns)",
			ERC7201("DECIMAL_DELEGATION_MIGRATION_STORAGE_LOCATION")),
		probeFrozen("_frozenStakeMigrated[frozenStakeHash] (repo literal ns)",
			common.HexToHash("0x0087ae057eb38dff9d4851f20c0622c62ee2528e4d13ada7657e17acfe136200")),
	}
}

// hashedTokenID = keccak256(abi.encodePacked(token, tokenId)) (see _getHashedTokenId).
func hashedTokenID(token common.Address, tokenID *big.Int) common.Hash {
	buf := make([]byte, 0, 20+32)
	buf = append(buf, token.Bytes()...)
	buf = append(buf, leftPad32(tokenID)...)
	return common.BytesToHash(crypto.Keccak256(buf))
}
