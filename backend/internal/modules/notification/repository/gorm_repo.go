package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/notification/domain"

	"gorm.io/gorm"
)

type Model struct {
	ID        string    `gorm:"primaryKey;size:36"`
	UserID    string    `gorm:"size:36;index;not null"`
	Channel   string    `gorm:"size:50;not null"`
	Title     string    `gorm:"size:255;not null"`
	Body      string    `gorm:"type:text;not null"`
	Status    string    `gorm:"size:50;not null"`
	MetaJSON  string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"not null"`
}

func (Model) TableName() string { return "notifications" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error { return s.db.AutoMigrate(&Model{}) }

func (s *Store) Create(ctx context.Context, n *domain.Notification) error {
	return s.db.WithContext(ctx).Create(&Model{
		ID: n.ID, UserID: n.UserID, Channel: n.Channel, Title: n.Title, Body: n.Body,
		Status: n.Status, MetaJSON: n.MetaJSON, CreatedAt: time.Now().UTC(),
	}).Error
}

func (s *Store) ListByUser(ctx context.Context, userID string, limit int) ([]domain.Notification, error) {
	if limit <= 0 {
		limit = 50
	}
	var rows []Model
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Notification, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.Notification{
			ID: r.ID, UserID: r.UserID, Channel: r.Channel, Title: r.Title, Body: r.Body, Status: r.Status, MetaJSON: r.MetaJSON,
		})
	}
	return out, nil
}

func (s *Store) MarkRead(ctx context.Context, id, userID string) error {
	return s.db.WithContext(ctx).Model(&Model{}).Where("id = ? AND user_id = ?", id, userID).Update("status", "read").Error
}
