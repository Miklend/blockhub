package collector

import (
	"context"
	"fmt"
	"lib/models"
	"time"
)

func (bc *BlockCollector) CollectBlockByNumber(ctx context.Context, blockNumber uint64) (*models.BlockDTO, error) {
	bc.logger.Debugf("Starting collection for block #%d", blockNumber)
	startTime := time.Now()
	hexNumber := fmt.Sprintf("0x%x", blockNumber)
	// Fetch block data
	blockFetchStart := time.Now()
	block, err := bc.Client().BlockByNumber(ctx, hexNumber)
	if err != nil {
		bc.logger.Errorf("Failed to fetch block %d from network: %v", blockNumber, err)
		return nil, fmt.Errorf("failed to fetch block %d: %w", blockNumber, err)
	}
	blockFetchTime := time.Since(blockFetchStart)

	// Fetch receipts
	receiptsFetchStart := time.Now()
	receipts, err := bc.Client().ReceiptByBlockNumber(ctx, hexNumber)
	if err != nil {
		bc.logger.Errorf("Failed to fetch receipts for block %d (hash: %s): %v",
			blockNumber, block.Hash, err)
		return nil, fmt.Errorf("failed to fetch receipts for block %d: %w", blockNumber, err)
	}
	receiptsFetchTime := time.Since(receiptsFetchStart)

	bc.logger.Debugf("Receipts for block %d fetched successfully in %v: %d receipts",
		blockNumber, receiptsFetchTime, len(receipts))

	for i := range block.Transactions {
		tx := &block.Transactions[i]
		for _, r := range receipts {
			if tx.Hash == r.TransactionHash {
				tx.Receipt = r // Присваиваем ReceiptDTO
				break
			}
		}
	}

	totalTime := time.Since(startTime)

	bc.logger.Infof("Block %d collection completed in %v (block: %v, receipts: %v, metrics: %v) - %d tx, %d gas, miner: %s",
		blockNumber, totalTime, blockFetchTime, receiptsFetchTime, receiptsFetchTime,
		len(block.Transactions), block.GasUsed, block.Miner)

	return &block, nil
}
