package httpcore

import (
	"context"
	"fmt"
	"time"

	"crm-starterkit/config"
	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/module"
	"crm-starterkit/internal/core/storage"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	app         *fiber.App
	cfg         *config.Config
	log         *zap.Logger
	db          *gorm.DB
	jwtManager  *jwtutil.Manager
	loader      *module.Loader
	events      *event.Bus
	storage     storage.Driver
	recorderHub *middleware.RecorderHub
}

type Options struct {
	Config            *config.Config
	Logger            *zap.Logger
	DB                *gorm.DB
	JWTManager        *jwtutil.Manager
	Loader            *module.Loader
	Events            *event.Bus
	Storage           storage.Driver
	RecorderHub       *middleware.RecorderHub
	RequestRecorder   middleware.RequestRecorder
	PermissionChecker middleware.PermissionChecker
}

func NewServer(opts Options) *Server {
	app := fiber.New(fiber.Config{
		AppName:      opts.Config.App.Name,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return response.FromError(c, err)
		},
	})

	s := &Server{
		app:         app,
		cfg:         opts.Config,
		log:         opts.Logger,
		db:          opts.DB,
		jwtManager:  opts.JWTManager,
		loader:      opts.Loader,
		events:      opts.Events,
		storage:     opts.Storage,
		recorderHub: opts.RecorderHub,
	}

	s.registerMiddleware(opts)
	s.registerCoreRoutes()

	return s
}

func (s *Server) registerMiddleware(opts Options) {
	s.app.Use(middleware.Recover(s.log))
	s.app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Request-ID", uuid.NewString())
		return c.Next()
	})
	s.app.Use(middleware.CORS(s.cfg.CORS))
	s.app.Use(middleware.RateLimit(s.cfg.RateLimit))
	s.app.Use(middleware.OptionalAuth(s.jwtManager))
	recorder := opts.RequestRecorder
	if recorder == nil && opts.RecorderHub != nil {
		recorder = opts.RecorderHub.Record
	}
	s.app.Use(middleware.RequestLogger(s.log, recorder))
}

func (s *Server) registerCoreRoutes() {
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return response.OK(c, fiber.Map{
			"status":  "ok",
			"service": s.cfg.App.Name,
			"env":     s.cfg.App.Env,
		})
	})
	s.app.Get("/ready", func(c *fiber.Ctx) error {
		sqlDB, err := s.db.DB()
		if err != nil {
			return response.Internal(c, "database unavailable")
		}
		if err := sqlDB.Ping(); err != nil {
			return response.Internal(c, "database unavailable")
		}
		return response.OK(c, fiber.Map{"status": "ready"})
	})
}

func (s *Server) App() *fiber.App {
	return s.app
}

func (s *Server) APIGroup() fiber.Router {
	return s.app.Group("/api/v1")
}

func (s *Server) BootstrapModules(ctx context.Context) error {
	if s.loader == nil {
		return nil
	}
	return s.loader.InitAll(ctx, module.Dependencies{
		DB:          s.db,
		Router:      s.APIGroup(),
		JWT:         s.jwtManager,
		Validate:    validator.New(),
		Events:      s.events,
		Config:      s.cfg,
		Storage:     s.storage,
		RecorderHub: s.recorderHub,
	})
}

func (s *Server) Start() error {
	addr := s.cfg.Address()
	s.log.Info("starting http server", zap.String("addr", addr))
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	_ = ctx
	s.log.Info("shutting down http server")
	return s.app.Shutdown()
}

func (s *Server) ListenAndServe() error {
	if err := s.BootstrapModules(context.Background()); err != nil {
		return fmt.Errorf("bootstrap modules: %w", err)
	}
	return s.Start()
}
