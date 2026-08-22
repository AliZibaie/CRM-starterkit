package domain

import "context"

type Repository interface {
	CreateCategory(ctx context.Context, c *Category) error
	ListCategories(ctx context.Context) ([]Category, error)

	CreateItem(ctx context.Context, item *Item) error
	UpdateItem(ctx context.Context, item *Item) error
	FindItem(ctx context.Context, id string) (*Item, error)
	FindItemBySlug(ctx context.Context, slug string) (*Item, error)
	ListItems(ctx context.Context) ([]Item, error)
	ListActiveItems(ctx context.Context) ([]Item, error)

	FindCategoryBySlug(ctx context.Context, slug string) (*Category, error)
	ListItemsByCategory(ctx context.Context, categoryID string) ([]Item, error)

	UpsertAttribute(ctx context.Context, a *Attribute) error
	ListAttributes(ctx context.Context, itemID string) ([]Attribute, error)

	CreateVariant(ctx context.Context, v *Variant) error
	FindVariant(ctx context.Context, id string) (*Variant, error)
	ListVariants(ctx context.Context, itemID string) ([]Variant, error)

	SetInventory(ctx context.Context, inv *Inventory) error
	GetInventory(ctx context.Context, variantID string) (*Inventory, error)
	AdjustInventory(ctx context.Context, variantID string, delta int64) error
}
