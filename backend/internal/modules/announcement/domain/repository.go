package domain

import "context"

type Repository interface {
	Create(ctx context.Context, a *Announcement) error
	Update(ctx context.Context, a *Announcement) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*Announcement, error)
	List(ctx context.Context) ([]Announcement, error)
}
