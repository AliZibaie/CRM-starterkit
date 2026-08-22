package domain

import "time"

type RequestEntry struct {
	ID           string
	RequestID    string
	IP           string
	Method       string
	Path         string
	HeadersJSON  string
	Payload      string
	DurationMs   int64
	Status       int
	UserID       string
	UserAgent    string
	Country      string
	City         string
	Suspicious   bool
	CreatedAt    time.Time
}

type ListFilter struct {
	IP       string
	Path     string
	UserID   string
	Status   *int
	From     *time.Time
	To       *time.Time
	Limit    int
	Offset   int
}

type Stats struct {
	Total       int64
	Suspicious  int64
	ByStatus    map[string]int64
	ByCountry   map[string]int64
}
