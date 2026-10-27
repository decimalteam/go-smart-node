package keeper

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// validateCriticalEventEnvelope validates an event from a trusted financial
// contract before zero values or stale decoded structs can affect node state.
func validateCriticalEventEnvelope(a *abi.ABI, entry *ethtypes.Log) error {
	if len(entry.Topics) == 0 {
		return fmt.Errorf("financial contract emitted an event without a signature")
	}
	event, err := a.EventByID(entry.Topics[0])
	if err != nil {
		return fmt.Errorf("financial contract emitted unknown event %s: %w", entry.Topics[0], err)
	}
	// Do not skip empty data: events with mandatory fields must reject it.
	if _, err = event.Inputs.NonIndexed().Unpack(entry.Data); err != nil {
		return fmt.Errorf("decode financial event %s: %w", event.Name, err)
	}
	var indexed abi.Arguments
	for _, arg := range event.Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err = abi.ParseTopicsIntoMap(map[string]interface{}{}, indexed, entry.Topics[1:]); err != nil {
		return fmt.Errorf("decode financial event topics %s: %w", event.Name, err)
	}
	return nil
}
