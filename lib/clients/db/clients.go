package clientsDB

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickhouseClient interface {
	Contributors() []string
	ServerVersion() (*clickhouse.ServerVersion, error)
	Select(ctx context.Context, dest any, query string, args ...any) error
	Query(ctx context.Context, query string, args ...any) (driver.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) driver.Row
	PrepareBatch(ctx context.Context, query string, opts ...driver.PrepareBatchOption) (driver.Batch, error)
	Exec(ctx context.Context, query string, args ...any) error
	AsyncInsert(ctx context.Context, query string, wait bool, args ...any) error
	Ping(context.Context) error
	Stats() driver.Stats
	Close() error
}

type CashClient interface {
}
