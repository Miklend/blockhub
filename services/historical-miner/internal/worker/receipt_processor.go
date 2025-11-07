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
	jobsChan    <-chan *models.BlockDTO
	kafkaClient broker.BrokerClient
}

func NewReceiptProcessor(
	logger *logging.Logger,
	collector collector.BlockCollector,
	jobs <-chan *models.BlockDTO,
	kafkaClient broker.BrokerClient,
) *ReceiptProcessor {
	return &ReceiptProcessor{
		logger:      logger,
		collector:   collector,
		jobsChan:    jobs,
		kafkaClient: kafkaClient,
	}
}

func (p *ReceiptProcessor) ProcessReceipts(ctx context.Context) {
	for {
		select {
		case blockData, ok := <-p.jobsChan:
			if !ok {
				return
			}

			p.logger.Infof("SUCCESSFULLY PROCESSED: Block %d is complete. Total Transactions: %d, Gas Used: %d.",
				blockData.Number, len(blockData.Transactions), blockData.GasUsed)

			data, err := json.Marshal(blockData)
			if err != nil {
				p.logger.Errorf("failed to serialize block %s: %v", blockData.Hash, err)
				continue
			}

			p.sendWithRetry(ctx, blockData.Hash, data)

		case <-ctx.Done():
			p.logger.Info("Receipt Processor received shutdown signal, stopping")
			return
		}

	}
}

func (p *ReceiptProcessor) sendWithRetry(ctx context.Context, blockHash models.Hash, data []byte) {
	var err error
	key := []byte(blockHash)

	if blockHash == "" {
		var tempBlock models.BlockDTO
		if err := json.Unmarshal(data, &tempBlock); err == nil {
			key = []byte(fmt.Sprintf("%d", tempBlock.Number))
		} else {
			key = []byte(blockHash)
		}
	} else {
		key = []byte(blockHash)
	}

	m := models.MessageBroker{
		Key:   key,
		Value: data,
		Topic: topicKafka,
	}

	for attempt := 1; attempt <= maxKafkaRetries; attempt++ {
		err = p.kafkaClient.SendMessage(ctx, m)

		if err == nil {
			p.logger.Infof("Block %s sent to Kafka successfully (attempt %d)", blockHash, attempt)
			return
		}

		p.logger.Warnf("Failed to send block %s to Kafka (attempt %d/%d): %v", blockHash, attempt, maxKafkaRetries, err)

		if attempt < maxKafkaRetries {
			select {
			case <-ctx.Done():
				p.logger.Warnf("Context cancelled during Kafka retry for block %s", blockHash)
				return
			case <-time.After(kafkaRetryDelay):
			}
		}
	}
	p.logger.Errorf("FATAL: Failed to send block %s to Kafka after %d attempts. DROPPING MESSAGE.", blockHash, maxKafkaRetries)
}
