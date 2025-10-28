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
func (r *LogRepository) InsertLog(table string, log models.Log, blockHash string, blockNumber uint64, blockTimestamp uint64) error {
	ctx := context.Background()

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
func (r *LogRepository) InsertLogs(table string, logs []models.Log, blockHash string, blockNumber uint64, blockTimestamp uint64) error {
	if len(logs) == 0 {
		return nil
	}

	ctx := context.Background()

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
func (r *LogRepository) InsertLogsFromReceipt(table string, receipt models.Receipt, txHash string, txIndex uint32, blockHash string, blockNumber uint64, blockTimestamp uint64) error {
	if len(receipt.Logs) == 0 {
		return nil
	}

	ctx := context.Background()

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for logs from receipt insert: %v", err)
		return err
	}

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

	r.Logger.Debugf("Successfully inserted %d logs from receipt for transaction %s", len(receipt.Logs), txHash)
	return nil
}

// InsertLogsFromBlock вставляет логи из блока
func (r *LogRepository) InsertLogsFromBlock(table string, block models.Block) error {
	ctx := context.Background()

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for logs from block insert: %v", err)
		return err
	}

	logCount := 0
	for i, tx := range block.Transactions {
		if tx.Receipt != nil {
			for _, log := range tx.Receipt.Logs {
				row := convertLogToClickHouseRow(log, block.Hash, block.Number, block.Timestamp, uint32(i))
				err = batch.Append(row...)
				if err != nil {
					r.Logger.Errorf("Failed to append log to batch: %v", err)
					return err
				}
				logCount++
			}
		}
	}

	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for logs from block insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted %d logs from block %s", logCount, block.Hash)
	return nil
}

// convertLogToClickHouseRow конвертирует Log в строку для вставки в ClickHouse
func convertLogToClickHouseRow(log models.Log, blockHash string, blockNumber uint64, blockTimestamp uint64, txIndex uint32) []interface{} {
	timestamp := time.Unix(int64(blockTimestamp), 0)

	return []interface{}{
		blockNumber,          // block_number UInt64
		blockHash,            // block_hash FixedString(66)
		log.TransactionHash,   // transaction_hash FixedString(66)
		txIndex,              // transaction_index UInt32
		uint32(log.LogIndex), // log_index UInt32
		log.Address,          // address FixedString(42)
		log.Data,             // data String
		log.Topics,           // topics Array(FixedString(66))
		timestamp,            // block_timestamp DateTime64(3, 'UTC')
		// date и topic0 автоматически вычисляются
	}
}
