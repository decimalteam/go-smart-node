package stakescan_test

import (
	"math/big"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
	"github.com/ethereum/go-ethereum/common"
)

func TestReconstructChecks(t *testing.T) {
	C := stakescan.ChecksBase
	checksField := addH(C, 1)
	detailsField := addH(C, 2)
	s := stakescan.Storage{}

	type tc struct {
		checkDetails common.Hash
		signer       common.Address
		status       uint8
		typeChecks   uint8
		amount       *big.Int
		dueBlock     *big.Int
		token        common.Address
		creator      common.Address
	}
	cases := []tc{
		{ // DEL, not yet redeemed
			checkDetails: common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
			signer:       addr("0x460a000000000000000000000000000000000768"),
			status:       0, // None
			typeChecks:   0, // DEL
			amount:       new(big.Int).Mul(u(1000), new(big.Int).Exp(u(10), u(18), nil)),
			dueBlock:     u(32_388_569),
			token:        common.Address{}, // DEL has no token
			creator:      addr("0x50aA000000000000000000000000000000001F19"),
		},
		{ // Token, redeemed, large dueBlock
			checkDetails: common.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222"),
			signer:       addr("0xb8C0000000000000000000000000000000005700"),
			status:       1, // Redeemed
			typeChecks:   1, // Token
			amount:       u(52),
			dueBlock:     new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 252), big.NewInt(1)), // 2^252-1
			token:        addr("0x4E81000000000000000000000000000000005266"),
			creator:      addr("0xEC6d000000000000000000000000000000005b5c"),
		},
	}

	want := map[common.Hash]tc{}
	for _, c := range cases {
		checkHash := stakescan.CheckHashOf(c.checkDetails, c.signer)
		base := mapSlot(checkHash, checksField)
		s[base] = c.checkDetails
		packed := new(big.Int).SetBytes(c.signer.Bytes())
		packed.Or(packed, new(big.Int).Lsh(big.NewInt(int64(c.status)), 160)) // status at byte offset 20
		s[addH(base, 1)] = bigH(packed)

		db := mapSlot(c.checkDetails, detailsField)
		s[db] = bigH(u(int64(c.typeChecks)))
		setIfNonZero(s, addH(db, 1), c.amount)
		setIfNonZero(s, addH(db, 2), c.dueBlock)
		if c.token != (common.Address{}) {
			s[addH(db, 3)] = addrH(c.token)
		}
		s[addH(db, 4)] = addrH(c.creator)
		want[checkHash] = c
	}

	// Decoy address-shaped slots that are not check bases.
	s[common.HexToHash("0xabc1")] = addrH(addr("0x1234123412341234123412341234123412341234"))

	got := stakescan.ReconstructChecks(s, C)
	if len(got) != len(cases) {
		t.Fatalf("got %d checks, want %d", len(got), len(cases))
	}
	for _, g := range got {
		w, ok := want[g.CheckHash]
		if !ok {
			t.Fatalf("unexpected checkHash %s", g.CheckHash.Hex())
		}
		if g.Signer != w.signer || g.Status != w.status {
			t.Errorf("signer/status mismatch: got %s/%d want %s/%d", g.Signer.Hex(), g.Status, w.signer.Hex(), w.status)
		}
		if g.CheckDetailsHash != w.checkDetails {
			t.Errorf("checkDetails mismatch: got %s want %s", g.CheckDetailsHash.Hex(), w.checkDetails.Hex())
		}
		if g.TypeChecks != w.typeChecks {
			t.Errorf("typeChecks: got %d want %d", g.TypeChecks, w.typeChecks)
		}
		if g.Amount.Cmp(w.amount) != 0 {
			t.Errorf("amount: got %s want %s", g.Amount, w.amount)
		}
		if g.DueBlock.Cmp(w.dueBlock) != 0 {
			t.Errorf("dueBlock: got %s want %s", g.DueBlock, w.dueBlock)
		}
		if g.Token != w.token || g.Creator != w.creator {
			t.Errorf("token/creator mismatch: got %s/%s want %s/%s", g.Token.Hex(), g.Creator.Hex(), w.token.Hex(), w.creator.Hex())
		}
		if !g.DetailsFound {
			t.Errorf("detailsFound should be true for %s", g.CheckHash.Hex())
		}
	}
}
