package payment

import (
	"context"

	"crm-starterkit/config"
	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/payment/delivery/http"
	"crm-starterkit/internal/modules/payment/repository"
	"crm-starterkit/internal/modules/payment/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "payment" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	payCfg := config.PaymentConfig{}
	if deps.Config != nil {
		payCfg = deps.Config.Payment
	}
	service := usecase.NewService(store, payCfg, deps.Events)
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
	return []string{"internal/modules/payment/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "payment.init", Name: "Init Payment", Description: "Initialize payments", Module: "payment"},
		{Key: "payment.verify", Name: "Verify Payment", Description: "Verify payments", Module: "payment"},
		{Key: "payment.refund", Name: "Refund Payment", Description: "Refund payments", Module: "payment"},
	}
}
