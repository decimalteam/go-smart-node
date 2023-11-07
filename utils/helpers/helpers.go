package helpers

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/cometbft/cometbft/crypto/tmhash"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	bigE15 = new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil)
	bigE18 = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	sdkE15 = sdkmath.NewIntFromBigInt(bigE15)
	sdkE18 = sdkmath.NewIntFromBigInt(bigE18)
)

func BipToPip(bip sdkmath.Int) sdkmath.Int {
	return bip.Mul(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
}

// DecToIntWithE18 converts Dec to Int (dec*10^18 and truncate).
func DecToIntWithE18(dec sdk.Dec) sdkmath.Int {
	return dec.MulInt(sdkE18).TruncateInt()
}

// DecToDecWithE18 converts Dec to Dec*10^18.
func DecToDecWithE18(dec sdk.Dec) sdk.Dec {
	return dec.MulInt(sdkE18)
}

// EtherToWei converts number 1 to 1 * 10^18.
func EtherToWei(ether sdkmath.Int) sdkmath.Int {
	return ether.Mul(sdkE18)
}

// FinneyToWei convert number 1 to 1 * 10^15
func FinneyToWei(finney sdkmath.Int) sdkmath.Int {
	return finney.Mul(sdkE15)
}

// WeiToFinney converts number 1 * 10^15 to 1.
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
	var limit int64
	switch {
	case crr < 10 || crr > 100:
		limit = 0
	case crr < 20:
		limit = 100_000
	case crr < 30:
		limit = 90_000
	case crr < 40:
		limit = 80_000
	case crr < 50:
		limit = 70_000
	case crr < 60:
		limit = 60_000
	case crr < 70:
		limit = 50_000
	case crr < 80:
		limit = 40_000
	case crr < 90:
		limit = 30_000
	default:
		limit = 10_000
	}
	return EtherToWei(sdkmath.NewInt(limit))
}

// DurationToString converts provided duration to human readable string presentation.
func DurationToString(d time.Duration) string {
	ns := time.Duration(d.Nanoseconds())
	ms := float64(ns) / 1000000.0
	var unit string
	var amount string
	if ns < time.Microsecond {
		amount, unit = humanize.CommafWithDigits(float64(ns), 0), "ns"
	} else if ns < time.Millisecond {
		amount, unit = humanize.CommafWithDigits(ms*1000.0, 3), "μs"
	} else if ns < time.Second {
		amount, unit = humanize.CommafWithDigits(ms, 3), "ms"
	} else if ns < time.Minute {
		amount, unit = humanize.CommafWithDigits(ms/1000.0, 3), "s"
	} else if ns < time.Hour {
		amount, unit = humanize.CommafWithDigits(ms/60000.0, 3), "m"
	} else if ns < 24*time.Hour {
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

// CalcHashSHA256 returns sha256 hash (32 bytes) calculated from specified string.
func CalcHashSHA256(str string) []byte {
	h := tmhash.New()
	_, err := h.Write([]byte(str))
	if err != nil {
		panic(err)
	}
	return h.Sum(nil)
}
