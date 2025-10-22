package broker

import (
	"context"
	"lib/models"
)

type BrokerClient interface {
	// Producer методы
	SendMessage(ctx context.Context, msg models.Message) error
	SendMessages(ctx context.Context, msgs []models.Message) error

	// Consumer методы
	Subscribe(ctx context.Context, topic string, handler models.MessageHandler) error
	SubscribeWithGroup(ctx context.Context, topic, groupID string, handler models.MessageHandler) error

	// Admin методы
	CreateTopic(ctx context.Context, topic string, partitions int, replicationFactor int) error
	HealthCheck(ctx context.Context) error

	// Закрытие соединений
	Close() error
}
