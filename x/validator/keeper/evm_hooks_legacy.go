package keeper

import (
	"fmt"
	"math/big"
	"strings"

	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegation"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegationNft"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegationv7"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// The mainnet DecimalDelegation contract emits events whose Stake tuple has 7 fields
// (no holdStartTime), while the node binding expects the 8-field layout. The two
// layouts hash to different event IDs, so without translation every Stake-carrying
// event is unknown to the node and silently skipped. normalizeDelegationLogs rewrites
// such legacy logs into their 8-field equivalent (holdStartTime = 0, which makes
// resolveHoldStartTime fall back to the block time exactly as the pre-hold-start
// node did) so the existing dispatch below handles both layouts.

// normalizeDelegationLogs returns a copy of logs where every legacy-layout event emitted
// by the delegation contract is re-encoded in the current layout. Logs of other
// contracts and current-layout logs are passed through untouched. An event from the
// delegation contract that neither layout knows is an error: accepting the EVM side of
// such a tx while the node ignores it is exactly the desync this guards against.
// The receipt's own logs are never mutated (they are part of the tx result).
func normalizeDelegationLogs(logs []*ethtypes.Log, addressDelegation string, current, legacy *abi.ABI) ([]*ethtypes.Log, error) {
	out := make([]*ethtypes.Log, 0, len(logs))
	for _, log := range logs {
		if log == nil || len(log.Topics) == 0 || strings.ToLower(log.Address.String()) != addressDelegation {
			out = append(out, log)
			continue
		}
		if _, err := current.EventByID(log.Topics[0]); err == nil {
			out = append(out, log)
			continue
		}
		ev, err := legacy.EventByID(log.Topics[0])
		if err != nil {
			return nil, fmt.Errorf("delegation contract event %s is unknown to the node", log.Topics[0].Hex())
		}
		converted, err := convertLegacyDelegationLog(log, ev.Name, current, legacy)
		if err != nil {
			return nil, fmt.Errorf("convert legacy delegation event %s: %w", ev.Name, err)
		}
		out = append(out, converted)
	}
	return out, nil
}

func convertLegacyDelegationLog(log *ethtypes.Log, name string, current, legacy *abi.ABI) (*ethtypes.Log, error) {
	var args []interface{}
	switch name {
	case "StakeUpdated":
		var e delegationv7.DelegationStakeUpdated
		if err := contracts.UnpackLog(legacy, &e, name, log); err != nil {
			return nil, err
		}
		args = []interface{}{e.StakeId, stakeFromLegacy(e.Stake), e.IsNew}
	case "StakeHolded":
		var e delegationv7.DelegationStakeHolded
		if err := contracts.UnpackLog(legacy, &e, name, log); err != nil {
			return nil, err
		}
		args = []interface{}{e.StakeId, stakeFromLegacy(e.Stake), e.IsNew}
	case "StakeReset":
		var e delegationv7.DelegationStakeReset
		if err := contracts.UnpackLog(legacy, &e, name, log); err != nil {
			return nil, err
		}
		args = []interface{}{e.StakeId, stakeFromLegacy(e.Stake)}
	case "WithdrawRequest":
		var e delegationv7.DelegationWithdrawRequest
		if err := contracts.UnpackLog(legacy, &e, name, log); err != nil {
			return nil, err
		}
		args = []interface{}{e.StakeId, e.StakeIndex, frozenStakeFromLegacy(e.FrozenStake)}
	case "TransferRequest":
		var e delegationv7.DelegationTransferRequest
		if err := contracts.UnpackLog(legacy, &e, name, log); err != nil {
			return nil, err
		}
		args = []interface{}{e.StakeId, e.StakeIndex, frozenStakeFromLegacy(e.FrozenStake)}
	default:
		return nil, fmt.Errorf("no translation for legacy event %s", name)
	}

	ev, ok := current.Events[name]
	if !ok {
		return nil, fmt.Errorf("event %s missing from current ABI", name)
	}
	if len(ev.Inputs.NonIndexed()) != len(ev.Inputs) {
		return nil, fmt.Errorf("event %s has indexed inputs, translation unsupported", name)
	}
	data, err := ev.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}

	converted := *log
	converted.Topics = append(converted.Topics[:0:0], log.Topics...)
	converted.Topics[0] = ev.ID
	converted.Data = data
	return &converted, nil
}

func stakeFromLegacy(s delegationv7.IDecimalDelegationCommonStake) delegation.IDecimalDelegationCommonStake {
	return delegation.IDecimalDelegationCommonStake{
		Validator:     s.Validator,
		Delegator:     s.Delegator,
		Token:         s.Token,
		Amount:        s.Amount,
		TokenId:       s.TokenId,
		TokenType:     s.TokenType,
		HoldTimestamp: s.HoldTimestamp,
		HoldStartTime: new(big.Int),
	}
}

func frozenStakeFromLegacy(f delegationv7.IDecimalDelegationCommonFrozenStake) delegation.IDecimalDelegationCommonFrozenStake {
	return delegation.IDecimalDelegationCommonFrozenStake{
		Stake:             stakeFromLegacy(f.Stake),
		FreezeStatus:      f.FreezeStatus,
		FreezeType:        f.FreezeType,
		UnfreezeTimestamp: f.UnfreezeTimestamp,
	}
}

// rejectLegacyNftDelegation reverts a tx that stakes into the retired legacy
// DecimalDelegationNFT contract. The node does not track that contract, so a deposit
// there would lock the NFT with no stake, power or rewards on the node. Exits
// (withdraw/transfer requests and completions) stay allowed so NFTs already locked
// there can be recovered.
func (k Keeper) rejectLegacyNftDelegation(ctx sdk.Context, logs []*ethtypes.Log) error {
	addressDelegationNft, err := contracts.GetAddressFromContractCenter(ctx, k.evmKeeper, contracts.NameOfSlugForGetAddressDelegationNft)
	if err != nil || addressDelegationNft == "" {
		return nil
	}
	nftCenter, _ := delegationNft.DelegationNftMetaData.GetAbi()
	if isLegacyNftDeposit(logs, strings.ToLower(addressDelegationNft), nftCenter) {
		return errors.ValidatorNftDelegationInactive
	}
	return nil
}

// isLegacyNftDeposit reports whether logs stake (deposit) into the legacy NFT delegation
// contract: a StakeUpdated/StakeHolded of it with no exit event of it in the same tx. Its
// exits emit StakeUpdated too (the stake decrease on a withdraw/transfer request, the
// re-stake when a transfer completes), and those must pass.
func isLegacyNftDeposit(logs []*ethtypes.Log, addressDelegationNft string, nftCenter *abi.ABI) bool {
	staked := false
	for _, log := range logs {
		if log == nil || len(log.Topics) == 0 || strings.ToLower(log.Address.String()) != addressDelegationNft {
			continue
		}
		ev, err := nftCenter.EventByID(log.Topics[0])
		if err != nil {
			continue
		}
		switch ev.Name {
		case "StakeUpdated", "StakeHolded":
			staked = true
		case "WithdrawRequest", "TransferRequest", "WithdrawCompleted", "TransferCompleted":
			return false
		}
	}
	return staked
}
