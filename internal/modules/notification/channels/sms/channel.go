package smschan

import (
	"context"

	"crm-starterkit/internal/core/event"
	"crm-starterkit/internal/modules/notification/domain"
)

const eventSend = "sms.send"

type Channel struct {
	bus *event.Bus
}

func New(bus *event.Bus) *Channel { return &Channel{bus: bus} }

func (c *Channel) Name() string { return "sms" }

func (c *Channel) Send(ctx context.Context, payload domain.Payload) error {
	_ = ctx
	if c.bus == nil || payload.To == "" {
		return nil
	}
	c.bus.Publish(eventSend, map[string]any{"to": payload.To, "text": payload.Body})
	return nil
}
