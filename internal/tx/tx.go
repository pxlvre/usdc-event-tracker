// Package tx provides transaction and receipt retrieval utilities
package tx

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// GetAllTransactionInBlock retrieves all transaction receipts for a given block number.
// It uses the BlockReceipts method for efficient batch retrieval.
// Returns an empty slice if the block contains no transactions.
func GetAllTransactionInBlock(client *ethclient.Client, ctx context.Context, blockNumber uint64) ([]*types.Receipt, error) {
	blockNum := rpc.BlockNumber(blockNumber)
	
	receipts, err := client.BlockReceipts(ctx, rpc.BlockNumberOrHashWithNumber(blockNum))
	if err != nil {
		return nil, fmt.Errorf("failed to get receipts for block %d: %w", blockNumber, err)
	}
	
	return receipts, nil
}
