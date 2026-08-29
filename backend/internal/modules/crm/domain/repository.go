package domain

import (
	"context"
	"time"
)

type Repository interface {
	CreateLead(ctx context.Context, l *Lead) error
	UpdateLead(ctx context.Context, l *Lead) error
	FindLead(ctx context.Context, id string) (*Lead, error)
	ListLeads(ctx context.Context) ([]Lead, error)

	CreateContact(ctx context.Context, c *Contact) error
	UpdateContact(ctx context.Context, c *Contact) error
	DeleteContact(ctx context.Context, id string) error
	FindContact(ctx context.Context, id string) (*Contact, error)
	FindContactByEmail(ctx context.Context, email string) (*Contact, error)
	ListContacts(ctx context.Context) ([]Contact, error)

	UpsertPipeline(ctx context.Context, p *Pipeline) error
	ListPipelines(ctx context.Context) ([]Pipeline, error)

	CreateDeal(ctx context.Context, d *Deal) error
	UpdateDeal(ctx context.Context, d *Deal) error
	FindDeal(ctx context.Context, id string) (*Deal, error)
	ListDealsByContact(ctx context.Context, contactID string) ([]Deal, error)
	ListDealsByPipeline(ctx context.Context, pipelineID string) ([]Deal, error)

	CreateActivity(ctx context.Context, a *Activity) error
	ListActivities(ctx context.Context, entityType, entityID string) ([]Activity, error)

	CreateTask(ctx context.Context, t *Task) error
	UpdateTask(ctx context.Context, t *Task) error
	ListTasksByOwner(ctx context.Context, ownerID string) ([]Task, error)
	ListDueReminders(ctx context.Context, before time.Time) ([]Task, error)

	UpsertTag(ctx context.Context, t *Tag) error
	AttachTag(ctx context.Context, entityType, entityID, tagID string) error
	ListTags(ctx context.Context, entityType, entityID string) ([]string, error)

	UpsertSegment(ctx context.Context, s *Segment) error
	ListSegments(ctx context.Context) ([]Segment, error)

	LoadOrdersByEmail(ctx context.Context, email string) ([]map[string]any, error)
	LoadTicketsByEmail(ctx context.Context, email string) ([]map[string]any, error)
}
