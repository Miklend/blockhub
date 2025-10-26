package clickhouseRepo

import "lib/models"

func (c *ClickhouseRepo) InsertTx(table string, tx models.Tx) error {
	c.logger.Debugf("InsertTx: table=%s hash=%s", table, tx.Hash)
	// TODO: реализовать вставку одной транзакции
	return nil
}

func (c *ClickhouseRepo) InsertTxs(table string, txs []models.Tx) error {
	c.logger.Debugf("InsertTxs: table=%s count=%d", table, len(txs))
	// TODO: реализовать batch вставку
	return nil
}

func (c *ClickhouseRepo) FetchTx(table string, hashTx string) (models.Tx, error) {
	c.logger.Debugf("FetchTx: table=%s hash=%s", table, hashTx)
	// TODO: реализовать SELECT по hash
	return models.Tx{}, nil
}

func (c *ClickhouseRepo) FetchTxs(table string, hashTxs []string) ([]models.Tx, error) {
	c.logger.Debugf("FetchTxs: table=%s count=%d", table, len(hashTxs))
	// TODO: реализовать SELECT по множеству hash
	return nil, nil
}
