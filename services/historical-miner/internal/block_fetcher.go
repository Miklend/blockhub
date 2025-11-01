package worker

import (
	"context"
	"lib/blocks/collector"
	"time"

	"lib/models"
	"lib/utils/logging"
)

// структура фетчера
type BlockFetcher struct {
	logger          *logging.Logger
	collector       *collector.BlockCollector
	jobsChan        <-chan uint64
	receiptsJobChan chan<- *models.Block
}

// Конструктор фетчера
func NewBlockFetcher(
	logger *logging.Logger,
	collector *collector.BlockCollector,
	jobsChan <-chan uint64,
	receiptsJobChan chan<- *models.Block,
) *BlockFetcher {
	return &BlockFetcher{
		logger:          logger,
		collector:       collector,
		jobsChan:        jobsChan,
		receiptsJobChan: receiptsJobChan,
	}
}

func (f *BlockFetcher) ProcessBlocks(ctx context.Context) {
	maxRetries := 5
	for {
		select {
		case blockNumber, ok := <-f.jobsChan:
			if !ok {
				f.logger.Warn("Fetcher blocks channel closed. Stopping")
				return
			}

			var block []models.Block
			var blockErr error

			for attempt := 1; attempt < maxRetries; attempt++ {
				block, blockErr = f.collector.FetchBlocksBatch(ctx, []uint64{blockNumber})

				if blockErr == nil && len(block) > 0 {
					break
				}
				f.logger.Warnf("Block %d metadata is not available, attempt %d/%d: %v", blockNumber, attempt, maxRetries, blockErr)

				if attempt <= maxRetries {
					select {
					case <-ctx.Done():
						return
					case <-time.After(time.Duration(attempt) * 500 * time.Millisecond):
					}
				}
			}
			if blockErr != nil && len(block) > 0 {
				f.logger.Errorf("All attempts failed to fetch metadata for block %d", blockNumber)
				continue
			}

			blockData := &block[0]
			f.logger.Infof("Block %d fetched correctly,num of transactions:%d sending to receipts queue", blockData.Number, len(blockData.Transactions))

			select {
			case f.receiptsJobChan <- blockData:
			case <-ctx.Done():
				f.logger.Info("Block Fetcher received shutdown signal, stopping sending results")
				return
			}
		case <-ctx.Done():
			f.logger.Info("Block Fetcher received shutdown signal, stopping")
			return
		}

	}
}
