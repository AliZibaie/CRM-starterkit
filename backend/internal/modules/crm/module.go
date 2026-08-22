package crm

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/crm/delivery/http"
	"crm-starterkit/internal/modules/crm/repository"
	"crm-starterkit/internal/modules/crm/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "crm" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	service := usecase.NewService(store, deps.Events)
	service.BindEvents(deps.Events)
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
	return []string{"internal/modules/crm/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "crm.lead", Name: "CRM Leads", Description: "Manage leads", Module: "crm"},
		{Key: "crm.contact", Name: "CRM Contacts", Description: "Manage contacts", Module: "crm"},
		{Key: "crm.deal", Name: "CRM Deals", Description: "Manage deals", Module: "crm"},
		{Key: "crm.task", Name: "CRM Tasks", Description: "Manage tasks", Module: "crm"},
		{Key: "crm.view360", Name: "CRM Customer 360", Description: "View customer 360", Module: "crm"},
		{Key: "crm.manage", Name: "CRM Manage", Description: "Manage CRM resources", Module: "crm"},
	}
}
