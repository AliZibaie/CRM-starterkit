package requestlog

import (
	"context"
	"time"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/requestlog/delivery/http"
	"crm-starterkit/internal/modules/requestlog/repository"
	"crm-starterkit/internal/modules/requestlog/usecase"
	"crm-starterkit/internal/platform/ratelimiter"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "requestlog" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	threshold := 100
	window := time.Minute
	if deps.Config != nil && deps.Config.RateLimit.Max > 0 {
		threshold = deps.Config.RateLimit.Max
		window = deps.Config.RateLimit.Window
	}
	tracker := ratelimiter.New(window, threshold)
	service := usecase.NewService(store, tracker)
	if deps.RecorderHub != nil {
		deps.RecorderHub.Set(service.RecordFromMiddleware)
	}
	m.handler = httpdelivery.NewHandler(service)
	m.jwtManager = deps.JWT
	return nil
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	if m.handler == nil || m.jwtManager == nil {
		return
	}
	httpdelivery.RegisterRoutes(router, m.handler, m.jwtManager)
}

func (m *Module) RegisterMigrations() []string {
	return []string{"internal/modules/requestlog/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "requestlog.read", Name: "Read Request Logs", Description: "View request logs", Module: "requestlog"},
	}
}
