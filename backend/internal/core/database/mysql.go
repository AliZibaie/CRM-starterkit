package database

import (
	"crm-starterkit/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDriver struct{}

func (d *MySQLDriver) Name() string {
	return "mysql"
}

func (d *MySQLDriver) Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn, err := cfg.DSN()
	if err != nil {
		return nil, err
	}
	return gorm.Open(mysql.Open(dsn), gormConfig())
}
