package support

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/support/delivery/http"
	"crm-starterkit/internal/modules/support/realtime"
	"crm-starterkit/internal/modules/support/repository"
	"crm-starterkit/internal/modules/support/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "support" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	hub := realtime.NewHub()
	service := usecase.NewService(store, hub)
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
	return []string{"internal/modules/support/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "support.read", Name: "Read Support", Description: "View tickets and departments", Module: "support"},
		{Key: "support.department.manage", Name: "Manage Departments", Description: "Manage support departments", Module: "support"},
		{Key: "support.ticket.create", Name: "Create Ticket", Description: "Create support tickets", Module: "support"},
		{Key: "support.ticket.reply", Name: "Reply Ticket", Description: "Reply to tickets", Module: "support"},
		{Key: "support.ticket.assign", Name: "Assign Ticket", Description: "Assign tickets to agents", Module: "support"},
		{Key: "support.ticket.status", Name: "Change Ticket Status", Description: "Change ticket status", Module: "support"},
		{Key: "support.ticket.close", Name: "Close Ticket", Description: "Close tickets", Module: "support"},
	}
}
