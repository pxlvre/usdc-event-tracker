package tracker

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"usdc-event-tracker/internal/config"
	"usdc-event-tracker/internal/sinks"
	"usdc-event-tracker/internal/sinks/console"
	"usdc-event-tracker/internal/sinks/fs"
	"usdc-event-tracker/internal/tx"
	"usdc-event-tracker/internal/usdc"
)

// Tracker monitors blockchain for USDC events
type Tracker struct {
	client        *ethclient.Client
	config        *config.Config
	blockInterval time.Duration
	sinkManager   *sinks.Manager
}

// New creates a new Tracker instance
func New(client *ethclient.Client, cfg *config.Config) *Tracker {
	t := &Tracker{
		client:        client,
		config:        cfg,
		blockInterval: cfg.BlockInterval,
		sinkManager:   sinks.NewManager(),
	}
	
	// Initialize sinks based on configuration
	t.initializeSinks(cfg)
	
	return t
}

// initializeSinks sets up configured sinks
func (t *Tracker) initializeSinks(cfg *config.Config) {
	for _, sinkName := range cfg.Sink {
		switch sinkName {
		case "console":
			t.sinkManager.AddSink(console.New(cfg.USDCAddress))
		case "sql":
			// TODO: Add SQL sink implementation
			log.Printf("SQL sink not yet implemented")
		case "mongodb":
			// TODO: Add MongoDB sink implementation
			log.Printf("MongoDB sink not yet implemented")
		case "kafka":
			// TODO: Add Kafka sink implementation
			log.Printf("Kafka sink not yet implemented")
		case "filesystem":
			// Configure filesystem sink from environment
			fsConfig := fs.Config{
				OutputDir:  os.Getenv("FS_OUTPUT_DIR"),
				FilePrefix: os.Getenv("FS_FILE_PREFIX"),
				MaxEvents:  0, // Could be made configurable
			}
			
			// Set format from environment
			format := os.Getenv("FS_FORMAT")
			switch format {
			case "csv":
				fsConfig.Format = fs.FormatCSV
			case "text":
				fsConfig.Format = fs.FormatText
			case "jsonl":
				fsConfig.Format = fs.FormatJSONL
			default:
				fsConfig.Format = fs.FormatJSON
			}
			
			t.sinkManager.AddSink(fs.New(fsConfig))
		}
	}
}

// Start begins tracking blockchain events
func (t *Tracker) Start(ctx context.Context) error {
	if err := t.printConnectionInfo(ctx); err != nil {
		return fmt.Errorf("failed to get connection info: %w", err)
	}
	
	// Initialize all sinks
	if err := t.sinkManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize sinks: %w", err)
	}
	
	// Print active sinks
	t.printActiveSinks()
	
	defer t.sinkManager.Close()
	
	return t.monitorBlocks(ctx)
}

// printConnectionInfo displays network connection details
func (t *Tracker) printConnectionInfo(ctx context.Context) error {
	chainID, err := t.client.NetworkID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get network ID: %w", err)
	}

	fmt.Printf("🔗 Connected to Ethereum network (%s)\n", t.config.Network)
	fmt.Printf("   Chain ID: %s\n", chainID.String())
	fmt.Printf("   USDC Address: %s\n", t.config.USDCAddress)
	fmt.Printf("   Block Interval: %v\n", t.blockInterval)

	return nil
}

// printActiveSinks displays configured sinks
func (t *Tracker) printActiveSinks() {
	fmt.Printf("📤 Active sinks: ")
	for i, sink := range t.config.Sink {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s", sink)
	}
	fmt.Printf("\n\n")
}

// monitorBlocks continuously monitors new blocks
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

// processCurrentBlock processes the latest block for USDC events
func (t *Tracker) processCurrentBlock(ctx context.Context) error {
	blockNumber, err := t.client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get block number: %w", err)
	}

	// Only log to console if console sink is active
	if t.sinkManager.HasSink("console") {
		fmt.Printf("📦 Processing block #%d\n", blockNumber)
	}

	receipts, err := tx.GetAllTransactionInBlock(t.client, ctx, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to get receipts for block %d: %w", blockNumber, err)
	}

	if len(receipts) == 0 {
		if t.sinkManager.HasSink("console") {
			fmt.Printf("   No transactions in this block\n\n")
		}
		return nil
	}

	if t.sinkManager.HasSink("console") {
		fmt.Printf("   Total transactions: %d\n", len(receipts))
	}

	// Filter for USDC transactions
	usdcTxs := usdc.MapUSDCTxs(receipts, t.config.USDCAddress)
	
	// Convert to sink events
	events := t.convertToEvents(usdcTxs, blockNumber)
	
	// Send to all configured sinks
	if err := t.sinkManager.Write(ctx, events); err != nil {
		return fmt.Errorf("failed to write to sinks: %w", err)
	}

	return nil
}

// convertToEvents converts receipts to sink events
func (t *Tracker) convertToEvents(receipts []*types.Receipt, blockNumber uint64) []sinks.Event {
	events := make([]sinks.Event, 0, len(receipts))
	
	for _, receipt := range receipts {
		// Filter logs for USDC address only
		usdcLogs := make([]*types.Log, 0)
		for _, log := range receipt.Logs {
			if log.Address.Hex() == t.config.USDCAddress {
				usdcLogs = append(usdcLogs, log)
			}
		}
		
		events = append(events, sinks.Event{
			BlockNumber: blockNumber,
			Receipt:     receipt,
			Logs:        usdcLogs,
		})
	}
	
	return events
}