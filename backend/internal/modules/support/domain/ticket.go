package domain

import "time"

const (
	StatusOpen       = "open"
	StatusInProgress = "in_progress"
	StatusResolved   = "resolved"
	StatusClosed     = "closed"
)

type Department struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Ticket struct {
	ID           string     `json:"id"`
	Subject      string     `json:"subject"`
	Status       string     `json:"status"`
	Priority     string     `json:"priority"`
	DepartmentID string     `json:"department_id"`
	RequesterID  string     `json:"requester_id"`
	AssignedTo   string     `json:"assigned_to,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
}

type Message struct {
	ID         string    `json:"id"`
	TicketID   string    `json:"ticket_id"`
	SenderID   string    `json:"sender_id"`
	Body       string    `json:"body"`
	IsInternal bool      `json:"is_internal"`
	CreatedAt  time.Time `json:"created_at"`
}
