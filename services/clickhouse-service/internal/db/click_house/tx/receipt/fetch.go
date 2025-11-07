package receipt

import (
	"context"

	"clickhouse-service/internal/db/click_house/rowtypes"
	"lib/models"
)

// FetchReceipt получает квитанцию по хешу транзакции
func (r *ReceiptRepository) FetchReceipt(table string, txHash string) (models.Receipt, error) {
	ctx := context.Background()

	var result []rowtypes.ReceiptRow

	query := "SELECT * FROM " + table + " WHERE transaction_hash = ? LIMIT 1"
	err := r.Client.Select(ctx, &result, query, txHash)
	if err != nil {
		r.Logger.Errorf("Failed to fetch receipt for transaction %s: %v", txHash, err)
		return models.Receipt{}, err
	}

	if len(result) == 0 {
		return models.Receipt{}, nil
	}

	// Конвертируем результат в модель Receipt
	row := result[0]
	receipt := receiptRowToModel(row)

	r.Logger.Debugf("Successfully fetched receipt for transaction %s", txHash)
	return receipt, nil
}

// FetchReceipts получает квитанции по хешам транзакций
func (r *ReceiptRepository) FetchReceipts(table string, txHashes []string) ([]models.Receipt, error) {
	if len(txHashes) == 0 {
		return []models.Receipt{}, nil
	}

	ctx := context.Background()

	var result []rowtypes.ReceiptRow

	query := "SELECT * FROM " + table + " WHERE transaction_hash IN (?)"
	err := r.Client.Select(ctx, &result, query, txHashes)
	if err != nil {
		r.Logger.Errorf("Failed to fetch receipts: %v", err)
		return nil, err
	}

	// Конвертируем результаты в модели Receipt
	receipts := make([]models.Receipt, len(result))
	for i, row := range result {
		receipts[i] = receiptRowToModel(row)
	}

	r.Logger.Debugf("Successfully fetched %d receipts", len(receipts))
	return receipts, nil
}

// FetchReceiptsByBlock получает квитанции по хешу блока
func (r *ReceiptRepository) FetchReceiptsByBlock(table string, blockHash string) ([]models.Receipt, error) {
	ctx := context.Background()

	var result []rowtypes.ReceiptRow

	query := "SELECT * FROM " + table + " WHERE block_hash = ? ORDER BY transaction_index"
	err := r.Client.Select(ctx, &result, query, blockHash)
	if err != nil {
		r.Logger.Errorf("Failed to fetch receipts by block %s: %v", blockHash, err)
		return nil, err
	}

	// Конвертируем результаты в модели Receipt (аналогично FetchReceipts)
	receipts := make([]models.Receipt, len(result))
	for i, row := range result {
		receipts[i] = receiptRowToModel(row)
	}

	r.Logger.Debugf("Successfully fetched %d receipts for block %s", len(receipts), blockHash)
	return receipts, nil
}

// FetchReceiptsByBlockNumber получает квитанции по номеру блока
func (r *ReceiptRepository) FetchReceiptsByBlockNumber(table string, blockNumber uint64) ([]models.Receipt, error) {
	ctx := context.Background()

	var result []rowtypes.ReceiptRow

	query := "SELECT * FROM " + table + " WHERE block_number = ? ORDER BY transaction_index"
	err := r.Client.Select(ctx, &result, query, blockNumber)
	if err != nil {
		r.Logger.Errorf("Failed to fetch receipts by block number %d: %v", blockNumber, err)
		return nil, err
	}

	// Конвертируем результаты в модели Receipt (аналогично FetchReceipts)
	receipts := make([]models.Receipt, len(result))
	for i, row := range result {
		receipts[i] = receiptRowToModel(row)
	}

	r.Logger.Debugf("Successfully fetched %d receipts for block number %d", len(receipts), blockNumber)
	return receipts, nil
}

// FetchReceiptsByAddress получает квитанции по адресу
func (r *ReceiptRepository) FetchReceiptsByAddress(table string, address string, limit int) ([]models.Receipt, error) {
	ctx := context.Background()

	var result []rowtypes.ReceiptRow

	query := "SELECT * FROM " + table + " WHERE from = ? OR to = ? OR contract_address = ? ORDER BY block_timestamp DESC LIMIT ?"
	err := r.Client.Select(ctx, &result, query, address, address, address, limit)
	if err != nil {
		r.Logger.Errorf("Failed to fetch receipts by address %s: %v", address, err)
		return nil, err
	}

	// Конвертируем результаты в модели Receipt (аналогично FetchReceipts)
	receipts := make([]models.Receipt, len(result))
	for i, row := range result {
		receipts[i] = receiptRowToModel(row)
	}

	r.Logger.Debugf("Successfully fetched %d receipts for address %s", len(receipts), address)
	return receipts, nil
}

func receiptRowToModel(row rowtypes.ReceiptRow) models.Receipt {
	receipt := models.Receipt{
		TransactionHash:   row.TransactionHash,
		TransactionIndex:  uint(row.TransactionIndex),
		BlockHash:         row.BlockHash,
		BlockNumber:       uint(row.BlockNumber),
		From:              row.From,
		CumulativeGasUsed: uint(row.CumulativeGasUsed),
		GasUsed:           uint(row.GasUsed),
		EffectiveGasPrice: uint(row.EffectiveGasPrice),
		Status:            uint(row.Status),
		LogsBloom:         row.LogsBloom,
		BlockTimestamp:    row.BlockTimestamp,
	}

	if row.To != nil {
		receipt.To = row.To
	}

	if row.ContractAddress != nil {
		receipt.ContractAddress = row.ContractAddress
	}

	return receipt
}
