package domain

import "time"

const (
	TypePhysical = "physical_product"
	TypeService  = "service"
	TypeCourse   = "course_license"
	TypeDigital  = "digital"
)

type Category struct {
	ID        string
	Name      string
	Slug      string
	ParentID  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Item struct {
	ID          string
	CategoryID  string
	Type        string
	Name        string
	Slug        string
	Description string
	BasePrice   int64
	Currency    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Attribute struct {
	ID        string
	ItemID    string
	Key       string
	Value     string
	CreatedAt time.Time
}

type Variant struct {
	ID        string
	ItemID    string
	SKU       string
	Name      string
	Price     int64
	AttrsJSON string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Inventory struct {
	ID        string
	VariantID string
	Quantity  int64
	UpdatedAt time.Time
}
