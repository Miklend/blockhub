package collector

import (
	"blockhub/services/realtime-miner/internal/node"
	"context"
	collectorLib "lib/blocks/collector"
	"lib/models"
	"lib/utils/logging"

	"fmt"
)

type realtimeCollector struct {
	bc                *collectorLib.BlockCollector
	logger            *logging.Logger
	HistoricalJobChan chan<- uint64
}

func NewRealtimeCollector(bc *collectorLib.BlockCollector, historicalJobChan chan<- uint64) node.RtCollector {
	return &realtimeCollector{
		bc:                bc,
		logger:            bc.Logger(),
		HistoricalJobChan: historicalJobChan,
	}
}

func (rc *realtimeCollector) SubscribeNewBlocks(ctx context.Context, maxRetries int) (<-chan *models.BlockDTO, error) {

	rc.logger.Info("Starting new block subscription setup.")

	out := make(chan *models.BlockDTO, 100)
	client := rc.bc.Client()
	rc.logger.Info("Attempting to subscribe to blocks with receipts...")

	sub, err := client.SubscribeBlockWithReceipts(ctx, out)

	if err != nil {

		return nil, fmt.Errorf("failed to establish block subscription: %w", err)
	}

	go func() {
		select {
		case err := <-sub.Err():
			if err != nil {
				rc.logger.Errorf("Block subscription failed permanently: %v", err)
			}
		case <-ctx.Done():
			rc.logger.Debug("Context cancelled, closing subscription watch.")
		}
		sub.Unsubscribe()
	}()

	rc.logger.Info("Successfully started block subscription. Data will arrive in channel.")

	return out, nil
}

// func (rc *realtimeCollector) processSubscription(
// 	ctx context.Context,
// 	sub ethereum.Subscription,
// 	headers <-chan *types.Header,
// 	out chan<- *models.BlockDTO,
// 	maxBlockRetries int,
// ) {
// 	defer sub.Unsubscribe()
// 	rc.logger.Info("Block header stream started.")

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			rc.logger.Debug("Context cancelled inside processing loop.")
// 			return

// 		case err := <-sub.Err():
// 			rc.logger.Errorf("Subscription error during header stream: %v", err)
// 			return

// 		case header := <-headers:
// 			blockNumber := header.Number.Uint64()
// 			rc.logger.Debugf("New block header received: #%d", blockNumber)

// 			initialDelay := 500 * time.Millisecond
// 			rc.logger.Debugf("Waiting %v for block %d to be available...", initialDelay, blockNumber)

// 			select {
// 			case <-time.After(initialDelay):
// 			case <-ctx.Done():
// 				rc.logger.Debug("Context cancelled during initial delay")
// 				return
// 			}

// 			var block *models.BlockDTO
// 			var blockErr error
// 			for attempt := 1; attempt <= maxBlockRetries; attempt++ {
// 				block, blockErr = rc.bc.CollectBlockByNumber(ctx, blockNumber)
// 				if blockErr == nil {
// 					break
// 				}

// 				if strings.Contains(blockErr.Error(), "not found") ||
// 					strings.Contains(blockErr.Error(), "not available") {
// 					rc.logger.Debugf("Block %d not available yet (attempt %d/%d), waiting...",
// 						blockNumber, attempt, maxBlockRetries)

// 				} else {
// 					rc.logger.Warnf("Attempt %d/%d failed for block %d: %v",
// 						attempt, maxBlockRetries, blockNumber, blockErr)
// 				}

// 				if attempt < maxBlockRetries {
// 					retryDelay := time.Duration(attempt) * 500 * time.Millisecond
// 					rc.logger.Debugf("Waiting %v before retry %d for block %d",
// 						retryDelay, attempt+1, blockNumber)

// 					select {
// 					case <-time.After(retryDelay):
// 						// продолжаем retry
// 					case <-ctx.Done():
// 						rc.logger.Debug("Context cancelled during retry delay")
// 						return
// 					}
// 				}
// 			}

// 			if blockErr != nil {
// 				if strings.Contains(blockErr.Error(), "not found") {
// 					rc.logger.Warnf("Block %d still not available after %d attempts",
// 						blockNumber, maxBlockRetries)
// 					select {
// 					case rc.HistoricalJobChan <- blockNumber:
// 						rc.logger.Infof("Delegated block %d to historical miner", blockNumber)
// 					case <-ctx.Done():
// 						rc.logger.Debug("Context cancelled during transfer to historical miner")
// 						return
// 					default:
// 						rc.logger.Errorf("historical job channel full, dropping block %d", blockNumber)
// 					}
// 				} else {
// 					rc.logger.Errorf("All %d attempts failed for block %d: %v",
// 						maxBlockRetries, blockNumber, blockErr)
// 				}
// 				continue
// 			}

// 			rc.logger.Infof("Successfully processed block #%d", blockNumber)

// 			select {
// 			case out <- block:
// 				// Block sent successfully
// 			case <-ctx.Done():
// 				rc.logger.Debug("Context cancelled while sending block")
// 				return
// 			default:
// 				rc.logger.Warn("Block channel is full, dropping block")
// 			}
// 		}
// 	}
// }
