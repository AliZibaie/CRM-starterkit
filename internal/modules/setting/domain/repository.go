package domain

import "context"

type Repository interface {
	Get(ctx context.Context, key string) (*Setting, error)
	Set(ctx context.Context, s *Setting) error
	List(ctx context.Context) ([]Setting, error)
}
