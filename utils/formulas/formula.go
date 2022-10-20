package formulas

import (
	"math/big"

	sdkmath "cosmossdk.io/math"

	"bitbucket.org/decimalteam/go-smart-node/types/bigfloat"
)

// CalculatePurchaseReturn calculates amount of coin that user will receive by depositing given amount of BIP.
// Return = supply * ((1 + deposit / reserve) ^ (crr / 100) - 1)
func CalculatePurchaseReturn(supply sdkmath.Int, reserve sdkmath.Int, crr uint, deposit sdkmath.Int) sdkmath.Int {
	if deposit.Sign() == 0 {
		return sdkmath.NewInt(0)
	}

	if crr == 100 {
		return sdkmath.NewInt(1).Mul(supply).Mul(deposit).Quo(reserve)
	}

	tSupply := newFloat(0).SetInt(supply.BigInt())
	tReserve := newFloat(0).SetInt(reserve.BigInt())
	tDeposit := newFloat(0).SetInt(deposit.BigInt())

	res := newFloat(0).Quo(tDeposit, tReserve)          // deposit / reserve
	res.Add(res, newFloat(1))                           // 1 + (deposit / reserve)
	res = bigfloat.Pow(res, newFloat(float64(crr)/100)) // (1 + deposit / reserve) ^ (crr / 100)
	res.Sub(res, newFloat(1))                           // ((1 + deposit / reserve) ^ (crr / 100) - 1)
	res.Mul(res, tSupply)                               // supply * ((1 + deposit / reserve) ^ (crr / 100) - 1)

	converted, _ := res.Int(nil)
	result := sdkmath.NewIntFromBigInt(converted)

	return result
}

// CalculatePurchaseAmount is the reversed version of function CalculatePurchaseReturn.
// Deposit = reserve * (((wantReceive + supply) / supply)^(100/c) - 1)
func CalculatePurchaseAmount(supply sdkmath.Int, reserve sdkmath.Int, crr uint, wantReceive sdkmath.Int) sdkmath.Int {
	if wantReceive.Sign() == 0 {
		return sdkmath.NewInt(0)
	}

	if crr == 100 {
		return sdkmath.NewInt(1).Mul(wantReceive).Mul(reserve).Quo(supply)
	}

	tSupply := newFloat(0).SetInt(supply.BigInt())
	tReserve := newFloat(0).SetInt(reserve.BigInt())
	tWantReceive := newFloat(0).SetInt(wantReceive.BigInt())

	res := newFloat(0).Add(tWantReceive, tSupply)       // reserve + supply
	res.Quo(res, tSupply)                               // (reserve + supply) / supply
	res = bigfloat.Pow(res, newFloat(100/float64(crr))) // ((reserve + supply) / supply)^(100/c)
	res.Sub(res, newFloat(1))                           // (((reserve + supply) / supply)^(100/c) - 1)
	res.Mul(res, tReserve)                              // reserve * (((reserve + supply) / supply)^(100/c) - 1)

	converted, _ := res.Int(nil)
	result := sdkmath.NewIntFromBigInt(converted)

	return result
}

// CalculateSaleReturn returns amount of BIP user will receive by depositing given amount of coins.
// Return = reserve * (1 - (1 - sellAmount / supply) ^ (100 / crr))
func CalculateSaleReturn(supply sdkmath.Int, reserve sdkmath.Int, crr uint, sellAmount sdkmath.Int) sdkmath.Int {
	// special case for 0 sell amount
	if sellAmount.Sign() == 0 {
		return sdkmath.NewInt(0)
	}

	// special case for selling the entire supply
	if sellAmount.Equal(supply) {
		return reserve
	}

	if crr == 100 {
		return sdkmath.NewInt(1).Mul(reserve).Mul(sellAmount).Quo(supply)
	}

	tSupply := newFloat(0).SetInt(supply.BigInt())
	tReserve := newFloat(0).SetInt(reserve.BigInt())
	tSellAmount := newFloat(0).SetInt(sellAmount.BigInt())

	res := newFloat(0).Quo(tSellAmount, tSupply)          // sellAmount / supply
	res.Sub(newFloat(1), res)                             // (1 - sellAmount / supply)
	res = bigfloat.Pow(res, newFloat(100/(float64(crr)))) // (1 - sellAmount / supply) ^ (100 / crr)
	res.Sub(newFloat(1), res)                             // (1 - (1 - sellAmount / supply) ^ (1 / (crr / 100)))
	res.Mul(res, tReserve)                                // reserve * (1 - (1 - sellAmount / supply) ^ (1 / (crr / 100)))

	converted, _ := res.Int(nil)
	result := sdkmath.NewIntFromBigInt(converted)

	return result
}

// CalculateSaleAmount is the reversed version of function CalculateSaleReturn.
// Deposit = -(-1 + (-(wantReceive - reserve)/reserve)^(1/crr)) * supply
func CalculateSaleAmount(supply sdkmath.Int, reserve sdkmath.Int, crr uint, wantReceive sdkmath.Int) sdkmath.Int {
	if wantReceive.Sign() == 0 {
		return sdkmath.NewInt(0)
	}

	if crr == 100 {
		return sdkmath.NewInt(1).Mul(wantReceive).Mul(supply).Quo(reserve)
	}

	tSupply := newFloat(0).SetInt(supply.BigInt())
	tReserve := newFloat(0).SetInt(reserve.BigInt())
	tWantReceive := newFloat(0).SetInt(wantReceive.BigInt())

	res := newFloat(0).Sub(tWantReceive, tReserve)      // (wantReceive - reserve)
	res.Neg(res)                                        // -(wantReceive - reserve)
	res.Quo(res, tReserve)                              // -(wantReceive - reserve)/reserve
	res = bigfloat.Pow(res, newFloat(float64(crr)/100)) // (-(wantReceive - reserve)/reserve)^(crr/100)
	res.Add(res, newFloat(-1))                          // -1 + (-(wantReceive - reserve)/reserve)^(crr/100)
	res.Neg(res)                                        // -(-1 + (-(wantReceive - reserve)/reserve)^(crr/100))
	res.Mul(res, tSupply)                               // -(-1 + (-(wantReceive - reserve)/reserve)^(crr/100)) * supply

	converted, _ := res.Int(nil)
	result := sdkmath.NewIntFromBigInt(converted)

	return result
}

////////////////////////////////////////////////////////////////

const precision = 100

func newFloat(x float64) *big.Float {
	return big.NewFloat(x).SetPrec(precision)
}
