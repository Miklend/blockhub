package db

import (
	"context"
	"fmt"

	"lib/clients/db/clickhouse"
	"lib/clients/db/postgresql"
	"lib/clients/db/redis"
	"lib/models"
	"lib/utils/logging"
)

// ClientFactory создает клиенты для различных СУБД
type ClientFactory struct {
	logger *logging.Logger
}

// NewClientFactory создает новую фабрику клиентов
func NewClientFactory(logger *logging.Logger) *ClientFactory {
	return &ClientFactory{
		logger: logger,
	}
}

// CreateSQLClient создает SQL клиент в зависимости от типа базы данных
func (f *ClientFactory) CreateSQLClient(ctx context.Context, dbType DatabaseType, cfg models.Config) (SQLClient, error) {
	switch dbType {
	case PostgreSQL:
		return postgresql.NewClient(ctx, cfg.PostgreSQL, f.logger)
	case ClickHouse:
		return clickhouse.NewClient(ctx, cfg.Clickhouse, f.logger)
	default:
		return nil, fmt.Errorf("unsupported SQL database type: %s", dbType)
	}
}

// CreateNoSQLClient создает NoSQL клиент в зависимости от типа базы данных
func (f *ClientFactory) CreateNoSQLClient(ctx context.Context, dbType DatabaseType, cfg models.Config) (NoSQLClient, error) {
	switch dbType {
	case Redis:
		return redis.NewClient(ctx, cfg.Redis, f.logger)
	default:
		return nil, fmt.Errorf("unsupported NoSQL database type: %s", dbType)
	}
}

// CreateClient создает клиент любого типа
func (f *ClientFactory) CreateClient(ctx context.Context, dbType DatabaseType, cfg models.Config) (Client, error) {
	switch dbType {
	case PostgreSQL, ClickHouse:
		return f.CreateSQLClient(ctx, dbType, cfg)
	case Redis:
		return f.CreateNoSQLClient(ctx, dbType, cfg)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// Storage представляет централизованное хранилище всех клиентов
type Storage struct {
	clients map[DatabaseType]Client
	factory *ClientFactory
}

// NewStorage создает новое централизованное хранилище клиентов
func NewStorage(ctx context.Context, cfg models.Config, logger *logging.Logger) (*Storage, error) {
	factory := NewClientFactory(logger)
	clients := make(map[DatabaseType]Client)

	// Создаем клиенты для всех поддерживаемых типов баз данных
	dbTypes := []DatabaseType{PostgreSQL, ClickHouse, Redis}
	
	for _, dbType := range dbTypes {
		client, err := factory.CreateClient(ctx, dbType, cfg)
		if err != nil {
			// Закрываем уже созданные клиенты при ошибке
			for _, c := range clients {
				c.Close()
			}
			return nil, fmt.Errorf("failed to create %s client: %w", dbType, err)
		}
		clients[dbType] = client
	}

	return &Storage{
		clients: clients,
		factory: factory,
	}, nil
}

// GetSQLClient возвращает SQL клиент по типу базы данных
func (s *Storage) GetSQLClient(dbType DatabaseType) (SQLClient, error) {
	client, exists := s.clients[dbType]
	if !exists {
		return nil, fmt.Errorf("client for database type %s not found", dbType)
	}

	sqlClient, ok := client.(SQLClient)
	if !ok {
		return nil, fmt.Errorf("client for database type %s is not SQL client", dbType)
	}

	return sqlClient, nil
}

// GetNoSQLClient возвращает NoSQL клиент по типу базы данных
func (s *Storage) GetNoSQLClient(dbType DatabaseType) (NoSQLClient, error) {
	client, exists := s.clients[dbType]
	if !exists {
		return nil, fmt.Errorf("client for database type %s not found", dbType)
	}

	nosqlClient, ok := client.(NoSQLClient)
	if !ok {
		return nil, fmt.Errorf("client for database type %s is not NoSQL client", dbType)
	}

	return nosqlClient, nil
}

// GetClient возвращает клиент любого типа по типу базы данных
func (s *Storage) GetClient(dbType DatabaseType) (Client, error) {
	client, exists := s.clients[dbType]
	if !exists {
		return nil, fmt.Errorf("client for database type %s not found", dbType)
	}

	return client, nil
}

// Close закрывает все клиенты
func (s *Storage) Close() error {
	var lastErr error
	for _, client := range s.clients {
		if err := client.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Ping проверяет соединение со всеми клиентами
func (s *Storage) Ping(ctx context.Context) error {
	for dbType, client := range s.clients {
		if err := client.Ping(ctx); err != nil {
			return fmt.Errorf("ping failed for %s: %w", dbType, err)
		}
	}
	return nil
}

