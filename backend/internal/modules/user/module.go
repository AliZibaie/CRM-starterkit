package user

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/user/delivery/http"
	"crm-starterkit/internal/modules/user/repository"
	"crm-starterkit/internal/modules/user/usecase"
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
	return "user"
}

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	repo := repository.NewGormRepository(deps.DB)
	if err := repo.AutoMigrate(); err != nil {
		return err
	}
	m.service = usecase.NewService(repo)
	m.handler = httpdelivery.NewHandler(m.service, deps.Validate)
	m.jwtManager = deps.JWT

	var role struct {
		ID string `gorm:"column:id"`
	}
	err := deps.DB.WithContext(ctx).Table("roles").Select("id").Where("name = ?", "admin").Take(&role).Error
	if err == nil && role.ID != "" {
		if _, err := m.service.EnsureDefaultAdmin(ctx, "admin@example.com", "ChangeMe123!", role.ID); err != nil {
			return err
		}
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
	return []string{"internal/modules/user/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "user.create", Name: "Create User", Description: "Create users", Module: "user"},
		{Key: "user.read", Name: "Read User", Description: "View users", Module: "user"},
		{Key: "user.update", Name: "Update User", Description: "Update users", Module: "user"},
		{Key: "user.delete", Name: "Delete User", Description: "Delete users", Module: "user"},
		{Key: "user.assign_role", Name: "Assign User Role", Description: "Assign roles to users", Module: "user"},
	}
}
