// Package usdc provides USDC-specific utilities and filters
package usdc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"usdc-event-tracker/internal/erc20"
)

// NewAddress converts a hex string to an Ethereum address.
// The hex string can be with or without the 0x prefix.
func NewAddress(hexAddress string) common.Address {
	return common.HexToAddress(hexAddress)
}

// MapUSDCTxs filters a slice of receipts to return only those that interact with the USDC contract.
// It uses the provided USDC contract address to identify relevant transactions.
func MapUSDCTxs(receipts []*types.Receipt, usdcAddress string) []*types.Receipt {
	return FilterByAddress(receipts, usdcAddress)
}

// FilterByAddress filters receipts to return only those containing logs from a specific contract address.
// This is a generic filter that can be used for any contract, not just USDC.
func FilterByAddress(receipts []*types.Receipt, contractAddress string) []*types.Receipt {
	addr := NewAddress(contractAddress)
	filtered := make([]*types.Receipt, 0)

	for _, receipt := range receipts {
		if hasLogsFromAddress(receipt, addr) {
			filtered = append(filtered, receipt)
		}
	}

	return filtered
}

// hasLogsFromAddress checks if a receipt contains any logs from the specified address.
// Returns true if at least one log matches the address, false otherwise.
func hasLogsFromAddress(receipt *types.Receipt, address common.Address) bool {
	for _, log := range receipt.Logs {
		if log.Address == address {
			return true
		}
	}
	return false
}

// GetEventType identifies the ERC20 event type from a log topic hash.
// Returns the event name (e.g., "Transfer", "Approval") or "Unknown" if not recognized.
func GetEventType(topic common.Hash) string {
	event, found := erc20.GetEventBySignature(topic.Hex())
	if found {
		return string(event)
	}
	return "Unknown"
}
