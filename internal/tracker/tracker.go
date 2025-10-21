package tracker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"usdc-event-tracker/internal/config"
	"usdc-event-tracker/internal/logging"
	"usdc-event-tracker/internal/sinks"
	"usdc-event-tracker/internal/sinks/console"
	"usdc-event-tracker/internal/sinks/elasticsearch"
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
	logger        *logging.Logger
}

// New creates a new Tracker instance
func New(client *ethclient.Client, cfg *config.Config) *Tracker {
	t := &Tracker{
		client:        client,
		config:        cfg,
		blockInterval: cfg.BlockInterval,
		sinkManager:   sinks.NewManager(),
		logger:        logging.GetLogger("tracker"),
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
			t.logger.Warn("SQL sink not yet implemented", map[string]interface{}{"sink": "sql"})
		case "mongodb":
			// TODO: Add MongoDB sink implementation
			t.logger.Warn("MongoDB sink not yet implemented", map[string]interface{}{"sink": "mongodb"})
		case "kafka":
			// TODO: Add Kafka sink implementation
			t.logger.Warn("Kafka sink not yet implemented", map[string]interface{}{"sink": "kafka"})
		case "elasticsearch":
			esConfig := elasticsearch.NewConfig()
			t.sinkManager.AddSink(elasticsearch.New(esConfig))
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
		t.logger.Error("Failed to get network ID", err)
		return fmt.Errorf("failed to get network ID: %w", err)
	}

	t.logger.LogConnection(t.config.Network, chainID.String(), t.config.USDCAddress, t.config.WebhookURL)

	return nil
}

// printActiveSinks displays configured sinks
func (t *Tracker) printActiveSinks() {
	t.logger.Info("Active sinks initialized", map[string]interface{}{
		"sink_count": len(t.config.Sink),
		"sinks":      t.config.Sink,
	})
}

// monitorBlocks continuously monitors new blocks
func (t *Tracker) monitorBlocks(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := t.processCurrentBlock(ctx); err != nil {
				t.logger.Error("Error processing block", err)
			}
			time.Sleep(t.blockInterval)
		}
	}
}

// processCurrentBlock processes the latest block for USDC events
func (t *Tracker) processCurrentBlock(ctx context.Context) error {
	blockNumber, err := t.client.BlockNumber(ctx)
	if err != nil {
		t.logger.Error("Failed to get block number", err)
		return fmt.Errorf("failed to get block number: %w", err)
	}

	receipts, err := tx.GetAllTransactionInBlock(t.client, ctx, blockNumber)
	if err != nil {
		t.logger.Error("Failed to get receipts for block", err, map[string]interface{}{
			"block_number": blockNumber,
		})
		return fmt.Errorf("failed to get receipts for block %d: %w", blockNumber, err)
	}

	t.logger.LogBlockProcessing(blockNumber, len(receipts))

	if len(receipts) == 0 {
		return nil
	}

	// Filter for USDC transactions
	usdcTxs := usdc.MapUSDCTxs(receipts, t.config.USDCAddress)
	
	// Log USDC transactions found
	if len(usdcTxs) > 0 {
		t.logger.Info("USDC transactions found", map[string]interface{}{
			"block_number":    blockNumber,
			"total_txs":       len(receipts),
			"usdc_txs":        len(usdcTxs),
			"usdc_percentage": float64(len(usdcTxs)) / float64(len(receipts)) * 100,
		})
	}
	
	// Convert to sink events
	events := t.convertToEvents(usdcTxs, blockNumber)
	
	// Send to all configured sinks
	start := time.Now()
	if err := t.sinkManager.Write(ctx, events); err != nil {
		t.logger.Error("Failed to write to sinks", err, map[string]interface{}{
			"block_number": blockNumber,
			"event_count":  len(events),
		})
		return fmt.Errorf("failed to write to sinks: %w", err)
	}
	
	t.logger.Info("Sink write completed", map[string]interface{}{
		"block_number": blockNumber,
		"event_count":  len(events),
		"duration_ms":  time.Since(start).Milliseconds(),
	})

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