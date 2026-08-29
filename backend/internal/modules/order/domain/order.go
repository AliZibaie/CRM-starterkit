package domain

import "time"

const (
	StatusPending    = "pending"
	StatusPaid       = "paid"
	StatusProcessing = "processing"
	StatusShipped    = "shipped"
	StatusDelivered  = "delivered"
	StatusRefunded   = "refunded"
	StatusCancelled  = "cancelled"

	EventCreated         = "commerce.order.created"
	EventPaid            = "commerce.order.paid"
	EventRefunded        = "commerce.order.refunded"
	EventPaymentVerified = "commerce.payment.verified"
)

type Order struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Status         string    `json:"status"`
	Currency       string    `json:"currency"`
	Subtotal       int64     `json:"subtotal"`
	DiscountCode   string    `json:"discount_code,omitempty"`
	DiscountAmount int64     `json:"discount_amount"`
	TotalAmount    int64     `json:"total_amount"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID        string `json:"id"`
	OrderID   string `json:"order_id"`
	ItemID    string `json:"item_id"`
	VariantID string `json:"variant_id,omitempty"`
	Title     string `json:"title"`
	UnitPrice int64  `json:"unit_price"`
	Quantity  int64  `json:"quantity"`
	LineTotal int64  `json:"line_total"`
}

type Invoice struct {
	ID       string    `json:"id"`
	OrderID  string    `json:"order_id"`
	Number   string    `json:"number"`
	Amount   int64     `json:"amount"`
	Currency string    `json:"currency"`
	IssuedAt time.Time `json:"issued_at"`
	Status   string    `json:"status"`
}

type CartLine struct {
	ItemID    string
	VariantID string
	Title     string
	UnitPrice int64
	Quantity  int64
}
