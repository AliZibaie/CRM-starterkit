package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"crm-starterkit/config"
	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/modules/payment/domain"
	stripegw "crm-starterkit/internal/modules/payment/gateways/stripe"
	"crm-starterkit/internal/modules/payment/gateways/zarinpal"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusInitiated = "initiated"
	StatusPaid      = "paid"
	StatusFailed    = "failed"
	StatusRefunded  = "refunded"
)

var (
	ErrNotFound      = errors.New("payment not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrGatewayMissing = errors.New("gateway not configured")
	ErrBadOrder      = errors.New("order not payable")
)

type Service struct {
	repo     domain.Repository
	gateways map[string]domain.Gateway
	bus      *event.Bus
	callback string
	cfg      config.PaymentConfig
}

func NewService(repo domain.Repository, cfg config.PaymentConfig, bus *event.Bus) *Service {
	gws := map[string]domain.Gateway{}
	if cfg.ZarinpalMerchant != "" {
		gws["zarinpal"] = zarinpal.New(cfg.ZarinpalMerchant, cfg.ZarinpalSandbox)
	}
	if cfg.StripeSecretKey != "" {
		gws["stripe"] = stripegw.New(cfg.StripeSecretKey, cfg.StripeCurrency)
	}
	return &Service{
		repo:     repo,
		gateways: gws,
		bus:      bus,
		callback: strings.TrimRight(cfg.CallbackBaseURL, "/"),
		cfg:      cfg,
	}
}

func (s *Service) ListGatewayStatus() []domain.GatewayStatus {
	defs := []struct {
		name       string
		configured bool
	}{
		{"zarinpal", s.cfg.ZarinpalMerchant != ""},
		{"stripe", s.cfg.StripeSecretKey != ""},
	}
	out := make([]domain.GatewayStatus, 0, len(defs))
	for _, d := range defs {
		_, enabled := s.gateways[d.name]
		out = append(out, domain.GatewayStatus{
			Name:       d.name,
			Configured: d.configured,
			Enabled:    d.configured && enabled,
		})
	}
	return out
}

func (s *Service) ListFailedPayments(ctx context.Context, limit int) ([]domain.Payment, error) {
	return s.repo.ListByStatus(ctx, StatusFailed, limit)
}

func (s *Service) InitPayment(ctx context.Context, orderID, gatewayName, description string) (*domain.Payment, *domain.PayResult, error) {
	gatewayName = strings.ToLower(strings.TrimSpace(gatewayName))
	gw, ok := s.gateways[gatewayName]
	if !ok {
		return nil, nil, ErrGatewayMissing
	}
	amount, currency, status, err := s.repo.FindOrderAmount(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrBadOrder
		}
		return nil, nil, err
	}
	if status != "pending" || amount <= 0 {
		return nil, nil, ErrBadOrder
	}
	if description == "" {
		description = "order " + orderID
	}
	callback := fmt.Sprintf("%s/%s", s.callback, gatewayName)
	result, err := gw.Pay(ctx, domain.PayRequest{
		Amount: amount, Currency: currency, Description: description,
		CallbackURL: callback, OrderID: orderID,
	})
	if err != nil {
		return nil, nil, err
	}
	now := time.Now().UTC()
	p := &domain.Payment{
		ID: uuid.NewString(), OrderID: orderID, Gateway: gatewayName, Amount: amount, Currency: currency,
		Status: StatusInitiated, Authority: result.Authority, ProviderRef: result.ProviderRef,
		CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, nil, err
	}
	return p, result, nil
}

func (s *Service) VerifyPayment(ctx context.Context, gatewayName, authority, orderID string) (*domain.Payment, *domain.VerifyResult, error) {
	gatewayName = strings.ToLower(strings.TrimSpace(gatewayName))
	gw, ok := s.gateways[gatewayName]
	if !ok {
		return nil, nil, ErrGatewayMissing
	}
	var p *domain.Payment
	var err error
	if authority != "" {
		p, err = s.repo.FindByAuthority(ctx, authority)
	} else if orderID != "" {
		p, err = s.repo.FindByOrder(ctx, orderID)
	} else {
		return nil, nil, ErrInvalidInput
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrNotFound
		}
		return nil, nil, err
	}
	if authority == "" {
		authority = p.Authority
	}
	result, err := gw.Verify(ctx, domain.VerifyRequest{
		Amount: p.Amount, Authority: authority, OrderID: p.OrderID,
	})
	if err != nil {
		return nil, nil, err
	}
	p.UpdatedAt = time.Now().UTC()
	if result.Success {
		p.Status = StatusPaid
		p.ProviderRef = result.RefID
		if err := s.repo.Update(ctx, p); err != nil {
			return nil, nil, err
		}
		if s.bus != nil {
			s.bus.Publish(domain.EventPaymentVerified, map[string]any{
				"order_id": p.OrderID, "payment_id": p.ID, "gateway": p.Gateway, "ref_id": result.RefID,
			})
		}
	} else {
		p.Status = StatusFailed
		_ = s.repo.Update(ctx, p)
	}
	return p, result, nil
}

func (s *Service) HandleCallback(ctx context.Context, gatewayName string, authority, orderID, status string) (*domain.Payment, *domain.VerifyResult, error) {
	if strings.EqualFold(status, "NOK") || status == "canceled" || status == "1" {
		return nil, &domain.VerifyResult{Success: false}, ErrInvalidInput
	}
	if authority == "" {
		authority = orderID
	}
	return s.VerifyPayment(ctx, gatewayName, authority, orderID)
}

func (s *Service) Refund(ctx context.Context, orderID string) (*domain.Payment, error) {
	p, err := s.repo.FindByOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	gw, ok := s.gateways[p.Gateway]
	if !ok {
		return nil, ErrGatewayMissing
	}
	res, err := gw.Refund(ctx, domain.RefundRequest{Amount: p.Amount, ProviderRef: p.ProviderRef})
	if err != nil {
		return nil, err
	}
	if !res.Success {
		return nil, ErrInvalidInput
	}
	p.Status = StatusRefunded
	p.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
