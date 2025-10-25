package internal

import (
	"lib/utils/logging"
	"time"
)

// имитация поведения кафки
type MockKafkaClient struct {
	logger logging.Logger
}

// NewMockKafkaClient создает заглушку
func NewMockKafkaClient(logger logging.Logger) *MockKafkaClient {
	return &MockKafkaClient{
		logger: logger,
	}
}

// SendMessage имитирует отправку в кафку.
func (m *MockKafkaClient) SendMessage(key string, data []byte) error {

	// тест задержка
	time.Sleep(200 * time.Millisecond)

	m.logger.Debugf("MOCK KAFKA: Successfully simulated sending block with key %s. Data size: %d bytes.", key, len(data))

	return nil
}

func (m *MockKafkaClient) Close() error {
	m.logger.Info("MOCK KAFKA: Client closed.")
	return nil
}
