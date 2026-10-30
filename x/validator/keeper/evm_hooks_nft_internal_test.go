package keeper

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-smart-node/contracts/delegation"
)

func TestRewriteFrozenStakeForNFT(t *testing.T) {
	var (
		nft     = common.HexToAddress("0x83e6d6a33c34c9455ab8ecebb0c95bb73fa0193d")
		reserve = common.HexToAddress("0x00000000000000000000000000000000000000de") // reserve coin DRC20 (e.g. WDEL)
		srcVal  = common.HexToAddress("0x05babc196b96c22f763f7e30a84e74665a920f98")
		id      = [32]byte{0xaa}
	)
	newCoin := func(tok, val common.Address) delegation.IDecimalDelegationCommonStake {
		return delegation.IDecimalDelegationCommonStake{Token: tok, Validator: val}
	}

	t.Run("DRC721 withdraw: token+amount rewritten to reserve coin", func(t *testing.T) {
		stake := delegation.IDecimalDelegationCommonStake{TokenType: 2, Token: nft, Amount: big.NewInt(1)}
		coin := newCoin(reserve, srcVal)

		src, applied, err := rewriteFrozenStakeForNFT(&stake, id, exitSibling{coin: &coin, changed: big.NewInt(-500)})

		require.NoError(t, err)
		require.True(t, applied)
		require.Equal(t, reserve, stake.Token, "token rewritten to reserve coin")
		require.Equal(t, int64(500), stake.Amount.Int64(), "amount = |reserve delta|")
		require.Equal(t, srcVal, src, "source validator from coin StakeUpdated")
	})

	t.Run("DRC1155 withdraw: rate*count reserve delta", func(t *testing.T) {
		stake := delegation.IDecimalDelegationCommonStake{TokenType: 3, Token: nft, Amount: big.NewInt(7)}
		coin := newCoin(reserve, srcVal)

		_, applied, err := rewriteFrozenStakeForNFT(&stake, id, exitSibling{coin: &coin, changed: big.NewInt(-900)})

		require.NoError(t, err)
		require.True(t, applied)
		require.Equal(t, reserve, stake.Token)
		require.Equal(t, int64(900), stake.Amount.Int64())
	})

	t.Run("non-NFT (DRC20): no-op", func(t *testing.T) {
		stake := delegation.IDecimalDelegationCommonStake{TokenType: 1, Token: reserve, Amount: big.NewInt(42)}

		src, applied, err := rewriteFrozenStakeForNFT(&stake, id, exitSibling{})

		require.NoError(t, err)
		require.False(t, applied)
		require.Equal(t, reserve, stake.Token, "unchanged")
		require.Equal(t, int64(42), stake.Amount.Int64(), "unchanged")
		require.Equal(t, common.Address{}, src)
	})

	t.Run("NFT but sibling coin events missing: descriptive error, no mutation", func(t *testing.T) {
		stake := delegation.IDecimalDelegationCommonStake{TokenType: 2, Token: nft, Amount: big.NewInt(1)}

		_, applied, err := rewriteFrozenStakeForNFT(&stake, id, exitSibling{})

		require.Error(t, err)
		require.True(t, applied, "recognized as NFT even though data is missing")
		require.Equal(t, nft, stake.Token, "not mutated on error")
	})
}

// Two NFT exits of one coin stake in one tx (the stakeId does not depend on the NFT) must
// each get their own reserve amount, not the last one's.
func TestPairExitSiblingsByPosition(t *testing.T) {
	center, err := delegation.DelegationMetaData.GetAbi()
	require.NoError(t, err)
	var (
		addr    = common.HexToAddress("0xa16c34ed1c0601c0e749e17ebef19752a15faa01")
		other   = common.HexToAddress("0x00000000000000000000000000000000000000ff")
		nft     = common.HexToAddress("0x83e6d6a33c34c9455ab8ecebb0c95bb73fa0193d")
		reserve = common.HexToAddress("0x00000000000000000000000000000000000000de")
		val     = common.HexToAddress("0x05babc196b96c22f763f7e30a84e74665a920f98")
		id      = [32]byte{0xaa}
	)
	mk := func(from common.Address, name string, args ...interface{}) *ethtypes.Log {
		ev := center.Events[name]
		data, err := ev.Inputs.NonIndexed().Pack(args...)
		require.NoError(t, err)
		return &ethtypes.Log{Address: from, Topics: []common.Hash{ev.ID}, Data: data}
	}
	coin := delegation.IDecimalDelegationCommonStake{Validator: val, Token: reserve, Amount: big.NewInt(0), TokenId: big.NewInt(0), HoldTimestamp: big.NewInt(0), HoldStartTime: big.NewInt(0)}
	frozen := delegation.IDecimalDelegationCommonFrozenStake{
		Stake:             delegation.IDecimalDelegationCommonStake{Validator: val, Token: nft, Amount: big.NewInt(1), TokenId: big.NewInt(1), TokenType: 2, HoldTimestamp: big.NewInt(0), HoldStartTime: big.NewInt(0)},
		FreezeType:        1,
		UnfreezeTimestamp: big.NewInt(0),
	}
	logs := []*ethtypes.Log{
		mk(addr, "StakeAmountUpdated", id, big.NewInt(-10000)),
		mk(addr, "StakeUpdated", id, coin, false),
		mk(addr, "WithdrawRequest", id, big.NewInt(0), frozen),
		mk(other, "StakeAmountUpdated", id, big.NewInt(-7)), // foreign contract: ignored
		mk(addr, "StakeAmountUpdated", id, big.NewInt(-1)),
		mk(addr, "StakeUpdated", id, coin, false),
		mk(addr, "WithdrawRequest", id, big.NewInt(1), frozen),
		mk(addr, "WithdrawRequest", id, big.NewInt(2), frozen), // no own siblings
	}
	pairs := pairExitSiblings(logs, strings.ToLower(addr.String()), center)

	for i, want := range map[int]int64{2: 10000, 6: 1} {
		s := frozen.Stake
		_, applied, err := rewriteFrozenStakeForNFT(&s, id, pairs[i])
		require.NoError(t, err)
		require.True(t, applied)
		require.Equal(t, want, s.Amount.Int64(), "exit at log %d", i)
		require.Equal(t, reserve, s.Token)
	}
	s := frozen.Stake
	_, _, err = rewriteFrozenStakeForNFT(&s, id, pairs[7])
	require.Error(t, err, "a request whose siblings were consumed must not reuse them")
}
