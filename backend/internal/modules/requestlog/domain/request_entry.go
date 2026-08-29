package domain

import "time"

type RequestEntry struct {
	ID          string    `json:"id"`
	RequestID   string    `json:"request_id,omitempty"`
	IP          string    `json:"ip"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	HeadersJSON string    `json:"headers_json,omitempty"`
	Payload     string    `json:"payload,omitempty"`
	DurationMs  int64     `json:"duration_ms"`
	Status      int       `json:"status"`
	UserID      string    `json:"user_id,omitempty"`
	UserAgent   string    `json:"user_agent,omitempty"`
	Country     string    `json:"country,omitempty"`
	City        string    `json:"city,omitempty"`
	Suspicious  bool      `json:"suspicious"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListFilter struct {
	IP     string
	Path   string
	UserID string
	Status *int
	From   *time.Time
	To     *time.Time
	Limit  int
	Offset int
}

type Stats struct {
	Total          int64
	Suspicious     int64
	AvgDurationMs  int64
	ErrorCount     int64
	ByStatus       map[string]int64
	ByCountry      map[string]int64
}

type StatsResponse struct {
	Total           int64 `json:"total"`
	AvgDurationMs   int64 `json:"avg_duration_ms"`
	ErrorCount      int64 `json:"error_count"`
	SuspiciousCount int64 `json:"suspicious_count"`
}
