package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-starterkit/internal/modules/review/domain"

	"github.com/google/uuid"
)

var ErrInvalidInput = errors.New("invalid input")

type Service struct{ repo domain.Repository }

func NewService(repo domain.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, itemID, userID, comment string, rating int) (*domain.Review, error) {
	if itemID == "" || userID == "" || rating < 1 || rating > 5 {
		return nil, ErrInvalidInput
	}
	r := &domain.Review{
		ID: uuid.NewString(), ItemID: itemID, UserID: userID, Rating: rating,
		Comment: strings.TrimSpace(comment), IsVisible: true, CreatedAt: time.Now().UTC(),
	}
	return r, s.repo.Create(ctx, r)
}

func (s *Service) ListByItem(ctx context.Context, itemID string) ([]domain.Review, float64, error) {
	items, err := s.repo.ListByItem(ctx, itemID)
	if err != nil {
		return nil, 0, err
	}
	avg, err := s.repo.AverageRating(ctx, itemID)
	return items, avg, err
}

func (s *Service) ListAll(ctx context.Context, itemID string, rating *int, limit int) ([]domain.Review, error) {
	return s.repo.ListAll(ctx, domain.ListFilter{ItemID: itemID, Rating: rating, Limit: limit})
}

func (s *Service) SetVisibility(ctx context.Context, id string, visible bool) error {
	return s.repo.SetVisibility(ctx, id, visible)
}

func (s *Service) AdminDelete(ctx context.Context, id string) error {
	return s.repo.AdminDelete(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
