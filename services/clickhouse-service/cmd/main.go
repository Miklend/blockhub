package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	clickhouseRepo "clickhouse-service/internal/db/click_house"
	clickhouseClient "lib/clients/db/clickhouse"
	"lib/models"
	"lib/utils/logging"
)

func main() {
	// Инициализируем логгер
	logger := logging.GetLogger()
	logger.Info("Starting ClickHouse service...")

	// Получаем конфигурацию
	config := models.GetConfig(logger)
	logger.Info("Configuration loaded successfully")

	// Создаем контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализируем ClickHouse клиент
	clickhouseClient, err := clickhouseClient.NewClient(ctx, config.Clickhouse)
	if err != nil {
		logger.Fatalf("Failed to initialize ClickHouse client: %v", err)
	}
	logger.Info("ClickHouse client initialized successfully")

	// Создаем репозиторий
	repo := clickhouseRepo.NewClickhouseService(clickhouseClient, logger)
	logger.Info("ClickHouse repository initialized successfully")

	// Проверяем соединение
	if err := clickhouseClient.Ping(ctx); err != nil {
		logger.Fatalf("Failed to ping ClickHouse: %v", err)
	}
	logger.Info("ClickHouse connection verified")

	// Настраиваем graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Ждем сигнал завершения
	<-sigChan
	logger.Info("Received shutdown signal, closing connections...")

	// Закрываем соединения
	if err := repo.Close(); err != nil {
		logger.Errorf("Error closing repository: %v", err)
	}

	logger.Info("ClickHouse service stopped")
}
