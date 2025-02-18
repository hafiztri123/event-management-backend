package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hafiztri123/src/internal/pkg/config"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
	defaultExpiration time.Duration
}

func NewRedisCache(cfg *config.RedisConfig) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB: 0,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("[FAIL] failed to connect to redis: %v", err)
	}

	return &RedisCache{
		Client: client,
		defaultExpiration: time.Duration(cfg.DurationMinute) * time.Minute,
	}
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = c.defaultExpiration
	}

	json, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("[FAIL] failed to marshal value: %v", err)
	}

	return c.Client.Set(ctx, key, json, expiration).Err()
}

func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	data, err := c.Client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return fmt.Errorf("[FAIL] failed to get value: %v", err)
	}

	return json.Unmarshal(data, value)
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

func (c *RedisCache) Clear(ctx context.Context) error {
	return c.Client.FlushAll(ctx).Err()
}

