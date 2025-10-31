package main

import (
	"context"
	"lib/blocks/collector"
	fabricClient "lib/clients/fabric_client"
	"lib/clients/node"
	"lib/models"
	"lib/utils/logging"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	logger := logging.GetLogger()
	cfg := models.GetConfig(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	clients, err := fabricClient.NewProviderPool(cfg.ProviderHistorical, logger)
	if err != nil {
		logger.Fatalf("Failed to create client pool: %v", err)
	}

	// Параметры
	startBlock := uint64(23000000)
	endBlock := uint64(23000029)
	batchSize := uint64(10) // уменьшаем для теста

	// Правильное разделение на клиентов
	totalBlocks := endBlock - startBlock + 1
	blocksPerClient := totalBlocks / uint64(len(clients))

	var wg sync.WaitGroup

	for i, client := range clients {
		wg.Add(1)

		clientStart := startBlock + uint64(i)*blocksPerClient
		clientEnd := clientStart + blocksPerClient - 1

		// Последний клиент получает остаток
		if i == len(clients)-1 {
			clientEnd = endBlock
		}

		go func(client node.Provider, start, end uint64, clientNum int) {
			defer wg.Done()
			runJob(ctx, logger, client, start, end, batchSize, clientNum)
		}(client, clientStart, clientEnd, i)
	}

	// Ждем завершения ВСЕХ горутин
	wg.Wait()
	logger.Info("All jobs completed")
}

func runJob(ctx context.Context, logger *logging.Logger, client node.Provider,
	start, end, batchSize uint64, clientNum int) {

	blockCollector := collector.NewBlockCollector(client, 1, logger)
	current := start

	for current <= end {
		// Проверяем, не отменен ли контекст
		select {
		case <-ctx.Done():
			logger.Infof("Client %d: context canceled, stopping at block %d", clientNum, current)
			return
		default:
		}

		batchEnd := current + batchSize - 1
		if batchEnd > end {
			batchEnd = end
		}

		batchBlocks := make([]uint64, 0, batchSize)
		for b := current; b <= batchEnd; b++ {
			batchBlocks = append(batchBlocks, b)
		}

		fetchedBlocks, err := blockCollector.PostBatch(ctx, batchBlocks)
		if err != nil {
			logger.Errorf("Client %d failed on blocks %d-%d: %v", clientNum, current, batchEnd, err)
			// Решаем, продолжать ли при ошибках
			if ctx.Err() != nil { // если контекст отменен - выходим
				return
			}
		} else {
			logger.Infof("Client %d success on blocks %d-%d", clientNum, current, batchEnd)
			for j, block := range fetchedBlocks {
				logger.Infof("Client %d: Block %d: Hash: %s Transactions: %v",
					clientNum, batchBlocks[j], block.Hash, len(block.Transactions))
			}
		}

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
