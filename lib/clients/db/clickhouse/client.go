package clickhouseClient

import (
	"context"
	"fmt"
	"lib/clients/db"
	"lib/models"
	"lib/utils/logging"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type clickhouseClient struct {
	conn   clickhouse.Conn
	logger *logging.Logger
}

func NewClient(ctx context.Context, cfg models.Clickhouse, logger *logging.Logger) (db.ClickhouseClient, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	// Проверяем соединение
	if err := conn.Ping(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	return &clickhouseClient{
		conn:   conn,
		logger: logger,
	}, nil
}
