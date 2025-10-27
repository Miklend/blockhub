package tx

import (
	"context"
	"strconv"
	"time"

	clientsDB "lib/clients/db"
	"lib/models"
	"lib/utils/logging"
)

type TxRepository struct {
	Client clientsDB.ClickhouseClient
	Logger *logging.Logger
}

func NewTxRepository(client clientsDB.ClickhouseClient, logger *logging.Logger) *TxRepository {
	return &TxRepository{
		Client: client,
		Logger: logger,
	}
}

// InsertTx вставляет одну транзакцию в таблицу
func (r *TxRepository) InsertTx(table string, tx models.Tx) error {
	ctx := context.Background()

	// Для InsertTx нужны данные блока, используем пустые значения
	row := convertTxToClickHouseRow(tx, "", 0, 0)

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for transaction insert: %v", err)
		return err
	}

	err = batch.Append(row...)
	if err != nil {
		r.Logger.Errorf("Failed to append transaction to batch: %v", err)
		return err
	}

	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for transaction insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted transaction %s", tx.Hash)
	return nil
}

// InsertTxs вставляет массив транзакций в таблицу
func (r *TxRepository) InsertTxs(table string, txs []models.Tx) error {
	if len(txs) == 0 {
		return nil
	}

	ctx := context.Background()

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for transactions insert: %v", err)
		return err
	}

	for _, tx := range txs {
		row := convertTxToClickHouseRow(tx, "", 0, 0)
		err = batch.Append(row...)
		if err != nil {
			r.Logger.Errorf("Failed to append transaction %s to batch: %v", tx.Hash, err)
			return err
		}
	}

	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for transactions insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted %d transactions", len(txs))
	return nil
}

// InsertTxWithBlockData вставляет транзакцию с данными блока
func (r *TxRepository) InsertTxWithBlockData(table string, tx models.Tx, blockHash string, blockNumber uint64, blockTimestamp uint64) error {
	ctx := context.Background()

	row := convertTxToClickHouseRow(tx, blockHash, blockNumber, blockTimestamp)

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for transaction insert: %v", err)
		return err
	}

	err = batch.Append(row...)
	if err != nil {
		r.Logger.Errorf("Failed to append transaction to batch: %v", err)
		return err
	}

	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for transaction insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted transaction %s with block data", tx.Hash)
	return nil
}

// InsertTxsWithBlockData вставляет массив транзакций с данными блока
func (r *TxRepository) InsertTxsWithBlockData(table string, txs []models.Tx, blockHash string, blockNumber uint64, blockTimestamp uint64) error {
	if len(txs) == 0 {
		return nil
	}

	ctx := context.Background()

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for transactions insert: %v", err)
		return err
	}

	for _, tx := range txs {
		row := convertTxToClickHouseRow(tx, blockHash, blockNumber, blockTimestamp)
		err = batch.Append(row...)
		if err != nil {
			r.Logger.Errorf("Failed to append transaction %s to batch: %v", tx.Hash, err)
			return err
		}
	}

	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for transactions insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted %d transactions with block data", len(txs))
	return nil
}

// convertTxToClickHouseRow конвертирует Tx в строку для вставки в ClickHouse
func convertTxToClickHouseRow(tx models.Tx, blockHash string, blockNumber uint64, blockTimestamp uint64) []interface{} {
	timestamp := time.Unix(int64(blockTimestamp), 0)

	// Конвертируем to в указатель
	var to *string
	if tx.To != "" {
		to = &tx.To
	}

	// Конвертируем maxFeePerGas если есть
	var maxFeePerGas *uint64
	if tx.MaxFeePerGas != "" {
		if val, err := parseHexToUint64(tx.MaxFeePerGas); err == nil {
			maxFeePerGas = &val
		}
	}

	// Конвертируем maxPriorityFeePerGas если есть
	var maxPriorityFeePerGas *uint64
	if tx.MaxPriorityFeePerGas != "" {
		if val, err := parseHexToUint64(tx.MaxPriorityFeePerGas); err == nil {
			maxPriorityFeePerGas = &val
		}
	}

	// Конвертируем chainID
	chainID, _ := parseHexToUint64(tx.ChainID)

	// Конвертируем gasPrice
	gasPrice, _ := parseHexToUint64(tx.GasPrice)

	return []interface{}{
		tx.Hash,                     // hash
		blockHash,                   // block_hash
		blockNumber,                 // block_number
		uint32(tx.TransactionIndex), // transaction_index
		tx.From,                     // from
		to,                          // to
		tx.Value,                    // value
		tx.Gas,                      // gas
		gasPrice,                    // gas_price
		tx.Input,                    // input
		tx.Nonce,                    // nonce
		tx.Type,                     // type
		maxFeePerGas,                // max_fee_per_gas
		maxPriorityFeePerGas,        // max_priority_fee_per_gas
		chainID,                     // chain_id
		tx.V,                        // v
		tx.R,                        // r
		tx.S,                        // s
		"",                          // access_list (пустой для простоты)
		timestamp,                   // block_timestamp
		timestamp,                   // date (MATERIALIZED)
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
