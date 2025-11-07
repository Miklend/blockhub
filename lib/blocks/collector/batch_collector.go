package collector

import (
	"context"
	"fmt"
	"lib/models"
)

// FetchBlocksBatch загружает блоки по номерам
func (bc *BlockCollector) FetchBlocksBatch(ctx context.Context, numbers []uint64) (map[models.Hash]models.BlockDTO, error) {
	if len(numbers) == 0 {
		return nil, nil
	}

	hexNumbers := make([]string, len(numbers))
	for i, num := range numbers {
		hexNumbers[i] = fmt.Sprintf("0x%x", num)
	}

	blockMap, err := bc.Client().BatchBlockWithReceiptByNumber(ctx, hexNumbers)
	if err != nil {
		return nil, fmt.Errorf("fetch blocks with receipts batch failed: %w", err)
	}

	return blockMap, nil
}

// // FetchReceiptsBatch загружает квитанции по номерам блоков через eth_getBlockReceipts
// func (bc *BlockCollector) FetchReceiptsBatch(ctx context.Context, numbers []uint64) (map[uint64][]models.Receipt, error) {
// 	elems := make([]rpc.BatchElem, 0, len(numbers))
// 	for _, num := range numbers {
// 		var raw json.RawMessage
// 		elems = append(elems, rpc.BatchElem{
// 			Method: "eth_getBlockReceipts",
// 			Args:   []interface{}{fmt.Sprintf("0x%x", num)},
// 			Result: &raw,
// 		})
// 	}

// 	if _, err := bc.DoBatch(ctx, elems); err != nil {
// 		return nil, err
// 	}

// 	receiptsMap := make(map[uint64][]models.Receipt)
// 	for i, e := range elems {
// 		if e.Error != nil {
// 			bc.logger.Warnf("receipts fetch error (number %d): %v", numbers[i], e.Error)
// 			continue
// 		}
// 		raw := *e.Result.(*json.RawMessage)
// 		receiptsMap[numbers[i]] = metrics.ParseBlockReceiptsJSON(raw)
// 	}
// 	return receiptsMap, nil
// }

// // FetchBlocksAndReceiptsBatch загружает блоки и их квитанции
// func (bc *BlockCollector) FetchBlocksAndReceiptsBatch(ctx context.Context, numbers []uint64) ([]models.Block, error) {
// 	blocks, err := bc.FetchBlocksBatch(ctx, numbers)
// 	if err != nil {
// 		return nil, fmt.Errorf("fetch blocks failed: %w", err)
// 	}

// 	receiptsMap, err := bc.FetchReceiptsBatch(ctx, numbers)
// 	if err != nil {
// 		return nil, fmt.Errorf("fetch receipts failed: %w", err)
// 	}

// 	// Связываем транзакции с квитанциями
// 	for i := range blocks {
// 		block := &blocks[i]
// 		blockReceipts, ok := receiptsMap[block.Number]
// 		if !ok {
// 			continue
// 		}
// 		if len(blockReceipts) != len(block.Transactions) {
// 			bc.logger.Warnf("block %d: number of receipts (%d) does not match transactions (%d)", block.Number, len(blockReceipts), len(block.Transactions))
// 		}
// 		for j := range block.Transactions {
// 			if j < len(blockReceipts) {
// 				block.Transactions[j].Receipt = &blockReceipts[j]
// 			}
// 		}
// 	}

// 	return blocks, nil
// }
