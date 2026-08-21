package usecase

import (
	"context"
	"errors"
	"time"

	"crm-starterkit/internal/modules/wishlist/domain"

	"github.com/google/uuid"
)

var ErrInvalidInput = errors.New("invalid input")

type Service struct{ repo domain.Repository }

func NewService(repo domain.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Add(ctx context.Context, userID, itemID string) (*domain.WishlistItem, error) {
	if userID == "" || itemID == "" {
		return nil, ErrInvalidInput
	}
	item := &domain.WishlistItem{ID: uuid.NewString(), UserID: userID, ItemID: itemID, CreatedAt: time.Now().UTC()}
	return item, s.repo.Add(ctx, item)
}

func (s *Service) Remove(ctx context.Context, userID, itemID string) error {
	return s.repo.Remove(ctx, userID, itemID)
}

func (s *Service) List(ctx context.Context, userID string) ([]domain.WishlistItem, error) {
	return s.repo.List(ctx, userID)
}
