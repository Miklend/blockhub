package kafka

import (
	"context"
	"fmt"
	"lib/clients/broker"
	"lib/models"
	"lib/utils/logging"
	"log"

	"github.com/segmentio/kafka-go"
)

// KafkaBroker реализация Broker для Apache Kafka
type KafkaBroker struct {
	config  models.Broker
	writers map[string]*kafka.Writer
	readers map[string]*kafka.Reader
	admin   *kafka.Client
}

// NewKafkaBroker создает новый Kafka брокер
func NewKafkaBroker(cfg models.Broker, logger *logging.Logger) broker.BrokerClient {

	return &KafkaBroker{
		config:  cfg,
		writers: make(map[string]*kafka.Writer),
		readers: make(map[string]*kafka.Reader),
		admin: &kafka.Client{
			Addr: kafka.TCP(cfg.Brokers...),
		},
	}
}

// SendMessage отправляет одно сообщение
func (k *KafkaBroker) SendMessage(ctx context.Context, msg models.Message) error {
	writer := k.getWriter(msg.Topic)

	kafkaMsg := kafka.Message{
		Key:   msg.Key,
		Value: msg.Value,
		Topic: msg.Topic,
	}

	// Конвертируем headers
	for key, value := range msg.Headers {
		kafkaMsg.Headers = append(kafkaMsg.Headers, kafka.Header{
			Key:   key,
			Value: []byte(value),
		})
	}

	return writer.WriteMessages(ctx, kafkaMsg)
}

// SendMessages отправляет несколько сообщений
func (k *KafkaBroker) SendMessages(ctx context.Context, msgs []models.Message) error {
	if len(msgs) == 0 {
		return nil
	}

	topic := msgs[0].Topic
	writer := k.getWriter(topic)

	kafkaMsgs := make([]kafka.Message, len(msgs))
	for i, msg := range msgs {
		kafkaMsgs[i] = kafka.Message{
			Key:   msg.Key,
			Value: msg.Value,
			Topic: msg.Topic,
		}

		// Конвертируем headers
		for key, value := range msg.Headers {
			kafkaMsgs[i].Headers = append(kafkaMsgs[i].Headers, kafka.Header{
				Key:   key,
				Value: []byte(value),
			})
		}
	}

	return writer.WriteMessages(ctx, kafkaMsgs...)
}

// Subscribe подписывается на топик без consumer group
func (k *KafkaBroker) Subscribe(ctx context.Context, topic string, handler models.MessageHandler) error {
	reader := k.getReader(topic, "")

	go k.consumeLoop(ctx, reader, handler)
	return nil
}

// SubscribeWithGroup подписывается на топик с consumer group
func (k *KafkaBroker) SubscribeWithGroup(ctx context.Context, topic, groupID string, handler models.MessageHandler) error {
	reader := k.getReader(topic, groupID)

	go k.consumeLoop(ctx, reader, handler)
	return nil
}

// CreateTopic создает топик
func (k *KafkaBroker) CreateTopic(ctx context.Context, topic string, partitions, replicationFactor int) error {
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}

	response, err := k.admin.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{topicConfig},
	})

	if err != nil {
		return err
	}

	if response.Errors[topic] != nil {
		return response.Errors[topic]
	}

	return nil
}

// HealthCheck проверяет доступность брокера
func (k *KafkaBroker) HealthCheck(ctx context.Context) error {
	_, err := k.admin.Metadata(ctx, &kafka.MetadataRequest{})
	return err
}

// Close закрывает все соединения
func (k *KafkaBroker) Close() error {
	var errs []error

	// Закрываем writers
	for _, writer := range k.writers {
		if err := writer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	// Закрываем readers
	for _, reader := range k.readers {
		if err := reader.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing kafka connections: %v", errs)
	}

	return nil
}

// Вспомогательные методы

func (k *KafkaBroker) getWriter(topic string) *kafka.Writer {
	if writer, exists := k.writers[topic]; exists {
		return writer
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(k.config.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    k.config.BatchSize,
		BatchTimeout: k.config.BatchTimeout,
		Async:        k.config.Async,
	}

	k.writers[topic] = writer
	return writer
}

func (k *KafkaBroker) getReader(topic, groupID string) *kafka.Reader {
	key := topic + "_" + groupID
	if reader, exists := k.readers[key]; exists {
		return reader
	}

	readerConfig := kafka.ReaderConfig{
		Brokers: k.config.Brokers,
		Topic:   topic,
		GroupID: groupID,
	}

	if k.config.StartOffset > 0 {
		readerConfig.StartOffset = k.config.StartOffset
	}

	reader := kafka.NewReader(readerConfig)
	k.readers[key] = reader
	return reader
}

func (k *KafkaBroker) consumeLoop(ctx context.Context, reader *kafka.Reader, handler models.MessageHandler) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			kafkaMsg, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			msg := models.Message{
				Key:     kafkaMsg.Key,
				Value:   kafkaMsg.Value,
				Topic:   kafkaMsg.Topic,
				Headers: make(map[string]string),
			}

			// Конвертируем headers
			for _, header := range kafkaMsg.Headers {
				msg.Headers[header.Key] = string(header.Value)
			}

			if err := handler(ctx, msg); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}
}
