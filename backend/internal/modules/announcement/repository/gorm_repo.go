package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/announcement/domain"

	"gorm.io/gorm"
)

type Model struct {
	ID          string     `gorm:"primaryKey;size:36"`
	Title       string     `gorm:"size:255;not null"`
	Body        string     `gorm:"type:text;not null"`
	Audience    string     `gorm:"size:50;not null"`
	AudienceRef string     `gorm:"size:100"`
	Status      string     `gorm:"size:50;not null"`
	CreatedBy   string     `gorm:"size:36"`
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   time.Time  `gorm:"not null"`
	BroadcastAt *time.Time
}

func (Model) TableName() string { return "announcements" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error { return s.db.AutoMigrate(&Model{}) }

func (s *Store) Create(ctx context.Context, a *domain.Announcement) error {
	return s.db.WithContext(ctx).Create(toModel(a)).Error
}

func (s *Store) Update(ctx context.Context, a *domain.Announcement) error {
	return s.db.WithContext(ctx).Save(toModel(a)).Error
}

func (s *Store) FindByID(ctx context.Context, id string) (*domain.Announcement, error) {
	var m Model
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (s *Store) List(ctx context.Context) ([]domain.Announcement, error) {
	var rows []Model
	if err := s.db.WithContext(ctx).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Announcement, 0, len(rows))
	for i := range rows {
		out = append(out, *toDomain(&rows[i]))
	}
	return out, nil
}

func toModel(a *domain.Announcement) *Model {
	return &Model{
		ID: a.ID, Title: a.Title, Body: a.Body, Audience: a.Audience, AudienceRef: a.AudienceRef,
		Status: a.Status, CreatedBy: a.CreatedBy, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt, BroadcastAt: a.BroadcastAt,
	}
}

func toDomain(m *Model) *domain.Announcement {
	return &domain.Announcement{
		ID: m.ID, Title: m.Title, Body: m.Body, Audience: m.Audience, AudienceRef: m.AudienceRef,
		Status: m.Status, CreatedBy: m.CreatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt, BroadcastAt: m.BroadcastAt,
	}
}
