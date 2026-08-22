package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	modulesconfig "crm-starterkit"
	"crm-starterkit/config"
	"crm-starterkit/internal/core/cache"
	"crm-starterkit/internal/core/database"
	"crm-starterkit/internal/core/event"
	httpcore "crm-starterkit/internal/core/http"
	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/logger"
	"crm-starterkit/internal/core/module"
	"crm-starterkit/internal/core/storage"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/swagger"
	"go.uber.org/zap"

	_ "crm-starterkit/docs"
)

// @title CRM StarterKit API
// @version 1.0
// @description Modular monolith CRM backend
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	cacheDriver, err := cache.Resolve(cfg.Cache)
	if err != nil {
		log.Fatal("cache connection failed", zap.Error(err))
	}
	defer cacheDriver.Close()

	storageDriver, err := storage.Resolve(cfg.Storage)
	if err != nil {
		log.Fatal("storage init failed", zap.Error(err))
	}
	bus := event.NewBus()
	_ = cacheDriver
	recorderHub := middleware.NewRecorderHub()

	jwtManager := jwtutil.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.Issuer,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)

	registry := module.NewRegistry()
	loader := module.NewLoader(registry, log)
	if err := loader.Load(modulesconfig.EnabledModules); err != nil {
		log.Fatal("module load failed", zap.Error(err))
	}

	server := httpcore.NewServer(httpcore.Options{
		Config:      cfg,
		Logger:      log,
		DB:          db,
		JWTManager:  jwtManager,
		Loader:      loader,
		Events:      bus,
		Storage:     storageDriver,
		RecorderHub: recorderHub,
	})
	httpcore.MountCoreRoutes(server.App())
	server.App().Get("/swagger/*", swagger.HandlerDefault)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("server shutdown failed", zap.Error(err))
		}
	case err := <-errCh:
		if err != nil {
			log.Fatal("server stopped", zap.Error(err))
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
