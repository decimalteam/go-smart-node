package formulas

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/math"
)

const (
	precision = 100
)

func newFloat(x float64) *big.Float {
	return big.NewFloat(x).SetPrec(precision)
}

// Return = supply * ((1 + deposit / reserve) ^ (crr / 100) - 1)
// Рассчитывает сколько монет мы получим заплатив deposit DEL (Покупка формула 2)
func CalculatePurchaseReturn(supply sdk.Int, reserve sdk.Int, crr uint, deposit sdk.Int) sdk.Int {
	if deposit.Equal(sdk.NewInt(0)) {
		return sdk.NewInt(0)
	}

	if crr == 100 {
		return sdk.NewInt(1).Mul(supply).Mul(deposit).Quo(reserve)
	}

	tSupply := newFloat(0).SetInt(supply.BigInt())
	tReserve := newFloat(0).SetInt(reserve.BigInt())
	tDeposit := newFloat(0).SetInt(deposit.BigInt())

	res := newFloat(0).Quo(tDeposit, tReserve)      // deposit / reserve
	res.Add(res, newFloat(1))                       // 1 + (deposit / reserve)
	res = math.Pow(res, newFloat(float64(crr)/100)) // (1 + deposit / reserve) ^ (crr / 100)
	res.Sub(res, newFloat(1))                       // ((1 + deposit / reserve) ^ (crr / 100) - 1)
	res.Mul(res, tSupply)                           // supply * ((1 + deposit / reserve) ^ (crr / 100) - 1)

	converted, _ := res.Int(nil)
	result := sdk.NewIntFromBigInt(converted)

	return result
}

// reversed function CalculatePurchaseReturn
// deposit = reserve * (((wantReceive + supply) / supply)^(100 / crr) - 1)
// Рассчитывает сколько DEL надо заплатить , чтобы получить wantReceive монет (Покупка)

func CalculatePurchaseAmount(supply sdk.Int, reserve sdk.Int, crr uint, wantReceive sdk.Int) sdk.Int {
	if wantReceive.Equal(sdk.NewInt(0)) {
		return sdk.NewInt(0)
	}

	if crr == 100 {
		result := sdk.NewInt(1).Mul(wantReceive).Mul(reserve)

		return result.Quo(supply)
	}

	tSupply := newFloat(0).SetInt(supply.BigInt())
	tReserve := newFloat(0).SetInt(reserve.BigInt())
	tWantReceive := newFloat(0).SetInt(wantReceive.BigInt())

	res := newFloat(0).Add(tWantReceive, tSupply)   // reserve + supply
	res.Quo(res, tSupply)                           // (reserve + supply) / supply
	res = math.Pow(res, newFloat(100/float64(crr))) // ((reserve + supply) / supply)^(100/c)
	res.Sub(res, newFloat(1))                       // (((reserve + supply) / supply)^(100/c) - 1)
	res.Mul(res, tReserve)                          // reserve * (((reserve + supply) / supply)^(100/c) - 1)

	converted, _ := res.Int(nil)
	result := sdk.NewIntFromBigInt(converted)

	return result
}

// Return = reserve * (1 - (1 - sellAmount / supply) ^ (100 / crr))
// Рассчитывает сколько DEL вы получите, если продадите sellAmount монет. (Продажа)
func CalculateSaleReturn(supply sdk.Int, reserve sdk.Int, crr uint, sellAmount sdk.Int) sdk.Int {
	// special case for 0 sell amount
	if sellAmount.Equal(sdk.NewInt(0)) {
		return sdk.NewInt(0)
	}

	// special case for selling the entire supply
	if sellAmount.Equal(supply) {
		return reserve
	}

	if crr == 100 {
		ret := reserve.Mul(sellAmount)
		ret = ret.Quo(supply)

		return ret
	}

	tSupply := newFloat(0).SetInt(supply.BigInt())
	tReserve := newFloat(0).SetInt(reserve.BigInt())
	tSellAmount := newFloat(0).SetInt(sellAmount.BigInt())

	res := newFloat(0).Quo(tSellAmount, tSupply)      // sellAmount / supply
	res.Sub(newFloat(1), res)                         // (1 - sellAmount / supply)
	res = math.Pow(res, newFloat(100/(float64(crr)))) // (1 - sellAmount / supply) ^ (100 / crr)
	res64, _ := res.Float64()
	res.Sub(newFloat(1), newFloat(res64)) // (1 - (1 - sellAmount / supply) ^ (1 / (crr / 100)))
	res.Mul(res, tReserve)                // reserve * (1 - (1 - sellAmount / supply) ^ (1 / (crr / 100)))

	converted, _ := res.Int(nil)
	result := sdk.NewIntFromBigInt(converted)

	return result
}

// reversed function CalculateSaleReturn
// -(-1 + (-(wantReceive - reserve)/reserve)^(crr / 100)) * supply
// Рассчитывает сколько монет надо продать, чтобы получить wantReceive DEL. (Продажа 2)

func CalculateSaleAmount(supply sdk.Int, reserve sdk.Int, crr uint, wantReceive sdk.Int) sdk.Int {
	if wantReceive.Equal(sdk.NewInt(0)) {
		return sdk.NewInt(0)
	}

	if crr == 100 {
		ret := wantReceive.Mul(supply)
		ret = ret.Quo(reserve)

		return ret
	}

	tSupply := newFloat(0).SetInt(supply.BigInt())
	tReserve := newFloat(0).SetInt(reserve.BigInt())
	tWantReceive := newFloat(0).SetInt(wantReceive.BigInt())

	res := newFloat(0).Sub(tWantReceive, tReserve)  // (wantReceive - reserve)
	res.Neg(res)                                    // -(wantReceive - reserve)
	res.Quo(res, tReserve)                          // -(wantReceive - reserve)/reserve
	res = math.Pow(res, newFloat(float64(crr)/100)) // (-(wantReceive - reserve)/reserve)^(crr/100)
	res.Add(res, newFloat(-1))                      // -1 + (-(wantReceive - reserve)/reserve)^(crr/100)
	res.Neg(res)                                    // -(-1 + (-(wantReceive - reserve)/reserve)^(crr/100))
	res.Mul(res, tSupply)                           // -(-1 + (-(wantReceive - reserve)/reserve)^(crr/100)) * supply

	converted, _ := res.Int(nil)
	result := sdk.NewIntFromBigInt(converted)

	return result
}

func GetReserveLimitFromCRR(crr uint) sdk.Int {
	convert, _ := sdk.NewIntFromString("1000000000000000000")
	// CRR always between 10 and 100
	if 10 <= crr && crr <= 19 {
		return sdk.NewInt(100000).Mul(convert)
	} else if 20 <= crr && crr <= 29 {
		return sdk.NewInt(90000).Mul(convert)
	} else if 30 <= crr && crr <= 39 {
		return sdk.NewInt(80000).Mul(convert)
	} else if 40 <= crr && crr <= 49 {
		return sdk.NewInt(70000).Mul(convert)
	} else if 50 <= crr && crr <= 59 {
		return sdk.NewInt(60000).Mul(convert)
	} else if 60 <= crr && crr <= 69 {
		return sdk.NewInt(50000).Mul(convert)
	} else if 70 <= crr && crr <= 79 {
		return sdk.NewInt(40000).Mul(convert)
	} else if 80 <= crr && crr <= 89 {
		return sdk.NewInt(30000).Mul(convert)
	} else {
		return sdk.NewInt(10000).Mul(convert)
	}
}
