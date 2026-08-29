package domain

import "context"

type Repository interface {
	Create(ctx context.Context, r *Review) error
	ListByItem(ctx context.Context, itemID string) ([]Review, error)
	ListAll(ctx context.Context, f ListFilter) ([]Review, error)
	AverageRating(ctx context.Context, itemID string) (float64, error)
	Delete(ctx context.Context, id, userID string) error
	AdminDelete(ctx context.Context, id string) error
	SetVisibility(ctx context.Context, id string, visible bool) error
}
