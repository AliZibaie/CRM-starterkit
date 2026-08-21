package cache

import (
	"context"
	"fmt"
	"time"

	"crm-starterkit/config"

	"github.com/redis/go-redis/v9"
)

type RedisDriver struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

func NewRedisDriver(cfg config.CacheConfig) (*RedisDriver, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &RedisDriver{
		client: client,
		prefix: cfg.Prefix,
		ttl:    cfg.DefaultTTL,
	}, nil
}

func (d *RedisDriver) Name() string {
	return "redis"
}

func (d *RedisDriver) key(key string) string {
	return d.prefix + key
}

func (d *RedisDriver) Get(ctx context.Context, key string) ([]byte, error) {
	return d.client.Get(ctx, d.key(key)).Bytes()
}

func (d *RedisDriver) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = d.ttl
	}
	return d.client.Set(ctx, d.key(key), value, ttl).Err()
}

func (d *RedisDriver) Delete(ctx context.Context, key string) error {
	return d.client.Del(ctx, d.key(key)).Err()
}

func (d *RedisDriver) Exists(ctx context.Context, key string) (bool, error) {
	n, err := d.client.Exists(ctx, d.key(key)).Result()
	return n > 0, err
}

func (d *RedisDriver) Ping(ctx context.Context) error {
	return d.client.Ping(ctx).Err()
}

func (d *RedisDriver) Close() error {
	return d.client.Close()
}
