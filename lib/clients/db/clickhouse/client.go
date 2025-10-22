package clickhouse

import (
	"context"
	"fmt"

	"lib/clients/db"
	"lib/models"
	"lib/utils/logging"

	"github.com/ClickHouse/clickhouse-go/v2"
)

// client реализует интерфейс SQLClient для ClickHouse
type client struct {
	conn   clickhouse.Conn
	logger *logging.Logger
}

// NewClient создает новый ClickHouse клиент
func NewClient(ctx context.Context, cfg models.Clickhouse, logger *logging.Logger) (db.SQLClient, error) {
	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

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

	return &client{
		conn:   conn,
		logger: logger,
	}, nil
}

// Ping проверяет соединение с базой данных
func (c *client) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

// Close закрывает соединение
func (c *client) Close() error {
	return c.conn.Close()
}

// GetType возвращает тип базы данных
func (c *client) GetType() db.DatabaseType {
	return db.ClickHouse
}

// Exec выполняет SQL команду без возврата данных
func (c *client) Exec(ctx context.Context, query string, args ...any) error {
	return c.conn.Exec(ctx, query, args...)
}

// Query выполняет SQL запрос и возвращает строки
func (c *client) Query(ctx context.Context, query string, args ...any) (db.Rows, error) {
	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &rowsWrapper{rows: rows}, nil
}

// QueryRow выполняет SQL запрос и возвращает одну строку
func (c *client) QueryRow(ctx context.Context, query string, args ...any) db.Row {
	row := c.conn.QueryRow(ctx, query, args...)
	return &rowWrapper{row: row}
}

// Begin начинает транзакцию
func (c *client) Begin(ctx context.Context) (db.Transaction, error) {
	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &transactionWrapper{tx: tx}, nil
}

// PrepareBatch подготавливает батч для вставки данных
func (c *client) PrepareBatch(ctx context.Context, query string) (db.Batch, error) {
	batch, err := c.conn.PrepareBatch(ctx, query)
	if err != nil {
		return nil, err
	}
	return &batchWrapper{batch: batch}, nil
}

// rowsWrapper оборачивает clickhouse.Rows для соответствия интерфейсу db.Rows
type rowsWrapper struct {
	rows clickhouse.Rows
}

func (r *rowsWrapper) Next() bool {
	return r.rows.Next()
}

func (r *rowsWrapper) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *rowsWrapper) Err() error {
	return r.rows.Err()
}

func (r *rowsWrapper) Close() error {
	return r.rows.Close()
}

// rowWrapper оборачивает clickhouse.Row для соответствия интерфейсу db.Row
type rowWrapper struct {
	row clickhouse.Row
}

func (r *rowWrapper) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

// transactionWrapper оборачивает clickhouse.Tx для соответствия интерфейсу db.Transaction
type transactionWrapper struct {
	tx clickhouse.Tx
}

func (t *transactionWrapper) Exec(ctx context.Context, query string, args ...any) error {
	return t.tx.Exec(ctx, query, args...)
}

func (t *transactionWrapper) Query(ctx context.Context, query string, args ...any) (db.Rows, error) {
	rows, err := t.tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &rowsWrapper{rows: rows}, nil
}

func (t *transactionWrapper) QueryRow(ctx context.Context, query string, args ...any) db.Row {
	row := t.tx.QueryRow(ctx, query, args...)
	return &rowWrapper{row: row}
}

func (t *transactionWrapper) Commit(ctx context.Context) error {
	return t.tx.Commit()
}

func (t *transactionWrapper) Rollback(ctx context.Context) error {
	return t.tx.Rollback()
}

// batchWrapper оборачивает clickhouse.Batch для соответствия интерфейсу db.Batch
type batchWrapper struct {
	batch clickhouse.Batch
}

func (b *batchWrapper) Append(args ...any) error {
	return b.batch.Append(args...)
}

func (b *batchWrapper) Send() error {
	return b.batch.Send()
}