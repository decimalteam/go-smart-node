package keeper

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
		coinBy := map[[32]byte]delegation.IDecimalDelegationCommonStake{id: newCoin(reserve, srcVal)}
		amtBy := map[[32]byte]*big.Int{id: big.NewInt(-500)}

		src, applied, err := rewriteFrozenStakeForNFT(&stake, id, coinBy, amtBy)

		require.NoError(t, err)
		require.True(t, applied)
		require.Equal(t, reserve, stake.Token, "token rewritten to reserve coin")
		require.Equal(t, int64(500), stake.Amount.Int64(), "amount = |reserve delta|")
		require.Equal(t, srcVal, src, "source validator from coin StakeUpdated")
	})

	t.Run("DRC1155 withdraw: rate*count reserve delta", func(t *testing.T) {
		stake := delegation.IDecimalDelegationCommonStake{TokenType: 3, Token: nft, Amount: big.NewInt(7)}
		coinBy := map[[32]byte]delegation.IDecimalDelegationCommonStake{id: newCoin(reserve, srcVal)}
		amtBy := map[[32]byte]*big.Int{id: big.NewInt(-900)}

		_, applied, err := rewriteFrozenStakeForNFT(&stake, id, coinBy, amtBy)

		require.NoError(t, err)
		require.True(t, applied)
		require.Equal(t, reserve, stake.Token)
		require.Equal(t, int64(900), stake.Amount.Int64())
	})

	t.Run("non-NFT (DRC20): no-op", func(t *testing.T) {
		stake := delegation.IDecimalDelegationCommonStake{TokenType: 1, Token: reserve, Amount: big.NewInt(42)}

		src, applied, err := rewriteFrozenStakeForNFT(&stake, id,
			map[[32]byte]delegation.IDecimalDelegationCommonStake{}, map[[32]byte]*big.Int{})

		require.NoError(t, err)
		require.False(t, applied)
		require.Equal(t, reserve, stake.Token, "unchanged")
		require.Equal(t, int64(42), stake.Amount.Int64(), "unchanged")
		require.Equal(t, common.Address{}, src)
	})

	t.Run("NFT but sibling coin events missing: descriptive error, no mutation", func(t *testing.T) {
		stake := delegation.IDecimalDelegationCommonStake{TokenType: 2, Token: nft, Amount: big.NewInt(1)}

		_, applied, err := rewriteFrozenStakeForNFT(&stake, id,
			map[[32]byte]delegation.IDecimalDelegationCommonStake{}, map[[32]byte]*big.Int{})

		require.Error(t, err)
		require.True(t, applied, "recognized as NFT even though data is missing")
		require.Equal(t, nft, stake.Token, "not mutated on error")
	})
}
