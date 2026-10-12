package app

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

// StateSnapshot is a full, comparable picture of delegation state on both sides
// (offline tooling for before/after upgrade comparisons).
type StateSnapshot struct {
	Height int64 `json:"height"`
	Time   int64 `json:"time"`

	// node
	Validators      []SnapValidator   `json:"validators"`
	Delegations     []SnapDelegation  `json:"delegations"`
	LastTotalPower  string            `json:"last_total_power"`
	CustomStaked    map[string]string `json:"custom_coin_staked"`
	Undelegations   map[string]string `json:"undelegations_by_denom"`
	UndelegCount    int               `json:"undelegation_entries"`
	Redelegations   map[string]string `json:"redelegations_by_denom"`
	RedelegCount    int               `json:"redelegation_entries"`
	Supply          map[string]string `json:"supply"`

	// contract
	Contract SnapContract `json:"contract"`
}

type SnapValidator struct {
	Operator string `json:"operator"`
	Status   string `json:"status"`
	Online   bool   `json:"online"`
	Jailed   bool   `json:"jailed"`
	Power    int64  `json:"power"`
	Rewards  string `json:"rewards"`
	Count    uint32 `json:"delegations_count"`
}

type SnapDelegation struct {
	Validator string           `json:"validator"`
	Delegator string           `json:"delegator"`
	Denom     string           `json:"denom"`
	Type      string           `json:"type"`
	Amount    string           `json:"amount"`
	Holds     map[string]string `json:"holds,omitempty"` // end -> amount (start in key suffix)
}

type SnapContract struct {
	Address        string            `json:"address"`
	Impl           string            `json:"impl"`
	Owner          string            `json:"owner"`
	Paused         bool              `json:"paused"`
	Stakes         []SnapStake       `json:"stakes"` // non-zero coin stakes
	ZeroStakes     int               `json:"zero_stakes"`
	NFTStakes      int               `json:"nft_stakes"`
	FrozenLive     int               `json:"frozen_live"`
	FrozenOld      int               `json:"frozen_deprecated"`
	FrozenPending  map[string]string `json:"frozen_pending_by_token"` // freezeStatus==Frozen
	FrozenHash     string            `json:"frozen_queue_digest"`
	QueueLen       int               `json:"auto_unbond_queue_len"`
	QueueByToken   map[string]string `json:"auto_unbond_queue_by_token"`
	TokenBalances  map[string]string `json:"token_balances"` // balanceOf(contract) per staked token
	StakedByToken  map[string]string `json:"staked_by_token"`
}

type SnapStake struct {
	Validator     string `json:"validator"`
	Delegator     string `json:"delegator"`
	Token         string `json:"token"`
	Amount        string `json:"amount"`
	HoldTimestamp int64  `json:"hold_timestamp"`
	HoldStartTime int64  `json:"hold_start_time"`
	PenaltyIndex  string `json:"penalty_index"`
}

const erc20BalanceABI = `[{"type":"function","name":"balanceOf","stateMutability":"view","inputs":[{"name":"a","type":"address"}],"outputs":[{"name":"","type":"uint256"}]},
 {"type":"function","name":"paused","stateMutability":"view","inputs":[],"outputs":[{"name":"","type":"bool"}]},
 {"type":"function","name":"owner","stateMutability":"view","inputs":[],"outputs":[{"name":"","type":"address"}]}]`

func addBig(m map[string]string, k string, v *big.Int) {
	cur, _ := new(big.Int).SetString(m[k], 10)
	if cur == nil {
		cur = new(big.Int)
	}
	m[k] = cur.Add(cur, v).String()
}

// TakeStateSnapshot reads everything; it never writes.
func TakeStateSnapshot(ctx sdk.Context, app *DSC) (StateSnapshot, error) {
	s := StateSnapshot{Height: ctx.BlockHeight(), Time: ctx.BlockTime().Unix(),
		CustomStaked: map[string]string{}, Undelegations: map[string]string{}, Redelegations: map[string]string{}, Supply: map[string]string{}}
	vk := app.ValidatorKeeper

	counts := vk.GetAllDelegationsCount(ctx)
	for _, v := range vk.GetAllValidators(ctx) {
		s.Validators = append(s.Validators, SnapValidator{Operator: v.OperatorAddress, Status: v.Status.String(), Online: v.Online,
			Jailed: v.Jailed, Power: v.Stake, Rewards: v.Rewards.String(), Count: counts[v.OperatorAddress]})
	}
	sort.Slice(s.Validators, func(i, j int) bool { return s.Validators[i].Operator < s.Validators[j].Operator })
	for _, d := range vk.GetAllDelegations(ctx) {
		sd := SnapDelegation{Validator: d.Validator, Delegator: d.Delegator, Denom: d.Stake.Stake.Denom, Type: d.Stake.Type.String(), Amount: d.Stake.Stake.Amount.String()}
		if d.Stake.Type == validatortypes.StakeType_NFT {
			sd.Denom = d.Stake.ID
		}
		if len(d.Stake.Holds) > 0 {
			sd.Holds = map[string]string{}
			for i, h := range d.Stake.Holds {
				sd.Holds[fmt.Sprintf("%d@%d#%d", h.HoldEndTime, h.HoldStartTime, i)] = h.Amount.String()
			}
		}
		s.Delegations = append(s.Delegations, sd)
	}
	sort.Slice(s.Delegations, func(i, j int) bool {
		a, b := s.Delegations[i], s.Delegations[j]
		return a.Validator+a.Delegator+a.Denom < b.Validator+b.Delegator+b.Denom
	})
	s.LastTotalPower = vk.GetLastTotalPower(ctx).String()
	for denom, amt := range vk.GetAllCustomCoinsStaked(ctx) {
		s.CustomStaked[denom] = amt.String()
	}
	vk.IterateUndelegations(ctx, func(_ int64, ubd validatortypes.Undelegation) bool {
		for _, e := range ubd.Entries {
			s.UndelegCount++
			addBig(s.Undelegations, e.Stake.Stake.Denom, e.Stake.Stake.Amount.BigInt())
		}
		return false
	})
	vk.IterateRedelegations(ctx, func(_ int64, red validatortypes.Redelegation) bool {
		for _, e := range red.Entries {
			s.RedelegCount++
			addBig(s.Redelegations, e.Stake.Stake.Denom, e.Stake.Stake.Amount.BigInt())
		}
		return false
	})
	app.BankKeeper.IterateTotalSupply(ctx, func(c sdk.Coin) bool {
		s.Supply[c.Denom] = c.Amount.String()
		return false
	})

	// contract
	proxyHex, err := contracts.GetAddressFromContractCenter(ctx, &app.EvmKeeper, contracts.NameOfSlugForGetAddressDelegation)
	if err != nil {
		return s, err
	}
	proxy := common.HexToAddress(proxyHex)
	storage := readEVMContractStorage(ctx, app, proxy)
	res := stakescan.Reconstruct(storage, stakescan.DelegationBase)
	c := SnapContract{Address: proxy.Hex(), FrozenPending: map[string]string{}, QueueByToken: map[string]string{},
		TokenBalances: map[string]string{}, StakedByToken: map[string]string{}}
	c.Impl = common.BytesToAddress(app.EvmKeeper.GetState(ctx, proxy, eip1967ImplSlot).Bytes()).Hex()
	c.Owner = common.BytesToAddress(app.EvmKeeper.GetState(ctx, proxy, ozOwnableSlot).Bytes()).Hex()
	tokens := map[common.Address]bool{}
	for _, cs := range res.CoinStakes {
		if cs.TokenType != contractTokenTypeDRC20 && cs.TokenType != contractTokenTypeDEL {
			c.NFTStakes++
			continue
		}
		tokens[cs.Token] = true
		if cs.Amount.Sign() == 0 {
			c.ZeroStakes++
			continue
		}
		hs := int64(0)
		if cs.HoldStartTime != nil {
			hs = cs.HoldStartTime.Int64()
		}
		pi := "0"
		if cs.PenaltyIndex != nil {
			pi = cs.PenaltyIndex.String()
		}
		c.Stakes = append(c.Stakes, SnapStake{Validator: cs.Validator.Hex(), Delegator: cs.Delegator.Hex(), Token: cs.Token.Hex(),
			Amount: cs.Amount.String(), HoldTimestamp: cs.HoldTimestamp.Int64(), HoldStartTime: hs, PenaltyIndex: pi})
		addBig(c.StakedByToken, cs.Token.Hex(), cs.Amount)
	}
	c.NFTStakes += len(res.NFTStakes)
	sort.Slice(c.Stakes, func(i, j int) bool {
		a, b := c.Stakes[i], c.Stakes[j]
		if a.Validator+a.Delegator+a.Token != b.Validator+b.Delegator+b.Token {
			return a.Validator+a.Delegator+a.Token < b.Validator+b.Delegator+b.Token
		}
		return a.HoldTimestamp < b.HoldTimestamp
	})
	c.FrozenLive, c.FrozenOld = len(res.FrozenLive), len(res.FrozenDeprecated)
	// The effective queue: migrated prefix + not-yet-migrated tail of the old array.
	frozen := append([]stakescan.FrozenStake{}, res.FrozenLive...)
	if len(res.FrozenDeprecated) > len(frozen) {
		frozen = append(frozen, res.FrozenDeprecated[len(frozen):]...)
	}
	var digest strings.Builder
	for _, f := range frozen {
		fmt.Fprintf(&digest, "%d:%s:%s:%s:%s:%d:%d:%s;", f.Index, f.Stake.Validator.Hex(), f.Stake.Delegator.Hex(), f.Stake.Token.Hex(),
			f.Stake.Amount, f.FreezeStatus, f.FreezeType, f.UnfreezeTimestamp)
		if f.FreezeStatus == 1 {
			addBig(c.FrozenPending, f.Stake.Token.Hex(), f.Stake.Amount)
		}
		tokens[f.Stake.Token] = true
	}
	c.FrozenHash = fmt.Sprintf("%x", sha256Hex(digest.String()))
	queue := stakescan.ReadAutoUnbondQueueAt(storage, stakescan.AutoUnbondQueueSlotMain)
	c.QueueLen = len(queue)
	for _, e := range queue {
		addBig(c.QueueByToken, e.Token.Hex(), e.Amount)
	}

	a, err := abi.JSON(strings.NewReader(erc20BalanceABI))
	if err != nil {
		return s, err
	}
	from := common.Address(validatortypes.ModuleAddress)
	if res, err := app.EvmKeeper.CallEVM(ctx, a, from, proxy, false, "paused"); err == nil && !res.Failed() {
		if out, err := a.Unpack("paused", res.Ret); err == nil {
			c.Paused = out[0].(bool)
		}
	}
	toks := make([]common.Address, 0, len(tokens))
	for t := range tokens {
		toks = append(toks, t)
	}
	sort.Slice(toks, func(i, j int) bool { return toks[i].Hex() < toks[j].Hex() })
	for _, t := range toks {
		res, err := app.EvmKeeper.CallEVM(ctx, a, from, t, false, "balanceOf", proxy)
		if err != nil || res.Failed() {
			c.TokenBalances[t.Hex()] = "error"
			continue
		}
		out, err := a.Unpack("balanceOf", res.Ret)
		if err != nil {
			c.TokenBalances[t.Hex()] = "error"
			continue
		}
		c.TokenBalances[t.Hex()] = out[0].(*big.Int).String()
	}
	s.Contract = c
	return s, nil
}

func sha256Hex(s string) []byte {
	h := sha256.Sum256([]byte(s))
	return h[:]
}
