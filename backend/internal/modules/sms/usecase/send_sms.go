package usecase

import (
	"context"
	"errors"
	"strings"

	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/modules/sms/domain"
)

var ErrInvalidInput = errors.New("invalid sms input")

type Service struct {
	provider domain.Provider
}

func NewService(provider domain.Provider) *Service {
	return &Service{provider: provider}
}

func (s *Service) Send(ctx context.Context, to, text string) error {
	to = strings.TrimSpace(to)
	text = strings.TrimSpace(text)
	if to == "" || text == "" {
		return ErrInvalidInput
	}
	return s.provider.Send(ctx, to, text)
}

func (s *Service) BindEvents(bus *event.Bus) {
	if bus == nil || s.provider == nil {
		return
	}
	bus.Subscribe(domain.EventSend, func(payload any) {
		req, ok := payload.(domain.SendRequest)
		if !ok {
			if m, ok := payload.(map[string]any); ok {
				req.To, _ = m["to"].(string)
				req.Text, _ = m["text"].(string)
			}
		}
		if req.To == "" || req.Text == "" {
			return
		}
		_ = s.provider.Send(context.Background(), req.To, req.Text)
	})
}
