package main

import (
	"fmt"
	"math/rand"
	"time"

	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
)

type ActionGenerator interface {
	// generate some action with random parameters
	Generate() Action
	Update(ui UpdateInfo)
}

type Action interface {
	// indicates than account can send transaction without errors
	// i.e. enought balance, account is owner of something...
	CanPerform(sa *StormAccount) bool
	// generate transaction data
	GenerateTx(sa *StormAccount) ([]byte, error)
}

// EmptyAction is used if generator can't create valid action
type EmptyAction struct{}

func (ea *EmptyAction) CanPerform(sa *StormAccount) bool {
	return false
}

func (ea *EmptyAction) GenerateTx(sa *StormAccount) ([]byte, error) {
	return nil, fmt.Errorf("empty action")
}

type UpdateInfo struct {
	Coins     []string
	Addresses []string
	FullCoins []dscApi.Coin
}

// TODO: do we need thread safety?

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
