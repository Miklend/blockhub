package clickhouseRepo

import "lib/models"

func (c *ClickhouseRepo) InsertLog(table string, log models.Log) error {
	c.logger.Debugf("InsertLog: table=%s", table)
	// TODO: реализовать вставку
	return nil
}

func (c *ClickhouseRepo) InsertLogs(table string, logs []models.Log) error {
	c.logger.Debugf("InsertLogs: table=%s count=%d", table, len(logs))
	// TODO: реализовать batch вставку
	return nil
}

func (c *ClickhouseRepo) FetchLog(table string, hashLog string) (models.Log, error) {
	c.logger.Debugf("FetchLog: table=%s hash=%s", table, hashLog)
	// TODO: реализовать SELECT по hash
	return models.Log{}, nil
}

func (c *ClickhouseRepo) FetchLogs(table string, hashLogs []string) ([]models.Log, error) {
	c.logger.Debugf("FetchLogs: table=%s count=%d", table, len(hashLogs))
	// TODO: реализовать SELECT по множеству hash
	return nil, nil
}
