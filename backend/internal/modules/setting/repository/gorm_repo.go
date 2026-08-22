package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/setting/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Model struct {
	Key       string    `gorm:"column:setting_key;primaryKey;size:191"`
	Value     string    `gorm:"type:text"`
	UpdatedBy string    `gorm:"size:36"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (Model) TableName() string { return "system_settings" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error { return s.db.AutoMigrate(&Model{}) }

func (s *Store) Get(ctx context.Context, key string) (*domain.Setting, error) {
	var m Model
	if err := s.db.WithContext(ctx).First(&m, "setting_key = ?", key).Error; err != nil {
		return nil, err
	}
	return &domain.Setting{Key: m.Key, Value: m.Value, UpdatedBy: m.UpdatedBy, UpdatedAt: m.UpdatedAt}, nil
}

func (s *Store) Set(ctx context.Context, setting *domain.Setting) error {
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "setting_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_by", "updated_at"}),
	}).Create(&Model{Key: setting.Key, Value: setting.Value, UpdatedBy: setting.UpdatedBy, UpdatedAt: setting.UpdatedAt}).Error
}

func (s *Store) List(ctx context.Context) ([]domain.Setting, error) {
	var rows []Model
	if err := s.db.WithContext(ctx).Order("setting_key asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Setting, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.Setting{Key: r.Key, Value: r.Value, UpdatedBy: r.UpdatedBy, UpdatedAt: r.UpdatedAt})
	}
	return out, nil
}
