package usdc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Event signatures for USDC contract
const (
	TransferEventSig = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	ApprovalEventSig = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
)

// NewAddress creates a USDC contract address from a hex string
func NewAddress(hexAddress string) common.Address {
	return common.HexToAddress(hexAddress)
}

// MapUSDCTxs filters receipts to find only those that interact with USDC
func MapUSDCTxs(receipts []*types.Receipt, usdcAddress string) []*types.Receipt {
	return FilterByAddress(receipts, usdcAddress)
}

// FilterByAddress filters receipts by a specific contract address
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

// hasLogsFromAddress checks if a receipt contains logs from a specific address
func hasLogsFromAddress(receipt *types.Receipt, address common.Address) bool {
	for _, log := range receipt.Logs {
		if log.Address == address {
			return true
		}
	}
	return false
}

// GetEventType returns a human-readable event type from a log topic
func GetEventType(topic common.Hash) string {
	switch topic.Hex() {
	case TransferEventSig:
		return "Transfer"
	case ApprovalEventSig:
		return "Approval"
	default:
		return "Unknown"
	}
}
