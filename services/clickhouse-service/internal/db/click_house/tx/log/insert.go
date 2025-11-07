package log

import (
	"context"
	"time"

	clientsDB "lib/clients/db"
	"lib/models"
	"lib/utils/logging"
)

type LogRepository struct {
	Client clientsDB.ClickhouseClient
	Logger *logging.Logger
}

func NewLogRepository(client clientsDB.ClickhouseClient, logger *logging.Logger) *LogRepository {
	return &LogRepository{
		Client: client,
		Logger: logger,
	}
}

// InsertLog вставляет один лог в таблицу
// Примечание: требует, чтобы данные блока были доступны
func (r *LogRepository) InsertLog(table string, log models.Log) error {
	ctx := context.Background()

	// Для InsertLog данные блока должны быть получены извне
	// Используем пустые значения, если блок недоступен (может привести к ошибкам)
	blockHash := ""
	blockNumber := uint64(0)
	blockTimestamp := time.Time{}

	row := convertLogToClickHouseRow(log, blockHash, blockNumber, blockTimestamp, 0)

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for log insert: %v", err)
		return err
	}

	err = batch.Append(row...)
	if err != nil {
		r.Logger.Errorf("Failed to append log to batch: %v", err)
		return err
	}

	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for log insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted log for transaction %s", log.TransactionHash)
	return nil
}

// InsertLogs вставляет массив логов в таблицу
// Примечание: требует, чтобы данные блока были доступны
func (r *LogRepository) InsertLogs(table string, logs []models.Log) error {
	if len(logs) == 0 {
		return nil
	}

	ctx := context.Background()

	// Для InsertLogs данные блока должны быть получены извне
	// Используем пустые значения, если блок недоступен (может привести к ошибкам)
	blockHash := ""
	blockNumber := uint64(0)
	blockTimestamp := time.Time{}

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for logs insert: %v", err)
		return err
	}

	for i, log := range logs {
		row := convertLogToClickHouseRow(log, blockHash, blockNumber, blockTimestamp, uint32(i))
		err = batch.Append(row...)
		if err != nil {
			r.Logger.Errorf("Failed to append log %d to batch: %v", i, err)
			return err
		}
	}

	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for logs insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted %d logs", len(logs))
	return nil
}

// InsertLogsFromReceipt вставляет логи из квитанции
// Примечание: требует, чтобы данные транзакции и блока были доступны
func (r *LogRepository) InsertLogsFromReceipt(table string, receipt models.Receipt) error {
	if len(receipt.Logs) == 0 {
		return nil
	}

	ctx := context.Background()

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for logs from receipt insert: %v", err)
		return err
	}

	// Для InsertLogsFromReceipt данные транзакции и блока должны быть получены извне
	// Используем пустые значения, если данные недоступны (может привести к ошибкам)
	txIndex := uint32(0)
	blockHash := ""
	blockNumber := uint64(0)
	blockTimestamp := time.Time{}

	for i, log := range receipt.Logs {
		row := convertLogToClickHouseRow(log, blockHash, blockNumber, blockTimestamp, txIndex)
		err = batch.Append(row...)
		if err != nil {
			r.Logger.Errorf("Failed to append log %d to batch: %v", i, err)
			return err
		}
	}

	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for logs from receipt insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted %d logs from receipt", len(receipt.Logs))
	return nil
}

// InsertLogsFromBlock вставляет логи из блока
func (r *LogRepository) InsertLogsFromBlock(table string, block models.Block) error {
	if len(block.Transactions) > 0 {
		r.Logger.Warn("InsertLogsFromBlock: block.Transactions does not contain receipt data; skipping logs insertion")
	}

	r.Logger.Debugf("InsertLogsFromBlock skipped for block %s: no receipt data available", block.Hash)
	return nil
}

// convertLogToClickHouseRow конвертирует Log в строку для вставки в ClickHouse
func convertLogToClickHouseRow(log models.Log, blockHash string, blockNumber uint64, blockTimestamp time.Time, txIndex uint32) []interface{} {
	timestamp := log.BlockTimestamp
	if timestamp.IsZero() {
		timestamp = blockTimestamp
	}

	return []interface{}{
		blockNumber,
		blockHash,
		log.TransactionHash,
		txIndex,
		uint32(log.LogIndex),
		log.Address,
		log.Data,
		copyStringSlice(log.Topics),
		timestamp,
	}
}
