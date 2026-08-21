package domain

import "time"

type WishlistItem struct {
	ID        string
	UserID    string
	ItemID    string
	CreatedAt time.Time
}
