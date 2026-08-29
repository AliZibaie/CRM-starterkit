package usecase

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/modules/catalog/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInUse        = errors.New("item is referenced by active orders")
	slugRe          = regexp.MustCompile(`[^a-z0-9]+`)
)

type Service struct{ repo domain.Repository; bus *event.Bus }

func NewService(repo domain.Repository, bus *event.Bus) *Service { return &Service{repo: repo, bus: bus} }

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = slugRe.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func (s *Service) CreateCategory(ctx context.Context, name, parentID string) (*domain.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now().UTC()
	c := &domain.Category{ID: uuid.NewString(), Name: name, Slug: slugify(name), ParentID: parentID, CreatedAt: now, UpdatedAt: now}
	if err := s.repo.CreateCategory(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *Service) CreateItem(ctx context.Context, typ, name, description, categoryID, currency string, basePrice int64) (*domain.Item, error) {
	name = strings.TrimSpace(name)
	switch typ {
	case domain.TypePhysical, domain.TypeService, domain.TypeCourse, domain.TypeDigital:
	default:
		return nil, ErrInvalidInput
	}
	if name == "" || basePrice < 0 {
		return nil, ErrInvalidInput
	}
	if currency == "" {
		currency = "IRR"
	}
	now := time.Now().UTC()
	item := &domain.Item{
		ID: uuid.NewString(), CategoryID: categoryID, Type: typ, Name: name, Slug: slugify(name) + "-" + uuid.NewString()[:8],
		Description: description, BasePrice: basePrice, Currency: currency, IsActive: true, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.CreateItem(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) ListItems(ctx context.Context) ([]domain.Item, error) { return s.repo.ListItems(ctx) }

func (s *Service) ListPublicItems(ctx context.Context) ([]domain.Item, error) {
	return s.repo.ListActiveItems(ctx)
}

func (s *Service) GetItemBySlug(ctx context.Context, slug string) (*domain.Item, []domain.Attribute, []domain.Variant, error) {
	item, err := s.repo.FindItemBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, ErrNotFound
		}
		return nil, nil, nil, err
	}
	return s.GetItem(ctx, item.ID)
}

func (s *Service) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, []domain.Item, error) {
	cat, err := s.repo.FindCategoryBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrNotFound
		}
		return nil, nil, err
	}
	items, err := s.repo.ListItemsByCategory(ctx, cat.ID)
	return cat, items, err
}

func (s *Service) GetItem(ctx context.Context, id string) (*domain.Item, []domain.Attribute, []domain.Variant, error) {
	item, err := s.repo.FindItem(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, ErrNotFound
		}
		return nil, nil, nil, err
	}
	attrs, err := s.repo.ListAttributes(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	variants, err := s.repo.ListVariants(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	return item, attrs, variants, nil
}

func (s *Service) SetAttribute(ctx context.Context, itemID, key, value string) (*domain.Attribute, error) {
	if _, err := s.repo.FindItem(ctx, itemID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, ErrInvalidInput
	}
	a := &domain.Attribute{ID: uuid.NewString(), ItemID: itemID, Key: key, Value: value, CreatedAt: time.Now().UTC()}
	if err := s.repo.UpsertAttribute(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) CreateVariant(ctx context.Context, itemID, sku, name, attrsJSON string, price int64) (*domain.Variant, error) {
	if _, err := s.repo.FindItem(ctx, itemID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if strings.TrimSpace(sku) == "" || strings.TrimSpace(name) == "" || price < 0 {
		return nil, ErrInvalidInput
	}
	now := time.Now().UTC()
	v := &domain.Variant{ID: uuid.NewString(), ItemID: itemID, SKU: sku, Name: name, Price: price, AttrsJSON: attrsJSON, IsActive: true, CreatedAt: now, UpdatedAt: now}
	if err := s.repo.CreateVariant(ctx, v); err != nil {
		return nil, err
	}
	_ = s.repo.SetInventory(ctx, &domain.Inventory{ID: uuid.NewString(), VariantID: v.ID, Quantity: 0, UpdatedAt: now})
	return v, nil
}

func (s *Service) SetInventory(ctx context.Context, variantID string, qty int64) (*domain.Inventory, error) {
	if _, err := s.repo.FindVariant(ctx, variantID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if qty < 0 {
		return nil, ErrInvalidInput
	}
	inv := &domain.Inventory{ID: uuid.NewString(), VariantID: variantID, Quantity: qty, UpdatedAt: time.Now().UTC()}
	if err := s.repo.SetInventory(ctx, inv); err != nil {
		return nil, err
	}
	return inv, nil
}

func (s *Service) GetInventory(ctx context.Context, variantID string) (*domain.Inventory, error) {
	inv, err := s.repo.GetInventory(ctx, variantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return inv, nil
}

func (s *Service) UpdateItem(ctx context.Context, id string, typ, name, description, categoryID, currency string, basePrice int64, isActive bool, actorID string) (*domain.Item, error) {
	item, err := s.repo.FindItem(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	switch typ {
	case domain.TypePhysical, domain.TypeService, domain.TypeCourse, domain.TypeDigital:
	default:
		return nil, ErrInvalidInput
	}
	name = strings.TrimSpace(name)
	if name == "" || basePrice < 0 {
		return nil, ErrInvalidInput
	}
	if currency == "" {
		currency = item.Currency
	}
	item.Type = typ
	item.Name = name
	item.Description = description
	item.CategoryID = categoryID
	item.Currency = currency
	item.BasePrice = basePrice
	item.IsActive = isActive
	item.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return nil, err
	}
	if s.bus != nil && actorID != "" {
		s.bus.Publish("audit.change", map[string]any{
			"actor_id": actorID, "module": "catalog", "action": "update",
			"entity_type": "item", "entity_id": item.ID, "after": item,
		})
	}
	return item, nil
}

func (s *Service) DeleteItem(ctx context.Context, id string) error {
	if _, err := s.repo.FindItem(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	inUse, err := s.repo.ItemReferencedInOrders(ctx, id)
	if err != nil {
		return err
	}
	if inUse {
		return ErrInUse
	}
	return s.repo.DeleteItem(ctx, id)
}

func (s *Service) UpdateCategory(ctx context.Context, id, name, parentID string) (*domain.Category, error) {
	c, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidInput
	}
	c.Name = name
	c.ParentID = parentID
	c.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateCategory(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCategory(ctx context.Context, id string) error {
	if _, err := s.repo.FindCategoryByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	count, err := s.repo.CountItemsByCategory(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrInUse
	}
	return s.repo.DeleteCategory(ctx, id)
}
