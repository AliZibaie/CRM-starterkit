package domain

import "time"

const EventChange = "audit.change"

type Entry struct {
	ID         string
	ActorID    string
	Module     string
	Action     string
	EntityType string
	EntityID   string
	BeforeJSON string
	AfterJSON  string
	CreatedAt  time.Time
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
