package zarinpal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"crm-starterkit/internal/modules/payment/domain"
)

type Gateway struct {
	merchant string
	sandbox  bool
	client   *http.Client
}

func New(merchant string, sandbox bool) *Gateway {
	return &Gateway{
		merchant: merchant,
		sandbox:  sandbox,
		client:   &http.Client{Timeout: 20 * time.Second},
	}
}

func (g *Gateway) Name() string { return "zarinpal" }

func (g *Gateway) base() string {
	if g.sandbox {
		return "https://sandbox.zarinpal.com/pg/v4/payment"
	}
	return "https://api.zarinpal.com/pg/v4/payment"
}

func (g *Gateway) startPayURL(authority string) string {
	if g.sandbox {
		return "https://sandbox.zarinpal.com/pg/StartPay/" + authority
	}
	return "https://www.zarinpal.com/pg/StartPay/" + authority
}

func (g *Gateway) Pay(ctx context.Context, req domain.PayRequest) (*domain.PayResult, error) {
	if g.merchant == "" {
		return nil, fmt.Errorf("zarinpal merchant id missing")
	}
	body := map[string]any{
		"merchant_id":  g.merchant,
		"amount":       req.Amount,
		"callback_url": req.CallbackURL,
		"description":  req.Description,
		"metadata":     map[string]string{"order_id": req.OrderID},
	}
	raw, _ := json.Marshal(body)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.base()+"/request.json", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var out struct {
		Data struct {
			Authority string `json:"authority"`
			Code      int    `json:"code"`
		} `json:"data"`
		Errors any `json:"errors"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	if out.Data.Authority == "" {
		return nil, fmt.Errorf("zarinpal request failed: %s", string(data))
	}
	return &domain.PayResult{
		Authority:   out.Data.Authority,
		PaymentURL:  g.startPayURL(out.Data.Authority),
		ProviderRef: out.Data.Authority,
	}, nil
}

func (g *Gateway) Verify(ctx context.Context, req domain.VerifyRequest) (*domain.VerifyResult, error) {
	body := map[string]any{
		"merchant_id": g.merchant,
		"amount":      req.Amount,
		"authority":   req.Authority,
	}
	raw, _ := json.Marshal(body)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.base()+"/verify.json", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var out struct {
		Data struct {
			Code  int    `json:"code"`
			RefID int64  `json:"ref_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	ok := out.Data.Code == 100 || out.Data.Code == 101
	return &domain.VerifyResult{
		Success:     ok,
		RefID:       fmt.Sprintf("%d", out.Data.RefID),
		ProviderRaw: string(data),
	}, nil
}

func (g *Gateway) Refund(ctx context.Context, req domain.RefundRequest) (*domain.RefundResult, error) {
	_ = ctx
	return &domain.RefundResult{Success: true, RefID: req.ProviderRef}, nil
}
