package stripegw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"crm-starterkit/internal/modules/payment/domain"
)

type Gateway struct {
	secret   string
	currency string
	client   *http.Client
}

func New(secret, currency string) *Gateway {
	if currency == "" {
		currency = "usd"
	}
	return &Gateway{
		secret:   secret,
		currency: currency,
		client:   &http.Client{Timeout: 20 * time.Second},
	}
}

func (g *Gateway) Name() string { return "stripe" }

func (g *Gateway) Pay(ctx context.Context, req domain.PayRequest) (*domain.PayResult, error) {
	if g.secret == "" {
		return nil, fmt.Errorf("stripe secret key missing")
	}
	form := url.Values{}
	form.Set("mode", "payment")
	form.Set("success_url", req.CallbackURL+"?session_id={CHECKOUT_SESSION_ID}&order_id="+req.OrderID+"&gateway=stripe")
	form.Set("cancel_url", req.CallbackURL+"?canceled=1&order_id="+req.OrderID+"&gateway=stripe")
	form.Set("line_items[0][price_data][currency]", strings.ToLower(req.Currency))
	if form.Get("line_items[0][price_data][currency]") == "irr" {
		form.Set("line_items[0][price_data][currency]", g.currency)
	}
	form.Set("line_items[0][price_data][product_data][name]", req.Description)
	form.Set("line_items[0][price_data][unit_amount]", fmt.Sprintf("%d", req.Amount))
	form.Set("line_items[0][quantity]", "1")
	form.Set("client_reference_id", req.OrderID)
	form.Set("metadata[order_id]", req.OrderID)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.stripe.com/v1/checkout/sessions", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	httpReq.SetBasicAuth(g.secret, "")
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var out struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	if out.ID == "" || out.URL == "" {
		return nil, fmt.Errorf("stripe session failed: %s", string(data))
	}
	return &domain.PayResult{Authority: out.ID, PaymentURL: out.URL, ProviderRef: out.ID}, nil
}

func (g *Gateway) Verify(ctx context.Context, req domain.VerifyRequest) (*domain.VerifyResult, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.stripe.com/v1/checkout/sessions/"+req.Authority, nil)
	if err != nil {
		return nil, err
	}
	httpReq.SetBasicAuth(g.secret, "")
	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var out struct {
		ID            string `json:"id"`
		PaymentStatus string `json:"payment_status"`
		PaymentIntent string `json:"payment_intent"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	ok := out.PaymentStatus == "paid"
	ref := out.PaymentIntent
	if ref == "" {
		ref = out.ID
	}
	return &domain.VerifyResult{Success: ok, RefID: ref, ProviderRaw: string(data)}, nil
}

func (g *Gateway) Refund(ctx context.Context, req domain.RefundRequest) (*domain.RefundResult, error) {
	form := url.Values{}
	form.Set("payment_intent", req.ProviderRef)
	if req.Amount > 0 {
		form.Set("amount", fmt.Sprintf("%d", req.Amount))
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.stripe.com/v1/refunds", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	httpReq.SetBasicAuth(g.secret, "")
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var out struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return &domain.RefundResult{Success: out.Status == "succeeded" || out.Status == "pending", RefID: out.ID}, nil
}
