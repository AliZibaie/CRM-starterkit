package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/modules/announcement/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("announcement not found")
	ErrInvalidInput = errors.New("invalid input")
)

type Service struct {
	repo domain.Repository
	bus  *event.Bus
}

func NewService(repo domain.Repository, bus *event.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

func (s *Service) Create(ctx context.Context, in domain.CreateInput) (*domain.Announcement, error) {
	title := strings.TrimSpace(in.Title)
	body := strings.TrimSpace(in.Body)
	audience := strings.TrimSpace(in.Audience)
	if title == "" || body == "" {
		return nil, ErrInvalidInput
	}
	if audience == "" {
		audience = domain.AudienceAll
	}
	if audience != domain.AudienceAll && audience != domain.AudienceRole && audience != domain.AudienceGroup {
		return nil, ErrInvalidInput
	}
	if audience != domain.AudienceAll && strings.TrimSpace(in.AudienceRef) == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now().UTC()
	a := &domain.Announcement{
		ID: uuid.NewString(), Title: title, Body: body, Audience: audience,
		AudienceRef: strings.TrimSpace(in.AudienceRef), Status: "draft",
		CreatedBy: in.CreatedBy, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) Update(ctx context.Context, id string, in domain.CreateInput) (*domain.Announcement, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if a.Status == "broadcasted" {
		return nil, ErrInvalidInput
	}
	title := strings.TrimSpace(in.Title)
	body := strings.TrimSpace(in.Body)
	audience := strings.TrimSpace(in.Audience)
	if title == "" || body == "" {
		return nil, ErrInvalidInput
	}
	if audience == "" {
		audience = domain.AudienceAll
	}
	if audience != domain.AudienceAll && audience != domain.AudienceRole && audience != domain.AudienceGroup {
		return nil, ErrInvalidInput
	}
	if audience != domain.AudienceAll && strings.TrimSpace(in.AudienceRef) == "" {
		return nil, ErrInvalidInput
	}
	a.Title = title
	a.Body = body
	a.Audience = audience
	a.AudienceRef = strings.TrimSpace(in.AudienceRef)
	a.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if a.Status == "broadcasted" {
		return ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) ResolveAudience(a *domain.Announcement) map[string]any {
	return map[string]any{
		"type": a.Audience,
		"ref":  a.AudienceRef,
	}
}

func (s *Service) Broadcast(ctx context.Context, id string) (*domain.Announcement, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	now := time.Now().UTC()
	a.Status = "broadcasted"
	a.BroadcastAt = &now
	a.UpdatedAt = now
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	if s.bus != nil {
		s.bus.Publish(domain.EventBroadcast, domain.BroadcastPayload{
			AnnouncementID: a.ID, Title: a.Title, Body: a.Body,
			Audience: a.Audience, AudienceRef: a.AudienceRef,
		})
	}
	return a, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Announcement, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id string) (*domain.Announcement, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return a, nil
}
