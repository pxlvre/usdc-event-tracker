// Package config provides configuration management for the USDC event tracker
package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	// USDC contract addresses for different networks
	USDCMainnet   = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
	USDCSepolia   = "0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238"
	USDCArbitrum  = "0xaf88d065e77c8cC2239327C5EDb3A432268e5831"
	USDCAvalanche = "0xB97EF9Ef8734C71904D8002F8b6Bc66Dd9c48a6E"
	USDCLinea     = "0x176211869cA2b568f2A7D4EE941E073a821EE1ff"
	USDCPolygon   = "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359"
	USDCOptimism  = "0x0b2c639c533813f4aa9d7837caf62653d097ff85"
)

// Config holds the application configuration
type Config struct {
	WebhookURL    string
	BlockInterval time.Duration
	USDCAddress   string
	Network       string
	Sink          []string
}

// Load reads configuration from environment variables and returns a Config instance.
// It loads from .env file if present, otherwise uses system environment variables.
// Required: WEBHOOK_URL must be set.
// Defaults: NETWORK=sepolia, SINKS=console if not specified.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("WEBHOOK_URL environment variable is required")
	}

	// Get network from environment, default to sepolia
	network := strings.ToLower(os.Getenv("NETWORK"))
	if network == "" {
		network = "sepolia"
	}

	// Select USDC address based on network
	var usdcAddress string
	switch network {
	case "mainnet", "ethereum":
		usdcAddress = USDCMainnet
	case "sepolia":
		usdcAddress = USDCSepolia
	case "arbitrum":
		usdcAddress = USDCArbitrum
	case "avalanche":
		usdcAddress = USDCAvalanche
	case "linea":
		usdcAddress = USDCLinea
	case "polygon":
		usdcAddress = USDCPolygon
	case "optimism":
		usdcAddress = USDCOptimism
	default:
		log.Fatalf("Unsupported network: %s. Supported networks: mainnet, sepolia, arbitrum, avalanche, linea, polygon, optimism", network)
	}

	// Parse sinks from environment (comma-separated)
	// Supported: console, sql, mongodb, kafka, filesystem, elasticsearch
	var sinks []string
	sinksEnv := os.Getenv("SINKS")
	if sinksEnv == "" {
		// Default to console if no sinks specified
		sinks = []string{"console"}
	} else {
		for _, sink := range strings.Split(sinksEnv, ",") {
			trimmed := strings.TrimSpace(strings.ToLower(sink))
			switch trimmed {
			case "console", "sql", "mongodb", "kafka", "filesystem", "elasticsearch":
				sinks = append(sinks, trimmed)
			default:
				log.Printf("Warning: Unsupported sink '%s'. Supported sinks: console, sql, mongodb, kafka, filesystem, elasticsearch", trimmed)
			}
		}
		// If no valid sinks were added, default to console
		if len(sinks) == 0 {
			sinks = []string{"console"}
		}
	}

	return &Config{
		WebhookURL:    webhookURL,
		BlockInterval: 12 * time.Second, // Ethereum block time
		USDCAddress:   usdcAddress,
		Network:       network,
		Sink:          sinks,
	}
}
