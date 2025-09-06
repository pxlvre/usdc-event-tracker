package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	WebhookURL      string
	BlockInterval   time.Duration
	USDCAddress     string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("WEBHOOK_URL environment variable is required")
	}

	return &Config{
		WebhookURL:    webhookURL,
		BlockInterval: 12 * time.Second, // Ethereum block time
		USDCAddress:   "0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238", // Sepolia USDC
	}
}