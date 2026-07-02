package redenom

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
)

// writeCoinStake lays out one _stakes[stakeId] record in the delegation storage
// exactly as the Solidity contract does (validator@0, delegator@1, token@2,
// amount@3, tokenId@4, tokenType@5, holdTimestamp@6, holdStartTime@7), keyed so
// that ReconstructCoinStakes ID-verifies it.
func writeCoinStake(s stakescan.Storage, validator, delegator, token common.Address, amount *big.Int, tokenType int64) {
	id := stakescan.CoinStakeID(validator, delegator, token, big.NewInt(0), big.NewInt(0))
	base := stakescan.MappingSlotBytes32(id, stakescan.SlotAdd(stakescan.DelegationBase, 1))

	s[base] = common.BytesToHash(validator.Bytes())
	s[stakescan.SlotAdd(base, 1)] = common.BytesToHash(delegator.Bytes())
	s[stakescan.SlotAdd(base, 2)] = common.BytesToHash(token.Bytes())
	s[stakescan.SlotAdd(base, 3)] = common.BigToHash(amount)
	// tokenId@4 = 0 (absent), holdTimestamp@6 / holdStartTime@7 = 0 (absent)
	s[stakescan.SlotAdd(base, 5)] = common.BigToHash(big.NewInt(tokenType))
}

// writeFrozenDELStake appends one live _frozenStakes element (stride 10) with an
// embedded DEL stake for the given validator.
func writeFrozenDELStake(s stakescan.Storage, validator, delegator, token common.Address, amount *big.Int) {
	field := stakescan.FrozenArrayFieldSlot(stakescan.DelegationBase)
	n := new(big.Int).SetBytes(s[field].Bytes()).Uint64()
	elem := stakescan.SlotAdd(stakescan.ArrayElemStart(field), n*stakescan.FrozenStride)

	s[elem] = common.BytesToHash(validator.Bytes())
	s[stakescan.SlotAdd(elem, 1)] = common.BytesToHash(delegator.Bytes())
	s[stakescan.SlotAdd(elem, 2)] = common.BytesToHash(token.Bytes())
	s[stakescan.SlotAdd(elem, 3)] = common.BigToHash(amount)
	s[stakescan.SlotAdd(elem, 5)] = common.BigToHash(big.NewInt(4)) // TokenType.DEL
	s[stakescan.SlotAdd(elem, 8)] = common.BigToHash(big.NewInt(0x100)) // freezeType=Withdraw(1), status=Completed(0)
	s[field] = common.BigToHash(new(big.Int).SetUint64(n + 1))
}

// validatorReserveSlots derives the ValidatorReserve slots for (validator, token,
// tokenId=0) independently of the production helper: penaltyIndex at
// keccak(hashedTokenID(token,0) ‖ keccak(validator ‖ base+5)), reserve one word
// after (DecimalDelegationCommon.sol _getValidatorReserve).
func validatorReserveSlots(validator, token common.Address) (penaltySlot, reserveSlot common.Hash) {
	htid := common.BytesToHash(crypto.Keccak256(token.Bytes(), common.BigToHash(big.NewInt(0)).Bytes()))
	mid := stakescan.MappingSlotBytes32(common.BytesToHash(validator.Bytes()), stakescan.SlotAdd(stakescan.DelegationBase, 5))
	rbase := stakescan.MappingSlotBytes32(htid, mid)
	return rbase, stakescan.SlotAdd(rbase, 1)
}

// TestCollectDelegationWrites_ValidatorReserve is the R2 regression
// (redenom-full-audit-2026-07-02.md): the per-validator DEL stake aggregate
// _validatorTokens[validator][hashedTokenID(WDEL,0)].reserve must be divided like
// the stakes that feed it — once per validator (deduped across stakes, derived
// from active AND frozen DEL stakes) — while custom-coin reserves and the
// adjacent penaltyIndex word stay untouched, and zero/absent reserves emit no
// write.
func TestCollectDelegationWrites_ValidatorReserve(t *testing.T) {
	deleg := addrOf(0x11)
	wdel := addrOf(0x22)
	customCoin := addrOf(0x33)

	valDEL := addrOf(0xc1)    // two active DEL stakes -> one deduped reserve write
	valCustom := addrOf(0xc2) // custom-coin stake -> reserve untouched
	valFrozen := addrOf(0xc3) // frozen-only DEL stake -> reserve still derived+scaled
	valZero := addrOf(0xc4)   // DEL stake but reserve slot absent -> no write

	store := stakescan.Storage{}
	writeCoinStake(store, valDEL, addrOf(0xb1), wdel, big.NewInt(5000), tokenTypeDEL)
	writeCoinStake(store, valDEL, addrOf(0xb2), wdel, big.NewInt(3000), tokenTypeDEL)
	writeCoinStake(store, valCustom, addrOf(0xb3), customCoin, big.NewInt(7000), 1) // DRC20
	writeCoinStake(store, valZero, addrOf(0xb4), wdel, big.NewInt(4000), tokenTypeDEL)
	writeFrozenDELStake(store, valFrozen, addrOf(0xb5), wdel, big.NewInt(2000))

	delPenalty, delReserve := validatorReserveSlots(valDEL, wdel)
	_, customReserve := validatorReserveSlots(valCustom, customCoin)
	_, frozenReserve := validatorReserveSlots(valFrozen, wdel)
	_, zeroReserve := validatorReserveSlots(valZero, wdel)

	store[delPenalty] = common.BigToHash(big.NewInt(7))      // penaltyIndex: must survive
	store[delReserve] = common.BigToHash(big.NewInt(8000))   // 5000+3000
	store[customReserve] = common.BigToHash(big.NewInt(7000))
	store[frozenReserve] = common.BigToHash(big.NewInt(2000))

	rep := newReport(sdkmath.NewInt(1000))
	holders := map[common.Address]struct{}{}
	writes, err := collectDelegationWrites(store, deleg, wdel, big.NewInt(1000), &rep, noopLogger{}, holders)
	if err != nil {
		t.Fatalf("collectDelegationWrites: %v", err)
	}

	if rep.EVMValidatorReserveSlots != 2 {
		t.Errorf("EVMValidatorReserveSlots = %d, want 2 (valDEL deduped + valFrozen)", rep.EVMValidatorReserveSlots)
	}

	got := map[common.Hash][]*big.Int{}
	for _, w := range writes {
		got[w.slot] = append(got[w.slot], w.newVal)
	}

	if vs := got[delReserve]; len(vs) != 1 || vs[0].Cmp(big.NewInt(8)) != 0 {
		t.Errorf("valDEL reserve writes = %v, want exactly one write of 8 (floor(8000/1000), deduped)", vs)
	}
	if vs := got[frozenReserve]; len(vs) != 1 || vs[0].Cmp(big.NewInt(2)) != 0 {
		t.Errorf("valFrozen reserve writes = %v, want exactly one write of 2 (derived from frozen stake)", vs)
	}
	if vs := got[customReserve]; len(vs) != 0 {
		t.Errorf("custom-coin validator reserve was scaled (%v), must be untouched", vs)
	}
	if vs := got[delPenalty]; len(vs) != 0 {
		t.Errorf("penaltyIndex slot was written (%v), must be untouched", vs)
	}
	if vs := got[zeroReserve]; len(vs) != 0 {
		t.Errorf("zero/absent reserve got a write (%v), must emit none", vs)
	}
}
