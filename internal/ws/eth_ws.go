package ws

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

// NewClient creates a new Ethereum client connection
func NewClient(url string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	return client, nil
}
