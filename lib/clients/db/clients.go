package clientsDB

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/redis/go-redis/v9"
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
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	TTL(ctx context.Context, key string) *redis.DurationCmd
	SetNX(ctx context.Context, key string, value any, expiration time.Duration) *redis.BoolCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	Decr(ctx context.Context, key string) *redis.IntCmd
	Ping(ctx context.Context) *redis.StatusCmd
	Close() error
}
