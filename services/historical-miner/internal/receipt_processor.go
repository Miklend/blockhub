package worker

import (
	"context"
	"lib/blocks/collector"
	_ "lib/clients/broker"
	"lib/models"
	"lib/utils/logging"
	"time"
)

const OUTPUT_FILE = "Block_new.json"

type ReceiptProcessor struct {
	logger    *logging.Logger
	collector collector.BlockCollector
	jobsChan  <-chan *models.Block
	//kafkaClient broker.BrokerClient  //канал отправляющий полностью загруженные блоки
}

func NewReceiptProcessor(
	logger *logging.Logger,
	collector collector.BlockCollector,
	jobs <-chan *models.Block,
	//kafkaClient broker.BrokerClient,
) *ReceiptProcessor {
	return &ReceiptProcessor{
		logger:    logger,
		collector: collector,
		jobsChan:  jobs,
		//kafkaClient: kafkaClient,
	}
}

// основная логика воркера
func (p *ReceiptProcessor) ProcessReceipts(ctx context.Context) {
	maxRetries := 10 // Определите число повторных попыток
	for {
		select {
		case blockData, ok := <-p.jobsChan:
			if !ok {
				return
			}

			blockNumber := blockData.Number

			var receiptsMap map[uint64][]models.Receipt
			var receiptsErr error
			for attempt := 1; attempt <= maxRetries; attempt++ {
				receiptsMap, receiptsErr = p.collector.FetchReceiptsBatch(ctx, []uint64{blockNumber})

				if receiptsErr == nil && len(receiptsMap) > 0 {
					break
				}

				p.logger.Warnf("Failed to fetch receipts for block %d (attempt %d/%d): %v",
					blockNumber, attempt, maxRetries, receiptsErr)

				if attempt < maxRetries {
					select {
					case <-ctx.Done():
						return
					case <-time.After(time.Duration(attempt) * 1 * time.Second):
					}
				}
			}

			if receiptsErr != nil || len(receiptsMap) == 0 {
				p.logger.Errorf("All %d attempts failed to fetch receipts for block %d. Skipping block.",
					maxRetries, blockNumber)
				continue
			}

			blockReceipts, ok := receiptsMap[blockNumber]
			if !ok {
				p.logger.Warnf("Block %d: Receipts map key not found. Skipping.", blockNumber)
				continue
			}

			if len(blockReceipts) != len(blockData.Transactions) {
				p.logger.Warnf("Block %d: number of fetched receipts (%d) does not match transactions (%d). Skipping processing.",
					blockNumber, len(blockReceipts), len(blockData.Transactions))
				continue
			}

			for j := range blockData.Transactions {
				blockData.Transactions[j].Receipt = &blockReceipts[j]
			}
			p.logger.Infof("SUCCESSFULLY PROCESSED: Block %d is complete. Total Transactions: %d.",
				blockNumber, len(blockData.Transactions))

		case <-ctx.Done():
			p.logger.Info("Receipt Processor received shutdown signal, stopping")
			return
		}

	}
}
