package domain

import "time"

type Banner struct {
	ID        string
	Title     string
	ImageURL  string
	LinkURL   string
	Position  string
	IsActive  bool
	StartsAt  *time.Time
	EndsAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateInput struct {
	Title    string
	ImageURL string
	LinkURL  string
	Position string
}

type ScheduleInput struct {
	StartsAt *time.Time
	EndsAt   *time.Time
	IsActive *bool
}
