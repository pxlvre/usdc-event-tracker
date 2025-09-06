package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"usdc-event-tracker/internal/config"
	"usdc-event-tracker/internal/tracker"
	"usdc-event-tracker/internal/ws"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create Ethereum client
	client, err := ws.NewClient(cfg.WebhookURL)
	if err != nil {
		log.Fatalf("Failed to create Ethereum client: %v", err)
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
		log.Println("\n📛 Shutting down gracefully...")
		cancel()
	}()

	// Start tracking
	log.Println("🚀 Starting USDC Event Tracker...")
	if err := t.Start(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Fatalf("Tracker error: %v", err)
		}
	}

	log.Println("✅ Tracker stopped")
}
