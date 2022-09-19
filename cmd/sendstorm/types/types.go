package types

import (
	"fmt"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscWallet "bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
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

func (sa *StormAccount) BalanceForCoin(denom string) sdk.Int {
	for _, b := range sa.currentBalance {
		if b.Denom == denom {
			return b.Amount
		}
	}
	return sdk.NewInt(0)
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
	DelPrice sdk.Dec
	Params   feetypes.Params
}

func NewFeeConfiguration() *FeeConfiguration {
	return &FeeConfiguration{}
}

func (fc *FeeConfiguration) Update(api *dscApi.API) error {
	delPrice, params, err := api.GetFeeParams("del", "usd")
	if err != nil {
		return err
	}
	fc.DelPrice = delPrice
	fc.Params = params
	return nil
}
