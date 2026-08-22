package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Kavenegar struct {
	apiKey string
	sender string
	client *http.Client
}

func NewKavenegar(apiKey, sender string) *Kavenegar {
	return &Kavenegar{
		apiKey: apiKey,
		sender: sender,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (k *Kavenegar) Name() string { return "kavenegar" }

func (k *Kavenegar) Send(ctx context.Context, to, text string) error {
	if k.apiKey == "" {
		return fmt.Errorf("kavenegar api key missing")
	}
	endpoint := fmt.Sprintf("https://api.kavenegar.com/v1/%s/sms/send.json", k.apiKey)
	q := url.Values{}
	q.Set("receptor", to)
	q.Set("message", text)
	if k.sender != "" {
		q.Set("sender", k.sender)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+q.Encode(), nil)
	if err != nil {
		return err
	}
	resp, err := k.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("kavenegar failed: %s", string(body))
	}
	return nil
}
