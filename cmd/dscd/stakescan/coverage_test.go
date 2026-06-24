package stakescan_test

import (
	"math/big"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
	"github.com/ethereum/go-ethereum/common"
)

// TestERC7201 anchors the namespace derivation to the real on-chain contract constants:
// if the formula were wrong, none of the reconstruction would line up with live storage.
func TestERC7201(t *testing.T) {
	cases := []struct {
		name string
		want common.Hash
	}{
		{"DECIMAL_DELEGATION_COMMON_STORAGE_LOCATION", stakescan.DelegationBase},
		{"DECIMAL_CONTRACT_CENTER_STORAGE_LOCATION", stakescan.ContractCenterBase},
		{"DECIMAL_DELEGATION_MIGRATION_STORAGE_LOCATION", common.HexToHash("0x0087ae057eb38dff9d4851f20c0622c62ee2528e4d13ada7657e17acfe136200")},
	}
	for _, c := range cases {
		if got := stakescan.ERC7201(c.name); got != c.want {
			t.Errorf("ERC7201(%q) = %s, want %s", c.name, got.Hex(), c.want.Hex())
		}
	}
}

// TestAuditFullCoverage builds a synthetic storage containing one of each structure and
// asserts the audit accounts for every slot (0 unexplained).
func TestAuditFullCoverage(t *testing.T) {
	B := stakescan.DelegationBase
	s := stakescan.Storage{}

	// contractCenter (B+0)
	s[B] = addrH(addr("0xc108715A06f76CAA96fa2c943Ebf05159c29A87D"))

	// one coin stake in _stakes (B+1)
	cs := stakescan.Stake{
		Validator: addr("0x1111111111111111111111111111111111111111"),
		Delegator: addr("0x2222222222222222222222222222222222222222"),
		Token:     addr("0x3333333333333333333333333333333333333333"),
		Amount:    u(1000), TokenType: 1,
	}
	cs.StakeID = stakescan.CoinStakeID(cs.Validator, cs.Delegator, cs.Token, big.NewInt(0), big.NewInt(0))
	setStake(s, mapSlot(cs.StakeID, addH(B, 1)), cs)
	result := stakescan.Result{CoinStakes: []stakescan.CoinStake{{Stake: cs}}}

	// one deprecated frozen entry (B+3, legacy stride 9) + its migration flag
	fz := stakescan.FrozenStake{
		Stake: stakescan.Stake{
			Validator: addr("0x4444444444444444444444444444444444444444"),
			Delegator: addr("0x5555555555555555555555555555555555555555"),
			Token:     addr("0x6666666666666666666666666666666666666666"),
			Amount:    u(500), TokenType: 4,
		},
		FreezeStatus: 1, FreezeType: 1, UnfreezeTimestamp: u(1_700_000_000),
	}
	fz.Stake.StakeID = stakescan.CoinStakeID(fz.Stake.Validator, fz.Stake.Delegator, fz.Stake.Token, big.NewInt(0), big.NewInt(0))
	s[addH(B, 3)] = bigH(u(1))
	fbase := arrStart(addH(B, 3))
	setLegacyStake(s, fbase, fz.Stake)
	packed7 := new(big.Int).Lsh(u(int64(fz.FreezeType)), 8)
	packed7.Or(packed7, u(int64(fz.FreezeStatus)))
	s[addH(fbase, 7)] = bigH(packed7)
	s[addH(fbase, 8)] = bigH(fz.UnfreezeTimestamp)
	result.FrozenDeprecated = []stakescan.FrozenStake{fz}

	migNS := stakescan.ERC7201("DECIMAL_DELEGATION_MIGRATION_STORAGE_LOCATION")
	s[mapSlot(stakescan.FrozenStakeHash(fz.Stake, fz.UnfreezeTimestamp), migNS)] = bigH(u(1))

	// one auto-unbond entry
	auNS := stakescan.AutoUnbondQueueSlotMain
	s[auNS] = bigH(u(1))
	aubase := arrStart(auNS)
	s[aubase] = addrH(addr("0x7777777777777777777777777777777777777777"))
	s[addH(aubase, 1)] = addrH(addr("0x8888888888888888888888888888888888888888"))
	s[addH(aubase, 2)] = bigH(u(123))
	s[addH(aubase, 3)] = addrH(addr("0x9999999999999999999999999999999999999999"))
	au := []stakescan.AutoUnbondQueue{{FieldSlot: auNS, Length: 1, Entries: stakescan.ReadAutoUnbondQueueAt(s, auNS)}}

	// slot 0 (legacy) + an OZ fixed slot
	s[common.Hash{}] = bigH(u(1))
	s[common.HexToHash("0x9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300")] = addrH(addr("0xAbcDef0000000000000000000000000000000001"))

	rep := stakescan.Audit(s, B, result, au, 0)
	if rep.TotalSlots != len(s) {
		t.Fatalf("total slots %d != storage size %d", rep.TotalSlots, len(s))
	}
	if rep.Unexplained != 0 {
		t.Fatalf("expected 0 unexplained, got %d: %+v", rep.Unexplained, rep.UnexplainedSample)
	}
	for _, want := range []string{"contractCenter", "coinStake", "frozenDeprecated", "frozenStakeMigrated", "autoUnbond", "slot0.legacyInitialized", "oz.Ownable"} {
		if rep.ByCategory[want] == 0 {
			t.Errorf("expected category %q to have slots, got 0", want)
		}
	}
}
