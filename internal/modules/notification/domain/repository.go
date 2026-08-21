package domain

import "context"

type Repository interface {
	Create(ctx context.Context, n *Notification) error
	ListByUser(ctx context.Context, userID string, limit int) ([]Notification, error)
	MarkRead(ctx context.Context, id, userID string) error
}
