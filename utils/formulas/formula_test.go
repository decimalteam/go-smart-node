package formulas

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PurchaseReturnData struct {
	Supply  sdkmath.Int
	Reserve sdkmath.Int
	Crr     uint
	Deposit sdkmath.Int
	Result  sdkmath.Int
}

func TestCalculatePurchaseReturn(t *testing.T) {
	data := []PurchaseReturnData{
		{
			Supply:  sdk.NewInt(1000000),
			Reserve: sdk.NewInt(100),
			Crr:     40,
			Deposit: sdk.NewInt(100),
			Result:  sdk.NewInt(319507),
		},
		{
			Supply:  sdk.NewInt(1000000),
			Reserve: sdk.NewInt(100),
			Crr:     100,
			Deposit: sdk.NewInt(100),
			Result:  sdk.NewInt(1000000),
		},
		{
			Supply:  sdk.NewInt(1000000),
			Reserve: sdk.NewInt(100),
			Crr:     100,
			Deposit: sdk.NewInt(0),
			Result:  sdk.NewInt(0),
		},
	}

	for _, item := range data {
		result := CalculatePurchaseReturn(item.Supply, item.Reserve, item.Crr, item.Deposit)

		if !result.Equal(item.Result) {
			t.Errorf("CalculatePurchaseReturn result is not correct. Expected %s, got %s", item.Result, result)
		}
	}
}

type PurchaseAmountData struct {
	Supply      sdkmath.Int
	Reserve     sdkmath.Int
	Crr         uint
	WantReceive sdkmath.Int
	Deposit     sdkmath.Int
}

func TestCalculatePurchaseAmount(t *testing.T) {
	data := []PurchaseAmountData{
		{
			Supply:      sdk.NewInt(1000000),
			Reserve:     sdk.NewInt(1000000),
			Crr:         40,
			WantReceive: sdk.NewInt(100),
			Deposit:     sdk.NewInt(250),
		},
		{
			Supply:      sdk.NewInt(1000000),
			Reserve:     sdk.NewInt(1000000),
			Crr:         100,
			WantReceive: sdk.NewInt(100),
			Deposit:     sdk.NewInt(100),
		},
		{
			Supply:      sdk.NewInt(1000000),
			Reserve:     sdk.NewInt(1000000),
			Crr:         100,
			WantReceive: sdk.NewInt(0),
			Deposit:     sdk.NewInt(0),
		},
	}

	for _, item := range data {
		deposit := CalculatePurchaseAmount(item.Supply, item.Reserve, item.Crr, item.WantReceive)

		if !deposit.Equal(item.Deposit) {
			t.Errorf("CalculatePurchaseAmount Deposit is not correct. Expected %s, got %s", item.Deposit, deposit)
		}
	}
}

type CalculateSaleReturnData struct {
	Supply     sdkmath.Int
	Reserve    sdkmath.Int
	Crr        uint
	SellAmount sdkmath.Int
	Result     sdkmath.Int
}

func TestCalculateSaleReturn(t *testing.T) {
	data := []CalculateSaleReturnData{
		{
			Supply:     sdk.NewInt(1000000),
			Reserve:    sdk.NewInt(100),
			Crr:        40,
			SellAmount: sdk.NewInt(1000000),
			Result:     sdk.NewInt(100),
		},
		{
			Supply:     sdk.NewInt(1000000),
			Reserve:    sdk.NewInt(100),
			Crr:        10,
			SellAmount: sdk.NewInt(100000),
			Result:     sdk.NewInt(65),
		},
		{
			Supply:     sdk.NewInt(1000000),
			Reserve:    sdk.NewInt(100),
			Crr:        10,
			SellAmount: sdk.NewInt(0),
			Result:     sdk.NewInt(0),
		},
		{
			Supply:     sdk.NewInt(1000000),
			Reserve:    sdk.NewInt(1000000),
			Crr:        100,
			SellAmount: sdk.NewInt(100),
			Result:     sdk.NewInt(100),
		},
	}

	for _, item := range data {
		result := CalculateSaleReturn(item.Supply, item.Reserve, item.Crr, item.SellAmount)

		if !result.Equal(item.Result) {
			t.Errorf("CalculateSaleReturn result is not correct. Expected %s, got %s", item.Result, result)
		}
	}
}

type CalculateBuyDepositData struct {
	Supply      sdkmath.Int
	Reserve     sdkmath.Int
	Crr         uint
	WantReceive sdkmath.Int
	Result      sdkmath.Int
}

func TestCalculateBuyDeposit(t *testing.T) {
	data := []CalculateBuyDepositData{
		{
			Supply:      sdk.NewInt(1000000),
			Reserve:     sdk.NewInt(100),
			Crr:         40,
			WantReceive: sdk.NewInt(10),
			Result:      sdk.NewInt(41268),
		},
		{
			Supply:      sdk.NewInt(1000000),
			Reserve:     sdk.NewInt(100),
			Crr:         10,
			WantReceive: sdk.NewInt(100),
			Result:      sdk.NewInt(1000000),
		},
		{
			Supply:      sdk.NewInt(1000000),
			Reserve:     sdk.NewInt(1000000),
			Crr:         100,
			WantReceive: sdk.NewInt(100),
			Result:      sdk.NewInt(100),
		},
		{
			Supply:      sdk.NewInt(1000000),
			Reserve:     sdk.NewInt(1000000),
			Crr:         100,
			WantReceive: sdk.NewInt(0),
			Result:      sdk.NewInt(0),
		},
	}

	for _, item := range data {
		result := CalculateSaleAmount(item.Supply, item.Reserve, item.Crr, item.WantReceive)

		if !result.Equal(item.Result) {
			t.Errorf("CalculateSaleAmount result is not correct. Expected %s, got %s", item.Result, result)
		}
	}
}
