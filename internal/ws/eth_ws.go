// Package ws provides WebSocket client utilities for Ethereum connections
package ws

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

// NewClient creates a new Ethereum client connection using the provided URL.
// The URL can be HTTP, HTTPS, WS, or WSS endpoint.
// Returns an error if the connection cannot be established.
func NewClient(url string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	return client, nil
}
