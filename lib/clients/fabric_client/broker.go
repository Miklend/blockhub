package fabricClient

import (
	"fmt"
	"lib/clients/broker"
	"lib/clients/broker/kafka"
	"lib/models"
	"lib/utils/logging"
)

const (
	kafkaBrokerType = "kafka"
	// Можно добавить другие типы: RabbitMQBrokerType, NATSBrokerType и т.д.
)

// NewBroker создает брокер указанного типа
func NewBroker(cfg models.Broker, logger *logging.Logger) (broker.BrokerClient, error) {
	switch cfg.BrockerType {
	case kafkaBrokerType:
		return kafka.NewKafkaBroker(cfg, logger), nil
	default:
		return nil, fmt.Errorf("not found broker client")
	}
}
