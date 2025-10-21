package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"usdc-event-tracker/internal/config"
	"usdc-event-tracker/internal/logging"
	"usdc-event-tracker/internal/tracker"
	"usdc-event-tracker/internal/ws"
)

func main() {
	// Initialize structured logging
	logging.Init("main")
	logger := logging.GetLogger("main")

	// Load configuration
	cfg := config.Load()

	logger.Info("Starting USDC Event Tracker", map[string]interface{}{
		"network":     cfg.Network,
		"sinks":       cfg.Sink,
		"usdc_address": cfg.USDCAddress,
	})

	// Create Ethereum client
	client, err := ws.NewClient(cfg.WebhookURL)
	if err != nil {
		logger.Error("Failed to create Ethereum client", err)
		os.Exit(1)
	}
	defer client.Close()

	// Create and start tracker
	t := tracker.New(client, cfg)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Shutdown signal received, stopping gracefully")
		cancel()
	}()

	// Start tracking
	if err := t.Start(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			logger.Error("Tracker error", err)
			os.Exit(1)
		}
	}

	logger.Info("Tracker stopped successfully")
}
