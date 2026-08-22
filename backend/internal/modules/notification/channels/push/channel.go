package push

import (
	"context"

	"crm-starterkit/internal/modules/notification/domain"
)

type Channel struct{}

func New() *Channel { return &Channel{} }

func (c *Channel) Name() string { return "push" }

func (c *Channel) Send(ctx context.Context, payload domain.Payload) error {
	_ = ctx
	_ = payload
	return nil
}
