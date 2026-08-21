package domain

import "context"

type Repository interface {
	Create(ctx context.Context, b *Banner) error
	Update(ctx context.Context, b *Banner) error
	FindByID(ctx context.Context, id string) (*Banner, error)
	List(ctx context.Context) ([]Banner, error)
	ListByPosition(ctx context.Context, position string) ([]Banner, error)
	Delete(ctx context.Context, id string) error
}
