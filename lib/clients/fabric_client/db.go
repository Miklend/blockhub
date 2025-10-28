package fabricClient

import (
	"context"
	clientsDB "lib/clients/db"
	clickhouseClient "lib/clients/db/clickhouse"
	redisClient "lib/clients/db/redis"
	"lib/models"
)

// NewClickhouse создает клиента ClickHouse по конфигу
func NewClickhouse(ctx context.Context, cfg models.Clickhouse) (clientsDB.ClickhouseClient, error) {
	return clickhouseClient.NewClient(ctx, cfg)
}

// NewCash создает клиента Redis по конфигу
func NewCash(ctx context.Context, cfg models.Redis) (clientsDB.CashClient, error) {
	return redisClient.NewClient(ctx, cfg)
}
