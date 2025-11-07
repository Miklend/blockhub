package block

import (
	"context"

	"clickhouse-service/internal/db/click_house/rowtypes"
	"lib/models"
)

// FetchBlock получает блок по хешу
func (r *BlockRepository) FetchBlock(table string, hashBlock string) (models.Block, error) {
	ctx := context.Background()

	var result []rowtypes.BlockRow

	query := "SELECT * FROM " + table + " WHERE hash = ? LIMIT 1"
	err := r.Client.Select(ctx, &result, query, hashBlock)
	if err != nil {
		r.Logger.Errorf("Failed to fetch block %s: %v", hashBlock, err)
		return models.Block{}, err
	}

	if len(result) == 0 {
		return models.Block{}, nil
	}

	// Конвертируем результат в модель Block
	row := result[0]
	block := blockRowToModel(row)

	r.Logger.Debugf("Successfully fetched block %s (number: %d)", block.Hash, block.Number)
	return block, nil
}

// FetchBlocks получает блоки по хешам
func (r *BlockRepository) FetchBlocks(table string, hashBlocks []string) ([]models.Block, error) {
	if len(hashBlocks) == 0 {
		return []models.Block{}, nil
	}

	ctx := context.Background()

	var result []rowtypes.BlockRow

	query := "SELECT * FROM " + table + " WHERE hash IN (?)"
	err := r.Client.Select(ctx, &result, query, hashBlocks)
	if err != nil {
		r.Logger.Errorf("Failed to fetch blocks: %v", err)
		return nil, err
	}

	// Конвертируем результаты в модели Block
	blocks := make([]models.Block, len(result))
	for i, row := range result {
		blocks[i] = blockRowToModel(row)
	}

	r.Logger.Debugf("Successfully fetched %d blocks", len(blocks))
	return blocks, nil
}

// FetchBlockByNumber получает блок по номеру
func (r *BlockRepository) FetchBlockByNumber(table string, blockNumber uint64) (models.Block, error) {
	ctx := context.Background()

	var result []rowtypes.BlockRow

	query := "SELECT * FROM " + table + " WHERE number = ? LIMIT 1"
	err := r.Client.Select(ctx, &result, query, blockNumber)
	if err != nil {
		r.Logger.Errorf("Failed to fetch block by number %d: %v", blockNumber, err)
		return models.Block{}, err
	}

	if len(result) == 0 {
		return models.Block{}, nil
	}

	// Конвертируем результат в модель Block (аналогично FetchBlock)
	row := result[0]
	block := blockRowToModel(row)

	r.Logger.Debugf("Successfully fetched block by number %d (hash: %s)", blockNumber, block.Hash)
	return block, nil
}

// FetchBlocksByRange получает блоки в диапазоне номеров
func (r *BlockRepository) FetchBlocksByRange(table string, fromBlock, toBlock uint64) ([]models.Block, error) {
	ctx := context.Background()

	var result []rowtypes.BlockRow

	query := "SELECT * FROM " + table + " WHERE number >= ? AND number <= ? ORDER BY number"
	err := r.Client.Select(ctx, &result, query, fromBlock, toBlock)
	if err != nil {
		r.Logger.Errorf("Failed to fetch blocks by range %d-%d: %v", fromBlock, toBlock, err)
		return nil, err
	}

	// Конвертируем результаты в модели Block (аналогично FetchBlocks)
	blocks := make([]models.Block, len(result))
	for i, row := range result {
		blocks[i] = blockRowToModel(row)
	}

	r.Logger.Debugf("Successfully fetched %d blocks in range %d-%d", len(blocks), fromBlock, toBlock)
	return blocks, nil
}

func blockRowToModel(row rowtypes.BlockRow) models.Block {
	block := models.Block{
		Hash:             row.Hash,
		Number:           uint(row.Number),
		ParentHash:       row.ParentHash,
		Nonce:            uint(row.Nonce),
		Sha3Uncles:       row.Sha3Uncles,
		LogsBloom:        row.LogsBloom,
		TransactionsRoot: row.TransactionsRoot,
		StateRoot:        row.StateRoot,
		ReceiptsRoot:     row.ReceiptsRoot,
		Miner:            row.Miner,
		Difficulty:       prefixHex(row.Difficulty),
		TotalDifficulty:  prefixHex(row.TotalDifficulty),
		Size:             uint(row.Size),
		ExtraData:        row.ExtraData,
		GasLimit:         uint(row.GasLimit),
		GasUsed:          uint(row.GasUsed),
		Timestamp:        row.Timestamp,
		MixHash:          row.MixHash,
		Transactions:     FetchcopyStringSlice(row.Transactions),
		Uncles:           FetchcopyStringSlice(row.Uncles),
	}

	block.BaseFeePerGas = uintPtrFromUint64(row.BaseFeePerGas)

	return block
}

func prefixHex(value string) string {
	if value == "" {
		return value
	}
	if len(value) >= 2 && value[:2] == "0x" {
		return value
	}
	return "0x" + value
}

func FetchcopyStringSlice(src []string) []string {
	if len(src) == 0 {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

func uintPtrFromUint64(src *uint64) *uint {
	if src == nil {
		return nil
	}
	val := uint(*src)
	return &val
}
