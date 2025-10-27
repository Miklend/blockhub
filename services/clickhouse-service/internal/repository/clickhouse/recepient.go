package clickhouseRepo

import "lib/models"

func (c *ClickhouseRepo) InsertReceipt(table string, receipt models.Receipt) error {
	c.logger.Debugf("InsertReceipt: table=%s hash=%s", table, receipt.TxHash)
	// TODO: реализовать вставку
	return nil
}

func (c *ClickhouseRepo) InsertReceipts(table string, receipts []models.Receipt) error {
	c.logger.Debugf("InsertReceipts: table=%s count=%d", table, len(receipts))
	// TODO: реализовать batch вставку
	return nil
}

func (c *ClickhouseRepo) FetchReceipt(table string, hashReceipt string) (models.Receipt, error) {
	c.logger.Debugf("FetchReceipt: table=%s hash=%s", table, hashReceipt)
	// TODO: реализовать SELECT по hash
	return models.Receipt{}, nil
}

func (c *ClickhouseRepo) FetchReceipts(table string, hashReceipts []string) ([]models.Receipt, error) {
	c.logger.Debugf("FetchReceipts: table=%s count=%d", table, len(hashReceipts))
	// TODO: реализовать SELECT по множеству hash
	return nil, nil
}
