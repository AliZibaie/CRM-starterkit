package domain

import "time"

type Media struct {
	ID           string
	OwnerID      string
	Key          string
	URL          string
	ThumbnailKey string
	ThumbnailURL string
	ContentType  string
	Size         int64
	CreatedAt    time.Time
}
