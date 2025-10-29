package db

import "lib/models"

type DB interface {
	Close() error

	// Блоки
	InsertBlock(table string, block models.Block) error
	InsertBlocks(table string, blocks []models.Block) error
	FetchBlock(table string, hashBlock string) (models.Block, error)
	FetchBlocks(table string, hashBlocks []string) ([]models.Block, error)
	FetchBlockByNumber(table string, blockNumber uint64) (models.Block, error)
	FetchBlocksByRange(table string, fromBlock, toBlock uint64) ([]models.Block, error)

	// Транзакции
	InsertTx(table string, tx models.Tx) error
	InsertTxs(table string, txs []models.Tx) error
	InsertTxWithBlockData(table string, tx models.Tx) error
	InsertTxsWithBlockData(table string, txs []models.Tx) error
	FetchTx(table string, txHash string) (models.Tx, error)
	FetchTxs(table string, txHashes []string) ([]models.Tx, error)
	FetchTxsByBlock(table string, blockHash string) ([]models.Tx, error)
	FetchTxsByBlockNumber(table string, blockNumber uint64) ([]models.Tx, error)
	FetchTxsByAddress(table string, address string, limit int) ([]models.Tx, error)

	// Квитанции
	InsertReceipt(table string, receipt models.Receipt) error
	InsertReceipts(table string, receipts []models.Receipt) error
	InsertReceiptsFromBlock(table string, block models.Block) error
	FetchReceipt(table string, txHash string) (models.Receipt, error)
	FetchReceipts(table string, txHashes []string) ([]models.Receipt, error)
	FetchReceiptsByBlock(table string, blockHash string) ([]models.Receipt, error)
	FetchReceiptsByBlockNumber(table string, blockNumber uint64) ([]models.Receipt, error)
	FetchReceiptsByAddress(table string, address string, limit int) ([]models.Receipt, error)

	// Логи
	InsertLog(table string, log models.Log) error
	InsertLogs(table string, logs []models.Log) error
	InsertLogsFromReceipt(table string, receipt models.Receipt) error
	InsertLogsFromBlock(table string, block models.Block) error
	FetchLogsByTransaction(table string, txHash string) ([]models.Log, error)
	FetchLogsByBlock(table string, blockHash string) ([]models.Log, error)
	FetchLogsByBlockNumber(table string, blockNumber uint64) ([]models.Log, error)
	FetchLogsByAddress(table string, address string, limit int) ([]models.Log, error)
	FetchLogsByTopic(table string, topic string, limit int) ([]models.Log, error)
	FetchLogsByTopic0(table string, topic0 string, limit int) ([]models.Log, error)
	FetchLogsByAddressAndTopic(table string, address string, topic string, limit int) ([]models.Log, error)
}
