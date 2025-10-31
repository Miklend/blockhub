package main

import (
	"context"
	"lib/blocks/collector"
	fabricClient "lib/clients/fabric_client"
	"lib/clients/node"
	"lib/models"
	"lib/utils/logging"
	"os/signal"
	"syscall"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Logger initialized successfully")
	cfg := models.GetConfig(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	clients, err := fabricClient.NewProviderPool(cfg.ProviderHistorical, logger)
	if err != nil {
		logger.Fatalf("Failed to create client pool: %v", err)
	}

	if len(clients) < 2 {
		logger.Fatal("Configuration must contain at least 2 provider keys.")
	}

	// ⭐ Параметры задачи
	startBlock := uint64(23000000)
	endBlock := uint64(23000029) // Общий диапазон 200 блоков
	batchSize := uint64(30)

	// ⭐ Разделяем диапазон поровну на 2 ключа
	totalBlocks := endBlock - startBlock + 1
	halfBlocks := totalBlocks / 3

	// Диапазоны для ключей
	key1End := startBlock + halfBlocks - 1
	key2Start := key1End + 1

	go runJob(ctx, logger, clients[0], startBlock, key1End, batchSize)

	// ⭐ Запуск второй горутины (Ключ 2)
	go runJob(ctx, logger, clients[1], key2Start, endBlock, batchSize)

	go runJob(ctx, logger, clients[2], key2Start, endBlock, batchSize)

}

func runJob(ctx context.Context, logger *logging.Logger, client node.Provider, start, end, batchSize uint64) {
	blockCollector := collector.NewBlockCollector(client, 1, logger)

	current := start
	for current <= end {
		batchStart := current
		batchEnd := current + batchSize - 1

		// Коррекция последнего батча
		if batchEnd > end {
			batchEnd = end
		}

		batchBlocks := make([]uint64, 0, batchSize)
		for b := batchStart; b <= batchEnd; b++ {
			batchBlocks = append(batchBlocks, b)
		}

		// Вызов PostBatch
		fetchedBlocks, err := blockCollector.PostBatch(ctx, batchBlocks)

		if err != nil {
			logger.Errorf("Client failed on blocks %d-%d: %v", batchStart, batchEnd, err)
		} else {
			logger.Infof("Client success on blocks %d-%d", batchStart, batchEnd)
		}
		for j, block := range fetchedBlocks {
			logger.Infof("Recieved block %d: Hash: %s Transctions: %v", batchBlocks[j], block.Hash, len(block.Transactions))
		}

		// Переход к следующему батчу
		current = batchEnd + 1
	}
}

// func main() {
// 	var wg sync.WaitGroup
// 	logger := logging.GetLogger()
// 	logger.Info("Logger initialized successfully")

// 	cfg := models.GetConfig(logger)

// 	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
// 	defer stop()

// 	clientsPool, err := fabricClient.NewProviderPool(cfg.ProviderHistorical, logger)
// 	if err != nil {
// 		logger.Fatalf("Failed to create consolidated provider pool: %v", err)
// 	}
// 	if len(cfg.ProviderHistorical) == 0 {
// 		logger.Fatal("ProviderHistorical config is empty.")
// 	}
// 	totalLimiterRate := cfg.ProviderHistorical[0].Limiter

// 	if totalLimiterRate <= 0 {
// 		logger.Warnf("Limiter is set to %d. Using default safe rate of 5 batches/second.", totalLimiterRate)
// 		totalLimiterRate = 5
// 	}

// 	blockCollector := collector.NewBlockCollector(clientsPool, totalLimiterRate, logger)

// 	startBlock := uint64(20000000)
// 	totalBatches := 10
// 	batchSize := 5
// 	for i := 0; i < totalBatches; i++ {

// 		var blocksToFetch []uint64
// 		for j := 0; j < batchSize; j++ {
// 			blocksToFetch = append(blocksToFetch, startBlock+uint64(i*batchSize+j))
// 		}

// 		wg.Add(1)
// 		go func(batchID int, blocks []uint64) {
// 			defer wg.Done()

// 			ctxTimeout, cancel := context.WithTimeout(ctx, 15*time.Second)
// 			defer cancel()

// 			logger.Debugf("Batch %d: Requesting blocks %v", batchID, blocks)

// 			fetchedBlocks, err := blockCollector.PostBatch(ctxTimeout, blocks)

// 			if err != nil {
// 				logger.Errorf("Batch %d failed: %v", batchID, err)
// 			} else if len(fetchedBlocks) > 0 {
// 				logger.Infof("Batch %d successful. Fetched %d blocks. (First block: %d, block data)", batchID, len(fetchedBlocks), fetchedBlocks[0].Number,  models.Block.Transactions)
// 			}
// 		}(i+1, blocksToFetch)
// 	}

// 	// Ждем завершения всех горутин (основная часть теста)
// 	wg.Wait()
// 	logger.Info("All batch load test goroutines finished.")

// 	<-ctx.Done()
// 	logger.Info("Shutdown signal received, stopping services...")
//}
