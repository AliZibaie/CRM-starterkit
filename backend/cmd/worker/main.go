package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	modulesconfig "crm-starterkit"
	"crm-starterkit/config"
	"crm-starterkit/internal/core/database"
	"crm-starterkit/internal/core/logger"
	"crm-starterkit/internal/core/module"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load(envPath())
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.Log)
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	db, err := database.Open(cfg.Database)
	if err != nil {
		log.Fatal("database connection failed", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("database handle failed", zap.Error(err))
	}
	defer sqlDB.Close()
	_ = db

	registry := module.NewRegistry()
	loader := module.NewLoader(registry, log)
	if err := loader.Load(modulesconfig.EnabledModules); err != nil {
		log.Fatal("module load failed", zap.Error(err))
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Queue.RedisAddr,
			Password: cfg.Queue.RedisPassword,
			DB:       cfg.Queue.RedisDB,
		},
		asynq.Config{
			Concurrency: cfg.Worker.Concurrency,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error("task failed", zap.String("type", task.Type()), zap.Error(err))
			}),
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("core:ping", func(ctx context.Context, task *asynq.Task) error {
		log.Info("ping task received", zap.ByteString("payload", task.Payload()))
		return nil
	})

	log.Info("worker started",
		zap.Int("concurrency", cfg.Worker.Concurrency),
		zap.Strings("modules", registry.Names()),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Run(mux)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Info("worker shutting down", zap.String("signal", sig.String()))
		srv.Shutdown()
	case err := <-errCh:
		if err != nil {
			log.Fatal("worker stopped", zap.Error(err))
		}
	}
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
