package database

import (
	"crm-starterkit/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDriver struct{}

func (d *PostgresDriver) Name() string {
	return "postgres"
}

func (d *PostgresDriver) Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn, err := cfg.DSN()
	if err != nil {
		return nil, err
	}
	return gorm.Open(postgres.Open(dsn), gormConfig())
}
