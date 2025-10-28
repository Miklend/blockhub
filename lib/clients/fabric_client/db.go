package fabricClient

import (
	"context"
	"fmt"
	clickhouseClient "lib/clients/db/clickhouse"
	redisClient "lib/clients/db/redis"
	"lib/models"
)

const (
	dbTypeClickhouse = "clickhouse"
	dbTypeRedis      = "redis"
)

func NewDB(ctx context.Context, dbType string, cfg models.DB) (any, error) {
	switch dbType {
	case dbTypeClickhouse:
		return clickhouseClient.NewClient(ctx, models.Clickhouse{DB: cfg})
	case dbTypeRedis:
		return redisClient.NewClient(ctx, models.Redis{DB: cfg})
	default:
		return nil, fmt.Errorf("not found db client for type: %s", dbType)
	}
}
