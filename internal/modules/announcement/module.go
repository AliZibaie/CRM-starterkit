package announcement

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/announcement/delivery/http"
	"crm-starterkit/internal/modules/announcement/repository"
	"crm-starterkit/internal/modules/announcement/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "announcement" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	service := usecase.NewService(store, deps.Events)
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
	return []string{"internal/modules/announcement/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "announcement.create", Name: "Create Announcement", Description: "Create announcements", Module: "announcement"},
		{Key: "announcement.read", Name: "Read Announcement", Description: "View announcements", Module: "announcement"},
		{Key: "announcement.broadcast", Name: "Broadcast Announcement", Description: "Broadcast announcements", Module: "announcement"},
	}
}
