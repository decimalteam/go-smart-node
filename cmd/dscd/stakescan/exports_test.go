package stakescan

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// TestWDELBalanceSlot pins the WDEL balanceOf slot derivation against a value verified
// against mainnet genesis: WDEL.balanceOf[delegation] lives at this exact slot (the one
// holding the delegation contract's ~12.08B WDEL balance). balanceOf is the WETH9 mapping
// at field 3, so slot = keccak256(leftPad32(holder) ‖ leftPad32(3)).
func TestWDELBalanceSlot(t *testing.T) {
	delegation := common.HexToAddress("0xA16c34Ed1c0601c0E749e17ebeF19752A15FaA01")
	want := common.HexToHash("0xf548706b1976977cf0b62749ac17be8daa224e815785e034b8275b4835fdf362")
	if got := WDELBalanceSlot(delegation); got != want {
		t.Fatalf("WDELBalanceSlot(delegation) = %s, want %s", got.Hex(), want.Hex())
	}
	if WDELBalancesField != 3 {
		t.Fatalf("WDELBalancesField = %d, want 3", WDELBalancesField)
	}
	// WDELBalanceSlot must equal the generic MappingSlotAddress at field 3.
	if WDELBalanceSlot(delegation) != MappingSlotAddress(delegation, 3) {
		t.Fatal("WDELBalanceSlot != MappingSlotAddress(_, 3)")
	}
}

// TestNFTReserveBase pins the per-collection NFTReserve ERC-7201 base to the literal in
// nft-center/contracts/NFTReserve.sol (DECIMAL_NFT_RESERVE_STORAGE_LOCATION).
func TestNFTReserveBase(t *testing.T) {
	want := common.HexToHash("0x1e5e5282ace2fee35cb1cd46e4746f781dcfdf434921d26f79802de610904f00")
	if NFTReserveBase != want {
		t.Fatalf("NFTReserveBase = %s, want %s", NFTReserveBase.Hex(), want.Hex())
	}
}

// TestNFTReserveSlots checks the Reserve struct field offsets (token@+0, amount@+1,
// reserveType@+2) relative to the _reserve[tokenId] entry base.
func TestNFTReserveSlots(t *testing.T) {
	tid := big.NewInt(123456)
	eb := NFTReserveEntryBase(tid)
	if eb != MappingSlotUint256(tid, NFTReserveBase) {
		t.Fatal("NFTReserveEntryBase != MappingSlotUint256(tid, NFTReserveBase)")
	}
	if NFTReserveTokenSlot(tid) != eb {
		t.Fatal("token slot != entryBase+0")
	}
	if NFTReserveAmountSlot(tid) != SlotAdd(eb, 1) {
		t.Fatal("amount slot != entryBase+1")
	}
	if NFTReserveTypeSlot(tid) != SlotAdd(eb, 2) {
		t.Fatal("type slot != entryBase+2")
	}
}

// TestFrozenStakeAmountSlot checks the frozen-array amount slot math for both the live
// (field 9, stride 10) and deprecated (field 3, stride 9) layouts; amount is at element+3.
func TestFrozenStakeAmountSlot(t *testing.T) {
	base := DelegationBase
	for _, idx := range []uint64{0, 1, 7, 4321} {
		// live
		elem := SlotAdd(ArrayElemStart(FrozenArrayFieldSlot(base)), idx*FrozenStride)
		if got, want := FrozenStakeAmountSlot(base, idx, false), SlotAdd(elem, 3); got != want {
			t.Fatalf("live FrozenStakeAmountSlot(%d) = %s, want %s", idx, got.Hex(), want.Hex())
		}
		// deprecated
		delem := SlotAdd(ArrayElemStart(LegacyFrozenArrayFieldSlot(base)), idx*LegacyFrozenStride)
		if got, want := FrozenStakeAmountSlot(base, idx, true), SlotAdd(delem, 3); got != want {
			t.Fatalf("deprecated FrozenStakeAmountSlot(%d) = %s, want %s", idx, got.Hex(), want.Hex())
		}
	}
	if FrozenStride != 10 || LegacyFrozenStride != 9 {
		t.Fatalf("strides: live=%d (want 10) legacy=%d (want 9)", FrozenStride, LegacyFrozenStride)
	}
}

// TestAutoUnbondAmountSlot checks the auto-unbond amount slot (element+2, stride 5).
func TestAutoUnbondAmountSlot(t *testing.T) {
	field := common.HexToHash("0xd92fb98890f4a927a2850eac803b570950d9a3cb4225c0de7badfaf82a9fc800")
	for _, idx := range []uint64{0, 1, 9, 1000} {
		elem := SlotAdd(ArrayElemStart(field), idx*AutoUnbondStride)
		if got, want := AutoUnbondAmountSlot(field, idx), SlotAdd(elem, 2); got != want {
			t.Fatalf("AutoUnbondAmountSlot(%d) = %s, want %s", idx, got.Hex(), want.Hex())
		}
	}
	if AutoUnbondStride != 5 {
		t.Fatalf("AutoUnbondStride = %d, want 5", AutoUnbondStride)
	}
}

// TestCheckDetailsBase checks the CheckDetails base and amount slot derivation.
func TestCheckDetailsBase(t *testing.T) {
	h := common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
	db := CheckDetailsBase(h)
	if db != MappingSlotBytes32(h, SlotAdd(ChecksBase, 2)) {
		t.Fatal("CheckDetailsBase != MappingSlotBytes32(hash, ChecksBase+2)")
	}
	if CheckAmountSlot(db) != SlotAdd(db, 1) {
		t.Fatal("CheckAmountSlot != detailsBase+1")
	}
}
