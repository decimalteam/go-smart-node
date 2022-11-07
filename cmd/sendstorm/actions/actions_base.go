package actions

import (
	"fmt"
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Generate some action with random parameters
type ActionGenerator interface {
	Generate() Action
	Update(ui UpdateInfo)
}

type Action interface {
	// returns list of accounts than can make transaction
	// need to decrease count of invalid actions, required ownership of coin/nft, coin balance etc...
	ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount
	// generate signed transaction data
	GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error)
	// for debug puprposes
	String() string
}

// EmptyAction is used if generator can't create valid action
type EmptyAction struct{}

func (ea *EmptyAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	return []*stormTypes.StormAccount{}
}

func (ea *EmptyAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	return nil, fmt.Errorf("empty action")
}

func (ea *EmptyAction) String() string {
	return "EmptyAction{}"
}

// UpdateInfo contains all external updatable data for generators
type UpdateInfo struct {
	Coins                []string
	Addresses            []string
	FullCoins            []dscApi.Coin
	NFTs                 []*dscApi.NFTToken
	NFTSubTokenReserves  map[NFTSubTokenKey]sdk.Coin
	MultisigWallets      []dscApi.MultisigWallet
	MultisigTransactions []dscApi.MultisigTransactionInfo
	MultisigBalances     map[string]sdk.Coins
	Validators           []dscApi.Validator
	Delegations          []dscApi.Delegation
	Redelegations        []dscApi.Redelegation
	Undelegations        []dscApi.Undelegation
}

type NFTSubTokenKey struct {
	TokenID string
	ID      uint32
}

// TPS (transactions per second) limiter
type TPSLimiter struct {
	counter    int64
	limit      int64
	lastRefill time.Time
}

func NewTPSLimiter(limit int64) *TPSLimiter {
	return &TPSLimiter{
		counter:    limit,
		limit:      limit,
		lastRefill: time.Now(),
	}
}

func (t *TPSLimiter) CanMake() bool {
	if t.counter > 0 {
		t.counter--
		return true
	}
	if time.Now().Sub(t.lastRefill) < time.Second {
		return false
	}
	t.lastRefill = time.Now()
	t.counter = t.limit - 1
	return true
}

func (t *TPSLimiter) Limit() int64 {
	return t.limit
}

// extract part of stake for un/redelegation
func ExtractPartCoin(rnd *rand.Rand, coin sdk.Coin) sdk.Coin {
	fin := helpers.WeiToFinney(coin.Amount)
	if !fin.IsInt64() {
		return sdk.NewCoin(coin.Denom, sdk.ZeroInt())
	}
	amount := fin.Int64()
	if amount == 0 {
		return sdk.NewCoin(coin.Denom, sdk.ZeroInt())
	}
	// 50% - part, 50% - entire stake
	if rnd.Intn(2) == 0 {
		return sdk.NewCoin(coin.Denom, helpers.FinneyToWei(sdk.NewInt(RandomRange(rnd, 1, amount+1))))
	} else {
		return coin
	}
}
