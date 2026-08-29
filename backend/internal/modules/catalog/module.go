package catalog

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/catalog/delivery/http"
	"crm-starterkit/internal/modules/catalog/repository"
	"crm-starterkit/internal/modules/catalog/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "catalog" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	m.handler = httpdelivery.NewHandler(usecase.NewService(store, deps.Events), deps.Validate)
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
	return []string{"internal/modules/catalog/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "catalog.read", Name: "Read Catalog", Description: "View catalog", Module: "catalog"},
		{Key: "catalog.category.manage", Name: "Manage Categories", Description: "Manage categories", Module: "catalog"},
		{Key: "catalog.item.manage", Name: "Manage Items", Description: "Manage catalog items", Module: "catalog"},
		{Key: "catalog.inventory.manage", Name: "Manage Inventory", Description: "Manage inventory", Module: "catalog"},
	}
}
