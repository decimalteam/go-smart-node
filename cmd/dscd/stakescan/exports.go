package stakescan

// This file exposes a minimal exported API over the package's internal slot
// arithmetic so that external packages (e.g. an upgrade-time storage rewriter)
// can compute the exact storage slot to write for a given struct base, using
// the SAME offsets the decoder reads. None of the decoding logic lives here;
// these are thin wrappers around the unexported helpers in stakescan.go.

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// --- exported slot-arithmetic primitives (wrappers over the unexported forms) ---

// SlotAdd returns base + n (mod 2^256), the slot n words after base. It mirrors
// the internal slotAdd used throughout the decoder.
func SlotAdd(base common.Hash, n uint64) common.Hash { return slotAdd(base, n) }

// MappingSlotBytes32 returns the storage slot of mapping[key] for a
// mapping(bytes32 => ...) (or any 32-byte key) whose field slot is fieldSlot:
// keccak256(key ‖ fieldSlot).
func MappingSlotBytes32(key, fieldSlot common.Hash) common.Hash {
	return mappingSlotBytes32(key, fieldSlot)
}

// MappingSlotString returns the storage slot of mapping[key] for a
// mapping(string => ...) whose field slot is fieldSlot: keccak256(key ‖ fieldSlot).
func MappingSlotString(key string, fieldSlot common.Hash) common.Hash {
	return mappingSlotString(key, fieldSlot)
}

// ArrayElemStart returns the slot of element 0 of a Solidity dynamic array whose
// length word is at fieldSlot: keccak256(fieldSlot).
func ArrayElemStart(fieldSlot common.Hash) common.Hash { return arrayElemStart(fieldSlot) }

// --- exported DEL-amount slot accessors ---
//
// These compute the exact slot holding the DEL token amount inside a struct
// whose base slot is `base`, matching the +N offsets the decoder reads. The
// offsets are pinned to the decoder so a caller can rewrite an amount without
// re-deriving the layout:
//   - CoinStake.amount      -> base+3   (decodeStakeAt)
//   - NFTStake.amount       -> base+2   (ReconstructNFTStakes)
//   - NFTStake.reserveAmount-> base+6   (ReconstructNFTStakes)
//   - CheckDetails.amount   -> detailsBase+1 (ReconstructChecks)

// StakeAmountSlot returns the slot holding CoinStake.amount for a Stake struct
// whose base slot is base. Matches decodeStakeAt (amount at base+3).
func StakeAmountSlot(base common.Hash) common.Hash { return slotAdd(base, 3) }

// NFTStakeAmountSlot returns the slot holding NFTStake.amount for an NFTStake
// struct whose base slot is base. Matches ReconstructNFTStakes (amount at base+2).
func NFTStakeAmountSlot(base common.Hash) common.Hash { return slotAdd(base, 2) }

// NFTStakeReserveAmountSlot returns the slot holding NFTStake.reserveAmount (the
// reserved DEL backing the NFT stake) for an NFTStake struct whose base slot is
// base. Matches ReconstructNFTStakes (reserveAmount at base+6).
func NFTStakeReserveAmountSlot(base common.Hash) common.Hash { return slotAdd(base, 6) }

// CheckAmountSlot returns the slot holding CheckDetails.amount for a CheckDetails
// struct whose base slot is detailsBase. detailsBase is the value-mapping slot
// of _checkDetails[checkDetailsHash], i.e. MappingSlotBytes32(checkDetailsHash,
// SlotAdd(ChecksBase, 2)). Matches ReconstructChecks (amount at detailsBase+1).
func CheckAmountSlot(detailsBase common.Hash) common.Hash { return slotAdd(detailsBase, 1) }

// CheckDetailsBase returns the value-mapping slot of _checkDetails[checkDetailsHash]
// (the base of the CheckDetails struct), i.e. MappingSlotBytes32(checkDetailsHash,
// SlotAdd(ChecksBase, fCheckDetails)). Matches ReconstructChecks's `db`.
func CheckDetailsBase(checkDetailsHash common.Hash) common.Hash {
	return mappingSlotBytes32(checkDetailsHash, slotAdd(ChecksBase, fCheckDetails))
}

// --- frozen-array element slot math (mirrors ReadFrozenStakes / ReadLegacyFrozenStakes) ---
//
// The frozen amount lives at element+3 inside each FrozenStake's embedded Stake
// (decodeStakeAt: amount at +3), regardless of stride. The element base is
// ArrayElemStart(fieldSlot) + index*stride. These accessors expose the field
// offsets (fFrozen=9, fFrozenDep=3) and strides (10, 9) so a caller does not
// hardcode them.

// FrozenArrayFieldSlot returns the slot of the live _frozenStakes array length word
// (DelegationBase+fFrozen). Elements start at ArrayElemStart of this slot.
func FrozenArrayFieldSlot(base common.Hash) common.Hash { return slotAdd(base, fFrozen) }

// LegacyFrozenArrayFieldSlot returns the slot of the deprecated _frozenStakesDeprecated
// array length word (DelegationBase+fFrozenDep).
func LegacyFrozenArrayFieldSlot(base common.Hash) common.Hash { return slotAdd(base, fFrozenDep) }

// FrozenStride is the slot count of one live FrozenStake element.
const FrozenStride = frozenStride

// LegacyFrozenStride is the slot count of one deprecated FrozenStake element.
const LegacyFrozenStride = legacyFrozenStride

// FrozenStakeAmountSlot returns the slot holding the embedded Stake.amount of the
// frozen element at `index`. When deprecated is true it uses the legacy
// (_frozenStakesDeprecated) layout/stride, otherwise the live (_frozenStakes) one.
// The amount is always at element+3 (decodeStakeAt / decodeLegacyStakeAt).
func FrozenStakeAmountSlot(base common.Hash, index uint64, deprecated bool) common.Hash {
	field := slotAdd(base, fFrozen)
	stride := uint64(frozenStride)
	if deprecated {
		field = slotAdd(base, fFrozenDep)
		stride = legacyFrozenStride
	}
	elem := slotAdd(arrayElemStart(field), index*stride)
	return slotAdd(elem, 3)
}

// --- auto-unbond queue element slot math (mirrors ReadAutoUnbondQueueAt) ---

// AutoUnbondStride is the slot count of one AutoUnbondEntry element
// (validator, delegator, amount, token, holdTimestamp).
const AutoUnbondStride = autoUnbondStride

// AutoUnbondAmountSlot returns the slot holding AutoUnbondEntry.amount for the entry
// at `index` of an AutoUnbondEntry[] array whose length word is at fieldSlot. The
// element base is ArrayElemStart(fieldSlot) + index*5; amount is at element+2.
func AutoUnbondAmountSlot(fieldSlot common.Hash, index uint64) common.Hash {
	elem := slotAdd(arrayElemStart(fieldSlot), index*autoUnbondStride)
	return slotAdd(elem, 2)
}

// --- per-validator token reserve (ValidatorReserve) slot math (mirrors coverage.go) ---

// ValidatorReserveBase returns the base slot of
// _validatorTokens[validator][hashedTokenID(token, tokenId)] — a ValidatorReserve
// {penaltyIndex@+0, reserve@+1} (DecimalDelegationCommon.sol _getValidatorReserve):
// keccak256(hashedTokenID ‖ keccak256(leftPad32(validator) ‖ base+5)), with
// hashedTokenID = keccak256(abi.encodePacked(token, tokenId)) (_getHashedTokenId).
func ValidatorReserveBase(base common.Hash, validator, token common.Address, tokenID *big.Int) common.Hash {
	mid := mappingSlotBytes32(common.BytesToHash(validator.Bytes()), slotAdd(base, fValidatorTokens))
	return mappingSlotBytes32(hashedTokenID(token, tokenID), mid)
}

// ValidatorReserveAmountSlot returns the slot holding ValidatorReserve.reserve —
// the running per-(validator, token) stake aggregate maintained by
// _addValidatorReserve / _reduceValidatorReserve — i.e. ValidatorReserveBase+1.
// For a DEL stake the key token is WDEL, so the stored reserve is a DEL amount.
func ValidatorReserveAmountSlot(base common.Hash, validator, token common.Address, tokenID *big.Int) common.Hash {
	return slotAdd(ValidatorReserveBase(base, validator, token, tokenID), 1)
}

// --- WDEL / ERC20 mapping(address => uint256) balance slot math ---

// WDELBalancesField is the storage slot of WDEL's `balanceOf` mapping. WDEL is the
// canonical WETH9 wrapper (Solidity 0.4.18, flat sequential storage): slot 0 = name,
// 1 = symbol, 2 = decimals, 3 = balanceOf (mapping(address=>uint)), 4 = allowance
// (mapping(address=>mapping(address=>uint))). Verified against mainnet genesis: the
// keccak slot for balanceOf[delegation] holds the delegation contract's WDEL balance.
// totalSupply() is NOT stored — WETH9 returns address(this).balance — so there is no
// totalSupply slot to rewrite.
const WDELBalancesField uint64 = 3

// MappingSlotAddress returns the storage slot of mapping[addr] for a
// mapping(address => ...) whose field slot index is field: keccak256(leftPad32(addr) ‖
// leftPad32(field)). This matches Solidity's mapping slot derivation for a top-level
// (non-namespaced) state variable at the given declaration slot.
func MappingSlotAddress(addr common.Address, field uint64) common.Hash {
	key := common.LeftPadBytes(addr.Bytes(), 32)
	fld := common.BigToHash(new(big.Int).SetUint64(field)).Bytes()
	return common.BytesToHash(crypto.Keccak256(key, fld))
}

// WDELBalanceSlot returns the slot holding WDEL.balanceOf[holder]:
// keccak256(leftPad32(holder) ‖ leftPad32(3)).
func WDELBalanceSlot(holder common.Address) common.Hash {
	return MappingSlotAddress(holder, WDELBalancesField)
}

// --- per-collection NFT reserve (DRC721 / DRC1155 NFTReserve) slot math ---

// NFTReserveBase is the ERC-7201 storage base of the NFTReserve mixin embedded in every
// DRC721 / DRC1155 collection contract:
// keccak256(abi.encode(uint256(keccak256("DECIMAL_NFT_RESERVE_STORAGE_LOCATION")) - 1)) & ~0xff.
// Verified equal to the literal in nft-center/contracts/NFTReserve.sol. The storage struct
// is { mapping(uint256 tokenId => Reserve) _reserve;  bool _refundable; }, so _reserve is at
// field 0 (slot == base). Reserve is { address token; uint256 amount; ReserveType reserveType }
// → per-entry token@+0, amount@+1, reserveType@+2 (ReserveType: None=0, DEL=1, DRC20=2).
var NFTReserveBase = ERC7201("DECIMAL_NFT_RESERVE_STORAGE_LOCATION")

// NFTReserveEntryBase returns the base slot of _reserve[tokenId] (a Reserve struct) for a
// DRC721/DRC1155 collection: MappingSlotUint256(tokenId, NFTReserveBase).
func NFTReserveEntryBase(tokenID *big.Int) common.Hash {
	return MappingSlotUint256(tokenID, NFTReserveBase)
}

// NFTReserveAmountSlot returns the slot holding _reserve[tokenId].amount (the DEL reserve
// backing the NFT) for a DRC721/DRC1155 collection: NFTReserveEntryBase(tokenId)+1.
func NFTReserveAmountSlot(tokenID *big.Int) common.Hash {
	return slotAdd(NFTReserveEntryBase(tokenID), 1)
}

// NFTReserveTypeSlot returns the slot holding _reserve[tokenId].reserveType:
// NFTReserveEntryBase(tokenId)+2.
func NFTReserveTypeSlot(tokenID *big.Int) common.Hash {
	return slotAdd(NFTReserveEntryBase(tokenID), 2)
}

// NFTReserveTokenSlot returns the slot holding _reserve[tokenId].token:
// NFTReserveEntryBase(tokenId)+0.
func NFTReserveTokenSlot(tokenID *big.Int) common.Hash {
	return NFTReserveEntryBase(tokenID)
}

// MappingSlotUint256 returns the storage slot of mapping[key] for a mapping(uint256 => ...)
// whose field slot is fieldSlot: keccak256(leftPad32(key) ‖ fieldSlot).
func MappingSlotUint256(key *big.Int, fieldSlot common.Hash) common.Hash {
	return mappingSlotBytes32(common.BigToHash(key), fieldSlot)
}
