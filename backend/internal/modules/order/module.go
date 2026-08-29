package order

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/order/delivery/http"
	"crm-starterkit/internal/modules/order/repository"
	"crm-starterkit/internal/modules/order/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "order" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	_ = ctx
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	service := usecase.NewService(store, deps.Events)
	service.BindPaymentEvents(deps.Events)
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
	return []string{"internal/modules/order/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "order.create", Name: "Create Order", Description: "Create orders from cart", Module: "order"},
		{Key: "order.read", Name: "Read Order", Description: "View orders", Module: "order"},
		{Key: "order.checkout", Name: "Checkout Order", Description: "Checkout orders", Module: "order"},
		{Key: "order.discount", Name: "Apply Order Discount", Description: "Apply discount to order", Module: "order"},
		{Key: "order.refund", Name: "Refund Order", Description: "Refund orders", Module: "order"},
		{Key: "order.update", Name: "Update Order", Description: "Update order status", Module: "order"},
	}
}
