package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/catalog/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryModel struct {
	ID        string    `gorm:"primaryKey;size:36"`
	Name      string    `gorm:"size:150;not null"`
	Slug      string    `gorm:"uniqueIndex;size:150;not null"`
	ParentID  string    `gorm:"size:36;index"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (CategoryModel) TableName() string { return "catalog_categories" }

type ItemModel struct {
	ID          string    `gorm:"primaryKey;size:36"`
	CategoryID  string    `gorm:"size:36;index"`
	Type        string    `gorm:"size:50;index;not null"`
	Name        string    `gorm:"size:255;not null"`
	Slug        string    `gorm:"uniqueIndex;size:255;not null"`
	Description string    `gorm:"type:text"`
	BasePrice   int64     `gorm:"not null;default:0"`
	Currency    string    `gorm:"size:10;not null;default:IRR"`
	IsActive    bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (ItemModel) TableName() string { return "catalog_items" }

type AttributeModel struct {
	ID        string    `gorm:"primaryKey;size:36"`
	ItemID    string    `gorm:"size:36;uniqueIndex:idx_item_attr;not null"`
	Key       string    `gorm:"size:100;uniqueIndex:idx_item_attr;not null"`
	Value     string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (AttributeModel) TableName() string { return "catalog_attributes" }

type VariantModel struct {
	ID        string    `gorm:"primaryKey;size:36"`
	ItemID    string    `gorm:"size:36;index;not null"`
	SKU       string    `gorm:"uniqueIndex;size:100;not null"`
	Name      string    `gorm:"size:150;not null"`
	Price     int64     `gorm:"not null;default:0"`
	AttrsJSON string    `gorm:"type:text"`
	IsActive  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (VariantModel) TableName() string { return "catalog_variants" }

type InventoryModel struct {
	ID        string    `gorm:"primaryKey;size:36"`
	VariantID string    `gorm:"uniqueIndex;size:36;not null"`
	Quantity  int64     `gorm:"not null;default:0"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (InventoryModel) TableName() string { return "catalog_inventory" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error {
	return s.db.AutoMigrate(&CategoryModel{}, &ItemModel{}, &AttributeModel{}, &VariantModel{}, &InventoryModel{})
}

func (s *Store) CreateCategory(ctx context.Context, c *domain.Category) error {
	return s.db.WithContext(ctx).Create(&CategoryModel{ID: c.ID, Name: c.Name, Slug: c.Slug, ParentID: c.ParentID, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}).Error
}

func (s *Store) ListCategories(ctx context.Context) ([]domain.Category, error) {
	var rows []CategoryModel
	if err := s.db.WithContext(ctx).Order("name asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Category, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.Category{ID: r.ID, Name: r.Name, Slug: r.Slug, ParentID: r.ParentID, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt})
	}
	return out, nil
}

func (s *Store) CreateItem(ctx context.Context, item *domain.Item) error {
	return s.db.WithContext(ctx).Create(itemModel(item)).Error
}

func (s *Store) UpdateItem(ctx context.Context, item *domain.Item) error {
	return s.db.WithContext(ctx).Save(itemModel(item)).Error
}

func (s *Store) FindItemBySlug(ctx context.Context, slug string) (*domain.Item, error) {
	var m ItemModel
	if err := s.db.WithContext(ctx).First(&m, "slug = ? AND is_active = ?", slug, true).Error; err != nil {
		return nil, err
	}
	return itemDomain(&m), nil
}

func (s *Store) ListActiveItems(ctx context.Context) ([]domain.Item, error) {
	var rows []ItemModel
	if err := s.db.WithContext(ctx).Where("is_active = ?", true).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Item, 0, len(rows))
	for i := range rows {
		out = append(out, *itemDomain(&rows[i]))
	}
	return out, nil
}

func (s *Store) FindCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var m CategoryModel
	if err := s.db.WithContext(ctx).First(&m, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &domain.Category{ID: m.ID, Name: m.Name, Slug: m.Slug, ParentID: m.ParentID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (s *Store) ListItemsByCategory(ctx context.Context, categoryID string) ([]domain.Item, error) {
	var rows []ItemModel
	q := s.db.WithContext(ctx).Where("is_active = ?", true)
	if categoryID != "" {
		q = q.Where("category_id = ?", categoryID)
	}
	if err := q.Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Item, 0, len(rows))
	for i := range rows {
		out = append(out, *itemDomain(&rows[i]))
	}
	return out, nil
}

func (s *Store) FindItem(ctx context.Context, id string) (*domain.Item, error) {
	var m ItemModel
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return itemDomain(&m), nil
}

func (s *Store) ListItems(ctx context.Context) ([]domain.Item, error) {
	var rows []ItemModel
	if err := s.db.WithContext(ctx).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Item, 0, len(rows))
	for i := range rows {
		out = append(out, *itemDomain(&rows[i]))
	}
	return out, nil
}

func (s *Store) UpsertAttribute(ctx context.Context, a *domain.Attribute) error {
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "item_id"}, {Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&AttributeModel{ID: a.ID, ItemID: a.ItemID, Key: a.Key, Value: a.Value, CreatedAt: a.CreatedAt}).Error
}

func (s *Store) ListAttributes(ctx context.Context, itemID string) ([]domain.Attribute, error) {
	var rows []AttributeModel
	if err := s.db.WithContext(ctx).Where("item_id = ?", itemID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Attribute, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.Attribute{ID: r.ID, ItemID: r.ItemID, Key: r.Key, Value: r.Value, CreatedAt: r.CreatedAt})
	}
	return out, nil
}

func (s *Store) CreateVariant(ctx context.Context, v *domain.Variant) error {
	return s.db.WithContext(ctx).Create(&VariantModel{
		ID: v.ID, ItemID: v.ItemID, SKU: v.SKU, Name: v.Name, Price: v.Price, AttrsJSON: v.AttrsJSON,
		IsActive: v.IsActive, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt,
	}).Error
}

func (s *Store) FindVariant(ctx context.Context, id string) (*domain.Variant, error) {
	var m VariantModel
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &domain.Variant{ID: m.ID, ItemID: m.ItemID, SKU: m.SKU, Name: m.Name, Price: m.Price, AttrsJSON: m.AttrsJSON, IsActive: m.IsActive, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (s *Store) ListVariants(ctx context.Context, itemID string) ([]domain.Variant, error) {
	var rows []VariantModel
	if err := s.db.WithContext(ctx).Where("item_id = ?", itemID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Variant, 0, len(rows))
	for _, m := range rows {
		out = append(out, domain.Variant{ID: m.ID, ItemID: m.ItemID, SKU: m.SKU, Name: m.Name, Price: m.Price, AttrsJSON: m.AttrsJSON, IsActive: m.IsActive, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt})
	}
	return out, nil
}

func (s *Store) SetInventory(ctx context.Context, inv *domain.Inventory) error {
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "variant_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity", "updated_at"}),
	}).Create(&InventoryModel{ID: inv.ID, VariantID: inv.VariantID, Quantity: inv.Quantity, UpdatedAt: inv.UpdatedAt}).Error
}

func (s *Store) GetInventory(ctx context.Context, variantID string) (*domain.Inventory, error) {
	var m InventoryModel
	if err := s.db.WithContext(ctx).First(&m, "variant_id = ?", variantID).Error; err != nil {
		return nil, err
	}
	return &domain.Inventory{ID: m.ID, VariantID: m.VariantID, Quantity: m.Quantity, UpdatedAt: m.UpdatedAt}, nil
}

func (s *Store) AdjustInventory(ctx context.Context, variantID string, delta int64) error {
	return s.db.WithContext(ctx).Model(&InventoryModel{}).Where("variant_id = ?", variantID).
		UpdateColumn("quantity", gorm.Expr("quantity + ?", delta)).Error
}

func itemModel(i *domain.Item) *ItemModel {
	return &ItemModel{ID: i.ID, CategoryID: i.CategoryID, Type: i.Type, Name: i.Name, Slug: i.Slug, Description: i.Description, BasePrice: i.BasePrice, Currency: i.Currency, IsActive: i.IsActive, CreatedAt: i.CreatedAt, UpdatedAt: i.UpdatedAt}
}

func itemDomain(m *ItemModel) *domain.Item {
	return &domain.Item{ID: m.ID, CategoryID: m.CategoryID, Type: m.Type, Name: m.Name, Slug: m.Slug, Description: m.Description, BasePrice: m.BasePrice, Currency: m.Currency, IsActive: m.IsActive, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}
