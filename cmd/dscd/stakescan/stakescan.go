// Package stakescan reconstructs DecimalDelegation stakes from raw EVM contract
// storage. It is pure (no cosmos/DB dependencies): the input is a slot->value
// map and the output is decoded, cryptographically-verified stakes.
//
// Stakes live in Solidity mappings keyed by bytes32 stakeId, which are not
// natively enumerable. We instead scan every stored slot as a candidate struct
// base, decode the struct, recompute the stakeId from its own fields, and verify
// by re-deriving the mapping slot (keccak256(id ‖ fieldSlot)) and matching the
// candidate slot. No keccak inversion is needed and there are no false positives.
//
// Storage layouts and id formulas are documented in
// docs/superpowers/specs/2026-06-22-delegation-stakes-reconstruction-design.md
// and were confirmed against solc's storageLayout output.
package stakescan

import (
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ERC-7201 storage bases (independent of contract address / network).
var (
	// DelegationBase is DecimalDelegationCommon.DECIMAL_DELEGATION_COMMON_STORAGE_LOCATION.
	DelegationBase = common.HexToHash("0xc1dae510251b57b62087f142cc8746564a134308f8810580b66b08a29b6def00")
	// ContractCenterBase is DecimalContractCenter.DECIMAL_CONTRACT_CENTER_STORAGE_LOCATION.
	ContractCenterBase = common.HexToHash("0x00e2f18f4aa7606d53e2e7b6172fb5f051be6ee8f80cf63eec8ca91bb44b3100")
)

// Field offsets within DelegationStorage (slot = DelegationBase + offset).
const (
	fStakes       = 1 // mapping(bytes32 => Stake)
	fStakePenalty = 2 // mapping(bytes32 => uint256)
	fFrozenDep    = 3 // FrozenStake[] (legacy)
	fNFTStakes    = 6 // mapping(bytes32 => NFTStake)
	fNFTBacked    = 7 // mapping(bytes32 => uint256)
	fNFTPenalty   = 8 // mapping(bytes32 => uint256)
	fFrozen       = 9 // FrozenStake[] (live)
)

// frozenStride is the slot count of one FrozenStake element (8-slot Stake + packed
// status/type slot + unfreezeTimestamp).
const frozenStride = 10

// legacyFrozenStride is the slot count of one LegacyFrozenStake element, used by the
// deprecated _frozenStakesDeprecated array (base+3). The legacy Stake had 7 fields
// (no holdStartTime), so the element is 7-slot stake + packed status/type + unfreezeTimestamp.
const legacyFrozenStride = 9

// maxFrozenLen guards against a corrupt length word causing an unbounded loop.
const maxFrozenLen = 10_000_000

var max256 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

// Storage is raw contract storage: slot -> 32-byte value. Missing slots are zero
// (ethermint deletes zero-valued slots).
type Storage map[common.Hash]common.Hash

// Stake mirrors IDecimalDelegationCommon.Stake.
type Stake struct {
	StakeID       common.Hash    `json:"stakeId"`
	Validator     common.Address `json:"validator"`
	Delegator     common.Address `json:"delegator"`
	Token         common.Address `json:"token"`
	Amount        *big.Int       `json:"amount"`
	TokenID       *big.Int       `json:"tokenId"`
	TokenType     uint8          `json:"tokenType"`
	HoldTimestamp *big.Int       `json:"holdTimestamp"`
	HoldStartTime *big.Int       `json:"holdStartTime"`
}

// CoinStake is a reconstructed entry of the _stakes mapping.
type CoinStake struct {
	Stake
	PenaltyIndex *big.Int    `json:"penaltyIndex"`
	NFTBacked    *big.Int    `json:"nftBackedAmount"`
	BaseSlot     common.Hash `json:"baseSlot"`
}

// NFTStake is a reconstructed entry of the _nftStakes mapping.
type NFTStake struct {
	NFTStakeID    common.Hash    `json:"nftStakeId"`
	NFTContract   common.Address `json:"nftContract"`
	TokenID       *big.Int       `json:"tokenId"`
	Amount        *big.Int       `json:"amount"`
	NFTType       uint8          `json:"nftType"`
	Delegator     common.Address `json:"delegator"`
	Validator     common.Address `json:"validator"`
	ReserveToken  common.Address `json:"reserveToken"`
	ReserveAmount *big.Int       `json:"reserveAmount"`
	HoldTimestamp *big.Int       `json:"holdTimestamp"`
	IsActive      bool           `json:"isActive"`
	PenaltyIndex  *big.Int       `json:"penaltyIndex"`
	BaseSlot      common.Hash    `json:"baseSlot"`
}

// FrozenStake is a reconstructed element of a FrozenStake[] array.
type FrozenStake struct {
	Index             uint64   `json:"index"`
	Stake             Stake    `json:"stake"`
	FreezeStatus      uint8    `json:"freezeStatus"`
	FreezeType        uint8    `json:"freezeType"`
	UnfreezeTimestamp *big.Int `json:"unfreezeTimestamp"`
}

// Result aggregates everything reconstructed for one delegation contract.
type Result struct {
	CoinStakes       []CoinStake   `json:"coinStakes"`
	NFTStakes        []NFTStake    `json:"nftStakes"`
	FrozenLive       []FrozenStake `json:"frozenStakes"`
	FrozenDeprecated []FrozenStake `json:"frozenStakesDeprecated"`
	ScannedSlots     int           `json:"scannedSlots"`
}

// Reconstruct runs the full coin/NFT/frozen reconstruction for a delegation
// contract whose DelegationStorage base is `base` (normally DelegationBase).
func Reconstruct(s Storage, base common.Hash) Result {
	return Result{
		CoinStakes:       ReconstructCoinStakes(s, base),
		NFTStakes:        ReconstructNFTStakes(s, base),
		FrozenLive:       ReadFrozenStakes(s, base, fFrozen),
		FrozenDeprecated: ReadLegacyFrozenStakes(s, base, fFrozenDep),
		ScannedSlots:     len(s),
	}
}

// ReconstructCoinStakes scans the storage for verified _stakes entries.
func ReconstructCoinStakes(s Storage, base common.Hash) []CoinStake {
	stakesSlot := slotAdd(base, fStakes)
	penaltySlot := slotAdd(base, fStakePenalty)
	nftBackedSlot := slotAdd(base, fNFTBacked)

	var out []CoinStake
	for k, v := range s {
		if !isAddressShaped(v) {
			continue
		}
		st := decodeStakeAt(s, k)
		if mappingSlotBytes32(st.StakeID, stakesSlot) != k {
			continue
		}
		out = append(out, CoinStake{
			Stake:        st,
			BaseSlot:     k,
			PenaltyIndex: s.uintAt(mappingSlotBytes32(st.StakeID, penaltySlot)),
			NFTBacked:    s.uintAt(mappingSlotBytes32(st.StakeID, nftBackedSlot)),
		})
	}
	return out
}

// ReconstructNFTStakes scans the storage for verified _nftStakes entries.
func ReconstructNFTStakes(s Storage, base common.Hash) []NFTStake {
	nftSlot := slotAdd(base, fNFTStakes)
	nftPenaltySlot := slotAdd(base, fNFTPenalty)

	var out []NFTStake
	for k, v := range s {
		if !isAddressShaped(v) {
			continue
		}
		nftContract := common.BytesToAddress(v.Bytes())
		tokenID := s.uintAt(slotAdd(k, 1))
		amount := s.uintAt(slotAdd(k, 2))
		packed3 := s[slotAdd(k, 3)]
		nftType := fieldUint(packed3, 0, 1).Uint64()
		delegator := fieldAddr(packed3, 1)
		validator := s.addrAt(slotAdd(k, 4))
		reserveToken := s.addrAt(slotAdd(k, 5))
		reserveAmount := s.uintAt(slotAdd(k, 6))
		holdTimestamp := s.uintAt(slotAdd(k, 7))
		isActive := fieldUint(s[slotAdd(k, 8)], 0, 1).Sign() != 0

		id := NFTStakeID(nftContract, tokenID, delegator, validator, holdTimestamp)
		if mappingSlotBytes32(id, nftSlot) != k {
			continue
		}
		out = append(out, NFTStake{
			NFTStakeID: id, NFTContract: nftContract, TokenID: tokenID, Amount: amount,
			NFTType: uint8(nftType), Delegator: delegator, Validator: validator,
			ReserveToken: reserveToken, ReserveAmount: reserveAmount, HoldTimestamp: holdTimestamp,
			IsActive:     isActive,
			PenaltyIndex: s.uintAt(mappingSlotBytes32(id, nftPenaltySlot)),
			BaseSlot:     k,
		})
	}
	return out
}

// ReadFrozenStakes iterates the live _frozenStakes FrozenStake[] array (current
// 10-slot layout: 8-slot Stake, packed status/type at +8, unfreezeTimestamp at +9).
func ReadFrozenStakes(s Storage, base common.Hash, fieldOffset uint64) []FrozenStake {
	fieldSlot := slotAdd(base, fieldOffset)
	n, ok := arrayLen(s, fieldSlot)
	if !ok {
		return nil
	}
	start := arrayElemStart(fieldSlot)

	out := make([]FrozenStake, 0, n)
	for i := uint64(0); i < n; i++ {
		elem := slotAdd(start, i*frozenStride)
		packed8 := s[slotAdd(elem, 8)]
		out = append(out, FrozenStake{
			Index:             i,
			Stake:             decodeStakeAt(s, elem),
			FreezeStatus:      uint8(fieldUint(packed8, 0, 1).Uint64()),
			FreezeType:        uint8(fieldUint(packed8, 1, 1).Uint64()),
			UnfreezeTimestamp: s.uintAt(slotAdd(elem, 9)),
		})
	}
	return out
}

// ReadLegacyFrozenStakes iterates the deprecated _frozenStakesDeprecated array
// (legacy 9-slot layout: 7-slot LegacyStake without holdStartTime, packed
// status/type at +7, unfreezeTimestamp at +8). On networks where the layout
// migration has not run, this is where the real frozen withdraw/transfer queue lives.
func ReadLegacyFrozenStakes(s Storage, base common.Hash, fieldOffset uint64) []FrozenStake {
	fieldSlot := slotAdd(base, fieldOffset)
	n, ok := arrayLen(s, fieldSlot)
	if !ok {
		return nil
	}
	start := arrayElemStart(fieldSlot)

	out := make([]FrozenStake, 0, n)
	for i := uint64(0); i < n; i++ {
		elem := slotAdd(start, i*legacyFrozenStride)
		packed7 := s[slotAdd(elem, 7)]
		out = append(out, FrozenStake{
			Index:             i,
			Stake:             decodeLegacyStakeAt(s, elem),
			FreezeStatus:      uint8(fieldUint(packed7, 0, 1).Uint64()),
			FreezeType:        uint8(fieldUint(packed7, 1, 1).Uint64()),
			UnfreezeTimestamp: s.uintAt(slotAdd(elem, 8)),
		})
	}
	return out
}

// arrayLen reads a Solidity dynamic-array length word, guarding against a corrupt
// value that would otherwise drive an unbounded loop.
func arrayLen(s Storage, fieldSlot common.Hash) (uint64, bool) {
	length := s.uintAt(fieldSlot)
	if !length.IsUint64() {
		return 0, false
	}
	n := length.Uint64()
	if n == 0 || n > maxFrozenLen {
		return 0, false
	}
	return n, true
}

// autoUnbondStride is the slot count of one AutoUnbondEntry element
// (validator, delegator, amount, token, holdTimestamp — each a full slot).
const autoUnbondStride = 5

// AutoUnbondQueueSlotMain is the ERC-7201 namespace base under which the deployed
// DSC-main contract stores the AutoUnbondEntry[] queue (the array length is at this
// slot, elements at keccak256(slot)+i*5). It is NOT in any Solidity source — it was
// located empirically by FindAutoUnbondQueues and recorded here for reference. The
// command relies on the locator (robust across networks/upgrades), not this constant.
var AutoUnbondQueueSlotMain = common.HexToHash("0xd92fb98890f4a927a2850eac803b570950d9a3cb4225c0de7badfaf82a9fc800")

// AutoUnbondEntry mirrors IDecimalDelegation.AutoUnbondEntry: one queued
// auto-unbond (force-withdrawal) created by the validator module's
// ExecuteAutoUnbondEnqueue calling autoUnbondEnqueue on the contract.
type AutoUnbondEntry struct {
	Index         uint64         `json:"index"`
	Validator     common.Address `json:"validator"`
	Delegator     common.Address `json:"delegator"`
	Amount        *big.Int       `json:"amount"`
	Token         common.Address `json:"token"`
	HoldTimestamp *big.Int       `json:"holdTimestamp"`
}

// AutoUnbondQueue is a located AutoUnbondEntry[] dynamic array.
type AutoUnbondQueue struct {
	FieldSlot    common.Hash       `json:"fieldSlot"`
	Length       uint64            `json:"length"`
	ValidSampled int               `json:"validSampled"`
	TotalSampled int               `json:"totalSampled"`
	Entries      []AutoUnbondEntry `json:"entries"`
}

// ReadAutoUnbondQueueAt reads an AutoUnbondEntry[] dynamic array whose length word
// is at fieldSlot (elements at keccak256(fieldSlot) + i*5).
func ReadAutoUnbondQueueAt(s Storage, fieldSlot common.Hash) []AutoUnbondEntry {
	n, ok := arrayLen(s, fieldSlot)
	if !ok {
		return nil
	}
	start := arrayElemStart(fieldSlot)
	out := make([]AutoUnbondEntry, 0, n)
	for i := uint64(0); i < n; i++ {
		elem := slotAdd(start, i*autoUnbondStride)
		out = append(out, AutoUnbondEntry{
			Index:         i,
			Validator:     s.addrAt(elem),
			Delegator:     s.addrAt(slotAdd(elem, 1)),
			Amount:        s.uintAt(slotAdd(elem, 2)),
			Token:         s.addrAt(slotAdd(elem, 3)),
			HoldTimestamp: s.uintAt(slotAdd(elem, 4)),
		})
	}
	return out
}

// FindAutoUnbondQueues locates AutoUnbondEntry[] arrays by scanning every stored slot
// as a candidate length word and checking that its elements decode as entries with
// address-shaped validator/delegator/token. The deployed contract's queue slot is not
// in any Solidity source, so it is detected rather than hard-coded. minValidFraction is
// the share of sampled elements that must look valid (1.0 = all). Results are sorted by
// descending length.
func FindAutoUnbondQueues(s Storage, minValidFraction float64) []AutoUnbondQueue {
	var out []AutoUnbondQueue
	for f, v := range s {
		L := new(big.Int).SetBytes(v.Bytes())
		if !L.IsUint64() {
			continue
		}
		n := L.Uint64()
		if n == 0 || n > maxFrozenLen {
			continue
		}
		start := arrayElemStart(f)
		valid, sampled := sampleAutoUnbondValidity(s, start, n)
		if sampled == 0 || float64(valid)/float64(sampled) < minValidFraction {
			continue
		}
		out = append(out, AutoUnbondQueue{
			FieldSlot:    f,
			Length:       n,
			ValidSampled: valid,
			TotalSampled: sampled,
			Entries:      ReadAutoUnbondQueueAt(s, f),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Length > out[j].Length })
	return out
}

func sampleAutoUnbondValidity(s Storage, start common.Hash, n uint64) (valid, sampled int) {
	step := uint64(1)
	if n > 64 {
		step = n / 64
	}
	for i := uint64(0); i < n; i += step {
		elem := slotAdd(start, i*autoUnbondStride)
		sampled++
		if isAddressShaped(s[elem]) && isAddressShaped(s[slotAdd(elem, 1)]) && isAddressShaped(s[slotAdd(elem, 3)]) {
			valid++
		}
	}
	return valid, sampled
}

// ResolveAddressBySymbol reads ContractCenter._contractAddresses[symbol] (a
// mapping(string => address) at offset 0 of ContractCenterBase).
func ResolveAddressBySymbol(cc Storage, symbol string) common.Address {
	return cc.addrAt(mappingSlotString(symbol, ContractCenterBase))
}

// CoinStakeID = keccak256(abi.encodePacked(validator, delegator, token, tokenId, holdTimestamp)).
func CoinStakeID(validator, delegator, token common.Address, tokenID, holdTimestamp *big.Int) common.Hash {
	buf := make([]byte, 0, 20*3+32*2)
	buf = append(buf, validator.Bytes()...)
	buf = append(buf, delegator.Bytes()...)
	buf = append(buf, token.Bytes()...)
	buf = append(buf, leftPad32(tokenID)...)
	buf = append(buf, leftPad32(holdTimestamp)...)
	return common.BytesToHash(crypto.Keccak256(buf))
}

// NFTStakeID = keccak256(abi.encodePacked("NFT", nftContract, tokenId, delegator, validator, holdTimestamp)).
func NFTStakeID(nftContract common.Address, tokenID *big.Int, delegator, validator common.Address, holdTimestamp *big.Int) common.Hash {
	buf := make([]byte, 0, 3+20*3+32*2)
	buf = append(buf, []byte("NFT")...)
	buf = append(buf, nftContract.Bytes()...)
	buf = append(buf, leftPad32(tokenID)...)
	buf = append(buf, delegator.Bytes()...)
	buf = append(buf, validator.Bytes()...)
	buf = append(buf, leftPad32(holdTimestamp)...)
	return common.BytesToHash(crypto.Keccak256(buf))
}

// decodeStakeAt decodes an 8-slot Stake struct starting at base and recomputes its id.
func decodeStakeAt(s Storage, base common.Hash) Stake {
	validator := s.addrAt(base)
	delegator := s.addrAt(slotAdd(base, 1))
	token := s.addrAt(slotAdd(base, 2))
	amount := s.uintAt(slotAdd(base, 3))
	tokenID := s.uintAt(slotAdd(base, 4))
	tokenType := s.uintAt(slotAdd(base, 5)).Uint64()
	holdTimestamp := s.uintAt(slotAdd(base, 6))
	holdStartTime := s.uintAt(slotAdd(base, 7))
	return Stake{
		StakeID:       CoinStakeID(validator, delegator, token, tokenID, holdTimestamp),
		Validator:     validator,
		Delegator:     delegator,
		Token:         token,
		Amount:        amount,
		TokenID:       tokenID,
		TokenType:     uint8(tokenType),
		HoldTimestamp: holdTimestamp,
		HoldStartTime: holdStartTime,
	}
}

// decodeLegacyStakeAt decodes a 7-slot LegacyStake (no holdStartTime) and recomputes
// its id (the id formula does not include holdStartTime, so it is identical).
func decodeLegacyStakeAt(s Storage, base common.Hash) Stake {
	validator := s.addrAt(base)
	delegator := s.addrAt(slotAdd(base, 1))
	token := s.addrAt(slotAdd(base, 2))
	amount := s.uintAt(slotAdd(base, 3))
	tokenID := s.uintAt(slotAdd(base, 4))
	tokenType := s.uintAt(slotAdd(base, 5)).Uint64()
	holdTimestamp := s.uintAt(slotAdd(base, 6))
	return Stake{
		StakeID:       CoinStakeID(validator, delegator, token, tokenID, holdTimestamp),
		Validator:     validator,
		Delegator:     delegator,
		Token:         token,
		Amount:        amount,
		TokenID:       tokenID,
		TokenType:     uint8(tokenType),
		HoldTimestamp: holdTimestamp,
		HoldStartTime: big.NewInt(0),
	}
}

// --- storage decoding helpers ---

func (s Storage) addrAt(slot common.Hash) common.Address {
	return common.BytesToAddress(s[slot].Bytes()) // low 20 bytes
}

func (s Storage) uintAt(slot common.Hash) *big.Int {
	return new(big.Int).SetBytes(s[slot].Bytes())
}

// isAddressShaped reports whether a slot word is a left-padded 20-byte address
// (high 12 bytes zero, low 20 non-zero).
func isAddressShaped(w common.Hash) bool {
	for i := 0; i < 12; i++ {
		if w[i] != 0 {
			return false
		}
	}
	for i := 12; i < 32; i++ {
		if w[i] != 0 {
			return true
		}
	}
	return false
}

// fieldUint extracts a packed integer field of sizeBytes at byte offsetBytes from
// the least-significant end of the slot word.
func fieldUint(w common.Hash, offsetBytes, sizeBytes uint) *big.Int {
	x := new(big.Int).SetBytes(w[:])
	x.Rsh(x, offsetBytes*8)
	mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), sizeBytes*8), big.NewInt(1))
	return x.And(x, mask)
}

func fieldAddr(w common.Hash, offsetBytes uint) common.Address {
	return common.BigToAddress(fieldUint(w, offsetBytes, 20))
}

// --- slot arithmetic ---

func slotAdd(base common.Hash, n uint64) common.Hash {
	x := new(big.Int).SetBytes(base[:])
	x.Add(x, new(big.Int).SetUint64(n))
	x.And(x, max256)
	return common.BigToHash(x)
}

func mappingSlotBytes32(key, fieldSlot common.Hash) common.Hash {
	return common.BytesToHash(crypto.Keccak256(key.Bytes(), fieldSlot.Bytes()))
}

func mappingSlotString(key string, fieldSlot common.Hash) common.Hash {
	return common.BytesToHash(crypto.Keccak256([]byte(key), fieldSlot.Bytes()))
}

func arrayElemStart(fieldSlot common.Hash) common.Hash {
	return common.BytesToHash(crypto.Keccak256(fieldSlot.Bytes()))
}

func leftPad32(x *big.Int) []byte {
	if x == nil {
		return make([]byte, 32)
	}
	return common.BigToHash(x).Bytes()
}
