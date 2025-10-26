package clickhouseClient

import (
	"context"
	"fmt"
	clientsDB "lib/clients/db"
	"lib/models"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// Client wraps the ClickHouse connection to implement ClickhouseClient interface
type Client struct {
	conn driver.Conn
}

// NewClient creates a new ClickHouse client with the provided configuration
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

	return &Client{conn: conn}, nil
}

// Contributors returns the list of contributors
func (c *Client) Contributors() []string {
	return c.conn.Contributors()
}

// ServerVersion returns the server version information
func (c *Client) ServerVersion() (*clickhouse.ServerVersion, error) {
	return c.conn.ServerVersion()
}

// Select executes a query and scans the results into dest
func (c *Client) Select(ctx context.Context, dest any, query string, args ...any) error {
	return c.conn.Select(ctx, dest, query, args...)
}

// Query executes a query and returns the rows
func (c *Client) Query(ctx context.Context, query string, args ...any) (driver.Rows, error) {
	return c.conn.Query(ctx, query, args...)
}

// QueryRow executes a query and returns a single row
func (c *Client) QueryRow(ctx context.Context, query string, args ...any) driver.Row {
	return c.conn.QueryRow(ctx, query, args...)
}

// PrepareBatch prepares a batch for insert operations
func (c *Client) PrepareBatch(ctx context.Context, query string, opts ...driver.PrepareBatchOption) (driver.Batch, error) {
	return c.conn.PrepareBatch(ctx, query, opts...)
}

// Exec executes a query without returning results
func (c *Client) Exec(ctx context.Context, query string, args ...any) error {
	return c.conn.Exec(ctx, query, args...)
}

// AsyncInsert performs an async insert operation
func (c *Client) AsyncInsert(ctx context.Context, query string, wait bool, args ...any) error {
	return c.conn.AsyncInsert(ctx, query, wait, args...)
}

// Ping checks the connection to the server
func (c *Client) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

// Stats returns connection statistics
func (c *Client) Stats() driver.Stats {
	return c.conn.Stats()
}

// Close closes the connection
func (c *Client) Close() error {
	return c.conn.Close()
}