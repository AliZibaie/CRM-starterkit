package repository

import (
	"context"
	"errors"
	"time"

	"crm-starterkit/internal/modules/cart/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartModel struct {
	ID        string    `gorm:"primaryKey;size:36"`
	UserID    string    `gorm:"uniqueIndex;size:36;not null"`
	Currency  string    `gorm:"size:10;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (CartModel) TableName() string { return "carts" }

type CartItemModel struct {
	ID        string    `gorm:"primaryKey;size:36"`
	CartID    string    `gorm:"size:36;uniqueIndex:idx_cart_variant;not null"`
	ItemID    string    `gorm:"size:36;not null"`
	VariantID string    `gorm:"size:36;uniqueIndex:idx_cart_variant;not null"`
	Title     string    `gorm:"size:255;not null"`
	UnitPrice int64     `gorm:"not null"`
	Quantity  int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (CartItemModel) TableName() string { return "cart_items" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error {
	return s.db.AutoMigrate(&CartModel{}, &CartItemModel{})
}

func (s *Store) GetOrCreateCart(ctx context.Context, userID, currency string) (*domain.Cart, error) {
	var m CartModel
	err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&m).Error
	if err == nil {
		return &domain.Cart{ID: m.ID, UserID: m.UserID, Currency: m.Currency, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	now := time.Now().UTC()
	m = CartModel{ID: uuid.NewString(), UserID: userID, Currency: currency, CreatedAt: now, UpdatedAt: now}
	if err := s.db.WithContext(ctx).Create(&m).Error; err != nil {
		return nil, err
	}
	return &domain.Cart{ID: m.ID, UserID: m.UserID, Currency: m.Currency, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (s *Store) FindCartByUser(ctx context.Context, userID string) (*domain.Cart, error) {
	var m CartModel
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&m).Error; err != nil {
		return nil, err
	}
	return &domain.Cart{ID: m.ID, UserID: m.UserID, Currency: m.Currency, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (s *Store) FindCart(ctx context.Context, id string) (*domain.Cart, error) {
	var m CartModel
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &domain.Cart{ID: m.ID, UserID: m.UserID, Currency: m.Currency, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (s *Store) ListItems(ctx context.Context, cartID string) ([]domain.CartItem, error) {
	var rows []CartItemModel
	if err := s.db.WithContext(ctx).Where("cart_id = ?", cartID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.CartItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.CartItem{ID: r.ID, CartID: r.CartID, ItemID: r.ItemID, VariantID: r.VariantID, Title: r.Title, UnitPrice: r.UnitPrice, Quantity: r.Quantity, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt})
	}
	return out, nil
}

func (s *Store) UpsertItem(ctx context.Context, item *domain.CartItem) error {
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "cart_id"}, {Name: "variant_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity", "unit_price", "title", "updated_at"}),
	}).Create(&CartItemModel{
		ID: item.ID, CartID: item.CartID, ItemID: item.ItemID, VariantID: item.VariantID, Title: item.Title,
		UnitPrice: item.UnitPrice, Quantity: item.Quantity, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
	}).Error
}

func (s *Store) RemoveItem(ctx context.Context, cartID, itemID string) error {
	return s.db.WithContext(ctx).Where("cart_id = ? AND id = ?", cartID, itemID).Delete(&CartItemModel{}).Error
}

func (s *Store) Clear(ctx context.Context, cartID string) error {
	return s.db.WithContext(ctx).Where("cart_id = ?", cartID).Delete(&CartItemModel{}).Error
}

func (s *Store) LookupCatalog(ctx context.Context, itemID, variantID string) (*domain.CatalogSnapshot, error) {
	type row struct {
		ItemID    string
		VariantID string
		Title     string
		UnitPrice int64
		Currency  string
		Available int64
	}
	var r row
	q := `
SELECT i.id as item_id, v.id as variant_id,
       COALESCE(NULLIF(v.name,''), i.name) as title,
       CASE WHEN v.id IS NULL OR v.id = '' THEN i.base_price ELSE v.price END as unit_price,
       i.currency as currency,
       COALESCE(inv.quantity, 0) as available
FROM catalog_items i
LEFT JOIN catalog_variants v ON v.item_id = i.id AND v.id = ?
LEFT JOIN catalog_inventory inv ON inv.variant_id = v.id
WHERE i.id = ? AND i.is_active = ?`
	if err := s.db.WithContext(ctx).Raw(q, variantID, itemID, true).Scan(&r).Error; err != nil {
		return nil, err
	}
	if r.ItemID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	if variantID != "" && r.VariantID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &domain.CatalogSnapshot{ItemID: r.ItemID, VariantID: r.VariantID, Title: r.Title, UnitPrice: r.UnitPrice, Currency: r.Currency, Available: r.Available}, nil
}
