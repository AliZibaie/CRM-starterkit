package domain

import "context"

type Repository interface {
	Add(ctx context.Context, item *WishlistItem) error
	Remove(ctx context.Context, userID, itemID string) error
	List(ctx context.Context, userID string) ([]WishlistItem, error)
}
