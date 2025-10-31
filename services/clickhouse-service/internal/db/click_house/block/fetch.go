package block

import (
	"context"
	"strconv"

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
	block := models.Block{
		Hash:             row.Hash,
		Number:           row.Number,
		ParentHash:       row.ParentHash,
		Nonce:            row.Nonce,
		Sha3Uncles:       row.Sha3Uncles,
		LogsBloom:        row.LogsBloom,
		TransactionsRoot: row.TransactionsRoot,
		StateRoot:        row.StateRoot,
		ReceiptsRoot:     row.ReceiptsRoot,
		Miner:            row.Miner,
		Difficulty:       "0x" + row.Difficulty,
		Size:             row.Size,
		ExtraData:        row.ExtraData,
		GasLimit:         row.GasLimit,
		GasUsed:          row.GasUsed,
		Timestamp:        uint64(row.Timestamp.Unix()),
		MixHash:          row.MixHash,
		Uncles:           row.Uncles,
	}

	// Конвертируем baseFeePerGas если есть
	if row.BaseFeePerGas != nil {
		block.BaseFeePerGas = formatUint64ToHex(*row.BaseFeePerGas)
	}

	// Создаем транзакции (только хеши)
	block.Transactions = make([]models.Tx, len(row.Transactions))
	for i, txHash := range row.Transactions {
		block.Transactions[i] = models.Tx{Hash: txHash}
	}

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
		block := models.Block{
			Hash:             row.Hash,
			Number:           row.Number,
			ParentHash:       row.ParentHash,
			Nonce:            row.Nonce,
			Sha3Uncles:       row.Sha3Uncles,
			LogsBloom:        row.LogsBloom,
			TransactionsRoot: row.TransactionsRoot,
			StateRoot:        row.StateRoot,
			ReceiptsRoot:     row.ReceiptsRoot,
			Miner:            row.Miner,
			Difficulty:       "0x" + row.Difficulty,
			Size:             row.Size,
			ExtraData:        row.ExtraData,
			GasLimit:         row.GasLimit,
			GasUsed:          row.GasUsed,
			Timestamp:        uint64(row.Timestamp.Unix()),
			MixHash:          row.MixHash,
			Uncles:           row.Uncles,
		}

		// Конвертируем baseFeePerGas если есть
		if row.BaseFeePerGas != nil {
			block.BaseFeePerGas = formatUint64ToHex(*row.BaseFeePerGas)
		}

		// Создаем транзакции (только хеши)
		block.Transactions = make([]models.Tx, len(row.Transactions))
		for j, txHash := range row.Transactions {
			block.Transactions[j] = models.Tx{Hash: txHash}
		}

		blocks[i] = block
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
	block := models.Block{
		Hash:             row.Hash,
		Number:           row.Number,
		ParentHash:       row.ParentHash,
		Nonce:            row.Nonce,
		Sha3Uncles:       row.Sha3Uncles,
		LogsBloom:        row.LogsBloom,
		TransactionsRoot: row.TransactionsRoot,
		StateRoot:        row.StateRoot,
		ReceiptsRoot:     row.ReceiptsRoot,
		Miner:            row.Miner,
		Difficulty:       "0x" + row.Difficulty,
		Size:             row.Size,
		ExtraData:        row.ExtraData,
		GasLimit:         row.GasLimit,
		GasUsed:          row.GasUsed,
		Timestamp:        uint64(row.Timestamp.Unix()),
		MixHash:          row.MixHash,
		Uncles:           row.Uncles,
	}

	if row.BaseFeePerGas != nil {
		block.BaseFeePerGas = formatUint64ToHex(*row.BaseFeePerGas)
	}

	block.Transactions = make([]models.Tx, len(row.Transactions))
	for i, txHash := range row.Transactions {
		block.Transactions[i] = models.Tx{Hash: txHash}
	}

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
		block := models.Block{
			Hash:             row.Hash,
			Number:           row.Number,
			ParentHash:       row.ParentHash,
			Nonce:            row.Nonce,
			Sha3Uncles:       row.Sha3Uncles,
			LogsBloom:        row.LogsBloom,
			TransactionsRoot: row.TransactionsRoot,
			StateRoot:        row.StateRoot,
			ReceiptsRoot:     row.ReceiptsRoot,
			Miner:            row.Miner,
			Difficulty:       "0x" + row.Difficulty,
			Size:             row.Size,
			ExtraData:        row.ExtraData,
			GasLimit:         row.GasLimit,
			GasUsed:          row.GasUsed,
			Timestamp:        uint64(row.Timestamp.Unix()),
			MixHash:          row.MixHash,
			Uncles:           row.Uncles,
		}

		if row.BaseFeePerGas != nil {
			block.BaseFeePerGas = formatUint64ToHex(*row.BaseFeePerGas)
		}

		block.Transactions = make([]models.Tx, len(row.Transactions))
		for j, txHash := range row.Transactions {
			block.Transactions[j] = models.Tx{Hash: txHash}
		}

		blocks[i] = block
	}

	r.Logger.Debugf("Successfully fetched %d blocks in range %d-%d", len(blocks), fromBlock, toBlock)
	return blocks, nil
}

// parseHexToUint64Safe безопасно парсит hex строку в uint64
func parseHexToUint64Safe(hexStr string) uint64 {
	if hexStr == "" {
		return 0
	}
	// Убираем префикс 0x если есть
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	val, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0
	}
	return val
}
