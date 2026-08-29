package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-starterkit/internal/modules/order/domain"

	"gorm.io/gorm"
)

type OrderModel struct {
	ID             string    `gorm:"primaryKey;size:36"`
	UserID         string    `gorm:"size:36;index;not null"`
	Status         string    `gorm:"size:30;index;not null"`
	Currency       string    `gorm:"size:10;not null"`
	Subtotal       int64     `gorm:"not null"`
	DiscountCode   string    `gorm:"size:50"`
	DiscountAmount int64     `gorm:"not null;default:0"`
	TotalAmount    int64     `gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

func (OrderModel) TableName() string { return "orders" }

type OrderItemModel struct {
	ID        string `gorm:"primaryKey;size:36"`
	OrderID   string `gorm:"size:36;index;not null"`
	ItemID    string `gorm:"size:36;not null"`
	VariantID string `gorm:"size:36"`
	Title     string `gorm:"size:255;not null"`
	UnitPrice int64  `gorm:"not null"`
	Quantity  int64  `gorm:"not null"`
	LineTotal int64  `gorm:"not null"`
}

func (OrderItemModel) TableName() string { return "order_items" }

type InvoiceModel struct {
	ID       string    `gorm:"primaryKey;size:36"`
	OrderID  string    `gorm:"uniqueIndex;size:36;not null"`
	Number   string    `gorm:"uniqueIndex;size:50;not null"`
	Amount   int64     `gorm:"not null"`
	Currency string    `gorm:"size:10;not null"`
	IssuedAt time.Time `gorm:"not null"`
	Status   string    `gorm:"size:30;not null"`
}

func (InvoiceModel) TableName() string { return "invoices" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error {
	return s.db.AutoMigrate(&OrderModel{}, &OrderItemModel{}, &InvoiceModel{})
}

func (s *Store) CreateOrder(ctx context.Context, order *domain.Order, items []domain.OrderItem, invoice *domain.Invoice) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&OrderModel{
			ID: order.ID, UserID: order.UserID, Status: order.Status, Currency: order.Currency,
			Subtotal: order.Subtotal, DiscountCode: order.DiscountCode, DiscountAmount: order.DiscountAmount,
			TotalAmount: order.TotalAmount, CreatedAt: order.CreatedAt, UpdatedAt: order.UpdatedAt,
		}).Error; err != nil {
			return err
		}
		for _, it := range items {
			if err := tx.Create(&OrderItemModel{
				ID: it.ID, OrderID: it.OrderID, ItemID: it.ItemID, VariantID: it.VariantID, Title: it.Title,
				UnitPrice: it.UnitPrice, Quantity: it.Quantity, LineTotal: it.LineTotal,
			}).Error; err != nil {
				return err
			}
		}
		return tx.Create(&InvoiceModel{
			ID: invoice.ID, OrderID: invoice.OrderID, Number: invoice.Number, Amount: invoice.Amount,
			Currency: invoice.Currency, IssuedAt: invoice.IssuedAt, Status: invoice.Status,
		}).Error
	})
}

func (s *Store) UpdateOrder(ctx context.Context, order *domain.Order) error {
	return s.db.WithContext(ctx).Save(&OrderModel{
		ID: order.ID, UserID: order.UserID, Status: order.Status, Currency: order.Currency,
		Subtotal: order.Subtotal, DiscountCode: order.DiscountCode, DiscountAmount: order.DiscountAmount,
		TotalAmount: order.TotalAmount, CreatedAt: order.CreatedAt, UpdatedAt: order.UpdatedAt,
	}).Error
}

func (s *Store) FindOrder(ctx context.Context, id string) (*domain.Order, error) {
	var m OrderModel
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return orderDomain(&m), nil
}

func (s *Store) ListOrdersByUser(ctx context.Context, userID string) ([]domain.Order, error) {
	return s.ListOrders(ctx, "", userID, nil, nil)
}

func (s *Store) ListOrders(ctx context.Context, status, userID string, from, to *time.Time) ([]domain.Order, error) {
	q := s.db.WithContext(ctx).Model(&OrderModel{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}
	if from != nil {
		q = q.Where("created_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("created_at <= ?", *to)
	}
	var rows []OrderModel
	if err := q.Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Order, 0, len(rows))
	for i := range rows {
		out = append(out, *orderDomain(&rows[i]))
	}
	return out, nil
}

func (s *Store) ListOrderItems(ctx context.Context, orderID string) ([]domain.OrderItem, error) {
	var rows []OrderItemModel
	if err := s.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.OrderItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.OrderItem{ID: r.ID, OrderID: r.OrderID, ItemID: r.ItemID, VariantID: r.VariantID, Title: r.Title, UnitPrice: r.UnitPrice, Quantity: r.Quantity, LineTotal: r.LineTotal})
	}
	return out, nil
}

func (s *Store) FindInvoiceByOrder(ctx context.Context, orderID string) (*domain.Invoice, error) {
	var m InvoiceModel
	if err := s.db.WithContext(ctx).Where("order_id = ?", orderID).First(&m).Error; err != nil {
		return nil, err
	}
	return &domain.Invoice{ID: m.ID, OrderID: m.OrderID, Number: m.Number, Amount: m.Amount, Currency: m.Currency, IssuedAt: m.IssuedAt, Status: m.Status}, nil
}

func (s *Store) UpdateInvoice(ctx context.Context, inv *domain.Invoice) error {
	return s.db.WithContext(ctx).Save(&InvoiceModel{
		ID: inv.ID, OrderID: inv.OrderID, Number: inv.Number, Amount: inv.Amount, Currency: inv.Currency, IssuedAt: inv.IssuedAt, Status: inv.Status,
	}).Error
}

func (s *Store) LoadUserCart(ctx context.Context, userID string) (string, string, []domain.CartLine, error) {
	var cart struct {
		ID       string
		Currency string
	}
	if err := s.db.WithContext(ctx).Table("carts").Select("id, currency").Where("user_id = ?", userID).Take(&cart).Error; err != nil {
		return "", "", nil, err
	}
	var rows []struct {
		ItemID    string
		VariantID string
		Title     string
		UnitPrice int64
		Quantity  int64
	}
	if err := s.db.WithContext(ctx).Table("cart_items").Where("cart_id = ?", cart.ID).Find(&rows).Error; err != nil {
		return "", "", nil, err
	}
	lines := make([]domain.CartLine, 0, len(rows))
	for _, r := range rows {
		lines = append(lines, domain.CartLine{ItemID: r.ItemID, VariantID: r.VariantID, Title: r.Title, UnitPrice: r.UnitPrice, Quantity: r.Quantity})
	}
	return cart.ID, cart.Currency, lines, nil
}

func (s *Store) ClearCart(ctx context.Context, cartID string) error {
	return s.db.WithContext(ctx).Exec("DELETE FROM cart_items WHERE cart_id = ?", cartID).Error
}

func (s *Store) PreviewDiscount(ctx context.Context, code string, amount int64) (int64, error) {
	return s.calcDiscount(ctx, code, amount, false)
}

func (s *Store) ConsumeDiscount(ctx context.Context, code string, amount int64) (int64, error) {
	return s.calcDiscount(ctx, code, amount, true)
}

func (s *Store) calcDiscount(ctx context.Context, code string, amount int64, consume bool) (int64, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return 0, nil
	}
	var d struct {
		ID             string
		Type           string
		Value          int64
		MinOrderAmount int64
		MaxUses        int64
		UsedCount      int64
		IsActive       bool
		StartsAt       *time.Time
		EndsAt         *time.Time
	}
	if err := s.db.WithContext(ctx).Table("discount_codes").Where("code = ?", code).Take(&d).Error; err != nil {
		return 0, err
	}
	now := time.Now().UTC()
	if !d.IsActive || (d.StartsAt != nil && now.Before(*d.StartsAt)) || (d.EndsAt != nil && now.After(*d.EndsAt)) {
		return 0, errors.New("discount not applicable")
	}
	if d.MaxUses > 0 && d.UsedCount >= d.MaxUses {
		return 0, errors.New("discount not applicable")
	}
	if amount < d.MinOrderAmount {
		return 0, errors.New("discount not applicable")
	}
	var discount int64
	if d.Type == "percent" {
		discount = amount * d.Value / 100
	} else {
		discount = d.Value
	}
	if discount > amount {
		discount = amount
	}
	if consume {
		if err := s.db.WithContext(ctx).Table("discount_codes").Where("id = ?", d.ID).Updates(map[string]any{
			"used_count": gorm.Expr("used_count + 1"),
			"updated_at": now,
		}).Error; err != nil {
			return 0, err
		}
	}
	return discount, nil
}

func (s *Store) AdjustInventory(ctx context.Context, variantID string, delta int64) error {
	if variantID == "" {
		return nil
	}
	return s.db.WithContext(ctx).Exec("UPDATE catalog_inventory SET quantity = quantity + ?, updated_at = ? WHERE variant_id = ?", delta, time.Now().UTC(), variantID).Error
}

func orderDomain(m *OrderModel) *domain.Order {
	return &domain.Order{
		ID: m.ID, UserID: m.UserID, Status: m.Status, Currency: m.Currency, Subtotal: m.Subtotal,
		DiscountCode: m.DiscountCode, DiscountAmount: m.DiscountAmount, TotalAmount: m.TotalAmount,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}
