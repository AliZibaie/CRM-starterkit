package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/support/domain"

	"gorm.io/gorm"
)

type DepartmentModel struct {
	ID          string    `gorm:"primaryKey;size:36"`
	Name        string    `gorm:"uniqueIndex;size:100;not null"`
	Description string    `gorm:"size:255"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (DepartmentModel) TableName() string { return "support_departments" }

type TicketModel struct {
	ID           string     `gorm:"primaryKey;size:36"`
	Subject      string     `gorm:"size:255;not null"`
	Status       string     `gorm:"size:50;index;not null"`
	Priority     string     `gorm:"size:50;not null"`
	DepartmentID string     `gorm:"size:36;index"`
	RequesterID  string     `gorm:"size:36;index"`
	AssignedTo   string     `gorm:"size:36;index"`
	CreatedAt    time.Time  `gorm:"not null"`
	UpdatedAt    time.Time  `gorm:"not null"`
	ClosedAt     *time.Time
}

func (TicketModel) TableName() string { return "tickets" }

type MessageModel struct {
	ID         string    `gorm:"primaryKey;size:36"`
	TicketID   string    `gorm:"size:36;index;not null"`
	SenderID   string    `gorm:"size:36;not null"`
	Body       string    `gorm:"type:text;not null"`
	IsInternal bool      `gorm:"not null;default:false"`
	CreatedAt  time.Time `gorm:"not null"`
}

func (MessageModel) TableName() string { return "ticket_messages" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error {
	return s.db.AutoMigrate(&DepartmentModel{}, &TicketModel{}, &MessageModel{})
}

func (s *Store) CreateDepartment(ctx context.Context, d *domain.Department) error {
	return s.db.WithContext(ctx).Create(&DepartmentModel{
		ID: d.ID, Name: d.Name, Description: d.Description, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}).Error
}

func (s *Store) ListDepartments(ctx context.Context) ([]domain.Department, error) {
	var rows []DepartmentModel
	if err := s.db.WithContext(ctx).Order("name asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Department, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.Department{ID: r.ID, Name: r.Name, Description: r.Description, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt})
	}
	return out, nil
}

func (s *Store) FindDepartment(ctx context.Context, id string) (*domain.Department, error) {
	var r DepartmentModel
	if err := s.db.WithContext(ctx).First(&r, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &domain.Department{ID: r.ID, Name: r.Name, Description: r.Description, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}, nil
}

func (s *Store) CreateTicket(ctx context.Context, t *domain.Ticket) error {
	return s.db.WithContext(ctx).Create(ticketModel(t)).Error
}

func (s *Store) UpdateTicket(ctx context.Context, t *domain.Ticket) error {
	return s.db.WithContext(ctx).Save(ticketModel(t)).Error
}

func (s *Store) FindTicket(ctx context.Context, id string) (*domain.Ticket, error) {
	var m TicketModel
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return ticketDomain(&m), nil
}

func (s *Store) ListTickets(ctx context.Context) ([]domain.Ticket, error) {
	var rows []TicketModel
	if err := s.db.WithContext(ctx).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Ticket, 0, len(rows))
	for i := range rows {
		out = append(out, *ticketDomain(&rows[i]))
	}
	return out, nil
}

func (s *Store) CreateMessage(ctx context.Context, m *domain.Message) error {
	return s.db.WithContext(ctx).Create(&MessageModel{
		ID: m.ID, TicketID: m.TicketID, SenderID: m.SenderID, Body: m.Body, IsInternal: m.IsInternal, CreatedAt: m.CreatedAt,
	}).Error
}

func (s *Store) ListMessages(ctx context.Context, ticketID string) ([]domain.Message, error) {
	var rows []MessageModel
	if err := s.db.WithContext(ctx).Where("ticket_id = ?", ticketID).Order("created_at asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Message, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.Message{ID: r.ID, TicketID: r.TicketID, SenderID: r.SenderID, Body: r.Body, IsInternal: r.IsInternal, CreatedAt: r.CreatedAt})
	}
	return out, nil
}

func ticketModel(t *domain.Ticket) *TicketModel {
	return &TicketModel{
		ID: t.ID, Subject: t.Subject, Status: t.Status, Priority: t.Priority, DepartmentID: t.DepartmentID,
		RequesterID: t.RequesterID, AssignedTo: t.AssignedTo, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt, ClosedAt: t.ClosedAt,
	}
}

func ticketDomain(m *TicketModel) *domain.Ticket {
	return &domain.Ticket{
		ID: m.ID, Subject: m.Subject, Status: m.Status, Priority: m.Priority, DepartmentID: m.DepartmentID,
		RequesterID: m.RequesterID, AssignedTo: m.AssignedTo, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt, ClosedAt: m.ClosedAt,
	}
}
