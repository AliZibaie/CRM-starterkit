package domain

import "time"

const (
	StatusPending   = "pending"
	StatusPaid      = "paid"
	StatusRefunded  = "refunded"
	StatusCancelled = "cancelled"

	EventCreated  = "commerce.order.created"
	EventPaid     = "commerce.order.paid"
	EventRefunded = "commerce.order.refunded"
	EventPaymentVerified = "commerce.payment.verified"
)

type Order struct {
	ID             string
	UserID         string
	Status         string
	Currency       string
	Subtotal       int64
	DiscountCode   string
	DiscountAmount int64
	TotalAmount    int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type OrderItem struct {
	ID        string
	OrderID   string
	ItemID    string
	VariantID string
	Title     string
	UnitPrice int64
	Quantity  int64
	LineTotal int64
}

type Invoice struct {
	ID          string
	OrderID     string
	Number      string
	Amount      int64
	Currency    string
	IssuedAt    time.Time
	Status      string
}

type CartLine struct {
	ItemID    string
	VariantID string
	Title     string
	UnitPrice int64
	Quantity  int64
}
