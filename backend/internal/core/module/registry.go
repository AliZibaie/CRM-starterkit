package module

import (
	"fmt"
	"sync"
)

type Registry struct {
	mu      sync.RWMutex
	modules map[string]Module
	order   []string
}

func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

func (r *Registry) Register(mod Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	name := mod.Name()
	if name == "" {
		return fmt.Errorf("module name cannot be empty")
	}
	if _, exists := r.modules[name]; exists {
		return fmt.Errorf("module already registered: %s", name)
	}
	r.modules[name] = mod
	r.order = append(r.order, name)
	return nil
}

func (r *Registry) Get(name string) (Module, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	mod, ok := r.modules[name]
	return mod, ok
}

func (r *Registry) All() []Module {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Module, 0, len(r.order))
	for _, name := range r.order {
		out = append(out, r.modules[name])
	}
	return out
}

func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, len(r.order))
	copy(out, r.order)
	return out
}

func (r *Registry) Permissions() []Permission {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var perms []Permission
	for _, name := range r.order {
		perms = append(perms, r.modules[name].RegisterPermissions()...)
	}
	return perms
}

func (r *Registry) Migrations() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var paths []string
	for _, name := range r.order {
		paths = append(paths, r.modules[name].RegisterMigrations()...)
	}
	return paths
}
