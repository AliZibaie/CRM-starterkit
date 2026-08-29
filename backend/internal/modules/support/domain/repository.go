package domain

import "context"

type Repository interface {
	CreateDepartment(ctx context.Context, d *Department) error
	UpdateDepartment(ctx context.Context, d *Department) error
	DeleteDepartment(ctx context.Context, id string) error
	ListDepartments(ctx context.Context) ([]Department, error)
	FindDepartment(ctx context.Context, id string) (*Department, error)

	CreateTicket(ctx context.Context, t *Ticket) error
	UpdateTicket(ctx context.Context, t *Ticket) error
	FindTicket(ctx context.Context, id string) (*Ticket, error)
	ListTickets(ctx context.Context, status, departmentID, priority string) ([]Ticket, error)

	CreateMessage(ctx context.Context, m *Message) error
	ListMessages(ctx context.Context, ticketID string) ([]Message, error)
}
