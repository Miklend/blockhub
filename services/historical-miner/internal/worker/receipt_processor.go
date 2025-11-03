package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"lib/blocks/collector"
	"lib/clients/broker"
	"lib/models"
	"lib/utils/logging"
	"time"
)

const OUTPUT_FILE = "Block_new.json"
const maxKafkaRetries = 3
const kafkaRetryDelay = 500 * time.Millisecond
const topicKafka = "blocks"

type ReceiptProcessor struct {
	logger      *logging.Logger
	collector   collector.BlockCollector
	jobsChan    <-chan *models.Block
	kafkaClient broker.BrokerClient
}

func NewReceiptProcessor(
	logger *logging.Logger,
	collector collector.BlockCollector,
	jobs <-chan *models.Block,
	kafkaClient broker.BrokerClient,
) *ReceiptProcessor {
	return &ReceiptProcessor{
		logger:      logger,
		collector:   collector,
		jobsChan:    jobs,
		kafkaClient: kafkaClient,
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

			data, err := json.Marshal(blockData)
			if err != nil {
				p.logger.Errorf("failed to serialize block %s: %v", blockData.Hash, err)
				continue
			}
			p.sendWithRetry(ctx, blockNumber, data)
		case <-ctx.Done():
			p.logger.Info("Receipt Processor received shutdown signal, stopping")
			return
		}

	}
}

func (p *ReceiptProcessor) sendWithRetry(ctx context.Context, blockNumber uint64, data []byte) {
	var err error
	key := []byte(fmt.Sprintf("%d", blockNumber))

	m := models.MessageBroker{
		Key:   key,
		Value: data,       // Используем сериализованные данные
		Topic: topicKafka, // Указываем топик
	}
	for attempt := 1; attempt <= maxKafkaRetries; attempt++ {
		// Отправка в Kafka
		err = p.kafkaClient.SendMessage(ctx, m)

		if err == nil {
			p.logger.Infof("Block %s sent to Kafka successfully (attempt %d)", string(m.Key), attempt)
			return
		}

		p.logger.Warnf("Failed to send block %s to Kafka (attempt %d/%d): %v", string(m.Key), attempt, maxKafkaRetries, err)

		if attempt < maxKafkaRetries {
			select {
			case <-ctx.Done():
				p.logger.Warnf("Context cancelled during Kafka retry for block %s", string(m.Key))
				return
			case <-time.After(kafkaRetryDelay):
				// Ждем перед следующей попыткой
			}
		}
	}
	p.logger.Errorf("FATAL: Failed to send block %s to Kafka after %d attempts. DROPPING MESSAGE.", string(m.Key), maxKafkaRetries)
}
