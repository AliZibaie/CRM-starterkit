package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/discount/domain"

	"gorm.io/gorm"
)

type Model struct {
	ID             string     `gorm:"primaryKey;size:36"`
	Code           string     `gorm:"uniqueIndex;size:50;not null"`
	Type           string     `gorm:"size:20;not null"`
	Value          int64      `gorm:"not null"`
	MinOrderAmount int64      `gorm:"not null;default:0"`
	MaxUses        int64      `gorm:"not null;default:0"`
	UsedCount      int64      `gorm:"not null;default:0"`
	StartsAt       *time.Time
	EndsAt         *time.Time
	IsActive       bool      `gorm:"not null;default:true"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

func (Model) TableName() string { return "discount_codes" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error { return s.db.AutoMigrate(&Model{}) }

func (s *Store) Create(ctx context.Context, d *domain.Discount) error {
	return s.db.WithContext(ctx).Create(toModel(d)).Error
}

func (s *Store) FindByCode(ctx context.Context, code string) (*domain.Discount, error) {
	var m Model
	if err := s.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (s *Store) FindByID(ctx context.Context, id string) (*domain.Discount, error) {
	var m Model
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (s *Store) List(ctx context.Context) ([]domain.Discount, error) {
	var rows []Model
	if err := s.db.WithContext(ctx).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Discount, 0, len(rows))
	for i := range rows {
		out = append(out, *toDomain(&rows[i]))
	}
	return out, nil
}

func (s *Store) IncrementUsage(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Model(&Model{}).Where("id = ?", id).Updates(map[string]any{
		"used_count": gorm.Expr("used_count + 1"),
		"updated_at": time.Now().UTC(),
	}).Error
}

func toModel(d *domain.Discount) *Model {
	return &Model{
		ID: d.ID, Code: d.Code, Type: d.Type, Value: d.Value, MinOrderAmount: d.MinOrderAmount,
		MaxUses: d.MaxUses, UsedCount: d.UsedCount, StartsAt: d.StartsAt, EndsAt: d.EndsAt,
		IsActive: d.IsActive, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func toDomain(m *Model) *domain.Discount {
	return &domain.Discount{
		ID: m.ID, Code: m.Code, Type: m.Type, Value: m.Value, MinOrderAmount: m.MinOrderAmount,
		MaxUses: m.MaxUses, UsedCount: m.UsedCount, StartsAt: m.StartsAt, EndsAt: m.EndsAt,
		IsActive: m.IsActive, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}
