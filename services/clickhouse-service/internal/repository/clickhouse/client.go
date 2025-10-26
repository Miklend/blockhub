package clickhouseRepo

import (
	"clickhouse-service/internal/repository"
	clientsDB "lib/clients/db"
	"lib/utils/logging"
)

type ClickhouseRepo struct {
	client clientsDB.ClickhouseClient
	logger *logging.Logger
}

func NewClickhouseService(client clientsDB.ClickhouseClient, logger *logging.Logger) repository.Storage {
	return &ClickhouseRepo{
		client: client,
		logger: logger,
	}
}

func (c *ClickhouseRepo) Close() error {
	return c.client.Close()
}
