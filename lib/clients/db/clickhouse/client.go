package clickhouseClient

import (
	"context"
	"fmt"
	clientsDB "lib/clients/db"
	"lib/models"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func NewClient(ctx context.Context, cfg models.Clickhouse) (clientsDB.ClickhouseClient, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
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

	return conn, nil
}
