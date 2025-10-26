package clickhouseRepo

import (
	"clickhouse-service/internal/db"
	clientsDB "lib/clients/db"
	"lib/utils/logging"
)

type ClickhouseRepo struct {
	client clientsDB.ClickhouseClient
	logger *logging.Logger
}

func NewClickhouseService(client clientsDB.ClickhouseClient, logger *logging.Logger) db.DB {
	return &ClickhouseRepo{
		client: client,
		logger: logger,
	}
}

func (c *ClickhouseRepo) Close() error {
	return c.client.Close()
}
