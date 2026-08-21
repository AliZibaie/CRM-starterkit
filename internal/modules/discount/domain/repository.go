package domain

import "context"

type Repository interface {
	Create(ctx context.Context, d *Discount) error
	FindByCode(ctx context.Context, code string) (*Discount, error)
	FindByID(ctx context.Context, id string) (*Discount, error)
	List(ctx context.Context) ([]Discount, error)
	IncrementUsage(ctx context.Context, id string) error
}
