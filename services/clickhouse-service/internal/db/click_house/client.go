package clickhouseRepo

import (
	"clickhouse-service/internal/db"
	"clickhouse-service/internal/db/click_house/block"
	"clickhouse-service/internal/db/click_house/tx"
	"clickhouse-service/internal/db/click_house/tx/log"
	"clickhouse-service/internal/db/click_house/tx/receipt"
	clientsDB "lib/clients/db"
	"lib/models"
	"lib/utils/logging"
)

type ClickhouseRepo struct {
	Client clientsDB.ClickhouseClient
	Logger *logging.Logger

	// Репозитории для каждой сущности
	BlockRepo   *block.BlockRepository
	TxRepo      *tx.TxRepository
	ReceiptRepo *receipt.ReceiptRepository
	LogRepo     *log.LogRepository
}

func NewClickhouseService(client clientsDB.ClickhouseClient, logger *logging.Logger) db.DB {
	repo := &ClickhouseRepo{
		Client: client,
		Logger: logger,
	}

	// Инициализируем репозитории
	repo.BlockRepo = block.NewBlockRepository(client, logger)
	repo.TxRepo = tx.NewTxRepository(client, logger)
	repo.ReceiptRepo = receipt.NewReceiptRepository(client, logger)
	repo.LogRepo = log.NewLogRepository(client, logger)

	return repo
}

func (c *ClickhouseRepo) Close() error {
	return c.Client.Close()
}

// Блоки

func (c *ClickhouseRepo) InsertBlock(table string, block models.Block) error {
	return c.BlockRepo.InsertBlock(table, block)
}

func (c *ClickhouseRepo) InsertBlocks(table string, blocks []models.Block) error {
	return c.BlockRepo.InsertBlocks(table, blocks)
}

func (c *ClickhouseRepo) FetchBlock(table string, hashBlock string) (models.Block, error) {
	return c.BlockRepo.FetchBlock(table, hashBlock)
}

func (c *ClickhouseRepo) FetchBlocks(table string, hashBlocks []string) ([]models.Block, error) {
	return c.BlockRepo.FetchBlocks(table, hashBlocks)
}

func (c *ClickhouseRepo) FetchBlockByNumber(table string, blockNumber uint64) (models.Block, error) {
	return c.BlockRepo.FetchBlockByNumber(table, blockNumber)
}

func (c *ClickhouseRepo) FetchBlocksByRange(table string, fromBlock, toBlock uint64) ([]models.Block, error) {
	return c.BlockRepo.FetchBlocksByRange(table, fromBlock, toBlock)
}

// Транзакции

func (c *ClickhouseRepo) InsertTx(table string, tx models.Tx) error {
	return c.TxRepo.InsertTx(table, tx)
}

func (c *ClickhouseRepo) InsertTxs(table string, txs []models.Tx) error {
	return c.TxRepo.InsertTxs(table, txs)
}

func (c *ClickhouseRepo) InsertTxWithBlockData(table string, tx models.Tx) error {
	return c.TxRepo.InsertTxWithBlockData(table, tx)
}

func (c *ClickhouseRepo) InsertTxsWithBlockData(table string, txs []models.Tx) error {
	return c.TxRepo.InsertTxsWithBlockData(table, txs)
}

func (c *ClickhouseRepo) FetchTx(table string, txHash string) (models.Tx, error) {
	return c.TxRepo.FetchTx(table, txHash)
}

func (c *ClickhouseRepo) FetchTxs(table string, txHashes []string) ([]models.Tx, error) {
	return c.TxRepo.FetchTxs(table, txHashes)
}

func (c *ClickhouseRepo) FetchTxsByBlock(table string, blockHash string) ([]models.Tx, error) {
	return c.TxRepo.FetchTxsByBlock(table, blockHash)
}

func (c *ClickhouseRepo) FetchTxsByBlockNumber(table string, blockNumber uint64) ([]models.Tx, error) {
	return c.TxRepo.FetchTxsByBlockNumber(table, blockNumber)
}

func (c *ClickhouseRepo) FetchTxsByAddress(table string, address string, limit int) ([]models.Tx, error) {
	return c.TxRepo.FetchTxsByAddress(table, address, limit)
}

// Квитанции

func (c *ClickhouseRepo) InsertReceipt(table string, receipt models.Receipt) error {
	return c.ReceiptRepo.InsertReceipt(table, receipt)
}

func (c *ClickhouseRepo) InsertReceipts(table string, receipts []models.Receipt) error {
	return c.ReceiptRepo.InsertReceipts(table, receipts)
}

func (c *ClickhouseRepo) InsertReceiptsFromBlock(table string, block models.Block) error {
	return c.ReceiptRepo.InsertReceiptsFromBlock(table, block)
}

func (c *ClickhouseRepo) FetchReceipt(table string, txHash string) (models.Receipt, error) {
	return c.ReceiptRepo.FetchReceipt(table, txHash)
}

func (c *ClickhouseRepo) FetchReceipts(table string, txHashes []string) ([]models.Receipt, error) {
	return c.ReceiptRepo.FetchReceipts(table, txHashes)
}

func (c *ClickhouseRepo) FetchReceiptsByBlock(table string, blockHash string) ([]models.Receipt, error) {
	return c.ReceiptRepo.FetchReceiptsByBlock(table, blockHash)
}

func (c *ClickhouseRepo) FetchReceiptsByBlockNumber(table string, blockNumber uint64) ([]models.Receipt, error) {
	return c.ReceiptRepo.FetchReceiptsByBlockNumber(table, blockNumber)
}

func (c *ClickhouseRepo) FetchReceiptsByAddress(table string, address string, limit int) ([]models.Receipt, error) {
	return c.ReceiptRepo.FetchReceiptsByAddress(table, address, limit)
}

// Логи

func (c *ClickhouseRepo) InsertLog(table string, log models.Log) error {
	return c.LogRepo.InsertLog(table, log)
}

func (c *ClickhouseRepo) InsertLogs(table string, logs []models.Log) error {
	return c.LogRepo.InsertLogs(table, logs)
}

func (c *ClickhouseRepo) InsertLogsFromReceipt(table string, receipt models.Receipt) error {
	return c.LogRepo.InsertLogsFromReceipt(table, receipt)
}

func (c *ClickhouseRepo) InsertLogsFromBlock(table string, block models.Block) error {
	return c.LogRepo.InsertLogsFromBlock(table, block)
}

func (c *ClickhouseRepo) FetchLogsByTransaction(table string, txHash string) ([]models.Log, error) {
	return c.LogRepo.FetchLogsByTransaction(table, txHash)
}

func (c *ClickhouseRepo) FetchLogsByBlock(table string, blockHash string) ([]models.Log, error) {
	return c.LogRepo.FetchLogsByBlock(table, blockHash)
}

func (c *ClickhouseRepo) FetchLogsByBlockNumber(table string, blockNumber uint64) ([]models.Log, error) {
	return c.LogRepo.FetchLogsByBlockNumber(table, blockNumber)
}

func (c *ClickhouseRepo) FetchLogsByAddress(table string, address string, limit int) ([]models.Log, error) {
	return c.LogRepo.FetchLogsByAddress(table, address, limit)
}

func (c *ClickhouseRepo) FetchLogsByTopic(table string, topic string, limit int) ([]models.Log, error) {
	return c.LogRepo.FetchLogsByTopic(table, topic, limit)
}

func (c *ClickhouseRepo) FetchLogsByTopic0(table string, topic0 string, limit int) ([]models.Log, error) {
	return c.LogRepo.FetchLogsByTopic0(table, topic0, limit)
}

func (c *ClickhouseRepo) FetchLogsByAddressAndTopic(table string, address string, topic string, limit int) ([]models.Log, error) {
	return c.LogRepo.FetchLogsByAddressAndTopic(table, address, topic, limit)
}
