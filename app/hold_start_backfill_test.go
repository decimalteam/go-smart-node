package app

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestComputeStakeID_packsLikeSolidity(t *testing.T) {
	val := common.HexToAddress("0x1111111111111111111111111111111111111111")
	del := common.HexToAddress("0x2222222222222222222222222222222222222222")
	tok := common.HexToAddress("0x3333333333333333333333333333333333333333")
	holdTs := big.NewInt(1_800_000_000)

	// Reference: abi.encodePacked(address,address,address,uint256,uint256) = 20+20+20+32+32.
	want := make([]byte, 0, 124)
	want = append(want, val.Bytes()...)
	want = append(want, del.Bytes()...)
	want = append(want, tok.Bytes()...)
	want = append(want, common.LeftPadBytes(big.NewInt(0).Bytes(), 32)...) // tokenId = 0
	want = append(want, common.LeftPadBytes(holdTs.Bytes(), 32)...)
	require.Equal(t, common.BytesToHash(crypto.Keccak256(want)), computeStakeID(val, del, tok, big.NewInt(0), holdTs))
}

func TestHoldStartTimeSlot_formula(t *testing.T) {
	stakeID := common.HexToHash("0xabcdef00000000000000000000000000000000000000000000000000000000ff")
	base := new(big.Int)
	base.SetString("c1dae510251b57b62087f142cc8746564a134308f8810580b66b08a29b6def00", 16)
	mappingSlot := new(big.Int).Add(base, big.NewInt(1)) // _stakes is field 1
	start := new(big.Int).SetBytes(crypto.Keccak256(
		append(stakeID.Bytes(), common.BigToHash(mappingSlot).Bytes()...),
	))
	want := common.BigToHash(new(big.Int).Add(start, big.NewInt(7))) // holdStartTime offset 7
	require.Equal(t, want, holdStartTimeSlot(stakeID))
}

func TestSelectHoldStart_earliestNonZeroOrFloor(t *testing.T) {
	floor := int64(1_700_000_000)
	// Multiple holds same end -> earliest non-zero start.
	require.Equal(t, int64(1_690_000_000), selectHoldStart([]int64{1_695_000_000, 1_690_000_000}, floor))
	// A zero among them is replaced by the floor before comparison.
	require.Equal(t, int64(1_699_000_000), selectHoldStart([]int64{0, 1_699_000_000}, floor))
	// All zero -> floor.
	require.Equal(t, floor, selectHoldStart([]int64{0, 0}, floor))
}

func TestHelpers_matchSolidityA3Vector(t *testing.T) {
	// From the Hardhat A3 test (delegation/test/hold-start-time.ts) on a live deployment.
	val := common.HexToAddress("0x037Ae89141152505FA0101936748DefA38c76103")
	del := common.HexToAddress("0xA8b091a85A5E938B43EE6E0F3A7126792C7F2Ba1")
	tok := common.HexToAddress("0x820c446B27F91AB03C751Fead10529D10C3F3a3d")
	holdTs := big.NewInt(1816349537)

	stakeID := computeStakeID(val, del, tok, big.NewInt(0), holdTs)
	require.Equal(t,
		common.HexToHash("0xd4656b2a15ceccac1da82ab8462752f2d4df626b1d6c9058be90f7a6509dac9d"),
		stakeID, "stakeId must match the Solidity getHoldStakeId output")

	slot := holdStartTimeSlot(stakeID)
	require.Equal(t,
		common.HexToHash("0xfc237ebdb26ab7dc99307e2a094f264ab2da4a3c5aa4d89f56e43a5c9c038018"),
		slot, "holdStartTime slot must match the storage slot the Solidity test read")
}

func TestCoinHoldStakeID_orderingAndConsistency(t *testing.T) {
	// Use two distinct 20-byte sequences so validator != delegator in the EVM address
	// space, making a positional swap detectable.
	rawVal := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
		0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14,
	}
	rawDel := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33, 0x44,
		0x55, 0x66, 0x77, 0x88, 0x99, 0x10, 0x20, 0x30, 0x40, 0x50,
	}

	// Build bech32 strings using the SDK helpers — same bytes, different HRP.
	valBech32 := sdk.ValAddress(rawVal).String() // d0valoper…
	delBech32 := sdk.AccAddress(rawDel).String() // d0…

	valHex := common.BytesToAddress(rawVal)
	delHex := common.BytesToAddress(rawDel)

	tok := common.HexToAddress("0x820c446B27F91AB03C751Fead10529D10C3F3a3d")
	endTime := int64(1816349537)

	// (a) coinHoldStakeID must place the validator in the first computeStakeID arg.
	want := computeStakeID(valHex, delHex, tok, big.NewInt(0), big.NewInt(endTime))
	got, err := coinHoldStakeID(valBech32, delBech32, tok, endTime)
	require.NoError(t, err)
	require.Equal(t, want, got, "coinHoldStakeID must equal computeStakeID(valHex, delHex, ...)")

	// (b) Swapping the bech32 args must yield a different stakeId.
	// ValAddressFromBech32 enforces the d0valoper HRP, so passing a d0 address as the
	// first arg returns an error — that error itself is proof that argument order is
	// enforced. We additionally confirm the non-swapped result differs from the
	// computeStakeID with args reversed.
	wrongOrder := computeStakeID(delHex, valHex, tok, big.NewInt(0), big.NewInt(endTime))
	require.NotEqual(t, want, wrongOrder,
		"computeStakeID(delHex, valHex) must differ from computeStakeID(valHex, delHex) — swap is detectable")

	// Confirm coinHoldStakeID rejects a delegator bech32 (d0…) in the validator slot.
	_, swapErr := coinHoldStakeID(delBech32, valBech32, tok, endTime)
	require.Error(t, swapErr, "passing delegator bech32 as validator arg must return an error (HRP mismatch)")
}
