package wishlist

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/wishlist/delivery/http"
	"crm-starterkit/internal/modules/wishlist/repository"
	"crm-starterkit/internal/modules/wishlist/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "wishlist" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	m.handler = httpdelivery.NewHandler(usecase.NewService(store))
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
	return []string{"internal/modules/wishlist/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "wishlist.manage", Name: "Manage Wishlist", Description: "Manage wishlist", Module: "wishlist"},
	}
}
