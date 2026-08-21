package database

import (
	"crm-starterkit/config"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type SQLServerDriver struct{}

func (d *SQLServerDriver) Name() string {
	return "sqlserver"
}

func (d *SQLServerDriver) Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn, err := cfg.DSN()
	if err != nil {
		return nil, err
	}
	return gorm.Open(sqlserver.Open(dsn), gormConfig())
}
