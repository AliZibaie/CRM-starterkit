package domain

import (
	"context"
	"time"
)

type Filter struct {
	ActorID    string
	Module     string
	EntityType string
	EntityID   string
	From       *time.Time
	To         *time.Time
	Limit      int
}

type Repository interface {
	Create(ctx context.Context, e *Entry) error
	List(ctx context.Context, f Filter) ([]Entry, error)
}
