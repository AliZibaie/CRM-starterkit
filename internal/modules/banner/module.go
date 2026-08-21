package banner

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/banner/delivery/http"
	"crm-starterkit/internal/modules/banner/repository"
	"crm-starterkit/internal/modules/banner/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "banner" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	service := usecase.NewService(store)
	m.handler = httpdelivery.NewHandler(service, deps.Validate)
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
	return []string{"internal/modules/banner/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "banner.create", Name: "Create Banner", Description: "Create banners", Module: "banner"},
		{Key: "banner.read", Name: "Read Banner", Description: "View banners", Module: "banner"},
		{Key: "banner.schedule", Name: "Schedule Banner", Description: "Schedule banner display window", Module: "banner"},
		{Key: "banner.position", Name: "Set Banner Position", Description: "Set banner slot/position", Module: "banner"},
	}
}
