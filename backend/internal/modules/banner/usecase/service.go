package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-starterkit/internal/modules/banner/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("banner not found")
	ErrInvalidInput = errors.New("invalid input")
)

var allowedPositions = map[string]struct{}{
	"home_hero": {}, "home_sidebar": {}, "catalog_top": {}, "checkout_banner": {}, "footer": {},
}

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, in domain.CreateInput) (*domain.Banner, error) {
	title := strings.TrimSpace(in.Title)
	image := strings.TrimSpace(in.ImageURL)
	position := strings.TrimSpace(in.Position)
	if title == "" || image == "" || position == "" {
		return nil, ErrInvalidInput
	}
	if _, ok := allowedPositions[position]; !ok {
		return nil, ErrInvalidInput
	}
	now := time.Now().UTC()
	b := &domain.Banner{
		ID: uuid.NewString(), Title: title, ImageURL: image, LinkURL: strings.TrimSpace(in.LinkURL),
		Position: position, IsActive: true, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) Update(ctx context.Context, id string, in domain.CreateInput) (*domain.Banner, error) {
	b, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	title := strings.TrimSpace(in.Title)
	image := strings.TrimSpace(in.ImageURL)
	position := strings.TrimSpace(in.Position)
	if title == "" || image == "" || position == "" {
		return nil, ErrInvalidInput
	}
	if _, ok := allowedPositions[position]; !ok {
		return nil, ErrInvalidInput
	}
	b.Title = title
	b.ImageURL = image
	b.LinkURL = strings.TrimSpace(in.LinkURL)
	b.Position = position
	b.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) Schedule(ctx context.Context, id string, in domain.ScheduleInput) (*domain.Banner, error) {
	b, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if in.StartsAt != nil && in.EndsAt != nil && in.EndsAt.Before(*in.StartsAt) {
		return nil, ErrInvalidInput
	}
	b.StartsAt = in.StartsAt
	b.EndsAt = in.EndsAt
	if in.IsActive != nil {
		b.IsActive = *in.IsActive
	}
	b.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) SetPosition(ctx context.Context, id, position string) (*domain.Banner, error) {
	position = strings.TrimSpace(position)
	if _, ok := allowedPositions[position]; !ok {
		return nil, ErrInvalidInput
	}
	b, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	b.Position = position
	b.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Banner, error) {
	return s.repo.List(ctx)
}

func (s *Service) ListByPosition(ctx context.Context, position string) ([]domain.Banner, error) {
	return s.repo.ListByPosition(ctx, position)
}

func (s *Service) Positions() []string {
	out := make([]string, 0, len(allowedPositions))
	for p := range allowedPositions {
		out = append(out, p)
	}
	return out
}
