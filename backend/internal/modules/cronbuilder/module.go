package cronbuilder

import (
	"context"

	"crm-starterkit/internal/core/module"
	httpdelivery "crm-starterkit/internal/modules/cronbuilder/delivery/http"
	"crm-starterkit/internal/modules/cronbuilder/engine"
	"crm-starterkit/internal/modules/cronbuilder/repository"
	"crm-starterkit/internal/modules/cronbuilder/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	handler    *httpdelivery.Handler
	jwtManager *jwtutil.Manager
	scheduler  *engine.Scheduler
}

func New() *Module {
	return &Module{}
}

func (m *Module) Name() string { return "cronbuilder" }

func (m *Module) Init(ctx context.Context, deps module.Dependencies) error {
	store := repository.NewStore(deps.DB)
	if err := store.AutoMigrate(); err != nil {
		return err
	}
	intro := engine.NewWhitelistIntrospector()
	actions := engine.NewDefaultActionRegistry()
	logs := store.AsLogRepo()
	executor := engine.NewExecutor(deps.DB, intro, actions, logs, deps.Events)
	m.scheduler = engine.NewScheduler(executor)
	service := usecase.NewService(store, logs, intro, actions, executor, m.scheduler)
	active, err := store.ListActive(ctx)
	if err != nil {
		return err
	}
	if err := m.scheduler.LoadAll(active); err != nil {
		return err
	}
	m.scheduler.Start()
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
	return []string{"internal/modules/cronbuilder/migrations"}
}

func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "cronbuilder.create", Name: "Create Cron Job", Description: "Create cron jobs", Module: "cronbuilder"},
		{Key: "cronbuilder.read", Name: "Read Cron Jobs", Description: "View cron jobs and schema", Module: "cronbuilder"},
		{Key: "cronbuilder.run", Name: "Run Cron Job", Description: "Run cron jobs immediately", Module: "cronbuilder"},
	}
}
