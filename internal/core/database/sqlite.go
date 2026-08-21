package database

import (
	"os"
	"path/filepath"

	"crm-starterkit/config"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type SQLiteDriver struct{}

func (d *SQLiteDriver) Name() string {
	return "sqlite"
}

func (d *SQLiteDriver) Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn, err := cfg.DSN()
	if err != nil {
		return nil, err
	}
	if dir := filepath.Dir(dsn); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	return gorm.Open(sqlite.Open(dsn), gormConfig())
}
