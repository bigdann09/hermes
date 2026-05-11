package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bigdann09/notifications/internal/config"
	"github.com/redis/go-redis/v9"
)

var CacheClient *redis.Client
var (
	ErrCacheMiss error = errors.New("cache miss")
)

type Cache struct {
	client *redis.Client
	ctx    context.Context
}

func NewCache(cfg *config.CacheConfig) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil
	}

	CacheClient = client
	return &Cache{
		client: client,
		ctx:    ctx,
	}
}

func GetCache() *redis.Client {
	return CacheClient
}

func (c *Cache) Set(key string, value any, ttl time.Duration) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(c.ctx, key, payload, ttl).Err()
}

func (c *Cache) Get(key string, dest any) error {
	payload, err := c.client.Get(c.ctx, key).Bytes()
	if err != nil {
		return ErrCacheMiss
	}
	return json.Unmarshal(payload, dest)
}

func (c *Cache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}
