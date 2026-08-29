package domain

import "time"

type Media struct {
	ID           string    `json:"id"`
	OwnerID      string    `json:"owner_id,omitempty"`
	Key          string    `json:"key,omitempty"`
	Filename     string    `json:"filename,omitempty"`
	URL          string    `json:"url"`
	ThumbnailKey string    `json:"thumbnail_key,omitempty"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	ContentType  string    `json:"mime_type,omitempty"`
	Size         int64     `json:"size,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
