package main

import (
	"fmt"
	"math/rand"
	"time"

	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
)

// Generate some action with random parameters
type ActionGenerator interface {
	Generate() Action
	Update(ui UpdateInfo)
}

type Action interface {
	// returns list of accounts than can make transaction
	// need to decrease count of invalid actions, required ownership of coin/nft, coin balance etc...
	ChooseAccounts(saList []*StormAccount) []*StormAccount
	// generate signed transaction data
	GenerateTx(sa *StormAccount) ([]byte, error)
	// for debug puprposes
	String() string
}

// EmptyAction is used if generator can't create valid action
type EmptyAction struct{}

func (ea *EmptyAction) ChooseAccounts(saList []*StormAccount) []*StormAccount {
	return []*StormAccount{}
}

func (ea *EmptyAction) CanPerform(sa *StormAccount) bool {
	return false
}

func (ea *EmptyAction) GenerateTx(sa *StormAccount) ([]byte, error) {
	return nil, fmt.Errorf("empty action")
}

func (ea *EmptyAction) String() string {
	return "EmptyAction{}"
}

// UpdateInfo contains all external updatable data for generators
type UpdateInfo struct {
	Coins     []string
	Addresses []string
	FullCoins []dscApi.Coin
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
func randomRange(rnd *rand.Rand, bottom, up int64) int64 {
	return rnd.Int63n(up-bottom) + bottom
}

// returns random string length n
func randomString(rnd *rand.Rand, n int64, source string) string {
	var letters = []rune(source)
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rnd.Intn(len(letters))]
	}
	return string(s)
}

// TODO: need generics
func randomChoice(rnd *rand.Rand, list []string) string {
	return list[rnd.Intn(len(list))]
}
