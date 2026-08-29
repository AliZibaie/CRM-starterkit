package domain

import "time"

type Banner struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	ImageURL  string     `json:"image_url"`
	LinkURL   string     `json:"link_url,omitempty"`
	Position  string     `json:"position"`
	IsActive  bool       `json:"is_active"`
	StartsAt  *time.Time `json:"starts_at,omitempty"`
	EndsAt    *time.Time `json:"ends_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
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
