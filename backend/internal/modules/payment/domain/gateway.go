package domain

import (
	"context"
	"time"
)

const EventPaymentVerified = "commerce.payment.verified"

type PayRequest struct {
	Amount      int64
	Currency    string
	Description string
	CallbackURL string
	OrderID     string
	Metadata    map[string]string
}

type PayResult struct {
	Authority string
	PaymentURL string
	ProviderRef string
}

type VerifyRequest struct {
	Amount    int64
	Authority string
	OrderID   string
}

type VerifyResult struct {
	Success     bool
	RefID       string
	ProviderRaw string
}

type RefundRequest struct {
	Amount      int64
	ProviderRef string
	Reason      string
}

type RefundResult struct {
	Success bool
	RefID   string
}

type Gateway interface {
	Name() string
	Pay(ctx context.Context, req PayRequest) (*PayResult, error)
	Verify(ctx context.Context, req VerifyRequest) (*VerifyResult, error)
	Refund(ctx context.Context, req RefundRequest) (*RefundResult, error)
}

type Payment struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"order_id"`
	Gateway     string    `json:"gateway"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	Authority   string    `json:"authority,omitempty"`
	ProviderRef string    `json:"provider_ref,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GatewayStatus struct {
	Name       string `json:"name"`
	Enabled    bool   `json:"enabled"`
	Configured bool   `json:"configured"`
}
