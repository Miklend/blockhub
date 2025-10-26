package clickhouse

import (
	"clickhouse-service/internal/db"
	clientsDB "lib/clients/db"
	"lib/utils/logging"
)

type clickhouseClient struct {
	client clientsDB.ClickhouseClient
	logger *logging.Logger
}

func NewClickhouseService(client clientsDB.ClickhouseClient, logger *logging.Logger) db.DB {
	return &clickhouseClient{
		client: client,
		logger: logger,
	}
}

func (c *clickhouseClient) Close() error {
	return c.client.Close()
}
