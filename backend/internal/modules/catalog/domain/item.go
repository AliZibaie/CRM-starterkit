package domain

import "time"

const (
	TypePhysical = "physical_product"
	TypeService  = "service"
	TypeCourse   = "course_license"
	TypeDigital  = "digital"
)

type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ParentID  string    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Item struct {
	ID          string    `json:"id"`
	CategoryID  string    `json:"category_id,omitempty"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	BasePrice   int64     `json:"base_price"`
	Currency    string    `json:"currency"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Attribute struct {
	ID        string    `json:"id"`
	ItemID    string    `json:"item_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

type Variant struct {
	ID        string    `json:"id"`
	ItemID    string    `json:"item_id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	AttrsJSON string    `json:"attrs_json,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Inventory struct {
	ID        string    `json:"id"`
	VariantID string    `json:"variant_id"`
	Quantity  int64     `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`
}
