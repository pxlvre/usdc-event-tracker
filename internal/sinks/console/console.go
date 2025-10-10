// Package console implements a sink that outputs events to the console
package console

import (
	"context"
	"fmt"
	
	"github.com/ethereum/go-ethereum/core/types"
	"usdc-event-tracker/internal/erc20"
	"usdc-event-tracker/internal/sinks"
)

// ConsoleSink implements the Sink interface for console output.
// It formats and displays USDC events in a human-readable format.
type ConsoleSink struct {
	usdcAddress string
}

// New creates a new console sink configured for the specified USDC address.
func New(usdcAddress string) *ConsoleSink {
	return &ConsoleSink{
		usdcAddress: usdcAddress,
	}
}

// Name returns "console" as the sink identifier.
func (c *ConsoleSink) Name() string {
	return "console"
}

// Initialize prepares the console sink for use.
// For console output, this simply prints an initialization message.
func (c *ConsoleSink) Initialize() error {
	fmt.Println("📊 Console sink initialized")
	return nil
}

// Write formats and displays events to the console.
// Each event is displayed with transaction details and USDC-specific information.
func (c *ConsoleSink) Write(ctx context.Context, events []sinks.Event) error {
	if len(events) == 0 {
		fmt.Printf("   No USDC transactions found\n\n")
		return nil
	}

	fmt.Printf("\n   💰 USDC Transactions (%d found):\n", len(events))
	for i, event := range events {
		c.displayTransaction(i+1, event)
	}
	fmt.Println()
	
	return nil
}

// Close performs cleanup. This is a no-op for console output.
func (c *ConsoleSink) Close() error {
	return nil
}

// displayTransaction formats and displays a single transaction
func (c *ConsoleSink) displayTransaction(index int, event sinks.Event) {
	fmt.Printf("   [%d] Transaction Details:\n", index)
	fmt.Printf("       Block: #%d\n", event.BlockNumber)
	fmt.Printf("       Hash: %s\n", event.Receipt.TxHash.Hex())
	fmt.Printf("       Status: %s\n", c.getStatusText(event.Receipt.Status))
	fmt.Printf("       Gas Used: %d\n", event.Receipt.GasUsed)
	
	// Display USDC-specific events
	for _, log := range event.Logs {
		c.displayUSDCEvent(log)
	}
}

// getStatusText returns a formatted status string
func (c *ConsoleSink) getStatusText(status uint64) string {
	if status == 1 {
		return "✅ Success"
	}
	return "❌ Failed"
}

// displayUSDCEvent formats and displays a USDC event
func (c *ConsoleSink) displayUSDCEvent(log *types.Log) {
	if len(log.Topics) == 0 {
		return
	}

	event, found := erc20.GetEventBySignature(log.Topics[0].Hex())
	if !found {
		fmt.Printf("       Event: Unknown (Topic: %s)\n", log.Topics[0].Hex()[:10]+"...")
		return
	}

	fmt.Printf("       Event: %s\n", event)
	
	switch event {
	case erc20.Transfer:
		if len(log.Topics) >= 3 {
			fmt.Printf("         From: %s\n", log.Topics[1].Hex())
			fmt.Printf("         To: %s\n", log.Topics[2].Hex())
		}
	case erc20.Approval:
		if len(log.Topics) >= 3 {
			fmt.Printf("         Owner: %s\n", log.Topics[1].Hex())
			fmt.Printf("         Spender: %s\n", log.Topics[2].Hex())
		}
	}
}