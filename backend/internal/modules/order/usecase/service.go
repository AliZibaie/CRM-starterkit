package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/modules/order/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("order not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrEmptyCart    = errors.New("cart is empty")
	ErrBadStatus    = errors.New("invalid order status")
)

type Service struct {
	repo domain.Repository
	bus  *event.Bus
}

func NewService(repo domain.Repository, bus *event.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

func (s *Service) CreateFromCart(ctx context.Context, userID, discountCode string) (*domain.Order, []domain.OrderItem, *domain.Invoice, error) {
	cartID, currency, lines, err := s.repo.LoadUserCart(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, ErrEmptyCart
		}
		return nil, nil, nil, err
	}
	if len(lines) == 0 {
		return nil, nil, nil, ErrEmptyCart
	}
	var subtotal int64
	items := make([]domain.OrderItem, 0, len(lines))
	orderID := uuid.NewString()
	for _, line := range lines {
		lt := line.UnitPrice * line.Quantity
		subtotal += lt
		items = append(items, domain.OrderItem{
			ID: uuid.NewString(), OrderID: orderID, ItemID: line.ItemID, VariantID: line.VariantID,
			Title: line.Title, UnitPrice: line.UnitPrice, Quantity: line.Quantity, LineTotal: lt,
		})
	}
	discountAmount := int64(0)
	if discountCode != "" {
		discountAmount, err = s.repo.ConsumeDiscount(ctx, discountCode, subtotal)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	total := subtotal - discountAmount
	now := time.Now().UTC()
	order := &domain.Order{
		ID: orderID, UserID: userID, Status: domain.StatusPending, Currency: currency,
		Subtotal: subtotal, DiscountCode: discountCode, DiscountAmount: discountAmount,
		TotalAmount: total, CreatedAt: now, UpdatedAt: now,
	}
	invoice := &domain.Invoice{
		ID: uuid.NewString(), OrderID: orderID, Number: fmt.Sprintf("INV-%d", now.Unix()),
		Amount: total, Currency: currency, IssuedAt: now, Status: "issued",
	}
	if err := s.repo.CreateOrder(ctx, order, items, invoice); err != nil {
		return nil, nil, nil, err
	}
	for _, it := range items {
		_ = s.repo.AdjustInventory(ctx, it.VariantID, -it.Quantity)
	}
	_ = s.repo.ClearCart(ctx, cartID)
	if s.bus != nil {
		s.bus.Publish(domain.EventCreated, map[string]any{"order_id": order.ID, "user_id": userID, "total": total})
	}
	return order, items, invoice, nil
}

func (s *Service) ApplyDiscount(ctx context.Context, orderID, code string) (*domain.Order, error) {
	order, err := s.repo.FindOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if order.Status != domain.StatusPending {
		return nil, ErrBadStatus
	}
	amount, err := s.repo.ConsumeDiscount(ctx, code, order.Subtotal)
	if err != nil {
		return nil, err
	}
	order.DiscountCode = code
	order.DiscountAmount = amount
	order.TotalAmount = order.Subtotal - amount
	order.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return nil, err
	}
	if inv, err := s.repo.FindInvoiceByOrder(ctx, order.ID); err == nil {
		inv.Amount = order.TotalAmount
		_ = s.repo.UpdateInvoice(ctx, inv)
	}
	return order, nil
}

func (s *Service) Checkout(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := s.repo.FindOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if order.Status != domain.StatusPending {
		return nil, ErrBadStatus
	}
	if s.bus != nil {
		s.bus.Publish("commerce.order.checkout_requested", map[string]any{
			"order_id": order.ID, "amount": order.TotalAmount, "currency": order.Currency, "user_id": order.UserID,
		})
	}
	return order, nil
}

func (s *Service) MarkPaid(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := s.repo.FindOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	order.Status = domain.StatusPaid
	order.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return nil, err
	}
	if inv, err := s.repo.FindInvoiceByOrder(ctx, order.ID); err == nil {
		inv.Status = "paid"
		_ = s.repo.UpdateInvoice(ctx, inv)
	}
	if s.bus != nil {
		s.bus.Publish(domain.EventPaid, map[string]any{"order_id": order.ID, "user_id": order.UserID})
	}
	return order, nil
}

func (s *Service) Refund(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := s.repo.FindOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if order.Status != domain.StatusPaid {
		return nil, ErrBadStatus
	}
	items, err := s.repo.ListOrderItems(ctx, orderID)
	if err != nil {
		return nil, err
	}
	for _, it := range items {
		_ = s.repo.AdjustInventory(ctx, it.VariantID, it.Quantity)
	}
	order.Status = domain.StatusRefunded
	order.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return nil, err
	}
	if inv, err := s.repo.FindInvoiceByOrder(ctx, order.ID); err == nil {
		inv.Status = "refunded"
		_ = s.repo.UpdateInvoice(ctx, inv)
	}
	if s.bus != nil {
		s.bus.Publish(domain.EventRefunded, map[string]any{"order_id": order.ID, "amount": order.TotalAmount})
	}
	return order, nil
}

func (s *Service) Get(ctx context.Context, id string) (*domain.Order, []domain.OrderItem, *domain.Invoice, error) {
	order, err := s.repo.FindOrder(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, ErrNotFound
		}
		return nil, nil, nil, err
	}
	items, err := s.repo.ListOrderItems(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	inv, err := s.repo.FindInvoiceByOrder(ctx, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil, err
	}
	return order, items, inv, nil
}

func (s *Service) ListMine(ctx context.Context, userID string) ([]domain.Order, error) {
	return s.repo.ListOrdersByUser(ctx, userID)
}

func (s *Service) ListAll(ctx context.Context, status, userID string, from, to *time.Time) ([]domain.Order, error) {
	return s.repo.ListOrders(ctx, status, userID, from, to)
}

func (s *Service) UpdateStatus(ctx context.Context, orderID, status string) (*domain.Order, error) {
	switch status {
	case domain.StatusPending, domain.StatusPaid, domain.StatusProcessing, domain.StatusShipped,
		domain.StatusDelivered, domain.StatusRefunded, domain.StatusCancelled:
	default:
		return nil, ErrBadStatus
	}
	order, err := s.repo.FindOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	order.Status = status
	order.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) BindPaymentEvents(bus *event.Bus) {
	if bus == nil {
		return
	}
	bus.Subscribe(domain.EventPaymentVerified, func(payload any) {
		m, ok := payload.(map[string]any)
		if !ok {
			return
		}
		orderID, _ := m["order_id"].(string)
		if orderID == "" {
			return
		}
		_, _ = s.MarkPaid(context.Background(), orderID)
	})
}
