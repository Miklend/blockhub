package collector

import (
	"context"
	"fmt"
	"lib/blocks/metrics"
	"lib/models"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

func (bc *BlockCollector) CollectBlockByNumber(ctx context.Context, blockNumber uint64) (*models.Block, error) {
	bc.logger.Debugf("Starting collection for block #%d", blockNumber)
	startTime := time.Now()

	// Fetch block data
	blockFetchStart := time.Now()
	block, err := bc.Client().BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		bc.logger.Errorf("Failed to fetch block %d from network: %v", blockNumber, err)
		return nil, fmt.Errorf("failed to fetch block %d: %w", blockNumber, err)
	}
	blockFetchTime := time.Since(blockFetchStart)

	bc.logger.Debugf("Block %d fetched successfully in %v: %d transactions, %d gas used",
		blockNumber, blockFetchTime, len(block.Transactions()), block.GasUsed())

	// Fetch receipts
	receiptsFetchStart := time.Now()
	blockNrOrHash := rpc.BlockNumberOrHashWithHash(block.Hash(), false)
	receipts, err := bc.Client().BlockReceipts(ctx, blockNrOrHash)
	if err != nil {
		bc.logger.Errorf("Failed to fetch receipts for block %d (hash: %s): %v",
			blockNumber, block.Hash().Hex(), err)
		return nil, fmt.Errorf("failed to fetch receipts for block %d: %w", blockNumber, err)
	}
	receiptsFetchTime := time.Since(receiptsFetchStart)

	bc.logger.Debugf("Receipts for block %d fetched successfully in %v: %d receipts",
		blockNumber, receiptsFetchTime, len(receipts))

	// Calculate metrics
	metricsCalcStart := time.Now()
	metrics := metrics.NewBlockMetrics(block, receipts)
	metricsCalcTime := time.Since(metricsCalcStart)

	totalTime := time.Since(startTime)

	bc.logger.Infof("Block %d collection completed in %v (block: %v, receipts: %v, metrics: %v) - %d tx, %d gas, miner: %s",
		blockNumber, totalTime, blockFetchTime, receiptsFetchTime, metricsCalcTime,
		len(block.Transactions()), block.GasUsed(), block.Coinbase().Hex())

	return &metrics, nil
}
