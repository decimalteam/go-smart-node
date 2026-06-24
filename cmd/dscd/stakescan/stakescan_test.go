package stakescan_test

import (
	"math/big"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// --- independent forward-construction helpers (deliberately NOT the package's
// internal slot math, so a transcription bug in the package fails the round-trip) ---

func addH(h common.Hash, n uint64) common.Hash {
	x := new(big.Int).SetBytes(h[:])
	x.Add(x, new(big.Int).SetUint64(n))
	return common.BigToHash(x)
}

func mapSlot(key, field common.Hash) common.Hash {
	in := append(append([]byte{}, key.Bytes()...), field.Bytes()...)
	return common.BytesToHash(crypto.Keccak256(in))
}

func strMapSlot(key string, field common.Hash) common.Hash {
	in := append([]byte(key), field.Bytes()...)
	return common.BytesToHash(crypto.Keccak256(in))
}

func arrStart(field common.Hash) common.Hash {
	return common.BytesToHash(crypto.Keccak256(field.Bytes()))
}

func bigH(x *big.Int) common.Hash        { return common.BigToHash(x) }
func addrH(a common.Address) common.Hash { return common.BytesToHash(a.Bytes()) }
func u(n int64) *big.Int                 { return big.NewInt(n) }
func addr(hex string) common.Address     { return common.HexToAddress(hex) }

// setStake writes an 8-slot Stake at base, omitting zero slots (mimicking ethermint trimming).
func setStake(s stakescan.Storage, base common.Hash, st stakescan.Stake) {
	s[base] = addrH(st.Validator)
	s[addH(base, 1)] = addrH(st.Delegator)
	s[addH(base, 2)] = addrH(st.Token)
	setIfNonZero(s, addH(base, 3), st.Amount)
	setIfNonZero(s, addH(base, 4), st.TokenID)
	s[addH(base, 5)] = bigH(u(int64(st.TokenType)))
	setIfNonZero(s, addH(base, 6), st.HoldTimestamp)
	setIfNonZero(s, addH(base, 7), st.HoldStartTime)
}

func setIfNonZero(s stakescan.Storage, slot common.Hash, x *big.Int) {
	if x != nil && x.Sign() != 0 {
		s[slot] = bigH(x)
	}
}

// --- tests ---

// TestKeccakAnchor pins the hash function to Keccak-256 (not SHA3-256).
func TestKeccakAnchor(t *testing.T) {
	got := common.BytesToHash(crypto.Keccak256(nil)).Hex()
	want := "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
	if got != want {
		t.Fatalf("keccak256(\"\") = %s, want %s", got, want)
	}
}

func TestReconstructCoinStakes(t *testing.T) {
	B := stakescan.DelegationBase
	stakesField := addH(B, 1)
	penaltyField := addH(B, 2)
	nftBackedField := addH(B, 7)
	s := stakescan.Storage{}

	type tc struct {
		st           stakescan.Stake
		penaltyIndex *big.Int
		nftBacked    *big.Int
	}
	cases := []tc{
		{ // DRC20, non-hold: tokenId=0, holdTimestamp=0, holdStartTime=0 (all trimmed)
			st: stakescan.Stake{
				Validator: addr("0x1111111111111111111111111111111111111111"),
				Delegator: addr("0x2222222222222222222222222222222222222222"),
				Token:     addr("0x3333333333333333333333333333333333333333"),
				Amount:    new(big.Int).Mul(u(123), new(big.Int).Exp(u(10), u(18), nil)),
				TokenType: 1,
			},
			penaltyIndex: u(0),
		},
		{ // DEL with a penalty index and nft-backed amount
			st: stakescan.Stake{
				Validator: addr("0x4444444444444444444444444444444444444444"),
				Delegator: addr("0x5555555555555555555555555555555555555555"),
				Token:     addr("0x6666666666666666666666666666666666666666"),
				Amount:    u(1000),
				TokenType: 4,
			},
			penaltyIndex: u(987654321),
			nftBacked:    u(42),
		},
		{ // DRC721 hold stake: tokenId!=0, holdTimestamp!=0, holdStartTime!=0
			st: stakescan.Stake{
				Validator:     addr("0x7777777777777777777777777777777777777777"),
				Delegator:     addr("0x8888888888888888888888888888888888888888"),
				Token:         addr("0x9999999999999999999999999999999999999999"),
				Amount:        u(1),
				TokenID:       u(777),
				TokenType:     2,
				HoldTimestamp: u(1_700_000_000),
				HoldStartTime: u(1_699_000_000),
			},
			penaltyIndex: u(5),
		},
	}

	want := map[common.Hash]tc{}
	for _, c := range cases {
		id := stakescan.CoinStakeID(c.st.Validator, c.st.Delegator, c.st.Token, orZero(c.st.TokenID), orZero(c.st.HoldTimestamp))
		base := mapSlot(id, stakesField)
		setStake(s, base, c.st)
		setIfNonZero(s, mapSlot(id, penaltyField), c.penaltyIndex)
		setIfNonZero(s, mapSlot(id, nftBackedField), c.nftBacked)
		want[id] = c
	}

	// Decoy slots that look like addresses but are not stake bases.
	s[common.HexToHash("0xdead")] = addrH(addr("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"))
	s[common.HexToHash("0xbeef")] = addrH(addr("0x1234123412341234123412341234123412341234"))

	got := stakescan.ReconstructCoinStakes(s, B)
	if len(got) != len(cases) {
		t.Fatalf("got %d stakes, want %d", len(got), len(cases))
	}
	for _, g := range got {
		w, ok := want[g.StakeID]
		if !ok {
			t.Fatalf("unexpected stake id %s", g.StakeID.Hex())
		}
		if g.Validator != w.st.Validator || g.Delegator != w.st.Delegator || g.Token != w.st.Token {
			t.Errorf("addr mismatch for %s", g.StakeID.Hex())
		}
		if g.Amount.Cmp(orZero(w.st.Amount)) != 0 {
			t.Errorf("amount: got %s want %s", g.Amount, orZero(w.st.Amount))
		}
		if g.TokenID.Cmp(orZero(w.st.TokenID)) != 0 {
			t.Errorf("tokenId: got %s want %s", g.TokenID, orZero(w.st.TokenID))
		}
		if g.TokenType != w.st.TokenType {
			t.Errorf("tokenType: got %d want %d", g.TokenType, w.st.TokenType)
		}
		if g.HoldTimestamp.Cmp(orZero(w.st.HoldTimestamp)) != 0 {
			t.Errorf("holdTimestamp: got %s want %s", g.HoldTimestamp, orZero(w.st.HoldTimestamp))
		}
		if g.HoldStartTime.Cmp(orZero(w.st.HoldStartTime)) != 0 {
			t.Errorf("holdStartTime: got %s want %s", g.HoldStartTime, orZero(w.st.HoldStartTime))
		}
		if g.PenaltyIndex.Cmp(orZero(w.penaltyIndex)) != 0 {
			t.Errorf("penaltyIndex: got %s want %s", g.PenaltyIndex, orZero(w.penaltyIndex))
		}
		if g.NFTBacked.Cmp(orZero(w.nftBacked)) != 0 {
			t.Errorf("nftBacked: got %s want %s", g.NFTBacked, orZero(w.nftBacked))
		}
	}
}

func TestReconstructNFTStakes(t *testing.T) {
	B := stakescan.DelegationBase
	nftField := addH(B, 6)
	nftPenaltyField := addH(B, 8)
	s := stakescan.Storage{}

	nftContract := addr("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	delegator := addr("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	validator := addr("0xcccccccccccccccccccccccccccccccccccccccc")
	reserveToken := addr("0xdddddddddddddddddddddddddddddddddddddddd")
	tokenID := u(555)
	holdTimestamp := u(0)

	id := stakescan.NFTStakeID(nftContract, tokenID, delegator, validator, holdTimestamp)
	base := mapSlot(id, nftField)

	s[base] = addrH(nftContract)
	s[addH(base, 1)] = bigH(tokenID)
	s[addH(base, 2)] = bigH(u(7)) // amount (DRC1155)
	// slot 3 packs nftType(off0,1B) ‖ delegator(off1,20B): (delegator << 8) | nftType
	packed3 := new(big.Int).Lsh(new(big.Int).SetBytes(delegator.Bytes()), 8)
	packed3.Or(packed3, u(3)) // nftType = DRC1155 = 3
	s[addH(base, 3)] = bigH(packed3)
	s[addH(base, 4)] = addrH(validator)
	s[addH(base, 5)] = addrH(reserveToken)
	s[addH(base, 6)] = bigH(u(2500)) // reserveAmount
	// slot 7 holdTimestamp = 0 (trimmed)
	s[addH(base, 8)] = bigH(u(1)) // isActive = true
	s[mapSlot(id, nftPenaltyField)] = bigH(u(99))

	got := stakescan.ReconstructNFTStakes(s, B)
	if len(got) != 1 {
		t.Fatalf("got %d nft stakes, want 1", len(got))
	}
	n := got[0]
	if n.NFTStakeID != id {
		t.Errorf("id: got %s want %s", n.NFTStakeID.Hex(), id.Hex())
	}
	if n.NFTContract != nftContract || n.Delegator != delegator || n.Validator != validator || n.ReserveToken != reserveToken {
		t.Errorf("address fields mismatch: %+v", n)
	}
	if n.NFTType != 3 {
		t.Errorf("nftType: got %d want 3", n.NFTType)
	}
	if n.TokenID.Cmp(tokenID) != 0 || n.Amount.Cmp(u(7)) != 0 || n.ReserveAmount.Cmp(u(2500)) != 0 {
		t.Errorf("numeric fields mismatch: %+v", n)
	}
	if !n.IsActive {
		t.Errorf("isActive: got false want true")
	}
	if n.PenaltyIndex.Cmp(u(99)) != 0 {
		t.Errorf("penaltyIndex: got %s want 99", n.PenaltyIndex)
	}
}

func TestReadFrozenStakes(t *testing.T) {
	B := stakescan.DelegationBase
	frozenField := addH(B, 9)
	s := stakescan.Storage{}

	elems := []stakescan.FrozenStake{
		{
			Stake: stakescan.Stake{
				Validator: addr("0x1010101010101010101010101010101010101010"),
				Delegator: addr("0x2020202020202020202020202020202020202020"),
				Token:     addr("0x3030303030303030303030303030303030303030"),
				Amount:    u(500), TokenType: 1,
			},
			FreezeStatus: 1, FreezeType: 1, UnfreezeTimestamp: u(1_800_000_000),
		},
		{
			Stake: stakescan.Stake{
				Validator: addr("0x4040404040404040404040404040404040404040"),
				Delegator: addr("0x5050505050505050505050505050505050505050"),
				Token:     addr("0x6060606060606060606060606060606060606060"),
				Amount:    u(750), TokenType: 4,
			},
			FreezeStatus: 1, FreezeType: 2, UnfreezeTimestamp: u(1_900_000_000),
		},
	}

	s[frozenField] = bigH(u(int64(len(elems))))
	start := arrStart(frozenField)
	for i, e := range elems {
		base := addH(start, uint64(i)*10)
		setStake(s, base, e.Stake)
		packed8 := new(big.Int).Lsh(u(int64(e.FreezeType)), 8)
		packed8.Or(packed8, u(int64(e.FreezeStatus)))
		s[addH(base, 8)] = bigH(packed8)
		s[addH(base, 9)] = bigH(e.UnfreezeTimestamp)
	}

	got := stakescan.ReadFrozenStakes(s, B, 9)
	if len(got) != len(elems) {
		t.Fatalf("got %d frozen, want %d", len(got), len(elems))
	}
	for i, g := range got {
		w := elems[i]
		if g.Index != uint64(i) {
			t.Errorf("index: got %d want %d", g.Index, i)
		}
		if g.Stake.Validator != w.Stake.Validator || g.Stake.Token != w.Stake.Token {
			t.Errorf("elem %d stake addr mismatch", i)
		}
		if g.Stake.Amount.Cmp(w.Stake.Amount) != 0 {
			t.Errorf("elem %d amount: got %s want %s", i, g.Stake.Amount, w.Stake.Amount)
		}
		if g.FreezeStatus != w.FreezeStatus || g.FreezeType != w.FreezeType {
			t.Errorf("elem %d freeze status/type: got %d/%d want %d/%d", i, g.FreezeStatus, g.FreezeType, w.FreezeStatus, w.FreezeType)
		}
		if g.UnfreezeTimestamp.Cmp(w.UnfreezeTimestamp) != 0 {
			t.Errorf("elem %d unfreeze: got %s want %s", i, g.UnfreezeTimestamp, w.UnfreezeTimestamp)
		}
	}
}

// setLegacyStake writes a 7-slot LegacyStake (no holdStartTime) at base.
func setLegacyStake(s stakescan.Storage, base common.Hash, st stakescan.Stake) {
	s[base] = addrH(st.Validator)
	s[addH(base, 1)] = addrH(st.Delegator)
	s[addH(base, 2)] = addrH(st.Token)
	setIfNonZero(s, addH(base, 3), st.Amount)
	setIfNonZero(s, addH(base, 4), st.TokenID)
	s[addH(base, 5)] = bigH(u(int64(st.TokenType)))
	setIfNonZero(s, addH(base, 6), st.HoldTimestamp)
}

func TestReadLegacyFrozenStakes(t *testing.T) {
	B := stakescan.DelegationBase
	depField := addH(B, 3) // _frozenStakesDeprecated
	s := stakescan.Storage{}

	elems := []stakescan.FrozenStake{
		{
			Stake: stakescan.Stake{
				Validator: addr("0x7f7ef7539869e92a8825b062Df7C0090524B293F"),
				Delegator: addr("0x65E6112341231234123412341234123412341F5f"),
				Token:     addr("0x1c5Db575E2Ac0894077C0B6Df7c0090524B7C0B"),
				Amount:    new(big.Int).Mul(u(5104), new(big.Int).Exp(u(10), u(18), nil)),
				TokenType: 4, // DEL
			},
			FreezeStatus: 1, FreezeType: 1, UnfreezeTimestamp: u(1_723_669_392),
		},
		{ // second element must align via stride 9 (regression guard for the old bug)
			Stake: stakescan.Stake{
				Validator: addr("0xE0c691D6C1432283e863Ae7f4fd232e8673a7A38"),
				Delegator: addr("0x74441234123412341234123412341234123420aA"),
				Token:     addr("0xc5E3000000000000000000000000000000000915"),
				Amount:    new(big.Int).Mul(u(7000), new(big.Int).Exp(u(10), u(18), nil)),
				TokenType: 1, // DRC20
			},
			FreezeStatus: 1, FreezeType: 2, UnfreezeTimestamp: u(1_730_000_000),
		},
	}

	s[depField] = bigH(u(int64(len(elems))))
	start := arrStart(depField)
	for i, e := range elems {
		base := addH(start, uint64(i)*9) // legacy stride = 9
		setLegacyStake(s, base, e.Stake)
		packed7 := new(big.Int).Lsh(u(int64(e.FreezeType)), 8)
		packed7.Or(packed7, u(int64(e.FreezeStatus)))
		s[addH(base, 7)] = bigH(packed7)
		s[addH(base, 8)] = bigH(e.UnfreezeTimestamp)
	}

	got := stakescan.ReadLegacyFrozenStakes(s, B, 3)
	if len(got) != len(elems) {
		t.Fatalf("got %d legacy frozen, want %d", len(got), len(elems))
	}
	for i, g := range got {
		w := elems[i]
		if g.Stake.Validator != w.Stake.Validator || g.Stake.Delegator != w.Stake.Delegator || g.Stake.Token != w.Stake.Token {
			t.Errorf("elem %d addr mismatch: got v=%s d=%s t=%s", i, g.Stake.Validator.Hex(), g.Stake.Delegator.Hex(), g.Stake.Token.Hex())
		}
		if g.Stake.Amount.Cmp(w.Stake.Amount) != 0 {
			t.Errorf("elem %d amount: got %s want %s", i, g.Stake.Amount, w.Stake.Amount)
		}
		if g.Stake.TokenType != w.Stake.TokenType {
			t.Errorf("elem %d tokenType: got %d want %d", i, g.Stake.TokenType, w.Stake.TokenType)
		}
		if g.Stake.HoldStartTime.Sign() != 0 {
			t.Errorf("elem %d holdStartTime should be 0 for legacy, got %s", i, g.Stake.HoldStartTime)
		}
		if g.FreezeStatus != w.FreezeStatus || g.FreezeType != w.FreezeType {
			t.Errorf("elem %d freeze status/type: got %d/%d want %d/%d", i, g.FreezeStatus, g.FreezeType, w.FreezeStatus, w.FreezeType)
		}
		if g.UnfreezeTimestamp.Cmp(w.UnfreezeTimestamp) != 0 {
			t.Errorf("elem %d unfreeze: got %s want %s", i, g.UnfreezeTimestamp, w.UnfreezeTimestamp)
		}
	}
}

func TestAutoUnbondQueue(t *testing.T) {
	ns := stakescan.AutoUnbondQueueSlotMain // any field slot works; use the real one
	s := stakescan.Storage{}

	entries := []stakescan.AutoUnbondEntry{
		{
			Validator: addr("0x8fF2f220FB80b3F26cd96652aeCbf3129983c115"),
			Delegator: addr("0x1aC9000000000000000000000000000000006657"),
			Amount:    new(big.Int).Mul(u(200_000_000), new(big.Int).Exp(u(10), u(18), nil)),
			Token:     addr("0x4E81000000000000000000000000000000005266"),
		},
		{ // hold entry (holdTimestamp != 0)
			Validator:     addr("0x8fF2f220FB80b3F26cd96652aeCbf3129983c115"),
			Delegator:     addr("0x7584000000000000000000000000000000003c07"),
			Amount:        new(big.Int).Mul(u(4499), new(big.Int).Exp(u(10), u(18), nil)),
			Token:         addr("0xB7D6000000000000000000000000000000006226"),
			HoldTimestamp: u(1_841_711_307),
		},
	}

	s[ns] = bigH(u(int64(len(entries))))
	start := arrStart(ns)
	for i, e := range entries {
		base := addH(start, uint64(i)*5)
		s[base] = addrH(e.Validator)
		s[addH(base, 1)] = addrH(e.Delegator)
		setIfNonZero(s, addH(base, 2), e.Amount)
		s[addH(base, 3)] = addrH(e.Token)
		setIfNonZero(s, addH(base, 4), e.HoldTimestamp)
	}

	// Decoy: a legacy frozen array (stride 9) — must NOT be detected as a stride-5 queue.
	depField := addH(stakescan.DelegationBase, 3)
	s[depField] = bigH(u(3))
	fstart := arrStart(depField)
	for i := 0; i < 3; i++ {
		b := addH(fstart, uint64(i)*9)
		setLegacyStake(s, b, stakescan.Stake{
			Validator: addr("0x1111111111111111111111111111111111111111"),
			Delegator: addr("0x2222222222222222222222222222222222222222"),
			Token:     addr("0x3333333333333333333333333333333333333333"),
			Amount:    u(100), TokenType: 1,
		})
		s[addH(b, 7)] = bigH(u(0x0101))
		s[addH(b, 8)] = bigH(u(123))
	}

	// Direct read.
	got := stakescan.ReadAutoUnbondQueueAt(s, ns)
	if len(got) != len(entries) {
		t.Fatalf("ReadAutoUnbondQueueAt: got %d, want %d", len(got), len(entries))
	}
	for i, g := range got {
		w := entries[i]
		if g.Validator != w.Validator || g.Delegator != w.Delegator || g.Token != w.Token {
			t.Errorf("entry %d addr mismatch: %+v", i, g)
		}
		if g.Amount.Cmp(w.Amount) != 0 {
			t.Errorf("entry %d amount: got %s want %s", i, g.Amount, w.Amount)
		}
		if g.HoldTimestamp.Cmp(orZero(w.HoldTimestamp)) != 0 {
			t.Errorf("entry %d holdTimestamp: got %s want %s", i, g.HoldTimestamp, orZero(w.HoldTimestamp))
		}
	}

	// Locator must find the queue and must not mistake the frozen decoy for one.
	qs := stakescan.FindAutoUnbondQueues(s, 0.9)
	foundQueue := false
	for _, q := range qs {
		if q.FieldSlot == ns {
			foundQueue = true
			if q.Length != uint64(len(entries)) {
				t.Errorf("located queue length: got %d want %d", q.Length, len(entries))
			}
		}
		if q.FieldSlot == depField {
			t.Errorf("frozen decoy (stride 9) was wrongly detected as a stride-5 auto-unbond queue")
		}
	}
	if !foundQueue {
		t.Fatalf("FindAutoUnbondQueues did not locate the queue at %s", ns.Hex())
	}
}

func TestResolveAddressBySymbol(t *testing.T) {
	cc := stakescan.Storage{}
	want := addr("0xfeedfeedfeedfeedfeedfeedfeedfeedfeedfeed")
	// _contractAddresses (mapping(string=>address)) at offset 0 of ContractCenterBase.
	cc[strMapSlot("delegation", stakescan.ContractCenterBase)] = addrH(want)

	got := stakescan.ResolveAddressBySymbol(cc, "delegation")
	if got != want {
		t.Fatalf("resolved %s, want %s", got.Hex(), want.Hex())
	}
	if missing := stakescan.ResolveAddressBySymbol(cc, "nonexistent"); missing != (common.Address{}) {
		t.Fatalf("missing symbol should resolve to zero address, got %s", missing.Hex())
	}
}

func orZero(x *big.Int) *big.Int {
	if x == nil {
		return big.NewInt(0)
	}
	return x
}
