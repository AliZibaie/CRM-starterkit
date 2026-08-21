package domain

import "context"

type Repository interface {
	Create(ctx context.Context, m *Media) error
	FindByID(ctx context.Context, id string) (*Media, error)
	ListByOwner(ctx context.Context, ownerID string) ([]Media, error)
	Delete(ctx context.Context, id string) error
}
