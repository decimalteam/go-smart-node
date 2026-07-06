package redenom

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"

	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
)

type noopLogger struct{}

func (noopLogger) Info(string, ...interface{}) {}

func addrOf(b byte) common.Address {
	var a common.Address
	for i := range a {
		a[i] = b
	}
	return a
}

// writeNFTStake lays out one NFTStake record in the delegation storage exactly as the
// Solidity contract does (mirrors cmd/dscd/stakescan TestReconstructNFTStakes).
func writeNFTStake(s stakescan.Storage, nftContract, delegator, validator, reserveToken common.Address, tokenID, reserveAmount *big.Int) {
	nftField := stakescan.SlotAdd(stakescan.DelegationBase, 6)
	id := stakescan.NFTStakeID(nftContract, tokenID, delegator, validator, big.NewInt(0))
	base := stakescan.MappingSlotBytes32(id, nftField)

	put := func(slot common.Hash, v *big.Int) { s[slot] = common.BigToHash(v) }
	putAddr := func(slot common.Hash, a common.Address) { s[slot] = common.BytesToHash(a.Bytes()) }

	putAddr(base, nftContract)
	put(stakescan.SlotAdd(base, 1), tokenID)
	put(stakescan.SlotAdd(base, 2), big.NewInt(1)) // intrinsic amount
	// slot 3 packs nftType(off0,1B) ‖ delegator(off1,20B)
	packed3 := new(big.Int).Lsh(new(big.Int).SetBytes(delegator.Bytes()), 8)
	packed3.Or(packed3, big.NewInt(2)) // nftType = DRC721 = 2
	put(stakescan.SlotAdd(base, 3), packed3)
	putAddr(stakescan.SlotAdd(base, 4), validator)
	putAddr(stakescan.SlotAdd(base, 5), reserveToken)
	put(stakescan.SlotAdd(base, 6), reserveAmount)
	put(stakescan.SlotAdd(base, 8), big.NewInt(1)) // isActive
}

// TestCollectDelegationWrites_NFTReserveGate is the H2 regression: an NFT stake whose
// reserve is custom-coin-denominated (reserveToken != wdel) must NOT have its reserve
// amount divided; only DEL reserves (reserveToken == wdel) are scaled.
func TestCollectDelegationWrites_NFTReserveGate(t *testing.T) {
	deleg := addrOf(0x11)
	wdel := addrOf(0x22)
	customCoin := addrOf(0x33)

	store := stakescan.Storage{}
	// DEL-backed NFT stake: reserveToken == wdel -> scaled.
	writeNFTStake(store, addrOf(0xa1), addrOf(0xb1), addrOf(0xc1), wdel,
		big.NewInt(1), big.NewInt(5000))
	// Custom-coin-backed NFT stake: reserveToken == customCoin -> left untouched.
	writeNFTStake(store, addrOf(0xa2), addrOf(0xb2), addrOf(0xc2), customCoin,
		big.NewInt(2), big.NewInt(7000))

	rep := newReport(sdkmath.NewInt(1000))
	holders := map[common.Address]struct{}{}
	writes, err := collectDelegationWrites(store, deleg, wdel, big.NewInt(1000), &rep, noopLogger{}, holders)
	if err != nil {
		t.Fatalf("collectDelegationWrites: %v", err)
	}

	// Exactly one NFT reserve slot scaled (the DEL one), reported once.
	if rep.EVMNFTReserveSlots != 1 {
		t.Fatalf("EVMNFTReserveSlots = %d, want 1 (custom reserve must be skipped)", rep.EVMNFTReserveSlots)
	}

	// Identify the expected scaled slot for the DEL-backed stake.
	delID := stakescan.NFTStakeID(addrOf(0xa1), big.NewInt(1), addrOf(0xb1), addrOf(0xc1), big.NewInt(0))
	delBase := stakescan.MappingSlotBytes32(delID, stakescan.SlotAdd(stakescan.DelegationBase, 6))
	wantSlot := stakescan.NFTStakeReserveAmountSlot(delBase)

	customID := stakescan.NFTStakeID(addrOf(0xa2), big.NewInt(2), addrOf(0xb2), addrOf(0xc2), big.NewInt(0))
	customBase := stakescan.MappingSlotBytes32(customID, stakescan.SlotAdd(stakescan.DelegationBase, 6))
	customSlot := stakescan.NFTStakeReserveAmountSlot(customBase)

	var found bool
	for _, w := range writes {
		if w.slot == customSlot {
			t.Fatalf("custom-coin NFT reserve slot %s was scaled (must be skipped)", customSlot.Hex())
		}
		if w.slot == wantSlot {
			found = true
			if w.newVal.Cmp(big.NewInt(5)) != 0 { // floor(5000/1000)
				t.Errorf("DEL reserve scaled to %s, want 5", w.newVal)
			}
		}
	}
	if !found {
		t.Fatalf("DEL-backed NFT reserve slot %s was not scaled", wantSlot.Hex())
	}
}
