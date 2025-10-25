package main

import (
	"blockhub/services/realtime-miner/internal/node/collector"
	"blockhub/services/realtime-miner/internal/node/worker"
	"context"
	"fmt"
	collectorLib "lib/blocks/collector"
	fabricClient "lib/clients/fabric_client"
	"lib/models"
	"lib/utils/logging"
	"os/signal"
	"syscall"
)

func main() {
	// Инициализация логгера
	logger := logging.GetLogger()
	logger.Info("Logger initialized successfully")

	// Загрузка конфигурации
	cfg := models.GetConfig(logger)
	if cfg.ProviderRealTime.ApiKey == "" {
		logger.Errorf("REALTIME_API_KEY is empty")
		return
	} else {
		logger.Infof("REALTIME_API_KEY found")
	}

	fmt.Println(cfg.Broker.BrockerType)

	// Создаем контекст с отменой и автоматической подпиской на сигналы
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Инициализация Alchemy клиента
	ProviderConfig := models.Provider{
		ProviderType: cfg.ProviderRealTime.ProviderType,
		BaseURL:      cfg.ProviderRealTime.BaseURL,
		NetworkName:  cfg.ProviderRealTime.NetworkName,
		ApiKey:       cfg.ProviderRealTime.ApiKey,
	}

	maxRetries := cfg.ProviderRealTime.MaxRetries
	providerClient, err := fabricClient.NewProvider(ProviderConfig, logger)
	if err != nil {
		logger.Errorf("Failed to create provider Client: %v", err)
	}
	defer providerClient.Close()

	// Инициализация Kafka клиента
	brockerClient := fabricClient.NewBroker(cfg.Broker, logger)
	defer func() {
		if err := brockerClient.Close(); err != nil {
			logger.Errorf("Failed to close Kafka client: %v", err)
		}
	}()

	// Инициализация BlockCollector
	blockCollector := collectorLib.NewBlockCollector(providerClient, logger)

	// Инициализация RealtimeCollector
	realtimeCollector := collector.NewRealtimeCollector(blockCollector)

	// Подписка на новые блоки
	blocksChan, err := realtimeCollector.SubscribeNewBlocks(ctx, maxRetries)
	if err != nil {
		logger.Fatalf("Failed to subscribe to new blocks: %v", err)
	}
	logger.Info("Subscribed to new blocks, waiting for incoming data...")

	blockTransfer := worker.NewBlockTransfer(logger, brockerClient).(*worker.BlockTransfer)

	go blockTransfer.TransferBlocks(ctx, blocksChan)
	// Ждем завершения по сигналу
	<-ctx.Done()
	logger.Info("Shutdown signal received, stopping services...")
}
