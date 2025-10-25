package worker

import (
	"context"
	"encoding/json"
	"lib/models"
	"os"
	"time"

	"lib/clients/broker"

	"blockhub/services/realtime-miner/pkg"
	"lib/utils/logging"
)

const maxKafkaRetries = 3                      // Максимальное количество попыток отправки
const kafkaRetryDelay = 500 * time.Millisecond // Задержка между попытками

const topicKafka = "blocks"

type BlockTransfer struct {
	Logger      *logging.Logger
	KafkaClient broker.BrokerClient
}

// NewBlockTransfer создаёт новый worker для отправки блоков в Kafka
func NewBlockTransfer(logger *logging.Logger, kafkaClient broker.BrokerClient) pkg.Worker {
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

			// Сериализация блока в JSON
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

			// Создаем сообщение для брокера
			m := models.MessageBroker{
				Key:   []byte(block.Hash),
				Value: data,       // Используем сериализованные данные
				Topic: topicKafka, // Указываем топик
			}

			// Отправка в Kafka с повторными попытками
			bt.sendWithRetry(ctx, m)
		}
	}
}

// sendWithRetry повторные попытки отправки
func (bt *BlockTransfer) sendWithRetry(ctx context.Context, m models.MessageBroker) {
	var err error

	for attempt := 1; attempt <= maxKafkaRetries; attempt++ {
		// Отправка в Kafka
		err = bt.KafkaClient.SendMessage(ctx, m)

		if err == nil {
			bt.Logger.Infof("Block %s sent to Kafka successfully (attempt %d)", string(m.Key), attempt)
			return
		}

		bt.Logger.Warnf("Failed to send block %s to Kafka (attempt %d/%d): %v", string(m.Key), attempt, maxKafkaRetries, err)

		if attempt < maxKafkaRetries {
			select {
			case <-ctx.Done():
				bt.Logger.Warnf("Context cancelled during Kafka retry for block %s", string(m.Key))
				return
			case <-time.After(kafkaRetryDelay):
				// Ждем перед следующей попыткой
			}
		}
	}
	bt.Logger.Errorf("FATAL: Failed to send block %s to Kafka after %d attempts. DROPPING MESSAGE.", string(m.Key), maxKafkaRetries)
}
