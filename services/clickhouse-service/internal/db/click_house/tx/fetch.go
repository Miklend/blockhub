package tx

import (
	"context"

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
	tx := txRowToModel(row)

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
		txs[i] = txRowToModel(row)
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
		txs[i] = txRowToModel(row)
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
		txs[i] = txRowToModel(row)
	}

	r.Logger.Debugf("Successfully fetched %d transactions for block number %d", len(txs), blockNumber)
	return txs, nil
}

// FetchTxsByAddress получает транзакции по адресу
func (r *TxRepository) FetchTxsByAddress(table string, address string, limit int) ([]models.Tx, error) {
	ctx := context.Background()

	var result []rowtypes.TxRow

	query := "SELECT * FROM " + table + " WHERE from = ? OR to = ? ORDER BY block_timestamp DESC LIMIT ?"
	err := r.Client.Select(ctx, &result, query, address, address, limit)
	if err != nil {
		r.Logger.Errorf("Failed to fetch transactions by address %s: %v", address, err)
		return nil, err
	}

	// Конвертируем результаты в модели Tx (аналогично FetchTxs)
	txs := make([]models.Tx, len(result))
	for i, row := range result {
		txs[i] = txRowToModel(row)
	}

	r.Logger.Debugf("Successfully fetched %d transactions for address %s", len(txs), address)
	return txs, nil
}

func txRowToModel(row rowtypes.TxRow) models.Tx {
	tx := models.Tx{
		Hash:             row.Hash,
		BlockHash:        row.BlockHash,
		BlockNumber:      uint(row.BlockNumber),
		TransactionIndex: uint(row.TransactionIndex),
		From:             row.From,
		Value:            prefixHex(row.Value),
		Gas:              uint(row.Gas),
		GasPrice:         uint(row.GasPrice),
		Input:            row.Input,
		Nonce:            uint(row.Nonce),
		Type:             uint(row.Type),
		ChainID:          uint(row.ChainID),
		V:                row.V,
		R:                row.R,
		S:                row.S,
		AccessList:       row.AccessList,
		BlockTimestamp:   row.BlockTimestamp,
	}

	if row.To != nil {
		tx.To = row.To
	}

	tx.MaxFeePerGas = uintPtrFromUint64(row.MaxFeePerGas)
	tx.MaxPriorityFeePerGas = uintPtrFromUint64(row.MaxPriorityFeePerGas)

	return tx
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

func uintPtrFromUint64(src *uint64) *uint {
	if src == nil {
		return nil
	}
	val := uint(*src)
	return &val
}
