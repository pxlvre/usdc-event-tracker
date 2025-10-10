// Package erc20 provides ERC20 token event definitions and utilities
package erc20

// Event represents an ERC20 event type
type Event string

const (
	// Standard ERC20 events
	Transfer Event = "Transfer"
	Approval Event = "Approval"
)

// EventSignatures maps ERC20 events to their keccak256 signature hashes.
// These signatures are used to identify events in transaction logs.
var EventSignatures = map[Event]string{
	Transfer: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", // Transfer(address,address,uint256)
	Approval: "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", // Approval(address,address,uint256)
}

// GetEventBySignature looks up an event type by its signature hash.
// Returns the Event and true if found, or empty string and false if not found.
func GetEventBySignature(signature string) (Event, bool) {
	for event, sig := range EventSignatures {
		if sig == signature {
			return event, true
		}
	}
	return "", false
}