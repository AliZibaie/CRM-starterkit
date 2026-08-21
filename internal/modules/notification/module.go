package notification

import (
	"context"

	"crm-starterkit/config"
	"crm-starterkit/internal/core/module"
	"crm-starterkit/internal/modules/notification/channels/email"
	"crm-starterkit/internal/modules/notification/channels/in_app"
	"crm-starterkit/internal/modules/notification/channels/push"
	smschan "crm-starterkit/internal/modules/notification/channels/sms"
	httpdelivery "crm-starterkit/internal/modules/notification/delivery/http"
	"crm-starterkit/internal/modules/notification/domain"
	"crm-starterkit/internal/modules/notification/repository"
	"crm-starterkit/internal/modules/notification/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "notification" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	emailCfg := config.EmailConfig{}
	if deps.Config != nil {
		emailCfg = deps.Config.Email
	}
	channels := []domain.Channel{
		inapp.New(store),
		email.New(emailCfg),
		smschan.New(deps.Events),
		push.New(),
	}
	service := usecase.NewService(store, channels)
	service.BindEvents(deps.Events)
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
	return []string{"internal/modules/notification/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "notification.dispatch", Name: "Dispatch Notification", Description: "Send notifications", Module: "notification"},
		{Key: "notification.read", Name: "Read Notification", Description: "View notifications", Module: "notification"},
	}
}
