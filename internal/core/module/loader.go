package module

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type Loader struct {
	registry *Registry
	logger   *zap.Logger
}

func NewLoader(registry *Registry, logger *zap.Logger) *Loader {
	return &Loader{
		registry: registry,
		logger:   logger,
	}
}

func (l *Loader) Load(modules []Module) error {
	for _, mod := range modules {
		if err := l.registry.Register(mod); err != nil {
			return err
		}
		l.logger.Info("module registered", zap.String("module", mod.Name()))
	}
	return nil
}

func (l *Loader) InitAll(ctx context.Context, deps Dependencies) error {
	deps.Registry = l.registry
	for _, mod := range l.registry.All() {
		modDeps := deps
		modDeps.Router = deps.Router
		if err := mod.Init(ctx, modDeps); err != nil {
			return fmt.Errorf("init module %s: %w", mod.Name(), err)
		}
		mod.RegisterRoutes(deps.Router)
		l.logger.Info("module initialized", zap.String("module", mod.Name()))
	}
	return nil
}

func (l *Loader) Registry() *Registry {
	return l.registry
}
