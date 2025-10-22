package redis

import (
	"context"
	"fmt"

	"lib/clients/db"
	"lib/models"
	"lib/utils/logging"

	"github.com/redis/go-redis/v9"
)

// client реализует интерфейс NoSQLClient для Redis
type client struct {
	rdb    *redis.Client
	logger *logging.Logger
}

// NewClient создает новый Redis клиент
func NewClient(ctx context.Context, cfg models.Redis, logger *logging.Logger) (db.NoSQLClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Проверяем соединение
	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return &client{
		rdb:    rdb,
		logger: logger,
	}, nil
}

// Ping проверяет соединение с базой данных
func (c *client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Close закрывает соединение
func (c *client) Close() error {
	return c.rdb.Close()
}

// GetType возвращает тип базы данных
func (c *client) GetType() db.DatabaseType {
	return db.Redis
}

// Set устанавливает значение по ключу
func (c *client) Set(ctx context.Context, key string, value interface{}) error {
	return c.rdb.Set(ctx, key, value, 0).Err()
}

// Get получает значение по ключу
func (c *client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

// Del удаляет значение по ключу
func (c *client) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

// Exists проверяет существование ключа
func (c *client) Exists(ctx context.Context, key string) (bool, error) {
	result := c.rdb.Exists(ctx, key)
	return result.Val() > 0, result.Err()
}

// HSet устанавливает значение в хеш
func (c *client) HSet(ctx context.Context, key, field string, value interface{}) error {
	return c.rdb.HSet(ctx, key, field, value).Err()
}

// HGet получает значение из хеша
func (c *client) HGet(ctx context.Context, key, field string) (string, error) {
	return c.rdb.HGet(ctx, key, field).Result()
}