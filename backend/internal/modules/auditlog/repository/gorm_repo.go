package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/auditlog/domain"

	"gorm.io/gorm"
)

type Model struct {
	ID         string    `gorm:"primaryKey;size:36"`
	ActorID    string    `gorm:"size:36;index"`
	Module     string    `gorm:"size:100;index"`
	Action     string    `gorm:"size:100"`
	EntityType string    `gorm:"size:100;index"`
	EntityID   string    `gorm:"size:100;index"`
	BeforeJSON string    `gorm:"type:text"`
	AfterJSON  string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"index;not null"`
}

func (Model) TableName() string { return "audit_logs" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error { return s.db.AutoMigrate(&Model{}) }

func (s *Store) Create(ctx context.Context, e *domain.Entry) error {
	return s.db.WithContext(ctx).Create(&Model{
		ID: e.ID, ActorID: e.ActorID, Module: e.Module, Action: e.Action,
		EntityType: e.EntityType, EntityID: e.EntityID, BeforeJSON: e.BeforeJSON, AfterJSON: e.AfterJSON, CreatedAt: e.CreatedAt,
	}).Error
}

func (s *Store) List(ctx context.Context, f domain.Filter) ([]domain.Entry, error) {
	q := s.db.WithContext(ctx).Model(&Model{})
	if f.ActorID != "" {
		q = q.Where("actor_id = ?", f.ActorID)
	}
	if f.Module != "" {
		q = q.Where("module = ?", f.Module)
	}
	if f.EntityType != "" {
		q = q.Where("entity_type = ?", f.EntityType)
	}
	if f.EntityID != "" {
		q = q.Where("entity_id = ?", f.EntityID)
	}
	if f.From != nil {
		q = q.Where("created_at >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("created_at <= ?", *f.To)
	}
	limit := f.Limit
	if limit <= 0 {
		limit = 100
	}
	var rows []Model
	if err := q.Order("created_at desc").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Entry, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.Entry{
			ID: r.ID, ActorID: r.ActorID, Module: r.Module, Action: r.Action,
			EntityType: r.EntityType, EntityID: r.EntityID, BeforeJSON: r.BeforeJSON, AfterJSON: r.AfterJSON, CreatedAt: r.CreatedAt,
		})
	}
	return out, nil
}
