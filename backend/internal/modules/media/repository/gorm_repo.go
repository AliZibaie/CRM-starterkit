package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/media/domain"

	"gorm.io/gorm"
)

type Model struct {
	ID           string    `gorm:"primaryKey;size:36"`
	OwnerID      string    `gorm:"size:36;index"`
	Key          string    `gorm:"column:storage_key;size:500;not null"`
	URL          string    `gorm:"size:500;not null"`
	ThumbnailKey string    `gorm:"size:500"`
	ThumbnailURL string    `gorm:"size:500"`
	ContentType  string    `gorm:"size:120"`
	Size         int64     `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
}

func (Model) TableName() string { return "media_files" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error { return s.db.AutoMigrate(&Model{}) }

func (s *Store) Create(ctx context.Context, m *domain.Media) error {
	return s.db.WithContext(ctx).Create(&Model{
		ID: m.ID, OwnerID: m.OwnerID, Key: m.Key, URL: m.URL, ThumbnailKey: m.ThumbnailKey,
		ThumbnailURL: m.ThumbnailURL, ContentType: m.ContentType, Size: m.Size, CreatedAt: m.CreatedAt,
	}).Error
}

func (s *Store) FindByID(ctx context.Context, id string) (*domain.Media, error) {
	var m Model
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (s *Store) ListByOwner(ctx context.Context, ownerID string) ([]domain.Media, error) {
	var rows []Model
	if err := s.db.WithContext(ctx).Where("owner_id = ?", ownerID).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Media, 0, len(rows))
	for i := range rows {
		out = append(out, *toDomain(&rows[i]))
	}
	return out, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(&Model{}, "id = ?", id).Error
}

func toDomain(m *Model) *domain.Media {
	return &domain.Media{
		ID: m.ID, OwnerID: m.OwnerID, Key: m.Key, URL: m.URL, ThumbnailKey: m.ThumbnailKey,
		ThumbnailURL: m.ThumbnailURL, ContentType: m.ContentType, Size: m.Size, CreatedAt: m.CreatedAt,
	}
}
