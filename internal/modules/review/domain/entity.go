package domain

import "time"

type Review struct {
	ID        string
	ItemID    string
	UserID    string
	Rating    int
	Comment   string
	IsVisible bool
	CreatedAt time.Time
}
