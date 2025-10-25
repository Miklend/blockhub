package worker

import (
	"context"
	"lib/blocks/collector"
	"lib/models"
	"lib/utils/logging"
	"time"
)

type BlockProcessor struct {
	logger      *logging.Logger
	collector   collector.BlockCollector
	jobsChan    <-chan uint64        //канал получающий номера блоков
	resultsChan chan<- *models.Block //канал отправляющий полностью загруженные блоки
}

func NewBlockProcessor(
	logger *logging.Logger,
	collector collector.BlockCollector,
	jobs <-chan uint64,
	results chan<- *models.Block,
) *BlockProcessor {
	return &BlockProcessor{
		logger:      logger,
		collector:   collector,
		jobsChan:    jobs,
		resultsChan: results,
	}
}

// основная логика воркера
func (p *BlockProcessor) ProcessBlocks(ctx context.Context) {
	maxRetries := 5
	for {
		select {
		// Ожидание заполнения переменной номером блока из очереди
		case blockNumber, ok := <-p.jobsChan:
			if !ok {
				p.logger.Debug("Block processor: Jobs channel closed, stopping")
				return
			}

			var blockData *models.Block
			var blockErr error
			// Если все в порядке то считываются данные блока
			for attempt := 1; attempt <= maxRetries; attempt++ {
				blockData, blockErr = p.collector.CollectBlockByNumber(ctx, blockNumber)
				if blockErr == nil {
					break
				}

				p.logger.Debugf("Block %d is not availible yet (attempt %d/%d), retrying", blockNumber, attempt, maxRetries)

				if attempt <= maxRetries {
					select {
					case <-ctx.Done():
						return
					case <-time.After(time.Duration(attempt) * 500 * time.Millisecond):

					}
				}
			}
			if blockErr != nil {
				p.logger.Errorf("All %d attempts failed for block %d. Stopping work", maxRetries, blockNumber)
				continue
			}
			p.logger.Debugf("Processing block %d finished. Sending to result channel", blockData.Number)
			select {
			case p.resultsChan <- blockData:
			case <-ctx.Done():
				p.logger.Info("Block Processor recieved shutdown signal, stopping sending results")
				return
			}
		case <-ctx.Done():
			p.logger.Info("Block Processor recieved shutdown signal, stopping")
			return
		}

	}
}
