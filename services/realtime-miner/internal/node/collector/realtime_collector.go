package collector

import (
	"blockhub/services/realtime-miner/internal/node"
	"context"
	collectorLib "lib/blocks/collector"
	"lib/models"
	"lib/utils/logging"
	"math"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type realtimeCollector struct {
	bc     *collectorLib.BlockCollector
	logger *logging.Logger
}

func NewRealtimeCollector(bc *collectorLib.BlockCollector) node.RtCollector {
	return &realtimeCollector{
		bc:     bc,
		logger: bc.Logger(),
	}
}

func (rc *realtimeCollector) SubscribeNewBlocks(ctx context.Context, maxRetries int) (<-chan *models.Block, error) {

	rc.logger.Info("Starting new block subscription setup.")

	out := make(chan *models.Block, 100)

	go func() {
		defer close(out)

		retryCount := 0

		for {
			select {
			case <-ctx.Done():
				rc.logger.Debug("Context cancelled, stopping reconnection loop.")
				return
			default:
			}

			headers := make(chan *types.Header)
			client := rc.bc.Client()

			rc.logger.Infof("Attempting to subscribe (Retry %d)...", retryCount)

			sub, err := client.SubscribeNewHead(ctx, headers)

			if err != nil {
				rc.logger.Errorf("Failed to establish subscription: %v", err)

				delay := time.Duration(math.Pow(2, math.Min(float64(retryCount), 6.0))) * time.Second

				retryCount++

				rc.logger.Warnf("Subscription failed. Waiting %v before next attempt...", delay)

				select {
				case <-time.After(delay):
					continue
				case <-ctx.Done():
					return
				}
			}

			rc.logger.Info("Successfully subscribed to new block headers.")
			retryCount = 0
			rc.processSubscription(ctx, sub, headers, out, maxRetries)

			rc.logger.Warn("Subscription disconnected. Attempting to reconnect...")
		}
	}()

	return out, nil
}

func (rc *realtimeCollector) processSubscription(
	ctx context.Context,
	sub ethereum.Subscription,
	headers <-chan *types.Header,
	out chan<- *models.Block,
	maxBlockRetries int,
) {
	defer sub.Unsubscribe()
	rc.logger.Info("Block header stream started.")

	for {
		select {
		case <-ctx.Done():
			rc.logger.Debug("Context cancelled inside processing loop.")
			return

		case err := <-sub.Err():
			rc.logger.Errorf("Subscription error during header stream: %v", err)
			return

		case header := <-headers:
			blockNumber := header.Number.Uint64()
			rc.logger.Debugf("New block header received: #%d", blockNumber)

			initialDelay := 500 * time.Millisecond
			rc.logger.Debugf("Waiting %v for block %d to be available...", initialDelay, blockNumber)

			select {
			case <-time.After(initialDelay):
			case <-ctx.Done():
				rc.logger.Debug("Context cancelled during initial delay")
				return
			}

			var block *models.Block
			var blockErr error
			for attempt := 1; attempt <= maxBlockRetries; attempt++ {
				block, blockErr = rc.bc.CollectBlockByNumber(ctx, blockNumber)
				if blockErr == nil {
					break
				}

				if strings.Contains(blockErr.Error(), "not found") ||
					strings.Contains(blockErr.Error(), "not available") {
					rc.logger.Debugf("Block %d not available yet (attempt %d/%d), waiting...",
						blockNumber, attempt, maxBlockRetries)
				} else {
					rc.logger.Warnf("Attempt %d/%d failed for block %d: %v",
						attempt, maxBlockRetries, blockNumber, blockErr)
				}

				if attempt < maxBlockRetries {
					retryDelay := time.Duration(attempt) * 500 * time.Millisecond
					rc.logger.Debugf("Waiting %v before retry %d for block %d",
						retryDelay, attempt+1, blockNumber)

					select {
					case <-time.After(retryDelay):
						// продолжаем retry
					case <-ctx.Done():
						rc.logger.Debug("Context cancelled during retry delay")
						return
					}
				}
			}

			if blockErr != nil {
				if strings.Contains(blockErr.Error(), "not found") {
					rc.logger.Warnf("Block %d still not available after %d attempts",
						blockNumber, maxBlockRetries)
				} else {
					rc.logger.Errorf("All %d attempts failed for block %d: %v",
						maxBlockRetries, blockNumber, blockErr)
				}
				continue
			}

			rc.logger.Infof("Successfully processed block #%d", blockNumber)

			select {
			case out <- block:
				// Block sent successfully
			case <-ctx.Done():
				rc.logger.Debug("Context cancelled while sending block")
				return
			default:
				rc.logger.Warn("Block channel is full, dropping block")
			}
		}
	}
}
