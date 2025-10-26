package redisClient

import (
	"context"
	"fmt"
	clientsDB "lib/clients/db"
	"lib/models"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps the Redis connection to implement CashClient interface
type Client struct {
	client *redis.Client
}

// NewClient creates a new Redis client with the provided configuration
func NewClient(ctx context.Context, cfg models.Redis) (clientsDB.CashClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	// Проверяем соединение
	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return &Client{client: rdb}, nil
}

// Get returns the value of key
func (c *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.client.Get(ctx, key)
}

// Set sets the value of key
func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	return c.client.Set(ctx, key, value, expiration)
}

// Del removes the specified keys
func (c *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.client.Del(ctx, keys...)
}

// Exists checks if key exists
func (c *Client) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return c.client.Exists(ctx, keys...)
}

// Expire sets an expiration time on key
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return c.client.Expire(ctx, key, expiration)
}

// TTL returns the remaining time to live of a key
func (c *Client) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return c.client.TTL(ctx, key)
}

// SetNX sets the value of key, only if the key does not exist
func (c *Client) SetNX(ctx context.Context, key string, value any, expiration time.Duration) *redis.BoolCmd {
	return c.client.SetNX(ctx, key, value, expiration)
}

// Incr increments the number stored at key by one
func (c *Client) Incr(ctx context.Context, key string) *redis.IntCmd {
	return c.client.Incr(ctx, key)
}

// Decr decrements the number stored at key by one
func (c *Client) Decr(ctx context.Context, key string) *redis.IntCmd {
	return c.client.Decr(ctx, key)
}

// Ping checks the connection to the server
func (c *Client) Ping(ctx context.Context) *redis.StatusCmd {
	return c.client.Ping(ctx)
}

// Close closes the connection
func (c *Client) Close() error {
	return c.client.Close()
}