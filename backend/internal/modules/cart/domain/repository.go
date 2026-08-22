package domain

import "context"

type Repository interface {
	GetOrCreateCart(ctx context.Context, userID, currency string) (*Cart, error)
	FindCartByUser(ctx context.Context, userID string) (*Cart, error)
	FindCart(ctx context.Context, id string) (*Cart, error)
	ListItems(ctx context.Context, cartID string) ([]CartItem, error)
	UpsertItem(ctx context.Context, item *CartItem) error
	RemoveItem(ctx context.Context, cartID, itemID string) error
	Clear(ctx context.Context, cartID string) error
	LookupCatalog(ctx context.Context, itemID, variantID string) (*CatalogSnapshot, error)
}
