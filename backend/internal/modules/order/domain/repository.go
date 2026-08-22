package domain

import "context"

type Repository interface {
	CreateOrder(ctx context.Context, order *Order, items []OrderItem, invoice *Invoice) error
	UpdateOrder(ctx context.Context, order *Order) error
	FindOrder(ctx context.Context, id string) (*Order, error)
	ListOrdersByUser(ctx context.Context, userID string) ([]Order, error)
	ListOrderItems(ctx context.Context, orderID string) ([]OrderItem, error)
	FindInvoiceByOrder(ctx context.Context, orderID string) (*Invoice, error)
	UpdateInvoice(ctx context.Context, inv *Invoice) error

	LoadUserCart(ctx context.Context, userID string) (cartID string, currency string, lines []CartLine, err error)
	ClearCart(ctx context.Context, cartID string) error
	PreviewDiscount(ctx context.Context, code string, amount int64) (discountAmount int64, err error)
	ConsumeDiscount(ctx context.Context, code string, amount int64) (discountAmount int64, err error)
	AdjustInventory(ctx context.Context, variantID string, delta int64) error
}
