package main

// import (
// 	worker "blockhub/services/historical-miner/internal/worker"
// 	"context"
// 	"lib/blocks/collector"
// 	fabricClient "lib/clients/fabric_client"
// 	"lib/clients/node"
// 	"lib/models"
// 	"lib/utils/logging"
// 	"os/signal"
// 	"sync"
// 	"syscall"
// 	"time"
// )

// const (
// 	NUM_BLOCK_FETCHERS           = 5
// 	NUM_BLOCK_RECEIPT_PROCESSORS = 5
// 	RECEIPT_RATE_LIMIT           = 2.0
// )

// func main() {
// 	logger := logging.GetLogger()
// 	cfg := models.GetConfig(logger)

// 	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
// 	defer stop()

// 	clients, err := fabricClient.NewProviderPool(cfg.ProviderHistorical, logger)
// 	if err != nil {
// 		logger.Fatalf("Failed to create client pool: %v", err)
// 	}

// 	blockJobsChan := make(chan uint64, 1000)
// 	receiptJobChan := make(chan *models.BlockDTO, 1000)

// 	// Параметры
// 	startBlock := uint64(23000000)
// 	endBlock := uint64(23000019)
// 	batchSize := uint64(10) // уменьшаем для теста

// 	// Правильное разделение на клиентов
// 	totalBlocks := endBlock - startBlock + 1
// 	blocksPerClient := totalBlocks / uint64(len(clients))

// 	var wg sync.WaitGroup
// 	var wgWorkers sync.WaitGroup

// 	for i := 0; i < NUM_BLOCK_FETCHERS; i++ {
// 		wgWorkers.Add(1)
// 		client := clients[i%len(clients)]
// 		blockCollector := collector.NewBlockCollector(client, 100.0, logger)
// 		fetcher := worker.NewBlockFetcher(logger, blockCollector, blockJobsChan, receiptJobChan)

// 		go func(i int) {
// 			defer wgWorkers.Done()
// 			fetcher.ProcessBlocks(ctx)
// 			logger.Debugf("Block Fetcher %d stopped", i)
// 		}(i)
// 	}

// 	for i := 0; i < NUM_BLOCK_RECEIPT_PROCESSORS; i++ {
// 		wgWorkers.Add(1)
// 		clients := clients[i%len(clients)]
// 		kafkaClient := fabricClient.NewBroker(cfg.Broker, logger)
// 		receiptCollector := collector.NewBlockCollector(clients, RECEIPT_RATE_LIMIT, logger)
// 		processor := worker.NewReceiptProcessor(logger, *receiptCollector, receiptJobChan, kafkaClient)
// 		go func(i int) {
// 			defer wgWorkers.Done()
// 			processor.ProcessReceipts(ctx)
// 			logger.Debugf("ReceiptProcessor %d stopped", i)
// 		}(i)
// 	}

// 	for i, client := range clients {
// 		wg.Add(1)

// 		clientStart := startBlock + uint64(i)*blocksPerClient
// 		clientEnd := clientStart + blocksPerClient - 1

// 		// Последний клиент получает остаток
// 		if i == len(clients)-1 {
// 			clientEnd = endBlock
// 		}

// 		go func(client node.Provider, start, end uint64, clientNum int) {
// 			defer wg.Done()
// 			runMasterJob(ctx, logger, client, start, end, batchSize, clientNum, blockJobsChan)
// 		}(client, clientStart, clientEnd, i)
// 	}

// 	// Ждем завершения ВСЕХ горутин
// 	wg.Wait()
// 	logger.Info("All jobs completed")
// 	close(blockJobsChan)

// 	wgWorkers.Wait()
// 	logger.Info("All workers completed.")
// }

// func runMasterJob(ctx context.Context, logger *logging.Logger, client node.Provider,
// 	start, end, batchSize uint64, clientNum int, blockJobCh chan<- uint64) {

// 	current := start

// 	for current <= end {
// 		select {
// 		case <-ctx.Done():
// 			logger.Infof("Master Job %d: context canceled.", clientNum)
// 			return
// 		default:
// 		}

// 		batchEnd := current + batchSize - 1
// 		if batchEnd > end {
// 			batchEnd = end
// 		}

// 		time.Sleep(50 * time.Millisecond)

// 		for b := current; b <= batchEnd; b++ {
// 			select {
// 			case blockJobCh <- b:
// 			case <-ctx.Done():
// 				return
// 			}
// 		}

// 		current = batchEnd + 1
// 	}
// }
