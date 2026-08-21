package domain

import "time"

const (
	AudienceAll  = "all"
	AudienceRole = "role"
	AudienceGroup = "group"
)

const EventBroadcast = "announcement.broadcast"

type Announcement struct {
	ID          string
	Title       string
	Body        string
	Audience    string
	AudienceRef string
	Status      string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	BroadcastAt *time.Time
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
