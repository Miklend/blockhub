package main

import (
	worker "blockhub/services/historical-miner/internal"
	"context"
	"lib/blocks/collector"
	fabricClient "lib/clients/fabric_client"
	"lib/clients/node"
	"lib/models"
	"lib/utils/logging"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	NUM_BLOCK_FETCHERS           = 3
	NUM_BLOCK_RECEIPT_PROCESSORS = 5
	RECEIPT_RATE_LIMIT           = 2.0
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

	blockJobsChan := make(chan uint64, 1000)
	receiptJobChan := make(chan *models.Block, 1000)

	// Параметры
	startBlock := uint64(23000000)
	endBlock := uint64(23000029)
	batchSize := uint64(10) // уменьшаем для теста

	// Правильное разделение на клиентов
	totalBlocks := endBlock - startBlock + 1
	blocksPerClient := totalBlocks / uint64(len(clients))

	var wg sync.WaitGroup
	var wgWorkers sync.WaitGroup

	for i := 0; i < NUM_BLOCK_FETCHERS; i++ {
		wgWorkers.Add(1)
		client := clients[i%len(clients)]
		blockCollector := collector.NewBlockCollector(client, 100.0, logger)
		fetcher := worker.NewBlockFetcher(logger, blockCollector, blockJobsChan, receiptJobChan)

		go func(i int) {
			defer wgWorkers.Done()
			fetcher.ProcessBlocks(ctx)
			logger.Debugf("Block Fetcher %d stopped", i)
		}(i)
	}

	for i := 0; i < NUM_BLOCK_RECEIPT_PROCESSORS; i++ {
		wgWorkers.Add(1)
		clients := clients[i%len(clients)]

		receiptCollector := collector.NewBlockCollector(clients, RECEIPT_RATE_LIMIT, logger)
		processor := worker.NewReceiptProcessor(logger, *receiptCollector, receiptJobChan)
		go func(i int) {
			defer wgWorkers.Done()
			processor.ProcessReceipts(ctx)
			logger.Debugf("ReceiptProcessor %d stopped", i)
		}(i)
	}

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
			runMasterJob(ctx, logger, client, start, end, batchSize, clientNum, blockJobsChan)
		}(client, clientStart, clientEnd, i)
	}

	// Ждем завершения ВСЕХ горутин
	wg.Wait()
	logger.Info("All jobs completed")
	close(blockJobsChan)

	wgWorkers.Wait()
	logger.Info("All workers completed.")
}

func runMasterJob(ctx context.Context, logger *logging.Logger, client node.Provider,
	start, end, batchSize uint64, clientNum int, blockJobCh chan<- uint64) {

	current := start

	for current <= end {
		select {
		case <-ctx.Done():
			logger.Infof("Master Job %d: context canceled.", clientNum)
			return
		default:
		}

		batchEnd := current + batchSize - 1
		if batchEnd > end {
			batchEnd = end
		}

		time.Sleep(1 * time.Second)

		for b := current; b <= batchEnd; b++ {
			select {
			case blockJobCh <- b:
			case <-ctx.Done():
				return
			}
		}

		current = batchEnd + 1
	}
}

// func runJob(ctx context.Context, logger *logging.Logger, client node.Provider,
// 	start, end, batchSize uint64, clientNum int, rateLimit float64) {
// 	blockCollector := collector.NewBlockCollector(client, rateLimit, logger)
// 	current := start

// 	for current <= end {
// 		// Проверяем, не отменен ли контекст
// 		select {
// 		case <-ctx.Done():
// 			logger.Infof("Client %d: context canceled, stopping at block %d", clientNum, current)
// 			return
// 		default:
// 		}

// 		batchEnd := current + batchSize - 1
// 		if batchEnd > end {
// 			batchEnd = end
// 		}

// 		batchBlocks := make([]uint64, 0, batchSize)
// 		for b := current; b <= batchEnd; b++ {
// 			batchBlocks = append(batchBlocks, b)
// 		}

// 		delay := time.Duration(1.0 / rateLimit * float64(time.Second))
// 		select {
// 		case <-ctx.Done():
// 			logger.Infof("Client %d: context canceled during delay.", clientNum)
// 			return
// 		case <-time.After(delay):
// 			// Продолжаем работу
// 		}

// 		fetchedBlocks, err := blockCollector.FetchBlocksBatch(ctx, batchBlocks)
// 		if err != nil {
// 			logger.Errorf("Client %d failed on blocks %d-%d: %v", clientNum, current, batchEnd, err)
// 			// Решаем, продолжать ли при ошибках
// 			if ctx.Err() != nil { // если контекст отменен - выходим
// 				return
// 			}
// 		} else {
// 			logger.Infof("Client %d success on blocks %d-%d", clientNum, current, batchEnd)
// 			for j, block := range fetchedBlocks {
// 				logger.Infof("Client %d: Block %d: Hash: %s Transactions: %v",
// 					clientNum, batchBlocks[j], block.Hash, len(block.Transactions))
// 			}
// 		}

// 		current = batchEnd + 1
// 	}
// }
