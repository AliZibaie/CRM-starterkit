package domain

import "context"

type Repository interface {
	Create(ctx context.Context, e *RequestEntry) error
	List(ctx context.Context, f ListFilter) ([]RequestEntry, error)
	Stats(ctx context.Context, f ListFilter) (*Stats, error)
}
