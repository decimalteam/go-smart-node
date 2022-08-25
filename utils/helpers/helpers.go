package helpers

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	bigE15 = new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil)
	bigE18 = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	sdkE15 = sdk.NewIntFromBigInt(bigE15)
	sdkE18 = sdk.NewIntFromBigInt(bigE18)
)

func BipToPip(bip sdkmath.Int) sdkmath.Int {
	return bip.Mul(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
}

// EtherToWei convert number 1 to 1 * 10^18
func EtherToWei(ether sdkmath.Int) sdkmath.Int {
	return ether.Mul(sdkE18)
}

// FinneyToWei convert number 1 to 1 * 10^15
func FinneyToWei(finney sdkmath.Int) sdkmath.Int {
	return finney.Mul(sdkE15)
}

// WeiToFinney convert 1 * 10^15 to 1
func WeiToFinney(wei sdkmath.Int) sdkmath.Int {
	return wei.Quo(sdkE15)
}

// WeiToEther convert 1 * 10^18 to 1
func WeiToEther(wei sdkmath.Int) sdkmath.Int {
	return wei.Quo(sdkE18)
}

// JoinAccAddresses returns string containing all provided address joined with ",".
func JoinAccAddresses(values []sdk.AccAddress) string {
	var sb strings.Builder
	for i, v := range values {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(v.String())
	}
	return sb.String()
}

// JoinUints returns string containing all provided uint values joined with ",".
func JoinUints64(values []uint64) string {
	var sb strings.Builder
	for i, v := range values {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatUint(v, 10))
	}
	return sb.String()
}

// GetReserveLimitFromCRR returns coin reserve limit for specific CRR value.
func GetReserveLimitFromCRR(crr uint) sdkmath.Int {
	// CRR must be in range [10; 100]
	if crr < 10 || crr > 100 {
		return sdk.NewInt(0)
	} else if crr < 20 {
		return sdk.NewInt(100000).Mul(sdkE18)
	} else if crr < 30 {
		return sdk.NewInt(90000).Mul(sdkE18)
	} else if crr < 40 {
		return sdk.NewInt(80000).Mul(sdkE18)
	} else if crr < 50 {
		return sdk.NewInt(70000).Mul(sdkE18)
	} else if crr < 60 {
		return sdk.NewInt(60000).Mul(sdkE18)
	} else if crr < 70 {
		return sdk.NewInt(50000).Mul(sdkE18)
	} else if crr < 80 {
		return sdk.NewInt(40000).Mul(sdkE18)
	} else if crr < 90 {
		return sdk.NewInt(30000).Mul(sdkE18)
	} else {
		return sdk.NewInt(10000).Mul(sdkE18)
	}
}

// DurationToString converts provided duration to human readable string presentation.
func DurationToString(d time.Duration) string {
	ns := time.Duration(d.Nanoseconds())
	ms := float64(ns) / 1000000.0
	var unit string
	var amount string
	if ns < time.Microsecond {
		unit = "ns"
		amount, unit = humanize.CommafWithDigits(float64(ns), 0), "ns"
	} else if ns < time.Millisecond {
		unit = "μs"
		amount, unit = humanize.CommafWithDigits(ms*1000.0, 3), "μs"
	} else if ns < time.Second {
		unit = "ms"
		amount, unit = humanize.CommafWithDigits(ms, 3), "ms"
	} else if ns < time.Minute {
		unit = "s"
		amount, unit = humanize.CommafWithDigits(ms/1000.0, 3), "s"
	} else if ns < time.Hour {
		unit = "m"
		amount, unit = humanize.CommafWithDigits(ms/60000.0, 3), "m"
	} else if ns < 24*time.Hour {
		unit = "h"
		amount, unit = humanize.CommafWithDigits(ms/3600000.0, 3), "h"
	} else {
		days := ms / 86400000.0
		unit = "day"
		if days > 1 {
			unit = "days"
		}
		amount = humanize.CommafWithDigits(days, 3)
	}
	return fmt.Sprintf("%s %s", amount, unit)
}
