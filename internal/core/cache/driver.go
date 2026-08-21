package cache

import (
	"context"
	"fmt"
	"time"

	"crm-starterkit/config"
)

type Driver interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Ping(ctx context.Context) error
	Close() error
	Name() string
}

func Resolve(cfg config.CacheConfig) (Driver, error) {
	switch cfg.Driver {
	case "redis":
		return NewRedisDriver(cfg)
	case "memcached":
		return NewMemcachedDriver(cfg)
	default:
		return nil, fmt.Errorf("unsupported cache driver: %s", cfg.Driver)
	}
}
