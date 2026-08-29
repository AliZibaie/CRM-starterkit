package domain

import "time"

const (
	AudienceAll  = "all"
	AudienceRole = "role"
	AudienceGroup = "group"
)

const EventBroadcast = "announcement.broadcast"

type Announcement struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	Audience    string     `json:"audience"`
	AudienceRef string     `json:"audience_ref"`
	Status      string     `json:"status"`
	CreatedBy   string     `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	BroadcastAt *time.Time `json:"broadcast_at,omitempty"`
}

type CreateInput struct {
	Title       string
	Body        string
	Audience    string
	AudienceRef string
	CreatedBy   string
}

type BroadcastPayload struct {
	AnnouncementID string `json:"announcement_id"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	Audience       string `json:"audience"`
	AudienceRef    string `json:"audience_ref"`
}
