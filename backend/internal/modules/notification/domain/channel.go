package domain

import "context"

type Payload struct {
	UserID  string
	To      string
	Title   string
	Body    string
	Channel string
	Data    map[string]any
}

type Channel interface {
	Name() string
	Send(ctx context.Context, payload Payload) error
}

type Notification struct {
	ID      string
	UserID  string
	Channel string
	Title   string
	Body    string
	Status  string
	MetaJSON string
}
