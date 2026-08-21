package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/modules/auditlog/domain"

	"github.com/google/uuid"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Record(ctx context.Context, p domain.ChangePayload) error {
	before, _ := json.Marshal(p.Before)
	after, _ := json.Marshal(p.After)
	return s.repo.Create(ctx, &domain.Entry{
		ID: uuid.NewString(), ActorID: p.ActorID, Module: p.Module, Action: p.Action,
		EntityType: p.EntityType, EntityID: p.EntityID, BeforeJSON: string(before), AfterJSON: string(after),
		CreatedAt: time.Now().UTC(),
	})
}

func (s *Service) List(ctx context.Context, f domain.Filter) ([]domain.Entry, error) {
	return s.repo.List(ctx, f)
}

func (s *Service) BindEvents(bus *event.Bus) {
	if bus == nil {
		return
	}
	bus.Subscribe(domain.EventChange, func(payload any) {
		p := toPayload(payload)
		_ = s.Record(context.Background(), p)
	})
	mirror := func(module, action, entityType, idKey string) event.Handler {
		return func(payload any) {
			m, _ := payload.(map[string]any)
			entityID, _ := m[idKey].(string)
			actorID, _ := m["user_id"].(string)
			if actorID == "" {
				actorID, _ = m["owner_id"].(string)
			}
			if actorID == "" {
				actorID, _ = m["actor_id"].(string)
			}
			_ = s.Record(context.Background(), domain.ChangePayload{
				ActorID: actorID, Module: module, Action: action, EntityType: entityType, EntityID: entityID, After: m,
			})
		}
	}
	bus.Subscribe("commerce.order.created", mirror("order", "created", "order", "order_id"))
	bus.Subscribe("commerce.order.paid", mirror("order", "paid", "order", "order_id"))
	bus.Subscribe("commerce.order.refunded", mirror("order", "refunded", "order", "order_id"))
	bus.Subscribe("crm.lead.created", mirror("crm", "created", "lead", "lead_id"))
	bus.Subscribe("crm.deal.stage_moved", mirror("crm", "stage_moved", "deal", "deal_id"))
	bus.Subscribe("announcement.broadcast", func(payload any) {
		_ = s.Record(context.Background(), domain.ChangePayload{
			Module: "announcement", Action: "broadcast", EntityType: "announcement", After: fmt.Sprintf("%v", payload),
		})
	})
	bus.Subscribe("commerce.payment.verified", mirror("payment", "verified", "payment", "payment_id"))
}

func toPayload(payload any) domain.ChangePayload {
	switch p := payload.(type) {
	case domain.ChangePayload:
		return p
	case map[string]any:
		out := domain.ChangePayload{Before: p["before"], After: p["after"]}
		out.ActorID, _ = p["actor_id"].(string)
		out.Module, _ = p["module"].(string)
		out.Action, _ = p["action"].(string)
		out.EntityType, _ = p["entity_type"].(string)
		out.EntityID, _ = p["entity_id"].(string)
		return out
	default:
		return domain.ChangePayload{Module: "unknown", Action: "event", After: payload}
	}
}
