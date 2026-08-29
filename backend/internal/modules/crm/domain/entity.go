package domain

import "time"

const (
	EventLeadCreated    = "crm.lead.created"
	EventTaskReminder   = "crm.task.reminder"
	EventDealStageMoved = "crm.deal.stage_moved"
)

type Lead struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	OwnerID   string    `json:"owner_id,omitempty"`
	Source    string    `json:"source,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Contact struct {
	ID        string    `json:"id"`
	LeadID    string    `json:"lead_id,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Company   string    `json:"company,omitempty"`
	OwnerID   string    `json:"owner_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Pipeline struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	StagesJSON string    `json:"stages_json"`
	CreatedAt  time.Time `json:"created_at"`
}

type Deal struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	ContactID  string    `json:"contact_id,omitempty"`
	PipelineID string    `json:"pipeline_id"`
	Stage      string    `json:"stage"`
	OwnerID    string    `json:"owner_id,omitempty"`
	Amount     int64     `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Activity struct {
	ID         string    `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Type       string    `json:"type"`
	Body       string    `json:"body"`
	ActorID    string    `json:"actor_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type Task struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	EntityType string     `json:"entity_type,omitempty"`
	EntityID   string     `json:"entity_id,omitempty"`
	OwnerID    string     `json:"owner_id,omitempty"`
	Status     string     `json:"status"`
	DueAt      *time.Time `json:"due_at,omitempty"`
	ReminderAt *time.Time `json:"reminder_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EntityTag struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	TagID      string `json:"tag_id"`
}

type Segment struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	RulesJSON string    `json:"rules_json"`
	CreatedAt time.Time `json:"created_at"`
}

type Customer360 struct {
	Contact    *Contact         `json:"contact"`
	Orders     []map[string]any `json:"orders"`
	Tickets    []map[string]any `json:"tickets"`
	Activities []Activity       `json:"activities"`
	Deals      []Deal           `json:"deals"`
	Tags       []string         `json:"tags"`
}
