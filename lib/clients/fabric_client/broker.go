package fabricClient

import (
	"lib/clients/broker"
	"lib/clients/broker/kafka"
	"lib/models"
	"lib/utils/logging"
)

const (
	KafkaBrokerType = "kafka"
	// Можно добавить другие типы: RabbitMQBrokerType, NATSBrokerType и т.д.
)

// NewBroker создает брокер указанного типа
func NewBroker(cfg models.Broker, logger *logging.Logger) broker.BrokerClient {
	switch cfg.BrockerType {
	case KafkaBrokerType:
		return kafka.NewKafkaBroker(cfg, logger)
	default:
		return nil
	}
}
