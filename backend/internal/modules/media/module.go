package media

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/media/delivery/http"
	"crm-starterkit/internal/modules/media/repository"
	"crm-starterkit/internal/modules/media/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "media" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	service := usecase.NewService(store, deps.Storage)
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
	return []string{"internal/modules/media/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "media.upload", Name: "Upload Media", Description: "Upload media files", Module: "media"},
		{Key: "media.read", Name: "Read Media", Description: "View media files", Module: "media"},
		{Key: "media.delete", Name: "Delete Media", Description: "Delete media files", Module: "media"},
	}
}
