package database

import (
	"fmt"

	"crm-starterkit/config"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Driver interface {
	Connect(cfg config.DatabaseConfig) (*gorm.DB, error)
	Name() string
}

func Resolve(driverName string) (Driver, error) {
	switch driverName {
	case "postgres":
		return &PostgresDriver{}, nil
	case "mysql":
		return &MySQLDriver{}, nil
	case "sqlserver":
		return &SQLServerDriver{}, nil
	case "sqlite":
		return &SQLiteDriver{}, nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driverName)
	}
}

func Open(cfg config.DatabaseConfig) (*gorm.DB, error) {
	driver, err := Resolve(cfg.Driver)
	if err != nil {
		return nil, err
	}
	return driver.Connect(cfg)
}

func gormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}
}
