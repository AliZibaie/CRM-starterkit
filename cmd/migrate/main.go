package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	modulesconfig "crm-starterkit"
	"crm-starterkit/config"
	"crm-starterkit/internal/core/logger"
	"crm-starterkit/internal/core/module"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/database/sqlserver"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func main() {
	command := flag.String("command", "up", "migration command: up | down | force | version")
	steps := flag.Int("steps", 0, "number of steps for up/down")
	version := flag.Int("version", 0, "target version for force")
	flag.Parse()

	cfg, err := config.Load(envPath())
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.Log)
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	registry := module.NewRegistry()
	loader := module.NewLoader(registry, log)
	if err := loader.Load(modulesconfig.EnabledModules); err != nil {
		log.Fatal("module load failed", zap.Error(err))
	}

	dsn, err := cfg.Database.MigrateURL()
	if err != nil {
		log.Fatal("migrate dsn failed", zap.Error(err))
	}

	migratePath, err := filepath.Abs(cfg.Migrate.Path)
	if err != nil {
		log.Fatal("migrate path failed", zap.Error(err))
	}
	if err := os.MkdirAll(migratePath, 0o755); err != nil {
		log.Fatal("migrate path create failed", zap.Error(err))
	}

	sourceURL := fmt.Sprintf("file://%s", filepath.ToSlash(migratePath))
	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		log.Fatal("migrate init failed", zap.Error(err), zap.String("source", sourceURL))
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Error("migrate source close failed", zap.Error(srcErr))
		}
		if dbErr != nil {
			log.Error("migrate db close failed", zap.Error(dbErr))
		}
	}()

	modulePaths := registry.Migrations()
	log.Info("running migrations",
		zap.String("command", *command),
		zap.String("path", migratePath),
		zap.Strings("module_paths", modulePaths),
	)

	switch *command {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
	case "force":
		err = m.Force(*version)
	case "version":
		ver, dirty, vErr := m.Version()
		if vErr != nil && !errors.Is(vErr, migrate.ErrNilVersion) {
			log.Fatal("version failed", zap.Error(vErr))
		}
		log.Info("migration version", zap.Uint("version", ver), zap.Bool("dirty", dirty))
		return
	default:
		log.Fatal("unknown command", zap.String("command", *command))
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("migration failed", zap.Error(err))
	}

	log.Info("migration completed")
}

func envPath() string {
	if path := os.Getenv("ENV_FILE"); path != "" {
		return path
	}
	if _, err := os.Stat(".env"); err == nil {
		return ".env"
	}
	return ".env.example"
}
