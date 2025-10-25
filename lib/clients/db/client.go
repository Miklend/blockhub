package db

import (
	"context"
)

// Client интерфейс для работы с различными БД
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row
	Begin(ctx context.Context) (Tx, error)
	SendBatch(ctx context.Context, b *Batch) BatchResults
	CopyFrom(ctx context.Context, table Identifier, columns []string, r CopyFromSource) (int64, error)
	Close() error
}

// CommonConfig общая конфигурация для подключения
type CommonConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// CommandTag представляет результат выполнения команды
type CommandTag interface {
	String() string
}

// Rows представляет набор строк результата запроса
type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
}

// Row представляет одну строку результата запроса
type Row interface {
	Scan(dest ...interface{}) error
}

// Tx представляет транзакцию
type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row
}

// Batch представляет пакет запросов
type Batch struct {
	Queries []BatchQuery
}

// BatchQuery представляет один запрос в пакете
type BatchQuery struct {
	SQL  string
	Args []interface{}
}

// BatchResults представляет результаты выполнения пакета
type BatchResults interface {
	Close() error
	Exec() (CommandTag, error)
	Query() (Rows, error)
	QueryRow() Row
}

// Identifier представляет идентификатор таблицы
type Identifier []string

// CopyFromSource представляет источник данных для COPY FROM
type CopyFromSource interface {
	Next() bool
	Values() ([]interface{}, error)
	Err() error
}
