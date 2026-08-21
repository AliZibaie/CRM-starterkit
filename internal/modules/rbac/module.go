package rbac

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/rbac/delivery/http"
	"crm-starterkit/internal/modules/rbac/repository"
	"crm-starterkit/internal/modules/rbac/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	service    *usecase.Service
	jwtManager *jwtutil.Manager
}

func New() *Module {
	return &Module{}
}

func (m *Module) Name() string {
	return "rbac"
}

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	m.service = usecase.NewService(store)
	m.handler = httpdelivery.NewHandler(m.service, deps.Validate, deps.Registry)
	m.jwtManager = deps.JWT
	if deps.Registry != nil {
		if _, err := m.service.SyncPermissionsFromRegistry(ctx, deps.Registry.Permissions()); err != nil {
			return err
		}
	}
	if _, err := m.service.EnsureAdminRole(ctx); err != nil {
		return err
	}
	return nil
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	if m.handler == nil || m.jwtManager == nil {
		return
	}
	httpdelivery.RegisterRoutes(router, m.handler, m.jwtManager)
}

func (m *Module) RegisterMigrations() []string {
	return []string{"internal/modules/rbac/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "rbac.role.create", Name: "Create Role", Description: "Create roles", Module: "rbac"},
		{Key: "rbac.role.read", Name: "Read Role", Description: "View roles", Module: "rbac"},
		{Key: "rbac.permission.create", Name: "Create Permission", Description: "Create permissions", Module: "rbac"},
		{Key: "rbac.permission.read", Name: "Read Permission", Description: "View permissions", Module: "rbac"},
		{Key: "rbac.permission.edit", Name: "Edit Permission", Description: "Edit permissions", Module: "rbac"},
		{Key: "rbac.permission.sync", Name: "Sync Permissions", Description: "Sync permissions from module registry", Module: "rbac"},
		{Key: "rbac.assign.permission", Name: "Assign Permission", Description: "Assign permission to role", Module: "rbac"},
		{Key: "rbac.assign.role", Name: "Assign Role", Description: "Assign role to user", Module: "rbac"},
	}
}

func (m *Module) Service() *usecase.Service {
	return m.service
}
