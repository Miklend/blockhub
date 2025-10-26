package repository

import "lib/models"

type Storage interface {
	Close() error

	InsertBlock(table string, block models.Block) error
	InsertBlocks(table string, block []models.Block) error
	FetchBlock(table string, hashBlock string) (models.Block, error)
	FetchBlocks(table string, hashBlocks []string) ([]models.Block, error)

	InsertTx(table string, tx models.Tx) error
	InsertTxs(table string, txs []models.Tx) error
	FetchTx(table string, hashTx string) (models.Tx, error)
	FetchTxs(table string, hashTxs []string) ([]models.Tx, error)

	InsertReceipt(table string, receipt models.Receipt) error
	InsertReceipts(table string, receipts []models.Receipt) error
	FetchReceipt(table string, hashReceipt string) (models.Receipt, error)
	FetchReceipts(table string, hashReceipts []string) ([]models.Receipt, error)

	InsertLog(table string, log models.Log) error
	InsertLogs(table string, logs []models.Log) error
	FetchLog(table string, hashLog string) (models.Log, error)
	FetchLogs(table string, hashLogs []string) ([]models.Log, error)
}
