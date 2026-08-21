package cms

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/cms/delivery/http"
	"crm-starterkit/internal/modules/cms/repository"
	"crm-starterkit/internal/modules/cms/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "cms" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	m.handler = httpdelivery.NewHandler(usecase.NewService(store), deps.Validate)
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
	return []string{"internal/modules/cms/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "cms.manage", Name: "Manage CMS", Description: "Manage pages menus seo faqs", Module: "cms"},
	}
}
