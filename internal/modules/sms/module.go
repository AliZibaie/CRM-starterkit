package sms

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/sms/delivery/http"
	"crm-starterkit/internal/modules/sms/providers"
	"crm-starterkit/internal/modules/sms/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "sms" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	apiKey, sender := "", ""
	if deps.Config != nil {
		apiKey = deps.Config.SMS.KavenegarAPIKey
		sender = deps.Config.SMS.KavenegarSender
	}
	provider := providers.NewKavenegar(apiKey, sender)
	service := usecase.NewService(provider)
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

func (m *Module) RegisterMigrations() []string { return nil }

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "sms.send", Name: "Send SMS", Description: "Send SMS messages", Module: "sms"},
	}
}
