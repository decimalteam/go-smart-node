package ante

import (
	"fmt"
	"math"

	sdkStore "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Gas consumption descriptors.
const (
	GasCommissionDesc = "commission"
)

// This gas meter is fork of standart cosmos sdk basicGasMeter
// except one IMPORTANT LINE, see below
// CommissionGasMeter consume gas only during special cases in FeeDecorator.
// This allow use predicatable amount of gas in transaction, no gas consumption for read/write/etc operations
type commissionGasMeter struct {
	limit    sdk.Gas
	consumed sdk.Gas
}

func NewCommissionGasMeter(limit sdk.Gas) sdk.GasMeter {
	return &commissionGasMeter{
		limit:    limit,
		consumed: 0,
	}
}

func (g *commissionGasMeter) GasConsumed() sdk.Gas {
	return g.consumed
}

func (g *commissionGasMeter) Limit() sdk.Gas {
	return g.limit
}

func (g *commissionGasMeter) GasConsumedToLimit() sdk.Gas {
	if g.IsPastLimit() {
		return g.limit
	}
	return g.consumed
}

// addUint64Overflow performs the addition operation on two uint64 integers and
// returns a boolean on whether or not the result overflows.
func addUint64Overflow(a, b uint64) (uint64, bool) {
	if math.MaxUint64-a < b {
		return 0, true
	}

	return a + b, false
}

func (g *commissionGasMeter) ConsumeGas(amount sdk.Gas, descriptor string) {
	// IMPORTANT LINE
	// skip all non-commission calls
	if descriptor != GasCommissionDesc {
		return
	}
	var overflow bool
	g.consumed, overflow = addUint64Overflow(g.consumed, amount)
	if overflow {
		g.consumed = math.MaxUint64
		panic(sdk.ErrorGasOverflow{Descriptor: descriptor})
	}

	if g.consumed > g.limit {
		panic(sdk.ErrorOutOfGas{Descriptor: descriptor})
	}
}

// RefundGas will deduct the given amount from the gas consumed. If the amount is greater than the
// gas consumed, the function will panic.
//
// Use case: This functionality enables refunding gas to the transaction or block gas pools so that
// EVM-compatible chains can fully support the go-ethereum StateDb interface.
// See https://github.com/cosmos/cosmos-sdk/pull/9403 for reference.
func (g *commissionGasMeter) RefundGas(amount sdk.Gas, descriptor string) {
	if g.consumed < amount {
		panic(sdkStore.ErrorNegativeGasConsumed{Descriptor: descriptor})
	}

	g.consumed -= amount
}

func (g *commissionGasMeter) IsPastLimit() bool {
	return g.consumed > g.limit
}

func (g *commissionGasMeter) IsOutOfGas() bool {
	return g.consumed >= g.limit
}

func (g *commissionGasMeter) String() string {
	return fmt.Sprintf("CommissionGasMeter:\n  limit: %d\n  consumed: %d", g.limit, g.consumed)
}
