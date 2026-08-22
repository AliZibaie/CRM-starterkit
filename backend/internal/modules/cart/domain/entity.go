package domain

import "time"

type Cart struct {
	ID        string
	UserID    string
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartItem struct {
	ID        string
	CartID    string
	ItemID    string
	VariantID string
	Title     string
	UnitPrice int64
	Quantity  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CatalogSnapshot struct {
	ItemID    string
	VariantID string
	Title     string
	UnitPrice int64
	Currency  string
	Available int64
}
