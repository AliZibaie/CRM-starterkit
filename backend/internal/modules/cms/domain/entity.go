package domain

import "time"

type Page struct {
	ID        string    `json:"id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Block struct {
	ID        string `json:"id"`
	PageID    string `json:"page_id"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	SortOrder int    `json:"sort_order"`
}

type Menu struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location,omitempty"`
	ItemsJSON string    `json:"items_json"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SEOMeta struct {
	ID          string `json:"id"`
	EntityType  string `json:"entity_type"`
	EntityID    string `json:"entity_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
}

type FAQ struct {
	ID        string `json:"id"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
}
