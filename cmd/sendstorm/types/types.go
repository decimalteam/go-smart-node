package types

import (
	"fmt"
	"math/rand"
	"sync"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/strings"

	appAnte "bitbucket.org/decimalteam/go-smart-node/app/ante"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	dscWallet "bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
	"bitbucket.org/decimalteam/go-smart-node/utils/formulas"
	feetypes "bitbucket.org/decimalteam/go-smart-node/x/fee/types"
)

type StormAccount struct {
	account        *dscWallet.Account
	api            *dscApi.API
	currentBalance sdk.Coins
	dirty          bool // marks last transaction failure and need to update balance + nonce
	feeDenom       string
	mu             sync.Mutex
}

func NewStormAccount(mnemonic string, api *dscApi.API) (*StormAccount, error) {
	var result StormAccount
	var err error
	result.account, err = dscWallet.NewAccountFromMnemonicWords(mnemonic, "")
	if err != nil {
		return nil, err
	}
	result.api = api
	result.feeDenom = api.BaseCoin()
	result.dirty = true // need to get balance and nonce
	return &result, nil
}

func (sa *StormAccount) UpdateNumberSequence() error {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	an, as, err := sa.api.AccountNumberAndSequence(sa.account.Address())
	if err != nil {
		return fmt.Errorf("%w: AccountNumberAndSequence", err)
	}
	sa.account = sa.account.WithAccountNumber(an).WithSequence(as).WithChainID(sa.api.ChainID())
	sa.dirty = false
	return nil
}

func (sa *StormAccount) LastHeight() int64 {
	return sa.api.GetLastHeight()
}

func (sa *StormAccount) UpdateBalance() error {
	var err error
	sa.mu.Lock()
	defer sa.mu.Unlock()
	sa.currentBalance, err = sa.api.AddressBalance(sa.account.Address())
	if err != nil {
		return fmt.Errorf("%w: AddressBalance", err)
	}
	return nil
}

func (sa *StormAccount) MarkDirty() {
	sa.dirty = true
}

func (sa *StormAccount) IsDirty() bool {
	return sa.dirty
}

func (sa *StormAccount) IncrementSequence() {
	sa.account.IncrementSequence()
}

func (sa *StormAccount) BalanceForCoin(denom string) sdkmath.Int {
	for _, b := range sa.currentBalance {
		if b.Denom == denom {
			return b.Amount
		}
	}
	return sdkmath.NewInt(0)
}

func (sa *StormAccount) Address() string {
	return sa.account.Address()
}

func (sa *StormAccount) FeeDenom() string {
	return sa.feeDenom
}

func (sa *StormAccount) Account() *dscWallet.Account {
	return sa.account
}

/////////////////////////
// Fee configuration
/////////////////////////

type FeeConfiguration struct {
	BaseDenom            string
	DelPrice             sdk.Dec
	Params               feetypes.Params
	KnownCoins           []dscApi.Coin
	UseCustomCoinsForFee bool
}

func NewFeeConfiguration(customForFee bool) *FeeConfiguration {
	return &FeeConfiguration{
		UseCustomCoinsForFee: customForFee,
	}
}

func (fc *FeeConfiguration) Update(api *dscApi.API) error {
	err := api.GetParameters()
	if err != nil {
		return err
	}
	fc.BaseDenom = api.BaseCoin()

	delPrice, params, err := api.GetFeeParams(fc.BaseDenom, "usd")
	if err != nil {
		return err
	}
	fc.DelPrice = delPrice
	fc.Params = params

	fc.KnownCoins, err = api.Coins()
	if err != nil {
		return err
	}

	return nil
}

func (fc *FeeConfiguration) MakeTransaction(sa *StormAccount, msg sdk.Msg) ([]byte, error) {
	if !fc.UseCustomCoinsForFee || len(fc.KnownCoins) < 2 {
		tx, err := dscTx.BuildTransaction(sa.Account(), []sdk.Msg{msg}, "", fc.BaseDenom, sa.api.GetFeeCalculationOptions())
		if err != nil {
			return nil, err
		}

		err = tx.SignTransaction(sa.Account())
		if err != nil {
			return nil, err
		}
		return tx.BytesToSend()
	} else {
		// preparation
		tx, err := dscTx.BuildTransaction(sa.Account(), []sdk.Msg{msg}, "", fc.BaseDenom, sa.api.GetFeeCalculationOptions())
		if err != nil {
			return nil, err
		}
		err = tx.SignTransaction(sa.Account())
		if err != nil {
			return nil, err
		}
		bz, err := tx.BytesToSend()
		if err != nil {
			return nil, err
		}
		denoms := []string{}
		for _, coinInfo := range fc.KnownCoins {
			if !sa.BalanceForCoin(coinInfo.Denom).IsZero() {
				denoms = append(denoms, coinInfo.Denom)
			}
		}
		// + 50 bytes for denom length, amount length
		comms, err := calculateCommission(msg, int64(len(bz)+50), fc.BaseDenom, fc.KnownCoins, denoms, sa.api.GetFeeCalculationOptions())
		if err != nil {
			return nil, err
		}
		var commCandidates sdk.Coins
		for _, cm := range comms {
			if sa.BalanceForCoin(cm.Denom).GTE(cm.Amount) {
				commCandidates = commCandidates.Add(cm)
			}
		}
		if len(commCandidates) == 0 {
			return nil, fmt.Errorf("not enough balance for fee for account '%s'", sa.Address())
		}
		fee := commCandidates[rand.Intn(len(commCandidates))]
		// final build
		tx, err = dscTx.BuildTransaction(sa.Account(), []sdk.Msg{msg}, "", fc.BaseDenom, sa.api.GetFeeCalculationOptions())
		if err != nil {
			return nil, err
		}
		tx.SetFeeAmount(sdk.NewCoins(fee))
		err = tx.SignTransaction(sa.Account())
		if err != nil {
			return nil, err
		}
		return tx.BytesToSend()
	}
}

func calculateCommission(msg sdk.Msg, txBytesLen int64, baseDenom string, fullCoins []dscApi.Coin,
	denoms []string, opts *dscTx.FeeCalculationOptions) (sdk.Coins, error) {
	commmissionInBase, err := appAnte.CalculateFee(opts.AppCodec, []sdk.Msg{msg}, txBytesLen, opts.DelPrice, opts.FeeParams)
	if err != nil {
		return sdk.NewCoins(), err
	}
	result := sdk.NewCoins(sdk.NewCoin(baseDenom, commmissionInBase))
	for _, coinInfo := range fullCoins {
		if !strings.StringInSlice(coinInfo.Denom, denoms) {
			continue
		}
		if coinInfo.Denom == baseDenom {
			continue
		}
		amount := formulas.CalculateSaleAmount(coinInfo.Volume, coinInfo.Reserve, uint(coinInfo.CRR), commmissionInBase)
		result = append(result, sdk.NewCoin(coinInfo.Denom, amount))
	}
	return result, nil
}
