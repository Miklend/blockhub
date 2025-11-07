package block

import (
	"context"
	"strconv"

	clientsDB "lib/clients/db"
	"lib/models"
	"lib/utils/logging"
)

type BlockRepository struct {
	Client clientsDB.ClickhouseClient
	Logger *logging.Logger
}

func NewBlockRepository(client clientsDB.ClickhouseClient, logger *logging.Logger) *BlockRepository {
	return &BlockRepository{
		Client: client,
		Logger: logger,
	}
}

// InsertBlock вставляет один блок в таблицу
func (r *BlockRepository) InsertBlock(table string, block models.Block) error {
	ctx := context.Background()

	// Конвертируем блок в строку для ClickHouse
	row := convertBlockToClickHouseRow(block)

	// Подготавливаем batch для вставки
	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for block insert: %v", err)
		return err
	}

	// Добавляем блок в batch
	err = batch.Append(row...)
	if err != nil {
		r.Logger.Errorf("Failed to append block to batch: %v", err)
		return err
	}

	// Выполняем вставку
	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for block insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted block %s (number: %d)", block.Hash, block.Number)
	return nil
}

// InsertBlocks вставляет массив блоков в таблицу
func (r *BlockRepository) InsertBlocks(table string, blocks []models.Block) error {
	if len(blocks) == 0 {
		return nil
	}

	ctx := context.Background()

	// Подготавливаем batch для вставки
	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for blocks insert: %v", err)
		return err
	}

	// Конвертируем и добавляем все блоки в batch
	for _, block := range blocks {
		row := convertBlockToClickHouseRow(block)
		err = batch.Append(row...)
		if err != nil {
			r.Logger.Errorf("Failed to append block %s to batch: %v", block.Hash, err)
			return err
		}
	}

	// Выполняем вставку
	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for blocks insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted %d blocks", len(blocks))
	return nil
}

// convertBlockToClickHouseRow конвертирует Block в строку для вставки в ClickHouse
func convertBlockToClickHouseRow(block models.Block) []interface{} {
	// Извлекаем хеши транзакций
	txHashes := copyStringSlice(block.Transactions)

	// Конвертируем baseFeePerGas если есть
	var baseFeePerGas *uint64
	if block.BaseFeePerGas != nil {
		val := uint64(*block.BaseFeePerGas)
		baseFeePerGas = &val
	}

	// Конвертируем difficulty и totalDifficulty
	difficulty, _ := parseHexToUint256(block.Difficulty)
	totalDifficulty, _ := parseHexToUint256(block.TotalDifficulty)

	return []interface{}{
		block.Hash,                             // hash
		uint64(block.Number),                   // number
		block.ParentHash,                       // parent_hash
		formatUint64ToHex(uint64(block.Nonce)), // nonce (String)
		block.Sha3Uncles,                       // sha3_uncles
		block.LogsBloom,                        // logs_bloom
		block.TransactionsRoot,                 // transactions_root
		block.StateRoot,                        // state_root
		block.ReceiptsRoot,                     // receipts_root
		block.Miner,                            // miner
		difficulty,                             // difficulty
		totalDifficulty,                        // total_difficulty
		uint64(block.Size),                     // size
		block.ExtraData,                        // extra_data
		uint64(block.GasLimit),                 // gas_limit
		uint64(block.GasUsed),                  // gas_used
		baseFeePerGas,                          // base_fee_per_gas
		block.Timestamp,                        // timestamp
		block.MixHash,                          // mix_hash
		txHashes,                               // transactions
		copyStringSlice(block.Uncles),          // uncles
		block.Timestamp,                        // date (MATERIALIZED)
	}
}

// Вспомогательные функции
func parseHexToUint256(hexStr string) (string, error) {
	if hexStr == "" {
		return "", nil
	}
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	return hexStr, nil
}

func formatUint64ToHex(n uint64) string {
	return "0x" + strconv.FormatUint(n, 16)
}

func copyStringSlice(src []string) []string {
	if len(src) == 0 {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}
