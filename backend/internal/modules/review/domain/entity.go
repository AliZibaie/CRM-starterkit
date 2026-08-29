package domain

import "time"

type Review struct {
	ID        string    `json:"id"`
	ItemID    string    `json:"item_id"`
	UserID    string    `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	IsVisible bool      `json:"is_visible"`
	CreatedAt time.Time `json:"created_at"`
}

type ListFilter struct {
	ItemID string
	Rating *int
	Limit  int
}
