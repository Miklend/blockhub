package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"lib/blocks/metrics"
	"lib/models"
	"math/big"

	"github.com/ethereum/go-ethereum/rpc"
)

// PostBatch отправляет батч запросов на Alchemy и возвращает массив models.Block.
func (bc *BlockCollector) PostBatch(ctx context.Context, blocks []uint64) ([]models.Block, error) {
	if len(blocks) == 0 {
		return nil, fmt.Errorf("no block numbers provided")
	}

	var batch []rpc.BatchElem
	for _, n := range blocks {
		numHex := "0x" + big.NewInt(int64(n)).Text(16)

		batch = append(batch,
			rpc.BatchElem{
				Method: "eth_getBlockByNumber",
				Args:   []interface{}{numHex, true},
				Result: new(json.RawMessage),
			},
			rpc.BatchElem{
				Method: "eth_getBlockReceipts",
				Args:   []interface{}{map[string]string{"blockNumber": numHex}},
				Result: new(json.RawMessage),
			},
		)
	}

	bc.logger.Infof("sending batch request for %d blocks", len(blocks))

	// Отправляем батч запрос
	if err := bc.client.BatchCallContext(ctx, batch); err != nil {
		bc.logger.Errorf("batch request failed: %v", err)
		return nil, fmt.Errorf("batch call failed: %w", err)
	}

	var result []models.Block
	var hasErrors bool

	for i := 0; i < len(batch); i += 2 {
		blockIndex := i / 2
		blockNumber := blocks[blockIndex]

		// Проверка ошибок конкретных RPC-элементов
		if batch[i].Error != nil {
			bc.logger.WithFields(map[string]interface{}{
				"block":  blockNumber,
				"method": batch[i].Method,
			}).Errorf("block fetch failed: %v", batch[i].Error)
			hasErrors = true
			continue
		}
		if batch[i+1].Error != nil {
			bc.logger.WithFields(map[string]interface{}{
				"block":  blockNumber,
				"method": batch[i+1].Method,
			}).Errorf("receipts fetch failed: %v", batch[i+1].Error)
			hasErrors = true
			continue
		}

		rawBlock, ok1 := batch[i].Result.(*json.RawMessage)
		rawReceipts, ok2 := batch[i+1].Result.(*json.RawMessage)

		if !ok1 || !ok2 || rawBlock == nil || rawReceipts == nil {
			bc.logger.WithFields(map[string]interface{}{
				"block": blockNumber,
			}).Warn("invalid batch element — nil or unexpected type")
			hasErrors = true
			continue
		}

		if len(*rawBlock) == 0 {
			bc.logger.WithFields(map[string]interface{}{
				"block": blockNumber,
			}).Warn("empty block response")
			hasErrors = true
			continue
		}

		block := metrics.NewBlockFromJSON(*rawBlock, *rawReceipts)
		result = append(result, block)

		bc.logger.WithFields(map[string]interface{}{
			"block": block.Number,
			"txs":   len(block.Transactions),
		}).Debug("block parsed successfully")
	}

	// Итог
	if hasErrors {
		bc.logger.Warnf("batch for %d blocks completed with partial errors", len(blocks))
		return result, fmt.Errorf("some batch elements failed, see logs for details")
	}

	bc.logger.Infof("batch for %d blocks completed successfully", len(blocks))
	return result, nil
}
