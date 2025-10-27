package block

import (
	"context"
	"strconv"
	"time"

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
	// Конвертируем timestamp из Unix в time.Time
	timestamp := time.Unix(int64(block.Timestamp), 0)

	// Извлекаем хеши транзакций
	txHashes := make([]string, len(block.Transactions))
	for i, tx := range block.Transactions {
		txHashes[i] = tx.Hash
	}

	// Конвертируем baseFeePerGas если есть
	var baseFeePerGas *uint64
	if block.BaseFeePerGas != "" {
		if val, err := parseHexToUint64(block.BaseFeePerGas); err == nil {
			baseFeePerGas = &val
		}
	}

	// Конвертируем difficulty и totalDifficulty
	difficulty, _ := parseHexToUint256(block.Difficulty)
	totalDifficulty := difficulty // В реальном проекте нужно вычислять

	return []interface{}{
		block.Hash,                     // hash
		block.Number,                   // number
		block.ParentHash,               // parent_hash
		formatUint64ToHex(block.Nonce), // nonce
		block.Sha3Uncles,               // sha3_uncles
		block.LogsBloom,                // logs_bloom
		block.TransactionsRoot,         // transactions_root
		block.StateRoot,                // state_root
		block.ReceiptsRoot,             // receipts_root
		block.Miner,                    // miner
		difficulty,                     // difficulty
		totalDifficulty,                // total_difficulty
		block.Size,                     // size
		block.ExtraData,                // extra_data
		block.GasLimit,                 // gas_limit
		block.GasUsed,                  // gas_used
		baseFeePerGas,                  // base_fee_per_gas
		timestamp,                      // timestamp
		block.MixHash,                  // mix_hash
		txHashes,                       // transactions
		block.Uncles,                   // uncles
		timestamp,                      // date (MATERIALIZED)
	}
}

// Вспомогательные функции для парсинга
func parseHexToUint64(hexStr string) (uint64, error) {
	if hexStr == "" {
		return 0, nil
	}
	// Убираем префикс 0x если есть
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	return strconv.ParseUint(hexStr, 16, 64)
}

func parseHexToUint256(hexStr string) (string, error) {
	// Для UInt256 в ClickHouse используем строку
	return hexStr, nil
}

func formatUint64ToHex(n uint64) string {
	return "0x" + strconv.FormatUint(n, 16)
}
