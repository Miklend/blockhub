package receipt

import (
	"context"
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
// Примечание: требует, чтобы данные транзакции и блока были доступны
func (r *ReceiptRepository) InsertReceipt(table string, receipt models.Receipt) error {
	ctx := context.Background()

	// Для InsertReceipt данные транзакции и блока должны быть получены извне
	// Используем пустые значения, если данные недоступны (может привести к ошибкам)
	txHash := ""
	txIndex := uint32(0)
	blockHash := ""
	blockNumber := uint64(0)
	blockTimestamp := time.Time{}

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

	r.Logger.Debugf("Successfully inserted receipt")
	return nil
}

// InsertReceipts вставляет массив квитанций в таблицу
// Примечание: требует, чтобы данные транзакций и блоков были доступны
func (r *ReceiptRepository) InsertReceipts(table string, receipts []models.Receipt) error {
	if len(receipts) == 0 {
		return nil
	}

	ctx := context.Background()

	batch, err := r.Client.PrepareBatch(ctx, "INSERT INTO "+table+" VALUES")
	if err != nil {
		r.Logger.Errorf("Failed to prepare batch for receipts insert: %v", err)
		return err
	}

	// Для InsertReceipts данные транзакций и блоков должны быть получены извне
	// Используем пустые значения, если данные недоступны (может привести к ошибкам)
	blockHash := ""
	blockNumber := uint64(0)
	blockTimestamp := time.Time{}

	for i, receipt := range receipts {
		// Используем пустые значения для txHash и txIndex
		txHash := ""
		txIndex := uint32(i)
		row := convertReceiptToClickHouseRow(receipt, txHash, txIndex, blockHash, blockNumber, blockTimestamp)
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

	r.Logger.Warn("InsertReceiptsFromBlock: block.Transactions does not include receipt payloads; skipping insertion")
	return nil
}

// convertReceiptToClickHouseRow конвертирует Receipt в строку для вставки в ClickHouse
func convertReceiptToClickHouseRow(receipt models.Receipt, txHash string, txIndex uint32, blockHash string, blockNumber uint64, blockTimestamp time.Time) []interface{} {
	timestamp := blockTimestamp
	if timestamp.IsZero() {
		timestamp = receipt.BlockTimestamp
	}

	return []interface{}{
		txHash,
		txIndex,
		blockHash,
		blockNumber,
		receipt.From,
		receipt.To,
		receipt.ContractAddress,
		uint64(receipt.CumulativeGasUsed),
		uint64(receipt.GasUsed),
		uint64(receipt.EffectiveGasPrice),
		uint8(receipt.Status),
		receipt.LogsBloom,
		timestamp,
	}
}
