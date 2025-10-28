package tx

import (
	"context"
	"strconv"
	"time"

	"clickhouse-service/internal/db/click_house/rowtypes"
	"lib/models"
)

// FetchTx получает транзакцию по хешу
func (r *TxRepository) FetchTx(table string, txHash string) (models.Tx, error) {
	ctx := context.Background()

	var result []rowtypes.TxRow

	query := "SELECT * FROM " + table + " WHERE hash = ? LIMIT 1"
	err := r.Client.Select(ctx, &result, query, txHash)
	if err != nil {
		r.Logger.Errorf("Failed to fetch transaction %s: %v", txHash, err)
		return models.Tx{}, err
	}

	if len(result) == 0 {
		return models.Tx{}, nil
	}

	// Конвертируем результат в модель Tx
	row := result[0]
	tx := models.Tx{
		Hash:             row.Hash,
		From:             row.From,
		Gas:              row.Gas,
		GasPrice:         formatUint64ToHex(row.GasPrice),
		Input:            row.Input,
		Nonce:            row.Nonce,
		TransactionIndex: uint64(row.TransactionIndex),
		Value:            "0x" + row.Value,
		Type:             row.Type,
		V:                row.V,
		R:                row.R,
		S:                row.S,
	}

	// Конвертируем to если есть
	if row.To != nil {
		tx.To = *row.To
	}

	// Конвертируем maxFeePerGas если есть
	if row.MaxFeePerGas != nil {
		tx.MaxFeePerGas = formatUint64ToHex(*row.MaxFeePerGas)
	}

	// Конвертируем maxPriorityFeePerGas если есть
	if row.MaxPriorityFeePerGas != nil {
		tx.MaxPriorityFeePerGas = formatUint64ToHex(*row.MaxPriorityFeePerGas)
	}

	// Конвертируем chainID
	if row.ChainID > 0 {
		tx.ChainID = formatUint64ToHex(row.ChainID)
	}

	r.Logger.Debugf("Successfully fetched transaction %s", tx.Hash)
	return tx, nil
}

// FetchTxs получает транзакции по хешам
func (r *TxRepository) FetchTxs(table string, txHashes []string) ([]models.Tx, error) {
	if len(txHashes) == 0 {
		return []models.Tx{}, nil
	}

	ctx := context.Background()

	var result []rowtypes.TxRow

	query := "SELECT * FROM " + table + " WHERE hash IN (?)"
	err := r.Client.Select(ctx, &result, query, txHashes)
	if err != nil {
		r.Logger.Errorf("Failed to fetch transactions: %v", err)
		return nil, err
	}

	// Конвертируем результаты в модели Tx
	txs := make([]models.Tx, len(result))
	for i, row := range result {
		tx := models.Tx{
			Hash:             row.Hash,
			From:             row.From,
			Gas:              row.Gas,
			GasPrice:         formatUint64ToHex(row.GasPrice),
			Input:            row.Input,
			Nonce:            row.Nonce,
			TransactionIndex: uint64(row.TransactionIndex),
			Value:            "0x" + row.Value,
			Type:             row.Type,
			V:                row.V,
			R:                row.R,
			S:                row.S,
		}

		// Конвертируем to если есть
		if row.To != nil {
			tx.To = *row.To
		}

		// Конвертируем maxFeePerGas если есть
		if row.MaxFeePerGas != nil {
			tx.MaxFeePerGas = formatUint64ToHex(*row.MaxFeePerGas)
		}

		// Конвертируем maxPriorityFeePerGas если есть
		if row.MaxPriorityFeePerGas != nil {
			tx.MaxPriorityFeePerGas = formatUint64ToHex(*row.MaxPriorityFeePerGas)
		}

		// Конвертируем chainID
		if row.ChainID > 0 {
			tx.ChainID = formatUint64ToHex(row.ChainID)
		}

		txs[i] = tx
	}

	r.Logger.Debugf("Successfully fetched %d transactions", len(txs))
	return txs, nil
}

// FetchTxsByBlock получает транзакции по хешу блока
func (r *TxRepository) FetchTxsByBlock(table string, blockHash string) ([]models.Tx, error) {
	ctx := context.Background()

	var result []rowtypes.TxRow

	query := "SELECT * FROM " + table + " WHERE block_hash = ? ORDER BY transaction_index"
	err := r.Client.Select(ctx, &result, query, blockHash)
	if err != nil {
		r.Logger.Errorf("Failed to fetch transactions by block %s: %v", blockHash, err)
		return nil, err
	}

	// Конвертируем результаты в модели Tx (аналогично FetchTxs)
	txs := make([]models.Tx, len(result))
	for i, row := range result {
		tx := models.Tx{
			Hash:             row.Hash,
			From:             row.From,
			Gas:              row.Gas,
			GasPrice:         formatUint64ToHex(row.GasPrice),
			Input:            row.Input,
			Nonce:            row.Nonce,
			TransactionIndex: uint64(row.TransactionIndex),
			Value:            "0x" + row.Value,
			Type:             row.Type,
			V:                row.V,
			R:                row.R,
			S:                row.S,
		}

		if row.To != nil {
			tx.To = *row.To
		}

		if row.MaxFeePerGas != nil {
			tx.MaxFeePerGas = formatUint64ToHex(*row.MaxFeePerGas)
		}

		if row.MaxPriorityFeePerGas != nil {
			tx.MaxPriorityFeePerGas = formatUint64ToHex(*row.MaxPriorityFeePerGas)
		}

		if row.ChainID > 0 {
			tx.ChainID = formatUint64ToHex(row.ChainID)
		}

		txs[i] = tx
	}

	r.Logger.Debugf("Successfully fetched %d transactions for block %s", len(txs), blockHash)
	return txs, nil
}

// FetchTxsByBlockNumber получает транзакции по номеру блока
func (r *TxRepository) FetchTxsByBlockNumber(table string, blockNumber uint64) ([]models.Tx, error) {
	ctx := context.Background()

	var result []rowtypes.TxRow

	query := "SELECT * FROM " + table + " WHERE block_number = ? ORDER BY transaction_index"
	err := r.Client.Select(ctx, &result, query, blockNumber)
	if err != nil {
		r.Logger.Errorf("Failed to fetch transactions by block number %d: %v", blockNumber, err)
		return nil, err
	}

	// Конвертируем результаты в модели Tx (аналогично FetchTxs)
	txs := make([]models.Tx, len(result))
	for i, row := range result {
		tx := models.Tx{
			Hash:             row.Hash,
			From:             row.From,
			Gas:              row.Gas,
			GasPrice:         formatUint64ToHex(row.GasPrice),
			Input:            row.Input,
			Nonce:            row.Nonce,
			TransactionIndex: uint64(row.TransactionIndex),
			Value:            "0x" + row.Value,
			Type:             row.Type,
			V:                row.V,
			R:                row.R,
			S:                row.S,
		}

		if row.To != nil {
			tx.To = *row.To
		}

		if row.MaxFeePerGas != nil {
			tx.MaxFeePerGas = formatUint64ToHex(*row.MaxFeePerGas)
		}

		if row.MaxPriorityFeePerGas != nil {
			tx.MaxPriorityFeePerGas = formatUint64ToHex(*row.MaxPriorityFeePerGas)
		}

		if row.ChainID > 0 {
			tx.ChainID = formatUint64ToHex(row.ChainID)
		}

		txs[i] = tx
	}

	r.Logger.Debugf("Successfully fetched %d transactions for block number %d", len(txs), blockNumber)
	return txs, nil
}

// FetchTxsByAddress получает транзакции по адресу
func (r *TxRepository) FetchTxsByAddress(table string, address string, limit int) ([]models.Tx, error) {
	ctx := context.Background()

	var result []struct {
		Hash                 string    `ch:"hash"`
		BlockHash            string    `ch:"block_hash"`
		BlockNumber          uint64    `ch:"block_number"`
		TransactionIndex     uint32    `ch:"transaction_index"`
		From                 string    `ch:"from"`
		To                   *string   `ch:"to"`
		Value                string    `ch:"value"`
		Gas                  uint64    `ch:"gas"`
		GasPrice             uint64    `ch:"gas_price"`
		Input                string    `ch:"input"`
		Nonce                uint64    `ch:"nonce"`
		Type                 uint8     `ch:"type"`
		MaxFeePerGas         *uint64   `ch:"max_fee_per_gas"`
		MaxPriorityFeePerGas *uint64   `ch:"max_priority_fee_per_gas"`
		ChainID              uint64    `ch:"chain_id"`
		V                    string    `ch:"v"`
		R                    string    `ch:"r"`
		S                    string    `ch:"s"`
		AccessList           string    `ch:"access_list"`
		BlockTimestamp       time.Time `ch:"block_timestamp"`
	}

	query := "SELECT * FROM " + table + " WHERE from = ? OR to = ? ORDER BY block_timestamp DESC LIMIT ?"
	err := r.Client.Select(ctx, &result, query, address, address, limit)
	if err != nil {
		r.Logger.Errorf("Failed to fetch transactions by address %s: %v", address, err)
		return nil, err
	}

	// Конвертируем результаты в модели Tx (аналогично FetchTxs)
	txs := make([]models.Tx, len(result))
	for i, row := range result {
		tx := models.Tx{
			Hash:             row.Hash,
			From:             row.From,
			Gas:              row.Gas,
			GasPrice:         formatUint64ToHex(row.GasPrice),
			Input:            row.Input,
			Nonce:            row.Nonce,
			TransactionIndex: uint64(row.TransactionIndex),
			Value:            "0x" + row.Value,
			Type:             row.Type,
			V:                row.V,
			R:                row.R,
			S:                row.S,
		}

		if row.To != nil {
			tx.To = *row.To
		}

		if row.MaxFeePerGas != nil {
			tx.MaxFeePerGas = formatUint64ToHex(*row.MaxFeePerGas)
		}

		if row.MaxPriorityFeePerGas != nil {
			tx.MaxPriorityFeePerGas = formatUint64ToHex(*row.MaxPriorityFeePerGas)
		}

		if row.ChainID > 0 {
			tx.ChainID = formatUint64ToHex(row.ChainID)
		}

		txs[i] = tx
	}

	r.Logger.Debugf("Successfully fetched %d transactions for address %s", len(txs), address)
	return txs, nil
}

// formatUint64ToHex конвертирует uint64 в hex строку
func formatUint64ToHex(n uint64) string {
	return "0x" + strconv.FormatUint(n, 16)
}
