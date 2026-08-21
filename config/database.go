package config

import "fmt"

func (c DatabaseConfig) DSN() (string, error) {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			c.Host, c.User, c.Password, c.Name, c.Port, c.SSLMode, c.Timezone,
		), nil
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User, c.Password, c.Host, c.Port, c.Name,
		), nil
	case "sqlserver":
		return fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?database=%s",
			c.User, c.Password, c.Host, c.Port, c.Name,
		), nil
	case "sqlite":
		return c.SQLitePath, nil
	default:
		return "", fmt.Errorf("unsupported database driver: %s", c.Driver)
	}
}

func (c DatabaseConfig) MigrateURL() (string, error) {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
		), nil
	case "mysql":
		return fmt.Sprintf(
			"mysql://%s:%s@tcp(%s:%d)/%s",
			c.User, c.Password, c.Host, c.Port, c.Name,
		), nil
	case "sqlserver":
		return fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?database=%s",
			c.User, c.Password, c.Host, c.Port, c.Name,
		), nil
	case "sqlite":
		return fmt.Sprintf("sqlite://%s", c.SQLitePath), nil
	default:
		return "", fmt.Errorf("unsupported database driver: %s", c.Driver)
	}
}
