package internal

import (
	"context"
	"encoding/json"
	"lib/models"
	"os"
	"time"

	"blockhub/services/realtime-miner/pkg"
	"lib/utils/logging"
)

const maxKafkaRetries = 3                      // Максимальное количество попыток отправки
const kafkaRetryDelay = 500 * time.Millisecond // Задержка между попытками

type BlockTransfer struct {
	Logger      *logging.Logger
	KafkaClient internal.MockKafkaClient
}

// NewBlockTransfer создаёт новый worker для отправки блоков в Kafka
func NewBlockTransfer(logger *logging.Logger, kafkaClient MockKafkaClient) pkg.Worker {
	return &BlockTransfer{
		Logger:      logger,
		KafkaClient: kafkaClient,
	}
}

// TransferBlocks слушает канал передачи блоков и отправляет их в Kafka
func (bt *BlockTransfer) TransferBlocks(ctx context.Context, in <-chan *models.Block) error {
	for {
		select {
		case <-ctx.Done():
			bt.Logger.Infof("block transfer stopped")
			return nil

		case block, ok := <-in:
			if !ok {
				bt.Logger.Warnf("block channel closed")
				return nil
			}

			//Сериализация блока в JSON
			data, err := json.Marshal(block)
			if err != nil {
				bt.Logger.Errorf("failed to serialize block %s: %v", block.Hash, err)
				continue
			}

			tst := os.WriteFile("Block_new.json", data, 0644)
			if tst != nil {
				bt.Logger.Errorf("failed to write block %s: %v", block.Hash, err)
				continue
			}

			if block.Number == 0 {
				bt.Logger.Debug("Skipping block with number 0")
				continue
			}

			//Отправка в Kafka с повторными попытками
			bt.sendWithRetry(ctx, block.Hash, data)
		}
	}
}

// повторные попытки
func (bt *BlockTransfer) sendWithRetry(ctx context.Context, key string, data []byte) {
	var err error

	for attempt := 1; attempt <= maxKafkaRetries; attempt++ {

		// Отправка в Kafka(сейчас в симулированный)
		err = bt.KafkaClient.SendMessage(key, data)

		if err == nil {
			bt.Logger.Infof("Block %s sent to Kafka successfully (attempt %d)", key, attempt)
			return
		}

		bt.Logger.Warnf("Failed to send block %s to Kafka (attempt %d/%d): %v", key, attempt, maxKafkaRetries, err)

		if attempt < maxKafkaRetries {
			select {
			case <-ctx.Done():
				bt.Logger.Warnf("Context cancelled during Kafka retry for block %s", key)
				return
			case <-time.After(kafkaRetryDelay):
			}
		}
	}
	bt.Logger.Errorf("FATAL: Failed to send block %s to Kafka after %d attempts. DROPPING MESSAGE.", key, maxKafkaRetries)
}
