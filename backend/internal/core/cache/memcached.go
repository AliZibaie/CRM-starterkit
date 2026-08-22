package cache

import (
	"context"
	"fmt"
	"time"

	"crm-starterkit/config"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedDriver struct {
	client *memcache.Client
	prefix string
	ttl    time.Duration
}

func NewMemcachedDriver(cfg config.CacheConfig) (*MemcachedDriver, error) {
	client := memcache.New(cfg.Memcached.Host)
	client.Timeout = 5 * time.Second
	if err := client.Ping(); err != nil {
		return nil, fmt.Errorf("memcached ping: %w", err)
	}
	return &MemcachedDriver{
		client: client,
		prefix: cfg.Prefix,
		ttl:    cfg.DefaultTTL,
	}, nil
}

func (d *MemcachedDriver) Name() string {
	return "memcached"
}

func (d *MemcachedDriver) key(key string) string {
	return d.prefix + key
}

func (d *MemcachedDriver) Get(ctx context.Context, key string) ([]byte, error) {
	_ = ctx
	item, err := d.client.Get(d.key(key))
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func (d *MemcachedDriver) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	_ = ctx
	if ttl <= 0 {
		ttl = d.ttl
	}
	seconds := int32(ttl.Seconds())
	if seconds < 1 {
		seconds = 1
	}
	return d.client.Set(&memcache.Item{
		Key:        d.key(key),
		Value:      value,
		Expiration: seconds,
	})
}

func (d *MemcachedDriver) Delete(ctx context.Context, key string) error {
	_ = ctx
	return d.client.Delete(d.key(key))
}

func (d *MemcachedDriver) Exists(ctx context.Context, key string) (bool, error) {
	_ = ctx
	_, err := d.client.Get(d.key(key))
	if err == memcache.ErrCacheMiss {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *MemcachedDriver) Ping(ctx context.Context) error {
	_ = ctx
	return d.client.Ping()
}

func (d *MemcachedDriver) Close() error {
	return nil
}
