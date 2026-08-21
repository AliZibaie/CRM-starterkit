package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"crm-starterkit/config"
)

type LocalDriver struct {
	root      string
	publicURL string
}

func NewLocalDriver(cfg config.StorageConfig) (*LocalDriver, error) {
	if err := os.MkdirAll(cfg.LocalPath, 0o755); err != nil {
		return nil, err
	}
	return &LocalDriver{
		root:      cfg.LocalPath,
		publicURL: strings.TrimRight(cfg.PublicBaseURL, "/"),
	}, nil
}

func (d *LocalDriver) Name() string {
	return "local"
}

func (d *LocalDriver) Put(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (*Object, error) {
	_ = ctx
	_ = size
	_ = contentType
	path := filepath.Join(d.root, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	written, err := io.Copy(file, reader)
	if err != nil {
		return nil, err
	}
	return &Object{
		Key:         key,
		ContentType: contentType,
		Size:        written,
		URL:         d.URL(key),
	}, nil
}

func (d *LocalDriver) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	_ = ctx
	return os.Open(filepath.Join(d.root, filepath.FromSlash(key)))
}

func (d *LocalDriver) Delete(ctx context.Context, key string) error {
	_ = ctx
	return os.Remove(filepath.Join(d.root, filepath.FromSlash(key)))
}

func (d *LocalDriver) URL(key string) string {
	return fmt.Sprintf("%s/%s", d.publicURL, strings.TrimLeft(key, "/"))
}
