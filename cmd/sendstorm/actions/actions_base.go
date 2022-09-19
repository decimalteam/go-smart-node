package actions

import (
	"fmt"
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
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
	MultisigTransactions []dscApi.MultisigTransaction
	MultisigBalances     map[string]sdk.Coins
}

type NFTSubTokenKey struct {
	Denom   string
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

// helpers
const charsAll = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const charsAbc = "abcdefghijklmnopqrstuvwxyz"

// returns random number in range [low,up)
func RandomRange(rnd *rand.Rand, bottom, up int64) int64 {
	return rnd.Int63n(up-bottom) + bottom
}

// returns random string length n
func RandomString(rnd *rand.Rand, n int64, source string) string {
	var letters = []rune(source)
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rnd.Intn(len(letters))]
	}
	return string(s)
}

// TODO: need generics
func RandomChoice(rnd *rand.Rand, list []string) string {
	return list[rnd.Intn(len(list))]
}

// Return random sublist (copy)
func RandomSublist(rnd *rand.Rand, list []uint64) []uint64 {
	if len(list) == 0 {
		return []uint64{}
	}
	if len(list) == 1 {
		return []uint64{list[0]}
	}
	// random indexes to choose
	ids := make([]int, len(list))
	for i := range list {
		ids[i] = i
	}
	rnd.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})
	n := int(RandomRange(rnd, 1, int64(len(list)+1)))
	result := make([]uint64, n)
	for i := 0; i < n; i++ {
		result[i] = list[ids[i]]
	}
	return result
}

// Return random sublist (copy)
func RandomSublist32(rnd *rand.Rand, list []uint32) []uint32 {
	if len(list) == 0 {
		return []uint32{}
	}
	if len(list) == 1 {
		return []uint32{list[0]}
	}
	// random indexes to choose
	ids := make([]int, len(list))
	for i := range list {
		ids[i] = i
	}
	rnd.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})
	n := int(RandomRange(rnd, 1, int64(len(list)+1)))
	result := make([]uint32, n)
	for i := 0; i < n; i++ {
		result[i] = list[ids[i]]
	}
	return result
}
