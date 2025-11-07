package worker

import (
	"context"
	"lib/blocks/collector"
	"time"

	"lib/models"
	"lib/utils/logging"
)

type BlockFetcher struct {
	logger          *logging.Logger
	collector       *collector.BlockCollector
	jobsChan        <-chan uint64
	receiptsJobChan chan<- *models.BlockDTO
}

func NewBlockFetcher(
	logger *logging.Logger,
	collector *collector.BlockCollector,
	jobsChan <-chan uint64,
	receiptsJobChan chan<- *models.BlockDTO,
) *BlockFetcher {
	return &BlockFetcher{
		logger:          logger,
		collector:       collector,
		jobsChan:        jobsChan,
		receiptsJobChan: receiptsJobChan,
	}
}

func (f *BlockFetcher) ProcessBlocks(ctx context.Context) {
	maxRetries := 10
	for {
		select {
		case blockNumber, ok := <-f.jobsChan:
			if !ok {
				f.logger.Warn("Fetcher blocks channel closed. Stopping")
				return
			}

			var blocksMap map[models.Hash]models.BlockDTO
			var blockErr error

			for attempt := 1; attempt <= maxRetries; attempt++ {
				blocksMap, blockErr = f.collector.FetchBlocksBatch(ctx, []uint64{blockNumber})

				if blockErr == nil && len(blocksMap) == 1 {
					break
				}

				f.logger.Warnf("Block %d data is not available (attempt %d/%d). Error: %v",
					blockNumber, attempt, maxRetries, blockErr)

				if attempt < maxRetries {
					select {
					case <-ctx.Done():
						return
					case <-time.After(time.Duration(attempt) * 500 * time.Millisecond):
					}
				}
			}

			if blockErr != nil || len(blocksMap) != 1 {
				f.logger.Errorf("All attempts failed to fetch block data (or map size incorrect) for block %d. Skipping.", blockNumber)
				continue
			}
			var blockData *models.BlockDTO
			for _, block := range blocksMap {
				blockData = &block
				break
			}

			if blockData == nil {
				f.logger.Errorf("Fatal error: Extracted nil block data from map for block %d", blockNumber)
				continue
			}

			f.logger.Infof("Block %d fetched correctly, num of transactions:%d sending to receipts queue",
				blockData.Number, len(blockData.Transactions))

			select {
			case f.receiptsJobChan <- blockData:
			case <-ctx.Done():
				f.logger.Info("Block Fetcher received shutdown signal, stopping sending results")
				return
			default:
				f.logger.Warnf("Receipts job channel full, dropping block %d", blockData.Number)
			}
		case <-ctx.Done():
			f.logger.Info("Block Fetcher received shutdown signal, stopping")
			return
		}
	}
}
