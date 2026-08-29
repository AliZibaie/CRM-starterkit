package domain

import "context"

type Repository interface {
	Create(ctx context.Context, p *Payment) error
	Update(ctx context.Context, p *Payment) error
	FindByID(ctx context.Context, id string) (*Payment, error)
	FindByAuthority(ctx context.Context, authority string) (*Payment, error)
	FindByOrder(ctx context.Context, orderID string) (*Payment, error)
	FindOrderAmount(ctx context.Context, orderID string) (amount int64, currency string, status string, err error)
	ListByStatus(ctx context.Context, status string, limit int) ([]Payment, error)
}
