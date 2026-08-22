package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/banner/domain"

	"gorm.io/gorm"
)

type Model struct {
	ID        string     `gorm:"primaryKey;size:36"`
	Title     string     `gorm:"size:255;not null"`
	ImageURL  string     `gorm:"size:500;not null"`
	LinkURL   string     `gorm:"size:500"`
	Position  string     `gorm:"size:100;index;not null"`
	IsActive  bool       `gorm:"not null;default:true"`
	StartsAt  *time.Time
	EndsAt    *time.Time
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (Model) TableName() string { return "banners" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error { return s.db.AutoMigrate(&Model{}) }

func (s *Store) Create(ctx context.Context, b *domain.Banner) error {
	return s.db.WithContext(ctx).Create(toModel(b)).Error
}

func (s *Store) Update(ctx context.Context, b *domain.Banner) error {
	return s.db.WithContext(ctx).Save(toModel(b)).Error
}

func (s *Store) FindByID(ctx context.Context, id string) (*domain.Banner, error) {
	var m Model
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (s *Store) List(ctx context.Context) ([]domain.Banner, error) {
	var rows []Model
	if err := s.db.WithContext(ctx).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Banner, 0, len(rows))
	for i := range rows {
		out = append(out, *toDomain(&rows[i]))
	}
	return out, nil
}

func (s *Store) ListByPosition(ctx context.Context, position string) ([]domain.Banner, error) {
	var rows []Model
	if err := s.db.WithContext(ctx).Where("position = ?", position).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Banner, 0, len(rows))
	for i := range rows {
		out = append(out, *toDomain(&rows[i]))
	}
	return out, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(&Model{}, "id = ?", id).Error
}

func toModel(b *domain.Banner) *Model {
	return &Model{
		ID: b.ID, Title: b.Title, ImageURL: b.ImageURL, LinkURL: b.LinkURL, Position: b.Position,
		IsActive: b.IsActive, StartsAt: b.StartsAt, EndsAt: b.EndsAt, CreatedAt: b.CreatedAt, UpdatedAt: b.UpdatedAt,
	}
}

func toDomain(m *Model) *domain.Banner {
	return &domain.Banner{
		ID: m.ID, Title: m.Title, ImageURL: m.ImageURL, LinkURL: m.LinkURL, Position: m.Position,
		IsActive: m.IsActive, StartsAt: m.StartsAt, EndsAt: m.EndsAt, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}
