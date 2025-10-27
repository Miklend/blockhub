package receipt

import (
	"context"
	"strconv"
	"time"

	clientsDB "lib/clients/db"
	"lib/models"
	"lib/utils/logging"
)

type ReceiptRepository struct {
	Client clientsDB.ClickhouseClient
	Logger *logging.Logger
}

func NewReceiptRepository(client clientsDB.ClickhouseClient, logger *logging.Logger) *ReceiptRepository {
	return &ReceiptRepository{
		Client: client,
		Logger: logger,
	}
}

// InsertReceipt вставляет одну квитанцию в таблицу
func (r *ReceiptRepository) InsertReceipt(table string, receipt models.Receipt, txHash string, txIndex uint32, blockHash string, blockNumber uint64, blockTimestamp uint64) error {
	ctx := context.Background()

	row := convertReceiptToClickHouseRow(receipt, txHash, txIndex, blockHash, blockNumber, blockTimestamp)

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for receipt insert: %v", err)
		return err
	}

	err = batch.Append(row...)
	if err != nil {
		r.Logger.Errorf("Failed to append receipt to batch: %v", err)
		return err
	}

	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for receipt insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted receipt for transaction %s", txHash)
	return nil
}

// InsertReceipts вставляет массив квитанций в таблицу
func (r *ReceiptRepository) InsertReceipts(table string, receipts []models.Receipt, txHashes []string, txIndexes []uint32, blockHash string, blockNumber uint64, blockTimestamp uint64) error {
	if len(receipts) == 0 {
		return nil
	}

	ctx := context.Background()

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for receipts insert: %v", err)
		return err
	}

	for i, receipt := range receipts {
		row := convertReceiptToClickHouseRow(receipt, txHashes[i], txIndexes[i], blockHash, blockNumber, blockTimestamp)
		err = batch.Append(row...)
		if err != nil {
			r.Logger.Errorf("Failed to append receipt %d to batch: %v", i, err)
			return err
		}
	}

	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for receipts insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted %d receipts", len(receipts))
	return nil
}

// InsertReceiptsFromBlock вставляет квитанции из блока
func (r *ReceiptRepository) InsertReceiptsFromBlock(table string, block models.Block) error {
	if len(block.Transactions) == 0 {
		return nil
	}

	ctx := context.Background()

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for receipts from block insert: %v", err)
		return err
	}

	for i, tx := range block.Transactions {
		if tx.Receipt != nil {
			row := convertReceiptToClickHouseRow(*tx.Receipt, tx.Hash, uint32(i), block.Hash, block.Number, block.Timestamp)
			err = batch.Append(row...)
			if err != nil {
				r.Logger.Errorf("Failed to append receipt for transaction %s to batch: %v", tx.Hash, err)
				return err
			}
		}
	}

	err = batch.Send()
	if err != nil {
		r.Logger.Errorf("Failed to send batch for receipts from block insert: %v", err)
		return err
	}

	r.Logger.Debugf("Successfully inserted receipts from block %s", block.Hash)
	return nil
}

// convertReceiptToClickHouseRow конвертирует Receipt в строку для вставки в ClickHouse
func convertReceiptToClickHouseRow(receipt models.Receipt, txHash string, txIndex uint32, blockHash string, blockNumber uint64, blockTimestamp uint64) []interface{} {
	timestamp := time.Unix(int64(blockTimestamp), 0)

	// Конвертируем to в указатель
	var to *string
	if receipt.To != "" {
		to = &receipt.To
	}

	// Конвертируем contractAddress в указатель
	var contractAddress *string
	if receipt.ContractAddress != "" {
		contractAddress = &receipt.ContractAddress
	}

	// Конвертируем effectiveGasPrice
	effectiveGasPrice, _ := parseHexToUint64(receipt.EffectiveGasPrice)

	return []interface{}{
		txHash,                    // transaction_hash
		txIndex,                   // transaction_index
		blockHash,                 // block_hash
		blockNumber,               // block_number
		receipt.From,              // from
		to,                        // to
		contractAddress,           // contract_address
		receipt.CumulativeGasUsed, // cumulative_gas_used
		receipt.GasUsed,           // gas_used
		effectiveGasPrice,         // effective_gas_price
		uint8(receipt.Status),     // status
		receipt.LogsBloom,         // logs_bloom
		timestamp,                 // block_timestamp
		timestamp,                 // date (MATERIALIZED)
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
