package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/modules/setting/domain"

	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("setting not found")
	ErrInvalidInput = errors.New("invalid input")
)

type Service struct {
	repo domain.Repository
	bus  *event.Bus
}

func NewService(repo domain.Repository, bus *event.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

func (s *Service) Get(ctx context.Context, key string) (*domain.Setting, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrInvalidInput
	}
	item, err := s.repo.Get(ctx, key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}

func (s *Service) Set(ctx context.Context, key, value, updatedBy string) (*domain.Setting, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrInvalidInput
	}
	item := &domain.Setting{Key: key, Value: value, UpdatedBy: updatedBy, UpdatedAt: time.Now().UTC()}
	if err := s.repo.Set(ctx, item); err != nil {
		return nil, err
	}
	if s.bus != nil {
		s.bus.Publish("audit.change", map[string]any{
			"actor_id": updatedBy, "module": "setting", "action": "set",
			"entity_type": "setting", "entity_id": key, "after": value,
		})
	}
	return item, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Setting, error) {
	return s.repo.List(ctx)
}
