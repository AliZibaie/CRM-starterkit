package usecase

import (
	"context"
	"fmt"

	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/modules/notification/domain"
)

type Service struct {
	repo     domain.Repository
	channels map[string]domain.Channel
}

func NewService(repo domain.Repository, channels []domain.Channel) *Service {
	m := make(map[string]domain.Channel, len(channels))
	for _, ch := range channels {
		m[ch.Name()] = ch
	}
	return &Service{repo: repo, channels: m}
}

func (s *Service) Dispatch(ctx context.Context, payload domain.Payload) error {
	names := []string{payload.Channel}
	if payload.Channel == "" || payload.Channel == "all" {
		names = []string{"in_app", "email", "sms", "push"}
	}
	var firstErr error
	for _, name := range names {
		ch, ok := s.channels[name]
		if !ok {
			continue
		}
		p := payload
		p.Channel = name
		if err := ch.Send(ctx, p); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (s *Service) ListMine(ctx context.Context, userID string) ([]domain.Notification, error) {
	return s.repo.ListByUser(ctx, userID, 100)
}

func (s *Service) MarkRead(ctx context.Context, id, userID string) error {
	return s.repo.MarkRead(ctx, id, userID)
}

func (s *Service) BindEvents(bus *event.Bus) {
	if bus == nil {
		return
	}
	bus.Subscribe("announcement.broadcast", func(payload any) {
		title, body := "Announcement", ""
		data := map[string]any{"source": "announcement"}
		switch p := payload.(type) {
		case map[string]any:
			if t, ok := p["title"].(string); ok {
				title = t
			}
			if b, ok := p["body"].(string); ok {
				body = b
			}
			data = p
		default:
			if bp, ok := toBroadcastMap(payload); ok {
				title = bp["title"]
				body = bp["body"]
			} else {
				body = fmt.Sprintf("%v", payload)
			}
		}
		_ = s.Dispatch(context.Background(), domain.Payload{
			Title: title, Body: body, Channel: "in_app", Data: data,
		})
	})
	bus.Subscribe("commerce.order.created", func(payload any) {
		m, _ := payload.(map[string]any)
		userID, _ := m["user_id"].(string)
		orderID, _ := m["order_id"].(string)
		_ = s.Dispatch(context.Background(), domain.Payload{
			UserID: userID, Title: "Order created", Body: "Order " + orderID + " created",
			Channel: "in_app", Data: m,
		})
	})
	bus.Subscribe("commerce.order.paid", func(payload any) {
		m, _ := payload.(map[string]any)
		userID, _ := m["user_id"].(string)
		orderID, _ := m["order_id"].(string)
		_ = s.Dispatch(context.Background(), domain.Payload{
			UserID: userID, Title: "Payment received", Body: "Order " + orderID + " paid",
			Channel: "all", Data: m,
		})
	})
	bus.Subscribe("crm.lead.created", func(payload any) {
		m, _ := payload.(map[string]any)
		ownerID, _ := m["owner_id"].(string)
		_ = s.Dispatch(context.Background(), domain.Payload{
			UserID: ownerID, Title: "New lead", Body: "A new lead was assigned to you",
			Channel: "in_app", Data: m,
		})
	})
	bus.Subscribe("crm.task.reminder", func(payload any) {
		m, _ := payload.(map[string]any)
		ownerID, _ := m["owner_id"].(string)
		title, _ := m["title"].(string)
		_ = s.Dispatch(context.Background(), domain.Payload{
			UserID: ownerID, Title: "Task reminder", Body: title, Channel: "all", Data: m,
		})
	})
}
