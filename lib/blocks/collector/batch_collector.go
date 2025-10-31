package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"lib/blocks/metrics"
	"lib/models"

	"github.com/ethereum/go-ethereum/rpc"
)

// DoBatch выполняет произвольный RPC-батч
func (bc *BlockCollector) DoBatch(ctx context.Context, elems []rpc.BatchElem) ([]rpc.BatchElem, error) {
	if len(elems) == 0 {
		return nil, fmt.Errorf("empty batch request list")
	}
	if bc.limiter != nil {
		if err := bc.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter wait failed: %w", err)
		}
	}

	if err := bc.Client().BatchCallContext(ctx, elems); err != nil {
		return nil, fmt.Errorf("batch call failed: %w", err)
	}

	return elems, nil
}

// FetchBlocksBatch загружает блоки по номерам
func (bc *BlockCollector) FetchBlocksBatch(ctx context.Context, numbers []uint64) ([]models.Block, error) {
	elems := make([]rpc.BatchElem, 0, len(numbers))
	for _, num := range numbers {
		var raw json.RawMessage
		elems = append(elems, rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{fmt.Sprintf("0x%x", num), true},
			Result: &raw,
		})
	}

	if _, err := bc.DoBatch(ctx, elems); err != nil {
		return nil, err
	}

	blocks := make([]models.Block, 0, len(elems))
	for i, e := range elems {
		if e.Error != nil {
			bc.logger.Warnf("block fetch error (number %d): %v", numbers[i], e.Error)
			continue
		}
		raw := *e.Result.(*json.RawMessage)
		blocks = append(blocks, metrics.ParseBlockJSON(raw))
	}
	return blocks, nil
}

// FetchReceiptsBatch загружает квитанции по номерам блоков через eth_getBlockReceipts
func (bc *BlockCollector) FetchReceiptsBatch(ctx context.Context, numbers []uint64) (map[uint64][]models.Receipt, error) {
	elems := make([]rpc.BatchElem, 0, len(numbers))
	for _, num := range numbers {
		var raw json.RawMessage
		elems = append(elems, rpc.BatchElem{
			Method: "eth_getBlockReceipts",
			Args:   []interface{}{fmt.Sprintf("0x%x", num)},
			Result: &raw,
		})
	}

	if _, err := bc.DoBatch(ctx, elems); err != nil {
		return nil, err
	}

	receiptsMap := make(map[uint64][]models.Receipt)
	for i, e := range elems {
		if e.Error != nil {
			bc.logger.Warnf("receipts fetch error (number %d): %v", numbers[i], e.Error)
			continue
		}
		raw := *e.Result.(*json.RawMessage)
		receiptsMap[numbers[i]] = metrics.ParseBlockReceiptsJSON(raw)
	}
	return receiptsMap, nil
}

// FetchBlocksAndReceiptsBatch загружает блоки и их квитанции
func (bc *BlockCollector) FetchBlocksAndReceiptsBatch(ctx context.Context, numbers []uint64) ([]models.Block, error) {
	blocks, err := bc.FetchBlocksBatch(ctx, numbers)
	if err != nil {
		return nil, fmt.Errorf("fetch blocks failed: %w", err)
	}

	receiptsMap, err := bc.FetchReceiptsBatch(ctx, numbers)
	if err != nil {
		return nil, fmt.Errorf("fetch receipts failed: %w", err)
	}

	// Связываем транзакции с квитанциями
	for i := range blocks {
		block := &blocks[i]
		blockReceipts, ok := receiptsMap[block.Number]
		if !ok {
			continue
		}
		if len(blockReceipts) != len(block.Transactions) {
			bc.logger.Warnf("block %d: number of receipts (%d) does not match transactions (%d)", block.Number, len(blockReceipts), len(block.Transactions))
		}
		for j := range block.Transactions {
			if j < len(blockReceipts) {
				block.Transactions[j].Receipt = &blockReceipts[j]
			}
		}
	}

	return blocks, nil
}
