package tracker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"usdc-event-tracker/internal/config"
	"usdc-event-tracker/internal/tx"
	"usdc-event-tracker/internal/usdc"
)

type Tracker struct {
	client        *ethclient.Client
	config        *config.Config
	blockInterval time.Duration
}

func New(client *ethclient.Client, cfg *config.Config) *Tracker {
	return &Tracker{
		client:        client,
		config:        cfg,
		blockInterval: cfg.BlockInterval,
	}
}

func (t *Tracker) Start(ctx context.Context) error {
	if err := t.printConnectionInfo(ctx); err != nil {
		return fmt.Errorf("failed to get connection info: %w", err)
	}

	return t.monitorBlocks(ctx)
}

func (t *Tracker) printConnectionInfo(ctx context.Context) error {
	chainID, err := t.client.NetworkID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get network ID: %w", err)
	}

	fmt.Printf("🔗 Connected to Ethereum network (%s)\n", t.config.Network)
	fmt.Printf("   Chain ID: %s\n", chainID.String())
	fmt.Printf("   USDC Address: %s\n", t.config.USDCAddress)
	fmt.Printf("   Block Interval: %v\n\n", t.blockInterval)

	return nil
}

func (t *Tracker) monitorBlocks(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := t.processCurrentBlock(ctx); err != nil {
				log.Printf("Error processing block: %v", err)
			}
			time.Sleep(t.blockInterval)
		}
	}
}

func (t *Tracker) processCurrentBlock(ctx context.Context) error {
	blockNumber, err := t.client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get block number: %w", err)
	}

	fmt.Printf("📦 Processing block #%d\n", blockNumber)

	receipts, err := tx.GetAllTransactionInBlock(t.client, ctx, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to get receipts for block %d: %w", blockNumber, err)
	}

	if len(receipts) == 0 {
		fmt.Printf("   No transactions in this block\n\n")
		return nil
	}

	fmt.Printf("   Total transactions: %d\n", len(receipts))

	usdcTxs := usdc.MapUSDCTxs(receipts, t.config.USDCAddress)
	t.displayUSDCTransactions(usdcTxs, blockNumber)

	return nil
}

func (t *Tracker) displayUSDCTransactions(txs []*types.Receipt, blockNumber uint64) {
	if len(txs) == 0 {
		fmt.Printf("   No USDC transactions found\n\n")
		return
	}

	fmt.Printf("\n   💰 USDC Transactions (%d found):\n", len(txs))
	for i, receipt := range txs {
		t.displayTransaction(i+1, receipt)
	}
	fmt.Println()
}

func (t *Tracker) displayTransaction(index int, receipt *types.Receipt) {
	fmt.Printf("   [%d] Transaction Details:\n", index)
	fmt.Printf("       Hash: %s\n", receipt.TxHash.Hex())
	fmt.Printf("       Status: %s\n", getStatusText(receipt.Status))
	fmt.Printf("       Gas Used: %d\n", receipt.GasUsed)
	
	// Display USDC-specific events
	for _, log := range receipt.Logs {
		if log.Address.Hex() == t.config.USDCAddress {
			displayUSDCEvent(log)
		}
	}
}

func getStatusText(status uint64) string {
	if status == 1 {
		return "✅ Success"
	}
	return "❌ Failed"
}

func displayUSDCEvent(log *types.Log) {
	if len(log.Topics) == 0 {
		return
	}

	// Transfer event signature: 0xddf252ad...
	transferSig := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	approvalSig := "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"

	switch log.Topics[0].Hex() {
	case transferSig:
		fmt.Printf("       Event: Transfer\n")
		if len(log.Topics) >= 3 {
			fmt.Printf("         From: %s\n", log.Topics[1].Hex())
			fmt.Printf("         To: %s\n", log.Topics[2].Hex())
		}
	case approvalSig:
		fmt.Printf("       Event: Approval\n")
		if len(log.Topics) >= 3 {
			fmt.Printf("         Owner: %s\n", log.Topics[1].Hex())
			fmt.Printf("         Spender: %s\n", log.Topics[2].Hex())
		}
	default:
		fmt.Printf("       Event: Unknown (Topic: %s)\n", log.Topics[0].Hex()[:10]+"...")
	}
}