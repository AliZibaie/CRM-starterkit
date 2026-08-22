package inapp

import (
	"context"
	"encoding/json"

	"crm-starterkit/internal/modules/notification/domain"

	"github.com/google/uuid"
)

type Channel struct {
	repo domain.Repository
}

func New(repo domain.Repository) *Channel { return &Channel{repo: repo} }

func (c *Channel) Name() string { return "in_app" }

func (c *Channel) Send(ctx context.Context, payload domain.Payload) error {
	meta, _ := json.Marshal(payload.Data)
	n := &domain.Notification{
		ID: uuid.NewString(), UserID: payload.UserID, Channel: "in_app",
		Title: payload.Title, Body: payload.Body, Status: "unread", MetaJSON: string(meta),
	}
	return c.repo.Create(ctx, n)
}
