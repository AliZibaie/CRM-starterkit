package domain

import "time"

const (
	EventLeadCreated   = "crm.lead.created"
	EventTaskReminder  = "crm.task.reminder"
	EventDealStageMoved = "crm.deal.stage_moved"
)

type Lead struct {
	ID, Name, Email, Phone, Status, OwnerID, Source string
	CreatedAt, UpdatedAt                            time.Time
}

type Contact struct {
	ID, LeadID, UserID, Name, Email, Phone, OwnerID string
	CreatedAt, UpdatedAt                            time.Time
}

type Pipeline struct {
	ID, Name string
	StagesJSON string
	CreatedAt time.Time
}

type Deal struct {
	ID, Title, ContactID, PipelineID, Stage, OwnerID string
	Amount                                           int64
	CreatedAt, UpdatedAt                             time.Time
}

type Activity struct {
	ID, EntityType, EntityID, Type, Body, ActorID string
	CreatedAt                                     time.Time
}

type Task struct {
	ID, Title, EntityType, EntityID, OwnerID, Status string
	DueAt                                            *time.Time
	ReminderAt                                       *time.Time
	CreatedAt, UpdatedAt                             time.Time
}

type Tag struct {
	ID, Name string
}

type EntityTag struct {
	EntityType, EntityID, TagID string
}

type Segment struct {
	ID, Name, RulesJSON string
	CreatedAt           time.Time
}

type Customer360 struct {
	Contact    *Contact
	Orders     []map[string]any
	Tickets    []map[string]any
	Activities []Activity
	Deals      []Deal
	Tags       []string
}
