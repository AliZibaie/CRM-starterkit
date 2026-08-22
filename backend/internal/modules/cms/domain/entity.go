package domain

import "time"

type Page struct {
	ID        string
	Slug      string
	Title     string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Block struct {
	ID        string
	PageID    string
	Type      string
	Content   string
	SortOrder int
}

type Menu struct {
	ID        string
	Name      string
	Location  string
	ItemsJSON string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SEOMeta struct {
	ID          string
	EntityType  string
	EntityID    string
	Title       string
	Description string
	Keywords    string
}

type FAQ struct {
	ID        string
	Question  string
	Answer    string
	SortOrder int
	IsActive  bool
}
