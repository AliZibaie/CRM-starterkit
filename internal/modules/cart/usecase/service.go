package usecase

import (
	"context"
	"errors"
	"time"

	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/modules/cart/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const EventCartCleared = "commerce.cart.cleared"

var (
	ErrNotFound     = errors.New("cart not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrOutOfStock   = errors.New("out of stock")
)

type Service struct {
	repo domain.Repository
	bus  *event.Bus
}

func NewService(repo domain.Repository, bus *event.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

func (s *Service) GetCart(ctx context.Context, userID string) (*domain.Cart, []domain.CartItem, int64, error) {
	cart, err := s.repo.GetOrCreateCart(ctx, userID, "IRR")
	if err != nil {
		return nil, nil, 0, err
	}
	items, err := s.repo.ListItems(ctx, cart.ID)
	if err != nil {
		return nil, nil, 0, err
	}
	var total int64
	for _, it := range items {
		total += it.UnitPrice * it.Quantity
	}
	return cart, items, total, nil
}

func (s *Service) AddItem(ctx context.Context, userID, itemID, variantID string, qty int64) (*domain.CartItem, error) {
	if qty <= 0 || itemID == "" {
		return nil, ErrInvalidInput
	}
	snap, err := s.repo.LookupCatalog(ctx, itemID, variantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if snap.VariantID != "" && snap.Available < qty {
		return nil, ErrOutOfStock
	}
	cart, err := s.repo.GetOrCreateCart(ctx, userID, snap.Currency)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	item := &domain.CartItem{
		ID: uuid.NewString(), CartID: cart.ID, ItemID: snap.ItemID, VariantID: snap.VariantID,
		Title: snap.Title, UnitPrice: snap.UnitPrice, Quantity: qty, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.UpsertItem(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) RemoveItem(ctx context.Context, userID, cartItemID string) error {
	cart, err := s.repo.FindCartByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.RemoveItem(ctx, cart.ID, cartItemID)
}

func (s *Service) Clear(ctx context.Context, userID string) error {
	cart, err := s.repo.FindCartByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if err := s.repo.Clear(ctx, cart.ID); err != nil {
		return err
	}
	if s.bus != nil {
		s.bus.Publish(EventCartCleared, map[string]string{"user_id": userID, "cart_id": cart.ID})
	}
	return nil
}
