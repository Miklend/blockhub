package broker

import (
	"context"
	"lib/models"
)

type BrokerClient interface {
	// Producer методы
	SendMessage(ctx context.Context, msg models.MessageBroker) error
	SendMessages(ctx context.Context, msgs []models.MessageBroker) error

	// Consumer методы
	Subscribe(ctx context.Context, topic string, handler models.MessageHandlerBroker) error
	SubscribeWithGroup(ctx context.Context, topic, groupID string, handler models.MessageHandlerBroker) error

	// Admin методы
	CreateTopic(ctx context.Context, topic string, partitions int, replicationFactor int) error
	HealthCheck(ctx context.Context) error

	// Закрытие соединений
	Close() error
}
