package domain

import "time"

const EventChange = "audit.change"

type Entry struct {
	ID         string    `json:"id"`
	ActorID    string    `json:"user_id"`
	Module     string    `json:"module"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	BeforeJSON string    `json:"before_json,omitempty"`
	AfterJSON  string    `json:"after_json,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type ChangePayload struct {
	ActorID    string `json:"actor_id"`
	Module     string `json:"module"`
	Action     string `json:"action"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Before     any    `json:"before"`
	After      any    `json:"after"`
}
