package broker

import (
	"context"
	"lib/models"
)

type MockBrokerClient struct{}

func NewMockBrokerClient() BrokerClient {
	return &MockBrokerClient{}
}

func (m *MockBrokerClient) SendMessage(ctx context.Context, msg models.MessageBroker) error {
	// Сообщение успешно отправлено
	return nil
}

func (m *MockBrokerClient) SendMessages(ctx context.Context, msgs []models.MessageBroker) error {
	// Сообщения успешно отправлены
	return nil
}

func (m *MockBrokerClient) Subscribe(ctx context.Context, topic string, handler models.MessageHandlerBroker) error {
	// Подписка успешно создана
	return nil
}

func (m *MockBrokerClient) SubscribeWithGroup(ctx context.Context, topic, groupID string, handler models.MessageHandlerBroker) error {
	// Подписка с группой успешно создана
	return nil
}

func (m *MockBrokerClient) CreateTopic(ctx context.Context, topic string, partitions, replicationFactor int) error {
	// Топик успешно создан
	return nil
}

func (m *MockBrokerClient) HealthCheck(ctx context.Context) error {
	// Брокер доступен
	return nil
}

func (m *MockBrokerClient) Close() error {
	// Соединение успешно закрыто
	return nil
}
