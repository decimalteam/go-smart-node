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

//
type ActionReactor struct {
	wags []*WeightedAG
	wsum int64
}

type WeightedAG struct {
	Weight int64
	AG     ActionGenerator
}

//TODO: parameters for generator
func (ar *ActionReactor) Add(generatorName string, weight int64) error {
	var wag *WeightedAG = nil
	switch generatorName {
	case "CreateCoin":
		{
			wag = &WeightedAG{
				AG:     NewCreateCoinGenerator(3, 9, 100, 1000, 1000, 2000, 1000000, 2000000),
				Weight: weight,
			}
		}
	case "SendCoin":
		{
			wag = &WeightedAG{
				AG:     NewSendCoinGenerator(500, 20000),
				Weight: weight,
			}
		}
	case "BuyCoin":
		{
			wag = &WeightedAG{
				AG:     NewBuyCoinGenerator(500, 20000, "del"),
				Weight: weight,
			}
		}
	}
	if wag == nil {
		return fmt.Errorf("%s: unknown generator name", generatorName)
	}
	ar.wsum += weight
	ar.wags = append(ar.wags, wag)
	return nil
}

// choose generator and generate action
func (ar *ActionReactor) Generate() Action {
	w := rand.Int63n(ar.wsum)
	for _, wag := range ar.wags {
		if w < wag.Weight {
			return wag.AG.Generate()
		}
		w -= wag.Weight
	}
	// we can not be here, this is for stub
	return ar.wags[0].AG.Generate()
}

func (ar *ActionReactor) Update(ui UpdateInfo) {
	for _, wag := range ar.wags {
		wag.AG.Update(ui)
	}
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

func randomChoice(rnd *rand.Rand, list []string) string {
	return list[rnd.Intn(len(list))]
}
