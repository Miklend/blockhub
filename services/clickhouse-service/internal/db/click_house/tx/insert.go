package tx

import (
	"context"
	"strings"
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
	row := convertTxToClickHouseRow(tx, "", 0, time.Time{})

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
		row := convertTxToClickHouseRow(tx, "", 0, time.Time{})
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
// Примечание: требует, чтобы блок был доступен через блокчейн API или другой источник данных
func (r *TxRepository) InsertTxWithBlockData(table string, tx models.Tx) error {
	ctx := context.Background()

	// Для InsertTxWithBlockData данные блока должны быть получены извне
	// Используем пустые значения, если блок недоступен (может привести к ошибкам)
	blockHash := ""
	blockNumber := uint64(0)
	blockTimestamp := time.Time{}

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
// Примечание: требует, чтобы блок был доступен через блокчейн API или другой источник данных
func (r *TxRepository) InsertTxsWithBlockData(table string, txs []models.Tx) error {
	if len(txs) == 0 {
		return nil
	}

	ctx := context.Background()

	// Для InsertTxsWithBlockData данные блока должны быть получены извне
	// Используем пустые значения, если блок недоступен (может привести к ошибкам)
	blockHash := ""
	blockNumber := uint64(0)
	blockTimestamp := time.Time{}

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
func convertTxToClickHouseRow(tx models.Tx, blockHash string, blockNumber uint64, blockTimestamp time.Time) []interface{} {
	timestamp := blockTimestamp
	if timestamp.IsZero() {
		timestamp = tx.BlockTimestamp
	}

	var maxFeePerGas *uint64
	if tx.MaxFeePerGas != nil {
		val := uint64(*tx.MaxFeePerGas)
		maxFeePerGas = &val
	}

	var maxPriorityFeePerGas *uint64
	if tx.MaxPriorityFeePerGas != nil {
		val := uint64(*tx.MaxPriorityFeePerGas)
		maxPriorityFeePerGas = &val
	}

	value := trimHexPrefix(tx.Value)

	return []interface{}{
		tx.Hash,
		blockHash,
		blockNumber,
		uint32(tx.TransactionIndex),
		tx.From,
		tx.To,
		value,
		uint64(tx.Gas),
		uint64(tx.GasPrice),
		tx.Input,
		uint64(tx.Nonce),
		uint8(tx.Type),
		maxFeePerGas,
		maxPriorityFeePerGas,
		uint64(tx.ChainID),
		tx.V,
		tx.R,
		tx.S,
		tx.AccessList,
		timestamp,
		timestamp,
	}
}

func trimHexPrefix(value string) string {
	if len(value) >= 2 {
		lower := strings.ToLower(value[:2])
		if lower == "0x" {
			return value[2:]
		}
	}
	return value
}
