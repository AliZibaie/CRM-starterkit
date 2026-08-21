package module

import (
	"context"

	"crm-starterkit/config"
	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/storage"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Permission struct {
	Key         string
	Name        string
	Description string
	Module      string
}

type Dependencies struct {
	DB          *gorm.DB
	Router      fiber.Router
	JWT         *jwtutil.Manager
	Registry    *Registry
	Validate    *validator.Validator
	Events      *event.Bus
	Config      *config.Config
	Storage     storage.Driver
	RecorderHub *middleware.RecorderHub
}

type Module interface {
	Name() string
	Init(ctx context.Context, deps Dependencies) error
	RegisterRoutes(router fiber.Router)
	RegisterMigrations() []string
	RegisterPermissions() []Permission
}
