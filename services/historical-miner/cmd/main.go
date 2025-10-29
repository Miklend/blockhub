package main

import (
	"context"
	"lib/blocks/collector"
	fabricClient "lib/clients/fabric_client"
	"lib/models"
	"lib/utils/logging"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup
	logger := logging.GetLogger()
	logger.Info("Logger initialized successfully")

	cfg := models.GetConfig(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	clientsPool, err := fabricClient.NewProviderPool(cfg.ProviderHistorical, logger)
	if err != nil {
		logger.Fatalf("Failed to create consolidated provider pool: %v", err)
	}
	if len(cfg.ProviderHistorical) == 0 {
		logger.Fatal("ProviderHistorical config is empty.")
	}
	totalLimiterRate := cfg.ProviderHistorical[0].Limiter

	if totalLimiterRate <= 0 {
		logger.Warnf("Limiter is set to %d. Using default safe rate of 5 batches/second.", totalLimiterRate)
		totalLimiterRate = 5
	}

	blockCollector := collector.NewBlockCollector(clientsPool, totalLimiterRate, logger)

	startBlock := uint64(20000000)
	totalBatches := 10
	batchSize := 5
	for i := 0; i < totalBatches; i++ {

		var blocksToFetch []uint64
		for j := 0; j < batchSize; j++ {
			blocksToFetch = append(blocksToFetch, startBlock+uint64(i*batchSize+j))
		}

		wg.Add(1)
		go func(batchID int, blocks []uint64) {
			defer wg.Done()

			ctxTimeout, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()

			logger.Debugf("Batch %d: Requesting blocks %v", batchID, blocks)

			fetchedBlocks, err := blockCollector.PostBatch(ctxTimeout, blocks)

			if err != nil {
				logger.Errorf("Batch %d failed: %v", batchID, err)
			} else if len(fetchedBlocks) > 0 {
				logger.Infof("Batch %d successful. Fetched %d blocks. (First block: %d)", batchID, len(fetchedBlocks), fetchedBlocks[0].Number)
			}
		}(i+1, blocksToFetch)
	}

	// Ждем завершения всех горутин (основная часть теста)
	wg.Wait()
	logger.Info("All batch load test goroutines finished.")

	<-ctx.Done()
	logger.Info("Shutdown signal received, stopping services...")
}
