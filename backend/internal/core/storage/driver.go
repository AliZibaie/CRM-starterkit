package storage

import (
	"context"
	"fmt"
	"io"

	"crm-starterkit/config"
)

type Object struct {
	Key         string
	ContentType string
	Size        int64
	URL         string
}

type Driver interface {
	Put(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (*Object, error)
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	URL(key string) string
	Name() string
}

func Resolve(cfg config.StorageConfig) (Driver, error) {
	switch cfg.Driver {
	case "local":
		return NewLocalDriver(cfg)
	case "s3":
		return NewS3Driver(cfg)
	default:
		return nil, fmt.Errorf("unsupported storage driver: %s", cfg.Driver)
	}
}
