package log

import (
	"context"
	"time"

	"lib/models"
)

// FetchLogsByTransaction получает логи по хешу транзакции
func (r *LogRepository) FetchLogsByTransaction(table string, txHash string) ([]models.Log, error) {
	ctx := context.Background()

	var result []struct {
		BlockNumber      uint64    `ch:"block_number"`
		BlockHash        string    `ch:"block_hash"`
		TransactionHash  string    `ch:"transaction_hash"`
		TransactionIndex uint32    `ch:"transaction_index"`
		LogIndex         uint32    `ch:"log_index"`
		Address          string    `ch:"address"`
		Data             string    `ch:"data"`
		Topics           []string  `ch:"topics"`
		BlockTimestamp   time.Time `ch:"block_timestamp"`
		Topic0           string    `ch:"topic0"`
	}

	query := "SELECT * FROM " + table + " WHERE transaction_hash = ? ORDER BY log_index"
	err := r.Client.Select(ctx, &result, query, txHash)
	if err != nil {
		r.Logger.Errorf("Failed to fetch logs for transaction %s: %v", txHash, err)
		return nil, err
	}

	// Конвертируем результаты в модели Log
	logs := make([]models.Log, len(result))
	for i, row := range result {
		logs[i] = models.Log{
			Address:         row.Address,
			Topics:          row.Topics,
			Data:            row.Data,
			TransactionHash: row.TransactionHash,
			LogIndex:        uint64(row.LogIndex),
			Removed:         false, // По умолчанию false
		}
	}

	r.Logger.Debugf("Successfully fetched %d logs for transaction %s", len(logs), txHash)
	return logs, nil
}

// FetchLogsByBlock получает логи по хешу блока
func (r *LogRepository) FetchLogsByBlock(table string, blockHash string) ([]models.Log, error) {
	ctx := context.Background()

	var result []struct {
		BlockNumber      uint64    `ch:"block_number"`
		BlockHash        string    `ch:"block_hash"`
		TransactionHash  string    `ch:"transaction_hash"`
		TransactionIndex uint32    `ch:"transaction_index"`
		LogIndex         uint32    `ch:"log_index"`
		Address          string    `ch:"address"`
		Data             string    `ch:"data"`
		Topics           []string  `ch:"topics"`
		BlockTimestamp   time.Time `ch:"block_timestamp"`
		Topic0           string    `ch:"topic0"`
	}

	query := "SELECT * FROM " + table + " WHERE block_hash = ? ORDER BY transaction_index, log_index"
	err := r.Client.Select(ctx, &result, query, blockHash)
	if err != nil {
		r.Logger.Errorf("Failed to fetch logs for block %s: %v", blockHash, err)
		return nil, err
	}

	// Конвертируем результаты в модели Log
	logs := make([]models.Log, len(result))
	for i, row := range result {
		logs[i] = models.Log{
			Address:         row.Address,
			Topics:          row.Topics,
			Data:            row.Data,
			TransactionHash: row.TransactionHash,
			LogIndex:        uint64(row.LogIndex),
			Removed:         false, // По умолчанию false
		}
	}

	r.Logger.Debugf("Successfully fetched %d logs for block %s", len(logs), blockHash)
	return logs, nil
}

// FetchLogsByBlockNumber получает логи по номеру блока
func (r *LogRepository) FetchLogsByBlockNumber(table string, blockNumber uint64) ([]models.Log, error) {
	ctx := context.Background()

	var result []struct {
		BlockNumber      uint64    `ch:"block_number"`
		BlockHash        string    `ch:"block_hash"`
		TransactionHash  string    `ch:"transaction_hash"`
		TransactionIndex uint32    `ch:"transaction_index"`
		LogIndex         uint32    `ch:"log_index"`
		Address          string    `ch:"address"`
		Data             string    `ch:"data"`
		Topics           []string  `ch:"topics"`
		BlockTimestamp   time.Time `ch:"block_timestamp"`
		Topic0           string    `ch:"topic0"`
	}

	query := "SELECT * FROM " + table + " WHERE block_number = ? ORDER BY transaction_index, log_index"
	err := r.Client.Select(ctx, &result, query, blockNumber)
	if err != nil {
		r.Logger.Errorf("Failed to fetch logs for block number %d: %v", blockNumber, err)
		return nil, err
	}

	// Конвертируем результаты в модели Log
	logs := make([]models.Log, len(result))
	for i, row := range result {
		logs[i] = models.Log{
			Address:         row.Address,
			Topics:          row.Topics,
			Data:            row.Data,
			TransactionHash: row.TransactionHash,
			LogIndex:        uint64(row.LogIndex),
			Removed:         false, // По умолчанию false
		}
	}

	r.Logger.Debugf("Successfully fetched %d logs for block number %d", len(logs), blockNumber)
	return logs, nil
}

// FetchLogsByAddress получает логи по адресу
func (r *LogRepository) FetchLogsByAddress(table string, address string, limit int) ([]models.Log, error) {
	ctx := context.Background()

	var result []struct {
		BlockNumber      uint64    `ch:"block_number"`
		BlockHash        string    `ch:"block_hash"`
		TransactionHash  string    `ch:"transaction_hash"`
		TransactionIndex uint32    `ch:"transaction_index"`
		LogIndex         uint32    `ch:"log_index"`
		Address          string    `ch:"address"`
		Data             string    `ch:"data"`
		Topics           []string  `ch:"topics"`
		BlockTimestamp   time.Time `ch:"block_timestamp"`
		Topic0           string    `ch:"topic0"`
	}

	query := "SELECT * FROM " + table + " WHERE address = ? ORDER BY block_timestamp DESC LIMIT ?"
	err := r.Client.Select(ctx, &result, query, address, limit)
	if err != nil {
		r.Logger.Errorf("Failed to fetch logs for address %s: %v", address, err)
		return nil, err
	}

	// Конвертируем результаты в модели Log
	logs := make([]models.Log, len(result))
	for i, row := range result {
		logs[i] = models.Log{
			Address:         row.Address,
			Topics:          row.Topics,
			Data:            row.Data,
			TransactionHash: row.TransactionHash,
			LogIndex:        uint64(row.LogIndex),
			Removed:         false, // По умолчанию false
		}
	}

	r.Logger.Debugf("Successfully fetched %d logs for address %s", len(logs), address)
	return logs, nil
}

// FetchLogsByTopic получает логи по топику
func (r *LogRepository) FetchLogsByTopic(table string, topic string, limit int) ([]models.Log, error) {
	ctx := context.Background()

	var result []struct {
		BlockNumber      uint64    `ch:"block_number"`
		BlockHash        string    `ch:"block_hash"`
		TransactionHash  string    `ch:"transaction_hash"`
		TransactionIndex uint32    `ch:"transaction_index"`
		LogIndex         uint32    `ch:"log_index"`
		Address          string    `ch:"address"`
		Data             string    `ch:"data"`
		Topics           []string  `ch:"topics"`
		BlockTimestamp   time.Time `ch:"block_timestamp"`
		Topic0           string    `ch:"topic0"`
	}

	query := "SELECT * FROM " + table + " WHERE has(topics, ?) ORDER BY block_timestamp DESC LIMIT ?"
	err := r.Client.Select(ctx, &result, query, topic, limit)
	if err != nil {
		r.Logger.Errorf("Failed to fetch logs for topic %s: %v", topic, err)
		return nil, err
	}

	// Конвертируем результаты в модели Log
	logs := make([]models.Log, len(result))
	for i, row := range result {
		logs[i] = models.Log{
			Address:         row.Address,
			Topics:          row.Topics,
			Data:            row.Data,
			TransactionHash: row.TransactionHash,
			LogIndex:        uint64(row.LogIndex),
			Removed:         false, // По умолчанию false
		}
	}

	r.Logger.Debugf("Successfully fetched %d logs for topic %s", len(logs), topic)
	return logs, nil
}

// FetchLogsByTopic0 получает логи по первому топику (topic0)
func (r *LogRepository) FetchLogsByTopic0(table string, topic0 string, limit int) ([]models.Log, error) {
	ctx := context.Background()

	var result []struct {
		BlockNumber      uint64    `ch:"block_number"`
		BlockHash        string    `ch:"block_hash"`
		TransactionHash  string    `ch:"transaction_hash"`
		TransactionIndex uint32    `ch:"transaction_index"`
		LogIndex         uint32    `ch:"log_index"`
		Address          string    `ch:"address"`
		Data             string    `ch:"data"`
		Topics           []string  `ch:"topics"`
		BlockTimestamp   time.Time `ch:"block_timestamp"`
		Topic0           string    `ch:"topic0"`
	}

	query := "SELECT * FROM " + table + " WHERE topic0 = ? ORDER BY block_timestamp DESC LIMIT ?"
	err := r.Client.Select(ctx, &result, query, topic0, limit)
	if err != nil {
		r.Logger.Errorf("Failed to fetch logs for topic0 %s: %v", topic0, err)
		return nil, err
	}

	// Конвертируем результаты в модели Log
	logs := make([]models.Log, len(result))
	for i, row := range result {
		logs[i] = models.Log{
			Address:         row.Address,
			Topics:          row.Topics,
			Data:            row.Data,
			TransactionHash: row.TransactionHash,
			LogIndex:        uint64(row.LogIndex),
			Removed:         false, // По умолчанию false
		}
	}

	r.Logger.Debugf("Successfully fetched %d logs for topic0 %s", len(logs), topic0)
	return logs, nil
}

// FetchLogsByAddressAndTopic получает логи по адресу и топику
func (r *LogRepository) FetchLogsByAddressAndTopic(table string, address string, topic string, limit int) ([]models.Log, error) {
	ctx := context.Background()

	var result []struct {
		BlockNumber      uint64    `ch:"block_number"`
		BlockHash        string    `ch:"block_hash"`
		TransactionHash  string    `ch:"transaction_hash"`
		TransactionIndex uint32    `ch:"transaction_index"`
		LogIndex         uint32    `ch:"log_index"`
		Address          string    `ch:"address"`
		Data             string    `ch:"data"`
		Topics           []string  `ch:"topics"`
		BlockTimestamp   time.Time `ch:"block_timestamp"`
		Topic0           string    `ch:"topic0"`
	}

	query := "SELECT * FROM " + table + " WHERE address = ? AND has(topics, ?) ORDER BY block_timestamp DESC LIMIT ?"
	err := r.Client.Select(ctx, &result, query, address, topic, limit)
	if err != nil {
		r.Logger.Errorf("Failed to fetch logs for address %s and topic %s: %v", address, topic, err)
		return nil, err
	}

	// Конвертируем результаты в модели Log
	logs := make([]models.Log, len(result))
	for i, row := range result {
		logs[i] = models.Log{
			Address:         row.Address,
			Topics:          row.Topics,
			Data:            row.Data,
			TransactionHash: row.TransactionHash,
			LogIndex:        uint64(row.LogIndex),
			Removed:         false, // По умолчанию false
		}
	}

	r.Logger.Debugf("Successfully fetched %d logs for address %s and topic %s", len(logs), address, topic)
	return logs, nil
}
