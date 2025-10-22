package db

import "context"

// DatabaseType определяет тип базы данных
type DatabaseType string

const (
	PostgreSQL DatabaseType = "postgresql"
	ClickHouse DatabaseType = "clickhouse"
	Redis      DatabaseType = "redis"
)

// Client описывает универсальный интерфейс для работы с различными СУБД
type Client interface {
	// Ping проверяет соединение с базой данных
	Ping(ctx context.Context) error
	
	// Close закрывает соединение
	Close() error
	
	// GetType возвращает тип базы данных
	GetType() DatabaseType
}

// SQLClient описывает интерфейс для SQL-подобных хранилищ (PostgreSQL, ClickHouse)
type SQLClient interface {
	Client
	
	// Exec выполняет SQL команду без возврата данных
	Exec(ctx context.Context, query string, args ...any) error
	
	// Query выполняет SQL запрос и возвращает строки
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	
	// QueryRow выполняет SQL запрос и возвращает одну строку
	QueryRow(ctx context.Context, query string, args ...any) Row
	
	// Begin начинает транзакцию
	Begin(ctx context.Context) (Transaction, error)
	
	// PrepareBatch подготавливает батч для вставки данных
	PrepareBatch(ctx context.Context, query string) (Batch, error)
}

// NoSQLClient описывает интерфейс для NoSQL хранилищ (Redis)
type NoSQLClient interface {
	Client
	
	// Set устанавливает значение по ключу
	Set(ctx context.Context, key string, value interface{}) error
	
	// Get получает значение по ключу
	Get(ctx context.Context, key string) (string, error)
	
	// Del удаляет значение по ключу
	Del(ctx context.Context, key string) error
	
	// Exists проверяет существование ключа
	Exists(ctx context.Context, key string) (bool, error)
	
	// HSet устанавливает значение в хеш
	HSet(ctx context.Context, key, field string, value interface{}) error
	
	// HGet получает значение из хеша
	HGet(ctx context.Context, key, field string) (string, error)
}

// Rows — интерфейс для итерации по строкам результата SQL запроса
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close() error
}

// Row — интерфейс для работы с одной строкой результата SQL запроса
type Row interface {
	Scan(dest ...any) error
}

// Transaction — интерфейс для работы с транзакциями
type Transaction interface {
	Exec(ctx context.Context, query string, args ...any) error
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Batch — интерфейс для батч-вставок
type Batch interface {
	Append(args ...any) error
	Send() error
}

