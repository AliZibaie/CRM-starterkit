package domain

import "time"

const (
	StatusOpen       = "open"
	StatusInProgress = "in_progress"
	StatusResolved   = "resolved"
	StatusClosed     = "closed"
)

type Department struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Ticket struct {
	ID           string
	Subject      string
	Status       string
	Priority     string
	DepartmentID string
	RequesterID  string
	AssignedTo   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ClosedAt     *time.Time
}

type Message struct {
	ID         string
	TicketID   string
	SenderID   string
	Body       string
	IsInternal bool
	CreatedAt  time.Time
}
