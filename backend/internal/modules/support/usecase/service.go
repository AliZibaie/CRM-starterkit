package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-starterkit/internal/modules/support/domain"
	"crm-starterkit/internal/modules/support/realtime"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrClosedTicket = errors.New("ticket is closed")
)

type Service struct {
	repo domain.Repository
	hub  *realtime.Hub
}

func NewService(repo domain.Repository, hub *realtime.Hub) *Service {
	return &Service{repo: repo, hub: hub}
}

func (s *Service) Hub() *realtime.Hub { return s.hub }

func (s *Service) CreateDepartment(ctx context.Context, name, description string) (*domain.Department, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidInput
	}
	now := time.Now().UTC()
	d := &domain.Department{ID: uuid.NewString(), Name: name, Description: strings.TrimSpace(description), CreatedAt: now, UpdatedAt: now}
	if err := s.repo.CreateDepartment(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) ListDepartments(ctx context.Context) ([]domain.Department, error) {
	return s.repo.ListDepartments(ctx)
}

func (s *Service) UpdateDepartment(ctx context.Context, id, name, description string) (*domain.Department, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidInput
	}
	d, err := s.repo.FindDepartment(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	d.Name = name
	d.Description = strings.TrimSpace(description)
	d.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateDepartment(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) DeleteDepartment(ctx context.Context, id string) error {
	if _, err := s.repo.FindDepartment(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.DeleteDepartment(ctx, id)
}

func (s *Service) CreateTicket(ctx context.Context, subject, priority, departmentID, requesterID string) (*domain.Ticket, error) {
	subject = strings.TrimSpace(subject)
	if subject == "" || requesterID == "" {
		return nil, ErrInvalidInput
	}
	if priority == "" {
		priority = "normal"
	}
	now := time.Now().UTC()
	t := &domain.Ticket{
		ID: uuid.NewString(), Subject: subject, Status: domain.StatusOpen, Priority: priority,
		DepartmentID: departmentID, RequesterID: requesterID, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.CreateTicket(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) Reply(ctx context.Context, ticketID, senderID, body string, internal bool) (*domain.Message, error) {
	body = strings.TrimSpace(body)
	if body == "" || senderID == "" {
		return nil, ErrInvalidInput
	}
	t, err := s.repo.FindTicket(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if t.Status == domain.StatusClosed {
		return nil, ErrClosedTicket
	}
	msg := &domain.Message{
		ID: uuid.NewString(), TicketID: ticketID, SenderID: senderID, Body: body,
		IsInternal: internal, CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.CreateMessage(ctx, msg); err != nil {
		return nil, err
	}
	if t.Status == domain.StatusOpen {
		t.Status = domain.StatusInProgress
		t.UpdatedAt = time.Now().UTC()
		_ = s.repo.UpdateTicket(ctx, t)
	}
	if s.hub != nil {
		s.hub.Broadcast(ticketID, map[string]any{"type": "message", "message": msg})
	}
	return msg, nil
}

func (s *Service) AssignAgent(ctx context.Context, ticketID, agentID string) (*domain.Ticket, error) {
	if agentID == "" {
		return nil, ErrInvalidInput
	}
	t, err := s.repo.FindTicket(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	t.AssignedTo = agentID
	if t.Status == domain.StatusOpen {
		t.Status = domain.StatusInProgress
	}
	t.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateTicket(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ChangeStatus(ctx context.Context, ticketID, status string) (*domain.Ticket, error) {
	switch status {
	case domain.StatusOpen, domain.StatusInProgress, domain.StatusResolved, domain.StatusClosed:
	default:
		return nil, ErrInvalidInput
	}
	t, err := s.repo.FindTicket(ctx, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	t.Status = status
	t.UpdatedAt = time.Now().UTC()
	if status == domain.StatusClosed {
		now := time.Now().UTC()
		t.ClosedAt = &now
	}
	if err := s.repo.UpdateTicket(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) CloseTicket(ctx context.Context, ticketID string) (*domain.Ticket, error) {
	return s.ChangeStatus(ctx, ticketID, domain.StatusClosed)
}

func (s *Service) ListTickets(ctx context.Context, status, departmentID, priority string) ([]domain.Ticket, error) {
	return s.repo.ListTickets(ctx, status, departmentID, priority)
}

func (s *Service) GetTicket(ctx context.Context, id string) (*domain.Ticket, []domain.Message, error) {
	t, err := s.repo.FindTicket(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrNotFound
		}
		return nil, nil, err
	}
	msgs, err := s.repo.ListMessages(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return t, msgs, nil
}
